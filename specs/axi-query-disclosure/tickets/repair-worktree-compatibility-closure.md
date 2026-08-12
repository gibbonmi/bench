# Repair worktree compatibility closure

Blocked by: none
Writes: `internal/worktree/list.go`, `internal/worktree/list_actions_test.go`, `internal/worktree/testdata/pre-disclosure-argv-pairs.json`

## What to build

Preserve the pre-migration `worktree list` bytes and exits for every argv partition outside the three single-token help spellings, and collapse orphan classification to one owner that carries the exact clean path with its row.

## Acceptance

- [ ] [WC1] (covers QD2, QD6) single-token `--help`, `-h`, and bare `help` answer usage/0; multi-token help, separator-bearing argv, the empty token, and ordinary rejected tokens retain their checked-in pre-disclosure stdout, stderr, and exit exactly.
- [ ] [WC2] (covers QD1) active and orphan rows retain their exact actions and order without independently re-deriving the orphan predicate or pairing paths through a parallel index.

