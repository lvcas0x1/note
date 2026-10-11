// End-to-end check on a temporary copy of ~/Documents/github/docsify.
// The original folder is only read. Steps: migrate the sidebar, mark every
// article as a draft and confirm nothing is published, un-draft a sample that
// covers Mermaid, Swagger, draw.io, and a local MP4, then check the server
// version and a static export under /notes/ in Chromium.
import { cp, mkdir, readFile, readdir, writeFile } from 'node:fs/promises';
import { homedir } from 'node:os';
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

const source = join(homedir(), 'Documents', 'github', 'docsify');
const dir = await workdir();
const content = join(dir, 'docsify');
await cp(source, content, { recursive: true });
const bin = buildAll(dir);
const data = join(dir, 'data');
const report = [];
const log = (line) => {
  console.log(line);
  report.push(line);
};

// 1. Migrate the sidebar.
notes(bin, ['migrate-sidebar', join(content, '_sidebar.md')]);

// 2. Mark every article as a draft; the build must publish nothing.
const mdFiles = (await readdir(content, { recursive: true })).filter((f) => f.endsWith('.md'));
for (const f of mdFiles) {
  const p = join(content, f);
  await writeFile(p, `---\ndraft: true\n---\n${await readFile(p, 'utf8')}`);
}
const allDrafts = notes(bin, ['build', '--content', content, '--data', data]);
log(`all drafts: ${allDrafts.trim()}`);
if (!/built 0 articles and 0 assets/.test(allDrafts)) throw new Error('drafts were published');

// 3. Un-draft a sample, converting docsify-only syntax to this site's embeds.
const sample = {
  'auth/jwt.md': (s) => s, // Mermaid fences
  'swagger/swagger.md': (s) =>
    s.replace('[swagger](/swagger/swagger.json)', '![swagger](/swagger/swagger.json)'),
  '_index.md': (s) =>
    s.replace(
      "[filename](/_plugin/vendor/drawio/test.drawio ':include :type=code')",
      '![drawio](/_plugin/vendor/drawio/test.drawio)',
    ),
  'golang/basic/basic.md': (s) => `${s}\n\n![mp4](/videos/demo.mp4)\n`,
};
await mkdir(join(content, 'videos'), { recursive: true });
await cp(
  join(repo, 'apps/web/public/fixtures/not-a-video.mp4'),
  join(content, 'videos', 'demo.mp4'),
);
for (const [f, edit] of Object.entries(sample)) {
  const p = join(content, f);
  const body = (await readFile(p, 'utf8')).replace(
    '---\ndraft: true\n---\n',
    '---\ndate: 2026-10-05\n---\n',
  );
  await writeFile(p, edit(body));
}
const sampleBuild = notes(bin, ['build', '--content', content, '--data', data]);
log(`sample: ${sampleBuild.trim()}`);

async function checkSite(label, base, origin, allowed, sidebarPath) {
  const { browser, page, record } = await openBrowser();
  const ok = [];
  const check = (name, pass) => {
    ok.push(pass);
    log(`  ${pass ? 'PASS' : 'FAIL'} ${label}: ${name}`);
  };
  try {
    // Sidebar: open the JWT article from its category.
    await page.goto(`${base}/`);
    await page.getByRole('treeitem', { name: 'Auth' }).click();
    await page.getByRole('link', { name: 'JWT' }).click();
    await page.locator('header h1').waitFor();
    check('sidebar opens an article', page.url().includes(`${sidebarPath}a/auth/jwt`));
    await page.locator('.embed svg').first().waitFor({ timeout: 20_000 });
    check('Mermaid renders', (await page.locator('.embed svg').count()) > 0);

    // Calendar: October 5 lists the sample articles; open one.
    const oct5 = page.getByRole('gridcell', { name: /^October 5, 2026/ });
    await oct5.waitFor({ timeout: 10_000 });
    log(`  ${label}: calendar cell label "${await oct5.getAttribute('aria-label')}"`);
    await page.getByRole('gridcell', { name: /October 5, 2026, 4 articles/ }).click();
    await page.getByRole('heading', { name: 'October 5, 2026', level: 1 }).waitFor();
    const titles = await page.locator('section ul li a').allTextContents();
    check(`calendar day lists 4 articles (${titles.join(', ')})`, titles.length === 4);
    await page.locator('section').getByRole('link', { name: 'Swagger' }).click();
    await page.locator('.embed .swagger-ui .opblock').first().waitFor({ timeout: 30_000 });
    check('Swagger 2.0 renders from the calendar', true);

    // Search: find the draw.io article and open it.
    const box = page.getByRole('combobox', { name: 'Search notes' });
    await box.fill('DrawIO');
    await page.getByRole('option').first().waitFor({ timeout: 20_000 });
    await page.getByRole('option').first().click();
    await page.locator('.embed .mxgraph svg').first().waitFor({ timeout: 30_000 });
    check('search opens the draw.io article and it renders', true);

    // Direct load (and reload) of the MP4 article.
    await page.goto(`${base}/a/golang/basic/basic${label === 'static' ? '/' : ''}`);
    await page.reload();
    await page.locator('.embed video').waitFor({ timeout: 20_000 });
    check('direct load + reload shows the video article', true);

    // A draft stays hidden.
    await page.goto(`${base}/a/kafka/deploy/deploy${label === 'static' ? '/' : ''}`);
    await page
      .getByText(/Article not found|Page not found/)
      .first()
      .waitFor({ timeout: 10_000 });
    check('a draft is not reachable', true);

    const violations = await cspViolations(page, record);
    const foreign = foreignRequests(record, origin, allowed);
    check(`no CSP violations (${violations.join('; ')})`, violations.length === 0);
    check(`no third-party requests (${foreign.join(', ')})`, foreign.length === 0);
  } catch (err) {
    check(`unexpected error: ${String(err).split('\n')[0]}`, false);
  } finally {
    await browser.close();
  }
  return ok.every(Boolean);
}

const server = await serve(bin, data);
const serverOK = await checkSite('server', server.url, server.url, [], '/');
server.stop();

const site = join(dir, 'site');
log(
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
  ]).trim(),
);
const pages = await servePages(site, '/notes/');
const staticOK = await checkSite(
  'static',
  `${pages.url}/notes`,
  pages.url,
  ['https://videos.example.com'],
  '/notes/',
);
pages.stop();

const reportFile = join(dir, 'report.txt');
await writeFile(reportFile, `${report.join('\n')}\n`);
console.log(`report: ${reportFile}`);
if (!serverOK || !staticOK) process.exit(1);
console.log('OK: docsify end-to-end check passed');
