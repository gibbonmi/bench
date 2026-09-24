# Repair the LE-B1 review findings

Blocked by: 4-record-the-shift-boundaries.md
Writes: internal/shift/record.go, internal/shift/record_test.go, internal/shift/result.go, specs/local-shift-evidence/spec.md
Covers: LE107

## What to build

Chunk: LE-B1.

The spec names the recovery kinds `none` and `worktree`, so a shift with no recovery pointer writes kind `none`. Add row LE107 and its test.

Split the recovery pointer in one helper beside its constructor. Build the test line match from the record constants. Add a test for the teardown-failure exit and a test for the acquire-failure exit.

## Acceptance

- [ ] A green shift that left no recovery pointer carries kind `none` and no key.
- [ ] A shift whose teardown fails before the release carries cleanup `none`.
- [ ] A shift whose worktree acquire fails still ends its span with outcome `usage`.
