## Outcome

FT290 landed at `e558624e` from source `16970a22..b9b24d23`, and the landing gate was green. `bench test` now proves what a named check ran and shows each failure diagnostic. It also filters the system suite, names the run executable, lists the check fixtures and the check inventory, and explains each `--changed` selection. Ten tickets in five chunks covered 56 acceptance rows. The delivery closed `roadmap/FT290.md`, and `decisions/run-binary-provenance.md` keeps the unresolved half.

## Gate-stage timings

- landing: commit e558624eb20b848ed4c2b6680d04fa8cc0f17557, trace 6e1abd6e6314185e632c8d9ff33d7479
- gofmt: 145 ms
- vet: 1550 ms
- test: 277299 ms
- race: 3131 ms
- system: 87195 ms
- shellcheck: 610 ms

## Ticket-versus-spec-slice and delegate performance

Ten fresh Opus/medium authors each committed their ticket first-pass inside the planned chunk, and each named probe bit. Fourteen fresh Opus/medium repair sessions closed the review findings. Ticket 8 took three rounds and ticket 10 took four rounds. Sonnet/high ran every review axis, and the review record holds no refuted finding.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| TP-C1a Standards | S1 named_check.go:57,61 | held | 0.6 | finding | Sonnet / high / reviewer |
| TP-C1a Standards | S2 tests_run_test.go duplicated fixture | held | 0.5 | finding | Sonnet / high / reviewer |
| TP-C1a Standards | S3 check_row_test.go comment claims | held | 0.5 | finding | Sonnet / high / reviewer |
| TP-C1a Coverage | C1 named_check.go:60-61 first block | held | 0.6 | finding | Sonnet / high / reviewer |
| TP-C1a Coverage | C2 named_check.go:62 | held | 0.6 | finding | Sonnet / high / reviewer |
| TP-C1a Coverage | C3 testreport.go:72 | held | 0.4 | finding | Sonnet / high / reviewer |
| TP-C1a Coverage | C4 named_check.go:60 | held | 0.4 | finding | Sonnet / high / reviewer |
| TP-C1b Standards | S1 outcome.go comment aging | held | 0.5 | finding | Sonnet / high / reviewer |
| TP-C1b Coverage | C1 prose --full drops findings | held | 0.9 | finding | Sonnet / high / reviewer |
| TP-C1b Coverage | C2 prose directory exclusion row | held | 0.8 | finding | Sonnet / high / reviewer |
| TP-C2 Standards | S1 usage line typed twice | held | 0.7 | finding | Sonnet / high / reviewer |
| TP-C2 Standards | S2 cmd/bench help repeats the usage | held | 0.5 | finding | Sonnet / high / reviewer |
| TP-C2 Standards | S3 unknown_check_test.go:255-258 | held | 0.5 | finding | Sonnet / high / reviewer |
| TP-C2 Coverage | C1 grammar line untested | held | 0.6 | finding | Sonnet / high / reviewer |
| TP-C2 Coverage | C2 refusal identity fallback | held | 0.5 | finding | Sonnet / high / reviewer |
| TP-C3 Standards | S1 fifth canary root join | held | 0.5 | finding | Sonnet / high / reviewer |
| TP-C3 Standards | S2 header comments name unrecorded reds | held | 0.4 | finding | Sonnet / high / reviewer |
| TP-C3 Spec | P1 --checks exit on invalid inventory | held | 0.4 | finding | Sonnet / high / reviewer |
| TP-C3 Coverage | C1 fixtures sort order | held | 0.8 | finding | Sonnet / high / reviewer |
| TP-C3 Coverage | C2 fixtures.go:343-347 | held | 0.7 | finding | Sonnet / high / reviewer |
| TP-C3 Coverage | C3 command.go:72,75 TP39 | held | 0.4 | finding | Sonnet / high / reviewer |
| TP-C3 confirming Coverage | C1 rootRecords swallows a real error | held | 0.7 | finding | Sonnet / high / reviewer |
| TP-C3 second confirming Coverage | C5 FixturePins callers silent | held | 0.8 | finding | Sonnet / high / reviewer |
| TP-C4 Standards | S1 embed branch versus causePrecedence | held | 0.4 | finding | Sonnet / high / reviewer |
| TP-C4 Standards | S2 changed form signalled by nil | held | 0.5 | finding | Sonnet / high / reviewer |
| TP-C4 Coverage | C1 unselected dependency skip | held | 0.9 | finding | Sonnet / high / reviewer |
| TP-C4 Coverage | C2 TestImports cause | held | 0.9 | finding | Sonnet / high / reviewer |
| TP-C4 Coverage | C3 changed beats embed | held | 0.8 | finding | Sonnet / high / reviewer |
| TP-C4 Coverage | C4 empty diff header | held | 0.9 | finding | Sonnet / high / reviewer |
| TP-C4 Coverage | C5 own path in XTestImports | held | 0.8 | finding | Sonnet / high / reviewer |
| TP-C4 confirming Coverage | C6 embedded sub-package not selected (inherited) | held | 0.7 | finding | Sonnet / high / reviewer |
| TP-C4 third confirming Standards | S3 provenance map duplicates the occurrence | held | 0.7 | finding | Sonnet / high / reviewer |
| TP-C4 third confirming Standards | S4 Sources names the wrong spec line | held | 0.7 | finding | Sonnet / high / reviewer |
| TP-C4 third confirming Standards | S5 FT290 map claims an occurrence record | held | 0.6 | finding | Sonnet / high / reviewer |

Brier mean: 0.169 over 34 pairs; 0 abstained.

## Coordinator catches

- The TP-C1a checkpoint reported a stale plan digest, because the chunk record came before the amendment record.
- The first authors recorded the probe exit code as 0. The checkpoint needs the mutated-run exit, which is 1.
- A named check graded the installed binary's source. The coordinator charged a worktree build and `./dist/bench` for each named check.
- One charge named an `--omit-file` probe that empties a file. The coordinator changed it to a path swap before verification.
- The first landing refused the folded `main` paths. The coordinator proved a kit defect, and a light-path fix landed at `4ca00c25`.
- The second landing went red, because two decision maps cited `roadmap/FT290.md`, and the closure deletes that file.
- The coordinator first deferred the inherited C6 defect. By reviewer direction, `--auto-approve` fixes it at once, at `dc26a362`.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-split-named-check-owner.md | 0 | none |
| 10-explain-changed-selection.md | 4 | delegate-error, other, check-gap, one-source |
| 2-count-tests-run.md | 1 | one-source |
| 3-prove-named-check-ran.md | 1 | delegate-error |
| 4-print-prose-check-result.md | 1 | delegate-error |
| 5-show-each-failure-diagnostic.md | 1 | delegate-error |
| 6-filter-system-suite.md | 1 | one-source |
| 7-name-running-executable.md | 1 | spec-row |
| 8-list-check-fixtures.md | 3 | one-source, delegate-error, delegate-error |
| 9-list-check-inventory.md | 1 | spec-row |

## Agent-experience improvements

### Bench CLI

- Make `bench probe` print the mutated-run exit, and make `bench record verification` read it, as the census entry for 97 raw calls proposes.
  Feeds: new
- Make the completion checkpoint grade the tree after the delivery closure, so a reference to a closed roadmap row reds before the landing.
  Feeds: new
- Make a named `--check` grade the calling worktree's source, or name the graded source in its result.
  Feeds: new

### Skills

- Put the amendment-before-chunk rule, the probe exit rule, and the named-check build rule into the author charge template in `craft-delegate`.
  Feeds: none
- Document `--reviewer`, `--consultant`, and `--auto-approve` as reusable `bench-implement-spec` options.
  Feeds: none

### Process

- Sweep each decision map for references to a roadmap row that the delivery closes, before the completion record.
  Feeds: none
- Give the landing refusal for folded `main` paths a cause that names the kit, not the build.
  Feeds: new