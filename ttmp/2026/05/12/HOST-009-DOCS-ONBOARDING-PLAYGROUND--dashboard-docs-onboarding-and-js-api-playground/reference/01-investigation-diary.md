---
Title: Investigation Diary
Ticket: HOST-009-DOCS-ONBOARDING-PLAYGROUND
DocType: reference
Status: active
Topics: [frontend, documentation, design-system, developer-experience]
---

# Investigation Diary

## 2026-05-12 — Ticket setup and initial brainstorm

User request: create a new ticket for dashboard documentation pages. The dashboard should become welcoming and intuitive for users, not only a management console. Users need to manage sites and agents, but more importantly they need to learn the hosted JavaScript APIs and try them out. User specifically suggested a live playground and a full runthrough walkthrough.

Created ticket:

```text
HOST-009-DOCS-ONBOARDING-PLAYGROUND
```

Created docs:

```text
design-doc/01-dashboard-docs-and-playground-brainstorm.md
reference/01-investigation-diary.md
```

Initial repo evidence:

- Existing CLI/developer docs live under `cmd/go-go-host/doc/`.
- The most relevant source documents are `developer-guide.md`, `js-api-reference.md`, `deploy-workflow.md`, `agent-guide.md`, and `agent-setup.md`.
- A durable runnable example exists under `examples/hello-beta/`.
- The hosted JS runtime modules live under `internal/sitejs/`, especially `web`, `uidsl`, and `dbguard`.

Initial design direction: the dashboard needs a first-class Docs/Learn area that combines prose, examples, copyable bundles, and a safe live playground. The playground should start as a client-side bundle authoring and preview environment, then graduate into controlled server-side dry-run/deploy flows.

## 2026-05-12 — reMarkable upload

Uploaded the brainstorm/design document to reMarkable:

```text
/ai/2026/05/12/HOST-009-DOCS-ONBOARDING-PLAYGROUND/HOST-009_Dashboard_Docs_Playground_Brainstorm.pdf
```

Command output:

```text
OK: uploaded HOST-009_Dashboard_Docs_Playground_Brainstorm.pdf -> /ai/2026/05/12/HOST-009-DOCS-ONBOARDING-PLAYGROUND
```

---

# Implementation Diary

## Goal

Build a MarkdownRenderer molecule, DocsIndexPage, and DocViewPage that render existing CLI docs from `cmd/*/doc/*.md` as dashboard documentation, with OS1-styled syntax highlighting and clipboard copy buttons on code blocks.

## Step 1: Add react-markdown and remark-gfm dependencies

Added `react-markdown` and `remark-gfm` as production dependencies to `web/admin`. These provide the markdown parsing and GFM (tables, strikethrough, autolinks) support.

**Commit (code):** `55695e6` — "Add react-markdown and remark-gfm"

### What I did
- Ran `pnpm add react-markdown remark-gfm` in `web/admin/`

### Why
- react-markdown is the standard React markdown renderer; remark-gfm adds table/strikethrough support

### What was tricky to build
- Nothing notable for this step

### What warrants a second pair of eyes
- N/A

### What should be done in the future
- Consider lazy-loading MarkdownRenderer if bundle size becomes a concern

## Step 2: Create MarkdownRenderer molecule

Created the `MarkdownRenderer` molecule under `web/admin/src/components/molecules/MarkdownRenderer/` with TSX, CSS, index, and Storybook stories. The CSS normalizes all markdown output (headings, prose, code, tables, lists, blockquotes, images) into the OS1 dashboard font scale with bordered/boxed code blocks and neutral body text.

**Commit (code):** `bf99367` — "Add MarkdownRenderer molecule and stories"

### What I did
- Created `MarkdownRenderer.tsx` with `react-markdown` + `remark-gfm`
- Created `MarkdownRenderer.css` with OS1-themed prose styling
- Created `index.ts` barrel export
- Added export to `molecules/index.ts`
- Created `MarkdownRenderer.stories.tsx` with 6 stories: Short, CodeHeavy, Table, Headings, Blockquote, FullDoc
- Captured Storybook screenshots

### Why
- A reusable molecule keeps docs rendering consistent across DocsIndexPage, DocViewPage, and future playground inline docs

### What was tricky to build
- Balancing the font scale hierarchy (h1–h6) within the small OS1 dashboard scale without making headings too dominant

### What warrants a second pair of eyes
- The heading size mapping: h1=15px (lg), h2=14px (md), h3-h6=13px (sm). Verify this reads well with real docs.

### What should be done in the future
- Add an optional table-of-contents sidebar for long docs
- Consider dark-mode / inverted syntax theme

## Step 3: Add syntax highlighting and code copy button

User explicitly requested syntax highlighting for code blocks and a clipboard copy button.

Added `rehype-highlight` + `highlight.js` for server-side-compatible syntax highlighting, and created `highlight-os1.css` with a light-background, neutral-ink, accent-colors-on-token-types theme that follows the OS1 dashboard rule (color accents on structured surfaces, not inline prose).

Added a `CodeBlockWithCopy` component that wraps each `<pre>` with a positioned copy button. The button fades in on hover, calls `navigator.clipboard.writeText()`, shows a checkmark for 1.5s, then reverts.

**Commit (code):** `284e862` — "Add Docs pages with MarkdownRenderer, syntax highlighting, and code copy buttons"

### What I did
- Added `rehype-highlight` and `highlight.js` dependencies
- Updated `MarkdownRenderer.tsx` to use `rehype-highlight` and override `pre` with `CodeBlockWithCopy`
- Created `highlight-os1.css` with token-type color mapping:
  - keywords: dark teal (#b028)
  - strings: deep blue (#047)
  - numbers: muted red-brown (#931)
  - types: warm brown (#530)
  - functions: blue-teal (#247)
- Added CSS for `.markdown-renderer__code-wrap`, `.markdown-renderer__copy-btn` with hover-reveal, active press, and copied-state styling
- Captured screenshots of code-heavy and full-doc stories

### Why
- Syntax highlighting makes code examples readable at a glance
- Copy button removes friction from following tutorials (users can paste CLI commands directly)
- OS1-themed highlight colors avoid the "modern SaaS neon syntax" look

### What was tricky to build
- TypeScript's `useRef<ReturnType<typeof setTimeout>>()` needs explicit `| null` and initial `null` value to satisfy strict mode; `clearTimeout(undefined)` errors because it expects exactly 1 argument
- react-markdown's `components` prop needs careful typing when overriding `pre`; the `children` prop comes from the rendered `<code>` element

### What warrants a second pair of eyes
- The copy button uses `navigator.clipboard.writeText()` which requires a secure context (HTTPS or localhost). In production over HTTPS this is fine; in HTTP-only dev it may silently fail.

### What should be done in the future
- Add a fallback for non-secure contexts (e.g. `document.execCommand('copy')`)
- Consider adding a language label badge on code blocks (e.g. "js", "json", "bash")

## Step 4: Create docs-data module

Created `web/admin/src/services/docs/docs-data.ts` that imports all 13 markdown files from `cmd/go-go-host/doc/` and `cmd/go-go-host-agent/doc/` as raw strings via Vite's `?raw` suffix. The module parses YAML frontmatter (Title, Slug, Short, SectionType) and exports:
- `docs: DocEntry[]` — sorted catalogue
- `docBySlug(slug)` — lookup
- `docsBySection()` — grouped for index page

Added `web/admin/src/vite-env.d.ts` with `declare module '*.md?raw'` for TypeScript.

### What I did
- Created `docs-data.ts` with 13 `?raw` imports
- Added frontmatter parser and body stripper
- Created typed `DocEntry` and `DocSection` interfaces

### Why
- Vite `?raw` imports embed the markdown at build time, no API needed
- Frontmatter parsing preserves existing metadata (Title, Slug, Short) from CLI docs

### What was tricky to build
- Root `.gitignore` excludes `data/` directories, so the module had to live under `services/docs/` instead of `data/`
- After moving the file deeper, relative paths to `cmd/*/doc/` needed an extra `../` level

### What warrants a second pair of eyes
- The `?raw` imports embed docs at build time; doc changes require a rebuild. This is acceptable for now but means docs are not dynamically updatable.

### What should be done in the future
- Consider an API endpoint that serves docs if dynamic updates are needed
- Add a `docs-data.test.ts` unit test for frontmatter parsing

## Step 5: Create DocsIndexPage and DocViewPage

Created two page components:
- `DocsIndexPage`: shows docs grouped by SectionType (Tutorials, Reference, etc.) with link cards
- `DocViewPage`: renders a single doc with breadcrumb, MarkdownRenderer, and prev/next navigation

Both pages use the OS1 dashboard-panel pattern and the dashboard font scale.

### What I did
- Created `DocsIndexPage.tsx` + `.css` + `.stories.tsx`
- Created `DocViewPage.tsx` + `.css` + `.stories.tsx` (5 stories: DeveloperGuide, JsApiReference, DeployWorkflow, AgentGuide, NotFound)
- Added routes `/app/orgs/:orgId/docs` and `/app/orgs/:orgId/docs/:slug` to `routes.tsx`
- Added `docs` section to `OrgSidebar.tsx` type and section list
- Updated `OrgLayout.tsx` to detect `/docs` in pathname

### Why
- The docs index gives a scannable overview grouped by purpose; the view page renders the actual content
- Sidebar integration makes docs a first-class nav item, not a hidden page

### What was tricky to build
- Storybook decorator for DocViewPage needs a `MemoryRouter` with the correct slug in the initial entry URL; the `Wrapper` component also needs to accept `children` prop for the decorator pattern
- `EmptyState` uses `body` not `message` prop; caught by `tsc`

### What warrants a second pair of eyes
- The prev/next navigation currently iterates the full flat `docs[]` array; it could be improved to navigate within the same section

### What should be done in the future
- Add a sidebar table of contents for long docs
- Add a search/filter input on the index page
- Add "agent" badge for agent-sourced docs

## Step 6: Screenshots and commit

Captured Storybook screenshots for:
- `host-009-markdown-short.png`
- `host-009-markdown-code-heavy.png`
- `host-009-markdown-table.png`
- `host-009-markdown-full-doc.png`
- `host-009-markdown-code-heavy-highlight.png`
- `host-009-markdown-copy-btn.png`
- `host-009-docs-index.png`
- `host-009-doc-view-js-api.png`

All committed in `284e862`.

## 2026-05-12 — Integrated docs via API

User requested integrating the docs into the main page with API route and MSW, instead of static Vite `?raw` imports.

Steps taken:

1. **Go backend**: Added `GET /api/v1/docs` and `GET /api/v1/docs/{slug}` endpoints in `internal/httpapi/docfs/docfs.go`. The package reads the existing `embed.FS` from `cmd/go-go-host/doc` and `cmd/go-go-host-agent/doc`, parses YAML frontmatter, and serves a sorted catalogue (list without body) and individual docs (with body). Exported `DocFS()` accessor functions on the existing `doc.go` packages.

2. **RTK Query**: Added `DocEntry` and `DocSection` types, `listDocs` and `getDoc` endpoints with `Docs` tag type, and exported `useListDocsQuery` / `useGetDocQuery` hooks.

3. **MSW**: Added `docs` fixture array with all 13 doc entries, plus `GET /api/v1/docs` and `GET /api/v1/docs/:slug` handlers. The slug handler returns a mock body with markdown for the view page.

4. **Pages rewrite**: `DocsIndexPage` and `DocViewPage` now use `useListDocsQuery()` and `useGetDocQuery()` respectively, with loading/error states.

5. **Removed static imports**: Deleted `web/admin/src/services/docs/docs-data.ts` and the `vite-env.d.ts` `?raw` type declarations. Bundle size dropped from 1148 KB to 1094 KB.

Commits:
- `ddcaf29` — Add Go docs API endpoints
- `5b7361d` — Integrate docs via API: RTK Query, MSW, page rewrites

## 2026-05-12 — Real app integration via devctl

Used `devctl up` to start the full dev stack (Postgres, Keycloak, go-go-hostd, web-admin Vite dev, Storybook).

Steps:
1. Started devctl: `devctl up --force --skip-build --skip-prepare`
2. Logged in as alice/alice via Keycloak
3. Created org "Demo Org" via API
4. Navigated to `/app/orgs/<orgId>/docs` — docs index renders with all 13 docs, grouped by Tutorials and Reference, "agent" badges on agent-sourced docs
5. Clicked into JS API Reference — full markdown body renders with syntax highlighting and copy buttons
6. Zero console errors

Fixed slug collision handling in `internal/httpapi/docfs/docfs.go`:
- Old approach: prefix every slug with `host-` or `agent-` → double-prefixes like `host-host-agent-guide`
- New approach: use slug from frontmatter as-is, only append `-agent` suffix when a slug collision is detected between host and agent docs
- Updated MSW fixtures and Storybook story slugs to match

Commits:
- `67a49c8` — Fix doc slug collision handling and verify real app integration
