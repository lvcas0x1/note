// Runs in UTC−10 (Pacific/Honolulu): article dates must not shift by a day.
process.env.TZ = 'Pacific/Honolulu';

import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { MonthCalendar } from '../src/components/calendar/MonthCalendar.tsx';
import { addMonths, localToday, monthRange } from '../src/lib/dates.ts';

describe('calendar in UTC−10', () => {
  it('uses the local date for today', () => {
    expect(new Date().getTimezoneOffset()).toBe(600);
    // 2026-10-12 05:00 UTC is still 2026-10-11 in Honolulu.
    expect(localToday(new Date(Date.UTC(2026, 9, 12, 5)))).toBe('2026-10-11');
  });

  it('marks article dates on the same day', () => {
    render(
      <MonthCalendar
        today="2026-10-11"
        counts={{ '2026-10-01': 1, '2026-10-31': 3 }}
        onSelectDate={() => {}}
      />,
    );
    expect(screen.getByRole('gridcell', { name: 'October 1, 2026, 1 article' })).toHaveProperty(
      'dataset.marked',
      'true',
    );
    expect(screen.getByRole('gridcell', { name: 'October 31, 2026, 3 articles' })).toHaveProperty(
      'dataset.marked',
      'true',
    );
    expect(
      screen.getByRole('gridcell', { name: 'October 2, 2026' }).dataset.marked,
    ).toBeUndefined();
  });

  it('computes month ranges without shifting', () => {
    expect(monthRange('2026-02')).toEqual({ from: '2026-02-01', to: '2026-02-28' });
    expect(addMonths('2026-12', 1)).toBe('2027-01');
    expect(addMonths('2026-01', -1)).toBe('2025-12');
  });
});
