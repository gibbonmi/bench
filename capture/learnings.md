# Learnings — usage journal

- 2026-08-19 — /bench-write-spec (FT229) took 2 iterations to accept. Stage
  that missed: spec authoring. What review caught: (1) two coverage rows
  promised that a diagnostic "names the command that rebuilds the binary"
  without pinning the command, so the cheapest wrong build prints plain
  `go build` — which `projects/benchkit.md` explicitly forbids — and both rows
  pass; (2) the spec invented a standing conformance check holding the
  tickets-only count at zero, which would red on every future light-path
  landing, because the same spec's close step requires that folder to survive
  its own landing gate run. The check was also an addition beyond the decision
  source, which asked for a deletion plus a status row. Why missed: the author
  applied the rubric's exact-predicate test to behavior cells but not to a
  quoted remedy string, and added a durable check without walking the in-flight
  tree states the rest of the spec creates. Proposed rule change: a row
  promising that output "names the command" carries the exact command token,
  never the family; and a spec adding a standing check enumerates the in-flight
  states that check must tolerate, drawn from the same spec's own lifecycle.
