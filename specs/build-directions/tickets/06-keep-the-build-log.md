# Keep one git-ignored build log

Blocked by: 05-scope-review-rounds.md
Writes: .agents/commands/bench-implement-spec.md, .gitignore, .bench/BENCH-reference.md, docs/field-guide.html, CHANGELOG.md, projects/benchkit.md, internal/landing/landing.go, internal/landing/landing_test.go, internal/landing/local_capture.go (new), internal/landing/local_capture_test.go (new), internal/landing/state_test.go, internal/adopt/init.go, internal/adopt/adopt_test.go, internal/conformance/local_capture_tracking_test.go, internal/conformance/tier_live_tree_test.go, internal/worktree/build_outputs.go, internal/worktree/land_local_capture_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_check_floor_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/guidance-prose-budgets, tests/canary/line-routing, tests/canary/skill-description-budgets, tests/canary/workflow-guidance-anchors
Covers: BD39, BD40, BD41, BD42, BD43, BD44, BD45, BD46, BD47, BD48

## What to build

The orchestrator keeps one git-ignored build log, and every reader of local capture accepts it.
The Land section of `/bench-implement-spec` takes the three build-log sentences that the spec quotes as a new paragraph.
The kit `.gitignore` takes the line `capture/build-log/` in its capture block.

`LocalCapturePath` moves from `internal/landing/landing.go` to a new `internal/landing/local_capture.go` with a new exported constant `landing.BuildLogDir`.
The predicate accepts that directory, its slash-less form, and each clean path beneath it.
It refuses a control byte, a `..` segment, a doubled slash, and a sibling name that only shares the prefix.
A new `internal/landing/local_capture_test.go` holds the predicate table test and a real-checkout fingerprint test.

`internal/adopt/init.go` adds the constant to `localCaptureIgnoreLines`, so `bench init` writes the ignore line.
The three want strings in `internal/adopt/adopt_test.go` change in place, because that file cannot grow.
`internal/conformance/local_capture_tracking_test.go` takes TestKitGitignoreIgnoresTheBuildLog, and `tier_live_tree_test.go` classifies it as a live-tree test.
New Require rows in `registry_chunk_chain.go` pin the three build-log sentences, with independent mirrors.
The implement-spec line budget moves from 91 to 93.
One `CHANGELOG.md` entry records the change.

## Acceptance

- [ ] Implement-spec holds the three build-log sentences, and TestChunkChainAnchors goes red when any one is removed (BD39, BD40, BD41).
- [ ] The kit `.gitignore` holds a line equal to `landing.BuildLogDir`, and TestKitGitignoreIgnoresTheBuildLog goes red without it (BD42).
- [ ] `bench init` writes the build-log line after the three capture lines, and every init line satisfies the landing predicate (BD43, BD44).
- [ ] The predicate accepts the build-log directory and its clean paths and refuses each unsafe or sibling name (BD45, BD46, BD47).
- [ ] `CheckoutFingerprint` returns the same value before and after an ignored build-log file appears (BD48).
- [ ] `internal/adopt/adopt_test.go` keeps its line count, and implement-spec has 93 lines under its budget row.
