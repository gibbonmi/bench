# Close the ownership closures

Blocked by: read-preflight-rows-from-the-parser.md
Writes: internal/canary/inventory.go, internal/canary/inventory_test.go, internal/tickets/registry_data.go (new), internal/tickets/registry_data_test.go (new), internal/preflight/gather.go, internal/preflight/decision.go, internal/preflight/decision_test.go, internal/preflight/command_build_test.go
Covers: TG14, TG15, TG16, TG19, TG26, TG39

## What to build

Three ownership rows complete the verdict table in one diff: `fixture-closure`,
`registry-closure`, and `kit-pin`, in that order. Each derives from a producer,
and none reads a copied list. `Decide` stays pure: the gatherer collects every
fact the three rows grade.

`internal/canary` exports the enumeration that maps a repository path to the
fixtures that pin it. A fixture pins a path when its `BASE` list, its `files/`
overlay, or its `MUTATE.json` names that path. The `BASE` walk follows `@`
includes. The `fixture-closure` row reds a ticket that writes a pinned path
and omits the owning fixture directory.

Declared data binds each bound package to the files a ticket must co-name. The
table lives in `internal/tickets` and follows the shape of
`internal/anchors/registry_data.go`. The command rows include the help
projection and the envelope cases. Seed rows also bind the dispatcher, the
renderer, and the terminal-lifecycle owner. The `registry-closure` row reds a
ticket that writes a bound package and omits a bound file.

The gatherer parses the build constraints of each written Go test file. When
the system tag is present, the `kit-pin` row requires the literal `BENCH_KIT`
in the ticket body.

With the span complete, the six grammar rows render in the fixed order. In
build mode with no tickets directory, all six render not-applicable.

## Acceptance

- [ ] TG14 — a ticket that writes a pinned path reds unless it names the fixture.
- [ ] TG15 — a planted synthetic fixture appears in the live pin enumeration.
- [ ] TG16 — a ticket that writes a bound package reds unless it names every bound file.
- [ ] TG19 — a written system-tagged test file reds unless the body states `BENCH_KIT`.
- [ ] TG26 — `bench preflight build` renders the six rows, each red with its own detail.
- [ ] TG39 — with no tickets directory in build mode, the six rows render not-applicable.
- [ ] The `kit-pin` row stays green for a test file with no system tag.

## Delegate charge

You work in the Bench repo on the `ticket-grammar` spec. Read
`specs/ticket-grammar/spec.md` first. Then read `internal/canary/inventory.go`
and `internal/canary/mutation.go` in full; `basePaths` lives in `mutation.go`.
Read `internal/preflight/gather.go` and `internal/preflight/decision.go` next.
Read `internal/anchors/registry_data.go` for the declared-data shape, and
`internal/conformance/axi_query_registry_test.go` for the approved query list.

Export one enumeration from `internal/canary` that answers which fixtures pin
a given repository path. Reuse the existing `basePaths` helper rather than a
second walk. Follow `@` includes.

Create `internal/tickets/registry_data.go`. Declare one row per bound package:
a package prefix and the files a ticket must co-name. Include the help
projection and the envelope cases for each command row. Add seed rows for the
dispatcher (`cmd/bench`), the renderer (`internal/toon`), and the
terminal-lifecycle owner (`internal/terminal`). Keep the rows sorted by
prefix.

Gather every graded fact in `gather.go`: the pin map, the binding matches, and
the build-constraint tags. Keep `Decide` free of I/O. Add the
`fixture-closure`, `registry-closure`, and `kit-pin` rows to `Decide` in that
order, after `writes-resolve` and before `rows-owned`. Name the offending
path and the omitted files in each detail.

Add `TestFixturePinsEnumeratesLiveInventory` in `internal/canary`; plant a
synthetic fixture in a temporary root and assert the new pin. Add
`TestRegistryBindingRowsAreSorted` in `internal/tickets`. Add
`TestFixtureClosureNamesUnnamedFixture`,
`TestRegistryClosureNamesOmittedRegistry`, `TestKitPinRequiresBenchKit`, and
`TestSixRowsNotApplicableWithoutTickets` in `internal/preflight`. Add a
command test that asserts the six-row render with per-row details. Assert the
exact detail text in each test.

Run only `bench worktree exec ft174-ticket-grammar -- go test ./internal/canary/ ./internal/tickets/ ./internal/preflight/`.
Do not commit. Do not edit the spec.
