# Link the adopt test fixture once for each package run

Blocked by: none
Writes: internal/adopt/golden_test.go (new), internal/adopt/link_transaction_test.go, internal/adopt/lifecycle_consumption_test.go, internal/adopt/lifecycle_special_test.go, internal/adopt/lifecycle_unlink_test.go, internal/adopt/lifecycle_preview_test.go
Covers: none

## What to build

The adopt test package runs one real Link into a golden consumer repository once for each package run. A test that needs a linked repository as its start state takes a copy of the golden repository and does not run its own Link. A copy is plain file writes with no synchronized writes, so it costs milliseconds. A real Link costs seconds.

Today about 30 call sites run a real Link, and most of them only prepare the start state. The parent-test fixture `linkedBase` and the copy fixture `linkedCopy` already show the pattern at the scope of one parent test. This ticket moves the base to package scope. The package keeps one owner for the golden repository and one owner for the copy.

A test whose subject is Link itself keeps its own real Link. Examples are a test that asserts on the Link output, the Link exit, or an arrangement that exists before the Link. A test that arranges the repository before its Link keeps that Link. The exception is an arrangement that can move after a copy with no change to the proof. A test that asserts on file identities from the original Link keeps its own Link.

This ticket changes only test files. It changes no production file, and it makes no test run in parallel. The remaining Link-as-setup sites in other adopt test files are a later ticket.

## Acceptance

- [ ] One real Link builds the golden repository once for each package run, and the fixtures that start from a linked repository copy it.
- [ ] Each migrated test proves the same behavior as before. Its subject operation runs for real on its copy, and its assertions stay as strong.
- [ ] Each test whose subject is Link keeps its own real Link.
- [ ] The golden repository has one owner, and the copy has one owner. No second copy helper exists.
- [ ] `go test -count=1 ./internal/adopt/` passes, and the author reports its elapsed time on the base commit and on the ticket tip.
