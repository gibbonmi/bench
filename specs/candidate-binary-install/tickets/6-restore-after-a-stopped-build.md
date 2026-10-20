# Restore the main executable after a declared build stops

Blocked by: 2-install-the-candidate.md, 3-hold-the-refresh.md, 5-refuse-an-uninstalled-writer.md
Writes: internal/worktree/build.go, internal/worktree/build_restore_test.go (new), internal/worktree/joins.go, internal/worktree/list.go, internal/worktree/journey_test.go, internal/worktree/parallel_census_test.go, internal/usage/worktree.go, internal/commitment/repository/window.go, internal/intent/intent.go, internal/intent/transaction.go, internal/intent/transaction_test.go, cmd/bench/help_inventory_test.go, internal/conformance/refusal_route_bypass_test.go, internal/worktree/build_restore.go (new), internal/spec/spec.go, internal/refusalroute/faces_gate.go
Covers: CB37, CB38, CB39, CB40, CB41, CB42, CB43, CB49, CB64

## What to build

Deliver CB-C4: `bench worktree build --restore <slug>` clears a stopped build's declared ledger fields under the installed candidate and reinstalls the `main` executable.

Add the `--restore` flag to the build grammar, with the spec slug as its positional argument. Put the restore form in the same `internal/usage/worktree.go` constant family as ticket 2's install form, and keep `usage.WorktreeBuild` unchanged. Take the home and the clock from the package's ambient values, not from a direct environment read, so the effect census stays green.

The verb refuses before any write, in this order:

1. The repository declares no Bench build inputs.
2. The default-branch copy of `specs/<slug>/spec.md` declares no persisted field. Read it through ticket 2's exported default-branch declaration reader. Print ticket 1's `persisted-fields-invalid` error when the scan refuses.
3. The running executable's strict ledger read fails. Its route is `bench worktree build --install <target>` while the target worktree exists, then a rerun of the restore. The restore looks the label up through ticket 2's tolerant read entry and prints the literal `<target>` slot only when that lookup fails.
4. The ledger holds an open build of any outcome, per ticket 3's any-open-build answer, not its declared-build answer. The refusal names `bench commitment block`, and it tells the operator to run the restore before `bench worktree release`.

As in ticket 2, each refusal route literal stays in `build.go` or composes through the refusal-route registry, and the `build.go` entry of the route-bypass allowlist names the restore routes.

Then copy the ledger to `$BENCH_HOME/restore/<repo-key>/`, with `poolkey.Key` as the repo key, and refuse when the copy fails. Remove each declared `intent-ledger` field that the ledger holds through one exported raw-field removal in `internal/intent`. The removal edits the raw document under the ledger lock, deletes the last path segment from each element along the path, and runs the strict decode on the result before the replace. The same removal serves a top-level path such as `commitment.gaps` and a nested path such as `commitment.gaps.chunk`, which has no `omitempty` tag. No typed path-to-field mapping exists beside it. `intent.go` has 41 lines of headroom. The restore's unit tests take injected home values and seams, because the journey-harness census forbids `t.Setenv`, `exec.Command`, and `gittest.` there.

Last, rebuild the primary `dist/bench` through the existing `buildSubject` seam, and run `doctor --fix` under the rebuilt executable with ticket 2's child environment. Raise the package test-count pin by one for each new top-level test.

Two items from the CB-C3 review also land here, because this ticket writes both files. First, collapse `wrapperDir` and `manifestPath` in `internal/worktree/build.go` onto `brokermanifest.Path`, so one source owns the primary manifest location. Second, correct two `internal/commitment/repository/window.go` comments. The `CheckEntry` doc omits the case where the primary checkout does not resolve. The declared-build count comment says only that an unreadable ledger lists no build.

## Acceptance

- [ ] `TestBuildRestore` shows the ledger copy under `$BENCH_HOME/restore/<repo-key>/` before the ledger changes (CB37).
- [ ] The restore removes the `gaps` key for a declared `intent-ledger:commitment.gaps`, and removes `chunk` from each `gaps` element for a declared `intent-ledger:commitment.gaps.chunk` (CB38, CB39).
- [ ] After the restore, the primary `dist/bench` verifies against the primary sources, and the manifest digest equals its digest (CB40).
- [ ] The restore refuses before any write on a failed strict ledger read, on an open build whose chunk gap is open, and on a spec with no declaration (CB41, CB42, CB43).
- [ ] After the rebuild, the restore calls the same child seam once with ticket 2's argv, working directory, and environment, under the rebuilt executable (CB64).
- [ ] The restore of a spec whose declared field the ledger does not hold completes and leaves a verified primary executable (CB49).
- [ ] `build.go` derives the primary manifest location only through `brokermanifest.Path`, and the two `window.go` comments match the code.

## Checkpoint verification

Run CB-C4's verification command from the spec's completion plan. Probe the refusal order (strict read before the open-build test) and the nested removal with one mutant each, and record the reds. The orchestrator then freezes the CB-C4 delta for independent review.
