import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import type { AdminAgent, AdminCapability, AdminDeployment, AdminDomain, AdminOrg, AdminQuota, AdminRuntimeSummary, AdminSite, AdminUser, Agent, AuditEvent, ConfigResponse, CreateAgentRequest, CreateOrgRequest, CreateSiteRequest, Deployment, MeResponse, Org, RevokeAgentRequest, RuntimeStatus, Site, UploadDeploymentResponse } from './types';

export interface UploadDeploymentRequest { siteId: string; file: File; message?: string; channel?: string; }

export const goGoHostApi = createApi({
  reducerPath: 'goGoHostApi',
  baseQuery: fetchBaseQuery({ baseUrl: '/api/v1' }),
  tagTypes: ['Me', 'Org', 'Site', 'Deployment', 'Runtime', 'AdminRuntime', 'AdminInventory', 'Agent', 'Audit', 'Config'],
  endpoints: (build) => ({
    getConfig: build.query<ConfigResponse, void>({ query: () => '/config', providesTags: ['Config'] }),
    getMe: build.query<MeResponse, void>({ query: () => '/me', providesTags: ['Me', 'Org'] }),
    createOrg: build.mutation<Org, CreateOrgRequest>({
      query: (body) => ({ url: '/orgs', method: 'POST', body }),
      invalidatesTags: ['Me', 'Org'],
    }),
    listSites: build.query<Site[], string>({ query: (orgId) => `/orgs/${orgId}/sites`, providesTags: (_r, _e, orgId) => [{ type: 'Site', id: `ORG:${orgId}` }] }),
    createSite: build.mutation<Site, CreateSiteRequest>({
      query: ({ orgId, slug, name }) => ({ url: `/orgs/${orgId}/sites`, method: 'POST', body: { slug, name } }),
      invalidatesTags: (_r, _e, { orgId }) => [{ type: 'Site', id: `ORG:${orgId}` }, 'Me'],
    }),
    getRuntime: build.query<RuntimeStatus, string>({ query: (siteId) => `/sites/${siteId}/runtime`, providesTags: (_r, _e, siteId) => [{ type: 'Runtime', id: siteId }] }),
    getAdminRuntimeSummary: build.query<AdminRuntimeSummary, void>({ query: () => '/admin/runtimes/summary', providesTags: ['AdminRuntime'] }),
    restartAdminRuntime: build.mutation<RuntimeStatus, string>({
      query: (siteId) => ({ url: `/admin/runtimes/${siteId}/restart`, method: 'POST' }),
      invalidatesTags: ['AdminRuntime', 'Audit', { type: 'AdminInventory', id: 'AUDIT' }],
    }),
    stopAdminRuntime: build.mutation<RuntimeStatus, string>({
      query: (siteId) => ({ url: `/admin/runtimes/${siteId}/stop`, method: 'POST' }),
      invalidatesTags: ['AdminRuntime', 'Audit', { type: 'AdminInventory', id: 'AUDIT' }],
    }),
    listAdminOrgs: build.query<AdminOrg[], void>({ query: () => '/admin/orgs', providesTags: [{ type: 'AdminInventory', id: 'ORGS' }] }),
    listAdminUsers: build.query<AdminUser[], void>({ query: () => '/admin/users', providesTags: [{ type: 'AdminInventory', id: 'USERS' }] }),
    listAdminSites: build.query<AdminSite[], void>({ query: () => '/admin/sites', providesTags: [{ type: 'AdminInventory', id: 'SITES' }] }),
    listAdminDeployments: build.query<AdminDeployment[], { orgId?: string; siteId?: string; status?: string; limit?: number } | void>({
      query: (params) => ({ url: '/admin/deployments', params: params ?? undefined }),
      providesTags: [{ type: 'AdminInventory', id: 'DEPLOYMENTS' }],
    }),
    getAdminDeployment: build.query<AdminDeployment, string>({ query: (deploymentId) => `/admin/deployments/${deploymentId}`, providesTags: (_r, _e, deploymentId) => [{ type: 'AdminInventory', id: `DEPLOYMENT:${deploymentId}` }] }),
    listAdminAgents: build.query<AdminAgent[], { orgId?: string; status?: string } | void>({
      query: (params) => ({ url: '/admin/agents', params: params ?? undefined }),
      providesTags: [{ type: 'AdminInventory', id: 'AGENTS' }],
    }),
    listAdminAudit: build.query<AuditEvent[], { orgId?: string; action?: string; actorType?: string; actorId?: string; resourceId?: string; limit?: number } | void>({
      query: (params) => ({ url: '/admin/audit', params: params ?? undefined }),
      providesTags: [{ type: 'AdminInventory', id: 'AUDIT' }],
    }),
    listAdminQuotas: build.query<AdminQuota[], void>({ query: () => '/admin/quotas', providesTags: [{ type: 'AdminInventory', id: 'QUOTAS' }] }),
    listAdminCapabilities: build.query<AdminCapability[], void>({ query: () => '/admin/capabilities', providesTags: [{ type: 'AdminInventory', id: 'CAPABILITIES' }] }),
    listAdminDomains: build.query<AdminDomain[], void>({ query: () => '/admin/domains', providesTags: [{ type: 'AdminInventory', id: 'DOMAINS' }] }),
    listDeployments: build.query<Deployment[], string>({ query: (siteId) => `/sites/${siteId}/deployments`, providesTags: (_r, _e, siteId) => [{ type: 'Deployment', id: `SITE:${siteId}` }] }),
    getDeployment: build.query<Deployment, string>({ query: (deploymentId) => `/deployments/${deploymentId}`, providesTags: (_r, _e, deploymentId) => [{ type: 'Deployment', id: deploymentId }] }),
    uploadDeployment: build.mutation<UploadDeploymentResponse, UploadDeploymentRequest>({
      queryFn: async ({ siteId, file, message, channel }) => {
        const form = new FormData();
        form.append('bundle', file);
        if (message) form.append('message', message);
        if (channel) form.append('channel', channel);
        try {
          const response = await fetch(`/api/v1/sites/${siteId}/deployments`, { method: 'POST', body: form });
          const data = await response.json();
          if (!response.ok && !(data && data.deployment && data.report)) return { error: { status: response.status, data } };
          return { data: data as UploadDeploymentResponse };
        } catch (error) {
          return { error: { status: 'FETCH_ERROR', error: error instanceof Error ? error.message : String(error) } };
        }
      },
      invalidatesTags: (_r, _e, { siteId }) => [{ type: 'Deployment', id: `SITE:${siteId}` }, { type: 'Runtime', id: siteId }],
    }),
    activateDeployment: build.mutation<Deployment, string>({
      query: (deploymentId) => ({ url: `/deployments/${deploymentId}/activate`, method: 'POST' }),
      invalidatesTags: (r, _e, deploymentId) => [{ type: 'Deployment', id: deploymentId }, { type: 'Deployment', id: `SITE:${r?.siteId ?? 'unknown'}` }, { type: 'Runtime', id: r?.siteId ?? 'unknown' }],
    }),
    rollbackDeployment: build.mutation<Deployment, string>({
      query: (siteId) => ({ url: `/sites/${siteId}/rollback`, method: 'POST' }),
      invalidatesTags: (r, _e, siteId) => [{ type: 'Deployment', id: `SITE:${siteId}` }, { type: 'Runtime', id: siteId }, { type: 'Deployment', id: r?.id ?? 'unknown' }],
    }),
    listAgents: build.query<Agent[], string>({ query: (orgId) => `/orgs/${orgId}/agents`, providesTags: (_r, _e, orgId) => [{ type: 'Agent', id: `ORG:${orgId}` }] }),
    createAgent: build.mutation<Agent, CreateAgentRequest>({
      query: ({ orgId, name }) => ({ url: `/orgs/${orgId}/agents`, method: 'POST', body: { name } }),
      invalidatesTags: (_r, _e, { orgId }) => [{ type: 'Agent', id: `ORG:${orgId}` }, 'Audit'],
    }),
    revokeAgent: build.mutation<{ status: string; agentId: string }, RevokeAgentRequest>({
      query: ({ orgId, agentId }) => ({ url: `/orgs/${orgId}/agents/${agentId}/revoke`, method: 'POST' }),
      invalidatesTags: (_r, _e, { orgId }) => [{ type: 'Agent', id: `ORG:${orgId}` }, 'Audit'],
    }),
    listAudit: build.query<AuditEvent[], { orgId: string; action?: string; actorType?: string; actorId?: string; resourceId?: string; limit?: number }>({
      query: ({ orgId, ...params }) => ({ url: `/orgs/${orgId}/audit`, params }),
      providesTags: ['Audit'],
    }),
  }),
});

export const { useGetConfigQuery, useGetMeQuery, useCreateOrgMutation, useListSitesQuery, useCreateSiteMutation, useGetRuntimeQuery, useGetAdminRuntimeSummaryQuery, useRestartAdminRuntimeMutation, useStopAdminRuntimeMutation, useListAdminOrgsQuery, useListAdminUsersQuery, useListAdminSitesQuery, useListAdminDeploymentsQuery, useGetAdminDeploymentQuery, useListAdminAgentsQuery, useListAdminAuditQuery, useListAdminQuotasQuery, useListAdminCapabilitiesQuery, useListAdminDomainsQuery, useListDeploymentsQuery, useGetDeploymentQuery, useUploadDeploymentMutation, useActivateDeploymentMutation, useRollbackDeploymentMutation, useListAgentsQuery, useCreateAgentMutation, useRevokeAgentMutation, useListAuditQuery } = goGoHostApi;
