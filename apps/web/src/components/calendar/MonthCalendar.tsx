import { useState } from 'react';
import {
  addMonths,
  daysInMonth,
  firstWeekday,
  formatDate,
  formatMonth,
  localToday,
  monthOf,
} from '../../lib/dates.ts';

export interface MonthCalendarProps {
  /** Article counts by date (YYYY-MM-DD) for the month shown. */
  counts: Readonly<Record<string, number>>;
  /** Highlighted date, e.g. the open article's date. The calendar shows its month. */
  selected?: string | undefined;
  /** Today's date; defaults to the browser's local date. */
  today?: string | undefined;
  /** Called with "YYYY-MM" whenever the shown month changes, including the first render. */
  onMonthChange?: ((month: string) => void) | undefined;
  onSelectDate: (date: string) => void;
}

const weekdays = ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa'];

/** Month calendar that marks days with articles. */
export function MonthCalendar({
  counts,
  selected,
  today,
  onMonthChange,
  onSelectDate,
}: MonthCalendarProps) {
  const todayDate = today ?? localToday();
  const [month, setMonth] = useState(() => monthOf(selected ?? todayDate));

  // Follow the selected date to its month (state derived from props during render).
  const [followed, setFollowed] = useState(selected);
  if (selected !== followed) {
    setFollowed(selected);
    if (selected && monthOf(selected) !== month) setMonth(monthOf(selected));
  }

  // Report the shown month to the parent so it can load that month's counts.
  const [reported, setReported] = useState<string | null>(null);
  if (reported !== month) {
    setReported(month);
    onMonthChange?.(month);
  }

  const lead = firstWeekday(month);
  const total = daysInMonth(month);
  const cells: (string | null)[] = [
    ...Array.from({ length: lead }, () => null),
    ...Array.from({ length: total }, (_d, i) => `${month}-${String(i + 1).padStart(2, '0')}`),
  ];

  return (
    <div className="select-none text-sm" aria-label="Calendar" role="group">
      <div className="mb-2 flex items-center justify-between">
        <button
          type="button"
          className="rounded px-2 py-1 hover:bg-line"
          aria-label="Previous month"
          onClick={() => setMonth(addMonths(month, -1))}
        >
          ‹
        </button>
        <h2 className="font-semibold" aria-live="polite">
          {formatMonth(month)}
        </h2>
        <button
          type="button"
          className="rounded px-2 py-1 hover:bg-line"
          aria-label="Next month"
          onClick={() => setMonth(addMonths(month, 1))}
        >
          ›
        </button>
      </div>
      <div
        className="grid grid-cols-7 gap-0.5 text-center"
        role="grid"
        aria-label={formatMonth(month)}
      >
        {weekdays.map((d) => (
          <div key={d} role="columnheader" className="py-1 text-xs text-muted">
            {d}
          </div>
        ))}
        {cells.map((date, i) => {
          if (!date) return <div key={`pad-${i}`} role="presentation" />;
          const count = counts[date] ?? 0;
          const isToday = date === todayDate;
          const isSelected = date === selected;
          const label = `${formatDate(date)}${count > 0 ? `, ${count} ${count === 1 ? 'article' : 'articles'}` : ''}${isToday ? ', today' : ''}`;
          return (
            <button
              key={date}
              type="button"
              role="gridcell"
              aria-label={label}
              aria-selected={isSelected}
              aria-current={isToday ? 'date' : undefined}
              data-marked={count > 0 || undefined}
              onClick={() => onSelectDate(date)}
              className={`relative rounded py-1 tabular-nums ${isSelected ? 'bg-accent text-white' : count > 0 ? 'font-semibold text-ink hover:bg-accent-soft' : 'text-muted hover:bg-line'} ${isToday && !isSelected ? 'ring-1 ring-accent' : ''}`}
            >
              {Number(date.slice(8))}
              {count > 0 ? (
                <span
                  aria-hidden="true"
                  className={`absolute bottom-0.5 left-1/2 size-1 -translate-x-1/2 rounded-full ${isSelected ? 'bg-white' : 'bg-accent'}`}
                />
              ) : null}
            </button>
          );
        })}
      </div>
    </div>
  );
}
