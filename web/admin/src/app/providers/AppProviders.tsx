import type { PropsWithChildren } from 'react';
import { Provider } from 'react-redux';
import { makeStore } from '../store';
import { SidebarExtraProvider } from './SidebarExtraProvider';

const store = makeStore();

export function AppProviders({ children }: PropsWithChildren) {
  return (
    <Provider store={store}>
      <SidebarExtraProvider>
        <div data-widget="hypercard" className="theme-macos1">{children}</div>
      </SidebarExtraProvider>
    </Provider>
  );
}
