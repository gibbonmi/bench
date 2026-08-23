# FT144 post-approval spec edits and the seam-move case

Status: ready

## Destination

One permission rule for both phases: what a build or review phase may do
under batch approval when its finding lands on an approved spec. This includes
the build-phase case where a story's intent stands but its pinned seam moves.
This map is the decision source for the kit edit (built under
`craft-synthesis`) to `/bench-review-implementation`, `/bench-implement-spec`,
and the workflow prose. The seam-move permission must reconcile with the
build phase's own fence and veto-flag prose.

## #1: The permission rule

Blocked by: none
Type: Grill

### Question

When a phase's finding lands on an approved spec rather than on the code,
what may it do under batch approval? Option (a): factual corrections — a
citation resolving to nothing, or a described mechanism the tree contradicts —
may be made post-approval. They use the in-line veto-flag convention (the
`**Post-approval correction, flagged:**` marker). Anything changing what gets
built still stops for sign-off. Option (b): all post-approval spec edits
stop, and the review persists findings to `reviews/<spec-slug>.md` instead.

The build-phase pair poses the same choice in its own dress. Either the
existing route is right, and a build whose seam moves pays the round-trip
back to `/bench-write-spec`. Or the workflow gets a named lighter case for
"intent stands, seam moves".

### Answer

Resolved 2026-08-02: rule (a). Amended 2026-08-03 after doc review to an
intent-based boundary, one rule for both phases. Under batch approval, two
kinds of change are permitted post-approval. A factual correction (a citation
resolving to nothing, or a described mechanism the tree contradicts) is one.
A seam move that preserves the story's intent and observable outcome is the
other.

Each permitted change is marked in-line with the
`**Post-approval correction, flagged:**` marker and flagged for veto. Anything
that changes a story's intent or observable outcome still stops for sign-off.

The round-trip alternative (b) is rejected on its stop-all-edits half only —
it pays a full cycle on every false citation. The review phase's persistence
of findings to `reviews/<spec-slug>.md` is untouched by this ruling.

Timing rule, added 2026-08-03: during an active spec-build run, a
post-approval correction to the staged spec waits for a run boundary (before
`start`, or after `integrate` or `abandon`). A mid-run spec commit trips the
lifecycle's staged-spec mismatch refusal, which has no exemption. This rule is
the decided disposition rather than a lifecycle change.

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
  prompt — is already directed by the roadmap row and needed no ruling. It
  rides the row, not this map.

## Sources
