# Anchor admission prune

Status: staged

Decision source: reviewer-confirmed conversation, 2026-10-16. The reviewer accepted these decisions:

- An anchor earns a place only for one of four reasons. It forbids a stale or wrong form (class FORBID). It names a label or a heading that a parser or a template contract uses (class LABEL). It names a phase, a skill, or a verb as a routing reference (class ROUTE). It carries a platform invariant in `.bench/BENCH.md` (class INVARIANT).
- The wording of a rule sentence is review-owned. No anchor pins it.
- Canary coverage stays at one planted needle for each check family, as ADR 0002 posture 6 states.
- `.bench/BENCH.md` keeps its 29 Forbid anchors and its 9 core Require anchors in Roles, the Bench CLI contract, the four invariants, and the three predicates. It gains one anchor for each of invariants 1, 2, and 3. Its 44 Workflow-section Require anchors go. Under the same rule, its three Require anchors outside the core sections go too: one in Capture and two in How to talk to me.
- A Require needle that pins a `bench …` command line with an argument or a flag stays in this cut (class COMMAND). A sentence that names a verb with no argument counts as a rule sentence.

Verification log: 4 iteration(s) to accept — the spec review round took 2 iterations, and the ticket review round took 2.

## Problem

A rewording of kit guidance turns the gate red even when no rule changes. The anchor registry pins whole sentences of rule wording. A typical skill carries one anchor for every two to four lines. Each pinned sentence also has a test-authored rule and, usually, a canary fixture. One wording fix therefore edits the guidance, a registry row, a test rule, and one to three fixtures.

The registry has no admission rule. The kit guidance tells authors to add a Require row for each owner sentence. So the registry grows with each guidance change, and it never shrinks. A second hard-coded registry in the conformance package restates many of the same diagnostics. Some hand-coded pins enforce the same fact twice.

## Solution

The registry admits only the five classes above: FORBID, LABEL, ROUTE, INVARIANT, and COMMAND. Each Require row declares its class, and a derived test refuses a Require row that has no class or that breaks its class shape. The rows outside the classes leave the registry together with their test-authored rules and their canary fixtures.

The second registry folds into the main registry. Each hand-coded pin that is outside the classes goes, and each pin that enforces a fact twice collapses to one source. ADR 0012 states the admission rule. The guidance that told authors to add a Require row for an owner sentence now says that rule wording is review-owned.

## User stories

Line: opus / medium.
Implementation-line reason: AP-C2 is hardest, because its gate-run second registry must fold with no lost Forbid; AP-C5 adds a class field and retires a kind.
A derived list drives the other deletions, and existing tests grade each one; AP-C6 guidance runs at high effort under the leverage override.
Harder chunks: AP-C2, AP-C5.

### The registry admits only the classes

1. As a maintainer, I want a reworded rule sentence to leave the gate green, so that a guidance edit costs one file.
2. As a maintainer, I want every Forbid row to stay, so that a stale verb still turns the gate red.
3. As a template author, I want a label row to stay when a parser reads it, so that agents write parsable artifacts.
4. As a phase author, I want a route row to stay, so that a command still names its skill, phase, or verb.
5. As a reviewer, I want the core invariant rows on `.bench/BENCH.md` to stay, so that the shipped platform rules stay pinned.
6. As a reviewer, I want an anchor on the key clause of invariant 1, so that its deletion turns the gate red.
7. As a reviewer, I want an anchor on the key clause of invariant 2, so that its deletion turns the gate red.
8. As a reviewer, I want an anchor on the key clause of invariant 3, so that its deletion turns the gate red.
9. As a reviewer, I want the BENCH.md Workflow rule sentences to leave the registry, so that build policy changes at review cost.
10. As a maintainer, I want a `bench …` command row with arguments to stay, so that a command line cannot drift yet.
11. As a maintainer, I want each Require row to declare its class, so that no rule-sentence row enters silently.
12. As a maintainer, I want a Require row with no class to turn the gate red, so that the rule has an executable owner.
13. As a maintainer, I want a LABEL row to name its reader, so that a reviewer can verify the reason.
14. As a maintainer, I want a ROUTE row to be short and to name a route token, so that no sentence passes as a route.
15. As a maintainer, I want INVARIANT rows only in the BENCH.md core sections, so that the class cannot spread.
16. As a maintainer, I want a COMMAND row to hold a `bench` span with an argument, so that the class stays narrow.
17. As a maintainer, I want no Require needle above the STE sentence bound, so that no needle pins two rule sentences.

### The removed rows leave with everything that mirrors them

18. As a maintainer, I want each test rule of a removed row to leave with it, so that no test expects a lost diagnostic.
19. As a maintainer, I want a registry file with no rows to be deleted, so that no empty topic slice remains.
20. As a maintainer, I want each fixture of a removed row to be deleted, so that the fixture-bite test stays green.
21. As a maintainer, I want each canary family to keep a fixture, so that ADR 0002 posture 6 still holds.
22. As a maintainer, I want an ownerless fixture checked before deletion, so that no fixture of a kept check goes.
23. As a reviewer, I want the two fixture-list tests to change only for removed rules, so that kept rules keep their proof.
24. As a maintainer, I want the anchors help fixture to match the kept AGENTS.md rows, so that the help test stays exact.

### One source for each anchor fact

25. As a maintainer, I want the retained-workflow check to read the main registry, so that no second list exists.
26. As a maintainer, I want its stale phrases to become registry Forbid rows, so that the fold loses no Forbid.
27. As a maintainer, I want the line-binding check to own the model identifiers alone, so that a tier change edits one file.
28. As a maintainer, I want each hand-coded rule-sentence pin to go, so that rule wording is review-owned everywhere.
29. As a maintainer, I want each admitted hand-coded pin to stay, so that the admitted facts keep their proof.
30. As a maintainer, I want each marker to have one enforcing check, so that one fact has one source.
31. As a maintainer, I want a derived count of retained-workflow rows, so that a registry edit does not edit a test.
32. As a maintainer, I want the dangling-phrase pin to go, so that no check points at a missing target.

### Match kinds stay live

33. As a maintainer, I want the unused step-scoped kind retired, so that no dead match kind remains.
34. As a maintainer, I want each surviving kind to keep a row, so that the kind tests keep a real subject.
35. As a maintainer, I want the emphasis row and its help test to stay together, so that the help test keeps its subject.

### The query and the lane stay well-formed

36. As an agent, I want `bench anchors <file>` to print an empty state for a file with no rows, so that I see no false failure.
37. As an agent, I want `bench anchors <dir>` to list only files that keep rows, so that the listing stays exact.
38. As a maintainer, I want a non-Markdown file with no rows to take the safe lane class, so that the lane never narrows.

### The decision and the guidance state the rule

39. As a teammate who just walked in, I want ADR 0012 to state the admission rule, so that I know what earns an anchor.
40. As a guidance author, I want the gate skill to stop asking for owner-sentence rows, so that the registry stops growing.
41. As a slicer, I want the retirement pass to ask only for Forbid rows, so that a retirement adds no rule-sentence row.
42. As a reviewer, I want the finding discipline to grade rule wording directly, so that a repair seeks no removed row.
43. As a spec author, I want the needle-length rule to allow a drop, so that the spec skill asks for no extra anchors.
44. As a delegate charger, I want the probe rule to cover only admitted anchors, so that no charge probes a removed row.
45. As a teammate, I want the profile to describe the guidance family by class, so that the profile matches the registry.

### Reviewed exclusions

46. As a reviewer, I want the BENCH.md and AGENTS.md trims to wait for guidance-weight, so that this prune changes no other wording.
47. As a reviewer, I want the COMMAND class to stay until a grammar check replaces it, so that this cut adds no seam.

### Kept command needles fit the bound

48. As a maintainer, I want each over-bound COMMAND needle re-cut to its `bench …` sentence, so that a kept row passes the sentence bound.

### Kept Forbid checks keep their scope

49. As a maintainer, I want the integration-source stale phrases kept as Forbid rows, so that a retired phrase still turns the gate red.

## Implementation decisions

### The admission classes

The registry `Anchor` type gains one class field. The class takes one of five values: FORBID, LABEL, ROUTE, INVARIANT, and COMMAND. A Forbid kind implies the FORBID class, so a Forbid row does not declare it. Each Require row declares one of the other four classes. The zero value means unclassified.

A LABEL row also names its reader: the Go parser or the template that reads the label. The reader is a short text value for review, not a path that a check resolves.

The admission test lives in the anchors package. It walks `Entries()` and applies one shape rule for each class:

- No Require row is unclassified.
- An INVARIANT row guards `.bench/BENCH.md`, and its section is Roles, The Bench CLI contract, or the four-invariants section. A row whose section field is empty uses the located H2 of its needle.
- A LABEL row names a non-empty reader.
- A ROUTE needle has six words or fewer, and it contains a phase, skill, or verb token: `/bench-`, `$bench-`, `craft-`, `bench-craft-`, or `bench `.
- A COMMAND needle contains a `bench ` span followed by an argument or a flag.
- No Require needle has more words than `prose.MaxSentenceWords`.

The word bound comes from `internal/prose`, so the test restates no bound. The test also counts words the way `internal/prose` counts them, where a code span is one word. If `internal/prose` exports no counter for one sentence, the AP-C5 ticket exports one there.

By a whitespace count, 11 COMMAND needles at the spec base are longer than the bound, and one holds a newline. The `internal/prose` count is lower, because a code span is one word. Each pruning ticket finds the needles in its registry file that are over the bound under the `internal/prose` count. It re-cuts each one to the sentence that holds its `bench …` span. The row stays and its needle shrinks, as the craft-spec rubric prescribes for a long needle. This re-cut is a non-behavioral call, flagged for reviewer veto; the alternative is a COMMAND exemption from the bound.

This test is the cheapest red-capable owner. A new rule-sentence row fails the unclassified rule at once. A misclassified row fails its shape rule.

### The derived snapshot

`assets/classified.md` is the classification of each registry row at the spec base. Each ticket re-derives its slice against the tree at its own base, and the tree wins. The snapshot is evidence for review, not an input that a check reads.

### The second registry

The retained-workflow check keeps its document reads. It takes its expected diagnostics from the main registry, by its group and its diagnostic prefix. Each stale phrase that it forbids moves into the registry as a Forbid row on the same file. Its hard-coded model identifiers leave, because the line-binding check already owns that fact. The check then holds no list that the registry also holds.

The integration-source stale-phrase check takes its file list from the integration-source Require rows, and AP-C1 removes 9 of those 10 rows. So AP-C1 moves each stale phrase of that check into the registry. Each phrase becomes a Forbid row on each file that the check scanned at the ticket base. The check and its count pin then leave, because the registry holds the same fact.

### Hand-coded pins

A hand-coded pin stays only when its text is in one of the five classes, or when a program reads the text. The structured-phase clause check, the shared-rule single-source check, the skills-index checks, and the `BENCH_MODEL` pin stay. A marker that both the registry and the shared-rule check enforce keeps only the shared-rule check.

That check asserts each marker in two ways: present in `.bench/BENCH.md`, and absent from `AGENTS.md` and `README.md`. One of its markers is the Roles needle of a core row. So the registry drops that row, and `.bench/BENCH.md` keeps 8 core Require rows in the registry and the Roles sentence through the shared-rule check. This is a non-behavioral call, flagged for reviewer veto.

The exported marker constants stay, because that check reads them. That check pins some Workflow markers of `.bench/BENCH.md` as a cross-file single source, so those few marker phrases stay pinned there and nowhere else.

### Removal census

A row whose removal leaves no test that still expects it cannot turn red, because its row, its test rule, and its fixtures leave together. So each pruning ticket record carries a removal census. The census lists each removed row and shows that the row is DROP in `assets/classified.md` at the ticket base, re-derived against the tree. For the retained-workflow fold, the census also shows that each stale phrase now has a registry Forbid row. The review grades the census.

### Match kinds

All nine step-scoped Require rows are rule sentences, so they leave. The kind tests need one row of each declared kind, so the nine rows leave in AP-C5 together with the kind, not in AP-C3. The kind, its step field, and its step narrowing leave the anchors package. The `bench anchors` table then has no `step` column, and `CHANGELOG.md` records that change. This is a non-behavioral call, flagged for reviewer veto. The alternative is to keep one step-scoped row and the column.

### Fixtures

A fixture leaves when its expected diagnostic belongs to a removed row. The fixture-bite test derives its set from the tree, so it needs no edit. The handoff-fixture list test and the cadence mutation test hold independent lists. Their entries for removed rules leave with those rules. This is the reviewer-decided retirement of those rules, not a weakening of a test to pass.

`assets/fixture-census.md` lists, for each ticket, the fixture directories that the ticket deletes. Each ticket names `tests/canary/workflow-guidance-anchors` as one `Writes:` prefix, and its author re-derives the list against the tree at the ticket base.

The anchors help test drops every second `AGENTS.md` Require row and expects a failure. After AP-C2 and AP-C1, `AGENTS.md` keeps one Require row, so the test drops nothing. The AP-C1 ticket reworks that test and its fixture repository in the same commit.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| AP-C2 / `1-fold-retained-workflow.md` | The retained-workflow check reads the main registry, and its rule rows leave. | AP02, AP03, AP04, AP10, AP25, AP26, AP27, AP31, AP46 | `bench test --package ./internal/anchors`, `bench test --check docs-currency-workflow`, `bench test --package ./internal/conformance`, `bench test --package ./cmd/bench`, `bench test --check root-conformance` | yes |
| AP-C1 / `2-prune-core-registry.md` | The core registry file keeps only admitted rows, BENCH.md keeps its core rows, the three invariant anchors exist, and the integration-source check keeps its scope. | AP02, AP03, AP04, AP05, AP06, AP07, AP08, AP10, AP18, AP20, AP21, AP22, AP23, AP24, AP46, AP47 | `bench test --package ./internal/anchors`, `bench test --package ./internal/conformance`, `bench test --check docs-currency-workflow`, `bench test --check conformance-canary-families`, `bench test --package ./cmd/bench`, `bench test --check root-conformance` | no |
| AP-C3 / `3-prune-build-registries.md`, `4-prune-phase-registries.md` | The build and phase topic registries keep only admitted rows, and empty topic files are gone. | AP02, AP03, AP04, AP10, AP18, AP19, AP20, AP46 | `bench test --package ./internal/anchors`, `bench test --check docs-currency-workflow`, `bench test --package ./internal/conformance`, `bench test --package ./cmd/bench`, `bench test --check root-conformance` | no |
| AP-C4 / `5-prune-hand-pins.md` | The hand-coded pins keep only admitted facts, each with one source. | AP28, AP29, AP30, AP32 | `bench test --package ./internal/conformance`, `bench test --check docs-currency-workflow`, `bench test --check root-conformance` | no |
| AP-C5 / `6-retire-step-kind.md`, `7-admit-by-class.md` | The step kind is retired, each Require row declares its class, and the admission test refuses an unclassified row. | AP01, AP09, AP11, AP12, AP13, AP14, AP15, AP16, AP17, AP33, AP34, AP35, AP36, AP37, AP38 | `bench test --package ./internal/anchors`, `bench test --package ./cmd/bench`, `bench test --package ./internal/gate`, `bench test --package ./internal/conformance`, `bench test --package ./internal/prose`, `bench test --check docs-currency-workflow`, `bench test --check root-conformance` | yes |
| AP-C6 / `8-state-admission-rule.md` | ADR 0012 and the guidance state the admission rule. | AP39, AP40, AP41, AP42, AP43, AP44, AP45 | `bench test --check docs-currency-workflow`, `bench test --check prose-mechanics`, `bench test --check guidance-prose-budgets`, `bench test --package ./internal/anchors`, `bench test --package ./internal/conformance`, `bench test --package ./cmd/bench`, `bench test --check root-conformance` | no |

AP-C2 goes first, because the retained-workflow check hard-codes diagnostics of rows in the core and topic registries. A prune of those rows before the fold turns that check red. AP-C1 follows and holds the BENCH.md rows and the invariant anchors. AP-C5 waits for AP-C1 to AP-C4, because the admission test reds the live tree until the last unadmitted row goes. AP-C6 waits for AP-C5.

Each ticket that writes `internal/anchors` names the five command-binding files. The preflight anchor-closure rule also makes a ticket that writes a guarded guidance file name every registry file that pins it. AP-C6 edits guidance files that many registry rows pin, so its `Writes:` line names `internal/anchors` files too. So all eight tickets run in series.

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"anchor-admission-prune-build-20261017","orchestrator_session":"claude:aap-orchestrator-20261017","author_limit":1,"assignments":{"1-fold-retained-workflow.md":[{"session":"claude:aap_t01","assignment":"d537417a1810707c0a2deed6ddb6640e","model":"opus","effort":"medium","source":"62f5ec4919813f2eb83c019e9c4234ed41c7cb53","native_ref":"claude-agent:aap_t01"}],"2-prune-core-registry.md":[{"session":"claude:aap_t02","assignment":"d537417a1810707c0a2deed6ddb6640e","model":"opus","effort":"medium","source":"45cf3764f0ebdc2f2ad427fc8b400ad72261fe7d","native_ref":"claude-agent:aap_t02"}],"3-prune-build-registries.md":[],"4-prune-phase-registries.md":[],"5-prune-hand-pins.md":[],"6-retire-step-kind.md":[],"7-admit-by-class.md":[],"8-state-admission-rule.md":[]}},"chunks":[{"id":"AP-C2","tickets":["1-fold-retained-workflow.md"],"verification":[{"id":"t1-anchors","command":"bench test --package ./internal/anchors","ticket":"1-fold-retained-workflow.md"},{"id":"t1-docs-currency-workflow","command":"bench test --check docs-currency-workflow","ticket":"1-fold-retained-workflow.md"},{"id":"t1-conformance","command":"bench test --package ./internal/conformance","ticket":"1-fold-retained-workflow.md"},{"id":"t1-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"1-fold-retained-workflow.md"},{"id":"t1-root-conformance","command":"bench test --check root-conformance","ticket":"1-fold-retained-workflow.md"}]},{"id":"AP-C1","tickets":["2-prune-core-registry.md"],"verification":[{"id":"t2-anchors","command":"bench test --package ./internal/anchors","ticket":"2-prune-core-registry.md"},{"id":"t2-conformance","command":"bench test --package ./internal/conformance","ticket":"2-prune-core-registry.md"},{"id":"t2-docs-currency-workflow","command":"bench test --check docs-currency-workflow","ticket":"2-prune-core-registry.md"},{"id":"t2-conformance-canary-families","command":"bench test --check conformance-canary-families","ticket":"2-prune-core-registry.md"},{"id":"t2-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"2-prune-core-registry.md"},{"id":"t2-root-conformance","command":"bench test --check root-conformance","ticket":"2-prune-core-registry.md"}]},{"id":"AP-C3","tickets":["3-prune-build-registries.md","4-prune-phase-registries.md"],"verification":[{"id":"t3-anchors","command":"bench test --package ./internal/anchors","ticket":"3-prune-build-registries.md"},{"id":"t3-docs-currency-workflow","command":"bench test --check docs-currency-workflow","ticket":"3-prune-build-registries.md"},{"id":"t3-conformance","command":"bench test --package ./internal/conformance","ticket":"3-prune-build-registries.md"},{"id":"t3-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"3-prune-build-registries.md"},{"id":"t3-root-conformance","command":"bench test --check root-conformance","ticket":"3-prune-build-registries.md"},{"id":"t4-anchors","command":"bench test --package ./internal/anchors","ticket":"4-prune-phase-registries.md"},{"id":"t4-docs-currency-workflow","command":"bench test --check docs-currency-workflow","ticket":"4-prune-phase-registries.md"},{"id":"t4-conformance","command":"bench test --package ./internal/conformance","ticket":"4-prune-phase-registries.md"},{"id":"t4-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"4-prune-phase-registries.md"},{"id":"t4-root-conformance","command":"bench test --check root-conformance","ticket":"4-prune-phase-registries.md"}]},{"id":"AP-C4","tickets":["5-prune-hand-pins.md"],"verification":[{"id":"t5-conformance","command":"bench test --package ./internal/conformance","ticket":"5-prune-hand-pins.md"},{"id":"t5-docs-currency-workflow","command":"bench test --check docs-currency-workflow","ticket":"5-prune-hand-pins.md"},{"id":"t5-root-conformance","command":"bench test --check root-conformance","ticket":"5-prune-hand-pins.md"}]},{"id":"AP-C5","tickets":["6-retire-step-kind.md","7-admit-by-class.md"],"verification":[{"id":"t6-anchors","command":"bench test --package ./internal/anchors","ticket":"6-retire-step-kind.md"},{"id":"t6-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"6-retire-step-kind.md"},{"id":"t6-gate","command":"bench test --package ./internal/gate","ticket":"6-retire-step-kind.md"},{"id":"t6-conformance","command":"bench test --package ./internal/conformance","ticket":"6-retire-step-kind.md"},{"id":"t6-prose","command":"bench test --package ./internal/prose","ticket":"6-retire-step-kind.md"},{"id":"t6-docs-currency-workflow","command":"bench test --check docs-currency-workflow","ticket":"6-retire-step-kind.md"},{"id":"t6-root-conformance","command":"bench test --check root-conformance","ticket":"6-retire-step-kind.md"},{"id":"t7-anchors","command":"bench test --package ./internal/anchors","ticket":"7-admit-by-class.md"},{"id":"t7-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"7-admit-by-class.md"},{"id":"t7-gate","command":"bench test --package ./internal/gate","ticket":"7-admit-by-class.md"},{"id":"t7-conformance","command":"bench test --package ./internal/conformance","ticket":"7-admit-by-class.md"},{"id":"t7-prose","command":"bench test --package ./internal/prose","ticket":"7-admit-by-class.md"},{"id":"t7-docs-currency-workflow","command":"bench test --check docs-currency-workflow","ticket":"7-admit-by-class.md"},{"id":"t7-root-conformance","command":"bench test --check root-conformance","ticket":"7-admit-by-class.md"}]},{"id":"AP-C6","tickets":["8-state-admission-rule.md"],"verification":[{"id":"t8-docs-currency-workflow","command":"bench test --check docs-currency-workflow","ticket":"8-state-admission-rule.md"},{"id":"t8-prose-mechanics","command":"bench test --check prose-mechanics","ticket":"8-state-admission-rule.md"},{"id":"t8-guidance-prose-budgets","command":"bench test --check guidance-prose-budgets","ticket":"8-state-admission-rule.md"},{"id":"t8-anchors","command":"bench test --package ./internal/anchors","ticket":"8-state-admission-rule.md"},{"id":"t8-conformance","command":"bench test --package ./internal/conformance","ticket":"8-state-admission-rule.md"},{"id":"t8-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"8-state-admission-rule.md"},{"id":"t8-root-conformance","command":"bench test --check root-conformance","ticket":"8-state-admission-rule.md"}]}],"final_verification":[{"id":"coverage-check","command":"bench coverage --check specs/anchor-admission-prune/spec.md"},{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"conformance","command":"bench test --package ./internal/conformance"},{"id":"docs-currency-workflow","command":"bench test --check docs-currency-workflow"},{"id":"conformance-canary-families","command":"bench test --check conformance-canary-families"},{"id":"conformance-meta","command":"bench test --check conformance-meta"},{"id":"root-conformance","command":"bench test --check root-conformance"}]}
```

## Testing decisions

- A good admission test walks the live registry and applies the class shape rules. It takes its word bound from `internal/prose`, and it takes no count from the code under test.
- Each removed row needs no new test. The row, its rule, and its fixtures leave together, and the existing tests stay green. The removal census in each pruning ticket record proves that each removed row was DROP.
- Each kept row keeps its test-authored rule and, where one exists, its fixture. A kept rule must still turn red when its needle is deleted.
- Each new invariant anchor gets one test-authored rule in the core registry test. The rule proves that a deletion of the clause raises the diagnostic.
- The canary family check proves posture 6. The fixture-bite test proves that each kept fixture still bites.
- `bench anchors` empty states and the lane path class use their existing tests in `cmd/bench` and `internal/gate`.

### Seam diagram

    trigger: bench gate, bench commit lane, bench anchors
        │
        ▼
    guidance files  ──▶  [ anchors.Entries() → EvaluateGroup / admission test ]  ──▶  diagnostics
                      ◀ tests attach here: per-topic rule tests, admission test, retained-workflow check, fixture-bite test

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| AP01 | 1,11 | A Require row with no class makes the admission test fail | planned TestAdmissionClasses in internal/anchors | A registry that accepts a new rule-sentence row stays green without the class rule. |
| AP02 | 2 | Every Forbid-kind row at the ticket base is still in the registry after the ticket | review-owned removal census in each pruning ticket record | A prune that deletes a Forbid row with its test rule stays green, so the census must show each removed row as DROP. |
| AP03 | 3 | Each LABEL row at the ticket base is still in the registry after the ticket | review-owned removal census in each pruning ticket record | A deleted label row and its test rule leave together and stay green, so the census must show each removed row as DROP. |
| AP04 | 4 | Each ROUTE row at the ticket base is still in the registry after the ticket | review-owned removal census in each pruning ticket record | A deleted route row and its test rule leave together and stay green, so the census must show each removed row as DROP. |
| AP05 | 5 | The 8 core Require rows on `.bench/BENCH.md` are still in the registry, and the shared-rule check still requires the Roles marker in `.bench/BENCH.md` | registry_data_test.go core tuple test and `internal/conformance/validity_checks_test.go` (`checkSharedRuleSingleSource`) | A prune that drops a core row breaks the exact tuple list, and a deleted Roles sentence raises the shared-rule diagnostic. |
| AP06 | 6 | Deletion of "Never edit, skip, weaken, or delete a test or a check to make it pass." raises the invariant diagnostic | `internal/anchors/registry_data_test.go` (`TestInvariantClauseAnchorsRedOnRemoval`) | No row reds when the invariant 1 clause goes. |
| AP07 | 7 | Deletion of "Before any multi-cycle stage, state the model, the effort, the iteration cap, and a one-clause justification in one line." raises the invariant diagnostic | `internal/anchors/registry_data_test.go` (`TestInvariantClauseAnchorsRedOnRemoval`) | No row reds when the invariant 2 clause goes. |
| AP08 | 8 | Deletion of `Put no file paths and no code snippets in an ADR; they rot.` raises the invariant diagnostic | `internal/anchors/registry_data_test.go` (`TestInvariantClauseAnchorsRedOnRemoval`) | No row reds when the invariant 3 clause goes. |
| AP09 | 9 | No Require row guards a `.bench/BENCH.md` line outside Roles, The Bench CLI contract, and the four-invariants section | planned TestAdmissionClasses INVARIANT section rule | A kept Workflow, Capture, or How to talk to me row fails the INVARIANT section rule. |
| AP10 | 10 | Each COMMAND row at the ticket base is still in the registry after the ticket | review-owned removal census in each pruning ticket record | A deleted command row and its test rule leave together and stay green, so the census must show each removed row as DROP. |
| AP11 | 12 | The admission test fails on the live tree when one unclassified Require row is added | planned TestAdmissionClasses mutation case | A test that reads only kept rows cannot see a new row. |
| AP12 | 13 | A LABEL row with an empty reader makes the admission test fail | planned TestAdmissionClasses table case | A label row without a reader passes review unchecked. |
| AP13 | 14 | A ROUTE needle of seven words makes the admission test fail | planned TestAdmissionClasses table case | A rule sentence passes as a route without the word bound. |
| AP14 | 15 | An INVARIANT row that guards a file other than `.bench/BENCH.md` makes the admission test fail | planned TestAdmissionClasses table case | The class spreads to other files without the file rule. |
| AP15 | 15 | An INVARIANT row in the Workflow section of `.bench/BENCH.md` makes the admission test fail | planned TestAdmissionClasses table case | A Workflow rule sentence passes as an invariant without the section rule. |
| AP16 | 16 | A COMMAND needle with no `bench ` span followed by an argument or a flag makes the admission test fail | planned TestAdmissionClasses table case | A rule sentence passes as a command without the span rule. |
| AP17 | 17 | A Require needle of `prose.MaxSentenceWords` plus one words makes the admission test fail | planned TestAdmissionClasses table case | A multi-sentence needle passes without the bound. |
| AP18 | 18 | The per-topic tests hold no rule whose row is gone | per-topic registry tests in internal/anchors | A rule for a removed row expects a diagnostic that no row emits, and its test fails. |
| AP19 | 19 | No registry topic file has zero rows | review-owned file census in the AP-C3 ticket record | `Entries()` flattens the topic slices, so no executable reader sees an empty slice; the census names each deleted file. |
| AP20 | 20 | Every kept fixture bites, and no fixture expects a removed diagnostic | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A fixture for a removed row no longer bites, and the test fails. |
| AP21 | 21 | Each registered canary family has at least one fixture directory | conformance-canary-families check | A prune that deletes the last fixture of a family reds the family check. |
| AP22 | 22 | Each of the fixtures whose expectation matches no literal keeps its fixture when a kept check owns it | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A deleted fixture of a kept check is review-visible, and a kept fixture without an owner fails to bite. |
| AP23 | 23 | The handoff-fixture list test and the cadence mutation test name only kept rules | `internal/conformance/fixture_bite_test.go` (`TestSpecTicketHandoffWorkflowFixturesAreComplete`) and `internal/conformance/fixture_bite_test.go` (`TestWorkflowCadenceAnchorsRejectDeletionAndSwap`) | A list entry for a removed rule fails, because its diagnostic no longer exists. |
| AP24 | 24 | The anchors help fixture repository places each kept AGENTS.md needle and no other | `cmd/bench/anchor_help_test.go` (`TestAnchorsReportsAbsentNeedles`) | A fixture that keeps a removed needle or lacks a kept needle fails the exact placement. |
| AP25 | 25 | The retained-workflow check holds no diagnostic literal that the registry also holds | `internal/conformance/retained_workflow_test.go` (`TestRetainedWorkflowReadsRegistry`) | A surviving hard-coded diagnostic list stays green without a copy rule. |
| AP26 | 26 | Each stale phrase that the retained-workflow check forbade at the ticket base is a registry Forbid row | review-owned removal census in the AP-C2 ticket record | A fold that drops a phrase and its rule together stays green, so the census must pair each old phrase with its new Forbid row. |
| AP27 | 27 | The retained-workflow check holds no model identifier, and the line-binding check still compares the profile prose with the binding | review-owned census in the AP-C2 ticket record, and the existing line-binding canary | A second identifier list survives only by review omission, and a stale profile line still reds the line-binding check. |
| AP28 | 28 | No hand-coded conformance pin pins a rule sentence outside the classes | review-owned pin census in the AP-C4 ticket record | A review that skips a pin leaves it, so the census lists each of the 25 sites with its disposition. |
| AP29 | 29 | Each hand-coded LABEL, ROUTE, FORBID, or INVARIANT pin still reds when its text goes | existing tests of each kept pin in internal/conformance | A kept pin that was deleted by mistake loses its test, and the test fails to compile or fails. |
| AP30 | 30 | Each BENCH.md marker has exactly one enforcing check | review-owned pin census in the AP-C4 ticket record | Some markers sit in the four-invariants section, where an INVARIANT row passes the section rule, so review must confirm that no registry twin remains. |
| AP31 | 31 | The retained-workflow test derives its expected row count from the registry | `internal/conformance/retained_workflow_test.go` (`TestRetainedWorkflowReadsRegistry`) | A hard-coded count fails when one registry row changes. |
| AP32 | 32 | No conformance pin requires "implementation-line recommendation contract" | review-owned pin census in the AP-C4 ticket record | The dangling pin is a rule-sentence pin, so the census names its removal. |
| AP33 | 33 | The anchors package declares no step-scoped kind, no step field, and no step narrowing | review-owned retirement census in the AP-C5 ticket record | The kind tests walk a hard-coded kind list, so a build that keeps the kind and drops it from that list stays green; the census shows each step symbol gone. |
| AP34 | 34 | Each surviving match kind has at least one registry row | `internal/anchors/anchor_harness_diagnostics_test.go` (`TestEvaluatePathAnchorKinds`) | A kind with no row fails the kind-coverage test. |
| AP35 | 35 | The case-folded emphasis row and its help test still agree | `cmd/bench/anchor_help_test.go` (`TestAnchorsReportsCaseFoldedEmphasisViolation`) | A removed emphasis row breaks the help test subject. |
| AP36 | 36 | `bench anchors <file>` on a file with no rows prints a zero-row table at exit 0 | cmd/bench anchors_command test | A file with no rows that errors or prints nothing fails the empty-state case. |
| AP37 | 37 | `bench anchors <dir>` lists only files that keep at least one row | cmd/bench anchors_dir_test.go | A listing that keeps a file with no rows breaks the derived count. |
| AP38 | 38 | A non-Markdown file with no rows selects the unknown path class in the lane | internal/gate lane_select_test.go | A file that falls to no class narrows the lane below its declared checks. |
| AP39 | 39 | ADR 0012 states the five classes, review-owned rule wording, and one needle per family | review-owned ADR text | The ADR text has no executable reader, so review grades it. |
| AP40 | 40 | The gate-authoring skill asks for a Forbid row only for a stale or wrong form | review-owned guidance text | No check reads this guidance, so review grades it. |
| AP41 | 41 | The slicing retirement pass asks for a Forbid row only for a stale or wrong form | review-owned guidance text | No check reads this guidance, so review grades it. |
| AP42 | 42 | The finding discipline says that review grades rule wording directly | review-owned guidance text | No check reads this guidance, so review grades it. |
| AP43 | 43 | The spec skill allows a drop of a long needle that holds rule wording | review-owned guidance text | No check reads this guidance, so review grades it. |
| AP44 | 44 | The delegation discipline applies the probe rule to admitted anchors | review-owned guidance text | No check reads this guidance, so review grades it. |
| AP45 | 45 | The project profile describes the workflow-guidance family by its admitted classes | review-owned profile text | No check reads this text, so review grades it. |
| AP47 | 49 | Each stale phrase that the integration-source check forbade at the ticket base is a registry Forbid row on each file that the check scanned | review-owned removal census in the AP-C1 ticket record | A prune that drops the check and its phrases together stays green, so the census must pair each old phrase and file with its new Forbid row. |
| AP46 | 48 | Each COMMAND needle that was longer than `prose.MaxSentenceWords` is re-cut to the one sentence that holds its `bench …` span | review-owned re-cut census in each pruning ticket record, then AP17 in AP-C5 | A ticket that keeps a long needle leaves AP-C5 red on the live tree, and the census names each re-cut. |

Not covered: story 46 — the trims belong to the guidance-weight re-slice, which is outside this spec.
Not covered: story 47 — the COMMAND class is a kept state, and AP10 and AP16 already observe it.

### Edge inventory

- A Require needle that holds a backticked command with no argument, such as `bench status`: it is a rule sentence, so it leaves. Rows AP16 and AP01 observe it.
- A LABEL that only a template uses, with no Go parser: it stays as LABEL, and its reader names the template. Row AP12 observes it.
- A needle that wraps across lines in the guidance: the matcher collapses whitespace, so the invariant needles match. Rows AP06 to AP08 observe it.
- A guarded document that loses every row: it leaves the `bench anchors <dir>` listing. Rows AP36 to AP38 observe it.
- A fixture whose expected text matches no literal: it needs an owner check. Row AP22 observes it.
- **Won't handle** a paraphrased restatement of a removed rule — review owns rule wording under the amended ADR 0012.
- **Won't handle** a grammar check for `bench …` command spans — it is a separate capability, priced under Out of scope.
- **Won't handle** CRLF guidance files — the matcher already normalizes them, and this prune changes no matcher behavior.

## Ownership fences

These paths are the union of the ticket `Writes:` lines, plus the review pickup:

- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `.agents/skills/bench-craft-gate/SKILL.md`
- `.agents/skills/bench-craft-review/references/finding-discipline.md`
- `.agents/skills/bench-craft-spec/SKILL.md`
- `.agents/skills/bench-craft-tickets/references/slicing-checks.md`
- `CHANGELOG.md`
- `cmd/bench/anchor_help_test.go`
- `cmd/bench/anchors_command.go`
- `cmd/bench/anchors_dir_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/testdata/anchors/fixture-repo/AGENTS.md`
- `docs/adr/0012-document-anchors-use-a-declarative-substring-registry.md`
- `internal/anchors`
- `internal/anchors/anchor_harness_diagnostics_test.go`
- `internal/anchors/registry_calibration.go`
- `internal/anchors/registry_calibration_test.go`
- `internal/anchors/registry_charge_binding.go`
- `internal/anchors/registry_charge_binding_test.go`
- `internal/anchors/registry_check_floor.go`
- `internal/anchors/registry_check_floor_test.go`
- `internal/anchors/registry_chunk_chain.go`
- `internal/anchors/registry_chunk_chain_test.go`
- `internal/anchors/registry_commitment.go`
- `internal/anchors/registry_craft_research.go`
- `internal/anchors/registry_craft_research_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_decision_maps.go`
- `internal/anchors/registry_decision_maps_test.go`
- `internal/anchors/registry_front_door.go`
- `internal/anchors/registry_front_door_test.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/anchors/registry_review_rules.go`
- `internal/anchors/registry_review_rules_test.go`
- `internal/anchors/registry_spec_trace.go`
- `internal/anchors/registry_spec_trace_test.go`
- `internal/anchors/registry_ticket_passes.go`
- `internal/anchors/registry_ticket_passes_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/charge_evidence_guidance_test.go`
- `internal/conformance/checks_test.go`
- `internal/conformance/commitment_guidance_test.go`
- `internal/conformance/data_handling_test.go`
- `internal/conformance/docs_workflow_checks_test.go`
- `internal/conformance/docs_workflow_helpers_test.go`
- `internal/conformance/fixture_bite_test.go`
- `internal/conformance/ft311_preparation_test.go`
- `internal/conformance/implementation_continuation_test.go`
- `internal/conformance/recurrence_maintenance_contract_test.go`
- `internal/conformance/retained_workflow_test.go`
- `internal/conformance/skills_index_checks_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/conformance/tier_live_tree_test.go`
- `internal/conformance/validity_checks_test.go`
- `internal/gate/lane_select_test.go`
- `internal/gate/lane_test.go`
- `internal/prose`
- `projects/benchkit.md`
- `tests/canary/guidance-prose-budgets`
- `tests/canary/line-routing`
- `tests/canary/skill-description-budgets`
- `tests/canary/skills-index-command-adapters`
- `tests/canary/workflow-guidance-anchors`
- `reviews/anchor-admission-prune.md`

## Out of scope

- A grammar check that parses each `bench …` span in guidance against `bench <verb> --help` and replaces the COMMAND class: about 6 edits, 2 gate runs.
- The BENCH.md and AGENTS.md text trims from the 2026-10-16 Bench-file review: these go into the guidance-weight re-slice, about 12 edits, 2 gate runs.

## Further notes

### Sequence and staged-spec conflicts

This prune lands before `specs/guidance-weight` re-slices. The reviewer confirmed that order on 2026-10-16. These staged specs overlap the prune:

- `specs/guidance-weight` moves 21 BENCH.md Require rows, repoints 5 Forbid rows, and adds 3 Require rows. The admission rule refuses GW04 and GW21. The reviewer decided on 2026-10-16 that GW03 gets no anchor. Its fence holds about 283 workflow-guidance fixtures, and its map drift line triggers on this prune.
- `specs/preflight-spec-staleness` ticket 05 keeps co-owned anchor and canary assertions. A build of it before this prune keeps rows that this prune deletes.
- `specs/markdown-block-reader` changes only the matcher. Its fence overlaps four registry files.
- `specs/exact-planning-readers` adds a reader-specific registry file. Its needles need a class check before it builds.
- `specs/process-lifetime` tickets 01 and 07 name the retained-workflow registry file. A build of it before this prune edits rows that this prune folds or deletes.

### Roadmap rows for the drain

The drain updates FT100, FT326, FT368, and FT369 to cite this prune.

### Flagged additions

- The class field on the registry type. It is the executable owner of the admission rule.
- The retirement of the step-scoped kind. This is a non-behavioral call, flagged for reviewer veto.
- The fold of the retained-workflow lists into the registry, under the one-source rule in AGENTS.md.
- The re-cut of the over-bound COMMAND needles. This is a non-behavioral call, flagged for reviewer veto.
- The Roles marker single-sourced in the shared-rule check, so 8 core rows stay in the registry. This is a non-behavioral call, flagged for reviewer veto.
- The removal of the `step` column from the `bench anchors` table with the step kind. This is a non-behavioral call, flagged for reviewer veto.

### Source-sentence-to-row table

| source decision | rows |
| --- | --- |
| Forbid anchors stay | AP02, AP26, AP47 |
| Label anchors stay | AP03, AP12 |
| Route anchors stay | AP04, AP13 |
| BENCH.md core stays | AP05, AP14, AP15 |
| Three invariant anchors | AP06, AP07, AP08 |
| Workflow rows go | AP09 |
| Command rows stay in this cut | AP10, AP16, AP46 |
| Rule wording is review-owned | AP01, AP11, AP17, AP18, AP28, AP32, AP39 to AP45 |
| One needle per family | AP20, AP21, AP22, AP23 |

### Pre-review proof checklist

- `Cited symbols`: `Entries`, `EvaluateGroup`, `FilesBelow`, `checkRetainedWorkflow`, `checkSharedRuleSingleSource`, `prose.MaxSentenceWords`, and the named tests resolve at the spec base.
- `Import edges`: the admission test in `internal/anchors` imports `internal/prose`. `go list -deps ./internal/prose` names no `internal/anchors` edge, so no cycle forms.
- `Source-row clauses and occurrences`: the decision bullets above quote each clause.
- `Promised field labels`: none.
- `Changed-function callers`: the one production `Entries` caller is `internal/gate/lane_select.go`. The other callers are test helpers in `cmd/bench`, `internal/conformance`, and `internal/anchors`. The slicing runs `bench consumers` for each.
- `Copy survival`: AP25 and AP31 fail if a copy of a registry diagnostic list survives.
- `Rendered-shape readers`: none.
- `Pin operators`: the per-topic tests use exact diagnostic equality. The help test uses exact line placement.
- `Entry reads`: none.
- `Derived expectations`: the admission test derives its word bound from `internal/prose`. The retained-workflow test derives its count from `Entries()`.
- `Consolidated rules`: the retained-workflow fold and the marker collapse, rows AP25, AP26, AP30, and AP31.
- `Quantified obligations`: "every Forbid row stays" (AP02), "each family keeps a fixture" (AP21), and "no Require row is unclassified" (AP01).
- `Workflow-step writes`: none.

## Planning claim inventory

This inventory declares the existing test seams that the coverage map cites. The spec predates the inventory, so it declares no premise, change, or decision item.

```bench-planning-facts
{"version":1,"authority":"authored-spec","completeness":"declared-complete","baseline":"source-base","scope":[],"changes":[],"premises":[],"seams":[{"id":"seam-1","location":{"path":"specs/anchor-admission-prune/spec.md","line":232,"needle":"| AP24 |","source":"spec"},"path":"cmd/bench/anchor_help_test.go","package":"github.com/gibbonmi/bench/cmd/bench","function":"TestAnchorsReportsAbsentNeedles","state":"existing","rows":["AP24"]},{"id":"seam-2","location":{"path":"specs/anchor-admission-prune/spec.md","line":243,"needle":"| AP35 |","source":"spec"},"path":"cmd/bench/anchor_help_test.go","package":"github.com/gibbonmi/bench/cmd/bench","function":"TestAnchorsReportsCaseFoldedEmphasisViolation","state":"existing","rows":["AP35"]},{"id":"seam-3","location":{"path":"specs/anchor-admission-prune/spec.md","line":242,"needle":"| AP34 |","source":"spec"},"path":"internal/anchors/anchor_harness_diagnostics_test.go","package":"github.com/gibbonmi/bench/internal/anchors","function":"TestEvaluatePathAnchorKinds","state":"existing","rows":["AP34"]},{"id":"seam-4","location":{"path":"specs/anchor-admission-prune/spec.md","line":228,"needle":"| AP20 |","source":"spec"},"path":"internal/conformance/fixture_bite_test.go","package":"github.com/gibbonmi/bench/internal/conformance","function":"TestEveryRetainedFixtureBitesThroughRegisteredOwner","state":"existing","rows":["AP20","AP22"]},{"id":"seam-5","location":{"path":"specs/anchor-admission-prune/spec.md","line":231,"needle":"| AP23 |","source":"spec"},"path":"internal/conformance/fixture_bite_test.go","package":"github.com/gibbonmi/bench/internal/conformance","function":"TestSpecTicketHandoffWorkflowFixturesAreComplete","state":"existing","rows":["AP23"]},{"id":"seam-6","location":{"path":"specs/anchor-admission-prune/spec.md","line":231,"needle":"| AP23 |","source":"spec"},"path":"internal/conformance/fixture_bite_test.go","package":"github.com/gibbonmi/bench/internal/conformance","function":"TestWorkflowCadenceAnchorsRejectDeletionAndSwap","state":"existing","rows":["AP23"]},{"id":"seam-7","location":{"path":"specs/anchor-admission-prune/spec.md","line":233,"needle":"| AP25 |","source":"spec"},"path":"internal/conformance/retained_workflow_test.go","package":"github.com/gibbonmi/bench/internal/conformance","function":"TestRetainedWorkflowReadsRegistry","state":"existing","rows":["AP25","AP31"]},{"id":"seam-8","location":{"path":"specs/anchor-admission-prune/spec.md","line":214,"needle":"| AP06 |","source":"spec"},"path":"internal/anchors/registry_data_test.go","package":"github.com/gibbonmi/bench/internal/anchors","function":"TestInvariantClauseAnchorsRedOnRemoval","state":"existing","rows":["AP06","AP07","AP08"]}],"decisions":[],"manual_claims":[]}
```
