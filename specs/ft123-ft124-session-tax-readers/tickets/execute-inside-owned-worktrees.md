# Execute inside owned worktrees

Blocked by: Resolve owned worktree targets

## What to build

Sessions can run a direct child argv in a resolved owned worktree with transparent
streams, environment, signals, and exit status, and can discover both worktree
commands through every shipped help and routing surface.

## Acceptance

- [x] `bench worktree exec` requires `--`, invokes argv without a shell, and
  preserves stdin, environment, cwd, stdout, stderr, and exact child exit status.
- [x] Non-TTY stdin and exits 0, 1, 2, 37, and 130 pass through; SIGINT leaves no
  descendant and does not mutate worktree state.
- [x] General and worktree help name `path`, `exec`, and the worktree-local
  `bash bin/bench.sh gate --fresh` route.
- [x] Compiled routing, real-kit and linked-repository wrappers, runtime
  contracts, CLI inventory, and the subcommand registry agree on the worktree
  surfaces.
