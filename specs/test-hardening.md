# test/hardening batch (FT53)

Status: implemented

## Problem

Five low-severity hardening findings from the 2026-07-08 platform assessment (§3, §4, §5, backlog item 12) sit in the shell/Go substrate that enforces the invariants. Individually each is friction or a latent false-pass; together they are the residue of two closed classes (the FT29 false-empty sweep, the assessment's "signals must recommend actions that work" theme). None is reachable as a user-visible bug today, but each is a defense-in-depth gap in code whose whole job is to be trustworthy when the happy path breaks:

- The `prepush.sh` ref-read loop drops a newline-less final line — a security-critical loop that trusts git always LF-terminates.
- `resolve_script_path` chases symlinks with no hop cap — a circular symlink to the wrapper hangs the CLI instead of erroring.
- `RegisteredWorktrees` discards the `git worktree list` classify error — a git failure reads as "no worktrees," so `bench worktree clean` prints "nothing to clean" at exit 0 on a broken repo. This is the last surviving instance of the false-empty class FT29 swept.
- The concurrent-acquire test's overlap comment overstates its barrier: it claims "a test-owned barrier, not a timed poll," but overlap *detection* is an event-keyed 60s deadline coupled to the shell's ~60s self-timeout, and the coupling is undocumented.
- (Assessment §3 / item 12) A missing pre-push hook in a linked repo should raise a `bench status` signal. **This one is already shipped** — see Implementation decisions.

## Solution

Close the four live gaps with the smallest diff each admits, every behavior change carrying a red-capable test, and correct the one comment. Confirm the fifth (the status signal) is already satisfied and cut it from the build.

This is a reviewer-directed batch drain standing in for a decision map: every decision below was pre-made by the orchestrator and is transcribed here, not reopened. Each is flagged **[defaulted — veto surface]** so the reviewer can reject it post-hoc.

## User stories

1. As a repo maintainer, I want the pre-push hook to process a newline-less final ref line, so that a hand-crafted or non-LF-terminated push input can never slip an unchecked ref (a direct-to-default push, or a gate-pin-drifting `.bench` tree) past the backstop. **[defaulted — veto surface: defense-in-depth; git LF-terminates today so this is unreachable via git.]** Line: claude-sonnet-5 / low. This is a one-token shell change at a known seam already driven by a contract suite, so the cheap tier fully covers it.

2. As an agent invoking `bench` through a symlink, I want `resolve_script_path` to cap its symlink hops (~40, readlink-style) and exit non-zero with a structured error on a cycle, so that a circular symlink to the wrapper fails fast instead of hanging the CLI forever. **[defaulted — veto surface: hop cap set to ~40 per readlink convention.]** Line: claude-sonnet-5 / medium. Mechanical shell at a known seam, but the timeout-guarded circular-symlink test needs care to prove the hang is gone, so medium effort.

3. As a maintainer running `bench worktree clean` (or reading `bench status`) in a repo where `git worktree list` fails, I want the command to surface the git error and exit non-zero, so that a git failure never masquerades as "nothing to clean" / a silent all-clear at exit 0. **[defaulted — veto surface: closes the last false-empty instance FT29 swept; all three callers (clean, status, dashboard) move off the swallowing accessor.]** Line: claude-sonnet-5 / medium. Go at a known seam with a fixed caller set, but the git-failure red signal's induction method is the one non-mechanical choice, so medium.

4. As the next reader of the concurrent-acquire contract test, I want its overlap comment to state honestly that release is barriered while overlap *detection* is an event-keyed 60s deadline coupled to the shell's ~60s self-timeout, so that I don't trust a barrier the test doesn't provide. **[defaulted — veto surface: comment-only; documents the deadline/self-timeout coupling.]** Line: claude-sonnet-5 / low. A comment correction with no behavior change, so the cheap tier at low effort suffices.

5. As a maintainer of a freshly cloned linked repo, I want `bench status` to raise a signal when the pre-push backstop is missing or unmanaged, with a remedy that actually reinstalls it, so that I don't silently lose the default-branch guard. **[ALREADY SHIPPED — see Implementation decisions; recommend cutting from the build.]** Line (were it built): claude-opus-4-8 / medium, per the profile's gate/conformance cached row, because the row touches the gate-tested status contract. This is here only as veto surface: the assessment requested it, but the tree already satisfies it.

## Implementation decisions

**Story 1 — `internal/adopt/prepush.sh`.** Add the `|| [ -n "$line" ]` tail to the `while IFS= read -r line` loop (~line 36) so a final line with no trailing newline is still appended to `ref_lines` and checked. `read` returns non-zero at EOF but leaves the partial line in `$line`; the tail admits that last iteration. No change to the check bodies.

**Story 2 — `bin/bench.sh` `resolve_script_path`** (~lines 35–44). Add a hop counter to the `while [[ -L "$source" ]]` loop; on exceeding ~40 hops, print a structured error to stderr (naming the offending path and the cause: "symlink cycle / too many levels") and exit non-zero rather than looping. Match the readlink-`-f` convention (bounded hops → ELOOP-style failure). Keep the GNU-free portability the function was written for — no dependence on `readlink -f`.

**Story 3 — `internal/worktree/classifier.go` and its three callers.** `RegisteredWorktrees` (line 25) currently discards the error from `ClassifyRegisteredWorktrees`. The erroring variant already exists and returns the git error correctly; the change is that callers stop swallowing it:
- `internal/worktree/clean.go:92` (`cleanOutOfPoolWorktrees`) — on classify failure, return non-zero with an error to stderr; never fall through to "nothing to clean" (line 135) at exit 0.
- `internal/status/status.go:269` (`appendWorktree`) — on classify failure, surface it rather than silently emitting no worktree row (fail visible, consistent with the false-empty rule).
- `internal/dashboard/dashboard.go:91` — handle the error rather than rendering an empty worktree pane as truth.

The swallowing `RegisteredWorktrees` wrapper is deleted once no caller remains (contract-and-delete). The unit test `TestClassifyRegisteredWorktrees` in `worktree_test.go` is the prior seam.

*Open decision within story 3 [defaulted — veto surface]:* the git-failure red signal's induction. Recommended: a unit test at the caller's inner function (the FT29/`structure_test.go` `gitOpError` pattern — inject the error deterministically), cheaper and more deterministic than a PATH-shimmed `git` that exits non-zero for `worktree list` at the contract level. Reviewer confirms the induction.

**Story 4 — `internal/contract/runtime/runtime_worktree_test.go`** (~lines 375–380). Reword the `testRuntimeWorktreeConcurrentAcquire` header comment: state that the release is barriered (the shell holds until the test drops the go-file) but overlap *detection* is an event-keyed 60s deadline (`overlapDeadline`, line 412), and document that this window is coupled to the shell's own ~60s self-timeout — the two must not converge. Comment only; no code change.

**Story 5 — already implemented; recommend cut.** `internal/status/status.go` already carries `appendGuards` (lines 322–353): a sev-3 `guards` row firing on `PrePushAbsent`/`PrePushForeign`/`PrePushDiverted`, gated to the primary checkout of a routed repo (`.bench/lines.env` present), quiet when bench-managed, remedy `bench link` (re-running link reinstalls the hook — the effective, non-no-op action the assessment's theme-3 demands; doctor only *reports* via `reportPrePush`, so `link` is the correct owner). It fits the five-row budget and the severity ladder (worktree < guards < drain), and is gate-tested by the **guards-signal contract family** in `internal/contract/runtime/runtime_status_test.go` (`testRuntimeStatusGuardsSignal` pins the row + ladder + managed-clears-it + unrouted-quiet; `testRuntimeStatusGuardsPrimaryOnly` pins primary-checkout-only). It shipped in commit `5fc4ed1` (2026-07-07), one day before the assessment that requested it — the §3-low/item-12 finding is stale. **No work remains; the orchestrator's item-5 decision is already the tree's behavior.** Left in scope only as veto surface: reviewer confirms the cut.

## Testing decisions

- **A good test here** drives the real shell/Go surface at its existing contract seam and observes the externally visible effect (hook exit code, CLI exit code + message, ladder ordering) — not internals.
- **Seams and prior art:**
  - Story 1 → `internal/contract/surface/prepush_test.go`. `runPrePush(t, f, <stdin>)` drives the installed hook with crafted stdin; `refLine(...)` builds ref lines. A newline-less final-line variant is the new input.
  - Story 2 → a surface/runtime contract that builds a circular symlink to the wrapper and invokes it **under a timeout**, asserting non-zero exit + the structured error rather than a hang. Prior art: the runtime contract suites run `dist/bench` in throwaway fixtures.
  - Story 3 → `internal/worktree/worktree_test.go` (`TestClassifyRegisteredWorktrees`) for the classifier propagation, plus a caller-level assertion that `bench worktree clean` exits non-zero on an induced git failure. Prior art: FT29's `structure_test.go` `gitOpError` injection.
  - Story 4 → not TDD-able (comment only).
  - Story 5 → already covered by the guards-signal family (above); no new test.
- **Gate command:** `.bench/gate.sh` (the project gate). The compiled-core and runtime/surface fragments must stay green; runtime/surface contracts run the freshly rebuilt `dist/bench`, so rebuild before hand-running them.

### Seam diagram

Story 1 — pre-push newline-less final ref line:

    trigger: git push (or runPrePush contract helper) writes ref lines to hook stdin
        │
        ▼
    "…drifted\n…clean"  ──▶  [ prepush.sh while-read loop (ref_lines[]) ]  ──▶  exit 0/1 + stderr
     (no trailing \n)   ──▶  [   + drift/protected-branch checks       ]
                                  ◀ tests attach here: prepush_test.go runPrePush feeds
                                    a final line with no "\n"; assert the last ref is still checked

Story 2 — symlink-hop cap in the wrapper resolver:

    trigger: invoking `bench` through a symlink (or a cyclic symlink to it)
        │
        ▼
    BASH_SOURCE  ──▶  [ resolve_script_path while [[ -L ]] loop + hop cap ]  ──▶  real path
     (cyclic)    ──▶  [                                                   ]  ──▶  err + exit≠0
                          ◀ tests attach here: create a self-referential symlink, invoke under
                            a timeout; assert structured error + non-zero, not a hang

Story 3 — classify error surfaced by callers:

    trigger: bench worktree clean / bench status / bench dashboard, git worktree list fails
        │
        ▼
    root  ──▶  [ ClassifyRegisteredWorktrees → (registered, err) ]  ──▶  callers
              [   callers propagate err, not swallow             ]  ──▶  clean: err + exit≠0
                  ◀ tests attach here: induce a git worktree-list failure; assert the caller
                    errors (exit≠0 / visible), never prints "nothing to clean" at exit 0

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | newline-less final ref line is still checked (drift or protected-branch push blocked) | `prepush_test.go` `runPrePush` | red: a new `runPrePush` case feeding a final ref line with no trailing `\n` (a drifted-tree or `refs/heads/main` line) exits 0 against the current loop, because the unterminated line is dropped | the degenerate loop (without the newline-tail guard) skips the last iteration, so the offending ref is never checked; the assertion demands exit≠0 |
| 2 | a circular symlink to the wrapper errors and exits non-zero, not hangs | new runtime/surface contract, run under a timeout | red: a self-referential-symlink invocation of the current wrapper never returns within the timeout (test times out) | the uncapped `while [[ -L ]]` loop never terminates on a cycle; a hop cap forces ELOOP-style exit, which the timeout-bounded assertion observes as a clean non-zero |
| 3 | `bench worktree clean` on a `git worktree list` failure exits non-zero with an error, not "nothing to clean" at exit 0 | `worktree_test.go` classifier test + a caller-level assertion (induction per open decision) | red: with the git query induced to fail, the current caller (swallowing `RegisteredWorktrees`) prints "nothing to clean" and exits 0 | swallowing the classify error makes an empty list indistinguishable from a real failure; propagating it lets the caller fail visibly, which the assertion requires |
| 4 | overlap comment states the event-keyed deadline + self-timeout coupling honestly | `runtime_worktree_test.go` header comment | not TDD-able (comment only) — no behavior changes; correctness is read at review, not by the gate | — |
| 5 | `bench status` raises a guards row (remedy `bench link`) when the pre-push hook is missing/unmanaged in a routed primary checkout | `runtime_status_test.go` guards-signal family | already covered — `testRuntimeStatusGuardsSignal` + `testRuntimeStatusGuardsPrimaryOnly` already pin the row, its remedy, the ladder position, and primary-checkout gating; no red signal available because the behavior already exists | the existing tests already go red on removal of `appendGuards`; the behavior is shipped, so this row documents coverage rather than driving new work |

**Degenerate-implementation check.** Story 1's always-drop loop, story 2's uncapped loop, and story 3's swallow-the-error accessor are each exactly the current (wrong-for-hostile-input) implementation, and each named red signal fails on it — the map pins behavior, not belief. Stories 4 and 5 change no behavior.

### Edge inventory

Walked against the profile's shell-CLI hostile-input checklist; several classes are literally this batch's subject.

- **hand-edited / non-LF-terminated last line** → story 1 coverage row (the subject).
- **invocation through a symlink** → story 2 coverage row (the cyclic case; the ordinary symlink case is already exercised by the existing `resolve_script_path` walk).
- **required tool missing from PATH (no `readlink -f`)** → story 2 keeps the GNU-free chase; the cap adds no `readlink -f` dependency, so this class stays satisfied. **Won't handle** beyond that: a missing `git` on PATH is out of this batch — it is the CLI-wide bootstrap concern, not one of these five seams.
- **absent vs present-but-empty input** → story 3 makes "git failed" (error) distinct from "no worktrees" (empty list at exit 0); the two are now separate observable behaviors.
- **invocation through every shipped surface** → story 3 fixes all three callers (clean, status, dashboard) so the same classify error is surfaced identically through each; story 1's fix is in the installed hook, reached by real git push and by the `runPrePush` contract driver alike.
- **control bytes in git-sourced text** — **Won't handle**: none of these five seams renders git text through `toon.Table`; the class is owned elsewhere and unchanged here.
- **paths/dir names with spaces or glob chars** — **Won't handle** as a new case: the story-2 fixture and story-3 callers inherit the existing path-quoting the wrapper and Go callers already use; this batch adds no new unquoted expansion.
- **unquoted multi-word arguments (`$*` vs `$1`)** — **Won't handle**: no argument parsing changes in any of the five.
- **interrupt (SIGINT) mid-loop; re-run idempotency; cwd deeper than repo root** — **Won't handle**: none of the five touches loop leases, scratch state, or root-relative cwd assumptions; the pre-push and resolver changes are idempotent by construction.

## Out of scope

- **Story 5 (status pre-push signal) as build work** — already shipped in `5fc4ed1` and gate-tested by the guards-signal family; retained in the stories only as veto surface. Cost to (re)build were it absent: `2 edits, 1 gate run`. Recommend cut.
- **The reclaim race (§4 med) and the salvage-orphan loop (§4 med)** — separate capabilities with their own assessment backlog rows (items 4 and 3), each needing its own concurrency/UX decision, not the rest of this hardening batch. Reclaim race: `~4 edits, 2 gate runs`. Salvage loop: `~3 edits, 2 gate runs`.
- **Pre-push default-branch honesty (§3 med, item 8)** — the fabricated-`main` fallback is partly addressed already (the hook now resolves `origin/HEAD` live) but the link-time warning + doctor comparison row is a distinct feature. `~3 edits, 2 gate runs`.
- **The one-source collapse and CLI-hygiene batches (items 9, 10)** — landed separately (FT50 for the collapse; CLI hygiene is its own row); not this batch's subject.
