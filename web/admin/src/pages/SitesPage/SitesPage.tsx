import { useNavigate, useParams } from 'react-router-dom';
import { ErrorCallout, LoadingBlock } from '../../components/atoms';
import { SitesTable } from '../../components/organisms';
import { useListSitesQuery } from '../../services/goGoHostApi';
import { fixtures } from '../../services/msw/fixtures';
import './SitesPage.css';

export function SitesPage() {
  const { orgId } = useParams();
  const navigate = useNavigate();
  const sites = useListSitesQuery(orgId ?? '', { skip: !orgId });
  if (sites.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={5} /></section>;
  if (sites.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load sites" error="The site list request failed." /></section>;
  return <section className="dashboard-panel sites-page"><header><div><h1>Sites</h1><p>Manage hosted Goja sites for this organization.</p></div><button type="button" data-part="btn" onClick={() => navigate(`/app/orgs/${orgId}/sites/new`)}>New site</button></header><SitesTable sites={sites.data ?? []} runtimes={{ site_123: fixtures.runtimeReady, site_456: fixtures.runtimeStopped }} onOpenSite={(siteId) => navigate(`/app/orgs/${orgId}/sites/${siteId}`)} /></section>;
}
