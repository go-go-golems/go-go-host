import type { PropsWithChildren } from 'react';
import { Provider } from 'react-redux';
import { makeStore } from '../store';

const store = makeStore();

export function AppProviders({ children }: PropsWithChildren) {
  return <Provider store={store}><div data-widget="hypercard" className="theme-macos1">{children}</div></Provider>;
}
