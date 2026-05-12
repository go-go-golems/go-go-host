import type { AuditEvent } from '../../../services/types';
import { EmptyState } from '../../atoms/EmptyState';
import { Timestamp } from '../../atoms/Timestamp';
import './AuditTimeline.css';
export function AuditTimeline({ events, selectedId, onSelect }: { events: AuditEvent[]; selectedId?: string; onSelect?: (id: string) => void }) {
  if (!events.length) return <EmptyState title="No audit events" body="Try adjusting filters or perform an org action." />;
  return <ol className="audit-timeline">{events.map((event) => <li key={event.id}><button type="button" data-selected={event.id === selectedId} onClick={() => onSelect?.(event.id)}><Timestamp value={event.createdAt} /><strong>{event.action}</strong><span>{event.actorType} {event.actorId}</span><span>{event.resourceType} {event.resourceId}</span></button>{event.id === selectedId ? <pre>{event.metadataJson}</pre> : null}</li>)}</ol>;
}
