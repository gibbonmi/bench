# Validate registry completeness and inertness

Blocked by: declare-axi-query-root-metadata.md, declare-axi-operational-root-metadata.md, declare-axi-spec-nested-metadata.md, declare-axi-worktree-nested-metadata.md, declare-axi-release-tool-nested-metadata.md
Ownership fence: `cmd/bench/command_registry.go`, `cmd/bench/command_registry_test.go`, `internal/conformance`, `projects/benchkit.md`
Integration surfaces: root and nested declarations→declare-axi-query-root-metadata.md, declare-axi-operational-root-metadata.md, declare-axi-spec-nested-metadata.md, declare-axi-worktree-nested-metadata.md, declare-axi-release-tool-nested-metadata.md; conformance check inventory→`internal/conformance/registry/registry.go` (`Checks`); profile advertisement of that inventory→`projects/benchkit.md` (the `conformance check | input source` table); compatibility oracle→`internal/axi/compatibility` exercised by RV1/inert; exact approved set→assert-approved-axi-disposition.md; guidance→document-ten-principle-axi.md
Contracts: the complete declaration set and the raw four-observation command capture cross `cmd/bench/command_registry.go`→`internal/conformance`, membership is every production root and nested member, order follows the registry, and a missing, duplicate, unclassified, or independently advertised member refuses, asserted by RV1 against the real declaration table and the real oracle rather than a fixture registry
Closure: RV1/missing, RV1/duplicate, RV1/unclassified, RV1/derivation-inputs, RV1/parallel, RV1/inert

## What to build

Registry validation becomes enforcement. The whole-registry validator in
`cmd/bench/command_registry.go` refuses a production member with no declaration,
a duplicate declaration name, or a member whose AXI disposition is unset — the
totality the five declaration tickets deliberately do not each assert. A new
conformance check registers in `internal/conformance/registry/registry.go` and
derives the declaration set independently from `cmd/bench/main.go` and
`cmd/bench/command_registry.go` by AST parse, in the same style as the existing
`subcommand-routing` derivation, so a second hand-maintained AXI membership list
anywhere in the tree is a refusal rather than a duplicate source. Registering the
check obliges its row in the `projects/benchkit.md` advertisement table, which is
why that file is on this ticket's fence.

Metadata stays byte-inert, and this ticket proves it rather than asserting it.

Tree condition at refresh time: this spec follows `axi-compatibility-oracle`, so
`internal/axi/compatibility` must exist and expose the four-observation
comparator that spec's `compare-four-observations.md` ticket lands (stdout,
stderr, exit, acceptance over the production `Command` seam), reachable as a
package-level exported entry point of package `compatibility`, with
`cmd/bench/axi_compatibility_test.go` present as its existing caller. Bind
RV1/inert to that exported symbol's real name at refresh time; if the package or
the exported `compatibility.Compare` symbol is absent (`go doc ./internal/axi/compatibility Compare`), stop and report rather than writing a local
byte-comparison substitute — a hand-rolled comparison would be exactly the
parallel oracle this spec's no-public-delta seam exists to avoid. This ticket
writes nothing under `internal/axi/compatibility`.

## Acceptance

- [ ] [RV1] (covers CR6) whole-registry validation refuses a missing, duplicate, unclassified, or independently advertised member, and declaring metadata changes no command's stdout, stderr, exit, or acceptance.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RV1/missing | delete the `guards` declaration record so a production member has none | `TestAXIRegistryRefusesIncompleteMembership` in `cmd/bench` (this ticket authors it) | run `go test ./cmd/bench -run TestAXIRegistryRefusesIncompleteMembership -timeout 120s`; expect the validator refusal `production member "guards" has no axi declaration` instead of a nil error; bound is the `-timeout 120s` binary deadline over the in-process registry, which spawns no process |
| RV1/duplicate | add a second declaration record named `maps` beside the first | `TestAXIRegistryRefusesDuplicateDeclaration` in `cmd/bench` (this ticket authors it) | run `go test ./cmd/bench -run TestAXIRegistryRefusesDuplicateDeclaration -timeout 120s`; expect the refusal `duplicate axi declaration "maps"` instead of the last record silently winning; bound is the `-timeout 120s` binary deadline over the in-process registry |
| RV1/unclassified | clear the `Disposition` field on the `bench spec build promote` record | `TestAXIRegistryRefusesUnclassifiedMember` in `cmd/bench` (this ticket authors it) | run `go test ./cmd/bench -run TestAXIRegistryRefusesUnclassifiedMember -timeout 120s`; expect the refusal `declaration "bench spec build promote": axi disposition is unset` rather than a default to non-AXI; bound is the `-timeout 120s` binary deadline over the in-process registry |
| RV1/derivation-inputs | drop the disposition field from the record the conformance AST derivation builds per member, so its comparison key is the member name alone (input mutation of the value that authorizes AXI enforcement, which every control-flow-only mutation of the check leaves green) | `TestAXIDispositionDerivationCarriesEveryDeclaredField` in `internal/conformance` (this ticket authors it) | run `go test ./internal/conformance -run TestAXIDispositionDerivationCarriesEveryDeclaredField -timeout 300s`; expect the field assertion `derived record for "maps" has empty Disposition, want "axi"` from the parse of `cmd/bench/command_registry.go`; bound is the `-timeout 300s` binary deadline over a `go/parser` read of two files, with no process spawn |
| RV1/parallel | add a second hand-maintained AXI member list to `projects/benchkit.md` that no longer derives from the declaration table | `TestAXIMembershipHasSoleProductionOwner` in `internal/conformance` (this ticket authors it) | run `go test ./internal/conformance -run TestAXIMembershipHasSoleProductionOwner -timeout 300s`; expect the refusal `parallel axi membership advertisement in projects/benchkit.md; the declaration table is the single source` naming the offending line; bound is the `-timeout 300s` binary deadline over file reads with no process spawn |
| RV1/inert | make the `bench maps` renderer print the declared renderer-family metadata as a trailing stdout line | `TestAXIRegistryMetadataIsByteInert` in `cmd/bench` (this ticket authors it), comparing through the `internal/axi/compatibility` four-observation comparator | run `go test ./cmd/bench -run TestAXIRegistryMetadataIsByteInert -timeout 300s`; expect the comparator's stdout delta `case maps/default: stdout differs (+1 line "renderer=flat-table")` with stderr, exit, and acceptance still equal; the compared run is bounded by an explicit `context.WithTimeout(t.Context(), 60*time.Second)` around each capture so a hung child fails the case instead of the binary deadline, and the `-timeout 300s` binary deadline is the outer bound |
