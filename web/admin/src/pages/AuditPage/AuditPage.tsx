import { type FormEvent, useMemo, useState } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';
import { EmptyState, ErrorCallout, LoadingBlock } from '../../components/atoms';
import { AuditTimeline } from '../../components/organisms';
import { apiErrorMessage } from '../../services/errors';
import { useListAuditQuery } from '../../services/goGoHostApi';
import './AuditPage.css';

export function AuditPage() {
  const { orgId } = useParams();
  const [searchParams, setSearchParams] = useSearchParams();
  const [selectedId, setSelectedId] = useState<string | undefined>();
  const [actionDraft, setActionDraft] = useState(searchParams.get('action') ?? '');
  const [actorTypeDraft, setActorTypeDraft] = useState(searchParams.get('actor_type') ?? '');
  const query = useMemo(() => ({ orgId: orgId ?? '', action: searchParams.get('action') || undefined, actorType: searchParams.get('actor_type') || undefined, actorId: searchParams.get('actor_id') || undefined, resourceId: searchParams.get('resource_id') || undefined, limit: 100 }), [orgId, searchParams]);
  const audit = useListAuditQuery(query, { skip: !orgId });
  function onFilter(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const next = new URLSearchParams(searchParams);
    actionDraft.trim() ? next.set('action', actionDraft.trim()) : next.delete('action');
    actorTypeDraft.trim() ? next.set('actor_type', actorTypeDraft.trim()) : next.delete('actor_type');
    setSearchParams(next);
  }
  if (audit.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={5} /></section>;
  return <div className="audit-page">
    <section className="dashboard-panel audit-page__header"><h1>Audit</h1><p>Filter organization audit events by action and actor type. Filters are preserved in the URL.</p><form onSubmit={onFilter}><input value={actionDraft} onChange={(event) => setActionDraft(event.target.value)} placeholder="deployment.activate" aria-label="Action filter" /><input value={actorTypeDraft} onChange={(event) => setActorTypeDraft(event.target.value)} placeholder="user" aria-label="Actor type filter" /><button type="submit" data-part="btn">Apply filters</button><button type="button" data-part="btn" onClick={() => { setActionDraft(''); setActorTypeDraft(''); setSearchParams(new URLSearchParams()); }}>Clear</button></form></section>
    <section className="dashboard-panel">{audit.error ? <ErrorCallout title="Unable to load audit" error={apiErrorMessage(audit.error)} /> : audit.data?.length ? <AuditTimeline events={audit.data} selectedId={selectedId} onSelect={setSelectedId} /> : <EmptyState title="No audit events" body="No events matched the current filters." />}</section>
  </div>;
}
