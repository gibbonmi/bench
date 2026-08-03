# Route terminal-run restart through the live-marker fallback

Blocked by: none
Ownership fence: `internal/specbuild/assign.go`, `internal/specbuild/start_test.go`
Assumptions: `startRun`'s empty-`previousGreen` fallback at `internal/specbuild/assign.go:266-268` is the one existing marker-read site; `retainTerminalAttempt` (`internal/specbuild/precondition.go:278`) already marshals the whole record including `Base` into `History`; claims re-derived from the tree at pickup

## What to build

A reviewer restarting an abandoned spec build after a sibling build promoted on
the same branch gets a successful `Start` instead of "project-green marker
conflicts with another tip". The terminal-restart branch (`assign.go:242`)
passes an empty `previousGreen` so `startRun`'s live-marker fallback fires;
`run.Base` is passed nowhere — it survives only in the retained attempt
history. Restart with the marker absent inherits the fallback's absent-marker
bootstrap semantics. Fresh (non-restart) starts are untouched.

## Acceptance

- [ ] [RM1] After a sibling promotion advances branch and green marker past the abandoned run's recorded base, `Start` of that run succeeds end-to-end against the real authorization gate.
- [ ] [RM2] On restart, the recording gate observes `expected` equal to the live marker (not `run.Base`), and the retained history entry still carries the old run's `base` field.
- [ ] [RM3] Restart with no green marker present bootstraps fresh, exactly as a fresh start does.
- [ ] [RM4] The fresh-path marker tests in `start_test.go` run unchanged and green.
- [ ] [RM5] An abandoned empty run whose tip advanced still restarts into a new attempt (composition control, green at introduction — pins that restart consumes the recompose refusal).

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RM1 | revert the operand: pass `run.Base` as `previousGreen` at the terminal-restart call | the sibling-promotion restart test | apply, run `go test ./internal/specbuild -run TestRestart`, expect the marker-conflict failure |
| RM2 | drop the `retainTerminalAttempt` history append from the restart call | the retained-history assertion | apply, run the recording-gate restart test, expect the missing-history failure |
| RM3 | pass a non-empty synthetic `previousGreen` on the absent-marker path | the marker-absent restart test | apply, run it, expect the "does not match expected prior tip" failure |
