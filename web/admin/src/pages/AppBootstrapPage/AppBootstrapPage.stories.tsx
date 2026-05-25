import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { AppBootstrapPage } from './AppBootstrapPage';

const meta = {
  title: 'Pages/AppBootstrapPage',
  component: AppBootstrapPage,
} satisfies Meta<typeof AppBootstrapPage>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
export const NoOrganizations: Story = {
  parameters: {
    msw: {
      handlers: [http.get('/api/v1/me', () => HttpResponse.json({ user: { id: 'usr_empty', email: 'empty@dev.local', displayName: 'Empty' }, memberships: [], platformAdmin: false }))],
    },
  },
};
export const SessionError: Story = {
  parameters: {
    msw: {
      handlers: [http.get('/api/v1/me', () => HttpResponse.json({ error: 'unauthorized' }, { status: 401 }))],
    },
  },
};
