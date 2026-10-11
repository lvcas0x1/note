import type { QueryableDB } from './static-search-core.ts';

/**
 * Opens database bytes in memory with the official SQLite WASM build, using
 * the cookbook's recipe (https://sqlite.org/wasm/doc/trunk/cookbook.md):
 * wasm.allocFromTypedArray + capi.sqlite3_deserialize with FREEONCLOSE.
 */
export async function openInMemory(bytes: ArrayBuffer): Promise<QueryableDB> {
  const { default: sqlite3InitModule } = await import('@sqlite.org/sqlite-wasm');
  const sqlite3 = await sqlite3InitModule();
  const p = sqlite3.wasm.allocFromTypedArray(new Uint8Array(bytes));
  const db = new sqlite3.oo1.DB();
  const pointer = db.pointer;
  if (pointer === undefined) throw new Error('SQLite database is not open');
  const rc = sqlite3.capi.sqlite3_deserialize(
    pointer,
    'main',
    p,
    bytes.byteLength,
    bytes.byteLength,
    sqlite3.capi.SQLITE_DESERIALIZE_FREEONCLOSE | sqlite3.capi.SQLITE_DESERIALIZE_READONLY,
  );
  db.checkRc(rc);
  return {
    selectObjects: (sql, bind) => db.selectObjects(sql, bind),
  };
}
