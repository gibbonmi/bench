# Kit-guidance fold

Status: staged

Roadmap: FT158, FT236, FT259, FT269, FT273

Decision source: named reviewed artifacts — the five drained roadmap detail files `roadmap/FT158.md`, `roadmap/FT236.md`, `roadmap/FT259.md`, `roadmap/FT269.md`, and `roadmap/FT273.md`, which the reviewer named together on 2026-09-01

Verification log: 2 iteration(s) to accept — one round returned revise; the author folded two blocking findings and seven acceptance findings

## Problem

Five drained roadmap rows each name one guidance gap, and each gap cost a repair cycle
in August 2026. The gaps sit in five kit files, and every one is a `kit-edit` row.

A gate author wrote a check on a workflow output and anchored only the consumer. One
swapped producer command restarted the macOS runners while the gate stayed green.

A coordinator accepted repairs that amended the coverage map, and no ticket owned the
amended rows. The final check then failed on `rows-owned`.

A kit-guidance diff went to a second harness only when the author offered the pass.
The cross-family reviewer had no named `bench worktree exec` form for an empty heredoc.

A write charge named the anchors unit test as its probe, and that test drives a
synthetic tree. An omission in the live skill file stayed green. Twelve occurrences
later, charges still named the wrong probe, the wrong root pass, or an incomplete
fence.

A review axis filed a finding its own citation refuted. Another cited a line it had not
read, and a third claimed an environment variable absent without a producer read. The
coordinator paid a read to dismiss each one.

The gate sees a dropped anchor, a broken fixture, and a budget overrun. The gate does not
see a semantic defect in guidance prose, so a second harness's read is the only check.

## Solution

The kit closes each gap at the place the reader already opens. Every new normative
sentence ships as a rule, an anchor tuple, a red-on-removal registry test, and a
live-mirror canary fixture. That is the shape of the spec-authoring-discipline landing.

`craft-gate` states that a check on an indirected value grades both ends, and that the
author asks which single edit defeats a new check.

`/bench-review-implementation` requires one repair ticket before the repair-scoped
re-review when accepted repairs amend the map. The same phase file makes the
cross-harness falsification pass standing for a kit-guidance diff, with one explicit
disposition per finding. The cross-harness recipes name the `bench worktree exec` form
with an empty quoted heredoc.

`references/delegation-discipline.md` gains the probe-oracle rules and the fence rules
from FT273. `references/map-discipline.md` extends the literal-bytes search to a moved
sentence.

`craft-review` gains `references/finding-discipline.md`, which holds the six per-axis
finding rules from FT236. The skill file points at the reference, and its budget row
rises to 122.

## User stories

### Group A — the indirected-value rule in `craft-gate` (FT269)

Line: opus / medium. The group adds two anchored sentences to a file at its budget, and
medium is the standing implementer effort for guidance prose.

1. As a gate author, I want the both-ends rule for an indirected value in `craft-gate`,
   so that a name proves nothing alone.
2. As a gate author, I want `craft-gate` to ask which single edit defeats a new check, so
   that I find the hole first.
3. As a kit maintainer, I want `craft-gate` SKILL.md at 120 lines or fewer, so that its
   glob budget row holds.

### Group B — the repair ticket owner (FT259)

Line: opus / medium. The group adds one anchored paragraph to an unbudgeted phase file
beside a pinned contract, so it needs the tier that reads the pinning test first.

4. As a coordinator, I want one repair ticket before the repair-scoped re-review when
   repairs amend the map, so that `rows-owned` stays green.
5. As a coordinator, I want that ticket to cite each amended row in `Covers:`, so that
   preflight reads the owner from the parsed field.
6. As a kit maintainer, I want the rule in the review phase file alone, so that the build
   phase file and `craft-tickets` take no line.
7. As a kit maintainer, I want the six pinned substrings of the convergence contract
   untouched, so that the live-tree test stays green.

### Group C — the standing falsification pass (FT158)

Line: opus / medium. The group edits two unbudgeted guidance files, one file at an exact
budget, and one guard test table, and medium carries a small Go table edit.

8. As a reviewer, I want a kit-guidance diff to take a standing cross-harness falsification
   pass, so that a second harness reads every guidance change.
9. As a reviewer, I want the kit-guidance set named as `.agents/` plus `.bench/BENCH.md`,
   so that a coordinator decides by path.
10. As a coordinator, I want each falsification finding to take accept, merge, or dismiss,
    so that no finding drops in silence.
11. As a coordinator, I want an accepted falsification finding to take the existing
    repair-routing label, so that the kit keeps one repair vocabulary.
12. As a coordinator, I want the recipes to name the exec form with an empty quoted
    heredoc, so that the reviewer stays inside the boundary.
13. As a kit maintainer, I want the two bare recipe lines byte-identical, so that their
    anchors and fixture stay green.
14. As a coordinator, I want the guard to allow an empty quoted heredoc inside an exec span,
    so that the named form runs.
15. As a coordinator, I want the guard to refuse a `/dev/null` redirection inside an exec
    span, so that the recipe's reason is tested.
16. As a builder, I want the `--full` paragraph to stop calling the falsification pass an
    offer, so that the two phase files agree.
17. As a kit maintainer, I want the build phase file at 75 lines, so that its exact budget
    row holds.

### Group D — the probe oracle and the fence rules (FT273)

Line: opus / medium. The group adds nine anchored sentences to two unbudgeted reference
files and rewrites two bullets in place, so it needs the tier that reads the anchor
registry before it edits.

18. As a coordinator, I want an anchor charge to name `bench test --check <owning-check>`
    as its probe, so that the probe grades the live file.
19. As a coordinator, I want a skippable-test charge to name the capability-skip check in
    its focused checks, so that no skip reaches the landing gate.
20. As a coordinator, I want the reference to say the conformance package run is not the
    root pass, so that no charge names it.
21. As a coordinator, I want to confirm the mutated bytes before I read a probe verdict, so
    that a no-op edit never passes.
22. As a coordinator, I want a `PATH` or environment bind to fence the serial-census ceiling
    file, so that the ceiling never sits outside a fence.
23. As a coordinator, I want a live-tree test charge to fence the live-tree inventory, so
    that no mid-build fence amendment repeats.
24. As a coordinator, I want a grammar charge to enumerate shared fixture owners and
    exact-record assertion families, so that no shared fixture goes unfenced.
25. As a write delegate, I want the spec's fence in my charge and a report before any
    out-of-fence edit, so that no write escapes.
26. As a kit maintainer, I want the repair-charge template and its five anchors untouched,
    so that FT164 stays its owner.
27. As a spec author, I want the literal-bytes search to cover a moved sentence over
    `tests/` and `internal/conformance`, so that no pinned sentence moves unseen.
28. As a kit maintainer, I want `craft-delegate` and `craft-spec` SKILL.md to take no net
    line, so that both stay inside their budget rows.
29. As a kit maintainer, I want each FT273 rule once across the two reference files, so
    that the one-source rule holds.

### Group E — the finding discipline reference in `craft-review` (FT236)

Line: opus / medium. The group creates one reference file, adds one pointer line to a
skill at its budget, and raises one budget row, so it needs the tier that reads the
budget parser.

30. As a review axis, I want the reference to call a generated script's string expectation
    the mutation catch, so that I never discount it.
31. As a review axis, I want to cite the line I read this pass or the symbol, so that every
    citation is checkable.
32. As a review axis, I want a test-deleting Standards finding to name the surviving
    assertion or file as coverage, so that coverage never drops.
33. As a review axis, I want a real run before I report a strong finding, so that the repo
    refutes me first.
34. As a Coverage axis, I want an environment-variable finding to cite the producer before
    it claims absence, so that a wrapper-set variable reads as present.
35. As a reviewer, I want an unreachable row seam disposed as a seam-column amendment, so
    that a helper-seam test is never a partial row.
36. As a kit reader, I want `references/finding-discipline.md` present with its charge-time
    lead sentence, so that the reference reads as a charged file.
37. As a kit reader, I want `craft-review` SKILL.md to point at the reference, so that an
    axis finds the rules from the skill.
38. As a kit maintainer, I want an exact budget row of 122 for `craft-review` SKILL.md, so
    that the pointer line fits.
39. As a kit maintainer, I want line 84 of the review phase file untouched, so that its
    anchor and fixture stay green.
40. As a kit maintainer, I want the refute paragraph to hold one refute rule and no real-run
    copy, so that each source holds one rule.

## Implementation decisions

**Each new normative sentence ships as a rule, an anchor, a registry test, and a live-mirror
fixture.** The sentence lands in its guidance file. One tuple in the anchor registry
requires or forbids the exact needle. One registry test function per ticket proves each
tuple reds on removal against a synthetic tree and stays silent on the live root. One
fixture under the workflow-guidance-anchors family carries `BASE`, `MUTATE.json`, and
`EXPECT`, in the live-mirror shape the spec-authoring-discipline landing used. The
build does not use the older `files/` snapshot shape for a new fixture.

The fixture is the row's seam, because the fixture-bite test already exists. The registry
test function is the build's own red proof, and each ticket names its function.

**Every fixture needle sits on one physical line.** The anchor evaluator matches under
collapsed whitespace, so a wrapped needle passes the registry test. The fixture
materializer refuses an `old` value that spans a line wrap, so the fixture-bite test reds
it. The build confirms the wrap after each edit.

**New tuples join the newest anchor group.** Each new tuple takes the group of the most
recent tuples in the registry. A tuple that scopes to a section uses the H2 heading it
lands under.

**`craft-gate` lands at exactly 120 lines.** Rule (a) is a standalone sentence at the end
of the first paragraph under `Run the real path`. Rule (b) is a standalone sentence at
the end of the second paragraph under `Prove it bites`. Each addition costs one net line
after a reflow of the ragged break in the `Run the real path` paragraph. No reader
pins those bytes. No budget row moves for this file.

**The repair ticket rule has one home.** The sentence pair lands as one paragraph under
`Review modes` in the review phase file, between the re-review paragraph and the landing
paragraph. The build phase file sits at its exact budget, `craft-tickets` sits one line
under its budget, and both stay untouched. The ticket template already carries `Covers:`
and `Blocked by:`, so a repair ticket is an ordinary ticket to the parser.

**The convergence contract bounds the edit.** A live-tree test pins six normalized
substrings of the review phase file and two of the build phase file. The test collapses
whitespace before it matches, so an insert or a reflow survives and a rewritten sentence
reds. The build inserts paragraphs and rewrites no existing sentence in `Review modes`.

**The falsification pass keeps the row's disposition labels.** The review phase file
already gives every review finding one repair-routing label. A falsification finding
first takes accept, merge, or dismiss. An accepted one then joins the review findings
and takes the existing label. One sentence states that bridge, so the two vocabularies
compose instead of competing.

**A kit-guidance diff is a path predicate.** The standing rule names the set as any file
under `.agents/` or the file `.bench/BENCH.md`. That is the widest of the three sets the
tree already uses, and a coordinator can decide it by path.

**The standing-pass paragraph lands under `Review modes`.** It sits after the repair
ticket paragraph and before the landing paragraph. Its four tuples scope to that section,
as the two repair ticket tuples do.

**The cross-harness recipes gain a form and lose nothing.** The two bare recipe lines are
anchored byte for byte and pinned by a fixture. The build adds one paragraph after them.
The paragraph names the exec form with an empty quoted heredoc and states the reason:
the guard refuses any non-heredoc redirection inside an exec span. Two new cases in the
ordinary guard test table make that reason a tested fact.

**The build phase file takes a same-count rewrite of its `--full` paragraph.** The
anchored ask-before-adding sentence takes the leading clause `Outside the kit-guidance
set,` and keeps its needle, so its tuple stays green. The last sentence now says that
tier escalation is an offer and that the review phase owns the standing pass. A forbid
tuple over the whole file keeps the retired sentence out.

**FT273 folds into two reference files and adds no line to either skill file.** Rules 1, 2,
5, 6, and 7 land as bullets under `In the charge`, and rule 5 names the ceiling file
`internal/worktree/parallel_census_test.go`. Rule 8 is an in-place generalization
of the out-of-fence bullet from "registry" to "write". Rules 3 and 4 land under
`Probes`, and rule 3 extends the live-tree probe bullet. Rule 10 extends the
literal-bytes bullet in `map-discipline.md` from "deletes" to "deletes or moves" and
names `tests/` and `internal/conformance`. Rule 9 adds no sentence.

**The kit keeps the term "exact-record assertion families".** The tree holds one instance
today, and the charge-side rule needs a category name that does not rot.

**FT236 lands in a new reference file.** `craft-review` SKILL.md sits at its budget, and
its one existing reference holds the Standards smell baseline alone. The six rules land
in `references/finding-discipline.md` with a charge-time lead sentence. One pointer
sentence lands under `What a finding must cite`. The exact budget row for the skill file
rises to 122, which matches the two existing lifted rows.

**The reference file has four H2 sections, and each rule tuple scopes to its section.**
The sections are `What a string expectation proves` for rule 1, `What a citation points
at` for rule 2, `Where an axis under-reads` for rules 3, 4, and 5, and `When a seam
cannot reach the state` for rule 6. The lead sentence sits above the first section, and
its tuple scopes to no section.

**The recall half of the citation rule stays where it is.** The review phase file already
states that a finding cites what its axis read now. That line is anchored and fixture
pinned. The location clause lands in the reference alone, and lines 82 to 84 of the
phase file already route the reader there.

**The tickets land serially on one integration source.** All five tickets write the anchor
registry and its test file. Each ticket is blocked by the previous one, so the frontier
holds one ticket at a time.

## Testing decisions

- A good test drives one published behavior. The published behaviors are four: the
  registry verdict over a guidance file, the fixture bite, the guard verdict, and the
  budget verdict.
- The prose seams are the registry red-on-removal tests in `internal/anchors` and the
  fixture-bite test in `internal/conformance`. Their prior art is the
  spec-authoring-discipline landing.
- The guard seam is the ordinary allow-and-refuse table in `internal/benchguard`. Its
  prior art is the existing exec heredoc case in that table.
- The gate observes the feature through the `docs-currency-workflow` check, the
  `workflow-guidance-anchors` canary family, and the `guidance-prose-budgets` check.

### Seam diagram

    trigger: bench gate, test phase
        │
        ▼
    guidance file text ──▶ [ anchors registry ] ──▶ diagnostic or silence
                      ◀ tests attach here: a registry test mutates one needle in a
                        synthetic tree, and a fixture pins the exact diagnostic

    trigger: Claude or Codex Bash tool call
        │
        ▼
    command text ──▶ [ benchguard exec span ] ──▶ allow, or a named refusal
                      ◀ tests attach here: the ordinary table drives one command
                        string and asserts the verdict

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| KG1 | 1 | `craft-gate` `Run the real path` states that a check on a workflow output, a config key, or an environment variable grades the producer and the consumer, or their binding, in the same change | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/craft-gate-indirected-value-both-ends` | Deleting the binding clause removes the needle bytes, and the anchor reds |
| KG2 | 2 | `craft-gate` `Prove it bites` states that the author asks which single edit defeats a new check while the gate stays green | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/craft-gate-single-edit-defeat` | Swapping `single edit` for `change` drops the one-edit bound, and the anchor reds |
| KG3 | 3 | `.agents/skills/bench-craft-gate/SKILL.md` holds 120 lines or fewer | `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetsHoldOnTheLiveTree`) | The fold lands at exactly 120, so one extra line reds the check |
| KG4 | 4 | The review phase file requires one repair ticket before the repair-scoped re-review when accepted repairs amend the map | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/review-repair-ticket-owner` | Swapping `before` for `after` restores the BG37 sequence, and the anchor reds |
| KG5 | 5 | The review phase file requires the repair ticket to cite each amended row in `Covers:` and to record the accepted repairs | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/review-repair-ticket-covers` | Dropping the `Covers:` clause lets a prose citation pass, which `rows-owned` reds |
| KG6 | 6 | The repair ticket rule appears in the review phase file alone | review-owned: the Standards axis reads the diff against the fences | Only a reader can see a second copy in the build phase file or the delegation reference |
| KG7 | 7 | The six pinned substrings of the review convergence contract still match the live tree | `internal/conformance/docs_workflow_helpers_test.go` (`TestReviewConvergenceContractCurrentDocs`) | A rewritten sentence in the re-review paragraph breaks one pinned substring, because the test collapses whitespace before it matches |
| KG8 | 8 | The review phase file states that a diff that changes kit guidance takes a standing cross-harness falsification pass | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/review-standing-falsification` | Swapping `standing` for `offered` removes the needle bytes, and the anchor reds |
| KG9 | 9 | The standing rule names the kit-guidance set as any file under `.agents/` or the file `.bench/BENCH.md` | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/review-kit-guidance-set` | A set narrowed to one directory drops the needle, and the anchor reds |
| KG10 | 10 | The review phase file gives each falsification finding one explicit disposition of accept, merge, or dismiss | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/review-falsification-dispositions` | Omitting the three labels leaves a disposition with no values, and the anchor reds |
| KG11 | 11 | The review phase file states that an accepted falsification finding joins the review findings and takes the repair-routing disposition | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/review-falsification-accept-routing` | Dropping the bridge leaves two competing vocabularies, and the anchor reds |
| KG12 | 12 | The cross-harness recipes name `bench worktree exec <target> --` with an empty quoted heredoc as the reviewer form | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/cross-harness-reviewer-exec-heredoc` | Reverting the form to a `/dev/null` redirection drops the needle, and the anchor reds |
| KG13 | 13 | The two bare Claude and Codex recipe lines stay byte-identical | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/cross-harness-reviewer-recipes` | A rewritten recipe line keeps the diagnostic alive on the restored tree, so the existing fixture's restoration run reds |
| KG14 | 14 | The guard allows `bench worktree exec L -- cat <<'EOF'` followed by a bare `EOF` line | `internal/benchguard/benchguard_test.go` (`TestClassifySpanScopedFollowOns`) | A guard that keys the heredoc exception on a non-empty body refuses the named form |
| KG15 | 15 | The guard refuses `bench worktree exec L -- cat < /dev/null` and names the redirection | `internal/benchguard/benchguard_test.go` (`TestClassifySpanScopedFollowOns`) | A guard that lets a redirection through the exec span makes the recipe's reason false |
| KG16 | 16 | The build phase file no longer holds the sentence that calls both passes offers | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/implement-spec-offer-retired` | A forbid tuple reds the moment the retired sentence returns |
| KG17 | 17 | `.agents/commands/bench-implement-spec.md` holds 75 lines or fewer | `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetsHoldOnTheLiveTree`) | A rewrite that adds a line reds the exact budget row |
| KG18 | 18 | `delegation-discipline.md` `In the charge` requires an anchor-adding charge to name `bench test --check <owning-check>` as its probe | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/delegate-anchor-probe-owning-check` | Swapping the owning check for the anchors package restores the FT214 miss, and the anchor reds |
| KG19 | 19 | `delegation-discipline.md` `In the charge` requires the capability-skip check in a skippable test's focused checks | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/delegate-skip-ownership-check` | Dropping the focused-checks clause lets a skip reach the landing gate, and the anchor reds |
| KG20 | 20 | `delegation-discipline.md` `Probes` states that `bench test --package ./internal/conformance` is not the root conformance pass | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/delegate-root-conformance-pass` | Swapping `is not` for `is` restores the FT213 miss, and the anchor reds |
| KG21 | 21 | `delegation-discipline.md` `Probes` requires the coordinator to confirm the mutated bytes against the copy aside before it reads the verdict | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/delegate-probe-mutated-bytes` | Dropping the before-the-verdict clause lets a no-op edit read as a pass, and the anchor reds |
| KG22 | 22 | `delegation-discipline.md` `In the charge` fences the serial-census ceiling file `internal/worktree/parallel_census_test.go` for a `PATH` or process-environment bind | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/delegate-serial-ceiling-fence` | Weakening `includes` to `may include` removes the needle bytes, and the anchor reds |
| KG23 | 23 | `delegation-discipline.md` `In the charge` fences `internal/conformance/tier_test.go` for a live-tree test | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/delegate-live-tree-inventory-fence` | Dropping the path restores the PL37 mid-build amendment, and the anchor reds |
| KG24 | 24 | `delegation-discipline.md` `In the charge` requires a grammar charge to enumerate the shared fixture owners and the exact-record assertion families | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/delegate-grammar-fence-inventory` | Keeping only the fixture conjunct leaves an assertion file unfenced, and the anchor reds |
| KG25 | 25 | `delegation-discipline.md` `In the charge` requires the delegate to report an out-of-fence write before it edits | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/delegate-out-of-fence-write` | Reverting `write` to `registry` narrows the rule to the pre-fold reading, and the anchor reds |
| KG26 | 26 | The `Repair-charge template` section keeps its five field needles | `internal/anchors/registry_data_test.go` (`TestRepairChargeTemplateAnchorsRedOnRemoval`) | A re-authored template drops one of the five needles the existing test pins |
| KG27 | 27 | `map-discipline.md` `Before the map locks` requires the literal-bytes search on a deleted or moved sentence over `tests/` and `internal/conformance` | `internal/anchors/registry_data_test.go` (`TestMapDisciplineTwoAudienceAndTransactionAnchorsRedOnRemoval`) with `tests/canary/workflow-guidance-anchors/map-discipline-moved-bytes-sweep` | Dropping `or moves` restores the 2026-09-01 miss, and the anchor reds |
| KG28 | 28 | `.agents/skills/bench-craft-delegate/SKILL.md` holds 122 lines or fewer | `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetsHoldOnTheLiveTree`) | The file sits at its bound, so one stray line reds the check |
| KG29 | 29 | Each FT273 rule appears once across the two reference files | review-owned: the Standards axis reads the diff | Only a reader can see a second copy of a rule in otherwise correct prose |
| KG30 | 30 | `finding-discipline.md` states that a generated script's independently authored string expectation is the mutation catch | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/review-string-expectation-catch` | Weakening `is the mutation catch` to a weaker signal removes the needle bytes, and the anchor reds |
| KG31 | 31 | `finding-discipline.md` requires a finding to cite the line the axis read this pass, or the symbol instead | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/review-citation-location` | Dropping the symbol arm leaves a location rule with no escape, and the anchor reds |
| KG32 | 32 | `finding-discipline.md` makes a test-deleting Standards finding name the surviving assertion or file as coverage | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/review-test-deletion-coverage` | Softening `names` to `may name` removes the needle bytes, and the anchor reds |
| KG33 | 33 | `finding-discipline.md` requires a real run before an axis reports a strong finding | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/review-strong-finding-run` | Dropping `with a real run` leaves the generic refute rule alone, and the anchor reds |
| KG34 | 34 | `finding-discipline.md` makes an environment-variable Coverage finding cite the producer before it claims absence | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/review-env-var-producer` | Dropping the before-the-claim clause lets an absence claim ship unread, and the anchor reds |
| KG35 | 35 | `finding-discipline.md` disposes an unreachable seam as an amendment of the row's seam column | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/review-seam-amendment` | Swapping `amends` for `notes` turns the amendment into a comment, and the anchor reds |
| KG36 | 36 | `finding-discipline.md` exists and carries its charge-time lead sentence | `internal/anchors/registry_data_test.go` (`TestReferenceFileAnchorsRedOnAbsence`) | A missing or empty file raises the file-missing or needle-missing diagnostic the test asserts |
| KG37 | 37 | `craft-review` SKILL.md points at `references/finding-discipline.md` under `What a finding must cite` | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/review-finding-discipline-pointer` | Removing the pointer leaves the reference unread, and the anchor reds |
| KG38 | 38 | `.agents/skills/bench-craft-review/SKILL.md` holds 122 lines or fewer under its new exact row | `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetsHoldOnTheLiveTree`) | The check parses the profile table, so a missing row leaves the glob bound of 120 and the pointer line reds it |
| KG39 | 39 | The citation sentence at line 84 of the review phase file stays byte-identical | `internal/conformance/gate_entry_test.go` (`TestRootConformance`) | A rewrite of line 84 drops its anchored needle, and the live-root evaluation reds the `docs-currency-workflow` check |
| KG41 | 16 | The `--full` paragraph scopes its ask-before-adding sentence to a diff outside the kit-guidance set | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) over `tests/canary/workflow-guidance-anchors/implement-spec-offer-scope` | Dropping the scope clause leaves the file asking before a standing pass, and the anchor reds |
| KG42 | 28 | `.agents/skills/bench-craft-spec/SKILL.md` holds 152 lines or fewer | `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetsHoldOnTheLiveTree`) | The file sits one line under its bound, so two stray lines red the check |
| KG40 | 40 | The `Refute before you report` section holds one refute rule and no copy of the real-run clause | review-owned: the Standards axis reads the skill and the reference together | A sectioned anchor cannot forbid a paraphrase, so only a reader can see the second copy |

### Edge inventory

- **The needle-wrap edge.** The anchor evaluator collapses whitespace, so a wrapped
  needle passes the registry test. The fixture materializer refuses an `old` value that
  spans a line wrap, so the fixture-bite test reds it. A probe's `sed` on a wrapped needle
  changes nothing. KG1 is the reference example.
- **The budget-ceiling edges.** `craft-gate` lands at exactly 120 after one reflow. The
  build phase file takes a same-count rewrite. `craft-delegate` and `craft-spec` take no
  line. `craft-review` rises to 122. Rows KG3, KG17, KG28, and KG38 hold these edges.
- **The anchored-neighbour edges.** The two recipe lines and line 84 of the review phase
  file are anchored, and the build rewrites neither. The anchored `--full` sentence takes
  a leading scope clause and keeps its needle. Rows KG13, KG16, KG39, and KG41 hold these
  edges.
- **The pinned-contract edge.** The convergence contract test pins six normalized
  substrings of the review phase file. An insert or a reflow survives, and a rewritten
  sentence reds. Row KG7 holds this edge.
- **The absent-versus-empty edge.** A missing reference file and an empty one both red the
  require tuple on its lead sentence, with two different diagnostics. Row KG36 holds this
  edge. Every new file ends with a newline.
- **The fixture-shape edge.** Two fixture shapes exist in the family. Every new fixture uses
  the live-mirror shape, so a needle edit in a guidance file needs no snapshot update.
- **The forbid-needle edge.** Row KG16 uses a forbid tuple, because absence is the property.
  Its fixture restores the retired sentence to make the diagnostic bite.
- **The exec-span either-side edge.** The guard refuses a redirection only inside an exec
  span, and the bare recipes outside a span stay valid. Rows KG14 and KG15 hold one side
  each.
- **The vocabulary edge.** The falsification labels and the repair-routing labels are two
  sets with one bridge sentence. Rows KG10 and KG11 hold the two sets.
- **The hand-edited-markdown edges.** A non-ASCII space inside a needle unanchors it, and
  the registry test writes the needle independently in Go, so the live-root run reds it.
  The new reference file carries no fence or comment block, so no unterminated delimiter
  can swallow it.

**Won't handle** — a machine check that a falsification pass ran on a guidance diff — no
artifact records the pass. The review phase's evidence step still asks for it.

**Won't handle** — a Go definition of the kit-guidance set — the rule is a path predicate
in prose. The review phase still reads the diff paths.

**Won't handle** — a ticket-parser change for a repair ticket — `Covers:` and `Blocked by:`
already carry it. The `rows-owned` check still grades the citation.

**Won't handle** — a registry for the exact-record assertion families — the tree holds one
instance. The charge author still enumerates them by hand.

**Won't handle** — a `craft-review` pointer at the single-edit rule in `craft-gate` — the
review skill already asks what authenticates the verifier. The gate author still reads
`craft-gate`.

**Won't handle** — a system-test case for the empty heredoc through the installed hook —
the ordinary guard table is the one predicate owner. The existing system journey still
covers a heredoc with a body.

## Ownership fences

The ticket binding registry binds `internal/anchors` to the five command-registry files,
so the `registry-closure` check forces every ticket to name them. The fence lists them
for that reason, and the build writes no byte in them.

- `.agents/skills/bench-craft-gate/SKILL.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `.agents/skills/bench-craft-spec/references/map-discipline.md`
- `.agents/skills/bench-craft-review/SKILL.md`
- `.agents/skills/bench-craft-review/references/finding-discipline.md`
- `projects/benchkit.md`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/benchguard/benchguard_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/main_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_test.go`
- `tests/canary/workflow-guidance-anchors`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `specs/kit-guidance-fold`
- `reviews/kit-guidance-fold.md`

## Out of scope

- A `CONTEXT.md` entry for "kit-guidance diff" with an Avoid list. Estimate: 2 edits, 1
  gate run.
- A `bench` verb that records a falsification pass as an artifact the final check reads.
  Estimate: 8 edits, 3 gate runs.
- A ticket-parser field that marks a ticket as a repair ticket. Estimate: 5 edits, 2 gate
  runs.
- A system-test case that drives the empty heredoc through the installed Claude hook.
  Estimate: 3 edits, 2 gate runs.
- A budget-row raise for `craft-gate`, `craft-delegate`, or `craft-spec`. Estimate: 1
  edit, 1 gate run.

## Further notes

### Reviewer decisions taken under the recommended answer

The reviewer named the five rows and asked for one spec. Each call below goes past the row
bodies. The reviewer can veto any one at sign-off.

1. **The `craft-review` budget row rises to 122.** The skill file sits at 120 of 120, and
   the pointer line needs one. A reflow would cost substantive prose. Story 38 and row
   KG38 hold it.
2. **The falsification pass keeps accept, merge, and dismiss.** The row body states the
   labels. One bridge sentence routes an accepted finding into the existing repair-routing
   labels. Stories 10 and 11 hold it.
3. **The kit-guidance set is `.agents/` plus `.bench/BENCH.md`.** No tree definition
   exists, and this is the widest of the three sets the tree uses. Story 9 holds it.
4. **The build phase file's `--full` paragraph takes a same-count rewrite.** The anchored
   sentence gains a scope clause and keeps its needle. Story 16 holds it.
5. **The repair ticket rule has one home, and no carve-out from the approved breakdown.**
   The approval rule scopes to assignment, and a retroactive ticket is never assigned.
   Story 6 holds it.
6. **FT273 rule 8 is an in-place generalization, and rule 10 extends the existing bullet.**
   A second bullet on the same fact would be a second source. Stories 25 and 27 hold it.
7. **The term "exact-record assertion families" stays.** Story 24 holds it.
8. **`craft-gate` ships at exactly 120 with no budget raise.** Story 3 holds it.
9. **The FT236 rules land in a new file, `references/finding-discipline.md`.** The
   existing reference holds the smell baseline alone. Story 36 holds it.

### Reads

The author read the spec discipline, the map discipline, the ticket template, the synthesis
skill, the delegate skill, the kit profile, and the precedent spec. The author read the
anchor kinds and groups, the first sixty registry tuples, and one live-mirror fixture.

The author read the five landing sites. These are the `Review modes` and disposition
paragraphs of the review phase file, and the full cross-harness recipes. They also
include the `--full` paragraph of the build phase file and the `In the charge` and
`Probes` bullets of the delegation reference. The last site is the citation and refute
sections of the review skill. The author read the `rowsOwnedCheck` function, the guard's non-heredoc refusal line, the
anchored `--full` tuple, and the `Roadmap:` grammar.

Five read-only research delegates, one per row, read the rest at opus and medium effort.
Their reads cover:

- the registry test functions named in the map and the fixture-bite test
- the budget parser and the convergence contract test
- the guard's allow and refuse tables
- the `TestRootConformance` skip and the serial-census ceiling constant
- the live-tree inventory and the repaired producer anchor from the FT269 occurrence

The author spot-checked four of their citations against the tree and found each fact
held. One delegate line citation for the guard was off by eighty lines, and the author
replaced it with a read.

The author did not read `internal/preflight/gather.go`, `internal/tickets`,
`internal/conformance/tier_test.go`, or `internal/worktree/parallel_census_test.go`.
Each claim about them rests on a delegate read.

### Reader sweep

- The `Roadmap:` grammar accepts one `FT<n>` token for the retire line's remainder text.
  A five-row value falls back to the generic placeholder and refuses nothing. The
  precedent spec used a three-row value.
- The two cross-harness recipe lines have three readers: two anchor tuples and the fixture
  `cross-harness-reviewer-recipes`. The build adds after them.
- The `--full` paragraph's anchored sentence has one reader, its tuple. The sentence after
  it has no reader. The build rewrites only the second.
- Line 84 of the review phase file has two readers: its tuple and the fixture
  `review-universal-claim-bar`. The `Review modes` section has one reader, the
  convergence contract test with six substrings. The build inserts and reflows nothing.
- The `craft-gate` paragraphs at the two landing sites have no reader. The registry
  names `craft-gate` only inside two command files.
- `delegation-discipline.md` and `map-discipline.md` have readers only in the anchor
  registry and its test. No fixture targets `delegation-discipline.md` today. The bytes
  the build rewrites in place have no reader.
- The reader-sweep sentence in `craft-spec` SKILL.md has three readers, which is why rule
  10 lands in the reference file instead.
- The `craft-review` citation and refute sections have no anchor reader. Two prose
  documents quote one sentence as commentary, and the build keeps that sentence.
- The budget table in the kit profile has one parser, the budget check. The new exact row
  joins four existing exact skill rows.
- The `workflow-guidance-anchors` family has no count pin. The fixture-bite test names six
  fixtures by name, and the build touches none of them.

### Executed-root trace

Every canary row reaches the executed root through the fixture-bite test in
`internal/conformance`, which enumerates the canary tree from the producer. The registry
binds the family to `docs-currency-workflow`, so a new fixture joins with no registry edit.
The gate runs `internal/anchors`, `internal/benchguard`, and `internal/conformance` as
ordinary test packages, so each cited function sits inside an executed phase. Each
diagnostic state is mutable: `MUTATE.json` rewrites one needle in a real guidance file, and
the forbid fixture restores one retired sentence.

### Flagged additions

Each item below goes past the five row bodies. Review demands a row for each item, or
removes the item.

1. The bridge sentence from an accepted falsification finding to the repair-routing
   labels. Row KG11.
2. The kit-guidance set definition. Row KG9.
3. Two guard table cases for the exec span. Rows KG14 and KG15.
4. The same-count rewrite of the `--full` paragraph, with a forbid tuple and a scope
   clause on the anchored sentence. Rows KG16, KG17, and KG41.
5. The `craft-review` budget row. Row KG38.
6. The new reference file and its pointer. Rows KG36 and KG37.
7. Rows KG6, KG29, and KG40 are review-owned. No check can see a duplicated rule.
8. The fences cover the two fixtures outside the anchor family that pin the kit profile.
   They also cover the five command-registry files the ticket binding forces into every
   `Writes:` list. A fence entry keeps the `Writes:` list and the fence in agreement. Row
   KG38.

### Source sentences and their rows

| source sentence | rows |
|---|---|
| FT269 — a check on an indirected value grades the producer and the consumer, or the binding | KG1 |
| FT269 — the author asks which single edit defeats the check while it stays green | KG2 |
| FT269 — one kit edit over a single owner surface | KG3 |
| FT259 — the coordinator writes a retroactive repair ticket after repairs land and before re-review | KG4, KG6, KG7 |
| FT259 — the ticket owns the amended rows and records the accepted repairs | KG5 |
| FT158 — the pass becomes standing for a diff that changes kit guidance | KG8, KG9, KG16, KG17, KG41 |
| FT158 — each finding takes an explicit accept, merge, or dismiss disposition | KG10, KG11 |
| FT158 — the reviewer runs the pass through `bench worktree exec`, with the empty heredoc named | KG12, KG13, KG14, KG15 |
| FT273 — an anchor charge names `bench test --check <owning-check>` | KG18 |
| FT273 — a skippable-test charge names the capability-skip check | KG19 |
| FT273 — the conformance package run is not the root pass | KG20 |
| FT273 — a probe reaches the seam and confirms the mutated bytes | KG21 |
| FT273 — a `PATH` or environment bind fences the ceiling file | KG22 |
| FT273 — a live-tree test fences the inventory | KG23 |
| FT273 — a grammar change enumerates fixture owners and assertion families | KG24 |
| FT273 — every write charge receives the fence and reports an out-of-fence write | KG25 |
| FT273 — FT164 remains the repair-template owner | KG26 |
| FT273 — the reader sweep covers the literal bytes of every moved sentence | KG27 |
| FT273 — the skill core sits at its lifted budget bound | KG28, KG29, KG42 |
| FT236 — a string expectation on a generated script is the mutation catch | KG30 |
| FT236 — a finding cites the line read this pass or the symbol | KG31, KG39 |
| FT236 — a test-deleting Standards finding is a coverage question first | KG32 |
| FT236 — an axis refutes a strong finding with a real run | KG33, KG40 |
| FT236 — an environment-variable finding cites the producer | KG34 |
| FT236 — an unreachable seam is a spec amendment of the seam column | KG35 |
| FT236 — one kit edit over a single owner surface | KG36, KG37, KG38 |
