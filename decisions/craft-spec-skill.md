# Craft-Spec Skill (FT23)

## #1: What does the skill own, and what does the command keep?

Blocked by: —
Type: Grill

### Question
The acceptance-coverage-map row schema (behavior / seam / red signal) and the
edge-inventory classes are owned by the `/bench-write-spec` command file,
which model-invoked skills cannot auto-load: `craft-tdd` and `craft-review`
point at it, but an ad-hoc TDD pass or self-review fires without reach to the
schema. Second consecutive assessment flagged it.

### Answer
New skill `bench-craft-spec` owns the spec-authoring discipline that other
contexts need: the coverage-map row schema, the edge-inventory classes, and
the story-sizing rules. `/bench-write-spec` keeps the phase mechanics (entry
contract, roadmap interaction, staging, sign-off) and composes the skill by
pointer for the schema — one source, the skill. `craft-tdd` and `craft-review`
repoint from the command file to the skill. Wiring follows the standing skill
conventions: `.agents/skills/bench-craft-spec/SKILL.md`, `.claude/skills`
symlink, skills-index row in BENCH-reference, description with concrete
trigger phrases. Rejected: duplicating the schema into craft-tdd/craft-review
(the drift this repo's one-source rule exists to prevent); moving the whole
write-spec phase into a skill (phases are reviewer-invoked contracts, not
model-triggered guidance).

## Handoff

1. **Module boundaries.** New `.agents/skills/bench-craft-spec/` owns schema +
   edge classes + sizing; `bench-write-spec.md` keeps phase mechanics and
   loses the schema text; `craft-tdd`/`craft-review` update one pointer each.
2. **Contracts.** The skill is model-invocable (no
   `disable-model-invocation`); its description triggers on spec authoring,
   coverage rows, and edge inventories; the command's schema section becomes a
   pointer to the skill file path.
3. **Deep vs thin.** The skill body is moved text plus a short framing — the
   depth already exists in the command file; this is an ownership move.
4. **Black-box assertables.** Conformance can assert: schema text appears in
   exactly one file (the skill); the command, craft-tdd, and craft-review
   contain the pointer; symlink and index rows present (existing skills-index
   `--check` covers sync).
5. **Gate attachment.** The one-source conformance family (FT18 built the
   schema-owner check — it moves to point at the skill) plus the skills-index
   check.
6. **Hostile-input owners.** n/a — prose artifact, no input surface.
7. **Uncertainty flags.** Whether the existing FT18 conformance check pins the
   schema's current owner by path — if so, the check moves in the same diff or
   the gate goes red between commits; implementer sequences it atomically.
8. **Rejected alternatives.** Schema duplication; whole-phase skill;
   leaving reachability to the command pointer (two assessments proved
   nobody follows it under model invocation).
9. **Domain watch-outs.** Leverage artifact — every future spec author loads
   it; wording defects multiply. Top-tier authoring per the craft-line
   override.

Dependency order: n/a — single spec.
