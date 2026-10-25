# Research: a reduced schedule for landings that change no Go file

Recommendation: derive the reader table from a traced whole run, because the
trace found every in-process reader with no hand list. Close the child-process
gap separately, because the trace cannot see it.

Scope: the kit gate, the landing, the conformance registry, the record classes,
and the Go toolchain facts for a census.

Evidence status: verified by reads and measurements. Each fact below cites a
path and line in the tree at
`a1528ae2509ac6e01fde77b8d612c7e97fc40989`, a Go 1.25.14 toolchain source, or a
measurement run on 2026-10-25.

Consumed by: the `no-go-schedule` decision map and the spec that compiles it.
Drift: a change to a cited path, a new landing mix, or a new Go toolchain.
Retire when: the spec that compiles this map lands.

## Facts

### Q1. How often does a landing change no Go file, and which paths change?

A landing is a first-parent merge on `main`. From 2026-10-18 to 2026-10-25,
`main` received 71 landings. Of these, 48 changed no `.go` file.

Of the 48, 30 changed only record paths: `specs/`, `roadmap/`, `ROADMAP.md`,
`capture/`, `reviews/`, `decisions/`, `CHANGELOG.md`, and
`.bench/commitment.json`. The other 18 also changed guidance under
`.agents/commands/`, `.agents/skills/`, `.bench/BENCH.md`, or `docs/`.

Method: `git log --first-parent --merges --since=2026-10-18 main`, then
`git diff --name-only <merge>^1 <merge>` for each merge. Invalidation trigger: a
new landing.

### Q2. Where does the whole gate spend its time?

The landing gate of 2026-10-25T18:39Z ran its six phases one after another, in
about 268 seconds. The `test` phase took 195 seconds, `system` took 47 seconds,
`race` took 8 seconds, and `vet` took 1.5 seconds. Source:
`.logs/gate-20261025T183906.460364287Z-465407.jsonl`, a local, ignored log.

The same run's stream lists 117 test packages with a sum of 669 package-seconds.
The five slowest packages are `internal/testreport` (74 s), `internal/worktree`
(73 s), `internal/gate` (63 s), `internal/probe` (50 s), and
`internal/systemtest` (46 s). Source:
`.logs/gate-20261025T183906.460364287Z-465407.out`.

The kit phase table is `gofmt`, `vet`, `test`, `race`, `system`, and
`shellcheck` (`projects/benchkit.md:340-349`).

### Q3. What does the conformance entry cost, and what do its inputs bind?

The registry holds 46 rows: 2 meta rows and 44 ordinary checks
(`internal/conformance/registry/checks.go:4-51`). Of the 44, 22 bind the
`catch-all` input, and one of the 22 is the ship-tier
`release-evidence-probe` (`checks.go:20`). `catch-all` binds a check to the
complete subject tree (`internal/conformance/registry/registry.go:80-81`).

The ordinary `test` phase runs the complete dev-tier check set through one
entry, `TestRootConformance` (`.bench/BENCH-reference.md:543-547`,
`internal/conformance/gate_entry_test.go:16`).

On 2026-10-25, a run of `TestRootConformance` alone passed in 21.2 seconds of
test time. The wall time, with the compile, was 34.7 seconds on a 16-core host.
The timing file listed 45 executed checks.

### Q4. Which reduced schedule and record exist today?

The delta set is the one reduced schedule. It runs `gofmt`, `vet`, the changed
test packages, and the root-conformance entry (`internal/gate/delta_checkpoint.go:213-223`).
Only a chunk checkpoint with a test-only change since a whole-green reviewed tip
takes it (`delta_checkpoint.go:59-87`).

A delta green writes the `delta record` class. Its reuse refusal stops every
reuse except by a later chunk checkpoint of the same base and subject
(`internal/gate/verdict_registry.go:45-50`, `delta_checkpoint.go:176-181`). No
landing accepts it (`.bench/BENCH-reference.md:647-650`).

The commit lane already maps paths to checks through the path-class table
(`internal/gate/lane_select.go:137-204`). A path that no class claims is
`unknown` and selects every declared check (`lane_select.go:114-116`). The lane
record authorizes no landing (`CONTEXT.md:25-28`).

The profile states that the gate has no component partition and that the delta
set is the one reduced schedule (`projects/benchkit.md:393-397`). It also states
that a reuse needs a current exact green (`projects/benchkit.md:476-479`).

### Q5. What does the green marker authorize?

The gate advances `refs/bench/green/<branch>` to the tip it graded green
(`CONTEXT.md:435-438`). The milestone verification reads the marker as the
revision whose tree last held green (`internal/commitment/repository/verification.go:37-47`).
The landing requires the marker to be an ancestor of its destination
(`internal/worktree/land_identity.go:136-145`,
`internal/gate/authorization/authorization.go:80-123`). The project-green
validation requires the marker at the tip and an exact green of its tree
(`internal/gate/engine.go:68-88`).

### Q6. Which mechanisms can find the readers of a path?

The Go test cache reuses a result only while the files that the test opened in
its module stay unchanged. `-count=1` turns the cache off
(`$GOROOT/src/cmd/go/internal/test/test.go:135`).

The test binary records each in-process open, stat, and environment read when
it gets `-test.testlogfile`. The flag is "for use only by cmd/go"
(`$GOROOT/src/testing/testing.go:479`). `cmd/go` passes it only when the cache
is on (`$GOROOT/src/cmd/go/internal/test/test.go:1546-1549`). The gate runs
`-count=1` (`projects/benchkit.md:346-348`), so a gate run records no test log.

The test log does not see a read that a child process makes. Some tests run a
child in the live tree: for example, `cmd/bench/commands_brief_test.go:68-69`
runs `go build` with the repository root as its directory.

Tests reach the live tree by several routes: a relative `../..` path, `git
rev-parse --show-toplevel`, and package helpers. One search found 27 test
packages that use these routes. The search cannot find a route that a helper
hides, so a static scan does not bound the readers.

The conformance package already refuses an unregistered live-tree test through
an AST scan of its own test files (`internal/conformance/tier_test.go:232-259`).
That scan covers one package.

The profile states that the gate has no per-package loop
(`projects/benchkit.md:393`).

### Q7. Which terms does the glossary already hold?

`census` is the hook-recorded count of raw calls (`CONTEXT.md:323-327`).
`census entry` is an exported function of the worktree package
(`CONTEXT.md:66-68`). The reader sweep entry keeps the word "census" for the
raw-call census and the architecture census (`CONTEXT.md:297-300`). `delta set`
and `delta record` name the one reduced schedule and its record today.

## Tested results

### T1. Which test packages read a record or guidance path in process?

On 2026-10-25, one traced pass ran `go test -trimpath -count=1 ./...` in the
shaping worktree with an `-exec` wrapper. The wrapper gave each test binary its
own `-test.testlogfile`. The pass wrote 112 package logs.

Three packages failed under the wrapper: `internal/preflight/staleness`,
`internal/subprocess/protection`, and `internal/worktree`. The wrapper did not
carry the gate environment. The failed packages still wrote their logs.

A read counts when it opens or stats a non-Go path of the worktree outside the
package's own directory. With this rule, 26 packages read the live tree.

| family | reader packages | gate package-seconds |
|---|---|---|
| record | 9 | 191 of 669 |
| record or guidance | 16 | 232 of 669 |

The 9 record readers are `cmd/bench`, `internal/conformance`,
`internal/consumers`, `internal/dashboard`, `internal/freshness`,
`internal/intent`, `internal/runbinary`, `internal/status`, and
`internal/testreport`. The guidance readers add `internal/adopt`,
`internal/adopt/repairtest`, `internal/anchors`, `internal/chargeevidence`,
`internal/commit`, `internal/prose`, and `internal/tickets`.

Five record readers each opened 143 files under `specs/`, so they read the whole
tree or a copy of it. The slowest reader is `internal/testreport`, at 74 seconds
in the gate run of Q2.

The package-seconds come from the gate stream of Q2. The wrapper, the parser,
and the raw output are scratch files of the shaping session.

### T2. How long are the runs of consecutive no-Go landings?

In landing order from 2026-10-18, the 48 no-Go landings form 16 runs. Their
lengths are 1 (four runs), 2 (five runs), 3 (four runs), 4, 7, and 11.

A bound of five reduced landings since the last whole green would reduce 46 of
the 48. A bound of three would reduce 44.

### T3. Which packages read a covered path through a child process?

This host has no `strace` and no `bpftrace`, so no tool traced a child. A
denial probe ran instead, in a disposable second worktree at the same commit.
It removed the read permission from the 539 tracked files of the covered
families. It then ran the traced pass of T1 and the system suite through
`bench test --check system`.

The test-phase pass failed in 14 packages. One of them is outside the
in-process reader list: `internal/releaseevidence`. Its test
`TestBuildReleaseEvidenceIncludesRegisteredPackageEvidence` runs
`scripts/build-release-evidence.mjs` through `node`, and the script reads
`CHANGELOG.md` (`scripts/build-release-evidence.mjs:65`).

The three packages that failed in T1 failed the same tests in the probe. The
probe therefore shows no added reader in them.

The system suite failed 28 of 274 tests with every covered file unreadable.
The failing journeys include `bench setup --yes`, `link`, and the stripped
distribution. With only the record files unreadable, the suite failed 3 of 279
tests. Each of the 3 starts at `newCandidateInstall`
(`internal/systemtest/candidate_install_test.go:52`). Its kit copy failed on
`.bench/commitment.json` in process (`internal/systemtest/owner_test.go:442`).

The probe ended with every permission restored and a clean status, and the
probe worktree was released.

## Inferences

- A no-Go landing leaves every Go source, `go.mod`, and `go.sum` unchanged. Only
  a test that reads a changed non-Go path can change its outcome, if the
  environment is the same. The lane class `go-build-input` already sends an
  embed target to the Go phases (`lane_select.go:146-151`).
- The complete conformance entry costs less than one fifth of the whole gate.
  If the reduced schedule runs the whole entry, the 22 catch-all checks need no
  per-check census.
- Guidance paths feed the shipped surface and the adoption journey of the
  system suite (`projects/benchkit.md:399-404`). A guidance census therefore
  names more packages, and the system phase, than a record census.

- A reader that walks the tree reads paths that did not exist at the last whole
  run. A table at path-family grain covers a new path in a known family.
- The reader set of a covered landing runs about 232 package-seconds of the
  `test` phase, plus `internal/releaseevidence` and the `system` phase. The
  whole `test` phase runs 669 package-seconds. On the 16-core host, the reader
  set takes about 130 seconds, and the whole gate takes about 268 seconds.
- A denial probe shows only reads whose failure changes an outcome. A reader
  that tolerates a failed read stays hidden from it.

## Proposals

- Derive the reader table from the test logs of each whole run, at path-family
  grain. If a log does not parse, derive no table, and run the whole gate.
- Name the child-process readers through a separate decision, because no test
  log records them.

## Contradictions

- The reviewer reopened the partition decision for no-Go landings on 2026-10-25
  (`roadmap/FT449.md:7-11`). These statements of the old rule stand until the
  spec amends them, as ticket #14 decides:
  - `projects/benchkit.md:351-353`, `:393-397`, and `:476-479`
  - `.bench/BENCH-reference.md:548`
  - `docs/adr/0017-the-worktree-commit-runs-the-fast-lane.md:17-20` and `:30`
- The `oracle` digest of a full verdict frames the tree
  (`internal/gate/subject.go:73`), so it cannot compare two landings. Ticket #8
  defines the oracle identity without the tree.

## Not verified

- A child-process read that tolerates a failed read. No tracer was available.
- The in-process reads of the tagged system suite at path grain. The traced
  pass did not set `-tags=system`, so T3 shows only that the suite reads both
  families.
- The reads of the three packages that failed under the wrapper, after their
  failure point.
- The cost of the reduced schedule on a runner with two to four cores.
- The test-time spread of the whole gate over more than three runs.

## Validation plan

1. Compare a traced pass with the declared census before the spec locks.
2. In the spec, plant a reader of a record path in a test package outside the
   census, and prove that the census check turns red.
3. After the build lands, compare the reduced-landing time with the whole-gate
   time over ten landings.

## Verification record

- [x] The output opens with the recommendation, the scope, and the evidence status.
- [x] Facts, inferences, tested results, and proposals are in separate sections.
- [x] Each question has one section.
- [x] A table compares the reader sets of the two family scopes.
- [x] Every contradiction and every residual unknown is kept.
- [ ] A diagram: no relation needs one.
- [x] The output states what it could not verify.
- [x] The output ends with a validation plan.
- [x] Every load-bearing claim cites a path and line or a dated measurement.
