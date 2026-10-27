# Show background outcomes

Blocked by: 07-start-a-background-grade-after-a-lane-pass.md
Writes: internal/status, internal/commit/commit.go, internal/commit/background.go (new), internal/commit/background_test.go (new), internal/gate/backgroundstate (new), CHANGELOG.md
Covers: BG84, BG85, BG99

## What to build

Show the author each background outcome where the author already looks.

- `bench status` shows the latest background outcome of each worktree, with its log path. The row obeys the dashboard's five-row budget and its severity ladder (`internal/status/status.go`).
- The next `bench commit` in a worktree prints the latest background outcome of that worktree and its log path, before it starts the next grade.

A background record is never a gate verdict, so neither surface reads it as green.
The wording of both lines is the author's choice.
`commit.go` has 400 lines, at its budget, so the print goes in `internal/commit/background.go`.

## Acceptance

- [ ] `bench status` shows the latest background outcome of each worktree.
- [ ] The next `bench commit` prints the latest background outcome and its log path.
- [ ] A worktree with a red background outcome and a green gate cache still shows the gate cache as its verdict.
