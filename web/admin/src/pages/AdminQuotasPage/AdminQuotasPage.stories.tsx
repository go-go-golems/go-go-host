import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { AdminQuotasPage } from './AdminQuotasPage';
const Wrapped = () => <MemoryRouter initialEntries={['/admin/quotas']}><Routes><Route path="/admin/quotas" element={<AdminQuotasPage />} /></Routes></MemoryRouter>;
const meta = { title: 'Admin Pages/AdminQuotasPage', component: Wrapped } satisfies Meta<typeof Wrapped>;
export default meta; type Story = StoryObj<typeof meta>;
export const Populated: Story = {};
export const Empty: Story = { parameters: { msw: { handlers: [http.get('/api/v1/admin/quotas', () => HttpResponse.json([]))] } } };
