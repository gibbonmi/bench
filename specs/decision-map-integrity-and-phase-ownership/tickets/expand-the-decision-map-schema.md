# Expand the decision-map schema

Blocked by: none

## What to build

`internal/maps` exposes the new canonical schema, direct active/compiled
candidate discovery, strict structural parsing, and the paste-ready template
without changing existing map-query consumers yet.

## Acceptance

- [x] The parser and template derive headings, status values, ticket types, and
  field names from one schema owner.
- [x] A well-formed shaping map parses, while every required-shape omission,
  duplicate terminal section, unsupported Handoff, and unsupported type earns a
  distinct diagnostic.
- [x] Research, Prototype, Grill, and Task decision tickets all parse.
- [x] Discovery includes direct active and compiled candidates and excludes
  hidden Markdown, README case variants, and nested asset Markdown.
- [x] CRLF and missing-final-newline maps retain the same meaning, and fenced
  examples do not become schema.
