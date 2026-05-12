import type { AgentKey } from '../../../services/types';
import { EmptyState, Timestamp } from '../../atoms';
import './AgentKeysTable.css';

export function AgentKeysTable({ keys, onRevoke }: { keys: AgentKey[]; onRevoke?: (keyId: string) => void }) {
  if (!keys.length) return <EmptyState title="No keys" body="This agent has not enrolled any signing keys yet." />;
  return <table className="agent-keys-table"><thead><tr><th>Key</th><th>Status</th><th>Created</th><th>Last used</th><th>Revoked</th><th>Actions</th></tr></thead><tbody>{keys.map((key) => <tr key={key.id}><td><strong>{key.fingerprint}</strong><br /><small>{key.id}</small></td><td><span className={`agent-keys-table__status agent-keys-table__status--${key.status}`}>{key.status}</span></td><td><Timestamp value={key.createdAt} /></td><td>{key.lastUsedAt ? <Timestamp value={key.lastUsedAt} /> : '—'}</td><td>{key.revokedAt ? <Timestamp value={key.revokedAt} /> : '—'}</td><td>{key.status === 'active' ? <button type="button" data-part="btn" onClick={() => onRevoke?.(key.id)}>Revoke key</button> : '—'}</td></tr>)}</tbody></table>;
}
