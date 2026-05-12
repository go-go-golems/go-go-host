import type { Agent } from '../../../services/types';
import { AgentStatusBadge } from '../../molecules/AgentStatusBadge';
import { EmptyState } from '../../atoms/EmptyState';
import { Timestamp } from '../../atoms/Timestamp';
import './AgentsTable.css';
export function AgentsTable({ agents, onRevoke }: { agents: Agent[]; onRevoke?: (agentId: string) => void }) {
  if (!agents.length) return <EmptyState title="No agents yet" body="Create an agent record when automation setup is ready." />;
  return <table className="agents-table"><thead><tr><th>Name</th><th>Status</th><th>Created</th><th>Last seen</th><th>Actions</th></tr></thead><tbody>{agents.map((agent) => <tr key={agent.id}><td><strong>{agent.name}</strong><br /><small>{agent.id}</small></td><td><AgentStatusBadge status={agent.status} /></td><td><Timestamp value={agent.createdAt} /></td><td><Timestamp value={agent.lastSeenAt} /></td><td>{agent.status === 'active' ? <button type="button" data-part="btn" onClick={() => onRevoke?.(agent.id)}>Revoke</button> : '—'}</td></tr>)}</tbody></table>;
}
