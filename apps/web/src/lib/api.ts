import {
  ArticleDetail,
  CalendarResponse,
  DayResponse,
  Problem,
  SearchResponse,
  SidebarResponse,
} from '@notes/schemas/api';
import type { z } from 'zod';
import { log } from './log.ts';

/** The API answered with an error status (an RFC 9457 problem). */
export class ApiError extends Error {
  readonly status: number;
  readonly problem: Problem | undefined;
  constructor(status: number, problem: Problem | undefined) {
    super(problem?.detail ?? `request failed with status ${status}`);
    this.name = 'ApiError';
    this.status = status;
    this.problem = problem;
  }
}

/** A response did not match its schema; the data is not used. */
export class SchemaError extends Error {
  readonly url: string;
  constructor(url: string, cause: unknown) {
    super(`response from ${url} does not match its schema`, { cause });
    this.name = 'SchemaError';
    this.url = url;
  }
}

/** Fetches url and validates the JSON body with schema. */
export async function fetchJson<T extends z.ZodType>(
  url: string,
  schema: T,
  signal?: AbortSignal,
): Promise<z.output<T>> {
  const res = await fetch(url, { signal: signal ?? null, headers: { Accept: 'application/json' } });
  if (!res.ok) {
    const body: unknown = await res.json().catch(() => undefined);
    const problem = Problem.safeParse(body);
    throw new ApiError(res.status, problem.success ? problem.data : undefined);
  }
  const body: unknown = await res.json();
  const parsed = schema.safeParse(body);
  if (!parsed.success) {
    log.error({ url, issues: parsed.error.issues }, 'response failed schema validation');
    throw new SchemaError(url, parsed.error);
  }
  return parsed.data;
}

/** Encodes an article ID for the path segment; slashes become %2F. */
export function encodeId(id: string): string {
  return encodeURIComponent(id);
}

/** Client for the server version's HTTP API. */
export const api = {
  sidebar: (signal?: AbortSignal) => fetchJson('/api/v1/sidebar', SidebarResponse, signal),
  article: (id: string, signal?: AbortSignal) =>
    fetchJson(`/api/v1/articles/${encodeId(id)}`, ArticleDetail, signal),
  calendar: (from: string, to: string, signal?: AbortSignal) =>
    fetchJson(`/api/v1/calendar?${new URLSearchParams({ from, to })}`, CalendarResponse, signal),
  day: (date: string, signal?: AbortSignal) =>
    fetchJson(`/api/v1/days/${encodeURIComponent(date)}`, DayResponse, signal),
  search: (q: string, signal?: AbortSignal) =>
    fetchJson(`/api/v1/search?${new URLSearchParams({ q })}`, SearchResponse, signal),
};
