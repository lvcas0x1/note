# Migrating the docsify notes

This guide moves the notes in `~/Documents/github/docsify` into `content/`.
The docsify folder itself is never changed. All steps below were run on a
temporary copy by `npm run e2e:docsify -w @notes/web`
(`apps/web/e2e/docsify.mjs`), which passes for both the server version and a
static export.

## Steps

1. Copy the docsify folder to `content/`.
2. Delete the docsify-only files. `notes build` would otherwise publish them
   as articles:
   - `_sidebar.md`
   - `_navbar.md`
   - `_404.md`
   - `index.html`
   - `_plugin/` (except files that articles reference, such as
     `_plugin/vendor/drawio/test.drawio`)
3. Convert the sidebar:

   ```sh
   cd server && go run ./cmd/notes migrate-sidebar ../content/_sidebar.md --root ../content
   ```

   This writes `content/sidebar.json`. Folder links map to `README.md`, or
   to the folder's only `.md` file. For example,
   `[Session Manager](aws/sessionmanager/)` becomes
   `aws/sessionmanager/sessionmanager.md`.

4. Add `draft: true` front matter to every article. The build then publishes
   nothing; the e2e run confirmed `built 0 articles and 0 assets`.
5. Run `notes stamp-dates --content ../content` to write each article's
   file date into its front matter. Then review the articles one by one. For each, add a `date`, fix the
   docsify-only syntax below, and remove `draft`.
6. Run `npm run build:content` and fix any reported errors.

## Build warnings on the docsify notes

A full build of the 67 articles, with no drafts, succeeds and publishes 205
assets. It reports one kind of warning, 67 times:

| Warning                                                        | Count              | Fix                                                                                                                                                 |
| -------------------------------------------------------------- | ------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| `no date in front matter; using file modification date <date>` | 67 (every article) | Add `date: YYYY-MM-DD`. The modification date changes when files are copied or checked out, so the calendar is wrong until each article has a date. |

## Docsify-only syntax that needs manual edits

Found by searching the docsify notes. Each item is rendered as plain Markdown
until it is edited.

| Syntax                                                                 | Where                         | Replace with                                                                                         |
| ---------------------------------------------------------------------- | ----------------------------- | ---------------------------------------------------------------------------------------------------- |
| `[swagger](/swagger/swagger.json)` (docsify-swagger link)              | `swagger/swagger.md`          | `![swagger](/swagger/swagger.json)`                                                                  |
| `[filename](/_plugin/vendor/drawio/test.drawio ':include :type=code')` | `_index.md`                   | `![drawio](/_plugin/vendor/drawio/test.drawio)`                                                      |
| `![img](./_media/ai-icon.png ':size=25%')` (`:size` option)            | `_index.md`                   | `![img](./_media/ai-icon.png)`. The size option is not supported and is dropped.                     |
| `<!-- tabs:start -->` … `<!-- tabs:end -->` (docsify-tabs)             | `_index.md`, `aws/cli/cli.md` | Headings or separate sections. The comments are ignored, so the tab contents show one after another. |
| ` ```charty ` fences (16)                                              | `_index.md`                   | ` ```echarts ` with an ECharts option in JSON.                                                       |
| `> [!NOTE]`, `> [!TIP]`, … (GitHub-style alerts, 4)                    | `_index.md`                   | Plain blockquotes, or keep them as is; they render as ordinary quotes.                               |
| `:id=` link option                                                     | `_index.md`                   | A plain link.                                                                                        |

Mermaid fences (22 in total, for example in `auth/jwt.md` and `auth/saml.md`)
work unchanged.

## What the e2e check covers

On the server version and on a static export under `/notes/`:

- Opening an article from the sidebar.
- Mermaid rendering.
- The calendar day list for 2026-10-05 (4 articles) and opening Swagger from it.
- Search for "DrawIO", opening the result, and draw.io rendering.
- Direct load and reload of the article with an MP4.
- A draft that cannot be reached.
- No CSP violations, and no request to another origin. In the static
  export, the configured video host is the only allowed exception.
