# Craft-spec skill — own the coverage-map schema, edge classes, and story sizing

Status: staged

## Problem

The acceptance-coverage-map row schema, the edge-inventory classes, and the story-sizing rules live inside the `/bench-write-spec` command file. Model-invoked skills cannot auto-load a command file, so an ad-hoc TDD pass (`craft-tdd`) or a self-review (`craft-review`) that reaches for the schema fires with a pointer at prose it cannot open. Two consecutive kit assessments flagged that nobody follows the pointer under model invocation — the schema is effectively unreachable exactly when it is needed off the reviewer-run phase.

## Solution

A new model-invocable skill, `bench-craft-spec`, becomes the single owner of the spec-authoring discipline other contexts need: the coverage-map row schema (five fields plus the red-signal definition and its classifications), the canonical edge-inventory classes, and the story-sizing rules. `/bench-write-spec` keeps every phase mechanic — entry contract, seam-picking, roadmap and staging interaction, routing, retire, sign-off, and the spec template — and its schema and edge-class definitions collapse to a one-line pointer at the skill. `craft-tdd` and `craft-review` repoint their existing pointers from the command to the skill. The skill is wired like every other `bench-craft-*` skill: `index:` frontmatter, a generated skills-index row, and a `.claude/skills` symlink. The FT18 schema-owner conformance anchors move to guard the skill as the new owner, and pointer anchors are added for the command and the two skills — all in the same commit as the prose, so the gate never reds between commits.

## User stories

1. As a spec author running an ad-hoc TDD pass or a self-review, I want the acceptance-coverage-map row schema — the five fields, the red-signal definition, and the already-covered / not-TDD-able classifications, together with the degenerate-implementation and quantifier-enumeration rules and the black-box-assertable-to-row rule — to live in the model-invocable `bench-craft-spec` skill, so that the schema loads on trigger without reaching into the reviewer-only command file. Line: claude-fable-5 / high. This is leverage prose every future spec author loads, so it takes the craft-line authoring override.
2. As a spec author, I want the canonical edge-inventory classes — error path, empty or absent input, boundary values, malformed input, interrupted or partial state, re-run idempotency, hostile environment — plus the every-edge-lands-as-a-row-or-Won't-handle rule, the amputated-caller check, and the official-format compatibility-probe rule owned by the same skill, so that `craft-tdd` and `craft-review` reach one canonical class list instead of the command file. Line: claude-fable-5 / high. It is the same leverage surface as the schema and is authored in the same pass.
3. As a spec author, I want the story-sizing rules — user stories as the breadth floor, the sub-threshold no-deferral rule, and the separate-capability-not-the-rest-of-this-one cut test — owned by the skill, so that sizing discipline is reachable when the phase file is not loaded. Line: claude-fable-5 / high. It is leverage prose folded into the same authoring pass.
4. As any harness, I want `bench-craft-spec` discoverable — name and description frontmatter with concrete spec-authoring triggers, an `index:` field, its generated row in the BENCH-reference skills index, and a `.claude/skills/bench-craft-spec` symlink — so that Claude Code auto-loads it and every other harness finds it in the index. Line: claude-fable-5 / high. The description is leverage prose that decides whether the skill fires, and the mechanical symlink and index regeneration ride along in the same pass rather than earning a separate tier.
5. As the reviewer running `/bench-write-spec`, I want the command's step-4 schema definition and step-5 edge-class list collapsed to a one-line pointer at `bench-craft-spec`, while it keeps the map-required entry contract, the template skeleton, the phase-step names, the seam-diagram convention, and the profile hostile-input-checklist instruction, so that the phase mechanics stay put and the schema has exactly one owner. Line: claude-fable-5 / high. Trimming the phase prose without dropping an anchored string is judgment the gate cannot fully grade.
6. As a spec author in a TDD pass, I want `craft-tdd`'s row-schema-and-red-signal pointer repointed from `/bench-write-spec` to `bench-craft-spec`, so that the pointer names a surface the model can actually load. Line: claude-fable-5 / high. It is a one-line prose repoint on a leverage skill.
7. As a reviewer running the Coverage axis, I want `craft-review`'s edge-class pointer repointed from `/bench-write-spec` to `bench-craft-spec`, so that the Coverage axis reaches the class list without the command file. Line: claude-fable-5 / high. It is a one-line prose repoint on a leverage skill.
8. As the gate, I want the schema-owner conformance anchors to guard `bench-craft-spec` as the new owner and pointer anchors added for the command, `craft-tdd`, and `craft-review`, with every anchor edit landing in the same commit as the prose it checks and no existing anchored string dropped, so that `TestRootConformance` stays green at a single commit with no red between commits. Line: claude-opus-4-8 / medium. Correctness of the oracle move outweighs speed, and the atomic sequencing is the map's one flagged uncertainty.

## Implementation decisions

- **New owner.** `.agents/skills/bench-craft-spec/SKILL.md` is created as a model-invocable craft skill (no `disable-model-invocation`). It owns three bodies of knowledge moved verbatim-plus-framing out of `/bench-write-spec`: the coverage-map row schema and red-signal definition (from the command's acceptance-coverage-map step), the edge-inventory classes (from the walk-the-edge-inventory step), and the story-sizing rules (breadth floor plus the price-every-cut deferral rules). The depth already exists in the command; this is an ownership move, not new invention.
- **Frontmatter.** `name: bench-craft-spec`; a `description:` whose trigger phrases name spec authoring, coverage-map rows, red signals, edge inventories, and story sizing; and an `index:` phrase for the generated skills-index row. `craft-skills` and `craft-adr` govern the authoring.
- **Command keeps phase mechanics.** `.agents/commands/bench-write-spec.md` retains its Entry orientation and Exit handoff, the map-required refusal, seam-picking off the Handoff, routing, staging, retire (promote-then-delete, `spec-retire:`, ROADMAP-row removal), sign-off, and the full spec Template. Its step-4 schema definition and step-5 edge-class list are replaced by a one-line pointer to `bench-craft-spec`. The template's section headers and column labels, the seam-diagram block, the `hostile-input checklist` instruction, and every currently-anchored phase-step string stay in place.
- **Two repoints.** `craft-tdd`'s row-schema-and-red-signal owner sentence and `craft-review`'s edge-class owner sentence change their owner token from `/bench-write-spec` to `bench-craft-spec`; the surrounding anchor phrases are preserved.
- **Wiring.** A committed `.claude/skills/bench-craft-spec` symlink to `../../.agents/skills/bench-craft-spec`, and the skills-index block in `.bench/BENCH-reference.md` regenerated from frontmatter with `.bench/skills-index.sh --write`.
- **Conformance move (atomic).** In `internal/conformance/docs_workflow_helpers_test.go`, `checkWorkflowAnchors` gains anchors requiring the schema, edge-class, and sizing text in `.agents/skills/bench-craft-spec/SKILL.md`, and pointer anchors requiring the string `bench-craft-spec` in the command, `craft-tdd`, and `craft-review`. No existing anchored string is removed — the command's phase-step and template anchors (and their `acceptance-coverage-anchor` / `edge-inventory-anchor` canary fixtures) stay valid. Every anchor edit ships in the same commit as the prose it guards.

## Testing decisions

- **What a good test is here.** The oracle is the gate's observable verdict against the tree — a conformance diagnostic, an index-sync failure, a broken mirror — never a reading of the prose's wording. A row is red-first: add the guarding anchor, run the check, watch it fail because the guarded content is absent, then author the content.
- **Seams and prior art.** All rows attach at the conformance suite (`checkWorkflowAnchors`, `checkSkillsIndex`, `checkClaudeSkillMirror` under `TestRootConformance`), the `.bench/skills-index.sh --check`/`--write` contract (already covered by `checkSkillsIndexGenerateVerify`), and the `workflow-guidance-anchors` / `skills-index-command-adapters` / `load-validity-metadata` canary families. Prior art: the existing anchor table in `checkWorkflowAnchors`, and the `acceptance-coverage-anchor`, `edge-inventory-anchor`, `unindexed-skill`, and `claude-skills-unmirrored` fixtures.
- **Gate command.** `.bench/gate.sh` (the project gate).

### Seam diagram

    trigger: bench gate  →  gate-phases conformance phase
        │
        ▼
    .agents/skills/bench-craft-spec/SKILL.md   ──▶
    .agents/commands/bench-write-spec.md        ──▶  [ RunConformance ]  ──▶  exit 0 / diagnostics
    .agents/skills/bench-craft-tdd/SKILL.md     ──▶  [ checkWorkflow   ]
    .agents/skills/bench-craft-review/SKILL.md  ──▶  [ Anchors +       ]
    .bench/BENCH-reference.md (index block)     ──▶  [ SkillsIndex +   ]
    .claude/skills/bench-craft-spec (symlink)   ──▶  [ ClaudeMirror    ]
                      ◀ tests attach here: run the conformance suite (or bench gate)
                        against the tree and assert green; assert red on a half-applied move

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | skill body carries the coverage-map row schema and red-signal definition | conformance suite, checkWorkflowAnchors | go test -run TestRootConformance reds: skill file missing the acceptance-coverage anchor (`why it catches the failure`) | the skill-owner anchor fails until the schema text is authored into the skill |
| 2 | skill body carries the canonical edge-inventory classes | conformance suite, checkWorkflowAnchors | go test -run TestRootConformance reds: skill file missing the edge-class anchor (`re-run idempotency`) | the edge-class anchor fails until the class list is authored into the skill |
| 3 | skill body carries the story-sizing rules | conformance suite, checkWorkflowAnchors | go test -run TestRootConformance reds: skill file missing the sizing anchor (`separate capability`) | the sizing anchor fails until the sizing rules are authored into the skill |
| 4 | skill is indexed and Claude-mirrored | checkSkillsIndex and checkClaudeSkillMirror | gate reds: skills index missing an entry for `bench-craft-spec`, and the skill missing from `.claude/skills` | the index-sync and mirror checks red until the index row is regenerated and the symlink exists |
| 5 | command collapses its schema section to a pointer and keeps its phase anchors | conformance suite, checkWorkflowAnchors | go test -run TestRootConformance reds: command file missing the `bench-craft-spec` pointer anchor | the command-pointer anchor fails until the command names the skill, and a dropped phase anchor reds its own existing require |
| 6 | craft-tdd points at the skill for the row schema | conformance suite, checkWorkflowAnchors | go test -run TestRootConformance reds: craft-tdd missing the `bench-craft-spec` pointer anchor | the tdd-pointer anchor fails until craft-tdd names the skill |
| 7 | craft-review points at the skill for the edge classes | conformance suite, checkWorkflowAnchors | go test -run TestRootConformance reds: craft-review missing the `bench-craft-spec` pointer anchor | the review-pointer anchor fails until craft-review names the skill |
| 8 | anchor edits and prose land in one commit, no red between commits | full gate at the single staged commit | a deliberately split staging reds go test -run TestRootConformance: skill anchor added before the skill body exists | proves the atomic sequencing, since a split leaves either the skill-owner or a pointer anchor unsatisfied at the intermediate commit |

Cheapest wrong implementation check: an always-green stub that leaves the schema in the command and only adds the skill file without the moved body would pass no row 1-3 anchor (the skill lacks the text) and would leave the pointer anchors of rows 5-7 red; a stub that adds anchors without content reds rows 1-3; a non-atomic split reds row 8. No degenerate implementation passes the map.

### Edge inventory

Edge classes walked per behavior; each resolves to a row above or a Won't-handle line here.

- error path — a half-applied move (prose moved, anchor stale, or anchor added before content) → row 8.
- interrupted or partial state — same atomic-commit concern → row 8.
- invocation through a symlink — the `.claude/skills` entry is a symlink that must resolve to `SKILL.md` → row 4 (`checkClaudeSkillMirror` asserts resolution).
- re-run idempotency — **Won't handle** — `.bench/skills-index.sh --write` idempotency and check-after-write are already guarded by the existing `checkSkillsIndexGenerateVerify`; this change adds no new generation path.
- empty or absent input, boundary values, malformed input, hostile environment — **Won't handle** — `bench-craft-spec` is a prose artifact with no runtime input surface, so these classes have no target (Handoff item 6: n/a).
- stale command copies inside the two workflow-anchor canary fixtures — **Won't handle** — canary EXPECT matching is substring-based, so the new command-pointer anchor firing an extra diagnostic on a stale fixture copy is tolerated; re-syncing the copies is cosmetic hygiene, not a bite change.
- a dedicated canary fixture for the skill-owned schema anchor — **Won't handle** — the require-present mechanism is already bite-proven by the `workflow-guidance-anchors` fixtures, so a per-string fixture proves no new mechanism.
- Codex-adapter or slash-menu false positive — **Won't handle** — `bench-craft-spec` has no matching `.agents/commands/` file, so `checkCodexCommandAdapters` and the slash-duplication guard do not iterate it; it is treated exactly like the other `bench-craft-*` skills (verified against both checks).
- description does not fire on model invocation — **Won't handle** — frontmatter presence is gate-checked, but trigger-phrase quality is `craft-skills` review-graded, not gate-observable; the description is authored under that skill in story 4.

## Out of scope

- **A standing single-source guard for the schema** — a new conformance check asserting the schema strings appear in exactly one `.agents/` file, forbidding future re-duplication back into the command or the two skills. It is a separate hardening capability (a new guard plus its canary bite fixture), not the rest of this move: this spec achieves single-source by relocating the text and repointing, and the closed map left the exactly-one-file assertion optional, naming only the repointed anchor and the skills-index check as the attachment. Estimate to build later: 4 edits, 3 gate runs.
- **Re-syncing the two stale workflow-anchor canary fixture copies** of `bench-write-spec.md` to the collapsed version — a cosmetic hygiene pass on the fixture tree, separable because substring EXPECT matching keeps both fixtures biting without it. Estimate: 2 edits, 2 gate runs.
