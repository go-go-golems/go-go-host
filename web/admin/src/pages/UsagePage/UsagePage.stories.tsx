import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { UsagePage } from './UsagePage';
const Wrapped = () => <MemoryRouter initialEntries={['/app/orgs/org_123/usage']}><Routes><Route path="/app/orgs/:orgId/usage" element={<UsagePage />} /></Routes></MemoryRouter>;
const meta = { title: 'Pages/UsagePage', component: Wrapped } satisfies Meta<typeof Wrapped>;
export default meta; type Story = StoryObj<typeof meta>;
export const PreviewCounters: Story = {};
export const NoSites: Story = { parameters: { msw: { handlers: [http.get('/api/v1/orgs/:orgId/sites', () => HttpResponse.json([]))] } } };
