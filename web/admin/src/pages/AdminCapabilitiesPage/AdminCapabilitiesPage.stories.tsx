import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { AdminCapabilitiesPage } from './AdminCapabilitiesPage';
const Wrapped = () => <MemoryRouter initialEntries={['/admin/capabilities']}><Routes><Route path="/admin/capabilities" element={<AdminCapabilitiesPage />} /></Routes></MemoryRouter>;
const meta = { title: 'Admin Pages/AdminCapabilitiesPage', component: Wrapped } satisfies Meta<typeof Wrapped>;
export default meta; type Story = StoryObj<typeof meta>;
export const Populated: Story = {};
export const Empty: Story = { parameters: { msw: { handlers: [http.get('/api/v1/admin/capabilities', () => HttpResponse.json([]))] } } };
