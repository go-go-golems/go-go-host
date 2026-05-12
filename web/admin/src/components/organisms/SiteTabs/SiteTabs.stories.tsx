import type { Meta, StoryObj } from '@storybook/react';
import { MemoryRouter } from 'react-router-dom';
import { SiteTabs } from './SiteTabs';
const meta = { title: 'Organisms/SiteTabs', component: SiteTabs, decorators: [(Story) => <MemoryRouter initialEntries={['/app/orgs/org_123/sites/site_123']}><Story /></MemoryRouter>], args: { basePath: '/app/orgs/org_123/sites/site_123' } } satisfies Meta<typeof SiteTabs>;
export default meta; type Story = StoryObj<typeof meta>;
export const OverviewActive: Story = {};
