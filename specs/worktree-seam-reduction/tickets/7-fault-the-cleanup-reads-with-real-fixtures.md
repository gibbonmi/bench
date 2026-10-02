# 7. Fault the cleanup reads with real fixtures

Blocked by: 6-land-the-stubbed-landing-tests-for-real.md
Writes: internal/worktree/joins.go, internal/worktree/effects.go, internal/worktree/clean.go, internal/worktree/clean_landed.go, internal/worktree/live_binary.go, internal/worktree/lifecycle.go, internal/worktree/worktree_test.go, internal/worktree/live_binary_test.go, internal/worktree/clean_landed_hostile_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS46, WS47, WS48, WS49, WS50

## What to build

Chunk: SR-C5.

Remove `liveBinaryWarnings` from the joins value. The residue guard and the section
warning write to the ambient warnings writer, which the verb entry sets to its stderr.
Convert the retirement warning test and the four live-binary tests to that writer.

Remove `planLandedExplicit` and the dead `planLandedExplicitWithOptions` declaration. The
special-path test grades the real planner's shape reason.

The ignored-stat test depends on a probe. Plant an ignored file in a directory without
search permission, then run `bench probe` on `clean.go` with the stat error branch
omitted. If the test turns red, remove `ignoredLstat`. If the probe stays green, keep the
field and its test unchanged, and record one `bench learning` entry. Under the root user,
the test calls `capability.Capability` with `capability.Privilege`.

Each converted test keeps its name. `worktree_test.go` and `lifecycle.go` are over their
line budgets, so the change does not grow them.

## Acceptance

- [ ] A retirement that cannot remove the handoff section prints the file and the line on the verb's stderr.
- [ ] The residue guard warns on the captured writer for the live binary and stays silent for a foreign binary.
- [ ] `clean --landed` retains each special path with the real planner's shape reason.
- [ ] An ignored file in a directory without search permission makes the plan retain as `uncertain`, or the learning entry records the failed probe.
- [ ] The joins value declares no `liveBinaryWarnings` and no `planLandedExplicit` field.
