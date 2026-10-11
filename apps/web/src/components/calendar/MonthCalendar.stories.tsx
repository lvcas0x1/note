import type { Meta, StoryObj } from '@storybook/react-vite';
import { expect, fn, userEvent } from 'storybook/test';
import { MonthCalendar } from './MonthCalendar.tsx';

const meta = {
  component: MonthCalendar,
  decorators: [
    (Story) => (
      <div className="w-72 p-3">
        <Story />
      </div>
    ),
  ],
  args: {
    today: '2026-10-11',
    counts: { '2026-10-05': 2, '2026-10-11': 1, '2026-11-03': 4 },
    onSelectDate: fn(),
    onMonthChange: fn(),
  },
} satisfies Meta<typeof MonthCalendar>;
export default meta;

type Story = StoryObj<typeof meta>;

export const InitialMonth: Story = {
  play: async ({ canvas, args }) => {
    await expect(canvas.getByRole('heading', { name: 'October 2026' })).toBeVisible();
    const today = canvas.getByRole('gridcell', { name: /October 11, 2026.*today/ });
    await expect(today).toHaveAttribute('aria-current', 'date');
    await expect(args.onMonthChange).toHaveBeenCalledWith('2026-10');
  },
};

export const MarkedAndUnmarkedDays: Story = {
  play: async ({ canvas, args }) => {
    const marked = canvas.getByRole('gridcell', { name: 'October 5, 2026, 2 articles' });
    await expect(marked).toHaveAttribute('data-marked', 'true');
    const unmarked = canvas.getByRole('gridcell', { name: 'October 6, 2026' });
    await expect(unmarked).not.toHaveAttribute('data-marked');
    await userEvent.click(marked);
    await expect(args.onSelectDate).toHaveBeenCalledWith('2026-10-05');
  },
};

export const ChangeMonth: Story = {
  play: async ({ canvas, args }) => {
    await userEvent.click(canvas.getByRole('button', { name: 'Next month' }));
    await expect(canvas.getByRole('heading', { name: 'November 2026' })).toBeVisible();
    await expect(args.onMonthChange).toHaveBeenLastCalledWith('2026-11');
    await expect(
      canvas.getByRole('gridcell', { name: 'November 3, 2026, 4 articles' }),
    ).toBeVisible();
    await userEvent.click(canvas.getByRole('button', { name: 'Previous month' }));
    await userEvent.click(canvas.getByRole('button', { name: 'Previous month' }));
    await expect(canvas.getByRole('heading', { name: 'September 2026' })).toBeVisible();
  },
};

export const FollowsSelectedArticle: Story = {
  args: { selected: '2025-03-14', counts: { '2025-03-14': 1 } },
  play: async ({ canvas }) => {
    await expect(canvas.getByRole('heading', { name: 'March 2025' })).toBeVisible();
    await expect(canvas.getByRole('gridcell', { name: /March 14, 2025/ })).toHaveAttribute(
      'aria-selected',
      'true',
    );
  },
};
