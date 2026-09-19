# 5. Grade the named-check test on a private root

Blocked by: 4-run-build-scripts-on-kit-copy.md
Writes: internal/testreport/check_test.go
Covers: TD22

## What to build

Chunk: TD-C2.

`TestNamedCheckRunsOnlyRegisteredDevScope` grades a kit copy from the shared git test scaffold. It compares the timing file of that copy before and after the named check. It reads no file under the live checkout's git directory, so a conformance run in another package cannot change its verdict.

## Acceptance

- [ ] The test reads no path under the live checkout's git directory.
- [ ] The test still fails when the named check writes the timing file of the root it grades.
