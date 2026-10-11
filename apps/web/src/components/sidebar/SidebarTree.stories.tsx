import type { Meta, StoryObj } from '@storybook/react-vite';
import { expect, userEvent, waitFor } from 'storybook/test';
import { fakeSidebar, seeded } from '../../test/fixtures.ts';
import { withRouter } from '../../test/router.tsx';
import { SidebarTree } from './SidebarTree.tsx';

const nodes = fakeSidebar(seeded(3));
const deep = nodes[2]?.children?.[1]?.children?.[2];

const indent = (el: Element | undefined) =>
  parseFloat(getComputedStyle(el?.firstElementChild ?? document.body).paddingLeft);

const meta = {
  component: SidebarTree,
  decorators: [
    (Story) => (
      <div className="h-64 w-72 overflow-y-auto border border-line p-2" data-testid="scroller">
        <Story />
      </div>
    ),
    withRouter,
  ],
  args: { nodes },
} satisfies Meta<typeof SidebarTree>;
export default meta;

type Story = StoryObj<typeof meta>;

export const Collapsed: Story = {
  play: async ({ canvas }) => {
    const cat = canvas.getByRole('treeitem', { name: 'Category 1' });
    await expect(cat).toHaveAttribute('aria-expanded', 'false');
    await userEvent.click(canvas.getByText('Category 1'));
    await expect(cat).toHaveAttribute('aria-expanded', 'true');
    await expect(canvas.getByRole('treeitem', { name: 'Group 1.1' })).toBeVisible();
  },
};

export const ThreeLevels: Story = {
  play: async ({ canvas }) => {
    await userEvent.click(canvas.getByText('Category 1'));
    await userEvent.click(canvas.getByText('Group 1.1'));
    const level1 = canvas.getByRole('treeitem', { name: 'Category 1' });
    const level2 = canvas.getByRole('treeitem', { name: 'Group 1.1' });
    const leaf = canvas
      .getAllByRole('treeitem')
      .find((el) => el.getAttribute('aria-level') === '3');
    await expect(level1).toHaveAttribute('aria-level', '1');
    await expect(level2).toHaveAttribute('aria-level', '2');
    await expect(leaf).toBeDefined();
    await expect(indent(level1)).toBeLessThan(indent(level2));
    await expect(indent(level2)).toBeLessThan(indent(leaf));
  },
};

export const ClickArticle: Story = {
  play: async ({ canvas }) => {
    await userEvent.click(canvas.getByText('Category 1'));
    await userEvent.click(canvas.getByText('Group 1.1'));
    const first = nodes[0]?.children?.[0]?.children?.[0];
    await userEvent.click(canvas.getByRole('link', { name: first?.title ?? '' }));
    await expect(canvas.getByTestId('location')).toHaveTextContent('/a/c1/g1/a1');
  },
};

export const KeyboardOpen: Story = {
  play: async ({ canvas }) => {
    const first = canvas.getByRole('treeitem', { name: 'Category 1' });
    first.focus();
    await userEvent.keyboard('{ArrowRight}'); // expand Category 1
    await expect(first).toHaveAttribute('aria-expanded', 'true');
    await userEvent.keyboard('{ArrowDown}{ArrowRight}'); // focus and expand Group 1.1
    await userEvent.keyboard('{ArrowDown}{ArrowDown}'); // second article
    await expect(canvas.getAllByRole('treeitem')[3]).toHaveFocus();
    await userEvent.keyboard('{Enter}');
    await expect(canvas.getByTestId('location')).toHaveTextContent('/a/c1/g1/a2');
    await userEvent.keyboard('{ArrowLeft}'); // to parent group
    await expect(canvas.getByRole('treeitem', { name: 'Group 1.1' })).toHaveFocus();
    await userEvent.keyboard('{ArrowLeft}'); // collapse it
    await expect(canvas.getByRole('treeitem', { name: 'Group 1.1' })).toHaveAttribute(
      'aria-expanded',
      'false',
    );
  },
};

/** Opening a deep article (from search or a URL) locates it in the tree. */
export const LocateDeepArticle: Story = {
  args: { currentPath: deep?.path },
  play: async ({ canvas }) => {
    await expect(canvas.getByRole('treeitem', { name: 'Category 3' })).toHaveAttribute(
      'aria-expanded',
      'true',
    );
    await expect(canvas.getByRole('treeitem', { name: 'Group 3.2' })).toHaveAttribute(
      'aria-expanded',
      'true',
    );
    const link = canvas.getByRole('link', { name: deep?.title ?? '' });
    await expect(link).toHaveAttribute('aria-current', 'page');
    const scroller = canvas.getByTestId('scroller').getBoundingClientRect();
    await waitFor(async () => {
      const box = link.getBoundingClientRect();
      await expect(box.top).toBeGreaterThanOrEqual(scroller.top - 1);
      await expect(box.bottom).toBeLessThanOrEqual(scroller.bottom + 1);
    });
  },
};

/** An article that is not in the sidebar highlights nothing. */
export const ArticleNotInSidebar: Story = {
  args: { currentPath: 'elsewhere/hidden.md' },
  play: async ({ canvas }) => {
    await expect(canvas.queryByRole('link', { current: 'page' })).toBeNull();
    for (const item of canvas.getAllByRole('treeitem')) {
      await expect(item).toHaveAttribute('aria-selected', 'false');
    }
  },
};
