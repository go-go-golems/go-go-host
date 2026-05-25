---
Title: "JavaScript API Reference for Hosted Sites"
Slug: "js-api-reference"
Short: "Detailed reference for go-go-host JavaScript modules, route handlers, UI DSL, database access, platform context, and capability policy."
Topics:
  - go-go-host
  - javascript
  - api-reference
  - goja
  - capabilities
Commands:
  - deploy
  - deployments
Flags:
  - site-id
  - path
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

The hosted JavaScript API is deliberately small. It gives an app enough power to handle HTTP requests, render HTML, store local data, inspect runtime quota state, and serve static assets. It does not give the app raw filesystem access, subprocess execution, process environment access, or arbitrary host control. That boundary is the point of the platform: hosted apps should feel productive without becoming trusted host programs.

This reference describes the API as it exists today. Treat it as both documentation and a contract for bundle authors. If a function is not listed here, do not assume it is supported in hosted v1.

## Module overview

| Module | Import | Purpose |
| --- | --- | --- |
| Express router | `require("express")` | Register HTTP routes for the site. |
| UI DSL | `require("ui.dsl")` or `require("ui")` | Build HTML documents and fragments as JavaScript values. |
| Database | `require("database")` or `require("db")` | Query the site's preconfigured SQLite database. |
| DB guard | `require("db.guard")` | Inspect database quota state and register limit callbacks. |
| Utilities | `require("path")`, `require("time")`, `require("timer")` | Safe runtime utility modules from go-go-goja middleware. |

Unavailable by design:

```js
require("fs");   // blocked in hosted v1
require("exec"); // blocked in hosted v1
```

## Bundle manifest

The manifest file is named `go-go-host.json` and must be at the archive root. It is read during upload validation before scripts run.

```json
{
  "name": "docs-site",
  "scriptsDir": "scripts",
  "assetsDir": "assets",
  "entrypoint": "app.js",
  "smokePath": "/",
  "capabilities": ["express", "ui.dsl", "database", "assets"],
  "allowedPaths": ["**"],
  "channel": "default"
}
```

| Field | Required | Meaning |
| --- | --- | --- |
| `name` | No | Human-readable bundle name. |
| `scriptsDir` | Yes | Directory containing JavaScript files. Every `.js` file under it is loaded in sorted order. |
| `assetsDir` | No | Directory served under `/assets` when assets are enabled. |
| `entrypoint` | No | Present for convention and future tooling. Current loading walks all `.js` files. |
| `smokePath` | No | Path requested during dry-run validation. Defaults to `/` if empty. |
| `capabilities` | No | Capabilities requested by the bundle. Site policy may deny them. |
| `allowedPaths` | No | Additional archive-entry allowlist declared by the bundle manifest. Agent `allowedBundlePaths` are also enforced against uploaded archive entries, so restricted agents need patterns for required files such as `go-go-host.json`, `scripts/**`, and `assets/**`. |
| `channel` | No | Deployment channel label, commonly `default`. |

The validator rejects paths that are absolute, contain `..`, or escape the bundle. This is true for manifest paths and archive entries.

## The Express module

Import and create an app:

```js
const express = require("express");
const app = express.app();
```

Register routes:

```js
app.get("/", handler);
app.post("/items", handler);
app.put("/items/:id", handler);
app.patch("/items/:id", handler);
app.delete("/items/:id", handler);
app.all("/health", handler);
```

A handler receives `(req, res)`. If the handler returns a value and has not already sent a response, the host renders or serializes that value.

```js
app.get("/hello", (req, res) => {
  return "hello";
});

app.get("/json", (req, res) => {
  return res.json({ ok: true });
});
```

### Static assets

Assets are normally mounted by the host at `/assets` from the manifest `assetsDir`. The app object also exposes:

```js
app.static(prefix, dir);
```

Use app-level static mounts sparingly. The preferred v1 pattern is to declare `assetsDir` in the manifest and request the `assets` capability.

## Request object

The request object is a plain JavaScript object derived from the incoming HTTP request.

```js
app.get("/inspect/:name", (req, res) => {
  return res.json({
    method: req.method,
    url: req.url,
    path: req.path,
    query: req.query,
    params: req.params,
    headers: req.headers,
    cookies: req.cookies,
    session: req.session,
    ip: req.ip,
    body: req.body,
    rawBody: req.rawBody,
    platform: req.platform
  });
});
```

| Property | Type | Notes |
| --- | --- | --- |
| `method` | string | HTTP method. |
| `url` | string | Path plus query string. |
| `path` | string | URL path. |
| `query` | object | Query parameters. Single values are strings; repeated values are arrays. |
| `params` | object | Route parameters when the route pattern captures them. |
| `headers` | object | Request headers as joined strings. |
| `cookies` | object | Cookie name/value map. |
| `session` | object | Session DTO. Present even when unused. |
| `ip` | string | Remote address best effort. |
| `body` | any | Parsed body when content type is supported. |
| `rawBody` | string | Raw request body. |
| `platform` | object | go-go-host request metadata. |

### Platform context

`req.platform` is how hosted code learns where it is running without reading environment variables.

```js
app.get("/platform", (req, res) => res.json(req.platform));
```

Typical shape:

```json
{
  "requestId": "20260512T120000.000000000Z",
  "orgId": "org_...",
  "siteId": "site_...",
  "deploymentId": "dep_...",
  "host": "hello.localhost"
}
```

Use this for diagnostics and tenant-aware display. Do not treat it as an authentication claim from the user; it is platform metadata about the serving runtime.

## Response object

The response object controls status, headers, and output.

```js
res.status(201);
res.set("X-App", "docs-site");
res.type("application/json");
res.send("plain text");
res.json({ ok: true });
res.html(ui.page({ title: "Hi" }, ui.main(ui.h1("Hi"))));
res.redirect("/new-location");
res.redirect(301, "/permanent-location");
res.end();
```

Important behavior:

- `res.send(string)` writes text or HTML depending on the string shape and headers.
- `res.send(object)` behaves like JSON serialization.
- `res.html(value)` uses the configured UI renderer.
- After a response has been sent, later send calls are ignored.
- A route may either call `res.*` or return a value; do not rely on both.

## UI DSL

The UI DSL represents HTML as data. This makes simple server-rendered pages pleasant to write without introducing a template language.

```js
const ui = require("ui.dsl");

app.get("/", (req, res) => {
  return ui.page(
    { title: "Catalog" },
    ui.main(
      ui.h1("Catalog"),
      ui.p("Choose an item."),
      ui.ul(
        ui.li(ui.a({ href: "/items/1" }, "Item 1")),
        ui.li(ui.a({ href: "/items/2" }, "Item 2"))
      )
    )
  );
});
```

### Element construction

Every element function accepts optional attributes as the first argument, followed by children.

```js
ui.a({ href: "/docs", class: "button" }, "Read docs")
ui.input({ type: "email", name: "email", required: true })
ui.div({ id: "app", dataRole: "panel" }, ui.p("Hello"))
```

Boolean attributes are rendered when truthy. Null and undefined children are ignored. Arrays and fragments are flattened.

### Document helpers

```js
ui.page({ title: "Title" }, ...children)
ui.fragment(...children)
ui.text(value)
ui.raw("<span>trusted html</span>")
ui.render(node)
```

Use `ui.raw` only for trusted content. It bypasses escaping by design.

### Supported tags

```text
html head body title meta link script style
main div span section article header footer nav
h1 h2 h3 h4 p a form input button select option textarea label
ul ol li table thead tbody tr th td
strong em small pre code
img br hr time
svg path rect line polyline circle
```

## Database API

The database module exposes the site's preconfigured SQLite database. It is guarded by quota settings and does not allow hosted code to choose an arbitrary database file.

```js
const db = require("database");
// or: const db = require("db");
```

Common operations:

```js
db.exec("CREATE TABLE IF NOT EXISTS notes (id INTEGER PRIMARY KEY, body TEXT NOT NULL)");
db.exec("INSERT INTO notes (body) VALUES (?)", "hello");
const rows = db.query("SELECT id, body FROM notes ORDER BY id DESC LIMIT ?", 10);
```

The database module comes from go-go-goja's database module, configured by the host. The practical contract for hosted apps is:

- Use `exec(sql, ...args)` for statements that modify schema or data.
- Use `query(sql, ...args)` for statements that return rows.
- Expect rows as JavaScript objects keyed by column name.
- Do not call `configure`; it is disabled in hosted v1.

This should fail:

```js
try {
  db.configure("sqlite3", ":memory:");
} catch (err) {
  // expected in hosted v1
}
```

## DB guard API

The DB guard observes database size and enforces soft/hard quota policy configured by the platform.

```js
const guard = require("db.guard");

app.get("/db-guard", (req, res) => {
  return res.json({
    stats: guard.stats(),
    overLimit: guard.isOverLimit(),
    last: guard.lastResult()
  });
});
```

Available functions:

| Function | Purpose |
| --- | --- |
| `guard.stats()` | Return current DB file size and quota measurements. |
| `guard.checkNow(reason)` | Force a quota check and return the result. |
| `guard.isOverLimit()` | Return whether the last known state is over limit. |
| `guard.lastResult()` | Return the last check result. |
| `guard.onLimitExceeded(fn)` | Register a callback for limit events. |
| `guard.configure(options)` | Low-level hook exists, but site quota remains the platform authority. |

Although `configure` exists in the module, app authors should treat the server-side site quota as authoritative. A future release may further restrict app-level guard reconfiguration.

## Utility modules

The runtime enables safe utility modules from go-go-goja middleware:

```js
const path = require("path");
const time = require("time");
const timer = require("timer");
```

Use them for local computation, formatting, and timers. Do not use timers to create unbounded background work; request timeout enforcement exists at the HTTP layer, and deeper interruption semantics may evolve.

## Error handling pattern

A good hosted route catches expected application errors and lets unexpected errors fail visibly during validation or runtime diagnostics.

```js
app.post("/notes", (req, res) => {
  if (!req.body || !req.body.body) {
    return res.status(400).json({ error: "body is required" });
  }
  db.exec("INSERT INTO notes (body) VALUES (?)", req.body.body);
  return res.status(201).json({ ok: true });
});
```

If a script throws during startup, deployment validation fails. If a handler throws during a request, the host returns an error and records runtime error counters.

## Capability reference

Capabilities are requested by the manifest and checked against site policy during validation.

| Capability | Enables | Denial symptom |
| --- | --- | --- |
| `express` | Route registration. | Bundle requesting it fails validation if site policy denies it. |
| `ui.dsl` | HTML DSL module. | `require("ui.dsl")` is not part of permitted policy. |
| `database` | `require("database")`. | Validation denies requested capability. |
| `db` | `require("db")` alias. | Validation denies requested capability. |
| `assets` | Static assets from `assetsDir`. | Assets are not mounted and validation may deny requested capability. |
| `time` / `timer` | Runtime time utilities. | Validation denies requested capability. |
| `sqlite` | Compatibility label. | Used for policy/reporting compatibility. |

A bundle should request only the capabilities it uses. That makes validation reports meaningful and keeps the review surface small.

## Full smoke app

Use this app to test the runtime surface end to end.

```js
const express = require("express");
const ui = require("ui.dsl");
const db = require("database");
const guard = require("db.guard");
const app = express.app();

db.exec("CREATE TABLE IF NOT EXISTS visits (id INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT NOT NULL)");

app.get("/", (req, res) => {
  db.exec("INSERT INTO visits (path) VALUES (?)", req.path);
  const rows = db.query("SELECT COUNT(*) AS count FROM visits");
  return ui.page({ title: "Smoke" }, ui.main(ui.h1("Smoke"), ui.p("Visits: " + rows[0].count)));
});

app.get("/json", (req, res) => res.json({ ok: true, query: req.query }));
app.get("/platform", (req, res) => res.json(req.platform));
app.get("/db-guard", (req, res) => res.json({ stats: guard.stats(), last: guard.lastResult() }));
app.get("/forbidden", (req, res) => {
  const out = {};
  for (const name of ["fs", "exec"]) {
    try { require(name); out[name] = "unexpectedly available"; }
    catch (_) { out[name] = "blocked"; }
  }
  return res.json(out);
});
```

After activation:

```bash
curl -fsS -H "Host: $SITE_HOST" "$API/"
curl -fsS -H "Host: $SITE_HOST" "$API/json?x=1" | jq .
curl -fsS -H "Host: $SITE_HOST" "$API/platform" | jq .
curl -fsS -H "Host: $SITE_HOST" "$API/db-guard" | jq .
curl -fsS -H "Host: $SITE_HOST" "$API/forbidden" | jq .
```

## Troubleshooting

| Problem | Likely cause | Fix |
| --- | --- | --- |
| `missing go-go-host.json manifest` | The archive was created from the wrong directory. | Run `tar -C "$APP_DIR" -czf bundle.tar.gz .` so the manifest is at archive root. |
| `capability "x" is not permitted` | The site policy denies a requested capability. | Enable the capability on the site Settings page or remove it from the manifest. |
| Dry-run smoke fails | The app did not register the smoke route or threw while handling it. | Set `smokePath` to a route that returns 2xx/3xx and inspect script errors. |
| `db.configure` fails | Hosted DBs are preconfigured and locked. | Use `db.query` and `db.exec`; do not open your own DB. |
| Static assets 404 | `assetsDir` is wrong or `assets` capability is disabled. | Fix manifest path and enable/request `assets`. |
| `fs` or `exec` cannot be required | This is expected in hosted v1. | Use supported APIs; request a platform feature if you need a new safe capability. |

## See Also

- `go-go-host help developer-guide`
- `go-go-host help deploy-workflow`
- `go-go-host help rollback-workflow`
- `go-go-host help agent-setup`
- `go-go-host-agent help agent-guide`
