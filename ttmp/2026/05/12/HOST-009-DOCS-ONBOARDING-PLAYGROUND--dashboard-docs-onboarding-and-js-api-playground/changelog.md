# Changelog

## 2026-05-12

- Initial workspace created


## 2026-05-12

Created initial dashboard docs/onboarding/playground brainstorm, including Learn navigation, full runthrough, JS API docs pages, playground MVP, and future dry-run design.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/cmd/go-go-host/doc/js-api-reference.md — Existing CLI-side JS API reference to adapt into dashboard docs
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/examples/hello-beta — Runnable example bundle that should seed the playground quickstart
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-009-DOCS-ONBOARDING-PLAYGROUND--dashboard-docs-onboarding-and-js-api-playground/design-doc/01-dashboard-docs-and-playground-brainstorm.md — Initial design brainstorm for user-facing Learn


## 2026-05-12

Uploaded the dashboard docs/playground brainstorm to reMarkable.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-009-DOCS-ONBOARDING-PLAYGROUND--dashboard-docs-onboarding-and-js-api-playground/design-doc/01-dashboard-docs-and-playground-brainstorm.md — Source design document uploaded as HOST-009_Dashboard_Docs_Playground_Brainstorm.pdf


## 2026-05-12

Added MarkdownRenderer molecule with syntax highlighting (rehype-highlight), clipboard copy button on code blocks, DocsIndexPage, DocViewPage, docs-data module importing cmd/*/doc/*.md via Vite ?raw, Docs nav in OrgSidebar, routing.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/app/routes.tsx — Routes for /docs and /docs/:slug
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/molecules/MarkdownRenderer/MarkdownRenderer.tsx — MarkdownRenderer molecule with rehype-highlight and code copy button
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/DocViewPage/DocViewPage.tsx — Doc view page rendering markdown with prev/next navigation
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/DocsIndexPage/DocsIndexPage.tsx — Docs index page showing grouped doc catalogue
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/services/docs/docs-data.ts — Docs data module importing CLI docs via Vite ?raw


## 2026-05-12

Integrated docs via Go API endpoints, RTK Query hooks, MSW handlers, and rewrote DocsIndexPage/DocViewPage to use API data. Removed static ?raw imports.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/cmd/go-go-host-agent/doc/doc.go — Exported DocFS() accessor
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/cmd/go-go-host/doc/doc.go — Exported DocFS() accessor
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/internal/httpapi/docfs/docfs.go — Go docs API handler reading embedded CLI docs
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/DocViewPage/DocViewPage.tsx — Rewritten to use useGetDocQuery
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/DocsIndexPage/DocsIndexPage.tsx — Rewritten to use useListDocsQuery
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/services/goGoHostApi.ts — Added listDocs and getDoc RTK Query endpoints

