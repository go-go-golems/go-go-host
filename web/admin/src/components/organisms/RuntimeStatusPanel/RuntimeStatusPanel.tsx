import type { RuntimeStatus } from '../../../services/types';
import { RuntimeBadge } from '../../molecules/RuntimeBadge';
import { MetricCard } from '../../molecules/MetricCard';
import { Timestamp } from '../../atoms/Timestamp';
import './RuntimeStatusPanel.css';
export function RuntimeStatusPanel({ runtime }: { runtime: RuntimeStatus }) {
  const requests = runtime.requestsTotal ?? 0; const errors = runtime.errorsTotal ?? 0;
  const rate = requests > 0 ? `${((errors / requests) * 100).toFixed(2)}%` : '0%';
  return <section className="runtime-status-panel"><header><h2>Runtime</h2><RuntimeBadge runtime={runtime} /></header><div className="runtime-status-panel__grid"><MetricCard label="Requests" value={requests} /><MetricCard label="Errors" value={errors} tone={errors ? 'danger' : 'default'} /><MetricCard label="Error rate" value={rate} /></div><dl><dt>Deployment</dt><dd>{runtime.deploymentId || '—'}</dd><dt>Hosts</dt><dd>{runtime.hosts?.join(', ') || '—'}</dd><dt>Started</dt><dd><Timestamp value={runtime.startedAt} /></dd><dt>Last error</dt><dd>{runtime.lastError || '—'}</dd></dl></section>;
}
