// Search query rules shared by the server and the static site
// (see server/internal/search and fixtures/search/cases.json).

export const MIN_QUERY_CHARS = 3;
export const MAX_QUERY_CHARS = 200;
export const MAX_RESULTS = 20;

/** Counts Unicode code points, as the Go server counts runes. */
export function codePoints(s: string): number {
  let n = 0;
  for (const _ of s) n++;
  return n;
}

/** Reports whether the trimmed query is long enough to search. */
export function isSearchable(q: string): boolean {
  return codePoints(q.trim()) >= MIN_QUERY_CHARS;
}

export type QueryError = 'too_short' | 'too_long';

/**
 * Turns user input into an FTS5 MATCH expression, exactly as the Go server
 * does (server/internal/search.MatchExpression): every term is quoted so FTS5
 * operators are literal, and terms shorter than three characters are dropped.
 */
export function matchExpression(q: string): { match: string } | { error: QueryError } {
  if (codePoints(q) > MAX_QUERY_CHARS) return { error: 'too_long' };
  const trimmed = q.trim();
  if (codePoints(trimmed) < MIN_QUERY_CHARS) return { error: 'too_short' };
  const terms = trimmed
    .split(/\s+/u)
    .filter((t) => codePoints(t) >= MIN_QUERY_CHARS)
    .map((t) => `"${t.replaceAll('"', '""')}"`);
  return terms.length === 0 ? { error: 'too_short' } : { match: terms.join(' ') };
}

/** Snippet markers inserted by FTS5 snippet(); replaced by <mark> after escaping. */
export const MARK_OPEN = '\u0002';
export const MARK_CLOSE = '\u0003';

const escapes: Record<string, string> = { '&': '&amp;', "'": '&#39;', '<': '&lt;', '>': '&gt;', '"': '&#34;' };

/** Escapes a raw snippet like Go's html.EscapeString and turns markers into <mark>. */
export function snippetHtml(raw: string): string {
  return raw
    .replace(/[&'<>"]/g, (c) => escapes[c] ?? c)
    .replaceAll(MARK_OPEN, '<mark>')
    .replaceAll(MARK_CLOSE, '</mark>');
}

/** The search query, shared with the Go server (bm25 weights: id, title, summary, tags, body). */
export const SEARCH_SQL = `SELECT a.id AS id, a.title AS title, a.date AS date, a.summary AS summary, a.category AS category,
       snippet(articles_fts, -1, ?, ?, '…', 16) AS snippet
FROM articles_fts JOIN articles a ON a.id = articles_fts.id
WHERE articles_fts MATCH ?
ORDER BY bm25(articles_fts, 0.0, 10.0, 3.0, 3.0, 1.0), a.id
LIMIT ?`;
