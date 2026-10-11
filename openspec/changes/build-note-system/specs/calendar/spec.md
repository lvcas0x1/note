# Spec Delta

## Purpose

Lets readers browse articles by date. A month calendar marks the days that
have articles and opens a list of that day's articles.

## ADDED Requirements

### Requirement: Month calendar
The sidebar SHALL show a month calendar at its top-left. When the page loads, it SHALL show the current month in the browser's local time zone. Previous-month and next-month controls SHALL change the month shown.

#### Scenario: Initial month
- **WHEN** a user opens the main page on 2026-10-11
- **THEN** the calendar shows October 2026 and marks the 11th as today

#### Scenario: Change month
- **WHEN** the user clicks the next-month control
- **THEN** the calendar shows the next month and its article markers

### Requirement: Article date markers
Each day that has at least one article SHALL be marked and SHALL show the article count in its accessible label. Days without articles SHALL NOT be marked. An article's date is a calendar date with no time zone, so it SHALL appear on the same day in every time zone.

#### Scenario: Day with articles
- **WHEN** two articles have `date: 2026-10-05`
- **THEN** October 5 is marked and its accessible label includes "2 articles"

#### Scenario: Day without articles
- **WHEN** no article has the date 2026-10-06
- **THEN** October 6 is not marked

### Requirement: Day article list
Clicking a day SHALL show, in the content area, a list of every article with that date. Each entry SHALL show the title and summary. The list SHALL be sorted by title. The URL SHALL change to a route that identifies the day, so the list can be bookmarked.

#### Scenario: Click a marked day
- **WHEN** the user clicks October 5 and it has two articles
- **THEN** the content area lists both articles with title and summary, and the URL identifies 2026-10-05

#### Scenario: Click an empty day
- **WHEN** the user clicks a day with no articles
- **THEN** the content area shows an empty-state message for that date

### Requirement: Open article from day list
Clicking a title in the day list SHALL open the full article in the content area. Back navigation SHALL return to the day list.

#### Scenario: Open and return
- **WHEN** the user clicks an article title in the day list and then presses Back
- **THEN** the full article opens, and after Back the day list is shown again

### Requirement: Calendar follows open article
When an article opens, the calendar SHALL switch to the article's month and highlight its date.

#### Scenario: Open older article
- **WHEN** the user opens an article dated 2025-03-14 from the sidebar
- **THEN** the calendar shows March 2025 and highlights the 14th
