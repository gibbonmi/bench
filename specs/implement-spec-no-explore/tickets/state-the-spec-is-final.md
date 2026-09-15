# State that the spec and ticket are final in the Build section

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md
Covers: none

## What to build

Newer agents tend to re-derive a solution during implementation, even when the
spec and the ticket already settled the approach. The `/bench-implement-spec`
command file states the ticket work in its "Build" section. The file does not
tell the retained author that the spec is final. The file does not tell the
author to implement the ticket as written instead of weighing other
approaches.

The "Build" section gains one sentence, placed after its first paragraph's
statement of ticket order and TDD seams. The sentence states that the spec
and the ticket are final. It states that the retained author does not
evaluate other approaches. It also states that the author implements the
ticket as written, runs its focused checks, and stops. No other section
changes.

## Acceptance

- [ ] The Build section of `.agents/commands/bench-implement-spec.md` states that the spec and the ticket are final and that the retained author does not evaluate other approaches.
- [ ] The same sentence states that the author implements the ticket as written, runs its focused checks, and stops.
- [ ] `bench gate-prose .agents/commands/bench-implement-spec.md` exits 0.
