/**
 * Maps a link in article HTML to an in-app route, or returns null for links
 * that should behave normally (other sites, assets, in-page anchors).
 * Markdown links to other articles look like "../go/basics.md#intro".
 */
export function inAppRoute(href: string, articleId: string, basePath = '/'): string | null {
  if (!href || href.startsWith('#') || /^[a-z][a-z\d+.-]*:/i.test(href) || href.startsWith('//'))
    return null;
  const [pathPart = '', hash = ''] = href.split('#', 2);
  const base = basePath.replace(/\/$/, '');
  let path = pathPart.split('?')[0] ?? '';
  if (base && path.startsWith(`${base}/`)) path = path.slice(base.length);
  if (path.startsWith('/assets/') || path === '') return null;
  if (path.startsWith('/a/') || path.startsWith('/d/') || path === '/')
    return path + (hash ? `#${hash}` : '');
  if (!/\.md$/i.test(path)) return null;
  const dir = articleId.includes('/') ? articleId.slice(0, articleId.lastIndexOf('/')) : '';
  const parts = (path.startsWith('/') ? path.slice(1) : `${dir}/${path}`).split('/');
  const out: string[] = [];
  for (const part of parts) {
    if (part === '' || part === '.') continue;
    if (part === '..') out.pop();
    else out.push(part);
  }
  return `/a/${out.join('/').replace(/\.md$/i, '')}${hash ? `#${hash}` : ''}`;
}
