// Fails when a workspace imports a package that its own package.json does not declare.
// npm hoists dependencies, so such imports would otherwise work by accident.
import { readdir, readFile } from 'node:fs/promises';
import { builtinModules } from 'node:module';
import { join } from 'node:path';

const root = JSON.parse(await readFile('package.json', 'utf8'));
const skip = new Set(['node_modules', 'dist', 'storybook-static', 'vendor']);
const importRe =
  /(?:^|[\s;])(?:import|export)\s[^'"`]*?from\s*['"]([^'"]+)['"]|import\(\s*['"]([^'"]+)['"]\s*\)|^\s*import\s+['"]([^'"]+)['"]/gm;
const builtins = new Set(builtinModules);

/**
 * @param {string} spec
 * @returns {string | null}
 */
function packageName(spec) {
  if (spec.startsWith('.') || spec.startsWith('/') || spec.startsWith('node:')) return null;
  if (spec.startsWith('virtual:') || spec.startsWith('~') || spec.startsWith('@/')) return null;
  const parts = spec.split('/');
  const name = spec.startsWith('@') ? parts.slice(0, 2).join('/') : (parts[0] ?? spec);
  return builtins.has(name) ? null : name;
}

/**
 * @param {string} dir
 * @returns {AsyncGenerator<string>}
 */
async function* walk(dir) {
  for (const entry of await readdir(dir, { withFileTypes: true })) {
    if (skip.has(entry.name)) continue;
    const path = join(dir, entry.name);
    if (entry.isDirectory()) yield* walk(path);
    else if (/\.(tsx?|mts|jsx?|mjs)$/.test(entry.name) && !entry.name.endsWith('.d.ts')) yield path;
  }
}

/** @returns {Promise<string[]>} */
async function workspaceDirs() {
  /** @type {string[]} */
  const dirs = [];
  for (const pattern of root.workspaces) {
    const base = pattern.replace(/\/\*$/, '');
    try {
      for (const entry of await readdir(base, { withFileTypes: true })) {
        if (entry.isDirectory()) dirs.push(join(base, entry.name));
      }
    } catch {
      // workspace folder does not exist yet
    }
  }
  return dirs;
}

/** @type {string[]} */
const problems = [];
for (const dir of await workspaceDirs()) {
  let pkg;
  try {
    pkg = JSON.parse(await readFile(join(dir, 'package.json'), 'utf8'));
  } catch {
    continue;
  }
  const declared = new Set([
    pkg.name,
    ...Object.keys(pkg.dependencies ?? {}),
    ...Object.keys(pkg.devDependencies ?? {}),
    ...Object.keys(pkg.peerDependencies ?? {}),
    ...Object.keys(pkg.optionalDependencies ?? {}),
  ]);
  for await (const file of walk(dir)) {
    const text = await readFile(file, 'utf8');
    for (const match of text.matchAll(importRe)) {
      const spec = match[1] ?? match[2] ?? match[3];
      const name = spec && packageName(spec);
      if (name && !declared.has(name))
        problems.push(`${file}: imports "${name}", which ${dir}/package.json does not declare`);
    }
  }
}
if (problems.length > 0) {
  console.error(problems.join('\n'));
  process.exit(1);
}
