# Migrate live-spec path consumers

Blocked by: Add safe live-spec primitives

## What to build

Move status orphan pairing, status ROADMAP reconciliation, and roadmap context
parsing onto `internal/spec`'s live-path primitives. Add deadline-backed
handoff and roadmap-context contracts for a special-file spec candidate and
remove every independent production derivation enumerated by review.

## Acceptance

- [x] Story 19 and its acceptance-coverage row are green.
- [x] Handoff returns by its deadline and roadmap context retains named empty-metadata evidence for a FIFO spec.
- [x] Production code outside `internal/spec` no longer constructs or parses the live folder-spec form.
