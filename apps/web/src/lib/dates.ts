// Calendar dates are plain YYYY-MM-DD strings with no time zone. All arithmetic
// uses UTC so a date never shifts by a day; only "today" reads the local clock.

const pad = (n: number) => String(n).padStart(2, '0');

/** Splits "YYYY-MM" or "YYYY-MM-DD" into numbers. */
function parts(s: string): { y: number; m: number; d: number } {
  const [y = 1970, m = 1, d = 1] = s.split('-').map(Number);
  return { y, m, d };
}

/** Today's date in the browser's local time zone. */
export function localToday(now = new Date()): string {
  return `${now.getFullYear()}-${pad(now.getMonth() + 1)}-${pad(now.getDate())}`;
}

/** The month of a date: "2026-10-05" → "2026-10". */
export function monthOf(date: string): string {
  return date.slice(0, 7);
}

/** Adds n months to "YYYY-MM". */
export function addMonths(month: string, n: number): string {
  const { y, m } = parts(month);
  const d = new Date(Date.UTC(y, m - 1 + n, 1));
  return `${d.getUTCFullYear()}-${pad(d.getUTCMonth() + 1)}`;
}

/** Number of days in "YYYY-MM". */
export function daysInMonth(month: string): number {
  const { y, m } = parts(month);
  return new Date(Date.UTC(y, m, 0)).getUTCDate();
}

/** First and last date of "YYYY-MM". */
export function monthRange(month: string): { from: string; to: string } {
  return { from: `${month}-01`, to: `${month}-${pad(daysInMonth(month))}` };
}

/** Weekday of the first day of "YYYY-MM", 0 = Sunday. */
export function firstWeekday(month: string): number {
  const { y, m } = parts(month);
  return new Date(Date.UTC(y, m - 1, 1)).getUTCDay();
}

/** Formats a date for people, e.g. "October 5, 2026", without time-zone conversion. */
export function formatDate(date: string): string {
  const { y, m, d } = parts(date);
  return new Intl.DateTimeFormat('en-US', { dateStyle: 'long', timeZone: 'UTC' }).format(
    new Date(Date.UTC(y, m - 1, d)),
  );
}

/** Formats "YYYY-MM" as "October 2026". */
export function formatMonth(month: string): string {
  const { y, m } = parts(month);
  return new Intl.DateTimeFormat('en-US', {
    month: 'long',
    year: 'numeric',
    timeZone: 'UTC',
  }).format(new Date(Date.UTC(y, m - 1, 1)));
}
