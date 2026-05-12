import { http, HttpResponse } from 'msw';
import { fixtures } from './fixtures';

export const handlers = [
  http.get('/api/v1/config', () => HttpResponse.json({ baseDomain: 'localhost', publicBaseUrl: 'http://127.0.0.1:8080', devAuth: true })),
  http.get('/api/v1/me', () => HttpResponse.json(fixtures.me)),
  http.get('/api/v1/admin/runtimes/summary', () => HttpResponse.json(fixtures.adminRuntimeSummary)),
  http.post('/api/v1/orgs', async ({ request }) => {
    const body = await request.json() as { slug?: string; name?: string };
    if (!body.slug || !body.name) return HttpResponse.json({ error: 'slug and name are required' }, { status: 400 });
    return HttpResponse.json({ id: `org_${body.slug}`, slug: body.slug, name: body.name }, { status: 201 });
  }),
  http.get('/api/v1/orgs/:orgId/sites', () => HttpResponse.json(fixtures.sites)),
  http.post('/api/v1/orgs/:orgId/sites', async ({ params, request }) => {
    const body = await request.json() as { slug?: string; name?: string };
    if (!body.slug || !body.name) return HttpResponse.json({ error: 'slug and name are required' }, { status: 400 });
    return HttpResponse.json({ id: `site_${body.slug}`, orgId: String(params.orgId), slug: body.slug, name: body.name, primaryHost: `${body.slug}.localhost`, status: 'active', activeDeploymentId: '' }, { status: 201 });
  }),
  http.get('/api/v1/sites/:siteId/runtime', () => HttpResponse.json(fixtures.runtimeReady)),
  http.get('/api/v1/sites/:siteId/deployments', () => HttpResponse.json(fixtures.deployments)),
  http.post('/api/v1/sites/:siteId/deployments', async ({ params }) => HttpResponse.json({ deployment: fixtures.deployments[0], report: { valid: true, files: 3, bytes: 1024 }, manifest: { name: 'hello', entry: 'scripts/app.js', siteId: params.siteId } }, { status: 201 })),
  http.get('/api/v1/deployments/:deploymentId', ({ params }) => HttpResponse.json(fixtures.deployments.find((d) => d.id === params.deploymentId) ?? fixtures.deployments[0])),
  http.post('/api/v1/deployments/:deploymentId/activate', ({ params }) => HttpResponse.json({ ...fixtures.deployments[0], id: String(params.deploymentId), status: 'active' })),
  http.post('/api/v1/sites/:siteId/rollback', () => HttpResponse.json({ ...fixtures.deployments[1], status: 'active' })), 
  http.get('/api/v1/orgs/:orgId/agents', () => HttpResponse.json(fixtures.agents)),
  http.post('/api/v1/orgs/:orgId/agents', async ({ params, request }) => {
    const body = await request.json() as { name?: string };
    if (!body.name) return HttpResponse.json({ error: 'name is required' }, { status: 400 });
    return HttpResponse.json({ id: `agt_${body.name.replace(/[^a-z0-9]/gi, '_').toLowerCase()}`, orgId: String(params.orgId), name: body.name, status: 'active', createdByUserId: 'usr_123', createdAt: new Date('2026-05-11T23:30:00Z').toISOString() }, { status: 201 });
  }),
  http.post('/api/v1/orgs/:orgId/agents/:agentId/revoke', ({ params }) => HttpResponse.json({ status: 'revoked', agentId: String(params.agentId) })),
  http.get('/api/v1/orgs/:orgId/audit', ({ request }) => {
    const url = new URL(request.url);
    const action = url.searchParams.get('action');
    const events = action ? fixtures.audit.filter((event) => event.action.includes(action)) : fixtures.audit;
    return HttpResponse.json(events);
  }),
];
