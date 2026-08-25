# The worktree test floor

Status: staged

Roadmap: FT215

Decision source: ready compiled map `specs/worktree-test-floor/decisions/worktree-test-floor.md`, resolved 2026-08-25

Verification log: 1 iteration to accept — the round accepted the eleven-ticket chain; the closing ticket carries the live-tree census so every ticket lands green. Amended 2026-08-25 after the build measured a 73 s wall. The reviewer dropped the five-slowest pin and reopened the `BENCH_HOME` seam. The verbs take an explicit home, and the `BENCH_HOME`-bound tests join the parallel set.

## Problem

The gate's `test` phase runs packages in parallel, so its wall is the slowest
package plus the build. `internal/worktree` is that package at 69–87 s. Its
suite is fully serial: 334 top-level tests at about 218 ms each, and every one
spawns `git` several times. No test calls `t.Parallel()`. The next slowest
package takes 14 s, so this one package sets the floor of every full run.

## Solution

Every `internal/worktree` test that binds no environment and changes no
directory runs in parallel. A census derives the serial set from call edges to
the harness's `bindEnv` and `chdir` helpers. It turns red when an eligible
test lacks `t.Parallel()`. The harness's shared state stays correct under
parallel siblings, and a race sentinel proves it under `-race`.

Four verbs resolve their root from the working directory today: `clean`,
`reclaim`, `list`, and `resume clean`. They take an explicit root like the
other seven verbs, so their tests need no `chdir`. No test is removed, merged,
or weakened, and the gate's phase table does not change.

Every verb reads `BENCH_HOME` inside the package today, so a test that runs a
verb in-process binds the process environment and stays serial. The verbs
take an explicit home in the shape of the explicit root, and a child a verb
starts carries that home on its environment. The test fixtures pass a private
home, so the `BENCH_HOME`-bound tests join the parallel set.

## User stories

### Eligible tests run in parallel

Line: opus / medium.

The change is a harness refactor with a race hazard. Only `-race` and the
census catch it, so the cached mid routing applies.

1. As a kit maintainer, I want every eligible `internal/worktree` test to call `t.Parallel()`, so that the package wall drops.
2. As a kit maintainer, I want a census that turns red when an eligible test lacks `t.Parallel()`, so that the cut cannot erode.
3. As a kit maintainer, I want the census to derive the serial set from call edges, so that an env-bound test is serial by construction.
4. As a kit maintainer, I want a test that binds environment or changes directory to stay serial, so that Go's rules hold.
5. As a kit maintainer, I want the harness's shared state to stay correct under parallel siblings, so that the race phase stays green.
6. As a kit maintainer, I want a race sentinel that runs parallel journeys under `-race`, so that a harness data race turns the gate red.
7. As a kit maintainer, I want the build's retro to record the package's serial sum and wall, so that the destination is measured.

### Four verbs take an explicit root

Line: opus / medium.

The change composes the root-taking seam seven verbs already use.

8. As a kit maintainer, I want `clean`, `reclaim`, `list`, and `resume clean` to take an explicit root, so that their tests need no `chdir`.
9. As an operator, I want those four verbs to behave as before from any directory inside the repo, so that no grammar or output changes.

### Every verb takes an explicit home

Line: opus / high.

The change hoists a process-environment read to thirteen verb boundaries and
must carry the home into every child the verbs start. A leak reaches the
operator's real home, so the top of the mid tier applies.

14. As a kit maintainer, I want every verb to take an explicit home after its root, so that an in-process test binds no environment.
15. As an operator, I want a child a verb starts to carry that home, so that a private home reaches the broker and the gate.
16. As a kit maintainer, I want the shared fixtures to bind no environment, so that the tests they serve are eligible.
17. As a reviewer, I want a ceiling pin on the serial set, so that a fixture cannot fall back to a process bind.

### What the gate proves does not change

Line: sonnet / low.

Two pins keep the cut honest, and each is an exact static check.

10. As a reviewer, I want the package's test count to stay at least 334, so that no test is removed or merged for wall-clock.
11. As a reviewer, I want the harness effect census to keep forbidding effects outside the harness files, so that parallel tests cannot bypass the harness.
12. As a reviewer, I want the `internal/freshness` and `internal/runbinary` floors left to FT246, so that this spec stays one capability.
13. As a reviewer, I want no sub-package split and no export of internal identifiers, so that the package's surface does not change.

## Implementation decisions

**Eligibility is a static predicate.** A top-level test is serial when its
body, or any test-file function it reaches through call edges, calls `bindEnv`
or `chdir`. Every other top-level test is eligible and must call
`t.Parallel()`. The census parses the package's `_test.go` files with the Go
AST and builds the call graph over test-file functions. It reports each
eligible test that lacks the call.

It also reports each serial test that carries the call. The census
lives beside the existing harness effect census and shares its file walk.

**The harness stays the one effect owner.** The effect log already holds a
mutex. The run-binary selector's journey log gains one, and its `sync.Once`
selection stays. `descendant` keeps its per-test cleanup.

Spawned journeys pass
`Dir` and `Env` on the child. A journey that needs a private `BENCH_HOME` or a
stub `PATH` for the child alone stays eligible. A journey that needs the
process environment for an in-process call binds it and stays serial.

**Four verbs take a root.** `CleanCommand`, `ReclaimCommand`, `ListCommand`,
and `ResumeCleanCommand` gain a leading `root string` parameter, in the shape
of `LandCommand` and `ReleaseCommand`. Their callers resolve the root once at
the command boundary, as the other seven verbs' callers do. The verbs' grammar,
output, and exit codes do not change.

**The race sentinel.** The race-test registry gains one worktree sentinel that
runs two parallel journeys against two repositories and reads the effect log.
The gate's `race` phase runs it.

**Every verb takes a home.** `effects.go` exports `Home()`, the one
`BENCH_HOME` read. Each public verb gains a `home string` parameter after
`root`, and each caller resolves the home once at the command boundary, as it
resolves the root. Some package functions read the home and serve another
package: `Pool`, `Acquire`, `Create`, `ConservativeCleanup`, and their kin.
Each keeps its form for that caller and gains an `At` form that takes the
home. A function this package alone calls takes the home outright. The tests
call the form that takes the home. The verbs' grammar, output, and exit
codes do not change.

**A child carries the home.** A verb starts children: the landing's gate and
broker, a `bench` child, and a git child that reads the pool. Each receives
`BENCH_HOME=<home>` on its environment. The verb sets that entry from the
resolved value, not from the process environment. The harness's home-residue guard turns red when
a child writes below the operator's real home.

**The fixtures stop binding.** `landingFixtureAtHome` writes the tally path
into the gate scripts as a literal and declares no environment in
`gate-inputs.json`, so it binds neither `BENCH_HOME` nor `LAND_GATE_TALLY`.
`newOwnedAssignment`, `newPendingAssignment`, `mustCreate`, `newReclaimPool`,
`reauthorizeFixture`, and `newResidueGuardFixture` pass the home they own. A
test that grades the boundary default itself keeps its bind and stays serial.

**The serial ceiling.** The census reports the serial set's size, and a pin
requires it at or below the count the closing ticket observes. A fixture that
falls back to a process bind raises the count above the pin.

**The count pin.** The census requires at least 334 top-level tests in the
package. A new test raises the count; a removal below the pin turns the gate
red.

**Measurement.** The build's retro records the package's serial sum and wall
before and after, from `go test -count=1 -json`. The destination is a wall
near 20 s, and the retro states the observed number.

## Testing decisions

The highest seam that shows each failure is the census itself: a package test
that walks the test files and reports. The census has unit tests against
synthetic file sets in a temporary directory. A planted omission turns it red
without an edit to the live tree. The harness's shared state is tested
under `-race` through the sentinel. The four verbs are tested by a call with an
explicit root while the working directory is elsewhere. The gate's `test` and
`race` phases observe all of it.

### Seam diagram

    trigger: go test ./internal/worktree (the gate's test phase)
        │
        ▼
    _test.go files ──▶ [ parallel census: AST call graph, bindEnv/chdir edges ] ──▶ eligible set, serial set
                            │ eligible test without t.Parallel()
                            ▼
                          red, naming file:line
                      ◀ tests attach here: synthetic file sets in a temp dir

    trigger: go test -race (the gate's race phase)
        │
        ▼
    two parallel journeys ──▶ [ harness: effect log, run-binary selector ] ──▶ no race report
                                  ◀ tests attach here: the worktree race sentinel

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| WF01 | 1 | The parallel census passes on the live tree. | package census test | An eligible test left serial turns the census red. |
| WF02 | 2 | A synthetic eligible test without `t.Parallel()` makes the census report its file and line. | census unit on a temp dir | A census that only counts calls misses the omission. |
| WF03 | 3 | A synthetic helper that calls `bindEnv` marks a test that calls the helper serial. | census unit on a temp dir | A census that reads only the test body marks the caller eligible. |
| WF04 | 3 | A synthetic test that calls `chdir` inside a subtest closure is serial. | census unit on a temp dir | A census that skips closures marks the test eligible. |
| WF05 | 4 | A synthetic serial test that calls `t.Parallel()` makes the census report the pair. | census unit on a temp dir | Go panics at run time, and the census names it statically first. |
| WF06 | 5 | Two goroutines that call the run-binary selector record two journey lines. | selector unit under `-race` | An unguarded append loses a line or races. |
| WF07 | 6 | The race-test registry names one worktree sentinel that runs two parallel journeys. | racetests registry test | A sentinel outside the registry never runs under `-race`. |
| WF09 | 8 | `CleanCommand` with an explicit root acts on that root while the working directory is a temp dir. | package unit | A verb that still resolves the working directory acts on the wrong repo. |
| WF10 | 8 | `ReclaimCommand`, `ListCommand`, and `ResumeCleanCommand` accept an explicit root in the same way. | package unit | One verb left on the working directory keeps its tests serial. |
| WF11 | 9 | `bench worktree clean --help` and `bench worktree list` from a subdirectory print the same bytes as before. | command registry test | A caller that passes the subdirectory as root changes the output. |
| WF12 | 10 | The census requires at least 334 top-level tests in the package. | package census test | A removal below the pin passes silently. |
| WF13 | 11 | The harness effect census still reports an `exec.Command` outside the harness files. | existing census test | A relaxed regex lets a parallel test spawn outside the harness. |
| WF14 | 13 | Every `internal/worktree` test file keeps `package worktree`. | package census test | An external test package is a sub-package split by another name. |
| WF15 | 14 | `CreateCommand` with an explicit home writes its pool below that home while the process `BENCH_HOME` names a different directory. | package unit | A verb that still reads the environment writes to the wrong home. |
| WF16 | 15 | A landing with an explicit home runs a gate script that records `$BENCH_HOME`, and the record names that home. | landing journey | A child that inherits the process environment grades under the wrong home. |
| WF17 | 16 | The landing fixture binds no environment, and a test it serves calls `t.Parallel()` under the census. | package census test | A fixture that binds keeps every landing test serial. |
| WF18 | 17 | The census reports the serial set's size, and a pin refuses a count above the ceiling. | package census test | A fixture that falls back to a bind passes silently. |

Not covered: story 7 — retro evidence, not gate behavior.
Not covered: story 12 — an exclusion; FT246 owns those packages.

### Edge inventory

- A test file with a build tag is walked like any other, and its tests obey the predicate.
- A helper in a harness file (`journey_test.go`, `main_test.go`, `test_run_test.go`) joins the call graph.
- A test that calls `t.Setenv` directly is already red under the effect census.
- `TestMain` is not a top-level test and is neither eligible nor serial.
- A table-driven test whose subtests call `t.Parallel()` and whose parent binds no env is eligible, and the parent calls `t.Parallel()` too.
- A serial parent with parallel subtests is allowed when no subtest binds env or changes directory.
- Two parallel journeys select the run binary once, through the `sync.Once`.
- A special file in the package directory (a FIFO or a socket) is skipped by the walk, as the effect census skips it today.
- An eligible test that calls `t.Parallel()` after a `bindEnv` in a helper is serial, and the census reports the pair.

**Won't handle** a sub-package split — every test file is `package worktree`, and the internal surface stays.

- A test that grades the boundary default (`Home()` without `BENCH_HOME`) binds the environment and stays serial.
- A `bench` child the landing starts reads `BENCH_HOME` from its own environment, so the verb sets it there from the resolved home.

**Won't handle** conversion of in-process verb calls to child processes — 250 tests would change seam, and the census makes the serial set explicit instead.

**Won't handle** a wall-clock assertion in the gate — timing is retro evidence, and the census is the gate's pin.

## Ownership fences

- `internal/worktree/`
- `internal/racetests/racetests.go`
- `internal/gittest/gittest.go`
- `cmd/bench/main.go`
- `cmd/bench/command_registry_test.go`
- `internal/sessioninspect/sessioninspect.go`
- `internal/harness/worktree.go`
- `internal/shift/`
- `internal/status/status.go`
- `internal/dashboard/dashboard.go`
- `CHANGELOG.md`

The tickets run serially on one source. The root ticket lands first, the
harness and census ticket second, and the parallel marks follow file by file.
The home seam lands after the marks: the verb ticket first, the fixture ticket
second, the direct-bind tickets in parallel, and the ceiling ticket last.

## Out of scope

- The `internal/freshness` and `internal/runbinary` floors: 10 edits, 2 gate runs, owned by FT246.
- Conversion of in-process verb tests to child-process journeys: 250 edits, 6 gate runs.
- A sub-package split with exported test seams: 80 edits, 4 gate runs, and it changes the package surface.
- The gate's phase table and what green means. The cadence spec pins the `test` argv.

## Further notes

The five slowest tests at spec time, from the asset:

- the resume-landing destructive-state refusal, 4.97 s
- the evicted-receipt refusal, 2.80 s
- the lifecycle fault boundaries, 2.61 s
- the public real-Git landing journey, 2.28 s
- the recovery-preserves-layers test, 2.17 s

The cost is count, not one slow test, so the parallel marks pay across the
whole file set.
