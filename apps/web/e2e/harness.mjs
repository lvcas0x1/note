// Shared helpers for end-to-end checks: build the Go binary and web app, run
// `notes build/serve/export`, serve static output like GitHub Pages, and drive
// Chromium with Playwright while recording CSP violations and request origins.
import { execFileSync, spawn } from 'node:child_process';
import { createServer } from 'node:http';
import { createServer as createNetServer } from 'node:net';
import { mkdtemp, readFile, stat } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { extname, join, normalize, resolve } from 'node:path';
import { chromium } from 'playwright';

export const repo = resolve(import.meta.dirname, '..', '..', '..');

const children = new Set();
// Never leave a server running when a check fails.
for (const sig of ['exit', 'SIGINT', 'SIGTERM']) {
  process.on(sig, () => {
    for (const c of children) c.kill();
    if (sig !== 'exit') process.exit(1);
  });
}

/** Returns a free local TCP port. */
export async function freePort() {
  const srv = createNetServer();
  await new Promise((r) => srv.listen(0, '127.0.0.1', r));
  const { port } = srv.address();
  await new Promise((r) => srv.close(r));
  return port;
}

export async function workdir() {
  return mkdtemp(join(tmpdir(), 'notes-e2e-'));
}

export function buildAll(dir) {
  execFileSync('npm', ['run', 'build', '-w', '@notes/web'], { cwd: repo, stdio: 'inherit' });
  const bin = join(dir, 'notes');
  execFileSync('go', ['build', '-o', bin, './cmd/notes'], {
    cwd: join(repo, 'server'),
    stdio: 'inherit',
    env: { ...process.env, CGO_ENABLED: '0' },
  });
  return bin;
}

export function notes(bin, args, env = {}) {
  return execFileSync(bin, args, {
    encoding: 'utf8',
    env: { ...process.env, ...env },
    stdio: ['ignore', 'pipe', 'pipe'],
  });
}

export async function serve(bin, dataDir, port) {
  port ??= await freePort();
  const child = spawn(bin, ['serve'], {
    env: { ...process.env, NOTES_DATA: dataDir, NOTES_LISTEN: `127.0.0.1:${port}` },
    stdio: 'ignore',
  });
  children.add(child);
  const url = `http://127.0.0.1:${port}`;
  for (let i = 0; i < 100; i++) {
    try {
      if ((await fetch(`${url}/healthz`)).ok) {
        return {
          url,
          stop: () => {
            child.kill();
            children.delete(child);
          },
        };
      }
    } catch {
      // not up yet
    }
    await new Promise((r) => setTimeout(r, 100));
  }
  child.kill();
  throw new Error('notes serve did not become healthy');
}

const types = {
  '.html': 'text/html; charset=utf-8',
  '.js': 'text/javascript',
  '.css': 'text/css',
  '.json': 'application/json',
  '.wasm': 'application/wasm',
  '.svg': 'image/svg+xml',
  '.png': 'image/png',
  '.woff2': 'font/woff2',
  '.drawio': 'application/xml',
  '.excalidraw': 'application/json',
  '.db': 'application/octet-stream',
  '.mp4': 'video/mp4',
};

/** Serves a static export under prefix like GitHub Pages: folder index.html, else 404.html with status 404. */
export async function servePages(root, prefix, port) {
  port ??= await freePort();
  const server = createServer(async (req, res) => {
    const path = decodeURIComponent(new URL(req.url, 'http://x').pathname);
    if (!path.startsWith(prefix)) {
      res.writeHead(404).end();
      return;
    }
    let file = normalize(join(root, path.slice(prefix.length)));
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
  });
  await new Promise((r) => server.listen(port, '127.0.0.1', r));
  return { url: `http://127.0.0.1:${port}`, stop: () => server.close() };
}

/** Opens a browser page that records CSP violations, console errors, and every request URL. */
export async function openBrowser() {
  const browser = await chromium.launch();
  const page = await browser.newPage();
  const record = { csp: [], requests: [], errors: [] };
  await page.addInitScript(() => {
    window.notesCspViolations = [];
    document.addEventListener('securitypolicyviolation', (e) =>
      window.notesCspViolations.push(
        `${e.violatedDirective} blocked ${e.blockedURI} at ${e.sourceFile}:${e.lineNumber}`,
      ),
    );
  });
  page.on('request', (r) => record.requests.push(r.url()));
  page.on('console', (m) => {
    if (m.type() === 'error' && /Content Security Policy/i.test(m.text()))
      record.csp.push(m.text());
  });
  page.on('pageerror', (e) => record.errors.push(String(e)));
  return { browser, page, record };
}

export async function cspViolations(page, record) {
  const inPage = await page.evaluate(() => window.notesCspViolations ?? []);
  return [...inPage, ...record.csp];
}

export function foreignRequests(record, origin, allowed = []) {
  return record.requests.filter((u) => {
    if (u.startsWith('data:') || u.startsWith('blob:')) return false;
    const o = new URL(u).origin;
    return o !== origin && !allowed.includes(o);
  });
}
