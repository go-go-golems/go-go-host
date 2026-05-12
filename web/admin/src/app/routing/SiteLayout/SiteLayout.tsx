import { Outlet, useNavigate, useParams } from 'react-router-dom';
import { EmptyState, ErrorCallout, LoadingBlock } from '../../../components/atoms';
import { SiteHeader, SiteTabs } from '../../../components/organisms';
import { apiErrorMessage } from '../../../services/errors';
import { useListSitesQuery } from '../../../services/goGoHostApi';
import './SiteLayout.css';

export function SiteLayout() {
  const { orgId, siteId } = useParams();
  const navigate = useNavigate();
  const sites = useListSitesQuery(orgId ?? '', { skip: !orgId });
  const site = sites.data?.find((candidate) => candidate.id === siteId);
  if (sites.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={5} /></section>;
  if (sites.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load site" error={apiErrorMessage(sites.error)} /></section>;
  if (!site) return <section className="dashboard-panel"><EmptyState title="Site not found" body="The selected site was not returned by the organization site list." action={<button type="button" data-part="btn" onClick={() => navigate(`/app/orgs/${orgId}/sites`)}>Back to sites</button>} /></section>;
  const basePath = `/app/orgs/${orgId}/sites/${site.id}`;
  return <div className="site-layout"><SiteHeader site={site} onBack={() => navigate(`/app/orgs/${orgId}/sites`)} /><SiteTabs basePath={basePath} /><Outlet context={{ site }} /></div>;
}
