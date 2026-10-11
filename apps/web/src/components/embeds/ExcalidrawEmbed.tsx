import { use, useEffect, useRef } from 'react';
import { readConfig } from '../../lib/config.ts';
import { fetchText, promiseCache } from './cache.ts';
import { latte } from '../../lib/palette.ts';
import type { EmbedProps } from './types.ts';

declare global {
  interface Window {
    EXCALIDRAW_ASSET_PATH?: string | string[];
  }
}

interface Scene {
  type?: unknown;
  elements?: unknown;
  appState?: Record<string, unknown>;
  files?: unknown;
}

const svgs = promiseCache<string>();

function isScene(v: unknown): v is Scene & { elements: unknown[] } {
  return (
    typeof v === 'object' &&
    v !== null &&
    'type' in v &&
    v.type === 'excalidraw' &&
    'elements' in v &&
    Array.isArray(v.elements)
  );
}

type ExportInput = Parameters<typeof import('@excalidraw/excalidraw').exportToSvg>[0];

async function renderScene(src: string): Promise<string> {
  const scene: unknown = JSON.parse(await fetchText(src));
  if (!isScene(scene)) throw new Error('not an Excalidraw scene');
  // Fonts are served from the site, not from the default CDN.
  window.EXCALIDRAW_ASSET_PATH = `${location.origin}${readConfig().basePath}vendor/excalidraw/`;
  const { exportToSvg } = await import('@excalidraw/excalidraw');
  const svg = await exportToSvg({
    // The scene file is the author's own export; Excalidraw checks element shapes itself.
    // oxlint-disable-next-line typescript/no-unsafe-type-assertion
    elements: scene.elements as ExportInput['elements'],
    appState: {
      ...scene.appState,
      exportBackground: true,
      viewBackgroundColor:
        typeof scene.appState?.viewBackgroundColor === 'string'
          ? scene.appState.viewBackgroundColor
          : latte.base,
    },
    // oxlint-disable-next-line typescript/no-unsafe-type-assertion
    files: (scene.files ?? null) as ExportInput['files'],
  });
  svg.removeAttribute('width');
  svg.removeAttribute('height');
  svg.setAttribute('role', 'img');
  return svg.outerHTML;
}

/** Renders an .excalidraw scene as a read-only SVG that fits the content width. The editor is never loaded. */
export default function ExcalidrawEmbed({ src }: EmbedProps) {
  const svg = use(svgs(src, () => renderScene(src)));
  const el = useRef<HTMLDivElement>(null);
  useEffect(() => {
    const node = el.current;
    if (!node) return;
    // The SVG comes from Excalidraw's exporter, not from article HTML.
    node.innerHTML = svg;
  }, [svg]);
  return <div ref={el} className="[&>svg]:h-auto [&>svg]:max-w-full" />;
}
