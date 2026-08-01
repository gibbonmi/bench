# Refuse a reduced verdict for release evidence

Blocked by: Expand the verdict record with the reduced shape

Ownership fence: `internal/preprelease`, `internal/contract/surface/preprelease`
Assumptions: `bench prep-release` already refuses without a current dev-green verdict, so this is a fail-closed default added at that existing precondition seam

## What to build

`bench prep-release` refuses a reduced verdict and names that as the reason. The
release path is the first consumer that must not accept a partial verdict, and
today's precondition asks only whether the recorded status is green — which a
reduced record satisfies while answering for a fraction of the tree.

The refusal message carries the reason rather than a bare failure, because the
maintainer's next action depends on knowing that the verdict was narrow rather than
absent: the fix is `bench gate --fresh`, not a rerun of whatever they just did. Any
other release-evidence precondition that leans on whole-tree green gets the same
treatment, so the fast path cannot leak into release authority through a second
door.

## Acceptance

- [ ] [R21] `bench prep-release` refuses when the current verdict is reduced and its message names the reduced verdict as the reason.
