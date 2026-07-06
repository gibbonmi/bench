---
name: craft-synthesis
description: The discipline for safely folding a candidate change into the Bench kit — respect closed decisions, assess by the gap it fills, pass three quality loops (legibility, consistency, dogfood), and propose rather than merge. Use whenever evaluating a change to the kit itself, from either source. Reach for this from /bench-update-kit, from any learnings-sourced roadmap item queued by /bench-what-next, or any time you're deciding whether a kit change earns its place.
index: evaluating a change to the kit itself
---

# Synthesis — fold a change into the kit without bloating it

This is the shared middle of `/bench-update-kit` (upstream changes) and the
learnings-sourced roadmap items `/bench-what-next` queues (learnings changes).
Those paths gather candidates and record results; this skill is what happens in
between — the same discipline regardless of where a candidate came from. The root virtue is **anti-sediment**: the kit growing is a *cost*, and a change
earns its place only by filling a real gap, never because it reads well in isolation.

The candidates arrive already tagged by origin (`upstream` or `learnings`). Keep the
tag on every one through every step, so the proposal shows where each came from. It
**proposes; the reviewer merges.**

## Respect closed decisions

Drop every candidate the baseline already settled and that hasn't materially changed.
A rejected item does not reopen on a cosmetic change — only a material change reopens
it, and then you show the diff that justifies reopening. For `upstream`, that's a
source repo whose technique actually moved; for `learnings`, an entry already
resolved (pruned from the journal — its verdict lives in the CHANGELOG and the
integration commit) is not re-litigated. A loop that re-opens settled calls every
run is worse than no loop.

## Assess — propose, don't decide

For each surviving candidate, classify in one line:

- **Map** — fills a gap the existing layers can't.
- **Fold** — absorb into an existing piece.
- **Recommend** — note it; don't build it.
- **Skip** — with the reason.

The bar is the anti-sediment bar from `craft-skills`: gap-filling, not
goodness-in-isolation. Present the proposed set with each item's origin tag.

## Three quality loops — all must pass before adoption

Run in order; a change that fails a loop is pruned or sent back, not shipped.

1. **Legibility loop.** Run `craft-skills` against each change. No-op?
   Duplicates an existing piece? Pushes the kit past its legibility ceiling? Cut or
   fold it.
2. **Consistency loop.** Apply to a working copy, then re-run the staleness audit:
   grep for invariant drift, broken cross-references, stale paths, app-specific
   leakage into core files, an out-of-date provenance table. Fix every hit.
3. **Dogfood loop — the oracle.** Run a real shift on a real repo with the changed
   kit: `/bench-write-spec` a small task, `bench shift`, confirm the gate gates, the hooks
   fire, `bench gate` ends green. A change that reads well but breaks a real run is
   rejected. If you can't run a dogfood shift, the synthesis is **not complete** — say
   so rather than adopting on paper. The kit does not grade its own update; a run does.
   A candidate that changes a skill or command trigger needs the dogfood run in a
   fresh session, because this session loaded those surfaces before the edit.
   Proportionality: a prose-only change — no hook, gate, CLI, or adapter touched —
   may substitute a green `bench gate` plus a read of every surface the prose
   steers for the full shift; say which verification ran. Anything that touches
   behavior always dogfoods.

## Propose; the reviewer merges

Never auto-apply. After the loops pass and the reviewer has signed off, hand back to
the calling path to record the result — `/bench-update-kit` to the provenance table and
CHANGELOG; a learnings-sourced build into the fixed artifact, the CHANGELOG verdict
line, and the retired roadmap row. The merge is the reviewer's.
