import { ErrorCallout, LoadingBlock, Timestamp } from '../../components/atoms';
import { apiErrorMessage } from '../../services/errors';
import { useListAdminQuotasQuery } from '../../services/goGoHostApi';
import '../AdminOrgsPage/AdminOrgsPage.css';

function mb(bytes: number) { return `${(bytes / 1024 / 1024).toFixed(0)} MiB`; }
export function AdminQuotasPage() {
  const quotas = useListAdminQuotasQuery();
  if (quotas.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={6} /></section>;
  if (quotas.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load quotas" error={apiErrorMessage(quotas.error)} /></section>;
  const rows = quotas.data ?? [];
  return <section className="dashboard-panel admin-inventory-page"><header><h1>Quotas</h1><p>Read-only current site quota policy and runtime usage counters. Editable defaults and overrides are planned next.</p></header><table><thead><tr><th>Site</th><th>Org</th><th>Bundle</th><th>DB soft</th><th>DB hard</th><th>Timeout</th><th>Usage</th><th>Updated</th></tr></thead><tbody>{rows.map((q) => <tr key={q.siteId}><td><strong>{q.siteSlug}</strong><br /><small>{q.primaryHost}</small></td><td>{q.orgName}<br /><small>{q.orgSlug}</small></td><td>{mb(q.bundleMaxBytes)}</td><td>{mb(q.dbSoftMaxBytes)}</td><td>{mb(q.dbHardMaxBytes)}</td><td>{q.requestTimeoutMs} ms</td><td>{q.requestsTotal} req / {q.errorsTotal} err</td><td><Timestamp value={q.updatedAt} /></td></tr>)}</tbody></table>{rows.length === 0 ? <p>No quota rows found.</p> : null}</section>;
}
