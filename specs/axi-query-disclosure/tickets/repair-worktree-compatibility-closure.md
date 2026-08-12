# Repair worktree compatibility closure

Blocked by: none
Writes: `internal/worktree/list.go`, `internal/worktree/list_actions_test.go`, `internal/worktree/testdata/pre-disclosure-argv-pairs.json`, `internal/usage/parse.go`, `internal/usage/parse_test.go`

## What to build

Preserve the pre-migration `worktree list` bytes and exits for every argv partition outside the three sole-argument help spellings. Add a grammar-scoped `usage.Parse` option that restricts help recognition to a sole argument for this command only; do not add a local parser bypass. Collapse orphan classification to one owner that carries the exact clean path with its row, and preserve the producers' ordering contract.

## Acceptance

- [ ] [WC1] (covers QD2, QD6) a shared-parser grammar option, exercised by `worktree list` only, makes sole-argument `--help`, `-h`, and `help` answer usage/0. The checked-in old/new matrix enumerates `--help extra`, `-h extra`, `help extra`, bare `--`, `-- x`, the empty token, one ordinary extra, and multiple ordinary extras; every row preserves its pre-disclosure stdout, stderr, and exit exactly.
- [ ] [WC2] (covers QD1) active and orphan rows retain their exact actions without independently re-deriving the orphan predicate or pairing paths through a parallel index. Assignments remain in the intent ledger's serialized ID order, followed by Git registrations in their producer order; the public test consumes that producer order and neither assumes creation order nor independently sorts it.
