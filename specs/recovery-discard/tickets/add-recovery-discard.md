# Add a deliberate per-ref recovery discard

Blocked by: enrich-recovery-plan.md
Ownership fence: `internal/worktree`, `internal/usage/worktree.go`, `internal/contract/runtime/runtime_worktree_test.go`, `internal/conformance/docs_workflow_helpers_test.go`, `bin/bench.sh`
Contracts: the plan fingerprint crosses `internal/worktree/classifier.go`→`internal/worktree/resume.go` carrying the new destructive authority, asserted by RD2 against a fingerprint produced by real `PlanRecovery`; the discard-eligible verdict crosses the classifier→`RecoveryCommand`, asserted by RD3 and RD5 over the whole enumerated action set rather than one sampled action
Assumptions: `enrich-recovery-plan.md` has landed, so orphaned and absent are already distinct plan verdicts; `LandedInDefault` and the retire path are unchanged; the Git objects behind a discarded ref stay reachable through the reflog, so no extra audit record is written; claims re-derived from the tree at pickup

## What to build

`bench worktree recovery <ref> --discard <fingerprint>` retires a recovery ref
whose payload the landedness proof does not accept. It mirrors the `--apply`
contract exactly: plan first, the exact fingerprint the plan just reported, one
ref per invocation, no glob and no set-selecting predicate. Discarding deletes
the ref and closes the assignment row through the same compaction `--apply`
uses, without asserting the payload landed.

Three refusals fail closed. A stale fingerprint changes nothing. A ref whose
plan action is `retire` is refused, because `--apply` is its route and the two
claims must stay distinguishable in the receipt. A fingerprint that is not 64
lowercase hex characters is refused before planning — and it is refused by the
*same* parse the `--apply` arm already reaches, generalized over the verb, not
by a second copy of the format knowledge.

Two contracts change with the new authority. The fingerprint's domain tag and
effect string must name the destructive discard, so a fingerprint planned under
the old retire-only semantics cannot authorize it. And `ApplyRecovery`'s silent
exit-zero return on a `retain`-action plan becomes a refusal for **both** verbs:
a caller who supplied a fingerprint asked for an action, and reporting success
when nothing happened reads as "the work is gone" when it is not. This is a
deliberate behaviour change to `--apply`.

An orphaned ref — one that exists with no owning row — is discardable on the
strength of the ref alone, and discarding it touches no intent record. An absent
ref is a no-op success, so a re-run of a completed discard converges rather than
failing.

## Acceptance

- [ ] [RD1] discarding an unlanded payload with the exact planned fingerprint deletes the recovery ref and closes the assignment row through the same compaction `--apply` uses.
- [ ] [RD2] a stale fingerprint refuses, leaving the ref and the assignment row byte-identical.
- [ ] [RD3] `--discard` against a ref whose plan action is `retire` refuses, and the refusal is reachable for every plan action the classifier can emit, enumerated rather than sampled.
- [ ] [RD4] a fingerprint that is not 64 lowercase hex characters refuses before planning when supplied to `--discard`, reaching the same single parse the `--apply` arm reaches.
- [ ] [RD5] discarding an orphaned ref deletes the ref, writes no intent record, and a `--discard` against an absent ref is a no-op success.
- [ ] [RD6] a successful discard emits a terminal action in the receipt distinct from the retire path's, so the two claims about the same disappearance are separable.
- [ ] [RD7] supplying a fingerprint with a `retain`-action plan refuses for `--apply` as well as `--discard`, replacing the current silent success.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RD1 | delete the ref but skip the assignment-row compaction | the discard end-state test | return before the compaction branch in the discard path, run `go test ./internal/worktree -run Discard -timeout 120s`, expect the closed-row assertion to fail |
| RD2 | compare only the fingerprint's first eight characters | the stale-fingerprint refusal test | truncate both sides of the comparison, run `go test ./internal/worktree -run Discard -timeout 120s`, expect the stale-refusal assertion to fail |
| RD3 | accept a `retire`-action plan in the discard arm | the retire-action refusal test | drop the action guard, run `go test ./internal/worktree -run Discard -timeout 120s`, expect the refusal assertion to fail |
| RD4 | give `--discard` its own arm that skips the hex/length/lowercase check | the malformed-fingerprint refusal test | add the duplicate arm in `RecoveryCommand`, run `go test ./internal/worktree -run Recovery -timeout 120s`, expect the exit-2 assertion to fail |
| RD5 | make the absent-ref discard return the invocation error instead of success | the absent no-op test | return the error, run `go test ./internal/worktree -run Discard -timeout 120s`, expect the no-op-success assertion to fail |
| RD6 | emit the existing retired action for a discard | the receipt-action test | reuse `RecoveryRetired` in the discard path, run `go test ./internal/worktree -run Discard -timeout 120s`, expect the distinct-action assertion to fail |
| RD7 | restore the silent nil return for a `retain`-action plan | the fingerprint-with-retain refusal test | return the plan and nil before the action switch, run `go test ./internal/worktree -run Recovery -timeout 120s`, expect the refusal assertion to fail for both verbs |
