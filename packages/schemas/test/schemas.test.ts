import { describe, expect, it } from 'vitest';
import { FrontMatter, Sidebar } from '../src/index.ts';

describe('Sidebar', () => {
  it('accepts three-level nesting and keeps order', () => {
    const input = [
      {
        title: 'Cloud',
        children: [
          {
            title: 'AWS',
            path: 'aws/README.md',
            children: [{ title: 'S3', path: 'aws/s3/s3.md' }],
          },
        ],
      },
      { title: 'Go', path: 'golang/basic/basic.md' },
    ];
    expect(Sidebar.parse(input)).toEqual(input);
  });

  it('rejects a node without a title', () => {
    const result = Sidebar.safeParse([{ path: 'a.md' }]);
    expect(result.success).toBe(false);
  });

  it('rejects an empty title', () => {
    expect(Sidebar.safeParse([{ title: '  ', path: 'a.md' }]).success).toBe(false);
  });

  it('rejects a node with neither path nor children', () => {
    expect(Sidebar.safeParse([{ title: 'Lonely' }]).success).toBe(false);
  });

  it('rejects a nested invalid node', () => {
    expect(Sidebar.safeParse([{ title: 'A', children: [{ title: 'B' }] }]).success).toBe(false);
  });

  it('rejects a path that is not a relative .md file', () => {
    expect(Sidebar.safeParse([{ title: 'A', path: '/abs.md' }]).success).toBe(false);
    expect(Sidebar.safeParse([{ title: 'A', path: 'a.txt' }]).success).toBe(false);
  });
});

describe('FrontMatter', () => {
  it('accepts all recognized keys', () => {
    const fm = FrontMatter.parse({
      title: 'T',
      date: '2026-10-11',
      summary: 'S',
      tags: ['go'],
      draft: true,
    });
    expect(fm).toMatchObject({ date: '2026-10-11', draft: true });
  });

  it('defaults draft to false', () => {
    expect(FrontMatter.parse({}).draft).toBe(false);
  });

  it.each(['2026-02-30', '2026-13-01', '2026-1-5', '20261011', '2026-10-11T00:00:00Z'])(
    'rejects the impossible or malformed date %s',
    (date) => {
      expect(FrontMatter.safeParse({ date }).success).toBe(false);
    },
  );

  it('rejects a non-boolean draft', () => {
    expect(FrontMatter.safeParse({ draft: 'yes' }).success).toBe(false);
  });

  it('keeps unknown keys so the build can warn about them', () => {
    expect(FrontMatter.parse({ layout: 'post' })).toHaveProperty('layout', 'post');
  });
});
