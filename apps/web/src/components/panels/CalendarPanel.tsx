import { useState } from 'react';
import { useMatch, useNavigate } from 'react-router';
import { monthRange } from '../../lib/dates.ts';
import { useArticle, useCalendar } from '../../lib/queries.ts';
import { MonthCalendar } from '../calendar/MonthCalendar.tsx';

/** Calendar wired to the data source and the router. It follows the open article or day. */
export function CalendarPanel({ today }: { today?: string | undefined }) {
  const navigate = useNavigate();
  const articleMatch = useMatch('/a/*');
  const dayMatch = useMatch('/d/:date');
  const articleId = (articleMatch?.params['*'] ?? '').replace(/\/+$/, '');
  const article = useArticle(articleId);
  const selected = dayMatch?.params.date ?? (articleId ? article.data?.date : undefined);

  const [month, setMonth] = useState<string | null>(null);
  const range = month ? monthRange(month) : { from: '', to: '' };
  const calendar = useCalendar(range.from, range.to);
  const counts = Object.fromEntries((calendar.data ?? []).map((d) => [d.date, d.count]));

  return (
    <MonthCalendar
      counts={counts}
      selected={selected}
      today={today}
      onMonthChange={setMonth}
      onSelectDate={(date) => void navigate(`/d/${date}`)}
    />
  );
}
