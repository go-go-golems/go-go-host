import type { Meta, StoryObj } from '@storybook/react';
import { NoOrgsPage } from './NoOrgsPage';
const meta = { title: 'Pages/NoOrgsPage', component: NoOrgsPage } satisfies Meta<typeof NoOrgsPage>;
export default meta; type Story = StoryObj<typeof meta>;
export const EmptyMemberships: Story = {};
