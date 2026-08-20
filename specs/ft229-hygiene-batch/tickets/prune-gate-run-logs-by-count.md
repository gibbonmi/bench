# Prune gate run logs by count

Blocked by: none
Writes: internal/gate/run_log.go

## What to build

`.logs` gains one record per gate run and nothing removes them. The gate becomes
the pruning owner: after a run, it retains the newest 20 records and removes the
rest, so no separate chore has to be remembered. Only files matching the gate's
own record name shape are considered, and the record the current run is writing
is never a candidate. A pruning failure warns on stderr and never fails the
gate.

## Acceptance

- [ ] a gate run over a log directory holding more than 20 records leaves exactly the newest 20 (H27).
- [ ] the record the current run is writing survives its own pruning (H28).
- [ ] a file in `.logs` outside the gate's record name shape survives pruning (H29).
- [ ] a pruning failure writes one stderr warning and leaves the gate's verdict and exit code unchanged (H34).
