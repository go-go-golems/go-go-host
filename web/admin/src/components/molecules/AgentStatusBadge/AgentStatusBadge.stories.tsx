import type { Meta, StoryObj } from '@storybook/react';
import { AgentStatusBadge } from './AgentStatusBadge';
const meta = { title: 'Molecules/AgentStatusBadge', component: AgentStatusBadge, args: { status: 'active' } } satisfies Meta<typeof AgentStatusBadge>;
export default meta; type Story = StoryObj<typeof meta>;
export const Active: Story = { args: { status: 'active' } };
export const Revoked: Story = { args: { status: 'revoked' } };
