import { RoleBadge } from '../../atoms/RoleBadge';
import type { Membership } from '../../../services/types';
import './OrgSwitcher.css';
export interface OrgSwitcherProps { memberships: Membership[]; selectedOrgId?: string; onSelect?: (orgId: string) => void; }
export function OrgSwitcher({ memberships, selectedOrgId, onSelect }: OrgSwitcherProps) {
  if (memberships.length === 0) return <span className="org-switcher org-switcher--empty">No organizations</span>;
  return <label className="org-switcher">Org <select value={selectedOrgId ?? memberships[0]?.orgId} onChange={(e) => onSelect?.(e.target.value)}>{memberships.map((m) => <option key={m.orgId} value={m.orgId}>{m.orgName}</option>)}</select><RoleBadge role={memberships.find((m) => m.orgId === (selectedOrgId ?? memberships[0]?.orgId))?.role ?? memberships[0].role} /></label>;
}
