#!/usr/bin/env node
// Optional local OIDC browser smoke. Requires `devctl up --force` with the
// Keycloak profile. Uses the Playwright dependency from web/admin so `make
// web-install` is enough to prepare it. Gated so normal CI does not run it
// accidentally.

if (process.env.GO_GO_HOST_OIDC_E2E !== '1') {
  console.log('Skipping OIDC E2E; set GO_GO_HOST_OIDC_E2E=1 to run.');
  process.exit(0);
}

import { createRequire } from 'node:module';

const require = createRequire(import.meta.url);
const { chromium } = require('../web/admin/node_modules/playwright');

const baseURL = process.env.GO_GO_HOST_BASE_URL || 'http://127.0.0.1:8080';
const username = process.env.GO_GO_HOST_OIDC_USER || 'platform-admin';
const password = process.env.GO_GO_HOST_OIDC_PASSWORD || 'admin';

const browser = await chromium.launch({ headless: process.env.HEADLESS !== '0' });
const page = await browser.newPage();
try {
  await page.goto(`${baseURL}/admin`, { waitUntil: 'networkidle' });
  await page.getByRole('textbox', { name: /username|email/i }).fill(username);
  await page.getByRole('textbox', { name: /^password$/i }).fill(password);
  await page.getByRole('button', { name: /sign in/i }).click();
  await page.waitForURL((url) => url.href.startsWith(`${baseURL}/admin`), { timeout: 30000 });
  await page.getByRole('heading', { name: /^Platform admin$/ }).waitFor({ timeout: 30000 });

  const me = await page.evaluate(async () => {
    const raw = localStorage.getItem('go-go-host.oidc.tokens');
    if (!raw) throw new Error('missing stored OIDC tokens');
    const tokens = JSON.parse(raw);
    const response = await fetch('/api/v1/me', { headers: { Authorization: `Bearer ${tokens.idToken}` } });
    if (!response.ok) throw new Error(`/api/v1/me failed after browser login: ${response.status}`);
    return response.json();
  });
  if (!me.platformAdmin) throw new Error(`expected ${username} to be platform admin: ${JSON.stringify(me)}`);
  console.log(`OIDC E2E ok: ${me.user.email} platformAdmin=${me.platformAdmin}`);
} finally {
  await browser.close();
}
