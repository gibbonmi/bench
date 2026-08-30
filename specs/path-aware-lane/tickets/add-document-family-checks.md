# Add the document-family checks to the kit lane

Blocked by: select-kit-lane-checks.md
Writes: internal/gate/lane_select.go, internal/gate/lane_select_test.go, internal/gate/lane.go, internal/gate/lane_test.go, internal/conformance/profile_lane_table_test.go, projects/benchkit.md

## What to build

A commit that changes a roadmap row, a decision map, a retro, or the kit profile
reds at the commit rather than at the landing. The document check runs inside the
lane, so the defect never reaches a full gate run.

The package `internal/gate` extends the path-class table with four rows. The
class names are exactly the `InputSource` values that
`internal/conformance/registry` declares: `roadmap-board`,
`decision-documents`, `capture-retros`, and `benchkit-profile`. The predicates
are `ROADMAP.md` or the `roadmap/` prefix, the `decisions/` prefix or a
`specs/<slug>/decisions/` prefix, the `capture/retros/` prefix, and
`projects/benchkit.md`. A roadmap Markdown path matches both `markdown` and
`roadmap-board`, so both classes contribute their checks.

`BenchkitLane` declares one check per dev-tier registry check whose `Inputs`
equals one of those four values. Each check carries the registry check's own
name and the argv `<run binary> test --check <name>`, so the lane builds no
second executable. The row list comes from the registry table, so a check the
registry adds joins the lane with no second list. A document check that exits
non-zero fails the lane as `lane{outcome=fail,check=<name>}`, by the existing
rule.

Each class binds to its checks through the registry's `Inputs` binding, so the
family-to-check fact has one source. The class-name vocabulary that
`select-kit-lane-checks.md` fixed stays unchanged, and this ticket only adds
rows.

The file `projects/benchkit.md` gains one lane row per document check, with the
`selected by` cell that `checkProfileLaneTable` renders from the class table. The
check's live-tree row grades the new advertisement.

The lane rows run `RunLane` with `BenchkitLane` and a stub run-binary factory,
by the precedent of `internal/worktree/test_run_test.go`. The stub records its
argv and exits non-zero for the one row that proves a red document check. The
class rows extend the existing table test over `SelectLane`.

## Acceptance

- [ ] PL28: `BenchkitLane` holds one `<run binary> test --check <name>` check per dev-tier `registry.Checks` entry whose `Inputs` is a document class, and no other `test --check` row. The test enumerates its expectation from `registry.Checks`.
- [ ] PL29: one `M` entry for `roadmap/FT1.md` returns `prose` and `roadmap-detail-integrity`, and the classes `markdown,roadmap-board`.
- [ ] PL30: one `M` entry for `specs/x/decisions/map.md` returns `prose` and `decision-map-integrity`.
- [ ] PL31: one `M` entry for `capture/retros/x.md` returns `prose` and `retro-improvement-markers`.
- [ ] PL32: one `M` entry for `projects/benchkit.md` returns `prose`, `guidance-prose-budgets`, and `profile-lane-table`.
- [ ] PL33: every document class name is a valid registry `InputSource`, and `registry.Checks` binds at least one dev-tier check to it.
- [ ] PL34: every document check in `BenchkitLane` carries the run-binary token as `argv[0]`.
- [ ] PL41: each lane check named after a gate check runs that check's bound over the composed checkout, per the spec's counterpart table. The Standards axis reads each argv against the phase table and the registry.
- [ ] PL35: `RunLane` with a stub run binary that exits 1 on `test --check roadmap-detail-integrity`, and a tree that changes `ROADMAP.md`, fails with `Check` equal to `roadmap-detail-integrity`.
- [ ] PL37: `checkProfileLaneTable` over the live tree returns no diagnostic, with the `selected by` column and the document rows in place.
- [ ] the gate `test` phase stays green for the whole `internal/gate` package.
- [ ] the gate `test` phase stays green for the whole `internal/conformance` package.
