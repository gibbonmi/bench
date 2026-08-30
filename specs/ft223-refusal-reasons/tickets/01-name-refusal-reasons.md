# Name the operator action and the open reason in authorization refusals

Blocked by: none
Writes: internal/gate/authorization/authorization.go, internal/landing/landing.go, internal/landing/merge.go, internal/landing/landing_reviewed_test.go, internal/landing/state_test.go, internal/commit/dry_run_test.go, internal/worktree/land_specless_test.go, internal/status/status.go, internal/status/status_gatecache_test.go, ROADMAP.md, roadmap/FT223.md (deleted)

## What to build

`bench commit` refuses a composed tree with
`prospective authorization refused: <kind>`, where `<kind>` is an internal
`authorization.Kind`. The operator reads the kind and has no next action. This
ticket makes three refusal surfaces name what happened and what to do next.
The gate's own output does not change. The failures table and the `gate: red`
trailer stay exactly as the gate prints them, and the refusal line follows
them.

**The inherited refusal.** `authorization.Inherited` means the gate ran red on
the composed tree, and no green baseline marker proves the red belongs to this
diff. The refusal says that in operator terms and names `bench gate --fresh`
as the next action. The same rewording applies to `authorization.Candidate`.
There, the gate ran red on the composed tree and the green baseline attributes
the red to this diff. The next action is to read the failures above and fix
them.

Every refusal line keeps the literal prefix `prospective authorization refused: <kind>`
and adds the operator sentence after it. Two tests depend on that prefix.
`TestLandReviewedAuthorizationKindTablePreservesState` in
`internal/landing/landing_reviewed_test.go` asserts the kind word for every
non-green kind. `TestLandCommandSpecLessGateRefusalPublishesNothing` in
`internal/worktree/land_specless_test.go` asserts the prefix through the
landing's `refused{detail=...}` line. Both stay green as written.

**The infrastructure refusal.** `authorization.Infrastructure` currently
prints only the kind. The cause lives in the gate's inspection `Reason`, for
example `declared environment unavailable`, and never reaches the operator.
Add a `Reason` field to `authorization.Result`. `AuthorizeWithWriters` fills it
from the execution's inspection reason when the kind is `Infrastructure`, and
leaves it empty otherwise.

The landing refusal prints the reason after the kind, so the operator reads
`prospective authorization refused: infrastructure (declared environment unavailable)`
and the next action `bench doctor`. When the reason is empty, the refusal
prints the kind alone with its next action. The lane outcomes `lane pass` and
`lane fail` keep their current text; `internal/landing/lane_test.go` pins it.

One function in `internal/landing/landing.go` renders every refusal from a
`Result`. All three callers use it: the composed-tree and reviewed-tree
authorizations in `landing.go`, and the merge authorization in
`internal/landing/merge.go`. Do not add a second string builder.

**The stale status row.** After a landing, `bench status` prints
`stale (gated tree <a>, work tree <b>)` even when the two trees match. The
reason for the staleness stays hidden. That reason is `GateInfo.Reason`,
filled from `gate.Inspection.Reason`. Extend `staleGateDetailAction` in
`internal/status/status.go` so the detail ends with the reason when one is
present: `stale (gated tree <a>, work tree <b>; <reason>)`. When the reason
is empty, the detail keeps its current shape. The existing expectations in
`internal/status/status_gatecache_test.go` that pin the current shape stay
green because their fixtures carry no reason.

The landing retires the FT223 row: the index line in `ROADMAP.md` and the
detail file `roadmap/FT223.md` leave together. `ROADMAP.md` names FT223 only
on its index line and in no `## Dependencies` edge, so no other row changes.

## Acceptance

- [ ] An inherited refusal from `bench commit` prints the failures table and `gate: red` as the gate prints them, then one refusal line. That line says the gate ran red on the composed tree and that no green baseline attributes the red to this diff. It names `bench gate --fresh` as the next action.
- [ ] A candidate refusal prints one refusal line. That line says the gate ran red on the composed tree and that the green baseline attributes the red to this diff. It names the failures above as the thing to fix.
- [ ] An infrastructure refusal with a non-empty reason prints the reason in parentheses after the kind and names `bench doctor` as the next action.
- [ ] An infrastructure refusal with an empty reason prints the kind alone with its next action.
- [ ] `authorization.Result.Reason` is empty for every kind other than `Infrastructure`, and a test in `internal/gate/authorization/` shows it filled when the gate subject is open.
- [ ] `internal/landing/lane_test.go` stays green without a change to its `lane pass` expectation.
- [ ] A merge authorization refusal through `internal/landing/merge.go` prints the same rendered line as the composed-tree refusal for the same `Result`.
- [ ] `TestLandReviewedAuthorizationKindTablePreservesState` and `TestLandCommandSpecLessGateRefusalPublishesNothing` stay green without a change to their assertions.
- [ ] `bench status` prints `stale (gated tree <a>, work tree <b>; <reason>)` when the gate inspection carries a reason, and the current shape when it does not.
- [ ] `internal/commit/dry_run_test.go` asserts the operator-facing inherited text, and the delegate records that expectation red on the tree before the fix.
- [ ] `go test ./internal/gate/authorization/ ./internal/landing/ ./internal/commit/ ./internal/status/ -parallel 2` passes.
- [ ] `ROADMAP.md` carries no FT223 line and `roadmap/FT223.md` is absent.
