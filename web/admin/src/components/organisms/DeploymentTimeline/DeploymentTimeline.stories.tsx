import type { Meta, StoryObj } from '@storybook/react';
import { DeploymentTimeline } from './DeploymentTimeline';
import { fixtures } from '../../../services/msw/fixtures';
const meta = { title: 'Organisms/DeploymentTimeline', component: DeploymentTimeline, args: { deployments: fixtures.deployments } } satisfies Meta<typeof DeploymentTimeline>;
export default meta; type Story = StoryObj<typeof meta>;
export const Mixed: Story = {};
export const Empty: Story = { args: { deployments: [] } };
