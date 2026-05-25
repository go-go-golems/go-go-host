import type { Preview } from '@storybook/react';
import '@go-go-golems/os-core/theme';
import '@go-go-golems/os-core/desktop-theme-macos1';
import '../src/app/macos1-bridge.css';
import { initialize, mswLoader } from 'msw-storybook-addon';
import { MockAppProviders } from '../src/app/providers/MockAppProviders';
import { handlers } from '../src/services/msw/handlers';

initialize({
  onUnhandledRequest(request, print) {
    const url = new URL(request.url);
    if (url.pathname.startsWith('/api/')) {
      print.warning();
    }
  },
});

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
