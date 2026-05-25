import type { Meta, StoryObj } from '@storybook/react';
import { AppShell } from './AppShell';
import { OrgSidebar } from '../OrgSidebar';
import { fixtures } from '../../../services/msw/fixtures';
const meta = { title: 'Organisms/AppShell', component: AppShell, args: { memberships: fixtures.me.memberships, selectedOrgId: 'org_123', userLabel: 'alice@dev.local', devAuth: true, sidebar: <OrgSidebar active="sites" />, children: <div className="dashboard-panel"><h1>Dashboard content</h1><p>Retro macOS1 shell using go-go-os-core theme scope.</p></div> } } satisfies Meta<typeof AppShell>;
export default meta; type Story = StoryObj<typeof meta>;
export const Default: Story = {};
export const NoSidebar: Story = { args: { sidebar: undefined } };
export const NoOrganizations: Story = { args: { memberships: [], selectedOrgId: undefined } };
