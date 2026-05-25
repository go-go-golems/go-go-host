import type { Meta, StoryObj } from '@storybook/react';
import { RoleBadge } from './RoleBadge';
const meta = { title: 'Atoms/RoleBadge', component: RoleBadge, args: { role: 'org_owner' } } satisfies Meta<typeof RoleBadge>;
export default meta; type Story = StoryObj<typeof meta>;
export const Owner: Story = { args: { role: 'org_owner' } };
export const Developer: Story = { args: { role: 'org_developer' } };
export const Viewer: Story = { args: { role: 'org_viewer' } };
