# Refuse a map whose row IDs carry two tags

Blocked by: none
Writes: internal/coverage/coverage.go, internal/coverage/coverage_schema_test.go
Covers: CE12, CE13

## What to build

A coverage map declares one row-ID tag. The delivered membership check in the
preflight scopes to the tag of the first row, so a row with a foreign tag escapes
that check. The one row-ID validator in the coverage package now refuses a map
whose row IDs carry more than one tag.

The validator collects the alphabetic tag of each row ID. It reports a violation
when the set holds more than one tag. The violation names each tag it found, so a
build delegate sees which rows disagree. A malformed row ID keeps its delivered
message and contributes no tag.

Extend the one validator. No consumer re-derives the map structure.

## Acceptance

- [ ] CE12 — a map whose row IDs carry two tags is a violation.
- [ ] CE13 — the mixed-tag violation names each tag it found.
- [ ] a map whose row IDs carry one tag stays green.
- [ ] `bench gate` stays green.

## Delegate charge

You work in the Bench repo on the `citation-execution-proof` spec. Read
`specs/citation-execution-proof/spec.md` first. Then read the row-ID validator in
`Check`, in `internal/coverage/coverage.go`. Also read
`internal/coverage/coverage_schema_test.go`.

Extend that one validator. Collect the alphabetic tag prefix of each well-formed
row ID. Report one violation when the map holds more than one tag. Name every tag
you found in the message, in a stable order.

Leave a malformed row ID on its delivered message. Do not let a malformed ID
contribute a tag. Do not add a second row-ID parser. Do not change the `rowIDs`
or `ParseSpec` signatures.

Add `TestMixedTagRowIDs` in `internal/coverage/coverage_schema_test.go`. Assert the
exact violation text for a two-tag map. Also assert that a one-tag map stays
silent.

Run only `bench worktree exec <label> -- go test ./internal/coverage/`. Do not
commit. Do not edit the spec.
