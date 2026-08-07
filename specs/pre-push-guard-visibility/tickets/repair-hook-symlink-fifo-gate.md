# Repair hook symlink FIFO gate

Blocked by: none
Ownership fence: `internal/adopt/link_hook.go`, `internal/adopt/link_hook_test.go`
Integration surfaces: hook inspection mode gate→`internal/adopt/link_hook.go` + HFG1; hook inspection contract→`internal/adopt/link_hook_test.go` + HFG1
Contracts: none crosses
Closure: HFG1/symlink-read-through, HFG2/dead-prepush-status

## What to build

Close the read-through in `InspectPrePush`'s special-file gate: a symlink at the
hook path currently bypasses the non-regular refusal and `os.ReadFile` follows
it, so a symlink to a writerless FIFO hangs every ambient surface (`bench
status`, guards, the SessionStart banner). Resolve the symlink's target mode
before reading — a symlink whose target is not a regular file classifies
`PrePushForeign` without any read. Add the deadlined symlink-to-FIFO fixture
alongside the existing direct-FIFO test. In the same fence, delete the dead
exported `PrePushStatus` type and its doc comment — its classifier was removed
and nothing produces or consumes it.

## Acceptance

- [ ] [HFG1] (covers local) `InspectPrePush` on a symlink pointing at a writerless FIFO returns `PrePushForeign` within the test deadline instead of blocking on the read.
- [ ] [HFG2] (covers local) The `PrePushStatus` type no longer exists in `internal/adopt`.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| HFG1/symlink-read-through | drop the symlink-target mode resolution so the read follows the link | hook inspection contract | plant `ln -s <fifo> .git/hooks/pre-push`, run `InspectPrePush` under the deadline, expect the deadline red |
| HFG2/dead-prepush-status | restore the `PrePushStatus` declaration | Standards review | enumerate `PrePushStatus` references, expect the dead-exported-type finding |
