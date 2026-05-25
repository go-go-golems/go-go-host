import type { Meta, StoryObj } from '@storybook/react';
import { AdminSidebar } from './AdminSidebar';

const meta = { title: 'Organisms/AdminSidebar', component: AdminSidebar } satisfies Meta<typeof AdminSidebar>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Overview: Story = { args: { active: 'overview' } };
export const Runtimes: Story = { args: { active: 'runtimes' } };
