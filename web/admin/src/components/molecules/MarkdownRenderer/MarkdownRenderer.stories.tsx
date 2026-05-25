import type { Meta, StoryObj } from '@storybook/react';
import { MarkdownRenderer } from './MarkdownRenderer';

const shortMd = `# Getting started

A **go-go-host** app is a small JavaScript program packaged as a bundle and run inside a Goja runtime. The bundle contains a manifest, one or more scripts, and optional static assets.

This guide walks you through creating your first site and deploying a hello-world bundle.
`;

const codeHeavyMd = `## The Express module

Register HTTP routes for the site:

\`\`\`js
const express = require("express");
const app = express.app();

app.get("/", (req, res) => {
  return res.json({ hello: "world" });
});

app.get("/greet/:name", (req, res) => {
  return res.json({ greeting: "Hello, " + req.params.name + "!" });
});
\`\`\`

You can also return a UI DSL node instead of calling \`res.json()\`:

\`\`\`js
const ui = require("ui.dsl");

app.get("/", (req, res) => {
  return ui.page(
    { title: "Hello" },
    ui.main(ui.h1("Hello from go-go-host"))
  );
});
\`\`\`

Unavailable by design:

\`\`\`js
require("fs");   // blocked in hosted v1
require("exec"); // blocked in hosted v1
\`\`\`
`;

const tableMd = `## Module overview

| Module | Import | Purpose |
| --- | --- | --- |
| Express router | \`require("express")\` | Register HTTP routes for the site. |
| UI DSL | \`require("ui.dsl")\` or \`require("ui")\` | Build HTML documents and fragments. |
| Database | \`require("database")\` or \`require("db")\` | Query the site's preconfigured SQLite database. |
| DB guard | \`require("db.guard")\` | Inspect database quota state. |
| Assets | declared in manifest | Static files under \`/assets\`. |
| Config | \`req.platform.config\` | Non-secret operator settings. |
`;

const headingsMd = `# JavaScript API Reference

## Module overview

### Express router

The Express module provides an HTTP router familiar from Node.js Express.

#### Route registration

Routes are registered before the runtime starts serving requests.

##### Method signatures

\`\`\`js
app.get(path, handler)
app.post(path, handler)
\`\`\`

### UI DSL

#### Element construction

##### Attributes

##### Children

## Database API

### Query methods

#### db.query

#### db.exec

### Error handling

## Capability reference

| Capability | Effect | Missing behavior |
| --- | --- | --- |
| \`express\` | Route registration | Bundle rejected |
| \`ui.dsl\` | HTML DSL module | Module unavailable |
| \`database\` | SQLite access | Module unavailable |
| \`assets\` | Static files | Assets not mounted |
`;

const blockquoteMd = `## Important notes

> The hosted JavaScript API is deliberately small. It gives an app enough power to handle HTTP requests, render HTML, store local data, and serve static assets. It does not give the app raw filesystem access or subprocess execution.

> **Warning:** Do not treat \`req.platform\` as an authentication claim from the user. It is platform metadata about the serving runtime.

Configuration belongs outside the bundle. Use the site settings API for non-secret runtime configuration, and keep secrets in the platform secret store.
`;

const fullDocMd = `# JavaScript API Reference for Hosted Sites

The hosted JavaScript API is deliberately small. It gives an app enough power to handle HTTP requests, render HTML, store local data, inspect runtime quota state, and serve static assets. It does not give the app raw filesystem access, subprocess execution, process environment access, or arbitrary host control. That boundary is the point of the platform: hosted apps should feel productive without becoming trusted host programs.

This reference describes the API as it exists today. Treat it as both documentation and a contract for bundle authors. If a function is not listed here, do not assume it is supported in hosted v1.

## Module overview

| Module | Import | Purpose |
| --- | --- | --- |
| Express router | \`require("express")\` | Register HTTP routes for the site. |
| UI DSL | \`require("ui.dsl")\` or \`require("ui")\` | Build HTML documents and fragments as JavaScript values. |
| Database | \`require("database")\` or \`require("db")\` | Query the site's preconfigured SQLite database. |
| DB guard | \`require("db.guard")\` | Inspect database quota state and register limit callbacks. |
| Utilities | \`require("path")\`, \`require("time")\`, \`require("timer")\` | Safe runtime utility modules from go-go-goja middleware. |

Unavailable by design:

\`\`\`js
require("fs");   // blocked in hosted v1
require("exec"); // blocked in hosted v1
\`\`\`

## The Express module

The Express-like router is the primary entry point for hosted apps.

\`\`\`js
const express = require("express");
const app = express.app();

app.get("/", (req, res) => {
  return ui.page({ title: "Hello" }, ui.main(ui.h1("Hello from go-go-host")));
});
\`\`\`

### Request object

\`\`\`js
app.get("/inspect", (req, res) => res.json({
  method: req.method,
  url: req.url,
  path: req.path,
  query: req.query,
  params: req.params,
  headers: req.headers,
  platform: req.platform
}));
\`\`\`

### Platform context

\`req.platform\` provides runtime metadata:

| Field | Type | Description |
| --- | --- | --- |
| \`requestId\` | string | Unique per-request ID |
| \`orgId\` | string | Organization ID |
| \`siteId\` | string | Site ID |
| \`deploymentId\` | string | Active deployment ID |
| \`host\` | string | Request Host header |

## UI DSL

\`\`\`js
const ui = require("ui.dsl");

ui.page({ title: "Catalog" },
  ui.main(
    ui.h1("Catalog"),
    ui.p("Choose an item."),
    ui.ul(
      ui.li(ui.a({ href: "/items/1" }, "Item 1")),
      ui.li(ui.a({ href: "/items/2" }, "Item 2"))
    )
  )
);
\`\`\`

## Error handling

\`\`\`js
app.get("/safe", (req, res) => {
  try {
    const rows = db.query("SELECT count(*) AS count FROM visits");
    return ui.page({ title: "Safe" }, ui.main(ui.p("Visits: " + rows[0].count)));
  } catch (err) {
    return res.status(500).json({ error: String(err) });
  }
});
\`\`\`

> **Warning:** Use \`ui.raw()\` only for trusted content. It bypasses escaping by design.

---

*See also: Developer Guide, Deploy Workflow, Agent Guide*
`;

const meta = {
  title: 'Molecules/MarkdownRenderer',
  component: MarkdownRenderer,
  args: { content: shortMd },
} satisfies Meta<typeof MarkdownRenderer>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Short: Story = {
  args: { content: shortMd },
};

export const CodeHeavy: Story = {
  args: { content: codeHeavyMd },
};

export const Table: Story = {
  args: { content: tableMd },
};

export const Headings: Story = {
  args: { content: headingsMd },
};

export const Blockquote: Story = {
  args: { content: blockquoteMd },
};

export const FullDoc: Story = {
  args: { content: fullDocMd },
};
