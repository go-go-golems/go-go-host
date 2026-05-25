import type { Meta, StoryObj } from '@storybook/react';
import { ValidationSummary } from './ValidationSummary';
const valid = { valid: true, files: 12, bytes: 44203, effectiveCapabilities: ['time', 'timer'] };
const invalid = { valid: false, files: 2, bytes: 500, errors: ['missing go-go-host.json manifest', 'capability "exec" is not permitted'] };
const meta = { title: 'Molecules/ValidationSummary', component: ValidationSummary, args: { report: valid } } satisfies Meta<typeof ValidationSummary>;
export default meta; type Story = StoryObj<typeof meta>;
export const Valid: Story = {};
export const Invalid: Story = { args: { report: invalid } };
export const WithJson: Story = { args: { report: invalid, showJson: true } };
