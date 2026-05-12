import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { AdminOverviewPage } from './AdminOverviewPage';
import { fixtures } from '../../services/msw/fixtures';

const Wrapped = () => <MemoryRouter initialEntries={['/admin/overview']}><Routes><Route path="/admin/overview" element={<AdminOverviewPage />} /></Routes></MemoryRouter>;
const meta = { title: 'Admin Pages/AdminOverviewPage', component: Wrapped } satisfies Meta<typeof Wrapped>;
export default meta;
type Story = StoryObj<typeof meta>;

export const WithRuntimes: Story = {};
export const Empty: Story = { parameters: { msw: { handlers: [http.get('/api/v1/admin/runtimes/summary', () => HttpResponse.json({ activeSites: 0, hosts: [], runtimes: [] }))] } } };
export const Forbidden: Story = { parameters: { msw: { handlers: [http.get('/api/v1/admin/runtimes/summary', () => HttpResponse.json({ error: 'platform admin required' }, { status: 403 }))] } } };
export const Degraded: Story = { parameters: { msw: { handlers: [http.get('/api/v1/admin/runtimes/summary', () => HttpResponse.json({ ...fixtures.adminRuntimeSummary, activeSites: 0 }))] } } };
