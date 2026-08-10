# Migrate the spec-build empty state

Blocked by: migrate-specbuild-aggregates.md
Ownership fence: `internal/specbuild`, `cmd/bench/specbuild.go`, `cmd/bench/specbuild_test.go`
Integration surfaces: typed empty carrier→`internal/axi/empty.go` exercised by SBE1; no-run status producer→`internal/specbuild/state.go` exercised by SBE1; compact and full status renderers→`cmd/bench/specbuild.go` exercised by SBE1; aggregate sibling blocker→migrate-specbuild-aggregates.md; registry empty declaration→declare-empty-dispositions.md; legacy carrier contraction→contract-aggregate-empty-routes.md
Contracts: the spec-build empty classification crosses `internal/specbuild/state.go`→`internal/axi/empty.go` and on to `cmd/bench/specbuild.go`; its class is the one-row `state=empty` projection, distinct from the zero-row table class used by the `assignments[0]`/`review[0]` detail blocks in the same output; order is the compact `spec_build` row first, then the detail tables; absence of a run is a success (exit 0), never a refusal, asserted by SBE1 against the real `Service.Status` producer
Closure: SBE1/one-row-empty, SBE1/empty-next, SBE1/detail-assignments-empty, SBE1/detail-review-empty, SBE1/exit-success, SBE1/route

## What to build

`bench spec build status <slug>` with no retained run keeps its exact empty class: one
`spec_build[1]{slug,state,subject,next}` row whose state is the literal `empty`, and under
`--full` the two zero-row detail tables `assignments[0]` and `review[0]` beside it. The
migration routes that classification through the shared typed empty carrier so the class is
declared and observable, without collapsing the one-row empty into a zero-row table or the
zero-row detail tables into a one-row empty.

AE8 is an already-covered row for this family: `TestSpecBuildStatusRendersDefinitiveEmptyProjection`
(`cmd/bench`) and `TestStatusHasDefinitiveEmptyAndActiveProjections` (`internal/specbuild`)
stay exactly as they are and remain the named existing controls. This ticket adds the subject
mutation the row lacks — a new `TestSpecBuildEmptyClassReachesTheTypedRoute` in `cmd/bench`
asserting that the no-run status was classified through the shared carrier rather than by a
local literal.

This ticket is deliberately separate from `migrate-specbuild-aggregates.md`: it lands green
on its own and shares no gate red with the count aggregates. It is blocked by that ticket
only because the two write the same files.

Tree condition that must hold when this ticket is refreshed: `internal/axi/empty.go` exists
and declares the exported empty classification type `EmptyClass` with distinct constants for
the zero-row table class and the one-row lifecycle-empty class. If that path or the symbol
is absent, stop and report rather than build — the prerequisite `axi-carriers-and-registry`
build has not landed.

## Acceptance

- [ ] [SBE1] (covers AE8) spec-build status with no retained run renders its one-row `state=empty` projection, its zero-row `--full` detail tables, and exit 0 through the shared typed empty route.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SBE1/one-row-empty | in `Service.Status`, return a zero-row projection for the not-found case instead of `Status{Slug: slug, State: "empty", ...}` | `TestSpecBuildStatusRendersDefinitiveEmptyProjection` (`cmd/bench`) | run `go test ./cmd/bench -run TestSpecBuildStatusRendersDefinitiveEmptyProjection -count=1 -timeout 120s`; the `strings.Contains(out, "spec_build[1]{slug,state,subject,next}:")` assertion fails against a rendered `spec_build[0]{...}` header; the fixture is a `t.TempDir()` `git init` repository holding one staged `specs/demo/spec.md`, and no lifecycle git call is made on the not-found path |
| SBE1/empty-next | in `Service.Status`, return the empty projection with an empty `Next` cell | `TestSpecBuildStatusRendersDefinitiveEmptyProjection` (`cmd/bench`) | run `go test ./cmd/bench -run TestSpecBuildStatusRendersDefinitiveEmptyProjection -count=1 -timeout 120s`; extend the existing `demo,empty` containment check to the full row — the assertion that the fourth cell is `bench spec build start demo` fails with an empty cell, so the empty state stops naming its recovery command; same `t.TempDir()` fixture, no git call on this path |
| SBE1/detail-assignments-empty | in `renderFullStatus`, skip the `assignments` table entirely when `full.Assignments` is empty | `TestSpecBuildStatusRendersDefinitiveEmptyProjection` (`cmd/bench`) | run `go test ./cmd/bench -run TestSpecBuildStatusRendersDefinitiveEmptyProjection -count=1 -timeout 120s`; the `strings.Contains(full, "assignments[0]")` assertion fails because the block is absent; same `t.TempDir()` fixture |
| SBE1/detail-review-empty | in `renderFullStatus`, return early on `full.Review == nil` without emitting the zero-row `review` table | `TestSpecBuildStatusRendersDefinitiveEmptyProjection` (`cmd/bench`) | run `go test ./cmd/bench -run TestSpecBuildStatusRendersDefinitiveEmptyProjection -count=1 -timeout 120s`; the `strings.Contains(full, "review[0]")` assertion fails because the block is absent; same `t.TempDir()` fixture |
| SBE1/exit-success | in `specBuildCommand`, route the not-found status through `buildError` so the empty state exits 1 | `TestSpecBuildStatusRendersDefinitiveEmptyProjection` (`cmd/bench`) | run `go test ./cmd/bench -run TestSpecBuildStatusRendersDefinitiveEmptyProjection -count=1 -timeout 120s`; the `code != 0` check fails with exit 1 and a structured error line in place of the projection; same `t.TempDir()` fixture |
| SBE1/route | classify the no-run case with the local `"empty"` string literal and never construct the shared typed empty value | `TestSpecBuildEmptyClassReachesTheTypedRoute` (`cmd/bench`, new) | run `go test ./cmd/bench -run TestSpecBuildEmptyClassReachesTheTypedRoute -count=1 -timeout 120s`; the assertion that the no-run status carried the one-row lifecycle-empty `axi.EmptyClass` constant fails with no classification observed, even though the rendered bytes are unchanged; same `t.TempDir()` fixture, no git call on the not-found path |
