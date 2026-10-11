import type { Meta, StoryObj } from '@storybook/react-vite';
import { codePoints } from '@notes/schemas/query';
import { expect, userEvent, waitFor, within } from 'storybook/test';
import { AppProviders } from '../../app/App.tsx';
import { ApiError } from '../../lib/api.ts';
import { fakeSource } from '../../test/fake-source.ts';
import { withRouter } from '../../test/router.tsx';
import { SearchBar } from './SearchBar.tsx';

const hit = (id: string, title: string, snippet: string) => ({
  id,
  title,
  date: '2026-10-05',
  summary: '',
  category: ['Cloud', 'Messaging'],
  snippet,
});

const searched: string[] = [];
const source = fakeSource([], [], (q) => {
  searched.push(q);
  if (q === 'busy') throw new ApiError(429, { status: 429, title: 'Too Many Requests' });
  if (q.toLowerCase().includes('kafka')) {
    return [
      hit('kafka/deploy', 'Kafka deploy', '<mark>Kafka</mark> deploy'),
      hit('notes/streams', 'Stream notes', '…compared <mark>kafka</mark> with…'),
    ];
  }
  if (q === '本語の') return [hit('jp/memo', '日本語メモ', '日<mark>本語の</mark>テキスト')];
  return [];
});

const meta = {
  component: SearchBar,
  decorators: [
    (Story) => (
      <AppProviders source={source}>
        <div className="mx-auto max-w-xl p-6">
          <Story />
        </div>
      </AppProviders>
    ),
    withRouter,
  ],
} satisfies Meta<typeof SearchBar>;
export default meta;

type Story = StoryObj<typeof meta>;

const box = (c: ReturnType<typeof within>) => c.getByRole('combobox', { name: 'Search notes' });

export const ShowResults: Story = {
  play: async ({ canvas }) => {
    await userEvent.type(box(canvas), 'kafka');
    const options = await canvas.findAllByRole('option');
    await expect(options).toHaveLength(2);
    await expect(options[0]).toHaveTextContent('Kafka deploy');
    await expect(options[0]).toHaveTextContent('October 5, 2026');
    await expect(options[0]).toHaveTextContent('Cloud / Messaging');
    await expect(options[0]?.querySelector('mark')?.textContent).toBe('Kafka');
  },
};

export const KeyboardShortcuts: Story = {
  play: async ({ canvas }) => {
    (document.activeElement as HTMLElement | null)?.blur();
    await userEvent.keyboard('/');
    await expect(box(canvas)).toHaveFocus();
    await expect(box(canvas)).toHaveValue('');
    box(canvas).blur();
    await userEvent.keyboard('{Control>}k{/Control}');
    await expect(box(canvas)).toHaveFocus();
  },
};

export const QueryTooShort: Story = {
  play: async ({ canvas }) => {
    searched.length = 0;
    await userEvent.type(box(canvas), ' ab ');
    await expect(await canvas.findByText('Type at least 3 characters to search.')).toBeVisible();
    await new Promise((r) => setTimeout(r, 400));
    await expect(searched).toEqual([]);
  },
};

export const Debounce: Story = {
  play: async ({ canvas }) => {
    searched.length = 0;
    await userEvent.type(box(canvas), 'kafka streams', { delay: 20 });
    await waitFor(() => expect(searched.length).toBeGreaterThan(0));
    await new Promise((r) => setTimeout(r, 400));
    // Typing faster than the debounce sends only the final query.
    await expect(searched).toEqual(['kafka streams']);
  },
};

export const JapaneseSubstring: Story = {
  play: async ({ canvas }) => {
    await userEvent.type(box(canvas), '本語の');
    await expect(await canvas.findByRole('option', { name: /日本語メモ/ })).toBeVisible();
  },
};

export const NoResults: Story = {
  play: async ({ canvas }) => {
    await userEvent.type(box(canvas), 'zebracorn');
    await expect(await canvas.findByText(/No articles match/)).toBeVisible();
  },
};

export const SpecialCharacters: Story = {
  play: async ({ canvas }) => {
    await userEvent.type(box(canvas), '"foo* OR (bar');
    await expect(await canvas.findByText(/No articles match/)).toBeVisible();
  },
};

export const KeyboardSelection: Story = {
  play: async ({ canvas }) => {
    await userEvent.type(box(canvas), 'kafka');
    await canvas.findAllByRole('option');
    await userEvent.keyboard('{ArrowDown}{ArrowDown}');
    await expect(canvas.getAllByRole('option')[1]).toHaveAttribute('aria-selected', 'true');
    await userEvent.keyboard('{Enter}');
    await expect(canvas.getByTestId('location')).toHaveTextContent('/a/notes/streams');
    await expect(canvas.queryByRole('listbox')).toBeNull();
  },
};

export const EscapeCloses: Story = {
  play: async ({ canvas }) => {
    await userEvent.type(box(canvas), 'kafka');
    await canvas.findAllByRole('option');
    await userEvent.keyboard('{Escape}');
    await expect(canvas.queryByRole('listbox')).toBeNull();
    await expect(box(canvas)).toHaveAttribute('aria-expanded', 'false');
  },
};

export const ClickResult: Story = {
  play: async ({ canvas }) => {
    await userEvent.type(box(canvas), 'kafka');
    await userEvent.click(await canvas.findByRole('option', { name: /Kafka deploy/ }));
    await expect(canvas.getByTestId('location')).toHaveTextContent('/a/kafka/deploy');
  },
};

export const RateLimited: Story = {
  play: async ({ canvas }) => {
    await userEvent.type(box(canvas), 'busy');
    await expect(await canvas.findByText(/Too many searches/)).toBeVisible();
  },
};

export const InputCap: Story = {
  play: async ({ canvas }) => {
    const long = 'あ'.repeat(205);
    await userEvent.click(box(canvas));
    await userEvent.paste(long);
    await expect(codePoints((box(canvas) as HTMLInputElement).value)).toBe(200);
  },
};
