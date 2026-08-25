# Review pickup: lifecycle-refusal-names-component

Frozen base `45626d91`, reviewed tip `ac8221b3`. Raw findings: 6. Repair targets after de-duplication: 5.

## Standards

Count: 3. Worst: the active-state predicate is derived twice across the ticket fence.

- `internal/worktree/path.go:56` re-derives `selected.State != intent.StateActive` where ticket 01 named that fact `landingActiveState` (`identity_component.go:89`). Cites AGENTS.md, "Two derivations of the same fact must collapse into one source." Disposition: auto-fix.
- `CONTEXT.md:54` holds a 31-word descriptive sentence. Cites `ste-prose.md`, the 25-word bound for a descriptive sentence. Disposition: auto-fix.
- `identity_component.go:25` `recovers` has no production consumer. Refuted: the spec's Implementation decisions require each entry to carry "whether it carries a recovery command", and the registry test consumes the field. Disposition: no-op.

## Spec

Count: 1. Worst: the bundle validator changed its precedence for a double fault.

- `internal/worktree/ownership.go:306` `validateCreationBundle` runs the branch-checkout check before the owner-ID and marker-path comparison. A wrong `OwnerID` plus a detached HEAD now yields `assignment branch is not checked out`; before the diff it named the owner marker. No spec line permits the change, and the Problem section names this failure mode. Disposition: auto-fix. Repair predicate: restore the old order (marker, branch, registration, lock) with one source of each predicate. Split `identityBundleRefusal` into a marker step and a registration-and-lock step, and compose both callers from those steps.

## Coverage

Count: 2. Worst: the two-component precedence the edge inventory promises has no test.

- Edge inventory, "Two components fail at once inside one bundle; the first in registry order is named." Every fixture in `identityComponentFixtures()` mutates one dimension. Input: a re-pointed registration plus an unlocked worktree. Expected: `worktree registration does not match assignment <id>`. Test that should exist: one landing case that applies two fixture mutations and asserts the earlier component. Disposition: auto-fix.
- `internal/worktree/ownership.go:396` wraps the bundle refusal with `; checkout retained` for `bench worktree release`. The diff changed that sentence (the assignment id joined it), and no test pins it; LR19 covers only the unknown-request case. Input: a release of an assignment with a rewritten owner marker. Expected: `owner marker does not match assignment <id>; checkout retained` on stderr. Disposition: auto-fix.
