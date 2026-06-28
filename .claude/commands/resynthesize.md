---
description: Re-synthesize Bench from its two upstream sources (Matt Pocock's skills and kunchenguid's tooling). Pull both, diff against what Bench already incorporates, assess each change for whether it improves the current iteration, and adopt the worthwhile ones only after three quality loops. Run periodically. Proposes; never auto-merges. Maintenance, not a workflow phase.
disable-model-invocation: true
---

# /resynthesize — re-run the synthesis against upstream

Bench is a synthesis of two moving repos. This re-runs that synthesis against their
latest state, so upstream improvements reach Bench without you re-reading both repos
by hand. It **proposes; you merge.** It does not reopen a decision Bench already
closed unless the upstream thing materially changed.

## 1. Pull and diff

- Pull `mattpocock/skills` and kunchenguid's repos (`axi`, `gnhf`, `treehouse`,
  `firstmate`, `no-mistakes`) at HEAD. Record the refs you pulled.
- Read Bench's current synthesis record: the provenance table in `README.md` and
  the entries in `CHANGELOG.md` from prior runs. This is the baseline — what Bench
  already maps, folds, recommends, or has deliberately rejected.
- Diff upstream against the baseline into three lists: **new** (skills or techniques
  not in the baseline), **changed** (a mapped/folded item whose source file moved),
  **gone** (something Bench uses that was deprecated upstream).
- Read `.bench/learnings.md` — the usage journal. Its open entries are a *second*
  source of candidate changes, this time from how the kit actually behaved in real
  work, not from upstream. Treat each open learning as a proposed change to assess
  alongside the upstream deltas. (This is the self-improvement path: the kit gets
  better from its own use, but only through this reviewed, signed-off loop — a
  learning never edits a rule on its own.)

## 2. Respect closed decisions

Before assessing anything, drop every delta the baseline already decided that hasn't
materially changed upstream. If Bench rejected firstmate's fleet orchestrator, a
commit that only edits its README does not reopen it. A closed decision reopens only
on a material change — and when it does, say so and show the diff that justifies
reopening. This guard is the point: a loop that re-litigates settled calls every run
is worse than no loop. (It's also the failure mode that got the ECC harness rejected.)

## 3. Assess — propose, don't decide

For each surviving delta, classify in one line: **Map** (fills a gap the existing
layers can't), **Fold** (absorb into an existing piece), **Recommend** (note it,
don't build), **Skip** (with the reason). The bar is the anti-sediment bar from
`writing-great-skills`: a change earns its place only if it fills a real gap, never
because it's good in isolation. Present the proposed set as a proposal.

## 4. Three quality loops — all must pass before adoption

Run these in order. A change that fails a loop is pruned or sent back, not shipped.

1. **Legibility loop.** Run `writing-great-skills` against each proposed change. Is
   it a no-op? Does it duplicate an existing skill or command? Does it push the kit
   past its legibility ceiling (the working-command surface, the skill count)? If so,
   cut or fold it. A bigger kit is a cost, not a feature.
2. **Consistency loop.** Apply to a working copy, then re-run the staleness audit:
   grep for invariant-count drift, broken cross-references, stale paths, app-specific
   leakage into core files, and an out-of-date provenance table. Fix every hit. The
   kit must be consistent *after* the change, not just at the changed spot.
3. **Dogfood loop — the oracle.** This is the external check and the one that
   actually decides. Run a real shift on a real repo with the changed kit: `/spec` a
   small task, `bench shift`, confirm the gate gates, the hooks fire, and `bench gate`
   ends green. A change that reads well but breaks a real run is rejected. If you
   can't run a dogfood shift, the re-synthesis is **not complete** — say so rather
   than adopting on paper. This is Bench's own thesis turned on itself: the kit does
   not grade its own update; a run does.

## 5. Apply and record

Only after all three loops pass and I've signed off: apply, update the provenance
table, and append a `CHANGELOG.md` entry — date, upstream refs, what was adopted,
what was rejected and why. That entry is the baseline the next `/resynthesize` reads,
so closed decisions stay closed and the kit's evolution stays legible. For any
`.bench/learnings.md` entry you acted on, mark it resolved (promoted or dismissed,
with one line of why) so it isn't re-reviewed. The merge is mine; never auto-apply.