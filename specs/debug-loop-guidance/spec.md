# Phase-specific action and evidence guidance

Status: staged

Decision source: specs/debug-loop-guidance/decisions/debug-loop-guidance.md

Verification log: Consolidated after the reviewer accepted the assessment on 2026-09-15. Mechanical authoring checks ran. No delegated review ran.

## Problem

Debug currently redirects repair authorship and excludes write delegates from its phase.
Other workflow owners leave parts of their next-action and evidence sequence implicit.
The finding rule also requires a real run for mandatory standards that have no automated check.

## Solution

Repair debug authorship at its current owner.
Add a local action, evidence, and stop sequence to spec writing, ticket slicing, implementation, and semantic review.
Each complete ticket includes its own guidance checks and future fresh-session adoption task.
The project gate remains the oracle.

## User stories

Line: gpt-5.6-sol / high.
Implementation-line reason: DG-C5 must reconcile two evidence routes without changing finding authority.
The source is precise and the anchor seam exists.
Fresh-session behavior remains review-owned, so guidance leverage selects the configured mid model at high effort.
Harder chunks: DG-C5.

### Retain debug repair authorship

1. As a debug author, I want retained repair authorship, so that the repair keeps its diagnostic context.
2. As a debug author, I want read-only diagnostic help, so that consultation preserves repair ownership.
3. As a ticket author, I want in-fence debug access, so that I can verify my own repair.
4. As a ticket author, I want out-of-fence diagnostic steps, so that the next author receives evidence.
5. As a ticket author, I want preserved dirty work at the fence, so that a blocked repair loses no work.
6. As a coordinator, I want a complete diagnostic report, so that I can validate the reslice.
7. As a reviewer, I want an observed in-fence debug repair, so that adoption rests on behavior.
8. As a reviewer, I want an observed out-of-fence handoff, so that the fence holds during adoption.

### Guide spec authoring with concrete evidence

9. As a spec author, I want concrete outcome scenarios, so that each requirement has an observable target.
10. As a spec author, I want current behavior and owner evidence, so that my design starts from the tree.
11. As a spec author, I want a cheapest-wrong evidence target, so that a trivial substitute cannot satisfy the requirement.
12. As a spec author, I want existing sufficient seams, so that I avoid unnecessary design work.
13. As a spec author, I want bounded actions with inspected results, so that evidence guides my next action.
14. As a reviewer, I want unresolved behavior returned, so that the author preserves my decision authority.
15. As a spec author, I want planned evidence for new features, so that specification needs no premature implementation.
16. As a reviewer, I want fresh-session spec evidence, so that the authoring guidance works on a new feature.

### Slice independently verifiable outcomes

17. As a ticket author, I want complete outcome slices, so that each ticket delivers useful behavior.
18. As a ticket author, I want a concrete acceptance scenario, so that completion has an observable meaning.
19. As a ticket author, I want checks usable before successors, so that my ticket can finish independently.
20. As a ticket author, I want predecessor value contracts, so that each blocker has a reason.
21. As a coordinator, I want serial order for shared writes, so that independent outcomes remain separate.
22. As a reviewer, I want justified split and merge decisions, so that the graph contains complete outcomes.
23. As a spec author, I want slicing without implementation, so that the graph precedes the build.
24. As a reviewer, I want fresh-session slicing evidence, so that shared writes do not erase useful outcomes.

### Drive implementation with focused evidence

25. As an implementer, I want a named acceptance target, so that my first action has a verification route.
26. As an implementer, I want the existing TDD sequence, so that I obtain a meaningful behavioral red.
27. As an implementer, I want minimal compiled setup, so that a missing declaration cannot masquerade as behavioral failure.
28. As an implementer, I want honest row classifications, so that I avoid fabricated red results.
29. As an implementer, I want focused reruns after material actions, so that each result guides the next action.
30. As an implementer, I want coherent actions with related edits, so that verification does not require artificial edit fragments.
31. As a reviewer, I want the existing wrong-spec stop, so that the author cannot silently change approved behavior.
32. As a reviewer, I want fresh-session implementation evidence, so that material actions receive focused verification.

### Ground semantic findings in appropriate evidence

33. As a review axis, I want current binding source evidence, so that each finding rests on the frozen subject.
34. As a review axis, I want real-run refutation for runnable claims, so that an easy counterexample defeats a false finding.
35. As a review axis, I want exact-source mandatory-standard evidence, so that an absent automated check does not excuse a violation.
36. As a review axis, I want contrary-source inspection, so that a precise citation still receives refutation.
37. As a reviewer, I want explicit uncertainty, so that an unsupported concern cannot become a repair target.
38. As an implementer, I want findings through the existing disposition, so that repair authorship remains with me.
39. As a reviewer, I want fresh-session review evidence, so that both finding evidence routes work.
40. As a reviewer, I want reconciled adoption evidence, so that each phase closes its own evidence obligation.
41. As a maintainer, I want unchanged prose budgets, so that the guidance stays within its approved size.
42. As a debug author, I want the local six-phase procedure, so that concrete reproduction guidance remains available.
43. As a Coverage reviewer, I want an independent bypass attempt, so that supplied positive checks do not hide a violating state.

### Trial one reviewer across all three axes

44. As a coordinator, I want an explicit unified-review mode, so that one independent reviewer can own all three axes without weakening the default.
45. As an implementation maintainer, I want review misses tied back to the implementation command. I want the workflow prose to improve when it contributed to an issue.

## Implementation decisions

### Existing owners

| Outcome | Rule owner | Integration contract |
| --- | --- | --- |
| Retain debug repair authorship | `.agents/commands/bench-debug.md` | Its ticket supplies guidance, checks, and adoption evidence before its successor |
| Guide spec authoring with concrete evidence | `.agents/skills/bench-craft-spec/SKILL.md` | Its ticket supplies guidance, checks, and adoption evidence before its successor |
| Slice independently verifiable outcomes | `.agents/skills/bench-craft-tickets/SKILL.md` | Its ticket supplies guidance, checks, and adoption evidence before its successor |
| Drive implementation with focused evidence | `.agents/commands/bench-implement-spec.md` | Its ticket supplies guidance, checks, and adoption evidence before its successor |
| Ground semantic findings in appropriate evidence | `.agents/skills/bench-craft-review/SKILL.md` | Its ticket supplies guidance, checks, and adoption evidence before its successor |
| Trial one reviewer across all three axes | `.agents/commands/bench-review-implementation.md` and `internal/reviewrecord` | An explicit plan mode changes reviewer cardinality while retaining axis evidence and participant exclusions |

The debug integration section owns its revised author and fence rules.
The delegate skill points to that section instead of repeating the report contract.
The debug session may obtain read-only diagnostic help without transferring repair authorship.
The coordinator retains report validation and the existing reslice route.

The spec skill owns the authoring sequence.
Seam alternatives remain available when current evidence leaves the seam unresolved.
Craft-seams retains its existing procedure for a genuinely uncertain seam.
The ticket skill owns outcome sizing and the keep, split, or merge decision.
Neither authoring phase implements a feature merely to obtain a red.

The implementation command owns the action and rerun sequence.
It references craft-tdd for row classifications and minimal compiled setup.
It references craft-line for continuation and stop conditions.
The existing bounded repair policy still governs work after initial review.

A material action changes behavior, a verification target, or a premise that determines the next action.
It can contain several related edits.
Inspect the focused result before the next material action.

The review skill owns candidate selection, source derivation, evidence inspection, and the continue-or-stop sequence.
It points to finding-discipline for the runnable-versus-source evidence distinction.
The reference owns that distinction, including the mandatory-standard exception.
Replace its unconditional real-run sentence and reconcile its existing anchor, unit expectation, and mutation fixture.
When a claimed enforcement can admit a bypass, Coverage constructs an independent counterexample that keeps the claimed positive evidence satisfied while violating the requirement.
Replaying only the author's supplied mutations does not satisfy this refutation step.

The completion plan may opt into one independent reviewer for all three axes.
The omitted mode keeps the existing three-distinct-session rule.
The unified reviewer still returns three separately attributable axis results and cannot be the orchestrator or an implementation author.
It re-derives each axis from its primary source within one session; separate fresh contexts remain the default outside this explicit mode.
Each unified result also states whether the issue or miss exposes an improvement to `.agents/commands/bench-implement-spec.md`.

For the exception, cite the binding rule and violating source.
Inspect applicable exceptions and contrary evidence before retaining the finding.
State why executable refutation is unavailable.
A preference remains optional advice under the bounded repair policy.

### Preserve one owner per rule

Refer to existing TDD, repair, and seam rules instead of restating them in the changed phase guidance.
The spec states their required behavior for acceptance, but it does not require a second guidance sentence.
Reuse an existing anchor and fixture when they already protect the required predicate at its owner.
Add an anchor only for a new predicate or a required owner reference.
Keep existing checks and prove each changed fixture still catches its named mutation.

### Local scope and source integrity

Keep the debug phases and loop-constructions reference at their existing locations.
Do not add a universal loop owner, trial flag, metric schema, or comparative speed claim.
Do not change gate authority, authorship authority, or the repair allowance.
Use the current prose budgets through concise replacement text.
The compiled decision map contains the reviewer refinements and retained debug decisions.
The old trial artifacts retire only after that preservation step.

## Implementation chunks

Each ticket is one complete chunk and one serial green commit checkpoint.
Shared registry, fixture, changelog, and evidence writes require this order.
Those shared writes do not merge the independently useful phase outcomes.
Each successor consumes the accepted predecessor tip and its current checks.
Each chunk receives the existing three-axis checkpoint before its successor.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| DG-C1 / `1-retain-debug-authorship.md` | Retain debug repair authorship with its adoption evidence | DG1, DG2, DG3, DG4, DG5, DG6, DG7, DG8, DG42 | workflow anchors, fixture bite, prose budgets, owning adoption task | no |
| DG-CR / `1a-enable-unified-review-trial.md` | Enable an opt-in unified reviewer without changing the default | DG44, DG45 | review-record and delegated-checkpoint tests, prose mechanics | no |
| DG-C2 / `2-guide-spec-evidence.md` | Guide spec authoring with concrete evidence with its adoption evidence | DG9, DG10, DG11, DG12, DG13, DG14, DG15, DG16 | workflow anchors, fixture bite, prose budgets, owning adoption task | no |
| DG-C3 / `3-slice-complete-outcomes.md` | Slice independently verifiable outcomes with its adoption evidence | DG17, DG18, DG19, DG20, DG21, DG22, DG23, DG24 | workflow anchors, fixture bite, prose budgets, owning adoption task | no |
| DG-C4 / `4-drive-implementation-evidence.md` | Drive implementation with focused evidence with its adoption evidence | DG25, DG26, DG27, DG28, DG29, DG30, DG31, DG32 | workflow anchors, fixture bite, prose budgets, owning adoption task | no |
| DG-C5 / `5-ground-semantic-findings.md` | Ground semantic findings in appropriate evidence with its adoption evidence | DG33, DG34, DG35, DG36, DG37, DG38, DG39, DG40, DG41, DG43 | workflow anchors, fixture bite, prose budgets, owning adoption task | yes |

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"debug-loop-guidance-full-20260916","orchestrator_session":"/root","author_limit":1,"assignments":{"1-retain-debug-authorship.md":[{"session":"/root/candidate_a","assignment":"debug-loop-candidate-a","model":"gpt-6-astra","effort":"high","source":"9deb0a7af31712427ff47d6fd0515e8458dbde0d","native_ref":"codex:session/candidate_a@ae21e4a538d7713b147a8a888221063c2ba2ba02"}],"1a-enable-unified-review-trial.md":[{"session":"/root/unified_review_author","assignment":"debug-loop-unified-review-trial","model":"gpt-5.6-sol","effort":"medium","source":"f68d62773d77e7138d8c92ef9a88f5dcaeb3e732","native_ref":"codex:collaboration/spawn_agent/unified_review_author"}],"2-guide-spec-evidence.md":[{"session":"/root/dgc2_author","assignment":"debug-loop-dgc2","model":"gpt-5.6-sol","effort":"medium","source":"8b2073e3104e3b0ddb8a52aec6c7506acf27c7fc","native_ref":"codex:collaboration/spawn_agent/dgc2_author"},{"session":"/root/dgc2_prose_luna","assignment":"debug-loop-dgc2-prose-compaction","model":"gpt-5.6-luna","effort":"medium","source":"1c70037d3a20d61b71b4c0873bf44850e5fa9816","native_ref":"codex:collaboration/spawn_agent/dgc2_prose_luna"},{"session":"/root/dgc2_workflow_terra","assignment":"debug-loop-dgc2-workflow-fix","model":"gpt-5.6-terra","effort":"medium","source":"848f8014dac2efec258433c3a755293bdd846ab6","native_ref":"codex:collaboration/spawn_agent/dgc2_workflow_terra"}],"3-slice-complete-outcomes.md":[],"4-drive-implementation-evidence.md":[],"5-ground-semantic-findings.md":[]}},"chunks":[{"id":"DG-C1","tickets":["1-retain-debug-authorship.md"],"verification":[{"id":"anchors","command":"bench test --check docs-currency-workflow","probe":"restore the blanket write-delegate debug ban","ticket":"1-retain-debug-authorship.md"},{"id":"bite","command":"bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner","ticket":"1-retain-debug-authorship.md"},{"id":"budgets","command":"bench test --check guidance-prose-budgets","ticket":"1-retain-debug-authorship.md"}]},{"id":"DG-CR","tickets":["1a-enable-unified-review-trial.md"],"verification":[{"id":"reviewrecord","command":"bench test --package ./internal/reviewrecord --run 'TestDelegated.*Review'","probe":"omit the explicit unified-review mode while reusing one reviewer","ticket":"1a-enable-unified-review-trial.md"},{"id":"checkpoint","command":"bench test --package ./internal/gate --run TestDelegatedDistinctAxes","ticket":"1a-enable-unified-review-trial.md"},{"id":"prose","command":"bench test --check prose-mechanics","ticket":"1a-enable-unified-review-trial.md"}]},{"id":"DG-C2","tickets":["2-guide-spec-evidence.md"],"verification":[{"id":"anchors","command":"bench test --check docs-currency-workflow","probe":"require an executable red for a new-feature specification","ticket":"2-guide-spec-evidence.md"},{"id":"bite","command":"bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner","ticket":"2-guide-spec-evidence.md"},{"id":"budgets","command":"bench test --check guidance-prose-budgets","ticket":"2-guide-spec-evidence.md"}]},{"id":"DG-C3","tickets":["3-slice-complete-outcomes.md"],"verification":[{"id":"anchors","command":"bench test --check docs-currency-workflow","probe":"merge useful outcomes solely because their writes overlap","ticket":"3-slice-complete-outcomes.md"},{"id":"bite","command":"bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner","ticket":"3-slice-complete-outcomes.md"},{"id":"budgets","command":"bench test --check guidance-prose-budgets","ticket":"3-slice-complete-outcomes.md"}]},{"id":"DG-C4","tickets":["4-drive-implementation-evidence.md"],"verification":[{"id":"anchors","command":"bench test --check docs-currency-workflow","probe":"remove the focused rerun after each material action","ticket":"4-drive-implementation-evidence.md"},{"id":"bite","command":"bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner","ticket":"4-drive-implementation-evidence.md"},{"id":"budgets","command":"bench test --check guidance-prose-budgets","ticket":"4-drive-implementation-evidence.md"}]},{"id":"DG-C5","tickets":["5-ground-semantic-findings.md"],"verification":[{"id":"anchors","command":"bench test --check docs-currency-workflow","probe":"restore unconditional real-run evidence for mandatory standards","ticket":"5-ground-semantic-findings.md"},{"id":"bite","command":"bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner","ticket":"5-ground-semantic-findings.md"},{"id":"budgets","command":"bench test --check guidance-prose-budgets","ticket":"5-ground-semantic-findings.md"}]}],"final_verification":[{"id":"acceptance","command":"bench coverage --check specs/debug-loop-guidance/spec.md"},{"id":"integration","command":"bench test --check docs-currency-workflow"},{"id":"budgets","command":"bench test --check guidance-prose-budgets"}]}
```

## Testing decisions

### Existing sufficient seams

The anchor registry detects missing or reversed guidance at its actual owner.
The existing fixture runner executes the registered check and proves restoration.
These checks prove instruction presence, not agent adoption.
The separate review-owned rows require fresh-session task evidence in `reviews/debug-loop-guidance.md`.
Each evidence excerpt names its source tip, task, first action, observed evidence, stop behavior, and handoff.
A missing or failed excerpt leaves its row open.

Use `internal/anchors/registry_debug_loop.go` for the new anchor group.

Append that group through `internal/anchors/registry_data.go` and register its fixture ownership in `internal/conformance/registry_test.go`.
Keep the existing fixture harness and one diagnostic per omitted rule.

For each automated guidance row, first identify an existing anchor and fixture that protect its predicate.
Reuse that coverage when it is sufficient.
For a missing predicate, add `tests/canary/workflow-guidance-anchors/dg-N/`, where N is the row number.
Each added fixture removes or reverses the owner rule or reference and checks its specific diagnostic.
Record each reused or added fixture beside its coverage row in the review pickup.

For DG6, exercise each missing report field in the fixture family.
For DG34, also update `review-strong-finding-run` to test the runnable arm.
Show the mutated check red and the restored check green.
This demonstrated omission justifies an independent expectation under the project standard.

An anchor addition uses `bench test --check docs-currency-workflow` as its executed root.
`TestEveryRetainedFixtureBitesThroughRegisteredOwner` uses `runFixtureBite` to invoke that registered owner.
`TestGuidanceProseBudgetsHoldOnTheLiveTree` protects DG41 through `guidance-prose-budgets`.
The prose-budget table remains unchanged.
For DG42, compare the old and new debug Phase 1–6 sections and the loop-constructions reference byte for byte.
Retain that comparison beside the debug adoption evidence.

### Seam diagram

```text
changed phase owner --> current anchor registry --> docs-currency-workflow verdict
                               ^
                         omission fixture
fresh session --> phase task --> review-owned evidence in the existing pickup
```

### Acceptance coverage map

Each behavior below states an acceptance predicate, not a mandatory duplicate sentence in a phase owner.
The named canary is the planned fixture when existing coverage does not already protect that predicate.
DG27 and DG28 exercise existing TDD rules through adoption, without new guidance copies.
DG41 uses the existing budget check instead of a new anchor.
Review-owned rows test real behavior and do not infer it from prose presence.

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| DG1 | 1 | The session that owns the debug loop writes its in-scope repair. | `docs-currency-workflow`, canary `dg-1` at `.agents/commands/bench-debug.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A repair dispatch transfers the loop to another author. |
| DG2 | 2 | Debug delegates only read-only diagnostic work. | `docs-currency-workflow`, canary `dg-2` at `.agents/commands/bench-debug.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A scoped fix delegates production writes. |
| DG3 | 3 | A write delegate runs debug through Phase 6 for an in-fence defect. | `docs-currency-workflow`, canary `dg-3` at `.agents/commands/bench-debug.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The phase refuses every write delegate. |
| DG4 | 4 | For an out-of-fence defect, run Phases 1 through 3 before the diagnostic handoff. | `docs-currency-workflow`, canary `dg-4` at `.agents/commands/bench-debug.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The delegate returns a theory without a reproduced symptom. |
| DG5 | 5 | At the fence, stop implementation edits and preserve the in-fence dirty work. | `docs-currency-workflow`, canary `dg-5` at `.agents/commands/bench-debug.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The delegate repairs outside its fence or discards unfinished work. |
| DG6 | 6 | The blocked report includes the command, red output digest, ranked hypotheses, failing surface, and in-fence dirty paths. | `docs-currency-workflow`, canary `dg-6` at `.agents/commands/bench-debug.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A diagnostic receiver lacks the command or one required report field. |
| DG7 | 7 | A fresh-session in-fence debug task retains its author through a verified repair. | review-owned: `reviews/debug-loop-guidance.md`, DG-C1 adoption evidence | An anchor passes while a real task transfers authorship. |
| DG8 | 8 | A fresh-session out-of-fence debug task returns the bounded diagnostic report without a repair outside its fence. | review-owned: `reviews/debug-loop-guidance.md`, DG-C1 adoption evidence | The diagnostic task writes beyond its approved fence. |
| DG9 | 9 | For each approved outcome, state a concrete scenario before choosing its verification seam. | `docs-currency-workflow`, canary `dg-9` at `.agents/skills/bench-craft-spec/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | An outcome family replaces an observable scenario. |
| DG10 | 10 | Inspect the current behavior and its relevant owner before selecting the next authoring action. | `docs-currency-workflow`, canary `dg-10` at `.agents/skills/bench-craft-spec/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A guessed owner drives the design. |
| DG11 | 11 | Name evidence that exposes the cheapest wrong result for the scenario. | `docs-currency-workflow`, canary `dg-11` at `.agents/skills/bench-craft-spec/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A happy-path assertion accepts a no-op implementation. |
| DG12 | 12 | Use a sufficient existing seam before exploring alternatives. | `docs-currency-workflow`, canary `dg-12` at `.agents/skills/bench-craft-spec/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A known seam triggers speculative design work. |
| DG13 | 13 | Choose a bounded authoring action and inspect its result before continuing. | `docs-currency-workflow`, canary `dg-13` at `.agents/skills/bench-craft-spec/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The author repeats analysis without inspecting an artifact. |
| DG14 | 14 | Return unresolved intended behavior to the reviewer before dependent authoring continues. | `docs-currency-workflow`, canary `dg-14` at `.agents/skills/bench-craft-spec/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The author silently decides an open behavior fork. |
| DG15 | 15 | A new-feature spec can plan its evidence without an existing executable red. | `docs-currency-workflow`, canary `dg-15` at `.agents/skills/bench-craft-spec/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The author invents implementation to obtain a red during specification. |
| DG16 | 16 | A fresh-session spec task produces a concrete scenario and evidence plan for a feature with no executable implementation. | review-owned: `reviews/debug-loop-guidance.md`, DG-C2 adoption evidence | The guidance exists but a fresh session demands an executable red. |
| DG17 | 17 | Start each ticket with its delivered outcome and smallest complete behavior, tests, and integration. | `docs-currency-workflow`, canary `dg-17` at `.agents/skills/bench-craft-tickets/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A ticket delivers only tests or only an implementation layer. |
| DG18 | 18 | State a concrete acceptance scenario before locking the ticket. | `docs-currency-workflow`, canary `dg-18` at `.agents/skills/bench-craft-tickets/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A title substitutes for observable completion. |
| DG19 | 19 | Name checks that establish completion while successor tickets remain unbuilt. | `docs-currency-workflow`, canary `dg-19` at `.agents/skills/bench-craft-tickets/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A predecessor needs a future consumer to demonstrate its result. |
| DG20 | 20 | Record each real dependency and the value its predecessor supplies. | `docs-currency-workflow`, canary `dg-20` at `.agents/skills/bench-craft-tickets/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A blocker names order without a value contract. |
| DG21 | 21 | Shared writes determine serial order without automatically merging useful outcomes. | `docs-currency-workflow`, canary `dg-21` at `.agents/skills/bench-craft-tickets/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Two useful outcomes merge only because they touch one registry. |
| DG22 | 22 | Split independently useful outcomes and merge fragments that cannot deliver or verify anything alone. | `docs-currency-workflow`, canary `dg-22` at `.agents/skills/bench-craft-tickets/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The graph preserves a fragment with no standalone result. |
| DG23 | 23 | Ticket slicing plans evidence without requiring implementation or an existing executable red. | `docs-currency-workflow`, canary `dg-23` at `.agents/skills/bench-craft-tickets/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The author implements code merely to approve the graph. |
| DG24 | 24 | A fresh-session slicing task separates useful outcomes that share writes and supplies checks usable before successors exist. | review-owned: `reviews/debug-loop-guidance.md`, DG-C3 adoption evidence | The fresh author merges shared writes or defers all tests. |
| DG25 | 25 | Before an approved slice, identify its acceptance target and existing verification route. | `docs-currency-workflow`, canary `dg-25` at `.agents/commands/bench-implement-spec.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The author edits without identifying the target check. |
| DG26 | 26 | Use craft-tdd for the behavioral-red sequence at an approved TDD seam. | `docs-currency-workflow`, canary `dg-26` at `.agents/commands/bench-implement-spec.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Implementation invents a second red-first policy. |
| DG27 | 27 | Minimal declarations may precede a behavioral red when the test needs them to compile. | review-owned: `reviews/debug-loop-guidance.md`, DG-C4 compiled-setup evidence | A compile error stands in for behavioral failure. |
| DG28 | 28 | Preserve the already covered and not TDD-able classifications under craft-tdd. | review-owned: `reviews/debug-loop-guidance.md`, DG-C4 row-classification evidence | The author manufactures a false red or silently skips a row. |
| DG29 | 29 | After each material implementation action, rerun its focused evidence and inspect the result. | `docs-currency-workflow`, canary `dg-29` at `.agents/commands/bench-implement-spec.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Several material actions accumulate without a discriminating rerun. |
| DG30 | 30 | An implementation action can contain several related edits before its focused rerun. | `docs-currency-workflow`, canary `dg-30` at `.agents/commands/bench-implement-spec.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A one-edit-per-run rule fragments a coherent repair. |
| DG31 | 31 | When evidence contradicts the approved behavior or seam, use the existing stop-short route. | `docs-currency-workflow`, canary `dg-31` at `.agents/commands/bench-implement-spec.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The author changes the target to make a check green. |
| DG32 | 32 | A fresh-session implementation task records a focused result after each of two material actions. | review-owned: `reviews/debug-loop-guidance.md`, DG-C4 adoption evidence | The task records only a final green check. |
| DG33 | 33 | Derive each candidate finding from its current binding source and the frozen implementation. | `docs-currency-workflow`, canary `dg-33` at `.agents/skills/bench-craft-review/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The axis trusts the ticket account without deriving the requirement. |
| DG34 | 34 | For a runnable defect claim, attempt refutation with a real run before reporting a strong finding. | `docs-currency-workflow`, canary `dg-34` at `.agents/skills/bench-craft-review/references/finding-discipline.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The exception removes real-run refutation for runnable claims. |
| DG35 | 35 | For a mandatory standard without an automated check, cite the exact requirement and the violating source. | `docs-currency-workflow`, canary `dg-35` at `.agents/skills/bench-craft-review/references/finding-discipline.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | An unconditional run requirement suppresses a proven standard violation. |
| DG36 | 36 | Inspect contrary source evidence before retaining a mandatory-standard finding without an automated check. | `docs-currency-workflow`, canary `dg-36` at `.agents/skills/bench-craft-review/references/finding-discipline.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A precise quotation appears without an attempt to disprove the claim. |
| DG37 | 37 | An unsupported concern remains uncertain and does not become a confirmed defect. | `docs-currency-workflow`, canary `dg-37` at `.agents/skills/bench-craft-review/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A reviewer upgrades a hunch into a repair target. |
| DG38 | 38 | Route accepted findings through the existing repair disposition. | `docs-currency-workflow`, canary `dg-38` at `.agents/skills/bench-craft-review/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A read-only review edits the repair or changes the allowance. |
| DG39 | 39 | A fresh-session review task retains a source-proven mandatory-standard finding and refutes a runnable claim with a real run. | review-owned: `reviews/debug-loop-guidance.md`, DG-C5 adoption evidence | A fresh axis suppresses the prose defect or skips executable refutation. |
| DG40 | 40 | Final reconciliation cites the completed adoption evidence from all five owning tickets. | review-owned: `reviews/debug-loop-guidance.md`, DG-C5 adoption evidence | The last ticket silently absorbs unfinished earlier adoption work. |
| DG41 | 41 | The changed guidance preserves the existing prose budgets. | `guidance-prose-budgets`: `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetsHoldOnTheLiveTree`) | A budget increase pays for unnecessary guidance growth. |
| DG42 | 42 | Debug keeps its Phase 1 through Phase 6 procedure and local loop-constructions reference. | review-owned: `reviews/debug-loop-guidance.md`, DG-C1 preservation evidence | A universal loop replaces the concrete local procedure. |
| DG43 | 43 | Coverage constructs an independent bypass that preserves claimed positive evidence while violating the requirement when such a state is possible. | `docs-currency-workflow`, canary `dg-43` at `.agents/skills/bench-craft-review/SKILL.md`  `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | The axis replays supplied mutations and misses a contradictory state that still satisfies the positive check. |
| DG44 | 44 | An explicit plan mode permits one independent reviewer to supply all three separately recorded axes while the omitted mode retains distinct sessions. | `internal/reviewrecord/delegated_test.go` (`TestDelegatedRecordVersions`), `internal/gate/delegated_checkpoint_test.go` (`TestDelegatedDistinctAxes`) | The opt-in weakens the default, accepts a participant, or collapses the axes into one result. |
| DG45 | 45 | A unified review reports whether an issue or miss exposes an improvement to the implementation-command prose. | `prose-mechanics`, review-owned: `reviews/debug-loop-guidance.md`, DG-CR trial evidence | A review identifies a workflow-caused miss without evaluating the command that directed the implementation. |

### Fresh-session adoption tasks

Run these tasks during their matching implementation tickets after loading that ticket's changed guidance in a fresh session.
Use a real disposable repository and authorized tasks.
Retain native evidence in the existing review pickup.
Do not count a synthetic transcript or an anchor test as adoption.
Use separate variants when a required stop prevents the successful path.

| Ticket | Concrete task and evidence | Required stop or alternate case |
| --- | --- | --- |
| 1 | Give a delegate a fenced total function with an empty-list defect. Retain reproduction, the same author's repair, and regression evidence. | Give it a second defect whose cause is outside the fence. Retain Phases 1–3, the complete report, and preserved dirty work without out-of-fence edits. |
| 2 | Request a new list-summary feature with no executable check. Retain exact inputs, expected output, current-owner evidence, a sufficient seam, and the inspected authoring result. | Leave the empty-list result unspecified in a variant. Retain the reviewer decision request before dependent design continues. |
| 3 | Supply independently useful summary and export outcomes that share a formatter. Retain separate complete tickets, predecessor value, and checks usable before successors exist. | Supply a test-only fragment. Retain its merger into the behavior slice that makes it useful. Do not implement either feature. |
| 4 | Supply two approved behavior changes at one compiled seam. Retain one focused result after each material action, with several related edits allowed within an action. | Omit a declaration, supply one passing row, and supply one row that cannot execute first. Retain setup, a behavioral red, both classifications, and the unavailable-row reason. |
| 5 | Supply a runnable defect claim that a check refutes, duplicated policy that violates the mandatory one-source standard, and enforcement whose supplied deletion mutation misses an additive contradiction. Retain the real run, exact rule/source evidence, and an independently constructed additive bypass. | Include an unsupported concern and optional advice. Retain uncertainty, applicable-exception checks, and advice outside finding totals. |

The implementation task also has a variant whose acceptance target contradicts the approved source.
Retain its existing wrong-spec exit before dependent work continues.
The review task states why executable refutation is unavailable for the mandatory-standard finding.
Confirmed findings retain the existing disposition and repair author.

DG7 and DG8 own the debug cases, and DG16 owns the spec cases.
DG24 owns the slicing cases.
DG27, DG28, and DG32 own setup, classifications, and the two material-action results.
DG39 owns both review evidence routes and their dispositions.
DG43 owns the independent bypass attempt within the review task.
Missing observations leave their matching row open.

The last ticket reconciles these existing excerpts rather than performing earlier tickets' adoption work.

### Edge inventory

Audience: The guidance serves each repository that links the kit. The conformance fixtures belong to this repository.

| Input or state | Disposition |
| --- | --- |
| In-fence defect | DG3 and DG7 retain the author through repair. |
| Out-of-fence defect | DG4–DG6 and DG8 preserve the diagnostic handoff. |
| Feature has no executable implementation | DG15 and DG16 require planned evidence without a false red. |
| Existing seam is sufficient | DG12 avoids unnecessary alternatives. |
| Intended behavior remains unresolved | DG14 stops dependent authoring. |
| Shared writes with useful independent outcomes | DG21 and DG24 serialize without forced merger. |
| A fragment has no result before its consumer | DG22 requires a merge. |
| A compiled test lacks declarations | DG27 permits minimal setup before behavioral failure. |
| Row already passes or cannot run before implementation | DG28 preserves the two craft-tdd classifications. |
| A material action contains several related edits | DG29 and DG30 require one coherent verification point. |
| Standard has no automated check | DG35 and DG36 require exact-source evidence and attempted refutation. |
| Claim has a runnable check | DG34 retains real-run refutation. |
| Evidence is missing or contradictory | DG37 preserves uncertainty. |
| Adoption evidence is missing, stale, or failed | DG7, DG8, DG16, DG24, DG32, DG39, and DG40 remain open. |

The critical debug repair adds one opt-in anchor normalization for case and ordinary Markdown emphasis.
It is not a general Markdown parser.
No new directory reader, executable hop, environment variable, or package-variable substitution enters this build.
The hostile-input checklist therefore adds no path, empty-directory, symlink, process, or external-service behavior.
The existing gate and fixture readers retain those contracts.

Won't handle: Full Markdown parsing, links, entities, HTML formatting, code-span interpretation, or semantic paraphrase detection.
Won't handle: Debug phase redesign — retain "No red-capable command, no Phase 2" and "Do not proceed until you have reproduced and minimised".
Won't handle: Trial comparison — ordinary review evidence remains the in-scope adoption record.

## Ownership fences

- `.agents/commands/bench-debug.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-review/SKILL.md`
- `.agents/skills/bench-craft-review/references/finding-discipline.md`
- `.agents/skills/bench-craft-spec/SKILL.md`
- `.agents/skills/bench-craft-tickets/SKILL.md`
- `CHANGELOG.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/anchors_command.go`
- `cmd/bench/anchor_help_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry.go`
- `internal/anchors/anchor_harness_diagnostics_test.go`
- `internal/anchors/match.go`
- `internal/anchors/match_test.go`
- `internal/anchors/locate.go`
- `internal/anchors/locate_test.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/gate/delegated_checkpoint_test.go`
- `internal/gate/review_checkpoint_test.go`
- `internal/reviewrecord/coverage.go`
- `internal/reviewrecord/delegated.go`
- `internal/reviewrecord/delegated_test.go`
- `internal/reviewrecord/plan.go`
- `internal/reviewrecord/record.go`
- `projects/benchkit.md`
- `reviews/debug-loop-guidance.md`
- `tests/canary/claude-agent-definitions`
- `tests/canary/guidance-prose-budgets`
- `tests/canary/line-routing`
- `tests/canary/skill-description-budgets`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/workflow-guidance-anchors`

These fences are the union of the five ticket Writes fields.

The command registry paths close the current `internal/tickets/registry_data.go` binding for the anchors package.
Their inclusion does not authorize a new CLI command or a product change.
Keep each unchanged unless its existing inventory requires a mechanical update.

The fixture directories close the current pinning rules for guidance owners.

The review-record paths implement only the explicit unified-review opt-in and retain the distinct-session default.

Build-time guidance rewrites exclude `specs/*/spec.md` and `specs/*/tickets/` by name.

## Out of scope

The reviewer excluded these capabilities from this build.
Each estimate describes a separate future capability, not deferred work needed by these tickets.

| Capability | Estimate | Basis |
| --- | --- | --- |
| Comparative trial controls | 5 edits, 1 gate run | Three phase controls, one trial owner, and one evidence integration. |
| New metric schema | 3 edits, 1 gate run | Schema, writer, and compatibility test. |
| Universal loop extraction | 7 edits, 1 gate run | One common owner, five consumers, and anchor integration. |

Budget increases, changed authority, and causal speed claims are prohibited rather than deferred implementation work.

## Further notes

### Source-sentence-to-row trace

| Source clause | Rows |
| --- | --- |
| Ticket 2: "Preserve debug's core instructions in place." | DG42 |
| Ticket 2: "Give spec writing, implementation, and semantic review their own explicit action, evidence, and stop sequence." | DG9–DG16, DG25–DG39, DG43 |
| Ticket 3: "Preserve the already-approved diagnostic handoff at the delegate fence." | DG1–DG8 |
| Ticket 3: "Use fresh-session task evidence for adoption, without trial flags or a comparative speed claim." | DG7, DG8, DG16, DG24, DG32, DG39, DG40 |
| Ticket 4: "Start with the delivered outcome and identify its smallest complete vertical slice." | DG17–DG19 |
| Ticket 4: "Record real dependencies and the value each predecessor supplies." | DG20 |
| Ticket 4: "Use shared writes to determine order, without treating them as automatic grounds to merge tickets." | DG21, DG22, DG24 |
| Ticket 4: "Do not require implementation during slicing." | DG23 |
| Reviewer refinements: existing finding owner, retained runnable refutation, and independent bypass construction | DG33–DG39, DG43 |
| Reviewer refinements: material-action reruns without one-edit-per-run | DG29, DG30, DG32 |
| Reviewer refinements: prefer a sufficient existing seam | DG12 |
| Research: minimal compiled setup and existing row classifications | DG26–DG28 |
| Scope: no raised prose budgets | DG41 |

### Source and enforcement reads

The author read the compiled map, its four resolved tickets, and its research asset.
The author reopened the three primary papers linked in that asset on 2026-09-15.
Those studies do not establish a Bench speed gain.
The historical Claude conversation was not independently reread; the current reviewer decisions supply the authority.

The source inventory includes all five guidance owners, craft-tdd, craft-delegate, craft-line, finding-discipline, and bounded-repair-policy.
The author also read both spec and ticket templates, map-discipline, STE rules, and the project profile.
The source-owner evidence resolved the seam within the declared twelve-file budget.

The enforcement inventory includes the following files.

- `internal/anchors/registry.go`, `registry_data.go`, `registry_ticket_passes.go`, and `registry_retained_workflow.go`.
- `internal/conformance/registry_test.go`, `fixture_bite_test.go`, `prose_budget_test.go`, and `ticket_grammar_test.go`.
- `internal/tickets/registry_data.go` and the `review-strong-finding-run` mutation fixture.

A repository-wide hidden-file search enumerated the moved decision paths and the retired trial paths.
Only the moved index and research asset required reference updates outside the retired directory.
The debug owner also has a pinned implicit-invocation fixture and Claude agent fixtures.
The finding rule has pins in `registry_data.go`, `registry_data_test.go`, and `review-strong-finding-run`.
The ticket fences include those consumers.

### Proof checklist

- Cited symbols: `Entries`, `EvaluateGroup`, `runFixtureBite`, and the two named test functions resolve in the enforcement inventory.
- Import edges: none change.
- Source-row clauses and occurrences: the trace table assigns the resolved clauses. The reader sweep identifies the moved-path and finding-rule occurrences.
- Promised field labels: the diagnostic report uses its existing command, digest, hypotheses, surface, and dirty-path content. No new schema enters the kit.
- Changed-function callers: none.
- Copy survival: the map moves as one unit. DG42 rejects extraction of debug's local procedure.
- Flagged additions: none beyond the dated reviewer refinements in the compiled decision source.

### Artifact status

This authoring run executes no implementation adoption task and no semantic review.
The mechanical validation record appears in the author return.
The staged spec and ticket graph require the reviewer's implementation approval.
