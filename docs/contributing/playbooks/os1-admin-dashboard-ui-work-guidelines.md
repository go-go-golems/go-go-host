---
Title: OS1 Admin Dashboard UI Work Guidelines
Ticket: HOST-008-OS1-DASHBOARD-VISUAL-POLISH
Status: active
Topics:
    - frontend
    - design-system
    - storybook
    - debugging
DocType: playbook
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/go-go-os-frontend/packages/os-core/src/theme/desktop/primitives.css
      Note: Reference OS-core primitive styling for buttons
    - Path: web/admin/src/app/macos1-bridge.css
      Note: Defines dashboard-wide OS1 bridge tokens
    - Path: web/admin/src/components/atoms/JsonEditor/JsonEditor.tsx
      Note: Canonical editable JSON field component referenced by the playbook
    - Path: web/admin/src/pages/SiteSettingsPage/SiteSettingsPage.tsx
      Note: Representative implementation of checkbox widgets
ExternalSources: []
Summary: ""
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: ""
WhenToUse: ""
---


# OS1 Admin Dashboard UI Work Guidelines

## Purpose

This playbook explains how to design and implement new `go-go-host` admin/dashboard pages so they look like one coherent OS1 application from the first commit. It captures the lessons from the first dashboard visual-polish pass: importing the OS theme is not enough; every page must use the same shell model, spacing model, typography scale, state styling, Storybook/MSW workflow, and screenshot review loop.

The target visual language is **classic Mac / OS1 admin console**: black outlines, striped title bars, compact bitmap-style typography, dense but readable tables, explicit panels, a white page background, and careful use of subtle state color on badges/buttons/panels. The goal is not ornamental nostalgia. The goal is a legible operator UI where every control looks intentional and every page follows the same rules.

## Environment assumptions

Work from the app repo:

```bash
cd /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host
```

Frontend workspace:

```bash
cd web/admin
```

Theme packages already imported by the app:

```ts
import '@go-go-golems/os-core/theme';
import '@go-go-golems/os-core/desktop-theme-macos1';
```

The React app is wrapped in the OS1 theme scope:

```tsx
<div data-widget="hypercard" className="theme-macos1">
  {children}
</div>
```

Useful local commands:

```bash
pnpm build
pnpm storybook
```

Primary screenshot tools:

- Storybook iframe pages for deterministic component states.
- Playwright screenshots for full-page visual review.
- Ticket-local screenshots under:

```text
ttmp/2026/05/12/HOST-008-OS1-DASHBOARD-VISUAL-POLISH--os1-dashboard-visual-polish-and-widget-reuse/sources/screenshots/
```

## The design model

Every dashboard page should be designed as a windowed OS workspace, not as a generic web-admin page.

The mental model is:

```text
OS1 theme root
  └── AppShell
      ├── sticky menu bar
      └── centered desktop work area
          ├── navigation window/palette
          └── page content windows
              ├── title/header
              ├── toolbar/form/filter rows
              ├── tables, editors, panels
              └── status/help text
```

A page should answer three visual questions immediately:

1. **Where am I?** The window title and page heading identify the current area.
2. **What can I do?** Primary actions sit in predictable toolbar/header positions.
3. **What needs attention?** Warnings, dangerous options, unavailable features, and safe states are visible through badges, buttons, and panels rather than colored body text.

## Core rules

### 1. Do not build a modern SaaS page inside retro colors

Avoid default web-admin habits:

- oversized cards,
- pill-heavy controls,
- large border radii,
- airy whitespace with no structural meaning,
- unstyled browser checkboxes,
- giant muted helper text,
- random font sizes from browser defaults,
- tables without strong grid structure.

Use OS1 structure instead:

- 2px black borders for major surfaces,
- striped title bars for windows,
- compact spacing,
- dense tables,
- square controls,
- explicit status badges,
- minimal but meaningful color.

### 2. Use the OS theme variables and parts

Prefer variables from the OS theme:

```css
var(--hc-color-bg)
var(--hc-color-fg)
var(--hc-color-border)
var(--hc-color-muted)
var(--hc-color-desktop-bg)  /* overridden to white in the admin dashboard */
var(--hc-window-body-bg)
var(--hc-window-shadow)
var(--hc-window-border-radius)
var(--hc-font-family)
```

Prefer `data-part` styling where an OS-core primitive exists:

```tsx
<button data-part="btn">Save</button>
<input data-part="field-input" />
<select data-part="field-select" />
```

For window-like local structures, reuse OS-core part names when appropriate:

```tsx
<section data-part="windowing-window" data-state="focused">
  <div data-part="windowing-window-title-bar" data-state="focused">
    <span data-part="windowing-close-button" />
    <span data-part="windowing-window-title">Navigation</span>
  </div>
  <div data-part="windowing-window-body">...</div>
</section>
```

Do not copy large chunks of OS-core CSS into page CSS. Use the imported theme and override only local layout needs.

### 2a. Keep the page background white

The admin dashboard uses OS1 window chrome, but the page behind those windows should be white, not the default light-blue/checker desktop from the reusable OS shell. This keeps the hosted-control-plane dashboard looking like a document/operator console rather than a toy desktop.

The dashboard bridge should explicitly override the desktop background token:

```css
[data-widget="hypercard"].theme-macos1 {
  --hc-color-desktop-bg: #fff;
  background: #fff;
}

.app-shell,
.app-shell__body {
  background: #fff;
  background-image: none;
}
```

Use white for the page/workspace background. Use OS1 borders, title bars, and shadows to create structure. Do not use the light-blue desktop fill as the global admin background.

### 3. Use the dashboard font scale

Keep typography compact and consistent. The dashboard bridge defines a small OS1 scale:

```css
--ggh-font-xs: 10px;
--ggh-font-sm: 11px;
--ggh-font-md: 12px;
--ggh-font-lg: 16px;
--ggh-font-code: 11px;
```

Use it like this:

```css
.page-root {
  font-size: var(--ggh-font-sm);
}

.page-root h1 { font-size: var(--ggh-font-lg); }
.page-root h2 { font-size: 14px; }
.page-root h3 { font-size: var(--ggh-font-md); }
.page-root code,
.page-root pre { font-size: var(--ggh-font-code); }
```

Avoid ad-hoc `0.875rem`, `1.25rem`, and browser-default sizes unless there is a deliberate reason.

### 4. Make headers readable before clever

The generic dashboard window header has split layout for pages with a title and right-side actions. Some pages, especially settings/instructions pages, need stacked explanatory text. Override locally when necessary:

```css
.my-page .dashboard-panel > header:first-child {
  display: grid !important;
  grid-template-columns: 1fr;
  gap: 0.35rem;
  align-items: start;
}
```

Do not let helper text drift to the far right just because a generic flex header inherited `justify-content: space-between`.

### 5. Keep whitespace structural

Whitespace should mean grouping, not accidental expansion.

Common failure:

```css
.dashboard-panel {
  display: grid;
  min-height: 28rem;
}
```

CSS Grid may stretch rows and create huge gaps between header and content. Use:

```css
.dashboard-panel {
  display: grid;
  grid-auto-rows: max-content;
  align-content: start;
  gap: 1rem;
}
```

Before accepting a page, inspect a full-page screenshot and ask:

- Is empty space bounded by a window or section?
- Does it separate related groups?
- Would a user think the page failed to load?
- Is a table pushed down by accidental grid/flex stretching?

### 6. Use color accents on badges, buttons, and panels — not inline prose

Color should annotate UI state. It should not turn body copy into a rainbow. Inline explanatory text should usually remain black/neutral; use plain `strong` emphasis for important terms.

Preferred color locations:

| Tone | Use for |
|---|---|
| safe | status pills, success buttons, verified-domain badges, enabled capability badges |
| info | informational badges, selected rows, copy-success buttons |
| warning | warning callouts, pending badges, manual placeholder panels |
| danger | destructive buttons, error callouts, revoked/failed/rejected badges |

Good:

```tsx
<StatusPill status="verified" tone="success" />
<button data-part="btn" data-variant="danger">Revoke</button>
<div className="warning-panel">Auto-activation requires a trusted pipeline.</div>
```

Avoid:

```tsx
<p>Store <span className="green-highlight">non-secret values</span> here.</p>
```

If a term in prose needs emphasis, use neutral emphasis instead:

```tsx
<p>Store <strong>non-secret values</strong> here.</p>
```

The current dashboard CSS keeps legacy `__highlight` classes neutral/bold for compatibility, but new pages should avoid colored inline text highlights.

### 7. Use OS widgets for boolean controls

Do not use native browser checkboxes directly in visible page UI unless they are wrapped and styled intentionally.

Use OS-core checkbox widgets:

```tsx
import { Checkbox } from '@go-go-golems/os-core';

<Checkbox
  label={enabled ? 'Enabled' : 'Disabled'}
  checked={enabled}
  disabled={isSaving}
  onChange={() => void toggleEnabled(!enabled)}
/>
```

For policy rows, the checkbox should be the control, and the status pill should be the readout:

```text
Capability | Policy checkbox | Status pill | Config
```

Do not show both a checkbox and an unrelated Enable/Disable button unless there are two distinct actions.

### 8. Use CodeMirror for editable JSON

Editable JSON fields should use the `JsonEditor` atom, not `<textarea>`.

Use:

```tsx
import { JsonEditor } from '../../components/atoms';

<JsonEditor
  value={configValue}
  onChange={setConfigValue}
  ariaLabel="Site config JSON value"
/>
```

The `JsonEditor` atom provides:

- JSON language mode,
- line numbers,
- history,
- default keymap,
- tab indentation,
- line wrapping,
- OS1 border/background styling,
- semantic syntax highlighting.

Keep JSON parsing and save behavior in the page:

```tsx
try {
  const value = JSON.parse(configValue);
  await upsertConfig({ siteId, key, value }).unwrap();
} catch (error) {
  setConfigError(error instanceof SyntaxError ? error.message : apiErrorMessage(error));
}
```

Bundle-size note: CodeMirror is not tiny. If more heavy editors are added, consider lazy-loading the settings/editor surface.

### 9. Tables should be dense, bordered, and predictable

Tables are operator surfaces. They should not be airy cards.

Use:

```css
table {
  width: 100%;
  border-collapse: collapse;
  border: 2px solid var(--hc-color-border);
  background: var(--hc-color-bg);
}

th,
td {
  border: 1px solid var(--hc-color-border);
  padding: 0.5rem;
  vertical-align: top;
  text-align: left;
}

th {
  background: var(--hc-color-fg);
  color: var(--hc-color-bg);
  text-transform: uppercase;
  letter-spacing: 0.03em;
}
```

Wrap long IDs/hosts deliberately:

```css
td code {
  overflow-wrap: anywhere;
}
```

### 10. Prefer page-specific polish before broad global rules

A bridge file such as `macos1-bridge.css` is useful for legacy mapping, but avoid dumping every future UI rule there. For a new page:

1. Use shared dashboard primitives when available.
2. Add page-local CSS for layout details.
3. Promote repeated patterns to shared primitives only after two or three pages need them.

This avoids global CSS surprises.

## New page workflow

### Step 1: Start from a Storybook story

Every page should have a Storybook story with realistic MSW data.

Pattern:

```tsx
import type { Meta, StoryObj } from '@storybook/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { MyPage } from './MyPage';

const meta = { title: 'Pages/MyPage', component: MyPage } satisfies Meta<typeof MyPage>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Populated: Story = {};
export const Empty: Story = { /* MSW override */ };
export const LoadError: Story = { /* MSW override */ };
```

For route params, use `MemoryRouter`. For outlet context, add a small story shell.

### Step 2: Add/verify MSW handlers first

Before styling, make the story render real data. Check the page's RTK Query endpoints in:

```text
web/admin/src/services/goGoHostApi.ts
```

Then add handlers in:

```text
web/admin/src/services/msw/handlers.ts
web/admin/src/services/msw/fixtures.ts
```

Storybook MSW should warn only for unhandled `/api/...` requests. Static Vite/CSS requests are intentionally ignored.

### Step 3: Compose the page as windows and sections

Use the page root plus `.dashboard-panel` sections:

```tsx
<div className="my-page">
  <section className="dashboard-panel my-page__intro">
    <header>
      <h1>Page title</h1>
      <p>Short operator explanation with semantic highlights.</p>
    </header>
  </section>

  <section className="dashboard-panel my-page__section">
    <header>
      <h2>Section title</h2>
      <p>What this section controls.</p>
    </header>
    ...
  </section>
</div>
```

Do not cram unrelated controls into one giant panel. Prefer multiple compact windows/sections.

### Step 4: Use the right controls

| Need | Preferred control |
|---|---|
| Boolean policy | `Checkbox` from `@go-go-golems/os-core` |
| Primary action | `<button data-part="btn">` |
| Text field | `<input data-part="field-input">` |
| Select | `<select data-part="field-select">` |
| Editable JSON | `JsonEditor` atom |
| Read-only JSON | `CodeBlock` or a future read-only syntax viewer |
| Status | `StatusPill` with explicit tone |
| Date | `Timestamp` |
| Copy action | `CopyButton` |

### Step 5: Normalize page CSS

Start with:

```css
.my-page {
  display: grid;
  gap: 1rem;
  font-size: var(--ggh-font-sm);
}

.my-page__section {
  overflow-x: auto;
}
```

Then add page-specific grid/table/form details.

Do not rely on browser default margins. Set margins deliberately.

### Step 6: Validate with build and screenshot

Run:

```bash
cd web/admin
pnpm build
pnpm storybook
```

Use Playwright to open the story iframe and capture a screenshot:

```text
http://127.0.0.1:<storybook-port>/iframe.html?id=<story-id>&viewMode=story
```

Save the screenshot into the ticket:

```text
ttmp/.../sources/screenshots/<ticket>-<page>-<state>.png
```

At minimum capture:

- populated state,
- empty state,
- error state if visually distinct,
- one narrow/mobile-ish viewport if layout is non-trivial.

### Step 7: Review visually before moving on

Use this checklist:

- [ ] Does the page look like OS1, not modern SaaS?
- [ ] Are all major surfaces windowed or clearly grouped?
- [ ] Is there accidental whitespace?
- [ ] Are font sizes consistent?
- [ ] Are tables dense and readable?
- [ ] Are dangerous/unavailable states shown through badges/buttons/panels rather than colored body text?
- [ ] Is the page/workspace background white?
- [ ] Are controls square and theme-consistent?
- [ ] Are JSON fields using CodeMirror if editable?
- [ ] Are boolean controls using OS checkbox widgets?
- [ ] Does Storybook render without unhandled `/api/...` MSW warnings?
- [ ] Does `pnpm build` pass?

## Common failure modes

### Huge whitespace between title and content

Cause: grid/flex stretching, often from `min-height` or default `align-content: stretch`.

Fix:

```css
grid-auto-rows: max-content;
align-content: start;
min-height: 0;
```

### Storybook says a story does not exist

Check the real story ID in:

```text
http://127.0.0.1:<port>/index.json
```

Story IDs are generated from title and export name. For example:

```text
Pages/SitesPage + Populated -> pages-sitespage--populated
```

### MSW warns about CSS files

That should be suppressed by Storybook preview config. Only `/api/...` unhandled requests matter for mocks.

### CodeMirror makes the bundle large

This is expected. Consider lazy loading if the warning becomes a real production concern:

```tsx
const JsonEditor = lazy(() => import('../../components/atoms/JsonEditor'));
```

Do not prematurely split until the route-level UX needs it.

### Page-specific header text is pushed to the far right

The generic dashboard header supports title-left/action-right layouts. Override the page header to stack content:

```css
.my-page .dashboard-panel > header:first-child {
  display: grid !important;
  grid-template-columns: 1fr;
}
```

## Exit criteria for new admin UI work

A new dashboard/admin page is ready when:

1. It has Storybook populated/empty/error states.
2. All required `/api/...` handlers exist in MSW.
3. It uses OS1 theme variables and controls.
4. Boolean controls use OS checkbox widgets.
5. Editable JSON uses `JsonEditor`.
6. Font sizes follow the dashboard scale.
7. Color accents are used for warnings/safe/danger states on badges, buttons, and panels — not inline prose.
8. Full-page screenshot has been reviewed and saved to the ticket.
9. `pnpm build` passes.
10. Diary/changelog/tasks are updated.

## Files to check before designing a new page

```text
web/admin/src/app/macos1-bridge.css
web/admin/src/components/organisms/AppShell/AppShell.tsx
web/admin/src/components/organisms/AppShell/AppShell.css
web/admin/src/components/atoms/JsonEditor/JsonEditor.tsx
web/admin/src/services/goGoHostApi.ts
web/admin/src/services/msw/handlers.ts
web/admin/src/services/msw/fixtures.ts
```

Reference OS package source:

```text
/home/manuel/code/wesen/go-go-golems/go-go-os-frontend/packages/os-core/src/theme/desktop/shell.css
/home/manuel/code/wesen/go-go-golems/go-go-os-frontend/packages/os-core/src/theme/desktop/primitives.css
/home/manuel/code/wesen/go-go-golems/go-go-os-frontend/packages/os-core/src/components/widgets/
```

## One-sentence rule

A new admin page should first look like a small, intentional OS1 operator window, and only then should it expose the product-specific data and actions.
