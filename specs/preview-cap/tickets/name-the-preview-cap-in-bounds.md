# Name the preview cap in the bounds registry

Blocked by: none
Writes: internal/bounds/bounds.go, internal/sanitize/sanitize.go, internal/sanitize/sanitize_test.go, internal/conformance/bounds_policy_test.go, internal/shift/objective.go, internal/shift/objective_test.go, CHANGELOG.md
Covers: none

## What to build

`sanitize.Preview` caps a preview at a bare literal of 120 code points, and the file
writes that literal two times. The cap is a resource policy, so it belongs in the
`internal/bounds` registry with the other bounds.

Add `PreviewRuneLimit = 240` to the `internal/bounds` const block. `Preview` reads the
constant for the comparison and for the slice bound, and the doc comment names the
constant. The new value gives a reader two times more of an operator string before the
byte-count suffix replaces the rest.

The `bounds-policy` check owns the new entry. The check requires the name in the
registry, and it requires `internal/sanitize/sanitize.go` to consume
`bounds.PreviewRuneLimit`. A future edit that returns the literal to `Preview` makes the
check red.

## Acceptance

- [ ] `Preview` returns a 240-rune string unchanged.
- [ ] `Preview` caps a 241-rune string at 240 runes and appends `… (482 bytes)`.
- [ ] `bench test --check bounds-policy` is green.
- [ ] `bench test --check bounds-policy` is red when `sanitize.go` does not name `bounds.PreviewRuneLimit`.
- [ ] `bench test --package ./internal/sanitize` is green.
- [ ] `bench test --package ./internal/testreport` is green.
