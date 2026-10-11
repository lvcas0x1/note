import type { Meta, StoryObj } from '@storybook/react-vite';
import { expect, waitFor, within } from 'storybook/test';
import { fakeArticle, seeded } from '../../test/fixtures.ts';
import { withRouter } from '../../test/router.tsx';
import { ArticleView } from '../article/ArticleView.tsx';

const base = { ...fakeArticle(seeded(9)), id: 'demo/embeds', title: 'Embeds' };
const article = (html: string) => ({
  ...base,
  html: `<p>Before the embed.</p>${html}<p>After the embed.</p>`,
});

const meta = {
  title: 'Embeds',
  component: ArticleView,
  decorators: [
    (Story) => (
      <div className="max-w-3xl p-6">
        <Story />
      </div>
    ),
    withRouter,
  ],
  args: { article: article('') },
} satisfies Meta<typeof ArticleView>;
export default meta;

type Story = StoryObj<typeof meta>;

const rest = async (canvas: ReturnType<typeof within>) => {
  await expect(canvas.getByText('Before the embed.')).toBeVisible();
  await expect(canvas.getByText('After the embed.')).toBeVisible();
};

export const Mermaid: Story = {
  args: { article: article('<pre data-embed="mermaid">flowchart LR; A--&gt;B</pre>') },
  play: async ({ canvas, canvasElement }) => {
    await waitFor(() => expect(canvasElement.querySelector('.embed svg')).not.toBeNull(), {
      timeout: 10_000,
    });
    const svg = canvasElement.querySelector('.embed svg');
    await expect(svg?.textContent).toContain('A');
    await expect(svg?.textContent).toContain('B');
    await rest(canvas);
  },
};

export const MermaidInvalid: Story = {
  args: { article: article('<pre data-embed="mermaid">this is not mermaid ((</pre>') },
  play: async ({ canvas }) => {
    const alert = await canvas.findByRole('alert', {}, { timeout: 10_000 });
    await expect(alert).toHaveTextContent('Mermaid diagram could not be shown');
    await rest(canvas);
  },
};

export const ECharts: Story = {
  args: {
    article: article(
      '<pre data-embed="echarts">{"xAxis":{"type":"category","data":["a","b","c"]},"yAxis":{"type":"value"},"series":[{"type":"bar","data":[3,1,2]}]}</pre>',
    ),
  },
  play: async ({ canvas, canvasElement }) => {
    await waitFor(() => expect(canvasElement.querySelector('.embed canvas')).not.toBeNull(), {
      timeout: 10_000,
    });
    await rest(canvas);
  },
};

export const Drawio: Story = {
  args: {
    article: article('<div data-embed="drawio" data-src="/fixtures/two-pages.drawio"></div>'),
  },
  play: async ({ canvas, canvasElement }) => {
    await waitFor(() => expect(canvasElement.querySelector('.embed .mxgraph svg')).not.toBeNull(), {
      timeout: 15_000,
    });
    const embed = canvasElement.querySelector<HTMLElement>('.embed');
    await expect(embed?.textContent).toContain('Client');
    // Read-only: no edit control and no link to an external editor.
    await expect(
      embed?.querySelector('[title="Edit"], [title="Edit Diagram"], a[href*="diagrams.net"]'),
    ).toBeNull();
    await rest(canvas);
  },
};

export const Excalidraw: Story = {
  args: {
    article: article('<div data-embed="excalidraw" data-src="/fixtures/shapes.excalidraw"></div>'),
  },
  play: async ({ canvas, canvasElement }) => {
    await waitFor(
      () => expect(canvasElement.querySelector('.embed svg[role="img"]')).not.toBeNull(),
      { timeout: 20_000 },
    );
    const svg = canvasElement.querySelector('.embed svg[role="img"]');
    await expect(svg?.querySelectorAll('path, rect, ellipse').length).toBeGreaterThan(0);
    await expect(canvasElement.querySelector('.excalidraw')).toBeNull(); // the editor is never mounted
    await rest(canvas);
  },
};

export const ExcalidrawCorrupt: Story = {
  args: {
    article: article('<div data-embed="excalidraw" data-src="/fixtures/corrupt.excalidraw"></div>'),
  },
  play: async ({ canvas }) => {
    const alert = await canvas.findByRole('alert', {}, { timeout: 10_000 });
    await expect(alert).toHaveTextContent('Excalidraw drawing could not be shown');
    await expect(alert).toHaveTextContent('/fixtures/corrupt.excalidraw');
    await rest(canvas);
  },
};

export const Swagger2: Story = {
  args: { article: article('<div data-embed="openapi" data-src="/fixtures/swagger2.json"></div>') },
  play: async ({ canvas }) => {
    const paths = await canvas.findAllByText('/pets', {}, { timeout: 15_000 });
    await expect(paths.length).toBeGreaterThanOrEqual(2); // GET and POST
    await expect(canvas.queryByRole('button', { name: /try it out/i })).toBeNull();
    await rest(canvas);
  },
};

export const OpenAPI31: Story = {
  args: {
    article: article('<div data-embed="openapi" data-src="/fixtures/openapi31.json"></div>'),
  },
  play: async ({ canvas }) => {
    await expect(await canvas.findByText('/api/v1/sidebar', {}, { timeout: 15_000 })).toBeVisible();
    await expect(canvas.queryByRole('button', { name: /try it out/i })).toBeNull();
  },
};

export const Video: Story = {
  args: { article: article('<div data-embed="mp4" data-src="/fixtures/not-a-video.mp4"></div>') },
  play: async ({ canvas, canvasElement }) => {
    const video = await waitFor(() => {
      const v = canvasElement.querySelector('video');
      if (!v) throw new Error('no video yet');
      return v;
    });
    await expect(video).toHaveAttribute('controls');
    await expect(video).toHaveAttribute('preload', 'metadata');
    await expect(video).toHaveAttribute('playsinline');
    await expect(video).not.toHaveAttribute('autoplay');
    await expect(video.paused).toBe(true);
    // The file cannot be played, so the download link appears.
    await expect(
      await canvas.findByText(
        'This browser cannot play the video.',
        { exact: false },
        { timeout: 10_000 },
      ),
    ).toBeVisible();
    await expect(
      canvas.getAllByRole('link', { name: 'Download the video' }).at(-1),
    ).toHaveAttribute('href', '/fixtures/not-a-video.mp4');
    await rest(canvas);
  },
};

/** A real H.264 video: no autoplay, metadata only until play, then plays inline. */
export const VideoPlays: Story = {
  args: { article: article('<div data-embed="mp4" data-src="/fixtures/demo-h264.mp4"></div>') },
  play: async ({ canvas, canvasElement }) => {
    const video = await waitFor(() => {
      const v = canvasElement.querySelector('video');
      if (!v) throw new Error('no video yet');
      return v;
    });
    // Metadata (duration) loads, but playback does not start on its own.
    await waitFor(() => expect(video.readyState).toBeGreaterThanOrEqual(1), { timeout: 10_000 });
    await expect(Math.round(video.duration)).toBe(4);
    await new Promise((r) => setTimeout(r, 500));
    await expect(video.paused).toBe(true);
    await expect(video.currentTime).toBe(0);

    video.muted = true; // browsers allow scripted play of muted media
    await video.play();
    await waitFor(() => expect(video.currentTime).toBeGreaterThan(0.3), { timeout: 10_000 });
    await expect(video.paused).toBe(false);
    await expect(video.error).toBeNull();
    await expect(
      canvas.queryByText('This browser cannot play the video.', { exact: false }),
    ).toBeNull();
    video.pause();
    await rest(canvas);
  },
};

export const UnreachableSource: Story = {
  args: { article: article('<div data-embed="drawio" data-src="/fixtures/missing.drawio"></div>') },
  play: async ({ canvas }) => {
    const alert = await canvas.findByRole('alert', {}, { timeout: 10_000 });
    await expect(alert).toHaveTextContent('/fixtures/missing.drawio');
    await rest(canvas);
  },
};
