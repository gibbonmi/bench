# Prepare complete review diffs with stable file pages

Blocked by: 1-order-review-preparation.md
Writes: internal/preflight, internal/diff, internal/git, internal/consumers, internal/chargeevidence, .agents/skills/bench-craft-delegate/references/charge-evidence-format.md, internal/anchors, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/conformance/injected_ports_registry_test.go, cmd/bench/preflight_version_test.go
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

Implement the spec's descriptor table and file-identity rule, including patches
without `---` or `+++` lines. Preserve ambient rename behavior and verbatim bytes.
Move the consumers C-quote decoder into `internal/git`. Both callers use that
pure decoder. Preserve the existing hunk interpretation.

RE1 supplies the author-record-before-charge contract for this chunk's review.
RE10 and RE13 through RE16 in `Covers:` associate this chunk with its later
orchestrator checkpoint. The ticket author delivers the executable evidence.
After this ticket commits green, the orchestrator records the actual reviews
and control comparison in the spec. Those obligations remain outside this
author's acceptance criteria.

## Acceptance

- [ ] Growing an earlier patch across a page boundary preserves a later patch's digests.
- [ ] Separate spec, ticket, and review-record edits preserve untouched code patch digests.
- [ ] Inserting an earlier file preserves the retained path-to-digest mapping.
- [ ] A changed patch changes its source digest.
- [ ] Equal-content files remain separate manifest members.
- [ ] Retrieved fragments reconstruct every RE5 fixture byte for byte.
- [ ] A later-chunk charge uses its explicit predecessor base.
- [ ] A large Unicode patch reconstructs through the existing page protocol.
- [ ] Selection by the spec's file identity retrieves only that file's complete patch.
- [ ] Headerless patches have the exact role and path that the spec defines.
- [ ] Repeated-path patches remain separate sources in patch order.
- [ ] An unpartitionable body refuses the charge without partial publication.
- [ ] Generated fragments retain their actual producer provenance under `TestEvidenceReviewProvenanceRows`.

Use planned `TestReviewFilePageStability`, `TestReviewFileIdentity`,
`TestReviewFileReconstruction`, and `TestReviewFileSelectedStream`.
Place command scenarios in new files under `internal/preflight/evidencecmd`.
Keep file growth within the existing budget. Reuse package fixture and
traversal helpers. Pure partition cases belong under `internal/diff`.

Use RE5's deterministic fixture and independent baseline contract.
Include pure renames, empty additions and deletions, binary patches, and mode changes
in the selected-stream test. Include a file-to-symlink change with two patches
for one path. Test rename detection both enabled and disabled.

RE5 also covers external-diff and forced-color outputs through the existing
failed-publication posture. Test exact reconstruction or refusal, never a partial
artifact. Preserve ambient Git output.
Use equal-length paths and multi-page equal content for the RE6 case.
Require equal later page digests while each path retains its own membership.

Do not rewrite production sources during ordinary tests. Prove the named
spec mutations through `bench probe` and require restoration.

Run focused diff, git, consumers, preflight, and chargeevidence suites with `-parallel 2`.
Retain binding, navigation, export, provenance, movement, and response-bound tests.
Run applicable command-registry and injected-port checks. Regenerate the
format reference after its description changes. Run its prose check too.

The package invariant is complete reconstruction with one source of paging
and path-decoding knowledge. The orchestrator performs the real control under
the spec's RE2 checkpoint obligations. It records both frozen pins before an
evidence-only spec update. Existing record and review rules apply to that update.
