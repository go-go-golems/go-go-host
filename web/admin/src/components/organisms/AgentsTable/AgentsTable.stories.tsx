import type { Meta, StoryObj } from '@storybook/react';
import { AgentsTable } from './AgentsTable';
import { fixtures } from '../../../services/msw/fixtures';
const meta = { title: 'Organisms/AgentsTable', component: AgentsTable, args: { agents: fixtures.agents } } satisfies Meta<typeof AgentsTable>;
export default meta; type Story = StoryObj<typeof meta>;
export const Populated: Story = {};
export const Empty: Story = { args: { agents: [] } };
export const Revoked: Story = { args: { agents: [{ ...fixtures.agents[0], status: 'revoked' }] } };
