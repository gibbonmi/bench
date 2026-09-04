# Resolve the bare push destination

Blocked by: none
Writes: internal/git/push_destination.go (new), internal/git/push_destination_test.go (new)
Covers: PG37, PG38

## What to build

Verify the premise first: read `ResolvedDefault` in internal/git/default_branch.go
and `CheckedOutBranch` in internal/git/status.go. Read the `Output` and `OK`
probe helpers in the same package. `CheckedOutBranch` returns the literal
`HEAD` on a detached head, so the new code tests that literal.

Add `BarePushDestination(root string) (string, bool)` in a new file. It returns
the branch that a bare `git push` targets. Under `push.default` `simple`,
`current`, or an unset value, it returns the checked-out branch. Under
`upstream` or `tracking`, it returns the upstream branch name. Under `matching`
or `nothing`, on a detached `HEAD`, outside a repository, or after a probe
error, it returns `("", false)`.

The function shares one contract with two siblings. The classifier ticket adds
the `Checker` field `BareDestination func() (string, bool)`. The guard ticket
wires this function into that field. This ticket writes no caller.

## Acceptance

- [ ] `BarePushDestination` returns `topic` for a checked-out `topic` when `push.default` is unset.
- [ ] It returns the upstream branch name when `push.default` is `upstream`.
- [ ] It returns the upstream branch name when `push.default` is `tracking`.
- [ ] It reports no destination on a detached `HEAD`.
- [ ] It reports no destination under `matching` and under `nothing`.
- [ ] It reports no destination for a directory outside a git repository.
- [ ] Self-probe: return the checked-out branch under `matching`, and report the no-destination test red.
