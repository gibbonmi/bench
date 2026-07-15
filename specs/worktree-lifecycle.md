# Worktree lifecycle reaches released by default

Status: implemented

## Problem

Worktrees accumulate structurally. Running the gate in any worktree creates
`dist/` (gitignored), and the cleanup planner retains on *any* ignored file, so
every normally used tree needs a manual `clean --discard-ignored` before it can
fully release. A session that ends without releasing leaves a lease file behind,
and the planner retains on lease-file *presence* without ever probing the
recorded owner PID, so a dead session's tree is re-retained at every session
start, forever. And the population is invisible between session starts: bare
`bench worktree` creates an assignment subshell rather than listing anything.
Live evidence 2026-07-15: three stale trees, none doing work, two of them fully
landed and blocked only by `dist/` or a dead lease.

## Solution

Three classification fixes on existing machinery — no new destructive paths.
The planner consults the lifecycle-owned liveness probe (live → retain, dead →
the lease stops blocking, unknown → retain); an ignored inventory lying
entirely under a repo-declared build-output set (`.bench/build-outputs.json`)
stops blocking release and automatic cleanup; and a new read-only
`bench worktree list` makes the population visible. Because explicit release,
explicit clean, and `PlanAutomatic` share one classify path, `bench resume`
inherits both fixes: a dead session's landed, clean tree now reconciles at
session start, and live or unproven-dead sessions retain exactly as today.

Decisions compiled from `decisions/worktree-lifecycle.md` (all three forks
reviewer-closed 2026-07-15).

## User stories

1. As a session, I want `release --request` on a landed, tracked-clean tree
   whose lease records a dead PID to release and compact the record, so a
   crashed or abandoned session stops retaining its tree forever.
   Line: sonnet / high. Gate-observable planner logic at a known seam, but it
   relaxes a destructive-safety guard, so the cheap tier gets high effort.
2. As a reviewer, I want a lease whose recorded PID is alive to retain exactly
   as today, so a live session is never released out from under its owner.
   Line: sonnet / high. This is the FT58 live-owner rule and the safety half of
   story 1, proven at the same seam.
3. As a reviewer, I want an unreadable or garbage lease to retain, so liveness
   is only ever relaxed on positive proof of a dead owner.
   Line: sonnet / high. Fail-closed twin of story 1 at the same seam.
4. As a session, I want a tree whose only ignored paths fall under the declared
   build-output set to release (and auto-clean) without `--discard-ignored`,
   so gate build output stops forcing a manual step on every used worktree.
   Line: sonnet / high. Containment logic is small, but it widens a destructive
   default, so it gets the same treatment as story 1.
5. As a reviewer, I want any ignored path outside the declared set — or an
   over-limit/truncated inventory that cannot prove containment — to retain as
   today, so unknown ignored files (secrets, notes) are never deleted by
   default. Line: sonnet / high. The fail-closed twin of story 4.
6. As a reviewer, I want a missing declaration to mean an empty set and a
   malformed or hostile one (traversal, absolute path, control bytes, bad
   schema) to retain with an explicit malformed-declaration reason, so the
   config can only narrow safety by being absent, never by being wrong.
   Line: sonnet / medium. Parser with exact fixtures; fully gate-observable.
7. As a session, I want `bench resume` at session start to reconcile a dead
   session's landed, clean tree (counted in the existing summary counters,
   format unchanged) while still retaining live-lease trees and reporting
   preserved unlanded work, so the population stops growing without any new
   risk to live or unmerged work. Line: sonnet / high. Resume inherits the
   planner change; the test proves the inheritance end to end.
8. As a reviewer, I want `bench worktree list` to show every assignment (id,
   label, state, tree present/missing, lease live/dead/none/unknown, landed,
   ignored count) plus foreign registered worktrees, with a self-describing
   empty table, so the population is visible without starting a session.
   Line: sonnet / medium. Read-only TOON query with strong prior art in the
   AXI surface.
9. As a session, I want `bench worktree list` to follow the AXI posture —
   unknown arguments exit 2 with usage, `-h/--help` exits 0, errors are
   structured on stdout with exit 1 — so agents can drive it like every other
   query. Line: sonnet / low. Mechanical conformance to the existing contract.
10. As a user, I want `bin/bench.sh` usage and dispatch to carry `worktree
    list` to the Go core, so the documented surface and the real one agree.
    Line: sonnet / low. Thin shell plumbing at the known dispatch seam.
11. As a teammate, I want ADR 0005 amended so that "non-live" is determined by
    the recorded owner identity rather than lease-file presence and declared
    build outputs are discardable residue, so the written decision matches the
    shipped behavior. Line: opus / medium. Prose the gate cannot grade; a
    one-passage amendment transcribing an already-made decision, so it stays
    mid rather than invoking the top-tier doc-authoring override.
12. As a reviewer, I want this repo's `.bench/build-outputs.json` to declare
    `dist/`, so the fix actually bites here and the file is the single source
    for what the gate is allowed to leave behind. Line: sonnet / low. One
    config file validated by story 6's parser.

## Implementation decisions

- **One liveness probe, lifecycle-owned.** `internal/worktree/lifecycle.go`
  already parses `<pid> <utc-time>` leases and owns dead-PID reclaim logic; it
  gains a three-state probe (live/dead/unknown) and the planner consumes it.
  The planner today retains on `os.Lstat` of the lease path; that branch
  becomes: live → retain (`ReasonLiveLease`), unknown → retain
  (`ReasonUncertain`), dead → no lease-sourced retain. No second PID parser.
- **Containment, not discard, is the new logic.** The planner already builds a
  full ignored-path inventory with `OverLimit`/`AtLeast` truncation flags. A
  new check answers "is every inventoried path under a declared entry"; only a
  provably contained inventory skips the ignored retain. Removal still rides
  the existing plan/apply transaction and receipts still record the inventory,
  so evidence and idempotency are unchanged.
- **`.bench/build-outputs.json`** follows the `gate-inputs.json` schema style:
  `{"schema":1,"paths":["dist/"]}`. Entries are repo-relative; a trailing-slash
  entry covers its subtree. Absent file → empty set. Any violation (schema,
  traversal, absolute, control bytes) → the whole declaration is malformed →
  retain with an explicit reason. The parser is the single owner of these
  rules.
- **`bench worktree list`** is a read-only AXI query in `internal/worktree`,
  dispatched from the existing `case "worktree"` in `cmd/bench/main.go` and
  named in `bin/bench.sh` usage. It joins the assignment ledger with
  registered worktrees; landedness reuses the existing default-branch
  resolution and reports `unknown` rather than fabricating `main` (FT86 rule).
- **`bench resume` output format is unchanged** — reconciled dead-session
  trees land in the existing removed/recovered counters, so the three
  resume-summary goldens (FT94) are not touched.
- **Structure note:** `internal/worktree` is already one file over its
  directory budget; the list query adds one more file. Accepted as part of
  this row rather than forcing an unrelated split — the split belongs to the
  standing structure backlog.

## Testing decisions

- A good test drives a public command (`ReleaseCommand`, `CleanCommand`,
  `ConservativeCleanup`, the list command) against a fixture repo and asserts
  observable outcomes — exit codes, receipts, ledger records, directory
  existence, TOON rows — never planner internals.
- Seams: `internal/worktree` package tests for every planner/liveness/config
  behavior (prior art: `newOwnedAssignment`, `commitInWorktree`,
  `requireTest`, `ApplyExplicitWithOptions`, and the lease fixtures in the
  lifecycle tests — dead PIDs come from a spawned-and-exited child, the same
  trick the reclaim tests use). One runtime contract fragment in
  `internal/contract/runtime` drives the built binary for the `list` CLI
  contract and shell dispatch (rebuild `dist/bench` first when hand-running).
- Gate: `.bench/gate.sh` — the worktree package runs under the compiled-core
  checks and the contract fragment under the contract phase.

### Seam diagram

Planner classification (stories 1–7, 12):

    trigger: release --request / clean / resume (PlanAutomatic)
        │
        ▼
    lease file ─────────▶ [ liveness probe (lifecycle) ] ──▶ live|dead|unknown
    ignored inventory ──▶ [ cleanup planner            ] ──▶ CleanupPlan
    build-outputs.json ─▶ [   (one classify path)      ]      (action+reason)
    landed/tracked ─────▶ [                            ] ──▶ receipt + ledger
                      ◀ tests attach here: fixture repo + lease/config/ignored
                        fixtures in, plan action, exit code, receipt, record
                        presence, and tree existence out

List query (stories 8–10):

    trigger: bench worktree list (shell dispatch → Go core)
        │
        ▼
    assignment ledger ──▶ [ list query (read-only) ] ──▶ TOON worktrees[N]
    registered trees ───▶ [                        ] ──▶ exit 0/1/2
                      ◀ tests attach here: built binary in a fixture repo,
                        stdout TOON rows and exit code asserted

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | dead-lease landed clean tree releases; record compacted | worktree | today `ReleaseCommand` retains with "assignment has a live or ambiguous lease" | passes only if a dead PID stops the lease-sourced retain |
| 2 | live-lease tree retains on release | worktree | not TDD-able as a new red — guards current behavior; asserted alongside story 1 so the relaxation cannot overshoot | fails if the probe treats any parseable lease as dead |
| 3 | garbage/unreadable lease retains | worktree | fails without the fix only in reason code (uncertain vs live-lease); asserted as fail-closed regression | fails if the probe defaults unknown to dead |
| 4 | dist/-only ignored tree releases without --discard-ignored | worktree | today retains with "ignored residuals require --discard-ignored" | passes only if declared containment lifts the ignored retain |
| 5 | undeclared ignored path (and an over-limit inventory) retains | worktree | not TDD-able as a new red — guards current behavior against the story-4 relaxation | fails if containment matches too broadly or ignores truncation flags |
| 6 | absent config = empty set; malformed/traversal/absolute/control-byte config retains with malformed reason | worktree | absent-vs-malformed distinction does not exist today; malformed fixture red until parser lands | fails if a hostile declaration widens instead of narrowing |
| 7 | resume reconciles a dead-session landed clean tree; retains live-lease; reports preserved work | worktree | today `ConservativeCleanup` counts the dead-lease tree under retained live-lease | fails unless PlanAutomatic inherits the probe and the safety branches survive |
| 8 | list emits one row per assignment + foreign trees; empty table at zero | runtime contract | `bench worktree list` today exits with usage — the subcommand does not exist | any missing column, row class, or empty-state fails the pinned block |
| 9 | unknown args exit 2, help exits 0, structured error exit 1 | runtime contract | same command absent today; posture rows red until wired | fails if the query deviates from the AXI exit contract |
| 10 | shell usage names list and dispatch reaches the core | runtime contract | shell path is part of the story-8/9 fragment (built binary invoked via `bin/bench.sh`) | fails if usage and dispatch drift from the Go surface |
| 11 | ADR 0005 matches shipped behavior | — | not TDD-able — prose; reviewed on the review axis | — |
| 12 | this repo declares dist/ | worktree | story-4 fixture uses the repo's real declaration file shape; a missing kit declaration fails the story-7 end-to-end fixture | fails if the kit ships the mechanism without the declaration |

### Edge inventory

- Paths with spaces/glob characters — coverage rows 4–6 fixtures include a
  space-bearing ignored path and declaration entry.
- Control bytes in git-sourced text — list output: TOON emitter refuses;
  covered by the existing emitter contract, list adds a hostile-label fixture
  in row 8.
- Absent vs present-but-empty declaration — row 6 asserts both (absent = empty
  set; present-but-empty `paths` = empty set, not malformed).
- Hand-edited declaration without trailing newline — row 6 fixture variant.
- PID recycling — **Won't handle:** a recycled PID can only read as live,
  which retains; never destructive in the wrong direction.
- Interrupt mid-cleanup — **Won't handle:** removal rides the existing
  plan/apply transaction whose interrupt behavior is already gate-covered;
  this change adds no new destructive machinery.
- cwd deeper than repo root — list resolves the root the same way sibling
  worktree commands do; covered by invoking the row-8 fragment from a subdir.
- Re-run idempotency — second release after a dead-lease reconcile
  short-circuits on the terminal receipt (existing FT93 machinery); asserted
  in row 1.
- Invocation through every shipped surface — row 10 drives `bin/bench.sh`;
  the linked-repo by-path CLI shares that wrapper (existing surface contract).
- Foreign/identity-mismatched registrations, nested-dirty, embedded repos —
  untouched retain branches; guarded by existing planner tests, re-run by the
  same package.

## Out of scope

- **Consumer seeding of `build-outputs.json` via `bench link`** — separate
  capability: the declaration is per-repo content like `IDEAS.md`, and
  shipping consumer scaffolding is FT92's track. ~4 edits, ~2 gate runs.
- **Sweeping this repo's three live stale trees** — operational cleanup, not
  code; two release with this fix at the next resume, the third (5 unmerged
  draft commits) needs the reviewer's discard call. 0 edits.
- **Stale `.git/bench-cleanup-*.lock` residue (33 files)** — belongs to
  FT87's concurrency-safe repair-cleanup clause, already on the roadmap.
