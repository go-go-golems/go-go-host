import type { Meta, StoryObj } from '@storybook/react';
import { OrgSidebar } from './OrgSidebar';
const meta = { title: 'Organisms/OrgSidebar', component: OrgSidebar, args: { active: 'sites' } } satisfies Meta<typeof OrgSidebar>;
export default meta; type Story = StoryObj<typeof meta>;
export const Sites: Story = {};
export const Agents: Story = { args: { active: 'agents' } };
export const Audit: Story = { args: { active: 'audit' } };
