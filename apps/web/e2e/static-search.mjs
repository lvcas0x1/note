// On a static export, the search index downloads only on the first search and
// only once, results appear, and there is no CSP violation. The second run adds
// 'wasm-unsafe-eval' to show that the result does not depend on it: SQLite WASM
// runs in a dedicated worker, which the page's <meta> policy does not govern.
import { readdir, readFile, writeFile, mkdir } from 'node:fs/promises';
import { join } from 'node:path';
import { buildAll, cspViolations, notes, openBrowser, servePages, workdir } from './harness.mjs';

const dir = await workdir();
const content = join(dir, 'content');
await mkdir(join(content, 'kafka'), { recursive: true });
await writeFile(
  join(content, 'sidebar.json'),
  JSON.stringify([
    {
      title: 'Kafka',
      children: [
        { title: 'Deploy', path: 'kafka/deploy.md' },
        { title: 'Memo', path: 'memo.md' },
      ],
    },
  ]),
);
await writeFile(
  join(content, 'kafka', 'deploy.md'),
  '---\ndate: 2026-10-05\n---\n# Kafka deploy\n\nBrokers and partitions.\n',
);
await writeFile(
  join(content, 'memo.md'),
  '---\ndate: 2026-10-11\n---\n# 日本語メモ\n\nこれは日本語のテキストです。\n',
);

const bin = buildAll(dir);
const data = join(dir, 'data');
notes(bin, ['build', '--content', content, '--data', data]);

async function exportSite(name) {
  const site = join(dir, name);
  notes(bin, ['export', '--data', data, '--out', site, '--base-path', '/notes/']);
  return site;
}

async function htmlFiles(root) {
  return (await readdir(root, { recursive: true }))
    .filter((f) => f.endsWith('.html'))
    .map((f) => join(root, f));
}

async function run(site, label) {
  const pages = await servePages(site, '/notes/');
  const { browser, page, record } = await openBrowser();
  try {
    await page.goto(`${pages.url}/notes/a/kafka/deploy/`);
    await page.locator('header h1', { hasText: 'Kafka deploy' }).waitFor();
    await page.waitForTimeout(500);
    const before = record.requests.filter((u) => u.endsWith('/data/search.db')).length;
    const box = page.getByRole('combobox', { name: 'Search notes' });
    await box.fill('kafka');
    let first = [];
    try {
      await page.getByRole('option').first().waitFor({ timeout: 15_000 });
      first = await page.getByRole('option').allTextContents();
    } catch {
      // no results
    }
    await box.fill('本語の');
    let second = [];
    try {
      await page.getByRole('option', { name: /日本語メモ/ }).waitFor({ timeout: 10_000 });
      second = await page.getByRole('option').allTextContents();
    } catch {
      // no results
    }
    const after = record.requests.filter((u) => u.endsWith('/data/search.db')).length;
    const violations = await cspViolations(page, record);
    console.log(
      `${label}: search.db downloads before search=${before} after two searches=${after}; results: ${first.length} then ${second.length}; CSP violations: ${violations.length}`,
    );
    for (const v of violations) console.log(`  CSP: ${v}`);
    return { before, after, first, second, violations };
  } finally {
    await browser.close();
    pages.stop();
  }
}

const normal = await run(await exportSite('site'), "exported CSP (no 'wasm-unsafe-eval')");

const widened = await exportSite('site-wasm');
for (const f of await htmlFiles(widened)) {
  await writeFile(
    f,
    (await readFile(f, 'utf8')).replace(
      "script-src 'self';",
      "script-src 'self' 'wasm-unsafe-eval';",
    ),
  );
}
const withWasm = await run(widened, "with 'wasm-unsafe-eval' added");

const ok =
  normal.before === 0 &&
  normal.after === 1 &&
  normal.first.some((t) => t.includes('Kafka deploy')) &&
  normal.second.some((t) => t.includes('日本語メモ')) &&
  normal.violations.length === 0;
const same =
  withWasm.first.length === normal.first.length && withWasm.second.length === normal.second.length;
console.log(`'wasm-unsafe-eval' changes the result: ${same ? 'no' : 'yes'}`);
if (!ok) process.exit(1);
console.log('OK: lazy single download, results, no CSP violations');
