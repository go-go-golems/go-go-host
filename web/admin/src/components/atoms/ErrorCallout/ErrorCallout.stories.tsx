import type { Meta, StoryObj } from '@storybook/react';
import { ErrorCallout } from './ErrorCallout';
const meta = { title: 'Atoms/ErrorCallout', component: ErrorCallout, args: { title: 'Unable to load sites', error: 'The API returned 403 forbidden.' } } satisfies Meta<typeof ErrorCallout>;
export default meta; type Story = StoryObj<typeof meta>;
export const AuthError: Story = {};
export const WithRetry: Story = { args: { onRetry: () => undefined } };
