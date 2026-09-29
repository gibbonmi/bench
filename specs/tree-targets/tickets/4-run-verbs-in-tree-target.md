# 4. Run a tree-scoped verb in its named tree target

Blocked by: 3-name-the-graded-tree.md
Writes: internal/treetarget/ (new), internal/worktree/tree_target.go (new), internal/worktree/tree_target_test.go (new), internal/worktree/exec.go, cmd/bench/tree_scope.go (new), cmd/bench/tree_scope_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/systemtest/, internal/canonicalpath/canonicalpath.go, internal/canonicalpath/canonicalpath_test.go
Covers: TT8, TT9, TT26, TT27, TT28, TT29, TT30, TT31, TT32, TT33, TT34, TT35, TT36, TT37, TT38, TT39, TT40, TT41, TT49, TT53, TT54, TT55, TT57, TT58, TT59, TT62

## What to build

Chunk: TT-C4.

The label lookup follows step 5 of the spec's tree-target section: exactly one active match wins, by the reviewer's decision of 2026-09-29. `canonicalpath.Resolve` makes a relative path absolute before it resolves symlinks, so a symlinked working directory gives the physical spelling (TT62).

`TestTreeTargetOnlyAsFirstArgument` also drives a repository verb with a late `--in`, such as `bench version x --in primary`, which reaches the verb's own grammar. The TT-C1 review found that no test pins this case.

`--in <label|primary>` as the first argument of a tree-scoped verb runs that verb as one child in the tree target. Resolve the value in the seven-step order of the spec's tree-target section. Add an exported exact-label lookup in `internal/worktree/tree_target.go`. It keeps the state, missing-tree, and creation-bundle checks of `resolveAssignmentIn`, and it accepts no id, prefix, or path. After the label lookup fails, it calls `targetPath` and answers a typed path-shape error.

Start the child with a new exported runner beside `runWorktreeChild`. It uses `execEnv` and prints no `worktree:` line. The child directory is `canonicalpath.Resolve` of the target root. The child argv is the executable, the verb, and the remaining arguments. The parent prints no row, applies no bound, and returns the child's exit code.

For a target without build inputs, the executable is the wrapper that `BENCH_WRAPPER` names, or else the running executable. `internal/treetarget` takes that executable as an argument, so its tests pass a marker script. Ticket 5 adds the kit worktree build.

Root help renders each tree-scoped row as `bench <name> [--in <label|primary>]<suffix>`, from the scope field. Change the whole expectation of `TestHelpInventoryIsComplete` in place. Add system rows with `BENCH_KIT` set to the kit root.

## Acceptance

- [ ] `bench status --in alpha` from the primary checkout prints the row `alpha,<worktree HEAD>,false` in the system suite.
- [ ] With its directory inside `alpha`, `bench status --in primary` prints the row `primary,<primary HEAD>,false`.
- [ ] `bench coverage --in alpha specs/only-in-alpha/spec.md` reads the spec in `alpha`.
- [ ] A missing value, an empty value, `--help`, an absolute path, `./alpha`, and `~` each exit 2 with the spec's usage lines, and no marker exists.
- [ ] An unknown, a released, or a colliding label exits 1 through the worktree target refusal, and no marker exists.
- [ ] A value with U+0001 refuses with no raw U+0001 on stderr.
- [ ] The label `team/alpha` of an active assignment starts the child.
- [ ] `bench gate --fresh --in primary` exits 2 with the gate usage.
- [ ] A child that prints `x` and exits 3 gives stdout `x`, exit 3, and no `worktree:` line.
- [ ] The wrapper in `BENCH_WRAPPER` runs as the child, and without it the running executable runs.
- [ ] A target without build inputs runs the wrapper even when `dist/bench` exists.
- [ ] The child `PWD` equals `canonicalpath.Resolve` of the recorded path when the home is a symlink.
- [ ] `bench worktree exec alpha -- <executable> status` exits 0 with the row target `alpha`, and so does a run with its directory in `alpha` and no `--in`.
- [ ] Root help shows the insertion on tree-scoped rows only.
