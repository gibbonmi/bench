# Ambient feedback — surface the state that drives the next action

> **GRILLED & RESOLVED (2026-06-30).** All tickets (#1–#9) decided and recorded below.
> Grown to two related features — the ambient-feedback **surface** (#1–#6) and the
> **roadmap/icebox** capture-and-forget sink (#7–#9, linked via the surface footer in
> #8). **Decided: two specs, roadmap (B) first.** B is **built** (`specs/roadmap.md`,
> `bench idea`/`bench roadmap`). Next action: `/spec` on A (the surface).

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
- **Large-file check** — exists as `bench structure` (a manual CLI subcommand), NOT a
  hook and NOT wired to any cadence. See #6 (resolved).

## #1: What is the feedback surface, and on what cadence?

Type: Grill

### Question
Where does ambient feedback appear and how often — a SessionStart dashboard (once,
cold), a post-run report (after every tool/run, like the large-file idea), an
on-demand `bench status`, or a combination? Cadence trades freshness against noise;
"after every run" is the exact wall-of-text risk the user just flagged.

### Answer
Resolved (grill, 2026-06-30). **One renderer, two triggers:** a single source computes
the feedback, surfaced automatically at **SessionStart** (cold, once per session) and
callable **on-demand** as `bench status`. **No per-run cadence** — "after every run" is
rejected as the wall-of-text risk, and #6 established there is no post-run hook to ride
anyway, so building one would be opting into the noise. SessionStart already fires a
dashboard hook, making the cold surface cheap to wire; the single renderer the hook and
the user both call gives one source of truth (feeds #5).

## #2: Which signals are relevant, and how are they prioritized?

Blocked by: #1
Type: Grill

### Question
Candidate signals: growing files/dirs (size deltas), open-learnings count, task/todo
status, gate red/green, uncommitted/unpushed changes, open shifts/worktrees, stale
decision maps. Which subset earns a place, and how is the list ranked so the most
action-relevant item leads? A dashboard that shows everything shows nothing.

### Answer
Resolved (grill, 2026-06-30). **Six signals, ranked by a severity ladder** — the
highest-severity present signal becomes the lead recommendation; the rest render as a
compact list (budget/threshold mechanics are #4).

| Rank | Signal | Source (cheap) | Implied action |
|---|---|---|---|
| — | **Gate red/green** | **cached last-shift result** (not a cold run) | red → fix before commit |
| 1 | Uncommitted / unpushed changes | `git status --porcelain`, `git cherry` | commit on green / push |
| 2 | Active shift or stray worktree | `git worktree list` | resume or clean up |
| 3 | Open-learnings count | `grep -c '[open]' .bench/learnings.md` | `/resynthesize` |
| 4 | Structural debt | `bench structure` | split (seams) |
| 5 | Unresolved decision map | open tickets in `decisions/*.md` | `/grill` → `/spec` |

**Gate signal (kept, via cache).** Running the full gate cold is too slow for a
SessionStart hook, so the renderer never runs it — instead the **shift loop writes a
cached last-gate result** (`bench shift` already runs the gate after each iteration).
The renderer reads that cache. A **red** cached gate is the top blocker (outranks
everything — it is the #3 thesis example). Open mechanism detail for spec/#5: what
writes the cache and what it stores (status + commit sha + timestamp), and how the
renderer flags **staleness** when the working tree has changed since the cached sha
(green-but-stale must not read as a clean bill).

**One drop:** in-session task list — ephemeral, not on disk, unreadable by a shell
renderer.

## #3: How does it compute the recommended *next action*, not just dump metrics?

Blocked by: #2
Type: Grill

### Question
The thesis is "assist the next action," not "report numbers." Given the signals, how
does the surface derive a recommendation (e.g. "gate red → fix before commit", "4 open
learnings → run /resynthesize", "spec written, no build → /build")? Rules table,
priority ladder, or model judgment at render time?

### Answer
Resolved (grill, 2026-06-30). **Deterministic rules table evaluated by the renderer —
no model at render time.** Each signal carries a fixed `(severity, action-string)`; the
renderer emits the ranked signals as structured output (axi/TOON) plus one lead
recommendation = the highest-severity present signal's action. Model judgment is
deliberately **not** in the renderer: it is plain shell inside a hook (no model in the
loop), and determinism makes it cheap, testable, and canary-able. Judgment lives one
layer up — the **agent** reads the dashboard and tailors the next-action call to live
context (the hand-off / recommend-next-action learning). Net: the renderer always "ends
in a recommended action" for the human at SessionStart, and hands the agent structured
signals to reason over.

## #4: How does it stay concise and non-noisy?

Blocked by: #1, #2
Type: Grill

### Question
This feature is feedback, and the user's standing rule is terse + scannable + no wall
of text. What keeps it small — a hard line/row budget, show-only-on-change, severity
threshold (surface a signal only when it crosses a bound, like a file growing past N)?

### Answer
Resolved (grill, 2026-06-30). **Three mechanics together:**
1. **Show-only-on-signal** — a signal with nothing to report emits no row (no "✓ all
   clean" spam). Per-signal trigger bounds: gate row only when **red or stale** (green
   + fresh is silent), git only when dirty/unpushed, learnings only when count ≥ 1
   (configurable floor), structure only on a violation, maps only on an unresolved
   ticket.
2. **Hard row budget** — cap at the **top 5 by severity**; if more fire, show 5 and a
   `+k more` tail so the worst case stays bounded.
3. **One-line lead** — the single highest-severity recommendation as a headline above
   the compact list.

**All-clear:** when nothing fires, collapse to one terse line (`bench: clean — nothing
pending`), never a table of green checks.

## #7: Capture an out-of-scope idea without committing to it (the roadmap/icebox)

Type: Grill

### Question
While building, the user gets an idea that is out of scope for the current build. There
is no kit mechanism to capture it without commitment — decision maps imply intent to
resolve, specs are committed work, issues are a backlog you'll work. Need a pure
capture-and-forget sink: captured, visible when the user chooses to look, and exerting
**zero** workflow obligation. What is the capture command, the storage, and where it
lives?

### Answer
Resolved (grill, 2026-06-30; shipped in `specs/roadmap.md`). `bench idea "<text>"` parks
the idea — one shot, appends and exits, no prompt, no grill, no spec. Storage is a single
append-only `ROADMAP.md` at the repo root, one dated line per entry (`- YYYY-MM-DD
<text>`), committed and product-facing. No IDs, no status. `bench roadmap` prints it on
demand. Root over `.bench/icebox.md` because visibility is the whole point — "I can see
it."

## #8: Does the roadmap appear on the ambient surface, and at what severity?

Blocked by: #2, #4, #7
Type: Grill

### Question
The surface's whole job is to surface the *next action*; a parked idea is explicitly
**not** one — surfacing it as an action item would defeat the capture-without-nagging
purpose. Does the roadmap show on the dashboard at all, and if so how — a passive
count, never in the severity ladder, never the lead?

### Answer
Resolved (grill, 2026-06-30; shipped in feature A). Footer count only. When `ROADMAP.md`
is non-empty the surface shows one passive line (`N idea(s) parked — bench roadmap`) at
**severity zero**: never in the severity ladder, never the lead recommendation, and not
counted against the five-row budget. Empty roadmap → no line. This threads "I can see it"
without "I'll work on it" — the surface acknowledges the roadmap without ever nudging the
user to act on it.

## #9: Promotion — how a parked idea graduates into committed work

Blocked by: #7

Type: Grill

### Question
When the user decides to act on a parked idea, how does it leave the roadmap and enter
the real workflow (`/start-ideation`, `/spec`, or `/to-issues`)? Is capture append-only
with manual pruning, or does the roadmap carry status/lifecycle? The commitment point
is the graduation, not the capture.

### Answer
Resolved (grill, 2026-06-30). **Capture append-only; no status/lifecycle in the
roadmap** — once committed, an idea's real state lives in its map/spec/issue, so the
roadmap stays a dumb sink (avoids two sources of truth).

**`/start-ideation` is the promotion seam.** When `/start-ideation` is invoked **on its
own** — i.e. with no specific idea already in hand from the conversation — it reads
`ROADMAP.md` and **asks the user which parked items, if any, they want to pull up** to
ideate on. The chosen item seeds the ideation session. When `/start-ideation` is already
carrying a fresh idea from the conversation, it proceeds with that and does not
interrupt with the roadmap prompt. `/spec` and `/to-issues` remain valid manual entry
points for an idea clear enough to skip ideation.

**On pull, the line is auto-removed** from `ROADMAP.md` at the moment `/start-ideation`
**creates the ideation map** from the pulled idea (decided 2026-06-30). Not at
selection: an abandoned pull that never writes a map keeps its line. The roadmap thus
holds only un-promoted ideas, with no manual cleanup.

**Build implication:** this changes `/start-ideation`'s command file — it must learn to
read the roadmap and offer items on a cold invocation. Captured here so the spec covers
it; the command edit is a build task, not part of this map.

## Open scoping decision (the user's call) — one spec or two?

The map resolved into **two related but separable features**:

- **A — the ambient-feedback surface** (#1–#6): one renderer, surfaced at SessionStart
  + on-demand `bench status`; six signals on a severity ladder; show-only-on-signal +
  5-row budget + one-line lead; gate read from a shift-written cache.
- **B — the roadmap/icebox** (#7–#9): `bench idea` capture → append-only `ROADMAP.md`
  → `bench roadmap` view; promotion via `/start-ideation` on a cold invoke.

They touch at exactly one seam: B's non-empty roadmap shows as a zero-severity **footer
count** on A's surface (#8). Otherwise each builds and ships without the other.

**Decided (2026-06-30): two specs, roadmap (B) first.** B is small and independently
useful and solves the live capture pain; it ships in one short build. A (the larger
surface) follows, and reads B's `ROADMAP.md` for its footer as its closing wire-up
(B-first means the footer has data to read the moment A lands).

## #5: Compose an existing seam or add a `bench status` command?

Blocked by: #1
Type: Grill

### Question
Build on the SessionStart hook + task list + the (to-confirm) large-file trigger, or
introduce a new `bench status` CLI subcommand as the single source the hook and the
user both call? Invariant 4 (compose before inventing) and the legibility ceiling both
bear on this — a new surface must fill a gap the existing layers can't.

### Answer
Resolved (follows from #1 + #6, 2026-06-30). **Both — compose the trigger, invent the
renderer.** Add a new `bench status` subcommand as the single renderer (the one source
of truth from #1), and **compose the existing SessionStart hook** by having it call
`bench status` rather than standing up a new surface type. #6 established there is no
post-run trigger to extend, so "compose" here means wiring the existing SessionStart
hook plus the existing signal sources (`bench structure`, git, `learnings.md`) *into*
the renderer — not extending a trigger that doesn't exist. The new subcommand clears the
invariant-4 bar because it fills a gap no existing layer covers: one deterministic
source both the hook and the user invoke.

## #6: Confirm the large-file trigger — where, and what does it report?

Type: Research

### Question
The user says "we have a trigger after every run to check for large files," but it's
not in this repo's `.claude/settings.json`. Is it a global hook, another harness, or
not-yet-built? Its existence and shape determine whether ambient feedback extends it
or replaces it.

### Answer
Resolved (research, 2026-06-30). It is **`bench structure`** (`bin/bench.sh:491`), a
manual CLI subcommand — **not** a hook, and **not** wired to run after every run. No
PostToolUse hook exists in this repo's `.claude/` or the global `~/.claude/`; the only
hooks are Stop (`stop.sh`), PreToolUse git-guard (`block-dangerous-git.sh`), and a
global SessionStart (gl-axi). `bench structure` flags files >400 lines
(`BENCH_MAX_LINES`) and dirs >12 source files (`BENCH_MAX_DIR_FILES`), but only over
source-code extensions (`py|ts|tsx|js|...`) — it would **not** catch a growing
`learnings.md`, decision map, or other markdown/JSON. Its help text labels it "(wire
into the gate)": the check is built but bound to no automatic cadence.

Consequence for the map: there is **no existing post-run trigger to extend**. Ambient
feedback either wires `bench structure` (and other signals) into a cadence — the gate,
SessionStart, or a new surface — or stands up its own. This sharpens #1 (no free
"after every run" seam exists) and #5 (composing means *wiring*, not *extending*).
