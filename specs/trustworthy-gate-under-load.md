# trustworthy-gate-under-load

Status: staged

## Problem

FT88's deadline fix bought load tolerance by giving up fail-fast: the R14 gate-owner
wait and the losing-racer repair wait now burn a flat 60s before reporting anything,
so a genuinely dead or wedged child costs a minute per subtest instead of seconds.
Three further defects still tax every landing:

- `TestRuntimeGateContracts/bench_gate_rebuilt_self-host_contract` fails under gate
  load with `TempDir RemoveAll cleanup: unlinkat .../.bench-contract-env: directory
  not empty` — a child still writing when Go's cleanup runs. It passes in 0.4s
  isolated and has never been reproduced on demand.
- Killing a running gate at the wrapper leaves the gate script's process group alive:
  `bench canary` reparents to init and keeps spawning nested `gate-phases` children.
  The next `bench gate` then reports the bare `gate execution already in progress`,
  which names nothing the operator can check, so a live orphan and a phantom lock look
  identical from the outside.
- The conformance phase reports `go test failed` with no output attached. Four of the
  five reds in the FT76 recurrence said only that, so each cost a manual re-run of the
  inner command to learn what actually broke.

`GATE-REPORT.md` is the closed diagnosis behind all of this. It is a committed
report that has outlived its fix and now reads as a second source of truth beside the
project profile.

## Solution

Give the gate one **owner record** — a plain, non-durable file written after lock
acquisition and before the pending durable write — and make it serve both jobs at
once: the marker waits get a fast liveness leg that fails in seconds on a real wedge
while only the fsync-stretchable leg carries the 60s tolerance, and the
already-in-progress refusal names the PID that holds the lock instead of asserting a
condition the operator cannot verify. Then make `gate-run` tear its process group down
on a signal so no orphan is created in the first place, reap spawned groups in the
contract fixture harness so no child outlives the test that spawned it, attach the
failing command's output tail to every conformance diag, and delete `GATE-REPORT.md`
after promoting its durable content into the profile and learnings.

**No decision map backs this spec.** It was written same-session under the batch-drain
override in `/bench-write-spec`, compiled from the `ROADMAP.md` FT88 row (a reviewed
`/bench-what-next` drain) and `GATE-REPORT.md`. Every decision the author defaulted is
flagged **(defaulted — veto point)** below. A mid-tier falsification pass ran against
the first draft and returned BLOCK; its findings are folded in, and the decisions it
surfaced as silently made are now flagged.

## User stories

1. As the gate operator, I want `gate-run` to write an owner record naming its PID and
   start time into the git directory immediately after it acquires the gate lock and
   **before** the pending durable write, and to remove that record when the run ends,
   so that a gate run is externally observable before the write that host I/O can
   stall.
   Line: `gpt-5.6-terra` / medium. This adds a new artifact to the oracle's execution
   path, where a wrong write, a wrong ordering, or a missed removal would either forge
   liveness or make every later run report a phantom owner.

2. As the gate operator, I want a second `bench gate` that cannot take the lock to
   report the owning PID and whether that process is still alive, instead of the bare
   `gate execution already in progress`, so that an orphaned gate tree is diagnosable
   from the refusal message alone.
   Line: `gpt-5.6-terra` / medium. The message is the operator's only handle on a
   contended lock, and it must stay honest when the record is absent, stale, or
   unreadable.

3. As the gate operator, I want `gate-run` to cancel its run and tear down the gate
   script's process group when it receives SIGINT, SIGTERM, or SIGHUP, so that killing
   a gate kills the whole tree rather than orphaning `canary` and its nested
   `gate-phases` children.
   Line: `gpt-5.6-terra` / medium. Signal handling on the oracle's own entry point
   decides whether an interrupted gate leaves a recorded verdict and a clean machine.

4. As the next debugger, I want one shared two-leg marker wait — a fast leg with its
   own short deadline and a slow leg allowing 60s — that reports which leg missed, so
   that a wedge before the durable write fails in seconds while a stalled fsync still
   passes.
   Line: `gpt-5.6-terra` / medium. The helper carries concurrency the gate cannot
   observe end to end, so it is built as an injectable unit with its own tests rather
   than as inline loop edits.

5. As the next debugger, I want the R14 gate-owner wait to use that helper with the
   owner record as its fast leg (5s) and the gate script's started-marker as its slow
   leg (60s), keeping the existing exit fast-fail and diagnostic dump, so that a gate
   child that never reaches the lock is reported in seconds.
   Line: `gpt-5.6-terra` / medium. It rewires the helper that already owns the single
   `Wait` and the in-loop exit check, where a mistake reintroduces a double-`Wait`.

6. As the next debugger, I want the losing-racer repair wait to use the same helper,
   with the repair child's pre-fetch start marker as its fast leg (10s) and the
   existing post-`fh.sync()` ready marker as its slow leg, so that a repair child that
   dies or never starts is reported in seconds rather than after a full minute.
   Line: `gpt-5.6-terra` / medium. The start hook is a blocking handshake, not a
   passive marker, so wiring it wrong hangs a currently-green test for the full slow
   leg.

7. As the gate operator, I want the contract fixture's command runner to reap the
   whole spawned process group before it returns, and `isolatedEnv` to remove
   `.bench-contract-env` itself ahead of Go's `TempDir` cleanup, so that no child
   outliving a subtest can fail an otherwise green gate at cleanup time.
   Line: `gpt-5.6-terra` / medium. It changes the harness every contract test runs
   through, so a wrong reap turns one flaky test into a broadly broken suite.

8. As the next debugger, I want every conformance core-check diag that reports a failed
   subprocess to carry the tail of that subprocess's output through one shared
   formatter, so that a red conformance phase attributes itself without a manual
   re-run of the inner command.
   Line: `gpt-5.6-terra` / medium. This is conformance-check logic, which the project
   profile routes to the mid tier because the oracle's correctness outranks speed.

9. As the teammate who just walked in, I want `GATE-REPORT.md` deleted, its durable
   content promoted — the three traps and the WSL2 fsync-stall hazard into
   `projects/benchkit.md`, the reproduction economics into `.bench/learnings.md` — and
   its six surviving references rewritten in the same commit, so that the repository
   carries one source for how the gate behaves under load and no reference dangles.
   Line: `gpt-5.6-sol` / high. This is the profile's cached doc-authoring leverage
   routing, not a fresh escalation; see the note under Implementation decisions.

## Implementation decisions

- **One owner record, three consumers.** `<git-dir>/bench-gate-owner` is written with a
  plain write (no fsync, no temp-and-rename) right after `Acquire` succeeds and before
  the pending `durableReplace`, and removed on every exit path. It is deliberately
  *not* durable: its only job is to be observable cheaply, and a crash that loses it
  is indistinguishable from the crash that ended the run. Stories 1, 2, and 5 all read
  this one file rather than each inventing a liveness signal.
- **Path and name are `<git-dir>/bench-gate-owner`** (defaulted — veto point). It sits
  beside the existing `bench-gate.lock`, `bench-last-gate`, and `bench-gate-pin`, so
  it inherits their gitignore-free, git-dir-resident convention and collides with
  none of them.
- **The owner write goes through the `gateEngine` interface**, like every other file
  operation in `internal/gate`. This is what makes story 1's ordering claim testable:
  the existing fake engine (`internal/gate/fault_engine_test.go`) can record the call
  sequence and assert the owner write precedes the first `durableReplace`. Without
  routing it through the engine, an implementation that writes the record *after* the
  fsync — which defeats the entire point of the record — is externally
  indistinguishable from the correct one.
- **Record contents** are the owner PID and an RFC3339 start time (defaulted — veto
  point). The pending verdict record already carries `OwnerPID` and `StartedAt`, but it
  is written *by* the fsync this whole arm is working around, so it cannot be the
  liveness source.
- **A dead holder cannot hold this lock.** The gate lock is a POSIX `fcntl` record
  lock (`internal/gate/verdict.go`), released by the kernel when the owning process
  dies. So "stale-holder detection" as the roadmap phrases it has no reclaim to
  perform: a contended lock always has a *live* holder, and the real defect is that
  the holder may be an orphan the operator believes they killed. Story 2 therefore
  reports rather than reclaims, and story 3 stops the orphan being created
  (defaulted — veto point; the alternative reading is a lock-takeover path, which
  would put a second writer on the verdict record and is not proposed).
- **Refusal message shape**: `gate execution already in progress` stays as the first
  line, with the owner PID and its liveness appended when the record is readable. An
  absent or unparseable record degrades to today's bare message — the refusal must
  never fail on a diagnostic read.
- **Signal teardown reuses `runProcessGroupCommand`.** `RunCommand` builds a
  `signal.NotifyContext` and passes it to `Execute`, which already runs the gate script
  under `Setpgid` and already SIGINTs the group with a 2s grace on cancellation.
- **SIGHUP is added to the set, which `gate-phases` does not have** (defaulted — veto
  point). `phases.go` registers SIGINT and SIGTERM only. SIGHUP is included because
  the orphan scenario starts with a closed terminal as often as with a kill, and the
  same decision should then propagate to `gate-phases` so the two stay consistent.
  Vetoing SIGHUP leaves both at SIGINT+SIGTERM and costs nothing else in this spec.
- **A cancelled gate stays pending, not red.** The existing `ctx.Err()` branch in
  `executeWithEngine` already declines to record a verdict for a cancelled run; story
  3 must not change that.
- **The two-leg wait is one injectable helper**, not two inline loop edits. It takes
  the two marker paths, the two deadlines, a stat function, and a liveness/exit
  channel, and returns which leg missed. Building it this way is what gives the
  roadmap's headline arm a real red instead of a review-only grade; the first draft
  had it as inline edits with no biting test, which the falsification pass flagged as
  the largest coverage concession in the spec.
- **Fast-leg deadlines are 5s for R14 and 10s for the repair child** (defaulted — veto
  point, and a deviation from the roadmap). `ROADMAP.md`'s FT88 row says a true wedge
  should fail in **~2s**; 5s is proposed instead because the R14 fast leg spans spawn,
  exec, git subject build, and lock acquisition, none of which fsync but all of which
  stretch under CPU load. The repair child gets 10s because Node's cold interpreter
  start-up is itself inside the window that flaked, per `GATE-REPORT.md`'s third-flake
  section. Vetoing back to 2s is a one-constant change.
- **The slow leg stays 60s at both sites**, unchanged from the shipped fix.
- **The repair start hook is a blocking handshake, not a passive marker.**
  `waitForTestRelease` in `bin/bench-repair-binary.mjs` creates the named file and then
  *blocks until the test removes it*. Story 6's fast leg must therefore remove the
  start marker the moment it observes it, exactly as `testRepairEarliestInterrupt`
  already does. Getting this wrong hangs the child until the 60s slow leg expires.
- **No exit fast-fail is added to the repair wait** — the retired
  `load-tolerant-marker-deadlines` spec closed that decision (a watcher goroutine
  would reintroduce a double-`Wait`), and story 6's fast leg gets the same benefit
  without new concurrency.
- **The contract harness reaps by process group.** The fixture's command runner sets
  `Setpgid` and, after `Wait`, signals the group so no grandchild survives the call.
  `isolatedEnv` additionally registers a `t.Cleanup` that removes
  `.bench-contract-env`, retrying **up to 3 times over 2s** (defaulted — veto point);
  because it is registered after `t.TempDir()`, it runs first, so the retry absorbs a
  late writer instead of the `RemoveAll` failing the test. Both live in the one
  harness, not repeated per test.
- **One conformance diag formatter.** A single helper takes a label and a `Probe` and
  returns the label plus a bounded, control-byte-stripped tail of the probe's output.
  Every subprocess-failure diag reachable from the gate routes through it. The
  gate-reachable sites are enumerated, not sampled: `npm pack --dry-run failed`, `go
  build setup failed`, both `go build failed` sites, `go vet failed`, `go list
  failed`, `go test failed`, `worktree cleanup race test failed`, and `worktree
  cleanup race test did not run`.
- **`crossCompileMatrix` is out of the formatter's coverage** (defaulted — veto point).
  Its real implementation is behind `//go:build stress`; the default build the gate
  runs returns `nil`, so a row asserting its diags would grade nothing. It is listed
  under Won't handle rather than claimed.
- **Tail bound is the last 40 lines, truncated to 4 KiB** (defaulted — veto point).
  Large enough for a Go test failure block, small enough that a pathological red does
  not bury the phase output.
- **`GATE-REPORT.md` promotion split**: the three traps (never plain-`go build` the
  dist binary, never write in the repo mid-gate, never kill only the wrapper) and the
  WSL2 host-I/O fsync hazard go to `projects/benchkit.md` — the traps to the
  cold-session notes, the hazard to the hostile-input checklist. The reproduction
  economics (only the real gate reproduces it; the disproven load shapes) go to
  `.bench/learnings.md`. The evidence table is not promoted (defaulted — veto point):
  it records how the decision was reached, which invariant 3 puts in git.
- **Story 9 rewrites the six surviving references in the same commit**: `ROADMAP.md`
  lines 26, 43, and 329, and `session-handoff.md` lines 24, 42, and 61. Deleting the
  file without these leaves six dangling references and nothing in the gate notices.
- **Story 9's top-tier line is a cached routing, not a fresh escalation.** The profile
  grants top+high standing for doc and guidance authoring; the no-standing-opt-out
  rule governs bumps *away* from a cached row. Flagged here so the reviewer can see
  the claim rather than infer it.

## Testing decisions

- What a good test is here: drive a real `bench` child through the shell wrapper and
  observe files, exit codes, and process state from outside — the runtime contract
  package's established shape. Two exceptions are taken deliberately and named below.
- Seams and prior art:
  - **`internal/contract/runtime`** for stories 2 and 3.
    `proveCancelledCommit` in `runtime_gate_action_proof_test.go` is direct prior art:
    it spawns a gate under `Setpgid`, waits on a marker, signals, and asserts the
    group is gone.
  - **`internal/gate`'s fake engine** for story 1's ordering. This is below the
    package boundary, taken because the ordering that makes the record worth having is
    not externally observable without fault injection the repo does not have. Prior
    art: `internal/gate/fault_engine_test.go`.
  - **The two-leg wait helper itself** for story 4, driven with a fake stat function
    and a fake exit channel. Prior art: none — this is the first test-support unit in
    the contract package with its own tests, taken because the alternative is grading
    the roadmap's headline arm by review alone.
  - **`internal/contract`'s fixture harness** for story 7, observed by a contract test
    that spawns a lingering grandchild.
  - **`internal/conformance`** for story 8: a unit test on the formatter, plus a
    source-level enumeration test over `package_core_checks_test.go` asserting no bare
    subprocess-failure diag literal survives. A fixture-bite test is *not* used:
    `checkGoCore` returns early without a `go.mod`, so biting it needs a real Go module
    fixture that would then also run `go vet`, `go test`, and the race probe inside the
    fixture — minutes of gate time for one assertion.
- Gate command: `bash bin/bench.sh gate`.
- **Load-window acceptance for the untriggerable residual.** The `.bench-contract-env`
  cleanup failure cannot be produced from inside the guest, so story 7's second row is
  validated the way the shipped deadline fix was: up to 3 `bash bin/bench.sh gate` runs
  while the reviewer spins host-side load, all green.

### Seam diagram

**Stories 1–3 — the gate owner record and signal teardown**

    trigger: `bash bin/bench.sh gate` (execs into `dist/bench gate-run`)
        │
        ▼
    root, git dir   ─▶  [ gate-run: Acquire lock       ]  ──▶  <git-dir>/bench-gate-owner (pid, start)
    SIGINT/TERM/HUP ─▶  [   write owner record (engine) ]  ──▶  pending verdict (durable, fsyncs)
                        [   run .bench/gate.sh (pgroup) ]  ──▶  gate stdout/stderr, exit code
                        [   on signal: cancel + kill pg ]  ──▶  owner record removed
                        [   on Acquire failure:         ]  ──▶  stderr: refusal + owner pid/liveness
              ◀ tests attach here (outside): a runtime contract test spawns a blocking
                gate, spawns a second and reads its stderr, signals the first and
                asserts the gate-script group is gone
              ◀ tests attach here (inside): the fake gateEngine records call order and
                asserts the owner write precedes the first durableReplace

**Stories 4–6 — the two-leg marker wait**

    trigger: contract subtest (R14 gate owner / losing-racer repair)
        │
        ▼
    fast marker path ─▶  [ two-leg wait helper          ]  ──▶  ok
    slow marker path ─▶  [  fast deadline: 5s / 10s     ]  ──▶  miss("fast", diagnostics)
    stat fn, exit ch ─▶  [  slow deadline: 60s          ]  ──▶  miss("slow", diagnostics)
              ◀ tests attach here: the helper is driven directly with a fake stat
                function and a fake exit channel — no child process needed

**Story 7 — the contract fixture harness**

    trigger: any contract test running a command through the fixture
        │
        ▼
    argv, env  ──▶  [ run in its own process group ]  ──▶  exit code, stdout, stderr
                    [ reap the group after Wait    ]  ──▶  no surviving grandchild
                    [ cleanup: remove .bench-contract-env, 3 tries / 2s ]
              ◀ tests attach here: a contract test runs a command that
                spawns a detached sleeper, then asserts nothing survives

**Story 8 — conformance diag attribution**

    trigger: gate conformance phase → TestRootConformance → core checks
        │
        ▼
    label, Probe  ──▶  [ one diag formatter          ]  ──▶  "<label> failed:\n<bounded, sanitized tail>"
              ◀ tests attach here: a unit test drives the formatter with a
                10 000-line probe containing ESC and BEL; a source-level
                enumeration test asserts no bare failure literal survives

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | The owner record write happens before the first `durableReplace` call | `internal/gate` fake engine | new test recording engine call order — red today, there is no owner write at all, and red against an implementation that writes after the pending record | This is the one assertion that separates the correct implementation from the degenerate one that writes the record behind the fsync it exists to precede |
| 1 | While a gate is running, `<git-dir>/bench-gate-owner` exists and names the live `gate-run` PID | `internal/contract/runtime` | new contract test against a blocking `.bench/gate.sh`: stat the record and check the PID is alive — red today, the file does not exist | The record is the only pre-fsync liveness signal; without the write the test finds no file |
| 1 | The owner record is gone after the gate run exits, green or red | `internal/contract/runtime` | same test asserts absence after the run — red against an implementation that writes but never removes | Excludes a write with no cleanup, which would make every later run report a phantom owner |
| 2 | A second `bench gate` blocked by a live owner prints the owning PID and its liveness alongside the refusal | `internal/contract/runtime` | new contract test: spawn a blocking gate, run a second one, assert its stderr contains the first child's PID — red today, the message names nothing | Without the record read the message cannot contain the PID at all |
| 2 | With the owner record absent or unparseable, the refusal still prints and still exits non-zero | `internal/contract/runtime` | same test with the record deleted mid-run — red against an implementation that errors on the diagnostic read | A diagnostic read that can fail the refusal turns benign contention into an operational failure |
| 3 | SIGTERM to `gate-run` leaves no member of the gate script's process group alive | `internal/contract/runtime` | new contract test modelled on `proveCancelledCommit`: record the gate script's PGID, SIGTERM the owner, assert `kill(-pgid, 0)` fails within the grace — red today, `gate-run` installs no handler | Without the handler Go's default terminate kills `gate-run` immediately and the group is orphaned, so the probe still succeeds |
| 3 | A signalled gate records no red verdict | `internal/contract/runtime` | same test asserts `gate.Inspect` reports pending, not red | Guards the existing cancellation semantics against a handler that falls through to the record path |
| 4 | A fast marker that never appears is reported as a fast-leg miss at the fast deadline, not the slow one | two-leg wait helper | new unit test with a fake stat function that never returns the fast marker; asserts the miss names the fast leg and returns within the fast deadline — red today, no helper exists | A single flat 60s wait with a reworded message returns at the wrong time and fails this assertion |
| 4 | A fast marker that appears promptly and a slow marker that appears at 30s both succeed | two-leg wait helper | same unit test with a scripted stat function — red against an implementation that applies the fast deadline to both markers | This is the load tolerance the shipped fix bought; an over-eager split would regress it |
| 4 | A child that exits before either marker is reported immediately, not at a deadline | two-leg wait helper | same unit test driving the fake exit channel — red against an implementation that only polls | Preserves the exit fast-fail the previous spec added to R14 |
| 5 | R14 uses the helper with the owner record as its fast leg and keeps the diagnostic dump on a miss | R14 helper | not TDD-able as an observed red — the helper reports via `t.Fatal` and the repo has no test-of-test-helper seam; graded by `/bench-review-implementation` against this row, with story 4's unit tests carrying the logic | Honest classification: the wiring is review-graded, but the behavior it wires is covered above, so the untested surface is one call site rather than the whole arm |
| 6 | The losing-racer test drives the start hook as a handshake — observe, then remove — so the child proceeds | losing-racer test body | the existing losing-racer subtest must stay green under the gate; wiring the hook without the removal hangs the child and reds it at 60s | The gate itself is the oracle for the mistake this decision exists to prevent |
| 7 | A contract command that spawns a detached grandchild leaves nothing alive once the runner returns | `internal/contract` fixture harness | new contract test spawning a sleeper through the fixture, then probing its PID — red today, the runner sets no `Setpgid` and calls plain `cmd.Run()` | Without the group reap the sleeper survives, which is the mechanism behind the `.bench-contract-env` cleanup failure |
| 7 | A subtest whose child is still writing into `.bench-contract-env` at teardown does not fail cleanup | `internal/contract` fixture harness | not TDD-able as an observed red — the original failure has never been reproduced on demand and did not fire in any run of the diagnosis session; graded by review plus the load-window acceptance in Testing decisions | Honest classification: the group reap above is the mechanism fix with a real red, and the retry is defence for a residual the repo cannot trigger |
| 8 | A failing probe's diag carries a bounded, control-byte-free tail of its output | `internal/conformance` | new unit test driving the formatter with a 10 000-line probe containing ESC and BEL — red today, no formatter exists | An unbounded or unsanitized tail buries the phase output and re-introduces the control-byte class the profile's checklist names |
| 8 | No bare subprocess-failure diag literal survives in the core checks | `internal/conformance` | new source-level enumeration test over `package_core_checks_test.go` covering all nine gate-reachable sites — red today, all nine are bare strings | The quantifier is "every diag"; enumerating the sites is what stops one being fixed and the rest left bare |
| 9 | `GATE-REPORT.md` is deleted and no reference to it survives in the tree | project gate | `rg 'GATE-REPORT' .` returns nothing — red today; this is a shell check, not a gate check, because no conformance sweep validates referenced file paths | Honest classification: the existing stale-reference sweep matches slash-command and adapter tokens only, so it cannot see this; review grades that the traps and hazard reached the profile |

### Edge inventory

- **Absent vs present-but-empty owner record** — story 2's second row; an empty file
  must read as "no owner named", not as PID 0.
- **Owner record naming a PID that has since been reused** — **Won't handle**: PID
  reuse cannot be excluded without a start-time comparison against `/proc`, and the
  record is diagnostic text, not an authorization input.
- **Control bytes in probe output** (ESC, BEL from a coloured toolchain) — story 8's
  first row.
- **File with no trailing newline** (the owner record hand-edited) — parse is
  whitespace-trimmed; covered by the unparseable-record row.
- **Interrupt mid-run leaving scratch state** — story 3's rows cover the process
  group and the pending verdict; the owner record's removal on the signal path is
  covered by story 1's absence row, which runs on every exit path.
- **Signal arriving before the lock is acquired** — **Won't handle**: nothing has been
  spawned and nothing recorded, so the Go default terminate is already correct.
- **Repeated signals (SIGINT twice)** — `signal.NotifyContext` restores the default
  handler after the first, so a second signal hard-kills; existing `gate-phases`
  behaviour, inherited unchanged.
- **Paths with spaces or glob characters in the git dir** — the owner record is
  written with `filepath.Join` and never through a shell; no new exposure.
- **Non-Linux hosts** — **Won't handle** beyond what exists: `Setpgid` and `fcntl`
  record locks are already the repo's Unix-only assumption.
- **A gate run under `bench shift`'s in-process path** (`RunAndRecordContext`) — the
  owner record is written by `Execute`, so the in-process caller gets it too; that
  caller owns its own teardown and must not be given a second signal handler.
- **Concurrent contract subtests reaping each other's groups** — story 7's reap is
  scoped to the group the runner itself created, never a discovered one.
- **`crossCompileMatrix` diags** — **Won't handle**: the gate's default build compiles
  in a `nil`-returning stub, so any row over it would grade nothing. Raise with the
  stress build if that ever joins the gate.
- **A diag site added after this spec lands** — story 8's enumeration test is the
  guard: a new bare failure literal in the core checks turns it red.

## Out of scope

- **A conformance check that validates referenced repository file paths exist** — the
  gate check story 9's red signal wanted and did not find. It is a separate capability
  (a documentation-integrity sweep over the whole tree, not a gate-trustworthiness
  fix), and it needs its own decisions about which reference syntaxes count —
  4 edits, 3 gate runs.
- **Core-count-aware gate/phase concurrency** — FT91's first arm and a separate
  capability (it changes how much work the gate schedules, not what it reports); the
  roadmap routes it to `/bench-shape-idea` because narrowing the oracle's parallelism
  is a reviewer decision — 3 edits, 3 gate runs.
- **Propagating SIGHUP to `gate-phases`** — if the reviewer accepts SIGHUP for
  `gate-run`, the same handler should widen there for consistency; held out only
  because it is contingent on that veto point — 1 edit, 1 gate run.
- **Confirming the fsync hypothesis with a goroutine dump** — the 60s deadlines absorb
  the stall, so the experiment is no longer reachable by this route; the permanent
  diagnostics already name the blocked line if a stall ever exceeds 60s — 0 edits.
- **Raising the sibling marker waits that have never flaked** (the 15s install-test
  deadline and others) — inherited unchanged from the retired
  `load-tolerant-marker-deadlines` spec; raise them if they ever fire — 1 edit,
  1 gate run.
