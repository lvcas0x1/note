import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { renderHook, waitFor } from '@testing-library/react';
import type { ReactNode } from 'react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { ApiError, SchemaError, encodeId } from '../src/lib/api.ts';
import { log } from '../src/lib/log.ts';
import { useArticle } from '../src/lib/queries.ts';

function wrapper({ children }: { children: ReactNode }) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return <QueryClientProvider client={client}>{children}</QueryClientProvider>;
}

function mockFetch(status: number, body: unknown) {
  const fn = vi.fn<typeof fetch>(async () => new Response(JSON.stringify(body), { status }));
  vi.stubGlobal('fetch', fn);
  return fn;
}

const article = {
  id: 'aws/s3/s3',
  title: 'S3',
  date: '2026-10-05',
  summary: 'Buckets',
  category: ['Cloud'],
  tags: [],
  inSidebar: true,
  html: '<p>x</p>',
};

describe('useArticle', () => {
  afterEach(() => vi.unstubAllGlobals());

  it('returns a valid article and encodes slashes in the ID', async () => {
    const fn = mockFetch(200, article);
    const { result } = renderHook(() => useArticle('aws/s3/s3'), { wrapper });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(article);
    expect(fn).toHaveBeenCalledWith('/api/v1/articles/aws%2Fs3%2Fs3', expect.anything());
  });

  it('turns a response missing title into an error state with no data', async () => {
    const { title: _omit, ...withoutTitle } = article;
    mockFetch(200, withoutTitle);
    const errors = vi.spyOn(log, 'error').mockImplementation(() => {});
    const { result } = renderHook(() => useArticle('aws/s3/s3'), { wrapper });
    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(result.current.data).toBeUndefined();
    expect(result.current.error).toBeInstanceOf(SchemaError);
    expect(errors).toHaveBeenCalled();
  });

  it('turns a problem response into an ApiError', async () => {
    mockFetch(404, { status: 404, title: 'Not Found', detail: 'article x does not exist' });
    const { result } = renderHook(() => useArticle('x'), { wrapper });
    await waitFor(() => expect(result.current.isError).toBe(true));
    const err = result.current.error as ApiError;
    expect(err).toBeInstanceOf(ApiError);
    expect(err.status).toBe(404);
    expect(err.message).toContain('does not exist');
  });
});

describe('encodeId', () => {
  it('escapes reserved characters', () => {
    expect(encodeId('a b/c?d')).toBe('a%20b%2Fc%3Fd');
  });
});
