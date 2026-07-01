# edge case coverage

## Problem

The feature-build path validates work against declared breadth and no stage owns
discovering the cases nobody declared. User stories are happy-path shaped, the
acceptance coverage map inherits that shape, `craft-tdd` forbids testing wider
than the stories, and `/bench-review-implementation` audits only the rows that
exist. The result is visible in this repo's history: hostile-input defects
(paths with spaces, missing trailing newlines, interrupts, symlinks) reached the
contracts only as post-defect regressions.

## Solution

Give the pipeline a case-generation obligation at spec time and a case-gap
review axis at review time. `/bench-write-spec` walks an edge inventory per
mapped behavior and turns each edge into a coverage row or an explicit
"won't handle" line. `craft-tdd` makes stories the breadth floor instead of the
ceiling and scopes the over-fit rule to semantics, not input selection.
`/bench-review-implementation` gains a third Coverage axis that names breaking
inputs no row or test exercises. `craft-seams` caveats the highest-seam rule
with failure-mode observability. The project profile carries a domain
hostile-input checklist the spec phase reads. Gate anchors plus canaries keep
the new contracts from silently rotting.

## User stories

1. As a reviewer, I want every non-trivial spec's Testing decisions to walk an
   edge inventory per mapped behavior — error path, empty/absent input,
   boundary, malformed input, interrupted/partial state, re-run idempotency,
   hostile environment — so untested edges are decisions on the page.
2. As a reviewer, I want each inventoried edge to become either a coverage row
   or an explicit "won't handle" line, so nothing is silently untested.
3. As a spec-writing agent, I want `/bench-write-spec` to consult the project
   profile's hostile-input checklist when one exists, so domain-specific edge
   classes recur instead of being rediscovered per defect.
4. As an implementing agent, I want `craft-tdd` to state that user stories are
   the breadth floor, not the ceiling — at a marked seam I enumerate
   failure-mode variants and propose them as rows the reviewer vetoes, rather
   than silently skipping them.
5. As an implementing agent, I want `craft-tdd` (and the matching line in
   `/bench-implement-spec`) to qualify "let the gate catch regressions": the
   gate only catches what some test observes, and off-seam behavior no seam can
   observe is a seam-set defect to surface, not to skip.
6. As a reviewer, I want `/bench-review-implementation` to run a third
   **Coverage** axis that names concrete inputs or states that would break the
   diff and that no acceptance row or existing test exercises, so coverage gaps
   surface pre-merge instead of as post-defect regressions.
7. As a spec-writing agent, I want `craft-seams` to caveat "use the highest
   seam" with observability — the highest seam at which the failure modes can
   still go red — so error paths don't get dropped for the sake of seam height.
8. As a Bench maintainer, I want this repo's profile to carry the shell-domain
   hostile-input checklist, and `/bench-setup-repo` to prompt new repos for
   their own, so the checklist exists where step 3 looks for it.
9. As a Bench maintainer, I want the gate to anchor the new content contracts
   and canaries to prove the new anchor checks bite, so the contracts cannot
   silently disappear.
10. As a Bench maintainer, I want this change dogfooded — this spec itself
    carries an edge inventory, and the next real feature exercises the full
    changed path — so adoption rests on behavior, not attractive prose.

## Implementation decisions

- **The edge inventory lives inside Testing decisions.** Edge cases that get
  tests become rows in the existing acceptance coverage map (story column may
  read "edge of N"); edges deliberately not handled become one-line
  "Won't handle" entries directly under the map. Both are reviewer veto
  surface, mirroring the Out of scope convention: exclusions are decisions on
  the page, not silent omissions.
- **`craft-tdd` rewording is a scope correction, not a new rule.** The over-fit
  guard protects *what correct means* (reviewer-chosen semantics at
  reviewer-chosen seams); it never restricted *which inputs get exercised*.
  Replace the "breadth comes from the spec's user stories, not what the agent
  could imagine" sentence with floor-not-ceiling language plus the
  propose-for-veto loop.
- **Coverage is a third parallel review axis, kept separate.** Same sub-agent
  pattern and word budget as Standards and Spec, reported under its own
  heading with its own count. Advisory like the others; it does not gate.
- **Anchors stay structural.** Extend gate check j with `require_anchor` lines
  for the new needles; semantic completeness stays a review/dogfood
  responsibility, exactly as the acceptance-coverage anchors decided.
- **Two new canaries, needle-absent with file present.** One for the spec-side
  contract (edge inventory anchor), one for the review-side contract (Coverage
  axis anchor). Fixtures keep the anchor file present so the EXPECT targets the
  "missing anchor" message and survives the vacuous-baseline check.
- **The profile owns the domain checklist.** `projects/benchkit.md` gets the
  shell-CLI hostile-input checklist distilled from this repo's defect history;
  `/bench-setup-repo`'s interview gains one prompt so linked repos start with
  their own.
- **`/bench-debug` is untouched.** The bug path already builds red loops; this
  change is the feature path learning what the bug path already knows.

## Testing decisions

- **Good tests here** exercise the kit content surface as future agents read
  it: the commands and skills carry the instructions, the gate goes red when an
  anchor disappears, and the canaries prove the anchor checks bite.
- **Seams (all existing, zero new):** kit content surface
  (`/bench-write-spec`, `craft-tdd`, `craft-seams`,
  `/bench-review-implementation`, `/bench-implement-spec`,
  `/bench-setup-repo`), the project gate's kit-conformance layer, and the
  canary harness in `tests/canary/`.
- **Prior art:** the acceptance-coverage anchors and the
  `acceptance-coverage-anchor` canary from `specs/tdd-acceptance-coverage.md` —
  this spec extends that exact mechanism.
- **Gate command:** `bench gate`.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1, 2 | `/bench-write-spec` requires an edge inventory per mapped behavior, each edge landing as a coverage row or a "won't handle" line. | Kit content surface: `/bench-write-spec` | Observed red before implementation: `rg -n "edge inventory\|won't handle" .agents/commands/bench-write-spec.md` exits 1. | The command has no edge-inventory language today, so the check fails until the template and process require it. |
| 3, 8 | The spec phase reads the profile's hostile-input checklist; this repo's profile carries the shell checklist; setup-repo prompts for one. | Kit content surface: `/bench-write-spec`, `projects/benchkit.md`, `/bench-setup-repo` | Observed red before implementation: `rg -n "hostile-input checklist" projects/benchkit.md .agents/commands/bench-write-spec.md .agents/commands/bench-setup-repo.md` exits 1. | No surface names the checklist today, so the wiring between spec phase and profile is provably absent. |
| 4 | `craft-tdd` states stories are the breadth floor and requires enumerate-and-propose at marked seams. | Kit content surface: `craft-tdd` | Observed red before implementation: `rg -in "floor\|enumerate" .agents/skills/bench-craft-tdd/SKILL.md` exits 1. | The skill currently caps breadth at the stories; the probe fails until the floor language replaces the cap. |
| 5 | `craft-tdd` and `/bench-implement-spec` qualify the off-seam rule with observability. | Kit content surface: `craft-tdd`, `/bench-implement-spec` | Covered by the story-4 probe for `craft-tdd` (the qualification lands in the same rewritten section); `/bench-implement-spec`'s line is verified in review, not anchored. | A rewrite that restores the floor language but keeps the circular off-seam rule is caught by reviewer read of the same section; anchoring one prose clause per file is the agreed granularity. |
| 6 | `/bench-review-implementation` runs a third Coverage axis naming unexercised breaking inputs. | Kit content surface: `/bench-review-implementation` | Observed red before implementation: `rg -n "## Coverage\|no acceptance row" .agents/commands/bench-review-implementation.md` exits 1. | The review command has two axes today; the probe fails until the third axis exists with its own heading. |
| 7 | `craft-seams` caveats highest-seam with failure-mode observability. | Kit content surface: `craft-seams` | Observed red before implementation: `rg -n "failure mode" .agents/skills/bench-craft-seams/SKILL.md` exits 1. | The skill's "observable behavior" line predates the caveat; the sharper probe fails until the highest-observable-seam language exists. |
| 9 | `bench gate` fails when the new anchors disappear, and canaries prove those checks bite. | Project gate: kit conformance + canary harness | Observed red before implementation: `rg -n "edge inventory\|Coverage axis" .bench/gate.sh` exits 1 and `ls tests/canary/ \| rg "edge\|coverage-axis"` exits 1. | Neither the anchor checks nor the canary fixtures exist yet, so the structural guard is provably absent. |
| 10 | This spec carries its own edge inventory; the next real feature runs the changed path end to end. | Dogfood workflow | Not TDD-able before implementation: the end-to-end proof requires the changed kit to exist first. The spec's own inventory is checkable now (below). | Prevents adopting prose that reads well but fails in use — the same bar the acceptance-map change met. |

### Edge inventory (this spec's own)

- **Anchor file missing entirely** (vs needle missing): already covered —
  `require_anchor` errs distinctly on a missing file.
- **New canary matching an empty fixture** (vacuous EXPECT): already covered —
  the gate's vacuous-baseline check rejects it.
- **Interrupted/partial edit state, re-run idempotency:** not applicable —
  prose edits, no runtime state.
- **Won't handle — needle quoted in an unrelated file:** `require_anchor`
  checks one named file, so a spec quoting the phrase elsewhere is harmless.
- **Won't handle — anchor present but semantically gutted prose:** the
  structural-anchor decision explicitly leaves semantics to review and
  dogfood; a semantic parser is out of scope (below).

## Out of scope

- **Retrofitting edge inventories into historical specs** — separate cleanup of
  records that describe past workflow state, ~1 hour.
- **A semantic checker for edge-inventory completeness** — a distinct
  conformance product needing false-positive design, ~1–2 hours; same cut the
  acceptance-map spec made and parked.
- **A domain checklist library beyond shell** (web/API/UI hostile-input
  checklists shipped as profile templates) — separate capability worth
  preserving; parked on the roadmap, ~1–2 hours.
- **Changes to `/bench-debug`** — the bug path already owns red-loop
  discipline, ~30–45 minutes if ever wanted.
