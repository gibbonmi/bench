# 9. Record the lane span with its failing check

Blocked by: keep-the-started-phase-line-after-a-kill.md, check-that-each-registered-seam-starts-a-span.md
Line: opus / medium
Rows: OT17
Writes: internal/gate/lane.go, internal/gate/lane_test.go, internal/gate/lane_record_test.go, internal/otelrecord/registry.go

## What to build

The lane writes a span at its own seam. A red lane's span carries the first
failing check and that check's first diagnostic line. The FT232 tripwire reads
the lane's red from this one attribute pair. A record of the phase and the exit
alone drops the check.

The span carries only the attributes that ticket 3 declares, so no diagnostic
payload beyond the first line enters the record.

This ticket adds the lane seam to the registry from ticket 8. Tickets 9, 10,
11, and 12 each write `internal/otelrecord/registry.go`, so the coordinator
runs them in that order.

The lane tests in `internal/gate` set `BENCH_HOME` to a temporary directory and
read the record lines back.

## Acceptance

- [ ] OT17: a red lane writes a span that carries the first failing check and its first diagnostic line.
- [ ] a green lane writes a span that names no failing check.
- [ ] the registry names the lane seam, and the conformance check from ticket 8 passes.
