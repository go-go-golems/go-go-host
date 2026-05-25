import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { userEvent, within, expect } from '@storybook/test';
import { CreateSitePage } from './CreateSitePage';

const Wrapped = () => <MemoryRouter initialEntries={['/app/orgs/org_123/sites/new']}><Routes><Route path="/app/orgs/:orgId/sites/new" element={<CreateSitePage />} /><Route path="/app/orgs/:orgId/sites/:siteId" element={<div className="dashboard-panel">Created site detail placeholder</div>} /></Routes></MemoryRouter>;
const meta = { title: 'Pages/CreateSitePage', component: Wrapped } satisfies Meta<typeof Wrapped>;
export default meta; type Story = StoryObj<typeof meta>;
export const FormShell: Story = {};
export const InvalidSlug: Story = { play: async ({ canvasElement }) => { const canvas = within(canvasElement); await userEvent.type(canvas.getByLabelText('Slug'), '-Bad Slug-'); await userEvent.type(canvas.getByLabelText('Name'), 'Bad Slug Demo'); await userEvent.click(canvas.getByRole('button', { name: 'Create site' })); await expect(canvas.getByText('Fix site details')).toBeInTheDocument(); } };
export const SuccessfulCreate: Story = { play: async ({ canvasElement }) => { const canvas = within(canvasElement); await userEvent.type(canvas.getByLabelText('Slug'), 'fresh-site'); await userEvent.type(canvas.getByLabelText('Name'), 'Fresh Site'); await userEvent.click(canvas.getByRole('button', { name: 'Create site' })); await expect(await canvas.findByText('Created site detail placeholder')).toBeInTheDocument(); } };
export const Forbidden: Story = { parameters: { msw: { handlers: [http.post('/api/v1/orgs/:orgId/sites', () => HttpResponse.json({ error: 'permission denied' }, { status: 403 }))] } }, play: async ({ canvasElement }) => { const canvas = within(canvasElement); await userEvent.type(canvas.getByLabelText('Slug'), 'blocked-site'); await userEvent.type(canvas.getByLabelText('Name'), 'Blocked Site'); await userEvent.click(canvas.getByRole('button', { name: 'Create site' })); await expect(await canvas.findByText('Unable to create site')).toBeInTheDocument(); } };
