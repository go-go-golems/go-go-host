import { type FormEvent, useMemo, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { EmptyState, ErrorCallout, LoadingBlock } from '../../components/atoms';
import { AuditTimeline } from '../../components/organisms';
import { apiErrorMessage } from '../../services/errors';
import { useListAdminAuditQuery } from '../../services/goGoHostApi';
import '../AuditPage/AuditPage.css';

export function AdminAuditPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [selectedId, setSelectedId] = useState<string | undefined>();
  const [actionDraft, setActionDraft] = useState(searchParams.get('action') ?? '');
  const [actorTypeDraft, setActorTypeDraft] = useState(searchParams.get('actorType') ?? '');
  const [orgDraft, setOrgDraft] = useState(searchParams.get('orgId') ?? '');
  const query = useMemo(() => ({ orgId: searchParams.get('orgId') || undefined, action: searchParams.get('action') || undefined, actorType: searchParams.get('actorType') || undefined, actorId: searchParams.get('actorId') || undefined, resourceId: searchParams.get('resourceId') || undefined, limit: 100 }), [searchParams]);
  const audit = useListAdminAuditQuery(query);
  function onFilter(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const next = new URLSearchParams(searchParams);
    actionDraft.trim() ? next.set('action', actionDraft.trim()) : next.delete('action');
    actorTypeDraft.trim() ? next.set('actorType', actorTypeDraft.trim()) : next.delete('actorType');
    orgDraft.trim() ? next.set('orgId', orgDraft.trim()) : next.delete('orgId');
    setSearchParams(next);
  }
  if (audit.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={5} /></section>;
  return <div className="audit-page">
    <section className="dashboard-panel audit-page__header"><h1>Global audit</h1><p>Platform-wide audit stream with URL-backed filters.</p><form onSubmit={onFilter}><input value={orgDraft} onChange={(event) => setOrgDraft(event.target.value)} placeholder="org_..." aria-label="Org filter" /><input value={actionDraft} onChange={(event) => setActionDraft(event.target.value)} placeholder="deployment.activate" aria-label="Action filter" /><input value={actorTypeDraft} onChange={(event) => setActorTypeDraft(event.target.value)} placeholder="user" aria-label="Actor type filter" /><button type="submit" data-part="btn">Apply filters</button><button type="button" data-part="btn" onClick={() => { setOrgDraft(''); setActionDraft(''); setActorTypeDraft(''); setSearchParams(new URLSearchParams()); }}>Clear</button></form></section>
    <section className="dashboard-panel">{audit.error ? <ErrorCallout title="Unable to load global audit" error={apiErrorMessage(audit.error)} /> : audit.data?.length ? <AuditTimeline events={audit.data} selectedId={selectedId} onSelect={setSelectedId} /> : <EmptyState title="No audit events" body="No events matched the current filters." />}</section>
  </div>;
}
