# Grade a landing subject in the background

Blocked by: 08-queue-and-supersede-background-grades.md
Writes: internal/gate/background_grade.go (new), internal/gate/background_landing.go (new), internal/gate/background_landing_test.go (new), internal/gate/complete_checkpoint.go, internal/gate/reader_plan.go, internal/gate/backgroundstate (new), internal/landing/published, CHANGELOG.md
Covers: BG67, BG68, BG69, BG70, BG71, BG72, BG73, BG74

## What to build

A kit worktree without an open build of a staged spec gets a landing grade.

- The subject is the tree that the landing would publish from the tip, composed on the destination at the start of the grade.
- A tickets-only deliverable takes the folder close and the delivery closure of `published.Tree` (`internal/landing/published/published.go:19-38`).
- A spec-less landing applies no publish transform, so its subject is the tip composed on the destination.
- The grade runs as the authenticated broker process, as `runUnderBroker` authenticates it (`internal/gate/complete_checkpoint.go:42-91`).
- A broker whose digest differs from its manifest makes the grade record `aborted` and run no oracle.
- The grade takes the phase schedule of the destination commit.

A lane pass queues no landing grade when the landing plan of the tip selects the reader set.
The lane pass evaluates that plan against the destination at the time of the lane pass.
The plan reads the green marker, the reader hold, and the reader count, as a landing does.

The rule that derives the landing deliverable is the author's choice. A wrong deliverable only costs a lookup miss.

## Acceptance

- [ ] A landing grade writes its green under the key of the published tree composed on the destination.
- [ ] A spec-less landing grade keys the tip composed on the destination, with no transform.
- [ ] The grade key of a landing grade frames the digest of the installed broker.
- [ ] A broker whose digest differs from its manifest makes the grade record `aborted` and run no oracle.
- [ ] A landing grade frames the phase schedule of the destination commit.
- [ ] A tip whose landing plan selects the reader set queues no grade, and a tip whose plan selects the whole gate queues one.
- [ ] A landing grade with a wrong deliverable writes its green only under the key that it ran.
