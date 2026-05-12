import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import type { Agent, AuditEvent, ConfigResponse, Deployment, MeResponse, RuntimeStatus, Site } from './types';

export const goGoHostApi = createApi({
  reducerPath: 'goGoHostApi',
  baseQuery: fetchBaseQuery({ baseUrl: '/api/v1' }),
  tagTypes: ['Me', 'Org', 'Site', 'Deployment', 'Runtime', 'Agent', 'Audit', 'Config'],
  endpoints: (build) => ({
    getConfig: build.query<ConfigResponse, void>({ query: () => '/config', providesTags: ['Config'] }),
    getMe: build.query<MeResponse, void>({ query: () => '/me', providesTags: ['Me', 'Org'] }),
    listSites: build.query<Site[], string>({ query: (orgId) => `/orgs/${orgId}/sites`, providesTags: (_r, _e, orgId) => [{ type: 'Site', id: `ORG:${orgId}` }] }),
    getRuntime: build.query<RuntimeStatus, string>({ query: (siteId) => `/sites/${siteId}/runtime`, providesTags: (_r, _e, siteId) => [{ type: 'Runtime', id: siteId }] }),
    listDeployments: build.query<Deployment[], string>({ query: (siteId) => `/sites/${siteId}/deployments`, providesTags: (_r, _e, siteId) => [{ type: 'Deployment', id: `SITE:${siteId}` }] }),
    listAgents: build.query<Agent[], string>({ query: (orgId) => `/orgs/${orgId}/agents`, providesTags: ['Agent'] }),
    listAudit: build.query<AuditEvent[], { orgId: string; action?: string; limit?: number }>({
      query: ({ orgId, ...params }) => ({ url: `/orgs/${orgId}/audit`, params }),
      providesTags: ['Audit'],
    }),
  }),
});

export const { useGetConfigQuery, useGetMeQuery, useListSitesQuery, useGetRuntimeQuery, useListDeploymentsQuery, useListAgentsQuery, useListAuditQuery } = goGoHostApi;
