import { useParams } from 'react-router-dom';
import { EmptyState, ErrorCallout, LoadingBlock } from '../../components/atoms';
import { MembersTable } from '../../components/organisms';
import { apiErrorMessage } from '../../services/errors';
import { useGetMeQuery } from '../../services/goGoHostApi';
export function MembersPage() { const { orgId } = useParams(); const me = useGetMeQuery(); if (me.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={5} /></section>; if (me.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load memberships" error={apiErrorMessage(me.error)} /></section>; const memberships = me.data?.memberships ?? []; return <div className="members-page" style={{ display: 'grid', gap: '1rem' }}><section className="dashboard-panel"><h1>Members</h1><p>Membership mutation APIs are pending. This page shows your current memberships and roles from /api/v1/me.</p></section><section className="dashboard-panel">{memberships.length ? <MembersTable memberships={memberships} selectedOrgId={orgId} /> : <EmptyState title="No memberships" body="Create or join an organization to see members." />}</section></div>; }
