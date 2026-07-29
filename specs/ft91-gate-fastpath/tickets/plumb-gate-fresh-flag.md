# Plumb bench gate --fresh through the wrapper

Blocked by: Short-circuit gate execution on a reusable green

## What to build

Story 10 of `specs/ft91-gate-fastpath/spec.md`: `bench gate --fresh` forces a
real run past a reusable green. `gate_command` and `run_gate` in
`bin/bench.sh` accept and forward the flag; `RunCommand`
(`internal/gate/gate.go`) parses it without colliding with its optional root
positional; any other unknown flag still exits 2 with usage on stderr and an
untouched oracle. `bench commit` gains no such flag.

## Acceptance

- [ ] `bench gate --fresh` on a tree holding a reusable green pays a real run
      (marker count 2 at the gate execute seam).
- [ ] The extended wrapper contract in `runtime_gate_test.go` proves the flag
      arrives at the binary instead of dying at `gate_command` arity checking.
- [ ] Any other unknown flag still exits 2 with usage on stderr and no oracle
      side effects.
