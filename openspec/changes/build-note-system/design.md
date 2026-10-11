# Design

## Context

- The repository is empty apart from `openspec/`. It is not a git repository.
- The test content is `~/Documents/github/docsify`, used read-only:
  - 67 Markdown files and a nested `_sidebar.md`.
  - No front matter, so every article falls back to its modification date.
  - 22 `mermaid` fences, Swagger 2.0 JSON, and one draw.io file
    (`_plugin/vendor/drawio/test.drawio`).
  - The site sits behind a password today, so it may contain notes that were
    never meant to be public.
- Decisions the user has made:
  - Go backend.
  - No authentication. The site may later be published as a public blog.
  - Keyword search only. LLM search is deferred.
  - Content comes from files plus a build step. There is no editor, and embeds
    are view-only.
  - Videos are local files, with no app-imposed size limit.
  - TypeScript 7 only, with Oxlint.
  - Gin and GORM.
  - Deployment: A+B. The Go server behind a CDN or reverse proxy (A), plus a
    static export for GitHub Pages (B).
  - In the static export, videos are hosted outside the site (the user plans
    S3 + CloudFront) and uploaded by hand. Cloud configuration is out of scope
    for this project.
  - The Markdown source stays local or in a private repository. Only the
    export output is pushed to the public Pages repository.

## Goals / Non-Goals

**Goals:**
- One origin. The browser talks only to the site's origin and loads nothing
  from third parties at runtime. The only exception is static-export videos,
  which come from the configured video origin.
- One build, two delivery modes (server and static), with the same UI code and
  the same search results.
- Safe to expose publicly:
  - Read-only.
  - No cookies.
  - Drafts never leave the build.
  - Bounded search cost.
  - Strict security headers.
- Content builds are reproducible, fully validated, and swapped in atomically.
- One source of truth for each contract:
  - Zod for the content file formats.
  - Go code for the API, which generates OpenAPI. The web app's Zod response
    schemas are written by hand and a test checks them against OpenAPI.

**Non-Goals:**
- Horizontal scaling. One host, one Go process.
- Terminating TLS in Go. TLS ends at the CDN or reverse proxy.
- Hot reload of content on file save. The operator runs `notes build`. A
  `--watch` flag can come later.
- Server-side rendering of diagrams to images.

## Decisions

### D1. Topology: single Go process behind a CDN or reverse proxy

```
Reader ──HTTPS──▶ CDN / reverse proxy (TLS, caching)
                    ──HTTP──▶ Go `notes serve` 127.0.0.1:8080 (Gin engine)
                               ├─ /api/v1/*   → Huma handlers via humagin
                               ├─ /assets/*   → current build's assets (Range support)
                               ├─ /healthz
                               └─ NoRoute     → embedded SPA (apps/web/dist), history fallback
```

- **Framework:** Gin v1.12.0, which requires Go 1.25+
  ([Gin docs](https://gin-gonic.com/en/docs/)). Gin's default logger and
  recovery middleware are replaced with the `internal/logx` slog logger (D9)
  and a recovery handler that returns problem+json.
- **Listen address:** `127.0.0.1:8080` by default (note-api: loopback by
  default). To expose the server, set `NOTES_LISTEN` explicitly, or put a
  reverse proxy on the same host.
- **Client IP:** `engine.SetTrustedProxies` is set from `NOTES_TRUSTED_PROXIES`
  (empty by default). Forwarding headers are honored only from those proxies.
  The IP-based search rate limit depends on this.
- **Timeouts:** `http.Server` sets `ReadHeaderTimeout` and `IdleTimeout`. It
  leaves `WriteTimeout` unset, so long video transfers are never cut off
  (note-api: asset delivery). Slow-client protection for large downloads is
  left to the CDN or reverse proxy.
- **CDN caching:**
  - `/assets/*` uses content-hashed paths with
    `Cache-Control: public, max-age=31536000, immutable`.
  - API responses get a short `max-age`.
  - Cloudflare caches files up to 512 MB on the Free, Pro, and Business plans,
    and 5 GB on Enterprise. `.mp4` is cached by default
    ([Cloudflare default cache behavior](https://developers.cloudflare.com/cache/concepts/default-cache-behavior/)).
    Larger videos are therefore served from the origin through the CDN. Range
    requests make that workable.
- **No authentication:** no cookies and no visitor data are stored. The
  original requirements for OAuth login and Better Auth are withdrawn.

### D2. Public-exposure protections

- **Drafts:** articles with `draft: true` are dropped at build time (D3), not
  hidden at runtime, so the served database cannot leak them.
- **Search abuse limits:**
  - Queries longer than 200 characters return 400.
  - Requests are rate-limited per IP with a token bucket from
    `golang.org/x/time/rate` (v0.16.0). The default is 60 requests per
    minute, set by `NOTES_SEARCH_RATE`. The limiter map is bounded and evicts
    idle entries.
  - Each query runs with a 2 s context timeout.
- **Security headers (Gin middleware):**
  - CSP:

    ```
    default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline';
    img-src 'self' data: blob:; media-src 'self'; connect-src 'self';
    font-src 'self'; object-src 'none'; frame-ancestors 'none'; base-uri 'self'
    ```

    `'unsafe-inline'` for styles is needed because Mermaid, ECharts, and
    swagger-ui inject inline styles. This is my assessment and has not been
    verified against each library's docs. Task 9.2 confirms it from CSP
    violation reports.
  - `Strict-Transport-Security: max-age=31536000`
  - `X-Content-Type-Options: nosniff`
  - `Referrer-Policy: strict-origin-when-cross-origin`
- **Write methods:** only `GET` and `HEAD` are routed. Every other method
  returns 405.

### D3. Content pipeline in Go (`notes build`)

1. Validate `sidebar.json` against `schemas/sidebar.schema.json`. That file is
   generated from the Zod schema with `z.toJSONSchema()`, draft 2020-12
   ([Zod JSON Schema](https://zod.dev/json-schema)), and checked in. Go embeds
   it and validates with `santhosh-tekuri/jsonschema/v6`, which reports the
   JSON pointers that the error messages need. Front matter is validated in Go
   against the same generated schema.
2. Parse each Markdown file with goldmark, using the GFM and Footnote
   extensions and `go.abhg.dev/goldmark/frontmatter`. Drop `draft: true`
   articles, along with their sidebar nodes and any categories left empty.
   Apply the metadata fallbacks.
3. A goldmark AST transformer rewrites embed images and fences into
   placeholders:
   - `<div data-embed="drawio|excalidraw|openapi|mp4" data-src="/assets/…">`
   - `<pre data-embed="mermaid|echarts">escaped source</pre>`

   Other code blocks are highlighted at build time with chroma (class-based
   CSS). `goldmark-highlighting` is not used because its latest release is a
   2023 pseudo-version.
4. Resolve every asset:
   - **Local files:** checked to exist inside the content root. Symlinks are
     resolved before the check. There is no size limit.
   - **Remote `drawio`, `excalidraw`, and `openapi` URLs:** downloaded with
     `net/http`. Only `http`/`https` are allowed, with a 30 s timeout and a
     10 MB cap enforced through `io.LimitReader`. A non-2xx response fails
     the build. The downloaded file is stored as a local asset.
   - **`mp4`:** must be a local `.mp4` file.
   - **Collection:** only assets referenced by published articles are
     collected. Assets are copied under content-hashed names
     (`assets/<sha256-prefix>/<basename>`).
5. Sanitize the HTML with bluemonday's UGC policy, extended to allow
   `data-embed`, `data-src`, and the chroma classes.
6. Write `data/builds/<UTC timestamp>/content.db` and `assets/`, then
   atomically replace the `data/current` symlink with `rename(2)`. The server
   watches `data/current` with fsnotify and reopens a read-only pool. A
   failed build never touches `current`. Only the newest build is kept.

- **Article ID:** the content-relative path without `.md`, for example
  `aws/s3/s3`. Web routes are `/a/<id>` for articles and `/d/<YYYY-MM-DD>` for
  days.
- **Remote embeds:** downloading them at build time rather than loading them at
  runtime means readers never contact third-party hosts. It also makes the
  build reproducible, and a dead URL fails the build instead of breaking the
  live page.

### D4. Storage and search: SQLite through GORM

- **Driver stack:** GORM v1.31.2 with the official dialector
  `gorm.io/driver/sqlite` v1.6.0, configured as
  `sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: …})`. The
  `database/sql` driver underneath is the pure-Go `modernc.org/sqlite`
  v1.60.1 ([pkg.go.dev](https://pkg.go.dev/gorm.io/driver/sqlite)).
- **Local spike** (`CGO_ENABLED=0`): `AutoMigrate`, GORM queries, and an
  FTS5 trigram table queried through `db.Raw` all worked, including the
  Japanese substring `本語の`. The `modernc.org/sqlite` docs do not mention
  FTS5, so the spike is the only evidence.
- **Schema:** each build writes a new database, so `AutoMigrate` on the
  `Article`, `SidebarDoc`, and `Asset` models is enough. The server opens it
  read-only (`mode=ro`).
- **FTS:**
  - The table is
    `articles_fts(title, summary, tags, body, tokenize='trigram')`, created
    with `db.Exec`. The trigram tokenizer is why the search spec has a
    three-character minimum.
  - Search uses `db.Raw` with bound parameters only.
  - Ranking uses `bm25(articles_fts, 10.0, 3.0, 3.0, 1.0)`, and snippets come
    from `snippet()`.
  - Each user term is wrapped in double quotes, with embedded quotes doubled,
    so FTS5 operators are never interpreted.
- **Rejected alternatives:**
  - `glebarez/sqlite`: its latest release (v1.11.0, 2024-03) pins
    `modernc.org/sqlite` v1.23.1
    ([GORM docs](https://gorm.io/docs/connecting_to_the_database.html)).
  - `mattn/go-sqlite3`: needs cgo.
- **DuckDB and Qdrant:** not used. There is no analytical or vector workload.

### D5. API: Huma v2 on Gin

- Gin has no built-in OpenAPI generation. Handlers are therefore registered
  through Huma v2.39.1 with the `humagin` adapter, which produces OpenAPI 3.1
  and RFC 9457 problem+json responses
  ([routers](https://huma.rocks/features/bring-your-own-router/),
  [errors](https://huma.rocks/features/response-errors/)).
- **Local spike:** Huma returns **422** for validation errors. The API
  overrides `huma.NewError` so that these return **400**, as the spec
  requires.
- **Endpoints:**
  - `GET /api/v1/sidebar`
  - `GET /api/v1/articles/{id}` (wildcard)
  - `GET /api/v1/calendar?from&to`
  - `GET /api/v1/days/{date}`
  - `GET /api/v1/search?q&limit`
  - `GET /api/v1/openapi.json`
  - `GET /healthz`

  The OpenAPI document declares no security scheme.
- **Assets:** `/assets/*` is a plain Gin route, not part of the API document.
  It is served with `http.ServeContent`, which handles `Range` and
  `If-Range` ([net/http](https://pkg.go.dev/net/http#ServeContent)). The
  route serves only paths listed in the current build's `assets` table, so
  any other path returns 404 and no filesystem traversal is possible.
- **Drift check:** `notes openapi > openapi/openapi.json` writes the
  document, which is checked in. `npm run check` regenerates it and fails on a
  diff.
- **Rejected alternatives:**
  - Comment-annotation generators: they can drift from the code.
  - oapi-codegen: two sources to keep in sync.

### D6. Web app

- **Stack:** React 19, Vite 8, `@vitejs/plugin-react`, React Router 8 in SPA
  mode, and TanStack Query 5.
- **API client:** hand-written. Zod response schemas live in
  `packages/schemas` (`api.ts`), and a small fetch client (5 endpoints)
  parses every response with them.
  - A Vitest test converts each schema with `z.toJSONSchema()` and compares it
    with the matching component in `openapi/openapi.json` (properties, types,
    required fields, formats). `npm run check` runs it, so drift fails the
    check.
  - Why not generate the client: `@hey-api/openapi-ts` 0.99.0 builds its
    output with the TypeScript compiler API (`ts.SyntaxKind`), which
    TypeScript 7.0 does not ship. Running it under TS 7 fails with
    `TypeError: Cannot read properties of undefined (reading 'AnyKeyword')`.
    Hey API declares `typescript` as a peer dependency, so an npm
    `overrides` entry cannot give it its own TS 6. Both were checked locally.
- **Styling:** Tailwind CSS 4 through `@tailwindcss/vite`
  ([Tailwind + Vite](https://tailwindcss.com/docs/installation/using-vite)).
  Light theme only: `color-scheme: light`, and no `dark:` variants (checked by
  a lint script, D8).
- **Palette:** Catppuccin Latte ([palette](https://catppuccin.com/palette/),
  v1.8.0), with mauve as the accent.
  - The tokens are in `styles/index.css` and `src/lib/palette.ts`. A test
    keeps the two equal and checks that every text/background pair used is
    at least 4.5:1 (WCAG AA).
  - The page background is `#f4f5f8` (Latte base mixed 30% with white),
    which matches the owner's terminal: Catppuccin Latte at 90% opacity.
    The sidebar and code blocks use Latte base `#eff1f5`.
  - Mauve text is used only on the page (4.97:1) and base (4.79:1)
    backgrounds. On tinted backgrounds, text uses the base text color.
  - The code-block Copy control is an icon button (copy, then a check when
    copied). It has an `aria-label` and a tooltip.
  - Code highlighting uses `catppuccin-latte-deep`, defined in
    `server/internal/content/highlight.go`.
    - The official Latte accents are pastel. On mantle, strings (green) are
      2.75:1, numbers (peach) 2.45:1, and comments (overlay0) 2.14:1.
    - The style therefore keeps the Latte hues but darkens each in OKLab to
      at least 4.6:1 on mantle. These darkened values are derived, not
      official Catppuccin colors.
    - Comments use subtext1 in italics. A Go test checks every code color at
      4.5:1 or more.
    - The fence language `curl` is highlighted as bash.
  - Mermaid and ECharts use matching theme variables.
- **Code font:** `JetBrainsMonoNL Nerd Font` when installed (the terminal's
  font), else JetBrains Mono from `@fontsource-variable/jetbrains-mono`
  5.3.0 (OFL-1.1), served from the site. Ligatures are off, to match the NL
  variant.
  - Vite never inlines fonts as `data:` URLs (`assetsInlineLimit`), because
    the CSP allows `font-src 'self'` only.
  - Verified in Chromium: the local font renders when installed, and the
    bundled font loads otherwise, with no CSP violations.
  - draw.io, Excalidraw, and Swagger UI keep the colors stored in their
    files and their own styles.
- **Layout:**
  - The content area has 5% horizontal padding and no maximum width.
  - Articles show a sticky table of contents on the right, from 1,024 px of
    content width. It is built from the body's h1–h6 headings, which have
    ids from goldmark's auto heading IDs. The current section is tracked on
    scroll.
- **Embeds:** `EmbedHost` mounts one `React.lazy` renderer per
  `[data-embed]` type, each inside its own error boundary. All renderers are
  read-only, with no editor UI.
  - **Mermaid:** `mermaid.render` with `securityLevel: 'strict'`.
  - **ECharts:** `echarts.init` and a `ResizeObserver`.
  - **Excalidraw:** `exportToSvg` from `@excalidraw/excalidraw`
    ([export utils](https://docs.excalidraw.com/docs/@excalidraw/excalidraw/api/utils/export)).
    Only the SVG is inserted; the Excalidraw editor component is never
    mounted.
  - **OpenAPI:** `swagger-ui-react` with `supportedSubmitMethods: []`. It
    supports Swagger 2.0, and its peer range `react >=16.8 <20` allows
    React 19.
  - **draw.io:** `viewer-static.min.js`, vendored from the draw.io repository
    at tag `v32.4.1` (`src/main/webapp/js/viewer-static.min.js`, confirmed
    with the GitHub API). draw.io is Apache-2.0, and the license notice is
    kept. The script is served from `/vendor/drawio/`, so it is same-origin.
    `GraphViewer.processElements` renders a `data-mxgraph` div with
    `toolbar: "zoom pages"`, `editable: false`, and `lightbox: false`
    ([draw.io embed](https://www.drawio.com/doc/faq/embed-html)). That page
    does not cover self-hosting, so task 9.4 verifies it. The docsify site
    already vendors `viewer.min.js`, which suggests it works.
  - **MP4:**

    ```html
    <video controls preload="metadata" playsinline src=…>
      <a href=… download>…</a>
    </video>
    ```

    There is no autoplay. `preload="metadata"` is what the spec recommends
    ([MDN video](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/video)).

### D7. Repository layout

The layout follows the Go server layout
([go.dev module layout](https://go.dev/doc/modules/layout)) for Go, and npm workspaces
([npm workspaces](https://docs.npmjs.com/cli/v11/using-npm/workspaces)) for
TypeScript. The root `package.json` declares
`"workspaces": ["apps/*", "packages/*"]` and
`"packageManager": "npm@11.19.0"` with `engines.npm` `>=11.19.0`:

```
server/            Go module: cmd/notes/; internal/{content,store,search,api,httpx,export,config,logx}
apps/web/          React SPA + Storybook (.storybook/); public/vendor/drawio/
packages/schemas/  Zod schemas (sidebar, front matter) + JSON Schema generator
openapi/           generated openapi.json (checked in)
schemas/           generated sidebar/frontmatter JSON Schemas (checked in)
content/           sample content for dev; Faker fixtures live under testdata/
data/              runtime: builds/, current → (git-ignored)
```

### D8. Toolchain and versions (checked 2026-10-11)

Versions are pinned exactly. Sources:

| What | Source |
|---|---|
| Go | https://go.dev/dl/?mode=json |
| Node.js | https://endoflife.date/api/nodejs.json |
| Go modules | `https://proxy.golang.org/<module>/@latest` |
| npm packages | `npm view <pkg> version` |
| golangci-lint, draw.io | GitHub releases API |

**Languages and runtimes**

| Item | Version | Notes |
|---|---|---|
| Go | 1.27.2 | |
| TypeScript | 7.0.2 | Only TS 7. No TS 6 compatibility package is used (see below) |
| Node.js | 24 LTS | Web build and dev tooling only. Node 26 becomes LTS on 2026-10-28; upgrade then |
| npm | 11.19.0 or later | The npm bundled with Node 24.21.0. Chosen by the user over the latest npm 12.2.0, so no global npm upgrade is needed. npm 11.19.0 also lists package install scripts for approval (`npm install-scripts`) instead of running them silently; observed during `npm install` |
| SQLite | 3.53.4 | Bundled in modernc.org/sqlite |

**Go modules (`server/`)**

| Module | Version | Use |
|---|---|---|
| github.com/gin-gonic/gin | v1.12.0 | HTTP framework |
| github.com/danielgtaylor/huma/v2 | v2.39.1 | OpenAPI 3.1, problem+json (`humagin`) |
| gorm.io/gorm | v1.31.2 | ORM |
| gorm.io/driver/sqlite | v1.6.0 | GORM SQLite dialector |
| modernc.org/sqlite | v1.60.1 | Pure-Go SQLite driver (FTS5) |
| github.com/yuin/goldmark | v1.8.6 | Markdown |
| go.abhg.dev/goldmark/frontmatter | v0.3.0 | Front matter |
| github.com/alecthomas/chroma/v2 | v2.27.0 | Syntax highlighting |
| github.com/microcosm-cc/bluemonday | v1.0.27 | HTML sanitizing |
| github.com/santhosh-tekuri/jsonschema/v6 | v6.0.3 | JSON Schema validation |
| golang.org/x/time | v0.16.0 | Search rate limiting |
| github.com/fsnotify/fsnotify | v1.10.1 | Build swap detection |
| github.com/spf13/cobra | v1.10.2 | CLI |
| github.com/getkin/kin-openapi | v0.149.0 | OpenAPI validation in tests |
| golangci-lint | v2.14.0 | Go lint (tool) |

**Web app (`apps/web`)**

| Package | Version |
|---|---|
| react / react-dom / @types/react / @types/react-dom | 19.3.0 |
| vite / @vitejs/plugin-react | 8.3.4 / 6.1.2 |
| react-router | 8.4.0 |
| @tanstack/react-query | 5.104.1 |
| tailwindcss / @tailwindcss/vite | 4.3.3 |
| zod | 4.6.5 |
| pino (browser build) | 10.4.0 |
| mermaid | 12.1.0 |
| echarts | 6.1.0 |
| @excalidraw/excalidraw | 0.18.1 |
| swagger-ui-react / @types/swagger-ui-react | 5.33.1 / 5.18.0 |
| draw.io `viewer-static.min.js` (vendored) | v32.4.1 |
| @sqlite.org/sqlite-wasm (static-export search) | 3.53.4-build2 (same SQLite 3.53.4 as the Go side) |

**Shared and dev tooling (root, `packages/schemas`)**

| Package | Version |
|---|---|
| typescript | 7.0.2 |
| zod | 4.6.5 |
| @faker-js/faker | 10.6.0 |
| vitest / @vitest/browser-playwright | 5.0.3 |
| playwright | 1.64.0 |
| storybook / @storybook/react-vite / @storybook/addon-vitest / @storybook/addon-docs | 10.6.1 |
| oxlint / oxlint-tsgolint | 1.87.0 / 7.0.2004 |
| prettier | 3.9.10 |
| @types/node | 24.19.2 (matches Node 24) |

**TypeScript 7 only:**
- typescript-eslint supports only TypeScript `>=4.8.4 <6.1.0`
  ([typescript-eslint](https://typescript-eslint.io/users/dependency-versions/)).
  ESLint is therefore replaced by Oxlint.
- Oxlint's type-aware rules run on `typescript-go` through tsgolint. They
  require TS 7.0+, are in "early stages", and cover 59 of 61 typescript-eslint
  type-aware rules
  ([Oxlint type-aware](https://oxc.rs/docs/guide/usage/linter/type-aware.html)).
  The react plugin provides `rules-of-hooks` and `exhaustive-deps`.
- Oxlint has no Storybook plugin
  ([Oxlint rules](https://oxc.rs/docs/guide/usage/linter/rules.html)), so
  Storybook conventions are enforced by the story tests instead.
- The rule that forbids `dark:` classes is a small `npm run check` script
  (grep), not a lint plugin.
- Vite transpiles TypeScript with Oxc and does not need the `typescript`
  package ([Vite](https://vite.dev/guide/features)).
- Storybook's default `reactDocgen: 'react-docgen'` does its own analysis.
  `react-docgen-typescript`, which invokes the TS compiler, is not used
  ([Storybook](https://storybook.js.org/docs/api/main-config/main-config-typescript)).
- `tsc --noEmit` (TS 7) is the type-check step.

### D9. Testing and logging

- **Vitest projects:**
  - `unit` (Node): schemas and web utilities.
  - `storybook` (browser mode, Playwright Chromium, through
    `@storybook/addon-vitest`): every story is a smoke test, and play
    functions hold the interaction assertions
    ([Storybook Vitest addon](https://storybook.js.org/docs/writing-tests/integrations/vitest-addon)).
  - Go uses `go test`.
- **Fixtures:** Faker, seeded, generates articles, including the
  1,000-article search set and date clusters for the calendar.
- **Real-content check:** `notes migrate-sidebar` and `notes build` against a
  temp copy of `~/Documents/github/docsify`.
- **Logging:**
  - Web: Pino browser build. Only errors and embed failures are logged, to
    the console.
  - Go: `log/slog` JSON with Pino-compatible keys (numeric `level`, `time` in
    epoch ms, `msg`). This key mapping is a design choice, not a Pino
    requirement.
  - Query strings are logged without the search text, to keep readers'
    searches private.

### D10. Static export (`notes export`)

- **Command:**
  `notes export --out <dir> [--base-path /] [--video-base-url <URL>]` reads
  `data/current` (the latest successful build) and writes the files listed
  below. GitHub Pages is "a static site hosting service that takes HTML, CSS,
  and JavaScript files straight from a repository"
  ([What is GitHub Pages](https://docs.github.com/en/pages/getting-started-with-github-pages/what-is-github-pages)),
  so everything the server computes at runtime is precomputed here.

  | Output | Content |
  |---|---|
  | `<out>/` | The built web app (`apps/web/dist`) and `config.js` (`{ mode: "static", basePath, videoBaseUrl }`). It is a separate file, because the CSP forbids inline scripts |
  | `<out>/data/` | `sidebar.json`, `calendar.json` (every day with a count), `days/<date>.json`, `articles/<id>.json`. These are the same Go structs the Huma handlers return, serialized as they are |
  | `<out>/data/search.db` | A search-only SQLite database: the `articles_fts` table plus the ID, title, date, and category path |
  | `<out>/assets/`, `<out>/vendor/drawio/` | Published assets and the vendored draw.io viewer |
  | `<out>/a/<id>/index.html`, `<out>/d/<date>/index.html` | Copies of the app shell, so direct links and reloads work without server-side routing |
  | `<out>/404.html` | The app shell, which shows the not-found view. GitHub Pages documents custom 404 pages; the exact behavior is checked in task 5.7 |
  | `<out>-videos/` | Videos plus `videos-manifest.json` (relative path, size, SHA-256), outside the site output |

- **Defaults:**
  - `npm run export` builds everything and writes `./static` and
    `./static-videos`, replacing both. `notes export --out` defaults to
    `static`.
  - Output files are world-readable (0644 files, 0755 folders), so any web
    server can serve them.
  - Story test fixtures in `apps/web/public/fixtures` are removed from the
    web build.
- **Videos:** the export never uploads anything. The operator uploads
  `<out>-videos/` by hand. Video `src` becomes `videoBaseUrl` plus the
  content-hashed relative path, so re-uploads are cache-safe. Hosting
  (S3 + CloudFront is planned) is the operator's concern and out of scope
  here.
- **Guards** (from the [GitHub Pages limits](https://docs.github.com/en/pages/getting-started-with-github-pages/github-pages-limits)
  and the [file size limit](https://docs.github.com/en/repositories/working-with-files/managing-large-files/about-large-files-on-github)):
  - The total site output must be 1 GB or less.
  - Every output file must be 100 MiB or less.
  - A video embed requires `--video-base-url`.
- **Where the source lives:** on GitHub Free, a Pages repository must be
  public
  ([Creating a Pages site](https://docs.github.com/en/pages/getting-started-with-github-pages/creating-a-github-pages-site)).
  The Markdown source therefore stays local or in a private repository, and
  only the export output is pushed to the public Pages repository, deployed
  from a branch. This needs no GitHub Actions workflow.
- **CSP:** each HTML page gets
  `<meta http-equiv="Content-Security-Policy">`:

  ```
  default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline';
  img-src 'self' data: blob:; media-src 'self' <video origin>; connect-src 'self';
  font-src 'self'; worker-src 'self' blob:; object-src 'none'; base-uri 'self'; form-action 'none'
  ```

  `frame-ancestors` is not used, because browsers ignore it in `<meta>`.
  `'wasm-unsafe-eval'` is not needed, which `apps/web/e2e/static-search.mjs`
  verified. SQLite WASM runs in a dedicated worker, and the document's
  `<meta>` policy does not govern that worker. Search behaves the same with
  and without the source, with no CSP violations.
  - Zod 4 probes `new Function` once to decide whether to compile validators,
    and strict CSPs report that probe as a violation even though Zod catches
    it. The web app therefore sets `z.config({ jitless: true })` at startup,
    which skips the probe (see the comment in `zod/v4/core/util.js`).
- **Web app data source:** a `NotesSource` interface has two implementations.
  - `ApiSource`: the hand-written fetch client.
  - `StaticSource`: fetches `data/*.json` under the base path and validates
    them with the same Zod schemas.

  `config.js` selects the implementation, and the UI components do not know
  which one is in use. The server version serves its own `config.js` with
  `mode: "server"`.
- **In-browser search:** the official SQLite WASM build
  (`@sqlite.org/sqlite-wasm`) runs in a Web Worker, so the UI thread never
  blocks.
  - On the first search, the worker fetches `data/search.db` and opens it in
    memory.
  - It then runs the same FTS5 `MATCH`, `bm25` weights, and `snippet()` as the
    Go search.
  - A local spike with 3.53.4-build2 in Node confirmed that FTS5 with the
    trigram tokenizer works in this build: the `本語の` substring match, the
    title-weighted `bm25` order, and `snippet()` highlighting.
  - Loading the fetched bytes follows the SQLite WASM cookbook
    ([cookbook](https://sqlite.org/wasm/doc/trunk/cookbook.md)):
    `wasm.allocFromTypedArray`, then `capi.sqlite3_deserialize` with
    `SQLITE_DESERIALIZE_FREEONCLOSE` (plus `READONLY`).
    `apps/web/src/lib/sqlite.ts` implements it.
- **Shared query rules:** query normalization (trim, 3-character minimum,
  200-character maximum, quoting each term) is implemented in Go and TS. Both
  run `fixtures/search/cases.json` against the same fixture database
  (`fixtures/search/fixture.search.db`, written by
  `go test ./internal/search -update-searchdb`).

## Risks / Trade-offs

- [Draft notes from the password-protected docsify site could be published by
  mistake] → During migration, `draft: true` is added to every article first
  (task 11.2), and articles are un-drafted one by one.
- [No WriteTimeout means a slow client can hold a connection] → The CDN or
  reverse proxy buffers clients. `IdleTimeout` and `ReadHeaderTimeout` still
  apply.
- [Remote embeds are a snapshot from build time] → The latest version is
  fetched again on the next build. This is intended for reproducibility.
- [Videos larger than the CDN's cache limit (512 MB on Cloudflare Free, Pro,
  and Business) are served from the origin every time] → Acceptable at blog
  traffic levels. Consider re-encoding very large videos.
- [Oxlint type-aware linting is early stage, and there are no Storybook lint
  rules] → `tsc --noEmit` and the story tests are the main gates. Lint is an
  extra check.
- [npm hoists dependencies, so a package can use a dependency missing from its
  own `package.json`] → `npm run check` runs `npm ls --all` to catch
  dependency mismatches. Its effectiveness is my assessment and is not backed
  by official docs.
- [The trigram index is about 3× the text size] → Acceptable for hundreds to a
  few thousand articles.
- [Missing `date` falls back to mtime, which a copy or checkout can change] →
  The build warns. `notes stamp-dates` writes each article's creation date
  (or mtime where the OS does not record creation times) into its front
  matter once, so the date no longer depends on file times. `npm start`,
  `npm run dev`, and `npm run build:content` run it before every build, so
  new articles get their date on their first build.
- [`gorm.io/driver/sqlite` imports `mattn/go-sqlite3` even though it is not
  used, and was last released 2025-06] → CI builds with `CGO_ENABLED=0`. If
  the dialector stalls, swap the `Conn` to another pure-Go wrapper.
- [Huma returns 422 by default] → Override the error mapping (D5), covered by
  a test in task 5.1.
- [swagger-ui-react is heavy] → Lazy-loaded only where it is used.
- [The CSP allows inline styles] → Scripts stay `'self'` only. The CSP is
  tightened if the libraries allow it.

- [The static site has no `frame-ancestors` protection, because the directive
  is ignored in `<meta>` CSP
  ([MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/frame-ancestors))]
  → Accepted for a read-only blog. The server version still sends the header.
- [The search index download grows with the content (the trigram index is
  about 3× the text)] → The index is loaded lazily on first search, holds
  only the FTS table and minimal metadata, and is served as a cacheable
  static file.
- [Search logic could drift between Go and TS] → Shared golden cases in
  `fixtures/search/cases.json` run in both `go test` and Vitest against the
  same fixture database.
- [Static and server data shapes could drift] → The export writes the exact
  structs the API handlers return, and the web app validates both with the
  same Zod schemas.

## Migration Plan

1. Run the new system locally (loopback) alongside the docsify site, which
   stays unchanged.
2. Copy the docsify content into `content/`. Run `notes migrate-sidebar`, then
   mark every article `draft: true`. Fix docsify-only syntax: `:include`
   links and `[swagger](...)` links become `![openapi](...)`.
3. Review each article. Add `date`, remove `draft`, and rebuild.
4. To publish: put the server behind a CDN or reverse proxy, set
   `NOTES_LISTEN` and `NOTES_TRUSTED_PROXIES`, and verify the security
   headers.
5. To publish the static version: run `notes export`, upload the video
   folder to the video host by hand, push the site output to the public Pages
   repository, and enable Pages for that repository (deploy from a branch).
6. Rollback: stop the server or remove the CDN route, or revert the Pages
   repository. The docsify site stays
   as it is.

## Open Questions

- Should articles that are published but not listed in `sidebar.json` be
  indexed and searchable? Proposed default: yes, but they are not shown in the
  sidebar.
