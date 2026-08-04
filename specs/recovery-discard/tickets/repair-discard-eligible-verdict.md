# Issue a distinct discard-eligible verdict after verification

Blocked by: none
Ownership fence: `internal/worktree`, `internal/contract/runtime/runtime_worktree_test.go`, `specs/recovery-discard/tickets/enrich-recovery-plan.md`, `CHANGELOG.md`
Contracts: the discard-eligible verdict crosses `internal/worktree/classifier.go`→`internal/worktree/resume.go` as authorization to delete the named ref and is asserted by DE1 and DE2 over every unclassifiable exit `PlanRecovery` can take, not over one of them; the terminal action value crosses `internal/worktree/resume.go`→the receipt renderer and `internal/contract/runtime` and is asserted by DE5 against the real `bench worktree recovery` command
Assumptions: `verifyRecovery` and `LandedInDefault` stay unchanged and remain the only proof of what a recovery ref holds; the orphaned verdict stays discardable because story 2 exists for it and it has no other route out; spec.md:55 already grants the choice between a distinct planned action and a modifier on `retain`, so this ticket needs no spec change; claims re-derived from the tree at pickup

## What to build

A critical authority defect found by the authoritative review of candidate `64fda745`,
reproduced against a checkout of that commit before this ticket was written.

`PlanRecovery` initialises `plan.Action` to `RecoveryRetain` and then returns through that
value on five different exits: ambiguous rows, a non-recovered assignment, a
`verifyRecovery` failure, an unresolvable default branch, and payloads the landedness
proof refused. Only the last of those is the operator's to drop. `recoveryVerb.authorizes`
accepts `RecoveryRetain` wholesale for `--discard`, so the other four inherit destructive
authority — the same catch-all shape as the `foreign` defect the previous round fixed.

The observed reproduction: a recovered row that `intent.ValidateAssignment` accepts, whose
payload names an OID nothing resolves, plans `retain` with detail `recovery payload is
missing or unreachable`, and `bench worktree recovery <ref> --discard <fingerprint>` exits
zero, deletes the ref, and records the terminal action `discarded`. No hand-edited state
file is required: a payload commit that was garbage-collected or rewritten reaches it.

Issue a distinct discard-eligible verdict, set **only** on the exit where `verifyRecovery`
succeeded, the default branch resolved, and the landedness proof refused the payloads.
`--discard` authorizes that verdict and the orphaned verdict, and nothing else. `retain`
becomes what it should always have meant — the plan could not classify this ref — and
authorizes neither verb, which is what spec.md:57 already requires and what acceptance row
RD7 on `add-recovery-discard.md` already asserts.

A second authority defect in the same function, on the same catch-all shape. The
interrupted-convergence branch in `applyRecoveryVerb` runs *before* the authorization
check and sets `plan.Action` to `verb.terminal()`. So a recovered row whose ref no longer
resolves lets `--apply` record `retired` — a claim that the landedness proof accepted the
payload — without that proof ever running. Convergence must still close the row, because
that is story 6's whole point, but the claim it records cannot assert a proof that did not
happen.

`enrich-recovery-plan.md`'s acceptance row EP3 advertises the old vocabulary — that an
unlanded payload "still plans `retain`" — and this ticket falsifies it. Update that row to
name the discard-eligible verdict; leave the rest of the file alone.

## Acceptance

- [ ] [DE1] a recovered row whose payload no longer resolves plans as unclassifiable and `--discard` refuses it, deleting nothing, driven through the real `RecoveryCommand` with the fingerprint that plan reported.
- [ ] [DE2] every other exit `PlanRecovery` reaches without completing verification — an unresolvable default branch, and a row the landedness proof never got to judge — likewise refuses `--discard`, asserted by enumerating those exits rather than sampling one.
- [ ] [DE3] a recovered row whose payloads verify and are genuinely unlanded plans discard-eligible and still discards, so the narrowing removes no capability story 1 exists for.
- [ ] [DE4] a genuine row-less recovery ref inside the namespace still plans orphaned and still discards.
- [ ] [DE5] the receipt distinguishes the discard-eligible plan verdict from the terminal `discarded` claim, asserted through `bench worktree recovery` in `internal/contract/runtime`.
- [ ] [DE6] an interrupted `--apply` whose ref is already gone closes the row without recording a claim that the landedness proof accepted the payload.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DE1 | add the unclassifiable verdict back to the discard authorization set | the unresolvable-payload refusal test | add that verdict to `authorizes`, run `go test ./internal/worktree -run Discard -count=1 -timeout 180s`, expect the surviving-ref assertion to fail |
| DE2 | set the discard-eligible verdict at the top of `PlanRecovery` instead of after verification | the unclassifiable-exit enumeration test | move the initialiser, run `go test ./internal/worktree -run Recovery -count=1 -timeout 180s`, expect the pre-verification exit assertions to fail |
| DE3 | leave the post-proof exit on the unclassifiable verdict | the unlanded-still-discards test | drop the discard-eligible assignment, run `go test ./internal/worktree -run Discard -count=1 -timeout 180s`, expect the successful-discard assertion to fail |
| DE4 | require a completed verification before the orphaned verdict | the orphaned-still-discards test | gate the orphan arm on `verifyRecovery`, run `go test ./internal/worktree -run Recovery -count=1 -timeout 180s`, expect the orphaned-discard assertion to fail |
| DE5 | render the terminal claim for the plan verdict as well | the runtime recovery receipt test | collapse the two values in the renderer, run `go test ./internal/contract/runtime -run Worktree -count=1 -timeout 300s`, expect the receipt-vocabulary assertion to fail |
| DE6 | restore `verb.terminal()` on the convergence branch for both verbs | the interrupted-apply claim test | set the terminal action before the authorization check, run `go test ./internal/worktree -run Retry -count=1 -timeout 180s`, expect the recorded-claim assertion to fail |
