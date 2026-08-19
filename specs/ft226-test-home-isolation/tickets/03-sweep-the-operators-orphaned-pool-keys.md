# Sweep the operator's orphaned pool keys once

Blocked by: 01-bind-the-reauthorize-fixture-home.md
Writes: specs/ft226-test-home-isolation/ (verification log and this ticket's checkboxes only)

## What to build

A throwaway plan-before-apply script in the session scratchpad — never
committed — that sweeps `$HOME/.bench/worktrees` once. A target is a key
matching `^001-[0-9]+$` whose every child directory holds a regular `.git` file
with a `gitdir:` target that does not exist, and which holds nothing else at its
top level. The plan step prints the target count, five sampled targets with
their `gitdir:` lines, the non-target count, and every surviving key name that
does not match `001-<digits>`; the apply step removes only the listed targets
and prints the remaining key count. The counts, sample, and survivor list go
into the spec's verification log. This ticket writes nothing in the tree beyond
that log and its own checkboxes. Spec rows SW1, SW2, SW3.

Runs after ticket 01 so the gate runs that land this spec do not refill the
pool. The reviewer's sign-off on the spec is the approval for the destructive
apply step; the plan output is still shown before apply.

## Acceptance

- [ ] The plan prints before any removal: target count, five sampled targets
      with dangling `gitdir:` lines, non-target count, surviving non-`001`
      key names (`bench-…`, `project with spaces`, the dogfood keys are among
      them).
- [ ] A key whose child `.git` is a directory, a key with an extra top-level
      entry, a key whose pointer target exists, and a key not matching
      `001-<digits>` are all absent from the target list (each either observed
      in the real pool or planted in a scratch copy of it before the real run).
- [ ] After apply, the remaining key count equals the pre-apply count minus
      the target count, and every surviving key name from the plan is present.
- [ ] One subsequent `bench gate` run leaves the sorted `$HOME/.bench/worktrees`
      key listing byte-identical (SW3).
