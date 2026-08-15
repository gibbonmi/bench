# Extract the skills-index module and collapse the conformance parsers

Blocked by: none
Writes: internal/skillsindex, internal/conformance/skills_index_checks_test.go, internal/conformance/checks_test.go, internal/conformance/docs_workflow_checks_test.go, internal/conformance/tier_test.go

## What to build

`internal/skillsindex` becomes the one reader of the skills index — frontmatter field,
entries as data, line rendering, block check with today's attributed diagnostics (four
with the `(regenerate: bench skills-index --write)` hint, three without), and an
idempotent 0644-preserving write that refuses on an unparseable allowlist with the
spec's literal. `checkSkillsIndex`, `kitOnlySkillSources`, and `frontmatterField`
collapse to calls into it; `markerBlock` and `checkSkillsIndexGenerateVerify` (and its
call site) are deleted, their generate/verify contract now the module's own tests. A new
`Test*` guard, registered in `classifiedLiveTreeTests`, self-parses the two conformance
files and fails on any surviving marker, allowlist-path, or line-format literal outside
its own function — red before the collapse, green after. The `.bench/skills-index.sh`
script stays untouched here so this ticket lands green alone.

Contract carried to the next ticket: `skillsindex` exposes the check/write entry points
the verb will wrap; the hint literal already names `bench skills-index --write` one commit
before that verb routes (gate-safe: the cold-pickup sweep reads only markdown).

## Acceptance

- [ ] covers SI1, SI2, SI3, SI5: module tests on a temp root observe rendering (adapter skipped, `audience`-before-`source` kit-only marker), the fence rule, the ordered attributed diagnostics, write→check→idempotent write at 0644, and the write refusal literal with unchanged bytes.
- [ ] covers SI4: the four skills-index canaries still bite through `checkSkillsIndex`.
- [ ] covers SI6: the guard test is red on the tree before the collapse and green after; hoisting a literal trips it.
- [ ] covers SI10: no pre-existing assertion changes; the only test edits are the enumerated mechanical ones.
