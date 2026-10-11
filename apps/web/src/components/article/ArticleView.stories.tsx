import type { Meta, StoryObj } from '@storybook/react-vite';
import { expect, userEvent, waitFor } from 'storybook/test';
import { fakeArticle, seeded } from '../../test/fixtures.ts';
import { withRouter } from '../../test/router.tsx';
import { ArticleView } from './ArticleView.tsx';

const base = fakeArticle(seeded(5));
const article = {
  ...base,
  id: 'golang/basic/basic',
  title: 'Go basics',
  date: '2026-10-05',
  tags: ['golang'],
  category: ['Golang'],
  html:
    '<h2 id="hello">Hello</h2><p>See <a href="../tips/tips.md#errors">the tips</a> and <a href="https://go.dev">go.dev</a>.</p>' +
    '<pre class="chroma"><code class="language-go"><span class="kd">func</span> main() {}</code></pre>' +
    '<table><thead><tr><th>a</th></tr></thead><tbody><tr><td>1</td></tr></tbody></table>',
};

const meta = {
  component: ArticleView,
  decorators: [
    (Story) => (
      <div className="max-w-3xl p-6">
        <Story />
      </div>
    ),
    withRouter,
  ],
  args: { article },
} satisfies Meta<typeof ArticleView>;
export default meta;

type Story = StoryObj<typeof meta>;

export const ShowArticle: Story = {
  play: async ({ canvas }) => {
    await expect(canvas.getByRole('heading', { level: 1, name: 'Go basics' })).toBeVisible();
    await expect(canvas.getByText('October 5, 2026')).toBeVisible();
    await expect(canvas.getByText('golang')).toBeVisible();
    // Highlighted code keeps the chroma classes.
    const keyword = canvas.getByText('func');
    await expect(keyword).toHaveClass('kd');
    await expect(getComputedStyle(keyword).color).not.toBe(
      getComputedStyle(keyword.parentElement!.parentElement!).color,
    );
  },
};

/** An article without embeds downloads no renderer code. */
export const PlainArticleLoadsNoRenderers: Story = {
  play: async ({ canvas }) => {
    await expect(canvas.getByRole('heading', { level: 1, name: 'Go basics' })).toBeVisible();
    await new Promise((r) => setTimeout(r, 300));
    const loaded = performance.getEntriesByType('resource').map((e) => e.name);
    const renderers = loaded.filter((n) =>
      /Embed\.tsx|mermaid|echarts|excalidraw|swagger-ui|viewer-static/.test(n),
    );
    await expect(renderers).toEqual([]);
  },
};

export const CopyCode: Story = {
  play: async ({ canvas, canvasElement }) => {
    let copied = '';
    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: { writeText: async (t: string) => void (copied = t) },
    });
    const button = await canvas.findByRole('button', { name: 'Copy code' });
    // An icon, not text: the label is for screen readers and the tooltip.
    await expect(button.textContent).toBe('');
    await expect(button.querySelector('svg')).not.toBeNull();
    await expect(button).toHaveAttribute('title', 'Copy code');
    await userEvent.click(button);
    await waitFor(() => expect(button).toHaveAttribute('aria-label', 'Copied'));
    await expect(button).toHaveAttribute('title', 'Copied');
    await expect(button).toHaveAttribute('data-state', 'check');
    await expect(copied).toContain('func main() {}');
    await expect(copied).toBe('func main() {}');
    await expect(canvasElement.querySelectorAll('pre button')).toHaveLength(1);
  },
};

export const InternalLink: Story = {
  play: async ({ canvas }) => {
    await userEvent.click(canvas.getByRole('link', { name: 'the tips' }));
    await expect(canvas.getByTestId('location')).toHaveTextContent('/a/golang/tips/tips');
  },
};
