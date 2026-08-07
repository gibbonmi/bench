# Repair stale hooks with doctor

Blocked by: expose-hook-health-record.md, render-doctor-hook-health.md
Ownership fence: `internal/adopt/doctor.go`, `internal/contract/surface/doctor_test.go`
Integration surfaces: hook-health producer→expose-hook-health-record.md; doctor stale signal and opt-in repair→render-doctor-hook-health.md; repair route→`internal/adopt/doctor.go` + R1/R2/R3/R4/R5/R6/R7; doctor runtime contract→`internal/contract/surface/doctor_test.go` + R1/R2/R3/R4/R5/R6/R7
Contracts: stale, current, foreign, and absent hook-health records cross `internal/adopt/link_hook.go`→`internal/adopt/doctor.go`, asserted by R1/R3/R5/R6 against the real exported producer
Closure: R1/stale-red-remedy, R2/plain-read-only, R3/repair-current, R4/executable-mode, R5/current-noop, R5/foreign-refusal, R6/absent-no-install, R7/repeat-idempotency

## What to build

Make stale managed hooks red with an opt-in `bench doctor --fix` repair, while keeping plain doctor read-only and refusing unsafe repair targets. Keep the stale diagnostic and repair together: a diagnostic-only cut strands the doctor runtime contract red because its named `--fix` remedy leaves the re-read stale; a repair-only cut has no stale signal to authorize it.

## Acceptance

- [ ] [R1] A stale managed hook reds doctor and names `bench doctor --fix`.
- [ ] [R2] Plain `bench doctor` never changes hook bytes or mode.
- [ ] [R3] `bench doctor --fix` repairs a stale managed hook and a fresh record read reports current.
- [ ] [R4] A repaired hook is executable.
- [ ] [R5] `--fix` reports an already-current hook unchanged and refuses a foreign hook byte-for-byte.
- [ ] [R6] `--fix` does not install an absent hook.
- [ ] [R7] Running `bench doctor --fix` twice over a stale hook is idempotent in bytes, mode, and exit behavior.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| R1/stale-red-remedy | keep stale currency green | doctor runtime contract | plant the dead-protocol hook, run doctor, expect nonzero and the exact remedy |
| R2/plain-read-only | call the repair path from plain doctor | doctor runtime contract | snapshot bytes and mode, run plain doctor, expect exact equality |
| R3/repair-current | repair only the shim | doctor runtime contract | run `--fix`, re-read hook health, expect current |
| R4/executable-mode | write repaired bytes with mode `0644` | doctor runtime contract | repair the hook, stat it, expect executable mode |
| R5/current-noop | rewrite an already-current hook | doctor runtime contract | snapshot current bytes and metadata, run `--fix`, expect unchanged no-op output |
| R5/foreign-refusal | overwrite a marker-less hook | doctor runtime contract | plant foreign bytes, run `--fix`, expect refusal and byte equality |
| R6/absent-no-install | create a hook when none exists | doctor runtime contract | run `--fix` in an absent fixture, expect the path to remain absent |
| R7/repeat-idempotency | report a second repair after the first | doctor runtime contract | run `--fix` twice, compare bytes, mode, and exits, expect the second no-op |
