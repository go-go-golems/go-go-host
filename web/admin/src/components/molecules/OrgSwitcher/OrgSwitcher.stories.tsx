import type { Meta, StoryObj } from '@storybook/react';
import { OrgSwitcher } from './OrgSwitcher';
import { fixtures } from '../../../services/msw/fixtures';
const meta = { title: 'Molecules/OrgSwitcher', component: OrgSwitcher, args: { memberships: fixtures.me.memberships, selectedOrgId: 'org_123' } } satisfies Meta<typeof OrgSwitcher>;
export default meta; type Story = StoryObj<typeof meta>;
export const OneOrg: Story = {};
export const ManyOrgs: Story = { args: { memberships: [...fixtures.me.memberships, { orgId: 'org_456', orgSlug: 'ops', orgName: 'Ops Org', role: 'org_viewer' }] } };
export const NoOrgs: Story = { args: { memberships: [] } };
