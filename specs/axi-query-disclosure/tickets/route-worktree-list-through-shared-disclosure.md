# Route worktree list through shared usage and disclosure

Blocked by: record-ft173-log-leverage.md
Writes: `internal/worktree/`, `internal/usage/`

## What to build

Route `worktree list` through `usage.Parse` and append actions derived from every rendered row: active rows carry both `bench worktree path <id>` and `bench worktree exec <id> -- <command>`, while orphaned rows carry `bench worktree clean <path>`. Preserve row order and exact values, deduplicate exact templates, and keep terminal or empty results honestly action-free.

## Acceptance

- [ ] [WH1] (covers QD1) every active row yields both path and exec actions with its id, every orphaned row yields its own clean action with its path, and many-row fixtures reject first-match-only, sampled, reordered, or guessed derivation.
- [ ] [WH2] (covers QD1) empty and terminal worktree-list results append the honest zero-row help block.
- [ ] [WH3] (covers QD2) `--help`, `-h`, and bare `help` return shared usage on exit 0; an old/new argv matrix proves every other accepted or rejected spelling, stream, and exit is byte-equal.
- [ ] [WH4] (covers QD6) old-to-new public-command fixtures prove that each named worktree-list state changes only by its appended help block plus the WH3 help-spelling delta; primary bytes, streams, exits, and all other argv behavior remain byte-equal.
