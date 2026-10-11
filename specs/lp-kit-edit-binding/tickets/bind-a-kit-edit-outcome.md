# Bind a kit-edit outcome through a one-ticket folder

Blocked by: none
Writes: .bench/BENCH.md
Covers: none

## What to build

The next-outcome rule in `.bench/BENCH.md` gives a deliverable binding route only to a `Next: ticket` outcome and to a spec-stage close.
On 2026-10-11 the reviewer decided that a `Next: kit-edit` outcome uses the same route as a `Next: ticket` outcome.
When that outcome becomes eligible and has no bound deliverable, the agent writes and lands a one-ticket tickets-only folder.
Then the agent binds that folder through `bench commitment plan` and `bench commitment approve`, and cites the closing phase's approval as the decision.
The committed outcomes FT380, FT403, and FT379 have `Next: kit-edit` rows and need this route.

## Acceptance

- [ ] The next-outcome rule names `Next: kit-edit` beside `Next: ticket` as a row whose outcome writes and lands a one-ticket tickets-only folder first.
- [ ] The rule binds that folder through `bench commitment plan` and `bench commitment approve`, with the closing phase's approval as the decision.
- [ ] The rule keeps the spec-stage close binding, the binding-only limit of the plan, and the reviewer direction for a reorder or a new admission.
- [ ] No other file restates the rule.
