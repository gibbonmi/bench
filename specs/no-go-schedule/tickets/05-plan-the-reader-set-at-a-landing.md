# Plan the reader set for a gate run that carries a landing value

Blocked by: 04-run-the-reader-set-schedule.md
Writes: internal/gate/reader_plan.go (new), internal/gate/reader_plan_test.go (new), internal/gate/run_transaction.go, internal/git/history.go (new), internal/git/history_test.go (new), internal/packagesurface/embed_directives.go (new), internal/packagesurface/embed_directives_test.go (new)
Covers: NG37, NG38, NG39, NG40, NG41, NG42, NG43, NG44, NG45, NG46, NG47, NG48, NG49, NG50, NG51, NG52, NG53, NG54, NG73, NG74, NG75, NG76, NG77

## What to build

`gate.WithLanding(ctx, destination, marker)` puts the destination commit and the green marker in the gate context, as `gate.WithCompletion` does. This ticket adds no production caller; ticket 06 adds the one in `landReviewed`.

The reader plan reads that value and runs the cases of the spec's "The landing plan" section in order. The first case that holds sends the run to the whole gate. This ticket implements every case except case 9, the reader hold, which ticket 07 adds. When a case holds, the gate prints `gate: whole gate (reader set refused: <case>)` before the phases. Otherwise the plan builds the reader schedule of ticket 04 and runs it.

The reader packages are the union of the table's readers of each changed family, the child readers of each changed family, and `./internal/conformance`. The race flag holds when a package of the race-test table (`internal/gate/gate_go.go:34-41`) is among them.

The changed paths come from `ComposedChanges` between the destination and the graded tree, with both modes (`internal/gate/lane_select.go:26-33`). The embed case reads a new scan in `internal/packagesurface` that reads every Go file of the tree, test files included. It matches a changed path by the directive's exact path, by a path below it, or by its glob. The existing `EmbedTargets` keeps its contract for the lane. The shellcheck case reads `shellcheckFiles`, and the plan tests inject the inventory.

`internal/git/history.go` adds `FirstParentsAfter(root, marker, destination string) ([]string, error)`. It lists the first-parent commits after the green marker up to the destination, oldest first. The chain case asks `gate.ReaderRecorded` for the tree of each listed commit. Tickets 07 and 08 read the chain through the same function.

The census forbids a direct process or `git init` constructor in the gate's and the Git package's tests (`internal/conformance/ordinary_build_census_test.go:69-80`). Build the fixtures through the `gittest` helpers.

`internal/gate/run_transaction.go` holds 397 lines. Put the plan and its dispatch in `reader_plan.go`, and keep `run_transaction.go` at 400 lines or fewer.

## Acceptance

- [ ] A landing value over a covered change, with a valid table, an equal identity, and a clean chain, runs the reader set.
- [ ] The reader packages are the changed families' readers, their child readers, and `./internal/conformance`, each once.
- [ ] The race flag holds when a race-table package reads a changed family.
- [ ] A change to `AGENTS.md`, and a rename of `docs/x.md` to `x.md`, each run the whole gate.
- [ ] A change to an exact embed target, to a file below an embedded directory, to a glob match, and to a test-file embed target each run the whole gate.
- [ ] A change to a path of an injected shellcheck inventory runs the whole gate.
- [ ] An added link, a removed link, an added gitlink, and a removed gitlink in `docs/` each run the whole gate.
- [ ] A missing table, a different identity digest, a chain gap, and an absent marker each run the whole gate.
- [ ] Five reader-record commits after the marker run the whole gate, and four run the reader set.
- [ ] The same covered change without a landing value, and in a linked root, each run the whole gate.
- [ ] A fixture where two cases hold prints `gate: whole gate (reader set refused: <case>)` with the first case.
