# Extract the gate run into one transaction module

Blocked by: characterize-gate-contention-and-persistence-failures.md
Writes: internal/gate (production files; one new run-module file; test edits limited to mechanical renames)

## What to build

The expand half of story 2. The body of `executeSubjectWithRunBinary` — subject
acceptance, pre-lock reuse, git-dir + lock + owner file, under-lock drift and
reuse, run-binary selection, pending record, oracle dispatch, and the one
terminal outcome pairing retain/invalidate with the terminal record and
restoring pending on a failed replace — moves into one production file that is
the only file in `internal/gate` referencing `bench-gate.lock`,
`bench-gate-owner`, `interruptedRecord(`, `retainGreen(`, `invalidateEvidence(`,
and the terminal/timeout `verdictRecord` construction. `Execute`,
`ExecuteReusingFreshGreen`, `ExecuteTree`, `RunCommand`, `RunAndRecord`, and
`RunAndRecordContext` become thin adapters (resolve root/kit/mode, log, call
the run module). No exported signature changes. `gateEngine` still exists after
this ticket (the contract ticket removes it); the move must not change any
observable behavior — the exit test is the whole pre-existing suite plus the
two characterization tickets, unmodified. Naming and file placement are
discretion within `internal/gate/`. A defect found mid-move is reported in the
return note, never fixed here; the coordinator parks it (the ticket's fence
does not include `capture/`).

## Acceptance

- [ ] `rg` over `internal/gate/*.go` (excluding tests) shows the six references above in exactly one production file, and none in the exported entry points (covers GT6)
- [ ] Every characterized outcome (GC1–GC11, GT3–GT5) and every consumer-visible gate behavior is unchanged, with test assertions unmodified — only rename hunks in `*_test.go` (covers GT1)
- [ ] Every exported function in `internal/gate` keeps its pre-ticket signature and no consumer package needs an edit
