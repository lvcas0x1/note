// Runs after `vite build`. Vite copies everything in public/, which includes
// test fixtures for the stories; they must not ship, so remove them. Vite also
// empties the output folder, so restore the placeholder that keeps the
// embedded folder in version control.
import { rmSync, writeFileSync } from 'node:fs';

const out = new URL('../../../server/internal/webui/static/', import.meta.url);
rmSync(new URL('fixtures', out), { recursive: true, force: true });
writeFileSync(new URL('.gitkeep', out), '');
