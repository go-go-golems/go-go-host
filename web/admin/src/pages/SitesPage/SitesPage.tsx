import { useCallback, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { ErrorCallout, LoadingBlock } from '../../components/atoms';
import { SitesTable } from '../../components/organisms';
import { apiErrorMessage } from '../../services/errors';
import { useListSitesQuery } from '../../services/goGoHostApi';
import type { RuntimeStatus } from '../../services/types';
import { SiteRuntimeFanout } from './SiteRuntimeBadges';
import './SitesPage.css';

export function SitesPage() {
  const { orgId } = useParams();
  const navigate = useNavigate();
  const sites = useListSitesQuery(orgId ?? '', { skip: !orgId });
  const [runtimes, setRuntimes] = useState<Record<string, RuntimeStatus>>({});
  const onRuntime = useCallback((siteId: string, runtime: RuntimeStatus) => setRuntimes((current) => ({ ...current, [siteId]: runtime })), []);
  if (sites.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={5} /></section>;
  if (sites.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load sites" error={apiErrorMessage(sites.error)} /></section>;
  const siteList = sites.data ?? [];
  return <section className="dashboard-panel sites-page"><SiteRuntimeFanout sites={siteList} onRuntime={onRuntime} /><header><div><h1>Sites</h1><p>Manage hosted Goja sites for this organization.</p></div><button type="button" data-part="btn" onClick={() => navigate(`/app/orgs/${orgId}/sites/new`)}>New site</button></header><SitesTable sites={siteList} runtimes={runtimes} onOpenSite={(siteId) => navigate(`/app/orgs/${orgId}/sites/${siteId}`)} /></section>;
}
