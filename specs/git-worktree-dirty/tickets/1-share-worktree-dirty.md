# Share the worktree dirtiness query

Blocked by: none
Writes: internal/git/status.go, internal/git/facts_test.go, internal/preflight/gather_inputs.go, internal/worktree/lifecycle.go, internal/shift/loop.go
Covers: none

## What to build

Promote the repeated porcelain dirtiness query into the existing Git owner for FT302.
The current census finds three identical queries in preflight, worktree reuse, and shift startup.
Each caller keeps its current error handling and user response.
The query keeps Git's existing configuration and untracked-file policy.

## Acceptance

- [x] One Git helper supplies the dirtiness fact to all three callers.
- [x] Clean trees and ignored files remain clean; tracked changes and visible untracked files report dirtiness.
- [x] Staged changes and deletions report dirtiness, including paths with spaces and newlines.
- [x] Query errors reach each caller without a new fallback or diagnostic.
- [x] Existing preflight, worktree, and shift tests pass without changed expectations.

## Verification

Run the affected package tests and root conformance.
Omit the status query and require the helper test to fail on an untracked file.

The affected selection passed all 52 packages, including root conformance.
The probe replaced the status query with an empty result and returned `bit`, with one failed test.
`TestWorktreeDirty` reported `false` where the visible untracked file required `true`.
The probe restored the source.
