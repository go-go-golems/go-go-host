import { ErrorCallout, LoadingBlock, Timestamp } from '../../components/atoms';
import { StatusPill } from '../../components/atoms/StatusPill';
import { apiErrorMessage } from '../../services/errors';
import { useListAdminDomainsQuery } from '../../services/goGoHostApi';
import '../AdminOrgsPage/AdminOrgsPage.css';

export function AdminDomainsPage() {
  const domains = useListAdminDomainsQuery();
  if (domains.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={6} /></section>;
  if (domains.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load domains" error={apiErrorMessage(domains.error)} /></section>;
  const rows = domains.data ?? [];
  return <section className="dashboard-panel admin-inventory-page"><header><h1>Domains</h1><p>Base/custom domain inventory and verification placeholders. TLS automation is deferred.</p></header><table><thead><tr><th>Hostname</th><th>Site</th><th>Org</th><th>Status</th><th>Verification token</th><th>Verified</th><th>Created</th></tr></thead><tbody>{rows.map((d) => <tr key={d.id}><td><strong>{d.hostname}</strong></td><td>{d.siteSlug}<br /><code>{d.siteId}</code></td><td>{d.orgName}<br /><small>{d.orgSlug}</small></td><td><StatusPill status={d.status} tone={d.status === 'verified' ? 'success' : 'warning'} /></td><td><code>{d.verificationToken || '—'}</code></td><td><Timestamp value={d.verifiedAt} /></td><td><Timestamp value={d.createdAt} /></td></tr>)}</tbody></table>{rows.length === 0 ? <p>No custom domain rows found. Primary hosts are visible on the Sites page.</p> : null}</section>;
}
