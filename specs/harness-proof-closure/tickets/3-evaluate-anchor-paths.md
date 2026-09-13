# Evaluate anchor paths

Blocked by: none
Writes: internal/anchors/registry.go, internal/anchors/anchor_harness_diagnostics_test.go, cmd/bench/anchors_command.go, cmd/bench/anchor_help_test.go, internal/gate/lane_select.go, internal/gate/lane_select_test.go, projects/benchkit.md
Covers: HP6, HP7, HP8, HP9, HP10

## What to build

Add one anchors-package path evaluation that reads a registered file once and
derives its ordered locations and registry diagnostics from the same bytes. Make
the existing group checks and `bench anchors <path>` use that evaluation logic.
The command keeps its anchor rows. A satisfied registered path returns success.
An unsatisfied registered path returns exit 1 with the exact registry diagnostics.

Derive an anchor-registry lane class from the registry's file fields. A modified
or deleted registered path selects `docs-currency-workflow` beside the existing
path classes. A path with no registry entries keeps the empty successful anchor
query and does not gain the docs check from this class.

## Acceptance

- [ ] A satisfied registered path prints its ordered anchor locations and exits 0.
- [ ] A missing required anchor or present forbidden anchor prints its registry diagnostic and exits 1.
- [ ] An unregistered regular path prints an empty anchor table and exits 0.
- [ ] A modified registered path selects `docs-currency-workflow` in the lane.
- [ ] Deletion of a registered path still selects `docs-currency-workflow`.
