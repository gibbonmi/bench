# Migrate Bench projection consumers

Blocked by: 1-share-git-command-projections.md, 3-migrate-bench-and-write-readers.md
Writes: internal/benchguard/benchguard.go, internal/benchguard/reviewer.go, internal/benchguard/projection_contract_test.go, internal/conformance/refusal_route_guard_test.go, cmd/bench/guard_projection_test.go
Covers: GP04, GP11, GP12, GP13, GP62, GP63, GP64, GP65, GP66, GP68, GP69, GP70, GP73

## What to build

Move every Bench-side wrapper and exec-child reader onto `ShellChildren` and `WorktreeExec`, then retire the private Bench recognizers.
Ticket 1 supplies `ShellChildren` and `WorktreeExec`. Ticket 3 supplies `benchguard.go` with its reader migrated, and the two test files `benchguard/projection_contract_test.go` and `cmd/bench/guard_projection_test.go`, which this ticket extends.
In `benchguard.go`, `judge` and `spanInvokesBench` read wrapper children through `ShellChildren` and keep their search across every matching child, one level deep.
`judge` uses `WorktreeExec.Head` only after its own `isBench` test. Bench keeps its resolver, `judgeExec`, and the pool-reference guard.

In `reviewer.go`, `ReviewerVerb` keeps its recursive walk, as the reviewer decided on 2026-10-20.
At each level it calls `ShellChildren` and `WorktreeExec` and recurses into each child stream or child argv.
It keeps `isBench`, `namesVerb`, and the resolver, and it applies `isBench` before `WorktreeExec`.
An exec head without `--` gives no child and no verb match.
After both files migrate, remove `isWrapper`, `isExecHead`, and `wrapperFlag` in one step, because both files use them.

In `refusal_route_guard_test.go`, `runsBenchUnderExec` keeps its Bench identity test and its agent-route scope, and reads the exec child through `WorktreeExec`.
The core entry order stays pool reference, then reviewer verb, then `Classify` (`cmd/bench/guards.go:170-178`); do not edit `guards.go`.
`benchguard.go` must end at 400 lines or fewer.

## Acceptance

- [ ] Bench `Classify` refuses the follow-on in the second `-c` child of W4 (GP04).
- [ ] `TestClassifyResolvesBareAliasFromProcessPath` and `TestClassifySpanScopedFollowOns` stay unchanged and pass (GP11, GP12).
- [ ] A pool reference with a follow-on emits the pool refusal (GP13).
- [ ] `ReviewerVerb` refuses the nested-wrapper and nested-exec override commands (GP62, GP63).
- [ ] `TestGuardDeniesTheOverrideVerb` and `TestFollowOnGuardDeniesTheOverrideEnvelope` stay unchanged and pass (GP64, GP70).
- [ ] A pool reference with the reviewer verb emits the pool refusal, and the reviewer verb with a follow-on emits the reviewer-verb refusal (GP65, GP66).
- [ ] `ReviewerVerb` allows `x worktree exec t -- bench commitment override` (GP73).
- [ ] `TestAgentRouteGuardCheckBites` keeps its RR12 finding, and `runsBenchUnderExec` reads the child only through `WorktreeExec` (GP68, GP69).
- [ ] Replacing the reviewer walk's recursion into a wrapper child with a one-level check makes `TestReviewerVerbKeepsRecursiveReach` red.
- [ ] No `isWrapper`, `isExecHead`, or `wrapperFlag` remains in `internal/benchguard`.
