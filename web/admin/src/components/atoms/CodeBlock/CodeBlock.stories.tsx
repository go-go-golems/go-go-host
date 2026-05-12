import type { Meta, StoryObj } from '@storybook/react';
import { CodeBlock } from './CodeBlock';
const meta = { title: 'Atoms/CodeBlock', component: CodeBlock, args: { language: 'shell', code: 'curl -H "Host: hello.localhost" http://127.0.0.1:8080/' } } satisfies Meta<typeof CodeBlock>;
export default meta; type Story = StoryObj<typeof meta>;
export const Shell: Story = {};
export const Json: Story = { args: { language: 'json', code: JSON.stringify({ scriptsDir: 'scripts', smokePath: '/' }, null, 2) } };
