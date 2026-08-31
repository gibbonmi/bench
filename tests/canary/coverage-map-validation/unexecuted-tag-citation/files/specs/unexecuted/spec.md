# unexecuted citation (canary fixture)

This fixture cites two files the gate never executes. One carries a build tag no test
phase compiles with. The other sits in a package Go excludes from every recursive
pattern. Together they prove both execution diagnostics reach the gate.

## User stories
1. As a reviewer, I want a citation into a never-executed file to fail the gate.

### Acceptance coverage map
| story | behavior | seam | why it catches the failure |
|---|---|---|---|
| 1 | a behavior | `internal/example/stress_test.go` (`TestStressOnly`) | why it catches |
| 1 | another behavior | `testdata/fixture/fixture_test.go` (`TestIgnoredPackage`) | why it catches |
