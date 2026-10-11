import { use, useEffect, useRef } from 'react';
import { readConfig } from '../../lib/config.ts';
import { fetchText, promiseCache } from './cache.ts';
import { latte } from '../../lib/palette.ts';
import type { EmbedProps } from './types.ts';

interface GraphViewerApi {
  createViewerForElement: (el: HTMLElement, callback?: () => void) => void;
}

declare global {
  interface Window {
    GraphViewer?: GraphViewerApi;
    DRAWIO_BASE_URL?: string;
    DRAWIO_LIGHTBOX_URL?: string;
  }
}

let viewer: Promise<GraphViewerApi> | undefined;

/** Loads the vendored draw.io viewer from the site's own origin, once. */
function loadViewer(): Promise<GraphViewerApi> {
  viewer ??= new Promise<GraphViewerApi>((resolve, reject) => {
    const base = `${location.origin}${readConfig().basePath}vendor/drawio`;
    window.DRAWIO_BASE_URL = base;
    window.DRAWIO_LIGHTBOX_URL = base;
    const script = document.createElement('script');
    script.src = `${base}/viewer-static.min.js`;
    script.addEventListener('load', () =>
      window.GraphViewer
        ? resolve(window.GraphViewer)
        : reject(new Error('draw.io viewer did not load')),
    );
    script.addEventListener('error', () => reject(new Error('draw.io viewer could not be loaded')));
    document.head.append(script);
  }).catch((err: unknown) => {
    viewer = undefined;
    throw err;
  });
  return viewer;
}

const files = promiseCache<string>();

/** Renders a .drawio file read-only, with zoom and page switching and no edit controls. */
export default function DrawioEmbed({ src }: EmbedProps) {
  const xml = use(files(src, () => fetchText(src)));
  const graphViewer = use(loadViewer());
  const el = useRef<HTMLDivElement>(null);
  useEffect(() => {
    const node = el.current;
    if (!node) return undefined;
    const target = document.createElement('div');
    target.className = 'mxgraph max-w-full';
    target.dataset.mxgraph = JSON.stringify({
      xml,
      toolbar: 'zoom pages',
      'toolbar-nohide': true,
      nav: true,
      resize: true,
      editable: false,
      lightbox: false,
      highlight: latte.mauve,
    });
    node.append(target);
    graphViewer.createViewerForElement(target);
    return () => target.remove();
  }, [xml, graphViewer]);
  return <div ref={el} className="overflow-x-auto rounded-lg border border-line bg-page p-2" />;
}
