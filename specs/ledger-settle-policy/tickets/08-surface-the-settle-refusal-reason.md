# Surface the settle refusal reason

Blocked by: 05-migrate-the-assignment-mutators.md, 07-migrate-the-composition-adapter-onto-the-settle-policy.md
Writes: internal/landing/composition.go, internal/landing/composition_test.go, internal/landing/landing_helpers_test.go, internal/landing/settlepolicy/, internal/landing/merge.go, internal/landing/landing.go, internal/landing/merge_test.go, CONTEXT.md, tests/canary/docs-currency-token-diet/signal-vocabulary-drift, tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary, tests/canary/workflow-guidance-anchors/context-coverage-map-term, tests/canary/workflow-guidance-anchors/context-coverage-row-parts, tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary, tests/canary/workflow-guidance-anchors/context-decision-map-term, tests/canary/workflow-guidance-anchors/context-reader-sweep-term, tests/canary/workflow-guidance-anchors/context-ticket-vocabulary
Covers: LS28, LS29, LS30, LS31, LS32, LS33, LS34, LS35, LS38

## What to build

Verify the premise first. Read `Conflict`, `ConflictError`, and
`CompositionResult` in internal/landing/composition.go. Read the two
`ConflictError` producers in internal/landing/merge.go and
internal/landing/landing.go. Read `landingConflictRefusal` in
internal/worktree/land_refusal.go, which renders `conflict.Error()` and
`conflict.Paths` through one refusal face.

Read the settle policy and its four refusal reasons, which ticket
06-add-the-settle-policy-child.md adds. Read the adapter ticket
07-migrate-the-composition-adapter-onto-the-settle-policy.md leaves. Read the
three census wrappers tickets 02, 04, and 06 add. Read the glossary shape of an
existing term in CONTEXT.md.

Add the settle-refusal reason field to `Conflict`. The adapter carries the
policy's reason and its refusing paths onto that field. `ConflictError.Error`
reads exactly `composition conflict: <kind>; settle refused: <reason>
(<paths>)` when a reason exists. It keeps today's `composition conflict:
<kind>` when no reason exists. The paths are the refusing paths,
comma-separated in merge-tree order.

The kind stays at the front of the text, so every existing prefix match holds.
LS33 covers two kinds at the landing surface.
`TestLandCommandConflictRefusalNamesThePath` in
internal/worktree/land_surface_test.go covers `textual`, and
`TestMergeRefusesACaptureAddAddWithDisagreeingModes` in
internal/worktree/merge_test.go covers `mode`. Neither test changes, and this
ticket edits neither file.

This ticket is the last one that touches `internal/landing`, so it carries that
package's whole-package invariant. It writes no new file inside
`internal/landing/` itself. It is also the last ticket of the spec, so the
reviewer reads the three census wrappers together for LS35.

Add the term **settle verdict** to CONTEXT.md. Define it as the policy's answer
for one conflicted path: a side, a union, a removal, or a refusal with a
reason. Give it the Avoid list `resolution` and `merge result`. `bench anchors
CONTEXT.md` names seven required needles, and this entry touches none of them.

The eight fixture directories in the `Writes:` line are closure headroom. They
pin CONTEXT.md, and this ticket edits none of them.

Write the new test `TestConflictErrorNamesTheSettleRefusalReason` in
internal/landing/composition_test.go. Its first case is a conflict on `named`
beside `capture/learnings.md`. Its second case is a symlink at
`capture/learnings.md`. Its third case is a permission change against a content
change. Its fourth case is binary content on both sides. Its fifth case is
`capture.md` beside a symlinked `capture/learnings.md`.

Assert on every case that the text does not equal `composition conflict:
<kind>` and holds the separator `; settle refused: `. Keep
internal/landing/composition_test.go under the 400-line budget, which ticket
07-migrate-the-composition-adapter-onto-the-settle-policy.md reaches first.

## Acceptance

- [ ] With a conflict on `capture/learnings.md` beside one on `named`, the error reads `composition conflict: textual; settle refused: path outside the capture table (named)`.
- [ ] With a symlink at `capture/learnings.md`, the error holds `settle refused: non-regular mode` and names that path.
- [ ] With a permission change on one side and a content change on the other, the error holds `settle refused: mode disagreement`.
- [ ] With binary content on both sides, the error holds `settle refused: union content not text`.
- [ ] With `capture.md` beside a symlinked `capture/learnings.md`, the error names exactly one reason, and that reason is `path outside the capture table`.
- [ ] With a refusing settle, the error does not equal `composition conflict: <kind>` and holds the separator `; settle refused: `.
- [ ] The landing refusal surface still holds the text `composition conflict: textual` and the conflicted path.
- [ ] `TestMergeRefusesACaptureAddAddWithDisagreeingModes` passes with its test logic unchanged.
- [ ] Each of the three new packages holds a `purity_census_test.go` whose census names that package's own source.
- [ ] CONTEXT.md holds the term `settle verdict` with its Avoid list.
- [ ] `bench structure` reports `internal/landing/composition_test.go` under the 400-line budget.
- [ ] The pre-existing `internal/landing` and `internal/worktree` suites pass with their test logic unchanged, except the new reason assertions.
- [ ] Self-probe: store the reason without rendering it, and report the reason test red.
