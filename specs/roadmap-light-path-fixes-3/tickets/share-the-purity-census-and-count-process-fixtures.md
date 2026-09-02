# Share the purity census and count process fixtures

Blocked by: derive-the-canonical-path-in-one-leaf-package.md
Writes: internal/puritycensus (new), internal/worktree/lifecyclepolicy/purity_census_test.go, internal/worktree/reclaimpolicy/purity_census_test.go, internal/worktree/landingpolicy/purity_census_test.go, internal/canonicalpath/purity_census_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LQ12, LQ13, LQ14, LQ15, LQ16

## What to build

Verify the premise first: the three `TestPurePackageSourceCensus` functions are
byte-identical except for the `internal/bounds` alternative. Then add
`internal/puritycensus` with one entry that takes the testing handle, the
directory, and a policy. The policy carries the forbidden-import patterns, the
ambient-effect set, the `t.Parallel` ban, and the self-exempt file name.
`exec.Command` and `exec.CommandContext` stay in the ambient set, so a
process-backed git fixture counts.

Each policy package keeps a one-line wrapper that passes `"."`, and the
`internal/canonicalpath` census from the blocker ticket becomes the fourth
wrapper. Each wrapper also asserts the scanned set holds the package's own
source file. The shared forbidden-import list carries `internal/bounds`, which
widens `landingpolicy`'s list; that package imports no bounds today. Add a
helper-level test over in-memory sources, mirroring
`TestHarnessCensusRefusesChildStartOutsideHarness` in
internal/worktree/effect_census_test.go. The helper package scans its own
directory under the same policy.

## Acceptance

- [ ] The three policy censuses and the `canonicalpath` census pass through the helper with the shared list.
- [ ] The helper test yields three diagnostics with file and line for a forbidden import, an `os.Getenv` call, and a `t.Parallel(` call.
- [ ] The helper test reds a source that calls `exec.Command`.
- [ ] Each wrapper reds when the scanned set lacks the package's own source file.
- [ ] Self-probe: make the helper resolve the directory from `os.Getwd`, and report which wrapper assertion reds.
