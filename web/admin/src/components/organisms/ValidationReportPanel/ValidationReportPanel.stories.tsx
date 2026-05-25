import type { Meta, StoryObj } from '@storybook/react';
import { ValidationReportPanel } from './ValidationReportPanel';
import { fixtures } from '../../../services/msw/fixtures';
const meta = { title: 'Organisms/ValidationReportPanel', component: ValidationReportPanel, args: { deployment: fixtures.deployments[0] } } satisfies Meta<typeof ValidationReportPanel>;
export default meta; type Story = StoryObj<typeof meta>;
export const Valid: Story = {};
export const Rejected: Story = { args: { deployment: fixtures.deployments[2] } };
export const MalformedJson: Story = { args: { deployment: { ...fixtures.deployments[0], manifestJson: '{bad', validationJson: '{bad' } } };
