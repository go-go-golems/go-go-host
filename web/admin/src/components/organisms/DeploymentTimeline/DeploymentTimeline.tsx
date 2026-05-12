import type { Deployment } from '../../../services/types';
import { DeploymentStatusPill } from '../../molecules/DeploymentStatusPill';
import { Timestamp } from '../../atoms/Timestamp';
import { EmptyState } from '../../atoms/EmptyState';
import './DeploymentTimeline.css';
export interface DeploymentTimelineProps { deployments: Deployment[]; onSelect?: (deploymentId: string) => void; }
export function DeploymentTimeline({ deployments, onSelect }: DeploymentTimelineProps) {
  if (!deployments.length) return <EmptyState title="No deployments yet" body="Upload a bundle to create the first deployment." />;
  return <ol className="deployment-timeline">{deployments.map((dep) => <li key={dep.id}><button type="button" onClick={() => onSelect?.(dep.id)}><span>v{dep.version}</span><DeploymentStatusPill status={dep.status} /><span>{dep.id}</span><Timestamp value={dep.createdAt} /></button></li>)}</ol>;
}
