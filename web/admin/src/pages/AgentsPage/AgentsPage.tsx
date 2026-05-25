import { type FormEvent, useState } from 'react';
import { Checkbox } from '@go-go-golems/os-core';
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
    <section className="dashboard-panel agents-page__notice"><header><h1>Agents</h1><p>Agents are <span className="agents-page__highlight agents-page__highlight--info">machine identities</span> for CI deploys. <span className="agents-page__highlight agents-page__highlight--danger">Auto-activation is dangerous</span>: only grant it to <span className="agents-page__highlight agents-page__highlight--safe">trusted pipelines</span> that may promote validated code to live traffic.</p></header></section>
    <section className="dashboard-panel agents-page__create"><header><h2>Create agent</h2><p>Create the long-lived automation identity first, then enroll a signing key with the one-time token.</p></header><form onSubmit={onCreate}><label>Agent name<input data-part="field-input" value={name} onChange={(event) => setName(event.target.value)} placeholder="ci-bot" aria-label="Agent name" /></label><div className="agents-page__danger"><Checkbox label="Allow auto-activation for grants created with this agent" checked={canActivate} onChange={() => setCanActivate((value) => !value)} disabled={createState.isLoading} /><p>Use only for pipelines that are allowed to promote validated deployments to public traffic.</p></div><button type="submit" data-part="btn" disabled={!name.trim() || createState.isLoading}>{createState.isLoading ? 'Creating…' : 'Create agent'}</button></form>{createState.data?.enrollmentToken ? <p className="agents-page__token"><strong>Enrollment token:</strong> <code>{createState.data.enrollmentToken}</code></p> : null}{createState.error ? <ErrorCallout title="Unable to create agent" error={apiErrorMessage(createState.error)} /> : null}{revokeState.error ? <ErrorCallout title="Unable to revoke agent" error={apiErrorMessage(revokeState.error)} /> : null}{enrollmentTokenState.data?.enrollmentToken ? <p className="agents-page__token"><strong>Rotation enrollment token:</strong> <code>{enrollmentTokenState.data.enrollmentToken}</code></p> : null}{enrollmentTokenState.error ? <ErrorCallout title="Unable to create rotation token" error={apiErrorMessage(enrollmentTokenState.error)} /> : null}{revokeKeyState.error ? <ErrorCallout title="Unable to revoke key" error={apiErrorMessage(revokeKeyState.error)} /> : null}</section>
    <section className="dashboard-panel agents-page__section"><header><h2>Agent records</h2><p>Select an agent to inspect signing keys. Revocation disables future signed requests without deleting deployment history.</p></header>{agents.error ? <ErrorCallout title="Unable to load agents" error={apiErrorMessage(agents.error)} /> : agents.data?.length ? <AgentsTable agents={agents.data} onRevoke={onRevoke} onSelect={setSelectedAgentId} selectedAgentId={selectedAgentId} /> : <EmptyState title="No agents" body="Create an agent record to reserve an automation identity." />}</section>
    {selectedAgentId ? <section className="dashboard-panel agents-page__section"><header><h2>Signing keys</h2><p>Keys authenticate signed deployment requests. Rotation tokens are shown once and should be handled like secrets.</p></header><button type="button" data-part="btn" onClick={onCreateEnrollmentToken} disabled={enrollmentTokenState.isLoading}>{enrollmentTokenState.isLoading ? 'Creating token…' : 'Create rotation token'}</button>{selectedKeys.isLoading ? <LoadingBlock lines={3} /> : selectedKeys.error ? <ErrorCallout title="Unable to load keys" error={apiErrorMessage(selectedKeys.error)} /> : <AgentKeysTable keys={selectedKeys.data ?? []} onRevoke={onRevokeKey} />}</section> : null}
  </div>;
}
