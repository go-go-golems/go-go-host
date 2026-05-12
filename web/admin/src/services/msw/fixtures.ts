import type { Agent, AuditEvent, Deployment, MeResponse, RuntimeStatus, Site } from '../types';

export const fixtures = {
  me: {
    user: { id: 'usr_123', email: 'alice@dev.local', displayName: 'Alice' },
    memberships: [{ orgId: 'org_123', orgSlug: 'demo', orgName: 'Demo Org', role: 'org_owner' }],
    platformAdmin: false,
  } satisfies MeResponse,
  sites: [
    { id: 'site_123', orgId: 'org_123', slug: 'hello', name: 'Hello Site', primaryHost: 'hello.localhost', status: 'active', activeDeploymentId: 'dep_4' },
    { id: 'site_456', orgId: 'org_123', slug: 'docs', name: 'Docs Site', primaryHost: 'docs.localhost', status: 'provisioning', activeDeploymentId: '' },
  ] satisfies Site[],
  runtimeReady: { siteId: 'site_123', orgId: 'org_123', deploymentId: 'dep_4', hosts: ['hello.localhost'], status: 'ready', startedAt: '2026-05-11T22:20:00Z', requestsTotal: 1234, errorsTotal: 2 } satisfies RuntimeStatus,
  deployments: [
    { id: 'dep_4', siteId: 'site_123', version: 4, status: 'active', bundleRef: 'bundles/site_123/dep_4.tar.gz', unpackedPath: 'sites/site_123/deployments/dep_4', manifestJson: '{}', validationJson: '{"valid":true,"files":3,"bytes":1024}', createdByType: 'user', createdById: 'usr_123', createdAt: '2026-05-11T22:10:00Z', activatedAt: '2026-05-11T22:20:00Z' },
  ] satisfies Deployment[],
  agents: [{ id: 'agt_123', orgId: 'org_123', name: 'ci-bot', status: 'active', createdByUserId: 'usr_123', createdAt: '2026-05-11T22:30:00Z' }] satisfies Agent[],
  audit: [{ id: 'aud_123', orgId: 'org_123', actorType: 'user', actorId: 'usr_123', action: 'deployment.activate', resourceType: 'deployment', resourceId: 'dep_4', ipAddress: '', userAgent: '', metadataJson: '{}', createdAt: '2026-05-11T22:20:00Z' }] satisfies AuditEvent[],
};
