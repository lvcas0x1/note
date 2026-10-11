# Spec Delta

## Purpose

Lets readers find articles by keyword from a search bar that is always
visible, using full-text search over titles, summaries, tags, and bodies.

## ADDED Requirements

### Requirement: Persistent search bar
A search input SHALL stay visible at the top center of the content area on every page, including while the content scrolls. Pressing `/` or `Ctrl+K` / `Cmd+K` SHALL focus it.

#### Scenario: Search bar while scrolling
- **WHEN** the user scrolls to the bottom of a long article
- **THEN** the search bar is still visible at the top center of the content area

#### Scenario: Keyboard shortcut
- **WHEN** the user presses `/` while focus is not in a text field
- **THEN** the search input receives focus

### Requirement: Keyword full-text search
Search SHALL cover only published (non-draft) articles. It SHALL match article titles, summaries, tags, and body text, ignoring case. Partial words of three or more characters SHALL match, including Japanese and other CJK text. Results SHALL be ranked by relevance, and title matches SHALL rank above body-only matches.

#### Scenario: Title match ranks first
- **WHEN** the user searches "kafka", article A has "Kafka" in its title, and article B mentions it only in its body
- **THEN** A is listed before B

#### Scenario: Japanese substring
- **WHEN** the user searches "本語の" and an article body contains "日本語のテキスト"
- **THEN** that article is in the results

#### Scenario: Query too short
- **WHEN** the query has fewer than three characters after trimming
- **THEN** no search request is sent and a hint asks for more characters

### Requirement: Search results panel
While the user types, results SHALL appear in a dropdown under the search bar, updated after input pauses for 200 ms. Each result SHALL show the title, date, sidebar category path, and a snippet with the matched terms highlighted. At most 20 results SHALL be shown.

#### Scenario: Show results
- **WHEN** the user types "oauth" and stops typing
- **THEN** a dropdown lists matching articles with highlighted snippets within 500 ms on a content set of 1,000 articles

#### Scenario: No results
- **WHEN** no article matches
- **THEN** the dropdown shows a no-results message

### Requirement: Open search result
Choosing a result with a click or with arrow keys and Enter SHALL open the article and close the dropdown. Escape SHALL close the dropdown.

#### Scenario: Keyboard selection
- **WHEN** the user presses Down twice and then Enter in the results dropdown
- **THEN** the second result's article opens

### Requirement: Search query safety
Search input SHALL be treated as plain text. Operator characters SHALL NOT cause errors and SHALL NOT change the query structure.

#### Scenario: Special characters
- **WHEN** the user searches `"foo* OR (bar`
- **THEN** the search returns results or an empty list, not an error
