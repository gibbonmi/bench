## Standards

Finding count: 1. Worst issue: stale review-artifact references in all three live repair tickets.

- `auto-fix` — the three repair tickets cite `reviews/spec-authoring-and-light-path.md`, but the required transient pickup was deleted after its findings closed. Remove the provenance clauses while retaining each self-contained requirement; do not restore the review as durable history.

## Spec

Finding count: 2. Worst issue: live repair tickets cite a missing review source.

- `auto-fix` — remove the same three dangling review-path citations and leave the complete finding predicates in their owning tickets.
- `auto-fix` — `repair-profile-loop-routing.md` added two fixtures and updated the canary count, but the spec's per-ticket census map enumerates only the original eight fixture-producing tickets. Add a dedicated census row and acceptance `covers` tag for this repair, and tag its loop-routing acceptance to WF10.

## Coverage

Finding count: 0. Worst issue: none. The refreshed Coverage axis re-derived all 204 producer fixtures and found the prior profile-loop and empty-ticket gaps closed.

Raw findings: Standards 1, Spec 2, Coverage 0. De-duplicated repair targets: 2.
