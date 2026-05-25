import './StatusPill.css';

export type StatusTone = 'success' | 'danger' | 'warning' | 'info' | 'neutral';

const toneByStatus: Record<string, StatusTone> = {
  active: 'success',
  ready: 'success',
  validated: 'success',
  failed: 'danger',
  rejected: 'danger',
  stopped: 'neutral',
  superseded: 'neutral',
  uploaded: 'info',
  starting: 'info',
  provisioning: 'info',
  draining: 'warning',
  revoked: 'warning',
};

export interface StatusPillProps {
  status: string;
  tone?: StatusTone;
}

export function StatusPill({ status, tone = toneByStatus[status] ?? 'neutral' }: StatusPillProps) {
  return <span className="status-pill" data-tone={tone}>{status}</span>;
}
