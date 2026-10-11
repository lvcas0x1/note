import { describe, expect, it } from 'vitest';
import { articleRoute, locate, visibleNodes } from '../src/components/sidebar/tree.ts';

const nodes = [
  { title: 'A', children: [{ title: 'A1', children: [{ title: 'x', path: 'a/x.md' }] }] },
  { title: 'B', path: 'b.md' },
];

describe('tree helpers', () => {
  it('shows only expanded branches', () => {
    expect(visibleNodes(nodes, new Set()).map((n) => n.key)).toEqual(['0', '1']);
    expect(visibleNodes(nodes, new Set(['0', '0.0'])).map((n) => `${n.key}@${n.level}`)).toEqual([
      '0@1',
      '0.0@2',
      '0.0.0@3',
      '1@1',
    ]);
  });

  it('locates a node and its ancestors', () => {
    expect(locate(nodes, 'a/x.md')).toEqual({ key: '0.0.0', ancestors: ['0', '0.0'] });
    expect(locate(nodes, 'missing.md')).toBeNull();
  });

  it('maps paths to routes', () => {
    expect(articleRoute('aws/s3/s3.md')).toBe('/a/aws/s3/s3');
  });
});
