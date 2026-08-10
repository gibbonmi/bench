# Pin every default/full pair and empty class

Blocked by: compare-four-observations.md
Ownership fence: `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`, `specs/axi-compatibility-oracle/testdata`
Integration surfaces: four-observation comparator→compare-four-observations.md; pair and empty case fixtures→`specs/axi-compatibility-oracle/testdata`; class derivation for pairs and empties→`internal/axi/compatibility`; declared pair and empty owners→`decisions/byte-preserving-axi-foundation/assets/ft173-helper-compatibility-census.md` exercised unchanged by every EF1 row except the two below; dashboard empty-section body→`internal/dashboard/render.go` exercised unchanged by EF1/empty-dashboard-section; roadmap recommended-sequence gap line→`internal/roadmap/roadmap.go` exercised unchanged by EF1/empty-roadmap-recommended-sequence; bound-edge consumer→pin-truncation-bound-edges.md; byte-class consumer→pin-toon-byte-classes.md
Contracts: the paired case IDs `<member>-default` and `<member>-full`, and the empty-form case IDs `<member>-empty`, cross `decisions/byte-preserving-axi-foundation/assets/ft173-helper-compatibility-census.md`→`specs/axi-compatibility-oracle/testdata`; their type is one baseline observation record per case ID, membership is the six declared default/full pairs, the three empty forms that census declares — the TOON zero-row table, spec build's one-row `state=empty`, and status's prose clean line — and two more the tree produces without the census naming them, the dashboard's empty section body from `internal/dashboard/render.go` and roadmap's recommended-sequence gap line from `internal/roadmap/roadmap.go`, order is stable case ID ascending, and a pair missing either half refuses the load; asserted by EF1 against the really rebuilt candidate executable
Closure: EF1/pair-diff, EF1/pair-test, EF1/pair-outline, EF1/pair-roadmap-context, EF1/pair-worktree-clean, EF1/pair-spec-build-status, EF1/empty-toon-zero-row, EF1/empty-spec-build-state, EF1/empty-status-prose, EF1/empty-dashboard-section, EF1/empty-roadmap-recommended-sequence

## What to build

The comparator exists, so CO5's seam is now available; this ticket fills its
default/`--full` and empty half. It derives a `<member>-default` and
`<member>-full` case for each of the six pairs the helper census declares —
`bench diff`, `bench test`, `bench outline`, `bench roadmap --context`,
`bench worktree clean`, and `bench spec build status` — and one `<member>-empty`
case for each of five distinct empty forms. Three come from the census, which
keeps them apart deliberately: the TOON zero-row table `name[0]{fields}:`, spec
build's one-row `spec_build` with `state=empty`, and status's prose
`bench: clean — nothing pending`. Two more are the tree's, not the census's: a
dashboard card renders an explicit empty body — `internal/dashboard/render.go`
prints `<p class="empty">…</p>` inside the section rather than dropping the
section — and `bench roadmap` with no drain prints the recommended-sequence
section verbatim, or, when it is absent, the explicit gap line `nextAction`
returns in `internal/roadmap/roadmap.go` rather than nothing at all.

The five empty forms mean five different things, and the census is explicit that a
shared renderer must not normalize its three into one spelling. Each row below
breaks exactly one of those meanings.

Every mutation is applied to a scratch copy of the tree from which the candidate
executable is rebuilt through `scripts/go-build.sh`, so the mutated subject is
real production rendering code and the oracle observes it the way an agent would.
The rebuild runs under a 180s `context.WithTimeout` and every case child under
30s, so a mutation that hangs a renderer fails as a bounded deadline.

## Acceptance

- [ ] [EF1] (covers CO5) each of the six declared default/`--full` pairs and each of the five empty forms compares byte-exact against the pinned baseline, and a candidate rebuilt with any one of those renderings changed reports a raw stdout delta on the case that owns it.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| EF1/pair-diff | make `bench diff --full` emit the bounded default body by dropping the log and diff body append in the candidate rebuild | the exact-pair test | run `go test ./cmd/bench -run TestExactMatrixPreservesDefaultFullPairs/diff -timeout 900s`; it fails at the raw stdout equality assertion for case `root-diff-full`, reporting the missing diff body while case `root-diff-default` still matches; the rebuild is bounded at 180s and the case children at 30s |
| EF1/pair-test | make `bench test --full` keep `sanitize.Preview` instead of switching to the uncapped `sanitize.Controls` in the candidate rebuild | the exact-pair test | run `go test ./cmd/bench -run TestExactMatrixPreservesDefaultFullPairs/test -timeout 900s`; it fails at the raw stdout equality assertion for case `root-test-full`, reporting the `… (N bytes)` suffix present where the baseline has the complete diagnostic; bounded by the 180s rebuild and 30s case deadlines |
| EF1/pair-outline | keep the `bounds.OutlineRowLimit` slice applied under `--full` in the candidate rebuild | the exact-pair test | run `go test ./cmd/bench -run TestExactMatrixPreservesDefaultFullPairs/outline -timeout 900s`; it fails at the raw stdout equality assertion for case `root-outline-full`, reporting 200 emitted rows against the baseline's complete row set; bounded by the 180s rebuild and 30s case deadlines |
| EF1/pair-roadmap-context | make `bench roadmap --context --full` return the 4096-byte capped body in the candidate rebuild | the exact-pair test | run `go test ./cmd/bench -run TestExactMatrixPreservesDefaultFullPairs/roadmap_context -timeout 900s`; it fails at the raw stdout equality assertion for case `root-roadmap-context-full`, reporting `truncated=true` and the short body where the baseline carries the complete raw source; bounded by the 180s rebuild and 30s case deadlines |
| EF1/pair-worktree-clean | keep the 20-path default display under `bench worktree clean --full` in the candidate rebuild | the exact-pair test | run `go test ./cmd/bench -run TestExactMatrixPreservesDefaultFullPairs/worktree_clean -timeout 900s`; it fails at the raw stdout equality assertion for case `nested-worktree-clean-full`, reporting 20 `ignored_paths` rows against the baseline's up-to-1000 set; bounded by the 180s rebuild and 30s case deadlines |
| EF1/pair-spec-build-status | drop the empty assignments and review blocks `--full` appends in the candidate rebuild | the exact-pair test | run `go test ./cmd/bench -run TestExactMatrixPreservesDefaultFullPairs/spec_build_status -timeout 900s`; it fails at the raw stdout equality assertion for case `nested-spec-build-status-full`, reporting the two missing blocks while the default case still matches; bounded by the 180s rebuild and 30s case deadlines |
| EF1/empty-toon-zero-row | render a nil row slice as blank output instead of `name[0]{fields}:` in the candidate rebuild | the exact-empty test | run `go test ./cmd/bench -run TestExactMatrixPreservesEmptyClasses/toon_zero_row -timeout 900s`; it fails at the raw stdout equality assertion for case `root-learnings-empty`, reporting empty bytes where the baseline carries the schema-bearing zero-row header; bounded by the 180s rebuild and 30s case deadlines |
| EF1/empty-spec-build-state | replace spec build's one-row `state=empty` projection with a zero-row table in the candidate rebuild | the exact-empty test | run `go test ./cmd/bench -run TestExactMatrixPreservesEmptyClasses/spec_build_state -timeout 900s`; it fails at the raw stdout equality assertion for case `nested-spec-build-status-empty`, reporting `spec_build[0]` against the baseline's `spec_build[1]` row with `state=empty`; bounded by the 180s rebuild and 30s case deadlines |
| EF1/empty-status-prose | replace status's `bench: clean — nothing pending` line with a zero-row signals table in the candidate rebuild | the exact-empty test | run `go test ./cmd/bench -run TestExactMatrixPreservesEmptyClasses/status_prose -timeout 900s`; it fails at the raw stdout equality assertion for case `root-status-empty`, printing both lines; bounded by the 180s rebuild and 30s case deadlines |
| EF1/empty-dashboard-section | omit an empty dashboard section instead of rendering its empty body in the candidate rebuild | the exact-empty test | run `go test ./cmd/bench -run TestExactMatrixPreservesEmptyClasses/dashboard_section -timeout 900s`; it fails at the raw stdout equality assertion for case `root-dashboard-empty`, reporting the missing section heading in the `--stdout` body; bounded by the 180s rebuild and 30s case deadlines |
| EF1/empty-roadmap-recommended-sequence | drop the recommended-sequence block when roadmap has no drain in the candidate rebuild | the exact-empty test | run `go test ./cmd/bench -run TestExactMatrixPreservesEmptyClasses/roadmap_recommended_sequence -timeout 900s`; it fails at the raw stdout equality assertion for case `root-roadmap-empty`, reporting the absent block against the baseline's empty-but-present one; bounded by the 180s rebuild and 30s case deadlines |
