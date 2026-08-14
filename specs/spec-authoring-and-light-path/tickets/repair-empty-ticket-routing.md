# Route an empty ticket directory back to write-spec

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md, internal/anchors/registry_data.go, tests/canary/workflow-guidance-anchors/implement-spec-entry-validation/

## What to build

Close the Spec-axis WF4 finding recorded in `reviews/spec-authoring-and-light-path.md`: entry validation treats an absent ticket directory and a present directory containing no ticket files as the same ticketless state, and routes both to `/bench-write-spec` before preflight. Re-derive the existing anchor and fixture rather than adding a second policy source.

## Acceptance

- [ ] the command names both absent and no-ticket-file states and routes either to `/bench-write-spec`
- [ ] the existing entry-validation mutation fixture bites the completed predicate
