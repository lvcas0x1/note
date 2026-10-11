import { useId, useState, type ReactNode } from 'react';
import { useLocation } from 'react-router';

export interface AppShellProps {
  /** Month calendar at the top of the sidebar. */
  calendar: ReactNode;
  /** Category tree below the calendar. */
  tree: ReactNode;
  /** Search bar pinned to the top center of the content area. */
  search?: ReactNode;
  children: ReactNode;
}

/**
 * Page layout: a left sidebar (calendar above the category tree) and the
 * content area. Below 768 px of width the sidebar collapses behind a toggle.
 * Width is measured with a container query, so the layout follows the space
 * the shell gets.
 */
export function AppShell({ calendar, tree, search, children }: AppShellProps) {
  const sidebarId = useId();
  const { pathname } = useLocation();
  // The overlay is open only on the page where it was opened, so it closes
  // after navigating on a narrow screen.
  const [openOn, setOpenOn] = useState<string | null>(null);
  const open = openOn === pathname;
  const setOpen = (next: boolean) => setOpenOn(next ? pathname : null);

  return (
    <div className="@container h-dvh w-full">
      <div className="flex h-full">
        <aside
          id={sidebarId}
          aria-label="Sidebar"
          className={`${open ? 'flex' : 'hidden'} absolute inset-y-0 left-0 z-20 w-72 shrink-0 flex-col border-r border-line bg-surface shadow-lg @min-[768px]:static @min-[768px]:flex @min-[768px]:shadow-none`}
        >
          <div className="border-b border-line p-3">{calendar}</div>
          <nav aria-label="Categories" className="min-h-0 flex-1 overflow-y-auto p-2">
            {tree}
          </nav>
        </aside>
        {open ? (
          <button
            type="button"
            aria-label="Close sidebar"
            className="absolute inset-0 z-10 bg-ink/20 @min-[768px]:hidden"
            onClick={() => setOpen(false)}
          />
        ) : null}
        <main className="relative flex min-w-0 flex-1 flex-col overflow-y-auto">
          <div className="sticky top-0 z-10 flex items-center gap-2 border-b border-line bg-page/95 px-3 py-2 backdrop-blur">
            <button
              type="button"
              className="rounded-md border border-line px-2 py-1 text-sm @min-[768px]:hidden"
              aria-controls={sidebarId}
              aria-expanded={open}
              onClick={() => setOpen(!open)}
            >
              Menu
            </button>
            <div className="mx-auto w-full max-w-xl">{search}</div>
          </div>
          <div className="w-full flex-1 px-[5%] py-6">{children}</div>
        </main>
      </div>
    </div>
  );
}
