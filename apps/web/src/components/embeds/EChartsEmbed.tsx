import * as echarts from 'echarts';

// Catppuccin Latte chart theme: series colors and muted axes.
echarts.registerTheme('catppuccin-latte', {
  color: seriesColors,
  backgroundColor: 'transparent',
  textStyle: { color: latte.text },
  title: { textStyle: { color: latte.text }, subtextStyle: { color: latte.subtext1 } },
  legend: { textStyle: { color: latte.subtext1 } },
  categoryAxis: {
    axisLine: { lineStyle: { color: latte.overlay0 } },
    axisLabel: { color: latte.subtext1 },
    splitLine: { lineStyle: { color: [latte.surface0] } },
  },
  valueAxis: {
    axisLine: { lineStyle: { color: latte.overlay0 } },
    axisLabel: { color: latte.subtext1 },
    splitLine: { lineStyle: { color: [latte.surface0] } },
  },
  tooltip: {
    backgroundColor: latte.mantle,
    borderColor: latte.surface1,
    textStyle: { color: latte.text },
  },
});

function isOption(v: unknown): v is echarts.EChartsOption {
  return typeof v === 'object' && v !== null && !Array.isArray(v);
}
import { useEffect, useRef } from 'react';
import { latte, seriesColors } from '../../lib/palette.ts';
import type { EmbedProps } from './types.ts';

/** Renders an echarts fence: its body is the ECharts option as JSON. The chart resizes with the content. */
export default function EChartsEmbed({ body }: EmbedProps) {
  const option: unknown = JSON.parse(body);
  if (!isOption(option)) throw new Error('the ECharts option must be a JSON object');
  const el = useRef<HTMLDivElement>(null);
  useEffect(() => {
    const node = el.current;
    if (!node) return undefined;
    const chart = echarts.init(node, 'catppuccin-latte', { renderer: 'canvas' });
    chart.setOption(option);
    const observer = new ResizeObserver(() => chart.resize());
    observer.observe(node);
    return () => {
      observer.disconnect();
      chart.dispose();
    };
  }, [option]);
  return <div ref={el} className="h-80 w-full" />;
}
