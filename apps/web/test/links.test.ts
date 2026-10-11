import { describe, expect, it } from 'vitest';
import { inAppRoute } from '../src/components/article/links.ts';

describe('inAppRoute', () => {
  it.each([
    ['../tips/tips.md', '/a/golang/tips/tips'],
    ['../tips/tips.md#errors', '/a/golang/tips/tips#errors'],
    ['/aws/s3/s3.md', '/a/aws/s3/s3'],
    ['practical.md', '/a/golang/basic/practical'],
    ['/a/aws/s3/s3', '/a/aws/s3/s3'],
    ['/notes/a/x', '/a/x'],
  ])('%s → %s', (href, want) => {
    expect(inAppRoute(href, 'golang/basic/basic', '/notes/')).toBe(want);
  });

  it.each([
    'https://go.dev',
    'mailto:a@b.c',
    '#intro',
    '/assets/0123/a.png',
    'file.pdf',
    '//cdn.example.com/x',
  ])('leaves %s alone', (href) => {
    expect(inAppRoute(href, 'golang/basic/basic')).toBeNull();
  });
});
