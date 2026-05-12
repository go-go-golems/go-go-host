---
Title: OS1 Dashboard Visual Redesign Investigation
Ticket: HOST-008-OS1-DASHBOARD-VISUAL-POLISH
Status: active
Topics:
    - frontend
    - design-system
    - storybook
    - debugging
DocType: design-doc
Intent: ""
Owners: []
RelatedFiles:
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/go-go-os-frontend/packages/os-core/src/theme/desktop/primitives.css
      Note: Reference OS-core button/table/form primitive styling used for consistency
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/go-go-os-frontend/packages/os-core/src/theme/desktop/shell.css
      Note: Reference OS-core window frame/title/body data-part styles reused by the dashboard
    - Path: web/admin/src/app/macos1-bridge.css
      Note: Bridge layer mapping dashboard panels and legacy atoms onto macOS1 theme tokens
    - Path: web/admin/src/components/organisms/AppShell/AppShell.css
      Note: Global centered desktop workspace and menu/sidebar layout styling
    - Path: web/admin/src/components/organisms/AppShell/AppShell.tsx
      Note: Dashboard shell structure now wraps navigation as an OS1 window palette
    - Path: web/admin/src/components/organisms/SitesTable/SitesTable.css
      Note: Representative table styling to align data grids with the OS1 visual system
ExternalSources: []
Summary: ""
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: ""
WhenToUse: ""
---


# OS1 Dashboard Visual Redesign Investigation

## Executive summary

The `go-go-host` dashboard already imports the OS1/macOS theme package, but it does not yet compose the page like an OS1 application. It uses the colors and font, while still laying out the app as a simple web admin table: a top bar, a sidebar, and a large flat content panel. That is why it feels wrong. The reference examples at `go-go-os-examples.yolo.scapegoat.dev` look better because they use a constrained desktop workspace, window frames, striped title bars, deliberate spacing, and consistently dense controls.

This ticket starts the conversion from "admin page with retro tokens" to "unified OS1 dashboard." The first implementation slice keeps the existing React routing/data model, reuses `@go-go-golems/os-core` CSS/data-part conventions, and improves the global shell, sidebar, dashboard panel, and Sites table. It does not yet migrate the whole app to `DesktopShell`; that can be a later step once visual consistency is achieved safely.

## Problem statement

The current public dashboard works functionally, but visually it has three problems.

First, it lacks a strong OS composition. The examples site feels like a desktop application because every major area has a role: navigation column, menu strip, content window, panel, table. The current dashboard has pieces of that language but they are not arranged as a coherent system.

Second, spacing is inconsistent. The public Sites page has a narrow header, a left nav, a table, and then a giant blue-gray void. Some elements are crowded, while the whole page wastes space. This makes the beta feel unfinished even though the backend and deployment flow are working.

Third, styling is split between OS tokens and ad-hoc component CSS. Some components use `--hc-*` variables, others still use `--os-*`, and several preserve modern rounded-pill defaults. The result is neither modern SaaS nor classic OS1.

## Evidence

Screenshots saved under this ticket:

```text
sources/screenshots/host-008-current-app.png
sources/screenshots/host-008-current-auth-app-viewport.png
sources/screenshots/host-008-go-go-os-examples-home.png
sources/screenshots/host-008-go-go-os-examples-rich-widgets.png
sources/screenshots/host-008-appshell-after-local.png
```

Key files inspected:

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

Reference package inspected:

```text
/home/manuel/code/wesen/go-go-golems/go-go-os-frontend/packages/os-core
```

Installed package used by dashboard:

```text
web/admin/node_modules/@go-go-golems/os-core
```

The dashboard entrypoint already imports:

```ts
import '@go-go-golems/os-core/theme';
import '@go-go-golems/os-core/desktop-theme-macos1';
```

and the provider already wraps the app in:

```tsx
<div data-widget="hypercard" className="theme-macos1">
```

That means the visual problem is not package availability. It is composition and consistency.

## Design direction

The desired visual model is:

```text
sticky OS1 menu bar
  -> centered desktop work area
      -> left navigation palette/window
      -> main content window
          -> page header
          -> toolbar/actions
          -> OS1 table/panels/forms
```

This keeps the dashboard usable as a normal web app while borrowing enough structure from the desktop shell to feel intentional.

### Why not migrate immediately to `DesktopShell`?

`@go-go-golems/os-core/desktop-react` exports a full `DesktopShell`, `DesktopMenuBar`, `WindowSurface`, and related primitives. Those are powerful, but moving the whole dashboard into draggable windows would change application behavior, routing assumptions, focus handling, and potentially accessibility semantics.

For this ticket, a staged approach is safer:

1. **Stage A — Visual shell alignment**: use OS theme imports, `data-part` names, window CSS, spacing, and table/control primitives while preserving app routes.
2. **Stage B — Component extraction**: introduce dashboard-local primitives like `DashboardWindow`, `DashboardToolbar`, `DashboardSection`, and `DashboardTable`.
3. **Stage C — Optional desktop shell**: only if product direction wants a true draggable multi-window dashboard, migrate selected views to `DesktopShell`.

The first slice implements Stage A.

## First implementation slice

Changed files:

```text
web/admin/src/components/organisms/AppShell/AppShell.tsx
web/admin/src/components/organisms/AppShell/AppShell.css
web/admin/src/app/macos1-bridge.css
web/admin/src/components/organisms/OrgSidebar/OrgSidebar.css
web/admin/src/components/organisms/SitesTable/SitesTable.css
web/admin/src/components/molecules/SiteHostCopy/SiteHostCopy.css
```

### App shell

The app shell now treats the page as a centered desktop workspace instead of an unbounded blue canvas. The body uses a fixed max width, consistent gap, and top/bottom padding. The menu bar is sticky and has the same strong black border language as the OS examples.

### Sidebar

The sidebar is now a static OS1 window/palette. It uses the existing OS-core data-part contract:

```tsx
<aside data-part="windowing-window" data-state="focused">
  <div data-part="windowing-window-title-bar" data-state="focused">
    <span data-part="windowing-close-button" />
    <span data-part="windowing-window-title">Navigation</span>
  </div>
  <div data-part="windowing-window-body">...</div>
</aside>
```

This reuses the installed theme's window border, title bar, close box, and body semantics without making the sidebar draggable.

### Dashboard panels

`.dashboard-panel` now renders as a Mac-style content window, with:

- 2px black border,
- window shadow,
- striped title bar,
- consistent internal padding,
- header separator,
- normalized heading/p/body margins.

This is deliberately a bridge. A later cleanup should replace the broad `.dashboard-panel::before` styling with a real `DashboardWindow` component that can set page-specific titles.

### Tables and copy controls

The Sites table now has stronger OS1 table styling:

- 2px outer border,
- dense but readable cell padding,
- uppercase black header row,
- hover background,
- wrapped deployment IDs,
- host/copy controls that wrap instead of crowding.

## Validation

Build passed:

```bash
cd web/admin
pnpm build
```

Local Storybook was started and a post-change screenshot was captured:

```text
sources/screenshots/host-008-appshell-after-local.png
```

## Remaining work

This is not finished. It is the first shell/panel pass. The remaining pages still need review and tightening:

- Agents page and agent key enrollment forms,
- Audit timeline,
- Members table,
- Usage page,
- Site detail layout and tabs,
- Deployment upload panel,
- Deployment detail/validation report,
- Admin overview and inventory tables,
- No-orgs/bootstrap/auth callback states,
- loading/error/empty state consistency,
- Storybook stories that show full-page shells instead of isolated atoms only.

## Recommended next implementation plan

1. Add `DashboardWindow` component with explicit `title`, `toolbar`, and `children` slots.
2. Replace raw `<section className="dashboard-panel">` uses with `DashboardWindow` page by page.
3. Add `DashboardToolbar` and `DashboardSection` primitives for page headers, filters, and action rows.
4. Normalize all table components to a single OS1 table class or reusable table wrapper.
5. Normalize all forms to `data-part="field-input"` / `field-select` / `field-grid` patterns.
6. Run Storybook and capture before/after screenshots for each major page.
7. Build and deploy a new beta image, then capture public screenshots from `https://hosting.yolo.scapegoat.dev/app`.

## Open questions

1. Should the beta dashboard remain a classic web layout with OS1 chrome, or eventually become a true multi-window desktop using `DesktopShell`?
2. Should `@go-go-golems/os-widgets` be added as a dependency for richer components, or should the dashboard stay on `@go-go-golems/os-core` only?
3. Should every page title appear in the window title bar, the page body, or both?
4. Should the user dashboard and platform-admin dashboard share the same shell primitives or have subtly distinct app identities?
