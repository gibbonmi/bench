# Take an explicit root in four verbs

Line: opus / medium.

Blocked by: none
Writes: internal/worktree/worktree.go, internal/worktree/list.go, internal/worktree/pool_reclaim.go, internal/worktree/resume.go, cmd/bench/main.go, internal/sessioninspect/sessioninspect.go, cmd/bench/command_registry_test.go, internal/worktree/clean_operand_test.go, internal/worktree/list_actions_test.go, internal/worktree/pool_reclaim_test.go, internal/worktree/resume_test.go

## What to build

An operator runs `clean`, `reclaim`, `list`, and `resume clean` from any
directory inside the repository. The four verbs print the same bytes and return
the same exit codes as before.

Give `CleanCommand`, `ReclaimCommand`, `ListCommand`, and `ResumeCleanCommand` a
leading `root string` parameter. Compose the shape that `LandCommand` and
`ReleaseCommand` already use. Update the five call sites in `cmd/bench/main.go`
and `internal/sessioninspect/sessioninspect.go`. Each caller resolves the root
once at the command boundary, as the other seven verbs' callers do.

The tests of these four verbs pass an explicit root. They no longer change the
working directory for that reason.

## Acceptance

- [ ] WF09 pins `CleanCommand` on an explicit root while the working directory is a temporary directory.
- [ ] WF10 pins `ReclaimCommand`, `ListCommand`, and `ResumeCleanCommand` on an explicit root in the same way.
- [ ] WF11 pins the bytes of `bench worktree clean --help` and `bench worktree list` from a subdirectory.
