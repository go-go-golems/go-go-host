import type { Meta, StoryObj } from '@storybook/react';
import { LoadingBlock } from './LoadingBlock';
const meta = { title: 'Atoms/LoadingBlock', component: LoadingBlock, args: { lines: 3 } } satisfies Meta<typeof LoadingBlock>;
export default meta; type Story = StoryObj<typeof meta>;
export const Small: Story = { args: { lines: 2 } };
export const Large: Story = { args: { lines: 6 } };
