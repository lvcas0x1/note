import type { Meta, StoryObj } from '@storybook/react-vite';
import { useState } from 'react';
import { createMemoryRouter, Outlet, RouterProvider, useNavigate } from 'react-router';
import { expect, userEvent, waitFor, within } from 'storybook/test';
import { fakeSource } from '../test/fake-source.ts';
import { fakeArticle, fakeSidebar, seeded } from '../test/fixtures.ts';
import { LocationProbe } from '../test/router.tsx';
import { AppProviders, appRoutes } from './App.tsx';

const f = seeded(11);
const sidebar = fakeSidebar(f);
const deep = sidebar[2]?.children?.[1]?.children?.[2];
const long = Array.from({ length: 60 }, (_p, i) => {
  const level = i / 10 + 1;
  const heading = i % 10 === 0 ? `<h${level} id="section-${i}">Section ${i}</h${level}>` : '';
  return `${heading}<p>Paragraph ${i} ${f.lorem.paragraph()}</p>`;
}).join('');
const articles = [
  { ...fakeArticle(f), id: 'c1/g1/a1', title: 'Zebra notes', date: '2026-10-05', html: long },
  { ...fakeArticle(f), id: 'c1/g1/a2', title: 'apple notes', date: '2026-10-05' },
  { ...fakeArticle(f), id: 'c2/g1/a1', title: 'Older entry', date: '2025-03-14' },
  { ...fakeArticle(f), id: 'c3/g2/a3', title: deep?.title ?? 'Deep', date: '2026-09-30' },
];
const source = fakeSource(sidebar, articles, (q) =>
  q === 'deep note'
    ? [
        {
          id: 'c3/g2/a3',
          title: deep?.title ?? 'Deep',
          date: '2026-09-30',
          summary: '',
          category: ['Category 3', 'Group 3.2'],
          snippet: '<mark>deep</mark> note',
        },
      ]
    : [],
);

/** Story-only controls: the current location and a browser-style Back button. */
function Harness() {
  const navigate = useNavigate();
  return (
    <>
      <Outlet />
      <LocationProbe />
      <button type="button" className="sr-only" onClick={() => void navigate(-1)}>
        Browser back
      </button>
    </>
  );
}

function App({ route }: { route: string }) {
  const [router] = useState(() =>
    createMemoryRouter([{ element: <Harness />, children: appRoutes({ today: '2026-10-11' }) }], {
      initialEntries: [route],
    }),
  );
  return (
    <AppProviders source={source}>
      <RouterProvider router={router} />
    </AppProviders>
  );
}

const indent = (el: Element) => parseFloat(getComputedStyle(el).paddingLeft);

const meta = {
  component: App,
  decorators: [
    (Story) => (
      <div style={{ width: 1200, height: 720 }}>
        <Story />
      </div>
    ),
  ],
  args: { route: '/' },
} satisfies Meta<typeof App>;
export default meta;

type Story = StoryObj<typeof meta>;

/** Click a day, see its list sorted by title, open an article, and go back. */
export const DayListFlow: Story = {
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    await userEvent.click(
      await canvas.findByRole('gridcell', { name: 'October 5, 2026, 2 articles' }),
    );
    const heading = await canvas.findByRole('heading', { name: 'October 5, 2026', level: 1 });
    const list = within(heading.closest('section') ?? canvasElement);
    const titles = (await list.findAllByRole('link')).map((l) => l.textContent);
    await expect(titles).toEqual(['apple notes', 'Zebra notes']);

    await userEvent.click(list.getByRole('link', { name: 'Zebra notes' }));
    await expect(
      await canvas.findByRole('heading', { name: 'Zebra notes', level: 1 }),
    ).toBeVisible();
    await expect(canvas.getByTestId('location')).toHaveTextContent('/a/c1/g1/a1');

    await userEvent.click(canvas.getByRole('button', { name: 'Browser back' }));
    await expect(
      await canvas.findByRole('heading', { name: 'October 5, 2026', level: 1 }),
    ).toBeVisible();
    await expect(canvas.getByTestId('location')).toHaveTextContent('/d/2026-10-05');
  },
};

/** An empty day shows an empty-state message. */
export const EmptyDay: Story = {
  args: { route: '/d/2026-10-06' },
  play: async ({ canvas }) => {
    await expect(await canvas.findByText('No articles on this day')).toBeVisible();
  },
};

/** The calendar follows the open article's month and date. */
export const CalendarFollowsArticle: Story = {
  args: { route: '/a/c2/g1/a1' },
  play: async ({ canvas }) => {
    await expect(
      await canvas.findByRole('heading', { name: 'Older entry', level: 1 }),
    ).toBeVisible();
    await expect(await canvas.findByRole('heading', { name: 'March 2025' })).toBeVisible();
    await expect(canvas.getByRole('gridcell', { name: /March 14, 2025/ })).toHaveAttribute(
      'aria-selected',
      'true',
    );
  },
};

/** The search bar stays at the top center while a long article scrolls. */
export const SearchStaysVisibleWhileScrolling: Story = {
  args: { route: '/a/c1/g1/a1' },
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    await canvas.findByText(/Paragraph 59/);
    const main = canvasElement.querySelector('main');
    if (!main) throw new Error('no main');
    main.scrollTop = main.scrollHeight;
    await new Promise((r) => setTimeout(r, 100));
    await expect(main.scrollTop).toBeGreaterThan(500);
    const search = canvas.getByRole('combobox', { name: 'Search notes' }).getBoundingClientRect();
    const area = main.getBoundingClientRect();
    await expect(search.top).toBeGreaterThanOrEqual(area.top);
    await expect(search.bottom).toBeLessThan(area.top + 80);
    const center = (search.left + search.right) / 2;
    await expect(Math.abs(center - (area.left + area.right) / 2)).toBeLessThan(40);
  },
};

/** Opening a nested search result navigates and locates it in the sidebar. */
export const SearchOpensNestedResult: Story = {
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    await expect(await canvas.findByRole('treeitem', { name: 'Category 3' })).toHaveAttribute(
      'aria-expanded',
      'false',
    );
    await userEvent.type(canvas.getByRole('combobox', { name: 'Search notes' }), 'deep note');
    await userEvent.click(
      await canvas.findByRole('option', { name: new RegExp(deep?.title ?? 'Deep') }),
    );
    await expect(canvas.getByTestId('location')).toHaveTextContent('/a/c3/g2/a3');
    await waitFor(() =>
      expect(canvas.getByRole('treeitem', { name: 'Category 3' })).toHaveAttribute(
        'aria-expanded',
        'true',
      ),
    );
    await expect(canvas.getByRole('treeitem', { name: 'Group 3.2' })).toHaveAttribute(
      'aria-expanded',
      'true',
    );
    await expect(canvas.getByRole('link', { name: deep?.title ?? 'Deep' })).toHaveAttribute(
      'aria-current',
      'page',
    );
  },
};

/** The content area keeps only 5% padding on each side. */
export const ContentPadding: Story = {
  args: { route: '/a/c1/g1/a1' },
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    await canvas.findByText(/Paragraph 59/);
    const main = canvasElement.querySelector('main');
    const content = main?.lastElementChild;
    if (!main || !(content instanceof HTMLElement)) throw new Error('no content area');
    const style = getComputedStyle(content);
    const expected = content.getBoundingClientRect().width * 0.05;
    await expect(Math.abs(parseFloat(style.paddingLeft) - expected)).toBeLessThan(1);
    await expect(Math.abs(parseFloat(style.paddingRight) - expected)).toBeLessThan(1);
  },
};

/** A table of contents (h1–h6) sits right of the text and stays visible while the article scrolls. */
export const TableOfContentsSticky: Story = {
  args: { route: '/a/c1/g1/a1' },
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    await canvas.findByText(/Paragraph 59/);
    const toc = canvas.getByRole('navigation', { name: 'Table of contents' });
    const links = within(toc).getAllByRole('link');
    await expect(links.map((l) => l.textContent)).toEqual([
      'Section 0',
      'Section 10',
      'Section 20',
      'Section 30',
      'Section 40',
      'Section 50',
    ]);
    // Deeper headings are indented further (h1 → h6 here).
    await expect(indent(links[0]!)).toBeLessThan(indent(links[1]!));
    // Right of the text.
    const body = canvasElement.querySelector('.article-body');
    await expect(toc.getBoundingClientRect().left).toBeGreaterThan(
      body!.getBoundingClientRect().right,
    );

    const main = canvasElement.querySelector('main')!;
    main.scrollTop = main.scrollHeight;
    await new Promise((r) => setTimeout(r, 150));
    const area = main.getBoundingClientRect();
    const box = toc.getBoundingClientRect();
    await expect(box.top).toBeGreaterThanOrEqual(area.top);
    await expect(box.top).toBeLessThan(area.top + 120);
    await waitFor(() => expect(links.at(-1)).toHaveAttribute('aria-current', 'location'));

    await userEvent.click(links[2]!);
    const target = canvasElement.querySelector('#section-20')!.getBoundingClientRect();
    await expect(target.top).toBeGreaterThanOrEqual(area.top);
    await expect(target.top).toBeLessThan(area.top + 120);
    await expect(links[2]).toHaveAttribute('aria-current', 'location');
  },
};
