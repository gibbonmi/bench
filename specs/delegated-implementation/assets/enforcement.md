# Enforcement and reader inventory

Baseline: `3580a79f43e2938348975c83cbe00f8c0f881878`, current main at phase entry.
Assessment delivery: `15a48d99`.
Assessment retirement and closeout: `3580a79f`.
Subject: Proposed delegated mode only.

## Existing enforcement

| Read source | Observed contract | Disposition |
| --- | --- | --- |
| `internal/reviewrecord/record.go` | Record version 1 names one implementation session. CheckReviews excludes that session. | DI1–DI7 extend this owner |
| `internal/reviewrecord/parse.go` | Strict JSON, terminal metadata, native excerpts, and review supersession | DI3, DI11, DI12 preserve refusal semantics |
| `internal/reviewrecord/plan.go` | Source-bound spec and ticket bytes form the plan digest. Planned order respects ticket dependencies. | DI2 binds delegated identities here |
| `internal/reviewrecord/coverage.go` | Verification uses ImplementationSession. Chunks form an ancestral, continuous source chain. | DI4 and DI23 preserve the chain |
| `internal/reviewrecord/check.go` | Final reconciliation and verification require ImplementationSession. Commands and probes match exact obligations. | DI8–DI12 make ownership mode-specific |
| `internal/reviewrecord/files.go` | Record paths, bounded files, and source-tree reads retain one owner | Preserve for both versions |
| `internal/gate/checkpoint.go` | CheckTrees runs before oracle acceptance and separates checkpoint purposes | DI4–DI7 and DI12 exercise its real entry |
| `internal/gate/completion.go` | Prospective composition permits only the status transform and review record | DI10 retains this contract |
| `internal/preflight/review.go` | The completion projection derives source and plan digests from the prepared tip | Ticket 1 checks both versions through this consumer |
| `internal/worktree/merge.go` | Only clean owned sibling tips or default-branch commits enter the target lane | DI23 composes existing behavior |
| `internal/assessment/collection.go` | Explicit mappings name one assignment per BenchInputs value | DI14 and DI15 extend collection across batches |
| `internal/assessment/vocabulary.go` | Roles excludes orchestration | DI13 adds that role at its sole owner |
| `internal/assessment/validate.go` | Native event ownership is unique across attempts | DI15 preserves this rule across selectors |
| `internal/assessment/update.go` | Updates retain ordered arrays and known values | DI17 preserves failed and replaced work |
| `internal/assessment/record.go`, `cost.go`, `summary.go` | Native events normalize once. Estimates and actuals stay separate. Wall time uses interval union. | DI16 and DI18 reuse these calculations |
| `internal/assessment/command.go`, `store.go` | Record import collects before validation and atomic publication | DI14–DI19 use the command seam |
| `internal/assessment/comparison_roles.go`, `plan_validate.go` | Comparisons and line validation consume Roles | DI13 includes both consumers |

The native excerpt digest proves excerpt integrity, not native-session authenticity.
The gate cannot prove that an orchestrator disclosed every unrecorded invocation.
DI27 therefore stays review-owned and compares the disclosed account with native dispatch evidence.

## Fixture precedents and gate attachment

`internal/reviewrecord/recordtest/fixture.go` is a shared compiled fixture owner.
It constructs committed plans, chunks, verification, reviews, and final reconciliation.
It is not a private test helper imported across a package seam.
Extend it for delegated records and preserve its version 1 defaults.

`TestReviewRecordSource` in `internal/reviewrecord/source_test.go` proves ancestral continuity and plan identity.
`TestReviewRecordTerminal` proves native metadata, axis ownership, and supersession.

`TestReviewCheckpoint` in `internal/gate/review_checkpoint_test.go` refuses incomplete evidence before the oracle runs.
`TestReviewCheckpointReuse` distinguishes ordinary, chunk, and complete obligations.
`TestLandingCompletionEvidence` in `internal/landing/completion_evidence_test.go` proves final results and destination-delta refusal.

`TestMergePublishesTheMergeTreeAndMovesTheCheckout` in `internal/worktree/merge_test.go` supplies the composition precedent.
`internal/worktree/completion_fixture_test.go` shows the existing landing-evidence attachment.
The new journey extends worktree fixtures, using the shared recordtest owner.
It does not introduce another subprocess constructor or paid harness runner.

The assessment command, collection, history, and comparison tests provide canned native-event precedents.
Use fixed intervals and unequal synthetic rates, never current prices or elapsed time from separate live runs.
A two-assignment import must drive Collect through Command before publication.

The gate enters through `.bench/gate.sh` and its existing phase owner.
Ordinary Go tests execute the package fixtures.
The system phase remains available through `bench test --check system`.

The conformance registry binds workflow anchors to `docs-currency-workflow`.
`checkWorkflowAnchors` consumes the AfterImplementSpec anchor group.
`TestRetainedWorkflow` and its mutation helper prove the retained workflow's existing prose tripwires.
New delegated clauses use the same owner and demonstrate an omitted-clause red.
No prose tripwire claims to prove native dispatch behavior.

## Reader sweep

The sweep used `rg --hidden`, excluded `.git/`, and included scripts and dot-directories.
It searched implementation_session, ImplementationSession, reviewrecord, bench-review-record, and bench-completion-plan.
It also searched retained-author wording, assessment Roles callers, and native Bench-input selectors.

| Reader group | Named consumers | Disposition |
| --- | --- | --- |
| Evidence | reviewrecord parse, files, plan, coverage, check, CheckReviews, recordtest, source_test, record_test | Ticket 1 |
| Checkpoint | gate checkpoint, completion, review_checkpoint_test, completion_test | Ticket 1 adds delegated consumer fixtures |
| Landing | landing completion_evidence_test, worktree completion_fixture_test and land_fixtures_test | Tickets 1 and 3 add delegated journeys |
| System fixtures | owner_land_race_test and owner_landing_fixture_test | Keep version 1 fixture defaults, DI1 |
| Preparation | preflight review.go completionEvidenceTable and renderReviewPacket | Ticket 1 tests the unchanged projection |
| Assessment roles | Validate, ValidatePlan, comparisonRoles, command_test, record_test, comparison_test, comparison_report_test | Ticket 2, DI13 and DI19 |
| Assessment batches | Collect, collectSpans, collectCensus, mapped, selectMapping, Command, Store.Record, compatible | Ticket 2, DI14–DI17 and DI19 |
| Workflow | BENCH, BENCH-reference, implement, review, final-check, craft-line, craft-delegate, delegation-discipline, craft-tickets | Ticket 3 |
| Workflow advertisements | CONTEXT, benchkit profile, ADR 0021, field guide | Ticket 3 |
| Workflow pins | anchors registry_retained_workflow, registry_data, registry_ft311_preparation, registry_ft311_review_dispatch | Add the opt-in clauses without erasing default predicates |
| Pin checks | retained_workflow_test, implementation_continuation_test, docs_workflow_helpers_test, recurrence_maintenance_contract_test | Preserve default contracts, add opt-in tripwires |
| Workflow fixtures | tests/canary/workflow-guidance-anchors, including retained-integration-source and delegate-resume cases | Ticket 3 updates only affected fixtures |
| Historical material | CHANGELOG, docs/research/ft311-session-orchestration.md, docs/greenfield-build-sequence.md | Add a changelog entry only; historical outcomes remain true |
| Other staged specs | session-context-efficiency, overflow, queries, cleanup, measurement | Excluded; their approved default author contracts remain valid |
| Routing hooks and adapters | .bench/lines.env, check-agent-line hook, existing adapters | Excluded; configured membership and model defaults do not change |

The source digest excludes only the current spec's review pickup.
The identity declaration therefore belongs in the plan, not solely inside that pickup.
Plan changes from author replacement require existing amendment evidence and current verification.
No new artifact path becomes an unreviewed source exclusion.

## Registry closure and fences

The spec's ownership fence is the union of all ticket Writes fields.
`internal/tickets/registry_data.go` binds assessment, anchors, and worktree paths to command registry files.
Tickets 2 and 3 co-name those files even though no CLI grammar change is intended.
The command registry, help inventory, AXI inventory, and subcommand-routing table were inspected for these consumers.
They continue to advertise the same verbs.
No new check registry or process seam is planned.

## Source-to-row mapping

| Confirmed source clause | Rows |
| --- | --- |
| No-flag runs retain the single author | DI1, DI19, DI20 |
| The invoking session selects all configured tiers and effort | DI21 |
| Each ticket delegate retains code, tests, probes, and repairs | DI4, DI11, DI25, DI37, DI38, DI40, DI41 |
| Independent chunks may author concurrently within a declared limit | DI22, DI23 |
| Dependencies and integrated acceptance remain ordered | DI23, DI24, DI42, DI43 |
| Three distinct mid-tier axes exclude all authors and the orchestrator | DI5, DI6, DI7, DI12, DI39 |
| The orchestrator verifies, reconciles, integrates, and lands | DI8, DI9, DI10, DI30 |
| Replacement follows agreed triggers after the old writer stops | DI11, DI26, DI31–DI36 |
| Every participant, failed attempt, repair, review, verification, and diagnostic remains accounted for | DI13, DI14, DI15, DI17, DI27 |
| Estimates, actual charges, and unknowns remain separate | DI16, DI18 |
| Identities and evidence survive interruption and replacement | DI2, DI3, DI11, DI28 |
| Synthetic fixtures only, no paid comparisons or default changes | DI29 |

## Pre-review proof checklist

Cited symbols: Existing symbols resolve in the sources above. Planned test names are explicitly new.

Import edges: Shared recordtest is compiled fixture support. No private test helper crosses a package seam.

Source-row clauses and occurrences: The table above covers the confirmed conversation. Default wording readers appear in the reader sweep.

Promised field labels: execution, mode, run_id, orchestrator_session, author_limit, ticket, bench_input_batches, integration-verification.

Changed-function callers: CheckTrees is called by Check and the gate checkpoint. CheckReviews is called by source coverage and record tests.

Changed-function callers: Collect is called by assessment Command and collection tests. Roles feeds validation, comparison, and their fixture families.

Copy survival: none. Existing version 1 forms remain supported without a second identity derivation.

## Flagged engineering additions

The version 2 plan is the proposed source-bound identity seam.
The optional assessment batch field closes the existing single-assignment collection limit.
Neither addition expands the approved product behavior.
Future implementations must not infer transcript authenticity from a self-declared identity.

## Proposed glossary additions

Orchestrator: The invoking session that assigns ticket authors, integrates their source, reconciles acceptance, and lands the approved run.
Avoid: implementation_session, ticket author.

Ticket author: The recorded session responsible for one ticket's production changes, tests, probes, repairs, and author verification.
Avoid: reviewer, orchestrator.

Branch-local evidence: Results bound to an author branch before the contribution receives integrated acceptance.
Avoid: accepted checkpoint, final verification.

These definitions enter CONTEXT only during implementation.

## Concurrent main update

Main advanced to `8aac4f9d` during authoring.
The added preflight plan check uses ReadPlan at the prepared tip.
Checkpoint refusals now name the preflight recovery route and the missing fence.
The author inspected those changes and preserves them.
They do not change the proposed identity or source-chain contract.
Ticket 1 must exercise `internal/preflight/plan.go` through its build and review completion-plan checks as well.

## Reviewer correction

Each ticket receives its own retained delegate.
The chunk remains the review scope.
Delegated review always uses the configured mid tier at high effort.
The source-bound requirement now names its owning ticket because one chunk can have several authors.
Same-chunk ticket dependencies require green commits, while prerequisite chunks require accepted checkpoints.
