import { useEffect, useMemo, useState } from 'react';

export interface TocEntry {
  id: string;
  text: string;
  level: number;
}

/** Lists the h1–h6 headings of compiled article HTML, in order. */
export function headingsOf(html: string): TocEntry[] {
  if (typeof DOMParser === 'undefined' || html === '') return [];
  const doc = new DOMParser().parseFromString(html, 'text/html');
  return [...doc.querySelectorAll<HTMLElement>('h1, h2, h3, h4, h5, h6')]
    .filter((h) => h.id !== '')
    .map((h) => ({
      id: h.id,
      text: (h.textContent ?? '').trim(),
      level: Number(h.tagName.slice(1)),
    }));
}

/**
 * Table of contents for the article body (h1–h6). ArticleView keeps it fixed (sticky)
 * while the article scrolls; it highlights the section being read.
 */
export function TableOfContents({ html, root }: { html: string; root: HTMLElement | null }) {
  const entries = useMemo(() => headingsOf(html), [html]);
  const [active, setActive] = useState<string | null>(null);

  // Highlight the last heading that has scrolled past the top of the content area.
  useEffect(() => {
    if (!root || entries.length === 0) return undefined;
    const headings = entries
      .map((e) => root.querySelector<HTMLElement>(`#${CSS.escape(e.id)}`))
      .filter((h): h is HTMLElement => h !== null);
    const scroller = root.closest('main');
    const update = () => {
      const top = (scroller?.getBoundingClientRect().top ?? 0) + 80;
      let current: string | null = headings[0]?.id ?? null;
      for (const h of headings) {
        if (h.getBoundingClientRect().top <= top) current = h.id;
      }
      setActive(current);
    };
    const target = scroller ?? window;
    target.addEventListener('scroll', update, { passive: true });
    const frame = requestAnimationFrame(update);
    return () => {
      cancelAnimationFrame(frame);
      target.removeEventListener('scroll', update);
    };
  }, [root, entries]);

  if (entries.length === 0) return null;
  const minLevel = Math.min(...entries.map((e) => e.level));

  return (
    <nav aria-label="Table of contents" className="text-sm">
      <p className="mb-2 font-semibold text-ink">On this page</p>
      <ul className="space-y-1 border-l border-line">
        {entries.map((e) => (
          <li key={e.id}>
            <a
              href={`#${e.id}`}
              aria-current={e.id === active ? 'location' : undefined}
              onClick={(ev) => {
                ev.preventDefault();
                root
                  ?.querySelector<HTMLElement>(`#${CSS.escape(e.id)}`)
                  ?.scrollIntoView({ block: 'start' });
                history.replaceState(history.state, '', `#${e.id}`);
                setActive(e.id);
              }}
              className={`-ml-px block border-l-2 py-0.5 pr-2 ${e.id === active ? 'border-accent font-medium text-accent' : 'border-transparent text-muted hover:text-ink'}`}
              style={{ paddingLeft: `${(e.level - minLevel) * 0.75 + 0.75}rem` }}
            >
              {e.text}
            </a>
          </li>
        ))}
      </ul>
    </nav>
  );
}
