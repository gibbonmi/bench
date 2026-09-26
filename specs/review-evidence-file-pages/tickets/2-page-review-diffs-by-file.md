# Prepare complete review diffs with stable file pages

Blocked by: 1-order-review-preparation.md
Writes: internal/preflight, internal/diff, internal/consumers, internal/chargeevidence, .agents/skills/bench-craft-delegate/references/charge-evidence-format.md, internal/anchors, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/conformance/injected_ports_registry_test.go, specs/review-evidence-file-pages/spec.md
Covers: RE1, RE2, RE3, RE4, RE5, RE6, RE7, RE8, RE9, RE10, RE13, RE14, RE15, RE16

## What to build

Deliver chunk RE2. Verify the current collector, pager, and explicit-base behavior.
Publish ordered file patches through the existing diff owner. Keep structural
fragments separate. Reconstruct the original full diff without a second
monolithic payload.

Reuse `chargeevidence.Build` for paging and digests. Keep ordinals local to the
pack. Match retained bytes through current manifest membership and descriptors.
Preserve producer provenance. Extend the shared `preflighttest` fixture and
its readers rather than copy the harness.

The current patch-path parser belongs to consumers. If the diff owner needs
that rule, move its implementation and migrate its existing reader here.
RE1 supplies the author-record-before-charge contract for this chunk's review.
The orchestrator records narrow reads and findings in the spec. It runs one
independent full control on the same frozen pair before the adoption decision.

## Acceptance

- [ ] Growing an earlier patch across a page boundary preserves a later patch's digests.
- [ ] Separate spec, ticket, and review-record edits preserve untouched code patch digests.
- [ ] Inserting an earlier file preserves the retained path-to-digest mapping.
- [ ] A changed patch changes its source digest.
- [ ] Equal-content files remain separate manifest members.
- [ ] Retrieved fragments reconstruct every RE5 fixture byte for byte.
- [ ] A later-chunk charge uses its explicit predecessor base.
- [ ] A large Unicode patch reconstructs through the existing page protocol.
- [ ] Selecting one file retrieves its complete patch without another file's body.
- [ ] Generated fragments retain their actual producer provenance.
- [ ] This spec records all narrow rounds and the same-pair full control.
- [ ] The reviewer receives the comparison before permanent adoption.

Use planned `TestReviewFilePageStability`, `TestReviewFileIdentity`,
`TestReviewFileReconstruction`, and `TestReviewFileSelectedStream`.
Place command scenarios in new files under `internal/preflight/evidencecmd`.
Keep file growth within the existing budget. Reuse package fixture and
traversal helpers. Pure partition cases belong under `internal/diff`.

Capture original full-diff outputs for the fixture family before production
edits. Keep the permanent comparison against those baseline bytes. Do not
rewrite production sources during ordinary tests. Prove the named spec
mutations through `bench probe` and require restoration.

Run focused diff, consumers, preflight, and chargeevidence suites with `-parallel 2`.
Retain binding, navigation, export, provenance, movement, and response-bound tests.
Run applicable command-registry and injected-port checks. Regenerate the
format reference after its description changes. Run its prose check too.

The package invariant is complete reconstruction with one source of paging
and path-decoding knowledge. The control is real review evidence, not a
simulated result. Use independent sessions under the review phase's rules.
Record both frozen pins before an evidence-only spec update. Apply existing
record and review rules to that update.
