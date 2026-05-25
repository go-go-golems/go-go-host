import './AdminSidebar.css';

export type AdminSection = 'overview' | 'runtimes' | 'orgs' | 'users' | 'sites' | 'deployments' | 'agents' | 'audit' | 'quotas' | 'capabilities' | 'domains';

const sections: { id: AdminSection; label: string }[] = [
  { id: 'overview', label: 'Overview' },
  { id: 'runtimes', label: 'Runtimes' },
  { id: 'orgs', label: 'Orgs' },
  { id: 'users', label: 'Users' },
  { id: 'sites', label: 'Sites' },
  { id: 'deployments', label: 'Deployments' },
  { id: 'agents', label: 'Agents' },
  { id: 'audit', label: 'Audit' },
  { id: 'quotas', label: 'Quotas' },
  { id: 'capabilities', label: 'Capabilities' },
  { id: 'domains', label: 'Domains' },
];

export function AdminSidebar({ active = 'overview', onSelect }: { active?: AdminSection; onSelect?: (section: AdminSection) => void }) {
  return <nav className="admin-sidebar" data-part="admin-sidebar" aria-label="Platform admin sections">{sections.map((s) => <button key={s.id} type="button" data-active={s.id === active} onClick={() => onSelect?.(s.id)}>{s.label}</button>)}</nav>;
}
