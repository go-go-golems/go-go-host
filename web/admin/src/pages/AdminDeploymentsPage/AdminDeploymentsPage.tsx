import { useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { ErrorCallout, LoadingBlock, Timestamp } from '../../components/atoms';
import { DeploymentStatusPill } from '../../components/molecules/DeploymentStatusPill';
import { apiErrorMessage } from '../../services/errors';
import { useListAdminDeploymentsQuery } from '../../services/goGoHostApi';
import type { DeploymentStatus } from '../../services/types';
import '../AdminOrgsPage/AdminOrgsPage.css';
import './AdminDeploymentsPage.css';

const statusOptions: Array<DeploymentStatus | ''> = ['', 'uploaded', 'validated', 'rejected', 'active', 'superseded'];

export function AdminDeploymentsPage() {
  const [status, setStatus] = useState<DeploymentStatus | ''>('');
  const params = useMemo(() => ({ status: status || undefined, limit: 100 }), [status]);
  const deployments = useListAdminDeploymentsQuery(params);
  if (deployments.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={6} /></section>;
  if (deployments.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load deployment inventory" error={apiErrorMessage(deployments.error)} /></section>;
  const rows = deployments.data ?? [];
  return <section className="dashboard-panel admin-inventory-page admin-deployments-page"><header><div><h1>Deployments</h1><p>Global deployment inventory with status and actor filters.</p></div><label>Status <select value={status} onChange={(e) => setStatus(e.target.value as DeploymentStatus | '')}>{statusOptions.map((option) => <option key={option || 'all'} value={option}>{option || 'all'}</option>)}</select></label></header><table><thead><tr><th>Deployment</th><th>Site</th><th>Org</th><th>Status</th><th>Actor</th><th>Created</th><th>Bundle</th></tr></thead><tbody>{rows.map((deployment) => <tr key={deployment.id}><td><Link to={`/admin/deployments/${deployment.id}`}>{deployment.id}</Link><br /><small>v{deployment.version}</small></td><td>{deployment.siteSlug}<br /><small>{deployment.primaryHost}</small></td><td>{deployment.orgName}<br /><small>{deployment.orgSlug}</small></td><td><DeploymentStatusPill status={deployment.status} /></td><td>{deployment.createdByType}<br /><code>{deployment.createdById}</code></td><td><Timestamp value={deployment.createdAt} /></td><td><code>{deployment.bundleRef || '—'}</code></td></tr>)}</tbody></table>{rows.length === 0 ? <p>No deployments found.</p> : null}</section>;
}
