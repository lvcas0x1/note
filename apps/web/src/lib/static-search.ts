import { SearchResponse, type SearchHit } from '@notes/schemas/api';
import { ApiError, SchemaError } from './api.ts';

let worker: Worker | undefined;
let nextId = 0;
const pending = new Map<number, { resolve: (v: unknown) => void; reject: (e: unknown) => void }>();

function getWorker(): Worker {
  if (!worker) {
    worker = new Worker(new URL('./search.worker.ts', import.meta.url), { type: 'module' });
    worker.addEventListener(
      'message',
      (e: MessageEvent<{ id: number; hits?: unknown; error?: string; reason?: string }>) => {
        const p = pending.get(e.data.id);
        if (!p) return;
        pending.delete(e.data.id);
        if (e.data.error !== undefined) {
          // Rejected queries behave like the server's 400 problem responses.
          p.reject(
            e.data.reason
              ? new ApiError(400, { status: 400, title: 'Bad Request', detail: e.data.error })
              : new Error(e.data.error),
          );
        } else {
          p.resolve(e.data.hits);
        }
      },
    );
  }
  return worker;
}

/** Searches a static site in the browser. The index downloads on the first search. */
export async function searchStatic(
  basePath: string,
  q: string,
  signal?: AbortSignal,
): Promise<SearchHit[]> {
  const id = ++nextId;
  const raw = await new Promise((resolve, reject) => {
    pending.set(id, { resolve, reject });
    signal?.addEventListener('abort', () => {
      pending.delete(id);
      reject(signal.reason);
    });
    // Worker.postMessage has no target origin.
    // oxlint-disable-next-line unicorn/require-post-message-target-origin
    getWorker().postMessage({ id, basePath, q });
  });
  const parsed = SearchResponse.safeParse(raw);
  if (!parsed.success) throw new SchemaError('search.db', parsed.error);
  return parsed.data;
}
