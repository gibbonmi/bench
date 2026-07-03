# craft-delegate — the delegation discipline as a skill

## Problem

Delegation guidance is smeared thin across surfaces: invariant #1's delegate
clause in `.bench/BENCH.md` (done-claims, worktree isolation), the routing half
in `craft-line`, and spawn mechanics inside `/bench-review-implementation`.
Nothing fires when the agent actually reaches for a subagent, so delegate
prompts, scope, isolation, and done-claim verification get re-derived per
session — and the roadmap names this the third missing craft skill.

## Solution

A model-invoked `craft-delegate` skill
(`.agents/skills/bench-craft-delegate/SKILL.md`) carrying the operational
discipline: when to delegate at all, how to charge a delegate (a self-contained
prompt), how to scope it, when a worktree is mandatory, and how a done-claim is
verified. It operationalizes invariant #1's delegate clause without restating it
(the invariant stays canonical in `BENCH.md`) and points to `craft-line` for the
model/effort half instead of duplicating it. The review phase's spawn step gains
a pointer so the one place the kit already orchestrates delegates routes through
the skill.

## User stories

1. As an agent about to spawn a subagent — an axis review, a scoped build, a
   fan-out search — I want a craft-delegate skill that fires on that trigger, so
   the charge, scope, isolation, and verification follow one discipline.
   Line: claude-fable-5 / high. Guidance prose under the profile's leverage
   override for skill authoring.
2. As a reviewer, I want write-delegations isolated and done-claims verified as
   a matter of skill-guided habit, not just an invariant clause read at session
   start, so stray delegate edits can't land silently in files I own.
   Line: claude-fable-5 / high. Same edit class — the skill body is the work.
3. As an agent on any harness, I want the skill wired into the generated skills
   index and the Claude adapter symlink.
   Line: inline on the session model, low effort. Mechanical, gate-observed.
4. As a kit dev, I want the review phase's routing to craft-delegate and the
   skill's done-claim rule enforced as prose anchors.
   Line: inline on the session model, low effort. require_anchor pattern.

## Implementation decisions

- Skill carries: the delegate-or-inline decision (context isolation and
  parallelism buy delegation; small mechanical single-file work stays inline);
  the charge (a delegate has no conversation memory — objective, inputs by
  path, the seam, the return shape, the budget, all in the prompt); scope (one
  delegate, one coherent unit; reviewer-owned decisions never delegated);
  isolation (write-delegations in an isolated worktree, read-only ones bare);
  verification (a done-claim is a claim — gate plus `git status` before
  accepting, spot-check citations before trusting a summary); and a pointer to
  `craft-line` for model/effort/cap rather than a restated routing table.
- One contrastive pair: a self-contained charge beside a context-assuming one.
- `/bench-review-implementation` step 3 gains a one-line pointer (spawn per
  `craft-delegate`); its axis-charge routing to `craft-review` is untouched.
- Invariant #1's clause in `BENCH.md` is not moved or trimmed — invariants are
  canonical there and gate-enforced (`ss_markers`); the skill operationalizes,
  it does not restate.
- Wiring identical to the previous two skills: dir naming, index regeneration,
  Claude symlink.

## Testing decisions

- Same seam as the previous two skill specs: the gate's kit-content conformance
  layer and the docs-contract anchor list; red observed before the text lands.
- Gate command: `bench gate`.

### Seam diagram

    trigger: bench gate (every shift iteration, every final-check)
        │
        ▼
    .agents/skills/bench-craft-delegate/SKILL.md    ──▶ [ gate conformance layer ] ──▶ exit 0/1 +
    .agents/commands/bench-review-implementation.md ──▶ [ index sync · anchors   ]     attributed
                                                                                        stderr lines
                              ◀ tests attach here: run `bench gate` with the piece
                                absent (red, targeted message), then present (green)

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | new skill appears in the generated skills index with its trigger line | kit content surface (`.bench/skills-index.sh --check`) | `bash .bench/skills-index.sh --check` run red after creating the skill, before regenerating | a skill on disk but absent from the committed block is the drift the check names |
| 1 | skill carries frontmatter and the `craft-delegate` visible name | kit content surface (gate checks 3/3b) | already covered — existing checks fire on any nonconforming `bench-craft-*` dir | per-skill attributed messages |
| 2 | skill carries the done-claim verification rule | docs-contract anchor seam | `bench gate` run red with the anchor added before the skill file exists | the anchor errs with the file-missing message, then the needle message |
| 4 | review phase routes spawning through craft-delegate | docs-contract anchor seam | same red run (both anchors added together) | anchor greps the needle in the phase file |
| 3 | Claude adapter symlink resolves | kit content surface | not TDD-able — no per-symlink gate check (same call as the prior two skill specs); verified by hand | link contracts cover the `bench-craft-*` glob for consumer installs |

### Edge inventory

- error path / malformed input — nonconforming skill file: existing gate checks
  3/3b (already covered).
- empty/absent — skill on disk, index stale: coverage row story 1.
- re-run idempotency — index `--write` twice: gate docs contract (m), already
  covered.
- interrupted/partial state — **Won't handle:** a mid-edit working tree is
  graded at the next gate run; nothing meaningful to protect for prose.
- boundary / hostile environment — **Won't handle:** not meaningful for a
  markdown skill file.

## Out of scope

- **Delegation enforcement** (a hook that blocks un-isolated write-delegations)
  — separate capability on the hook surface; today only the model-tier half is
  enforceable at the Agent boundary. Estimate: ~15 edits, ~8 gate runs, plus a
  design question (how to detect a write-delegation) that needs shaping first.
- **Behavioral skill-firing tests** — decided out in `decisions/dogfooding.md`
  #3. Estimate if reopened: ~40 edits, ~20 gate runs.
