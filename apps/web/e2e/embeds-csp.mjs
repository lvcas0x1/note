// An article with every embed type renders with no CSP violation and
// no request to another origin, in the server version (CSP header) and in a
// static export (CSP meta tag; the video origin is the only allowed exception).
import { cp, mkdir, writeFile } from 'node:fs/promises';
import { join } from 'node:path';
import {
  buildAll,
  cspViolations,
  foreignRequests,
  notes,
  openBrowser,
  repo,
  serve,
  servePages,
  workdir,
} from './harness.mjs';

const dir = await workdir();
const content = join(dir, 'content');
await mkdir(join(content, 'demo', 'files'), { recursive: true });
for (const f of ['two-pages.drawio', 'shapes.excalidraw', 'swagger2.json', 'not-a-video.mp4']) {
  await cp(join(repo, 'apps/web/public/fixtures', f), join(content, 'demo', 'files', f));
}
await writeFile(
  join(content, 'sidebar.json'),
  JSON.stringify([{ title: 'Demo', children: [{ title: 'Embeds', path: 'demo/embeds.md' }] }]),
);
await writeFile(
  join(content, 'demo', 'embeds.md'),
  `---
date: 2026-10-05
---
# Every embed

\`\`\`mermaid
flowchart LR; A-->B
\`\`\`

\`\`\`echarts
{"xAxis":{"type":"category","data":["a","b"]},"yAxis":{"type":"value"},"series":[{"type":"bar","data":[1,2]}]}
\`\`\`

![drawio](files/two-pages.drawio)

![excalidraw](files/shapes.excalidraw)

![swagger](files/swagger2.json)

![mp4](files/not-a-video.mp4)

\`\`\`go
fmt.Println("hi")
\`\`\`
`,
);

const bin = buildAll(dir);
const data = join(dir, 'data');
notes(bin, ['build', '--content', content, '--data', data]);

async function check(label, url, origin, allowed = []) {
  const { browser, page, record } = await openBrowser();
  try {
    await page.goto(url, { waitUntil: 'load' });
    await page.locator('header h1', { hasText: 'Every embed' }).waitFor();
    const missing = [];
    for (const selector of [
      '.embed svg',
      '.embed canvas',
      '.embed .mxgraph svg',
      '.embed svg[role="img"]',
      '.embed .swagger-ui .opblock',
      '.embed video',
    ]) {
      try {
        await page.locator(selector).first().waitFor({ timeout: 30_000 });
      } catch {
        missing.push(selector);
        console.log(`  missing: ${selector}`);
      }
    }
    await page.waitForTimeout(1500);
    for (const alert of await page.getByRole('alert').allTextContents())
      console.log(`  embed error: ${alert}`);
    const violations = await cspViolations(page, record);
    const foreign = foreignRequests(record, origin, allowed);
    console.log(
      `${label}: ${record.requests.length} requests, ${violations.length} CSP violations, ${foreign.length} requests to other origins`,
    );
    for (const v of violations) console.log(`  CSP: ${v}`);
    for (const f of foreign) console.log(`  foreign: ${f}`);
    for (const e of record.errors) console.log(`  console: ${e.slice(0, 300)}`);
    return violations.length === 0 && foreign.length === 0 && missing.length === 0;
  } finally {
    await browser.close();
  }
}

const server = await serve(bin, data);
const serverOK = await check('server', `${server.url}/a/demo/embeds`, server.url);
server.stop();

const site = join(dir, 'site');
notes(bin, [
  'export',
  '--data',
  data,
  '--out',
  site,
  '--base-path',
  '/notes/',
  '--video-base-url',
  'https://videos.example.com/notes/',
]);
const pages = await servePages(site, '/notes/');
const staticOK = await check('static', `${pages.url}/notes/a/demo/embeds/`, pages.url, [
  'https://videos.example.com',
]);
pages.stop();

if (!serverOK || !staticOK) process.exit(1);
console.log('OK: no CSP violations and no third-party requests');
