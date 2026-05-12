import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { AdminDeploymentDetailPage } from './AdminDeploymentDetailPage';
const Wrapped = () => <MemoryRouter initialEntries={['/admin/deployments/dep_4']}><Routes><Route path="/admin/deployments/:deploymentId" element={<AdminDeploymentDetailPage />} /></Routes></MemoryRouter>;
const meta = { title: 'Admin Pages/AdminDeploymentDetailPage', component: Wrapped } satisfies Meta<typeof Wrapped>;
export default meta; type Story = StoryObj<typeof meta>;
export const Active: Story = {};
export const NotFound: Story = { parameters: { msw: { handlers: [http.get('/api/v1/admin/deployments/:deploymentId', () => HttpResponse.json({ error: 'deployment not found' }, { status: 404 }))] } } };
