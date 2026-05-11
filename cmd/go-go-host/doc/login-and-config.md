---
Title: "Login and CLI Configuration"
Slug: "login-and-config"
Short: "Store API URL and auth defaults for go-go-host CLI commands."
Topics:
  - go-go-host
  - login
  - configuration
Commands:
  - login
  - me
Flags:
  - api-url
  - dev-user
  - bearer-token
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

This tutorial covers CLI configuration for local development and bearer-token smoke testing. The current v1 CLI stores defaults in a local config file so subsequent commands do not need repeated API/auth flags.

## Configure dev auth

In dev mode, the daemon accepts `X-Go-Go-Host-User`. Store that value with `login`:

```bash
go-go-host login --api-url http://127.0.0.1:8080 --dev-user alice
```

Then verify the session context:

```bash
go-go-host me --output table
```

## Configure bearer auth

For non-dev auth smoke tests, store a bearer token instead:

```bash
go-go-host login --api-url https://host.example.com --bearer-token "$TOKEN"
```

Commands also accept `--bearer-token` directly, which overrides the stored value for that invocation.

## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| Commands still use the wrong API URL | A local config file is taking precedence | Pass `--api-url` explicitly or set `GO_GO_HOST_CLI_CONFIG` to a test config path |
| `me` returns unauthorized | Missing dev user or invalid bearer token | Re-run `login` with the correct auth option |
| Tests affect your normal config | The default config lives in the OS user config directory | Set `GO_GO_HOST_CLI_CONFIG=$(mktemp)` for smoke tests |

## See Also

- `getting-started`
- `deploy-workflow`
