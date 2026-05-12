import type { Meta, StoryObj } from '@storybook/react';
import { RuntimeStatusDot } from './RuntimeStatusDot';

const meta = { title: 'Atoms/RuntimeStatusDot', component: RuntimeStatusDot, args: { status: 'ready' } } satisfies Meta<typeof RuntimeStatusDot>;
export default meta;
type Story = StoryObj<typeof meta>;
export const Ready: Story = { args: { status: 'ready' } };
export const Starting: Story = { args: { status: 'starting' } };
export const Failed: Story = { args: { status: 'failed' } };
export const Stopped: Story = { args: { status: 'stopped' } };
export const Draining: Story = { args: { status: 'draining' } };
