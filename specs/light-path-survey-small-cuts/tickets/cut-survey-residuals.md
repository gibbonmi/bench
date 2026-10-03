# Cut the quality survey residuals

Blocked by: none
Writes: internal/reviewrecord/parse.go, internal/reviewrecord/coverage.go, internal/reviewrecord/record.go, internal/reviewrecord/plan.go, internal/reviewrecord/delegated.go, internal/preflight/evidencecmd/operations.go, internal/gate/gate_go.go, internal/gate/tag_census_test.go, internal/releasepreflight/types.go, internal/releasepreflight/decision.go, internal/releasepreflight/command.go, internal/releaseevidence/release_evidence.go, internal/releaseevidence/release_requirements.go, internal/releaseevidence/requirement_inspection.go, internal/preflight/decision.go, internal/preflight/proposal.go, internal/tickets/registry_data.go, internal/tickets/registry_data_test.go, internal/canary/inventory.go, internal/racetests/racetests.go, internal/census/census_test.go, projects/benchkit.md, ASSESSMENT.md
Covers: none

## What to build

Roadmap row FT367 owns these cuts. Each item is dead, duplicated, or stale.

Four production sites keep a hand-rolled string-slice `contains` beside the
standard `slices.Contains`. The sites are `internal/reviewrecord`,
`internal/preflight/evidencecmd`, `internal/gate`, and the
`releaseevidence.Contains` helper that `internal/releasepreflight` wraps. Each
site calls `slices.Contains` directly. The `gitguard` helper is out of scope,
because it reuses its own `indexOf`.

Four more helpers repeat the same loop under other names:
`containsStr` in `internal/preflight`, `holdsString` in `internal/tickets` and
in `internal/canary`, and `containsProfile` in `internal/releaseevidence`. Each
caller also calls `slices.Contains` directly.

The `internal/racetests` package has one importer, `internal/gate`. Its table
moves into the gate package, and the package goes. Each reference to the old
package path changes to the new owner.

`ASSESSMENT.md` finding H-A1 and its summary row say that the tag workflow
loops `npm publish`. The workflow now runs `bench release submit`. The
finding states the current workflow.

The `internal/terminal` package stays. It is a seed owner of the
ticket-grammar binding registry, so its cut needs a reviewer decision.

## Acceptance

- [ ] No production Go file defines a string-slice `contains` helper that `slices.Contains` replaces.
- [ ] The race phase runs the same seven tests from the gate package, and no Go file or guidance file names the `internal/racetests` package.
- [ ] `ASSESSMENT.md` H-A1 names `bench release submit` as the publish step of the tag workflow.
- [ ] The affected package suites and the root conformance pass stay green.
