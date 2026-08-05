# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `0a85ea6` (hotfix merged fast-forward), clean tree, 3 unpushed commits
Spec: `specs/exact-prospective-landing/spec.md` (Status: staged; active lifecycle run)
Gate: green at `0a85ea6` (full fresh run in the hotfix worktree over this exact tree)

## State

**Phase reached: the out-of-fence-repair workflow hotfix is merged to main; the
stranded exact-prospective-landing run is ready to resume through the new route.**

The hotfix (`0a85ea6`, reviewer-merged; built in owned worktree request
`hotfix-oof-repair-route-2026-08-05`, branch
`bench/assign/9c5bf5813ba3aaab3a743cca5612dd5e/082e223dd6adb2e83151a524afff9b56`,
now releasable) adds
`bench spec build assign --refresh <debug-receipt>`: a delegate blocked outside its
ticket fence returns a validated debug receipt; the repair rides ticket-commit →
promote-recompose → assign/checkpoint/integrate; refresh re-bases the blocked
assignment onto the repaired candidate preserving in-fence bytes behind
`refs/bench/specbuild/refresh/<digest>`. Regression suites:
`go test ./internal/specbuild -run 'TestRefresh|TestAbandonRemains'` (hostile cases)
and `go test ./internal/contract/runtime -run TestRuntimeSpecBuildOutOfFenceRepair`
(compiled-binary dogfood trace incl. promote). On main the route is deterministically
red: `assign --refresh` exits 2 (unknown argument).

The active lifecycle run `869a2d33cab96f962882762a9ea56fd21c952976248e64fa917f6da6c48dbca3`
is preserved untouched: candidate `33a220bf99c5edee487553a8b3617eab16cf1eb1` integrated,
adapter assignment `daa1f4183dcb0dd159cca7b738d0fe7c` owned and uncheckpointed with its
in-fence `internal/commit/commit.go` edit intact. Its repro
(`go test -count=1 ./internal/contract/runtime -run 'TestRuntimeCommitContracts/fresh_green_verdict_is_reused' -v`)
is red only inside that worktree; cause confirmed: `authorization.Authorize` discards
gate output (`internal/gate/authorization/authorization.go:155`) and `landing.Request`
carries no writer contract — a repair outside the adapter's fence. Do not apply any
stale abandonment fingerprint; the abandon route is the positive control, not the plan.

Closed decisions: no ninth lifecycle verb (refresh extends `assign`); the debug-receipt
schema is owned by `internal/specbuild/refresh.go` validation; abandon remains the
escape hatch only.

## Next command

`/bench-implement-spec --full specs/exact-prospective-landing/spec.md`

— resuming the stranded run through the new documented route: land the
landing/authorization repair ticket, promote-recompose, assign/checkpoint/integrate it,
then `bench spec build assign exact-prospective-landing --ticket <adapter-ticket>
--request <original-request> --refresh <debug-receipt>` before resuming the adapter
delegate.

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
