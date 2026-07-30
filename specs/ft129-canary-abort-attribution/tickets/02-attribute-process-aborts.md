# Attribute process aborts

Blocked by: Attribute Go test aborts

## What to build

The subprocess and canary runner-result contracts distinguish normal numeric
completion from spawn and signal aborts. Sweep reports those aborts before bite
classification while preserving completed-run semantics and stable aggregate
diagnostics.

## Acceptance

- [x] Spawn failures and signaled processes are reported as process aborts even
  when they produce no output.
- [x] Numeric exits 1 and 2 remain ordinary completed failures.
- [x] Exit zero, completed red without `EXPECT`, and completed red with `EXPECT`
  retain their existing classifications.
- [x] Process-abort classification wins when partial output already contains
  the fixture's `EXPECT`.
- [x] One panic, one process abort, and one ordinary did-not-bite result are all
  returned in fixture order under concurrent execution.
- [x] Repeated classification is byte-stable and includes neither raw stack
  output nor an unbounded child error.
