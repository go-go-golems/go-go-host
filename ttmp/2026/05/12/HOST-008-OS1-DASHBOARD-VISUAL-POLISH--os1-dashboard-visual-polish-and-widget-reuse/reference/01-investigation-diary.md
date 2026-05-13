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

## 2026-05-12 — CodeMirror JSON editor atom

Installed official CodeMirror packages in `web/admin`:

```text
@codemirror/commands
@codemirror/lang-json
@codemirror/state
@codemirror/view
```

Added reusable atom:

```text
web/admin/src/components/atoms/JsonEditor/JsonEditor.tsx
web/admin/src/components/atoms/JsonEditor/JsonEditor.css
web/admin/src/components/atoms/JsonEditor/index.ts
```

The component creates an `EditorView` with JSON language support, line numbers, history, default keymap, tab indentation, line wrapping, and a small OS1-compatible theme layer. It is controlled from React by syncing external `value` changes into the editor while reporting CodeMirror edits back through `onChange`.

Validation: `pnpm build` passed.

Commit: `36b846e Add CodeMirror JSON editor atom`.

## 2026-05-12 — Settings page implementation

Implemented the requested settings-page slice.

Code changes:

```text
web/admin/src/pages/SiteSettingsPage/SiteSettingsPage.tsx
web/admin/src/pages/SiteSettingsPage/SiteSettingsPage.css
web/admin/src/app/macos1-bridge.css
web/admin/src/components/atoms/JsonEditor/JsonEditor.tsx
web/admin/src/components/atoms/JsonEditor/JsonEditor.css
web/admin/package.json
web/admin/pnpm-lock.yaml
```

What changed:

- Replaced the editable JSON `<textarea>` with the reusable CodeMirror `JsonEditor`.
- Added explicit CodeMirror JSON syntax highlighting through `HighlightStyle`/`syntaxHighlighting`, covering property names, strings, numbers, booleans, nulls, and punctuation.
- Replaced capability Enable/Disable button semantics with `@go-go-golems/os-core` `Checkbox` widgets in the settings capability table. The checkbox toggles the same API mutation and is disabled while the mutation is in flight.
- Added subtle semantic highlights:
  - green for safe/non-secret/supported values,
  - blue for informational tokens,
  - yellow for warnings/deferred features,
  - red for unavailable/dangerous capabilities such as `exec`, `fs`, and secrets.
- Normalized settings-page font sizes around a small OS1 scale (`--ggh-font-xs`, `--ggh-font-sm`, `--ggh-font-md`, `--ggh-font-lg`, `--ggh-font-code`).
- Overrode settings page headers to stack title/body copy instead of inheriting the generic dashboard header's split left/right layout; this makes explanatory text easier to read.

Validation:

```bash
cd web/admin
pnpm build
```

Build passed. The build now emits a chunk-size warning because CodeMirror increases the admin bundle from roughly 445 KB minified JS to roughly 752 KB. This is acceptable for the current beta slice, but we should consider lazy-loading the settings page or the JSON editor before broad production use.

Playwright screenshot captured:

```text
sources/screenshots/host-008-site-settings-final.png
```

Commit: `3def444 Polish settings page OS1 controls`.

## 2026-05-12 — Admin UI guidelines playbook

User changed direction from continuing implementation to writing a reusable playbook/guideline document for future admin-backend UI work, then uploading it to reMarkable.

Created playbook:

```text
playbook/01-os1-admin-dashboard-ui-work-guidelines.md
```

The playbook captures the rules learned so far:

- use the OS1 shell/window mental model;
- do not merely apply retro colors to a modern SaaS layout;
- use OS-core theme variables and `data-part` contracts;
- use the dashboard font scale;
- keep whitespace structural;
- use semantic highlights subtly;
- use OS checkbox widgets for booleans;
- use CodeMirror `JsonEditor` for editable JSON;
- build Storybook/MSW states before styling;
- capture Playwright screenshots and review full pages before moving on.

Next step: upload the playbook to reMarkable.

## 2026-05-12 — reMarkable upload

Uploaded the guidelines playbook to reMarkable:

```text
/ai/2026/05/12/HOST-008-OS1-DASHBOARD-VISUAL-POLISH/HOST-008_OS1_Admin_UI_Guidelines.pdf
```

Command output:

```text
OK: uploaded HOST-008_OS1_Admin_UI_Guidelines.pdf -> /ai/2026/05/12/HOST-008-OS1-DASHBOARD-VISUAL-POLISH
```

## 2026-05-12 — Agents page cleanup scope

Continuing HOST-008 after writing the UI playbook. Next target: `AgentsPage`, because it is security-sensitive and currently violates several playbook rules:

- native checkbox for auto-activation instead of OS-core `Checkbox`,
- warning copy is a broad yellow strip but not semantically precise,
- form is a loose flex row rather than a compact OS1 control group,
- tables use different border weights than Sites/Settings tables,
- signing-key status styles use ad-hoc colors and large font sizes,
- page headers do not use the improved stacked explanatory-text pattern.

Plan:

1. Replace the auto-activation native checkbox with `@go-go-golems/os-core` `Checkbox`.
2. Add subtle semantic highlights for machine identity, auto-activation, trusted pipelines, enrollment tokens, and revocation.
3. Normalize Agents page font sizes and table density to the dashboard scale.
4. Align `AgentsTable` and `AgentKeysTable` with the same 2px outer border / black header pattern used by Sites and Settings.
5. Capture a Storybook screenshot of the populated Agents page.

## 2026-05-12 — Agents page implementation

Implemented the Agents page OS1 cleanup.

Changed files:

```text
web/admin/src/pages/AgentsPage/AgentsPage.tsx
web/admin/src/pages/AgentsPage/AgentsPage.css
web/admin/src/components/organisms/AgentsTable/AgentsTable.css
web/admin/src/components/organisms/AgentKeysTable/AgentKeysTable.css
```

What changed:

- Replaced the native auto-activation checkbox with `@go-go-golems/os-core` `Checkbox`.
- Reworked page copy with subtle semantic highlights:
  - machine identity = info,
  - auto-activation = danger,
  - trusted pipelines = safe.
- Converted the create form into a compact OS1 control group with bordered inputs and a warning panel for auto-activation.
- Added explicit sections for agent records and signing keys with stacked explanatory headers.
- Normalized Agents and Agent Keys tables to the 2px OS1 outer border / black header / dense cell pattern.
- Reworked selected row and key status states with subtle OS1-compatible color fills.

Validation:

```bash
cd web/admin
pnpm build
```

Build passed, with the already-known CodeMirror chunk-size warning.

Storybook screenshot captured via Playwright:

```text
sources/screenshots/host-008-agents-page-os1.png
```

## 2026-05-12 — Cross-page normalization and color-accent correction

User clarified an important visual rule: color accents should be used for badge colors and buttons, not really for inline body text.

Adjusted direction accordingly:

- updated the OS1 UI playbook to discourage colored inline prose highlights;
- kept legacy `__highlight` classes neutral/bold for compatibility instead of colored backgrounds;
- moved color emphasis into `StatusPill`, `RoleBadge`, `MetricCard` values, copy-button states, warning/error panels, and table status treatments.

Implemented a cross-page CSS normalization pass rather than page-by-page rewrites for every simple inventory page:

- normalized `dashboard-panel` header behavior, form controls, links, admin inventory tables, runtime panels, upload panels, timeline lists, code blocks, empty/error states, and badges;
- tightened Audit and Create Site page layouts;
- normalized Admin Overview/Runtimes/Inventory page font scale and table density;
- normalized DeploymentTimeline and AuditTimeline away from harsh invert hover and toward OS1 bordered rows;
- normalized RuntimeStatusPanel and DeploymentUploadPanel to 2px OS1 panel styling;
- converted text highlight classes in Settings/Agents to neutral bold emphasis.

Validation:

```bash
cd web/admin
pnpm build
```

Build passed with the known CodeMirror chunk-size warning.

Captured Storybook screenshots for representative remaining pages:

```text
sources/screenshots/host-008-admin-overview-os1-pass.png
sources/screenshots/host-008-admin-orgs-os1-pass.png
sources/screenshots/host-008-admin-deployments-os1-pass.png
sources/screenshots/host-008-admin-audit-os1-pass.png
sources/screenshots/host-008-audit-os1-pass.png
sources/screenshots/host-008-create-site-os1-pass.png
sources/screenshots/host-008-site-overview-os1-pass.png
sources/screenshots/host-008-deployments-os1-pass.png
sources/screenshots/host-008-usage-os1-pass.png
sources/screenshots/host-008-members-os1-pass.png
```

These screenshots do not mean every page is perfect, but they cover the main page families and verify the shared normalization pass applies broadly.

## 2026-05-12 — White page background correction

User asked whether all pages/organisms are done and clarified that the page background should be white, not light blue.

I checked coverage quickly:

- all organisms have CSS coverage;
- several simple pages still do not have dedicated page CSS and rely on shared dashboard normalization (`AuthCallbackPage`, `RuntimePage`, `UsagePage`, `MembersPage`, `DeploymentsPage`, `DeploymentDetailPage`, and several simple admin inventory pages);
- therefore the work is not honestly "all pages individually polished" yet, but the shared normalization covers the main page families and representative Storybook screenshots have been captured.

Implemented the background correction:

- set the OS1 theme root `--hc-color-desktop-bg` override to white in `macos1-bridge.css`;
- changed `.app-shell` and `.app-shell__body` backgrounds to white and disabled the inherited desktop checker/light-blue look.

Validation:

```bash
cd web/admin
pnpm build
```

Build passed with the known CodeMirror chunk-size warning.
