/** Props every embed renderer receives. */
export interface EmbedProps {
  /** Asset URL for file embeds (drawio, excalidraw, openapi, mp4). */
  src: string;
  /** Inline source for fence embeds (mermaid, echarts). */
  body: string;
}

export const embedKinds = ['drawio', 'excalidraw', 'openapi', 'mp4', 'mermaid', 'echarts'] as const;
export type EmbedKind = (typeof embedKinds)[number];

export function isEmbedKind(v: string | undefined): v is EmbedKind {
  return (embedKinds as readonly string[]).includes(v ?? '');
}

export const embedLabels: Record<EmbedKind, string> = {
  drawio: 'draw.io diagram',
  excalidraw: 'Excalidraw drawing',
  openapi: 'OpenAPI reference',
  mp4: 'video',
  mermaid: 'Mermaid diagram',
  echarts: 'ECharts chart',
};
