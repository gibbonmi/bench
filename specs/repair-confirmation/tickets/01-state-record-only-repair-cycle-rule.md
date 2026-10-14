# State the record-only repair cycle rule

Blocked by: none
Writes: .agents/skills/bench-craft-line/references/bounded-repair-policy.md, .bench/commitment.json, CHANGELOG.md, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/conformance/implementation_continuation_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns/, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary/
Covers: RC1, RC2, RC3

## What to build

The bounded repair policy states the reviewer's 2026-10-05 decision.
A correction that changes only the review record consumes no repair cycle, even when it changes a disposition.
The issuing axis reaffirms a changed disposition.

Add the two sentences after the sentence pair "Evidence-only corrections consume no repair cycle. Batch all cited corrections before verification."
Do not put a sentence between those two sentences, because one anchor needle holds both.
Keep the evidence-only definition and the issuing-axis sentence byte-identical.
Split the paragraph so that each part holds six sentences or fewer.
Keep the `Scope and allowance` heading unchanged, because `CONTEXT.md` links it.

Add one anchor row for each new sentence to `implementationContinuationAnchors`, in section `Classification and completion`.
Each diagnostic starts with `implementation continuation: ` and has a unique subject.
Add the independent twin of each row to `wantPredicates`.
Record a demonstrated red for each twin, as the working agreement requires.
The retained-workflow registry is at 396 of 400 lines, so add exactly two rows there.
The closure-only files over their structure budget must not grow.

This ticket names `.bench/commitment.json` for the spec-stage binding only, and it does not edit that file during the build.

## Acceptance

- [ ] RC1: the policy states the record-only no-cycle rule, disposition changes included. Its anchor and twin red when the sentence is absent.
- [ ] RC2: the policy states that the issuing axis reaffirms a changed disposition, and its anchor and twin red when the sentence is absent.
- [ ] RC3: the evidence-only definition keeps its exact needle, including the word disposition.
- [ ] The prose-mechanics check passes on the policy file.
