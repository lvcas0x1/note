// Fails when any source file uses a Tailwind `dark:` variant (light theme only).
import { readdir, readFile } from 'node:fs/promises';
import { join } from 'node:path';

const roots = ['apps', 'packages'];
const exts = /\.(tsx?|jsx?|css|html|mdx)$/;
const skip = new Set(['node_modules', 'dist', 'storybook-static', 'client', 'vendor']);
const pattern = /(^|[\s"'`{(])dark:[\w[-]/;

/**
 * @param {string} dir
 * @returns {AsyncGenerator<string>}
 */
async function* walk(dir) {
  let entries;
  try {
    entries = await readdir(dir, { withFileTypes: true });
  } catch {
    return;
  }
  for (const entry of entries) {
    if (skip.has(entry.name)) continue;
    const path = join(dir, entry.name);
    if (entry.isDirectory()) yield* walk(path);
    else if (exts.test(entry.name)) yield path;
  }
}

/** @type {string[]} */
const hits = [];
for (const root of roots) {
  for await (const file of walk(root)) {
    const lines = (await readFile(file, 'utf8')).split('\n');
    lines.forEach((line, i) => {
      if (pattern.test(line)) hits.push(`${file}:${i + 1}: ${line.trim()}`);
    });
  }
}
if (hits.length > 0) {
  console.error('dark: variants are not allowed (light theme only):\n' + hits.join('\n'));
  process.exit(1);
}
