import type { Meta, StoryObj } from '@storybook/react';
import { QuotaPanel } from './QuotaPanel';
const meta = { title: 'Organisms/QuotaPanel', component: QuotaPanel, args: { sitesTotal: 2, requestsTotal: 1234, errorsTotal: 2 } } satisfies Meta<typeof QuotaPanel>;
export default meta; type Story = StoryObj<typeof meta>;
export const PreviewCounters: Story = {};
export const Empty: Story = { args: { sitesTotal: 0, requestsTotal: 0, errorsTotal: 0 } };
