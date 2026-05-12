import type { Preview } from '@storybook/react';
import '@go-go-golems/os-core/theme.css';
import '@go-go-golems/os-core/themes/desktop.css';
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
