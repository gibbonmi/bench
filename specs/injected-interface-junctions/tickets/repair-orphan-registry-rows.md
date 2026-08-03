# Repair the injected-port check's orphan-row blind spot

Blocked by: none
Ownership fence: `internal/conformance/injected_ports_test.go`
Contracts: every registry row's port name crosses registry→derived inventory; a row naming a port the derivation no longer reports fires a distinct orphan-row diagnostic, asserted by OR1 against the real check
Assumptions: sol review round 1 findings SP2/C2 are the authority; the four existing messages stay unchanged; claims re-derived from the tree at pickup

## What to build

The check visits every registry row: a row whose port is absent from the
derived inventory fails with a fifth distinct message (orphan row — a stale
advertisement whose port or injection arm left the tree), so deleting an
injection arm together with its junction test cannot leave a green row while
other ports keep the inventory non-zero.

## Acceptance

- [ ] [OR1] a unit fixture with a registry row naming an underived port (other ports still derived) fails with the orphan-row message; the same fixture also proves the missing-test message still fires for a derived port.
- [ ] [OR2] the real registry passes the extended check (scoped root-conformance run green, cache-busted, execution proven not skipped).

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| OR1 | drop the orphan-row branch | the unit fixture | run the unit test, expect the missing orphan diagnostic to fail it |
| OR2 | add a fake registry row naming an underived port | the scoped conformance run | run scoped with `BENCH_CONFORMANCE_ROOT` set, expect the orphan-row red; revert |
