import { z } from 'zod';
import { IsoDate } from './index.ts';

// Response schemas for the note API. They are written by hand and checked
// against openapi/openapi.json by test/openapi-drift.test.ts. Each schema's
// meta id is the name of its OpenAPI component.

export type ApiNodeShape = { title: string; path?: string | undefined; children?: ApiNodeShape[] | undefined };

/** A sidebar node as the API returns it (component Node). */
export const ApiNode: z.ZodType<ApiNodeShape> = z
  .object({
    title: z.string(),
    path: z.string().optional(),
    get children() {
      return z.array(ApiNode).optional();
    },
  })
  .meta({ id: 'Node' });
export type ApiNode = z.infer<typeof ApiNode>;

export const SidebarResponse = z.array(ApiNode);
export type SidebarResponse = z.infer<typeof SidebarResponse>;

const summaryShape = {
  id: z.string(),
  title: z.string(),
  date: IsoDate,
  summary: z.string(),
  category: z.array(z.string()),
};

/** An article in day lists (component ArticleSummary). */
export const ArticleSummary = z.object(summaryShape).meta({ id: 'ArticleSummary' });
export type ArticleSummary = z.infer<typeof ArticleSummary>;

/** A full article (component ArticleDetail). */
export const ArticleDetail = z
  .object({ ...summaryShape, tags: z.array(z.string()), inSidebar: z.boolean(), html: z.string() })
  .meta({ id: 'ArticleDetail' });
export type ArticleDetail = z.infer<typeof ArticleDetail>;

/** Articles per day (component DayCount). */
export const DayCount = z.object({ date: IsoDate, count: z.int().min(1) }).meta({ id: 'DayCount' });
export type DayCount = z.infer<typeof DayCount>;

/** A search result; snippet is HTML with <mark> around matches (component Hit). */
export const SearchHit = z.object({ ...summaryShape, snippet: z.string() }).meta({ id: 'Hit' });
export type SearchHit = z.infer<typeof SearchHit>;

export const CalendarResponse = z.array(DayCount);
export type CalendarResponse = z.infer<typeof CalendarResponse>;
export const DayResponse = z.array(ArticleSummary);
export type DayResponse = z.infer<typeof DayResponse>;
export const SearchResponse = z.array(SearchHit);
export type SearchResponse = z.infer<typeof SearchResponse>;

/** An RFC 9457 problem body (component ErrorModel). */
export const Problem = z
  .object({
    type: z.string().optional(),
    title: z.string().optional(),
    status: z.int().optional(),
    detail: z.string().optional(),
    instance: z.string().optional(),
    errors: z
      .array(z.object({ location: z.string().optional(), message: z.string().optional(), value: z.unknown().optional() }).meta({ id: 'ErrorDetail' }))
      .optional(),
  })
  .meta({ id: 'ErrorModel' });
export type Problem = z.infer<typeof Problem>;

/** Every schema that must match an OpenAPI component, by component name. */
export const apiComponents = {
  Node: ApiNode,
  ArticleSummary,
  ArticleDetail,
  DayCount,
  Hit: SearchHit,
  ErrorModel: Problem,
} as const;
