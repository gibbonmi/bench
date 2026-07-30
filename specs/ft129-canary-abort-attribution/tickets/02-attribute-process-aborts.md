# Attribute process aborts

Blocked by: Attribute Go test aborts

## What to build

The subprocess and canary runner-result contracts distinguish normal numeric
completion from spawn and signal aborts. Sweep reports those aborts before bite
classification while preserving completed-run semantics and stable aggregate
diagnostics.

## Acceptance

- [ ] Spawn failures and signaled processes are reported as process aborts even
  when they produce no output.
- [ ] Numeric exits 1 and 2 remain ordinary completed failures.
- [ ] Exit zero, completed red without `EXPECT`, and completed red with `EXPECT`
  retain their existing classifications.
- [ ] Process-abort classification wins when partial output already contains
  the fixture's `EXPECT`.
- [ ] One panic, one process abort, and one ordinary did-not-bite result are all
  returned in fixture order under concurrent execution.
- [ ] Repeated classification is byte-stable and includes neither raw stack
  output nor an unbounded child error.
