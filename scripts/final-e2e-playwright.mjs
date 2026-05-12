#!/usr/bin/env node
// Final HOST-001 E2E skeleton for local operator runs.
// Usage: GO_GO_HOST_E2E=1 node scripts/final-e2e-playwright.mjs
// Requires an already running devctl stack and the Playwright package installed.

if (process.env.GO_GO_HOST_E2E !== '1') {
  console.log('Skipping: set GO_GO_HOST_E2E=1 to run the final browser smoke.');
  process.exit(0);
}

const { chromium } = await import('playwright');
const baseURL = process.env.GO_GO_HOST_BASE_URL || 'http://127.0.0.1:8080';
const browser = await chromium.launch({ headless: true });
try {
  const page = await browser.newPage({ extraHTTPHeaders: { 'X-Go-Go-Host-User': process.env.GO_GO_HOST_E2E_USER || 'dev-user' } });
  await page.goto(`${baseURL}/app`, { waitUntil: 'networkidle' });
  await page.getByText(/Sites|Create your first organization|No organizations/i).first().waitFor({ timeout: 15000 });
  await page.goto(`${baseURL}/admin`, { waitUntil: 'networkidle' });
  await page.getByText(/Overview|Runtimes|Platform/i).first().waitFor({ timeout: 15000 });
  console.log('Final browser smoke passed: /app and /admin rendered.');
} finally {
  await browser.close();
}
