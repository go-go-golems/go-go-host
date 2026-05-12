const express = require("express");
const ui = require("ui.dsl");
const db = require("database");
const guard = require("db.guard");

const app = express.app();

db.exec(`
  CREATE TABLE IF NOT EXISTS visits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL,
    host TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
  )
`);

app.get("/", (req, res) => {
  db.exec("INSERT INTO visits (path, host) VALUES (?, ?)", req.path, req.platform.host || "unknown");
  const rows = db.query("SELECT COUNT(*) AS count FROM visits");
  return ui.page(
    { title: "go-go-host beta demo" },
    ui.main(
      ui.h1("Hello from go-go-host beta"),
      ui.p("This is a live hosted Goja app served through wildcard DNS, wildcard TLS, Traefik, and go-go-host runtime routing."),
      ui.p("Host: " + (req.platform.host || "unknown")),
      ui.p("Site ID: " + (req.platform.siteId || "unknown")),
      ui.p("Deployment ID: " + (req.platform.deploymentId || "unknown")),
      ui.p("Visits recorded in the per-site SQLite DB: " + rows[0].count),
      ui.ul(
        ui.li(ui.a({ href: "/platform" }, "Platform JSON")),
        ui.li(ui.a({ href: "/db" }, "DB guard JSON")),
        ui.li(ui.a({ href: "/assets/style.css" }, "Static asset"))
      )
    )
  );
});

app.get("/platform", (req, res) => res.json(req.platform));
app.get("/db", (req, res) => res.json({ stats: guard.stats(), overLimit: guard.isOverLimit() }));
