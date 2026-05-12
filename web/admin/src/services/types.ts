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
export interface Agent { id: string; orgId: string; name: string; status: 'active' | 'revoked'; createdByUserId: string; createdAt: string; lastSeenAt?: string; }
export interface AuditEvent { id: string; orgId: string; actorType: string; actorId: string; action: string; resourceType: string; resourceId: string; ipAddress: string; userAgent: string; metadataJson: string; createdAt: string; }
export interface ValidationReport { valid: boolean; errors?: string[]; warnings?: string[]; files: number; bytes: number; requestedCapabilities?: string[]; effectiveCapabilities?: string[]; }
export interface UploadDeploymentResponse { deployment: Deployment; report: ValidationReport; manifest: Record<string, unknown>; }
