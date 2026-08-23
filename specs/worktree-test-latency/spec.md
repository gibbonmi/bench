# worktree-test-latency

Status: staged

Decision source: `specs/worktree-test-latency/decisions/worktree-test-latency.md`, ready on 2026-08-23, with its owned invocation census.

Verification log: 2 iteration(s) to accept. One read-only `gpt-5.6-terra`
high review returned REVISE because three rows lacked an inventory and ledger
owner. The fold added one proof inventory, missing-proof reds, and
`evidence/coverage-ledger.md`. Iteration 2 returned APPROVE with no material
findings. Coverage, review preflight, actual-tree conformance, prose mechanics,
and diff checks are green.

## Problem

Fresh worktree tests no longer provide a short Go feedback loop. Three recent
clean gates reported `internal/worktree` at 130.013, 125.779, and 125.790
seconds. The historical package span was 19.49 seconds.

Compilation is not the long pole. The package repeats selected-binary builds,
materializes broad real-Git fixtures, and serially executes many policy
partitions through public commands.

Process-global `BENCH_HOME`, environment changes, and working-directory changes
also prevent safe scheduling. Adding parallel tests before removing those
inputs would multiply filesystem contention and create state races.

The historical 31.90-second whole-suite floor is not an acceptable target.
It included an unrelated 30.03-second publication connection wait, which a
separate `$bench-debug` repair owns.

## Solution

The worktree test run selects or builds one Bench executable. Every real-binary
journey receives that exact path instead of building its own executable.

The parent `internal/worktree` package remains the serial command and real-Git
adapter. It gathers environment, repository, filesystem, process, and clock
facts once at each boundary.

Three pure packages own landing, lifecycle, and pool-reclaim decisions. Their
tests consume typed facts without Git repositories, descendant processes,
environment mutation, or working-directory mutation.

Representative real-Git journeys remain at every public command seam. The
build records a behavior ledger before it replaces repeated repository-backed
partitions with pure decision tests.

This first spec adds no scheduler and no `t.Parallel`. It records the reduced
demand and leaves measured parallelism plus the regression budget to the
second spec.

## User stories

### One selected executable and explicit inputs

Line: `gpt-5.6-sol` / medium. The user selected this build line, and the seams
are bounded by existing run-binary and command adapters.

1. As a maintainer, I want one Bench executable per worktree test run, so
   repeated build and seal publication do not dominate fresh tests.
2. As a gate operator, I want the tests to reuse the inherited selected
   executable, so a gate does not build a private second binary.
3. As a maintainer, I want a missing or invalid selected executable to fail
   before any public journey, so the suite never mixes executable identities.
4. As a maintainer, I want command boundaries to resolve home and repository
   inputs once, so decision tests can inject exact values.
5. As a maintainer, I want pure owner tests to avoid environment and directory
   mutation, so later scheduling cannot race on process-global state.

### Deep worktree decision owners

Line: `gpt-5.6-sol` / medium. The extraction preserves public behavior while
moving existing decisions behind typed facts.

6. As a maintainer, I want landing policy to decide from typed source,
   destination, resume, and residue facts, so policy partitions need no Git repository.
7. As a maintainer, I want lifecycle policy to decide ownership, lease,
   eligibility, age, ignored-output, and action outcomes in process.
8. As a maintainer, I want pool-reclaim policy to decide protection,
   reclaimability, drift, and action outcomes from typed key facts.
9. As a user, I want command output, exit status, mutation posture, and recovery
   actions unchanged, so the latency repair does not change workflow semantics.
10. As a maintainer, I want real Git retained where Git supplies behavior, so
    native registration, locking, conflicts, and deletion remain proven.
11. As a maintainer, I want each fact adapter covered with focused real Git, so
    typed policy cannot hide a bad repository translation.
12. As a reviewer, I want every replaced journey mapped to surviving behavior
    coverage, so no test disappears only because it is slow.

### Measured handoff to the second spec

Line: `gpt-5.6-sol` / medium. The work is measurement and enforcement evidence
over the reduced implementation, without a timing-policy decision.

13. As a developer, I want before-and-after worktree timings, so the first
    spec shows the demand it actually removed.
14. As a reviewer, I want executable, repository, descendant, environment, and
    directory counts beside timing, so faster output has an architectural explanation.
15. As the second-spec author, I want reproducible raw measurements, so worker
    width and latency budgets are priced from the landed architecture.
16. As a maintainer, I want `-count=1` and the single Go package driver retained,
    so the repair does not weaken freshness or add outer scheduling.
17. As a reviewer, I want this spec to add no `t.Parallel`, so parallelism waits
    for the post-reduction census and a separate approval.

## Implementation decisions

- **The parent package is the effect boundary.** `internal/worktree` keeps the
  public command functions and gathers Git, filesystem, environment, process,
  and clock facts. Public grammar and rendered records do not change.
- **Three child packages own pure decisions.** Use importable packages under
  `internal/worktree` for landing policy, lifecycle policy, and reclaim policy.
  Each exposes typed facts and a decision or plan. A child never imports the
  parent package.
- **Pure means no ambient effects.** The three owners import no `os/exec` or
  `internal/git`. They do not read environment variables, resolve the current
  directory, mutate package globals, or start descendants.
- **Adapters translate once.** Each command resolves repository root,
  `BENCH_HOME`, selected executable, filesystem state, and current time at its
  effect boundary. Lower owners receive those values explicitly.
- **Compatibility adapters are temporary.** An expand ticket may add explicit
  forms beside existing helpers. The contract ticket removes ambient internal
  forms after every caller migrates.
- **One test-run executable owner.** The worktree test harness consumes
  `BENCH_RUN_BINARY` when the gate supplies it. A direct package run builds and
  seals one executable, then passes that path to every binary journey.
- **One serial journey harness owns descendants.** Every worktree test that
  starts Git, Bench, Go, or another process routes through this harness. It
  owns disposable repositories, explicit environment, explicit directories,
  process cleanup, and the selected executable.
- **Journey retention follows behavior.** Landing retains representative
  publish and release, conflict refusal, interrupted resume, and hostile
  residue journeys. Lifecycle retains native create and remove, registration,
  lock, and recovery journeys. Reclaim retains real lease, registration,
  process-liveness, and deletion journeys.
- **One proof inventory closes retained coverage.** The parent test harness owns
  required proof identifiers for the named journeys and one fact adapter per
  policy owner. Tests mark a proof only after their named observation. The
  package run fails when any required proof remains absent.
- **The proof inventory is independently tested.** Removing one journey proof
  and one adapter proof must each make the package red with the missing
  identifier. Restoring each proof returns green.
- **Policy matrices move below adapters.** Ownership, eligibility, lease, age,
  ignored-output, drift, and action combinations use typed facts. A new policy
  partition does not earn another repository because its caller uses Git.
- **No first-spec concurrency.** This spec adds no scheduler and no
  `t.Parallel`. Go may schedule the new packages through its existing `-p`
  behavior, while the descendant-owning parent package stays serial.
- **Measurements have one evidence owner.** The build writes raw commands,
  commit, host conditions, package spans, whole-suite spans, and demand counts
  under `specs/worktree-test-latency/evidence/`.

## Testing decisions

- A good test proves unchanged public command behavior while using fewer
  executable builds, repositories, and descendants.
- Pure table tests attach to each typed decision owner. Adapter tests attach at
  fact translation, and representative journeys attach at public commands.
- The behavior ledger maps every removed repository-backed test to a retained
  pure, adapter, or public-journey test before deletion. Ticket 06 writes it to
  `specs/worktree-test-latency/evidence/coverage-ledger.md`.
- The selected-executable owner gets an injected-builder test. Repeated
  selection calls must return one identity and increment the builder once.
- A gate-inherited executable test requires zero private builds and the exact
  inherited path at multiple public journeys.
- The journey harness records repository and descendant starts. It fails when
  a worktree test starts either effect outside the harness.
- The proof inventory requires landing publish/release, conflict refusal
  without mutation, interrupted resume, and hostile residue journeys.
- It requires lifecycle native create/remove, registration/lock, and recovery journeys.
- It requires reclaim process-liveness, lease, registration, and deletion journeys.
- The same inventory requires one focused fact-adapter proof for landing,
  lifecycle, and reclaim. A missing proof names its owner and class.
- Ticket 06 compares every deleted or replaced repository-backed test function
  in its base-to-tip diff with the coverage ledger. Its verification log records
  the comparison and rejects an unlisted deletion.
- Pure-owner tests include structural assertions against descendant starts,
  environment mutation, and working-directory mutation.
- The ordinary gate remains the oracle. No gate phase, `-count=1` argument,
  package loop, or cache policy changes in this spec.
- Reference measurements use an idle WSL host, normal caches, and no concurrent
  gates. Timing is evidence in this spec, not an ordinary-host red threshold.

### Seam diagram

    bench worktree command or serial public journey
                         │
                         ▼
       [ parent effect adapter and journey harness ]
          │              │                 │
          │ typed facts  │ typed facts     │ typed facts
          ▼              ▼                 ▼
    [ landing policy ] [ lifecycle policy ] [ reclaim policy ]
          ▲              ▲                 ▲
          └──────── pure table tests attach here ─────────┘

    one test-run executable owner ──▶ serial journey harness ──▶ every Bench child
    Git/files/env/CWD adapters       ──▶ typed facts             ──▶ pure decisions

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| SB1 | 1 | A direct worktree test run builds and seals one Bench executable before public journeys, and every journey observes the same identity. | selected-executable owner with injected builder and journey identity log | A per-test builder increments more than once or produces a second path, so repeated construction cannot hide behind fixture helpers. |
| SB2 | 2 | A gate-supplied `BENCH_RUN_BINARY` reaches multiple public journeys unchanged and causes zero private builds. | selected-executable owner with inherited selection | A private fallback build changes the counter or identity even when each journey remains green. |
| SB3 | 3 | A missing, invalid, stale, or seal-less selected executable refuses before any public journey starts. | selected-executable validation with journey-start counter | A permissive fallback or late validation starts a journey or produces the wrong refusal before the counter remains zero. |
| EI1 | 4, 5 | Pure owners receive explicit home, root, time, and repository facts without reading environment or current directory. | owner APIs and pure-package source census | An ambient read or mutation appears in the owner package even when a test supplies a convenient default. |
| LP1 | 6, 9 | Landing decisions over publish, refusal, interruption, resume, and residue facts match command behavior. | landing policy tables plus command adapter tests | A pass-through wrapper cannot satisfy the decision table when one fact changes the required verdict. |
| LC1 | 7, 9 | Lifecycle decisions over ownership, lease, eligibility, age, ignored output, and action facts match command behavior. | lifecycle policy tables plus command adapter tests | The matrix varies each decision input, so a missing ownership or preservation branch turns one partition red. |
| RP1 | 8, 9 | Reclaim decisions protect live or uncertain keys and act only on typed, provably dead keys. | reclaim policy tables plus plan/apply adapter tests | A filesystem-only shortcut misclassifies a live, uncertain, drifted, or hostile key in a named partition. |
| RJ1 | 10 | Each named Git-supplied behavior remains covered by one representative serial real-Git journey. | parent test-run proof inventory and public command journey harness | Deleting a required journey leaves its proof identifier absent and makes the package red by name. |
| FA1 | 11 | One focused adapter per policy owner translates representative real state into the exact typed facts its owner consumes. | parent test-run proof inventory and focused adapter tests | Deleting or bypassing a required adapter proof leaves its owner identifier absent, even when pure tests use invented facts. |
| CV1 | 12 | Every removed repository-backed test has one surviving pure, adapter, or journey disposition before deletion. | `evidence/coverage-ledger.md` and ticket 06 base-to-tip deleted-test comparison | An unlisted deleted test function fails the recorded comparison instead of treating a smaller suite as success. |
| DM1 | 13, 14, 15 | Evidence records reproducible before-and-after timings and counts for executable builds, repositories, descendants, environment changes, and directory changes. | first-spec evidence artifact and recorded commands | A timing improvement without reduced demand, or a demand claim without raw results, leaves a required field absent. |
| GF1 | 16 | The gate still runs one ordinary `go test -count=1 ./...` driver with normal caches and Go-owned package scheduling. | existing gate phase argv tests | A removed freshness flag, package loop, or second driver changes the already enforced argv. |
| NP1 | 17 | The first-spec diff adds no `t.Parallel` and no scheduler to the worktree subtree. | changed-tree source census | Parallel execution cannot enter early under a helper name because the Go call and scheduler imports remain absent. |

### Edge inventory

- Selected executable: inherited, direct-build, missing, invalid, stale, and
  seal-less identities are covered before a public journey starts.
- Environment: empty, explicit, inherited, hostile, and fallback home values
  resolve at the adapter boundary. Pure owners receive only resolved values.
- Repository root: linked worktree, primary checkout, detached checkout,
  outside-repository, missing ref, and moved destination remain adapter concerns.
- Landing: publish, pre-publication refusal, conflict, incomplete reconcile,
  interrupted resume, hostile residue, and release failure retain public journeys.
- Lifecycle: active, pending, complete, orphaned, stale, future, malformed,
  locked, registered, dirty, ignored, and recovery states reach typed decisions.
- Reclaim: live process, dead process, stale lease, unreadable lease, current
  repository key, symlink, FIFO, malformed pointer, drift, and deletion failure remain covered.
- Paths and records: spaces, glob bytes, control bytes, symlinks, and special
  files retain sanitization and no-follow behavior at effect adapters.
- Process failure: start, exit, interrupt, timeout, and cleanup failures remain
  serial journey outcomes.
- Measurement: concurrent gates, cold caches, publication wait, and moving
  source invalidate a reference result instead of changing the target.
- **Won't handle** line: Parallel real-Git or Bench journeys — the second spec
  requires a new census and explicit product budget before any such width.
- **Won't handle** line: Pure-test `t.Parallel` — the second spec prices it from
  the landed first-spec workload.
- **Won't handle** line: Publication's 30-second connection wait — a separate
  `$bench-debug` repair owns that unrelated transport defect.
- **Won't handle** line: Windows-mounted repository timing — the approved timing
  target is the Linux-ext4 reference WSL host.

## Ownership fences

- `internal/worktree/`
- `CHANGELOG.md`
- `capture/session-handoff.md`
- `decisions/worktree-test-latency.md`
- `decisions/assets/worktree-test-invocation-census.md`
- `specs/worktree-test-latency/`

No writer may change gate scheduling, `internal/runbinary`, another production
package, or shared Bench rules without reviewer approval and a fence update.

## Out of scope

- Measured `t.Parallel`, worker width, and the slow-package regression budget
  belong to the second spec. Estimated cost: 8 edits and 2 gate runs.
- Publication's uncontrolled request wait belongs to `$bench-debug`.
  Estimated cost: 4 edits and 2 gate runs.
- The WSL/Codex Go bootstrap usability defect is a separate track.
  Estimated cost: 8 edits and 2 gate runs.
- New CLI grammar, output records, exit meanings, or recovery behavior are not
  part of test-demand reduction. Estimated cost: 6 edits and 2 gate runs.
- Removing `-count=1`, clearing Go caches, diff-scoped green verdicts, or a
  Bench-owned package scheduler are prohibited, not deferred capabilities.
- Replacing every public Git journey with mocks is prohibited. Representative
  Git-supplied behavior remains at each command seam.

## Further notes

The second spec starts only after this spec lands and its evidence is stable.
It must meet the ready map's three-run package and whole-suite envelopes without
raising limits or weakening freshness.

The whole-suite envelope becomes authoritative only after the publication wait
repair lands. Until then, evidence reports that wait separately and makes no
worktree improvement claim from its removal.
