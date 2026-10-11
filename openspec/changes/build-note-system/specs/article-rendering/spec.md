# Spec Delta

## Purpose

Shows a compiled article in the content area and turns its embed placeholders
into read-only diagrams, charts, API documentation, and video players.

## ADDED Requirements

### Requirement: Article view
The content area SHALL show the article title, date, tags, and compiled body. Code blocks SHALL be syntax-highlighted and SHALL have a copy button. Relative links to other articles SHALL open in the app without a full page reload.

#### Scenario: Show article
- **WHEN** the user opens an article
- **THEN** the content area shows its title, date, and body with highlighted code blocks

#### Scenario: Copy code
- **WHEN** the user clicks the copy button on a code block
- **THEN** the block text is copied to the clipboard and the button confirms it

#### Scenario: Unknown article
- **WHEN** the user opens a route for an article that does not exist
- **THEN** the content area shows a not-found message

### Requirement: Read-only embeds
Every embed SHALL be view-only. The system SHALL NOT show editing tools, edit buttons, or links that open an embed in an external editor.

#### Scenario: No editor UI
- **WHEN** an article with draw.io, Excalidraw, Mermaid, ECharts, and OpenAPI embeds is displayed
- **THEN** no embed shows an edit control, and no embed opens an external editor

### Requirement: No third-party runtime requests
While the system shows a page, the browser SHALL request scripts, styles, data, and media only from the site's own origin. The one exception is video media in a static export, which SHALL be requested only from the configured video base URL's origin. Embed renderers SHALL be served by the site.

#### Scenario: Page with every embed type
- **WHEN** an article with every embed type is displayed by the server version
- **THEN** the browser makes no request to any other origin

#### Scenario: Static site with video
- **WHEN** a static-export article with a video is played
- **THEN** the only request to another origin is for the video, and it goes to the configured video base URL's origin

### Requirement: Content area margins
The content area SHALL have padding of 5% of its width on the left and on the right, and no other width limit.

#### Scenario: Wide window
- **WHEN** an article is shown in a content area 1,000 px wide
- **THEN** the text starts 50 px from the left edge and ends 50 px from the right edge of the content area

### Requirement: Table of contents
An article view SHALL show a table of contents to the right of the text, listing the article's h1–h6 headings in order and indented by level. It SHALL stay visible while the article scrolls, SHALL mark the section being read, and selecting an entry SHALL scroll that heading into view below the search bar. An article without headings SHALL show no table of contents. In content areas narrower than 1,024 px it MAY be hidden.

#### Scenario: Long article
- **WHEN** a reader scrolls to the end of a long article with headings
- **THEN** the table of contents is still visible at the top of the content area and marks the last heading

#### Scenario: Jump to a heading
- **WHEN** the reader selects an entry in the table of contents
- **THEN** that heading scrolls into view below the search bar and the entry is marked as current

### Requirement: Light theme only
The UI SHALL use a single light color theme. It SHALL NOT switch to a dark theme, even when the operating system prefers dark mode.

#### Scenario: OS dark mode
- **WHEN** the operating system prefers a dark color scheme
- **THEN** the UI still renders with the light theme

### Requirement: draw.io embed
A draw.io embed SHALL render the diagram as a read-only viewer with zoom and page switching.

#### Scenario: Local draw.io file
- **WHEN** an article contains a draw.io embed that points to a local `.drawio` asset
- **THEN** the diagram renders inline in place of the embed

### Requirement: Excalidraw embed
An Excalidraw embed SHALL render the `.excalidraw` scene read-only, as a vector image that fits the content width.

#### Scenario: Excalidraw file
- **WHEN** an article contains an Excalidraw embed with a valid scene file
- **THEN** the scene renders as a read-only image

### Requirement: Mermaid embed
A Mermaid embed SHALL render the diagram source as an inline vector diagram.

#### Scenario: Mermaid flowchart
- **WHEN** an article contains a Mermaid fence with `flowchart LR; A-->B`
- **THEN** a diagram with nodes A and B renders in place of the fence

### Requirement: OpenAPI embed
An OpenAPI embed SHALL render the referenced OpenAPI 3.x or Swagger 2.0 document (JSON or YAML) as browsable API reference documentation. Calls to the documented API from the viewer SHALL be disabled.

#### Scenario: Swagger 2.0 JSON
- **WHEN** an article contains `![swagger](./swagger.json)` and the file is Swagger 2.0
- **THEN** the operations and models of the document render inline

### Requirement: ECharts embed
An ECharts embed SHALL render a chart from its JSON option. The chart SHALL resize with the content area.

#### Scenario: Bar chart
- **WHEN** an article contains an `echarts` fence with a bar chart option
- **THEN** a bar chart renders in place of the fence

### Requirement: MP4 video embed
A video embed SHALL render a player with the browser's standard controls. It SHALL NOT autoplay, SHALL load only metadata until the user presses play, and SHALL play inline on mobile. If the browser cannot play the file, it SHALL show a download link.

#### Scenario: Play video
- **WHEN** an article contains a video embed and the user presses play
- **THEN** the video plays inline in the article

#### Scenario: Seek
- **WHEN** the user drags the playback position to the middle of a video that is not fully downloaded
- **THEN** playback continues from the new position without downloading the whole file first

#### Scenario: No autoplay
- **WHEN** an article with a video embed is opened
- **THEN** the video does not start playing and only its metadata is loaded

### Requirement: Embed failure isolation
If an embed cannot load or render, the system SHALL show an inline error box for that embed with its type and source. The rest of the article SHALL still render.

#### Scenario: Invalid Mermaid syntax
- **WHEN** a Mermaid fence contains invalid syntax
- **THEN** that embed shows an error box and the rest of the article displays normally

#### Scenario: Corrupt Excalidraw file
- **WHEN** an Excalidraw asset is not a valid scene
- **THEN** that embed shows an error box naming the file

### Requirement: Lazy embed loading
Embed renderers SHALL be downloaded only when an article contains that embed type. Articles without embeds SHALL NOT download renderer code.

#### Scenario: Plain article
- **WHEN** the user opens an article without embeds
- **THEN** no Mermaid, draw.io, Excalidraw, OpenAPI, or ECharts renderer code is downloaded
