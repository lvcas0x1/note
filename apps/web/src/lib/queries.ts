import { isSearchable } from '@notes/schemas/query';
import { useQuery } from '@tanstack/react-query';
import { ApiError } from './api.ts';
import { useNotesSource } from './source-context.tsx';

// TanStack Query hooks for the note data. Every hook reads through the
// NotesSource (server API or static files). A failed request or a response
// that does not match its schema puts the query in the error state with no data.

// Retry only what can succeed later: network failures and server errors.
// Client errors (4xx) and schema mismatches are final.
const retryTransient = (count: number, error: unknown) =>
  count < 2 && ((error instanceof ApiError && error.status >= 500) || error instanceof TypeError);

export function useSidebar() {
  const source = useNotesSource();
  return useQuery({
    queryKey: ['sidebar'],
    queryFn: ({ signal }) => source.sidebar(signal),
    staleTime: 60_000,
    retry: retryTransient,
  });
}

export function useArticle(id: string) {
  const source = useNotesSource();
  return useQuery({
    queryKey: ['article', id],
    queryFn: ({ signal }) => source.article(id, signal),
    enabled: id !== '',
    retry: retryTransient,
  });
}

export function useCalendar(from: string, to: string) {
  const source = useNotesSource();
  return useQuery({
    queryKey: ['calendar', from, to],
    queryFn: ({ signal }) => source.calendar(from, to, signal),
    retry: retryTransient,
  });
}

export function useDay(date: string) {
  const source = useNotesSource();
  return useQuery({
    queryKey: ['day', date],
    queryFn: ({ signal }) => source.day(date, signal),
    retry: retryTransient,
  });
}

export function useSearch(q: string) {
  const source = useNotesSource();
  return useQuery({
    queryKey: ['search', q.trim()],
    queryFn: ({ signal }) => source.search(q.trim(), signal),
    enabled: isSearchable(q),
    retry: false,
  });
}
