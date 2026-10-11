import {
  MARK_CLOSE,
  MARK_OPEN,
  MAX_RESULTS,
  matchExpression,
  SEARCH_SQL,
  snippetHtml,
} from '@notes/schemas/query';
import type { SearchHit } from '@notes/schemas/api';

/** The part of SQLite WASM's oo1.DB this module uses. */
export interface QueryableDB {
  selectObjects(sql: string, bind: (string | number)[]): Record<string, unknown>[];
}

/** Error for queries the search rejects, matching the server's 400 cases. */
export class QueryRejected extends Error {
  readonly reason: 'too_short' | 'too_long';
  constructor(reason: 'too_short' | 'too_long') {
    super(
      reason === 'too_short'
        ? 'query must have at least 3 characters'
        : 'query must have at most 200 characters',
    );
    this.reason = reason;
  }
}

/** Runs the shared search query on a search.db opened with SQLite WASM. */
export function runSearch(db: QueryableDB, q: string, limit = MAX_RESULTS): SearchHit[] {
  const m = matchExpression(q);
  if ('error' in m) throw new QueryRejected(m.error);
  const rows = db.selectObjects(SEARCH_SQL, [
    MARK_OPEN,
    MARK_CLOSE,
    m.match,
    Math.min(Math.max(limit, 1), MAX_RESULTS),
  ]);
  return rows.map((r) => ({
    id: str(r.id),
    title: str(r.title),
    date: str(r.date),
    summary: str(r.summary),
    category: parseCategory(r.category),
    snippet: snippetHtml(str(r.snippet)),
  }));
}

function str(v: unknown): string {
  return typeof v === 'string' ? v : '';
}

function parseCategory(v: unknown): string[] {
  try {
    const parsed: unknown = JSON.parse(str(v));
    return Array.isArray(parsed) ? parsed.filter((x): x is string => typeof x === 'string') : [];
  } catch {
    return [];
  }
}
