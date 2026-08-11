# Add the prototype skill and the TDD/seams reference leaves

Blocked by: 01-add-craft-domain-skill.md
Ownership fence: `.agents/skills/prototype/SKILL.md`, `.claude/skills/prototype`, `.agents/skills/bench-craft-tdd`, `.agents/skills/bench-craft-seams`, `.bench/BENCH-reference.md`, `CHANGELOG.md`, `internal/conformance`, `internal/anchors/registry_data.go`, `tests/canary/workflow-guidance-anchors`
Integration surfaces: skills index block→`.bench/BENCH-reference.md`; Claude adapter→`.claude/skills/prototype`; adapter classifier→`internal/conformance` (`checkClaudeSkillMirror` in skills_index_checks_test.go currently rejects every non-`bench-craft-*` entry); leaf pointers→`.agents/skills/bench-craft-tdd`; dependency-category leaf→`.agents/skills/bench-craft-seams`; anchors→`internal/anchors/registry_data.go`
Contracts: the `index:` frontmatter trigger crossing `.agents/skills/prototype/SKILL.md`→`.bench/BENCH-reference.md`, asserted by PL1 against the real generator (`.bench/skills-index.sh --check`); the `.claude/skills/prototype` symlink crossing `.claude/skills`→`internal/conformance`, asserted by PL6 against the real mirror check
Closure: PL1/prototype-skill, PL2/tdd-leaves, PL3/seams-leaf, PL4/one-design-it-twice-owner, PL5/discard-not-retain, PL6/mirror-classifier, PL7/anchors-still-green

## What to build

Three additions. (1) A user-invoked `prototype` skill (≤120 lines): one named
question per prototype, trivial to run, state in memory unless persistence is
the question, relevant state surfaced, verdict recorded, artifact discarded —
upstream's throwaway-branch retention (mattpocock/skills @ 84fdeff) is
deliberately not adopted. Add the `.claude/skills/prototype` symlink and index
entry. (2) New `bench-craft-tdd/references/` leaves `tests.md` and `mocking.md`
teaching public-interface behavior, independent expected values, system-seam
mocking, and realistic failure shapes; `craft-tdd`'s SKILL.md points to both
from the branches that need them without duplicating their content (do not push
SKILL.md over 120 lines; ticket 03 owns its deeper rewrite). (3) One new
`bench-craft-seams/references/` deepening leaf covering the four dependency
classes — in-process, local-substitutable, remote-owned, true-external — and
their test strategies; extend the existing `references/design-it-twice.md`
rather than creating a second design-it-twice source. Mirror the leaf layout of
`bench-craft-seams/references/` and `bench-craft-review/references/`.
Preserve the existing craft-seams anchors (`failure modes`,
`structure.budgets`, the file-length-budget clause) and craft-tdd's six anchors.
`checkClaudeSkillMirror` (internal/conformance/skills_index_checks_test.go)
today refuses any `.claude/skills` entry not prefixed `bench-craft-*`; extend
the classifier to admit a user-invoked skill whose symlink resolves to a real
`.agents/skills/<name>/SKILL.md` while still refusing command-adapter
duplication, with a focused test for both postures. CHANGELOG entry for the
new skill and leaves. Won't-handle (explicit unused disposition): prototype
persistence and UI-specific prototype polish stay out per the spec's closed
decision — this ticket adds no persistence or design-source route.

## Acceptance

- [ ] [PL1] (covers PG6) `.agents/skills/prototype/SKILL.md` exists ≤120 lines with frontmatter, index entry, and Claude symlink; body carries named-question, trivial-run, in-memory-state, surfaced-state, recorded-verdict, and discard rules.
- [ ] [PL2] (covers PG4) `bench-craft-tdd/references/tests.md` and `references/mocking.md` exist and are each referenced from `craft-tdd`'s SKILL.md; the leaves teach the four listed behaviors.
- [ ] [PL3] (covers PG5) a `bench-craft-seams` reference leaf enumerates all four dependency classes with a test strategy each, and SKILL.md reaches it.
- [ ] [PL4] (covers local) `rg -l "design it twice|design-it-twice" .agents/skills` still resolves to one owning reference under craft-seams.
- [ ] [PL5] (covers local) the prototype skill states discard-after-verdict; no branch-retention route appears.
- [ ] [PL6] (covers local) `checkClaudeSkillMirror` admits `.claude/skills/prototype` resolving to its SKILL.md, still reds a command-adapter entry and a dangling link, and the focused test proves both.
- [ ] [PL7] (covers local) `go test ./internal/anchors ./internal/conformance` green after the edits; craft-tdd's and craft-seams' existing anchors resolve.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PL1/prototype-skill | strip the `index:` line | skills-index check | run `.bench/skills-index.sh --check`, expect the stale-index red |
| PL2/tdd-leaves | delete `references/mocking.md` | semantic review reread | the SKILL.md pointer dangles; review cites the missing leaf |
| PL3/seams-leaf | drop the remote-owned class from the leaf | semantic review reread | reviewer-graded: enumeration incomplete against PG5 |
| PL4/one-design-it-twice-owner | add a second design-it-twice section to the new leaf | semantic review reread | reviewer-graded: duplicated knowledge against the one-source standard |
| PL5/discard-not-retain | swap discard for keep-on-a-branch | semantic review reread | reviewer-graded against the closed non-adoption |
| PL6/mirror-classifier | point the prototype symlink at a missing directory | mirror-check focused test | break the link, run `go test ./internal/conformance -run SkillsIndex`, expect the broken-adapter red |
| PL7/anchors-still-green | delete a pinned craft-tdd sentence | docs-currency-workflow check | run the anchors check, expect the missing-needle red |
