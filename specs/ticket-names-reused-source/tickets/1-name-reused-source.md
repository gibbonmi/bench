# Name the reused source in a ticket that repeats a rule

Blocked by: none
Writes: .agents/skills/bench-craft-tickets/references/slicing-checks.md
Covers: none

## What to build

A slicer reads one slicing rule for a ticket that repeats an existing rule, such as an adapter, a validator, or a predicate. The rule tells the slicer to name the source of that rule in `What to build` and its package in `Writes:`. The author then extends that one source and does not paste a copy. This rule replaces the narrower second-caller rule, so the slicing checks hold one rule for this case.

## Acceptance

- [ ] The slicing rules hold one rule that names the source of a repeated rule in `What to build` and its package in `Writes:`.
- [ ] The slicing rules no longer hold the separate second-caller rule.
- [ ] The prose, prose-mechanics, and guidance-prose-budgets checks pass.
