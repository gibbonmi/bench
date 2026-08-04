# Confine discard to the recovery namespace and seal the change summary

Blocked by: none
Ownership fence: `internal/worktree`, `internal/contract/runtime/runtime_worktree_test.go`
Contracts: the orphaned verdict crosses `internal/worktree/classifier.go`→`internal/worktree/resume.go` as authorization to delete the named ref and is asserted by RN1 and RN2 over refs outside the recovery namespace rather than only over refs inside it; the plan's change summary crosses the classifier→`fingerprintRecovery` in the same file and is asserted by RN4 by mutating the summary value alone
Assumptions: `intent.RecoveryRefPrefix` and `recoveryRefNamespace` in `internal/intent` remain the one definition of where recovery refs live and this ticket reads them rather than restating the literal; the retire path and `LandedInDefault` stay unchanged

## What to build

Two authority defects found by the authoritative review of candidate `2f486147`.

**The critical one.** `PlanRecovery` classifies a ref with no owning assignment row using
one piece of evidence: whether the ref exists. Nothing confines it to
`refs/bench/recovery/`. The only namespace check lives inside `verifyRecovery`, which is
reached only when a row *does* exist — exactly the branch the row-less path skips. So
`bench worktree recovery refs/heads/main --discard <fingerprint>` plans `orphaned` and
deletes an ordinary branch. A destructive verb must not accept a ref the capability was
never about.

Confine the orphaned verdict to the recovery namespace, taking the namespace from
`internal/intent` rather than writing the literal a second time. A ref outside it is not
orphaned and is not discardable; give it a verdict that says so plainly and carries no
authorization. The spec's edge inventory noticed there was no envelope check on this path
and concluded only that control characters mattered — that reasoning was wrong, and the
fix is the envelope, not more escaping.

**The second.** The plan's derived change summary is what story 3 offers the operator as
judgment support before a destructive act, but `fingerprintRecovery` does not hash it.
Mutating `paths=1` to `paths=999` leaves the fingerprint identical, so a plan plus its
fingerprint can authorize a discard after the count the operator relied on has changed.
Seal `Changes` into the fingerprint alongside the facts already sealed.

## Acceptance

- [ ] [RN1] a ref outside `refs/bench/recovery/` never plans as orphaned, however ordinary or well-formed it is.
- [ ] [RN2] `--discard` against a ref outside the recovery namespace refuses and deletes nothing, asserted by driving a real branch ref through plan and discard.
- [ ] [RN3] a genuine row-less recovery ref inside the namespace still plans orphaned and still discards, so the fix narrows the verdict without removing the capability story 2 exists for.
- [ ] [RN4] two plans differing only in their change summary produce different fingerprints.
- [ ] [RN5] the namespace constraint reads its prefix from `internal/intent` rather than restating the literal path.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RN1 | drop the namespace guard so any existing ref classifies orphaned | the namespace-confinement test | remove the guard in `PlanRecovery`, run `go test ./internal/worktree -run Recovery -count=1 -timeout 180s`, expect the ordinary-branch assertion to fail |
| RN2 | authorize the discard verb for the out-of-namespace verdict | the out-of-namespace refusal test | add that verdict to the discard authorization set, run `go test ./internal/worktree -run Discard -count=1 -timeout 180s`, expect the surviving-ref assertion to fail |
| RN3 | widen the guard to reject every row-less ref including real recovery refs | the orphaned-still-works test | require an owning row unconditionally, run `go test ./internal/worktree -run Recovery -count=1 -timeout 180s`, expect the orphaned-verdict assertion to fail |
| RN4 | remove `Changes` from the fingerprint parts list | the change-summary sealing test | delete that part, run `go test ./internal/worktree -run Recovery -count=1 -timeout 180s`, expect the differing-fingerprint assertion to fail |
| RN5 | replace the `internal/intent` prefix reference with the literal `"refs/bench/recovery/"` | the single-source namespace test | inline the literal, run `go test ./internal/worktree -count=1 -timeout 180s`, expect the single-source assertion to fail |
