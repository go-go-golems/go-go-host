import React from 'react';
import ReactDOM from 'react-dom/client';
import '@go-go-golems/os-core/theme.css';
import '@go-go-golems/os-core/themes/desktop.css';
import { App } from './app/App';
import { AppProviders } from './app/providers/AppProviders';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <AppProviders>
      <App />
    </AppProviders>
  </React.StrictMode>,
);
