import { http, HttpResponse } from 'msw';
import { fixtures } from './fixtures';

export const handlers = [
  http.get('/api/v1/config', () => HttpResponse.json({ baseDomain: 'localhost', publicBaseUrl: 'http://127.0.0.1:8080', devAuth: true })),
  http.get('/api/v1/me', () => HttpResponse.json(fixtures.me)),
  http.get('/api/v1/orgs/:orgId/sites', () => HttpResponse.json(fixtures.sites)),
  http.post('/api/v1/orgs/:orgId/sites', async ({ params, request }) => {
    const body = await request.json() as { slug?: string; name?: string };
    if (!body.slug || !body.name) return HttpResponse.json({ error: 'slug and name are required' }, { status: 400 });
    return HttpResponse.json({ id: `site_${body.slug}`, orgId: String(params.orgId), slug: body.slug, name: body.name, primaryHost: `${body.slug}.localhost`, status: 'active', activeDeploymentId: '' }, { status: 201 });
  }),
  http.get('/api/v1/sites/:siteId/runtime', () => HttpResponse.json(fixtures.runtimeReady)),
  http.get('/api/v1/sites/:siteId/deployments', () => HttpResponse.json(fixtures.deployments)),
  http.get('/api/v1/orgs/:orgId/agents', () => HttpResponse.json(fixtures.agents)),
  http.get('/api/v1/orgs/:orgId/audit', () => HttpResponse.json(fixtures.audit)),
];
