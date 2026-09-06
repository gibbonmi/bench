# Review pickup: bench-probe

Frozen base: `963aa880f9365b6a34211e49fd067ef421bff399`. Reviewed tip: `5e7f952501b2dc5733e468f7408efb766b35bd55`. Three axes ran on `opus`: Standards and Spec at low effort, and Coverage at medium effort in its own worktree.

## Standards

Findings: 7 raw, 5 repair targets. Worst issue: the stub `go` harness in `internal/testreport/outcome_test.go` beside `writeCheckGo`.

- `no-op` — `installCannedGo` in `internal/testreport/outcome_test.go` beside `writeCheckGo` in `internal/testreport/check_test.go`. The two stubs encode different facts. One records the environment and replays one event. The other replays a canned stream with an exit code. The spec's testing decision names `writeCheckGo` as prior art to mirror.
- `auto-fix` — the gate lock path `.git/bench-gate.lock` is composed twice in `internal/probe/refusal_test.go`. One local helper collapses it. The cross-package derivation stays, because `internal/gate` does not export the path.
- `ask-user` — `resolveSubject` in `internal/probe/subject.go` re-derives the operand rule of `anchorQueryPath` in `cmd/bench/anchors_command.go`. A shared seam is a design call.
- `auto-fix` — comment provenance and narration in `internal/testreport/outcome_test.go` (the base-commit golden note and the prior-art note) and in `internal/probe/probe_test.go` (the FT290 note). The commit or the spec owns the red record.
- `auto-fix` — `awaitFile` in `internal/probe/refusal_test.go` discards a live argument.
- `no-op` — `runProbe` in `internal/probe/probe_test.go` forwards one call. It names intent at about thirty sites.
- `no-op` — `proseCheck` in `internal/probe/probe.go` restates the recorded decision in the spec's Further notes.

## Spec

Findings: 6 raw, 3 repair targets. Worst issue: the build amended row PB4 in its own commits.

- `auto-fix` — PB38 promises `restore-failed` at 2 over every kind, and `TestVerdictExitCodes` covers only the six-kind mapping. The override lives in `render`, and only the `failed` cause exercises it.
- `ask-user` — PB4 now reads `--package ./` and no `--run`. The reviewer accepted the original row. The reason is recorded under Build decisions.
- `ask-user` — the fence amendments. `cmd/bench/command_registry.go` took a pure move under a closure-headroom entry, and `internal/testreport/check_test.go` joined the fence. Both are recorded under Build decisions.
- `no-op` — `toon.TableTyped` for the verdict row, the `unreadable` reason, and the `proseCheck` constant each restate a recorded decision.

## Coverage

Findings: 3 raw, 3 repair targets. Worst issue: a render refusal over a failed restore drops the `restore-failed` verdict.

- `auto-fix` — a BEL-named subject with a failed restore prints the render error line and exits 1. The mutation stays on disk, and the copy is never named. The spec says the `restore-failed` verdict overrides the run's verdict and names the copy. The repair prints the preserved row after the render error line and exits 2. Row PB47 owns it.
- `ask-user` — a subject inside the checkout's administration directory mutates and restores. The spec closes the hostile-subject list at five reasons, so a sixth reason is a spec change.
- `ask-user` — a build failure in one package beside a failing test in another package reports `bit` at 0 under `--package ./...`. Ticket 02 fixes the precedence as a failing test row first, so the behavior obeys the spec. The reviewer decides whether a build failure anywhere must dominate.
