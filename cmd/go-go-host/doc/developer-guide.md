---
Title: "Developer Guide: Build and Deploy go-go-host Apps"
Slug: "developer-guide"
Short: "Learn the complete app lifecycle: bundle layout, JavaScript runtime APIs, deployment, activation, settings, and operations."
Topics:
  - go-go-host
  - developer-guide
  - deployments
  - javascript
  - sites
Commands:
  - org
  - site
  - deploy
  - deployments
  - maintenance
Flags:
  - api-url
  - dev-user
  - site-id
  - path
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

A go-go-host app is a small JavaScript program packaged as a bundle and run inside a Goja runtime. The bundle contains a manifest, one or more scripts, and optional static assets. The platform supplies the HTTP router, the per-site SQLite database, a small HTML DSL, request metadata, deployment history, audit events, and operator settings. Your job as an app developer is to write the app code, request the capabilities it needs, package it safely, and deploy it to a site.

The most useful mental model is this: a site is stable, deployments are immutable, and activation changes traffic. You can upload many bundles to a site, inspect validation reports, and choose which deployment should serve traffic. That separation is what lets you test, roll back, export, and audit the app without losing track of what code ran when.

## The shape of an app

Every app bundle has the same outer shape. The manifest tells go-go-host where scripts and assets live. The scripts register HTTP routes. The assets directory is served under `/assets` when the site policy allows the `assets` capability.

```text
my-site/
├── go-go-host.json
├── scripts/
│   └── app.js
└── assets/
    └── style.css
```

The smallest useful manifest looks like this:

```json
{
  "name": "hello-card",
  "scriptsDir": "scripts",
  "assetsDir": "assets",
  "entrypoint": "app.js",
  "smokePath": "/",
  "capabilities": ["express", "ui.dsl", "database", "assets"],
  "channel": "default"
}
```

The `smokePath` matters because uploads perform a dry-run runtime load and request that path before the deployment is accepted as validated. A bundle that cannot boot or cannot answer its smoke route is rejected before it reaches production traffic.

## A complete first app

This app uses the Express-like router, the UI DSL, and the per-site database. It records page visits and renders an HTML page. It also exposes JSON endpoints that are useful while learning the platform.

```js
const express = require("express");
const ui = require("ui.dsl");
const db = require("database");
const guard = require("db.guard");

const app = express.app();

db.exec(`
  CREATE TABLE IF NOT EXISTS visits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
  )
`);

app.get("/", (req, res) => {
  db.exec("INSERT INTO visits (path) VALUES (?)", req.path);
  const rows = db.query("SELECT COUNT(*) AS count FROM visits");

  return ui.page(
    { title: "Hello Card" },
    ui.main(
      ui.h1("Hello from go-go-host"),
      ui.p("This page was rendered inside a hosted Goja runtime."),
      ui.p("Visits: " + rows[0].count),
      ui.ul(
        ui.li(ui.a({ href: "/platform" }, "Platform context")),
        ui.li(ui.a({ href: "/db" }, "Database guard stats")),
        ui.li(ui.a({ href: "/assets/style.css" }, "Static asset"))
      )
    )
  );
});

app.get("/platform", (req, res) => {
  return res.json(req.platform);
});

app.get("/db", (req, res) => {
  return res.json({ stats: guard.stats(), overLimit: guard.isOverLimit() });
});
```

Several details are worth noticing. The code does not create an HTTP server; the host owns the server and calls your registered handlers. The code does not open a database file; the host gives it a preconfigured, quota-guarded database. The code returns a UI node from the route handler, and the host renders that node as HTML because the response has not already been sent.

## Build the bundle

Create the files, then archive from inside the bundle directory so the manifest sits at the archive root.

```bash
APP_DIR=$(mktemp -d)
mkdir -p "$APP_DIR/scripts" "$APP_DIR/assets"

cat > "$APP_DIR/go-go-host.json" <<'EOF'
{
  "name": "hello-card",
  "scriptsDir": "scripts",
  "assetsDir": "assets",
  "entrypoint": "app.js",
  "smokePath": "/",
  "capabilities": ["express", "ui.dsl", "database", "assets"],
  "channel": "default"
}
EOF

cat > "$APP_DIR/scripts/app.js" <<'EOF'
const express = require("express");
const ui = require("ui.dsl");
const db = require("database");
const guard = require("db.guard");
const app = express.app();
db.exec("CREATE TABLE IF NOT EXISTS visits (id INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT NOT NULL)");
app.get("/", (req, res) => {
  db.exec("INSERT INTO visits (path) VALUES (?)", req.path);
  const rows = db.query("SELECT COUNT(*) AS count FROM visits");
  return ui.page({ title: "Hello Card" }, ui.main(ui.h1("Hello Card"), ui.p("Visits: " + rows[0].count)));
});
app.get("/platform", (req, res) => res.json(req.platform));
app.get("/db", (req, res) => res.json({ stats: guard.stats(), overLimit: guard.isOverLimit() }));
EOF

echo 'body { font-family: system-ui; }' > "$APP_DIR/assets/style.css"
tar -C "$APP_DIR" -czf /tmp/hello-card.tar.gz .
```

## Create the site

In local development, dev auth identifies a user by the `X-Go-Go-Host-User` header. The CLI exposes that as `--dev-user`.

```bash
export API=http://127.0.0.1:8080
export USER=alice

ORG_JSON=$(go-go-host org create \
  --api-url "$API" \
  --dev-user "$USER" \
  --slug alice-labs \
  --name "Alice Labs" \
  --output json)

ORG_ID=$(echo "$ORG_JSON" | jq -r '.[0].id')

SITE_JSON=$(go-go-host site create \
  --api-url "$API" \
  --dev-user "$USER" \
  --org-id "$ORG_ID" \
  --slug hello-card \
  --name "Hello Card" \
  --output json)

SITE_ID=$(echo "$SITE_JSON" | jq -r '.[0].id')
SITE_HOST=$(echo "$SITE_JSON" | jq -r '.[0].primary_host')
```

The primary host is generated from the site slug and the configured base domain. In local development it is usually something like `hello-card.localhost`.

## Deploy and activate

Upload validates the bundle and stores it as an immutable deployment. Activation tells the runtime supervisor to load that deployment and route traffic to it.

```bash
DEPLOY_JSON=$(go-go-host deploy \
  --api-url "$API" \
  --dev-user "$USER" \
  --site-id "$SITE_ID" \
  --path /tmp/hello-card.tar.gz \
  --message "first deploy" \
  --channel default \
  --output json)

DEPLOYMENT_ID=$(echo "$DEPLOY_JSON" | jq -r '.[0].id')

go-go-host deployments activate \
  --api-url "$API" \
  --dev-user "$USER" \
  --deployment-id "$DEPLOYMENT_ID" \
  --output json
```

Now ask the daemon for the public host. The Host header is the important piece; it is how go-go-host decides which site runtime should handle the request.

```bash
curl -fsS -H "Host: $SITE_HOST" "$API/"
curl -fsS -H "Host: $SITE_HOST" "$API/platform" | jq .
```

## Configuration belongs outside the bundle

Do not bake every operational setting into JavaScript. Phase 11 added site configuration as a separate API surface. That means a deployment bundle can remain immutable while operators change non-secret settings around it.

Use the dashboard at:

```text
/app/orgs/{orgId}/sites/{siteId}/settings
```

or the API:

```bash
curl -X PUT "$API/api/v1/sites/$SITE_ID/config" \
  -H "X-Go-Go-Host-User: $USER" \
  -H 'Content-Type: application/json' \
  -d '{"key":"theme.title","value":{"text":"Hello Card"}}'
```

Secrets are intentionally not part of the v1 JavaScript API. There is no process environment passthrough, no plaintext secret read API, and no unrestricted filesystem module. This is a design boundary, not a missing convenience.

## Capabilities are the contract

The manifest may request capabilities, but the site policy decides whether those capabilities are allowed. The deployment report includes both requested and effective capability information. This matters because capability policy is how operators keep hosted code inside the intended sandbox.

| Capability | Purpose | Notes |
| --- | --- | --- |
| `express` | Route registration and HTTP request handling. | Needed for web apps. |
| `ui.dsl` | HTML node DSL and renderer. | Also available as `ui`. |
| `database` / `db` | Preconfigured SQLite access. | `configure()` is disabled. |
| `db.guard` | Database quota inspection. | Available for observing DB limits. |
| `assets` | Static files under `/assets`. | Depends on `assetsDir`. |
| `time`, `timer`, `path` | Utility modules. | Safe runtime middleware. |

The following are intentionally unavailable in hosted v1:

- `fs`, because unrestricted host file reads break tenant isolation.
- `exec`, because subprocess execution is outside the v1 trust model.
- process environment passthrough, because secrets need a separate encrypted design.

## Operations during development

A normal development loop uses these commands:

```bash
go-go-host deployments list --site-id "$SITE_ID" --dev-user "$USER" --output table
go-go-host rollback --site-id "$SITE_ID" --dev-user "$USER" --output json
go-go-host audit list --org-id "$ORG_ID" --dev-user "$USER" --limit 100 --output json
go-go-host maintenance export metadata --site-id "$SITE_ID" --dev-user "$USER" -o site.json
```

Use the dashboard for the same ideas visually:

- `/app` for org/site/deployment/runtime/settings/agent workflows.
- `/admin` for platform-wide inventory and policy views.
- Storybook on port `6007` for UI component review during development.

## A useful debugging sequence

When a deployment fails, debug in this order:

1. Inspect the validation report. If the bundle is rejected, activation is not the problem.
2. Check `go-go-host.json` paths. The manifest must be at the archive root and paths must not escape the bundle.
3. Check capabilities. A requested capability that site policy denies causes validation failure.
4. Check the smoke route. The runtime must answer the manifest `smokePath` with a 2xx or 3xx response.
5. Check runtime status and audit events after activation.

```bash
go-go-host deployments show --deployment-id "$DEPLOYMENT_ID" --dev-user "$USER" --output json | jq .
curl -fsS "$API/api/v1/sites/$SITE_ID/runtime" -H "X-Go-Go-Host-User: $USER" | jq .
go-go-host audit list --org-id "$ORG_ID" --dev-user "$USER" --limit 50 --output json | jq .
```

## Key points

- A site is the stable identity; deployments are immutable snapshots; activation chooses which deployment serves traffic.
- The JavaScript app registers routes but does not own the HTTP server.
- The database is preconfigured and quota-guarded; hosted code must not open arbitrary database files.
- Capabilities are both documentation and policy. Request only what your app needs.
- Operational settings belong in site config and domains, not in ad-hoc mutable bundle edits.

## See Also

- `go-go-host help js-api-reference`
- `go-go-host help deploy-workflow`
- `go-go-host help create-site-workflow`
- `go-go-host help rollback-workflow`
- `go-go-host-agent help agent-guide`
