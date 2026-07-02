# single-source the shared platform sections

## Problem

`specs/shared-rules.md` single-sourced the four invariants and the communication
rules into `.bench/BENCH.md`, and parked the rest of the overlap as a ranked cut.
That rest is still split: the worker/reviewer contract, the how-the-pieces-fit
framing, the process-proportionality and learnings-capture rules, and the skills
index live only in benchkit's `AGENTS.md` — a file consumers never receive. A
linked repo gets the thin managed block, so agents working in consumer repos run
without the platform rules the kit itself runs on. The kit's `AGENTS.md` also
leaks a project name (gl-axi) into content that describes the platform.

## Solution

Finish the move that `shared-rules` started: every rule that holds in *any*
bench-linked repo becomes canonical in `.bench/BENCH.md` (which `bench link`
ships), and benchkit's `AGENTS.md` shrinks to what a consumer's `AGENTS.md` is —
project-owned content plus pointers. The gate's index-sync checks follow the
index to its new home, the duplication guard extends to the moved rules, and a
fresh link gives Claude Code sessions the canonical guide directly.

## User stories

1. As an agent in a consumer repo, I want the worker/reviewer contract ("build on
   the reviewer's behalf; surface decisions that are theirs; never decide what
   ships") in `.bench/BENCH.md`, so I inherit the role definition, not just the
   one-liner "the reviewer owns merge decisions".
   Line: claude-fable-5 / high. This edits shared platform rules that every
   future session loads, which is the leverage override in craft-line and the
   profile's cached routing for guidance prose.
2. As an agent in a consumer repo, I want the how-the-pieces-fit framing in
   `BENCH.md` — skills are probabilistic guidance, commands are the canonical
   phases, the gate and hooks are enforcement with authority I don't have, the
   CLI is the operational layer I drive — so I understand the system, not just
   its file list.
   Line: claude-fable-5 / high. The fold must land cleanly in guidance prose
   that steers every future session, so the leverage override applies on top of
   the gate's blindness to fold quality.
3. As an agent starting a session in any linked repo, I want `BENCH.md` to tell
   me to read `CONTEXT.md` and `projects/<name>.md` when present, so cold pickup
   works the same everywhere.
   Line: claude-fable-5 / high. It rides inside the story-2 fold, so it takes
   the same line for the same reason.
4. As an agent in a consumer repo, I want the process-proportionality rule
   (right-size the process; deviating is the reviewer's call; standing rules
   stick) in `BENCH.md`, merged with its condensed Workflow paragraph so exactly
   one statement exists.
   Line: claude-fable-5 / high. The merged statement becomes the only copy of a
   rule every linked repo runs on, so the leverage override routes it top on
   top of the gate-invisible dedup.
5. As an agent in a consumer repo, I want the learnings-capture rule (when to
   append to `.bench/learnings.md`; capture, never self-rewrite rules) in
   `BENCH.md`, so the journal `bench init` scaffolds actually gets fed.
   Line: claude-fable-5 / high. The rewrite must keep the capture-versus-decide
   authority boundary exact in prose every agent inherits — compounding
   guidance, so the leverage override applies.
6. As an agent on a harness that doesn't auto-load skills, I want the skills
   index in `BENCH.md`, so any linked repo tells me which skill to read when.
   Line: claude-sonnet-4-6 / low. Moving the index is mechanical and the
   repointed index checks catch a missing or dangling entry in both directions.
7. As a maintainer, I want benchkit's `AGENTS.md` to reference the moved sections
   instead of restating them, keeping its required pointer to `BENCH.md`.
   Line: claude-fable-5 / high. The slimmed file stays the template every
   consumer's working agreement mirrors, so it is guidance prose under the
   leverage override even though only pointers and marker absence are
   gate-visible.
8. As a maintainer, I want the gate's index-sync checks (skill on disk ↔ index,
   index ↔ disk, command ↔ index) to enforce against `BENCH.md`, both directions.
   Line: claude-sonnet-4-6 / low. The edit swaps a filename in three existing
   greps and the canary meta-check goes red if the result stops biting.
9. As a maintainer, I want the shared-rule drift guard extended with markers for
   the newly moved rules, so none of them can be copied back into `AGENTS.md`.
   Line: claude-sonnet-4-6 / low. Adding entries to an existing marker list is a
   known shape and the drift canary already proves the mechanism.
10. As a maintainer, I want the `dangling-index` canary updated to the new index
    home, so the repointed check provably still bites.
    Line: claude-sonnet-4-6 / low. The fixture convention is established and the
    gate's vacuous-EXPECT baseline catches a canary that proves nothing.
11. As a consumer running a fresh `bench link`, I want the generated `CLAUDE.md`
    to import `@.bench/BENCH.md` alongside `@AGENTS.md`, so Claude Code loads the
    canonical rules instead of finding them one hop away.
    Line: claude-sonnet-4-6 / low. It is a one-line generator change asserted by
    a link contract that runs the real link in a throwaway repo.
12. As a consumer of the shipped kit, I want the gl-axi references generalized
    out of platform content — the skills-index line and the `bench-craft-cli`
    skill's description and body — so the kit doesn't name one project's tool.
    Line: claude-fable-5 / high. Rewording a skill description changes when the
    skill fires in every consumer repo — the leverage override's core case, with
    no gate check observing firing semantics.
13. As a teammate reading the docs cold, I want stale cross-references updated
    (e.g. README's "CLAUDE.md — one-line import of AGENTS.md"), so no doc
    describes the pre-move layout.
    Line: claude-opus-4-8 / medium. Command references are gate-checked but
    layout prose is not, so accuracy here rests on the model rather than the
    oracle.

## Implementation decisions

- **`BENCH.md` gains two things, folded not appended.** The worker/reviewer
  contract and pieces-fit framing land as a new section near the top (before the
  invariants, which they frame); proportionality + learnings merge into the
  existing Workflow section, replacing its condensed paragraph; the skills index
  becomes a new section. No content is duplicated between sections after the fold.
- **`AGENTS.md` keeps only project-owned framing** — which file is canonical for
  what, the pointer to the shared rules, and the gate-required pointer phrase.
  This intentionally makes the kit's `AGENTS.md` shaped like a consumer's.
- **Gate checks move with the index.** Checks 5a/5b/5c grep `BENCH.md` instead of
  `AGENTS.md`; error messages name the new file. The `ss_markers` list gains one
  stable phrase per moved rule (worker/reviewer, proportionality, learnings
  capture), enforced with the existing present-in-BENCH / absent-from-AGENTS
  shape. The existing pointer check is unchanged.
- **Skills index generalizes as it moves.** The index line drops "(gl-axi)"; the
  `bench-craft-cli` skill's frontmatter description and body references
  generalize to "an AXI-conformant CLI project" phrasing without changing when
  the skill fires.
- **Fresh-link `CLAUDE.md` imports both files.** The generated content becomes
  the two-import form. Relink does not retrofit an existing `CLAUDE.md` (it is
  project-owned once created).
- **Consumer path otherwise unchanged.** `BENCH.md` already ships via the link
  plan and npm `files[]`; no packaging change.

## Testing decisions

- **A good test here** asserts structural properties of the real shipped files —
  where the index lives, that moved rules exist in exactly one place, what a
  fresh link produces — never prose quality, which stays with review.
- **Seams (all existing, none new):**
  1. Kit-conformance checks in `.bench/gate.sh` (5a–5c, shared-rule markers).
     Prior art: the `shared-rules` drift check and the current index checks.
  2. Canary fixtures in `tests/canary/`. Prior art: `dangling-index`,
     `shared-rule-drift`.
  3. Link contracts in `.bench/gate-link-contracts.sh`. Prior art: the fresh-link
     assertions extended earlier today.
- **Gate command:** `.bench/gate.sh` (`bench gate`).

### Acceptance coverage map
| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1–6 | each moved section's marker phrase present in `BENCH.md` | gate `ss_markers` check | add markers before moving content → gate red "shared rule missing from canonical .bench/BENCH.md" (observed shape from shared-rules) | the check greps the shipped file; content that never moved fails it |
| 7, 9 | moved rules absent from `AGENTS.md`; pointer retained | gate `ss_markers` + pointer check | with markers added and `AGENTS.md` not yet slimmed → gate red "shared rule duplicated in AGENTS.md" | duplication is exactly what the marker pair detects |
| 8 | skill/command index sync enforced against `BENCH.md` | gate 5a/5b/5c | repoint checks before moving the index → gate red "skill ... not referenced" | the check reads the new home; an unmoved index fails it |
| 10 | repointed index check still bites | `dangling-index` canary | canary run with old fixture after repoint → "canary did not bite" from the gate's canary harness | the gate's meta-check fails when EXPECT no longer matches |
| 11 | fresh link writes two-import `CLAUDE.md` | link contract | new assertion against current `bench link` → red "fresh link CLAUDE.md does not import .bench/BENCH.md" | asserts generated content before the generator changes |
| 12 | no `gl-axi` in platform content | not TDD-able as a blanket grep — `projects/gl-axi.md` legitimately names it; verified by targeted grep of the index line and craft-cli skill during review | — | — |
| 13 | no stale layout descriptions | already covered for commands by gate check g; prose accuracy stays with review | — | — |

### Edge inventory
- `BENCH.md` missing entirely — covered: every repointed check greps it and errs red (grep on a missing file fails closed).
- Marker phrase copied back into `AGENTS.md` later — covered: story 9 rows; same mechanism the `shared-rule-drift` canary proves.
- Skill referenced in *both* files — **Won't handle:** `AGENTS.md` may freely mention skills for project reasons; only rule markers are duplication-guarded, index presence is not exclusivity.
- Canary fixture dot-dir handling — covered: fixture stores `.bench/` as `dot-bench/` per the existing canary convention; the updated `dangling-index` fixture follows it.
- Relink on a repo with an existing single-import `CLAUDE.md` — **Won't handle:** `CLAUDE.md` is project-owned after creation; retrofitting is a relink-semantics change parked below.
- Fence-wrapped marker phrases in `AGENTS.md` docs — **Won't handle:** the existing marker check is `grep -qF`, not fence-aware, and the moved sections carry no fenced examples; inherits the shared-rules behavior unchanged.

## Out of scope

1. **Auto-generating the skills index from skill frontmatter** — deferred for the
   decision, not the size: it makes the index a generated artifact and changes
   what the gate enforces (sync checks become generation-freshness checks), which
   is a reviewer call. Est ~5–8 min (4 edits, 2–3 gate runs). Parked on the
   roadmap.
2. **Retrofitting existing consumers' `CLAUDE.md` on relink** — deferred for the
   decision, not the size: it lets relink modify a project-owned file it didn't
   create, a change to relink's safety contract. Est ~3–5 min (2 edits, 1–2 gate
   runs).
