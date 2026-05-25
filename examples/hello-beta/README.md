# hello-beta

Small go-go-host beta smoke app. It exercises the runtime router, UI DSL rendering, static assets, per-site SQLite state, and DB guard stats.

Build a bundle from the repository root:

```bash
tar -C examples/hello-beta -czf /tmp/hello-beta.tar.gz .
```

Expected routes after deployment:

- `/` renders an HTML page and increments a SQLite-backed visit counter.
- `/platform` returns the go-go-host platform context JSON.
- `/db` returns DB guard/quota stats.
- `/assets/style.css` serves a static CSS asset.
