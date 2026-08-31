# unexecuted tag citation (canary fixture)

This fixture cites a stress-tagged test file, which the gate never executes, and
proves the execution diagnostic reaches the gate.

## User stories
1. As a reviewer, I want a citation into a never-executed file to fail the gate.

### Acceptance coverage map
| story | behavior | seam | why it catches the failure |
|---|---|---|---|
| 1 | a behavior | `internal/example/stress_test.go` (`TestStressOnly`) | why it catches |
