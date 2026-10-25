# Read each touched row against every source on its topic

Blocked by: none
Writes: .agents/commands/bench-drain.md, CHANGELOG.md
Covers: none

## What to build

A drain splits its read delegates by scope.
In one drain, the coordinator told the retro delegate that another scope drained the gate ideas.
The new gate rows came from the idea memo alone, so they missed the lessons of the FT445 retro and changelog.

The `## Delegate the evidence` section of `/bench-drain` now closes this gap.
The coordinator charges each read delegate to read a topic that another scope drains and to report the overlap.
Before retained batch authorship starts, the coordinator reads each new or merged row against every drained source on its topic.
The coordinator also reads the row against the shipped changelog entry of the build that each source follows.
A `CHANGELOG.md` entry under `[Unreleased]` records the change.

## Acceptance

- [ ] `## Delegate the evidence` tells the coordinator to charge each read delegate to read a topic that another scope drains and to report the overlap.
- [ ] `## Delegate the evidence` tells the coordinator to read each new or merged row against every drained source on its topic, not only the unit that opened it.
- [ ] `## Delegate the evidence` tells the coordinator to read each such row against the shipped changelog entry of the build that each source follows.
- [ ] The section adds no more than three sentences, and it restates no rule that another section owns.
- [ ] `CHANGELOG.md` has one entry under `[Unreleased]` that states the new drain behavior and the gap it closes.
