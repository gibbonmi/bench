# Leave temporary-file cleanup to the leaf test

Blocked by: none
Writes: internal/reviewrecord/replacement_test.go
Covers: none

## What to build

The review-record replacement test reads the record bytes without a sibling-file check in the record directory. The durable-file leaf test owns the temporary-file cleanup postcondition.

## Acceptance

- [ ] `internal/reviewrecord/replacement_test.go` does not call `durablefiletest.Siblings`.
- [ ] `TestReplaceTemporaryCleanup` fails when `durablefile.Replacer.Replace` omits temporary-file removal.
