# Lossless shift recovery and truthful result states (FT79)

## Destination

Every shift failure after agent mutation preserves the work — a durable recovery
ref by default, a locked worktree as fallback — and prints its exact location.
Staging, adapter, commit, signal, and teardown errors propagate instead of
vanishing into exit 1 or false success. Shift requires an objective, validates a
positive bounded iteration cap, runs under a wall deadline, and distinguishes
complete / incomplete / failed / interrupted / no-op with honest exit codes and a
structured result. Sources: `RR:R-03`, `RR:R-08`, `RC:H-05`.

## #1: What preserves post-mutation work — worktree or ref?

Type: Grill

### Question
`worktree.Release` hard-resets and cleans, so any failure after agent mutation
(commit, staging, teardown) erases uncommitted work today. Preserve the charged
worktree, or snapshot to a durable ref?

### Answer
**Ref primary, worktree fallback.** Snapshot the dirty tree (scratch excluded)
with plumbing — `write-tree`/`commit-tree` under a temp index, explicit synthetic
identity so missing Git identity cannot block the snapshot, parented on the shift
branch tip — to `refs/bench/recovery/<shift-branch>`, then release the pool
worktree normally. Only when the snapshot itself fails does the shift retain the
worktree, `git worktree lock` it with a reason, and print its path. Committed
iterations already survive on the shift branch. Keeps the pool usable and reuses
the existing recovery-ref envelope machinery. Rejected: always-retain (pool
starvation), both-always (cost without added recoverability).

## #2: What does adapter failure mean?

Type: Grill

### Question
Today adapter start/exit failures are ignored ("the gate decides"), so
`/bin/false` plus a green gate exits 0 with "objective likely met". What do they
mean now?

### Answer
**Stop the loop; split the outcome by evidence.** Spawn failure or nonzero exit
stops the shift after processing any mutations (gate → commit on green →
otherwise snapshot). The "no change = objective met" inference now additionally
requires adapter exit 0. Outcome: `failed` when zero iterations committed this
shift, `incomplete` when prior green commits exist. The gate stays the oracle for
*work*; the adapter exit becomes evidence for *progress*. Rejected: always
`failed` (understates committed work), retry-once (retry policy muddies the
evidence trail).

## #3: Outcome taxonomy and exit codes

Type: Grill

### Question
Exit codes are 0/1/130 today, and cap exhaustion exits 0. What code does each
outcome get?

### Answer
**0/1/2/3/4/130.** `0` complete · `1` failed · `2` usage/preflight (nothing ran)
· `3` incomplete (cap or wall exhausted, adapter failure after commits) · `4`
no-op · `130` interrupted. AXI-consistent for 0/1/2; distinct codes keep
scripting honest. Rejected: folding incomplete/no-op into 1; exit 0 on cap
exhaustion.

## #4: Where do the structured result and recovery pointer live?

Type: Grill

### Question
"Structured results" and "prints its exact location" need a surface that
survives the terminal. Which?

### Answer
**TOON `shift_result` block on stdout + the intent ledger.** The final block
carries outcome, exit code, branch, committed count, iterations used, recovery
reference (ref, retained worktree path, or none), and detail — the kit's TOON
emitter, control-byte-safe by construction. The same outcome + recovery pointer
is recorded on the shift's existing intent-ledger entry, so `bench status`'s
intent signal is the after-the-fact discovery surface. The full append-only
event schema is FT71, not this. Rejected: a new `<git-dir>` result file
(duplicates FT71), stdout-only (a closed terminal loses the pointer).

## #5: Wall deadline

Type: Grill

### Question
No wall deadline exists. Knob, default, and expiry behavior?

### Answer
**`BENCH_MAX_WALL`, Go duration, default 2h, kills in-flight.** Validated
positive and ≤ 24h. Expiry acts like a pulled line: kill the adapter process
group, cancel a running gate, snapshot mutations, exit 3 (incomplete, deadline
detail). Rejected: checkpoint-only checks (a hung adapter means no wall at all),
opt-in-only (the roadmap asks for a property, not an option).

## #6: Objective and cap validation

Type: Grill

### Question
A missing objective silently becomes "improve the codebase"; `envInt` accepts
invalid, zero, negative, and unbounded caps. Posture?

### Answer
**Require the objective; error on explicitly-set invalid values.** Empty
objective → exit 2 (the default dies); an objective carrying control bytes →
exit 2 at entry (TOON refuses them later anyway). Env knobs: unset or
empty-string keeps defaults (`BENCH_MAX_ITERS` 12, `BENCH_REFACTOR_ITERS` 4,
`BENCH_MAX_WALL` 2h); set-but-invalid (non-integer, < 1, over bound) exits 2
naming the variable and accepted range, before any worktree acquire. Bounds: 100
iterations, 24h wall. Rejected: clamp-and-warn (executes something the user
didn't ask for), unbounded positives ("bounded" is in the roadmap row verbatim).

## #7: The no-op boundary

Type: Grill

### Question
Zero commits, adapter exit 0, and `.bench/done.sh` passes on iteration 1 —
complete or no-op?

### Answer
**No-op (exit 4), detail notes the predicate.** Exit 4 whenever zero iterations
committed, even if done.sh passes; the result detail says the predicate was
already satisfied. Exit 0 always means work landed and the oracle blessed it.
Rejected: complete-on-predicate (scripts can't distinguish "already done" from
"did something").

## #8: Recovery-ref lifecycle, and the red-gate path

Type: Grill

### Question
Where do recovery refs live, who consumes them, and does the existing
red-gate preserve-worktree path survive?

### Answer
One rule, no special case: **every post-mutation failure — including a red-gate
iteration — snapshots and releases.** The red-gate path's retain-worktree
behavior is replaced; inspection uses the printed hint (`git worktree add`
or `git diff <branch>..<ref>` against the recovery ref). Refs live at
`refs/bench/recovery/<shift-branch>`, one per shift; creation fails closed on
conflict (the existing envelope conflict check — a timestamp-collision branch
name must not overwrite another shift's evidence). Consumption is manual via
the printed hint; `bench status` surfaces the pointer from the intent ledger;
`bench resume` and the worktree classifier must *retain, never delete* recovery
refs and locked fallback worktrees. Garbage collection of consumed refs is out
of scope (fog).

## #9: Error propagation — staging, refactor, teardown

Type: Grill

### Question
`stageTouched` ignores `git add` failures; a refactor commit failure returns 1
into the erase path; `teardown`/`Release` is a silent void. Posture per class?

### Answer
**All three propagate under the matrix rule.** A staging or commit failure (main
loop or refactor phase) snapshots and exits by the evidence split (`failed` with
zero commits, `incomplete` after commits). A teardown error exits 1 with the
result detail naming what is safe (the branch, any recovery ref) — the failure
is real even when the work is. The refactor phase's red-gate *rollback* stays
discard-by-design: a red refactor probe is not successful work.

## The recovery matrix (resolved under #1–#9)

| Failure | Behavior | Outcome / exit |
|---|---|---|
| missing Git identity at commit | snapshot (synthetic identity), release | failed 1 / incomplete 3 |
| failing commit hook | same | failed 1 / incomplete 3 |
| staging failure | same | failed 1 / incomplete 3 |
| adapter start/exit failure | process mutations, stop loop | failed 1 / incomplete 3 |
| red-gate iteration | snapshot, release | failed 1 / incomplete 3 |
| signal interruption | kill child, snapshot mutations, release | interrupted 130 |
| wall deadline | like signal, deadline detail | incomplete 3 |
| cap exhaustion | work stays on branch | incomplete 3 |
| no-op iteration (adapter clean) | nothing to preserve | no-op 4 |
| teardown error | report; branch/ref already safe | failed 1 |
| snapshot failure (any row above) | retain + lock worktree, print path | row's outcome |

## Not yet specified

- Garbage collection / retirement of consumed recovery refs.
- FT71's event schema will re-carry the outcome + recovery fields; alignment is
  FT71's problem.
- Adapter authors' contract doc for "exit 0 = clean" normalization in the shims.

## Out of scope

- FT71 versioned shift evidence (append-only event log).
- FT88 environment passlists and prompt-via-stdin.
- FT58 identity-safe lease reclamation (pool lock protocol beyond what recovery
  ordering needs).
- Locked assignment worktrees and their lifecycle — existing assignments stay
  untouched.
- Any gate-execution or verdict-cache change (R-02 territory).

## Handoff

1. **Module boundaries.**
   - `internal/shift` (deep) — outcome taxonomy and exit codes, objective/cap/wall
     validation, the matrix orchestration (when to snapshot, when to stop), the
     TOON `shift_result` emission, wall-deadline timer wired into the existing
     signal/checkpoint machinery.
   - `internal/worktree` — the snapshot primitive (temp-index `write-tree` +
     `commit-tree` with explicit synthetic identity + fail-closed ref creation,
     reusing the recovery-envelope machinery) and the retain-and-lock fallback;
     release ordering guarantees.
   - `internal/intent` — outcome + recovery fields on shift entries.
   - `internal/status` — renders the recovery pointer on the intent signal.
   - `cmd/bench/main.go`, adapters, hooks — unchanged shape (thin).

2. **Contracts.**
   - CLI: `bench shift <objective…>`; empty or control-byte objective → exit 2.
     Env: `BENCH_MAX_ITERS` int 1–100 default 12; `BENCH_REFACTOR_ITERS` int
     1–100 default 4; `BENCH_MAX_WALL` Go duration (0, 24h] default 2h; empty
     string = unset; invalid → exit 2 naming variable and range, before acquire.
   - Exit codes: 0 complete, 1 failed, 2 usage/preflight, 3 incomplete, 4 no-op,
     130 interrupted — per the matrix above.
   - `shift_result` TOON block: outcome, exit, branch, committed,
     iterations_used, recovery (`ref:<name>` | `worktree:<path>` | `none`),
     detail.
   - Recovery ref: `refs/bench/recovery/<shift-branch>`; snapshot commit of the
     full dirty tree (scratch excluded), parent = branch tip; hooks and identity
     config bypassed by plumbing construction; creation conflicts fail closed.
   - Intent entry carries outcome + recovery pointer; `bench status` renders it.

3. **Deep vs thin.** `internal/shift`'s loop and `internal/worktree`'s snapshot
   primitive are the two deep modules; the seam between them is "snapshot this
   root onto this ref, tell me what you preserved". `Command` arg parsing,
   adapters, and status rendering are thin pass-throughs.

4. **Black-box assertables.** Per matrix row: exit code; `shift_result` fields on
   stdout; `git show-ref` proves the recovery ref and `git ls-tree`/`cat-file`
   proves the mutated file is in its tree; branch commit count; intent-ledger
   JSON content; pool lease present/absent; worktree locked (fallback) or gone.

5. **Gate attachment.** Runtime contract fragments
   (`internal/contract/runtime/runtime_shift_test.go`, built-binary harness)
   with scripted adapters: `/bin/false`; an adapter that writes a file then
   exits nonzero; identity-stripped env plus a failing `pre-commit`/`commit-msg`
   hook fixture; a sleeping adapter under a tiny `BENCH_MAX_WALL`; an
   always-dirty adapter to exhaust the cap; SIGTERM mid-adapter. Two rows the
   shell contracts cannot reach — staging failure and teardown error — attach at
   in-process fault seams (the worktree package's existing `Fault` pattern) as
   Go unit tests: TDD there. Manual verify: none flagged; every row is
   observable black-box or via a fault seam.

6. **Hostile-input owners** (profile checklist → seam):
   - spaces/glob/newline paths in mutations → existing `ParsePorcelainZ` +
     `:(literal)`; the snapshot's temp-index `add -A` inherits whole-tree safety.
   - control bytes in git-sourced/objective text → CLI entry validation (exit 2)
     plus `toon.Table` refusal as backstop.
   - absent vs empty env var → empty = unset, both asserted.
   - required tool missing / bad `BENCH_AGENT` → existing preflight, now exit 2.
   - interrupt mid-loop → snapshot-then-release checkpoint path, exit 130.
   - destructive worktree state → existing fail-closed classifier; recovery refs
     and locked fallback worktrees classify as retained, never residue.
   - re-run idempotency → ref-creation conflict fails closed; a re-run shift
     gets a fresh branch/ref pair.
   - cwd deeper than root / symlink invocation → existing `git.Root` routing.

7. **Uncertainty flags.**
   - Snapshot mechanics: temp-index (`GIT_INDEX_FILE`) `add -A` → `write-tree` →
     `commit-tree` is the recommended shape; the spec-writer should probe that
     untracked files land in the tree in a *worktree* checkout before locking
     the coverage row (`git stash create` was rejected on untracked handling,
     verify the replacement actually covers it).
   - `bench resume` / classifier interplay: confirm the classifier treats a
     locked fallback worktree as `ReasonUnexpectedLock`-style retained today, or
     name the change needed.
   - Adapter exit-code semantics are harness-specific (a Stop-hook-blocked stop
     may surface as nonzero in some harnesses); the shims own normalization —
     confirm each shipped adapter exits 0 on a clean iteration.

8. **Rejected alternatives.** Always-retain worktree; snapshot-and-retain both;
   adapter retry-once; adapter failure always `failed`; folding incomplete/no-op
   into exit 1; exit 0 on cap exhaustion; complete-on-predicate with zero
   commits; clamp-and-warn caps; unbounded caps; checkpoint-only or opt-in-only
   deadline; `<git-dir>` result file; stdout-only recovery pointer; keeping the
   red-gate retain-worktree special case. Do not reopen.

9. **Domain watch-outs.**
   - `worktree.Release` hard-resets and cleans; the snapshot must be created
     *and verified* before any release path can run — ordering is load-bearing,
     and the fallback retain-and-lock is the only legal response to a snapshot
     failure.
   - `checkpoint()` exits via `os.Exit(130)`, which skips deferred cleanup —
     every exit path must run snapshot/teardown explicitly, never rely on defer.
   - Plumbing commits bypass hooks and identity config: snapshots succeed
     exactly where iteration commits fail. That asymmetry is the design.
   - Worktrees share the repo's ref namespace — a recovery ref is visible
     repo-wide the moment it exists, and survives worktree pruning.
   - Timestamp-derived shift-branch names can collide within one second across
     concurrent shifts; ref creation failing closed is what protects the first
     shift's evidence.

Dependency order: recommended as one spec; if sliced — (A) validation + outcome
taxonomy + exit codes + `shift_result` emission (no recovery machinery); (B) the
snapshot primitive + matrix wiring + intent/status surfacing. B depends on A.
Slicing stays the reviewer's call.
