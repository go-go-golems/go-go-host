---
Title: "Agent Guide: Build, Enroll, and Deploy from CI"
Slug: "agent-guide"
Short: "A complete guide for machine deploy agents: grants, keys, enrollment, signed deploys, activation, rotation, and troubleshooting."
Topics:
  - go-go-host-agent
  - agents
  - deployments
  - ci
  - security
Commands:
  - keygen
  - enroll
  - deploy
  - status
Flags:
  - config
  - token
  - bundle
  - site-id
  - channel
  - path
  - activate
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

A go-go-host agent is a machine identity that can deploy to a site without borrowing a human user's credentials. The human operator creates the agent, grants it access to one or more sites, and hands the machine a one-time enrollment token. The machine generates an Ed25519 key pair, enrolls the public key, and uses the private key to sign deploy-run requests.

This design gives the platform a useful security property: the agent can be powerful enough to deploy, but still narrower than a human account. Its grant names the site, channel, bundle path policy, and whether it may activate traffic. Its signing key can be rotated or revoked without deleting the whole agent.

## The three actors

There are three participants in a signed deployment.

| Actor | Responsibility |
| --- | --- |
| Human operator | Creates the org/site, creates the agent, and grants site access. |
| Agent machine | Stores the private key and runs `go-go-host-agent deploy`. |
| go-go-host daemon | Verifies signatures, grants, upload tokens, bundle policy, and activation rules. |

The agent never needs the human's bearer token or dev-user identity. It receives a one-time enrollment token, then uses signed requests from that point forward.

## What the human creates

The human side uses the `go-go-host` CLI, not the agent CLI. This command creates an agent, creates a site grant, and returns a one-time enrollment token.

```bash
go-go-host agents create \
  --api-url http://127.0.0.1:8080 \
  --dev-user alice \
  --org-id ORG_ID \
  --name ci-agent \
  --site-id SITE_ID \
  --channel default \
  --bundle-path '**' \
  --can-activate \
  --output json
```

The grant is the important part. It says what the agent may do.

| Grant field | Meaning |
| --- | --- |
| `siteId` | The only site this grant applies to. |
| `allowedChannels` | The deploy channels this agent may request, commonly `default`. |
| `allowedBundlePaths` | Logical bundle path policy. `**` means all bundle paths; this does not constrain files inside the tar/zip archive. |
| `canDeploy` | Allows deploy-run creation and bundle upload. |
| `canActivate` | Allows the agent to request `--activate` and promote traffic automatically. |
| `expiresAt` | Optional expiry for temporary automation. |

Use `--can-activate` carefully. It lets CI move traffic after upload validation succeeds. Without it, the agent can upload a validated deployment, but a human must activate it.

## Agent configuration file

The agent stores its local state in a JSON config file. Treat this file like a secret because it contains the private signing key.

```bash
export AGENT_CONFIG=./go-go-host-agent.json
```

A typical enrolled config contains:

```json
{
  "apiUrl": "http://127.0.0.1:8080",
  "agentId": "agt_...",
  "keyId": "ak_...",
  "privateKey": "...",
  "publicKey": "..."
}
```

Recommended filesystem handling:

```bash
umask 077
go-go-host-agent keygen --config "$AGENT_CONFIG" --api-url http://127.0.0.1:8080
```

Store it in your CI secret store or machine-local secret volume. Do not commit it.

## Enroll the machine

The machine first generates a key pair.

```bash
go-go-host-agent keygen \
  --config "$AGENT_CONFIG" \
  --api-url http://127.0.0.1:8080 \
  --output json
```

Then it exchanges the one-time token for key registration.

```bash
go-go-host-agent enroll \
  --config "$AGENT_CONFIG" \
  --api-url http://127.0.0.1:8080 \
  --token ENROLLMENT_TOKEN \
  --output json
```

The enrollment token is single-use. If enrollment succeeds and the file is lost, create a rotation enrollment token from the human dashboard or CLI and enroll a new key.

## Deploy from the agent

A signed deploy has two phases. First the agent signs a JSON deploy-run request. If the daemon accepts the signature and grant, it returns a short-lived upload token. Then the agent uploads the bundle using that upload token. The token is bound to the deploy run and cannot be reused.

```bash
go-go-host-agent deploy \
  --config "$AGENT_CONFIG" \
  --api-url http://127.0.0.1:8080 \
  --bundle ./site.tar.gz \
  --site-id SITE_ID \
  --channel default \
  --bundle-path bundles/site.tar.gz \
  --output json
```

If the grant includes `canActivate`, the agent can ask for scoped auto-activation:

```bash
go-go-host-agent deploy \
  --config "$AGENT_CONFIG" \
  --api-url http://127.0.0.1:8080 \
  --bundle ./site.tar.gz \
  --site-id SITE_ID \
  --channel default \
  --bundle-path bundles/site.tar.gz \
  --activate \
  --output json
```

Auto-activation still requires the uploaded bundle to validate. The server does not trust an upload-time activation flag; it checks the persisted deploy-run allowed actions.

## A CI job sketch

A real CI job usually has three inputs: the agent config secret, the target site ID, and the bundle produced by the build step.

```bash
set -euo pipefail

: "${GO_GO_HOST_AGENT_CONFIG:?missing config path}"
: "${GO_GO_HOST_SITE_ID:?missing site id}"
: "${GO_GO_HOST_API_URL:=https://host.example.com}"

npm ci
npm run build

tar -C dist-go-go-host -czf site.tar.gz .

go-go-host-agent status \
  --config "$GO_GO_HOST_AGENT_CONFIG" \
  --api-url "$GO_GO_HOST_API_URL" \
  --output json

go-go-host-agent deploy \
  --config "$GO_GO_HOST_AGENT_CONFIG" \
  --api-url "$GO_GO_HOST_API_URL" \
  --bundle site.tar.gz \
  --site-id "$GO_GO_HOST_SITE_ID" \
  --channel default \
  --bundle-path "bundles/${GITHUB_SHA:-manual}.tar.gz" \
  --activate \
  --output json
```

The path should be deterministic enough for audit and broad enough to match the grant. A common pattern is `bundles/<git-sha>.tar.gz` with a grant path like `bundles/**`.

## What is signed

The agent signs deploy-run creation, not the multipart upload body. The signed request includes the intent: site, channel, path, and requested actions. If the daemon accepts that intent, it returns an upload token for exactly that deploy run.

Signed requests include headers like:

```text
X-Go-Go-Agent-ID
X-Go-Go-Agent-Key-ID
X-Go-Go-Agent-Timestamp
X-Go-Go-Agent-Nonce
X-Go-Go-Agent-Signature
```

The daemon rejects:

- missing signature headers,
- timestamps outside the allowed skew window,
- replayed nonces,
- invalid signatures,
- revoked agents,
- revoked keys,
- grants for the wrong site,
- channels or paths outside the grant,
- reused upload tokens.

## Key rotation

Key rotation replaces the signing key without replacing the whole agent identity. The human operator creates a replacement enrollment token for the existing agent. The machine generates a new config or key, enrolls the replacement key, verifies deployment success, and then the old key can be revoked.

High-level sequence:

```text
human: create replacement token for agent
agent: keygen new config
agent: enroll with replacement token
agent: status/deploy smoke with new key
human: revoke old key
```

Use the dashboard Agents page for key inventory and revoke actions, or use the relevant human CLI/API endpoints. Keep at least one known-good active key until the replacement key has completed a deploy-run successfully.

## Troubleshooting signed deploys

| Symptom | Likely cause | Fix |
| --- | --- | --- |
| `agent_signature_missing` | Required signature headers were not sent. | Use `go-go-host-agent deploy`; do not hand-roll requests unless reproducing the canonical signature exactly. |
| `agent_timestamp_skew` | Machine clock differs from daemon clock. | Sync NTP on the CI runner. |
| `agent_nonce_replay` | A signed request was retried with the same nonce. | Run the CLI again; do not replay captured signed requests. |
| `agent_key_revoked` | The key was revoked or inactive. | Enroll a replacement key or ask an owner to reactivate policy through a new token. |
| `agent_grant_denied` | Site, channel, bundle path, or activation request does not match the grant. | Compare `--site-id`, `--channel`, `--bundle-path`, and `--activate` with the grant. |
| `upload_token_invalid` | Upload token expired, was reused, or does not belong to the run. | Start a fresh deploy command. |
| Bundle validates but does not serve traffic | The grant lacks `canActivate` or `--activate` was not used. | Ask a human to activate the deployment or grant scoped activation. |

Use JSON output while debugging:

```bash
go-go-host-agent status --config "$AGENT_CONFIG" --output json
go-go-host-agent deploy --config "$AGENT_CONFIG" --bundle ./site.tar.gz --site-id SITE_ID --channel default --bundle-path bundles/debug.tar.gz --output json
```

## What to tell a deployment agent

A deployment agent needs five facts and one file:

```text
API URL:        https://host.example.com
site ID:        site_...
channel:        default
bundle path:    bundles/<sha>.tar.gz, matching the human grant
activation:     whether --activate is allowed
config file:    enrolled go-go-host-agent JSON config
```

With those in place, the command is simple:

```bash
go-go-host-agent deploy \
  --config ./agent.json \
  --api-url https://host.example.com \
  --bundle ./site.tar.gz \
  --site-id site_... \
  --channel default \
  --bundle-path bundles/$GIT_SHA.tar.gz \
  --activate
```

## See Also

- `go-go-host-agent help agent-keygen-enroll-deploy`
- `go-go-host-agent help agent-signature-troubleshooting`
- `go-go-host help developer-guide`
- `go-go-host help js-api-reference`
