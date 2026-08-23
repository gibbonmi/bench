# Worktree fresh-test latency

Status: ready

## Destination

Restore the kit's fresh test run to a sub-minute developer feedback loop without
weakening `-count=1` or changing what green means. Reduce `internal/worktree`'s duplicated executable builds and real-Git
materialization. Replace process-global environment and working-directory inputs
with explicit values. Separate the landing, lifecycle, and pool/reclaim
responsibilities so Go can schedule them as independent packages. Introduce
measured bounded parallelism only after those global inputs are gone. Retain a
regression budget that makes renewed package growth visible before it becomes
the gate's long pole.

The whole-tree gate remains the oracle. This map decides delivery slicing, the
latency target, real-Git journey fidelity, parallelism posture, and what the
regression budget proves. Spec authoring will place exact test and implementation
seams after those choices close.

## #1: What workload and closed decisions constrain the repair?

Blocked by: none
Type: Research

### Question

Measure the current fresh test run, and distinguish compile cost from test cost.
Trace the package's process/global-state owners against the gate's settled
single-build and scheduling decisions. Identify whether WSL resource placement,
Go compilation, one hung test, or cumulative package architecture is the long
pole.

### Answer

Resolved 2026-08-23 on `0e17d428`. A warm compile-only
`internal/worktree` run is 0.32 s. One landing subtest is 2.11 s, and its complete
top-level journey is 5.29 s. The whole package is 137.86 s wall / 130.39 s
reported by Go.

The whole fresh test run remained unfinished when stopped at
163.62 s. This is cumulative process and filesystem demand, not compilation or
one hung test. The 2026-08-13 closed baseline was 31.9 s for the entire test
phase.

The package has 286 top-level tests and no `t.Parallel`. Its `TestMain` replaces
process-global `BENCH_HOME`; helpers and callers use `t.Setenv` and `os.Chdir`.
`newWorktreeRepo` initializes and commits a real repository for broad
equivalence partitions. `buildLandingBinary` invokes `scripts/go-build.sh` at
12 call sites even though `internal/runbinary` and the gate already define one
selected Bench executable per top-level run. The directory has 51 source files;
`land_test.go` is 2,133 lines.

The WSL repository, build cache, and module cache are on Linux ext4 under
`/home`, and the host showed no material CPU or memory pressure. The closed gate
architecture leaves package scheduling to Go's `-p`, runs one ordinary
`go test -count=1 ./...`, and has no Bench-owned test fan-out. This work must not
reintroduce an outer scheduler, clear the Go build cache, scope the gate by diff,
or remove `-count=1`.

## #2: Does this ship as one spec or as a measured two-spec sequence?

Blocked by: #1
Type: Grill

### Question

The deterministic demand reduction is independently useful. Select one Bench
binary, make environment and working directory explicit, split responsibility
owners, and move equivalence partitions below real Git. Parallelism and a stable
regression budget depend on the post-reduction workload. Should one spec bind all
of it, or should the second outcome be priced from a new baseline after the first
lands?

### Answer

Two sequential specs (reviewer, 2026-08-23). The first removes deterministic
demand and the process-global inputs that prevent safe scheduling. It delivers
one selected Bench binary per run, responsibility owners, and typed policy
partitions below representative real-Git journeys. The second measures that landed shape, then
introduces only the bounded parallelism it still needs and attaches the
regression budget to the reduced workload. The second spec must not price a
worker limit or time threshold from the architecture the first spec removes.

## #7: What is invoked where and when during the current fresh test run?

Blocked by: #1
Type: Research

### Question

Produce an invocation census for `internal/worktree` and the whole fresh test
run. Attribute top-level test spans and every observed child-process class. The
classes include nested `go` or Bench invocations, `scripts/go-build.sh`, Git
commands, shells, sleeps/waits, and any gate or successor test process. Distinguish useful
CPU from process/filesystem wait and identify light contention or serialized
idle windows. Compare that topology with the 2026-08-13 31.9-second baseline before you
recommend a new latency envelope. That baseline included a named 30-second
publication wait.

### Answer

Resolved 2026-08-23. The invocation census is
`decisions/assets/worktree-test-invocation-census.md`. There is no nested
`go test`, recursive whole gate, or one fixed wait inside `internal/worktree`.
The topology has 12 repeated selected-binary builds. They expand to 12
`go build` and 12 seal-producing `go list` calls. Those build-bearing tests
start 55 private Bench commands.

There are 123 static `newWorktreeRepo` call
sites and hundreds of Git helper references whose table loops multiply
execution.

The package has 106 `t.Setenv`, 15 `t.Chdir`, two direct `os.Chdir`, and zero
`t.Parallel` references. In the timing sample, 38 tests at or above one second
contributed 99.45 seconds and the slowest single test was 12.27 seconds. The
sample itself overlapped another Go verification and a moving checkout. Its
167.41-second wall time is contention evidence, not a baseline. The independent
clean gate observed `internal/worktree` at 136.319 seconds on the pre-FT113
source. FT113 changed only 13 production lines in `worktree.go`, leaving the
measured test topology intact.

The historical 31.9-second whole test phase was a different floor:
`internal/publication` spent 30.03 seconds waiting for a kernel connection retry
against `127.0.0.1:1`. That uncontrolled wait remains present and must not be
priced as acceptable worktree latency. The worktree repair therefore targets
serialized process/filesystem churn; the publication wait is a separate bug.

## #8: Is publication's uncontrolled 30-second wait a prerequisite bug?

Blocked by: #7
Type: Grill

### Question

The whole-suite latency target cannot honestly credit the kernel's connection
retry to worktree architecture. Should the existing publication test become a
separate `$bench-debug` repair? That repair would replace its background-context
and default-client request before either worktree spec prices a whole-suite
budget. The alternative is for this map to absorb that unrelated transport
behavior.

### Answer

Yes: treat it as a separate prerequisite bug (reviewer, 2026-08-23).
Repair publication's uncontrolled background-context/default-client request
through `$bench-debug` before either worktree spec attaches a whole-suite
timing budget. It is not part of the worktree specs, and no worktree change may
hide it or claim its removal as worktree latency improvement.

## #3: What latency envelope counts as restored?

Blocked by: #1, #7, #8
Type: Grill

### Question

Choose the observable target for the reference WSL development host. The target
has three parts: a maximum `internal/worktree` package span, a maximum whole
fresh test run, and the sample rule that distinguishes a regression from
ordinary timing noise. The target must
leave room for gate setup and the other packages rather than making 59 seconds a
nominal success.

### Answer

Use the reference WSL host target recommended in the grill (reviewer,
2026-08-23): across three fresh `go test -count=1` runs with normal
`GOCACHE`/`GOMODCACHE`, an idle host, and no concurrent gates,
`internal/worktree` has a median at or below 20 seconds and no run above 25
seconds. After the separate publication repair lands, the whole suite has a
median at or below 25 seconds. No run may exceed 35 seconds under the same
conditions.

The first spec records before/after demand measurements without pretending its
seam extraction alone must reach the final number. The second spec must meet
both envelopes without weakening `-count=1`.

## #4: Which behaviors still earn a real Git journey?

Blocked by: #1
Type: Grill

### Question

Decide the fidelity rule for the three named responsibility families. Real Git
is the local substitutable dependency at each public command seam. Its use for
every decision-table partition repeats repository setup and process waits.
Which representative success, refusal, interruption, and hostile-state paths
must remain public real-Git journeys? Which partitions may move to typed
in-process facts behind those commands?

### Answer

Retain real Git at each public command seam for behavior Git itself supplies
(reviewer, 2026-08-23). Landing keeps representative publish/release,
conflict/refusal-without-mutation, interrupted/resumed, and Git-dependent
hostile-residue journeys. Lifecycle keeps native create/remove and actual Git
registration and lock behavior. Pool/reclaim keeps journeys where a real
process, lease file, registration, or worktree deletion determines the result.

Each adapter that gathers repository facts gets focused real-Git coverage.
Once those facts are typed, combinatorial ownership, lease, eligibility,
ignored-output, age, and action partitions run in-process through the decision
owner. A new policy partition does not earn another repository solely because
its public caller ultimately uses Git.

## #5: How is parallel demand bounded after global state is removed?

Blocked by: #4
Type: Grill

### Question

Choose whether package-level scheduling alone is the default. Choose whether pure
in-process tests also use `t.Parallel`. Choose whether real Git journeys may run
concurrently under an explicit shared limit. The answer must preserve the closed
machine-wide resource posture and must not make WSL filesystem contention the
new source of flakes.

### Answer

Do not add a scheduler for this workload (reviewer, 2026-08-23). Keep every
public journey that starts Git, Bench, Go, or another descendant process in one
serial journey package. The landing, lifecycle, and pool/reclaim owner packages
may use `t.Parallel` only for typed, in-process policy tests. Those tests start
no descendants and mutate no process-global environment or working directory.

This boundary avoids multiplying nested process demand while letting pure
decision work overlap. Any future proposal to parallelize descendant-spawning
journeys must reopen shaping, re-census the then-current workload, and adopt a
measured shared product budget before introducing width. FT171's historical
worker-by-`GOMAXPROCS` budget is precedent, not an active scheduler available
to `internal/worktree`.

## #6: What regression budget becomes part of the oracle?

Blocked by: #2, #3, #5
Type: Grill

### Question

Choose the fail posture and owner for renewed test demand. The options are a
wall-clock package budget, a deterministic census of expensive constructors and
real-Git journeys, or a combined rule. In the combined rule, structural growth
fails the gate and timing remains measured evidence. The budget must catch architectural drift without making
ordinary host variance produce arbitrary red gates.

### Answer

Use the combined rule (reviewer, 2026-08-23). The ordinary gate hard-fails
structural demand drift. Selected Bench binary construction has one owner and
occurs once per top-level test run. Every test descendant start routes through
the serial journey harness, which contains no `t.Parallel`. Tests outside it
neither mutate process-global environment/CWD nor bypass their typed owner
seams to start descendants. Each enforcement must be independently
mutation-proven red without duplicating its production registry or parser.

Wall time remains required acceptance evidence rather than an ordinary-host
gate threshold. The second spec records three fresh reference-WSL runs and does
not complete unless they meet #3's package and whole-suite envelopes. A miss
reopens workload investigation; it does not authorize raising the limit or
weakening `-count=1`. This combines a deterministic everyday oracle with a
host-controlled slow-package regression budget without turning machine noise
into arbitrary red gates.

## Not yet specified

## Spec-writer discretion

- Exact package, interface, helper, and file names inside each reviewer-approved
  responsibility owner, provided no new public behavior or package cycle appears.
- How the first spec groups and presents its required post-reduction timing
  distribution, provided the raw reference-host measurements remain reproducible.
- Whether a pure owner test parallelizes at its top level or within subtests,
  provided it obeys #5's no-descendant and no-process-global-state boundary.

## Out of scope

- Removing or conditionally bypassing `-count=1`.
- Clearing `GOCACHE` or `GOMODCACHE` as part of the fresh-test measurement.
- Diff-scoped, component-scoped, or cached green verdicts.
- Weakening, deleting, or moving coverage solely to improve wall time.
- Replacing real Git at every public command seam with mocks.
- Reintroducing Bench-owned outer phase or package fan-out.
- Treating the separate Codex/WSL Go bootstrap defect as test-latency work.

## Sources

- Path: `internal/worktree/main_test.go`
  Supports: #1's process-global `BENCH_HOME` finding and #5's parallelism constraints.
  Drift: re-read if package test setup or the residue oracle changes.
- Path: `internal/worktree/worktree_test.go`
  Supports: #1's `newWorktreeRepo` and process-global working-directory findings, and #4's fixture producer.
  Drift: re-read if shared worktree test helpers move or stop materializing real Git repositories.
- Path: `internal/worktree/land_test.go`
  Supports: #1's repeated selected-binary builds and #4's current public landing journeys.
  Drift: re-measure and re-count after any landing-test split or binary-fixture change.
- Path: `internal/worktree/eligibility.go`
  Supports: #4's existing typed decision owner behind repository fact gathering.
  Drift: re-read if eligibility fact gathering or decision ownership moves.
- Path: `internal/gittest/gittest.go`
  Supports: #4's real Git materialization producer and dependency classification.
  Drift: re-read if repository initialization or identity setup changes.
- Path: `internal/runbinary/runbinary.go`
  Supports: #1 and #2's existing one-selected-executable owner.
  Drift: re-read if run binary selection, inheritance, or lifetime changes.
- Path: `projects/benchkit.md`
  Supports: #1, #2, and #5's settled gate shape: one selected executable and Go-owned package scheduling.
  Drift: re-read before spec authoring if the Gate section changes.
- Path: `decisions/gate-budget.md`
  Supports: #1 and #3's 2026-08-13 baseline and #5/#6's closed scheduling and resource posture.
  Drift: #26–#27 are the current closed state; replace their baseline only with a newer exact-subject census.
- Path: `decisions/gate-concurrency.md`
  Supports: #5 and #6's closed product-budget and no-weakening constraints.
  Drift: re-read if the machine-wide budget owner or canary arithmetic changes.
- Path: `decisions/assets/worktree-test-invocation-census.md`
  Supports: #7's process topology, timing attribution, historical comparison, and the factual premises of #3, #5, and #6.
  Drift: re-run the clean timing census after the first demand-reduction spec lands; re-count if worktree test helpers, run-binary ownership, or publication transport changes first.
- URL: `https://pkg.go.dev/cmd/go#hdr-Test_packages`
  Supports: #1 and #5's `-count=1`, package scheduling, and `t.Parallel` semantics; fetched 2026-08-23.
  Drift: mutable Go documentation; re-verify if the repository's declared toolchain changes from Go 1.25.
- URL: `https://learn.microsoft.com/windows/wsl/filesystems`
  Supports: #1's Linux-filesystem placement finding; fetched 2026-08-23.
  Drift: mutable WSL guidance; re-verify before making a new portability claim.
