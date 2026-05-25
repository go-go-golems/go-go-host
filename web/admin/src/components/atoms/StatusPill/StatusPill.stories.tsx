import type { Meta, StoryObj } from '@storybook/react';
import { StatusPill } from './StatusPill';

const meta = {
  title: 'Atoms/StatusPill',
  component: StatusPill,
  args: { status: 'ready' },
} satisfies Meta<typeof StatusPill>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Ready: Story = { args: { status: 'ready' } };
export const Failed: Story = { args: { status: 'failed' } };
export const Rejected: Story = { args: { status: 'rejected' } };
export const Superseded: Story = { args: { status: 'superseded' } };
