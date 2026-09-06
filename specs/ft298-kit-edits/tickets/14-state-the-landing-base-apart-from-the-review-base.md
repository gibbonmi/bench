# State the landing base apart from the review base

Blocked by: 13-walk-each-row-to-the-ticket-that-writes-its-seam.md
Writes: .agents/commands/bench-review-implementation.md
Covers: none

## What to build

Add one rule to `.agents/commands/bench-review-implementation.md` under `Process step 7`. The rule reads: The review base is the fold commit that merged `main` into the source. The landing base is that `main` tip itself. `bench worktree land --base` takes the `main` tip, and it refuses the fold commit. Write it in ASD-STE100, in the shape of its neighbouring items, and add no other sentence. Source: roadmap row FT298.

## Acceptance

- [ ] step 7 states the two bases in three sentences.
- [ ] the text names the `main` tip as the landing base.
- [ ] `bench gate-prose` passes on the edited file.
