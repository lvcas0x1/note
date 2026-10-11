import { createContext, useContext, type ReactNode } from 'react';
import { apiSource, type NotesSource } from './source.ts';

const SourceContext = createContext<NotesSource>(apiSource);

export function NotesSourceProvider({
  source,
  children,
}: {
  source: NotesSource;
  children: ReactNode;
}) {
  return <SourceContext value={source}>{children}</SourceContext>;
}

export function useNotesSource(): NotesSource {
  return useContext(SourceContext);
}
