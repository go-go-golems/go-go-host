import { JsonTree } from '../../atoms/JsonTree';
import './ManifestSummary.css';
export interface ManifestSummaryProps { manifest: Record<string, unknown>; }
export function ManifestSummary({ manifest }: ManifestSummaryProps) {
  return <section className="manifest-summary"><dl><dt>Scripts</dt><dd>{String(manifest.scriptsDir ?? '—')}</dd><dt>Assets</dt><dd>{String(manifest.assetsDir ?? '—')}</dd><dt>Smoke</dt><dd>{String(manifest.smokePath ?? '—')}</dd></dl><JsonTree value={manifest} /></section>;
}
