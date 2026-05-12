import type { Meta, StoryObj } from '@storybook/react';
import { SiteHeader } from './SiteHeader';
import { fixtures } from '../../../services/msw/fixtures';
const meta = { title: 'Organisms/SiteHeader', component: SiteHeader, args: { site: fixtures.sites[0] } } satisfies Meta<typeof SiteHeader>;
export default meta; type Story = StoryObj<typeof meta>;
export const Active: Story = {};
export const Provisioning: Story = { args: { site: fixtures.sites[1] } };
