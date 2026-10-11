# Spec Delta

## Purpose

Produces a self-contained static copy of the published site that a static host
such as GitHub Pages can serve without the Go server, with the same reading,
navigation, calendar, and search behavior as the server version.

## ADDED Requirements

### Requirement: Export command
The build tool SHALL provide an export command that writes, from the latest successful build, a static site to an output folder. The output SHALL contain the web app, the article, sidebar, and calendar data, the search index, the published assets, and the vendored embed renderers. It SHALL NOT contain draft articles, their assets, or Markdown sources.

#### Scenario: Export a build
- **WHEN** the operator runs the export command after a successful build
- **THEN** the output folder holds a complete static site and the command exits with code 0

#### Scenario: No build yet
- **WHEN** the operator runs the export command before any successful build
- **THEN** the command fails with an error that says to run the build first

#### Scenario: Drafts stay out
- **WHEN** the content has draft articles
- **THEN** no file in the output contains a draft's title, body, or assets, and no `.md` file is present

### Requirement: Same behavior as the server version
A static site SHALL offer the same sidebar navigation, calendar, day lists, article rendering, and search results as the server version for the same build. Data files SHALL have the same shape as the API responses and SHALL be validated with the same schemas.

#### Scenario: Search parity
- **WHEN** the same query from the shared search test cases runs against the server API and against the static site for the same build
- **THEN** both return the same articles in the same order with the same snippets

#### Scenario: Invalid data file
- **WHEN** a static data file does not match its schema
- **THEN** the web app shows an error state for that view, as it does for an invalid API response

### Requirement: Direct links on a static host
Every article URL and day URL SHALL load directly on a static host, including after a browser reload, without server-side routing. Unknown paths SHALL show the app's not-found view.

#### Scenario: Reload an article
- **WHEN** a reader opens `/a/aws/s3/s3/` on the static site and reloads the page
- **THEN** the article is shown

#### Scenario: Unknown path
- **WHEN** a reader opens a path that matches no article, day, or file
- **THEN** the not-found view is shown

### Requirement: Base path
The export command SHALL accept a base path so the site works both at a domain root and under a sub-path. All links, data requests, and asset URLs SHALL include the base path.

#### Scenario: Project site sub-path
- **WHEN** the site is exported with base path `/notes/` and served under `/notes/`
- **THEN** every page, data file, asset, and search request loads without a 404

### Requirement: External video hosting
In a static export, videos SHALL be written to a separate folder outside the site output, with a manifest listing each file's relative path, size, and SHA-256 hash. Video embeds SHALL point to the configured video base URL plus that relative path. If any published article has a video and no video base URL is set, the export SHALL fail.

#### Scenario: Video URL rewrite
- **WHEN** the export runs with video base URL `https://videos.example.com/` and an article embeds `./videos/demo.mp4`
- **THEN** the site output contains no `.mp4` file, the video folder contains the file at the manifest path, and the article's video source is `https://videos.example.com/` plus that path

#### Scenario: Missing video base URL
- **WHEN** a published article embeds a video and the export runs without a video base URL
- **THEN** the export fails and the error names the article

### Requirement: Static host size limits
The export SHALL fail if the site output is larger than 1 GB, or if any single output file is larger than 100 MiB. The error SHALL name the total or the file.

#### Scenario: Site too large
- **WHEN** the published assets make the site output larger than 1 GB
- **THEN** the export fails with an error that gives the total size

#### Scenario: File too large
- **WHEN** a published asset in the site output is larger than 100 MiB
- **THEN** the export fails and the error names the file

### Requirement: Static security policy
Every HTML page of the static site SHALL declare a Content Security Policy in a `<meta>` element. The policy SHALL allow only the site's own origin, plus the configured video base URL's origin for media. It SHALL NOT allow inline scripts.

#### Scenario: CSP meta tag
- **WHEN** a reader loads any static page
- **THEN** the page has a CSP `<meta>` element, and the browser reports no CSP violation while the reader uses every feature

### Requirement: Search without a server
Search on a static site SHALL run entirely in the browser. It SHALL load the search index only when the reader first searches, and SHALL meet the search capability's matching, ranking, highlighting, and length rules.

#### Scenario: Lazy index load
- **WHEN** a reader opens an article and does not use search
- **THEN** the search index is not downloaded

#### Scenario: First search
- **WHEN** a reader types a query of three or more characters for the first time
- **THEN** the index is downloaded once and results appear without any server request other than static files
