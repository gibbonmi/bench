# Test the hook input verdicts in harness and shellcommand

Blocked by: none
Writes: internal/harness/worktree_test.go (new), internal/shellcommand/shellcommand_test.go
Covers: none

## What to build

The quality survey of 2026-09-29, card 07, found two hook-layer packages with untested verdicts for malformed input. No test runs `internal/harness`, so its statement coverage is 0%. The `worktree-lifecycle.sh` hook passes each Claude worktree event to `harness.WorktreeCommand`. The hook states that a malformed event must refuse and must not bypass the lifecycle. In `internal/shellcommand`, no test runs the fallback split that `Parse` uses when the lexer finds an unbalanced quote. The Git and Bench guards read the `Parse` stream, so the fallback must keep each later command visible.

This ticket adds tests only. It changes no production code.

Add `internal/harness/worktree_test.go`. Each case calls `WorktreeCommand` and checks the exit code, the stderr verdict, and an empty stdout:

- An unknown action, a missing action, or an extra argument exits 2 with the usage line.
- An event of 1 MiB plus one byte exits 1 with the size-limit verdict, for `create` and for `remove`.
- An event of exactly 1 MiB passes the size limit, and the field check gives the verdict.
- An empty event or a whitespace-only event exits 1 with the empty-input verdict.
- An event that is not JSON, has trailing data, or has a field of the wrong type exits 1 with the invalid-JSON verdict.
- An event with no required field, a blank name, or a JSON `null` exits 1 with the required-field verdict.
- A `create` event whose `cwd` is outside a Git repository exits 1.
- A `remove` event whose repository hint is outside a Git repository exits 1. The hint is `CLAUDE_PROJECT_DIR` when it is set, and `worktree_path` when it is empty.
- A `remove` event from a session that owns no assignment reaches the release, and the release refuses it with exit 1.

Extend `internal/shellcommand/shellcommand_test.go` with `Parse` cases for an unbalanced input:

- An unterminated single quote before `&& git push --force`.
- An unterminated double quote before a newline and `git push`.
- A trailing backslash.

For each case, the test compares the full token stream and the simple-command spans with a literal. The literal shows that a later command stays a separate simple command.

## Acceptance

- [ ] `go test ./internal/harness/` runs each `WorktreeCommand` case above, and each case passes.
- [ ] `go test ./internal/shellcommand/` runs each unbalanced-input case above, and each case passes.
- [ ] A `bench probe` on the size-limit verdict in `internal/harness/worktree.go` turns the harness tests red.
- [ ] A `bench probe` on the fallback split in `internal/shellcommand/shellcommand.go` turns the shellcommand tests red.
- [ ] No production file changes.
