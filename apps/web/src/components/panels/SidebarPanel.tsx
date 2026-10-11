import { useMatch } from 'react-router';
import { useSidebar } from '../../lib/queries.ts';
import { SidebarTree } from '../sidebar/SidebarTree.tsx';

/** Sidebar tree wired to the data source; it locates the open article. */
export function SidebarPanel() {
  const sidebar = useSidebar();
  const match = useMatch('/a/*');
  const id = (match?.params['*'] ?? '').replace(/\/+$/, '');
  if (sidebar.isPending) return <p className="p-2 text-sm text-muted">Loading…</p>;
  if (sidebar.isError)
    return <p className="p-2 text-sm text-muted">The sidebar could not be loaded.</p>;
  return <SidebarTree nodes={sidebar.data} currentPath={id ? `${id}.md` : undefined} />;
}
