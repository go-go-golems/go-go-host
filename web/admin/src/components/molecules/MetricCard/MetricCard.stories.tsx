import type { Meta, StoryObj } from '@storybook/react';
import { MetricCard } from './MetricCard';
const meta = { title: 'Molecules/MetricCard', component: MetricCard, args: { label: 'Requests', value: '1,234', detail: 'since runtime start' } } satisfies Meta<typeof MetricCard>;
export default meta; type Story = StoryObj<typeof meta>;
export const Requests: Story = {};
export const Errors: Story = { args: { label: 'Errors', value: 2, tone: 'danger' } };
export const SuccessRate: Story = { args: { label: 'Success rate', value: '99.8%', tone: 'success' } };
