import React from 'react';
import ReactDOM from 'react-dom/client';
import '@go-go-golems/os-core/theme';
import '@go-go-golems/os-core/desktop-theme-macos1';
import './app/macos1-bridge.css';
import { App } from './app/App';
import { AppProviders } from './app/providers/AppProviders';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <AppProviders>
      <App />
    </AppProviders>
  </React.StrictMode>,
);
