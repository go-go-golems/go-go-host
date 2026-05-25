import type { Meta, StoryObj } from '@storybook/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { DocsIndexPage } from './DocsIndexPage';

const Wrapper = ({ children }: { children: React.ReactNode }) => (
  <MemoryRouter initialEntries={['/app/orgs/org-1/docs']}>
    <Routes>
      <Route path="/app/orgs/:orgId/docs" element={children} />
      <Route path="/app/orgs/:orgId/docs/:slug" element={children} />
    </Routes>
  </MemoryRouter>
);

const meta = {
  title: 'Pages/DocsIndexPage',
  component: DocsIndexPage,
  decorators: [(Story) => <Wrapper><Story /></Wrapper>],
} satisfies Meta<typeof DocsIndexPage>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
