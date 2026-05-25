import './ErrorCallout.css';
export interface ErrorCalloutProps { title: string; error: string; retryLabel?: string; onRetry?: () => void; }
export function ErrorCallout({ title, error, retryLabel = 'Retry', onRetry }: ErrorCalloutProps) {
  return <section className="error-callout" role="alert"><h2>{title}</h2><p>{error}</p>{onRetry ? <button type="button" onClick={onRetry}>{retryLabel}</button> : null}</section>;
}
