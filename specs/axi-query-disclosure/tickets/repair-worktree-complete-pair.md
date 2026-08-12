# Repair worktree complete pair

Blocked by: none
Writes: `internal/worktree/` (tests and testdata only)

## What to build

Cover the owned completed assignment at the public surface (review finding
R9): the existing terminal pair drives only a present foreign registration, so
a `listAssignmentRow` regression for `StateComplete` stays green. Add a public
old-to-new terminal pair that drives an owned completed assignment through
`ListCommand` and pins its rendered row and the honest empty help block.

## Acceptance

- [ ] [WC1] (covers QD6) a public `ListCommand` test creates an owned
  assignment, completes it through the production lifecycle, and asserts the
  checked-in old/new terminal pair: the `StateComplete` row's exact rendered
  values plus exactly one appended `help[0]{cmd,why}:` block and no actions.
- [ ] [WC2] (covers QD6) the pair names its state distinctly from the foreign
  terminal fixture — no shared oracle without an explicit alias naming both
  states.
