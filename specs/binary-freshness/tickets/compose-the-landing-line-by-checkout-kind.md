# Compose the landing line by checkout kind

Blocked by: add-seal-and-broker-rows-to-doctor.md
Writes: internal/worktree/land.go, internal/worktree/land_freshness_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: BF19, BF20

## What to build

Verify the premise first: the broker-source notice in internal/worktree/land.go
names `bench repair` unconditionally. Then compose the line from two existing
sources. When `adopt.KitSourceCheckout` is true for the destination, the line names
the `freshness.RebuildAction` sentence and `bench doctor --fix`. Otherwise it
names `bench repair`. The doctor ticket exports that predicate first, so this
ticket follows it. If the import from `internal/worktree` into
`internal/adopt` would cycle, report the seam and stop.

## Acceptance

- [ ] `TestLandCommandReportsInstallStepForABrokerChangingDiff` receives the rebuild sentence and `bench doctor --fix` in the kit source checkout.
- [ ] A new test with the predicate false receives `bench repair` and not the rebuild sentence.
- [ ] Self-probe: inline the rebuild text, and report which sweep reds it or that none does.
