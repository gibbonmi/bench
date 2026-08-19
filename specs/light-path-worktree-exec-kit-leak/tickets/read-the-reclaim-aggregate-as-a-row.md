# Read the reclaim aggregate as a row

Blocked by: none
Writes: internal/worktree/pool_reclaim_test.go

## What to build

`TestReclaimApplyRemovesExactlyThePlannedKeys` asserts on its aggregate row by
substring: it looks for `3,2,1,` followed by the bare fingerprint. The
fingerprint is random per run, and the TOON writer quotes a cell that reads as
a number, so a fingerprint of all digits and leading zeros comes back quoted
and the substring never matches. The test reds on roughly one run in a
thousand, on a tree that changed nothing near it, and the red names a counting
failure that did not happen.

Read the aggregate the way the package already reads its other rows: locate the
row under its header, split its four fields, and compare each through
`unquoteCell`, which is this package's one statement of "a TOON cell may
arrive quoted". The assertion then tests what it means to test — the counts and
the fingerprint — rather than one of the two spellings the writer may choose.

## Acceptance

- [ ] The test passes when the aggregate row's fingerprint is quoted.
- [ ] The test passes when the aggregate row's fingerprint is unquoted.
- [ ] The test still reds when the aggregate reports the wrong key, removed, or retained count.
- [ ] The test still reds when the aggregate reports a fingerprint other than the planned one.
