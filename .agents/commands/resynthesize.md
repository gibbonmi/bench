---
description: Improve the kit from its two distinct sources — upstream sync (pull Matt Pocock's skills + kunchenguid's tooling and diff for changes) and learnings integration (drain .bench/learnings.md from real use). Run either alone or both; each proposes changes that pass three quality loops before adoption. Proposes, never auto-merges. Maintenance, not a workflow phase.
disable-model-invocation: true
---

# /resynthesize — improve the kit, from two separate sources

Bench improves from two independent inputs, and they are not the same job:

- **Upstream sync** — *others'* changes. Pull the two source repos and diff them
  against what Bench already incorporates. External, periodic, version-driven.
- **Learnings integration** — *your own* changes. Drain `.bench/learnings.md`, the
  journal the agent appended to during real use. Internal, usage-driven, has nothing
  to do with upstream.

Pick one or both at the start; don't conflate them. Both feed the same review path
(closed-decisions guard → assess → three quality loops → apply), but they come from
different places and are recorded differently. It **proposes; you merge.**

## 0. Choose the scope

Ask (or take from my request) which to run: **upstream**, **learnings**, or **both**.
If I just say "resynthesize," default to both but report them as two separate sections
so I can see what came from where.

## 1a. Upstream sync (only if upstream is in scope)

- Pull `mattpocock/skills` and kunchenguid's repos (`axi`, `gnhf`, `treehouse`,
  `firstmate`, `no-mistakes`) at HEAD. Record the refs you pulled.
- Read the baseline: the provenance table in `README.md` and prior `CHANGELOG.md`
  entries — what Bench already maps, folds, recommends, or rejected.
- Diff upstream against the baseline into three lists: **new** (techniques not in the
  baseline), **changed** (a mapped/folded item whose source moved), **gone**
  (something Bench uses that was deprecated upstream).

## 1b. Learnings integration (only if learnings is in scope)

- Read `.bench/learnings.md`. Each **open** entry is a candidate change sourced from
  how the kit actually behaved in real work — a deviation, a should-have-asked, a
  recurring friction. This stream never touches the upstream repos; it's the kit
  learning from itself.
- Group entries that point at the same fix so one change resolves several, and note
  any that are one-off context rather than a general rule (those get dismissed, not
  promoted).

Keep the two streams labeled through the rest of the run so every proposed change
carries its origin (upstream vs learnings).

## 2. Respect closed decisions

Drop every candidate the baseline already decided that hasn't materially changed.
For **upstream**: a rejected item (e.g. firstmate's fleet orchestrator) does not
reopen on a cosmetic commit — only a material change reopens it, and then you show the
diff that justifies it. For **learnings**: an entry already marked resolved or
dismissed is not re-litigated. A loop that re-opens settled calls every run is worse
than no loop — it's the failure mode that got the ECC harness rejected.

## 3. Assess — propose, don't decide

For each surviving candidate from either stream, classify in one line: **Map** (fills
a gap the existing layers can't), **Fold** (absorb into an existing piece),
**Recommend** (note it, don't build), **Skip** (with the reason). The bar is the
anti-sediment bar from `bench-craft-skills`: a change earns its place only if it
fills a real gap, never because it's good in isolation. Present the proposed set with
each item tagged upstream or learnings.

## 4. Three quality loops — all must pass before adoption

Run in order; a change that fails a loop is pruned or sent back, not shipped.

1. **Legibility loop.** Run `bench-craft-skills` against each change. No-op?
   Duplicates an existing piece? Pushes the kit past its legibility ceiling? Cut or
   fold it. A bigger kit is a cost, not a feature.
2. **Consistency loop.** Apply to a working copy, then re-run the staleness audit:
   grep for invariant drift, broken cross-references, stale paths, app-specific
   leakage into core files, an out-of-date provenance table. Fix every hit.
3. **Dogfood loop — the oracle.** Run a real shift on a real repo with the changed
   kit: `/bench-spec` a small task, `bench shift`, confirm the gate gates, the hooks fire,
   `bench gate` ends green. A change that reads well but breaks a real run is
   rejected. If you can't run a dogfood shift, the re-synthesis is **not complete** —
   say so rather than adopting on paper. The kit does not grade its own update; a run
   does.

## 5. Apply and record

Only after the loops pass and I've signed off, apply — and record each stream where it
belongs:

- **Upstream** adoptions update the provenance table and append a `CHANGELOG.md` entry
  (date, refs pulled, what was adopted, what was rejected and why). That entry is the
  baseline the next upstream sync reads.
- **Learnings** promotions land in the rule/skill/command they fix, get a `CHANGELOG.md`
  line, and the source entries in `.bench/learnings.md` are marked resolved (promoted
  or dismissed, one line of why) so they're never re-reviewed.

The merge is mine; never auto-apply.