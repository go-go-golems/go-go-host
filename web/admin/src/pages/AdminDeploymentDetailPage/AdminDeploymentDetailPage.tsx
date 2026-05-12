import { Link, useParams } from 'react-router-dom';
import { ErrorCallout, JsonTree, LoadingBlock, Timestamp } from '../../components/atoms';
import { DeploymentStatusPill, ManifestSummary, ValidationSummary } from '../../components/molecules';
import { apiErrorMessage } from '../../services/errors';
import { useGetAdminDeploymentQuery } from '../../services/goGoHostApi';
import { parseManifest, parseValidationReport } from '../../services/json';
import './AdminDeploymentDetailPage.css';

export function AdminDeploymentDetailPage() {
  const { deploymentId } = useParams();
  const deployment = useGetAdminDeploymentQuery(deploymentId ?? '', { skip: !deploymentId });
  if (deployment.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={6} /></section>;
  if (deployment.error || !deployment.data) return <section className="dashboard-panel"><ErrorCallout title="Unable to load deployment" error={apiErrorMessage(deployment.error) || 'Deployment not found'} /></section>;
  const dep = deployment.data;
  const manifest = parseManifest(dep.manifestJson);
  const validation = parseValidationReport(dep.validationJson);
  return <div className="admin-deployment-detail-page">
    <section className="dashboard-panel admin-deployment-detail-page__header">
      <div><Link to="/admin/deployments">← Deployments</Link><h1>{dep.id}</h1><p>{dep.orgName} / {dep.siteSlug} · {dep.primaryHost}</p></div>
      <DeploymentStatusPill status={dep.status} />
    </section>
    <section className="dashboard-panel admin-deployment-detail-page__meta"><h2>Metadata</h2><dl><dt>Organization</dt><dd>{dep.orgName} <code>{dep.orgId}</code></dd><dt>Site</dt><dd>{dep.siteSlug} <code>{dep.siteId}</code></dd><dt>Version</dt><dd>{dep.version}</dd><dt>Created by</dt><dd>{dep.createdByType} <code>{dep.createdById}</code></dd><dt>Created</dt><dd><Timestamp value={dep.createdAt} /></dd><dt>Activated</dt><dd><Timestamp value={dep.activatedAt} /></dd><dt>Bundle</dt><dd><code>{dep.bundleRef || '—'}</code></dd><dt>Unpacked path</dt><dd><code>{dep.unpackedPath || '—'}</code></dd></dl></section>
    <section className="dashboard-panel"><h2>Manifest summary</h2><ManifestSummary manifest={manifest} /></section>
    <section className="dashboard-panel"><h2>Validation</h2><ValidationSummary report={validation} /></section>
    <section className="dashboard-panel"><h2>Manifest JSON</h2><JsonTree value={manifest} /></section>
    <section className="dashboard-panel"><h2>Validation JSON</h2><JsonTree value={validation} /></section>
  </div>;
}
