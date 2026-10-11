# Spec Delta

## Purpose

Turns a folder of Markdown files, a nested `sidebar.json`, and embedded asset
files into validated, compiled, and indexed content that the note system serves.

## ADDED Requirements

### Requirement: Content folder layout
The build SHALL read one content folder that holds `sidebar.json` at its root, Markdown articles (`.md`) in any subfolder, and asset files referenced by articles: images, `.drawio`, `.excalidraw`, OpenAPI `.json`/`.yaml`, and `.mp4` videos.

#### Scenario: Build a valid folder
- **WHEN** the operator runs the build command on a content folder that has a valid `sidebar.json` and Markdown articles
- **THEN** the build exits with code 0 and the articles become available in the running system

#### Scenario: Missing sidebar
- **WHEN** the content folder has no `sidebar.json`
- **THEN** the build fails with a non-zero exit code and an error that names the missing file

### Requirement: Sidebar JSON schema
`sidebar.json` SHALL contain an ordered array of nodes. Each node SHALL have a non-empty `title`. A node SHALL have a `path` to a Markdown file relative to the content root, a `children` array of nodes, or both. Nesting depth SHALL be unlimited. A published JSON Schema SHALL describe the format.

#### Scenario: Nested categories
- **WHEN** `sidebar.json` contains a category node with a child category that has article children
- **THEN** the build accepts it and keeps the nesting and order

#### Scenario: Invalid node
- **WHEN** a node has no `title`, or has neither `path` nor `children`
- **THEN** the build fails and the error gives the JSON path of the invalid node

#### Scenario: Broken path
- **WHEN** a node's `path` points to a file that does not exist
- **THEN** the build fails and the error names the node and the missing file

### Requirement: Article front matter
Each article MAY start with a YAML front matter block. Recognized keys SHALL be `title` (string), `date` (ISO 8601 calendar date `YYYY-MM-DD`), `summary` (string), `tags` (string array), and `draft` (boolean, default `false`). Unknown keys SHALL be ignored with a warning.

#### Scenario: Valid front matter
- **WHEN** an article has front matter with `title`, `date: 2026-10-11`, and `summary`
- **THEN** the indexed article uses those values

#### Scenario: Invalid date
- **WHEN** an article's `date` is not a valid `YYYY-MM-DD` date
- **THEN** the build fails and the error names the file and the bad value

### Requirement: Draft exclusion
An article with `draft: true` SHALL be excluded from the build output. It SHALL NOT appear in the content database, the search index, the calendar, or the API, and assets used only by draft articles SHALL NOT be published. A `sidebar.json` node that points to a draft article SHALL be omitted, with a warning. A category left with no children and no path SHALL be omitted too.

#### Scenario: Draft article
- **WHEN** an article has `draft: true` and the build succeeds
- **THEN** requesting its ID returns 404, it is not in search results or calendar counts, and its sidebar node is absent

#### Scenario: Asset used only by a draft
- **WHEN** an image is referenced only by a draft article
- **THEN** the image is not in the build output and requesting its asset URL returns 404

#### Scenario: Asset shared with a published article
- **WHEN** an image is referenced by both a draft and a published article
- **THEN** the image is published

### Requirement: Metadata fallbacks
Every article SHALL get a title, a date, and a summary. If `title` is missing, the build SHALL use the first H1 heading, then the sidebar node title, then the file name. If `date` is missing, it SHALL use the file's last-modified date and print a warning. If `summary` is missing, it SHALL use the first paragraph as plain text, cut to 200 characters.

#### Scenario: Docsify article without front matter
- **WHEN** an article has no front matter and starts with an H1 heading and a paragraph
- **THEN** the article's title is the H1 text, its date is the file's modification date, its summary is the first paragraph, and the build prints a missing-date warning naming the file

### Requirement: Markdown compilation
The build SHALL compile each article to HTML with CommonMark plus GitHub Flavored Markdown tables, task lists, strikethrough, autolinks, and footnotes. Fenced code SHALL keep its language class for syntax highlighting. Raw HTML from Markdown SHALL be sanitized so scripts and event-handler attributes are removed.

#### Scenario: GFM features
- **WHEN** an article contains a table, a task list, a footnote, and a fenced `go` code block
- **THEN** the compiled HTML contains a table, checkboxes, a linked footnote, and a code block marked as language `go`

#### Scenario: Script in Markdown
- **WHEN** an article contains `<script>` or an `onclick` attribute in raw HTML
- **THEN** the compiled HTML contains neither

### Requirement: Embed detection
The build SHALL turn embed syntax into embed placeholders that keep the type and source. Image syntax whose alt text is `drawio`, `excalidraw`, `openapi`, `swagger`, or `mp4` SHALL become an embed of that type (`swagger` is an alias for `openapi`). Fenced blocks with language `mermaid` or `echarts` SHALL become embeds whose source is the block body. Embeds SHALL be allowed anywhere a block element is allowed.

#### Scenario: draw.io image syntax
- **WHEN** an article contains `![drawio](./diagrams/net.drawio)`
- **THEN** the compiled article contains a draw.io embed placeholder that points to that asset

#### Scenario: MP4 video syntax
- **WHEN** an article contains `![mp4](./videos/demo.mp4)`
- **THEN** the compiled article contains a video embed placeholder that points to that asset

#### Scenario: Mermaid fence
- **WHEN** an article contains a fenced block with language `mermaid`
- **THEN** the compiled article contains a Mermaid embed placeholder with the block body

#### Scenario: Ordinary image
- **WHEN** an article contains `![photo](./img/a.png)`
- **THEN** the compiled article contains a normal image

### Requirement: Local asset validation
For each embed or image that references a local file, the build SHALL check that the file exists inside the content folder and SHALL publish it as an asset. A reference that escapes the content folder SHALL fail the build. The build SHALL NOT impose a file-size limit on local assets.

#### Scenario: Missing local asset
- **WHEN** an article references `![excalidraw](./missing.excalidraw)` and the file does not exist
- **THEN** the build fails and the error names the article and the missing file

#### Scenario: Path traversal
- **WHEN** an article references `![drawio](../../etc/passwd)`
- **THEN** the build fails with a path-outside-content error

#### Scenario: Large local video
- **WHEN** an article references a 2 GB local `.mp4` file
- **THEN** the build publishes it without a size error

### Requirement: Video embed source
An `mp4` embed SHALL reference a local file whose name ends in `.mp4`. A remote URL or another extension SHALL fail the build.

#### Scenario: Remote video URL
- **WHEN** an article contains `![mp4](https://example.com/a.mp4)`
- **THEN** the build fails with an error that videos must be local files

#### Scenario: Wrong extension
- **WHEN** an article contains `![mp4](./clip.mov)`
- **THEN** the build fails and the error names the file and the expected `.mp4` extension

### Requirement: Remote embed download
For `drawio`, `excalidraw`, and `openapi` embeds that use an `http` or `https` URL, the build SHALL download the file and publish it as a local asset. The compiled article SHALL reference only the local copy. A download that fails, returns a non-2xx status, exceeds 10 MB, or takes more than 30 seconds SHALL fail the build. Other URL schemes SHALL fail the build.

#### Scenario: Remote draw.io URL
- **WHEN** an article contains `![drawio](https://example.com/a.drawio)` and the URL returns the file
- **THEN** the compiled article contains a draw.io embed placeholder that points to a local asset, and the output contains no reference to `example.com`

#### Scenario: Unreachable remote URL
- **WHEN** a remote embed URL returns 404 or cannot be reached
- **THEN** the build fails and the error names the article and the URL

#### Scenario: Unsupported scheme
- **WHEN** an article contains `![openapi](file:///etc/hosts)`
- **THEN** the build fails with an unsupported-scheme error

### Requirement: ECharts option validation
An `echarts` block body SHALL be a JSON object that is used as the ECharts option. The build SHALL reject a block that is not valid JSON.

#### Scenario: Invalid echarts JSON
- **WHEN** an `echarts` block contains malformed JSON
- **THEN** the build fails and the error names the file and line of the block

### Requirement: Atomic rebuild
A build SHALL replace the served content all at once. Until a build succeeds, the system SHALL keep serving the previous content. A failed build SHALL leave the served content unchanged.

#### Scenario: Failed rebuild
- **WHEN** a build fails validation while the server is running
- **THEN** the server keeps serving the content from the last successful build

### Requirement: Date stamping
The build tool SHALL provide a command that writes `date: YYYY-MM-DD` into the front matter of every article that has no `date`. It SHALL use the file's creation time where the operating system records it and the modification time otherwise, in the local time zone. It SHALL NOT change articles that already have a date, SHALL keep the file's line endings, and SHALL offer a dry run that lists the dates without writing.

#### Scenario: Article without front matter
- **WHEN** the command runs on an article with no front matter that was created on 2026-09-22
- **THEN** the article starts with front matter containing `date: 2026-09-22`, and the build no longer warns about a missing date

#### Scenario: Existing date
- **WHEN** an article already has `date: 2025-01-02`
- **THEN** the command leaves the file unchanged

#### Scenario: Dry run
- **WHEN** the command runs with the dry-run option
- **THEN** it lists the articles and dates it would write and changes no file

### Requirement: Docsify sidebar migration
The build tool SHALL provide a command that converts a docsify `_sidebar.md` (nested Markdown link list) into a `sidebar.json` that passes validation. A link to a folder SHALL map to the folder's `README.md` if it has one. If the link has no children and the folder holds exactly one `.md` file, it SHALL map to that file. Otherwise it SHALL become a category node without a path. Article files SHALL NOT be rewritten.

#### Scenario: Convert docsify sidebar
- **WHEN** the operator runs the migration command on `~/Documents/github/docsify/_sidebar.md`
- **THEN** it writes a `sidebar.json` with the same nesting, titles, and order, and the build accepts it

#### Scenario: Leaf folder link
- **WHEN** the docsify sidebar has `[Session Manager](aws/sessionmanager/)` with no children and that folder holds only `sessionmanager.md`
- **THEN** the node's path is `aws/sessionmanager/sessionmanager.md`
