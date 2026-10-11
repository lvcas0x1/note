import type { Meta, StoryObj } from '@storybook/react-vite';
import { expect, userEvent, within } from 'storybook/test';
import { seeded } from '../../test/fixtures.ts';
import { withRouter } from '../../test/router.tsx';
import { AppShell } from './AppShell.tsx';

const f = seeded(7);

const meta = {
  component: AppShell,
  decorators: [withRouter],
  args: {
    calendar: <div className="h-48 rounded bg-page p-2 text-sm">Calendar</div>,
    tree: <p className="text-sm">Category tree</p>,
    search: <input aria-label="Search" className="w-full rounded border border-line px-2 py-1" />,
    children: (
      <article>
        <h1 className="text-2xl font-semibold">Article</h1>
        {Array.from({ length: 30 }, (_p, i) => (
          <p key={i} className="mt-3">
            {f.lorem.paragraph()}
          </p>
        ))}
      </article>
    ),
  },
} satisfies Meta<typeof AppShell>;
export default meta;

type Story = StoryObj<typeof meta>;

export const Desktop: Story = {
  decorators: [
    (Story) => (
      <div style={{ width: 1200, height: 700 }}>
        <Story />
      </div>
    ),
  ],
  play: async ({ canvas }) => {
    const sidebar = canvas.getByRole('complementary', { name: 'Sidebar' });
    await expect(sidebar).toBeVisible();
    await expect(canvas.queryByRole('button', { name: 'Menu' })).toBeNull();
    // Calendar is at the top-left, the tree below it, the content on the right.
    const cal = canvas.getByText('Calendar').getBoundingClientRect();
    const tree = canvas.getByText('Category tree').getBoundingClientRect();
    const heading = canvas.getByRole('heading', { name: 'Article' }).getBoundingClientRect();
    await expect(cal.top).toBeLessThan(tree.top);
    await expect(cal.left).toBeLessThan(heading.left);
  },
};

export const Mobile: Story = {
  decorators: [
    (Story) => (
      <div style={{ width: 375, height: 700 }}>
        <Story />
      </div>
    ),
  ],
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);
    const toggle = canvas.getByRole('button', { name: 'Menu' });
    await expect(canvas.queryByRole('complementary', { name: 'Sidebar' })).toBeNull();
    await expect(toggle).toHaveAttribute('aria-expanded', 'false');
    // The content uses the full width when the sidebar is collapsed.
    const main = canvasElement.querySelector('main');
    await expect(main?.getBoundingClientRect().width).toBeGreaterThan(360);

    await userEvent.click(toggle);
    await expect(toggle).toHaveAttribute('aria-expanded', 'true');
    await expect(canvas.getByRole('complementary', { name: 'Sidebar' })).toBeVisible();

    await userEvent.click(canvas.getByRole('button', { name: 'Close sidebar' }));
    await expect(canvas.queryByRole('complementary', { name: 'Sidebar' })).toBeNull();
  },
};
