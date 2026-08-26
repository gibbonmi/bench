# Bounded gate output

Status: staged

Decision source: reviewer-confirmed conversation, 2026-08-26. The reviewer asked for failure rows only on red and at most ten rows on green. Two late questions closed: filter the plain phase stream, and print the phase table with the verdict.

Verification log: 2 iterations to accept — round one returned eight blocking findings. The second iteration folded all twenty-six, and a scoped re-review accepted.

## Problem

A gate run prints every line of every phase. That is seventy `ok` package
lines, the race phase's `=== RUN` chatter, the system phase, six summaries,
and the verdict. A landing or a commit relays the same stream. A red run buries
its one failing check in more than a hundred lines, and a green run costs the
same read for no information. The AXI principles the kit teaches its own
queries do not hold for the three verbs an agent runs most.

## Solution

The gate engine buffers each phase's stream instead of relaying it. On green
it prints one `phases[N]{phase,verdict,elapsed_ms}` table, one
`capability-skips` line, and `gate: green`. On red it prints one
`failures[N]{phase,line}` table that holds only failure rows, then
`gate: red`, on stdout. A Go test phase's failure rows are its `--- FAIL:`
blocks, its `# <package>` build-error blocks, its `FAIL <package>` lines, and
its `panic:` lines. Every other phase's failure rows are every non-empty line
of its red stream. Each phase's rows cap at twenty, and one more row names
the complete stream.

The complete stream goes to `.logs/gate-<run>.out` beside the progress log,
under the same retention in runs. `bench commit` and `bench worktree land`
relay the engine's output, so they inherit the shape and keep their own
records. The fast lane keeps its own relay.

## User stories

### A red gate prints failure rows only

Line: opus / high.

The projection sits inside the oracle's reporting path, where a wrong filter
hides a red. The cached routing for foundational Go-seam work applies.

1. As an agent, I want a red gate's stdout to hold one `failures[N]{phase,line}` table and `gate: red`, so that I read the failure and nothing else.
2. As an agent, I want a green phase's stream printed nowhere on stdout when another phase is red, so that the table names only failures.
3. As an agent, I want a red Go test phase's rows to be its `--- FAIL:` lines, so that each failing test is one row.
4. As an agent, I want a FAIL block's diagnostic lines as rows under that phase, so that the reason travels with the name.
5. As an agent, I want a Go test phase's `FAIL <package>` and `panic:` lines as rows, so that a build failure or a panic shows.
6. As an agent, I want every non-empty line of a red gofmt, vet, or shellcheck stream as a row, so that findings are rows.
7. As an agent, I want each phase's rows capped at twenty plus one `+<k> more lines: <path>` row, so that a wide red stays bounded.
8. As an agent, I want an empty red stream to print one row `exit <code> with no output`, so that a silent red shows.
9. As an agent, I want a phase skipped because a need went red to add no row, so that consequences do not read as failures.
10. As an agent, I want every red diagnosis of the skip reporter as a row under phase `capability`, so that an unmade assertion shows.
11. As an agent, I want `gate: red` on stdout, so that the verdict and the rows share one stream.
12. As an agent, I want a cancelled run to keep its stragglers on stderr and print no verdict, so that an interrupt is not graded.
13. As an agent, I want the exit codes unchanged, so that every caller's contract holds.
14. As an agent, I want a control byte in a failure line removed before the row renders, so that the emitter never refuses the table.
15. As an agent, I want a stream's last line without a newline flushed as one row, so that a truncated tool output still shows.
16. As an agent, I want a `# <package>` build-error block's lines as rows, so that a compile error names its file and line.
17. As an agent, I want a red Go phase with no classified row to print its last twenty lines, so that a race report shows.
18. As an agent, I want the fast lane's relay unchanged, so that a worktree commit still prints its check's own lines.

### A green gate prints at most ten rows

Line: opus / medium.

The green shape composes the phase results the engine already holds and the
existing TOON emitter, so the cached mid routing applies.

19. As an agent, I want a green gate to print `phases[N]{phase,verdict,elapsed_ms}`, one `capability-skips` line, and `gate: green`, so that six phases fit in nine rows.
20. As an agent, I want a green six-phase run to print exactly nine stdout lines and no phase-stream line, so that the bound holds.
21. As an agent, I want a phase table above seven phases collapsed to one `phases: N/N green` row, so that the count never passes ten.
22. As an agent, I want a skipped optional phase to show `skipped` in its verdict cell, so that an absent tool is not a red.
23. As an agent, I want the capability classes folded into the one `capability-skips` line, so that three lines become one.
24. As an agent, I want `gate: green (fresh verdict reused for this tree)` unchanged, so that a reused verdict still reads as one line.
25. As an agent, I want `bench commit` and `bench worktree land` to relay the engine's stdout byte for byte with their records, so that all three read alike.
26. As an agent, I want the `elapsed_ms` cell taken from the phase timing the progress log records, so that the table and the log agree.

### The complete stream is retained

Line: opus / medium.

The retention composes the run log's existing naming and pruning, so the
cached mid routing applies.

27. As a reviewer, I want every phase line written to `.logs/gate-<run>.out` as it arrives, so that a killed run keeps its stream.
28. As a reviewer, I want twenty runs retained, each with its `.jsonl` and its `.out`, so that `.logs` stays readable.
29. As a reviewer, I want an unwritable `.logs` to leave the projection bounded with the row `+<k> more lines (stream unavailable)`, so that a logging failure never unbounds stdout.
30. As a reviewer, I want the engine to name the stream file once on stderr, so that a short red points at the whole.
31. As a reviewer, I want the tests that pin `gate: green` and the reused line to stay green, so that the contract holds.

### The shape is documented

Line: opus / high.

The profile and the reference are guidance prose, so the leverage override
routes them mid and high.

32. As a reviewer, I want the profile and the reference to describe the bounded shape and the `.out` file, so that a cold session knows.

## Implementation decisions

**The engine gets its own phase writer; the lane keeps its relay.** The fast
lane composes `prefixedPhaseWriters` today and never reaches the report, so
that function stays as it is. The engine gains `newPhaseStreams`, a phase
stream buffer. Each phase's stdout and stderr lines arrive in order and go to
the `.out` file at once. When a phase settles, the projection reads the
buffer: green prints nothing, red prints failure rows. The buffer holds
lines, not bytes, so a partial last line flushes as one line at close.

**One classifier owns Go test lines.** `internal/testreport` classifies
runner lines today. That predicate moves to a new low package,
`internal/testlines`, with one more: the failure rows of a red stream. A
`--- FAIL:` line at any indent starts a block. A `# <package>` line starts a
build-error block. A `WARNING: DATA RACE` line also starts a block. Each
block ends at the next runner line.

A race report's stack precedes its `--- FAIL:` block, and the new opener
turns the whole thing into rows instead of dropped lines (BG38). `FAIL
<package>` and `panic:` lines are rows on their own. When the classifier
yields no row for a red Go phase, the phase's last twenty non-empty lines
are the rows. `testreport` composes the moved predicate, so one source names
a runner line. A phase is a Go test phase when its argv starts with `go
test`; every other phase's red stream is failure rows line by line.

**The red shape.** The engine prints `failures[N]{phase,line}` through
`toon.Table`, in phase-table order, then `gate: red` on stdout. A red phase
with no lines prints `exit <code> with no output`. A phase with more than
twenty rows prints twenty, then `+<k> more lines: <path>`. The cap is twenty
so a one-phase red fits the stop hook's thirty-line tail.

Every red diagnosis the skip reporter produces is a row under phase
`capability`: an environment skip, a strict-mode capability skip, and an
unreadable skip log. `capability` is a reserved filing name; a manifest that
declares a phase with that name is refused before the run starts (BG39). A
phase skipped by a red need adds no row. Every cell
passes a control-byte filter that removes bytes below 0x20 except tab and
escapes nothing, so a backslash reaches the reader once. A long line renders
whole. The cap counts rows, not bytes.

A red phase can carry a `StartErr`: a bad working directory, an empty argv,
an exec failure, or a scheduler deadlock. Such a phase renders that error's
text as its row, instead of the no-output text (BG37).

**The green shape.** The engine prints `phases[N]{phase,verdict,elapsed_ms}`
in phase-table order, with `green` or `skipped` in the verdict cell and the
elapsed time the progress log records. Then one `capability-skips` line
carries the total and every class: `capability-skips: 6 (capability=6
environment=0; fifo=3 privilege=3)`. Then `gate: green`. A table above seven
phases collapses to `phases: N/N green`. The reused-verdict line is unchanged.

**The stream file.** The run log owner writes `gate-<run>.out` beside
`gate-<run>.jsonl`. The name round-trips through the same recognizer. The
prune counts runs, not files: it keeps the twenty newest runs and removes
each older run's pair together.

When the file opens, the engine prints
`gate: stream <path>` once on stderr. The stderr line `gate: progress log
<path>` stays as it is. When the file cannot open, the phase buffer still
bounds stdout and the more-row says `(stream unavailable)`.

**The verbs.** `bench commit` and `bench worktree land` run the gate through
the engine and relay its writers, so they change nothing. Their journeys
prove the relay with a fixture gate that prints a canned bounded shape. The
engine's own rows attach in-package at `phasesCommandAtKitWithSelection`
with a constructed run-binary selection. The exec path demands a sealed run
binary no unit test holds.

## Testing decisions

The highest seam that shows each failure is the engine's report: phase
results with captured streams in, stdout out. `aggregateAndReport` and the
phase buffer take fake results and fake streams in a unit test. The
in-package call `phasesCommandAtKitWithSelection` with fixture phase scripts
proves the composition and the exit codes.

The outer `bench gate` journey
with fixture scripts (`run_outcomes_test.go` prior art) proves the reused
line. The commit and landing journeys prove the relay with a canned fixture
gate. The testlines unit proves the classifier on captured `go test` output.
The gate's `test` phase observes all of it.

### Seam diagram

    trigger: bench gate / bench commit / bench worktree land
        │
        ▼
    phase streams ──▶ [ phase buffer: lines → .logs/gate-<run>.out ] ──▶ per-phase line buffer
                                │ phase settles
                                ▼
                     [ projection: green → nothing; red → testlines.FailureRows ]
                                │
                                ▼
                     failures[N]{phase,line} + gate: red   |   phases[N]{...} + capability-skips + gate: green
                                ◀ tests attach here: fake phase results + captured streams in, stdout out

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| BG01 | 1 | A run with one red phase prints `failures[N]{phase,line}` and `gate: red` on stdout and nothing else. | engine report unit | A relayed stream adds lines. |
| BG02 | 2 | A run with one green and one red phase prints no line from the green phase. | engine report unit | A filter keyed on the run, not the phase, prints both. |
| BG03 | 3 | A red `go test` stream with two `--- FAIL:` lines yields two rows with those lines. | testlines unit | A filter that drops runner lines drops the names. |
| BG04 | 4 | A FAIL block with two diagnostic lines yields two rows after the name row. | testlines unit | A name-only filter loses the reason. |
| BG05 | 5 | A stream with `FAIL\tgithub.com/x/y [build failed]` and `panic: boom` yields both as rows. | testlines unit | A block-only filter misses a build failure. |
| BG06 | 6 | A red `vet` phase with three non-empty lines yields three rows. | engine report unit | A Go-only filter prints nothing for vet. |
| BG07 | 7 | A red phase with fifty failure lines prints twenty rows and one `+30 more lines: <path>` row. | engine report unit | An uncapped table is unbounded. |
| BG08 | 8 | A red phase with an empty stream prints one row `exit 1 with no output`. | engine report unit | A silent red prints no row. |
| BG09 | 9 | A phase skipped because its need went red adds no row. | engine report unit | A skip row reads as a failure. |
| BG10 | 10 | An environment skip prints a row under phase `capability`. | engine report unit | A skip outside the table is lost. |
| BG11 | 11 | `gate: red` appears on stdout and not on stderr. | engine report unit | The old stderr line splits the answer. |
| BG12 | 12 | A cancelled run prints stragglers on stderr and no verdict on either stream. | engine report unit | A verdict grades an interrupted run. |
| BG13 | 13 | A red run of fixture phase scripts exits 1. | in-package phases command with fixture scripts | A changed code breaks every caller. |
| BG14 | 14 | A failure line with an ESC byte renders with the byte removed. | engine report unit | An unsanitized cell makes the emitter refuse, so the whole table is lost. |
| BG15 | 15 | A red stream whose last line lacks a newline yields that line as a row. | phase buffer unit | A newline-keyed buffer drops it. |
| BG16 | 19 | A green run of six phases prints `phases[6]{phase,verdict,elapsed_ms}`, one `capability-skips` line, and `gate: green`. | engine report unit | An omitted table leaves the old six lines. |
| BG17 | 21 | A green run of eight phases prints `phases: 8/8 green` and no table. | engine report unit | An uncollapsed table passes ten rows. |
| BG18 | 22 | A skipped optional phase shows `skipped` in its verdict cell and the run stays green. | engine report unit | A skip read as red fails the run. |
| BG19 | 23 | Three capability skips in two classes print one `capability-skips` line that names both classes. | engine report unit | Three lines pass the cap. |
| BG20 | 24 | A second run on an unchanged tree prints `gate: green (fresh verdict reused for this tree)` alone. | outer runner journey | A changed reuse line breaks its pin. |
| BG21 | 25 | `bench commit` on a fixture gate that prints a canned red table relays that table, `gate: red`, and its own refusal record. | commit journey with a canned fixture gate | A relay that filters loses the table or the record. |
| BG22 | 25 | `bench worktree land` on a fixture gate that prints a canned green shape relays it before `landed{...}`. | landing journey with a canned fixture gate | A relay that filters loses the table. |
| BG23 | 26 | A phase's `elapsed_ms` cell equals the `phase.finish` record's `elapsed_ms` in the progress log. | in-package phases command reading the log | A second timer disagrees with the log. |
| BG24 | 27 | A run killed after its first phase leaves that phase's lines in `.logs/gate-<run>.out`. | in-package phases command with a killed fixture | A buffer written at settle loses a killed run's lines. |
| BG25 | 28 | A twenty-first run leaves twenty `.jsonl` files and twenty `.out` files and removes the oldest pair. | run log prune unit | A prune that counts files keeps ten runs. |
| BG26 | 29 | An unwritable `.logs` leaves the red table at twenty rows plus `+<k> more lines (stream unavailable)`. | engine report unit with a refused stream file | A failed open unbounds the relay. |
| BG27 | 32 | The profile's gate section and the reference each carry the bounded-output sentence needle. | anchors registry test | A dropped sentence leaves the shape undocumented. |
| BG28 | 31 | The run-outcome test that pins the reused line and the adoption test that pins `gate: green` pass unchanged. | existing tests | A pin on the old shape reds the gate. |
| BG29 | 18 | A lane run with a red prose check prints that check's lines as it does today. | lane unit | A lane routed through the buffer prints nothing. |
| BG30 | 20 | A green run of six phases prints exactly nine stdout lines and no line with a `[phase]` prefix. | engine report unit | A relaying build adds the table and passes BG16. |
| BG31 | 16 | A stream with `# github.com/x/y` and `./x.go:12:3: undefined: y` yields both as rows. | testlines unit | A FAIL-block-only classifier drops a compile error. |
| BG32 | 17 | A red Go phase whose stream holds no `--- FAIL:`, `# `, `FAIL`, `panic:`, or `WARNING: DATA RACE` line prints its last twenty non-empty lines. | engine report unit | A classifier with no fallback prints an empty table. |
| BG33 | 10 | A strict-mode capability skip and an unreadable skip log each print a row under phase `capability`. | engine report unit | A diagnosis left on stderr leaves an empty red table. |
| BG34 | 13 | A green run of fixture phase scripts exits 0. | in-package phases command with fixture scripts | A changed code breaks every caller. |
| BG36 | 13 | A cancelled run of fixture phase scripts exits 130. | in-package phases command with fixture scripts | A changed code breaks every caller. |
| BG35 | 30 | A run whose stream file opens prints `gate: stream <path>` once on stderr. | in-package phases command reading stderr | A short red never names the whole stream. |
| BG37 | 1 | A phase red with a `StartErr` (a bad `Dir`, an empty argv, an exec failure, or a scheduler deadlock) renders that error's text as its row. | engine report unit | The generic no-output row hides the real reason. |
| BG38 | 17 | A red Go phase whose stream carries a `WARNING: DATA RACE` block followed by a `--- FAIL:` block yields rows for both. | testlines unit | A `--- FAIL:`-only opener drops the race report before it. |
| BG39 | 10 | A manifest phase named `capability` is refused before the run starts. | manifest validation unit | An unreserved name lets a phase's rows merge with the skip reporter's. |

Review repair (2026-08-26): the initial review found the red table printed the
`capability-skips` totals line on every run, not only on green, so BG01's
"and nothing else" did not hold. The fix moves that print to the green
branch; BG01 already covers the corrected behavior and gets no new row.

### Edge inventory

- A `--- FAIL:` line nested under a parent test (indented) starts a block like a top-level one.
- A `--- PASS:` or `=== RUN` line inside a red stream is a runner line and ends a block; it never becomes a row.
- A `FAIL` line alone on its line is a row; `ok` and `?` lines never are.
- A phase whose stderr carries the failure and whose stdout is empty yields its stderr lines as rows.
- Two phases red at once print their rows in phase-table order, never interleaved.
- A red phase with exactly twenty lines prints twenty rows and no more-row.
- A stream line longer than `sanitize.Preview`'s bound renders whole; the rows never call the preview.
- A run whose `.logs` directory is a symlink writes nothing there and reports the stream unavailable.
- A capability-class skip outside strict mode never enters the failure table; an environment skip always does.
- A gate script that does not exec `bench gate-phases` (the scaffolded gate, a `$BENCH_GATE` command) relays its own lines unchanged. The bounded shape is the engine's.
- An interrupted run discards its buffer; the `.out` holds what arrived when the run log exists, and nothing otherwise.
- The scaffolded gate keeps `gate: red` on stderr, so a linked project without the phase-table engine reads the verdict on the other stream.

**Won't handle** a `-json` test phase — the phase table pins the plain argv, and the filter reads plain output.

**Won't handle** the canary's `EXPECT` matching — the canary invokes each check's owner directly and reads no gate stdout.

**Won't handle** a `--full` flag that prints the stream on stdout — the `.out` file is the complete value, named once per run.

**Won't handle** the fast lane's shape — its four checks print little, and it never reaches the report.

## Ownership fences

- `internal/gate/`
- `internal/testlines/`
- `internal/testreport/`
- `internal/sanitize/`
- `internal/commit/`
- `internal/landing/`
- `internal/worktree/`
- `internal/systemtest/`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/toon/toon.go` (review-approved addition, 2026-08-26: `sanitize.Strip`
  composes the encoder's own control-byte predicate instead of re-deriving it)
- `projects/benchkit.md`
- `.bench/BENCH-reference.md`
- `CHANGELOG.md`
- `specs/bounded-gate-output/`

The classifier ticket lands first. The buffer and the red shape follow, then
the bound and the green shape, then the stream file, then the verb journeys
and the docs.

## Out of scope

- Go test phases on `-json` with `bench test`'s decoder: 12 edits, 3 gate runs, a phase-table change.
- A `--full` flag that relays the stream to stdout: 3 edits, 1 gate run.
- AXI approved-query status for the gate with a `help[N]` envelope: 4 edits, 1 gate run; the gate is a mutation, not a query.
- A bounded projection for `bench shift`'s iteration output: 6 edits, 2 gate runs.
- The scaffolded gate script's own `gate: green` and `gate: red` lines (`bench init`): 2 edits, 1 gate run; its `gate: red` stays on stderr, unlike the engine's.
- The two pre-phase refusals that print `gate: red` on stderr before any phase runs: 2 edits, 1 gate run; they stand for veto.
- A stop hook tail above thirty lines: 1 edit, 1 gate run.

## Further notes

The reviewer's two answers on 2026-08-26 are the decision source's late
questions: filter the plain stream rather than switch the argv, and print the
phase table with the verdict rather than one aggregate row. Four choices are the author's and stand for veto. They are the twenty-row
cap per phase, the collapse above seven phases, the stream file with its
twenty-run retention, and the `gate: stream <path>` line.

The initial review (2026-08-26) closed four judgment calls. Capability rows
stay uncapped, since they carry no backing stream file to point a reader at.
`sanitize.Strip` composes `toon`'s own control-byte predicate rather than
re-deriving it, which is the fence addition above. A `StartErr` red renders
that error's text (BG37). A `WARNING: DATA RACE` line opens a block the same
way `--- FAIL:` and `# ` do (BG38).
