# 3. Name the graded tree in each tree-scoped response

Blocked by: 1-declare-command-scope.md
Writes: internal/treetarget/ (new), cmd/bench/tree_scope.go (new), cmd/bench/tree_scope_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, cmd/bench/main_test.go, cmd/bench/selected_queries_test.go, cmd/bench/isolated_command_fixtures_test.go, cmd/bench/spill_support_test.go, cmd/bench/census_output_test.go, cmd/bench/preflight_version_test.go, cmd/bench/commit_chain_test.go, cmd/bench/census_output.go, internal/responsebound/owner.go, internal/responsebound/owner_test.go, internal/systemtest/
Covers: TT10, TT11, TT12, TT13, TT14, TT15, TT16, TT17, TT18, TT19, TT20, TT21, TT22, TT23, TT24, TT25, TT56

## What to build

Chunk: TT-C3.

This ticket waits for pending reviewer decision 1 in the spec, the row and the response bound.

Create the package `internal/treetarget` with an `Identify` function and a row renderer. `Identify` answers the target, head, and dirty values for a root, as the spec's identity-row section states. The renderer prints the block `tree[1]{target,head,dirty}:` through `toon.Table`.

The dispatcher computes the row before the verb runs. The bounded response owner writes it at finish as the first block, and only when the exit is not 2. A call with an exempt bound disposition prints the row on stderr after the verb returns, under the same exit rule. A help form, a repository-scoped verb, a grammar refusal at exit 2, and a call outside a repository print no row.

The new first block reds the fixtures that the spec's posture-change section lists. Add one test helper in `cmd/bench/tree_scope_test.go` that removes a leading row block. Call it once inside each of the four dispatch helpers that the section names. Then edit each of the seven direct sites in place. `main_test.go` and `command_registry_test.go` are over budget, so keep each at or under its current line count.

The system tests run with `BENCH_KIT` set to the kit root.

## Acceptance

- [ ] `bench roadmap` in a primary checkout starts with `tree[1]{target,head,dirty}:` and the row `primary,<40-character HEAD>,false`.
- [ ] `Identify` answers `alpha` for a worktree that an active assignment labeled `alpha` owns, and the row holds no byte of its path.
- [ ] `Identify` answers `unassigned` for an unowned linked worktree and for a released assignment.
- [ ] An unborn HEAD renders `none`, and an untracked file or a modified tracked file renders `dirty` `true`.
- [ ] A corrupt index renders `dirty` `unknown`.
- [ ] `bench dashboard --stdout` prints no `tree[` line on stdout, and its stderr carries the row block.
- [ ] `bench gate --help`, `bench version`, and `bench coverage x` outside a repository print no `tree[` line.
- [ ] `bench gate --brief` exits 2 with an empty stdout.
- [ ] A planted tree-scoped verb that spills keeps the row header as its first inline line.
