import { useCallback, useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import { ErrorCallout, LoadingBlock } from '../../components/atoms';
import { QuotaPanel, SitesTable } from '../../components/organisms';
import { apiErrorMessage } from '../../services/errors';
import { useListSitesQuery } from '../../services/goGoHostApi';
import type { RuntimeStatus } from '../../services/types';
import { SiteRuntimeFanout } from '../SitesPage/SiteRuntimeBadges';
export function UsagePage() {
  const { orgId } = useParams();
  const sites = useListSitesQuery(orgId ?? '', { skip: !orgId });
  const [runtimes, setRuntimes] = useState<Record<string, RuntimeStatus>>({});
  const onRuntime = useCallback((siteId: string, runtime: RuntimeStatus) => setRuntimes((current) => ({ ...current, [siteId]: runtime })), []);
  const totals = useMemo(() => Object.values(runtimes).reduce((acc, r) => ({ requests: acc.requests + (r.requestsTotal ?? 0), errors: acc.errors + (r.errorsTotal ?? 0) }), { requests: 0, errors: 0 }), [runtimes]);
  if (sites.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={5} /></section>;
  if (sites.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load usage" error={apiErrorMessage(sites.error)} /></section>;
  const siteList = sites.data ?? [];
  return <div className="usage-page" style={{ display: 'grid', gap: '1rem' }}><SiteRuntimeFanout sites={siteList} onRuntime={onRuntime} /><QuotaPanel sitesTotal={siteList.length} requestsTotal={totals.requests} errorsTotal={totals.errors} /><section className="dashboard-panel"><h2>Runtime counters by site</h2><SitesTable sites={siteList} runtimes={runtimes} /></section></div>;
}
