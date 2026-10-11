import {
  ArticleDetail,
  CalendarResponse,
  DayResponse,
  SidebarResponse,
  type SearchHit,
} from '@notes/schemas/api';
import { api, fetchJson } from './api.ts';
import type { NotesConfig } from './config.ts';

/** Where the web app reads note data from: the server API or static files. */
export interface NotesSource {
  sidebar(signal?: AbortSignal): Promise<SidebarResponse>;
  article(id: string, signal?: AbortSignal): Promise<ArticleDetail>;
  /** Days in [from, to] that have articles. */
  calendar(from: string, to: string, signal?: AbortSignal): Promise<CalendarResponse>;
  day(date: string, signal?: AbortSignal): Promise<DayResponse>;
  search(q: string, signal?: AbortSignal): Promise<SearchHit[]>;
}

/** Reads from the Go server's API. */
export const apiSource: NotesSource = api;

/** Searches a static site's index in the browser (see static-search.ts). */
export type StaticSearcher = (q: string, signal?: AbortSignal) => Promise<SearchHit[]>;

/**
 * Reads the files written by `notes export` under basePath. Responses are
 * validated with the same schemas as API responses.
 */
export function staticSource(basePath: string, searcher: StaticSearcher): NotesSource {
  const data = `${basePath.replace(/\/?$/, '/')}data/`;
  // One calendar.json download serves every month. It is shared between
  // queries, so it must not use any single caller's abort signal.
  let calendarCache: Promise<CalendarResponse> | undefined;
  const allDays = () => {
    calendarCache ??= fetchJson(`${data}calendar.json`, CalendarResponse).catch((err: unknown) => {
      calendarCache = undefined;
      throw err;
    });
    return calendarCache;
  };
  return {
    sidebar: (signal) => fetchJson(`${data}sidebar.json`, SidebarResponse, signal),
    article: (id, signal) =>
      fetchJson(`${data}articles/${articlePath(id)}.json`, ArticleDetail, signal),
    calendar: async (from, to) => (await allDays()).filter((d) => d.date >= from && d.date <= to),
    day: async (date, signal) => {
      // Days without articles have no file; they are an empty list, not an error.
      const days = await allDays();
      if (!days.some((d) => d.date === date)) return [];
      return fetchJson(`${data}days/${encodeURIComponent(date)}.json`, DayResponse, signal);
    },
    search: searcher,
  };
}

/** Encodes each segment of an article ID for a static file path. */
export function articlePath(id: string): string {
  return id.split('/').map(encodeURIComponent).join('/');
}

/** Picks the source for the runtime configuration. */
export function sourceFor(config: NotesConfig, searcher: StaticSearcher): NotesSource {
  return config.mode === 'static' ? staticSource(config.basePath, searcher) : apiSource;
}
