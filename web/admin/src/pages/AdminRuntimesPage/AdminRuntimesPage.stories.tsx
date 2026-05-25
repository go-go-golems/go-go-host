import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { AdminRuntimesPage } from './AdminRuntimesPage';

const Wrapped = () => <MemoryRouter initialEntries={['/admin/runtimes']}><Routes><Route path="/admin/runtimes" element={<AdminRuntimesPage />} /></Routes></MemoryRouter>;
const meta = { title: 'Admin Pages/AdminRuntimesPage', component: Wrapped } satisfies Meta<typeof Wrapped>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Populated: Story = {};
export const Empty: Story = { parameters: { msw: { handlers: [http.get('/api/v1/admin/runtimes/summary', () => HttpResponse.json({ activeSites: 0, hosts: [], runtimes: [] }))] } } };
export const LoadError: Story = { parameters: { msw: { handlers: [http.get('/api/v1/admin/runtimes/summary', () => HttpResponse.json({ error: 'database offline' }, { status: 500 }))] } } };
