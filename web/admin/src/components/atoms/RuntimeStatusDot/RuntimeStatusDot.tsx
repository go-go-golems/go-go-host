import './RuntimeStatusDot.css';

export type RuntimeState = 'starting' | 'ready' | 'failed' | 'stopped' | 'draining';

const toneByStatus: Record<RuntimeState, string> = {
  ready: 'success',
  failed: 'danger',
  stopped: 'neutral',
  starting: 'info',
  draining: 'warning',
};

export interface RuntimeStatusDotProps {
  status: RuntimeState;
  label?: string;
}

export function RuntimeStatusDot({ status, label = `runtime ${status}` }: RuntimeStatusDotProps) {
  return <span className="runtime-status-dot" data-tone={toneByStatus[status]} aria-label={label} title={label} />;
}
