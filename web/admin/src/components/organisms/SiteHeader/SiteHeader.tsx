import type { Site } from '../../../services/types';
import { SiteHostCopy } from '../../molecules';
import { StatusPill } from '../../atoms';
import './SiteHeader.css';
export function SiteHeader({ site, onBack }: { site: Site; onBack?: () => void }) {
  return <section className="site-header dashboard-panel"><button type="button" data-part="btn" onClick={onBack}>← Sites</button><div><h1>{site.name}</h1><p>{site.slug}</p></div><SiteHostCopy host={site.primaryHost} /><StatusPill status={site.status} /></section>;
}
