// @vitest-environment node
// Runs the shared search cases (fixtures/search/cases.json) against the
// search-only database with SQLite WASM; the Go search passes the same cases.
import { readFile } from 'node:fs/promises';
import { beforeAll, describe, expect, it } from 'vitest';
import { openInMemory } from '../src/lib/sqlite.ts';
import { QueryRejected, runSearch, type QueryableDB } from '../src/lib/static-search-core.ts';

interface Case {
  name: string;
  query: string;
  ids?: string[];
  error?: 'too_short' | 'too_long';
  snippets?: Record<string, string>;
}

const dir = new URL('../../../fixtures/search/', import.meta.url);
const { cases } = JSON.parse(await readFile(new URL('cases.json', dir), 'utf8')) as {
  cases: Case[];
};
let db: QueryableDB;

beforeAll(async () => {
  const bytes = await readFile(new URL('fixture.search.db', dir));
  db = await openInMemory(
    bytes.buffer.slice(bytes.byteOffset, bytes.byteOffset + bytes.byteLength),
  );
});

const errorCases = cases.filter((c) => c.error).map((c) => [c.name, c] as const);
const resultCases = cases.filter((c) => !c.error).map((c) => [c.name, c] as const);

function rejection(q: string): string | undefined {
  try {
    runSearch(db, q);
    return undefined;
  } catch (err) {
    return err instanceof QueryRejected ? err.reason : 'other error';
  }
}

describe('shared search cases on SQLite WASM', () => {
  it.each(errorCases)('%s', (_name, c) => {
    expect(rejection(c.query)).toBe(c.error);
  });

  it.each(resultCases)('%s', (_name, c) => {
    const hits = runSearch(db, c.query);
    expect(hits.map((h) => h.id)).toEqual(c.ids ?? []);
    const snippets = Object.fromEntries(hits.map((h) => [h.id, h.snippet]));
    expect(snippets).toMatchObject(c.snippets ?? {});
  });
});
