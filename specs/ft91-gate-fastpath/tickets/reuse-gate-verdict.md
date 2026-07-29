# Short-circuit gate execution on a reusable green

Blocked by: Extend freshness to 60 minutes and declare the toolchain tools

## What to build

Stories 8 and 9 of `specs/ft91-gate-fastpath/spec.md`: the reuse
short-circuit lands in `executeWithEngineAfterAcquire` — under the held lock,
after the under-lock subject re-check, before the pending replace. On
`ReusableGreen` it prints `gate: green (fresh verdict reused for this tree)`,
writes nothing (verdict record, `RecordedAt`, cache file byte-identical), and
returns the full green tuple (`GateExit` 0, `ActionExit` 0, an `Inspection`
projecting the reusable green). Everything up to and including lock
acquisition stays byte-identical. In the same change, the R17 fault-injection
family (`runner_test.go` / `runtime_gate_partial_proof_test.go`) is re-seeded
to non-reusable subjects (expired `RecordedAt` via the fake engine's clock)
for every faulted op the short-circuit now precedes, registry tuples updated;
pre-acquire contracts (`owner_record_test.go` lock/demotion) stay unedited.

The blocker is a file fence, not logic: this ticket and its blocker both edit
`internal/gate` and `gate_test.go`.

## Acceptance

- [ ] A second `Execute` on an unchanged tree within freshness runs nothing:
      marker count stays 1, stdout carries the reuse line, cache record
      byte-identical, returned tuple green with an `Inspection` projecting the
      reusable green (throwaway repo whose gate script appends a run marker).
- [ ] Red, expired, pending, subject-open, and tree-changed states each pay a
      real run.
- [ ] Every R17 faulted op is still reached on re-seeded non-reusable subjects
      (`opCounts`), and the pre-acquire demotion contracts pass unedited.
