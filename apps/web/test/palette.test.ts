// @vitest-environment node
import { readFile } from 'node:fs/promises';
import { describe, expect, it } from 'vitest';
import { latte, tokens } from '../src/lib/palette.ts';

const lin = (c: number) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4);

function luminance(hex: string): number {
  const [r = 0, g = 0, b = 0] = [1, 3, 5].map((i) => parseInt(hex.slice(i, i + 2), 16) / 255);
  return 0.2126 * lin(r) + 0.7152 * lin(g) + 0.0722 * lin(b);
}

function contrast(a: string, b: string): number {
  const [hi, lo] = [luminance(a), luminance(b)].toSorted((x, y) => y - x);
  return ((hi ?? 0) + 0.05) / ((lo ?? 0) + 0.05);
}

describe('Catppuccin Latte theme', () => {
  it('defines the same tokens in CSS and palette.ts', async () => {
    const css = await readFile(new URL('../src/styles/index.css', import.meta.url), 'utf8');
    for (const [name, value] of Object.entries(tokens)) {
      expect(css).toContain(`--color-${name}: ${value};`);
    }
  });

  // Every text/background pair the UI uses must meet WCAG AA for normal text (4.5:1).
  it.each([
    ['text on page', tokens.ink, tokens.page],
    ['text on sidebar', tokens.ink, tokens.surface],
    ['muted text on page', tokens.muted, tokens.page],
    ['muted text on sidebar', tokens.muted, tokens.surface],
    ['links and TOC (accent) on page', tokens.accent, tokens.page],
    ['accent on sidebar', tokens.accent, tokens.surface],
    ['white on selected calendar day', '#ffffff', tokens.accent],
    ['text on selected sidebar item and tags', tokens.ink, tokens['accent-soft']],
    ['text on search highlight', tokens.ink, tokens.mark],
    ['text on embed error box', tokens.ink, tokens['danger-soft']],
  ])('%s', (_name, fg, bg) => {
    expect(contrast(fg, bg)).toBeGreaterThanOrEqual(4.5);
  });

  it('uses the official Latte base colors', () => {
    expect(latte.base).toBe('#eff1f5');
    expect(latte.text).toBe('#4c4f69');
    expect(latte.mauve).toBe('#8839ef');
  });
});
