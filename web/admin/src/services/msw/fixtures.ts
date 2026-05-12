import type { AdminRuntimeSummary, Agent, AuditEvent, Deployment, MeResponse, RuntimeStatus, Site } from '../types';

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
  runtimeStopped: { siteId: 'site_456', status: 'stopped' } satisfies RuntimeStatus,
  runtimeFailed: { siteId: 'site_123', orgId: 'org_123', deploymentId: 'dep_bad', hosts: ['hello.localhost'], status: 'failed', lastError: 'dry-run smoke check failed', requestsTotal: 5, errorsTotal: 5 } satisfies RuntimeStatus,
  deployments: [
    { id: 'dep_4', siteId: 'site_123', version: 4, status: 'active', bundleRef: 'bundles/site_123/dep_4.tar.gz', unpackedPath: 'sites/site_123/deployments/dep_4', manifestJson: '{}', validationJson: '{"valid":true,"files":3,"bytes":1024}', createdByType: 'user', createdById: 'usr_123', createdAt: '2026-05-11T22:10:00Z', activatedAt: '2026-05-11T22:20:00Z' },
    { id: 'dep_3', siteId: 'site_123', version: 3, status: 'superseded', bundleRef: 'bundles/site_123/dep_3.tar.gz', unpackedPath: 'sites/site_123/deployments/dep_3', manifestJson: '{}', validationJson: '{"valid":true,"files":3,"bytes":900}', createdByType: 'user', createdById: 'usr_123', createdAt: '2026-05-11T21:10:00Z', activatedAt: '2026-05-11T21:20:00Z' },
    { id: 'dep_2', siteId: 'site_123', version: 2, status: 'rejected', bundleRef: 'bundles/site_123/dep_2.tar.gz', unpackedPath: 'sites/site_123/deployments/dep_2', manifestJson: '{}', validationJson: '{"valid":false,"files":2,"bytes":500,"errors":["missing go-go-host.json manifest"]}', createdByType: 'user', createdById: 'usr_123', createdAt: '2026-05-11T20:10:00Z' },
  ] satisfies Deployment[],
  agents: [{ id: 'agt_123', orgId: 'org_123', name: 'ci-bot', status: 'active', createdByUserId: 'usr_123', createdAt: '2026-05-11T22:30:00Z' }] satisfies Agent[],
  audit: [{ id: 'aud_123', orgId: 'org_123', actorType: 'user', actorId: 'usr_123', action: 'deployment.activate', resourceType: 'deployment', resourceId: 'dep_4', ipAddress: '', userAgent: '', metadataJson: '{}', createdAt: '2026-05-11T22:20:00Z' }] satisfies AuditEvent[],
  adminMe: {
    user: { id: 'usr_admin', email: 'admin@dev.local', displayName: 'Platform Admin' },
    memberships: [{ orgId: 'org_123', orgSlug: 'demo', orgName: 'Demo Org', role: 'org_owner' }],
    platformAdmin: true,
  } satisfies MeResponse,
  adminRuntimeSummary: { activeSites: 1, hosts: ['hello.localhost'], runtimes: [
    { siteId: 'site_123', orgId: 'org_123', deploymentId: 'dep_4', hosts: ['hello.localhost'], status: 'ready', startedAt: '2026-05-11T22:20:00Z', requestsTotal: 1234, errorsTotal: 2 },
    { siteId: 'site_456', orgId: 'org_123', deploymentId: 'dep_bad', hosts: ['docs.localhost'], status: 'failed', startedAt: '2026-05-11T22:25:00Z', lastError: 'dry-run smoke check failed', requestsTotal: 15, errorsTotal: 7 },
  ] } satisfies AdminRuntimeSummary,
};
