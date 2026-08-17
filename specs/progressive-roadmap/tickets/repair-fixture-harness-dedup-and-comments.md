# Repair: de-duplicate the split-board test fixture harness, drop PR-talk comments

Blocked by: none
Writes: internal/roadmap/tree_helpers_test.go, internal/roadmap/tree_test.go, internal/conformance/docs_workflow_checks_test.go, cmd/bench/command_registry_test.go, internal/status/status_test.go, internal/conformance/recurrence_maintenance_contract_test.go, internal/roadmap/context_test.go, internal/roadmap/tree.go

## What to build

Review findings Standards F1, F5, F3 (`reviews/progressive-roadmap.md`).

One split-board test fixture builder replaces the five near-duplicate copies:
`internal/roadmap/tree_helpers_test.go:50` (`writeSplitBoard`),
`internal/conformance/docs_workflow_checks_test.go:579` (`writeSplitRoadmap`),
`cmd/bench/command_registry_test.go:383`, and the two inline copies in
`internal/status/status_test.go`. `internal/roadmap/tree_test.go:218`'s local
`writeBoard` closure, which shadows the package-level `writeBoard` at
`tree_helpers_test.go:69`, is removed in favor of the package-level one.
`board()`'s `[2]string` row data clump (`tree_helpers_test.go:78`, space-index
panic risk on a heading with no space) becomes a small named struct.

Drop the PR-talk register comments at: `internal/status/status_test.go:177`,
`internal/conformance/docs_workflow_checks_test.go:42`,
`internal/conformance/recurrence_maintenance_contract_test.go:157`,
`internal/roadmap/tree.go:111`, `internal/roadmap/context_test.go:428` — per
`craft-comments`, no narration of what changed or where something used to
live; state what the code does now.

## Acceptance

- [ ] Exactly one fixture-harness helper builds a split board (index + row files) for tests, called from every prior duplicate site.
- [ ] `internal/roadmap/tree_test.go` has no local `writeBoard` shadowing the package-level one.
- [ ] The five named PR-talk comments read as durable description, not change narration.
- [ ] `go test ./...` and `bench gate` stay green.
