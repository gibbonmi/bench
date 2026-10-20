# Install a declared build's candidate as the repository executable and broker

Blocked by: 1-declare-persisted-fields.md
Writes: internal/worktree/build.go, internal/worktree/build_install_test.go (new), internal/worktree/build_test.go, internal/worktree/joins.go, internal/worktree/verb_runner_test.go, internal/worktree/verb_runner_check_test.go, internal/worktree/parallel_census_test.go, internal/usage/worktree.go, internal/intent/transaction.go, internal/intent/transaction_test.go, cmd/bench/help_inventory_test.go, internal/conformance/refusal_route_bypass_test.go
Covers: CB13, CB14, CB16, CB17, CB18, CB19, CB20, CB22, CB54, CB55, CB56

## What to build

Deliver the install half of CB-C2: `bench worktree build --install <target>` installs a declared build's candidate as the primary checkout's executable and broker.

Add the `--install` flag to the build grammar in `build.go`. Keep the existing `usage.WorktreeBuild` text unchanged, because `WorktreeBuildFor` derives every rebuild remedy from it. Put the install form in a separate constant in `internal/usage/worktree.go`. The help row in `cmd/bench/main.go` must not grow, because that file is over its line cap.

The verb refuses before any write, in the spec's order:

1. The repository declares no Bench build inputs.
2. The running executable's canonical physical path is not the primary checkout's `dist/bench`.
3. The wrapper's order selects a package or version-cache executable before the primary `dist/bench`. The refusal names that path.
4. The target has no bound assignment. The lookup reads the assignments and the commitment bindings through one new exported read-only entry in `internal/intent` over the existing tolerant read mode. That entry writes no ledger byte and takes no write.
5. The default-branch tip holds no copy of the bound spec, or that copy declares no persisted field. Use the CB-C1 header scan on the default-branch bytes.

Then build the target through the existing build owner, and keep the worktree copy. Publish a copy of the executable, its seal, and the broker manifest at the primary checkout through `freshness.Publish`, with the target tree as the seal root and the primary `bin/` as the manifest directory. Last, run the primary `bin/bench.sh doctor --fix` as a child. The child's working directory is the primary checkout. Its environment unsets `BENCH_KIT` and `BENCH_WRAPPER` and sets `BENCH_RUN_BINARY` to the primary executable. Reach the build, the publication, and the child through `joins.go` seams so the unit test drives them without a real build. A failed child leaves the published set in place, and the verb reports the child's refusal.

Each refusal route literal stays in `build.go`, which the route-bypass test allowlists, or it composes through the refusal-route registry. Raise the package test-count pin by one for each new top-level test.

CB-C3 and CB-C4 consume the install constant, the child environment, and the tolerant read entry.

## Acceptance

- [ ] `TestBuildInstall` shows the primary executable bytes equal the target's sealed bytes, and both copies have one SHA-256 digest (CB13, CB14).
- [ ] The install refuses and writes no file in a repository with no build inputs, for an unbound target, for a default-branch spec with no declaration, and when only the target tree's spec copy declares a field (CB16, CB17, CB18, CB19).
- [ ] A second install after a new commit replaces both copies with the new digest (CB20).
- [ ] With `$BENCH_HOME/cache/bin/<version>` holding an executable for this platform, the install refuses and names that path (CB22).
- [ ] With a ledger that holds an unknown commitment field, an install run by the `main` executable finds the bound assignment and publishes, and the ledger bytes stay unchanged (CB54, CB55).
- [ ] When the publication fails after the executable write, the primary executable, seal, and manifest keep their earlier bytes (CB56).
- [ ] The help inventory and every existing `bench worktree build <target>` remedy string stay unchanged.

## Checkpoint verification

Run the `worktree` command of CB-C2's verification plan. Probe the canonical-path refusal and the tolerant lookup with one omission mutant each, and record the red.
