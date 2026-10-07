# Warn and allow an unbound delegation model

Blocked by: none
Writes: internal/lines/lines.go, internal/lines/lines_agentline_test.go, internal/conformance/line_routing_exec_test.go, .bench/hooks/check-agent-line.sh, docs/adr/0002-accepted-enforcement-postures.md
Covers: none

## What to build

The agent-line guard denies an Agent delegation whose model matches no bound tier. A reviewer can therefore not run a delegate on a model outside the binding, for example to compare a new model against a bound tier. The reviewer decided on 2026-10-07 that this branch warns and allows the delegation.

Make the unbound-model branch of the verdict return exit 0 with a one-line warning. The warning names the model and the asking harness's own bound column, so the declared line stays visible. An allowed delegation records its intent, as every allowed delegation does. The guard keeps its two denials: an omitted or empty model in a routed repo, and a model declared on a fork. Update the guard manifest and the enforcement-posture ADR to state the new posture.

## Acceptance

- [ ] In a routed repo, an Agent delegation on a model that no tier binds runs.
- [ ] The guard's warning for that delegation names the model and the asking harness's bound column.
- [ ] The warning names no token from another harness's column.
- [ ] An Agent delegation with an omitted or empty model in a routed repo is still denied, and the guard records no intent for it.
- [ ] A fork delegation that declares a model is still denied.
- [ ] `bench guards` describes the guard as a deny of an omitted model and a warning for an unbound model.
