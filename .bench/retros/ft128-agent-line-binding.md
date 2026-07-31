# FT128 agent-line binding — implementation retro

Spec: `specs/ft128-agent-line-binding/spec.md` (9 stories, 37 coverage rows)
Shape: 5 build tickets + 1 reviewer-directed repair + 1 review-repair pass, 8 commits.

## The independent probe must differ in kind, not just in instance

The story-6 profile cross-check shipped proving token *membership* in the profile's
`Lines` section rather than *placement* in the right cell. The build delegate's
per-cell tests each deleted a token, and the coordinator's independent probe also
deleted a token — a different cell, but the same mutation kind. Both passed. A
semantic review swapped two cells instead, which leaves every token present, and
the gate stayed green while the profile contradicted the binding.

The delegation discipline already requires the coordinator to probe a behavior the
delegate did not author. That is necessary but not sufficient: a probe that varies
only the instance inherits the delegate's blind spot about which *class* of defect
is possible. For a quantifier claim, vary the mutation kind — delete, swap,
duplicate, reorder — before accepting the row.

## Ticket claims must be re-derived from the tree, not from the spec's account of it

The T3 ticket asserted three defects in `checkLineBinding`; the tree already
satisfied two of them, because the preceding ticket had gone further than the spec
anticipated. The delegate reported this honestly rather than staging a false red,
which is the right behavior — but the wasted framing came from the coordinator
writing ticket claims from the spec's description of the base instead of reading
the base at derivation time. Tickets are derived after earlier tickets land; their
account of "what is red today" ages fastest.

## The worktree porcelain assumes tickets land on the working branch

`bench worktree create` always roots at the default branch with no base-ref flag,
and `bench worktree release` refuses while an assignment branch has not landed
there. A build that chains tickets on an integration branch therefore cannot cut a
worktree at the chain tip, cannot advance a shared integration ref (`git branch -f`
is correctly blocked), and cannot retire any worktree until the reviewer merges.
The chain was formed instead by merging the previous ticket's assignment branch
into each new worktree — the manual form of a compare-and-swap integrate.

This is direct evidence for `decisions/spec-integration-gate-cadence.md`, whose
`bench spec build assign` / `checkpoint` / `integrate` family is exactly the missing
surface. Ticket #6's mandatory-parallelism rule was adoptable without any porcelain;
ticket #1's provisional-commit cadence was not, and remained the largest unclaimed
saving in the run.

## Two routes to the spec status transition are one too many

`bench spec implemented <spec>` and `bench commit --spec <slug>` both perform the
`Status: staged` → `implemented` transition. Running the former first makes the
latter fail with `no Status: staged line`, because each expects to own the flip.
The phase contract names `bench commit --spec` as the owner; the standalone
subcommand is a second author of the same transition.
