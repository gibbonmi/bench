# Deny a fork delegation that declares an ignored model

Blocked by: Resolve the line through the harness matrix

## What to build

The guard learns the delegation type from `tool_input.subagent_type` and stops
grading a fork backwards. When the value is exactly `fork` the harness runs the
delegation on the invoking session's model and ignores any declared one, so a
fork carrying a model denies and a fork carrying none is allowed with a warning
naming the inheritance. Comparison is exact-string on the discriminator only —
no provider lookup, no model discovery, no session-model cache — so an
unexpected or absent value never impersonates a fork and every non-fork path
keeps the postures it already settled.

Covers story 4.

## Acceptance

- [ ] A fork envelope declaring a bound model denies with exit 2.
- [ ] A fork envelope declaring no model allows with a warning stating that the
      delegation inherits this session's model, without naming that model.
- [ ] Every bound cell in the matrix is still allowed through the Claude guard,
      including the Codex column.
- [ ] A non-fork envelope with no model keeps its routed-complete exit 2, and an
      unrouted or incomplete binding keeps its fail-open rim.
- [ ] A missing, malformed, or unexpected `subagent_type` never impersonates a
      fork.
