import type { Membership } from '../../../services/types';
import { RoleBadge } from '../../atoms';
import './MembersTable.css';
export function MembersTable({ memberships, selectedOrgId }: { memberships: Membership[]; selectedOrgId?: string }) {
  return <table className="members-table"><thead><tr><th>Organization</th><th>Role</th><th>Org ID</th><th>Status</th></tr></thead><tbody>{memberships.map((m) => <tr key={m.orgId}><td><strong>{m.orgName}</strong><br /><small>{m.orgSlug}</small></td><td><RoleBadge role={m.role} /></td><td>{m.orgId}</td><td>{m.orgId === selectedOrgId ? 'current' : 'available'}</td></tr>)}</tbody></table>;
}
