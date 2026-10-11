/**
 * Caches one promise per key so `use()` sees a stable promise across renders.
 * A rejected promise stays cached: dropping it would make every re-render
 * start a new load and suspend forever instead of reaching the error boundary.
 */
export function promiseCache<T>() {
  const cache = new Map<string, Promise<T>>();
  return (key: string, load: () => Promise<T>): Promise<T> => {
    let p = cache.get(key);
    if (!p) {
      p = load();
      cache.set(key, p);
    }
    return p;
  };
}

/** Fetches a same-origin asset as text, failing on HTTP errors. */
export async function fetchText(url: string): Promise<string> {
  const res = await fetch(url);
  if (!res.ok) throw new Error(`HTTP ${res.status} for ${url}`);
  return res.text();
}
