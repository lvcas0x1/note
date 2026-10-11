import { z } from 'zod';

/** A calendar date in `YYYY-MM-DD` form that exists (rejects 2026-02-30). */
export const IsoDate = z.iso.date();

type SidebarNodeShape = {
  title: string;
  path?: string | undefined;
  children?: SidebarNodeShape[] | undefined;
};

const title = z.string().trim().min(1);
const path = z
  .string()
  .regex(/^[^/].*\.md$/, 'must be a content-relative path to a .md file');

/** A sidebar node: a title plus a path, children, or both. */
export const SidebarNode: z.ZodType<SidebarNodeShape> = z.union([
  z.object({
    title,
    path,
    get children() {
      return z.array(SidebarNode).optional();
    },
  }),
  z.object({
    title,
    path: path.optional(),
    get children() {
      return z.array(SidebarNode);
    },
  }),
]);
export type SidebarNode = z.infer<typeof SidebarNode>;

/** The whole `sidebar.json` file: an ordered array of nodes. */
export const Sidebar = z.array(SidebarNode.meta({ id: "SidebarNode" }));
export type Sidebar = z.infer<typeof Sidebar>;

/** Recognized article front matter keys. Unknown keys are allowed and reported by the build. */
export const FrontMatter = z.looseObject({
  title: z.string().trim().min(1).optional(),
  date: IsoDate.optional(),
  summary: z.string().optional(),
  tags: z.array(z.string().min(1)).optional(),
  draft: z.boolean().default(false),
});
export type FrontMatter = z.infer<typeof FrontMatter>;
