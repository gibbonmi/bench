# ambient-feedback surface — `bench status`

> Feature **A** of the ambient-feedback map (`decisions/ambient-feedback.md`, #1–#6,
> plus the #8 footer that wires in feature B). B (the roadmap) is already built; A reads
> its `ROADMAP.md` for the footer.

## Problem

The kit knows things the user has to keep asking for: is the gate red, are there
uncommitted changes, how many learnings are open, is a worktree stranded, is a decision
map still unresolved? Today you go looking. A cold session starts blind, and the
standing rule (terse, no wall of text) means the answer can't just be "dump everything."

## Solution

One renderer, `bench status`, that surfaces only what needs attention and leads with a
single recommended next action. It runs automatically at the start of a Claude Code
session (cold dashboard) and on demand any time. Signals appear only when they fire,
ranked by urgency, capped at five rows; when nothing needs attention it collapses to one
line. Parked roadmap ideas show as a passive footer count — visible, never nagging. The
gate's verdict comes from a cache the Stop hook writes as you work, so the dashboard is
fast and never runs the gate cold.

## User stories

1. As a user starting a cold session, I want a one-screen view of what needs attention
   so I don't have to interrogate the repo myself.
2. As a user, I want the view to lead with a single recommended next action, not a wall
   of metrics.
3. As a user, I want `bench status` runnable on demand with the same output, so I can
   re-check anytime.
4. As a user, I want a signal to appear only when it has something to report — no "all
   clean" checkmarks.
5. As a user, when nothing needs attention I want one terse line (`bench: clean —
   nothing pending`), not an empty table.
6. As a user, I want at most five signal rows, with a `+k more` tail when more fire, so
   it never becomes a wall.
7. As a user, I want signals ranked by urgency so the most action-relevant one leads.
8. As a user, I want a red gate surfaced as the top blocker, with "fix before commit."
9. As a user, I want gate state read from a cache, not a slow cold re-run, so SessionStart
   stays fast.
10. As a user, I want a green-but-stale gate to *not* read as a clean bill — it shows
    "stale — re-run."
11. As a user, I want a fresh green gate to be silent (no row).
12. As a user with no gate cache yet, I want the gate signal silent — no nag on a fresh
    repo.
13. As a user, I want uncommitted/unpushed changes surfaced with "commit / push."
14. As a user, I want a stray worktree or active shift surfaced with "resume or clean up."
15. As a user, I want open learnings (≥ a floor) surfaced with "/resynthesize" and the
    count.
16. As a user, I want structural debt surfaced with "split."
17. As a user, I want an unresolved decision map surfaced with "/bench-craft-grill → /bench-spec."
18. As a user, I want parked roadmap ideas shown as a passive footer count (`N idea(s)
    parked — bench roadmap`), never an action row, never inside the five-row budget.
19. As a user, I want the dashboard to appear automatically at SessionStart in Claude
    Code.
20. As a user, I want the gate cache written automatically as I work a shift, so the
    dashboard reflects my latest verdict.
21. As a user, I want the cache to live outside git (runtime state) so it never pollutes
    commits.
22. As a user on a harness without SessionStart hooks, I want `bench status` to still
    work when I run it by hand.
23. As a user, I want the learnings floor configurable via an env var.

## Implementation decisions

- **New `bench status` subcommand — the single renderer.** Deterministic plain shell, no
  model in the loop (#3): it reads repo state plus the gate cache, computes a fixed
  severity ladder, and emits scannable structured output per the `bench-craft-cli` skill — a one-line
  lead, then signal rows (`signal · detail · action`), then the footer. The agent reading
  it applies judgment; the renderer never does.
- **Severity ladder** (rank, lower = more urgent; lead = the lowest-rank present signal's
  action):

  | Rank | Signal | Trigger | Action shown |
  |---|---|---|---|
  | 0 | gate red | cache status=red **and** fresh (sha == HEAD) | fix before commit |
  | 1 | uncommitted / unpushed | `git status --porcelain` non-empty, or commits ahead of upstream | commit on green / push |
  | 2 | stray worktree / active shift | `git worktree list` shows > 1 | resume or clean up (`bench worktree`) |
  | 3 | open learnings | `grep -c '[open]' .bench/learnings.md` ≥ floor | `/resynthesize` (+count) |
  | 4 | structural debt | `bench structure` exits non-zero | split (bench-craft-seams) (+count) |
  | 5 | unresolved decision map | a `decisions/*.md` holds an open-ticket marker | `/bench-craft-grill` → `/bench-spec` |
  | 6 | gate stale | cache exists but sha != HEAD | re-run the gate |

- **Concision (#4):** show-only-on-signal; hard five-row budget with a `+k more` tail;
  one-line lead headline; an all-clear line when nothing fires. The footer is outside the
  budget.
- **Gate signal from a cache file in the git dir** —
  `$(git rev-parse --absolute-git-dir)/bench-last-gate` (format: `<status> <sha>
  <iso8601>`). Chosen over `.bench/last-gate`: the git dir is never tracked, so the cache
  can't pollute commits or read as dirty and needs no `.gitignore` management — story 21
  for free, and it works for consumers automatically. Resolution: **red** iff status=red
  and sha == HEAD; **stale** iff sha != HEAD (whatever the status — a stale green is never
  a clean bill, #2); **silent** iff (green and fresh) or no cache at all.
- **Cache writer: the Stop hook** (`.bench/hooks/stop.sh`), composing the existing
  completion-oracle seam. It already runs `bench gate` and branches on the verdict; it
  writes the cache in both branches. This is chosen over "the shift loop writes it"
  (#2's tentative phrasing) because `run_gate` **execs** the gate, so the in-process loop
  cannot act after a gate run — whereas the Stop hook runs the gate as a subprocess and
  already holds the verdict. *(Open for veto; see Out of scope re: the exec defect.)*
- **Cache is runtime state:** lives in the git dir (never tracked), never in
  `package.json` `files[]`.
- **SessionStart trigger:** a shared `.bench/hooks/session-start.sh` runs `bench status`;
  the Claude Code adapter wires it under `.claude/settings.json` `hooks.SessionStart`.
  `bench link` installs the shared hook and the adapter (the link contract already
  exercises shared-hook installation). This composes the existing hook architecture
  (`stop.sh`, `block-dangerous-git.sh`) rather than inventing a surface type (#5).
- **Footer (#8, the B wire-up):** when `ROADMAP.md` is non-empty, append one
  zero-severity line `N ideas parked — bench roadmap`. Never ranked, never the lead,
  never counted against the five-row budget.
- **Open-ticket marker** (rank 5) is a documented convention: a decision file is
  unresolved when it carries an open-answer placeholder (an answer body beginning `—
  (open` / `— (deferred`) or a `GRILL DEFERRED` banner. Crisp enough for a `grep`.
- **Thresholds via env:** `BENCH_LEARNINGS_FLOOR` (default 1); structural debt reuses
  `bench structure`'s existing `BENCH_MAX_LINES` / `BENCH_MAX_DIR_FILES`.

## Testing decisions

- **A good test** constructs repo state in a throwaway repo, runs `bench status`, and
  asserts the rendered output — external behavior at the CLI seam, never internals.
- **Seam:** `bench status` (the renderer). One primary seam, the highest that exercises
  the real behavior. **Prior art:** the `bench idea`/`roadmap` block (gate check 1f) and
  the `init`/`link` contract blocks already in `.bench/gate.sh` use exactly this
  throwaway-repo, exercise-the-real-CLI pattern.
- **Gate contract block** (add to `.bench/gate.sh`) asserting: the all-clear line on a
  clean repo; each signal appears with its action when its state is constructed; the lead
  is the highest-severity present signal; the five-row budget caps with `+k more`; the
  footer appears iff `ROADMAP.md` is non-empty and is never the lead; gate **red** (fresh
  red cache), gate **stale** (sha mismatch), and gate **silent** (fresh green cache, and
  no cache).
- **Cache write:** a focused assertion that the Stop hook writes `<git-dir>/bench-last-gate`
  in a format `bench status` reads back (round-trip), exercised with `BENCH_SHIFT=1`.
- **Canary fixture** (`tests/canary/`): a `bench status` that drops a signal (e.g. always
  prints clean) must make the gate go red, proving the renderer contract bites. Follow the
  existing fixture shape.
- **Gate command:** the project gate — `.bench/gate.sh` (`bench gate`).

## Out of scope

- **Fixing the `run_gate` exec defect** so the in-process shift loop iterates and could
  itself write the cache. A is designed to not depend on it (the Stop hook writes the
  cache). The exec — `run_gate` replaces the process, so the loop's `if run_gate` cannot
  continue or act afterward — is a pre-existing shift-loop bug with its own blast radius
  and its own test. Surfaced as a finding, not folded in silently. Est ~30–45 min + a
  loop test.
- **SessionStart auto-wiring for Codex / non-Claude harnesses.** `bench status` is
  harness-independent and runnable by hand or wired manually; the *automatic* trigger
  depends on each harness's hook model (Claude Code's `hooks.SessionStart` here). A
  distinct capability per harness once its hook contract is confirmed. Est ~20–30 min
  each.
- **Growth/size-delta tracking of arbitrary files** (the original "growing files" idea)
  beyond `bench structure`'s source-line check. A real new signal needing a persisted
  baseline to diff against. Est ~1 hr.
- **Per-signal mute / richer config** beyond the learnings floor and the existing
  structure thresholds. Separate capability. Est ~30 min.
