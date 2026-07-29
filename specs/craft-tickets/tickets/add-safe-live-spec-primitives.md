# Add safe live-spec primitives

Blocked by: none

## What to build

Make `internal/spec` the one owner of live folder-spec path construction and
ROADMAP-token parsing. Classify every `specs/*/spec.md` candidate before
`Facts` reads it, preserving malformed regular-file evidence while refusing
special, dangling, unreadable, and oversized control records without blocking.

## Acceptance

- [x] Story 18 and its acceptance-coverage row are green.
- [x] Story 19's path-construction and token-parser primitives have focused unit coverage.
- [x] Deadline-backed `Facts` coverage proves a FIFO candidate returns rather than blocking.
