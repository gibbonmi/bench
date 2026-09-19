# 4. Run the build-script tests on a private kit copy

Blocked by: none
Writes: internal/gittest/gittest.go, internal/gittest/gittest_test.go (new), cmd/bench/build_subject_mode_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: TD47, TD48

## What to build

Chunk: TD-C2.

Add a kit copy helper to the shared git test scaffold. It copies each file that git lists as tracked, or as untracked and not ignored, from the kit root into a fresh temporary directory. It then initializes a git repository there with one commit. The copy holds no ignored file.

`TestReleasePreflightBuildDoesNotRebindThePromotionBroker` and `TestGoBuildSubjectModePublishesTheStampedVersion` each run the real build script inside a kit copy. Each test still grades the operands that the script hands the builder, and each still reads the broker manifest from the copy's wrapper directory. The live checkout then gets no `dist/` path and no rewrite of its broker manifest.

The five registry files join `Writes:` through the binding registry closure only. The ticket expects no edit to them.

## Acceptance

- [ ] A kit copy holds every tracked file and every untracked file that is not ignored, its own git directory, and no `dist/` path.
- [ ] Both build-script tests pass on the copy.
- [ ] After both tests run, the live checkout holds no new `dist/` path and an unchanged `bin/bench-broker.manifest`. The ticket report records the probe run that shows it.
