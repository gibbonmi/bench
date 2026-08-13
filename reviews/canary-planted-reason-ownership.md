# Canary planted-reason ownership review pickup

Exact candidate: `b048801e..78090113` plus the seven uncommitted Ticket 9 documentation paths.

## Standards

Finding count: 2

Worst issue: the candidate changed an explicitly excluded ownership area.

- **Auto-fix — P1:** `internal/gate/capability_skips.go:29` changes comment prose even though `specs/canary-planted-reason-ownership/spec.md:361` authorizes no `internal/gate` edit. Restore the pre-candidate comment; do not expand the fence.
- **Auto-fix — P3:** `projects/benchkit.md:200` leaves a 101-character line in a cold-session profile whose neighboring prose is wrapped consistently. Rewrap without changing meaning.

## Spec

Finding count: 1

Worst issue: the platform-reference rewrite removes a live evidence contract beyond DOC3's retired-canary scope.

- **Auto-fix — P2:** `.bench/BENCH-reference.md:121-123` removes the still-live selected/inherited conformance evidence and mixed-verdict contract and overstates the phase table as the complete description of execution. `specs/canary-planted-reason-ownership/spec.md:154-156` authorizes removing the authenticated inner-canary and canary-phase-selector claims only. Restore the current selected/inherited evidence behavior documented by `internal/gate/verdict.go:611-660`, `internal/gate/decision.go:49-63`, and `internal/status/status.go:251-257`, without restoring retired canary execution language.

## Coverage

Finding count: 0

Worst issue: none.

The axis enumerated 182 fixtures across 11 family partitions, four explicit `CHECK` bindings, the two owning-package release replacements, the complete spec write set, and the generic plus project hostile-input inventories. Its initial fixture-retirement concern was withdrawn as a no-op after checking the approved current spec and decision source.
