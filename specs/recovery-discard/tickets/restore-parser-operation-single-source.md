# Return the parser's operation set to one source

Blocked by: repair-reclaim-history-and-partial-apply.md
Ownership fence: `internal/spec`, `internal/specbuild/reclaim_test.go`
Contracts: the parser's operation set crosses `internal/spec/build.go`→`internal/spec/build_test.go` and is asserted by PS1 by enumerating the grammar table itself rather than restating a count in the test's name
Assumptions: `repair-reclaim-history-and-partial-apply.md` has landed so `internal/specbuild/reclaim_test.go` is settled and this ticket only removes the relocated grammar cases from it; the reclaim grammar and dispatch already shipped and this ticket changes no production behaviour

## What to build

The canonical owner of the parser's operation set — `internal/spec/build_test.go` — still
names itself `TestParseBuildExposesExactlyEightOperations`, omits `reclaim` from its case
table, and probes an unknown operation with the token `"ninth"`. Meanwhile the reclaim
grammar cases sit in `internal/specbuild`, carrying a comment that explains they are there
because of an ownership fence.

That fence was drawn wrong by the coordinator and the compromise is now the defect: one
fact — which operations the parser exposes — has two homes, and the canonical one is
already stale. The comment records build provenance rather than timeless rationale, which
is not what a code comment is for.

Move the reclaim grammar cases to the parser's own test, fold `reclaim` into the existing
enumeration, and delete the relocated cases and their provenance comment from
`internal/specbuild/reclaim_test.go`. Leave that file's genuine service-level reclaim
coverage alone — only the `ParseBuild` cases move.

Prefer an enumeration derived from the grammar table over a count restated in a test name.
A name carrying a number goes stale the next time an operation is added, which is exactly
how this defect arrived.

## Acceptance

- [ ] [PS1] the parser's own test covers every operation the grammar table declares, derived from the table rather than from a hand-maintained list or a count in the test's name.
- [ ] [PS2] the bound-arity and blamed-token cases for `reclaim` live in `internal/spec` and still assert the blamed token, not merely the exit code.
- [ ] [PS3] `internal/specbuild/reclaim_test.go` no longer contains `ParseBuild` grammar cases or the fence-provenance comment, and its service-level reclaim coverage is unchanged.
- [ ] [PS4] adding a hypothetical operation to the grammar table without touching the test makes the enumeration fail, so the single source is enforced rather than merely tidied.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PS1 | hand-list the operations in the test instead of deriving them from the grammar table | the enumeration test | replace the derivation with a literal slice omitting `reclaim`, run `go test ./internal/spec -count=1 -timeout 120s`, expect the coverage assertion to fail |
| PS2 | assert only the exit code for the two-slug case | the blamed-token test | drop the message comparison, set `reclaim`'s `MaxArgs` to `-1`, run `go test ./internal/spec -count=1 -timeout 120s`, expect the mutation to survive green |
| PS3 | leave the relocated `ParseBuild` cases in place in both packages | the duplication check | restore them, run `rg -n "ParseBuild" internal/specbuild`, expect the parser's operation set to have two homes again |
| PS4 | add a tenth operation to the grammar table and run the suite unchanged | the enumeration test | insert the row, run `go test ./internal/spec -count=1 -timeout 120s`, expect the enumeration to fail naming the uncovered operation |
