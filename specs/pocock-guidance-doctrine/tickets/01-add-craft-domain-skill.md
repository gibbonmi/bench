# Add the craft-domain companion skill and charge it from grill, shaping, and spec

Blocked by: none
Ownership fence: `.agents/skills/bench-craft-domain/SKILL.md`, `.claude/skills/bench-craft-domain`, `.bench/BENCH-reference.md`, `.agents/commands/bench-shape-idea.md`, `.agents/commands/bench-write-spec.md`, `.agents/skills/bench-craft-grill/SKILL.md`, `CONTEXT.md`, `CHANGELOG.md`, `internal/anchors/registry_data.go`, `tests/canary/workflow-guidance-anchors`
Integration surfaces: skills index block→`.bench/BENCH-reference.md`; Claude adapter→`.claude/skills/bench-craft-domain`; grill charge→`.agents/skills/bench-craft-grill/SKILL.md`; shaping charge→`.agents/commands/bench-shape-idea.md`; spec charge→`.agents/commands/bench-write-spec.md`
Contracts: the `index:` frontmatter trigger crossing `.agents/skills/bench-craft-domain/SKILL.md`→`.bench/BENCH-reference.md`, asserted by DM2 against the real generator (`.bench/skills-index.sh --check`)
Closure: DM1/skill-exists, DM2/index-entry, DM3/adapter-symlink, DM4/three-charges, DM5/adr-owner-preserved, DM6/glossary-only, DM7/anchors-still-green

## What to build

A new companion craft skill `.agents/skills/bench-craft-domain/SKILL.md` (≤120
lines, YAML frontmatter with `name`, `description`, `index:`) owning: canonical
domain terms and Avoid lists, concrete scenarios at concept seams,
producer-derived equivalence partitions, code-versus-claim comparison, and
inline glossary-only `CONTEXT.md` maintenance. It is charged from exactly three
places — `craft-grill`, `/bench-shape-idea`, `/bench-write-spec` — by one-line
charges added to each; every other phase consumes `CONTEXT.md` ambiently.
Hard-to-reverse decisions route to `craft-adr` via a pointer inside
`craft-domain`; `craft-adr` itself is not edited. Add the `.claude/skills/`
symlink, regenerate the skills index (`.bench/skills-index.sh --write`), and add
a CHANGELOG entry. Mirror frontmatter shape on
`.agents/skills/bench-craft-seams/SKILL.md`.

## Acceptance

- [ ] [DM1] (covers PG1) `.agents/skills/bench-craft-domain/SKILL.md` exists as a real regular file ≤120 lines with conformant frontmatter.
- [ ] [DM2] (covers local) the generated skills index carries a domain-modeling entry sourced from the skill's `index:` frontmatter; `.bench/skills-index.sh --check` passes.
- [ ] [DM3] (covers local) `.claude/skills/bench-craft-domain` is a symlink resolving to `../../.agents/skills/bench-craft-domain`.
- [ ] [DM4] (covers PG2) `craft-grill`, `bench-shape-idea.md`, and `bench-write-spec.md` each charge the skill by name at their entry, and the skill body enumerates the five owned behaviors (terms/Avoid lists, concept-edge scenarios, producer-derived partitions, code-versus-claim, inline glossary updates). Semantic quality of scenarios is the reviewer-graded exception PG2 records.
- [ ] [DM5] (covers PG3) `craft-domain` points hard-to-reverse decisions at `craft-adr` and states no ADR format of its own; `.agents/skills/bench-craft-adr` is untouched in the diff.
- [ ] [DM6] (covers local) any `CONTEXT.md` edit stays glossary-only (terms and definitions, no workflow prose).
- [ ] [DM7] (covers local) `go test ./internal/anchors ./internal/conformance` is green after the edits — the existing anchors on `craft-grill`, `bench-shape-idea.md`, `bench-write-spec.md`, `CHANGELOG.md`, and `CONTEXT.md` still resolve (charge lines are additions; if a pinned sentence must move, migrate its anchor).

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DM1/skill-exists | delete the SKILL.md | skills-index check | remove the file, run `.bench/skills-index.sh --check`, expect the missing-source red |
| DM2/index-entry | strip the `index:` frontmatter line | skills-index check | remove the line, run `.bench/skills-index.sh --check`, expect the stale-index red |
| DM3/adapter-symlink | remove the `.claude/skills` symlink | semantic review reread | `test -L .claude/skills/bench-craft-domain` fails; review cites the missing adapter |
| DM4/three-charges | delete the charge line from one of the three surfaces | semantic review reread | `rg bench-craft-domain` over the three files returns two hits, not three |
| DM5/adr-owner-preserved | add an ADR template to craft-domain | semantic review reread | reviewer-graded: duplicate format owner cited against PG3 |
| DM6/glossary-only | add workflow prose to CONTEXT.md | semantic review reread | reviewer-graded against the glossary-only decision |
| DM7/anchors-still-green | delete a pinned sentence adjacent to a charge insertion | docs-currency-workflow check | run `go test ./internal/conformance -run Anchors`, expect the missing-needle red |
