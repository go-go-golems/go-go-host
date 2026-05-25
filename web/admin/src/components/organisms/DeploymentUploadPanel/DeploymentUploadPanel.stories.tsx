import type { Meta, StoryObj } from '@storybook/react';
import { http, HttpResponse } from 'msw';
import { DeploymentUploadPanel } from './DeploymentUploadPanel';
const meta = { title: 'Organisms/DeploymentUploadPanel', component: DeploymentUploadPanel, args: { siteId: 'site_123' } } satisfies Meta<typeof DeploymentUploadPanel>;
export default meta; type Story = StoryObj<typeof meta>;
export const Idle: Story = {};
export const ValidationRejected: Story = { parameters: { msw: { handlers: [http.post('/api/v1/sites/:siteId/deployments', () => HttpResponse.json({ deployment: { id: 'dep_bad', siteId: 'site_123', version: 5, status: 'rejected', bundleRef: '', unpackedPath: '', manifestJson: '{}', validationJson: '{"valid":false,"files":1,"bytes":20,"errors":["missing manifest"]}', createdByType: 'user', createdById: 'usr_123', createdAt: '2026-05-11T23:00:00Z' }, report: { valid: false, files: 1, bytes: 20, errors: ['missing manifest'] }, manifest: {} }, { status: 400 }))] } } };
