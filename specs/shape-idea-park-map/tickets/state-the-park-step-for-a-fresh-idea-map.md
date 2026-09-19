# State the park step for a fresh-idea map

Blocked by: none
Writes: .agents/commands/bench-shape-idea.md, CHANGELOG.md
Covers: none

## What to build

The shaping command tells the agent to park a map that starts from a fresh
idea. Such a map has no roadmap row, and the command never edits the roadmap.
The agent runs `bench idea` with the map path in the text, so that the next
`/bench-drain` gives the map a row. The rule sits with the roadmap row rules,
in one place.

## Acceptance

- [ ] The roadmap section of the shaping command states that a fresh-idea map has no row.
- [ ] The same section tells the agent to run `bench idea` with the map path before the exit.
- [ ] A map that starts from a pulled roadmap row gets no park step.
