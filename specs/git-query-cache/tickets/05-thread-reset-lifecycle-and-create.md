# Thread the verb's run through reset, lifecycle, and creation, and census each run construction

Blocked by: 04-thread-the-cleanup-and-release-paths.md
Writes: .bench/commitment.json, cmd/bench/commitment_test.go, internal/conformance/bounds_policy_test.go, internal/worktree/admin_readers_test.go, internal/worktree/clean.go, internal/worktree/effect_census_test.go, internal/worktree/effects.go, internal/worktree/land_run_test.go (new), internal/worktree/layers.go, internal/worktree/lifecycle.go, internal/worktree/lifecycle_acquire_test.go, internal/worktree/lifecycle_test.go, internal/worktree/ownership.go, internal/worktree/ownership_test.go, internal/worktree/parallel_census_test.go, internal/worktree/path.go, internal/worktree/pool_root.go, internal/worktree/pool_root_test.go, internal/worktree/protection_test.go, internal/worktree/release_registration_test.go, internal/worktree/reset.go, internal/worktree/reset_apply.go, internal/worktree/reset_apply_test.go, internal/worktree/reset_envelope.go, internal/worktree/reset_restore.go, internal/worktree/reset_restore_refusal_test.go, internal/worktree/reset_restore_test.go, internal/worktree/resume.go, internal/worktree/resume_reconcile_test.go, internal/worktree/resume_test.go, internal/worktree/snapshot_test.go, internal/worktree/subshell.go, internal/worktree/tree_target.go, internal/worktree/worktree_test.go, internal/worktree/land.go, internal/worktree/land_stack.go, internal/worktree/reauthorize.go, internal/worktree/show.go, internal/worktree/exec.go, internal/worktree/build.go, internal/worktree/build_restore.go, internal/worktree/clean_discard.go, internal/worktree/reset_repair_test.go, internal/worktree/worktree.go
Covers: GQ38, GQ39, GQ40, GQ42, GQ46

## What to build

Remove the last four helper constructions, so that a run construction sits only in `newAmbient` or in an exported entry.
`planReset` (`reset.go:143`) takes the run as its first parameter, as `func planReset(run *git.Run, `, and it drops its construction at `:158`.
The completion-plan probe inserts a planted helper before that exact text. The callers are `reset.go:66` and `reset_apply.go:23`.
`acquireAt` (`lifecycle.go:180`) and `releaseWith` (`lifecycle.go:258`) take the run and drop the constructions at `:196` and `:272`. `Acquire` and `Release` are the exported entries.
`validateCreationBundle` (`ownership.go:200`) takes the run and drops the construction at `:204`. Its callers are `ownership.go:118`, `:295`, `pool_root.go:100`, `:204`, and `path.go:95`.

`createAt` (`pool_root.go:80`) takes the run and passes it to `validateCreationBundle`. Its callers are `pool_root.go:65`, `tree_target.go:167`, and `subshell.go:70`.
Keep the declaration name `createAt` and its `RegisterSibling` use. The route inventory at `cmd/bench/commitment_test.go:126` reads them.
`lockCleanupRegistration` (`lifecycle.go:320`) and `releaseRegistration` (`:399`) read their locations through the run. Their callers are `lifecycle.go:367`, `reset_apply.go:18`, and `resume.go:116`.
`captureLayers` (`layers.go:71`) reads its admin directory (`:80`) through the run. Its callers are `clean.go:44`, `reset_envelope.go:11`, and `reset_restore.go:79`. `realIndexTree` reads `clean.go:127` through the run.

Each test caller passes a run or nil. The test callers are `lifecycle_acquire_test.go`, `pool_root_test.go`, `protection_test.go`, `snapshot_test.go`, `lifecycle_test.go`, and `release_registration_test.go`.
The other test callers are `admin_readers_test.go`, `ownership_test.go`, `resume_test.go`, `reset_apply_test.go`, `resume_reconcile_test.go`, `reset_restore_test.go`, and `reset_restore_refusal_test.go`.
Each helper takes the run as its first parameter after any `joins` value. No `XIn` wrapper is left.

Keep one `census.Drop(` call and one `handoffdoc.RemoveSection(` call in `lifecycle.go`, for the pins at `worktree_test.go:671` and `:679`.
Keep `bounds.LeaseStale` in `lifecycle.go` for the pin at `internal/conformance/bounds_policy_test.go:49`.
`ownership_test.go`, `resume_test.go`, `lifecycle_test.go`, and `parallel_census_test.go` are over their budgets. Change them in place and add no line.

Add `TestRunConstructionCensus` to `effect_census_test.go`. It reads each production file of the package with `go/ast` and names each run construction outside `newAmbient` and the exported functions.
A run construction is a `&git.Run{}` literal or a `new(git.Run)` call. The base holds both forms.
The census reports the file and the line of each such construction. Two cases plant one construction of each form in an unexported function, and the census names each one.
At the base this census names 9 sites. Tickets 03 and 04 remove 5 of them, and this ticket removes the last 4.
Raise `worktreeTestCount` from 879 to 880. `worktreeSerialCeiling` stays 22.

Correct the cache-scope comments at `effects.go:30-36` and `:48-50`. When this ticket lowers the landing count, set the `==` pin of `TestLandingReadsTheDestinationCommonDirOnce` to the new observed count.
The spec-stage close binds the commitment, so this ticket names `.bench/commitment.json`. At completion, the coordinator times the package five times with the spec command (GQ42).

## Acceptance

- [ ] A production file of `internal/worktree` builds a run only in `newAmbient` or in an exported function (GQ38).
- [ ] The census names the file and line of a planted run construction of each form in an unexported function (GQ39).
- [ ] At the checkpoint, `worktreeTestCount` reads 880 and `worktreeSerialCeiling` reads 22 (GQ40).
- [ ] At completion, five timed runs with the spec command give a smallest gain outside the baseline spread (GQ42, review-owned).
- [ ] The named `internal/worktree` comments state the run's current cache scope (GQ46, review-owned).
- [ ] An unexported helper with its own run construction in `reset.go` makes `TestRunConstructionCensus` fail.
- [ ] `TestLandingReadsTheDestinationCommonDirOnce` passes at its pinned count.

## Checkpoint verification

- `bench test --package ./internal/worktree`
- `bench test --package ./internal/conformance`
- `bench test --package ./cmd/bench`

Run `bench consumers` on each changed unexported function before the edit. Record the red of the probe mutation, then the green.
