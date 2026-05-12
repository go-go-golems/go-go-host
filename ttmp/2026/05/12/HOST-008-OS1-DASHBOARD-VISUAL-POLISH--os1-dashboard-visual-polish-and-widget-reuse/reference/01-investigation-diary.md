---
Title: Investigation Diary
Ticket: HOST-008-OS1-DASHBOARD-VISUAL-POLISH
DocType: reference
Status: active
Topics: [frontend, design-system, storybook, debugging]
---

# Investigation Diary

## 2026-05-12 — Ticket setup and first visual baseline

User request: the public dashboard visually looks bad despite already importing `@go-go-golems/os-core` assets; make it feel like a unified macOS/OS1 site with consistent styling and margins, using Playwright/screenshots and the available npm packages from `go-go-os-frontend`.

Created docmgr ticket:

```text
HOST-008-OS1-DASHBOARD-VISUAL-POLISH
```

Created documents:

```text
design-doc/01-os1-dashboard-visual-redesign-investigation.md
reference/01-investigation-diary.md
```

Captured baseline screenshots:

```text
sources/screenshots/host-008-current-app.png
sources/screenshots/host-008-current-auth-app-viewport.png
sources/screenshots/host-008-go-go-os-examples-home.png
sources/screenshots/host-008-go-go-os-examples-rich-widgets.png
```

Findings from `https://hosting.yolo.scapegoat.dev/app`:

- The unauthenticated/early-loading state can show a giant flat blue rectangle without enough chrome or content hierarchy.
- The authenticated Sites view has the right raw tokens — bitmap-ish font, black borders, white panels — but it is compositionally weak:
  - top bar, sidebar, and content table do not feel like one OS surface;
  - margins are too tight in some places and too huge/empty in others;
  - the main content is a flat panel rather than a Mac window surface;
  - the sidebar is just buttons on a white block, not a window/navigation palette;
  - table cells and copy controls crowd each other;
  - the active site table works functionally but does not read as deliberate OS1 design.

Comparison target screenshots from `https://go-go-os-examples.yolo.scapegoat.dev/`:

- `host-008-go-go-os-examples-home.png`
- `host-008-go-go-os-examples-rich-widgets.png`

The examples site looks more coherent because it uses:

- a constrained centered desktop work area,
- consistent left navigation column width,
- strong black window borders,
- striped active title bars,
- regular spacing between navigation, hero panels, and content panels,
- OS1 button/data-table primitives rather than ad-hoc modern pill styling.

## 2026-05-12 — Code inspection

Inspected dashboard frontend files:

```text
web/admin/src/main.tsx
web/admin/src/app/providers/AppProviders.tsx
web/admin/src/app/macos1-bridge.css
web/admin/src/components/organisms/AppShell/AppShell.tsx
web/admin/src/components/organisms/AppShell/AppShell.css
web/admin/src/components/organisms/OrgSidebar/OrgSidebar.tsx
web/admin/src/components/organisms/OrgSidebar/OrgSidebar.css
web/admin/src/pages/SitesPage/SitesPage.tsx
web/admin/src/components/organisms/SitesTable/SitesTable.tsx
web/admin/src/components/organisms/SitesTable/SitesTable.css
```

Important finding: the dashboard already imports the right base package:

```ts
import '@go-go-golems/os-core/theme';
import '@go-go-golems/os-core/desktop-theme-macos1';
```

and wraps the app in:

```tsx
<div data-widget="hypercard" className="theme-macos1">
```

So the problem is not absence of the OS1 package. The problem is incomplete composition. The app consumes variables such as `--hc-color-bg`, but it does not consistently use the structural parts from the OS desktop shell: window frame, title bar, window body, menu bar spacing, data-table density, and desktop workspace constraints.

The installed `@go-go-golems/os-core` package exports `desktop-react` primitives including `WindowSurface`, `DesktopMenuBar`, `DesktopShell`, and lower-level parts. For this first slice I avoided a full desktop-shell migration and instead reused the package's `data-part` contract plus CSS imported by `@go-go-golems/os-core/theme`.

## 2026-05-12 — First visual-polish slice

Implemented a first local polish slice focused on global shell consistency:

- rewrote `AppShell.css` to make the app a constrained centered OS desktop workspace;
- made the top bar sticky and consistently bordered like a menu bar;
- converted the org sidebar into a static OS1 window using `data-part="windowing-window"`, `data-part="windowing-window-title-bar"`, and `data-part="windowing-window-body"`;
- added a Navigation title bar to the sidebar in `AppShell.tsx`;
- expanded `macos1-bridge.css` so `.dashboard-panel` renders as a Mac-style content window with striped title bar and consistent padding;
- tightened typography, table header casing, button/input font inheritance, and general panel margins;
- normalized the Sites table border, padding, wrapping, hover state, and shadow;
- made site host copy controls wrap cleanly.

Validation:

```bash
cd web/admin
pnpm build
```

Result: build passed.

Started local Storybook on port 6008 and captured a post-change shell screenshot:

```text
sources/screenshots/host-008-appshell-after-local.png
```

This screenshot is not the final desired landing page, but it proves the direction: the navigation now reads as a Mac window/palette, the main dashboard panel has a title-bar treatment, the desktop area is constrained, and margins are far more deliberate than the public baseline.

## Current next steps

1. Apply the same treatment to remaining major pages and organisms, especially Agents, Audit, Deployment detail, Site settings, admin pages, forms, dialogs, and empty/loading/error states.
2. Add explicit dashboard-specific shell primitives rather than relying on broad `.dashboard-panel::before` CSS for all panels.
3. Consider adding `@go-go-golems/os-widgets` if we want richer prebuilt surfaces, but keep the dashboard mostly on `@go-go-golems/os-core` primitives for app chrome.
4. Capture authenticated local/full-page screenshots after wiring MSW or a local backend so the final visual test covers the real Sites page, not only `AppShell` Storybook.
5. Deploy a new dashboard image and recapture `https://hosting.yolo.scapegoat.dev/app` after rollout.

## 2026-05-12 — Full-page Storybook screenshot and grid whitespace fix

User asked to inspect a full page before moving further and correctly noticed a large blank band between the Sites page title and the Sites table.

Captured the first full-page Storybook Sites view:

```text
sources/screenshots/host-008-sitespage-full-after-local.png
```

Root cause: `.dashboard-panel` had been converted to `display: grid` and also had a large `min-height: 28rem`. CSS Grid's default alignment stretches auto tracks to fill available block space, so the header row expanded vertically and pushed the table downward. This looked like arbitrary whitespace between the title and table.

Fix applied in `web/admin/src/app/macos1-bridge.css`:

```css
.dashboard-panel {
  display: grid;
  grid-auto-rows: max-content;
  align-content: start;
  min-height: 0;
}
```

Captured the corrected full-page screenshot:

```text
sources/screenshots/host-008-sitespage-full-gap-fixed.png
```

Also confirmed the MSW warning noise was not only missing API handlers: `msw-storybook-addon` was warning about Storybook/Vite static module and CSS requests. Updated `.storybook/preview.tsx` so only unhandled `/api/...` requests warn. Added missing agent key/enrollment MSW routes for the Agents page:

```text
POST /api/v1/orgs/:orgId/agents/:agentId/enrollment-token
GET  /api/v1/orgs/:orgId/agents/:agentId/keys
POST /api/v1/orgs/:orgId/agents/:agentId/keys/:keyId/revoke
```

Validation: `pnpm build` passed after the CSS/MSW changes. The only remaining browser console error in the Storybook full-page screenshot was `favicon.ico` 404, not an API/MSW route issue.

## 2026-05-12 — Settings-page polish scope

User requested the next focused slice on the settings page:

- use proper checkbox widgets,
- add subtle color highlights in text where appropriate,
- unify font sizes,
- use CodeMirror with JSON syntax highlighting for JSON fields,
- add ticket tasks first, then implement with commits and diary updates.

Added four HOST-008 tasks for checkbox widgets, semantic highlights, font scale unification, and CodeMirror JSON editing.

Inspected `SiteSettingsPage.tsx` and `SiteSettingsPage.css`. Current state:

- the settings page is structurally useful but visually inconsistent;
- the only native checkbox in the page is inside the agent create flow, but settings has capability toggles that behave like boolean controls via Enable/Disable buttons;
- editable JSON config uses a plain `<textarea>`;
- JSON display uses `CodeBlock`, which is readable but not editable/syntax-highlighted;
- helper text uses mostly muted gray and misses semantic emphasis for safe vs unsupported capability/environment language;
- settings CSS still has `font: inherit` and mixed fallbacks instead of a small OS1 scale.

Plan for this slice:

1. Add a small reusable CodeMirror JSON editor component under `components/atoms/JsonEditor` using official CodeMirror packages.
2. Replace the settings JSON textarea with that editor.
3. Render capability policy rows with OS-core `Checkbox` widgets so enabled/disabled state reads as a boolean control, while still preserving async save behavior.
4. Add subtle highlight classes for safe, warning, danger, and identifier text.
5. Normalize settings-page font sizes and table/control density.
