// MP4 embed through the Go server, with a real H.264 video made by ffmpeg:
// no autoplay, inline playback, and seeking that continues from the new
// position through an HTTP range request instead of a full download.
import { execFileSync } from 'node:child_process';
import { mkdir, stat, writeFile } from 'node:fs/promises';
import { join } from 'node:path';
import { buildAll, cspViolations, notes, openBrowser, serve, workdir } from './harness.mjs';

const dir = await workdir();
const content = join(dir, 'content');
await mkdir(join(content, 'videos'), { recursive: true });
const video = join(content, 'videos', 'long.mp4');
// 60 s at about 2 Mbit/s (about 15 MB), so a seek to the end needs a range request.
execFileSync('ffmpeg', [
  '-hide_banner',
  '-loglevel',
  'error',
  '-y',
  '-f',
  'lavfi',
  '-i',
  'testsrc2=size=640x360:rate=25:duration=60',
  '-c:v',
  'libx264',
  '-pix_fmt',
  'yuv420p',
  '-b:v',
  '2M',
  '-g',
  '50',
  '-movflags',
  '+faststart',
  video,
]);
const size = (await stat(video)).size;
await writeFile(
  join(content, 'sidebar.json'),
  JSON.stringify([{ title: 'Video', path: 'video.md' }]),
);
await writeFile(
  join(content, 'video.md'),
  '---\ndate: 2026-10-05\n---\n# Video\n\n![mp4](videos/long.mp4)\n',
);

const bin = buildAll(dir);
const data = join(dir, 'data');
notes(bin, ['build', '--content', content, '--data', data]);
const server = await serve(bin, data);
const { browser, page, record } = await openBrowser();
const ranges = [];
page.on('response', (res) => {
  if (!res.url().endsWith('/long.mp4')) return;
  const range = res.request().headers().range ?? '';
  ranges.push(`${res.status()} ${range} ${res.headers()['content-type']}`);
});
// Count the bytes the browser actually received for the video (Content-Length
// of an open-ended range response is the whole file even when reading stops early).
const cdp = await page.context().newCDPSession(page);
await cdp.send('Network.enable');
const videoRequests = new Set();
let bytes = 0;
cdp.on('Network.requestWillBeSent', (e) => {
  if (e.request.url.endsWith('/long.mp4')) videoRequests.add(e.requestId);
});
cdp.on('Network.dataReceived', (e) => {
  if (videoRequests.has(e.requestId)) bytes += e.dataLength;
});
const results = [];
const check = (name, pass) => {
  results.push(pass);
  console.log(`${pass ? 'PASS' : 'FAIL'} ${name}`);
};
try {
  await page.goto(`${server.url}/a/video`);
  const v = page.locator('.embed video');
  await v.waitFor();
  await page.waitForFunction(() => document.querySelector('video')?.readyState >= 1);
  await page.waitForTimeout(500);
  const before = await v.evaluate((el) => ({
    paused: el.paused,
    time: el.currentTime,
    duration: el.duration,
    preload: el.preload,
    autoplay: el.autoplay,
  }));
  check(
    `no autoplay, metadata loaded (duration ${before.duration.toFixed(1)} s)`,
    before.paused &&
      before.time === 0 &&
      Math.round(before.duration) === 60 &&
      before.preload === 'metadata' &&
      !before.autoplay,
  );

  await v.evaluate((el) => {
    el.muted = true;
    return el.play();
  });
  await page.waitForFunction(() => document.querySelector('video').currentTime > 0.5);
  check('plays inline', await v.evaluate((el) => !el.paused && !document.fullscreenElement));

  const bytesBeforeSeek = bytes;
  await v.evaluate((el) => {
    el.currentTime = 50;
  });
  await page.waitForFunction(
    () => {
      const el = document.querySelector('video');
      return el.currentTime > 50.3 && !el.paused;
    },
    null,
    { timeout: 20_000 },
  );
  check('seeking to 50 s continues playback from there', true);
  const midRange = ranges.some((r) => /^206 bytes=([1-9]\d*)-/.test(r));
  check(`seek used a range request starting mid-file (${ranges.join(' | ')})`, midRange);
  check(
    `did not download the whole file first (${bytesBeforeSeek} of ${size} bytes before the seek)`,
    bytesBeforeSeek < size,
  );
  check(
    'served as video/mp4',
    ranges.every((r) => r.endsWith('video/mp4')),
  );
  check('no CSP violations', (await cspViolations(page, record)).length === 0);
} finally {
  await browser.close();
  server.stop();
}
if (!results.every(Boolean)) process.exit(1);
console.log('OK: video playback and seeking');
