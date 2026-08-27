# Derive the reporter directory without an entry

Blocked by: 10-single-source-the-cache-refusals-and-trim-the-comments.md

Writes: internal/gocache/gocache.go, internal/gocache/gocache_test.go, internal/gate/report.go, internal/gate/report_test.go, specs/go-build-cache-footprint/spec.md

## What to repair

Review finding (Coverage 3), reviewer decision 2026-08-27. A phase runner
launched from a plain shell carries no `GOCACHE` entry. The reporter then
prints an empty directory and logs an empty path. `FromEnv` answers the `GOCACHE`
entry when the slice carries one, and otherwise the `HOME` derivation. When
neither exists the reporter prints no line and logs no event, because it has
no directory to name.

Add row R18 to the spec's acceptance coverage map beside R17.

## Acceptance

- [ ] R18 — A reporter env with no `GOCACHE` entry names the `HOME`-derived directory in the line and the event.

Delivered outcome: every run that reaches a verdict names the Bench build
cache, whichever way the runner was launched.
