# Run a supplied check before the commit publishes

Blocked by: none
Writes: internal/landing/landing.go, internal/landing/check.go (new), internal/landing/attribution.go, internal/landing/check_test.go (new), internal/commit/commit.go, internal/commit/check.go (new), internal/commit/check_test.go (new), internal/commit/refusal_route_test.go, internal/refusalroute/faces_commit.go, internal/refusalroute/routetest/routetest.go, internal/conformance/refusal_route_guard_test.go, internal/commit/merge_parent_test.go, internal/commit/landing_test.go
Covers: PP10, PP11, PP12, PP13, PP14, PP15, PP16

## What to build

Give the landing one optional pre-publication check, and give the commit verb one entry that supplies it. No production caller supplies a check in this ticket. Ticket 03 binds the chain.

`landing.Request` gains the field `Check func(candidate string) error`. `Owner.Land` calls a non-nil check after `commit-tree` and before `updateRef` (`internal/landing/landing.go:184-193`). The check receives the candidate commit.
A check error returns `landing.CheckRefusal`. Land then moves no ref and runs no reconcile. A nil check keeps every current step, and `DryRun` never calls the check.
`internal/landing/landing.go` has 28 lines of headroom under the 400-line limit. If the seam does not fit, put the refusal type in `internal/landing/check.go`.

The commit package exports `RunChecked(args, stdout, stderr, check)`, where `check` has the signature `func(root, candidate string) error`. `Run` keeps its signature and calls `RunChecked` with a nil check, so the callers of `Run` do not change.
`landingFace` maps `landing.CheckRefusal` to a new face `commit-chain-red`, and `refuse` prints it at exit 1 with the caller's chain flags in the rerun. If `internal/commit/commit.go` passes its line limit, put `RunChecked` in `internal/commit/check.go`.

Declare `commit-chain-red` in `internal/refusalroute/faces_commit.go` beside `commit-red`: Verb Commit, Authority Agent, and a route of one instruction to repair each red chain step, then `commitRerun`.
`TestCommitFacesFollowTheirRoutes` requires a producing fixture for each commit face (`internal/commit/refusal_route_test.go:128`, `:252`). Add the fixture for the new face. That file has 17 lines of headroom, so put a long fixture body in `internal/commit/check_test.go` and register it from the fixture list.
The refusal-route guard grades every Agent route (`internal/conformance/refusal_route_guard_test.go:34-105`). The new route has the same step shape as `commit-red`, so the guard must stay green.

## Acceptance

- [ ] PP10: inside the check, the destination ref still names the expected base, and the candidate object exists.
- [ ] PP11: the candidate tree equals the composed tree, and its first parent equals the expected base.
- [ ] PP12: a check error returns `landing.CheckRefusal`, and the destination ref and the checkout fingerprint stay unchanged.
- [ ] PP13: `RunChecked` with a red check exits 1, prints `next=` with the caller's `--preflight-build` rerun, and leaves HEAD unchanged.
- [ ] PP14: a lane red with a check prints the lane refusal, and the check never runs.
- [ ] PP15: `TestLandComposesOnlyNamedPathsAndCASesExpectedBase` stays green with no check.
- [ ] PP16: a check that moves the destination ref leaves the winner in place, and Land returns the tip-moved refusal.
- [ ] `TestCommitFacesFollowTheirRoutes` and the refusal-route guard stay green with the new face.
