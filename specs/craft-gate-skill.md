# craft-gate — the skill for authoring the oracle

## Problem

The gate is the highest seam in every Bench repo — if it is weak, the whole
system is weak — yet nothing fires when someone adds, edits, weakens, or removes
a gate check. The kit has accumulated real oracle-authoring discipline (prove a
check bites, attribute failures, choose a fail posture, run the real path,
canary the meta-level), but it lives only in `gate.sh` comments, `specs/canary.md`,
and the profile — none of which load when an agent touches a check. Roadmap calls
this the highest-leverage missing craft skill.

## Solution

A model-invoked `craft-gate` skill (`.agents/skills/bench-craft-gate/SKILL.md`,
visible name `craft-gate`) that distills the kit's oracle-authoring discipline,
fires on any gate-check change, and is routed to from the two command phases
where gate authoring actually happens: `/bench-setup-repo` (scaffolding a first
gate) and `/bench-final-check` (a check that looks wrong). Wiring follows the
existing craft-skill conventions; the new cross-references are enforced as prose
anchors so they cannot silently rot.

## User stories

1. As a kit or project dev about to add, edit, weaken, or remove a gate check, I
   want a craft-gate skill that fires on that trigger, so oracle edits follow one
   discipline instead of ad-hoc judgment.
   Line: claude-fable-5 / high. Guidance prose compounds through every session
   that loads it, which is the profile's standing leverage override for skill
   authoring.
2. As a reviewer running `/bench-setup-repo` on a fresh repo, I want the
   gate-scaffolding section to route through craft-gate, so first gates start
   strong instead of accreting checks ad hoc.
   Line: claude-fable-5 / high. Same edit class as story 1 — command guidance
   prose under the leverage override.
3. As an agent inside `/bench-final-check` facing a check that itself looks
   wrong, I want the stop-and-say-so rule to point at craft-gate, so a governed
   check-change path exists instead of a route-around temptation.
   Line: claude-fable-5 / high. Same edit class as story 1.
4. As an agent on any harness, I want the skill wired into the generated skills
   index and the Claude adapter symlink, so it is reachable everywhere the other
   craft skills are.
   Line: inline on the session model, low effort. Mechanical wiring whose
   correctness the gate fully observes.
5. As a kit dev, I want the new cross-references enforced as prose anchors in
   the gate's docs contracts, so the routing survives future edits to those
   command files.
   Line: inline on the session model, low effort. One-pattern extension of the
   existing require_anchor list, gate-observable.

## Implementation decisions

- One new skill directory following the `bench-craft-*` naming contract; visible
  name `craft-gate`; model-invoked (keeps a `description`), like every craft
  skill.
- Content distills what the repo already knows rather than inventing doctrine:
  prove-red (a check you have never seen fail is a hope), attribution (distinct
  message per failure mode; guard checks on the surface they test), run the real
  path against a fixture (never a reimplementation of the checked logic — one
  source per fact), hermetic fixtures (mktemp repos, controlled inputs),
  fail-posture as an explicit decision (fail-open ergonomics vs fail-closed
  enforcement; best-effort only when deliberate), advertisement/enforcement
  pairing, the layered gate shape (parse → structure → conformance → behavior
  contracts → canary), speed as a budget (the gate runs every shift iteration),
  and the authority rule from the author's side (weakening or removing a check
  is a reviewer decision made in the open — never a step inside making
  something pass).
- One contrastive pair (the skill governs an output surface — check code):
  a real-path fixture check with a targeted message beside a reimplementation
  grep with a generic message.
- `/bench-setup-repo` Section A and `/bench-final-check`'s wrong-check rule each
  gain a one-line pointer to craft-gate; both anchored in
  `.bench/gate-docs-contracts.sh` via `require_anchor`.
- Skills index regenerated with `.bench/skills-index.sh --write` (the gate
  verifies committed == generated).

## Testing decisions

- What a good test is here: the gate's existing conformance layer *is* the test
  surface for kit content — frontmatter, `craft-*` visible-name mapping, index
  sync, stale-reference sweep, and the anchor contracts. New behavior gets a red
  signal at that seam before the change that turns it green.
- Seams: the kit content surface and the docs-contract anchor seam (prior art:
  gate checks 3/3b/5a and the `require_anchor` list in
  `.bench/gate-docs-contracts.sh`).
- Gate command: `bench gate`.

### Seam diagram

    trigger: bench gate (every shift iteration, every final-check)
        │
        ▼
    .agents/skills/bench-craft-gate/SKILL.md ──▶ [ gate conformance layer   ] ──▶ exit 0/1 +
    .agents/commands/bench-setup-repo.md     ──▶ [ frontmatter · craft-* map ]     attributed
    .agents/commands/bench-final-check.md    ──▶ [ index sync · anchors      ]     stderr lines
                              ◀ tests attach here: run `bench gate` against the
                                tree with the piece absent (red, targeted message),
                                then present (green)

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | new skill dir appears in the generated skills index with its `index:` trigger line | kit content surface (`.bench/skills-index.sh --check` inside the gate) | `bash .bench/skills-index.sh --check` run red after creating the skill file, before regenerating the committed block | the check diffs committed index against generated; a skill on disk but absent from the block is exactly the drift it names |
| 1 | skill carries frontmatter and the `craft-gate` visible name for the `bench-craft-gate` dir | kit content surface (gate checks 3 and 3b) | already covered — existing checks fire on any nonconforming new `bench-craft-*` dir | a missing fence or wrong visible name errs with the per-skill attributed message |
| 2 | `/bench-setup-repo` names craft-gate in its gate section | docs-contract anchor seam | `bench gate` run red with the new `require_anchor` rows added before the command file gains the reference | the anchor greps the exact needle in the exact file; absence is the failure it reports |
| 3 | `/bench-final-check` names craft-gate in its wrong-check rule | docs-contract anchor seam | same red run as story 2 (both anchors added together, observed red together) | same anchor mechanism, distinct per-file message |
| 4 | Claude adapter symlink resolves to the new skill dir | kit content surface | not TDD-able — no gate check observes per-skill symlinks in this repo; verified by hand (`ls -la .claude/skills/`) and by the linked-repo install path already gate-covered for the `bench-craft-*` glob | the link contracts prove the glob-driven install; a per-skill symlink check does not exist and is not worth adding for a convention `bench link` regenerates |
| 5 | the two new anchors themselves stay in the gate | docs-contract anchor seam | already covered once landed — the anchors run on every gate; their own red-capability was observed in the story-2 red run | an anchor deleted later turns the gate red at the next run with the file+needle message |

### Edge inventory

- error path — a nonconforming skill file: covered by existing gate checks 3/3b
  (coverage row, story 1).
- empty/absent input — skill file present but index not regenerated: coverage
  row story 1.
- malformed input — broken frontmatter: existing gate check 3, already covered.
- re-run idempotency — `skills-index.sh --write` twice: already covered by gate
  docs contract (m).
- interrupted/partial state — **Won't handle:** a half-written SKILL.md between
  edits is a working-tree state the gate grades at the next run; no
  mid-edit protection is meaningful for prose.
- boundary values / hostile environment — **Won't handle:** not meaningful for a
  markdown skill file; the hostile-input checklist targets shell CLI surfaces.

## Out of scope

- **Behavioral skill-firing tests** (does craft-gate actually fire when a check
  is edited) — separate capability, needs a running-agent harness; decided out
  in `decisions/dogfooding.md` #3. Estimate if reopened: ~40 edits, ~20 gate
  runs.
- **craft-review and craft-delegate skills** — same roadmap entry, next in this
  batch, each its own spec. Estimate: ~10 edits, ~4 gate runs each.
