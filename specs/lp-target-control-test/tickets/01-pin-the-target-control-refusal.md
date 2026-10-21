# Pin the target control-character refusal

Blocked by: none
Writes: internal/worktree/identifier_operand_test.go
Covers: none

## What to build

The worktree target resolver refuses a target that contains a control character. The
refusal reads `target contains control characters`, and the resolver applies it after the
ledger read for each target verb. No test in `internal/worktree` gives such a target to a
target verb, so a removal of the refusal stays green.

Add one control-character case to the resolver refusal cases that
`TestTargetVerbsNameTheResolverReason` runs. That test gives each case to the path, exec,
show, and build verbs through their public entries. Change no production code.

## Acceptance

- [ ] Each of the path, exec, show, and build verbs refuses a control-character target with exit 1, the line `target contains control characters`, and the route `next=bench worktree list`.
- [ ] A `bench probe` that omits the control-character check from the resolver turns the new case red.
