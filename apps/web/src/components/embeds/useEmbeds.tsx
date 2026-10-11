import { lazy, Suspense, useEffect, type ComponentType, type LazyExoticComponent } from 'react';
import { createRoot, type Root } from 'react-dom/client';
import { EmbedBoundary } from './EmbedBoundary.tsx';
import { embedLabels, isEmbedKind, type EmbedKind, type EmbedProps } from './types.ts';

// Each renderer is its own chunk, downloaded only when an article contains that embed type.
const renderers: Record<EmbedKind, LazyExoticComponent<ComponentType<EmbedProps>>> = {
  mermaid: lazy(() => import('./MermaidEmbed.tsx')),
  echarts: lazy(() => import('./EChartsEmbed.tsx')),
  drawio: lazy(() => import('./DrawioEmbed.tsx')),
  excalidraw: lazy(() => import('./ExcalidrawEmbed.tsx')),
  openapi: lazy(() => import('./OpenAPIEmbed.tsx')),
  mp4: lazy(() => import('./VideoEmbed.tsx')),
};

/**
 * Turns the embed placeholders in a rendered article body into renderers.
 * Each embed gets its own React root and error boundary.
 */
export function useEmbeds(body: HTMLElement | null, html: string) {
  useEffect(() => {
    if (!body || html === '') return undefined;
    const roots: Root[] = [];
    const added: HTMLElement[] = [];
    for (const el of body.querySelectorAll<HTMLElement>('[data-embed]')) {
      const kind = el.dataset.embed;
      if (!isEmbedKind(kind)) continue;
      const src = el.dataset.src ?? '';
      const source = el.tagName === 'PRE' ? (el.textContent ?? '') : '';
      let host: HTMLElement = el;
      if (el.tagName === 'PRE') {
        // Keep the source text in the DOM for search engines and copy; show the rendering instead.
        host = document.createElement('div');
        el.hidden = true;
        el.after(host);
        added.push(host);
      }
      host.classList.add('embed');
      host.setAttribute('aria-label', embedLabels[kind]);
      const Renderer = renderers[kind];
      const root = createRoot(host);
      root.render(
        <EmbedBoundary kind={kind} source={src || 'inline source'}>
          <Suspense fallback={<p className="text-sm text-muted">Loading {embedLabels[kind]}…</p>}>
            <Renderer src={src} body={source} />
          </Suspense>
        </EmbedBoundary>,
      );
      roots.push(root);
    }
    return () => {
      // Unmount after the current render to avoid unmounting during React's commit.
      const pending = [...roots];
      queueMicrotask(() => {
        for (const r of pending) r.unmount();
        for (const h of added) h.remove();
      });
    };
  }, [body, html]);
}
