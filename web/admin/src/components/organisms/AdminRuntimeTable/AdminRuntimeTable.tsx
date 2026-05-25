import { Link } from 'react-router-dom';
import { RuntimeBadge } from '../../molecules/RuntimeBadge';
import { EmptyState, Timestamp } from '../../atoms';
import type { RuntimeStatus } from '../../../services/types';
import './AdminRuntimeTable.css';

export interface AdminRuntimeTableProps {
  runtimes: RuntimeStatus[];
  onRestart?: (runtime: RuntimeStatus) => void;
  onStop?: (runtime: RuntimeStatus) => void;
  actionBusySiteId?: string;
}

export function AdminRuntimeTable({ runtimes, onRestart, onStop, actionBusySiteId }: AdminRuntimeTableProps) {
  if (runtimes.length === 0) return <EmptyState title="No runtimes yet" body="No site runtimes have reported status in this daemon process." />;
  return (
    <table className="admin-runtime-table">
      <thead><tr><th>Site</th><th>Org</th><th>Status</th><th>Deployment</th><th>Hosts</th><th>Requests</th><th>Errors</th><th>Started</th><th>Last error</th><th>Actions</th></tr></thead>
      <tbody>{runtimes.map((runtime) => (
        <tr key={runtime.siteId} data-state={runtime.status}>
          <td><code>{runtime.siteId}</code></td>
          <td><code>{runtime.orgId || '—'}</code></td>
          <td><RuntimeBadge runtime={runtime} compact /></td>
          <td>{runtime.deploymentId ? <Link to={`/admin/deployments/${runtime.deploymentId}`}>{runtime.deploymentId}</Link> : '—'}</td>
          <td>{runtime.hosts?.length ? runtime.hosts.join(', ') : '—'}</td>
          <td>{runtime.requestsTotal ?? 0}</td>
          <td>{runtime.errorsTotal ?? 0}</td>
          <td><Timestamp value={runtime.startedAt} /></td>
          <td className="admin-runtime-table__error">{runtime.lastError || '—'}</td>
          <td className="admin-runtime-table__actions"><button type="button" data-part="btn" disabled={!onRestart || actionBusySiteId === runtime.siteId} onClick={() => onRestart?.(runtime)}>Restart</button><button type="button" data-part="btn" disabled={!onStop || runtime.status === 'stopped' || actionBusySiteId === runtime.siteId} onClick={() => onStop?.(runtime)}>Stop</button></td>
        </tr>
      ))}</tbody>
    </table>
  );
}
