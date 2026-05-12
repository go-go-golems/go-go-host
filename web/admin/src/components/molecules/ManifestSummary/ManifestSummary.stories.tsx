import type { Meta, StoryObj } from '@storybook/react';
import { ManifestSummary } from './ManifestSummary';
const meta = { title: 'Molecules/ManifestSummary', component: ManifestSummary, args: { manifest: { scriptsDir: 'scripts', assetsDir: 'assets', smokePath: '/', capabilities: ['time', 'timer'] } } } satisfies Meta<typeof ManifestSummary>;
export default meta; type Story = StoryObj<typeof meta>;
export const Default: Story = {};
export const Minimal: Story = { args: { manifest: { scriptsDir: 'scripts' } } };
