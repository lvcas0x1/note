import { describe, expect, it } from 'vitest';
import { headingsOf } from '../src/components/article/TableOfContents.tsx';

describe('headingsOf', () => {
  it('lists h1–h6 with ids in document order', () => {
    const html =
      '<h1 id="a">A</h1><p>x</p><h2 id="b">B <code>c</code></h2><h3 id="c">C</h3><h4 id="d">D</h4><h5 id="e">E</h5><h6 id="f">F</h6><h2>no id</h2>';
    expect(headingsOf(html)).toEqual([
      { id: 'a', text: 'A', level: 1 },
      { id: 'b', text: 'B c', level: 2 },
      { id: 'c', text: 'C', level: 3 },
      { id: 'd', text: 'D', level: 4 },
      { id: 'e', text: 'E', level: 5 },
      { id: 'f', text: 'F', level: 6 },
    ]);
  });

  it('is empty for an article without headings', () => {
    expect(headingsOf('<p>plain</p>')).toEqual([]);
  });
});
