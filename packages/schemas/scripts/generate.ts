// Writes the JSON Schemas for sidebar.json and article front matter from the Zod
// schemas. With --check, fails instead of writing when the files are out of date.
import { readFile, writeFile, mkdir } from 'node:fs/promises';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { z } from 'zod';
import { FrontMatter, Sidebar } from '../src/index.ts';

const repoRoot = resolve(dirname(fileURLToPath(import.meta.url)), '../../..');
// The Go build embeds its own copy, because go:embed cannot read outside the module.
const outDirs = [resolve(repoRoot, 'schemas'), resolve(repoRoot, 'server/internal/content/schemas')];
const targets = {
  'sidebar.schema.json': Sidebar,
  'frontmatter.schema.json': FrontMatter,
};

const check = process.argv.includes('--check');
const stale: string[] = [];
for (const outDir of outDirs) {
  await mkdir(outDir, { recursive: true });
  for (const [file, schema] of Object.entries(targets)) {
    const json = z.toJSONSchema(schema, { target: 'draft-2020-12', io: 'input' });
    const text = `${JSON.stringify(json, null, 2)}\n`;
    const path = resolve(outDir, file);
    if (check) {
      const current = await readFile(path, 'utf8').catch(() => '');
      if (current !== text) stale.push(path.slice(repoRoot.length + 1));
    } else {
      await writeFile(path, text);
    }
  }
}
if (stale.length > 0) {
  console.error(
    `Out of date: ${stale.join(', ')}. Run "npm run generate -w @notes/schemas" and commit the result.`,
  );
  process.exit(1);
}
