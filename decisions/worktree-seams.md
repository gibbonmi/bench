# Worktree joins seams and the verb runner (FT356)

Status: shaping

## Destination

The worktree package keeps only the joins seams that make an unreproducible
failure reproducible, and its tests call every verb through one verb runner.
The scope is two specs from this map, in this order.

- The verb runner spec: one test helper calls a verb entry and returns a verb
  result. Every run wrapper, fingerprint extractor, and positional fixture tuple
  in the package moves onto it, and the census refuses a direct verb call.
- The seam reduction spec: the joins value keeps only the fields that pass the
  keep rule. Each verb entry reads the kit root, the Bench home, and the clock
  once. The census refuses a read below a verb entry, and the serial ceiling
  drops.

## Notes

Domain: the worktree package's joins value, its verb entries, and its test
harness. Charge `bench-craft-domain`, `bench-craft-seams`, and
`bench-craft-gate` for the specs.

The terms stay apart. A **joins seam** is one field of the joins value that a
test replaces. A **verb entry** is the exported command function of one
worktree verb. The **verb runner** is the one test helper that calls a verb
entry. A **verb result** is what the verb runner returns: the exit code, the
rows read through `axitest`, and the fingerprint. `CONTEXT.md` holds the Avoid
lists.

A map-owned asset stays in the map's assets folder,
decisions/worktree-seams/assets/.

## Decisions so far

- [Does FT356 ship as one spec or as two?](worktree-seams/tickets/1.md): two specs from this map; the verb runner spec lands first.
- [Which predicate keeps a joins seam?](worktree-seams/tickets/2.md): a field stays for a fixture-proof fault or for a real compile or gate run.
- [What happens to a stubbed test whose seam goes?](worktree-seams/tickets/3.md): the test converts to a real fixture, and the count pin holds.
- [Does the gate enforce the single read at the verb entry?](worktree-seams/tickets/4.md): the census refuses a read below a verb entry, and the serial ceiling drops.
- [How much of the test tree does the runner spec migrate?](worktree-seams/tickets/5.md): all run wrappers, fingerprint extractors, direct verb calls, and fixture tuples.
- [Where does the verb runner live?](worktree-seams/tickets/6.md): the worktree package only, on `axitest.DecodeDocument`.
- [Does the gate refuse a direct verb call in a test?](worktree-seams/tickets/7.md): yes, the census reports a verb entry call outside the verb runner file.
- [Does the runner spec carry a size target?](worktree-seams/tickets/8.md): no; the acceptance is structural.

## Not yet specified

## Spec-writer discretion

- The verb result fields beyond the exit code, the rows, and the fingerprint.
- The value that carries the kit root, the Bench home, and the clock below a
  verb entry.
- The chunk order of the verb runner migration.

## Out of scope

- The test harness of every package other than the worktree package.
- A line-count target for the worktree tests.
- A deletion or a merge of a top-level test.

## Sources

- Path: `roadmap/FT356.md`
  Supports: the destination and tickets #1, #2, and #5.
  Drift: a new occurrence or a body edit on the row.
- Path: `internal/worktree/joins.go`
  Supports: tickets #2, #9, and #10: the 29 fields and their defaults.
  Drift: a field added to or removed from `joins` or `defaultJoins`.
- Path: `internal/worktree/parallel_census_test.go`
  Supports: tickets #3, #4, and #7: the count pin of 664, the serial ceiling of 46, and the package-variable refusal.
  Drift: a change to `worktreeTestCount`, `worktreeSerialCeiling`, or `parallelCensus`.
- Path: `internal/worktree/effects.go`
  Supports: ticket #4: `Home` and `currentTime` are the home and clock reads that source files call directly.
  Drift: a change to `Home` or `currentTime`.
