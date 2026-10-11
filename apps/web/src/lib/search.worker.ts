/// <reference lib="webworker" />
// Static-site search: downloads data/search.db once, opens it with SQLite WASM,
// and answers queries off the main thread.
import type { QueryableDB } from './static-search-core.ts';

interface Request {
  id: number;
  basePath: string;
  q: string;
}

declare const self: DedicatedWorkerGlobalScope & { __NOTES_CONFIG__?: { basePath: string } };

let db: Promise<QueryableDB> | undefined;

function load(basePath: string): Promise<QueryableDB> {
  db ??= (async () => {
    // Built asset URLs (the WASM binary) resolve against the site's base path.
    self.__NOTES_CONFIG__ = { basePath };
    const res = await fetch(`${basePath}data/search.db`);
    if (!res.ok) throw new Error(`search index: HTTP ${res.status}`);
    const bytes = await res.arrayBuffer();
    const { openInMemory } = await import('./sqlite.ts');
    return openInMemory(bytes);
  })().catch((err: unknown) => {
    db = undefined;
    throw err;
  });
  return db;
}

self.addEventListener('message', (e: MessageEvent<Request>) => {
  const { id, basePath, q } = e.data;
  void (async () => {
    try {
      const { runSearch } = await import('./static-search-core.ts');
      const hits = runSearch(await load(basePath), q);
      // A dedicated worker's postMessage has no target origin.
      // oxlint-disable-next-line unicorn/require-post-message-target-origin
      self.postMessage({ id, hits });
    } catch (err) {
      const reason =
        err instanceof Error && 'reason' in err && typeof err.reason === 'string'
          ? err.reason
          : undefined;
      // oxlint-disable-next-line unicorn/require-post-message-target-origin
      self.postMessage({ id, error: err instanceof Error ? err.message : String(err), reason });
    }
  })();
});
