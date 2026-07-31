# Checkpoint and integrate attributed patches

Blocked by: Start runs and assign frontier tickets

Ownership fence: `internal/specbuild`
Assumptions: Start and Assign persist the exact candidate and assignment facts

## What to build

Extend `internal/specbuild` with canonical bounded receipt validation,
attributed checkpoint commits, exact-old-tip candidate integration, unchanged
disjoint sibling replay, retained provenance, and lifecycle-owned assignment
release. Re-read the live candidate and ticket facts before each transition;
delegate-supplied evidence never becomes project-green evidence.

## Acceptance

- [ ] [R10-R15] Checkpoint accepts only complete row, seam-check, independently produced probe, exact-tree, ownership, and assumption evidence; honest already-covered/not-TDD-able rows remain valid.
- [ ] [R16-R20] Integrate accepts verified checkpoints, replays unchanged disjoint siblings, refuses overlap or drift without mutation, uses compare-and-swap, and resumes release without replaying twice.
- [ ] [R54] The final-newline receipt posture is deterministic and malformed framing has one stable refusal.
- [ ] An independent probe varies mutation kind across omission, substitution, reordering, or duplication rather than repeating the delegate's defect class.
