# `bench test` projection faces (FT290)

Status: ready

## Destination

`bench test` gives a caller the facts that the caller now reads from test
source, from a second run, or not at all. The scope is one spec with three
outcomes.

- The inventory faces: `bench test --checks` lists the named checks, and
  `bench test --check <name> --fixtures` lists the fixtures a check owns.
- The result evidence: a `check` row and a `tests_run` count prove a run, the
  prose check prints a result on green, the failures table shows each
  diagnostic, `--check system` accepts `--run`, and the unknown-check refusal
  names the running executable.
- The cause cell: each `bench test --changed` packages row names why the
  selection holds that package.

## Notes

Domain: the `bench test` verb and its result projection. Charge
`bench-craft-domain` and `bench-craft-cli` for the spec.

Two executables are in play, and the terms stay apart. The **running
executable** is the Bench executable that runs the verb, and it owns the check
name set. The **run binary** is the source-bound executable that a run selects
for its Go child. Do not write "the binary" for one of them.

A map-owned asset stays in the map's assets folder,
decisions/ft290-test-projection/assets/.

## Decisions so far

- [Which occurrence items join the destination?](ft290-test-projection/tickets/1.md): system `--run` and the failures row unit join; the compile error is a shipped exclusion.
- [Which identity does the unknown-check refusal name?](ft290-test-projection/tickets/2.md): the running executable path and the source commit of its seal.
- [What proves that the selected check ran?](ft290-test-projection/tickets/3.md): a `check` row with a `tests_run` count, zero exits 1, and a `tests_run` cell on each packages row.
- [What does the fixtures face print?](ft290-test-projection/tickets/4.md): one row for each owned fixture and no test run. A check with no family prints an empty table at exit 0.
- [What does the prose check print on green?](ft290-test-projection/tickets/5.md): the `check` row with the graded subject count, no packages table, and the subjects under `--full`.
- [How does the changed form explain a widened set?](ft290-test-projection/tickets/6.md): one `selected_by` cause for each package, by the precedence `go-metadata`, `changed`, `embed`, `imports <package>`.
- [What does the check inventory print?](ft290-test-projection/tickets/7.md): one row for each named check with its kind and its family count, at exit 0.
- [What is the failures row unit?](ft290-test-projection/tickets/8.md): one row for each failed test with a `lines` count, and one row for each diagnostic line under `--full`.
- [Which named checks accept a run pattern?](ft290-test-projection/tickets/9.md): `system` only, and a pattern with no match exits 1.
- [Does the shaped scope split?](ft290-test-projection/tickets/10.md): the run binary provenance moves to its own map, and the other outcomes stay as one spec.

## Not yet specified

## Spec-writer discretion

- The internal owner of each new row producer, when the output stays as the tickets fix it.

## Out of scope

- The compile error on a non-compiling `--package`; commit `4ff47076` shipped it.
- The run binary provenance and the stale-test occurrence; `decisions/run-binary-provenance.md` owns them.
- A `--run` filter for a named check other than `system`.
- The pinned repository paths of a fixture in the `--fixtures` rows.

## Sources

- Path: `roadmap/FT290.md`
  Supports: the destination and the occurrence record behind tickets #1 to #9.
  Drift: a new occurrence or a body edit on the row.
- Path: `internal/testreport/command.go`
  Supports: tickets #2, #5, and #9. The name check comes before the run binary selection. The prose check returns no output on green, and `--check` refuses `--run`.
  Drift: a change to `parseFocusedRequest`, `runProseCheck`, or `runNamedCheck`.
- Path: `internal/testreport/testreport.go`
  Supports: tickets #3 and #8: the packages row has no test count, and one failures row covers one failed test.
  Drift: a change to `render`, `failures`, or `failureCell`.
- Path: `internal/testreport/selection.go`
  Supports: ticket #6: the four selection causes and the reverse-import closure.
  Drift: a change to `selectCurrentPackages`.
- Path: `internal/conformance/registry/registry.go`
  Supports: tickets #4 and #7: `CanaryFamilies` binds a check to its families, and `Names` owns the check set.
  Drift: a change to the family binding or to the check set owner.
