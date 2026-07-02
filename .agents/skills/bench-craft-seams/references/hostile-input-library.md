# Hostile-input library — edge-class templates per domain

Templates for the domain hostile-input checklist a project profile carries in
`projects/<name>.md`. `/bench-setup-repo` seeds a new profile from the sections
matching the stack; `/bench-write-spec`'s edge inventory reads the *profile's*
checklist, not this file — the profile is the tuned instance, this is the
quarry. Pick the matching sections, cut classes the stack can't meet, and add
the classes only this project knows about. A project spanning domains takes a
section per surface.

## Shell CLI

- paths and directory names containing spaces or glob characters
- hand-edited files whose last line lacks a trailing newline
- absent file vs present-but-empty file (distinct behaviors, both asserted)
- unquoted multi-word arguments (`$*` vs `$1`)
- required tool missing from PATH
- invocation through a symlink rather than the real path
- interrupt (SIGINT) mid-loop: leftover scratch state, locks, temp files
- re-run idempotency: second run of every state-changing command
- cwd deeper than the project root when the command assumes root

## HTTP API service

- auth: expired or revoked credential vs missing vs malformed (distinct
  responses, each asserted)
- malformed payload: invalid JSON vs valid JSON that fails the schema
- payload boundaries: empty body, oversized body, unicode/emoji in string fields
- replay: the same mutation submitted twice (retry storm, double-submit)
- concurrent writers to the same resource (lost update, version conflict)
- pagination: page past the end, page size 0 or negative, cursor to a deleted row
- upstream dependency failure: timeout vs 5xx vs connection refused (distinct
  handling, each asserted)
- clock edges: expiry at the boundary second, timestamps crossing midnight or DST
- rate limit reached mid-flow

## Web UI

- empty states: zero rows, first run, everything filtered out
- fetch states: slow response, failed response, retry after failure
- overflow: one very long unbroken string, thousands of rows
- input: pasted text with formatting and whitespace, IME/unicode entry, autofill
- double-click and double-submit on every mutating action
- stale view: acting on a record another session already changed or deleted
- navigation: back button mid-flow, deep link into an intermediate state,
  refresh with a half-filled form
- viewport: narrow screen, browser text zoom

## Background jobs / pipelines

- duplicate delivery of the same message or event
- out-of-order events (late arrival after a newer state landed)
- poison message: one unprocessable item must not stall the queue
- partial batch failure: some items succeed, the rest must not be lost or doubled
- replay/backfill overlapping live processing
- schema drift in upstream data: a new field, a missing field, a type change
- worker death mid-item: the item is neither lost nor processed twice
