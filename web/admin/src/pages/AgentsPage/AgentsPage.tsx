import { type FormEvent, useState } from 'react';
import { useParams } from 'react-router-dom';
import { EmptyState, ErrorCallout, LoadingBlock } from '../../components/atoms';
import { AgentsTable } from '../../components/organisms';
import { apiErrorMessage } from '../../services/errors';
import { useCreateAgentMutation, useListAgentsQuery, useRevokeAgentMutation } from '../../services/goGoHostApi';
import './AgentsPage.css';

export function AgentsPage() {
  const { orgId } = useParams();
  const agents = useListAgentsQuery(orgId ?? '', { skip: !orgId });
  const [name, setName] = useState('');
  const [createAgent, createState] = useCreateAgentMutation();
  const [revokeAgent, revokeState] = useRevokeAgentMutation();
  async function onCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!orgId || !name.trim()) return;
    await createAgent({ orgId, name: name.trim() }).unwrap();
    setName('');
  }
  async function onRevoke(agentId: string) {
    if (!orgId || !window.confirm('Revoke this agent? Existing enrollment keys/grants are not implemented yet, but this disables the record.')) return;
    await revokeAgent({ orgId, agentId }).unwrap();
  }
  if (agents.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={5} /></section>;
  return <div className="agents-page">
    <section className="dashboard-panel agents-page__notice"><h1>Agents</h1><p>Preview: agent enrollment keys, grants, and deploy-run tokens are not implemented yet. This page manages the initial agent records exposed by the API.</p></section>
    <section className="dashboard-panel agents-page__create"><h2>Create agent</h2><form onSubmit={onCreate}><input value={name} onChange={(event) => setName(event.target.value)} placeholder="ci-bot" aria-label="Agent name" /><button type="submit" data-part="btn" disabled={!name.trim() || createState.isLoading}>{createState.isLoading ? 'Creating…' : 'Create agent'}</button></form>{createState.error ? <ErrorCallout title="Unable to create agent" error={apiErrorMessage(createState.error)} /> : null}{revokeState.error ? <ErrorCallout title="Unable to revoke agent" error={apiErrorMessage(revokeState.error)} /> : null}</section>
    <section className="dashboard-panel">{agents.error ? <ErrorCallout title="Unable to load agents" error={apiErrorMessage(agents.error)} /> : agents.data?.length ? <AgentsTable agents={agents.data} onRevoke={onRevoke} /> : <EmptyState title="No agents" body="Create an agent record to reserve an automation identity." />}</section>
  </div>;
}
