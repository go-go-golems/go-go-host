import type { Meta, StoryObj } from '@storybook/react';
import { Timestamp } from './Timestamp';
const meta = { title: 'Atoms/Timestamp', component: Timestamp, args: { value: '2026-05-11T22:20:00Z' } } satisfies Meta<typeof Timestamp>;
export default meta; type Story = StoryObj<typeof meta>;
export const Absolute: Story = {};
export const Relative: Story = { args: { value: new Date(Date.now() - 5 * 60_000).toISOString(), mode: 'relative' } };
export const Empty: Story = { args: { value: '' } };
