import type { ApiNode } from '@notes/schemas/api';

/** A tree node with its position: key is the index path, e.g. "0.1.2". */
export interface FlatNode {
  key: string;
  node: ApiNode;
  level: number;
  parentKey: string | null;
  hasChildren: boolean;
}

/** Lists the nodes a reader can see: every node whose ancestors are all expanded. */
export function visibleNodes(nodes: ApiNode[], expanded: ReadonlySet<string>): FlatNode[] {
  const out: FlatNode[] = [];
  const walk = (list: ApiNode[], level: number, parentKey: string | null) => {
    list.forEach((node, i) => {
      const key = parentKey === null ? String(i) : `${parentKey}.${i}`;
      const hasChildren = (node.children?.length ?? 0) > 0;
      out.push({ key, node, level, parentKey, hasChildren });
      if (hasChildren && expanded.has(key)) walk(node.children ?? [], level + 1, key);
    });
  };
  walk(nodes, 1, null);
  return out;
}

/** Returns the key of the node whose path is `path`, and the keys of its ancestors. */
export function locate(
  nodes: ApiNode[],
  path: string,
): { key: string; ancestors: string[] } | null {
  const walk = (
    list: ApiNode[],
    parents: string[],
  ): { key: string; ancestors: string[] } | null => {
    for (const [i, node] of list.entries()) {
      const key = parents.length === 0 ? String(i) : `${parents.at(-1)}.${i}`;
      if (node.path === path) return { key, ancestors: parents };
      const found = walk(node.children ?? [], [...parents, key]);
      if (found) return found;
    }
    return null;
  };
  return path ? walk(nodes, []) : null;
}

/** Route for an article path: "aws/s3/s3.md" → "/a/aws/s3/s3". */
export function articleRoute(path: string): string {
  return `/a/${path.replace(/\.md$/i, '')}`;
}
