# A landing that changes no Go file runs a reduced schedule (FT449)

Status: ready

## Destination

A landing whose every changed path is inside the reader table runs the reader
set instead of the whole gate. A changed path outside the reader table runs the
whole gate. A periodic fresh run grades the whole gate. The whole gate stays the
oracle for every other change, and a green reader set never authorizes another
subject.

## Notes

Domain: the gate phase table, the landing, the conformance registry, and the
gate record classes. `CONTEXT.md` owns the terms. Consult `craft-grill`,
`craft-domain`, `craft-research`, and `craft-gate`.

The reviewer reopened the no-component-partition decision on 2026-10-25, but
only for landings that change no Go file. The reviewer also directed that the
research of this map runs inline, with no delegate.

The facts, the measurements, and their citations are in the
[research asset](no-go-schedule/assets/research.md). A Fable review at medium
effort found the gaps that #7, #8, #9, #11, #13, and #14 now close.

Standing preference for the spec: slice the tickets in the spec session itself,
with no delegate. Keep each ticket implementation under 150k tokens of context.

A map-owned asset stays in the map's assets folder,
decisions/no-go-schedule/assets/.

## Decisions so far

- [What are the canonical terms for the census and the reduced schedule?](no-go-schedule/tickets/1.md): the reader table and the reader set.
- [Which changed-path families does the first census cover?](no-go-schedule/tickets/2.md): the record families and the guidance families.
- [Does the reduced schedule apply only when the graded root is the Bench kit?](no-go-schedule/tickets/3.md): the kit only.
- [Which gate entry points run the reduced schedule?](no-go-schedule/tickets/4.md): `bench worktree land` only.
- [Does the reduced schedule run the complete conformance entry?](no-go-schedule/tickets/5.md): yes, every dev-tier check.
- [How does a reduced green record, and what does it authorize?](no-go-schedule/tickets/6.md): a reader record, and no green-marker advance.
- [Who owns the census, and how is it derived?](no-go-schedule/tickets/7.md): each whole run derives it from its own test logs.
- [What bounds a run of reduced landings before a whole run?](no-go-schedule/tickets/8.md): five reader landings, a chain gap, a red, or a new oracle.
- [What happens when the whole run after reduced landings is red?](no-go-schedule/tickets/9.md): the refusal names the suspects, and a red destination goes to `/bench-debug`.
- [What does the ambient dashboard show for a tip that a reader record covers?](no-go-schedule/tickets/12.md): a quiet `reader green (<n> of 5)` row with no command.
- [Which test packages read the live tree through a child process?](no-go-schedule/tickets/10.md): `internal/releaseevidence` and the `system` phase.
- [How does the reader table cover a read that a child process makes?](no-go-schedule/tickets/11.md): the `system` phase always runs. One list names the others.
- [Which changed paths force the whole gate, and which phases does the reader set run?](no-go-schedule/tickets/13.md): embeds, shell files, and links force it.
- [Does the spec amend ADR 0017 and the profile?](no-go-schedule/tickets/14.md): yes, with the reader-set exception.

## Not yet specified

## Spec-writer discretion

- The file layout of the reader tables and the reader records. Each commit
  keeps its own reader record (#8).
- The shim that gives each test binary its own test log, while the `test` phase
  stays one `go test` call.
- The field set and the validator of the reader record class.

## Out of scope

- A reader set for a linked repository (#3).
- A reader set for any gate entry point other than `bench worktree land` (#4).
- A per-check selection inside the conformance entry (#5).
- A tracer for child processes, such as `strace` (#11).
- The reuse of a merge verdict by a later landing, which FT314 owns.
- A landing queue that grades several sources under one gate, which FT450 owns.

## Sources

- Path: `roadmap/FT449.md`
  Supports: the destination and the three named decisions of the row.
  Drift: a change to the row.
- Path: `internal/conformance/registry/checks.go`
  Supports: the count of 44 ordinary checks and 22 catch-all inputs.
  Drift: a change to a registry row or its input source.
- Path: `internal/gate/delta_checkpoint.go`
  Supports: the shape of the one reduced schedule and its record today.
  Drift: a change to `deltaPlan`, `deltaPhases`, or `writeDeltaGreen`.
- Path: `internal/gate/lane_select.go`
  Supports: the path-class table that the commit lane uses today.
  Drift: a change to `laneClasses` or `documentFamilies`.
- Path: `internal/gate/verdict_registry.go`
  Supports: the record classes and their reuse refusals.
  Drift: a change to `verdictRecordClasses`.
- Path: `internal/gate/authorization/authorization.go`
  Supports: the green-marker rules that #6 and #8 keep.
  Drift: a change to `checkMarker`.
- Path: `internal/status/status.go`
  Supports: the stale gate row that #12 replaces for a reader-record tip.
  Drift: a change to `appendGateInfo` or `staleGateDetailAction`.
- Path: `internal/systemtest/candidate_install_test.go`
  Supports: the record reads of the `system` phase in #10 and #11.
  Drift: a change to the kit files the candidate-install journeys copy.
- Path: `internal/gate/subject.go`
  Supports: the subject identity that the oracle identity of #8 reuses without the tree.
  Drift: a change to `buildSubjectForTree`.
- Path: `internal/gate/phases.go`
  Supports: the shellcheck inventory that #13 sends to the whole gate.
  Drift: a change to `shellcheckFiles`.
- Path: `internal/gate/gate_go.go`
  Supports: the race-test table that #13 reads.
  Drift: a change to `raceTests`.
- Path: `docs/adr/0017-the-worktree-commit-runs-the-fast-lane.md`
  Supports: the whole-gate landing rule that #14 amends.
  Drift: a change to the ADR.
- Path: `decisions/no-go-schedule/assets/research.md`
  Supports: the measurements and the probes behind #2, #7, #8, #10, and #11.
  Drift: a new landing mix, a new reader package, or a new Go toolchain.
