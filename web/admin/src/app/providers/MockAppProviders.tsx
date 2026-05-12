import type { PropsWithChildren } from 'react';
import { Provider } from 'react-redux';
import { makeStore } from '../store';

export function MockAppProviders({ children }: PropsWithChildren) {
  return <Provider store={makeStore()}>{children}</Provider>;
}
