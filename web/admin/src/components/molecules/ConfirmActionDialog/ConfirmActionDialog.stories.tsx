import type { Meta, StoryObj } from '@storybook/react';
import { ConfirmActionDialog } from './ConfirmActionDialog';
const meta = { title: 'Molecules/ConfirmActionDialog', component: ConfirmActionDialog } satisfies Meta<typeof ConfirmActionDialog>;
export default meta; type Story = StoryObj<typeof meta>;
export const StopRuntime: Story = { args: { open: true, title: 'Stop runtime?', body: 'Stop runtime for site_123. This operator action will be audited.', confirmLabel: 'Stop runtime', onCancel: () => {}, onConfirm: () => {} } };
export const Busy: Story = { args: { ...StopRuntime.args, busy: true } };
