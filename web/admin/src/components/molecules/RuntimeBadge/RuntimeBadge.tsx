import { RuntimeStatusDot } from '../../atoms/RuntimeStatusDot';
import type { RuntimeStatus } from '../../../services/types';
import './RuntimeBadge.css';

export interface RuntimeBadgeProps {
  runtime: RuntimeStatus;
  compact?: boolean;
}

export function RuntimeBadge({ runtime, compact = false }: RuntimeBadgeProps) {
  const requests = runtime.requestsTotal ?? 0;
  const errors = runtime.errorsTotal ?? 0;
  return (
    <span className="runtime-badge" data-part="runtime-badge" data-state={runtime.status}>
      <RuntimeStatusDot status={runtime.status} />
      <span className="runtime-badge__status">{runtime.status}</span>
      {!compact ? <span className="runtime-badge__counts">{requests} req / {errors} err</span> : null}
    </span>
  );
}
