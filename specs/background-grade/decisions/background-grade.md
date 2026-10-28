# Background grade of lane-passed tips (FT451)

Status: ready

## Destination

After a lane pass, the gate grades the new tip in the background.
A chunk checkpoint or a landing of that tip then finds the background green under its grade key and runs no second gate.
Background grades run inside a decided machine bound measured at `-p=2`, and the grade of a superseded tip takes no new test slot.

## Notes

Commitment outcome: `FT451-background-grade` in milestone `quality-2`, admitted measure-first after FT448.
Domain: the grade key, the evidence store, the checkout cache and lock, the fast lane, the checkpoints, the landing broker, and the process lifetime owner.
A background grade is a gate run that a lane pass starts and that runs detached from that verb.
Avoid: "prepared run" for a background grade, "daemon", "precompute".

In this map, a consumer is the verb that looks up a background green: the chunk checkpoint or the landing.
Avoid "consumer repo" here; `CONTEXT.md` names that thing a linked repo.

Skills: craft-domain, craft-grill, craft-research, and craft-seams for the spec.
By reviewer direction, this map's research ran in one session with no delegated reads.
A map-owned asset stays in the map's assets folder, specs/background-grade/decisions/background-grade/assets/.
The research asset is `background-grade/assets/research.md`.

## Decisions so far

- [What does the gate reuse machinery do for a background grade?](background-grade/tickets/1.md): one shared green store, one lock per checkout, and worktree reads.
- [What does the resource census under ADR 24 show?](background-grade/tickets/2.md): a closed gate environment with no Bench width input, and recorded contention on this host.
- [Which consumers does a background grade serve, and does FT451 split?](background-grade/tickets/3.md): both consumers in one outcome, with no split.
- [What starts a background grade, and in which repositories?](background-grade/tickets/4.md): a publishing lane pass, in a kit worktree only.
- [Does the grade identity read the run's own run-binary selection?](background-grade/tickets/5.md): yes, through one derivation.
- [How does a consumer meet an in-flight background grade?](background-grade/tickets/6.md): a claim per key, and a consumer waits for its own key.
- [What machine bound holds the background grades?](background-grade/tickets/7.md): one at a time, at `-p=2 -parallel=2`, at the lowest priority.
- [What subject does a background grade compute?](background-grade/tickets/8.md): the next consumer's exact run, with no live worktree read.
- [How does a background red reach the author, and what does a consumer do with it?](background-grade/tickets/9.md): a signal only.
- [What happens to the prepared run and its guidance?](background-grade/tickets/10.md): FT451 retires `--prepare`.
- [Where does the census measurement run, and what passes it?](background-grade/tickets/11.md): the first build ticket, with a 1.2 times stop rule.
- [What stays out of scope?](background-grade/tickets/12.md): ten confirmed exclusions, and ticket 4 adds one more.
- [What happens to a background grade when its worktree lands, is released, or is cleaned?](background-grade/tickets/13.md): stop it.
- [Where does a background delta green go?](background-grade/tickets/14.md): the delta record of its worktree.
- [Does a background grade run for a landing that takes the reader set?](background-grade/tickets/15.md): no.
- [Who starts a queued grade, and when is a claim live?](background-grade/tickets/16.md): a chain of grade processes, and a live owner in this boot.

## Not yet specified

## Spec-writer discretion

- The claim mechanism, the queue storage, and the priority mechanism of each platform, provided they meet tickets 6, 7, 13, and 16.
- The form of the Bench-owned width input, provided that it is not a grade-key member.
- The rule that classifies a kit worktree as an open build of a staged spec and derives its landing deliverable. A wrong answer must only cost a lookup miss.
- The one stable identity value that a run with no run-binary selection frames under ticket 5.
- The wording of the background outcome line in `bench commit` and of its `bench status` row.
- Whether the background grade takes a new ADR, or amends ADR 0002, ADR 0017, ADR 0024, and ADR 0031. Each of those four states a rule that FT451 changes.

## Out of scope

- A new grade when the destination moves. A landing after a destination move runs its own gate.
- A background grade of the complete checkpoint subject.
- The FT255 shared pool and any change of a foreground width.
- The FT448 rerun cost of a repair record.
- The FT372 family of nested gates that outlive their probe.
- A landing stack under ADR 0032 as a consumer.
- A landing that reads a delta green.
- The GR44 edge between two foreground runs.
- A change of the ADR 0031 rule that the chunk checkpoint keeps its own process.
- A freshness window longer than 60 minutes for a background green.
- A background grade in a linked repo; a parked idea holds the opt-in.

## Sources

- Path: `.bench/commitment.json`
  Supports: the lookup and bound criteria of `FT451-background-grade`.
  Drift: a `bench commitment plan` change of the outcome.
- Path: `docs/adr/0031-one-grade-key-names-each-whole-gate-green.md`
  Supports: the grade key members and the executable of each consumer.
  Drift: a superseding ADR.
- Path: `docs/adr/0024-gate-phase-overlap-requires-a-resource-census.md`
  Supports: the census requirement before adoption.
  Drift: a superseding ADR.
- Path: `docs/adr/0030-one-lifetime-owner-and-process-user-protection.md`
  Supports: the lifetime rules of a grade that outlives its verb.
  Drift: a superseding ADR.
- Path: `internal/gate/run_transaction.go`
  Supports: the reuse order, the checkout lock, and the run-binary selection of a run.
  Drift: a change to `executeSubjectOwning` or the lock helpers.
- Path: `internal/gate/subject.go`
  Supports: the grade identity members, the live worktree reads, and the closed environment.
  Drift: a change to `buildSubjectForTree`.
- Path: `internal/gate/checkpoint.go`
  Supports: the chunk subject and the prepared run that ticket 10 retires.
  Drift: a change to `executeChunkCheckpoint` or `gradesSubjectAlone`.
- Path: `internal/gate/complete_checkpoint.go`
  Supports: the published tree and the broker run that the landing subject copies.
  Drift: a change to `executeCompleteCheckpoint` or `runUnderBroker`.
- Path: `roadmap/FT255.md`
  Supports: the `-p=2` measurement decision and the contention record.
  Drift: a drain edit of the row or an FT255 delivery.
