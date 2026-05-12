import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import type { Agent, AuditEvent, ConfigResponse, CreateOrgRequest, CreateSiteRequest, Deployment, MeResponse, Org, RuntimeStatus, Site, UploadDeploymentResponse } from './types';

export interface UploadDeploymentRequest { siteId: string; file: File; message?: string; channel?: string; }

export const goGoHostApi = createApi({
  reducerPath: 'goGoHostApi',
  baseQuery: fetchBaseQuery({ baseUrl: '/api/v1' }),
  tagTypes: ['Me', 'Org', 'Site', 'Deployment', 'Runtime', 'Agent', 'Audit', 'Config'],
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
    listAgents: build.query<Agent[], string>({ query: (orgId) => `/orgs/${orgId}/agents`, providesTags: ['Agent'] }),
    listAudit: build.query<AuditEvent[], { orgId: string; action?: string; limit?: number }>({
      query: ({ orgId, ...params }) => ({ url: `/orgs/${orgId}/audit`, params }),
      providesTags: ['Audit'],
    }),
  }),
});

export const { useGetConfigQuery, useGetMeQuery, useCreateOrgMutation, useListSitesQuery, useCreateSiteMutation, useGetRuntimeQuery, useListDeploymentsQuery, useGetDeploymentQuery, useUploadDeploymentMutation, useActivateDeploymentMutation, useRollbackDeploymentMutation, useListAgentsQuery, useListAuditQuery } = goGoHostApi;
