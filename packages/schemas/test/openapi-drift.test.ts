import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';
import { z } from 'zod';
import { apiComponents } from '../src/api.ts';

// Compares each hand-written response schema with its component in
// openapi/openapi.json: property names, required fields, types, date formats,
// and array item references.

type Json = Record<string, unknown>;
const openapi = JSON.parse(
  readFileSync(new URL('../../../openapi/openapi.json', import.meta.url), 'utf8'),
) as { components: { schemas: Record<string, Json> } };

function asJson(v: unknown): Json {
  return typeof v === 'object' && v !== null && !Array.isArray(v) ? (v as Json) : {};
}

function refName(ref: unknown): string | undefined {
  return typeof ref === 'string' ? ref.split('/').pop() : undefined;
}

function shapeOf(schema: Json, self: string): unknown {
  const s: Json = schema.$ref ? { ref: refName(schema.$ref) } : schema;
  if (s.ref) return { ref: s.ref };
  const type = Array.isArray(s.type) ? s.type.filter((t) => t !== 'null')[0] : s.type;
  if (type === 'array') return { type, items: shapeOf(asJson(s.items), self) };
  if (type === 'object' && s.properties) {
    const props = asJson(s.properties);
    return {
      type,
      required: [...((s.required as string[] | undefined) ?? [])].toSorted(),
      properties: Object.fromEntries(
        Object.keys(props)
          .toSorted()
          .map((k) => [k, shapeOf(asJson(props[k]), self)]),
      ),
    };
  }
  return { type: type ?? 'any', ...(s.format === 'date' ? { format: 'date' } : {}) };
}

function zodComponent(name: string, schema: z.ZodType): Json {
  const json = asJson(z.toJSONSchema(schema, { target: 'draft-2020-12', io: 'output', unrepresentable: 'any' }));
  const defs = asJson(json.$defs);
  // A schema with a meta id is emitted as a $ref to its own $def.
  return asJson(json.$ref ? defs[refName(json.$ref) ?? ''] : (defs[name] ?? json));
}

describe('API response schemas match openapi/openapi.json', () => {
  it.each(Object.entries(apiComponents))('%s', (name, schema) => {
    const component = openapi.components.schemas[name];
    expect(component).toBeDefined();
    expect(shapeOf(zodComponent(name, schema), name)).toEqual(shapeOf(asJson(component), name));
  });
});
