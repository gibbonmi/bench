# Run excludable phases against a stripped worktree

Blocked by: Declare the reduced gate scope

Ownership fence: `internal/gate/stripped_worktree.go`, `internal/gate/runner.go`, `internal/gate/phases.go`, `internal/contract/runtime/runtime_gate_stripped_test.go`
Assumptions: the gate owns no working-tree materializer today — `.bench/gate-prospective.sh` is a build-then-exec hook and `internal/gate/prospective.go` only hashes it — so this ticket adds one rather than composing an existing seam

## What to build

On every full gate run, the excludable phases execute against a materialized
worktree of the current tree with the allowlisted paths removed, and they execute
there with capabilities required. Included phases keep running against the real
root, so the checks that exist to grade the allowlisted files keep grading them.

Both halves are load-bearing and the second is the one that is easy to drop.
Stripping alone is not enforcement: this tree's idiom for a missing subject file is
`skipIfSubjectFileMissing`, which turns absence into a capability skip, and a
capability skip is informational in the dev tier. A phase whose subject file
vanished would therefore go permanently green rather than red — worse than no
feature at all, because it degrades every run instead of only the capture-confined
ones. Required capabilities is what converts that silent green into a red. This was
the first draft's fatal error; do not simplify it away.

Materialize through git so the result is a real repository — a contract test stages
its subject with `git ls-files` and fails outright against a plain directory copy —
and run the build phase inside it so its `dist/` output is produced there exactly as
on the primary root. Concurrent runs use distinct temporary roots, as the worktree
owner already requires, and an orphaned stripped worktree is retired by the existing
recovery path rather than by new cleanup.

## Acceptance

- [ ] [R08] A planted excludable phase that hard-reads an allowlisted path reds a full gate.
- [ ] [R09] A planted excludable phase that soft-skips on the missing path reds a full gate rather than reporting an informational skip.
- [ ] [R10] Included phases still find the allowlisted paths on that same run.
- [ ] [R11] The materialized stripped tree is a git repository, so a subject staged through `git ls-files` succeeds against it.
