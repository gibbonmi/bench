# Give the handoff rerun suffix one owner

Blocked by: none
Writes: internal/handoff/handoff.go, internal/handoff/state_file.go
Covers: none

## What to build

The handoff package has one constant that holds the rerun instruction ", then rerun bench handoff". The state-file repairs and the unknown-checkout refusal build their text from that constant. The refusal bytes that the command writes do not change.

## Acceptance

- [ ] The unknown-checkout refusal in the handoff package builds its repair text from the shared rerun constant, and no literal copy of the rerun instruction stays in the package.
- [ ] Each handoff refusal writes the same bytes as before the change, and the handoff package tests pass.
