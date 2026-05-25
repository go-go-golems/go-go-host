import { ErrorCallout, JsonTree, LoadingBlock, Timestamp } from '../../components/atoms';
import { StatusPill } from '../../components/atoms/StatusPill';
import { apiErrorMessage } from '../../services/errors';
import { useListAdminCapabilitiesQuery } from '../../services/goGoHostApi';
import { parseJson } from '../../services/json';
import '../AdminOrgsPage/AdminOrgsPage.css';

export function AdminCapabilitiesPage() {
  const caps = useListAdminCapabilitiesQuery();
  if (caps.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={6} /></section>;
  if (caps.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load capabilities" error={apiErrorMessage(caps.error)} /></section>;
  const rows = caps.data ?? [];
  return <section className="dashboard-panel admin-inventory-page"><header><h1>Capabilities</h1><p>Read-only effective capability policy. `exec` remains unavailable for hosted v1.</p></header><table><thead><tr><th>Site</th><th>Org</th><th>Capability</th><th>Enabled</th><th>Config</th><th>Updated</th></tr></thead><tbody>{rows.map((c) => <tr key={`${c.siteId}:${c.capability}`}><td><strong>{c.siteSlug}</strong><br /><code>{c.siteId}</code></td><td>{c.orgName}<br /><small>{c.orgSlug}</small></td><td><code>{c.capability}</code></td><td><StatusPill status={c.enabled ? 'enabled' : 'disabled'} tone={c.enabled ? 'success' : 'danger'} /></td><td><JsonTree value={parseJson(c.configJson, {})} /></td><td><Timestamp value={c.updatedAt} /></td></tr>)}</tbody></table>{rows.length === 0 ? <p>No per-site capability rows found; runtime defaults still apply.</p> : null}</section>;
}
