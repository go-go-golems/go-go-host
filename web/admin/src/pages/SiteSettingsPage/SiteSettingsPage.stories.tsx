import type { Meta, StoryObj } from '@storybook/react';
import { MemoryRouter, Outlet, Route, Routes } from 'react-router-dom';
import { SiteSettingsPage } from './SiteSettingsPage';
import { fixtures } from '../../services/msw/fixtures';

const meta: Meta<typeof SiteSettingsPage> = { title: 'Pages/SiteSettingsPage', component: SiteSettingsPage };
export default meta;
type Story = StoryObj<typeof SiteSettingsPage>;

function OutletContextShell() {
  return <Outlet context={{ site: fixtures.sites[0] }} />;
}

function StoryShell() {
  return <MemoryRouter initialEntries={[`/app/orgs/org_123/sites/${fixtures.sites[0].id}/settings`]}><Routes><Route path="/app/orgs/:orgId/sites/:siteId" element={<OutletContextShell />}><Route path="settings" element={<SiteSettingsPage />} /></Route></Routes></MemoryRouter>;
}

export const Default: Story = { render: () => <StoryShell /> };
