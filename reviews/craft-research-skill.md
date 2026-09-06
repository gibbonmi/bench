# Review pickup: craft-research-skill

Frozen base `3f29857f7c8a1e52807f48a16c27179586fe305c`, reviewed tip `2dcb445eee9e12a93de762910f8db30dc1182aa9`. Axes ran at `opus` (Standards and Spec at low, Coverage at medium). The Codex falsification pass ran at `gpt-5.6-terra`. Raw findings: 8. Repair targets after de-duplication: 5. Open reviewer question: 1.

## Standards

Count 2. Worst: the assessment command derives the unknowns rule in three places.

- S1 `auto-fix`. `.agents/commands/bench-assess.md` step 3 keeps the sentence that starts "Mark each re-verified claim with a ✓" and ends "as an unknown in the verification notes." Step 4's verification-notes bullet and the skill both own that rule. Rule: `AGENTS.md`, "Two derivations of the same fact must collapse into one source." Repair: cut the sentence from step 3.
- S2 `auto-fix`. `internal/anchors/registry_craft_research_test.go` opens with a ten-line doc comment that recites every needle listed ten lines below. Rule: `craft-comments`, match the file's comment density. Repair: keep the independence sentence, cut the recital.

## Spec

Count 1. Worst: CR7 places one element outside the report section.

- P1 `no-op`. CR7 lists "every question answered or unknown" under the report section. The skill states it under Sources and verification as the completion rule. The rule is present once, where it belongs.
- Codex F3 `no-op`. The destination precedence and the three metadata labels have no anchor needle. The spec's Further notes fix the anchored sentence list, and the precedence is not in it. CR6 is graded by the one-output needle, and the precedence stays review-owned.

## Coverage

Count 2. Worst: C1, a coverage row claims a catch the seam cannot make.

- C1 `ask-user`, open for the reviewer. Four of the five anchored phrases in `.agents/commands/bench-assess.md` also sit in the frontmatter `description:` line, and the anchor match is a whole-file substring. A body rewrite that drops a phrase stays green while the description keeps it, and the probe proved it. CR15's why-it-catches clause is therefore not true for a body-only change, and the condition predates this diff. Repair shape if accepted: `RequireInSection` rows for those four phrases, or body-unique needles. Not repaired here.
- C2 `auto-fix`. `internal/conformance/registry_test.go` names the anchor registry files that own the `workflow-guidance-anchors` fixtures, and the list lacks `internal/anchors/registry_craft_research.go`. Repair: one-line addition; the fence extends to that file under Build decisions.

## Codex falsification

Count 3. Accepted 2, dismissed 1 (F3 above).

- F1 `auto-fix`. `.agents/skills/bench-craft-spec/references/bootstrap-authority.md` says "Charge `craft-research` for the runnable probe that a compatibility claim needs." The skill excludes prototypes and gives the probe to the calling phase. Repair: reword the pointer so it points at the skill's probe rule, and update the anchored needle in the registry and its mutation rule.
- F2 `auto-fix`. `.agents/skills/bench-craft-research/SKILL.md` says "Route anything with write access or a done-claim through `craft-delegate`," and later says only the coordinator writes the output. Repair: reword the first sentence so it names a write delegate and a done-claim, not every write.

## Build decision recorded in the spec

The spec said each needle gets a fixture under the `workflow-guidance-anchors` family. Fourteen needles landed with two fixtures. The mutation table is the per-needle proof, and the fixture family holds one fixture per edited subject file. The repair adds the missing bootstrap-authority fixture. The decision is recorded in the spec's Build decisions for veto.
