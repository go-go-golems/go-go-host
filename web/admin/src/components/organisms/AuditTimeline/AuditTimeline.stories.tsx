import type { Meta, StoryObj } from '@storybook/react';
import { AuditTimeline } from './AuditTimeline';
import { fixtures } from '../../../services/msw/fixtures';
const meta = { title: 'Organisms/AuditTimeline', component: AuditTimeline, args: { events: fixtures.audit } } satisfies Meta<typeof AuditTimeline>;
export default meta; type Story = StoryObj<typeof meta>;
export const Populated: Story = {};
export const Empty: Story = { args: { events: [] } };
export const SelectedMetadata: Story = { args: { events: fixtures.audit, selectedId: 'aud_123' } };
