# Ambient feedback — surface the state that drives the next action

> **GRILL DEFERRED.** This map is captured but unresolved. Do not resolve the tickets
> from prose — run `/grill` on it at the **start of the next session**, one ticket at
> a time, then record answers here. Frontier is identified, not decided.

## Vision (the core thesis)

The kit's core job is to **assist the user in taking the next action toward their
vision** by surfacing relevant information at the right moment — not to make them go
looking for it. Right now the user has to ask "how many learnings are open? what's the
todo status? what files are growing?" The kit should *volunteer* that, concisely, as
ambient feedback. This is the productization of the session's concision /
stage-summary / recommend-next-action learnings: feedback that is relevant, scannable,
and ends in a recommended action.

## Grounding — what already exists (compose before inventing)

- **SessionStart hook** (global `~/.claude/`) — already prints a cold-session dashboard
  (currently the gl-axi repo state). The natural home for a benchkit dashboard.
- **Stop hook** (`stop.sh`) — the completion oracle; runs the gate on shift-stop.
- **Task list** (TaskCreate/TaskUpdate) — the persistent progress visual adopted this
  session.
- **`.bench/learnings.md`** — the journal; `grep -c '[open]'` = current open count (4).
- **Large-file "after every run" trigger** — *referenced by the user but NOT found in
  this repo's `.claude/` (no PostToolUse hook).* Location/existence to confirm in #6.

## #1: What is the feedback surface, and on what cadence?

Type: Grill

### Question
Where does ambient feedback appear and how often — a SessionStart dashboard (once,
cold), a post-run report (after every tool/run, like the large-file idea), an
on-demand `bench status`, or a combination? Cadence trades freshness against noise;
"after every run" is the exact wall-of-text risk the user just flagged.

### Answer
— (deferred to next session's grill)

## #2: Which signals are relevant, and how are they prioritized?

Blocked by: #1
Type: Grill

### Question
Candidate signals: growing files/dirs (size deltas), open-learnings count, task/todo
status, gate red/green, uncommitted/unpushed changes, open shifts/worktrees, stale
decision maps. Which subset earns a place, and how is the list ranked so the most
action-relevant item leads? A dashboard that shows everything shows nothing.

### Answer
— (deferred to next session's grill)

## #3: How does it compute the recommended *next action*, not just dump metrics?

Blocked by: #2
Type: Grill

### Question
The thesis is "assist the next action," not "report numbers." Given the signals, how
does the surface derive a recommendation (e.g. "gate red → fix before commit", "4 open
learnings → run /resynthesize", "spec written, no build → /build")? Rules table,
priority ladder, or model judgment at render time?

### Answer
— (deferred to next session's grill)

## #4: How does it stay concise and non-noisy?

Blocked by: #1, #2
Type: Grill

### Question
This feature is feedback, and the user's standing rule is terse + scannable + no wall
of text. What keeps it small — a hard line/row budget, show-only-on-change, severity
threshold (surface a signal only when it crosses a bound, like a file growing past N)?

### Answer
— (deferred to next session's grill)

## #5: Compose an existing seam or add a `bench status` command?

Blocked by: #1
Type: Grill

### Question
Build on the SessionStart hook + task list + the (to-confirm) large-file trigger, or
introduce a new `bench status` CLI subcommand as the single source the hook and the
user both call? Invariant 4 (compose before inventing) and the legibility ceiling both
bear on this — a new surface must fill a gap the existing layers can't.

### Answer
— (deferred to next session's grill)

## #6: Confirm the large-file trigger — where, and what does it report?

Type: Research

### Question
The user says "we have a trigger after every run to check for large files," but it's
not in this repo's `.claude/settings.json`. Is it a global hook, another harness, or
not-yet-built? Its existence and shape determine whether ambient feedback extends it
or replaces it.

### Answer
— (deferred to next session's grill)
