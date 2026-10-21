# Share Git command projections

Blocked by: none
Writes: internal/shellcommand/projection.go (new), internal/shellcommand/projection_test.go (new), internal/gitguard/scan.go, internal/gitguard/projection_contract_test.go (new)
Covers: GP01, GP02, GP03, GP05, GP06, GP07, GP08, GP09, GP10, GP72, GP74

## What to build

Add the shared wrapper and worktree-child projection, and move the Git guard onto it with its current verdicts.
Add `ShellChildren` and `WorktreeExec` in a new `internal/shellcommand/projection.go`, as the spec's Projection API states.
Keep `Parse`, `ProjectCommandWords`, and `ResolveRoutinePrefix` unchanged. Do not edit `shellcommand.go`, which has 3 lines of headroom.
`ShellChildren` returns the matching `-c` children of `sh`, `bash`, or `zsh` in argv order, reads each argv position independently, and does not recurse.
`WorktreeExec` trusts the caller's Bench identity, returns the child argv after the first later `--` unchanged, and sets `DifferentCheckout` only for a nonempty child.

Migrate `scanWords` and `worktreeExecChild` in `internal/gitguard/scan.go` to these facts, then remove the private `wrappers`, `wrapperCFlagRe`, and exec-child grammar.
The Git guard keeps its shell-keyword trim, its first matching shell child, its one-level traversal, its Bench basename test before `WorktreeExec`, `FindSubcommand`, option parsing, and repository authority.
It passes `DifferentCheckout` into its existing `elsewhere` decision.
Keep the `scan` signature, so `gitguard.go:135` needs no edit.

Before this ticket changes source, the coordinator captures the pre-migration reference executable that GP46 uses.
GP2 receives the reviewed projection API with all Git rows green, while its own consumers are still unbuilt.
`internal/gitguard` holds 10 of its 12 allowed files; this ticket adds one, so add no other new file there.

## Acceptance

- [ ] `TestResolveRoutinePrefix` keeps its executable indexes (GP01).
- [ ] `ShellChildren` returns one ordered child for each supported flag in family W, including `bash -c -c 'git status'`, and gives no child for an invalid index or a missing operand (GP02, GP03).
- [ ] `WorktreeExec` returns quoted child argv byte for byte, gives no child without `--`, and sets `DifferentCheckout` only for a nonempty child (GP05–GP08).
- [ ] `TestClassifyVerdicts` keeps its one-level nested wrapper and exec verdicts (GP09).
- [ ] A child push and a redirected push keep the unresolved-destination label, and omitting `DifferentCheckout` propagation makes `TestGitProjectionAuthority` red (GP10).
- [ ] `mybench worktree exec X -- git reset --hard` keeps the empty Git label (GP72).
- [ ] `bash -c 'echo ok' -c 'git reset --hard'` keeps the Git allow verdict, because Git reads only the first matching child (GP74).
- [ ] Omitting a supported shell basename or flag form makes `TestShellChildren` red, and byte-identical restoration returns the same green.
- [ ] No private wrapper or exec-child grammar remains in `internal/gitguard`, and each new file stays within 400 lines.
