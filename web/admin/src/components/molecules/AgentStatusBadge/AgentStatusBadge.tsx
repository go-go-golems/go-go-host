import { StatusPill } from '../../atoms/StatusPill';
export function AgentStatusBadge({ status }: { status: 'active' | 'revoked' }) { return <StatusPill status={status} />; }
