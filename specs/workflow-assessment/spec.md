# Measure complete workflow cost and quality

Status: staged
Decision source: `docs/adr/0021-benchmark-workflow-orchestration.md`
Verification log: 2 iteration(s) to accept — Sol/high verified the design repairs. The author folded its final independent fixture cases for A21, A33, and A38.

## Problem

The benchmark comparison changes when cached input, repairs, failed attempts, and review costs are counted. Bench and the harness also own different measurements. A local record must preserve those distinctions before it supports a workflow decision.

## Solution

Keep one machine-readable local record per run, with chunk and role breakdowns and links to native evidence. Collect during ordinary work. Compare runs only with explicit quality and provenance limits, and require an approved plan before paid comparison trials.

Build dependency: none for record and ordinary collection; completion-evidence for importing its new result format.

## User stories

Line: gpt-5.6-sol / high. One retained implementation session. Starting effort is a recommendation for approval.

Sol/high is recommended because the data owner has bounded inputs and synthetic accounting fixtures. Chunk 2 is harder because it correlates measurements from several producers.

1. As an operator, I want to retain local run records, so that I can assess and continue approved work accurately.
2. As an operator, I want to keep failed work in the total, so that I can assess and continue approved work accurately.
3. As an operator, I want to compare chunk and role costs, so that I can assess and continue approved work accurately.
4. As an operator, I want to count cache usage correctly, so that I can assess and continue approved work accurately.
5. As an operator, I want to avoid repeated usage counts, so that I can assess and continue approved work accurately.
6. As an operator, I want to handle counter resets, so that I can assess and continue approved work accurately.
7. As an operator, I want to keep unknown usage honest, so that I can assess and continue approved work accurately.
8. As an operator, I want to estimate cost with the token cache, so that I can assess and continue approved work accurately.
9. As an operator, I want to separate actual charges, so that I can assess and continue approved work accurately.
10. As an operator, I want to account for incomplete prices, so that I can assess and continue approved work accurately.
11. As an operator, I want to measure concurrent wall time, so that I can assess and continue approved work accurately.
12. As an operator, I want to preserve measurement provenance, so that I can assess and continue approved work accurately.
13. As an operator, I want to collect ordinary work, so that I can assess and continue approved work accurately.
14. As an operator, I want to retain known history on updates, so that I can assess and continue approved work accurately.
15. As an operator, I want to inspect records efficiently, so that I can assess and continue approved work accurately.
16. As an operator, I want to read imports safely, so that I can assess and continue approved work accurately.
17. As an operator, I want to preserve records on storage failure, so that I can assess and continue approved work accurately.
18. As an operator, I want to plan comparable trials, so that I can assess and continue approved work accurately.
19. As an operator, I want to avoid adopting from one familiar task, so that I can assess and continue approved work accurately.
20. As an operator, I want to see variation and quality, so that I can assess and continue approved work accurately.
21. As an operator, I want to limit causal claims, so that I can assess and continue approved work accurately.
22. As an operator, I want to keep adoption under user control, so that I can assess and continue approved work accurately.
23. As an operator, I want to detect incomplete native evidence, so that I can assess and continue approved work accurately.

24. As an operator, I want to include all accepted-change work, so that the required evidence remains complete.
25. As an operator, I want to append run evidence, so that the required evidence remains complete.
26. As an operator, I want to fill unavailable measurements, so that the required evidence remains complete.
27. As an operator, I want to import idempotently, so that the required evidence remains complete.
28. As an operator, I want to retain known attempts, so that the required evidence remains complete.
29. As an operator, I want to protect known measurements, so that the required evidence remains complete.
30. As an operator, I want to collect valid Bench timing, so that the required evidence remains complete.
31. As an operator, I want to collect valid Bench census, so that the required evidence remains complete.
32. As an operator, I want to collect valid native usage, so that the required evidence remains complete.
33. As an operator, I want to reject ambiguous evidence joins, so that the required evidence remains complete.
34. As an operator, I want to require held-out comparison tasks, so that the required evidence remains complete.
35. As an operator, I want to require repeated eligible work, so that the required evidence remains complete.
36. As an operator, I want to compare equivalent assurance, so that the required evidence remains complete.
37. As an operator, I want to preserve quality, so that the required evidence remains complete.
38. As an operator, I want to hold comparison conditions fixed, so that the required evidence remains complete.
39. As an operator, I want to isolate a causal capability, so that the required evidence remains complete.

## Implementation decisions

This spec delivers the local assessment-record portion of FT231. FT231 keeps its broader experiment and causal-measurement scope and is not retired by this build. The assessment is advisory. It cannot change a model default or replace the gate.

Add the narrow `bench assessment` surface: `list`, `record --input <file>` imports a normalized run, `show <run-id>` renders detail, and `compare --plan <file> --runs <id,...>` renders a comparison. Use `list` to list local runs. Bare invocation returns usage at exit two. Approve only list, show, and compare as AXI children.

Record remains an operational mutation. All query views follow AXI with typed empty results, compact TOON, structured errors, help spellings, and complete detail on demand. Record import is non-interactive and returns the run ID and stored path. No command launches a model.

Trial execution remains a separately approved action.

Store versioned JSON at `$BENCH_HOME/assessment/<repo-key>/<run-id>.json`, beside the census and OTEL stores and outside the disposable worktree pool. Resolve repo identity through the existing pool-key owner. Records survive worktree release and reclaim. They remain until explicit cleanup. Do not add automatic expiry. Initial explicit cleanup uses a user-directed file removal after naming exact records; this spec adds no cleanup verb.

Run fields are `version`, `run_id`, `repo_key`, `source`, `condition`, `task_id`, `holdout`, `started_at`, `ended_at`, `time_reference`, `state`, `attempts`, `evidence`, `quality`, `bench_inputs`, `harness_inputs`, `diagnostics`, and `trial`. An attempt names `attempt_id`, `chunk_id`, `role`, `session_id`, `model`, `effort`, `state`, `started_at`, `ended_at`, `time_reference`, `usage`, `cost`, `measures`, `intervals`, and native evidence references. Roles distinguish implementation, repair, verification, review, and diagnostic consultation. Failed and cancelled attempts remain. An imported update may append or fill previously unknown evidence, but cannot silently discard an existing attempt or change a known identity. Identical import is idempotent; conflicting content is a refusal.

Usage fields are `input_uncached`, `input_cached`, and `output`. A source may also supply `input_total` with its declared semantics. If total includes cached input, derive uncached as total minus cached. Never add total and cached together.

Keep absent measures unknown. Zero means the source measured zero. Record the native counter name and whether it is a delta or a cumulative snapshot. Deduplicate on native event identity.

Cumulative counters use monotonic session epochs; a reset starts a new epoch. Ambiguous counters remain unknown with a reason.

Cost fields distinguish `estimated` and `actual`. Estimated token cost is the sum of uncached input, cached input, and output multiplied by their separate rates. Store currency, unit scale, rate source, rate date, and pricing conditions with the estimate. Include applicable tool and other charges when known.

If any component is unknown, label the total partial and report known components separately. Do not assume an API estimate equals a subscription charge. Actual charges remain unknown unless an authoritative charge record is supplied. This spec pins no current price.

Bench supplies elapsed spans, iterations it observes, raw-command census, and diff metadata through existing owners. Harness input supplies tokens, tools, read paths, and turns when available. The normalized import includes producer provenance and a native evidence reference per measure. Reuse `internal/otelrecord` and `internal/census` for their data; extend their read projections only where needed.

Do not infer missing spans from zero rows, and report malformed or incomplete input coverage. Wall time uses run endpoints or the union of observed intervals. Sum disjoint effort durations only under a distinct label. Concurrent reviews do not multiply wall time.

The phase instructions request one ordinary-work record update at chunk review and final close. Missing harness metrics do not block implementation or manufacture token counts. Unsupported harnesses can submit the normalized format with unknown fields. Provide documented field mappings and small fixtures for the native formats actually consumed; refuse unrecognized counter semantics. A completion record can be an evidence reference once that spec lands. It is not copied into a competing assurance verdict.

A comparison plan pins tasks, held-out designation, repetitions, conditions, source and harness versions, model and effort, quality tolerances, and a budget. It includes the user's approval reference before it is used to run paid trials. The comparison command validates the plan and compares declared runs; it does not grant or authenticate user approval. Later default changes require repeated comparable held-out tasks outside the FT311 authoring case and equivalent assurance.

Report failures, repairs, verification, reviews, unknown measures, variation, and quality outcomes as well as time and cost. A pilot remains a pilot. For a causal claim about the kit, require the FT231 three conditions: no Bench, current Bench, and one changed capability. For a narrower descriptive comparison, state its limits.

The user decides adoption.

Use bounded regular-file reads, schema validation, atomic replacement, no symlink traversal, and path-safe IDs. Do not execute commands or fetch URLs found in imported records. Store only the metric and evidence data needed for assessment. Prefer native references and narrow excerpts to full session transcripts.

## Implementation chunks

Implement these tickets in the retained session. Each ticket is one initial review chunk. Commit the ticket on its lane pass, run the three delegated axes, and repair findings before the next chunk. Review line: decision #13, overridden by the user for this implementation as gpt-5.6-sol / high on each axis. Final acceptance reconciliation stays with the author.

| chunk / ticket | blocked by | delivered outcome | harder chunk |
| --- | --- | --- | --- |
| 1.md — Store and inspect complete run records | none | Deliver normalized import, durable per-run storage, and list/detail output together | no |
| 2.md — Collect measurements from ordinary work | 1.md | Join existing Bench spans and census with explicitly supplied native harness measures | yes |
| 3.md — Compare runs under a pinned plan | 2.md | Deliver plan validation and descriptive comparison with complete cost, quality, and variation | no |

## Testing decisions

Use the named existing owner and new tests below. A new test name is a planned seam, not a claim that the test exists. Read its nearest fixture before implementation. Demonstrate each required omission or behavioral mutation as a diagnostic red, restore it, and show green. A compile failure is not the required red. Semantic review judges prose quality and the sufficiency of the test.

### Seam diagram

```text
Bench spans + census / native measures -> normalized run -> local store
                                                               |
                         approved plan + selected runs -> comparison report
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| A1 | 1 | A record survives release of its worktree | New test: TestAssessmentRecordStorage in internal/assessment/record_test.go | Release a disposable fixture assignment and read its assessment record. |
| A2 | 2 | Failed and cancelled attempts remain in the run | New test: TestAssessmentRecordPreservesHistory in internal/assessment/record_test.go | Import a success after failure and verify the failure remains. |
| A3 | 3 | A report separates chunk and performer roles | New test: TestAssessmentCommandDetail in internal/assessment/command_test.go | Mix implementation repair review verification and advice in one fixture. |
| A4 | 4 | Inclusive input totals exclude cached tokens from uncached input | New test: TestAssessmentRecord in internal/assessment/record_test.go | Use total 100 cached 80 and require uncached 20. |
| A5 | 5 | Repeated native events do not duplicate usage | New test: TestAssessmentRecordDuplicateEvents in internal/assessment/record_test.go | Import the same event twice. |
| A6 | 6 | Usage sums both epochs across a cumulative counter reset | New test: TestAssessmentRecordEpochs in internal/assessment/record_test.go | Require pre-reset and post-reset usage; reject epoch regression and isolate ambiguity by attempt. TestAssessmentRecordReviewCounters covers cumulative ambiguity; TestAssessmentRecordDeltaEpochs preserves independent delta quantities. |
| A7 | 7 | Absent usage is reported as unknown | New test: TestAssessmentRecordUnknownAndCharges in internal/assessment/record_test.go | Compare absent with measured zero; TestAssessmentRecordReviewCounters also covers a missing field between cumulative snapshots. |
| A8 | 8 | The estimate charges uncached input cached input and output at their own rates | New test: TestAssessmentRecordPrices in internal/assessment/record_test.go | Use three unequal nonzero rates and quantities so omitting any category fails. |
| A9 | 9 | An estimated cost cannot populate actual charges | New test: TestAssessmentRecordUnknownAndCharges in internal/assessment/record_test.go | Import token counts and rates without a billing record. |
| A10 | 10 | A missing charge component makes total cost partial | New test: TestAssessmentRecordUnknownAndCharges in internal/assessment/record_test.go | Omit tool pricing from a run with tool charges. |
| A11 | 11 | Concurrent attempts do not multiply elapsed wall time | New test: TestAssessmentRecordWall in internal/assessment/record_test.go | Overlap two reviews inside one run interval. |
| A12 | 12 | Each measure names its producer and native evidence | New test: TestAssessmentRecordProvenance in internal/assessment/record_test.go | Refuse unattributed tokens; TestAssessmentRecordMeasureProvenance also rejects unattributed quality, timestamps, and applicable charges. |
| A13 | 13 | Phase guidance requests records without requiring paid trials | review-owned: Spec checks phase instructions | Reject collection guidance that launches an experiment or blocks builds for missing optional metrics. |
| A14 | 14 | A conflicting record update preserves prior data | New test: TestAssessmentRecordUpdates in internal/assessment/record_test.go | Change a known session identity and verify refusal before replacement. |
| A15 | 15 | The query renders typed empty and complete detail states | New test: TestAssessmentCommand in internal/assessment/command_test.go | Exercise empty store list detail errors and each help spelling. |
| A16 | 16 | Invalid or unsafe imports produce no stored run | New test: TestAssessmentCommandUnsafe in internal/assessment/command_test.go | Exercise traversal, FIFO, symlinks, oversized input, version, duplicate IDs and keys, and ESC/BEL in task and chunk IDs. |
| A17 | 17 | An interrupted update leaves the previous record readable | New test: TestAssessmentRecordUpdates in internal/assessment/record_test.go | Inject write and rename failures with an existing record. |
| A18 | 18 | A paid-trial plan names approval and pinned conditions | New test: TestAssessmentComparison in internal/assessment/comparison_test.go | Omit budget approval or quality tolerance from the plan. |
| A19 | 19 | A pilot cannot establish default-change evidence | New test: TestAssessmentComparison in internal/assessment/comparison_test.go | Supply only the FT311 authoring case or one repetition. |
| A20 | 20 | Comparison reports variation and failed quality outcomes | New test: TestAssessmentComparison in internal/assessment/comparison_test.go | Use mixed successful failed and incomplete runs. |
| A21 | 21 | A kit-causal comparison requires the three FT231 conditions | New test: TestAssessmentComparison in internal/assessment/comparison_test.go | Omit each required arm independently. Also add an unknown fourth arm. Require refusal for every case. |
| A22 | 22 | Assessment never changes executable model defaults | New test: TestAssessmentCommand in internal/assessment/command_test.go | Run each assessment command and compare binding bytes before and after. |
| A23 | 23 | Malformed or unfinished native input remains incomplete | New test: TestAssessmentCollectionNativeGaps and TestAssessmentCollectionSpanCoverage in internal/assessment/collection_test.go | Mix one valid span with a truncated span and omit an expected session log. |

| A24 | 24 | Total cost includes every role and terminal attempt state | New test: TestAssessmentRecordUnknownAndCharges in internal/assessment/record_test.go | Use distinct costs for implementation failures repairs reviews verification and advice. |
| A25 | 25 | An update can append an attempt without deleting earlier attempts | New test: TestAssessmentRecordUpdates in internal/assessment/record_test.go | Append a successful attempt to a failed run. |
| A26 | 26 | An update can fill an unknown measure from new native evidence | New test: TestAssessmentRecordUpdates in internal/assessment/record_test.go | Supply a previously absent output count with provenance. |
| A27 | 27 | An identical import leaves the stored record unchanged | New test: TestAssessmentRecordUpdates in internal/assessment/record_test.go | Import the same normalized run twice; compare bytes and require no replacement write. |
| A28 | 28 | An update that deletes an existing attempt is refused | New test: TestAssessmentRecordPreservesHistory in internal/assessment/record_test.go | Omit an earlier failed attempt from a replacement. |
| A29 | 29 | Conflicting known evidence is refused before replacement | New test: TestAssessmentRecordUpdates in internal/assessment/record_test.go | Change a measured token count while retaining its native event ID. |
| A30 | 30 | Selected OTEL spans reach the stored run with provenance | New test: TestAssessmentCollection in internal/assessment/collection_test.go | Join two known trace IDs and verify their derived elapsed value and references. |
| A31 | 31 | Selected census events reach the stored run with provenance | New test: TestAssessmentCollection in internal/assessment/collection_test.go | Join known raw-command events from the expected assignment. |
| A32 | 32 | Mapped harness counters reach the stored attempt with provenance | New test: TestAssessmentCollection in internal/assessment/collection_test.go | Import a valid session epoch and cache-bearing usage event. |
| A33 | 33 | A foreign or ambiguous assignment join is refused | New test: TestAssessmentCollection in internal/assessment/collection_test.go | Run separate cases for a foreign assignment and two matching events with no unique correlation. |
| A34 | 34 | Non-held-out tasks cannot establish default-change evidence | New test: TestAssessmentComparison in internal/assessment/comparison_test.go | Repeat the same familiar task while marking it ineligible. |
| A35 | 35 | Missing planned repetitions make default-change evidence incomplete | New test: TestAssessmentComparison in internal/assessment/comparison_test.go | Omit one required repetition from an eligible task. |
| A36 | 36 | Unequal acceptance or review obligations make adoption evidence ineligible | New test: TestAssessmentComparison in internal/assessment/comparison_test.go | Remove an axis from one condition. |
| A37 | 37 | Quality outside the approved tolerance makes adoption evidence ineligible | New test: TestAssessmentComparison in internal/assessment/comparison_test.go | Use a cheaper condition that exceeds the failure tolerance. |
| A38 | 38 | An undeclared condition difference makes the comparison ineligible | New test: TestAssessmentComparison in internal/assessment/comparison_test.go | Change each fixed field alone: task, revision, model, effort, harness, and execution limits. Require refusal for every case. |
| A39 | 39 | A causal arm changing two capabilities is ineligible | New test: TestAssessmentComparison in internal/assessment/comparison_test.go | Keep all three arms but change two capabilities in the experimental arm. |

### Edge inventory

The storage precedent is `internal/census` and `internal/otelrecord`, both outside the worktree pool. OTEL's current reader returns completed spans and skips malformed lines; assessment must expose incomplete coverage instead of turning that into a complete total. Its existing readers remain compatible. The public command registry owns help and AXI membership; update the skill table and conformance inventories together. A new package joins the existing package and injected-port registries.

Fixtures include absent and empty stores, unknown versus zero, and inclusive and exclusive cache semantics. They also cover duplicate events, counter resets, unfinished attempts, parallel reviews, failed imports, and loss of a native log. Absent and empty stores yield the same typed empty list. A missing requested run is a refusal. Tests use fixed synthetic prices and native log fragments with declared semantics. No live account, network, or paid model call is a gate fixture.

Scope cut: a broad experimental runner could ship separately after this record and comparison surface. This spec does not charge or execute that work. The record and import API can support it later without changing the adoption rule.

## Ownership fences

- `tests/canary/workflow-guidance-anchors/review-clean-terminal-result`

These paths are the union of ticket expectations. A directory entry is an exact prefix for that existing owner or fixture family. Expansion follows decision #5, with the plan updated before use. It cannot weaken existing guarantees.

- `internal/assessment` (new)
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/main.go`
- `.agents/skills/bench-craft-cli/SKILL.md`
- `.bench/BENCH-reference.md`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/registry/packages.go`
- `internal/conformance/injected_ports_registry_test.go`
- `internal/otelrecord/reader.go`
- `internal/otelrecord/reader_test.go`
- `internal/census/census.go`
- `internal/census/events.go`
- `internal/census/census_test.go`
- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/commands/bench-final-check.md`
- `roadmap/FT231.md`
- `specs/workflow-assessment/spec.md`
- `specs/workflow-assessment/tickets`
- `reviews/workflow-assessment.md` (new)

- `internal/conformance/subcommand_routing_table_test.go`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/docs-currency-token-diet/introduces-undeclared-command`
- `tests/canary/docs-currency-token-diet/stale-command-reference`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/dangling-index`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/skills-index-command-adapters/missing-index-field`
- `tests/canary/skills-index-command-adapters/stale-index-wording`
- `tests/canary/skills-index-command-adapters/unindexed-skill`
- `tests/canary/workflow-guidance-anchors/coverage-axis-anchor`
- `tests/canary/workflow-guidance-anchors/final-check-bare-leftover-clean-retired`
- `tests/canary/workflow-guidance-anchors/final-check-census-read-before-land`
- `tests/canary/workflow-guidance-anchors/final-check-landed-worktree-sweep`
- `tests/canary/workflow-guidance-anchors/final-check-light-path-changelog-heading`
- `tests/canary/workflow-guidance-anchors/final-check-scratch-branch-clean`
- `tests/canary/workflow-guidance-anchors/implement-spec-coverage-task-seeding`
- `tests/canary/workflow-guidance-anchors/implement-spec-entry-validation`
- `tests/canary/workflow-guidance-anchors/implement-spec-incapable-harness`
- `tests/canary/workflow-guidance-anchors/implement-spec-inline-exception`
- `tests/canary/workflow-guidance-anchors/implement-spec-mandatory-delegation-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-offer-retired`
- `tests/canary/workflow-guidance-anchors/implement-spec-offer-scope`
- `tests/canary/workflow-guidance-anchors/implement-spec-read-only-helper`
- `tests/canary/workflow-guidance-anchors/implement-spec-status-flip-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-worktree-before-preflight`
- `tests/canary/workflow-guidance-anchors/implement-spec-write-delegation`
- `tests/canary/workflow-guidance-anchors/line-anchor-missing`
- `tests/canary/workflow-guidance-anchors/prepared-build-approval`
- `tests/canary/workflow-guidance-anchors/prepared-build-freshness`
- `tests/canary/workflow-guidance-anchors/prepared-review-axis-returns`
- `tests/canary/workflow-guidance-anchors/prepared-review-blast-evidence`
- `tests/canary/workflow-guidance-anchors/prepared-review-capable-handoff`
- `tests/canary/workflow-guidance-anchors/prepared-review-inline-axis-route`
- `tests/canary/workflow-guidance-anchors/prepared-review-legacy-entry-points`
- `tests/canary/workflow-guidance-anchors/prepared-review-native-dispatch`
- `tests/canary/workflow-guidance-anchors/prepared-review-runtime-capability`
- `tests/canary/workflow-guidance-anchors/prepared-review-shared-evidence`
- `tests/canary/workflow-guidance-anchors/reference-agent-push-rule`
- `tests/canary/workflow-guidance-anchors/reference-bench-operational-layer`
- `tests/canary/workflow-guidance-anchors/reference-category-context`
- `tests/canary/workflow-guidance-anchors/reference-category-oracle`
- `tests/canary/workflow-guidance-anchors/reference-category-setup`
- `tests/canary/workflow-guidance-anchors/reference-category-work`
- `tests/canary/workflow-guidance-anchors/reference-gate-authority`
- `tests/canary/workflow-guidance-anchors/reference-kit-only-ship`
- `tests/canary/workflow-guidance-anchors/reference-no-path-fallback`
- `tests/canary/workflow-guidance-anchors/reference-progressive-loading-term`
- `tests/canary/workflow-guidance-anchors/reference-refusal-route-shape`
- `tests/canary/workflow-guidance-anchors/reference-retro-capture-owner`
- `tests/canary/workflow-guidance-anchors/reference-retro-drain-owner`
- `tests/canary/workflow-guidance-anchors/reference-skills-guidance`
- `tests/canary/workflow-guidance-anchors/reference-upgrade-route`
- `tests/canary/workflow-guidance-anchors/review-base-merged-main-tip`
- `tests/canary/workflow-guidance-anchors/review-falsification-accept-routing`
- `tests/canary/workflow-guidance-anchors/review-falsification-dispositions`
- `tests/canary/workflow-guidance-anchors/review-kit-guidance-set`
- `tests/canary/workflow-guidance-anchors/review-persistence-anchor`
- `tests/canary/workflow-guidance-anchors/review-preflight-explicit-base`
- `tests/canary/workflow-guidance-anchors/review-repair-ticket-covers`
- `tests/canary/workflow-guidance-anchors/review-repair-ticket-owner`
- `tests/canary/workflow-guidance-anchors/review-standing-falsification`
- `tests/canary/workflow-guidance-anchors/review-universal-claim-bar`

- `projects/benchkit.md`
- `internal/tickets/registry_data.go`



- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-owner`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-routing`
- `tests/canary/workflow-guidance-anchors/benchkit-spec-ownership`
- `tests/canary/workflow-guidance-anchors/benchkit-system-suite-route`

- `internal/otelrecord/attributes.go`
- `internal/commit/commit.go`
- `internal/commit/assessment_span_test.go`
- `internal/worktree/land.go`
- `internal/worktree/land_trace_test.go`

## Out of scope

This spec excludes a paid benchmark launcher, the complete FT231 program, automatic default adoption, and automatic expiry. It also excludes speculative harness adapters and recovery of unavailable charges.

## Further notes

Approval is pending for this spec and its ticket graph. Authoring these documents does not authorize their implementation. The source decisions remain closed. The current kit instructions have not yet been changed by these specs.

Flagged additions: the engineering mechanisms named under Implementation decisions, the new test seams, and the ticket graph. They implement the source outcomes. The review must reject an additional behavior with no coverage row.

Pre-review proof checklist:

- Cited symbols: existing owner names were read in this session. New test names are explicitly planned and must be created by the ticket that owns their row.
- Import edges: no new forbidden-import claim. New data owners list their actual consumers in the seam diagram.
- Source-row clauses and occurrences: the table below maps the closed source clauses. The index and resolved tickets are the authoritative occurrences. Historical recommendation assets do not override them.
- Promised field labels: the exact record fields, command arguments, and plan fields appear under Implementation decisions where this spec adds them. Guidance-only specs add no machine record.
- Changed-function callers: the reader and enforcement inventory appears under Edge inventory. Recheck function-level callers before changing an existing signature. No signature change is authorized by assumption.
- Copy survival: W16 owns live workflow-copy removal. Other specs add one owner and require consumers to call it rather than parse the same facts independently. Historical evidence remains readable.

| source clause | acceptance rows |
| --- | --- |
| #11: local machine-readable per-run records with native evidence | A1, A2, A3, A12, A13, A14, A15, A16, A17, A23 |
| #9 and #11: full cost including cache and unknown actual charges | A4, A5, A6, A7, A8, A9, A10, A11 |
| #9 and #11: approved plans and repeated held-out comparisons | A18, A19, A20, A21, A22 |

The normalized import can select Bench evidence through `bench_inputs`. It contains `trace_ids`, `census_event_ids`, and the expected `assignment_id`. Each selected event maps to an `attempt_id`, `chunk_id`, and `role`. Resolve only those identifiers from the repository's existing OTEL and census stores.

Do not sweep unrelated sessions by timestamp. Reject a foreign assignment or an ambiguous join. A native harness measure carries its session, epoch, event ID, counter semantics, and mapped attempt explicitly.

The successful join writes the selected native references and derived measures into the same normalized record. Keep selection and missing-input diagnostics in that record. Expose this through the existing `record --input` form. An input without selectors remains a valid normalized import with unknown Bench measures. This keeps ticket 1 useful before native collection lands.

Comparison returns descriptive results and an evidence-eligibility disposition. The disposition is not adoption approval. Missing comparisons remain unknown. A default-change disposition requires eligible held-out tasks, the plan's repetitions, equivalent acceptance and review obligations, and quality within the approved tolerance. Match task, revision, model, effort, harness, and execution limits except the plan's single declared experimental variable. A causal kit plan has exactly the three named conditions, with one capability difference in its experimental arm.

Review repair coverage: A24, A25, A26, A27, A28, A29, A30, A31, A32, A33, A34, A35, A36, A37, A38, A39. These rows refine the source clauses already mapped above.

Plan-edit authority is limited to chunk boundaries, dependencies, row assignments, Writes expectations, and necessary ownership expansion under decision #5. Preserve the approved behavior and pass criteria. A material acceptance change still needs the user's decision. Bulk rewrites exclude this spec and its tickets. Completion evidence must review a changed plan digest before it can satisfy the next checkpoint.

The A33, A38, and A21 fixture tables must execute every named case independently. A combined mutation cannot replace an individual case. These case refinements close the second assessment review.

Authoring close: this spec and its tickets are staged for user sign-off. The reviews above assess the proposed build. No implementation or implementation test result is claimed.

Reviewer amendment on 2026-09-11: apply decision #13 to the declared implementation model. Sol implementations use Astra/high review axes. This amendment follows the spec-authoring reviews recorded above.

## Completion verification plan

The stable chunk IDs are `1`, `2`, and `3`, in the approved ticket order.
Each checkpoint requires current author tests and all three independent review results.
The user approves Astra implementation in the retained session and Sol/high review axes.
The synthetic cross-attempt counter case extends A6 and A24 without changing their pass criteria.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "1",
      "tickets": [
        "1.md",
        "4.md"
      ],
      "verification": [
        {
          "id": "assessment",
          "command": "bench test --package ./internal/assessment"
        },
        {
          "id": "dispatcher",
          "command": "bench test --package ./cmd/bench"
        },
        {
          "id": "cache-probe",
          "command": "bench test --package ./internal/assessment --run TestAssessmentRecord",
          "probe": "omit cached-input subtraction"
        }
      ]
    },
    {
      "id": "2",
      "tickets": [
        "2.md",
        "5.md"
      ],
      "verification": [
        {
          "id": "assessment",
          "command": "bench test --package ./internal/assessment"
        },
        {
          "id": "dispatcher",
          "command": "bench test --package ./cmd/bench"
        },
        {
          "id": "cache-probe",
          "command": "bench test --package ./internal/assessment --run TestAssessmentRecord",
          "probe": "omit cached-input subtraction"
        }
      ]
    },
    {
      "id": "3",
      "tickets": [
        "3.md"
      ],
      "verification": [
        {
          "id": "assessment",
          "command": "bench test --package ./internal/assessment"
        },
        {
          "id": "dispatcher",
          "command": "bench test --package ./cmd/bench"
        },
        {
          "id": "cache-probe",
          "command": "bench test --package ./internal/assessment --run TestAssessmentRecord",
          "probe": "omit cached-input subtraction"
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "assessment",
      "command": "bench test --package ./internal/assessment"
    },
    {
      "id": "dispatcher",
      "command": "bench test --package ./cmd/bench"
    },
    {
      "id": "system",
      "command": "bench test --check system"
    }
  ]
}
```

## Chunk 1 author evidence

The synthetic suite observed behavioral reds for storage, inclusive cache arithmetic, duplicate events, counter resets, rates, provenance, and concurrent wall time.
The cross-attempt cumulative test observed USD 25 before repair and USD 15 after repair.
The large-integer test rejected loss of precision after the JSON comparison repair.
Cache subtraction, history preservation, and token pricing mutations each produced a diagnostic red and restored the source.
The remaining update and unknown-value cases exercised behavior already present when their rows were added.

The assessment suite and full dispatcher suite passed. The chunk has no paid comparison trial or model-default change.

## Chunk 1 repair evidence

User approval permits the higher implementation model in this retained session. The three independent reviewers used Sol at high effort. The additional cross-harness review was waived explicitly by the user. No comparison trial was authorized or launched.

Repair ticket 4 belongs to chunk 1 and blocks chunk 2. It covers the eight accepted findings recorded in reviews/workflow-assessment.md. These changes refine tests and implementation within the approved acceptance criteria.

Synthetic tests reproduced eight failures before repair. Cases covered missing cumulative fields, epoch regression, unrelated-attempt ambiguity, duplicate JSON keys, and ESC/BEL in task and chunk identifiers. Five further red cases reproduced missing provenance for quality, run time, attempt time, unknown tool charges, and unknown actual charges. The full assessment suite then passed.

Named omission probes bit and restored production. A3 omitted a typed role cell; A7 changed absent input to zero; A9 marked absent billing complete. A14 bypassed preservation; A15 omitted the summary; A17 bypassed write and rename failure seams. A25 refused array growth; A26 refused omitted-field fills; A27 bypassed the idempotent early return. Each reported one failing test and successful restoration.

The first A26 probe was silent because its changed nil branch did not own omitted JSON fields. The first A3 probe was silent because raw-record text masked the typed-table omission. Two subsequent probes were invalid because the strengthened baseline ignored TOON numeric-string quoting. Corrected fixtures and targeted mutations produced the biting results above. These unsuccessful diagnostic attempts remain part of ordinary-work accounting.

Synthetic lifecycle dogfood recorded synthetic-release from a disposable main-based assignment. After bench worktree release removed that assignment, the integration assignment successfully queried the same record. Earlier release attempts against unmerged source-based fixtures refused and preserved their assignments. No model ran in these lifecycle fixtures.

The second repair pass closes P4/C3 and P5. TestAssessmentRecordDeltaEpochs first failed with an unknown total for independent quantities 10 and 5. The fix limits epoch-order validation to cumulative counters. The same test and the full assessment suite then passed. Both authoritative timestamp field lists now include time_reference.

The broader reviewer-merge gate found a guidance-table ordering regression. Restoring the original first row made TestAXIGuidanceContractBites pass without changing its assertion.

## Chunk 2 author evidence

Chunk 1 closed at 0d255141 after three clean Sol reaffirmations. Its checkpoint passed gofmt, vet, tests, race checks, and the system suite. Shellcheck was skipped; six capability-dependent cases were skipped. Ticket 4 is now checked complete.

Collection selectors and their explicit mappings are defined in internal/assessment/collection.go. Bench inputs select trace IDs and census event IDs under one expected assignment. Harness inputs select a bounded native fragment and explicitly declare its session, event, epoch, sequence, and supported counter semantics.

Attempt measures retain referenced numeric observations. Attempt intervals retain observed start/end pairs with native references. Run diagnostics retain missing, malformed, or unfinished input coverage. These fields implement the existing native-evidence and elapsed-union requirements.

Initial command fixtures failed because native selectors were unsupported. After collection was added, four further cases failed: start/end closure, disjoint spans, nested spans, and malformed coverage. The shared OTEL decoder and interval-union owner now cover those cases.

Named probes bit and restored production. Bypassing trace assignment checks failed both foreign and ambiguous cases. Dropping the second selected trace failed the two-trace union case. Dropping the second census event failed its selected-count case. Bypassing census collection failed both census refusal partitions.

Omitting native diagnostics failed four missing/malformed/unfinished cases. Bypassing harness counter semantics failed the unsupported-semantics case. Native inputs and expected arithmetic use synthetic fixtures. No comparison trial or model launch occurred.

The existing OTEL and census suites passed after their read projections were extended. Historical readers retain their existing behavior. Phase guidance requests ordinary record updates at review and final close, without requiring optional harness metrics or authorizing paid experiments.

The census event projection lives in internal/census/events.go. The commit lane caught growth beyond census.go's 400-line bound, so this approved ownership expansion preserves that existing file budget.

## Chunk 2 repair plan

Ticket 5 belongs to chunk 2 and blocks chunk 3. It closes all nine initial Sol findings without changing acceptance criteria.

A30 adds these publication details:

- Authentic commit and landing spans with separate assignment provenance
- Diff-path values
- Native references

A30-A32 each add positive second-attempt routing. A32 and A7 add delta and partial native counters. A16 and A23 add parent-symlink and present-empty native fragments. Unknown fields stay unknown, and partial native counters mark coverage incomplete.

The necessary ownership expansion adds the OTEL assignment attribute and its commit and landing producers with synthetic producer tests. The assignment identifier stays separate from the published Git subject. Historical subject-bearing traces without assignment evidence remain unjoinable.

Chunk 2 repair validation reproduced five failures: both publication joins and each missing native counter category. A present-empty fragment then reproduced a separate classification failure. All six pass after repair.

Named probes bit for these partitions and restored production:

- Second-attempt attribution for each producer
- Diff-path projection for each publication shape
- Missing-counter diagnostics for each token category
- Delta semantics
- Parent-symlink protection
- Landing assignment emission
- Commit assignment identity

The first commit-assignment omission could not compile because it left an unused variable. A compiling wrong-assignment mutation replaced that diagnostic attempt. The unsuccessful attempt remains in cost accounting.

The second repair pass extends the existing full landing test to assert resolved assignment provenance. It removes the duplicate landing fixture. A33 adds conflicting duplicate selectors whose two attempt mappings are both valid. The existing guarantees remain unchanged.

The full landing assignment-copy omission produced one diagnostic red and restored production. The valid duplicate-selector mutation produced two reds, one for census and one for traces, and restored production. An initial test-field typo prevented test execution twice; the corrected fixture uses creation.Assignment.ID. Those diagnostic attempts remain recorded.

## Chunk 3 implementation details

The optional trial object binds a run to its plan and repetition. It records the harness version, limits, capabilities, assurance obligations, and native reference. The run source and actual attempt lines remain the owners of revision, model, and effort.

Comparison validation checks each actual run against its condition. Missing planned repetitions remain incomplete evidence. Repetition checks use observed slots and counts, so large declared repetition counts cannot cause unbounded iteration.
