# Make the falsification pass standing for guidance diffs

Blocked by: name-the-repair-ticket-owner-before-re-review.md
Writes: .agents/commands/bench-review-implementation.md, .agents/commands/bench-implement-spec.md, .agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/benchguard/benchguard_test.go, tests/canary/workflow-guidance-anchors, tests/canary/workflow-guidance-anchors/review-standing-falsification (new), tests/canary/workflow-guidance-anchors/review-kit-guidance-set (new), tests/canary/workflow-guidance-anchors/review-falsification-dispositions (new), tests/canary/workflow-guidance-anchors/review-falsification-accept-routing (new), tests/canary/workflow-guidance-anchors/cross-harness-reviewer-exec-heredoc (new), tests/canary/workflow-guidance-anchors/implement-spec-offer-retired (new), tests/canary/workflow-guidance-anchors/implement-spec-offer-scope (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: KG8, KG9, KG10, KG11, KG12, KG13, KG14, KG15, KG16, KG17, KG41

## What to build

A second harness reads every kit-guidance diff, because the cross-harness falsification
pass is standing. The two phase files agree, the recipes name a reviewer form that the
guard allows, and the guard tests that form.

Four new sentences land in `.agents/commands/bench-review-implementation.md`. The first
sentence states that a diff that changes kit guidance takes a standing cross-harness
falsification pass. The second sentence names the kit-guidance set as any file under
`.agents/` or the file `.bench/BENCH.md`. The third sentence gives each falsification
finding one explicit disposition of accept, merge, or dismiss. The fourth sentence bridges
the two vocabularies: an accepted falsification finding joins the review findings and
takes the existing repair-routing disposition.

The mutations fix four needle words. The first mutation swaps `standing` for `offered`.
The second mutation narrows the set to one directory. The third mutation omits the three
labels. The fourth mutation drops the bridge. Hold each of those words inside its needle.

The paragraph lands under `Review modes`, after the repair ticket paragraph and before
the landing paragraph. Give each of its four sentences one `RequireInSection` tuple over
that section.

One paragraph lands in
`.agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md`, after the two
bare recipe lines. The paragraph names `bench worktree exec <target> --` with an empty
quoted heredoc as the reviewer form. It states the reason: the guard refuses any
non-heredoc redirection inside an exec span. The mutation reverts the form to a
`/dev/null` redirection. That file carries one H1 and no H2, so its tuple takes `Require`
with no section.

The `--full` paragraph in `.agents/commands/bench-implement-spec.md` takes a same-count
rewrite of its last two sentences. The anchored ask-before-adding sentence takes the
leading clause `Outside the kit-guidance set,` and keeps its needle, so its tuple stays
green. The mutation for row KG41 drops that scope clause. The last sentence states that
tier escalation is an offer, and that the review phase owns the standing pass. Give the
retired sentence one `Forbid` tuple over the whole file, so the fixture reds the moment
the sentence returns.

Two cases join the ordinary guard table `TestClassifySpanScopedFollowOns` in
`internal/benchguard/benchguard_test.go`. The allow case drives
`bench worktree exec L -- cat <<'EOF'` followed by a bare `EOF` line. The refuse case
drives `bench worktree exec L -- cat < /dev/null`, and it asserts the refusal that names
the redirection. The two cases make the recipe's stated reason a tested fact.

Two anchored neighbours stay byte-identical. The two bare recipe lines keep their bytes,
because two tuples and the fixture `cross-harness-reviewer-recipes` pin them. Line 84 of
the review phase file keeps its bytes. The `Review modes` section takes an inserted
paragraph and no rewritten sentence, because `TestReviewConvergenceContractCurrentDocs`
pins six normalized substrings of that file.

Give each new sentence one tuple in `internal/anchors/registry_data.go`, and give the
retired sentence one `Forbid` tuple. Every tuple takes `AfterImplementSpec`, the group the
newest tuples in that file use. Keep each fixture needle on one physical line, because the fixture materializer refuses an `old` value that spans a line wrap.

Add `TestStandingFalsificationAnchorsRedOnRemoval` to
`internal/anchors/registry_data_test.go`. The function writes each needle and each
diagnostic independently of the registry. It proves that each tuple reds when a synthetic
tree drops the needle, or when it restores the retired sentence. It also proves that the
live root stays silent.

Add seven fixtures under `tests/canary/workflow-guidance-anchors/` in the live-mirror shape.
`BASE` names the live file path the tuple reads. `MUTATE.json` holds one old-and-new
needle swap. `EXPECT` holds the diagnostic the tuple raises. The fixture
`implement-spec-offer-retired` restores the retired sentence, because absence is the
property there. Do not use the older `files/` snapshot shape for any fixture.

The build phase file sits at its exact budget row of 75 lines, so the rewrite keeps the
same line count. The review phase file and the recipes file carry no budget row. No budget
row moves for this ticket.

This ticket appends its tuples and its test function to the two registry files. It edits
nothing the previous two tickets added, and the next ticket appends beside all three sets.

## Acceptance

- [ ] KG8 — the fixture `review-standing-falsification` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG9 — the fixture `review-kit-guidance-set` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG10 — the fixture `review-falsification-dispositions` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG11 — the fixture `review-falsification-accept-routing` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG12 — the fixture `cross-harness-reviewer-exec-heredoc` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG13 — the fixture `cross-harness-reviewer-recipes` stays green, because both bare recipe lines keep their bytes.
- [ ] KG14 — `TestClassifySpanScopedFollowOns` allows `bench worktree exec L -- cat <<'EOF'` with a bare `EOF` line.
- [ ] KG15 — `TestClassifySpanScopedFollowOns` refuses `bench worktree exec L -- cat < /dev/null` and names the redirection.
- [ ] KG16 — the fixture `implement-spec-offer-retired` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
- [ ] KG17 — `TestGuidanceProseBudgetsHoldOnTheLiveTree` stays green with the build phase file at 75 lines.
- [ ] KG41 — the fixture `implement-spec-offer-scope` bites through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.
