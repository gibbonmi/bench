# Learnings — usage journal

<!-- entries below -->


## 2026-08-24 — a ticket premise can go stale between spec and build

What happened: ticket 13 named a 795-line file that had held 209 lines since before the spec's base. The build stopped as a material acceptance shortfall, and the reviewer had to re-decide scope.
Right behavior: fail a stale premise before a charge spends tokens.
Proposed rule change: `bench preflight build` compares each ticket's stated line count against the tree and reds on a mismatch.

## 2026-08-24 — mint the worktree request token before the create command

What happened: the create command embedded `$(date +%s)`, and a second expansion printed a different value. The exact token was unrecoverable, and the landing had to try candidates.
Right behavior: mint the token into a variable first, then pass and record that one value.
Proposed rule change: `bench worktree create` echoes the request token back in its output.

## 2026-08-24 — `bench spec retire` dirties a checkout that cannot commit

What happened: the retire verb ran on the primary checkout and deleted the spec folder there, but `bench commit` refuses the primary. The discard hook then blocked the cleanup, and the state needed a stash.
Right behavior: run the retire inside a bench worktree from the start.
Proposed rule change: `bench spec retire` refuses the primary checkout the same way `bench commit` does.
