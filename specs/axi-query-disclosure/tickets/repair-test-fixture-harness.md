# Repair test fixture harness

Blocked by: none — all five blocker tickets are landed (b18fa260, d1b2ef19, 5d71b964, 69492c70, 93fbe8e9) and the fence amendment naming `internal/gittest/` is reviewer-approved in the spec's `## Ownership fences`
Writes: `internal/gittest/` (new), `internal/learnings/learnings_test.go`, `internal/maps/maps_test.go`, `internal/guards/guards_test.go`, `internal/worktree/` test files, `cmd/bench/command_registry_test.go`

## What to build

The contract step of the expand/migrate/contract sequence for review finding
R2: collapse the five independently derived minimal git test-repo scaffolds —
`learningsRepo`, `mapsRepo`, the twelve inline `git init` blocks in
`internal/guards/guards_test.go`, `newAXIEnvelopeRepo`, and (explicit
migration target, resolving the reviewer's ambiguity in the widening
direction; contestable) the pre-existing `newWorktreeRepo` — plus any fixtures
the blocker tickets add, onto one shared scaffold.

Coupled decisions for the reviewer, to be taken as one: (a) the fence
amendment naming the scaffold's home, argued against the
`repair-test-support-fence.md` precedent, whose whole content is REMOVING a
top-level test-support package — this proposal asks for the exception that
ticket declined; (b) the parked census-scope IDEA, because a non-`_test.go`
scaffold makes these repository constructors permanently invisible to
`architectureOwnedTest`. R1 (census scope) rides the parked IDEA whole; this
ticket carries no repo-spawn-reduction clause (no oracle exists for it).

## Acceptance

- [ ] [FH1] (covers local) exactly one scaffold source remains; all five
  named derivations and the blocker tickets' fixtures consume it; no test
  expectation bytes change.
- [ ] [FH2] (covers local) `bench preflight review axi-query-disclosure`
  reports `paths-authorized` green under the amended fences before a new
  candidate snapshot is captured.
