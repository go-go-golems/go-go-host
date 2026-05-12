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
  children: ReactNode;
  onOrgSelect?: (orgId: string) => void;
}

export function AppShell({ memberships, selectedOrgId, userLabel, devAuth, sidebar, children, onOrgSelect }: AppShellProps) {
  return (
    <div className="app-shell" data-part="app-shell">
      <header className="app-shell__menubar" data-part="menubar">
        <strong>go-go-host</strong>
        <OrgSwitcher memberships={memberships} selectedOrgId={selectedOrgId} onSelect={onOrgSelect} />
        <span className="app-shell__spacer" />
        {devAuth ? <span className="app-shell__badge">Dev auth ON</span> : null}
        <span>{userLabel}</span>
      </header>
      <div className="app-shell__body">
        {sidebar ? <aside className="app-shell__sidebar">{sidebar}</aside> : null}
        <main className="app-shell__content">{children}</main>
      </div>
    </div>
  );
}
