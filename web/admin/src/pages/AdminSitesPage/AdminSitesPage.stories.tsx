import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { AdminSitesPage } from './AdminSitesPage';
const Wrapped = () => <MemoryRouter initialEntries={['/admin/sites']}><Routes><Route path="/admin/sites" element={<AdminSitesPage />} /></Routes></MemoryRouter>;
const meta = { title: 'Admin Pages/AdminSitesPage', component: Wrapped } satisfies Meta<typeof Wrapped>;
export default meta; type Story = StoryObj<typeof meta>;
export const Populated: Story = {};
export const Empty: Story = { parameters: { msw: { handlers: [http.get('/api/v1/admin/sites', () => HttpResponse.json([]))] } } };
