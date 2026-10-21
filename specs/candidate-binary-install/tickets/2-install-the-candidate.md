# Install a declared build's candidate as the repository executable and broker

Blocked by: 1-declare-persisted-fields.md
Writes: internal/worktree/build.go, internal/worktree/build_install_test.go (new), internal/worktree/build_test.go, internal/worktree/joins.go, internal/worktree/list.go, internal/worktree/journey_test.go, internal/worktree/verb_runner_test.go, internal/worktree/verb_runner_check_test.go, internal/worktree/parallel_census_test.go, internal/usage/worktree.go, internal/intent/transaction.go, internal/intent/transaction_test.go, internal/commitment/repository/sources.go, internal/commitment/repository/window_test.go, cmd/bench/help_inventory_test.go, internal/conformance/refusal_route_bypass_test.go, internal/systemtest/candidate_install_test.go (new), internal/systemtest/owner_test.go, internal/commitment/commitmenttest/admission.go, internal/worktree/path.go, internal/worktree/exec.go
Covers: CB13, CB14, CB15, CB16, CB17, CB18, CB19, CB20, CB21, CB23, CB54, CB55, CB56, CB63

## What to build

Deliver the install half of CB-C2: `bench worktree build --install <target>` installs a declared build's candidate as the primary checkout's executable and broker.

Add the `--install` flag to the build grammar in `build.go`. Keep the existing `usage.WorktreeBuild` text unchanged, because `WorktreeBuildFor` derives every rebuild remedy from it. Put the install form in a separate constant in `internal/usage/worktree.go`; add it to `buildGrammar.Help`, which prints the `--help` output. The leaf help test in `cmd/bench/worktree_leaves_test.go` matches `"usage: " + usage.WorktreeBuild` in that output with a contains match, so the added form keeps it green. The help row in `cmd/bench/main.go` must not grow, because that file is over its line cap.

The verb refuses before any write, in the spec's order:

1. The repository declares no Bench build inputs.
2. The running executable's canonical physical path is not the primary checkout's `dist/bench`. The refusal names the running executable's path, so a worktree or version-cache selection is visible, and its route names the primary checkout. The Go code does not reproduce the wrapper's resolution order; `bin/bench.sh` stays its one source.
3. The target has no bound assignment. The lookup reads the assignments and the commitment bindings through one new exported read-only entry in `internal/intent` over the existing tolerant read mode. That entry writes no ledger byte. It recovers past an unknown commitment field; it drops an assignment record that it cannot decode, so it does not cover a declaration under `assignments`.
4. The default-branch tip holds no copy of the bound spec, or that copy declares no persisted field. Add one exported reader to `internal/commitment/repository/sources.go` that reads the declaration of a bound spec at the default-branch tip over the existing default-branch source read and the ticket 1 scan. The install uses it here, and tickets 3 and 6 consume it. Print the scan's `persisted-fields-invalid` error when it refuses.

Then build the target through the existing build owner, and keep the worktree copy. Publish a copy of the executable, its seal, and the broker manifest at the primary checkout through `freshness.Publish`, with the target tree as the seal root and the primary `bin/` as the manifest directory. Last, run the primary `bin/bench.sh doctor --fix` as a child. The child's working directory is the primary checkout. Its environment unsets `BENCH_KIT` and `BENCH_WRAPPER` and sets `BENCH_RUN_BINARY` to the primary executable. Reach the build, the publication, and the child through `joins.go` seams, because the journey-harness census forbids `t.Setenv`, `exec.Command`, and `gittest.` in these unit tests. A failed child leaves the published set in place, and the verb reports the child's refusal.

The verb prints one table with the worktree, the executable digest, and the manifest path; put its table name beside the existing table names in `list.go`. Its `next` line names the first ticket dispatch. Each refusal route literal stays in `build.go` or composes through the refusal-route registry. Raise the package test-count pin by one for each new top-level test.

Add the system test file with the `system` build tag. It builds the `main` executable and a candidate tree that adds one ledger field and one review-record field, binds the build on the fixture's default branch, and runs the install through the real wrapper. Run the suite only through `bench test --check system --run TestCandidateInstall`, which sets `BENCH_KIT` to the tree under test. Never set `BENCH_KIT` to the primary checkout. Shared fixture helpers may go in `owner_test.go`, which has 75 lines of headroom. This ticket's test function is `TestCandidateInstallPublishesTheCandidate`. Ticket 3 adds its own landing test functions to this file and does not edit this one.

Tickets 3, 5, and 6 consume the install constant, the child environment, and the tolerant read entry.

## Acceptance

- [ ] `TestBuildInstall` shows the primary executable bytes equal the target's sealed bytes, and both copies have one SHA-256 digest (CB13, CB14).
- [ ] The install refuses and writes no file in a repository with no build inputs, for an unbound target, for a default-branch spec with no declaration, and when only the target tree's spec copy declares a field (CB16, CB17, CB18, CB19).
- [ ] A second install after a new commit replaces both copies with the new digest (CB20).
- [ ] With a ledger that holds an unknown commitment field, an install run by the `main` executable finds the bound assignment and publishes, and the ledger bytes stay unchanged (CB54, CB55).
- [ ] When the publication fails after the executable write, the primary executable, seal, and manifest keep their earlier bytes (CB56).
- [ ] `TestCandidateInstallPublishesTheCandidate` shows the manifest `path` is the primary executable with that executable's digest after a child that starts beside a built worktree `dist/bench` (CB15).
- [ ] After the publication, the install calls its child seam once with argv `<primary>/bin/bench.sh doctor --fix`, the primary checkout as its working directory, no `BENCH_KIT` or `BENCH_WRAPPER`, and `BENCH_RUN_BINARY` set to the primary executable. A failing child makes the install exit non-zero, and the published set keeps the candidate bytes (CB63).
- [ ] An install whose running executable is a worktree `dist/bench` refuses before its build, names that executable's path, and writes no file (CB21).
- [ ] After the install, `bench status` from the primary checkout exits 0 with a ledger that holds the fixture's new ledger field, which the `main` executable does not know (CB23).
- [ ] The help inventory and every existing `bench worktree build <target>` remedy string stay unchanged.

## Checkpoint verification

Run the `worktree`, `system`, `cmd-bench`, `internal-conformance`, `internal-intent`, `internal-usage`, and `root-conformance` commands of CB-C2's verification plan. Probe the canonical-path refusal and the tolerant lookup with one omission mutant each, and record the red.
