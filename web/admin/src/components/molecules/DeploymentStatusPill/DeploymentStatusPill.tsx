import { StatusPill } from '../../atoms/StatusPill';
import type { DeploymentStatus } from '../../../services/types';
export function DeploymentStatusPill({ status }: { status: DeploymentStatus }) { return <StatusPill status={status} />; }
