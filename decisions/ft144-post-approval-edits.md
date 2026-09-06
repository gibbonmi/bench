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

## Notes

## Decisions so far

- [The permission rule](ft144-post-approval-edits/tickets/1.md): Resolved 2026-08-02: rule (a).

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
