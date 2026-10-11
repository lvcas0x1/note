import type { Meta, StoryObj } from '@storybook/react-vite';
import { expect } from 'storybook/test';
import { seeded } from '../test/fixtures.ts';
import { Placeholder } from './Placeholder.tsx';

const meta = {
  component: Placeholder,
  args: { title: 'Nothing here yet', children: seeded().lorem.sentence() },
} satisfies Meta<typeof Placeholder>;
export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  play: async ({ canvas }) => {
    await expect(canvas.getByText('Nothing here yet')).toBeVisible();
  },
};
