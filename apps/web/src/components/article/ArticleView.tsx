import type { ArticleDetail } from '@notes/schemas/api';
import { useEffect, useState, type MouseEvent } from 'react';
import { useNavigate } from 'react-router';
import { formatDate } from '../../lib/dates.ts';
import { useEmbeds } from '../embeds/useEmbeds.tsx';
import { createCopyButton } from './copyButton.ts';
import { inAppRoute } from './links.ts';
import { TableOfContents } from './TableOfContents.tsx';

export interface ArticleViewProps {
  article: ArticleDetail;
  basePath?: string;
}

/**
 * Shows an article: title, date, tags, and the compiled HTML body. The HTML
 * was sanitized at build time. Code blocks get a copy button, and links to
 * other articles open in the app without a page reload.
 */
export function ArticleView({ article, basePath = '/' }: ArticleViewProps) {
  const navigate = useNavigate();
  const [body, setBody] = useState<HTMLDivElement | null>(null);
  const html = article.html;
  useEmbeds(body, html);

  // Add a copy button to every code block after React writes the body HTML.
  useEffect(() => {
    const root = body;
    if (!root || html === '') return undefined;
    const buttons: HTMLButtonElement[] = [];
    for (const pre of root.querySelectorAll('pre:not([data-embed])')) {
      const button = createCopyButton(pre);
      pre.append(button);
      buttons.push(button);
    }
    return () => {
      for (const b of buttons) b.remove();
    };
  }, [body, html]);

  const onClick = (e: MouseEvent<HTMLDivElement>) => {
    if (e.defaultPrevented || e.button !== 0 || e.metaKey || e.ctrlKey || e.shiftKey || e.altKey)
      return;
    const anchor = e.target instanceof Element ? e.target.closest('a') : null;
    if (!anchor || anchor.target === '_blank') return;
    const route = inAppRoute(anchor.getAttribute('href') ?? '', article.id, basePath);
    if (!route) return;
    e.preventDefault();
    void navigate(route);
  };

  return (
    <div className="flex items-start gap-8">
      <article className="min-w-0 flex-1">
        <header className="mb-6">
          <h1 className="text-3xl font-semibold leading-tight">{article.title}</h1>
          <div className="mt-2 flex flex-wrap items-center gap-2 text-sm text-muted">
            <time dateTime={article.date}>{formatDate(article.date)}</time>
            {article.category.length > 0 ? <span>· {article.category.join(' / ')}</span> : null}
            {article.tags.map((t) => (
              <span key={t} className="rounded-full bg-accent-soft px-2 py-0.5 text-xs text-ink">
                {t}
              </span>
            ))}
          </div>
        </header>
        {/* oxlint-disable-next-line react/no-danger -- HTML is sanitized by the Go build (bluemonday) */}
        <div
          ref={setBody}
          className="article-body"
          onClick={onClick}
          dangerouslySetInnerHTML={{ __html: html }}
        />
      </article>
      <aside className="sticky top-16 hidden max-h-[calc(100dvh-5rem)] w-60 shrink-0 overflow-y-auto @min-[1024px]:block">
        <TableOfContents html={html} root={body} />
      </aside>
    </div>
  );
}
