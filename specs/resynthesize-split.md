# Split resynthesize into bench-update + bench-learn

Implements decision #4 of `decisions/bench-naming.md`: `resynthesize` becomes two
commands reading two different inputs, with their shared review discipline extracted
to a skill. This is the behavior-change slice the mechanical rename deferred.

## Problem

`/resynthesize` bundles two unrelated jobs behind a "choose the scope" switch:
syncing the kit against *upstream* repos, and draining the *local* learnings journal.
They share a review path but have nothing else in common — different inputs, different
records, different cadence. One command for two jobs means the user picks a sub-mode
every run, and the ambient `bench status` learnings nudge points at a maintenance
command whose name says nothing about learnings.

## Solution

Two commands, one per job, sharing the review discipline through a skill:

- **`/bench-update`** — others' changes. Pull the upstream repos, diff against the
  baseline, propose what to adopt, record adoptions in the provenance table + CHANGELOG.
- **`/bench-learn`** — your own changes. Drain `.bench/learnings.md`, propose
  promotions, land them in the rule/skill they fix, mark the entries resolved.
- **`bench-craft-synthesis`** (new skill) — the shared middle both run: respect closed
  decisions → assess (Map/Fold/Recommend/Skip) → three quality loops → propose, don't
  merge. Invoked by both commands, exactly as `bench-craft-grill` serves `/bench-ideate`
  and `/bench-spec`.

`resynthesize.md` is removed. The `bench status` learnings signal points at
`/bench-learn`.

## User stories

1. As a user, I want `/bench-update` to sync against upstream — pull the source repos,
   diff against the CHANGELOG/provenance baseline, classify new/changed/gone — so
   others' improvements reach the kit.
2. As a user, I want `/bench-learn` to drain `.bench/learnings.md` — read open entries,
   group ones pointing at the same fix, dismiss one-offs — so the kit learns from use.
3. As a user, I want both commands to run the identical review discipline (closed-
   decisions guard → assess → three quality loops → propose-don't-merge) via the
   `bench-craft-synthesis` skill, so the discipline is single-sourced and can't drift.
4. As a user, I want `/bench-update` to record adoptions in the provenance table + a
   CHANGELOG entry (the baseline the next sync reads), so upstream state is tracked.
5. As a user, I want `/bench-learn` to land promotions in the rule/skill/command they
   fix, add a CHANGELOG line, and mark the source learnings entries resolved, so they're
   never re-reviewed.
6. As a user, I want `resynthesize.md` removed, so there's one command per job and no
   stale third entry-point.
7. As the agent, I want `bench-craft-synthesis` to carry frontmatter and an AGENTS.md
   skills-index entry, so it loads and the gate's skill-sync passes.
8. As a user, I want the `bench status` learnings row to point at `/bench-learn` (not
   `/resynthesize`), so the ambient nudge names the command that drains the journal.
9. As the agent, I want the gate's status contract moved to `/bench-learn` in lockstep
   with the renderer, so the oracle still checks the live signal — fixture-following,
   not a weakened check.
10. As a user, I want `bench init`'s scaffolded `learnings.md` header to name
    `/bench-learn`, so a freshly-linked repo's journal points at the right command.
11. As a user, I want every doc reference updated by stream — *upstream* mentions
    (AGENTS.md, README, CHANGELOG header, CONTEXT, HANDOFF) → `/bench-update`,
    *learnings* mentions → `/bench-learn` — so no doc names the removed command.
12. As a user, I want the `bench-craft-skills` model-invoked example list to include
    `bench-craft-synthesis`, so the skill inventory stays accurate.
13. As a user, I want the learnings.md / CHANGELOG **dated entries** and all of
    `decisions/` left untouched, so the historical record isn't falsified.
14. As a user, I want `bench gate` green and zero `/resynthesize` on the active surface,
    so done is the oracle's verdict plus a clean sweep.

## Implementation decisions

**New skill `bench-craft-synthesis`** (model-invoked; frontmatter `name` +
`description`). Stream-agnostic: it takes candidate changes already tagged by origin
and runs them through closed-decisions guard → assess (Map/Fold/Recommend/Skip,
anti-sediment bar, delegating the legibility judgment to `bench-craft-skills`) → three
quality loops (legibility / consistency / dogfood-is-the-oracle) → propose, never
merge. The commands supply the candidates and do the recording; the skill owns the
discipline in between.

**Two commands** (`disable-model-invocation: true`, as `resynthesize` was). Each is
thin: its **input** + "run the candidates through `bench-craft-synthesis`" + its
**recording**. `bench-update` input = pull repos + diff vs baseline; recording =
provenance table + CHANGELOG. `bench-learn` input = drain `learnings.md`; recording =
land in the fixed artifact + CHANGELOG line + mark entries resolved.

**No combined runner.** The old "both" default becomes "invoke each." Running both is
two commands; there is no third command.

**Reference updates are per-stream, not a blind replace.** `/resynthesize` maps to
*two* targets depending on context (upstream→`/bench-update`, learnings→`/bench-learn`),
so each reference is classified and replaced individually — never a global sed (the
lesson from the rename slice's separator-slash corruption).

**Gate/CLI literals move with the signal.** `bin/bench.sh`'s learnings status row and
`.bench/gate.sh`'s status contract (the `/resynthesize` checks) both become
`/bench-learn`, in lockstep — fixture-following, assertion logic unchanged.

**`.claude/` untouched** — `.claude/skills` symlinks into `.agents/`, so the new skill
dir passes through.

**Historical surfaces are left intact:** learnings.md dated entries ("Resolved via
/resynthesize"), CHANGELOG dated entries, and all of `decisions/`. The learnings.md
*header* and CHANGELOG *header* (current-state instructions) do update.

## Testing decisions

**Seam: `bench gate` (the oracle).** The split is structurally correct when green:
- **command↔index conformance** (gate 5c) — `bench-update` + `bench-learn` referenced
  as `/name` in AGENTS.md; `resynthesize.md` gone, so its requirement disappears.
- **skill↔index conformance + frontmatter** (5a/5b/3) — `bench-craft-synthesis` indexed
  and fronted.
- **status contract** (1F) — the learnings signal now `/bench-learn`, renderer and
  check moved together; canary still bites (section 7 — update the `status-regressed`
  fixture's EXPECT only if it pinned the old string; the canary run reports it).
- **pack/init integrity** (4, 1d) — `.agents/` ships wholesale; `bench init` still
  scaffolds `learnings.md` (existence check unaffected by the header text change).

**Completeness sweep** (gate-green is necessary, not sufficient): grep the active
surface for `/resynthesize`, asserting zero — excluding `decisions/`, the learnings.md
and CHANGELOG dated entries, and `specs/bench-naming-mechanical.md`.

**What the gate does NOT test, stated honestly:** the new commands and skill are prose
instructions to an agent. Their *semantic* correctness (do they faithfully capture the
two streams + shared path) is the `/bench-review` pass. Their *behavioral* proof —
that following them produces a sound kit update — is the synthesis discipline's own
**dogfood loop**, exercised the next time the commands are actually run on a real
upstream sync or learnings drain. This docs slice builds the definitions; it does not
run them.

**Gate command:** `bash .bench/gate.sh`.

## Out of scope

- **Actually running `/bench-learn` to drain the current open learnings.** That's
  *using* the command (and needs per-promotion sign-off), not building it — a separate
  act. ~varies.
- **Capture wiring** (agent runs `bench idea` on "park this"). The other deferred
  behavior change from the map; its own slice. ~20m.
- **HANDOFF.md retirement.** A parked roadmap item about its content; its command-list
  reference *is* updated here. Separate.
