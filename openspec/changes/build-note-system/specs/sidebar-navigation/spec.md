# Spec Delta

## Purpose

Lets readers browse and open articles through the nested categories in
`sidebar.json` and shows where the open article sits in that tree.

## ADDED Requirements

### Requirement: Sidebar layout
The main page SHALL show a left sidebar and a content area to its right. The sidebar SHALL hold the calendar at its top and the category tree below it. The tree SHALL scroll on its own, apart from the content area.

#### Scenario: Main page layout
- **WHEN** a user opens the main page on a viewport at least 1024 px wide
- **THEN** the calendar is at the top-left, the category tree is below it, and the content area is on the right

#### Scenario: Narrow viewport
- **WHEN** the viewport is narrower than 768 px
- **THEN** the sidebar collapses behind a toggle button and the content area uses the full width

### Requirement: Nested category tree
The sidebar SHALL render `sidebar.json` as a tree with the same order and nesting. Nodes with children SHALL expand and collapse. A node with a path SHALL be a link to its article.

#### Scenario: Expand a category
- **WHEN** the user clicks a collapsed category
- **THEN** its children appear and the expand state is shown by an indicator and `aria-expanded="true"`

#### Scenario: Three-level nesting
- **WHEN** `sidebar.json` has a category inside a category inside a category
- **THEN** the sidebar shows all three levels with increasing indentation

### Requirement: Open article from sidebar
Clicking a node that has a path SHALL open its article in the content area and update the browser URL to a route that identifies the article. Back and forward navigation SHALL restore earlier articles.

#### Scenario: Click article link
- **WHEN** the user clicks an article node
- **THEN** the content area shows that article and the URL changes to the article's route

#### Scenario: Browser back
- **WHEN** the user opens article A, then article B, then presses Back
- **THEN** article A is shown again

### Requirement: Locate current article
When an article is open, the sidebar SHALL expand every ancestor category of the article's node, highlight the node, and scroll it into view. This SHALL also happen when the article was opened from the calendar, search, a link inside another article, or a direct URL.

#### Scenario: Open from search
- **WHEN** the user opens a deeply nested article from a search result while its categories are collapsed
- **THEN** the sidebar expands its ancestors, highlights the node, and scrolls it into view

#### Scenario: Article not in sidebar
- **WHEN** the open article is not referenced by any sidebar node
- **THEN** no node is highlighted and the article still displays

### Requirement: Keyboard access
The tree SHALL be operable with the keyboard: arrow keys move focus, Right and Left expand and collapse, and Enter opens the focused article.

#### Scenario: Keyboard open
- **WHEN** the user focuses an article node with the arrow keys and presses Enter
- **THEN** the article opens
