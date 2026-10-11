import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useState, type ReactNode } from 'react';
import { Outlet, type RouteObject } from 'react-router';
import { AppShell } from '../components/layout/AppShell.tsx';
import { CalendarPanel } from '../components/panels/CalendarPanel.tsx';
import { SidebarPanel } from '../components/panels/SidebarPanel.tsx';
import { SearchBar } from '../components/search/SearchBar.tsx';
import { NotesSourceProvider } from '../lib/source-context.tsx';
import type { NotesSource } from '../lib/source.ts';
import { ArticlePage } from '../pages/ArticlePage.tsx';
import { DayPage } from '../pages/DayPage.tsx';
import { HomePage } from '../pages/HomePage.tsx';
import { NotFoundPage } from '../pages/NotFoundPage.tsx';

/** Data providers shared by the app and its stories. */
export function AppProviders({ source, children }: { source: NotesSource; children: ReactNode }) {
  const [client] = useState(
    () => new QueryClient({ defaultOptions: { queries: { staleTime: 30_000 } } }),
  );
  return (
    <QueryClientProvider client={client}>
      <NotesSourceProvider source={source}>{children}</NotesSourceProvider>
    </QueryClientProvider>
  );
}

export interface LayoutOptions {
  today?: string;
  search?: ReactNode;
}

function Layout({ today, search }: LayoutOptions) {
  return (
    <AppShell
      calendar={<CalendarPanel today={today} />}
      tree={<SidebarPanel />}
      search={search ?? <SearchBar />}
    >
      <Outlet />
    </AppShell>
  );
}

/** The route table: start page, articles (/a/<id>), days (/d/<date>), and not found. */
export function appRoutes(options: LayoutOptions = {}): RouteObject[] {
  return [
    {
      element: <Layout {...options} />,
      children: [
        { index: true, element: <HomePage /> },
        { path: 'a/*', element: <ArticlePage /> },
        { path: 'd/:date', element: <DayPage /> },
        { path: '*', element: <NotFoundPage /> },
      ],
    },
  ];
}
