# Refuse a direct verb call outside the verb runner

Blocked by: 9-migrate-landing-composition.md, 10-migrate-landing-effects.md
Writes: internal/worktree/verb_call_census_test.go (new), internal/worktree/parallel_census_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: VR40, VR41, VR42, VR43, VR44, VR45, VR46, VR47, VR48, VR49, VR50, VR51, VR52, VR53, VR54, VR55, VR56, VR57, VR58

## What to build

Add `verb_call_census_test.go`. It reuses the parse helpers of `parallel_census_test.go` and adds the verb call census. A verb entry is an exported function of a non-test file with a parameter named `args` of type `[]string`. A joins form is a function that a verb entry's body calls with `defaultJoins()` as its first argument. A runner-private function is a top-level function of a runner file whose last result is an `error`.

The runner files are exactly `verb_runner_test.go` and `verb_runner_check_test.go`. The census reports each identifier that names a verb entry, a joins form, or a runner-private function in any other test file. It walks function bodies, closures, and package-level declarations. Each report names the file, the line, the enclosing declaration, and the reported name, and the reports come back sorted. The synthetic tests plant file sets with `plantTestFiles`. One live-tree test runs the census over the package and expects no report.

Raise `worktreeTestCount` by the tests this ticket adds. This ticket is the last one that touches the package, so it carries the package-wide end state. That state is the tuple scan, the VR41 command, the count pin, the serial ceiling, and the final differential run.

## Acceptance

- [ ] The census reports an entry call, a joins form call, an entry value, a subtest call, and a package-level reference.
- [ ] The census reports a runner-private function call outside the runner files.
- [ ] Each census report names its file and line.
- [ ] The census reports nothing in either runner file and nothing for an exported function without an `args []string` parameter.
- [ ] The census derives a synthetic entry from its signature and a synthetic joins form from the entry's body.
- [ ] The live-tree census reports no line.
- [ ] The tuple scan and the VR41 command print no line.
- [ ] `TestPackageTestCountPin` and `TestSerialSetStaysBelowTheCeiling` pass.
- [ ] The fresh test run's PASS and SKIP name sets match the spec base's sets plus the added tests.
