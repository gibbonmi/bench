# Parallelize drain reading and centralize writing

Blocked by: none
Writes: .agents/commands/bench-drain.md, CHANGELOG.md, internal/conformance/recurrence_maintenance_contract_test.go, specs/drain-delegation/tickets/01-parallelize-drain-reading-and-centralize-writing.md

## What to build

Each drain uses available delegates to read independent evidence in parallel.
The coordinator keeps the single snapshot, cross-source judgment, and final
verdict. One write delegate authors the tracked batch after all applicable
implement-now items land, or after read and decision resolution when none exists.

## Acceptance

- [ ] The drain assigns non-empty roadmap reconcile, idea and journal mapping, and retro analysis scopes to at most three read delegates.
- [ ] Each read delegate uses only the snapshot and its named paths, edits nothing, and returns a fixed evidence shape.
- [ ] The coordinator resolves duplicate incidents and reviewer decisions before any batch writer starts.
- [ ] An implement-now writer can overlap the remaining reads.
- [ ] The batch worktree starts after every implement-now item lands green, or after read and decision resolution when none exists.
- [ ] Each delegate uses `craft-delegate` and `craft-line`; each writer gets the required isolation and routing.
- [ ] One write delegate authors the complete tracked batch only when tracked changes remain.
- [ ] After approval, the coordinator removes ignored sources.
- [ ] The coordinator writes the ignored handoff last.
- [ ] The reviewer batch separates tracked work from proposed ignored-source removals and journal verdicts.
- [ ] The recurrence-maintenance check reds when the command loses the parallel-read or single-writer contract.
