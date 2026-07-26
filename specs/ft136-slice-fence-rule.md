# FT136 — the slice-fence rule

Status: implemented

## Problem

Delegate slices cut by behavior theme instead of ownership fence either
overreach (edit outside their stated boundary) or stop short (return with the
shared piece unwritten), and both cost a round trip. Worse, seam-slicing
optimizes local coherence: when no delegate owns a shared primitive, each
fence writes its own copy — exactly the knowledge-duplication defect the FT86
review existed to fix. Today no kit surface states the rule: `craft-spec` is
silent on how delegate slices are cut, `craft-review` hunts duplication
generally but not at fence boundaries specifically, and `craft-delegate`
charges slices it had no say in shaping.

## Solution

One rule, one source, three surfaces (map `decisions/cost-follows-project-size.md`,
ticket #5, closed): `craft-spec` owns the slicing rule — the slice boundary and
the ownership fence must be the same line, and shared primitives are named up
front and land as a deep-unit slice before the consuming seams. `craft-review`'s
Standards axis adds fence-boundary duplication as an explicit hunt.
`craft-delegate` gets a one-line charge-time cross-reference pointing at
`craft-spec`. The rule is tier-independent (ticket #4): the FT86 evidence —
zero merge conflicts on fence-aligned slices versus ~25 min / 184k tokens on
the theme-cut slice — supports it regardless of delegate tier. Kit edit under
the `craft-synthesis` discipline.

## User stories

1. As a spec author slicing a build for delegates, I want `craft-spec` to state
   that the slice boundary and the ownership fence must be the same line, so
   that a slice neither reaches outside its fence nor stops short of its own
   work.
   Line: `gpt-5.6-sol` / high. `craft-line`'s leverage override routes skill
   guidance prose top because the rule's wording compounds through every spec
   authored after it.
2. As a spec author, I want `craft-spec` to require that shared primitives are
   named up front and land as a deep-unit slice before the consuming seams, so
   that no two fences hand-derive their own copy of the same fact.
   Line: `gpt-5.6-sol` / high. This clause is the countermeasure to the
   knowledge-duplication risk the map's watch-outs name, so its precise wording
   is the value being bought.
3. As a spec author, I want the slicing rule to state its evidence posture in
   one clause — fence alignment is tier-independent — so that a future
   cheap-tier retest (map ticket #6) changes delegate routing, never the rule.
   Line: `gpt-5.6-sol` / high. Same skill, same leverage override; one clause,
   but it is the clause that keeps ticket #6 from reopening the rule.
4. As a reviewer running the Standards axis, I want `craft-review` to name
   duplication introduced at fence boundaries as an explicit hunt, so that the
   defect seam-slicing manufactures is looked for where it appears.
   Line: `gpt-5.6-sol` / high. Review guidance prose compounds the same way —
   a hunt exists only if its name recruits what the reviewer already knows.
5. As a coordinator charging a write-delegate, I want `craft-delegate` to carry
   a one-line cross-reference pointing at `craft-spec`'s slicing rule, so that
   a charge whose slice was never fence-checked gets caught at charge time.
   Line: `gpt-5.6-terra` / medium. A single pointer line with its content
   decided by the map is mechanical placement, and the reviewer has set
   mid-or-higher as the floor for this spec's build delegates.
6. As any session resolving a skill by its index entry, I want `craft-spec`'s
   frontmatter and the generated skills index to advertise the slicing rule,
   so that "how do I slice this for delegates" retrieves the owning skill.
   Line: `gpt-5.6-terra` / medium. Frontmatter phrasing plus a regeneration
   command is mechanical, and the index-equality check makes the gate the
   grader.

## Implementation decisions

- The rule lands in `craft-spec` as a new section after "Story sizing and
  scope cuts" (the existing wide-refactor paragraph stays where it is — it
  governs sequencing, not delegate fences). The section carries the fence
  rule, shared-primitives-first, the tier-independence clause, and one
  contrastive good/bad pair per `craft-skills`, sourced from the FT86
  evidence: fence-aligned package slices merged clean; the theme-cut
  "fail-closed" slice spanned four packages and reached into
  `internal/bounds`.
- `craft-review`'s addition rides the existing Standards bullet where
  knowledge duplication is already named — one sentence extending it to fence
  boundaries with a pointer to `craft-spec`, not a new axis or section.
- `craft-delegate`'s pointer lands in "The charge" (charge time is when the
  fence is checkable), one sentence, no rule content — the rule lives once.
- Edit `.agents/skills/` only; `.claude/skills/` are symlinks (Handoff item 1)
  and the mirroring conformance check guards the invariant.
- `craft-spec`'s `description:` and `index:` frontmatter extend to name
  slicing; `.bench/skills-index.sh --write` regenerates the index in the same
  diff.
- Build-delegate floor: every write-delegate for this spec runs mid tier or
  higher (reviewer directive, 2026-07-26). The stories' lines above satisfy
  it; no cheap-tier routing is permitted even for the mechanical stories.
- The whole edit is one write-delegation in one worktree, stories in order
  1→6 — the three files are one coherent unit (a rule and its two pointers),
  and splitting them across concurrent delegates would put the shared
  primitive (the rule text) behind a fence no pointer-writer owns, the exact
  failure the rule names. The single delegate therefore runs at the maximum
  declared line, `gpt-5.6-sol` / high; stories 5–6's `gpt-5.6-terra` lines
  record what those stories would cost routed alone and bind only if the
  build is ever split. Dogfood note for the synthesis loop: this ordering
  is the deep-unit-first rule applied to its own build.
- Kit edit, so `craft-synthesis` applies: respect closed decisions (tickets
  #4/#5 are closed; do not re-litigate single-skill placement), and run the
  three quality loops — legibility, consistency, dogfood — before handing the
  diff back.

## Testing decisions

- Semantic quality of guidance prose is not gate-observable (Handoff item 4:
  n/a — prose). It is graded by the top-tier falsification pass on this spec
  (reviewer-directed), the `craft-synthesis` loops during the build, and
  `/bench-review-implementation` after it.
- The gate grades structure: frontmatter parse/validity, skills-index
  equality, Claude skill mirroring, stale command references. Prior art: every
  previous kit-prose spec (most recently ft96-delegation-discipline) tested at
  exactly this seam.
- Gate command: `.bench/gate.sh` (the project gate, unchanged).

### Seam diagram

    trigger: /bench-implement-spec write-delegate edits three skill files
        │
        ▼
    .agents/skills/{craft-spec,craft-review,craft-delegate}/SKILL.md
                  ──▶  [ kit content surface ]  ──▶  loaded guidance prose
    .bench/skills-index.sh --write               ──▶  .bench/BENCH-reference.md index
                      ◀ tests attach here: `.bench/gate.sh` conformance phase
                        (parse/validity, index equality, Claude mirroring);
                        semantics attach to review, not the gate

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `craft-spec` states the fence rule (slice boundary = ownership fence) | kit content surface | not TDD-able — prose semantics; graded by the top-tier falsification pass, the synthesis loops, and review | no structural check can tell the rule from its absence; the graders named are the only surfaces that read meaning |
| 2 | `craft-spec` states shared-primitives-first (deep-unit slice lands before consuming seams) | kit content surface | not TDD-able — same posture as story 1 | same |
| 3 | the rule carries its tier-independence clause | kit content surface | not TDD-able — same posture as story 1 | same |
| 4 | `craft-review` Standards axis names the fence-boundary duplication hunt | kit content surface | not TDD-able — same posture as story 1 | same |
| 5 | `craft-delegate` charge section points at `craft-spec`'s rule, content-free | kit content surface | not TDD-able — same posture as story 1; the review's Spec axis additionally checks the pointer carries no rule content | a pointer that restates the rule is the one-source defect this spec exists to prevent |
| 6 | skills index equals regenerated output after the `index:` edit | conformance phase | already covered — the skills-index equality check reds any `index:` edit without regeneration (observed check class; the canary proves it bites) | the check compares committed index to generated output byte-for-byte |
| edge | frontmatter stays parseable after the description edit | conformance phase | already covered — conformance parse/validity checks | malformed YAML frontmatter reds the phase |
| edge | `.claude/skills/` mirror stays consistent | conformance phase | already covered — Claude skill mirroring check | an edit landing on the mirror side instead of `.agents/` reds the phase |
| edge | the craft-delegate pointer names a craft-spec section that exists | kit content surface | not TDD-able — the stale-reference sweep validates command tokens only, not section anchors; the review's Spec axis verifies the anchor resolves | a pointer at a renamed or absent section passes the gate silently, so only a grader that reads the target file catches it |

### Edge inventory

Runtime edge classes (error path, empty input, boundary values, malformed
input, interrupted state, idempotency, hostile environment) — **Won't handle**:
n/a, the deliverable is guidance prose with no runtime contract (Handoff
items 2, 4, 6). The structural edges that do apply — frontmatter validity,
index drift, mirror drift, stale references — are mapped as rows above. The
profile's hostile-input checklist targets the shell CLI surface and does not
apply to a prose-only diff; no new checklist section is proposed.

## Out of scope

- **Cheap-tier retest on a genuinely seam-shaped slice** (map ticket #6) — a
  measurement task with its own trigger, not part of this rule's text;
  opportunistic, ~0 edits, 1–2 gate runs when the slice appears.
- **FT91 conformance parallelization** — separate capability behind tickets
  #2/#3; its spec starts only after #3 closes. Estimate belongs to that spec.
- **FT101 monorepo scoping (docs and profile halves)** — deferred undesigned
  behind its revive trigger (a linked repo with more than one bounded
  context).
