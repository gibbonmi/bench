# Import complete native probe results

Blocked by: 1-record-native-assignments.md, 2-enforce-native-record-order.md
Writes: .agents/commands/bench-implement-spec.md, .agents/commands/bench-review-implementation.md, .bench/BENCH-reference.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/record_operations_test.go (new), cmd/bench/response_bound_test.go, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/probe/baseline.go, internal/probe/baseline_test.go, internal/probe/omit_file_test.go, internal/probe/outcome_test.go, internal/probe/probe.go, internal/probe/probe_test.go, internal/probe/probetest/fixture.go (new), internal/probe/refusal_test.go, internal/probe/result.go (new), internal/probe/result_test.go (new), internal/responsebound/lines_test.go, internal/responsebound/owner.go, internal/responsebound/owner_test.go, internal/responsebound/responseboundtest/spill.go, internal/reviewrecord/recordcmd/command.go, internal/reviewrecord/recordcmd/refusal_test.go, internal/reviewrecord/recordcmd/verification.go (new), internal/reviewrecord/recordcmd/verification_import_test.go (new), internal/reviewrecord/recordcmd/verification_test.go, internal/reviewrecord/recordtest/fixture.go, internal/treetarget/identify.go, internal/treetarget/identify_test.go, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/dangling-index, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/skills-index-command-adapters/missing-index-field, tests/canary/skills-index-command-adapters/stale-index-wording, tests/canary/skills-index-command-adapters/unindexed-skill, tests/canary/workflow-guidance-anchors/agents-handoff-section-rule, tests/canary/workflow-guidance-anchors/calibration-pickup-confidence, tests/canary/workflow-guidance-anchors/calibration-pickup-step-move, tests/canary/workflow-guidance-anchors/coverage-axis-anchor, tests/canary/workflow-guidance-anchors/delegated-axis-exclusions, tests/canary/workflow-guidance-anchors/delegated-dispatch-declaration, tests/canary/workflow-guidance-anchors/delegated-entry-refusals, tests/canary/workflow-guidance-anchors/delegated-resumption-contents, tests/canary/workflow-guidance-anchors/dg-25, tests/canary/workflow-guidance-anchors/dg-26, tests/canary/workflow-guidance-anchors/dg-29, tests/canary/workflow-guidance-anchors/dg-29-verification-target, tests/canary/workflow-guidance-anchors/dg-30, tests/canary/workflow-guidance-anchors/dg-31, tests/canary/workflow-guidance-anchors/dg-31-contradiction-trigger, tests/canary/workflow-guidance-anchors/implement-spec-adoption-freshness, tests/canary/workflow-guidance-anchors/implement-spec-coverage-task-seeding, tests/canary/workflow-guidance-anchors/implement-spec-cross-harness-pointer, tests/canary/workflow-guidance-anchors/implement-spec-entry-validation, tests/canary/workflow-guidance-anchors/implement-spec-incapable-harness, tests/canary/workflow-guidance-anchors/implement-spec-inline-exception, tests/canary/workflow-guidance-anchors/implement-spec-mandatory-delegation-anchor, tests/canary/workflow-guidance-anchors/implement-spec-prose-owner-transfer, tests/canary/workflow-guidance-anchors/implement-spec-read-only-helper, tests/canary/workflow-guidance-anchors/implement-spec-red-preflight-route, tests/canary/workflow-guidance-anchors/implement-spec-review-charge-omitted, tests/canary/workflow-guidance-anchors/implement-spec-review-charge-order, tests/canary/workflow-guidance-anchors/implement-spec-review-charge-reversed, tests/canary/workflow-guidance-anchors/implement-spec-status-flip-anchor, tests/canary/workflow-guidance-anchors/implement-spec-worktree-before-preflight, tests/canary/workflow-guidance-anchors/implement-spec-write-delegation, tests/canary/workflow-guidance-anchors/line-anchor-missing, tests/canary/workflow-guidance-anchors/prepared-build-approval, tests/canary/workflow-guidance-anchors/prepared-build-freshness, tests/canary/workflow-guidance-anchors/prepared-review-axis-returns, tests/canary/workflow-guidance-anchors/prepared-review-blast-evidence, tests/canary/workflow-guidance-anchors/prepared-review-capable-handoff, tests/canary/workflow-guidance-anchors/prepared-review-inline-axis-route, tests/canary/workflow-guidance-anchors/prepared-review-legacy-entry-points, tests/canary/workflow-guidance-anchors/prepared-review-native-dispatch, tests/canary/workflow-guidance-anchors/prepared-review-runtime-capability, tests/canary/workflow-guidance-anchors/prepared-review-shared-evidence, tests/canary/workflow-guidance-anchors/reference-agent-push-rule, tests/canary/workflow-guidance-anchors/reference-bench-operational-layer, tests/canary/workflow-guidance-anchors/reference-category-context, tests/canary/workflow-guidance-anchors/reference-category-oracle, tests/canary/workflow-guidance-anchors/reference-category-setup, tests/canary/workflow-guidance-anchors/reference-category-work, tests/canary/workflow-guidance-anchors/reference-gate-authority, tests/canary/workflow-guidance-anchors/reference-kit-only-ship, tests/canary/workflow-guidance-anchors/reference-no-path-fallback, tests/canary/workflow-guidance-anchors/reference-progressive-loading-term, tests/canary/workflow-guidance-anchors/reference-refusal-route-shape, tests/canary/workflow-guidance-anchors/reference-retro-capture-owner, tests/canary/workflow-guidance-anchors/reference-retro-drain-owner, tests/canary/workflow-guidance-anchors/reference-skills-guidance, tests/canary/workflow-guidance-anchors/reference-upgrade-route, tests/canary/workflow-guidance-anchors/review-base-merged-main-tip, tests/canary/workflow-guidance-anchors/review-clean-terminal-result, tests/canary/workflow-guidance-anchors/review-cross-harness-opt-in, tests/canary/workflow-guidance-anchors/review-falsification-accept-routing, tests/canary/workflow-guidance-anchors/review-falsification-dispositions, tests/canary/workflow-guidance-anchors/review-persistence-anchor, tests/canary/workflow-guidance-anchors/review-preflight-explicit-base, tests/canary/workflow-guidance-anchors/review-repair-ticket-covers, tests/canary/workflow-guidance-anchors/review-repair-ticket-owner, tests/canary/workflow-guidance-anchors/review-standing-falsification, tests/canary/workflow-guidance-anchors/review-universal-claim-bar
Covers: NR23, NR24, NR25, NR26, NR27, NR28, NR29, NR30, NR31, NR32, NR38, NR39, NR40, NR41, NR42, NR43, NR44, NR45

## What to build

Deliver producer-owned execution metadata and complete native result import through the actual probe and record command entries.
Capture Execute's normalized mutated-execution exit directly, keeping command exit separate and no-mutated-execution null.
Render one probe_execution table after the unchanged probe/selection pair and count every byte of the final report and restoration suffix.
All result paths use that projection without changing restoration precedence, preserved-copy effects, or command exits.

Tickets 1 and 2 supply validated dispatch, scalar verification, typed repair routes, and the current-digest writer guard.
Their independent chunk reviews close before this overlapping command, help, guidance, fixture, and journey slice starts.
Verification import fills fields through probe.ReadResult and the existing no-follow evidence reader, retaining exact selected bytes and digest.
It retains explicit performer, model, effort, requirement, selector, source policy, and checked command exit.
No scalar/import combination or inferred focused-execution exit enters a write.

Compose responsebound.RequireComplete and treetarget.ReadLead at their public owners, then decode the canonical result with the pinned strict TOON decoder.
Reject projection frames without following their paths, and accept one optional canonical initial tree identity or a complete spill body.
Check exact envelope shape, unique tables, field types, and trailing_bytes against the complete remaining suffix.
Keep responseboundtest's spill accessor as a thin projection of the single marker owner.
Neither framing, byte counts, excerpt digests, nor an imported result authenticate origin or executed behavior.

Extract the one two-start canned Go transport into probetest/fixture.go and migrate both caller families in this same green checkpoint.
Probe's five listed caller files retain their canned streams, hooks, expectations, and effect assertions through thin adapters.
The real cmd/bench journey starts from recordtest.NewLinked and composes Prepare, kittest.WriteTree, and freshness.Publish in its resolved linked Go root.
It authors its own event expectations; record import runs in that same linked root with the primary-checkout guard active.
The new owner has no package variable, private fallback, production policy, or runbinary factory swap.
Test-owned files, marker paths, PATH restoration, real-Go passthrough, and answer alternation retain their current meanings.

Extend TestNativeRecordCLIJourney with actual bounded producer and consumer runs through the production registry.
The ordinary output retains its canonical tree lead; the --full report must actually spill before its selected complete body is imported.
Assert original selected-file bytes and digest, then drive projection, failed spill, truncated envelope, truncated suffix, and count mismatch refusals.
Their outside-path sentinel stays unopened and each refused import keeps record bytes unchanged.
No planted handler, synthetic domain response, reconstructed tree lead, or copied record builder can close this journey.

This final ticket also checks the completed family after its own focused green results.
Every first-use assertion, independent expectation, fixture closure, and headroom obligation already applies in its introducing ticket.
Global reconciliation does not repair or defer earlier green checkpoints.

Keep every existing assertion and independent expectation effective.
Do not rewrite old tests to agree with a new projection or replace a real entry with a planted handler.
Use existing fixture owners and record each independent expectation's named, compiling mutation red during implementation.
For every probe, pin the source and diagnostic, prove byte-identical restoration, and rerun the same focused green check.
Invalid, compilation-only, or restore-failed evidence cannot close an acceptance row.

The Writes line includes each first-use command binding, anchor holder, and transitive fixture unit.
Keep all 78 listed units and their BASE, MUTATE.json, EXPECT, meanings, and sufficient assertions intact.
Guidance changes keep current anchor placement and close their independent expectations in this ticket.
No canonical registry grows a duplicate rule, and no fixture helper acquires a second derivation.

## Acceptance

- [ ] A biting probe reports focused-execution 1 while Command.Run returns 0; a refused or baseline-only probe reports null for mutated execution.
- [ ] All render paths retain the original probe and selection schemas and order, followed once by the owner execution table.
- [ ] trailing_bytes exactly counts every report, preserved-copy, and restoration diagnostic byte after that table, including a zero-byte suffix.
- [ ] Complete canonical tree-prefixed stdout and the actual selected spill body import through real Command.Run producer and consumer paths.
- [ ] Each retained excerpt and digest equals the original selected file bytes, without re-encoding, concatenation, or a reconstructed tree lead.
- [ ] The long --full report actually spills under the private Bench home; the test reads its complete body only after observing its production marker.
- [ ] Projected and spill-failed stdout, malformed reserved frames, duplicate envelopes, bad field types, and old schemas refuse before a write.
- [ ] Truncation within an envelope or suffix and either direction of trailing_bytes disagreement refuse with unchanged record bytes and unopened marker-path sentinels.
- [ ] Missing, malformed, or repeated tree identity handling retains the spec's zero-or-one canonical initial-envelope rule and current dirty/head representations.
- [ ] Explicit scalar verification keeps all original assertions, and imported command exit must equal the retained explicit exit flag.
- [ ] Imported silent and restoration-failed results remain failed checkpoint evidence; import does not run a probe or widen source authority.
- [ ] Copied flag help identifies the exact producer verb and owner field for each exit, outcome, and restoration value.
- [ ] probe_test.go, baseline_test.go, outcome_test.go, refusal_test.go, and omit_file_test.go all reach one probetest.InstallStarts after the old executable builder is removed.
- [ ] cmd/bench's actual journey uses that same runner, with independently authored event streams and expected command/focused exits.
- [ ] Runner real-Go passthrough, per-installation marker reset, alternation, hooks, supplied exits, quoted events, test lifetime, and PATH cleanup preserve existing behavior.
- [ ] Every existing probe assertion survives, including refused subject bytes, empty preservation home, no run child, restoration, and preserved-copy checks.
- [ ] Substituting successful command exit for focused exit compiles and reds the independent real journey assertion.
- [ ] Omitting tree-lead handling, projection refusal, or exact suffix comparison separately reds the corresponding real producer/consumer refusal witness.
- [ ] Altering a renderer suffix without updating its count reds producer/consumer completeness; do not synthesize a matching expectation.
- [ ] Omitting runner passthrough or answer alternation separately reds the actual focused producer path and preserves valid evidence rules on restoration.
- [ ] Before each probe verdict is accepted, its source and exact diagnostic are pinned, restoration is byte-identical, and the same focused check returns green.
- [ ] Invalid, compilation-only, skipped, silent, or restore-failed proof closes no obligation.
- [ ] New result, import, journey, and runner files meet the spec's headroom limits now; verification_test.go and over-budget registries do not grow.
- [ ] All framing and marker readers use their single public owners; no second TOON parser, production test-helper import, global swap, or duplicate runner survives.
- [ ] The final three-chunk reconciliation keeps all 45 predicates, every sufficient original assertion, and the FT317, conversion, raw-status, and exact-reader exclusions.
