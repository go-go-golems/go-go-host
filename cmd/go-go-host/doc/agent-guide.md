---
Title: "Agent Guide for Operators"
Slug: "agent-guide"
Short: "Create deployment agents, grant site access, hand off enrollment tokens, and operate key rotation/revoke workflows."
Topics:
  - go-go-host
  - agents
  - deployments
  - security
Commands:
  - agents
  - audit
Flags:
  - org-id
  - site-id
  - channel
  - path
  - can-activate
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

An agent is a machine identity for deployment automation. It is not a human user with a password or bearer token. It is an object in an organization, plus one or more signing keys, plus grants that say which sites it may deploy to. The operator controls the grant; the machine controls the private key.

This guide is written for the human operator who creates the agent and decides how much authority it should have. The companion machine-side guide is bundled in the agent binary as `go-go-host-agent help agent-guide`.

## Create a scoped agent

Create an agent with a grant for one site. The returned enrollment token is shown once, so copy it into your CI secret handoff process immediately.

```bash
go-go-host agents create \
  --api-url http://127.0.0.1:8080 \
  --dev-user alice \
  --org-id ORG_ID \
  --name ci-agent \
  --site-id SITE_ID \
  --channel default \
  --path 'bundles/**' \
  --output json
```

Add `--can-activate` only when the CI system should be allowed to promote traffic after a bundle validates.

```bash
go-go-host agents create \
  --org-id ORG_ID \
  --name production-ci \
  --site-id SITE_ID \
  --channel default \
  --path 'bundles/**' \
  --can-activate \
  --output json
```

The difference is significant. Without `canActivate`, an agent can upload a validated deployment but a human must activate it. With `canActivate`, the signed deploy run can include `activate`, and the server will promote the deployment automatically after validation succeeds.

## Hand off enrollment

Give the machine operator or CI secret store these values:

```text
apiUrl: http://127.0.0.1:8080
siteId: SITE_ID
channel: default
path pattern: bundles/<git-sha>.tar.gz
enrollmentToken: enroll_...
```

The machine then runs:

```bash
go-go-host-agent keygen --config ./agent.json --api-url http://127.0.0.1:8080
go-go-host-agent enroll --config ./agent.json --api-url http://127.0.0.1:8080 --token enroll_...
go-go-host-agent deploy --config ./agent.json --bundle ./site.tar.gz --site-id SITE_ID --channel default --path bundles/site.tar.gz
```

## Inspect and revoke keys

The dashboard Agents page shows signing keys, fingerprints, creation time, last-used time, and revoked state. Use key-level revoke when a single runner is compromised. Use agent revoke when the whole automation identity should stop.

Key-level revoke is safer during rotation because the agent object and its grants remain intact while you replace one key.

## Rotate a key

A replacement-key flow should prove the new key works before deleting the old key.

1. Create a replacement enrollment token for the existing agent.
2. Generate and enroll a new key on the machine.
3. Run `go-go-host-agent status` and a deploy smoke with the new config.
4. Revoke the old key.
5. Inspect audit events for enrollment, deploy, and revoke actions.

## Audit agent activity

Use audit filters to inspect what an agent did.

```bash
go-go-host audit list \
  --org-id ORG_ID \
  --actor-type agent \
  --limit 100 \
  --output json
```

Interesting actions include deployment upload, deployment activation, signature failures, grant denials, and key revocation. Failed signed requests are intentionally audited with stable reason codes so incident response can distinguish clock skew from replay attempts or revoked keys.

## See Also

- `go-go-host-agent help agent-guide`
- `go-go-host-agent help agent-signature-troubleshooting`
- `go-go-host help developer-guide`
- `go-go-host help js-api-reference`
