import type { Meta, StoryObj } from '@storybook/react';
import { AgentKeysTable } from './AgentKeysTable';

const meta = { title: 'Organisms/AgentKeysTable', component: AgentKeysTable } satisfies Meta<typeof AgentKeysTable>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Populated: Story = { args: { keys: [{ id: 'ak_123', agentId: 'agt_123', fingerprint: 'SHA256:abcdef0123456789', status: 'active', createdAt: new Date().toISOString(), lastUsedAt: new Date().toISOString() }, { id: 'ak_old', agentId: 'agt_123', fingerprint: 'SHA256:deadbeef00000000', status: 'revoked', createdAt: new Date(Date.now() - 86400000).toISOString(), revokedAt: new Date().toISOString() }] } };
export const Empty: Story = { args: { keys: [] } };
