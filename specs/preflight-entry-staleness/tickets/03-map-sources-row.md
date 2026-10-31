# Grade every decision-map source against the published tree at build entry

Blocked by: 01-writes-landed-row.md, 02-closure-prediction.md
Writes: internal/maps/validation.go, internal/maps/tree_validation.go, internal/maps/maps.go, internal/maps/freshness.go, internal/maps/maps_graph_test.go, internal/maps/maps_parse_test.go, internal/preflight/decision.go, internal/preflight/gather.go, internal/preflight/gather_inputs.go, internal/preflight/map_sources.go (new), internal/preflight/map_sources_test.go (new), internal/preflight/source_tip_test.go, internal/preflight/charge_cells_test.go, internal/preflight/preflighttest, decisions/architecture-planning.md
Covers: ES01, ES02, ES03, ES04, ES05, ES06, ES07, ES08, ES09, ES10, ES11, ES12, ES15, ES16, ES17, ES18, ES19, ES20, ES21, ES22, ES23, ES41, ES42, ES46, ES48, ES49, ES52, ES53

## What to build

A build preflight reds `map-sources` before the first chunk when a decision map cites a path that the published tree does not hold, the delivery closure of the spec included.
Read these spec sections before you start: "One Sources grammar with two resolvers", "The candidate maps at the tip", "The `map-sources` fact and its red", "This build's own map record", the acceptance coverage map, and the edge inventory.

Make these changes:

- In `internal/maps`, keep one Sources record grammar and one candidate rule. Export one function that takes the map path, the map bytes, and a resolver, and returns the Sources diagnostics with the map path prefix. Export the published-tree resolver and one candidate predicate. Keep `sourceDiagnostics` as a wrapper with its current signature.
- The lexical path rule runs before a resolver and gives it the cleaned slash path. The working-tree resolver keeps each current reason of `validateSourcePath`, the `escapes the repository` reason included. `maps_graph_test.go:79-87` calls `validateSourcePath` directly; keep that function or update those calls.
- The `tests/canary/decision-map-integrity/source-*` units and `internal/conformance/decision_map_integrity_test.go` pin the exact working-tree messages. Keep those messages byte for byte, so those files need no edit.
- `internal/maps` holds 12 Go files at its directory limit. Put the exports in existing files. `maps_parse_test.go` holds 450 lines over its 400-line budget, so add no line to it; if the candidate predicate needs a new test, put it in `maps_graph_test.go`.
- Make the file-system walk call the new candidate predicate. `scanActiveMap` (`internal/maps/maps.go:81`) and `discoverDecisionMapCandidates` keep their answers. If `sourceLocator` changes its name, update `sourcePathLocators` in `freshness.go`.
- In `internal/preflight`, add the build-mode fact `MapSourceDiagnostics`. Fill it from the source tip: list the tip tree under `decisions` and `specs`, keep each path that the candidate predicate admits, refuse a candidate entry whose mode is not `100644` or `100755`, read each candidate with `git.ReadControlBlob`, predict the closure of `SpecPath` with the store method of ticket 02, and grade each map with the published-tree resolver in candidate path order. Put this code in the new file `map_sources.go`, because `gather.go` has little headroom.
- Add the row `map-sources` directly after `binary-seal`, in build mode only, outside `ticketRow`, with no `exempt` flag. Update the `Decide` doc comment. Render the detail exactly as the spec states.
- Update the row-list and row-count readers: `bareRowCount` in build mode becomes 14, and `TestDecideChargeRows` gets `map-sources` after `binary-seal`.
- Add one shared helper in a new file under `internal/preflight/preflighttest` that writes a decision map with given Sources records. `fixture.go` has 6 lines of headroom. Compose `commitmenttest.SeedClosure` over the conformant seed for the closure cases, and make sure that the deliverable identity matches the spec bytes.
- Repoint the record at `decisions/architecture-planning.md:45-47` to the record that the spec states. After this ticket, a build preflight of this spec prints no `map-sources` red.

## Acceptance

- [ ] `map-sources` reds a closure deletion cited by a compiled map, an active map, and another spec's compiled map, with the spec's detail text (ES01, ES02, ES03, ES20).
- [ ] `map-sources` reds a source that the tip lacks and a source that only the checkout holds, keeps the red for an uncommitted map repoint, and stays green for a kept source (ES04, ES05, ES06, ES07).
- [ ] `map-sources` reds a malformed record, an unknown opener, and an escaping path, and stays green for a well-formed URL record (ES08, ES09, ES10, ES52).
- [ ] The working-tree grammar keeps its messages, and no production preflight file holds a Sources record literal (ES11, ES12).
- [ ] The prediction is empty without a binding and reds for an obligation-free binding and for a policy that does not parse (ES15, ES16, ES17).
- [ ] An absent folder and a folder that holds only skipped names are green; an over-limit map, a link source, and a link candidate map are red (ES18, ES19, ES22, ES23, ES53).
- [ ] Review mode prints no `map-sources` row, the row counts and the charge rows hold the new row, and `--audit-spec` reports `manual_required` true while the row is red (ES21, ES41, ES42, ES46).
- [ ] `decisions/architecture-planning.md` cites no `roadmap/FT375.md`, and this spec's build preflight prints no `map-sources` red (ES48, ES49).
