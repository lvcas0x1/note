import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { renderHook, waitFor } from '@testing-library/react';
import type { ReactNode } from 'react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { ApiError, SchemaError } from '../src/lib/api.ts';
import { log } from '../src/lib/log.ts';
import { useArticle } from '../src/lib/queries.ts';
import { NotesSourceProvider } from '../src/lib/source-context.tsx';
import { apiSource, articlePath, staticSource, type NotesSource } from '../src/lib/source.ts';

// One fixture build, served both as API responses and as exported files.
const sidebar = [{ title: 'Cloud', children: [{ title: 'S3', path: 'aws/s3/s3.md' }] }];
const article = {
  id: 'aws/s3/s3',
  title: 'S3',
  date: '2026-10-05',
  summary: 'Buckets',
  category: ['Cloud'],
  tags: ['aws'],
  inSidebar: true,
  html: '<p>x</p>',
};
const days = [
  { date: '2026-09-30', count: 1 },
  { date: '2026-10-05', count: 2 },
];
const day = [
  { id: 'aws/s3/s3', title: 'S3', date: '2026-10-05', summary: 'Buckets', category: ['Cloud'] },
  { id: 'b', title: 'T', date: '2026-10-05', summary: '', category: [] },
];

const routes: Record<string, unknown> = {
  '/api/v1/sidebar': sidebar,
  '/api/v1/articles/aws%2Fs3%2Fs3': article,
  '/api/v1/calendar?from=2026-10-01&to=2026-10-31': [days[1]],
  '/api/v1/days/2026-10-05': day,
  '/notes/data/sidebar.json': sidebar,
  '/notes/data/articles/aws/s3/s3.json': article,
  '/notes/data/calendar.json': days,
  '/notes/data/days/2026-10-05.json': day,
};

function serve(overrides: Record<string, unknown> = {}) {
  const all = { ...routes, ...overrides };
  vi.stubGlobal(
    'fetch',
    vi.fn<typeof fetch>(async (input) => {
      const url = typeof input === 'string' ? input : input instanceof URL ? input.href : input.url;
      if (!(url in all)) return new Response('<html>404</html>', { status: 404 });
      return new Response(JSON.stringify(all[url]), { status: 200 });
    }),
  );
}

const noSearch = async () => [];

describe('static and API sources return identical objects', () => {
  afterEach(() => vi.unstubAllGlobals());

  it.each([
    ['sidebar', (s: NotesSource) => s.sidebar()],
    ['article', (s: NotesSource) => s.article('aws/s3/s3')],
    ['calendar', (s: NotesSource) => s.calendar('2026-10-01', '2026-10-31')],
    ['day', (s: NotesSource) => s.day('2026-10-05')],
  ])('%s', async (_name, read) => {
    serve();
    const fromApi = await read(apiSource);
    const fromStatic = await read(staticSource('/notes/', noSearch));
    expect(fromStatic).toEqual(fromApi);
  });

  it('treats a day without articles as an empty list', async () => {
    serve();
    await expect(staticSource('/notes/', noSearch).day('2026-10-06')).resolves.toEqual([]);
  });

  it('reports a missing static article as a 404', async () => {
    serve();
    await expect(staticSource('/notes/', noSearch).article('nope')).rejects.toMatchObject({
      status: 404,
    });
    await expect(staticSource('/notes/', noSearch).article('nope')).rejects.toBeInstanceOf(
      ApiError,
    );
  });
});

describe('a mismatching static file', () => {
  afterEach(() => vi.unstubAllGlobals());

  it('produces the same error state as an invalid API response', async () => {
    const { title: _t, ...bad } = article;
    serve({ '/notes/data/articles/aws/s3/s3.json': bad, '/api/v1/articles/aws%2Fs3%2Fs3': bad });
    vi.spyOn(log, 'error').mockImplementation(() => {});
    for (const source of [staticSource('/notes/', noSearch), apiSource]) {
      const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
      const wrapper = ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={client}>
          <NotesSourceProvider source={source}>{children}</NotesSourceProvider>
        </QueryClientProvider>
      );
      const { result } = renderHook(() => useArticle('aws/s3/s3'), { wrapper });
      await waitFor(() => expect(result.current.isError).toBe(true));
      expect(result.current.data).toBeUndefined();
      expect(result.current.error).toBeInstanceOf(SchemaError);
    }
  });
});

describe('articlePath', () => {
  it('encodes segments but keeps slashes', () => {
    expect(articlePath('a b/c#d')).toBe('a%20b/c%23d');
  });
});

describe('static calendar', () => {
  afterEach(() => vi.unstubAllGlobals());

  it('is not broken when the first caller aborts', async () => {
    serve();
    const source = staticSource('/notes/', noSearch);
    const first = new AbortController();
    const p1 = source.calendar('2026-10-01', '2026-10-31', first.signal);
    first.abort();
    await p1.catch(() => undefined);
    await expect(source.calendar('2026-10-01', '2026-10-31')).resolves.toEqual([days[1]]);
  });
});
