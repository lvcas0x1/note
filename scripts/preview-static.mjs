// npm run preview:static: serves ./static like GitHub Pages so the exported
// site can be checked locally. A folder serves its index.html; a missing path
// serves 404.html with status 404. Opening static/index.html as a file does
// not work: the site loads scripts, data, and the search worker by URL.
//
// Settings: NOTES_BASE_PATH (default "/", must match the export),
//           PORT (default 4173). Listens on 127.0.0.1 only.
import { createServer } from 'node:http';
import { readFile, stat } from 'node:fs/promises';
import { extname, join, normalize, resolve } from 'node:path';

const root = resolve('static');
const base = (process.env.NOTES_BASE_PATH ?? '/').replace(/\/?$/, '/');
const port = Number(process.env.PORT ?? 4173);
/** @type {Record<string, string>} */
const types = {
  '.html': 'text/html; charset=utf-8',
  '.js': 'text/javascript; charset=utf-8',
  '.css': 'text/css; charset=utf-8',
  '.json': 'application/json',
  '.wasm': 'application/wasm',
  '.svg': 'image/svg+xml',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.gif': 'image/gif',
  '.webp': 'image/webp',
  '.woff2': 'font/woff2',
  '.drawio': 'application/xml',
  '.excalidraw': 'application/json',
  '.yaml': 'application/yaml',
  '.yml': 'application/yaml',
  '.db': 'application/octet-stream',
};

try {
  await stat(join(root, 'index.html'));
} catch {
  console.error('static/index.html not found. Run `npm run export` first.');
  process.exit(1);
}

createServer(async (req, res) => {
  const path = decodeURIComponent(new URL(req.url ?? '/', 'http://localhost').pathname);
  if (!path.startsWith(base)) {
    res.writeHead(302, { Location: base }).end();
    return;
  }
  let file = normalize(join(root, path.slice(base.length)));
  if (!file.startsWith(root)) {
    res.writeHead(403).end();
    return;
  }
  try {
    if ((await stat(file)).isDirectory()) file = join(file, 'index.html');
    const body = await readFile(file);
    res
      .writeHead(200, { 'Content-Type': types[extname(file)] ?? 'application/octet-stream' })
      .end(body);
  } catch {
    res
      .writeHead(404, { 'Content-Type': types['.html'] })
      .end(await readFile(join(root, '404.html')));
  }
}).listen(port, '127.0.0.1', () => {
  console.log(`Previewing static/ at http://127.0.0.1:${port}${base}  (Ctrl+C to stop)`);
});
