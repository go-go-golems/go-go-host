import { ErrorCallout, LoadingBlock, Timestamp } from '../../components/atoms';
import { apiErrorMessage } from '../../services/errors';
import { useListAdminOrgsQuery } from '../../services/goGoHostApi';
import './AdminOrgsPage.css';

export function AdminOrgsPage() {
  const orgs = useListAdminOrgsQuery();
  if (orgs.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={6} /></section>;
  if (orgs.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load org inventory" error={apiErrorMessage(orgs.error)} /></section>;
  const rows = orgs.data ?? [];
  return <section className="dashboard-panel admin-inventory-page"><header><h1>Organizations</h1><p>Cross-tenant organization inventory with membership, site, and deployment counts.</p></header><table><thead><tr><th>Org</th><th>ID</th><th>Members</th><th>Sites</th><th>Deployments</th><th>Created</th></tr></thead><tbody>{rows.map((org) => <tr key={org.id}><td><strong>{org.name}</strong><br /><small>{org.slug}</small></td><td><code>{org.id}</code></td><td>{org.memberCount}</td><td>{org.siteCount}</td><td>{org.deploymentCount}</td><td><Timestamp value={org.createdAt} /></td></tr>)}</tbody></table>{rows.length === 0 ? <p>No organizations found.</p> : null}</section>;
}
