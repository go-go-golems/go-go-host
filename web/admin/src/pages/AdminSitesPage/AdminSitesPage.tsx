import { ErrorCallout, LoadingBlock } from '../../components/atoms';
import { RuntimeBadge } from '../../components/molecules';
import { apiErrorMessage } from '../../services/errors';
import { useListAdminSitesQuery } from '../../services/goGoHostApi';
import '../AdminOrgsPage/AdminOrgsPage.css';

export function AdminSitesPage() {
  const sites = useListAdminSitesQuery();
  if (sites.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={6} /></section>;
  if (sites.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load site inventory" error={apiErrorMessage(sites.error)} /></section>;
  const rows = sites.data ?? [];
  return <section className="dashboard-panel admin-inventory-page"><header><h1>Sites</h1><p>Global hosted-site inventory with org, runtime, host, and counter state.</p></header><table><thead><tr><th>Site</th><th>Org</th><th>Host</th><th>Status</th><th>Runtime</th><th>Active deployment</th><th>Last error</th></tr></thead><tbody>{rows.map((site) => <tr key={site.id}><td><strong>{site.name}</strong><br /><small>{site.slug}</small><br /><code>{site.id}</code></td><td>{site.orgName}<br /><small>{site.orgSlug}</small></td><td>{site.primaryHost}</td><td>{site.status}</td><td><RuntimeBadge compact runtime={{ siteId: site.id, orgId: site.orgId, status: site.runtimeStatus, requestsTotal: site.requestsTotal, errorsTotal: site.errorsTotal }} /></td><td>{site.activeDeploymentId || '—'}</td><td>{site.lastError || '—'}</td></tr>)}</tbody></table>{rows.length === 0 ? <p>No sites found.</p> : null}</section>;
}
