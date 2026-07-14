# Lossless shift recovery and truthful result states

Status: staged

## Problem

A `bench shift` that fails after the agent has mutated the tree destroys the work.
`worktree.Release` hard-resets and cleans, so a commit, staging, teardown, gate, or
signal failure erases every uncommitted change. The result states also lie: cap
exhaustion exits `0`, a `/bin/false` adapter plus a green gate reports "objective
likely met" and exits `0`, a missing objective silently becomes "improve the
codebase", and `envInt` accepts zero, negative, non-integer, and unbounded caps.
There is no wall deadline, so a hung adapter runs unbounded. An operator cannot tell
from the exit code or the terminal whether work landed, was preserved, or vanished —
and a closed terminal loses the preserved location entirely.

## Solution

Every shift failure after agent mutation preserves the work — a durable recovery ref
by default, a locked worktree as fallback — and prints its exact location. Staging,
adapter, commit, signal, wall, and teardown errors propagate into honest outcomes
instead of vanishing into exit `1` or false success. A shift requires an objective,
validates positive bounded iteration and wall caps before acquiring anything, runs
under a wall deadline, and reports one of five outcomes — complete, incomplete,
failed, interrupted, no-op — through distinct exit codes and a structured
`shift_result` block on stdout, mirrored onto the intent ledger so `bench status`
can surface the recovery pointer after the terminal is gone.

The build routes to the cheap tier at elevated effort (the reviewer's standing call
for this spec): the seams and the gate signal are pinned here, so `gpt-5.6-luna`
(Claude Code alias `sonnet` / Sonnet 5 for in-session delegates) executes precisely
while the runtime gate observes every outcome. This is a deliberate deviation from
the profile's cached "gate/conformance → mid" routing; it is paired with high effort
and interactive TDD because a wrong oracle is the worst bug class in a kit.

## User stories

1. As an operator, I want any failure after the agent mutates the tree to snapshot
   the dirty tree to a durable recovery ref by default, so a crash never erases
   uncommitted work. Line: gpt-5.6-luna / high. This is the load-bearing preservation
   path and its ordering against `Release` is where a mistake silently loses work.
2. As an operator, I want the snapshot built with git plumbing under a temp index and
   an explicit synthetic identity, so a missing Git identity or a hostile hook cannot
   block preservation the way an ordinary commit would. Line: gpt-5.6-luna / high.
   The plumbing-vs-commit asymmetry is the whole reason snapshots survive where
   iteration commits fail, so it must be exact.
3. As an operator, I want the recovery ref to live at `refs/bench/recovery/<shift-branch>`
   with creation failing closed on any pre-existing ref, so a timestamp-collision
   branch name cannot overwrite another shift's evidence. Line: gpt-5.6-luna / medium.
   Fail-closed ref creation is a single `update-ref` with a zero old-oid and is fully
   gate-observable.
4. As an operator, I want the shift to fall back to retaining and `git worktree lock`-ing
   the charged worktree — and printing its path — only when the snapshot itself fails,
   so preservation degrades safely instead of releasing into the pool. Line:
   gpt-5.6-luna / high. The fallback is the only legal response to a snapshot failure
   and must never run the release-and-clean path.
5. As an operator, I want the exact recovery location printed on every preserving
   failure, so a closed or scrolled terminal is not the only record of where my work
   went. Line: gpt-5.6-luna / medium. The pointer rides the structured block and the
   ledger, both gate-observable.
6. As a script author, I want distinct exit codes — `0` complete, `1` failed, `2`
   usage/preflight, `3` incomplete, `4` no-op, `130` interrupted — so I can branch on
   the real outcome instead of guessing. Line: gpt-5.6-luna / high. The taxonomy is
   the contract every downstream row asserts against.
7. As an operator, I want adapter spawn failure or a nonzero adapter exit to stop the
   loop after processing any mutations, splitting the outcome by committed evidence
   (`failed` at zero commits, `incomplete` after commits), so a broken harness never
   reports false success. Line: gpt-5.6-luna / high. This inverts today's
   "the gate decides, ignore the adapter" posture and needs care at the evidence split.
8. As an operator, I want a `BENCH_MAX_WALL` deadline (Go duration, default `2h`,
   bounded `(0, 24h]`) that on expiry kills the adapter process group, cancels a
   running gate, snapshots mutations, and exits `3`, so a hung adapter cannot run
   forever. Line: gpt-5.6-luna / high. The timer wires into the existing
   signal/checkpoint machinery and must not rely on deferred cleanup.
9. As an operator, I want a shift with an empty objective to exit `2` before acquiring
   a worktree, so the loop never silently invents "improve the codebase". Line:
   gpt-5.6-luna / medium. A single entry guard, fully observable.
10. As an operator, I want an objective carrying control bytes to exit `2` at entry,
    so hostile text is rejected before it can reach the ledger or the TOON emitter.
    Line: gpt-5.6-luna / medium. Entry validation with the TOON refusal as backstop.
11. As an operator, I want `BENCH_MAX_ITERS`, `BENCH_REFACTOR_ITERS`, and
    `BENCH_MAX_WALL` that are set-but-invalid (non-integer, `< 1`, over bound) to exit
    `2` naming the variable and accepted range before any worktree acquire, while
    unset or empty-string keeps the defaults, so a typo fails loud instead of running
    something I did not ask for. Line: gpt-5.6-luna / high. The absent-vs-empty
    distinction and the pre-acquire ordering are both easy to get subtly wrong.
12. As an operator, I want cap exhaustion to leave the committed work on the branch and
    exit `3`, not `0`, so "ran out of iterations" is never mistaken for "done". Line:
    gpt-5.6-luna / medium. A single changed exit path.
13. As an operator, I want a shift that commits zero iterations — even when the adapter
    exits clean and `.bench/done.sh` passes on iteration 1 — to exit `4` with a detail
    noting the predicate was already satisfied, so "already done" is distinguishable
    from "did work". Line: gpt-5.6-luna / high. This redraws the current
    "objective likely met → exit 0" boundary and interacts with the done.sh path.
14. As an operator, I want a staging failure (`git add`) in the main or refactor phase
    to propagate — snapshot and split by evidence — instead of being ignored, so a
    partially-added tree is never committed or lost. Line: gpt-5.6-luna / high.
    Reaching this needs an in-process fault seam because the shell cannot force it.
15. As an operator, I want a teardown/`Release` error to exit `1` with a detail naming
    what is already safe (the branch, any recovery ref), so a cleanup failure is
    reported rather than swallowed by a silent void. Line: gpt-5.6-luna / medium.
    Also an in-process fault seam.
16. As an operator, I want a structured `shift_result` TOON block on stdout carrying
    outcome, exit code, branch, committed count, iterations used, recovery pointer,
    and detail, so tooling parses one control-byte-safe surface. Line:
    gpt-5.6-luna / medium. Straight TOON emission over the computed result.
17. As an operator returning later, I want the outcome and recovery pointer recorded on
    the shift's intent-ledger entry and rendered by `bench status`, so I can discover a
    preserved shift after the terminal is gone. Line: gpt-5.6-luna / medium. Ledger
    field addition plus a status render line.
18. As an operator, I want `bench resume` and the worktree classifier to retain — never
    delete — recovery refs and locked fallback worktrees, and a re-run shift to get a
    fresh branch/ref pair, so recovery evidence survives housekeeping and repeat runs.
    Line: gpt-5.6-luna / high. Retention across the resume path is a safety property and
    the re-run path must not collide with prior evidence.

## Implementation decisions

- **`internal/shift` is the deep owner of the outcome model.** A single result value
  (outcome, exit code, branch, committed count, iterations used, recovery pointer,
  detail) is computed at every exit path, emitted once as the `shift_result` TOON
  block, mirrored onto the intent entry, and returned as the process exit code.
  `Loop` stops returning bare `1`s; each path resolves to a taxonomy member. The
  `checkpoint()`/`os.Exit` paths (which skip deferred cleanup) run snapshot, emit, and
  record explicitly rather than via `defer`.
- **Outcome → exit code** is the map's matrix verbatim: `complete 0`, `failed 1`,
  `usage/preflight 2`, `incomplete 3`, `no-op 4`, `interrupted 130`. Every setup step
  before the first adapter run (objective/cap validation, adapter preflight, acquire,
  branch create, config, scratch write) resolves to `2` (nothing ran). Post-mutation
  failures resolve by the evidence split: `failed` when zero iterations committed this
  shift, `incomplete` when prior green commits exist. Teardown error is `1`. Cap and
  wall exhaustion are `3`. A zero-commit clean finish is `4`.
- **Validation runs before acquire.** Empty objective and control-byte objective exit
  `2` at entry. `BENCH_MAX_ITERS` and `BENCH_REFACTOR_ITERS` are integers in `[1,100]`
  default `12`/`4`; `BENCH_MAX_WALL` is a Go duration in `(0, 24h]` default `2h`.
  Unset or empty string keeps the default; set-but-invalid exits `2` naming the
  variable and the accepted range. `envInt` is replaced by validating parsers that
  distinguish unset/empty from present-but-invalid.
- **The snapshot primitive lives in `internal/worktree`,** reusing the recovery-envelope
  machinery (`temporaryIndex`, `commitTree` with the fixed `bench`/`bench@local`
  synthetic identity, the fail-closed `update-ref` with a zero old-oid). It takes the
  worktree root, the parent commit (the shift-branch tip), the target ref name, and an
  explicit set of paths to exclude; it runs `read-tree HEAD` → `add -A` →
  `rm --cached` the excluded paths → `write-tree` → `commit-tree` → fail-closed
  `update-ref`, and verifies the ref resolves to the new commit before returning. The
  exclude set is passed in by `internal/shift` (which owns the scratch names), so the
  primitive stays generic and the scratch policy stays single-sourced. Probed: in a
  linked worktree, `add -A` captures untracked and nested files; `rm --cached` drops
  the scratch; the second `update-ref` with a zero old-oid is rejected.
- **The retain-and-lock fallback** is a second `internal/worktree` entry: it
  `git worktree lock --reason <reason>`s the charged worktree and drops its lease
  file *without* restoring cleanliness, so the dirty work stays on disk. No classifier
  change is needed — a no-owner-marker locked worktree already retains as
  `ReasonUnexpectedLock` (`subshell.go`), and `Acquire` already skips it because it is
  not clean.
- **Release ordering is load-bearing.** The snapshot must be created and verified
  before any release path runs. On snapshot success the pool worktree releases
  normally; on snapshot failure only the retain-and-lock path is legal. This ordering
  replaces the current red-gate "preserve the worktree verbatim" special case with one
  uniform rule for every post-mutation failure.
- **Adapter exit becomes evidence for progress.** `runAdapter` captures the child's
  exit status. A spawn failure or nonzero exit stops the loop after processing that
  iteration's mutations (gate → commit on green → otherwise snapshot). The
  "no change = objective met" inference additionally requires adapter exit `0`. The
  gate stays the oracle for *work*; the adapter exit is evidence for *progress* only.
- **The wall deadline** is a timer in the session that, on expiry, sets a deadline flag
  and performs the same kill-process-group + cancel-gate + preserve steps as the signal
  handler; the checkpoint distinguishes deadline (exit `3`, deadline detail) from signal
  (exit `130`).
- **`internal/intent`** gains outcome and recovery fields on the shift `Entry`;
  `internal/status` renders the recovery pointer on the intent signal.
  `cmd/bench/main.go`, the shipped adapters, and the hooks keep their current thin
  shape — `run()` already returns the shift's int straight to `os.Exit`, so the new
  codes propagate unchanged. Harness-specific "exit 0 = clean" normalization in the
  shims is out of scope (see below); the loop's contract is proven with scripted
  adapters.

## Testing decisions

- **A good test here asserts observable outcome, not internals:** the process exit
  code, the `shift_result` fields on stdout, git ref/tree/branch state, intent-ledger
  JSON, and pool-lease / worktree-lock presence — never the loop's private control
  flow. This keeps the expectations independent of the implementation, which is what
  lets a named omission or mutation turn the gate red.
- **Two seams, matching the Handoff.** Most rows attach at **Seam A**, the built-binary
  runtime shift contract (`internal/contract/runtime/runtime_shift_test.go`) driving
  `bench shift` through throwaway fixture repos with scripted adapters. The two rows the
  shell cannot force — staging failure and teardown error — plus the snapshot
  primitive's own mechanics attach at **Seam B**, in-process Go tests using the
  worktree package's existing `Fault` pattern and direct primitive calls
  (`internal/worktree`). These are the TDD points.
- **Prior art to compose:** `runtime_shift_test.go` (existing shift contracts:
  `shiftFixture`, `f.BenchEnv`, `probe.RequireExit`, `shiftBranch`, `shiftWorktree`,
  `requireNoLease`, `requireRegisteredWorktree`) for Seam A;
  `recovery_retry_test.go` and `lifecycle_acquire_test.go` (the `Fault`/`LifecycleStep`
  seam) and `clean.go`'s temp-index helpers for Seam B.
- **Gate command:** the project gate, `.bench/gate.sh`. The runtime and worktree
  fragments run inside its Runtime-and-behavior-contracts phase; the Go unit tests run
  under the compiled-core checks.

### Seam diagram

**Seam A — built-binary runtime shift contract** (exit code, `shift_result`, git/intent state):

    trigger: a runtime contract fragment runs the built `bench` binary
        │
        ▼
    objective + env (BENCH_MAX_ITERS / _WALL / AGENT)   ──▶  [ bench shift          ]  ──▶  exit code (0/1/2/3/4/130)
    scripted adapter (/bin/false, write-then-nonzero,       [  Loop: validate →     ]  ──▶  shift_result TOON on stdout
      identity-stripped + failing hook, sleeper,            [  acquire → iterate →  ]  ──▶  refs/bench/recovery/<branch> (git show-ref)
      always-dirty, SIGTERM mid-adapter)                    [  snapshot/commit →    ]  ──▶  branch commit count, snapshot tree (ls-tree)
                                                            [  emit + record →      ]  ──▶  intent ledger JSON (outcome + recovery)
                                                            [  release/retain-lock  ]  ──▶  pool lease gone / worktree locked
                        ◀ tests attach here: BenchEnv runs the binary in a fixture repo,
                          then asserts RequireExit + RequireContains(shift_result …) and
                          inspects git refs, intent JSON, and worktree registration.

**Seam B — in-process worktree fault + snapshot primitive** (staging failure, teardown error, snapshot mechanics):

    trigger: a Go unit test calls the primitive / drives the loop with a fault injected
        │
        ▼
    worktree root + parent + ref + exclude set   ──▶  [ SnapshotDirty:            ]  ──▶  recovery ref OID (rev-parse)
    injected Fault(step) at StageTouched/Teardown     [  temp-index add -A →      ]  ──▶  tree contents (untracked in, scratch out)
                                                       [  rm --cached scratch →    ]  ──▶  error on pre-existing ref (fail closed)
                                                       [  write-tree → commit-tree ]
                                                       [  → fail-closed update-ref ]  ──▶  RetainAndLock: locked + lease dropped
                        ◀ tests attach here: call the primitive directly and assert the
                          tree/ref, or set the Fault var and assert the loop's outcome
                          value and preserved state.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | red-gate iteration after a mutation snapshots the dirty tree and releases the pool worktree | Seam A: scripted red gate + writing adapter, `BENCH_MAX_ITERS=1` | `git show-ref refs/bench/recovery/<branch>` resolves and its tree contains the mutated file; pool lease gone | if the failure erased the work (today's `Release`), no ref exists and `ls-tree` lacks the file |
| 1,5 | preserving failure prints and reports the recovery ref location | Seam A | `shift_result` recovery cell equals `ref:refs/bench/recovery/<branch>` | a silent erase or a stdout-only pointer leaves the cell `none` or absent |
| 2 | snapshot commit is built with synthetic identity under a failing-hook/identity-stripped env | Seam A: identity-stripped env + failing `pre-commit` fixture | recovery ref exists and resolves though the iteration commit would fail; `shift_result` shows `failed`/`incomplete` not exit 0 | if the snapshot used an ordinary commit it would block exactly where the iteration commit blocks, losing the work |
| 3 | recovery ref creation fails closed against a pre-existing ref | Seam B: pre-create the ref, call `SnapshotDirty` | primitive returns an error and does not move the existing ref | a plain `update-ref` (no zero old-oid) would overwrite another shift's evidence |
| 3,18 | a re-run shift gets a fresh branch/ref pair | Seam A: two shifts | second shift's branch and recovery ref differ from the first; neither is clobbered | a fixed or colliding ref name would overwrite the first shift's preserved work |
| 4 | when the snapshot fails, the worktree is retained, `git worktree lock`-ed, and its path printed | Seam B: `SnapshotDirty` fault + `RetainAndLock`; Seam A asserts the message | `git worktree list --porcelain` shows the path `locked`; `shift_result`/stdout names `worktree:<path>`; lease dropped | if it released, `worktree list` loses the path and the dirty work is gone |
| 4,18 | a locked fallback worktree is retained by resume, never deleted | Seam A: lock a pool worktree, run `bench resume-clean` | the worktree remains registered and locked after resume | the classifier already retains it as `ReasonUnexpectedLock`; a regression that deleted it would drop the registration |
| 6,16 | a complete shift reports outcome `complete`, exit 0, and a full `shift_result` block | Seam A: green gate, committing adapter | exit 0 and `shift_result` carries outcome/exit/branch/committed/iterations_used/recovery/detail | a missing or malformed block, or a wrong outcome, fails the field assertions |
| 6,7 | `/bin/false` adapter with zero commits stops the loop as `failed`, exit 1 | Seam A: `BENCH_AGENT=/bin/false`, green gate | exit 1 and `shift_result` outcome `failed`, committed 0 | today this exits 0 "objective likely met"; the assertion on exit 1 + `failed` catches the false success |
| 7 | adapter nonzero exit *after* a prior green commit stops the loop as `incomplete`, exit 3 | Seam A: adapter commits once then exits nonzero | exit 3, outcome `incomplete`, committed ≥ 1, prior commit on branch | folding this into `failed`/1 would understate the committed work |
| 8 | a sleeping adapter under a tiny `BENCH_MAX_WALL` is killed, mutations snapshotted, exit 3 | Seam A: sleeper adapter, `BENCH_MAX_WALL=1s` | exit 3, outcome `incomplete` with a deadline detail; no orphaned adapter; recovery ref if mutated | no wall means the shift hangs forever; a checkpoint-only wall never fires against a hung adapter |
| 9 | an empty objective exits 2 before acquiring a worktree | Seam A: `bench shift` with no objective | exit 2, no lease created, no branch | today it silently runs "improve the codebase"; the exit-2 + no-lease assertion catches the default |
| 10 | a control-byte objective exits 2 at entry | Seam A: objective containing `ESC` | exit 2 before acquire | an unvalidated objective would reach the ledger/TOON and either corrupt output or be refused late |
| 11 | set-but-invalid `BENCH_MAX_ITERS`/`_REFACTOR_ITERS`/`_WALL` exit 2 naming the variable and range, before acquire | Seam A: `BENCH_MAX_ITERS=0`, `=abc`, `BENCH_MAX_WALL=48h` | exit 2, stderr names the variable and accepted range, no lease | `envInt`'s silent fallback runs a shift the operator never authorized |
| 11 | unset or empty-string cap env keeps the default | Seam A: unset vs `BENCH_MAX_ITERS=""` | both run to the default cap (no exit 2) | conflating empty with invalid would reject a legitimately-unset knob |
| 12 | cap exhaustion leaves work on the branch and exits 3 | Seam A: always-dirty adapter, `BENCH_MAX_ITERS` small | exit 3, outcome `incomplete`, committed count equals the cap | today cap exhaustion exits 0, reading as done |
| 13 | zero commits with a clean adapter and a passing `done.sh` exits 4, detail notes the predicate | Seam A: no-op adapter + `.bench/done.sh` that passes | exit 4, outcome `no-op`, detail mentions the predicate | complete-on-predicate would exit 0 and hide that nothing was done |
| 14 | a staging failure snapshots and splits by evidence instead of committing a partial tree | Seam B: `Fault` at the stage step | loop resolves to `failed`/`incomplete` with a recovery ref, no partial commit | today `stageTouched` ignores `git add` failure and could commit or lose a partial tree |
| 15 | a teardown/`Release` error exits 1 with a detail naming what is safe | Seam B: `Fault` at teardown | exit 1, detail names the branch and any recovery ref | a silent teardown void would report success despite a real cleanup failure |
| 17 | outcome and recovery pointer are recorded on the intent entry and rendered by `bench status` | Seam A: run a preserving shift, then `bench status` | intent JSON entry carries outcome + recovery; `bench status` prints the pointer | a stdout-only pointer is lost when the terminal closes; the ledger assertion catches the omission |
| 1,18 | snapshot excludes scratch and captures space/glob/newline and untracked paths whole | Seam B: `SnapshotDirty` over a worktree with `step 1 [a].txt`, nested untracked, and scratch | snapshot tree contains the hostile-named and untracked paths, excludes `.bench-objective`/`.bench-notes.md` | a naive `add`/pathspec would drop untracked or nested files or ride scratch into the snapshot |

### Edge inventory

Walked against the profile's shell-CLI hostile-input checklist:

- **paths with spaces/glob/newline** → covered (row 1,18): `add -A` under the temp index
  inherits whole-tree safety; asserted with `step 1 [a].txt` and a nested untracked path.
- **control bytes in git-sourced / objective text** → covered (row 10): objective
  validation exits 2 at entry, with `toon.Table` refusal as the backstop for any
  git-sourced cell in `shift_result`.
- **absent vs present-but-empty env var** → covered (row 11): empty string = unset =
  default; both asserted.
- **required tool missing / bad `BENCH_AGENT`** → covered by the existing adapter
  preflight, now resolving to exit 2 (usage/preflight) rather than 1.
- **destructive worktree state; recovery refs + locked fallback retained** → covered
  (row 4,18): resume/classifier retain the locked worktree and never delete recovery
  refs.
- **interrupt (SIGINT/SIGTERM) mid-loop** → covered: SIGTERM mid-adapter kills the
  process group, snapshots mutations, and exits 130; the existing interrupt-cleanup
  contracts extend to assert the recovery ref and the 130 outcome.
- **re-run idempotency** → covered (row 3,18): fresh branch/ref pair per run; ref
  creation fails closed on conflict.
- **hand-edited file with no trailing newline** — **Won't handle**: the snapshot
  captures raw bytes verbatim through `write-tree`; a missing final newline changes no
  behavior and needs no dedicated row.
- **unquoted multi-word arguments (`$*` vs `$1`)** — **Won't handle**: `shift.Command`
  already joins positional args into the objective; FT79 adds no new argument surface.
- **invocation through a symlink; cwd deeper than the repo root** — **Won't handle**:
  routing stays on the existing `git.Root` resolution, which FT79 does not touch.
- **every shipped surface reaching one implementation** — **Won't handle**: FT79 leaves
  the adapters and hooks at their current thin shape; adapter exit-code normalization is
  out of scope below.

## Out of scope

- **Garbage collection / retirement of consumed recovery refs.** A separate capability
  (a lifecycle for `refs/bench/recovery/*` after the operator has recovered them), not
  part of preserving the work. Fog in the map; deferred deliberately. ~3 edits, 2 gate runs.
- **Adapter-author contract doc for "exit 0 = clean" normalization in the shims.** A
  documentation capability for harness authors, separate from the loop's evidence
  handling, which is proven here with scripted adapters. ~2 edits, 1 gate run.
- **FT71 versioned shift evidence (append-only event log).** A distinct persistence
  surface that will re-carry the outcome + recovery fields; alignment is FT71's problem.
  Not estimated here — it is its own roadmap item.
- **FT58 identity-safe lease reclamation.** Pool-lock protocol beyond the release
  ordering FT79 needs; a separate concurrency capability. Not estimated here.
- **Any gate-execution or verdict-cache change (R-02 territory).** FT79 changes nothing
  about how the gate runs or caches; out of scope by construction.
