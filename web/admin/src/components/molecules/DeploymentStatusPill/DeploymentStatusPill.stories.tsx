import type { Meta, StoryObj } from '@storybook/react';
import { DeploymentStatusPill } from './DeploymentStatusPill';
const meta = { title: 'Molecules/DeploymentStatusPill', component: DeploymentStatusPill, args: { status: 'active' } } satisfies Meta<typeof DeploymentStatusPill>;
export default meta; type Story = StoryObj<typeof meta>;
export const Active: Story = { args: { status: 'active' } };
export const Validated: Story = { args: { status: 'validated' } };
export const Rejected: Story = { args: { status: 'rejected' } };
export const Superseded: Story = { args: { status: 'superseded' } };
