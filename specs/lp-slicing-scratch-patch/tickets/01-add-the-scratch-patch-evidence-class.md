# Add the scratch-patch evidence class

Blocked by: none
Writes: .agents/skills/bench-craft-tickets/references/slicing-checks.md
Covers: none

## What to build

A slicer reads a seventh fence-expansion evidence class in the slicing checks: a test file that goes red in a scratch-patch test run. The run applies to a ticket that adds a refusal or retires a rule, a fixture shape, or a row. Before the slicer locks that ticket, the research delegate applies the change as a scratch patch in its own Bench worktree. The delegate runs the reached packages once with one `bench test` call and returns each red file. The existing keep rule then puts each red file in the ticket's `Writes:` line. The delegate discards the patch and lands nothing.

## Acceptance

- [ ] The fence-expansion evidence list holds a seventh class for a test file that goes red in a scratch-patch test run.
- [ ] The slicing checks state when the run applies, what the research delegate runs, what it returns, and that it lands nothing.
- [ ] The prose, docs-currency-workflow, guidance-prose-budgets, and skills-index-command-adapters checks pass.
