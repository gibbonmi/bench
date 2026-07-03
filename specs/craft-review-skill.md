# craft-review — the three-axis review judgment as a skill

## Problem

The three-axis review judgment (Standards / Spec / Coverage: what each axis
hunts, what a finding must cite, why the axes stay separate) lives only inside
`/bench-review-implementation` — a reviewer-invoked phase. Nothing fires when the
agent reviews anything *outside* that phase: a delegate's returned diff, a PR, a
self-review before commit. And the axis charges are buried in a procedure file,
so the phase's sub-agent prompts and any ad-hoc review each re-derive them.

## Solution

A model-invoked `craft-review` skill (`.agents/skills/bench-craft-review/SKILL.md`)
that is the one source for the axis charges and the finding standard. The phase
keeps the procedure — pinning the diff with `bench diff`, finding sources,
spawning one delegate per axis, aggregating without merging — and routes each
delegate's charge to the skill instead of restating it.

## User stories

1. As an agent reviewing any diff — the review phase, a delegate's done-claim, a
   PR, a pre-commit self-review — I want a craft-review skill that fires on that
   trigger, so findings meet one citation standard everywhere.
   Line: claude-fable-5 / high. Guidance prose under the profile's leverage
   override for skill authoring.
2. As a kit maintainer, I want the axis charges single-sourced in the skill with
   the phase reduced to procedure plus a pointer, so the two surfaces cannot
   drift apart.
   Line: claude-fable-5 / high. Same edit class — command and skill prose.
3. As an agent on any harness, I want the skill wired into the generated skills
   index and the Claude adapter symlink.
   Line: inline on the session model, low effort. Mechanical, gate-observed.
4. As a kit dev, I want the phase→skill routing and the skill's load-bearing
   coverage rule enforced as prose anchors, so the single-sourcing survives
   future edits.
   Line: inline on the session model, low effort. require_anchor pattern.

## Implementation decisions

- Skill carries: the advisory frame (review says good, the gate says done, the
  reviewer decides ships), the three axis charges, the separation rationale (one
  axis can mask another), the finding citation standard per axis, knowledge
  duplication as a Standards finding (the repo code standard review grades
  against), the refute-before-report rule for Coverage findings, and one
  contrastive finding pair.
- The phase's step 3 shrinks to procedure: spawn one delegate per axis, each
  charged from the skill; the procedural inputs (spec file, `bench coverage`
  rows, no-spec fallback) stay in the phase. The gate's existing anchors on the
  phase ("acceptance coverage map", "mapped behavior", "Coverage axis",
  "## Coverage", "bench diff") all remain satisfied by the phase text — no
  anchor is weakened or moved.
- Two new anchors: the phase must name `craft-review`; the skill must carry the
  adversarial coverage rule ("an edge nobody decided is").
- Wiring identical to craft-gate: `bench-craft-review` dir, visible name
  `craft-review`, index regenerated, Claude symlink.

## Testing decisions

- Same seam as craft-gate-skill: the gate's kit-content conformance layer and
  the docs-contract anchor list. Red observed at the seam before the text lands.
- Gate command: `bench gate`.

### Seam diagram

    trigger: bench gate (every shift iteration, every final-check)
        │
        ▼
    .agents/skills/bench-craft-review/SKILL.md   ──▶ [ gate conformance layer  ] ──▶ exit 0/1 +
    .agents/commands/bench-review-implementation.md ─▶ [ index sync · anchors  ]     attributed
                                                                                      stderr lines
                              ◀ tests attach here: run `bench gate` with the piece
                                absent (red, targeted message), then present (green)

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | new skill appears in the generated skills index with its trigger line | kit content surface (`.bench/skills-index.sh --check`) | `bash .bench/skills-index.sh --check` run red after creating the skill file, before regenerating | a skill on disk but absent from the committed block is exactly the drift the check names |
| 1 | skill carries frontmatter and the `craft-review` visible name | kit content surface (gate checks 3/3b) | already covered — existing checks fire on any nonconforming `bench-craft-*` dir | per-skill attributed frontmatter/name messages |
| 2 | the phase names craft-review as the source of the axis charges | docs-contract anchor seam | `bench gate` run red with the new anchor added before the phase gains the reference | the anchor greps the needle in the file; absence is the reported failure |
| 2 | the phase keeps its five existing anchored strings after the shrink | docs-contract anchor seam | already covered — the existing anchors run on every gate and were green before this change | any anchored string lost in the edit turns the gate red with its file+needle message |
| 4 | the skill carries the adversarial coverage rule | docs-contract anchor seam | same red run as story 2 (both anchors added together, observed red together) | same mechanism, distinct message naming the skill file |
| 3 | Claude adapter symlink resolves | kit content surface | not TDD-able — no per-symlink gate check exists (same call as craft-gate-skill spec); verified by hand | linked-repo installs are covered by the link contracts' `bench-craft-*` glob |

### Edge inventory

- error path / malformed input — nonconforming skill file: existing gate checks
  3/3b (already covered).
- empty/absent — skill on disk, index stale: coverage row story 1.
- re-run idempotency — `skills-index.sh --write` twice: gate docs contract (m),
  already covered.
- interrupted/partial state — **Won't handle:** a mid-edit working tree is
  graded at the next gate run; no mid-edit protection is meaningful for prose.
- boundary / hostile environment — **Won't handle:** not meaningful for a
  markdown skill file; the hostile-input checklist targets shell CLI surfaces.

## Out of scope

- **Reworking the review phase's delegate mechanics** (worktree isolation,
  delegate prompt shape, ~60k routing) — that is craft-delegate's ground, next
  in this batch. Estimate: ~10 edits, ~4 gate runs.
- **Behavioral skill-firing tests** — decided out in `decisions/dogfooding.md`
  #3. Estimate if reopened: ~40 edits, ~20 gate runs.
