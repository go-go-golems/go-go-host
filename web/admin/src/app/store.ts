import { configureStore } from '@reduxjs/toolkit';
import { goGoHostApi } from '../services/goGoHostApi';

export function makeStore() {
  return configureStore({
    reducer: { [goGoHostApi.reducerPath]: goGoHostApi.reducer },
    middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(goGoHostApi.middleware),
  });
}

export type AppStore = ReturnType<typeof makeStore>;
export type RootState = ReturnType<AppStore['getState']>;
export type AppDispatch = AppStore['dispatch'];
