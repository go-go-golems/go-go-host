import { useMemo, useState } from 'react';
import { ErrorCallout, LoadingBlock, Timestamp } from '../../components/atoms';
import { AgentStatusBadge } from '../../components/molecules/AgentStatusBadge';
import { apiErrorMessage } from '../../services/errors';
import { useListAdminAgentsQuery } from '../../services/goGoHostApi';
import '../AdminOrgsPage/AdminOrgsPage.css';
import './AdminAgentsPage.css';

export function AdminAgentsPage() {
  const [status, setStatus] = useState('');
  const params = useMemo(() => ({ status: status || undefined }), [status]);
  const agents = useListAdminAgentsQuery(params);
  if (agents.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={6} /></section>;
  if (agents.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load global agents" error={apiErrorMessage(agents.error)} /></section>;
  const rows = agents.data ?? [];
  return <section className="dashboard-panel admin-inventory-page admin-agents-page"><header><div><h1>Agents</h1><p>Global automation-agent inventory across organizations.</p></div><label>Status <select value={status} onChange={(event) => setStatus(event.target.value)}><option value="">all</option><option value="active">active</option><option value="revoked">revoked</option></select></label></header><table><thead><tr><th>Agent</th><th>Org</th><th>Status</th><th>Grants</th><th>Created by</th><th>Created</th><th>Last seen</th></tr></thead><tbody>{rows.map((agent) => <tr key={agent.id}><td><strong>{agent.name}</strong><br /><code>{agent.id}</code></td><td>{agent.orgName}<br /><small>{agent.orgSlug}</small></td><td><AgentStatusBadge status={agent.status} /></td><td>{agent.grantCount}</td><td><code>{agent.createdByUserId || '—'}</code></td><td><Timestamp value={agent.createdAt} /></td><td><Timestamp value={agent.lastSeenAt} /></td></tr>)}</tbody></table>{rows.length === 0 ? <p>No agents found.</p> : null}</section>;
}
