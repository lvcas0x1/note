# Spec Delta

## Purpose

Defines the public, read-only HTTP API and asset delivery that the web app uses,
publishes the API as an OpenAPI document, and sets the protections the server
needs when it is exposed to the internet.

## ADDED Requirements

### Requirement: Read-only note endpoints
The API SHALL provide read-only JSON endpoints to get the sidebar tree, get one article by ID (metadata and compiled HTML), list per-day article counts for a date range, list the articles for one date, and search articles. The API SHALL NOT create, change, or delete content.

#### Scenario: Get article
- **WHEN** a client requests an existing article by ID
- **THEN** the API returns HTTP 200 with the title, date, summary, tags, sidebar path, and compiled HTML

#### Scenario: Unknown article ID
- **WHEN** a client requests an article ID that does not exist
- **THEN** the API returns HTTP 404 with a problem body

#### Scenario: Write method
- **WHEN** a client sends POST, PUT, PATCH, or DELETE to a note endpoint
- **THEN** the API returns HTTP 405

### Requirement: Calendar range endpoint
The day-count endpoint SHALL accept inclusive `from` and `to` dates in `YYYY-MM-DD` format. The range SHALL be at most 366 days. It SHALL return only days that have at least one article.

#### Scenario: Month range
- **WHEN** a client requests counts from 2026-10-01 to 2026-10-31
- **THEN** the API returns a date and count for each October day that has articles

#### Scenario: Range too large
- **WHEN** `to` is more than 366 days after `from`, or `from` is after `to`
- **THEN** the API returns HTTP 400 with a problem body naming the invalid parameter

### Requirement: Uniform error format
Every error response SHALL use the RFC 9457 problem details JSON format (`application/problem+json`), with at least `status`, `title`, and `detail`.

#### Scenario: Invalid date parameter
- **WHEN** a client sends `date=2026-13-40` to the per-date endpoint
- **THEN** the API returns HTTP 400 with an `application/problem+json` body

### Requirement: Published OpenAPI document
The system SHALL publish an OpenAPI 3.1 document that describes every note endpoint, parameter, response schema, and error response. It SHALL declare no security scheme. The document SHALL be produced from the API code, so it cannot drift from the implementation.

#### Scenario: Fetch the document
- **WHEN** a client requests the OpenAPI document
- **THEN** it receives a valid OpenAPI 3.1 document that lists every note endpoint

#### Scenario: Drift check
- **WHEN** an endpoint changes but the committed OpenAPI document is not regenerated
- **THEN** the project's check command fails

### Requirement: Validated responses in the client
The web app SHALL validate every API response against schemas that are checked to match the OpenAPI document. A response that does not match SHALL show an error state and SHALL NOT render partial data.

#### Scenario: Schema mismatch
- **WHEN** the API returns an article response that is missing `title`
- **THEN** the web app shows an error state for that view and logs the validation failure

#### Scenario: Client schema drift
- **WHEN** a response schema in the OpenAPI document changes but the web app's schema does not
- **THEN** the project's check command fails

### Requirement: Health endpoint
The API SHALL provide a health endpoint. It SHALL report HTTP 200 when the content database is readable, and HTTP 503 otherwise. It SHALL return no note data.

#### Scenario: No build yet
- **WHEN** the server starts before any successful build exists
- **THEN** the health endpoint returns HTTP 503

### Requirement: No authentication
All pages, API endpoints, and published assets SHALL be readable without authentication. The system SHALL NOT set cookies or store any visitor data.

#### Scenario: Anonymous request
- **WHEN** a client with no cookies requests an article
- **THEN** the API returns HTTP 200 and the response sets no cookie

### Requirement: Loopback by default
The server SHALL listen on `127.0.0.1` by default. A different listen address SHALL be used only when set explicitly in configuration.

#### Scenario: Default start
- **WHEN** the server starts without a listen-address setting
- **THEN** it accepts connections on `127.0.0.1` and refuses them on other interfaces

#### Scenario: Explicit address
- **WHEN** the listen address is configured as `0.0.0.0:8080`
- **THEN** the server accepts connections on all interfaces

### Requirement: Search abuse limits
The search endpoint SHALL reject queries longer than 200 characters with HTTP 400. It SHALL limit each client IP to a configurable number of search requests per minute (default 60) and reply HTTP 429 with a `Retry-After` header when the limit is exceeded. The client IP SHALL be taken from forwarding headers only when the request comes from a configured trusted proxy.

#### Scenario: Query too long
- **WHEN** a client searches with a 201-character query
- **THEN** the API returns HTTP 400 with a problem body

#### Scenario: Rate limit exceeded
- **WHEN** one IP sends 61 search requests within one minute with the default setting
- **THEN** the 61st request returns HTTP 429 with a `Retry-After` header

#### Scenario: Spoofed forwarding header
- **WHEN** a request not from a trusted proxy sends `X-Forwarded-For` with a different IP
- **THEN** the rate limit uses the connection's IP, not the header value

### Requirement: Asset delivery
Published assets SHALL be served only from the current build's output, with the correct content type. The server SHALL support HTTP range requests and SHALL NOT impose a size limit or a total response timeout on asset responses.

#### Scenario: Video range request
- **WHEN** a client requests bytes 1000000-1999999 of a published `.mp4`
- **THEN** the server returns HTTP 206 with exactly those bytes and `Content-Type: video/mp4`

#### Scenario: Large file
- **WHEN** a client downloads a 2 GB published video over a slow connection
- **THEN** the server keeps sending until the transfer completes

#### Scenario: Unknown asset
- **WHEN** a client requests an asset path that is not in the current build
- **THEN** the server returns HTTP 404

### Requirement: Security headers
Every response SHALL include a Content Security Policy that allows scripts, styles, connections, images, and media only from the site's own origin and forbids framing. Responses SHALL also include `Strict-Transport-Security`, `X-Content-Type-Options: nosniff`, and `Referrer-Policy: strict-origin-when-cross-origin`.

#### Scenario: Headers on article page
- **WHEN** a client requests the main page
- **THEN** the response contains all listed headers, and the CSP lists no origin other than `'self'`
