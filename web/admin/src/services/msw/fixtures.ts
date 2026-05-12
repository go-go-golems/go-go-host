import type { AdminAgent, AdminCapability, AdminDeployment, AdminDomain, AdminOrg, AdminQuota, AdminRuntimeSummary, AdminSite, AdminUser, Agent, AuditEvent, Deployment, MeResponse, RuntimeStatus, Site, SiteCapability, SiteConfigItem, SiteDomain, SiteEnvironmentPlaceholder } from '../types';

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
  siteConfig: [
    { key: 'theme.title', value: { text: 'Hello Site' }, updatedAt: '2026-05-11T22:00:00Z' },
    { key: 'features.comments', value: false, updatedAt: '2026-05-11T22:01:00Z' },
  ] satisfies SiteConfigItem[],
  siteCapabilities: [
    { siteId: 'site_123', capability: 'express', enabled: true, config: {}, updatedAt: '2026-05-11T22:00:00Z' },
    { siteId: 'site_123', capability: 'assets', enabled: true, config: {}, updatedAt: '2026-05-11T22:00:00Z' },
    { siteId: 'site_123', capability: 'exec', enabled: false, config: { reason: 'never available in v1' }, updatedAt: '2026-05-11T22:00:00Z' },
  ] satisfies SiteCapability[],
  siteDomains: [
    { id: 'dom_site_123', siteId: 'site_123', hostname: 'hello.example.com', status: 'pending', verificationToken: 'ggh-verify-demo', createdAt: '2026-05-11T22:00:00Z' },
    { id: 'dom_site_456', siteId: 'site_123', hostname: 'hello.localhost', status: 'verified', verificationToken: '', verifiedAt: '2026-05-11T22:10:00Z', createdAt: '2026-05-11T21:00:00Z' },
  ] satisfies SiteDomain[],
  siteEnvironment: { siteId: 'site_123', status: 'design-placeholder', supported: ['non-secret site config via /config'], notSupported: ['process env passthrough', 'plaintext secret values in API responses'], message: 'Secrets/environment variables are intentionally not implemented in v1. Use non-secret site config only until encrypted secret storage and runtime injection are designed.' } satisfies SiteEnvironmentPlaceholder,
  adminMe: {
    user: { id: 'usr_admin', email: 'admin@dev.local', displayName: 'Platform Admin' },
    memberships: [{ orgId: 'org_123', orgSlug: 'demo', orgName: 'Demo Org', role: 'org_owner' }],
    platformAdmin: true,
  } satisfies MeResponse,
  adminRuntimeSummary: { activeSites: 1, hosts: ['hello.localhost'], runtimes: [
    { siteId: 'site_123', orgId: 'org_123', deploymentId: 'dep_4', hosts: ['hello.localhost'], status: 'ready', startedAt: '2026-05-11T22:20:00Z', requestsTotal: 1234, errorsTotal: 2 },
    { siteId: 'site_456', orgId: 'org_123', deploymentId: 'dep_bad', hosts: ['docs.localhost'], status: 'failed', startedAt: '2026-05-11T22:25:00Z', lastError: 'dry-run smoke check failed', requestsTotal: 15, errorsTotal: 7 },
  ] } satisfies AdminRuntimeSummary,
  adminOrgs: [
    { id: 'org_123', slug: 'demo', name: 'Demo Org', createdAt: '2026-05-11T20:00:00Z', memberCount: 3, siteCount: 2, deploymentCount: 4 },
    { id: 'org_456', slug: 'labs', name: 'Labs Org', createdAt: '2026-05-11T21:00:00Z', memberCount: 1, siteCount: 0, deploymentCount: 0 },
  ] satisfies AdminOrg[],
  adminUsers: [
    { id: 'usr_admin', email: 'admin@dev.local', displayName: 'Platform Admin', createdAt: '2026-05-11T19:00:00Z', platformAdmin: true, orgCount: 1 },
    { id: 'usr_123', email: 'alice@dev.local', displayName: 'Alice', createdAt: '2026-05-11T20:00:00Z', platformAdmin: false, orgCount: 1 },
  ] satisfies AdminUser[],
  adminSites: [
    { id: 'site_123', orgId: 'org_123', orgSlug: 'demo', orgName: 'Demo Org', slug: 'hello', name: 'Hello Site', primaryHost: 'hello.localhost', status: 'active', activeDeploymentId: 'dep_4', createdAt: '2026-05-11T20:10:00Z', runtimeStatus: 'ready', requestsTotal: 1234, errorsTotal: 2 },
    { id: 'site_456', orgId: 'org_123', orgSlug: 'demo', orgName: 'Demo Org', slug: 'docs', name: 'Docs Site', primaryHost: 'docs.localhost', status: 'active', activeDeploymentId: 'dep_bad', createdAt: '2026-05-11T20:15:00Z', runtimeStatus: 'failed', requestsTotal: 15, errorsTotal: 7, lastError: 'dry-run smoke check failed' },
  ] satisfies AdminSite[],
  adminAgents: [
    { id: 'agt_123', orgId: 'org_123', orgSlug: 'demo', orgName: 'Demo Org', name: 'ci-bot', status: 'active', createdByUserId: 'usr_123', createdAt: '2026-05-11T22:30:00Z', grantCount: 2 },
    { id: 'agt_456', orgId: 'org_123', orgSlug: 'demo', orgName: 'Demo Org', name: 'old-bot', status: 'revoked', createdByUserId: 'usr_admin', createdAt: '2026-05-11T21:30:00Z', lastSeenAt: '2026-05-11T22:00:00Z', grantCount: 0 },
  ] satisfies AdminAgent[],
  adminQuotas: [
    { siteId: 'site_123', siteSlug: 'hello', primaryHost: 'hello.localhost', orgId: 'org_123', orgSlug: 'demo', orgName: 'Demo Org', bundleMaxBytes: 52428800, dbSoftMaxBytes: 52428800, dbHardMaxBytes: 104857600, requestTimeoutMs: 2000, updatedAt: '2026-05-11T22:00:00Z', requestsTotal: 1234, errorsTotal: 2 },
    { siteId: 'site_456', siteSlug: 'docs', primaryHost: 'docs.localhost', orgId: 'org_123', orgSlug: 'demo', orgName: 'Demo Org', bundleMaxBytes: 52428800, dbSoftMaxBytes: 52428800, dbHardMaxBytes: 104857600, requestTimeoutMs: 2000, updatedAt: '2026-05-11T22:00:00Z', requestsTotal: 15, errorsTotal: 7 },
  ] satisfies AdminQuota[],
  adminCapabilities: [
    { siteId: 'site_123', siteSlug: 'hello', orgId: 'org_123', orgSlug: 'demo', orgName: 'Demo Org', capability: 'express', enabled: true, configJson: '{}', updatedAt: '2026-05-11T22:00:00Z' },
    { siteId: 'site_123', siteSlug: 'hello', orgId: 'org_123', orgSlug: 'demo', orgName: 'Demo Org', capability: 'exec', enabled: false, configJson: '{"reason":"never available in v1"}', updatedAt: '2026-05-11T22:00:00Z' },
  ] satisfies AdminCapability[],
  adminDomains: [
    { id: 'dom_123', siteId: 'site_123', siteSlug: 'hello', orgId: 'org_123', orgSlug: 'demo', orgName: 'Demo Org', hostname: 'hello.localhost', status: 'verified', verificationToken: '', verifiedAt: '2026-05-11T22:00:00Z', createdAt: '2026-05-11T20:00:00Z' },
    { id: 'dom_456', siteId: 'site_456', siteSlug: 'docs', orgId: 'org_123', orgSlug: 'demo', orgName: 'Demo Org', hostname: 'docs.example.com', status: 'pending', verificationToken: 'ggh-verify-demo', createdAt: '2026-05-11T20:30:00Z' },
  ] satisfies AdminDomain[],
  adminDeployments: [
    { id: 'dep_4', siteId: 'site_123', siteSlug: 'hello', primaryHost: 'hello.localhost', orgId: 'org_123', orgSlug: 'demo', orgName: 'Demo Org', version: 4, status: 'active', bundleRef: 'bundles/site_123/dep_4.tar.gz', unpackedPath: 'sites/site_123/deployments/dep_4', manifestJson: '{}', validationJson: '{"valid":true}', createdByType: 'user', createdById: 'usr_123', createdAt: '2026-05-11T22:10:00Z', activatedAt: '2026-05-11T22:20:00Z' },
    { id: 'dep_bad', siteId: 'site_456', siteSlug: 'docs', primaryHost: 'docs.localhost', orgId: 'org_123', orgSlug: 'demo', orgName: 'Demo Org', version: 2, status: 'rejected', bundleRef: 'bundles/site_456/dep_bad.tar.gz', unpackedPath: '', manifestJson: '{}', validationJson: '{"valid":false,"errors":["smoke failed"]}', createdByType: 'user', createdById: 'usr_123', createdAt: '2026-05-11T22:15:00Z' },
  ] satisfies AdminDeployment[],
};
