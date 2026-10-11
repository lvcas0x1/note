// Writes openapi/openapi.json from the Go API code. With --check, fails when the
// committed file differs from what the code produces.
import { execFileSync } from 'node:child_process';
import { readFile, writeFile } from 'node:fs/promises';

const target = 'openapi/openapi.json';
const generated = execFileSync('go', ['run', './cmd/notes', 'openapi'], {
  cwd: 'server',
  encoding: 'utf8',
});
if (process.argv.includes('--check')) {
  const committed = await readFile(target, 'utf8').catch(() => '');
  if (committed !== generated) {
    console.error(
      `${target} is out of date. Run "npm run generate:openapi" and commit the result.`,
    );
    process.exit(1);
  }
} else {
  await writeFile(target, generated);
}
