import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { AdminUsersPage } from './AdminUsersPage';
const Wrapped = () => <MemoryRouter initialEntries={['/admin/users']}><Routes><Route path="/admin/users" element={<AdminUsersPage />} /></Routes></MemoryRouter>;
const meta = { title: 'Admin Pages/AdminUsersPage', component: Wrapped } satisfies Meta<typeof Wrapped>;
export default meta; type Story = StoryObj<typeof meta>;
export const Populated: Story = {};
export const Empty: Story = { parameters: { msw: { handlers: [http.get('/api/v1/admin/users', () => HttpResponse.json([]))] } } };
