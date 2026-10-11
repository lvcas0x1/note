import type { ApiNode } from '@notes/schemas/api';
import { useEffect, useMemo, useRef, useState, type KeyboardEvent } from 'react';
import { Link, useNavigate } from 'react-router';
import { articleRoute, locate, visibleNodes, type FlatNode } from './tree.ts';

export interface SidebarTreeProps {
  nodes: ApiNode[];
  /** Content path of the open article (e.g. "aws/s3/s3.md"); its node is located and highlighted. */
  currentPath?: string | undefined;
}

/**
 * Nested category tree (WAI-ARIA tree pattern). Categories expand and
 * collapse; nodes with a path link to their article. The open article's
 * ancestors are expanded, its node highlighted and scrolled into view.
 */
export function SidebarTree({ nodes, currentPath = '' }: SidebarTreeProps) {
  const navigate = useNavigate();
  const current = useMemo(() => locate(nodes, currentPath), [nodes, currentPath]);
  const [expanded, setExpanded] = useState<ReadonlySet<string>>(() => new Set(current?.ancestors));
  const [focusKey, setFocusKey] = useState<string | null>(null);

  // Expand the ancestors when another article opens (state derived from props during render).
  const [locatedFor, setLocatedFor] = useState(current?.key);
  if (current?.key !== locatedFor) {
    setLocatedFor(current?.key);
    if (current) setExpanded((prev) => new Set([...prev, ...current.ancestors]));
  }

  const items = visibleNodes(nodes, expanded);
  const activeKey =
    focusKey && items.some((i) => i.key === focusKey) ? focusKey : (current?.key ?? items[0]?.key);
  const itemRefs = useRef(new Map<string, HTMLLIElement>());

  // Scroll the open article's node into view (the DOM is the external system here).
  useEffect(() => {
    if (current) itemRefs.current.get(current.key)?.scrollIntoView({ block: 'nearest' });
  }, [current]);

  const toggle = (key: string, open?: boolean) =>
    setExpanded((prev) => {
      const next = new Set(prev);
      if (open ?? !next.has(key)) next.add(key);
      else next.delete(key);
      return next;
    });

  const focus = (key: string | null | undefined) => {
    if (!key) return;
    setFocusKey(key);
    itemRefs.current.get(key)?.focus();
  };

  const onKeyDown = (e: KeyboardEvent, item: FlatNode) => {
    const index = items.findIndex((i) => i.key === item.key);
    const isOpen = expanded.has(item.key);
    switch (e.key) {
      case 'ArrowDown':
        focus(items[index + 1]?.key);
        break;
      case 'ArrowUp':
        focus(items[index - 1]?.key);
        break;
      case 'Home':
        focus(items[0]?.key);
        break;
      case 'End':
        focus(items.at(-1)?.key);
        break;
      case 'ArrowRight':
        if (item.hasChildren && !isOpen) toggle(item.key, true);
        else if (item.hasChildren) focus(items[index + 1]?.key);
        break;
      case 'ArrowLeft':
        if (item.hasChildren && isOpen) toggle(item.key, false);
        else focus(item.parentKey);
        break;
      case 'Enter':
      case ' ':
        if (item.node.path) void navigate(articleRoute(item.node.path));
        else if (item.hasChildren) toggle(item.key);
        break;
      default:
        return;
    }
    e.preventDefault();
  };

  return (
    <ul role="tree" aria-label="Categories" className="space-y-0.5 text-sm">
      {items.map((item) => {
        const isCurrent = item.key === current?.key;
        const isOpen = expanded.has(item.key);
        return (
          <li
            key={item.key}
            ref={(el) => {
              if (el) itemRefs.current.set(item.key, el);
              else itemRefs.current.delete(item.key);
            }}
            role="treeitem"
            aria-level={item.level}
            aria-expanded={item.hasChildren ? isOpen : undefined}
            aria-selected={isCurrent}
            tabIndex={item.key === activeKey ? 0 : -1}
            onKeyDown={(e) => onKeyDown(e, item)}
            onFocus={(e) => {
              if (e.target === e.currentTarget) setFocusKey(item.key);
            }}
            className="rounded-md outline-none focus-visible:ring-2 focus-visible:ring-accent"
          >
            <div
              className={`flex items-center gap-1 rounded-md py-1 pr-2 ${isCurrent ? 'bg-accent-soft font-medium text-ink shadow-[inset_3px_0_0_var(--color-accent)]' : 'hover:bg-line/60'}`}
              style={{ paddingLeft: `${(item.level - 1) * 0.875 + 0.25}rem` }}
            >
              {item.hasChildren ? (
                <button
                  type="button"
                  tabIndex={-1}
                  aria-hidden="true"
                  className="w-4 shrink-0 text-muted"
                  onClick={() => toggle(item.key)}
                >
                  {isOpen ? '▾' : '▸'}
                </button>
              ) : (
                <span className="w-4 shrink-0" aria-hidden="true" />
              )}
              {item.node.path ? (
                <Link
                  to={articleRoute(item.node.path)}
                  tabIndex={-1}
                  aria-current={isCurrent ? 'page' : undefined}
                  className="min-w-0 flex-1 truncate"
                >
                  {item.node.title}
                </Link>
              ) : (
                <span
                  className="min-w-0 flex-1 cursor-pointer truncate"
                  onClick={() => toggle(item.key)}
                >
                  {item.node.title}
                </span>
              )}
            </div>
          </li>
        );
      })}
    </ul>
  );
}
