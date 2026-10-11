import type { ApiNode, ArticleDetail, SearchHit } from '@notes/schemas/api';
import { ApiError } from '../lib/api.ts';
import type { NotesSource } from '../lib/source.ts';

const summary = ({ id, title, date, summary: s, category }: ArticleDetail) => ({
  id,
  title,
  date,
  summary: s,
  category,
});

/** An in-memory NotesSource for stories and tests. */
export function fakeSource(
  sidebar: ApiNode[],
  articles: ArticleDetail[],
  hits: (q: string) => SearchHit[] = () => [],
): NotesSource {
  const byId = new Map(articles.map((a) => [a.id, a]));
  return {
    sidebar: async () => sidebar,
    article: async (id) => {
      const a = byId.get(id);
      if (!a)
        throw new ApiError(404, {
          status: 404,
          title: 'Not Found',
          detail: `article ${id} does not exist`,
        });
      return a;
    },
    calendar: async (from, to) => {
      const counts = new Map<string, number>();
      for (const a of articles)
        if (a.date >= from && a.date <= to) counts.set(a.date, (counts.get(a.date) ?? 0) + 1);
      return [...counts]
        .toSorted(([a], [b]) => a.localeCompare(b))
        .map(([date, count]) => ({ date, count }));
    },
    day: async (date) =>
      articles
        .filter((a) => a.date === date)
        .map(summary)
        .toSorted((a, b) => a.title.localeCompare(b.title, undefined, { sensitivity: 'base' })),
    search: async (q) => hits(q),
  };
}
