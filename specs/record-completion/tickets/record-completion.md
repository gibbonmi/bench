# Record final completion evidence

Blocked by: none
Writes: internal/reviewrecord/, specs/record-completion/, CHANGELOG.md
Covers: none

## What to build

Add `bench record completion <slug> --source <commit>` at the command seam.
The command derives the completed state, source digest, reconciler, and planned
acceptance rows from the named source. It completes an existing review record
only after the retained evidence satisfies the final checkpoint.

## Acceptance

- [ ] The command refuses incomplete, failed, or stale chunk verification,
  review, reconciliation, and final verification evidence without a write.
- [ ] The command preserves existing evidence and surrounding prose when it
  writes the derived completion entry.
- [ ] A repeated call for the same source leaves the record bytes unchanged.
- [ ] The form declaration owns its grammar, help text, and output row.
