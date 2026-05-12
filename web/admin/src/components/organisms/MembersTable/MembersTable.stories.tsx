import type { Meta, StoryObj } from '@storybook/react';
import { MembersTable } from './MembersTable';
import { fixtures } from '../../../services/msw/fixtures';
const meta = { title: 'Organisms/MembersTable', component: MembersTable, args: { memberships: fixtures.me.memberships, selectedOrgId: 'org_123' } } satisfies Meta<typeof MembersTable>;
export default meta; type Story = StoryObj<typeof meta>;
export const CurrentOrg: Story = {};
export const RoleVariants: Story = { args: { selectedOrgId: 'org_123', memberships: [{ orgId: 'org_123', orgSlug: 'owner', orgName: 'Owner Org', role: 'org_owner' }, { orgId: 'org_456', orgSlug: 'dev', orgName: 'Dev Org', role: 'org_developer' }, { orgId: 'org_789', orgSlug: 'view', orgName: 'View Org', role: 'org_viewer' }] } };
