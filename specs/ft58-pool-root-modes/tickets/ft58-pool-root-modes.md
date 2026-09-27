# Refuse a pool root that Bench cannot tighten

Blocked by: none
Writes: internal/worktree/pool_root.go (new), internal/worktree/pool_root_test.go (new), internal/worktree/lifecycle.go, internal/worktree/ownership.go, internal/worktree/lifecycle_acquire_test.go
Covers: none

## What to build

Bench selects and creates the pool root for a pooled acquire and for an owned `bench worktree create`. Today, both paths ignore a failure to set the root to mode 0700, and both paths follow a pool root that is a symlink.

Both paths must prepare the root through one owner. That owner refuses a root that is a symlink or that is not a directory. It sets the mode of each root to 0700 through the `chmodPool` join, and it refuses the root when that operation fails. The refusal uses the existing lifecycle refusal grammar and names the pool root and the wanted mode. The operation that needs the root stops before it puts a checkout in the root.

## Acceptance

- [ ] When the `chmodPool` join fails for the pool root, the acquire and the create each refuse. The refusal names the root and mode 0700, and the root stays empty.
- [ ] When the pool root is a symlink, the acquire and the create each refuse it. The mode and the contents of the symlink target do not change.
- [ ] When the pool root exists with mode 0777, the acquire and the create each set it to 0700. When that change fails, the result is the refusal of the first row.
- [ ] When the pool root exists at mode 0700, a second acquire and a second create each succeed. The root stays the same directory at mode 0700.
