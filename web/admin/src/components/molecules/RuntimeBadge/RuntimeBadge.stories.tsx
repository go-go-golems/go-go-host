import type { Meta, StoryObj } from '@storybook/react';
import { RuntimeBadge } from './RuntimeBadge';
import { fixtures } from '../../../services/msw/fixtures';
const meta = { title: 'Molecules/RuntimeBadge', component: RuntimeBadge, args: { runtime: fixtures.runtimeReady } } satisfies Meta<typeof RuntimeBadge>;
export default meta; type Story = StoryObj<typeof meta>;
export const Ready: Story = {};
export const Stopped: Story = { args: { runtime: fixtures.runtimeStopped } };
export const Failed: Story = { args: { runtime: fixtures.runtimeFailed } };
export const Compact: Story = { args: { compact: true } };
