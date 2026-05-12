import { useNavigate, useOutletContext, useParams } from 'react-router-dom';
import { CodeBlock, EmptyState, ErrorCallout, LoadingBlock } from '../../components/atoms';
import { DeploymentTimeline, RuntimeStatusPanel, ValidationReportPanel } from '../../components/organisms';
import { apiErrorMessage } from '../../services/errors';
import { useGetRuntimeQuery, useListDeploymentsQuery } from '../../services/goGoHostApi';
import type { Site } from '../../services/types';
import './SiteOverviewPage.css';

export function SiteOverviewPage() {
  const { orgId, siteId } = useParams();
  const { site } = useOutletContext<{ site: Site }>();
  const navigate = useNavigate();
  const runtime = useGetRuntimeQuery(siteId ?? '', { skip: !siteId });
  const deployments = useListDeploymentsQuery(siteId ?? '', { skip: !siteId });
  const activeDeployment = deployments.data?.find((deployment) => deployment.id === site.activeDeploymentId) ?? deployments.data?.[0];
  if (runtime.isLoading || deployments.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={6} /></section>;
  return <div className="site-overview-page">
    {runtime.error ? <section className="dashboard-panel"><ErrorCallout title="Unable to load runtime" error={apiErrorMessage(runtime.error)} /></section> : runtime.data ? <RuntimeStatusPanel runtime={runtime.data} /> : null}
    <section className="dashboard-panel site-overview-page__deployments"><header><div><h2>Deployments</h2><p>Latest upload and activation history for this site.</p></div><button type="button" data-part="btn" onClick={() => navigate(`/app/orgs/${orgId}/sites/${site.id}/deployments`)}>Open deployments</button></header>{deployments.error ? <ErrorCallout title="Unable to load deployments" error={apiErrorMessage(deployments.error)} /> : <DeploymentTimeline deployments={deployments.data ?? []} onSelect={(deploymentId) => navigate(`/app/orgs/${orgId}/sites/${site.id}/deployments/${deploymentId}`)} />}</section>
    {activeDeployment ? <ValidationReportPanel deployment={activeDeployment} /> : <section className="dashboard-panel"><EmptyState title="No active deployment" body="Upload and activate a bundle to serve traffic." /></section>}
    <section className="dashboard-panel site-overview-page__debug"><h2>Site DTO</h2><CodeBlock code={JSON.stringify(site, null, 2)} /></section>
  </div>;
}
