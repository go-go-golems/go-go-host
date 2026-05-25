import type { Meta, StoryObj } from '@storybook/react';
import { JsonTree } from './JsonTree';
const meta = { title: 'Atoms/JsonTree', component: JsonTree, args: { value: { valid: true, files: 3, bytes: 1024 } } } satisfies Meta<typeof JsonTree>;
export default meta; type Story = StoryObj<typeof meta>;
export const ValidationReport: Story = {};
export const Manifest: Story = { args: { value: { scriptsDir: 'scripts', assetsDir: 'assets', capabilities: ['time', 'timer'] } } };
