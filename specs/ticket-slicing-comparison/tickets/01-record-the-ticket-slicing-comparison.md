# Record the ticket slicing comparison

Blocked by: none
Writes: docs/research/2026-10-10-ticket-slicing-comparison.md (new)
Covers: none

## What to build

The reviewer directed on 2026-10-10 a one-session measurement run. The run compares two ticket-slicing methods of `/bench-write-spec` on one built spec. The OLD method lets the spec fork slice every ticket, and one review round grades the spec and the tickets together. The NEW method slices one chunk at a time after per-chunk research, with a cheap review for each chunk.

Both arms sliced the FT349 spec at its pre-build commit. The coordinator graded each arm's `Writes:` lines against the files that the real build changed. The report records the subject, the method of each arm, the measurements, the missed paths by class, the threats to validity, and one recommendation. Two consultants reached that recommendation together.

## Acceptance

- [ ] `docs/research/2026-10-10-ticket-slicing-comparison.md` opens with the recommendation, the scope, and the evidence status.
- [ ] The report names the BASE commit and the landing commit, and it gives the method of each arm.
- [ ] The report holds a measurement table and a table of the missed paths by class.
- [ ] The report states the threats to validity and names what it read and what it did not read.
- [ ] `bench gate-prose` passes on the report.
