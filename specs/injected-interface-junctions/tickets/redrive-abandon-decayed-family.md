# Re-drive the decayed-family abandon tests through the real planner

Blocked by: none
Ownership fence: `internal/specbuild/abandon_test.go`
Contracts: the decayed-shape plan fingerprint and refusal cross `internal/specbuild`→`internal/worktree` through the real `worktree.PlanAbandon`/`ApplyAbandon`, asserted by AD1 and AD2 against the real planner
Assumptions: production is presumed correct — a red against unmodified production is a stop-and-surface finding, never forced green; the `decayedOwner` fake stays for the fast path but its comment is updated to name the junction tests that now compose the real contract; the unreadable-metadata fixture keeps its root-guard skip via `requireUnreadableMetadata`; claims re-derived from the tree at pickup

## What to build

The decayed/husk/unreadable abandon shapes — today graded only against
`decayedOwner`'s synthetic fingerprint — also driven through the real planner
(the existing `abandonOwner` real-delegation pattern) with privilege-free
fixtures: dir-without-.git (decayed), non-directory path (husk), chmod-0
metadata (unreadable, root-guarded).

## Acceptance

- [ ] [AD1] the decayed-directory and husk abandon paths complete through the real planner with the Service-observable outcomes the fake-driven twins assert.
- [ ] [AD2] the unreadable-metadata abandon path completes through the real planner under the root-guarded chmod fixture, restoring permissions in cleanup.
- [ ] [AD3] the `decayedOwner` comment names the junction tests that now compose the real contract instead of claiming the contract lives only in internal/worktree.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AD1 | soften the real planner locally so a shape `ClassifyPathShape` refuses is admitted | the junction tests | apply the mutation in a scratch tree, run `go test ./internal/specbuild -run <AD1 tests>`, expect the junction red while the fake-driven twins stay green; revert |
| AD2 | remove the chmod from the fixture so the metadata is readable | the junction test | apply the mutation, run the AD2 test, expect the asserted unreadable-path outcome to go missing |
| AD3 | restore the old comment text | grep tripwire in review | `rg "internal/worktree's contract" internal/specbuild/abandon_test.go` returns no hit after the ticket |
