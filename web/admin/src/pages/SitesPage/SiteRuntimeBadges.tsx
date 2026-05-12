import { useEffect } from 'react';
import { useGetRuntimeQuery } from '../../services/goGoHostApi';
import type { RuntimeStatus, Site } from '../../services/types';

function SiteRuntimeBadgeLoader({ site, onRuntime }: { site: Site; onRuntime: (siteId: string, runtime: RuntimeStatus) => void }) {
  const runtime = useGetRuntimeQuery(site.id);
  useEffect(() => { if (runtime.data) onRuntime(site.id, runtime.data); }, [onRuntime, runtime.data, site.id]);
  return null;
}

export function SiteRuntimeFanout({ sites, onRuntime }: { sites: Site[]; onRuntime: (siteId: string, runtime: RuntimeStatus) => void }) {
  return <>{sites.map((site) => <SiteRuntimeBadgeLoader key={site.id} site={site} onRuntime={onRuntime} />)}</>;
}
