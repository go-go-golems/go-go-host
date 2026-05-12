export type Role = 'org_owner' | 'org_developer' | 'org_viewer';
export type DeploymentStatus = 'uploaded' | 'validated' | 'rejected' | 'active' | 'superseded';
export type RuntimeState = 'starting' | 'ready' | 'failed' | 'stopped' | 'draining';

export interface ConfigResponse { baseDomain: string; publicBaseUrl: string; devAuth: boolean; }
export interface User { id: string; email: string; displayName: string; }
export interface Membership { orgId: string; orgSlug: string; orgName: string; role: Role; }
export interface MeResponse { user: User; memberships: Membership[]; platformAdmin: boolean; }
export interface Org { id: string; slug: string; name: string; }
export interface CreateOrgRequest { slug: string; name: string; }
export interface Site { id: string; orgId: string; slug: string; name: string; primaryHost: string; status: string; activeDeploymentId: string; }
export interface CreateSiteRequest { orgId: string; slug: string; name: string; }
export interface Deployment { id: string; siteId: string; version: number; status: DeploymentStatus; bundleRef: string; unpackedPath: string; manifestJson: string; validationJson: string; createdByType: string; createdById: string; createdAt: string; activatedAt?: string; }
export interface RuntimeStatus { siteId: string; orgId?: string; deploymentId?: string; hosts?: string[]; status: RuntimeState; startedAt?: string; lastError?: string; requestsTotal?: number; errorsTotal?: number; }
export interface AdminRuntimeSummary { activeSites: number; hosts: string[]; runtimes: RuntimeStatus[]; }
export interface AdminOrg { id: string; slug: string; name: string; createdAt: string; memberCount: number; siteCount: number; deploymentCount: number; }
export interface AdminUser { id: string; email: string; displayName: string; createdAt: string; lastLoginAt?: string; platformAdmin: boolean; orgCount: number; }
export interface AdminSite { id: string; orgId: string; orgSlug: string; orgName: string; slug: string; name: string; primaryHost: string; status: string; activeDeploymentId: string; createdAt: string; runtimeStatus: RuntimeState; requestsTotal: number; errorsTotal: number; lastError?: string; }
export interface AdminDeployment { id: string; siteId: string; siteSlug: string; primaryHost: string; orgId: string; orgSlug: string; orgName: string; version: number; status: DeploymentStatus; bundleRef: string; unpackedPath: string; manifestJson: string; validationJson: string; createdByType: string; createdById: string; createdAt: string; activatedAt?: string; }
export interface Agent { id: string; orgId: string; name: string; status: 'active' | 'revoked'; createdByUserId: string; createdAt: string; lastSeenAt?: string; }
export interface CreateAgentRequest { orgId: string; name: string; }
export interface RevokeAgentRequest { orgId: string; agentId: string; }
export interface AuditEvent { id: string; orgId: string; actorType: string; actorId: string; action: string; resourceType: string; resourceId: string; ipAddress: string; userAgent: string; metadataJson: string; createdAt: string; }
export interface ValidationReport { valid: boolean; errors?: string[]; warnings?: string[]; files: number; bytes: number; requestedCapabilities?: string[]; effectiveCapabilities?: string[]; }
export interface UploadDeploymentResponse { deployment: Deployment; report: ValidationReport; manifest: Record<string, unknown>; }
