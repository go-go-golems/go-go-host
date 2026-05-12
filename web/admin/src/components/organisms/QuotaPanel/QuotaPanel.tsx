import { MetricCard } from '../../molecules';
import './QuotaPanel.css';
export function QuotaPanel({ requestsTotal, errorsTotal, sitesTotal }: { requestsTotal: number; errorsTotal: number; sitesTotal: number }) {
  const errorRate = requestsTotal ? `${((errorsTotal / requestsTotal) * 100).toFixed(2)}%` : '0%';
  return <section className="quota-panel"><header><h2>Usage preview</h2><p>Runtime counters are aggregated from current site status. Error totals above zero are quota/health warning signals for operators.</p>{errorsTotal ? <p role="alert"><strong>Warning:</strong> one or more site runtimes reported errors. Check runtime pages and audit before activating more deployments.</p> : null}</header><div className="quota-panel__grid"><MetricCard label="Sites" value={sitesTotal} /><MetricCard label="Requests" value={requestsTotal} /><MetricCard label="Errors" value={errorsTotal} tone={errorsTotal ? 'danger' : 'default'} /><MetricCard label="Error rate" value={errorRate} /></div></section>;
}
