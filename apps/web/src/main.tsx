import './lib/zod-config.ts';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { RouterProvider } from 'react-router/dom';
import { AppProviders } from './app/App.tsx';
import { createRouter } from './app/routes.tsx';
import { readConfig } from './lib/config.ts';
import { sourceFor } from './lib/source.ts';
import './styles/index.css';

const config = readConfig();
// Static sites search in the browser; the searcher is loaded on first use.
const source = sourceFor(config, async (q, signal) =>
  (await import('./lib/static-search.ts')).searchStatic(config.basePath, q, signal),
);

const root = document.getElementById('root');
if (root) {
  createRoot(root).render(
    <StrictMode>
      <AppProviders source={source}>
        <RouterProvider router={createRouter(config.basePath)} />
      </AppProviders>
    </StrictMode>,
  );
}
