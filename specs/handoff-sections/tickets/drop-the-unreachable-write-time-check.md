# Drop the unreachable write-time check

Blocked by: refuse-a-state-that-breaks-the-grammar.md
Writes: internal/handoffdoc, internal/handoff/handoff.go, internal/handoff/state_scan.go, internal/handoff/state_scan_test.go
Covers: HS27, HS29

## What to build

Repair for the re-review finding Sp-R1. Verify the premise first.
`handoffdoc.Update` parses the document before it calls the mutator, and
the verb never rewrites State. So `ValidateState` grades a State that
`Parse` already accepted, and no production input reaches its refusals.

Then delete `ValidateState`, `StateError`, `errOrNil`, and `stateGrammar`.
Keep `scanFences` for the scan, and drop the line index that no caller
reads. Keep the two verb tests, and make each assert the parser refusal's
file and line in the printed error.

## Acceptance

- [ ] `rg ValidateState` over the tree finds nothing.
- [ ] The two verb tests assert the refusal names the file and line.
- [ ] The fenced-commit scan test stays green.
- [ ] Self-probe: drop the fence skip in the scan, and report the fenced-commit test red.
