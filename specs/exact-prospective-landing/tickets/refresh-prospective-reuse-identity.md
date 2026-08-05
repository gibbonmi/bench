# Refresh prospective reuse identity before the lock

Blocked by: reuse-exact-green-before-gate-lock.md, adopt-exact-landing-in-commit.md
Ownership fence: `internal/gate/gate.go`, `internal/gate/evaluation.go`, `internal/gate/prospective_test.go`, `internal/gate/runner_test.go`
Integration surfaces: prospective pre-lock reuse decision→`internal/gate/gate.go`; exact-tree subject generations→`internal/gate/evaluation.go`; unchanged-evidence lock avoidance→`internal/gate/prospective_test.go`; retained FT78 symlink-chain proof→`internal/gate/runner_test.go` and its existing shared proof helper
Contracts: an unpublished exact tree crosses `internal/landing`→`internal/gate.ExecuteTree` while declared local gate inputs remain rooted in the repository; PR1 proves that current resolved-input identity governs reuse, and PR2 proves that unchanged exact evidence still answers before gate-lock acquisition

## What to build

Repair exact prospective reuse so the pre-lock decision is made from a subject that
includes the current identity of repository-rooted declared inputs. A mutation behind
a declared symlink chain must invalidate the retained green, execute the prospective
gate, retain the resulting red under the refreshed oracle, and leave the destination
tip unchanged. Preserve the feature's pre-lock fast path for byte-identical exact
subjects: unchanged reusable green must still return without acquiring the gate lock
or running the gate again.

Do not weaken, rewrite, or special-case the shared FT78 symlink proof. Keep subject
capture single-sourced through the evaluation boundary; do not add an adapter-local
oracle calculation.

## Acceptance

- [ ] [PR1] `TestFT78Story3ProofLedger/R8/symlink-chain` deterministically rejects `bench commit` after the resolved target changes, observes a changed oracle, records exactly one additional red gate run, and leaves HEAD unchanged.
- [ ] [PR2] An unchanged exact prospective tree with reusable green still returns before gate-lock acquisition, emits the reuse announcement, preserves the retained verdict bytes, and does not run the gate again.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PR1 | make the pre-lock exact-tree decision consult the prospective Git tree while ignoring the repository-rooted resolved target | existing `TestFT78Story3ProofLedger/R8/symlink-chain` | seed green, mutate `inputs/target` behind the declared symlink chain, run `bench commit`, and expect the stale reuse to make the proof red |
| PR2 | move exact reusable-green inspection back behind gate-lock acquisition | existing exact prospective held-lock reuse proof in `internal/gate/prospective_test.go` | seed reusable green, hold the gate lock, execute the identical unpublished tree, and expect lock refusal or a second gate run |
