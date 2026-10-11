import mermaid from 'mermaid';
import { use } from 'react';
import { promiseCache } from './cache.ts';
import { latte, tokens } from '../../lib/palette.ts';
import type { EmbedProps } from './types.ts';

let ready = false;
let counter = 0;
const svgs = promiseCache<string>();

async function renderMermaid(source: string): Promise<string> {
  if (!ready) {
    // strict: no click handlers or HTML labels from the diagram source.
    mermaid.initialize({
      startOnLoad: false,
      securityLevel: 'strict',
      theme: 'base',
      themeVariables: {
        background: latte.base,
        primaryColor: latte.crust,
        primaryBorderColor: latte.mauve,
        primaryTextColor: latte.text,
        secondaryColor: latte.mantle,
        tertiaryColor: latte.base,
        lineColor: latte.subtext1,
        textColor: latte.text,
        noteBkgColor: tokens.mark,
        noteTextColor: latte.text,
        actorBkg: latte.crust,
        actorBorder: latte.mauve,
        signalColor: latte.text,
        fontFamily: 'inherit',
      },
    });
    ready = true;
  }
  const id = `mermaid-${++counter}`;
  try {
    return (await mermaid.render(id, source)).svg;
  } finally {
    // Mermaid leaves a temporary element behind when rendering fails.
    document.getElementById(`d${id}`)?.remove();
  }
}

/** Renders a Mermaid fence as an inline SVG diagram. */
export default function MermaidEmbed({ body }: EmbedProps) {
  const svg = use(svgs(body, () => renderMermaid(body)));
  return (
    <div
      className="flex justify-center overflow-x-auto"
      // Mermaid sanitizes its SVG output in strict mode.
      dangerouslySetInnerHTML={{ __html: svg }}
    />
  );
}
