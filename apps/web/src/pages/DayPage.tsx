import { Link, useParams } from 'react-router';
import { Placeholder } from '../components/Placeholder.tsx';
import { formatDate } from '../lib/dates.ts';
import { useDay } from '../lib/queries.ts';

/** Lists every article of one day with its title and summary, sorted by title. */
export function DayPage() {
  const date = useParams().date ?? '';
  const valid = /^\d{4}-\d{2}-\d{2}$/.test(date);
  const day = useDay(valid ? date : '');
  if (!valid)
    return <Placeholder title="Not a valid date">Dates look like 2026-10-05.</Placeholder>;

  return (
    <section aria-labelledby="day-title">
      <h1 id="day-title" className="text-2xl font-semibold">
        {formatDate(date)}
      </h1>
      {day.isPending ? (
        <p className="mt-4 text-muted">Loading…</p>
      ) : day.isError ? (
        <div className="mt-4">
          <Placeholder title="Could not load this day">Try again later.</Placeholder>
        </div>
      ) : day.data.length === 0 ? (
        <div className="mt-4">
          <Placeholder title="No articles on this day" />
        </div>
      ) : (
        <ul className="mt-4 divide-y divide-line">
          {day.data
            .toSorted((a, b) => a.title.localeCompare(b.title, undefined, { sensitivity: 'base' }))
            .map((a) => (
              <li key={a.id} className="py-3">
                <Link to={`/a/${a.id}`} className="text-lg font-medium text-accent hover:underline">
                  {a.title}
                </Link>
                {a.summary ? <p className="mt-1 text-muted">{a.summary}</p> : null}
              </li>
            ))}
        </ul>
      )}
    </section>
  );
}
