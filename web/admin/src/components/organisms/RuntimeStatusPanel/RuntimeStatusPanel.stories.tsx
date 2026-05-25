import type { Meta, StoryObj } from '@storybook/react';
import { RuntimeStatusPanel } from './RuntimeStatusPanel';
import { fixtures } from '../../../services/msw/fixtures';
const meta = { title: 'Organisms/RuntimeStatusPanel', component: RuntimeStatusPanel, args: { runtime: fixtures.runtimeReady } } satisfies Meta<typeof RuntimeStatusPanel>;
export default meta; type Story = StoryObj<typeof meta>;
export const Ready: Story = {};
export const Stopped: Story = { args: { runtime: fixtures.runtimeStopped } };
export const Failed: Story = { args: { runtime: fixtures.runtimeFailed } };
