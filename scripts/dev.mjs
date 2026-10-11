// npm run dev: writes a date into articles that have none, builds the content,
// then runs the Go server (API, assets, config.js on 127.0.0.1:8080) and the
// Vite dev server, which proxies to it. Ctrl+C stops both. After editing
// content, run `npm run build:content`; the server picks up the new build.
import { spawn, spawnSync } from 'node:child_process';
import { existsSync } from 'node:fs';

const env = {
  ...process.env,
  NOTES_DATA: process.env.NOTES_DATA ?? '../data',
  NOTES_CONTENT: process.env.NOTES_CONTENT ?? '../content',
};
for (const step of ['stamp-dates', 'build']) {
  const r = spawnSync('go', ['run', './cmd/notes', step], { cwd: 'server', env, stdio: 'inherit' });
  if (r.status !== 0 && existsSync('data/current')) {
    console.error(`notes ${step} failed; serving the last good build.`);
    break;
  }
  if (r.status !== 0) process.exit(r.status ?? 1);
}
const children = [
  spawn('go', ['run', './cmd/notes', 'serve'], { cwd: 'server', env, stdio: 'inherit' }),
  spawn('npm', ['run', 'dev', '-w', '@notes/web'], { stdio: 'inherit' }),
];
const stop = () => {
  for (const c of children) c.kill('SIGTERM');
};
process.on('SIGINT', stop);
process.on('SIGTERM', stop);
for (const c of children)
  c.on('exit', (code) => {
    stop();
    process.exitCode = code ?? 0;
  });
