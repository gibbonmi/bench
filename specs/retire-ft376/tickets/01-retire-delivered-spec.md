# Retire the delivered FT376 spec

Blocked by: none
Writes: specs/spec-stage-grader-trace, reviews/spec-stage-grader-trace.md
Covers: none

## What to build

The spec-stage-grader-trace spec is implemented, and FT376 is delivered. The `Before the map locks` section of `map-discipline.md` holds its rules, and the anchors and canaries hold its gate.
Then `bench spec retire spec-stage-grader-trace` removes the spec folder and its review record.

## Acceptance

- [ ] The spec folder and its review record are absent, and the landing commit subject ends with `spec-retire: spec-stage-grader-trace`.
