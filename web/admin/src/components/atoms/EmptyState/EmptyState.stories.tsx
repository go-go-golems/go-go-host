import type { Meta, StoryObj } from '@storybook/react';
import { EmptyState } from './EmptyState';
const meta = { title: 'Atoms/EmptyState', component: EmptyState, args: { title: 'No sites yet', body: 'Create a site or deploy from the CLI.' } } satisfies Meta<typeof EmptyState>;
export default meta; type Story = StoryObj<typeof meta>;
export const WithoutAction: Story = {};
export const WithAction: Story = { args: { action: <button>Create site</button> } };
