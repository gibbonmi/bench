# Worktree lifecycle reaches released by default (FT95)

## Destination

`specs/worktree-lifecycle.md`. Reviewer-prioritized 2026-07-15 ("major usability
issue"); all three forks were put to the reviewer directly and closed the same
day, so this map records decisions rather than opening them.

## #1: What may cleanup discard automatically when the tracked tree is clean?

**Decided: declared build outputs only.** A repo-owned declaration
(`.bench/build-outputs.json`, `gate-inputs.json` schema style) names ignored
paths that are rebuildable build output (this repo: `dist/`). An ignored
inventory that lies entirely within the declared set no longer blocks release
or automatic cleanup — removing the tree removes them, and the receipt's
existing ignored-inventory evidence records what went. Any ignored path
outside the set retains exactly as today. A truncated or over-limit inventory
cannot prove containment and retains. A missing declaration means an empty
set (today's behavior); a malformed or hostile one (traversal, absolute
paths, control bytes) fails closed to retain with an explicit reason.
Rejected: discard-any-ignored (deletes non-rebuildable ignored files —
secrets, notes — in consumer repos) and keep-retain-default (leaves the
manual `--discard-ignored` step that causes the accumulation).

## #2: How does resume treat an `active` record whose owning session is dead?

**Decided: probe the recorded owner, then reconcile.** The lease file already
records `<pid> <utc-time>` and `internal/worktree/lifecycle.go` already owns
dead-PID reclaim logic; the cleanup planner today never reads it — it retains
on lease-file *presence* (`subshell.go`), which is why a dead session's tree
is retained at every session start. The planner instead consults one
lifecycle-owned probe: **live** → retain (unchanged, FT58's live-owner rule);
**dead** → the lease no longer blocks, and the ordinary classification
(landed, tracked-clean, ignored ⊆ declared) decides; **unknown** (unreadable,
garbage, permission error) → retain. PID recycling can only produce a
false *live*, which retains — never destructive. Rejected: age thresholds
(nags long-lived live sessions, waits days on dead ones) and folding into
FT58 (blocks an ASAP fix on a MEDIUM row; FT58's lock-protocol scope is
untouched).

## #3: What is the list surface?

**Decided: a new `bench worktree list` subcommand; bare invocation keeps its
subshell meaning.** AXI/TOON read-only query: one row per assignment record
(id, label, state, tree present/missing, lease live/dead/none/unknown,
landed, ignored count) plus rows for registered non-root worktrees with no
assignment (source `foreign`). Self-describing empty table, structured
errors, exit 0/1/2, help exits 0. Rejected: bare-becomes-list (breaks the
documented subshell flow).

## Handoff

1. **Module boundaries.**
   - `internal/worktree` (deep) — the cleanup planner's lease-liveness
     consultation and declared-build-output containment check (one classify
     path serves explicit release, explicit clean, and `PlanAutomatic`, so
     resume inherits both fixes); the new list query. No new destructive
     machinery: removal and discard ride the existing plan/apply transaction.
   - `internal/worktree/lifecycle.go` — single owner of the liveness probe
     (live/dead/unknown), extending the existing lease/reclaim logic; the
     planner consumes it, never re-derives it.
   - `.bench/build-outputs.json` — new repo-owned declaration, one source for
     "what is rebuildable build output"; parser lives beside the planner.
   - `cmd/bench/main.go`, `bin/bench.sh` — thin: `list` dispatch and usage.
   - `docs/adr/0005` — amended: non-live is determined by the recorded owner
     identity, not lease-file presence; declared build outputs are
     discardable residue.
2. **Contracts.**
   - Planner: retain reasons keep their codes; a dead lease adds no retain; an
     ignored inventory with every path under a declared entry adds no retain;
     `OverLimit`/`AtLeast` inventories retain; unknown lease state retains.
   - Config: `{"schema":1,"paths":["dist/"]}`; entries are repo-relative,
     no `..`, no leading `/`, no control bytes; a directory entry covers
     everything under it; violations → the file is malformed → retain with an
     explicit malformed-declaration reason.
   - CLI: `bench worktree list` — TOON `worktrees[N]{id,label,state,source,tree,lease,landed,ignored}`;
     empty table at zero; unknown args exit 2 with usage; `-h/--help` exit 0;
     errors structured on stdout, exit 1.
   - `bench resume` summary format unchanged (dispositions shift between
     existing counters; no golden edits).
3. **Deep vs thin.** The planner and the liveness probe are the deep pair; the
   seam between them is "is this lease live, dead, or unknown". List, config
   parse, dispatch, and usage text are thin.
4. **Black-box assertables.** Release exit code and receipt action; assignment
   record presence/absence in the ledger; the worktree directory's existence;
   recovery ref untouched where preserved; `ResumeResult` counters; `worktree
   list` TOON rows and exit codes.
5. **Gate attachment.** `internal/worktree` package tests (fixture repos;
   prior art: `newOwnedAssignment`, `commitInWorktree`, `requireTest`,
   `ApplyExplicitWithOptions`, and the lease fixtures in the lifecycle tests)
   for every planner behavior; one runtime contract fragment
   (`internal/contract/runtime`, built binary) for the `list` CLI contract and
   dispatch. The project gate `.bench/gate.sh` is the oracle.
6. **Hostile-input owners.** Config parser owns traversal/absolute/control-byte
   entries and absent-vs-empty. The liveness probe owns garbage lease content,
   unreadable lease, and PID recycling (false-live retains). The planner owns
   truncated inventories. TOON emitter owns control bytes in list output
   (refusal, existing behavior).
7. **Uncertainty flags.** None — all three forks are reviewer-closed and the
   mechanisms extend code that already exists.
