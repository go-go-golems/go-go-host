import './OrgSidebar.css';
export type OrgSection = 'sites' | 'agents' | 'audit' | 'members' | 'usage' | 'docs';
const sections: { id: OrgSection; label: string }[] = [
  { id: 'sites', label: 'Sites' }, { id: 'agents', label: 'Agents' }, { id: 'docs', label: 'Docs' }, { id: 'audit', label: 'Audit' }, { id: 'members', label: 'Members' }, { id: 'usage', label: 'Usage' },
];
export function OrgSidebar({ active = 'sites', onSelect }: { active?: OrgSection; onSelect?: (section: OrgSection) => void }) {
  return <nav className="org-sidebar" data-part="org-sidebar">{sections.map((s) => <button key={s.id} type="button" data-active={s.id === active} onClick={() => onSelect?.(s.id)}>{s.label}</button>)}</nav>;
}
