import type { Meta, StoryObj } from '@storybook/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { CreateSitePage } from './CreateSitePage';
const Wrapped = () => <MemoryRouter initialEntries={['/app/orgs/org_123/sites/new']}><Routes><Route path="/app/orgs/:orgId/sites/new" element={<CreateSitePage />} /></Routes></MemoryRouter>;
const meta = { title: 'Pages/CreateSitePage', component: Wrapped } satisfies Meta<typeof Wrapped>;
export default meta; type Story = StoryObj<typeof meta>;
export const FormShell: Story = {};
