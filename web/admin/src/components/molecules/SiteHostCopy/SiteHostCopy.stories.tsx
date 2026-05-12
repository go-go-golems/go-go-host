import type { Meta, StoryObj } from '@storybook/react';
import { SiteHostCopy } from './SiteHostCopy';
const meta = { title: 'Molecules/SiteHostCopy', component: SiteHostCopy, args: { host: 'hello.localhost' } } satisfies Meta<typeof SiteHostCopy>;
export default meta; type Story = StoryObj<typeof meta>;
export const Localhost: Story = {};
export const PublicDomain: Story = { args: { host: 'hello.example.com', publicBaseUrl: 'https://host.example.com' } };
