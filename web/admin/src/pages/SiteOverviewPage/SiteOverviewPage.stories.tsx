import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { SiteLayout } from '../../app/routing/SiteLayout';
import { SiteOverviewPage } from './SiteOverviewPage';
import { fixtures } from '../../services/msw/fixtures';

const Wrapped = ({ initialPath = '/app/orgs/org_123/sites/site_123' }: { initialPath?: string }) => <MemoryRouter initialEntries={[initialPath]}><Routes><Route path="/app/orgs/:orgId/sites/:siteId" element={<SiteLayout />}><Route index element={<SiteOverviewPage />} /><Route path="deployments" element={<div className="dashboard-panel">Deployments page placeholder</div>} /><Route path="deployments/:deploymentId" element={<div className="dashboard-panel">Deployment detail placeholder</div>} /></Route><Route path="/app/orgs/:orgId/sites" element={<div className="dashboard-panel">Sites list placeholder</div>} /></Routes></MemoryRouter>;
const meta = { title: 'Pages/SiteOverviewPage', component: Wrapped } satisfies Meta<typeof Wrapped>;
export default meta; type Story = StoryObj<typeof meta>;
export const Ready: Story = {};
export const StoppedRuntime: Story = { args: { initialPath: '/app/orgs/org_123/sites/site_456' }, parameters: { msw: { handlers: [http.get('/api/v1/sites/:siteId/runtime', () => HttpResponse.json(fixtures.runtimeStopped)), http.get('/api/v1/sites/:siteId/deployments', () => HttpResponse.json([]))] } } };
export const RuntimeLoadError: Story = { parameters: { msw: { handlers: [http.get('/api/v1/sites/:siteId/runtime', () => HttpResponse.json({ error: 'runtime recorder unavailable' }, { status: 500 }))] } } };
export const MissingSite: Story = { args: { initialPath: '/app/orgs/org_123/sites/site_missing' } };