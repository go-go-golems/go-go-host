import { CopyButton } from '../../atoms/CopyButton';
import './SiteHostCopy.css';
export interface SiteHostCopyProps { host: string; publicBaseUrl?: string; }
export function SiteHostCopy({ host, publicBaseUrl = 'http://127.0.0.1:8080' }: SiteHostCopyProps) {
  const curl = `curl -H 'Host: ${host}' ${publicBaseUrl}/`;
  return <span className="site-host-copy" data-part="site-host-copy"><code>{host}</code><CopyButton value={host} label="Copy host" /><CopyButton value={curl} label="Copy curl" /></span>;
}
