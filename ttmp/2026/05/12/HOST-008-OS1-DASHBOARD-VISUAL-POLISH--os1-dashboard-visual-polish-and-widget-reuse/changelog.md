# Changelog

## 2026-05-12

- Initial workspace created


## 2026-05-12

Created OS1 dashboard visual polish ticket, captured public/reference screenshots, documented current-state gaps, and applied the first shell/sidebar/panel/table CSS cleanup slice using existing @go-go-golems/os-core theme parts.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/app/macos1-bridge.css — Dashboard panels now render with Mac-style borders
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/organisms/AppShell/AppShell.css — Centered desktop workspace
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/organisms/AppShell/AppShell.tsx — Sidebar now uses OS-core window data-part structure
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/organisms/SitesTable/SitesTable.css — Sites table spacing and borders normalized for OS1 style


## 2026-05-12

Polished the site settings page with OS-core checkbox widgets, subtle semantic text highlights, unified OS1 font sizing, and a CodeMirror JSON editor with syntax highlighting.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/package.json — Adds CodeMirror dependencies
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/atoms/JsonEditor/JsonEditor.tsx — Reusable CodeMirror JSON editor atom with syntax highlighting
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/SiteSettingsPage/SiteSettingsPage.css — Settings-specific OS1 font scale
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/SiteSettingsPage/SiteSettingsPage.tsx — Settings page now uses JsonEditor


## 2026-05-12

Added a reusable OS1 admin dashboard UI work guidelines playbook for future page design and implementation.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-008-OS1-DASHBOARD-VISUAL-POLISH--os1-dashboard-visual-polish-and-widget-reuse/playbook/01-os1-admin-dashboard-ui-work-guidelines.md — Reusable design and implementation guidelines for future admin dashboard UI work


## 2026-05-12

Uploaded the OS1 admin UI guidelines playbook to reMarkable.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-008-OS1-DASHBOARD-VISUAL-POLISH--os1-dashboard-visual-polish-and-widget-reuse/playbook/01-os1-admin-dashboard-ui-work-guidelines.md — Source document uploaded as HOST-008_OS1_Admin_UI_Guidelines.pdf


## 2026-05-12

Applied the OS1 UI guidelines to the Agents page: checkbox widget for auto-activation, semantic highlights, compact OS1 create form, and normalized agent/signing-key tables.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/organisms/AgentKeysTable/AgentKeysTable.css — Signing-key table normalized with subtle status highlights
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/organisms/AgentsTable/AgentsTable.css — Agent table normalized to OS1 dense table styling
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/AgentsPage/AgentsPage.css — Agents page OS1 form
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/pages/AgentsPage/AgentsPage.tsx — Agents page now uses OS-core Checkbox and semantic highlighted explanatory copy


## 2026-05-12

Applied cross-page OS1 normalization and updated color guidance so accents are used for badges/buttons/panels rather than inline body text.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/ttmp/2026/05/12/HOST-008-OS1-DASHBOARD-VISUAL-POLISH--os1-dashboard-visual-polish-and-widget-reuse/playbook/01-os1-admin-dashboard-ui-work-guidelines.md — Updated color guidance to avoid colored inline prose highlights
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/app/macos1-bridge.css — Shared cross-page OS1 form
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/atoms/CodeBlock/CodeBlock.css — Read-only code blocks aligned with OS1 surface styling
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/atoms/StatusPill/StatusPill.css — Status accents moved into badge backgrounds
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/organisms/AuditTimeline/AuditTimeline.css — Timeline rows normalized to OS1 bordered rows
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/organisms/DeploymentTimeline/DeploymentTimeline.css — Deployment history rows normalized to OS1 bordered rows


## 2026-05-12

Changed the dashboard page background from light blue to white and recorded remaining coverage caveats.

### Related Files

- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/app/macos1-bridge.css — Overrides the macOS1 desktop background token to white for the admin dashboard
- /home/manuel/workspaces/2026-05-11/go-go-host-v1/go-go-host/web/admin/src/components/organisms/AppShell/AppShell.css — Forces the shell and work area background to white instead of light-blue desktop fill

