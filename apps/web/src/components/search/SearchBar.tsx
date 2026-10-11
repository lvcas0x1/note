import { codePoints, isSearchable, MAX_QUERY_CHARS } from '@notes/schemas/query';
import { useEffect, useId, useRef, useState, type KeyboardEvent } from 'react';
import { useNavigate } from 'react-router';
import { ApiError } from '../../lib/api.ts';
import { formatDate } from '../../lib/dates.ts';
import { useSearch } from '../../lib/queries.ts';
import { useDebounced } from '../../lib/useDebounced.ts';

/** Debounce before a typed query is sent (search spec: 200 ms after typing pauses). */
export const SEARCH_DEBOUNCE_MS = 200;

function truncate(s: string, max: number): string {
  if (codePoints(s) <= max) return s;
  return Array.from(s).slice(0, max).join('');
}

function isTextField(el: Element | null): boolean {
  return (
    el instanceof HTMLElement &&
    (el.isContentEditable || /^(INPUT|TEXTAREA|SELECT)$/.test(el.tagName))
  );
}

/**
 * Search box with a results dropdown (WAI-ARIA combobox pattern). `/` or
 * Ctrl/Cmd+K focuses it. Results show title, date, category, and a snippet
 * with matches highlighted.
 */
export function SearchBar() {
  const navigate = useNavigate();
  const input = useRef<HTMLInputElement>(null);
  const listId = useId();
  const [text, setText] = useState('');
  const [open, setOpen] = useState(false);
  const [active, setActive] = useState(-1);
  const query = useDebounced(text, SEARCH_DEBOUNCE_MS);
  const search = useSearch(query);
  const hits = search.data ?? [];
  const searchable = isSearchable(text);

  // Keyboard shortcuts: "/" outside text fields, Ctrl+K or Cmd+K anywhere.
  useEffect(() => {
    const onKey = (e: globalThis.KeyboardEvent) => {
      const slash =
        e.key === '/' &&
        !e.metaKey &&
        !e.ctrlKey &&
        !e.altKey &&
        !isTextField(document.activeElement);
      const cmdK = e.key.toLowerCase() === 'k' && (e.metaKey || e.ctrlKey);
      if (slash || cmdK) {
        e.preventDefault();
        input.current?.focus();
        input.current?.select();
      }
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, []);

  const go = (id: string) => {
    setOpen(false);
    setActive(-1);
    input.current?.blur();
    void navigate(`/a/${id}`);
  };

  const onKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Escape') {
      setOpen(false);
      setActive(-1);
      return;
    }
    if (!open && (e.key === 'ArrowDown' || e.key === 'ArrowUp')) setOpen(true);
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setActive((i) => Math.min(i + 1, hits.length - 1));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setActive((i) => Math.max(i - 1, -1));
    } else if (e.key === 'Enter') {
      const hit = hits[active >= 0 ? active : 0];
      if (hit && open) {
        e.preventDefault();
        go(hit.id);
      }
    }
  };

  const rateLimited = search.error instanceof ApiError && search.error.status === 429;
  const showPanel = open && text.trim() !== '';
  const settling = query !== text || search.isFetching;

  return (
    <div className="relative w-full">
      <input
        ref={input}
        type="search"
        role="combobox"
        aria-label="Search notes"
        aria-expanded={showPanel}
        aria-controls={listId}
        aria-autocomplete="list"
        aria-activedescendant={showPanel && active >= 0 ? `${listId}-${active}` : undefined}
        placeholder="Search  (press / )"
        maxLength={MAX_QUERY_CHARS * 2}
        value={text}
        onChange={(e) => {
          setText(truncate(e.target.value, MAX_QUERY_CHARS));
          setOpen(true);
          setActive(-1);
        }}
        onFocus={() => setOpen(true)}
        onBlur={() => setOpen(false)}
        onKeyDown={onKeyDown}
        className="w-full rounded-lg border border-line bg-page px-3 py-1.5 text-sm shadow-sm outline-none focus:border-accent focus:ring-2 focus:ring-accent-soft"
      />
      {showPanel ? (
        <div className="absolute inset-x-0 top-full z-30 mt-1 max-h-[70vh] overflow-y-auto rounded-lg border border-line bg-page shadow-lg">
          {!searchable ? (
            <p className="p-3 text-sm text-muted">Type at least 3 characters to search.</p>
          ) : rateLimited ? (
            <p role="status" className="p-3 text-sm text-muted">
              Too many searches. Please wait a moment and try again.
            </p>
          ) : search.isError ? (
            <p role="status" className="p-3 text-sm text-muted">
              Search is not available right now.
            </p>
          ) : hits.length === 0 && !settling ? (
            <p role="status" className="p-3 text-sm text-muted">
              No articles match “{text.trim()}”.
            </p>
          ) : (
            <ul id={listId} role="listbox" aria-label="Search results" aria-busy={settling}>
              {hits.map((h, i) => (
                <li
                  key={h.id}
                  id={`${listId}-${i}`}
                  role="option"
                  aria-selected={i === active}
                  // Prevent the input's blur from closing the list before the click lands.
                  onMouseDown={(e) => e.preventDefault()}
                  onClick={() => go(h.id)}
                  className={`cursor-pointer border-b border-line px-3 py-2 last:border-b-0 ${i === active ? 'bg-accent-soft' : 'hover:bg-surface'}`}
                >
                  <div className="flex items-baseline justify-between gap-3">
                    <span className="font-medium text-ink">{h.title}</span>
                    <time className="shrink-0 text-xs text-muted" dateTime={h.date}>
                      {formatDate(h.date)}
                    </time>
                  </div>
                  {h.category.length > 0 ? (
                    <div className="text-xs text-muted">{h.category.join(' / ')}</div>
                  ) : null}
                  <div
                    className="mt-1 line-clamp-2 text-sm text-muted [&_mark]:rounded [&_mark]:bg-mark [&_mark]:px-0.5 [&_mark]:text-ink"
                    // Snippets are HTML-escaped by the search, with <mark> around matches.
                    dangerouslySetInnerHTML={{ __html: h.snippet }}
                  />
                </li>
              ))}
            </ul>
          )}
        </div>
      ) : null}
    </div>
  );
}
