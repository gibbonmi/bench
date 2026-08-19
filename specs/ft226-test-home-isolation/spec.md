# ft226-test-home-isolation

Status: staged

Decision source: roadmap/FT226.md (reviewed drain commit cd355f45, 2026-08-19) with the audit evidence it cites — docs/audits/2026-08-bench-capability/results-fable-high/reconciliation-ledger.md entry L-06 and action-items.yaml item A6.

Verification log: 2 iteration(s) to accept — one independent `opus` high-effort read-only round; iteration 1 returned REVISE on two blocking findings (the fence omitted `capture/session-handoff.md`, which every phase close writes; the residue unit of report was unpinned and contradicted ticket 02), folded together with six prose findings (SIGINT helper child disposed, IS2/IS3 predicates tightened, SW3 compares sorted listings, EV1 timing, ticket 02 `Writes:`); the round's seam and ticket-graph verdicts were agreement. Build probes (a)–(c) and the sweep counts are appended here by the build.

## Problem

The kit's own tests write into the operator's real `BENCH_HOME`. `Pool(root)`
resolves the worktree pool under `$BENCH_HOME/worktrees/<key>`, and the
`internal/worktree` reauthorize tests build their fixture through
`reauthorizeFixture` → `mustCreate` → `Create` without binding a temporary
`BENCH_HOME`, unlike every sibling fixture in the package. Each fixture repo is
a `t.TempDir` named `001`, so every run adds ten `001-<digits>` pool keys to the
operator's `~/.bench/worktrees` and never removes them: 759 orphans at audit
time, 1,690 (63 MB) when this spec was written, roughly ten per gate run. Under
a subject-pointing `BENCH_HOME` the same leak appeared as pollution of an
external subject (audit E006).

Measured for this spec: `go test -count=1 ./...` with `BENCH_HOME` pointing at
an empty sentinel directory passes and leaves exactly ten `worktrees/001-*`
entries in the sentinel, all attributable to `TestReauthorize*` through the
leaked worktrees' `.git` pointers. No other package writes there.

## Solution

Two in-package changes and one operator sweep.

`reauthorizeFixture` binds `BENCH_HOME` under its test's temporary root for the
test's duration, exactly as its sibling fixture `newOwnedAssignment` does, so
its pool worktrees live and die with the test.

The `internal/worktree` test binary gains a `TestMain` that binds `BENCH_HOME`
to a process-private directory created empty before the first test and removed
after the last, and that turns the package red when anything remains under that
directory after the run, printing each residue path and the `gitdir:` target of
every leaked worktree so the message names the originating test. This makes the
operator's pool structurally unreachable from any driver of the package — the
gate's `test` and `race` phases, `bench test`, or a hand `go test` — and makes a
future fixture that forgets to bind fail loudly instead of accumulating
silently. The red is proven once by mutation in the build.

The existing orphans are swept once from the operator's real pool as part of the
same build, plan-before-apply, under a predicate that cannot touch a live pool:
a key named `001-<digits>` whose every child worktree carries a `.git` pointer
file with a dangling `gitdir:` target.

## User stories

### Isolation and detection

Line: `opus` / medium. Exact spec, known shape, one package; the detector is
itself oracle logic whose wrong implementation would sail through, so the
cheap row's bump applies.

1. As an operator, I want a gate run to leave `$HOME/.bench/worktrees` exactly
   as it found it, so kit development stops filling my machine with fixture
   worktrees.
2. As an operator, I want every driver of the worktree package tests — the
   gate's `test` and `race` phases, `bench test`, a hand `go test` — to be
   unable to reach my real pool, so the protection does not depend on which
   driver runs.
3. As a maintainer, I want `reauthorizeFixture` to bind a temporary `BENCH_HOME`
   under its test's temporary root like its sibling fixtures, so its pool
   worktrees are removed with the test's `t.TempDir`.
4. As a maintainer, I want the worktree test binary to run under a
   process-private `BENCH_HOME` created empty before the first test and removed
   after the last, so a fixture that forgets to bind still cannot touch the
   operator's pool.
5. As a maintainer, I want the package to exit red when anything remains under
   that process-private home after the tests, naming each residue path and the
   originating test's temporary root, so a future leak fails loudly.
6. As a maintainer, I want that red proven by mutation — removing the fixture
   binding turns the package red naming a `TestReauthorize…` root, restoring it
   turns the package green — so the detector is known to bite.
7. As a maintainer, I want the fallback test that empties `BENCH_HOME` to keep
   computing `$HOME/.bench/...` as a pure path without writing, so `HOME`
   isolation stays out of scope (reviewed exclusion).
8. As the reviewer, I want the build to record a full ordinary-suite run under
   an empty sentinel `BENCH_HOME` that leaves the sentinel empty, so the claim
   covers kit tests as a whole and not only the one package with a retained
   oracle.

### One-off sweep

Line: `opus` / medium. Exact predicate, known shape, but the target is the
operator's machine outside the tree where no gate observes a mistake — the
uncovered-gate bump applies and the reviewer's sign-off on this spec is the
standing approval for the destructive step.

9. As an operator, I want the orphaned `001-<digits>` pool keys already in my
   real pool removed once, so the accumulated directories and megabytes go away.
10. As an operator, I want the sweep to remove only keys whose every child
    worktree's `.git` pointer is dangling, leaving every other key (`bench-…`,
    `project with spaces`, the dogfood keys) untouched, so it cannot destroy a
    real worktree.
11. As an operator, I want the sweep to print its target count, a sample of
    targets with their dangling `gitdir:` lines, and the surviving keys before
    applying, and the remaining count after, so the destructive step is
    plan-before-apply and leaves evidence.

## Implementation decisions

- **Fixture binding follows the package convention.** `reauthorizeFixture`
  gains the same per-test binding its sibling `newOwnedAssignment` uses:
  `BENCH_HOME` set under the fixture's own root before the first `Create`.
  Collapsing the package's ~60 explicit per-test bindings into the shared
  `newWorktreeRepo` is a separate refactor (Out of scope) — several tests build
  two fixture repos in one scope and rebinding on the second call would move
  the first repo's pool out from under the test.
- **One `TestMain`, one residue predicate.** `internal/worktree/main_test.go`
  owns the process-private home: create it empty with `os.MkdirTemp`, record
  the inherited value, export the private path, run, compute residue, remove,
  exit. Creation failure is fail-closed: the binary exits non-zero naming the
  error before any test runs rather than running against the inherited home.
  The residue predicate is one unexported function over a directory. Its unit
  of report is the top-level entry: one line per top-level entry under the
  private home (a file, or a directory such as `worktrees/` whether empty or
  not), plus, for each leaked pool worktree found beneath it that holds a
  `.git` pointer file, that worktree's `gitdir:` line. It is unit-tested on a
  synthetic home, and `TestMain` composes it rather than re-deriving it.
- **Exit code combines.** The binary exits non-zero when `m.Run` is non-zero, when
  residue exists, or when the private home cannot be removed; a clean green run
  exits zero. The residue report goes to stderr with one heading line stating
  the count and the private home path, then one line per residue entry.
- **Private home location is the OS temp directory.** Same lifetime class as
  `t.TempDir`; an interrupted run leaves it there like any interrupted test's
  temp root, never under the operator's `BENCH_HOME`.
- **No production change.** `Pool`, `benchHome`, the wrapper's `BENCH_HOME`
  export, and the gate's phase environment are untouched. Non-goal from A6.
- **Sweep is a one-off, plan-before-apply, outside the tree.** The build runs a
  throwaway script from the scratchpad against `$HOME/.bench/worktrees`: a
  target is a key matching `^001-[0-9]+$` whose every child directory holds a
  regular `.git` file with a `gitdir:` target that does not exist, and which
  holds nothing else at its top level. The plan prints before any removal; the
  apply step removes only listed targets. Counts and the surviving key names
  are recorded in this spec's verification log and the landing evidence. No
  reusable verb is added (Out of scope).

## Testing decisions

- The external behavior a good test exercises: a package run, whatever driver
  launches it, leaves the operator's pool listing byte-identical; a test that
  creates a pool worktree under a home it did not bind itself turns the package
  red with a message that names it.
- Seams: `TestMain` in `internal/worktree` (new; prior art
  `internal/systemtest/owner_test.go`'s owner `TestMain` with its post-run
  verify and cleanup), the unexported residue predicate (unit seam), and the
  reauthorize fixture (prior art `newOwnedAssignment` in `resume_test.go`).
- Gate seam: the ordinary `test` phase (`go test -count=1 ./...`) and the
  `race` phase both run the worktree test binary and therefore its `TestMain`;
  no gate or shell edit. The hostile-shell checklist adds no new shell boundary;
  the walk below records the classes that apply.
- Mutation probes recorded in the verification log, not retained: (a) remove
  the fixture binding → `go test -count=1 ./internal/worktree -run Reauthorize`
  exits non-zero with residue naming a `TestReauthorize…` temporary root;
  restore → green. (b) `GOTMPDIR=/tmp TMPDIR=/nonexistent go test
  -count=1 ./internal/worktree -run TestPool` → exits non-zero with the
  `TestMain` creation error and no test output. `GOTMPDIR` keeps the go driver's
  own work directory valid; without it the driver fails before the binary
  compiles and the probe never reaches `TestMain` (reviewer-approved correction,
  build finding P1). (c) a temporary `t.Fatal` in any test with
  a clean home → package FAIL (exit code combination).

### Seam diagram

    driver: gate test/race phase · bench test · go test
        │
        ▼
    [ TestMain: MkdirTemp → export BENCH_HOME → m.Run → residue → RemoveAll → exit ]
        │                                          ◀ unit test: residue predicate on a synthetic home
        ▼                                          ◀ probe: fixture binding removed → red names the test
    tests: fixture binds BENCH_HOME under its own t.TempDir ──▶ Create ──▶ pool under the test root
           (reauthorizeFixture joins newOwnedAssignment's convention)

    operator sweep (once): plan(targets under $HOME/.bench/worktrees) ──▶ sample ──▶ apply ──▶ counts

### Acceptance coverage map
| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| IS1 | 3, 6 | `reauthorizeFixture` binds `BENCH_HOME` under the test's temporary root before creating its worktree, so the created path is under that home | mutation probe (a): binding removed → `go test ./internal/worktree -run Reauthorize` exits non-zero with residue naming a `TestReauthorize…` root; restored → green | a fixture still creating under the process home is exactly the residue the probe's red names; a detector that never reds leaves the probe green |
| IS2 | 1, 2, 4 | at test start `BENCH_HOME` is the directory `TestMain` created: it exists, is not the inherited value, is not under it, is not under `$HOME/.bench`, and `Pool` of a fresh root resolves under it | `TestPackageHomeIsPrivate` | a `TestMain` that forgets to export, exports after `m.Run`, exports the inherited value, or creates its home beneath the operator's shows the operator's path to the test |
| IS3 | 4 | the private home starts empty and no longer exists after the run | `TestPackageHomeIsPrivate` asserts `os.ReadDir` of the exported home returns no entries at test start, plus verification-log evidence that the path is gone after `go test -count=1 ./internal/worktree` | a `TestMain` that reuses a fixed path or skips removal leaves the path behind between runs |
| DT1 | 5 | any top-level entry left under the private home after `m.Run` is residue: the package exits non-zero and prints one line per top-level entry, including a stray non-pool file and an empty `worktrees/` directory | `TestHomeResidueListsLeakedWorktreesWithOrigin` on a synthetic home planting a `worktrees/001-1/<id>/.git` pointer and a stray top-level file (two entries), a home holding only an empty `worktrees/` directory (one entry), and an empty home (none) | a predicate that only scans `worktrees/`, only counts pool keys, or walks for files and skips directories misses one of the planted shapes |
| DT2 | 5 | each residue worktree's `.git` pointer target is printed so the report names the originating test's temporary root | `TestHomeResidueListsLeakedWorktreesWithOrigin` asserts the planted `gitdir:` target appears in the report | a report of bare pool-key names leaves the offender unnamed and the next author guessing |
| DT3 | 5 | the binary's exit is non-zero when `m.Run` fails, when residue exists, or when removal fails, and zero only when all three are clean | probes (a) and (c) plus the ordinary green run | a `TestMain` that exits with the residue verdict alone masks a failing test; one that exits with `m.Run` alone masks the leak |
| DT4 | 4 | when the private home cannot be created the binary exits non-zero naming the error before any test runs | probe (b) | ignoring the `MkdirTemp` error runs the whole package against the inherited home — the operator's pool |
| EX1 | 7 | `TestPoolDefaultBenchHome` still computes a `$HOME/.bench/worktrees/…` path with `BENCH_HOME` empty and writes nothing | existing test unchanged; `t.Setenv` restores the private home afterwards | binding `HOME` too would change what this test proves; the row pins that the package does not |
| EV1 | 8 | `go test -count=1 ./...` with `BENCH_HOME` set to an empty sentinel directory passes and leaves the sentinel with no entries | verification-log evidence (before/after `find` of the sentinel), captured at ticket 01 before the `TestMain` lands — a post-02 re-run covers the remaining packages only, since the worktree binary no longer consumes the sentinel | the only kit-wide oracle for "kit tests" is this run; the gate-level private home is priced in Out of scope |
| SW1 | 9, 10 | a sweep target is a key matching `^001-[0-9]+$` whose every child directory holds a regular `.git` file with a dangling `gitdir:` target and nothing else at top level, and every other key survives | plan listing in the verification log: target count, non-target count, surviving non-`001` key names | a predicate on the name alone would take a live `001` pool; a predicate on dangling pointers alone would take a foreign key the operator still uses |
| SW2 | 11 | the plan (count, five sampled targets with their `gitdir:` lines, survivors) prints before any removal, and the post-apply remaining count equals the pre-apply count minus targets | verification log | an apply-first script leaves no evidence that the predicate selected what it claimed |
| SW3 | 1, 9 | after the sweep, one subsequent gate run leaves the sorted `$HOME/.bench/worktrees` key listing identical | verification log (sorted `ls` listing before and after `bench gate`, compared byte for byte) | the operator-visible acceptance from A6: the leak is closed under the real driver, not only under a sentinel |

### Edge inventory

- **Error path:** private home creation failure → exit non-zero before tests
  (DT4). Removal failure → reported and non-zero (DT3).
- **Empty/absent input:** an empty private home after the run is the green
  case; a home with only an empty `worktrees/` directory is residue (DT1 —
  nothing legitimate materializes it).
- **Boundary values:** exit-code combination over the three verdicts (DT3).
- **Malformed input:** a residue worktree without a readable `.git` pointer is
  still listed by path (DT1); its origin line is simply absent.
- **Interrupted / partial state:** an interrupted run leaves the private home
  in the OS temp directory, the same lifetime class as an interrupted test's
  `t.TempDir`; never under the operator's pool.
- **Re-run idempotency:** each process creates a fresh private home; `-count=n`
  and `-run` filters change nothing.
- **Process-boundary lifecycle:** `go test` grades the binary's exit code;
  `TestMain` must combine, not replace (DT3, probe c).
- **Hostile environment:** a `TMPDIR` with spaces or glob characters gives the
  private home such a path; only a test that failed to bind its own home ever
  resolves a pool there, and that test is the red. Tests that bind hostile
  `BENCH_HOME` values (spaces, ESC, tab) keep doing so under their own roots.
  `t.Setenv` restores the private home after each of them.
- **Sweep hostile state:** a key whose child `.git` is a directory (a full
  repository rather than a pointer), a key with an extra top-level entry, a key
  whose pointer target exists, and a key not matching `001-<digits>` all
  survive (SW1).
- **Won't handle:** binding `HOME` for the package — `TestPoolDefaultBenchHome`
  needs the real fallback, git reads the operator's global config, and no
  worktree test writes under `$HOME/.bench` with `BENCH_HOME` empty today; a
  future one is a review catch (EX1).
- **Won't handle:** short-circuiting `TestMain` for the SIGINT helper child that
  `ownership_test.go` re-execs from the test binary with `BENCH_SIGINT_HELPER=1`
  — the child receives its own private home, and the cleanup transaction it
  drives reads `root` and `path` and never `BENCH_HOME`; ticket 02 runs
  `-run TestActualSIGINT` under the new `TestMain`, and a red there takes the
  `systemtest` owner's helper short-circuit as the fix.
- **Won't handle:** a test that deliberately binds `BENCH_HOME` inside the
  private home — none exists; it would read as residue and the red names it.
- **Won't handle:** residue-report escaping beyond `%q` — the report is stderr
  from a test binary read by a person or the gate log, not a TOON sink.
- **Won't handle:** concurrent operator `bench worktree` activity during a
  test run — the package never reads the operator's pool, so there is nothing
  to race with.

## Ownership fences

One writer at a time; the three tickets are serial. Reviewer disposition:
approve, merge, or split at sign-off.

- `internal/worktree/reauthorize_test.go`
- `internal/worktree/main_test.go`
- `specs/ft226-test-home-isolation/`
- `capture/session-handoff.md`

The sweep ticket writes nothing in the tree except its ticket checkboxes and
this spec's verification log; its target is `$HOME/.bench/worktrees`.

## Out of scope

- **Gate-level private `BENCH_HOME`** for the `test` and `race` phases with a
  residue red across every package: ~5 edits (`internal/gate/phases.go`,
  `internal/gate/runner.go`, their tests, `projects/benchkit.md` gate section),
  2 gate runs. The sentinel run (EV1) shows one package writes today; this spec
  guards that package under every driver, which the gate-level guard would not.
- **A `bench worktree` verb that sweeps dangling pool keys** under
  `$BENCH_HOME/worktrees`: ~6 edits (command, predicate, TOON rows, tests,
  help, profile), 2 gate runs. The one-off sweep here is operator hygiene, not
  a capability.
- **Collapsing the package's explicit per-test `BENCH_HOME` bindings into
  `newWorktreeRepo`**: ~20 test-file edits, 1–2 gate runs, with the
  two-fixture-repo tests needing a stable per-test home first.
- **`HOME` isolation** for the package (see Won't handle).

## Build verification log

**EV1 — the ordinary suite under an empty sentinel `BENCH_HOME`.** Before ticket
01, `BENCH_HOME=<sentinel> go test -count=1 ./internal/worktree` passed and left
exactly ten `worktrees/001-<digits>` entries in the sentinel. After ticket 01,
`BENCH_HOME=<sentinel> go test -count=1 ./...` passed and left the sentinel with
zero entries (`find <sentinel> -mindepth 1` empty). Captured at ticket 01, before
the `TestMain` landed, so the sentinel was the only oracle in play.

**Probe (a) — the fixture binding removed.** Run twice: once by the write
delegate, once independently by the coordinator. With the `t.Setenv` line deleted
from `reauthorizeFixture`, `go test -count=1 ./internal/worktree -run Reauthorize`
exits 1 and prints one residue entry (the top-level `worktrees` directory under
the private home) with ten `gitdir:` origin lines, each naming a
`TestReauthorizeCommand…` temporary root. Restoring the line returns exit 0 with a
clean tree. The detector bites.

**Probe (b) — the private home cannot be created.** The command the spec named,
`TMPDIR=/nonexistent go test -count=1 ./internal/worktree -run TestPool`, does not
exercise `TestMain`: the go driver fails creating its own work directory before
the test binary compiles, so the exit is 1 but the message is the toolchain's.
The form that does exercise it keeps the toolchain's work directory valid:
`GOTMPDIR=/tmp TMPDIR=/nonexistent go test -count=1 ./internal/worktree -run TestPool`
exits 1 with `private BENCH_HOME: stat /nonexistent: no such file or directory`
and no test output, as does `TMPDIR=/nonexistent ./worktree.test -test.run TestPool`
against a prebuilt binary. DT4 holds; the probe command in Testing decisions is
what the finding corrects, and the correction is the reviewer's to veto.

**Probe (c) — exit-code combination.** A temporary `t.Fatal` in `TestPool` under a
clean private home makes the package FAIL with exit 1; removed, the package
passes with exit 0. The residue verdict does not mask a failing test.

**IS3 lifetime.** After `go test -count=1 ./internal/worktree`,
`ls -d /tmp/bench-worktree-home-*` finds no directory: the private home is removed.

**Driver runs, all exit 0.** `go test -count=1 ./internal/worktree`;
`go test -count=1 -v ./internal/worktree -run TestPoolDefaultBenchHome` (EX1,
unchanged and passing); `go test -count=1 ./internal/worktree -run TestActualSIGINT`
(the re-exec'd helper child takes its own private home, so the short-circuit the
spec priced as a contingency was not needed); `go test -race -count=1 -v
./internal/worktree -run TestConcurrentCleanupRecordsOneTransaction`. Both ticket
commits ran the full six-phase gate green.

**SW1 predicate, proven before the real run.** A scratch pool planted all four
hostile shapes: a key whose child `.git` is a directory, a key with an extra
top-level entry, a key whose `gitdir:` target exists, and a key named `001-abc`.
All four survived the plan, along with the `bench-…` and `project with spaces`
keys; only the fully-dangling `001-<digits>` key was targeted.

**SW2 plan and apply.** Plan against `$HOME/.bench/worktrees`: 1,719 keys total,
1,710 targets, 9 non-targets, 91 MB. Every sampled target carried one child with a
dangling `gitdir:` naming a `TestReauthorizeCommand…` temporary root. No `001-*`
key survived the predicate. The nine survivors, printed before the apply, were
`bench-2826441890` (this build's own integration worktree),
`bench-dogfood-20260710-459PYf-2003174707`,
`bench-dogfood-gl-axi-20260712-3579930476`,
`bench-dogfood-gl-axi-20260712-r2-722800923`,
`bench-luna-dogfood-mx8bDj-2174529179`,
`bench-mandatory-delegation-dogfood-1057428268`,
`c65757ef284c706c00622d7e181844b8-e4f58d5d86f77bdf927e119a15eb421c-1941573793`,
`ft87-refresh-probe.US84JO-3910988020`, and `project with spaces-2317837872`.
After the apply the script reported 9 keys remaining — 1,719 minus 1,710 — and
every named survivor was present. The pool fell to 51 MB. The reviewer ran the
destructive step; the auto-mode classifier denied it to the agent.

**Ten keys arrived after the apply, and they confirm the diagnosis.** A concurrent
session landed two commits on `main` (`e1b44e62`, `a9e8e232`) during this build.
`main` does not carry ticket 01, so those two gate runs leaked ten fresh
`001-<digits>` keys with `TestReauthorizeCommand…` origins — the same signature,
from the only tree still able to produce it. They remain in the pool and are swept
by a second pass once this spec lands.

**SW3 — one gate run leaves the pool listing identical.** Sorted `ls` of
`$HOME/.bench/worktrees` before and after `bench gate` in the integration worktree,
which carries both tickets: identical, 19 keys, no diff. The leak is closed under
the real driver, not only under a sentinel. That gate reported red solely because
`BENCH_CONFORMANCE_ROOT` was unset for `TestRootConformance`, an absent-environment
verdict; every test phase ran green, and both ticket commits gated green through
`bench commit`, which stages that environment.

## Further notes

The residue report is the guidance for the next fixture author: it names the
leaked path and the `gitdir:` target that carries the test name, so the fix
(bind `BENCH_HOME` under the fixture root) is legible from the red alone. No
profile note is added; a second statement of the rule would drift from the
`TestMain` comment that owns it.
