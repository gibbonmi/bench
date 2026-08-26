# Review pickup — bounded-gate-output (FT185)

Base `81e3415f750a127b745d857c078f4fbe5304fa0d`, tip `55ddad223e91e3341c747fe46801faa162108ffd`.
Initial review, three axes in parallel.

## Standards

Count: 7. Worst: `sanitize.Strip` re-derives `toon.Representable`'s
control-byte threshold instead of composing it. The two rules can now drift.
Only `Representable` is pinned to the encoder's own test.

1. `internal/sanitize/sanitize.go:72` and `internal/toon/toon.go:91` derive
   the same "byte the encoder refuses" fact twice. Only `Representable` is
   pinned to the encoder's own test. — **ask-user** (a fix composes across
   the `internal/toon` boundary, outside this spec's ownership fence).
2. Four test comments narrate the diff or cite "Ticket 05" instead of
   describing the code for its next reader (`internal/gate/report_test.go:315`,
   `:364`, `:657`; `internal/gate/run_stream_test.go:102`). — **auto-fix**.
3. `internal/gate/phase_stream.go:24`'s `stderr io.Writer` field is stored,
   never read, and carries a comment describing a duty nothing performs. —
   **auto-fix**.
4. `gateLogRetainedRecords` (`internal/gate/run_log.go:25`) now counts runs,
   not records; its name and a test comment (`run_log_prune_test.go:355`)
   still say "records". — **auto-fix**.
5. `internal/testlines/testlines.go:23,26` re-list four prefixes `RunnerLine`
   (`:7`) already owns as a duplicated literal list. — **auto-fix**.
6. `projects/benchkit.md:403` and `.bench/BENCH-reference.md:277` carry
   near-identical paragraphs. Story 32 and BG27 name both files, so this is
   spec-sanctioned. — **no-op**.
7. A manifest phase named `capability` collides with the reserved filing
   name. `internal/gate/report.go:27` asserts this is impossible, but
   `manifest.go` reserves nothing. — folded into Coverage finding 5 below.

## Spec

Count: 6. Acceptance rows genuinely unmet: BG01. Worst: `StartErr` (a bad
`Dir`, empty argv, an exec failure, or a scheduler deadlock) is now silently
unreachable. Every such red renders the same uninformative
`exit 1 with no output` row the old `phaseSummary` never used for these cases.

1. **BG01 violated.** The row requires "`failures[N]{phase,line}` and
   `gate: red` on stdout **and nothing else**." `reportCapabilitySkips`
   (`internal/gate/capability_skips.go:186`) unconditionally prints the
   `capability-skips: N (...)` totals line, on every run, red or green. So a
   one-red-phase run prints that line before the table. The spec's own
   "green shape" ties this line to the green branch only. — **auto-fix**:
   gate the totals-line print to the green branch.
2. `phaseResult.StartErr` (set at `runner.go:271,320,329,374`) is never read
   by the report. Every such red collapses to `exit 1 with no output`. —
   **ask-user** (no acceptance row pins the exact row text; needs a design
   choice).
3. Edge inventory: "A `FAIL` line alone on its line is a row" reads against
   `testlines.FailureRows`. A bare `FAIL` yields no row there
   (`testlines.go:28`, pinned by its own test). This plausibly means the
   terminal `FAIL <package>` line, already covered by BG05. — **no-op**
   (ambiguous, not contradicted on the more likely reading).
4. BG24's declared seam is "in-package phases command with a killed
   fixture." Its only test is a buffer-level unit
   (`run_stream_test.go:42`). — folded into Coverage findings 3–4 below.
5. Capability rows are uncapped (`report.go:151` appends `capabilityRows`
   whole). The twenty-row cap is scoped to phase rows only. This was a
   flagged, accepted judgment call from ticket 03. — **ask-user** (re-surfaced
   once, not a new defect).
6. `phases: N/N green`'s numerator excludes skipped phases by design
   (`greenPhases`, `report.go:135`). BG17's all-green case still passes. —
   **no-op** (already-accepted judgment call from ticket 04).

## Coverage

Count: 6. Worst: the kit's own `race` phase loses its entire race report. Go
emits `WARNING: DATA RACE` and the stack before the `--- FAIL:` line. The
classifier's block-open rule never captures those lines, so the two rows it
does produce name no race and no stream file.

1. **A real `go test -race` failure prints no race evidence.** Verified with
   a captured live race stream.

   `internal/testlines/testlines.go:23-24` opens a block only at
   `--- FAIL:`/`# `. Every DATA RACE line arriving first is dropped.

   The fallback in `report.go:215` never fires, because two rows are already
   produced. BG32's premise (a stream of only DATA RACE lines) is not a
   shape `go test` emits. — **ask-user** (a classifier rule is a spec-level
   change).

2. `StartErr` unreachable — same finding as Spec 2, same disposition.

3. The one production line wiring the retained stream
   (`internal/gate/runner.go:148`, `streams.retain(gateRunStreamFile(ctx))`)
   has no test that drives it end to end and then reads
   `.logs/gate-<run>.out`. Every existing test retains a hand-built file or
   context instead. — **auto-fix**: add an in-package killed-fixture test
   (this is also BG24's declared seam, from Spec finding 4).

4. `TestUnwritableLogDirLeavesTheTableBoundedAndTheStreamUnavailable`
   (`run_stream_test.go:148`) chmods `.logs` before `beginGateRunLog`. So the
   `.jsonl` open fails too, and the run log never starts. It does not
   exercise the case BG26 is meant to prove: the `.jsonl` opens, but the
   `.out` open is refused. — **auto-fix**: pre-create the `.out` path so only
   its `O_EXCL` open fails.

5. A manifest phase named `capability` merges its rows with the skip
   reporter's rows, with only its own half capped. `manifest.go`'s
   `validPhaseName` reserves nothing. — **auto-fix**: refuse a manifest phase
   named `capability`.

6. `+<k> more lines: <path>` counts classified rows dropped, not stream
   lines, for a Go test phase. The two coincide for every existing example
   (a `vet` phase). — **no-op** (wording only).

## Disposition summary

**auto-fix** (8, de-duplicated):

- BG01's misplaced capability-skips line.
- The `capability`-named-phase collision, in the manifest and the report.
- The missing end-to-end `.out`-retention test. This closes BG24's declared
  seam.
- BG26's fixture ordering.
- Four comment-register cleanups.
- The unused `phaseStreams.stderr` field.
- The `gateLogRetainedRecords` naming.
- The `testlines` duplicated prefix list.

**ask-user** (4, de-duplicated):

- The `sanitize.Strip` and `toon.Representable` duplication. A fix crosses
  the ownership fence.
- The `StartErr` diagnosis loss. This needs a row-text decision.
- The race-phase classifier gap. This is a spec-level change tied to BG32's
  premise.
- The uncapped capability rows. This is a re-surfaced accepted call.

**no-op** (6): docs duplication (spec-sanctioned), bare-`FAIL` ambiguity,
`phases: N/N green` numerator (accepted), more-lines wording, and two
already-accepted ticket judgment calls.
