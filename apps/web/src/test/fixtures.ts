import { faker } from '@faker-js/faker';
import type { ApiNode, ArticleDetail, ArticleSummary } from '@notes/schemas/api';

/** Seeded fake data so stories and tests render the same content every run. */
export function seeded(seed = 42) {
  faker.seed(seed);
  return faker;
}

export function fakeSummary(f = seeded(), date = '2026-10-05'): ArticleSummary {
  const title = f.lorem.words({ min: 2, max: 5 }).replace(/^./, (c) => c.toUpperCase());
  return {
    id: `${f.word.noun()}/${f.helpers.slugify(title).toLowerCase()}`,
    title,
    date,
    summary: f.lorem.sentence({ min: 8, max: 16 }),
    category: [f.word.noun()],
  };
}

export function fakeArticle(f = seeded()): ArticleDetail {
  const s = fakeSummary(f);
  return {
    ...s,
    tags: [f.word.noun(), f.word.noun()],
    inSidebar: true,
    html: `<h2 id="intro">Intro</h2><p>${f.lorem.paragraph()}</p><pre class="chroma"><code class="language-go">fmt.Println("hi")</code></pre>`,
  };
}

/** A three-level sidebar tree: categories, sub-categories, articles. */
export function fakeSidebar(f = seeded()): ApiNode[] {
  return Array.from({ length: 3 }, (_, i) => ({
    title: `Category ${i + 1}`,
    children: Array.from({ length: 2 }, (_g, j) => ({
      title: `Group ${i + 1}.${j + 1}`,
      children: Array.from({ length: 3 }, (_a, k) => ({
        title: f.lorem.words(2).replace(/^./, (c) => c.toUpperCase()),
        path: `c${i + 1}/g${j + 1}/a${k + 1}.md`,
      })),
    })),
  }));
}
