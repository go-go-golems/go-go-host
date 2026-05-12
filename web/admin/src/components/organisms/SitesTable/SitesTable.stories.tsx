import type { Meta, StoryObj } from '@storybook/react';
import { SitesTable } from './SitesTable';
import { fixtures } from '../../../services/msw/fixtures';
const meta = { title: 'Organisms/SitesTable', component: SitesTable, args: { sites: fixtures.sites, runtimes: { site_123: fixtures.runtimeReady, site_456: fixtures.runtimeStopped } } } satisfies Meta<typeof SitesTable>;
export default meta; type Story = StoryObj<typeof meta>;
export const Populated: Story = {};
export const Empty: Story = { args: { sites: [], runtimes: {} } };
export const FailedRuntime: Story = { args: { sites: fixtures.sites.slice(0, 1), runtimes: { site_123: fixtures.runtimeFailed } } };
