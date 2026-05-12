import { http, HttpResponse } from 'msw';
import { fixtures } from './fixtures';

export const handlers = [
  http.get('/api/v1/config', () => HttpResponse.json({ baseDomain: 'localhost', publicBaseUrl: 'http://127.0.0.1:8080', devAuth: true })),
  http.get('/api/v1/me', () => HttpResponse.json(fixtures.me)),
  http.get('/api/v1/orgs/:orgId/sites', () => HttpResponse.json(fixtures.sites)),
  http.get('/api/v1/sites/:siteId/runtime', () => HttpResponse.json(fixtures.runtimeReady)),
  http.get('/api/v1/sites/:siteId/deployments', () => HttpResponse.json(fixtures.deployments)),
  http.get('/api/v1/orgs/:orgId/agents', () => HttpResponse.json(fixtures.agents)),
  http.get('/api/v1/orgs/:orgId/audit', () => HttpResponse.json(fixtures.audit)),
];
