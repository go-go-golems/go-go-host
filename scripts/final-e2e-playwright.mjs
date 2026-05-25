#!/usr/bin/env node
// Final HOST-001 E2E for local operator runs.
// Usage:
//   devctl up --force
//   GO_GO_HOST_E2E=1 GO_GO_HOST_TEST_BUNDLE=/tmp/go-go-host-test-bundle.tar.gz node scripts/final-e2e-playwright.mjs
// Requires the Playwright package to be available in the environment.

if (process.env.GO_GO_HOST_E2E !== '1') {
  console.log('Skipping: set GO_GO_HOST_E2E=1 to run the final browser/API smoke.');
  process.exit(0);
}

const { chromium } = await import('playwright');
const fs = await import('node:fs/promises');
const { spawnSync } = await import('node:child_process');
const baseURL = process.env.GO_GO_HOST_BASE_URL || 'http://127.0.0.1:8080';
const user = process.env.GO_GO_HOST_E2E_USER || `e2e-${Date.now()}`;
const bundlePath = process.env.GO_GO_HOST_TEST_BUNDLE || '/tmp/go-go-host-test-bundle.tar.gz';
const suffix = String(Date.now());

async function api(path, options = {}) {
  const headers = { 'X-Go-Go-Host-User': user, ...(options.headers || {}) };
  const response = await fetch(`${baseURL}${path}`, { ...options, headers });
  const text = await response.text();
  if (!response.ok) throw new Error(`${options.method || 'GET'} ${path}: ${response.status} ${text}`);
  return text ? JSON.parse(text) : null;
}

async function main() {
  const browser = await chromium.launch({ headless: true });
  try {
    const page = await browser.newPage({ extraHTTPHeaders: { 'X-Go-Go-Host-User': user } });
    await page.goto(`${baseURL}/app`, { waitUntil: 'networkidle' });
    await page.getByText(/Sites|Create your first organization|No organizations/i).first().waitFor({ timeout: 15000 });

    const org = await api('/api/v1/orgs', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ slug: `e2e-org-${suffix}`, name: 'E2E Org' }) });
    const site = await api(`/api/v1/orgs/${org.id}/sites`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ slug: `e2e-site-${suffix}`, name: 'E2E Site' }) });

    const form = new FormData();
    form.append('bundle', new Blob([await fs.readFile(bundlePath)]), 'bundle.tar.gz');
    const uploadResponse = await fetch(`${baseURL}/api/v1/sites/${site.id}/deployments`, { method: 'POST', headers: { 'X-Go-Go-Host-User': user }, body: form });
    const uploadText = await uploadResponse.text();
    if (!uploadResponse.ok) throw new Error(`bundle upload failed: ${uploadResponse.status} ${uploadText}`);
    const upload = JSON.parse(uploadText);
    const deployment = await api(`/api/v1/deployments/${upload.deployment.id}/activate`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: '{}' });
    await api(`/api/v1/sites/${site.id}/rollback`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: '{}' }).catch(() => null);

    const publicResponse = await fetch(`${baseURL}/`, { headers: { Host: site.primaryHost } });
    if (!publicResponse.ok) throw new Error(`public host failed: ${publicResponse.status} ${await publicResponse.text()}`);

    const agent = await api(`/api/v1/orgs/${org.id}/agents`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ name: 'e2e-agent', siteId: site.id, allowedChannels: ['default'], allowedPaths: ['**'] }) });
    const tmpConfig = `/tmp/go-go-host-e2e-agent-${suffix}.json`;
    spawnSync('go', ['run', './cmd/go-go-host-agent', 'keygen', '--config', tmpConfig, '--api-url', baseURL], { stdio: 'inherit', check: true });
    spawnSync('go', ['run', './cmd/go-go-host-agent', 'enroll', '--config', tmpConfig, '--api-url', baseURL, '--token', agent.enrollmentToken], { stdio: 'inherit', check: true });
    spawnSync('go', ['run', './cmd/go-go-host-agent', 'deploy', '--config', tmpConfig, '--bundle', bundlePath, '--site-id', site.id, '--channel', 'default', '--path', 'bundles/e2e.tar.gz'], { stdio: 'inherit', check: true });

    const audit = await api(`/api/v1/orgs/${org.id}/audit?limit=100`);
    for (const action of ['deployment.upload', 'deployment.activate', 'site.create']) {
      if (!audit.some((event) => event.action === action)) throw new Error(`missing audit action ${action}`);
    }

    await page.goto(`${baseURL}/app/orgs/${org.id}/sites/${site.id}`, { waitUntil: 'networkidle' });
    await page.getByText(/Deployments|Runtime|Site DTO/i).first().waitFor({ timeout: 15000 });
    console.log(JSON.stringify({ ok: true, orgId: org.id, siteId: site.id, host: site.primaryHost, deploymentId: deployment.id }, null, 2));
  } finally {
    await browser.close();
  }
}

await main();
