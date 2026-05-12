import type { Preview } from '@storybook/react';
import '../src/app/macos1-bridge.css';
import { initialize, mswLoader } from 'msw-storybook-addon';
import { MockAppProviders } from '../src/app/providers/MockAppProviders';
import { handlers } from '../src/services/msw/handlers';

initialize();

const preview: Preview = {
  decorators: [
    (Story) => (
      <MockAppProviders>
        <Story />
      </MockAppProviders>
    ),
  ],
  loaders: [mswLoader],
  parameters: {
    msw: { handlers },
    controls: { matchers: { color: /(background|color)$/i, date: /Date$/i } },
  },
};

export default preview;
