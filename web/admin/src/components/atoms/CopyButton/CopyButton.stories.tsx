import type { Meta, StoryObj } from '@storybook/react';
import { CopyButton } from './CopyButton';
const meta = { title: 'Atoms/CopyButton', component: CopyButton, args: { value: 'hello.localhost', label: 'Copy host' } } satisfies Meta<typeof CopyButton>;
export default meta; type Story = StoryObj<typeof meta>;
export const Default: Story = {};
