import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

// The Go server (`notes serve`) answers the API, assets, and config.js in development.
const goServer = 'http://127.0.0.1:8080';

export default defineConfig({
  plugins: [react(), tailwindcss()],
  experimental: {
    // Chunks and their CSS load from the base path set at runtime by config.js,
    // so one build works at the domain root (server) and under a sub-path
    // such as /notes/ (static export). HTML URLs stay root-absolute and are
    // rewritten by `notes export`.
    renderBuiltUrl(filename, { hostType }) {
      if (hostType === 'js') {
        return {
          runtime: `((globalThis.__NOTES_CONFIG__ && globalThis.__NOTES_CONFIG__.basePath) || "/") + ${JSON.stringify(filename)}`,
        };
      }
      return { relative: hostType === 'css' };
    },
  },
  build: {
    // The Go binary embeds this folder (server/internal/webui).
    outDir: '../../server/internal/webui/static',
    emptyOutDir: true,
    // /assets/ is reserved for content assets served by the Go server.
    assetsDir: 'static',
    // Never inline fonts as data: URLs; the CSP allows fonts from 'self' only.
    assetsInlineLimit: (file) => (/\.(woff2?|ttf|otf)$/.test(file) ? false : undefined),
  },
  // Pre-bundle dependencies that stories and tests import, so Vite does not
  // re-optimize and reload in the middle of a browser test run.
  optimizeDeps: {
    include: [
      'react',
      'react-dom',
      'react-dom/client',
      'react-router',
      'react-router/dom',
      '@tanstack/react-query',
      '@faker-js/faker',
      'zod',
      'pino',
      'mermaid',
      'echarts',
      '@excalidraw/excalidraw',
      'swagger-ui-react',
    ],
    exclude: ['@sqlite.org/sqlite-wasm'],
  },
  // SQLite WASM locates its .wasm file relative to its own module; pre-bundling would break that.
  worker: { format: 'es' },
  server: {
    proxy: {
      '/api': goServer,
      '/assets': goServer,
      '/config.js': goServer,
      '/healthz': goServer,
    },
  },
});
