# Review pickup: prospective-artifact-recovery

Base: c2b767d81179ce0bed5a8774accdd279e0f1c85d
Reviewed tip: 88441b81aee86656431aa5ef7445f94d375b4921

## Standards

Count: 7. Worst: the owner-record wire shape has four derivations across the test files.

- S1 `internal/gate/prospectiveartifact/prospectiveartifact.go:149-152` — The `EPERM` retain arm is subsumed by the `!ESRCH` arm; "only ESRCH is dead" has two derivations. Rule: AGENTS.md "one source per fact". Fix: keep one `!ESRCH` guard. Disposition: auto-fix.
- S2 `prospectiveartifact_test.go:33`, `prospective_owner_test.go:390`, `owner_artifact_recovery_test.go:160`, `:389` — The record shape (schema, owner_pid, common_dir, 0600, strict decode) is derived four times against `ownerRecord`. Rule: one source per fact. Fix: export the record type, the prefix, and one read/plant helper; decode through it. Disposition: auto-fix.
- S3 `prospective_owner_test.go:335`, `owner_artifact_recovery_test.go:274`, `:378`, `land_freshness_test.go:184` — The `bench-prospective-artifact-` prefix literal repeats `bundlePrefix`. Rule: one source per fact. Fix: export the prefix. Disposition: auto-fix.
- S4 `prospectiveartifact_test.go:108-112` — The body of `requireCheckoutUnregistered` (`:220-226`) is inlined. Fix: call the helper. Disposition: auto-fix.
- S5 `internal/worktree/land_freshness_test.go:184` — The `bench-gate-subject-` arm matches a prefix no producer creates. Fix: drop the arm. Disposition: auto-fix.
- S6 `internal/gate/prospective.go:26` — `prospectiveRunBinaryOwner` has no production caller (Lazy Element); its comment describes the callee's policy. Fix: inline at the two test sites and move the comment. Disposition: auto-fix.
- S7 `internal/gate/lane_test.go:246` — "still on the defer-only helper" narrates a prior state absent from the tree. Rule: craft-comments. Fix: state the mutation only. Disposition: auto-fix.

## Spec

Count: 5. Worst: an authored run binary can escape the bundle.

- P1 `internal/gate/prospective.go:30-40` — Spec: "contains a fixed checkout child and every owner-authored run-binary directory." The `ReuseOrOwn` branch carries no `TempRoot`, so a candidate without its own kit authors `bench-run-*` outside the bundle. Fix: set `TempRoot` on both branches, as `lane.go:253` does. Disposition: auto-fix.
- P2 `internal/systemtest/owner_artifact_recovery_test.go:300-305` — Spec seam "concurrent child-process authorization" is unreachable: `acquireExecutionLock` (`run_transaction.go:219-228`) refuses a second authorization per repository. The tests use a live OS process beside a planted registered bundle; the sweep input is identical. Non-behavioral contradiction; flag for reviewer veto. Disposition: no-op.
- P3 PAR04 `owner_artifact_recovery_test.go:87-95` — Partial: the second fresh sweep asserts zero bundles, not "changes no path". Fix: snapshot the temp root and the Git worktree list before and after the second sweep. Disposition: auto-fix.
- P4 PAR21 `TestCloseConfinesRemovalToItsBundle` — Partial: the hostile path has no Git registration, so the "real Git lifecycle seam" is untested. Fix: register the checkout in the hostile-path bundle. Disposition: auto-fix.
- P5 PAR33 `TestOpenRefusesRecoveryWhenRegistrationRemovalFails` — Partial: it never asserts that no new bundle or checkout was created. Fix: assert the bundle count and the Git worktree list are unchanged after the refusal. Disposition: auto-fix.

## Coverage

Count: 3. Worst: a symlinked temporary root leaks the Git registration forever.

- C1 `prospectiveartifact.go:64-66`, `:223` — Input: a temporary root with a symlink component (macOS `/var` → `/private/var`). Git registers the resolved path, so the `filepath.Clean` comparison in `removeBundle` never matches. `RemoveAll(root)` then deletes the bundle the next sweep needs, and story 3 fails. Fix: resolve the base with `EvalSymlinks` in `Open` and `sweep`. Add a PAR03 sibling row through a symlinked base; assert against raw `git worktree list` output, not the production comparison. Disposition: auto-fix.
- C2 `readRecord:164` — No test plants a 0644 or 0400 record; deleting the mode clause turns nothing red. Fix: add the case to `TestSweepRetainsABundleWithANonRegularRecord`. Disposition: auto-fix.
- C3 `Open:75-81` — SIGKILL between `MkdirTemp` and `writeRecord` leaves a recordless prefix directory that no sweep removes. The Won't-handle names a legacy unrecorded tree, not this window. Recommendation: accept as residual (an empty directory). Disposition: ask-user.

## Repair targets (de-duplicated)

1. `prospective.go` run-binary owner: P1 + S6.
2. `prospectiveartifact.go` and its tests: C1 + S1 + C2 + P4 + P5 + S4.
3. Test single-sourcing: S2 + S3 + S5 + S7.
4. Process rows: P3.
