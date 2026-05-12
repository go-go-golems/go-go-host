import { type FormEvent, useState } from 'react';
import { useParams } from 'react-router-dom';
import { EmptyState, ErrorCallout, LoadingBlock } from '../../components/atoms';
import { AgentKeysTable, AgentsTable } from '../../components/organisms';
import { apiErrorMessage } from '../../services/errors';
import { useCreateAgentEnrollmentTokenMutation, useCreateAgentMutation, useListAgentKeysQuery, useListAgentsQuery, useRevokeAgentKeyMutation, useRevokeAgentMutation } from '../../services/goGoHostApi';
import './AgentsPage.css';

export function AgentsPage() {
  const { orgId } = useParams();
  const agents = useListAgentsQuery(orgId ?? '', { skip: !orgId });
  const [name, setName] = useState('');
  const [selectedAgentId, setSelectedAgentId] = useState('');
  const [canActivate, setCanActivate] = useState(false);
  const [createAgent, createState] = useCreateAgentMutation();
  const [createEnrollmentToken, enrollmentTokenState] = useCreateAgentEnrollmentTokenMutation();
  const [revokeAgent, revokeState] = useRevokeAgentMutation();
  const [revokeAgentKey, revokeKeyState] = useRevokeAgentKeyMutation();
  const selectedKeys = useListAgentKeysQuery({ orgId: orgId ?? '', agentId: selectedAgentId }, { skip: !orgId || !selectedAgentId });
  async function onCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!orgId || !name.trim()) return;
    if (canActivate && !window.confirm('Grant this agent auto-activation? This lets the agent promote validated code to live traffic for granted sites.')) return;
    await createAgent({ orgId, name: name.trim(), canActivate }).unwrap();
    setName('');
    setCanActivate(false);
  }
  async function onRevoke(agentId: string) {
    if (!orgId || !window.confirm('Revoke this agent? This disables the agent record.')) return;
    await revokeAgent({ orgId, agentId }).unwrap();
  }
  async function onCreateEnrollmentToken() {
    if (!orgId || !selectedAgentId || !window.confirm('Create a new one-time enrollment token for key rotation? The token will be shown once.')) return;
    await createEnrollmentToken({ orgId, agentId: selectedAgentId }).unwrap();
  }
  async function onRevokeKey(keyId: string) {
    if (!orgId || !selectedAgentId || !window.confirm('Revoke this signing key? Existing deployments stay intact, but future signed requests with this key will fail.')) return;
    await revokeAgentKey({ orgId, agentId: selectedAgentId, keyId, reason: 'revoked from dashboard' }).unwrap();
  }
  if (agents.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={5} /></section>;
  return <div className="agents-page">
    <section className="dashboard-panel agents-page__notice"><h1>Agents</h1><p>Agents are machine identities for CI deploys. Auto-activation is dangerous: only grant it to trusted pipelines that may promote validated code to live traffic.</p></section>
    <section className="dashboard-panel agents-page__create"><h2>Create agent</h2><form onSubmit={onCreate}><input value={name} onChange={(event) => setName(event.target.value)} placeholder="ci-bot" aria-label="Agent name" /><label className="agents-page__danger"><input type="checkbox" checked={canActivate} onChange={(event) => setCanActivate(event.target.checked)} /> Allow auto-activation for grants created with this agent</label><button type="submit" data-part="btn" disabled={!name.trim() || createState.isLoading}>{createState.isLoading ? 'Creating…' : 'Create agent'}</button></form>{createState.data?.enrollmentToken ? <p><strong>Enrollment token:</strong> <code>{createState.data.enrollmentToken}</code></p> : null}{createState.error ? <ErrorCallout title="Unable to create agent" error={apiErrorMessage(createState.error)} /> : null}{revokeState.error ? <ErrorCallout title="Unable to revoke agent" error={apiErrorMessage(revokeState.error)} /> : null}{enrollmentTokenState.data?.enrollmentToken ? <p><strong>Rotation enrollment token:</strong> <code>{enrollmentTokenState.data.enrollmentToken}</code></p> : null}{enrollmentTokenState.error ? <ErrorCallout title="Unable to create rotation token" error={apiErrorMessage(enrollmentTokenState.error)} /> : null}{revokeKeyState.error ? <ErrorCallout title="Unable to revoke key" error={apiErrorMessage(revokeKeyState.error)} /> : null}</section>
    <section className="dashboard-panel">{agents.error ? <ErrorCallout title="Unable to load agents" error={apiErrorMessage(agents.error)} /> : agents.data?.length ? <AgentsTable agents={agents.data} onRevoke={onRevoke} onSelect={setSelectedAgentId} selectedAgentId={selectedAgentId} /> : <EmptyState title="No agents" body="Create an agent record to reserve an automation identity." />}</section>
    {selectedAgentId ? <section className="dashboard-panel"><h2>Signing keys</h2><button type="button" data-part="btn" onClick={onCreateEnrollmentToken} disabled={enrollmentTokenState.isLoading}>{enrollmentTokenState.isLoading ? 'Creating token…' : 'Create rotation token'}</button>{selectedKeys.isLoading ? <LoadingBlock lines={3} /> : selectedKeys.error ? <ErrorCallout title="Unable to load keys" error={apiErrorMessage(selectedKeys.error)} /> : <AgentKeysTable keys={selectedKeys.data ?? []} onRevoke={onRevokeKey} />}</section> : null}
  </div>;
}
