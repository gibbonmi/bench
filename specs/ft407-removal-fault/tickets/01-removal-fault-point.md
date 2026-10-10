# Add a removal fault point to the adoption transaction seam

Blocked by: none
Writes: internal/adopt/transaction.go, internal/adopt/transaction/publication.go, internal/adopt/transaction/transaction.go, internal/systemtest/lifecycle_test.go, ROADMAP.md, roadmap/FT407.md
Covers: none

## What to build

The `BENCH_LINK_FAULT` seam fires only when a change is published by a rename. A change that removes its destination never meets the seam.
So no test can interrupt `bench unlink` after its file removals and before its final manifest write.

The transaction store gains a removal hook that sits beside its rename hook. The adoption publication owner routes both hooks through one fault check.
The existing grammar is unchanged: an integer, `last`, and `interrupt:<n>` each name the nth change of a promotion.
That change can now be a rename or a removal.

One system test in `internal/systemtest/lifecycle_test.go` runs the selected binary with `BENCH_KIT` set to the owner kit.
It links a repository, trims the unlink manifest to two managed rows, and interrupts `bench unlink` at the removal of the pre-push hook.
That is the last change of the first promotion. The test then calls `Undo` on the journal that the interrupted process left in its `TMPDIR`.

This ticket retires roadmap row FT407. The same commit removes the FT407 index line from `ROADMAP.md` and deletes `roadmap/FT407.md`.

## Acceptance

- [ ] An interrupt at the removal of the hook exits 97 and leaves the manifest unchanged.
- [ ] After `Undo`, both removed managed files, the stripped `AGENTS.md`, and the hook hold their pre-unlink content.
- [ ] After `Undo`, the link manifest holds its pre-unlink content.
- [ ] The existing rename faults and interrupts in the lifecycle system tests keep their results.
- [ ] `ROADMAP.md` has no FT407 line, and `roadmap/FT407.md` is gone.
