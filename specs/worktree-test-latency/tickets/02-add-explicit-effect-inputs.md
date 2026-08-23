# Add explicit worktree effect inputs

Blocked by: none
Writes: internal/worktree/

## What to build

Add explicit root, home, time, environment, directory, and executable inputs at
the worktree effect boundary. Keep public commands as thin process-context adapters.

Provide temporary compatibility forms for callers that still need migration.
Do not change command grammar, rendered output, or exit status.

## Acceptance

- [ ] EI1: Lower owners receive resolved values without reading environment or current directory.
- [ ] GF1: The ordinary gate driver and `-count=1` remain unchanged.
- [ ] Public command behavior remains byte-compatible for unchanged scenarios.
