import type { RuntimeStatus, Site } from '../../../services/types';
import { RuntimeBadge } from '../../molecules/RuntimeBadge';
import { SiteHostCopy } from '../../molecules/SiteHostCopy';
import { EmptyState } from '../../atoms/EmptyState';
import { StatusPill } from '../../atoms/StatusPill';
import './SitesTable.css';
export interface SitesTableProps { sites: Site[]; runtimes?: Record<string, RuntimeStatus>; onOpenSite?: (siteId: string) => void; }
export function SitesTable({ sites, runtimes = {}, onOpenSite }: SitesTableProps) {
  if (sites.length === 0) return <EmptyState title="No sites yet" body="Create a site or deploy from the CLI." />;
  return <table className="sites-table" data-part="sites-table"><thead><tr><th>Name</th><th>Host</th><th>Runtime</th><th>Active deployment</th><th>Status</th></tr></thead><tbody>{sites.map((site) => <tr key={site.id} onClick={() => onOpenSite?.(site.id)}><td><strong>{site.name}</strong><br /><small>{site.slug}</small></td><td><SiteHostCopy host={site.primaryHost} /></td><td>{runtimes[site.id] ? <RuntimeBadge runtime={runtimes[site.id]} compact /> : '—'}</td><td>{site.activeDeploymentId || '—'}</td><td><StatusPill status={site.status} /></td></tr>)}</tbody></table>;
}
