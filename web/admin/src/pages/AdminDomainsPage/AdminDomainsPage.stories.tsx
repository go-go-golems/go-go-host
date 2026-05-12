import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { AdminDomainsPage } from './AdminDomainsPage';
const Wrapped = () => <MemoryRouter initialEntries={['/admin/domains']}><Routes><Route path="/admin/domains" element={<AdminDomainsPage />} /></Routes></MemoryRouter>;
const meta = { title: 'Admin Pages/AdminDomainsPage', component: Wrapped } satisfies Meta<typeof Wrapped>;
export default meta; type Story = StoryObj<typeof meta>;
export const Populated: Story = {};
export const Empty: Story = { parameters: { msw: { handlers: [http.get('/api/v1/admin/domains', () => HttpResponse.json([]))] } } };
