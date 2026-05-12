import './RoleBadge.css';
export type Role = 'org_owner' | 'org_developer' | 'org_viewer';
export function RoleBadge({ role }: { role: Role }) {
  return <span className="role-badge" data-role={role}>{role.replace('org_', '')}</span>;
}
