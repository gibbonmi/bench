# Review: kit-guidance-fold

Frozen pair: base `21293e39`, reviewed tip `01f873f9`. Raw findings: Standards 10,
Spec 5, Coverage 4, falsification 3. Repair targets after collapse: 9. This file is
pickup state; the repair commit deletes it.

## Standards

Count: 10. Worst issue: the `Run the real path` paragraph in `craft-gate` holds seven
sentences, and the prose check splits it at a label-shaped line, so the bound is not
graded.

- `auto-fix` — `.agents/skills/bench-craft-gate/SKILL.md:43-52`: seven sentences in one
  paragraph; the STE bound is six. Split the paragraph and reflow, so the file stays at
  120 lines and the KG1 fixture `old` stays on one physical line.
- `auto-fix` — `.agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md:23-24`:
  a 24-word procedural sentence; the bound is 20.
- `auto-fix` — `internal/anchors/registry_data_test.go:974-976`, `:1026-1027`,
  `:1087-1088`: doc-comment sentences of 27, 29, and 34 words; the bound is 25.
- `auto-fix` — `.agents/skills/bench-craft-delegate/references/delegation-discipline.md:73-74`:
  "fences" as a verb; the sibling rule says "includes ... in its fence".
- `auto-fix` — `.agents/skills/bench-craft-review/references/finding-discipline.md:4-5`:
  the lead says `craft-review` keeps the citation standard and the refute step, and the
  file then states a citation rule and a refute rule.
- `auto-fix` — `internal/anchors/registry_data_test.go:1151`: the restated count "six"
  drifts when a seventh rule lands; drop the number.
- `auto-fix` — `cross-harness-reviewers.md:26`: the exec form drops `--effort <level>`
  that the bare line 13 carries; restore it, and update the tuple, the test rule, and
  the fixture `review`-side bytes.
- `no-op` — the third sentence of the live-tree probe bullet restates the bullet in
  command form; the kit prefers the exact command token.

## Spec

Count: 5, all `no-op`. Worst issue: story 39 and row KG39 locate the pinned citation
sentence at line 84, and the inserts moved it to line 94; the seam grades the bytes,
not the line.

- `auto-fix` (folded into the repair) — `specs/kit-guidance-fold/spec.md`: story 39 and
  row KG39 name the sentence, not a line number.

## Coverage

Count: 4. Worst issue: the present-but-empty arm of row KG36 is asserted by no test and
no fixture.

- `auto-fix` — `internal/anchors/registry_data_test.go:686-700`: add the arm that writes
  the reference with its lead removed and asserts the dropped-lead diagnostic fires
  while the file-missing diagnostic does not.
- `ask-user`, taken as a Won't handle line — nothing grades that
  `internal/conformance/tier_test.go` and `internal/worktree/parallel_census_test.go`
  still exist; the anchor pins the spelling, not the file. The spec gains one Won't
  handle line, and the reviewer can veto it at the landing.
- `no-op` — the KG14 guard row drives `cat`, not the recipe's own command words.
- `no-op` — the `<<EOF` and `<<-'EOF'` spellings inside an exec span are undecided
  edges the recipes do not name.

## Falsification

Count: 3. Outcomes: accept, dismiss, accept.

- accept, `auto-fix` — `.agents/commands/bench-review-implementation.md:49`: the sentence
  uses "disposition" for accept, merge, or dismiss, and `CONTEXT.md:221` defines
  disposition as the repair-routing label. Reword to "outcome", update the tuple, the
  test rule, and the fixture, and amend story 10 and row KG10.
- dismiss — "exact-record assertion families" is undefined; reviewer decision 7 and a
  Won't handle line keep the term and the hand enumeration.
- accept, `auto-fix` — `cross-harness-reviewers.md:23-28`: the exec form is Claude-only;
  add the parallel Codex line.
