# Developer guide

How the note site is built, tested, and deployed. For writing content and
running the app, see [README.md](README.md). Design decisions and specs are
in `openspec/changes/build-note-system/`.

## Tech stack and versions

Versions are pinned exactly (checked 2026-10-11 against go.dev, nodejs.org
release data, proxy.golang.org, and the npm registry).

### Languages and runtimes

| Item       | Version                 | Notes                                                                        |
| ---------- | ----------------------- | ---------------------------------------------------------------------------- |
| Go         | 1.27.2                  | server and build tool; builds with `CGO_ENABLED=0`                           |
| Node.js    | 24 LTS (24.15 or later) | web app build and dev tooling only                                           |
| npm        | 11.19.0 or later        | bundled with Node 24.21.0; `.npmrc` sets `engine-strict`                     |
| TypeScript | 7.0.2                   | `tsc` only; no TypeScript 6 compatibility package                            |
| SQLite     | 3.53.4                  | bundled in `modernc.org/sqlite` (Go) and `@sqlite.org/sqlite-wasm` (browser) |

### Go modules (`server/`)

| Module                                                     | Version          | Use                                                                                               |
| ---------------------------------------------------------- | ---------------- | ------------------------------------------------------------------------------------------------- |
| github.com/gin-gonic/gin                                   | v1.12.0          | HTTP framework                                                                                    |
| github.com/danielgtaylor/huma/v2                           | v2.39.1          | API handlers, OpenAPI 3.1, RFC 9457 errors (`humagin`)                                            |
| gorm.io/gorm / gorm.io/driver/sqlite                       | v1.31.2 / v1.6.0 | ORM over `modernc.org/sqlite` (`DriverName: "sqlite"`)                                            |
| modernc.org/sqlite                                         | v1.60.1          | pure-Go SQLite, FTS5 trigram index                                                                |
| github.com/yuin/goldmark, go.abhg.dev/goldmark/frontmatter | v1.8.6, v0.3.0   | Markdown and front matter                                                                         |
| github.com/alecthomas/chroma/v2                            | v2.27.0          | build-time syntax highlighting (`catppuccin-latte-deep` style in `internal/content/highlight.go`) |
| github.com/microcosm-cc/bluemonday                         | v1.0.27          | HTML sanitizing                                                                                   |
| github.com/santhosh-tekuri/jsonschema/v6                   | v6.0.3           | sidebar and front matter validation                                                               |
| golang.org/x/time                                          | v0.16.0          | search rate limiting                                                                              |
| github.com/fsnotify/fsnotify                               | v1.10.1          | picking up new builds                                                                             |
| github.com/spf13/cobra                                     | v1.10.2          | CLI                                                                                               |
| github.com/getkin/kin-openapi                              | v0.149.0         | OpenAPI validation in tests                                                                       |

### Web app (`apps/web`) and shared packages

| Package                             | Version                                                                                                                                                                                        |
| ----------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| react, react-dom                    | 19.3.0 (forced for all dependencies with npm `overrides`; see below)                                                                                                                           |
| vite / @vitejs/plugin-react         | 8.3.4 / 6.1.2                                                                                                                                                                                  |
| react-router                        | 8.4.0                                                                                                                                                                                          |
| @tanstack/react-query               | 5.104.1                                                                                                                                                                                        |
| tailwindcss / @tailwindcss/vite     | 4.3.3 (light theme only)                                                                                                                                                                       |
| zod                                 | 4.6.5 (also in `packages/schemas`)                                                                                                                                                             |
| pino                                | 10.4.0 (browser logger)                                                                                                                                                                        |
| mermaid                             | 12.1.0                                                                                                                                                                                         |
| echarts                             | 6.1.0                                                                                                                                                                                          |
| @excalidraw/excalidraw              | 0.18.1 (fonts copied to `public/vendor/excalidraw` by `scripts/vendor-assets.mjs`)                                                                                                             |
| swagger-ui-react                    | 5.33.1                                                                                                                                                                                         |
| @sqlite.org/sqlite-wasm             | 3.53.4-build2 (static-site search in a Web Worker)                                                                                                                                             |
| @fontsource-variable/jetbrains-mono | 5.3.0 (OFL-1.1). Code font fallback; the stack prefers a locally installed `JetBrainsMonoNL Nerd Font`. Fonts are never inlined as `data:` URLs, because the CSP allows `font-src 'self'` only |
| draw.io `viewer-static.min.js`      | v32.4.1, vendored in `public/vendor/drawio` (Apache-2.0)                                                                                                                                       |

### Tooling

| Tool                                                                             | Version                                                 |
| -------------------------------------------------------------------------------- | ------------------------------------------------------- |
| storybook, @storybook/react-vite, @storybook/addon-vitest, @storybook/addon-docs | 10.6.1                                                  |
| vitest, @vitest/browser-playwright                                               | 5.0.3                                                   |
| playwright                                                                       | 1.64.0 (Chromium for story tests and end-to-end checks) |
| oxlint / oxlint-tsgolint                                                         | 1.87.0 / 7.0.2004 (type-aware lint on TypeScript 7)     |
| prettier                                                                         | 3.9.10                                                  |
| @faker-js/faker                                                                  | 10.6.0 (seeded fixtures)                                |
| golangci-lint                                                                    | v2.14.0 (binary install)                                |
| ffmpeg                                                                           | any recent version; only for `e2e:video`                |

## Setup

```sh
npm install
npx playwright install chromium     # for story tests and end-to-end checks
# golangci-lint: https://golangci-lint.run/docs/welcome/install/local/
```

## Commands

| Command                              | What it does                                                                                                                                           |
| ------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `npm run dev`                        | Runs `stamp-dates` and `build`, then the Go server on :8080 plus Vite on http://localhost:5173 (proxies `/api`, `/assets`, `/config.js`)               |
| `npm start`                          | Builds the web app into `bin/notes`, runs `stamp-dates` and `build` on `content/`, serves on `NOTES_LISTEN`                                            |
| `npm run build:content`              | `notes stamp-dates` then `notes build --content ../content --data ../data`                                                                             |
| `npm run export`                     | Builds the web app and content, then `notes export` to `./static` and `./static-videos` (replaced each run); `NOTES_BASE_PATH`, `NOTES_VIDEO_BASE_URL` |
| `npm run preview:static`             | Serves `./static` on http://127.0.0.1:4173 like GitHub Pages (folder `index.html`, `404.html` fallback); `NOTES_BASE_PATH`, `PORT`                     |
| `npm run build:go`                   | `CGO_ENABLED=0 go build ./...`                                                                                                                         |
| `npm run generate:openapi`           | Writes `openapi/openapi.json` from the Go API                                                                                                          |
| `npm run generate -w @notes/schemas` | Writes the JSON Schemas from the Zod schemas                                                                                                           |

### `notes` CLI

| Command                                           | What it does                                                                                                                                                                      |
| ------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `notes build --content --data`                    | Validates and compiles content into `data/builds/<timestamp>/`, then swaps `data/current` atomically. Keeps only the newest build.                                                |
| `notes serve`                                     | Serves the SPA, API (`/api/v1/*`), assets (`/assets/*`, with Range support), `/healthz`, and `/config.js`                                                                         |
| `notes export --out --base-path --video-base-url` | Static copy for static hosts (data JSON, `search.db`, app shells, `404.html`, CSP `<meta>`). `--out` defaults to `static` (env `NOTES_OUT`); files are written 0644, folders 0755 |
| `notes stamp-dates --content [--dry-run]`         | Writes `date:` into articles without one: creation time (macOS, Windows), otherwise modification time, in the local time zone. Existing dates are not touched.                    |
| `notes migrate-sidebar <_sidebar.md>`             | Converts a docsify sidebar (see `docs/migration.md`)                                                                                                                              |
| `notes openapi`                                   | Prints the OpenAPI 3.1 document                                                                                                                                                   |

### Environment variables

| Variable                | Default          | Meaning                                                         |
| ----------------------- | ---------------- | --------------------------------------------------------------- |
| `NOTES_LISTEN`          | `127.0.0.1:8080` | Listen address; loopback by default                             |
| `NOTES_TRUSTED_PROXIES` | (none)           | Comma-separated IPs or CIDRs whose `X-Forwarded-For` is trusted |
| `NOTES_SEARCH_RATE`     | `60`             | Search requests per client IP per minute (429 above it)         |
| `NOTES_DATA`            | `data`           | Builds folder                                                   |
| `NOTES_CONTENT`         | `content`        | Content folder                                                  |

## Checks and tests

| Script                              | What it checks                                                                                                                                               |
| ----------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `npm run lint`                      | Oxlint (type-aware) and golangci-lint                                                                                                                        |
| `npm run typecheck`                 | `tsc -b --force` (TypeScript 7)                                                                                                                              |
| `npm test`                          | Vitest (unit tests in jsdom and Node, Storybook stories in Chromium) and `go test ./...`                                                                     |
| `npm run check`                     | JSON Schemas, OpenAPI document, and API response schemas are current; no `dark:` classes; every workspace import is declared in its `package.json`; Prettier |
| `npm run e2e:embeds -w @notes/web`  | Every embed type renders with no CSP violation and no third-party request (server and static export)                                                         |
| `npm run e2e:search -w @notes/web`  | Static search: lazy single download of `search.db`, results, CSP                                                                                             |
| `npm run e2e:video -w @notes/web`   | Real H.264 video (made with ffmpeg): no autoplay, inline playback, seeking with range requests                                                               |
| `npm run e2e:docsify -w @notes/web` | End-to-end check on a copy of `~/Documents/github/docsify`                                                                                                   |

CI (`.github/workflows/ci.yml`) runs the Go build, check, lint, typecheck,
and tests.

Shared fixtures:

- `fixtures/search/cases.json` and `fixture.search.db` are used by both the
  Go search and the browser search, so their results stay identical.
- Regenerate them with `go test ./internal/search -update-cases -update-searchdb`.

## Architecture notes

- **Content build** (`server/internal/content`):
  1. goldmark parses Markdown. An AST transformer turns the embed prefixes
     into `data-embed` placeholders.
  2. Assets are resolved and content-hashed. Remote embeds and images are
     downloaded at build time (30 seconds and 10 MB at most).
  3. chroma highlights code, and bluemonday sanitizes the HTML.
  4. Drafts stop after their front matter, so their assets are never
     published.
- **Storage and search** (`internal/store`, `internal/search`):
  - Each build is a fresh SQLite database (GORM `AutoMigrate`) with an FTS5
    `trigram` table.
  - Search quotes every term, so FTS5 operators are never interpreted.
  - Ranking uses `bm25` weights, with title highest.
  - Snippets are HTML-escaped, with `<mark>` around matches.
- **API** (`internal/api`):
  - Huma on Gin, read-only.
  - Validation errors are mapped from 422 to 400.
  - Arrays are never null.
- **Web app** (`apps/web`):
  - A `NotesSource` interface has two implementations: the API in the
    server version, and static JSON plus SQLite WASM search in the static
    version. `config.js` selects one.
  - Embeds load lazily per type, each inside an error boundary.
  - Built chunk URLs resolve against the runtime base path (Vite
    `renderBuiltUrl`).
- **Security**:
  - No authentication or cookies. Only GET and HEAD are routed.
  - A strict CSP: a header in the server version, a `<meta>` tag in the
    static version.
  - Zod runs with `jitless`, so it never probes `eval`.
  - The static CSP needs no `'wasm-unsafe-eval'`, because SQLite WASM runs
    in a worker.
- **Known compromise**: the npm `overrides` force React 19.3.0 because some
  of swagger-ui-react's dependencies declare React 18 at most. It works in
  the story tests and end-to-end checks, but those packages don't officially
  support React 19.

## Deployment

### Server version

1. Run `notes serve` on loopback behind a CDN or reverse proxy that
   terminates TLS.
2. Set `NOTES_TRUSTED_PROXIES` to the proxy's address.
3. Check the security headers: `curl -sI https://your.site/ | grep -iE 'content-security|strict-transport|nosniff|referrer'`.
4. Large videos: CDNs cap the size they cache. Cloudflare caches files up
   to 512 MB on the Free, Pro, and Business plans
   ([docs](https://developers.cloudflare.com/cache/concepts/default-cache-behavior/)).
   Range requests still allow seeking.

### Static version (GitHub Pages)

1. Keep `content/` out of the public Pages repository. On GitHub Free, a
   Pages repository must be public.
2. Run `npm run build -w @notes/web`, then `notes export`.
3. Upload `site-videos/` by hand.
4. Push `site/` and deploy from a branch.

The export enforces the limits: 1 GB per site and 100 MiB per file
([Pages limits](https://docs.github.com/en/pages/getting-started-with-github-pages/github-pages-limits)).

How missing pages are served:

- GitHub Pages shows a custom error page from a `404.html` file when someone
  opens a page that does not exist
  ([custom 404](https://docs.github.com/en/pages/getting-started-with-github-pages/creating-a-custom-404-page-for-your-github-pages-site)).
- `notes export` writes `404.html` as a copy of the app shell, so unknown
  URLs show the app's not-found view.
- The status code sent with `404.html` is not documented. Every article and
  day URL therefore also gets its own `index.html` copy (`/a/<id>/` and
  `/d/<date>/`), so those URLs never depend on the 404 page.

## Repository layout

```
server/            Go: cmd/notes and internal/ (content, store, search, api, httpx, export, config, logx, webui)
apps/web/          React web app, Storybook stories (src/**/*.stories.tsx), end-to-end scripts (e2e/)
packages/schemas/  Zod schemas (content formats, API responses, search query rules)
schemas/           JSON Schemas generated from packages/schemas
openapi/           OpenAPI document generated from the Go API
fixtures/search/   search cases shared by the Go and browser search
content/           the notes (user content)
docs/              guides (docsify migration)
scripts/           npm script helpers (dev, start, export, preview, checks)
openspec/          OpenSpec change and specs for this project
.claude/           OpenSpec commands and skills for Claude Code
```

### What Git ignores

Everything below is generated and is rebuilt by the npm scripts, so it is not
committed (see `.gitignore`):

| Path                                                                    | Made by                                 |
| ----------------------------------------------------------------------- | --------------------------------------- |
| `node_modules/`                                                         | `npm install`                           |
| `dist/`, `bin/`, `*.tsbuildinfo`, `server/internal/webui/static/*`      | `npm start`, `npm run build:go`, `tsc`  |
| `apps/web/public/vendor/excalidraw/`                                    | the web build (copies Excalidraw fonts) |
| `data/`                                                                 | `npm run build:content`                 |
| `static/`, `static-videos/`                                             | `npm run export`                        |
| `coverage/`, `test-results/`, `playwright-report/`, `storybook-static/` | tests and Storybook                     |

Also ignored: `.env*` (except `.env.example`), `.claude/settings.local.json`,
editor folders, `.DS_Store`, and logs.

`content/` is committed, including drafts. Keep this repository private if the
notes are private; publish only `static/`.
