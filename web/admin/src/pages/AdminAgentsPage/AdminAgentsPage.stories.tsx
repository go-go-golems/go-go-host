import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { AdminAgentsPage } from './AdminAgentsPage';
const Wrapped = () => <MemoryRouter initialEntries={['/admin/agents']}><Routes><Route path="/admin/agents" element={<AdminAgentsPage />} /></Routes></MemoryRouter>;
const meta = { title: 'Admin Pages/AdminAgentsPage', component: Wrapped } satisfies Meta<typeof Wrapped>;
export default meta; type Story = StoryObj<typeof meta>;
export const Populated: Story = {};
export const Empty: Story = { parameters: { msw: { handlers: [http.get('/api/v1/admin/agents', () => HttpResponse.json([]))] } } };
