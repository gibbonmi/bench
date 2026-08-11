# Repair: split the oversized test file and make PF17's seam test real

Blocked by: repair-fence-grammar-edges.md
Ownership fence: `internal/preflight/`
Integration surfaces: structure ceiling→existing `bench structure` 400-line rule (observed, not edited); split files→`internal/preflight/`
Contracts: none crosses — a test-file reorganization plus two test additions inside one package
Closure: RT1/under-ceiling, RT2/newline-seam-real, RT3/status-branch-tested

## What to build

Review findings S2, V1, V3. Split `internal/preflight/command_test.go`
(772 lines vs the 400 ceiling) into cohesive files along the mode seam —
review-contract tests, build/bootstrap-contract tests, shared seeding harness —
per `craft-seams`' split-or-grant rule (split, not grant: the seams are clean).
Replace the tautological `TestDecideTrailingNewlineParity`
(`Decide(x)==Decide(x)`) with direct gatherer-level tests over unterminated
content: `fenceTokens` and the ticket-token scan each parse a final line
lacking `\n` identically to its terminated form — PF17's declared seam made
real. Add the missing test for the `spec status not readable` bootstrap branch
(`gather.go:115-117`). No production-code changes.

## Acceptance

- [ ] [RT1] (covers local) `bench structure` reports no FILE TOO LONG for `internal/preflight/`, and the full package suite passes unchanged after the split.
- [ ] [RT2] (covers PF17) unterminated-final-line fence and ticket content parse identically to terminated content, asserted directly at the gatherer seam.
- [ ] [RT3] (covers local) the status-not-readable bootstrap branch is exercised by a test reaching its structured error.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RT1/under-ceiling | (structural) merge the split files back over 400 lines | `bench structure` | run it, expect the FILE TOO LONG report |
| RT2/newline-seam-real | key the fence scanner to `\n` so the unterminated last line drops | the new gatherer-level parity test | `go test ./internal/preflight -run Newline`, expect the dropped-token failure |
| RT3/status-branch-tested | swallow the not-readable error into a generic bootstrap red | the new status-branch test | run it, expect the missing-diagnostic failure |
