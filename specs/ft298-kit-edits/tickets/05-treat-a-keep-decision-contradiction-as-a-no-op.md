# Treat a keep decision contradiction as a no op

Blocked by: 04-probe-the-fixture-source-on-the-coverage-axis.md
Writes: .agents/skills/bench-craft-review/references/finding-discipline.md
Covers: none

## What to build

Add one rule to `.agents/skills/bench-craft-review/references/finding-discipline.md` under `When a ticket already decided`. The rule reads: A finding that contradicts a ticket's explicit keep decision is a `no-op`. The coordinator cites the ticket line in the disposition. Write it in ASD-STE100, in the shape of its neighbouring items, and add no other sentence. Source: roadmap row FT298.

## Acceptance

- [ ] a new `## When a ticket already decided` heading holds the one list item.
- [ ] the item names the `no-op` disposition and the ticket line citation.
- [ ] `bench gate-prose` passes on the edited file.
