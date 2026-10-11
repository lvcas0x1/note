// npm start: builds the web app into the Go binary (bin/notes), writes a date
// into articles that have none, builds the content, and serves it on
// NOTES_LISTEN (default 127.0.0.1:8080).
import { execFileSync, spawn } from 'node:child_process';

const env = {
  ...process.env,
  NOTES_DATA: process.env.NOTES_DATA ?? 'data',
  NOTES_CONTENT: process.env.NOTES_CONTENT ?? 'content',
};
execFileSync('npm', ['run', 'build', '-w', '@notes/web'], { stdio: 'inherit' });
execFileSync('go', ['build', '-o', '../bin/notes', './cmd/notes'], {
  cwd: 'server',
  stdio: 'inherit',
  env: { ...process.env, CGO_ENABLED: '0' },
});
execFileSync('bin/notes', ['stamp-dates'], { stdio: 'inherit', env });
execFileSync('bin/notes', ['build'], { stdio: 'inherit', env });
const server = spawn('bin/notes', ['serve'], { stdio: 'inherit', env });
for (const sig of /** @type {const} */ (['SIGINT', 'SIGTERM']))
  process.on(sig, () => server.kill(sig));
server.on('exit', (code) => (process.exitCode = code ?? 0));
