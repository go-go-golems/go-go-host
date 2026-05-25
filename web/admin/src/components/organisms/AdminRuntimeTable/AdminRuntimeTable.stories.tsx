import type { Meta, StoryObj } from '@storybook/react';
import { MemoryRouter } from 'react-router-dom';
import { AdminRuntimeTable } from './AdminRuntimeTable';
import { fixtures } from '../../../services/msw/fixtures';

const meta = { title: 'Organisms/AdminRuntimeTable', component: AdminRuntimeTable, decorators: [(Story) => <MemoryRouter><Story /></MemoryRouter>] } satisfies Meta<typeof AdminRuntimeTable>;
export default meta;
type Story = StoryObj<typeof meta>;

export const HealthyAndFailed: Story = { args: { runtimes: fixtures.adminRuntimeSummary.runtimes } };
export const Empty: Story = { args: { runtimes: [] } };
export const FailedOnly: Story = { args: { runtimes: [fixtures.runtimeFailed] } };
