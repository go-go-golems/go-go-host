export type Role = 'org_owner' | 'org_developer' | 'org_viewer';
export type DeploymentStatus = 'uploaded' | 'validated' | 'rejected' | 'active' | 'superseded';
export type RuntimeState = 'starting' | 'ready' | 'failed' | 'stopped' | 'draining';

export interface OIDCConfig { issuer: string; clientId: string; deviceClientId?: string; scopes?: string[]; redirectPath?: string; logoutRedirectPath?: string; }
export interface ConfigResponse { baseDomain: string; publicBaseUrl: string; devAuth: boolean; oidc?: OIDCConfig; }
export interface User { id: string; email: string; displayName: string; }
export interface Membership { orgId: string; orgSlug: string; orgName: string; role: Role; }
export interface MeResponse { user: User; memberships: Membership[]; platformAdmin: boolean; }
export interface Org { id: string; slug: string; name: string; }
export interface CreateOrgRequest { slug: string; name: string; }
export interface Site { id: string; orgId: string; slug: string; name: string; primaryHost: string; status: string; activeDeploymentId: string; }
export interface CreateSiteRequest { orgId: string; slug: string; name: string; }
export interface SiteConfigItem { key: string; value: unknown; updatedAt: string; }
export interface UpsertSiteConfigRequest { siteId: string; key: string; value: unknown; }
export interface DeleteSiteConfigRequest { siteId: string; key: string; }
export interface SiteCapability { siteId: string; capability: string; enabled: boolean; config: unknown; updatedAt: string; }
export interface UpsertSiteCapabilityRequest { siteId: string; capability: string; enabled: boolean; config?: unknown; }
export interface SiteDomain { id: string; siteId: string; hostname: string; status: string; verificationToken: string; verifiedAt?: string; createdAt: string; }
export interface AddSiteDomainRequest { siteId: string; hostname: string; }
export interface VerifySiteDomainRequest { siteId: string; domainId: string; }
export interface DeleteSiteDomainRequest { siteId: string; domainId: string; }
export interface SiteEnvironmentPlaceholder { siteId: string; status: string; supported: string[]; notSupported: string[]; message: string; }
export interface Deployment { id: string; siteId: string; version: number; status: DeploymentStatus; bundleRef: string; unpackedPath: string; manifestJson: string; validationJson: string; createdByType: string; createdById: string; createdAt: string; activatedAt?: string; bundleSha256?: string; }
export interface RuntimeStatus { siteId: string; orgId?: string; deploymentId?: string; hosts?: string[]; status: RuntimeState; startedAt?: string; lastError?: string; requestsTotal?: number; errorsTotal?: number; }
export interface AdminRuntimeSummary { activeSites: number; hosts: string[]; runtimes: RuntimeStatus[]; }
export interface AdminOrg { id: string; slug: string; name: string; createdAt: string; memberCount: number; siteCount: number; deploymentCount: number; }
export interface AdminUser { id: string; email: string; displayName: string; createdAt: string; lastLoginAt?: string; platformAdmin: boolean; orgCount: number; }
export interface AdminSite { id: string; orgId: string; orgSlug: string; orgName: string; slug: string; name: string; primaryHost: string; status: string; activeDeploymentId: string; createdAt: string; runtimeStatus: RuntimeState; requestsTotal: number; errorsTotal: number; lastError?: string; }
export interface AdminDeployment { id: string; siteId: string; siteSlug: string; primaryHost: string; orgId: string; orgSlug: string; orgName: string; version: number; status: DeploymentStatus; bundleRef: string; unpackedPath: string; manifestJson: string; validationJson: string; createdByType: string; createdById: string; createdAt: string; activatedAt?: string; bundleSha256?: string; }
export interface AdminAgent { id: string; orgId: string; orgSlug: string; orgName: string; name: string; status: 'active' | 'revoked'; createdByUserId: string; createdAt: string; lastSeenAt?: string; grantCount: number; }
export interface AdminQuota { siteId: string; siteSlug: string; primaryHost: string; orgId: string; orgSlug: string; orgName: string; bundleMaxBytes: number; dbSoftMaxBytes: number; dbHardMaxBytes: number; requestTimeoutMs: number; updatedAt: string; requestsTotal: number; errorsTotal: number; }
export interface AdminCapability { siteId: string; siteSlug: string; orgId: string; orgSlug: string; orgName: string; capability: string; enabled: boolean; configJson: string; updatedAt: string; }
export interface AdminDomain { id: string; siteId: string; siteSlug: string; orgId: string; orgSlug: string; orgName: string; hostname: string; status: string; verificationToken: string; verifiedAt?: string; createdAt: string; }
export interface Agent { id: string; orgId: string; name: string; status: 'active' | 'revoked'; createdByUserId: string; createdAt: string; lastSeenAt?: string; }
export interface AgentKey { id: string; agentId: string; fingerprint: string; status: 'active' | 'revoked'; createdAt: string; revokedAt?: string; lastUsedAt?: string; }
export interface AgentGrant { agentId: string; siteId: string; canDeploy: boolean; canRollback: boolean; canActivate: boolean; allowedChannels: string[]; allowedBundlePaths?: string[]; allowedPaths: string[]; expiresAt?: string; createdAt: string; updatedAt: string; }
export interface CreateAgentRequest { orgId: string; name: string; siteId?: string; allowedChannels?: string[]; allowedBundlePaths?: string[]; allowedPaths?: string[]; canActivate?: boolean; }
export interface CreateAgentResponse { agent: Agent; enrollmentToken?: string; grant?: AgentGrant; }
export interface RevokeAgentRequest { orgId: string; agentId: string; }
export interface RevokeAgentKeyRequest { orgId: string; agentId: string; keyId: string; reason?: string; }
export interface CreateAgentEnrollmentTokenRequest { orgId: string; agentId: string; }
export interface CreateAgentEnrollmentTokenResponse { enrollmentToken: string; }
export interface AuditEvent { id: string; orgId: string; actorType: string; actorId: string; action: string; resourceType: string; resourceId: string; ipAddress: string; userAgent: string; metadataJson: string; createdAt: string; }
export interface ValidationReport { valid: boolean; errors?: string[]; warnings?: string[]; files: number; bytes: number; requestedCapabilities?: string[]; effectiveCapabilities?: string[]; }
export interface UploadDeploymentResponse { deployment: Deployment; report: ValidationReport; manifest: Record<string, unknown>; }

/* ── Docs ─────────────────────────────────────────────────── */
export type DocSection = 'Tutorial' | 'GeneralTopic' | 'Example' | 'Application' | '';
export interface DocEntry {
  slug: string;
  title: string;
  short: string;
  section: DocSection;
  source: 'host' | 'agent';
  body?: string;
}
