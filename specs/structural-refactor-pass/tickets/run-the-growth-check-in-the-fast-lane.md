# Run the growth check in the fast lane

Blocked by: add-the-structure-growth-mode.md
Writes: internal/gate/lane.go, internal/gate/lane_output.go (new), internal/gate/lane_test.go, internal/gate/lane_run_test.go (new), internal/gate/lane_select.go, internal/gate/lane_select_test.go, internal/gate/authorization/authorization.go, internal/commit/lane_test.go, internal/commit/lane_structure_test.go (new), internal/conformance/profile_lane_table_test.go, projects/benchkit.md, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership, tests/canary/workflow-guidance-anchors/benchkit-system-suite-route
Covers: SR52, SR53, SR54, SR55, SR58, SR59

## What to build

The line is opus / medium. Give the built-in lane a Bench-owned `structure`
check after `build`. Its argv names the growth flag `--growth` and a base
token. The lane request gains a base field, the caller sets it, and the lane
resolver replaces the base token with that base. The go-source class selects
the check, and a path outside every class already selects every check. The
check runs in the fast lane only, on `bench commit` and `bench worktree
merge`; the landing gate is unchanged, per decision (k).

The profile lane table in `projects/benchkit.md` gains the `structure` row
with `bench structure --growth <base>` selected by go-source. The profile
check renders the base token as `<base>`, so the advertisement matches the
lane.

Pins change on purpose, per decision (n): `TestBenchkitLaneTable` gains the
structure row, `TestResolveLane` gains the base replacement, and three rows
of `TestSelectLaneByClass` gain `structure`. `TestLaneClassesNameOnlyDeclaredChecks`
derives its expectation and needs no edit.

Split two over-budget files, per decision (o). The lane diagnostics and the
tap writer move to the new file `internal/gate/lane_output.go`. The moved
symbols are `laneDiagnostics`, `record`, `firstLine`, `laneWriters`,
`discardIfNil`, `laneTapWriter`, `Write`, and `flush`. The lane execution
rows move to the new file `internal/gate/lane_run_test.go`, with
`laneRecordsArgv` and `laneRefusesArgv`; `laneCheck` stays in
`internal/gate/lane_test.go`. Every existing test function keeps its name.

This ticket reads the growth contract from the sibling ticket's diff: the
`--growth` flag spelling, the base operand, and the `FILE GREW` row label.
`capture/restructure-backlog.md` already holds the census, and this ticket
does not edit it. The review compares its rows to the `bench structure`
output at 8eea2d15.

## Acceptance

- [ ] A `bench commit` that raises an over-budget Go file's count fails the lane with `lane{outcome=fail,check=structure}` and the `FILE GREW` line.
- [ ] A `bench commit` that lowers that count passes the lane.
- [ ] The built-in lane carries a Bench-owned `structure` check after `build` whose argv names the growth flag and the base token.
- [ ] The go-source class selects the `structure` check.
- [ ] The lane resolver replaces the base token with the request's base.
- [ ] The profile lane table carries the `structure` row with `bench structure --growth <base>` selected by go-source, and `bench test --check profile-lane-table` stays green.
- [ ] `capture/restructure-backlog.md` holds one row with a disposition for every file `bench structure` reports over budget at 8eea2d15.
- [ ] `bench structure --growth <part-1 tip>` over the Part 2 diff prints the ok line at exit 0.
- [ ] Self-probe: omit the base-token replacement in the resolver, and report the observed red.
