# Repair the review findings

Blocked by: add-the-finding-discipline-reference-to-craft-review.md
Writes: .agents/skills/bench-craft-gate/SKILL.md, .agents/commands/bench-review-implementation.md, .agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, .agents/skills/bench-craft-review/references/finding-discipline.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go, tests/canary/workflow-guidance-anchors, tests/canary/workflow-guidance-anchors/review-falsification-dispositions, tests/canary/workflow-guidance-anchors/cross-harness-reviewer-exec-heredoc, specs/kit-guidance-fold/spec.md, specs/kit-guidance-fold/tickets/repair-the-review-findings.md (new)
Covers: KG10, KG12, KG36, KG39

## What to build

The review of this spec accepted nine repairs. This ticket lands them, and it amends the
four coverage rows the repairs move.

The review phase file gives a falsification finding an outcome, not a disposition.
`CONTEXT.md` reserves disposition for the repair-routing label, so the two labels no
longer share one word. The tuple, the test rule, and the fixture take the new word.

The cross-harness recipes carry the exec form for both families. The Claude exec line
regains `--effort <level>`, and a parallel Codex exec line lands after it. The exec child
already runs inside the worktree, so the Codex line takes `-C .`. The intro sentence
drops to 20 words and keeps its reason. The two bare recipe lines stay byte-identical.

The `Run the real path` section of `craft-gate` holds two paragraphs. The first paragraph
ends at the second-derivation rule, and the second paragraph states the
advertisement-and-enforcement rule. A tighter wrap recovers the new blank line, so the
file stays at 120 lines.

The registry test file states each doc-comment sentence in 25 words or fewer. The
finding-discipline doc comment names no rule count.

`delegation-discipline.md` says that a live-tree charge includes the inventory file in
its fence. The fixture span keeps its bytes.

`finding-discipline.md` names what `craft-review` keeps: the three-axis split, the smell
baseline, and the universal-claim rule. The reference alone holds the citation standard
and the refute step. The anchored lead sentence keeps its bytes.

`TestReferenceFileAnchorsRedOnAbsence` grades a third tree per reference. The tree
carries the reference file with a title and no lead sentence. The dropped-lead diagnostic
fires, and the file-missing diagnostic stays silent.

Story 39 and row KG39 name the citation sentence, because the seam grades the bytes and
not the line number. One Won't handle line records that no check grades the existence of
the two fenced file paths.

## Acceptance

- [ ] The review phase file gives each falsification finding one explicit outcome of accept, merge, or dismiss, and the fixture `review-falsification-dispositions` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] The recipes hold a Claude exec line and a Codex exec line with an empty quoted heredoc, and the fixture `cross-harness-reviewer-exec-heredoc` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] `bench test --check guidance-prose-budgets` and `bench test --check prose-mechanics` stay green with `craft-gate` SKILL.md at 120 lines and its `Run the real path` section in two paragraphs.
- [ ] Every doc-comment sentence in `internal/anchors/registry_data_test.go` holds 25 words or fewer, and no doc comment states a rule count.
- [ ] `delegation-discipline.md` says a live-tree charge includes the inventory file in its fence, and the fixture `delegate-live-tree-inventory-fence` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] `finding-discipline.md` keeps its anchored lead sentence and names the three-axis split, the smell baseline, and the universal-claim rule.
- [ ] `TestReferenceFileAnchorsRedOnAbsence` reds on a tree that carries a reference file with no lead sentence, and it stays silent on the live root.
- [ ] Story 39 and row KG39 name the citation sentence, and the spec holds one Won't handle line for the two fenced file paths.
- [ ] `bench preflight build kit-guidance-fold` accepts this ticket, and `reviews/kit-guidance-fold.md` is absent.
