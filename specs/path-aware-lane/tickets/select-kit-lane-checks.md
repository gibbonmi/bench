# Select the kit lane's checks by path class

Blocked by: export-embed-target-derivation.md, derive-composed-changes.md
Writes: internal/gate/lane_select.go, internal/gate/lane_select_test.go, internal/gate/lane.go, internal/gate/lane_test.go, internal/gate/lane_record_test.go, internal/gate/gate_prose_test.go, internal/gate/authorization/authorization.go, internal/gate/authorization/lane_test.go, internal/commit/lane_test.go, internal/conformance/profile_lane_table_test.go, projects/benchkit.md

## What to build

A worktree `bench commit` on the kit root runs only the checks its composed
changes select. A Markdown-only commit pays no Go toolchain, and a Go-only
commit pays no prose check.

The package `internal/gate` owns the path-class table in
`internal/gate/lane_select.go`. Each row binds a class name to a path predicate
and to check names. This ticket declares four rows:

- `go-source` — a `.go` suffix
- `go-build-input` — `go.mod`, `go.sum`, or an embed target
- `markdown` — a `.md` suffix
- `prose-policy` — `.bench/prose-exclusions`

The `go-build-input` predicate reads the embed targets from
`packagesurface.EmbedTargets`. A path matches every class whose predicate holds.
A path that matches no class, a `120000` mode on either side, or a `160000`
mode is `unknown`.

The class-name vocabulary is a shared contract.
`add-document-family-checks.md` extends this table with four more rows. Those
rows carry the registry's `InputSource` names. So keep the row shape open to a
class that names a lane check by name.

The class table and `SelectLane` live in `internal/gate/lane_select.go` beside
the `Lane` type, so `internal/gate/lane.go` does not grow. `SelectLane(checks,
changes)` returns the declared checks the classes select, in declared order, and
the class names in table order. The `unknown` class selects
every declared check, and a duplicate never enters the result. `RunLane` applies
`SelectLane` when the lane is selective, and it leaves a manifest lane alone.
`LaneResult` gains the selected check names and the selected class names, so the
authority prints the lane line without a second selection. The lane record
schema does not change, and the run still touches neither the verdict cache nor
the evidence store.

The type `gate.Lane` gains `Selective`. `LaneForCommit` sets it true for the
kit's built-in lane and false for a manifest lane, so the selection has one
switch. The package `internal/gate/authorization` prints
`lane{outcome=pass,checks=<names>,classes=<classes>}` for a selective lane and
keeps `lane{outcome=pass,checks=<names>}` for a manifest lane. A fail keeps
`lane{outcome=fail,check=<name>}`.

The file `projects/benchkit.md` gains a third lane-table column, `selected by`,
that lists the classes per check. `internal/gate` exports the class table in
table order, and `checkProfileLaneTable` renders each `selected by` cell from
that one source. So a stale cell reds and no second class list exists.

The selection tests are a table test over `SelectLane` with literal change
entries. The end-to-end rows run `RunLane` with `BenchkitLane` and a stub
run-binary factory, by the precedent of `internal/worktree/test_run_test.go`.
The stub records its argv, and the fixture holds no `go.mod`, so a selected
`vet` or `build` reds the lane and proves the selection. The commit rows use the
manifest fixture that `laneRepo` builds. The profile rows call
`checkProfileLaneTable` on a fixture root and on the live tree.

## Acceptance

- [ ] PL9: `SelectLane` over the kit lane with one `M` entry for `internal/x/y.go` returns `gofmt`, `vet`, and `build` in that order, and the class `go-source`.
- [ ] PL10: one `M` entry for `go.sum` returns `vet` and `build`, and the class `go-build-input`.
- [ ] PL11: one `M` entry for `internal/adopt/prepush.sh`, with that path among the embed targets, returns `vet` and `build`, and the class `go-build-input`.
- [ ] PL12: one `M` entry for `docs/note.md` returns `prose` alone and the class `markdown`.
- [ ] PL13: one `M` entry for `.bench/prose-exclusions` returns `prose` alone and the class `prose-policy`.
- [ ] PL14: `bench gate-prose <root>` with no path and an exclusion row without a reason exits 1 and prints `malformed exclusion row`.
- [ ] PL15: entries for `a.go` and `b.md` return `gofmt`, `prose`, `vet`, and `build` in declared order, and the classes `go-source,markdown`.
- [ ] PL16: one `M` entry for `bin/bench.sh` returns every declared check and the class `unknown`.
- [ ] PL17: one `M` entry for `x.go` with source mode `120000`, and one with destination mode `120000`, each return every declared check and the class `unknown`.
- [ ] PL18: one `M` entry with destination mode `160000` returns every declared check and the class `unknown`.
- [ ] PL19: `RunLane` with `BenchkitLane`, `Selective` true, a stub run binary, and a tree that changes `note.md` alone passes. The stub recorded one `gate-prose` invocation and no `gate-go` invocation.
- [ ] PL20: that run's stdout holds `lane{outcome=pass,checks=prose,classes=markdown}`.
- [ ] PL21: the manifest fixture commit of `note.md` alone prints `lane{outcome=pass,checks=check,gofmt,prose,build}` and no `classes=` cell.
- [ ] PL22: `LaneForCommit` at a root equal to `BENCH_KIT` returns `Selective` true, and at a manifest root returns `Selective` false.
- [ ] PL23: `bench commit --dry-run` with the manifest fixture prints the lane line and moves no ref.
- [ ] PL24: after a selective `RunLane` the lane record names the tree, the lane, and the outcome. The gate cache and the evidence store stay absent.
- [ ] PL36: `checkProfileLaneTable` over a fixture profile whose `gofmt` row reads `selected by` `markdown` returns one diagnostic that names `gofmt`.
- [ ] PL38: `SelectLane` over any entry set returns a subsequence of the declared lane.
- [ ] PL39: `BenchkitPhases` keeps the `test` phase argv `go test -trimpath -count=1 ./...`.
- [ ] PL42: one `M` entry for `templates/a.txt` with the embed target list `templates/*` returns every declared check and the class `unknown`.
- [ ] PL47: a `D` entry for `internal/x/y.go` returns `gofmt`, `vet`, and `build`, and a `D` entry for the embed target `internal/adopt/prepush.sh` returns `vet` and `build`.
- [ ] the gate `test` phase stays green for the whole `internal/commit` package.
