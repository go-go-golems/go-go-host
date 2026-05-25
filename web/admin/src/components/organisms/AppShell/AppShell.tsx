import type { ReactNode } from 'react';
import { OrgSwitcher } from '../../molecules/OrgSwitcher';
import type { Membership } from '../../../services/types';
import './AppShell.css';

export interface AppShellProps {
  memberships: Membership[];
  selectedOrgId?: string;
  userLabel: string;
  devAuth?: boolean;
  sidebar?: ReactNode;
  /** Extra content rendered as a second window below the main sidebar. */
  sidebarExtra?: ReactNode;
  sidebarExtraTitle?: string;
  children: ReactNode;
  onOrgSelect?: (orgId: string) => void;
  onLogout?: () => void;
}

export function AppShell({ memberships, selectedOrgId, userLabel, devAuth, sidebar, sidebarExtra, sidebarExtraTitle = 'Contents', children, onOrgSelect, onLogout }: AppShellProps) {
  return (
    <div className="app-shell" data-part="app-shell">
      <header className="app-shell__menubar" data-part="menubar">
        <strong>go-go-host</strong>
        <OrgSwitcher memberships={memberships} selectedOrgId={selectedOrgId} onSelect={onOrgSelect} />
        <span className="app-shell__spacer" />
        {devAuth ? <span className="app-shell__badge">Dev auth ON</span> : null}
        <span>{userLabel}</span>
        {onLogout ? <button type="button" data-part="btn" onClick={onLogout}>Sign out</button> : null}
      </header>
      <div className="app-shell__body" data-part="windowing-icon-layer">
        {sidebar ? (
          <aside className="app-shell__sidebar" aria-label="Navigation">
            <div className="app-shell__sidebar-window" data-part="windowing-window" data-state="focused">
              <div className="app-shell__sidebar-title" data-part="windowing-window-title-bar" data-state="focused">
                <span aria-hidden="true" data-part="windowing-close-button" />
                <span data-part="windowing-window-title">Navigation</span>
              </div>
              <div className="app-shell__sidebar-body" data-part="windowing-window-body">{sidebar}</div>
            </div>
            {sidebarExtra ? (
              <div className="app-shell__sidebar-window" data-part="windowing-window" data-state="focused">
                <div className="app-shell__sidebar-title" data-part="windowing-window-title-bar" data-state="focused">
                  <span aria-hidden="true" data-part="windowing-close-button" />
                  <span data-part="windowing-window-title">{sidebarExtraTitle}</span>
                </div>
                <div className="app-shell__sidebar-body" data-part="windowing-window-body">{sidebarExtra}</div>
              </div>
            ) : null}
          </aside>
        ) : null}
        <main className="app-shell__content">{children}</main>
      </div>
    </div>
  );
}
