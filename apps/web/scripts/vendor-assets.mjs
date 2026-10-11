// Copies third-party runtime assets that are served from the site's own origin:
// Excalidraw fonts (otherwise fetched from esm.sh). The draw.io viewer is
// committed under public/vendor/drawio.
import { cp, mkdir } from 'node:fs/promises';
import { createRequire } from 'node:module';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const here = dirname(fileURLToPath(import.meta.url));
const require = createRequire(import.meta.url);
// The package does not export package.json; find its folder from the entry point.
const entry = require.resolve('@excalidraw/excalidraw');
const excalidraw = entry.slice(0, entry.lastIndexOf('/dist/'));
const target = join(here, '..', 'public', 'vendor', 'excalidraw', 'fonts');
await mkdir(target, { recursive: true });
await cp(join(excalidraw, 'dist', 'prod', 'fonts'), target, { recursive: true });
