import type { Meta, StoryObj } from '@storybook/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { DocViewPage } from './DocViewPage';

const Wrapper = ({ slug, children }: { slug: string; children: React.ReactNode }) => (
  <MemoryRouter initialEntries={[`/app/orgs/org-1/docs/${slug}`]}>
    <Routes>
      <Route path="/app/orgs/:orgId/docs/:slug" element={children} />
    </Routes>
  </MemoryRouter>
);

const meta = {
  title: 'Pages/DocViewPage',
  component: DocViewPage,
  decorators: [(Story, context) => <Wrapper slug={context.parameters?.slug ?? 'developer-guide'}><Story /></Wrapper>],
} satisfies Meta<typeof DocViewPage>;

export default meta;
type Story = StoryObj<typeof meta>;

export const DeveloperGuide: Story = {
  parameters: { slug: 'developer-guide' },
};

export const JsApiReference: Story = {
  parameters: { slug: 'js-api-reference' },
};

export const DeployWorkflow: Story = {
  parameters: { slug: 'deploy-workflow' },
};

export const AgentGuide: Story = {
  parameters: { slug: 'agent-guide' },
};

export const AgentGettingStarted: Story = {
  parameters: { slug: 'agent-getting-started' },
};

export const NotFound: Story = {
  parameters: { slug: 'nonexistent-doc' },
};
