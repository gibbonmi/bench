# FT144 post-approval spec edits and the seam-move case

Status: ready

## Destination

One permission rule for both phases: what a build or review phase may do
under batch approval when its finding lands on an approved spec, including
the build-phase case where a story's intent stands but its pinned seam moves.
This map is the decision source for the kit edit (built under
`craft-synthesis`) to `/bench-review-implementation`, `/bench-implement-spec`
(whose fence and veto-flag prose the seam-move permission must reconcile
with), and the workflow prose.

## #1: The permission rule

Blocked by: none
Type: Grill

### Question

When a phase's finding lands on an approved spec rather than on the code,
what may it do under batch approval? (a) factual corrections — a citation
resolving to nothing, or a described mechanism the tree contradicts — may be
made post-approval under the in-line veto-flag convention (the
`**Post-approval correction, flagged:**` marker), while anything changing
what gets built stops for sign-off; or (b) all post-approval spec edits stop
and the review persists findings to `reviews/<spec-slug>.md` instead. The
build-phase pair is the same choice in its own dress: either the existing
route is right and a build whose seam moves pays the round-trip back to
`/bench-write-spec`, or the workflow gets a named lighter case for "intent
stands, seam moves".

### Answer

Resolved 2026-08-02: rule (a). Amended 2026-08-03 after doc review to an
intent-based boundary, one rule for both phases: under batch approval, a
factual correction (a citation resolving to nothing, or a described mechanism
the tree contradicts) and a seam move that preserves the story's intent and
observable outcome are permitted, each marked in-line with the
`**Post-approval correction, flagged:**` marker and flagged for veto;
anything that changes a story's intent or observable outcome stops for
sign-off. The round-trip alternative (b) is rejected on its stop-all-edits
half only — it pays a full cycle on every false citation; the review phase's
persistence of findings to `reviews/<spec-slug>.md` is untouched by this
ruling.

Timing rule, added 2026-08-03: during an active spec-build run, a
post-approval correction to the staged spec waits for a run boundary (before
`start`, or after `integrate` or `abandon`) — a mid-run spec commit trips the
lifecycle's staged-spec mismatch refusal, which has no exemption, and this
rule is the decided disposition rather than a lifecycle change.

## Not yet specified

## Spec-writer discretion

- Exact wording and placement of the rule in `/bench-review-implementation`
  and the workflow prose, provided the factual-versus-behavioral boundary and
  the mandatory flag survive verbatim in meaning.

## Out of scope

- Weakening the spec sign-off gate itself: absent batch approval, spec
  sign-off remains a hard stop; this rule governs post-approval findings
  only.
- FT144's other kit edit — the `craft-spec` two-audience edge-inventory
  prompt — is already directed by the roadmap row and needed no ruling; it
  rides the row, not this map.

## Sources
