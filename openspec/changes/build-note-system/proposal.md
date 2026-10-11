# Proposal

## Why

The current notes site (`~/Documents/github/docsify`) is a docsify app with a
client-side SHA-256 password gate, a Markdown sidebar, and many vendored plugins.
It has no way to browse notes by date, and diagrams, charts, and API docs come
from separate plugins with different syntaxes. This change replaces it with a
note system that can later be published as a public, read-only blog. The new
system has a JSON sidebar, a date calendar, one embed syntax, and server-side
search.

## What Changes

- New greenfield monorepo in this repository: a Go server and build CLI, and a
  React + Vite web app.
- No authentication, user accounts, or login page. The site is read-only for
  every visitor. Requirements 1 (OAuth login) and 18 (Better Auth) from the
  original request are withdrawn.
- Built for public exposure:
  - Articles marked `draft: true` never leave the build.
  - The server listens on loopback by default and is meant to sit behind a
    CDN or reverse proxy.
  - The search API limits query length and rate.
  - The site loads no third-party scripts or content at runtime.
- Content is Markdown files plus a nested `sidebar.json` in a content folder.
  A build command validates them, compiles Markdown to HTML, and indexes
  articles in SQLite.
- Every article has a `date`. A calendar in the top-left corner marks dates
  that have articles. Clicking a date lists the titles and summaries for that
  day. Clicking a title opens the article.
- A nested, collapsible left sidebar locates and highlights the open article.
- Read-only embeds in Markdown:
  - draw.io: `![drawio](...)`
  - Excalidraw: `![excalidraw](...)`
  - OpenAPI/Swagger: `![openapi](...)`
  - MP4 video: `![mp4](...)`
  - Mermaid: ` ```mermaid ` fences
  - Apache ECharts: ` ```echarts ` fences (optional)

  Remote diagram and API-doc URLs are downloaded at build time and served
  locally. Videos are local files.
- A search bar, always visible at the top center of the content area, runs
  keyword full-text search (SQLite FTS5). LLM search is deferred.
- An OpenAPI 3.1 document for the note API. The web app's typed client and Zod
  response schemas are checked against it.
- Supporting tooling:
  - Storybook for every UI component.
  - Vitest for tests.
  - Pino for logging in the web app.
  - Faker for test data.
  - Tailwind CSS, light theme only.
- Two ways to publish from the same build:
  - **Server:** the Go server, behind a CDN or reverse proxy.
  - **Static:** a static export for hosts such as GitHub Pages, with
    in-browser search. The static export writes videos to a separate folder
    for manual upload to external storage, served from a configured base URL.
    The Markdown source stays out of the published repository.
- A migration command converts the docsify `_sidebar.md` into `sidebar.json`,
  so the existing docsify notes can be used to verify the system.

## Capabilities

### New Capabilities

- `content-build`: the content folder layout, the `sidebar.json` and front
  matter schemas, draft exclusion, validation, Markdown-to-HTML compilation,
  embed detection, asset collection (including build-time download of remote
  embeds), and indexing.
- `sidebar-navigation`: the nested sidebar, opening articles from it, and
  locating the current article.
- `calendar`: the date calendar, the per-day article list with summaries, and
  opening an article from the list.
- `article-rendering`: displaying a compiled article and rendering its draw.io,
  Excalidraw, OpenAPI, Mermaid, ECharts, and MP4 video embeds read-only.
- `search`: the persistent search bar and keyword full-text search.
- `static-export`: the static copy of the published site for static hosts:
  data files, in-browser search, direct links, base path, external video URLs,
  size limits, and the meta-tag security policy.
- `note-api`: the HTTP API contract, its published OpenAPI document, asset and
  video delivery, and the protections for public exposure.

### Modified Capabilities

None. The project has no existing specs.

## Impact

- New code:
  - `server/` (Go module)
  - `apps/web` (TypeScript/React)
  - `packages/schemas` (shared Zod schemas)
  - Root workspace and tooling configuration
- New runtime dependencies: Go 1.27 and SQLite (embedded, one content
  database per build). Node.js 24 LTS and npm are needed only to build the web
  app and to run the dev tooling.
- Go server stack: Gin (HTTP framework) with Huma (OpenAPI), and GORM (ORM)
  on SQLite.
- The existing docsify site is not modified. It is used read-only as test
  content.
- Out of scope:
  - Authentication and user management
  - LLM or semantic search
  - DuckDB and Qdrant
  - An in-browser editor or any diagram editing UI
  - Comments and other user-generated content
  - A dark theme
  - Production deployment automation, and any cloud configuration (for example
    the S3 and CloudFront setup for videos)
