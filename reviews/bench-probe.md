# Review pickup: bench-probe

Frozen base: `963aa880f9365b6a34211e49fd067ef421bff399`. Reviewed tip: `5e7f952501b2dc5733e468f7408efb766b35bd55`. Three axes ran on `opus`: Standards and Spec at low effort, and Coverage at medium effort in its own worktree.

## Standards

Findings: 7 raw, 5 repair targets, 3 repaired at ticket 07. Worst issue: the stub `go` harness in `internal/testreport/outcome_test.go` beside `writeCheckGo`.

- `no-op` — `installCannedGo` in `internal/testreport/outcome_test.go` beside `writeCheckGo` in `internal/testreport/check_test.go`. The two stubs encode different facts. One records the environment and replays one event. The other replays a canned stream with an exit code. The spec's testing decision names `writeCheckGo` as prior art to mirror.
- `ask-user` — `resolveSubject` in `internal/probe/subject.go` re-derives the operand rule of `anchorQueryPath` in `cmd/bench/anchors_command.go`. A shared seam is a design call.
- `no-op` — `runProbe` in `internal/probe/probe_test.go` forwards one call. It names intent at about thirty sites.
- `no-op` — `proseCheck` in `internal/probe/probe.go` restates the recorded decision in the spec's Further notes.

## Spec

Findings: 6 raw, 3 repair targets, 1 repaired at ticket 07. Worst issue: the build amended row PB4 in its own commits.

- `ask-user` — PB4 now reads `--package ./` and no `--run`. The reviewer accepted the original row. The reason is recorded under Build decisions.
- `ask-user` — the fence amendments. `cmd/bench/command_registry.go` took a pure move under a closure-headroom entry, and `internal/testreport/check_test.go` joined the fence. Both are recorded under Build decisions.
- `no-op` — `toon.TableTyped` for the verdict row, the `unreadable` reason, and the `proseCheck` constant each restate a recorded decision.

## Coverage

Findings: 3 raw, 3 repair targets, 1 repaired at ticket 07. Worst issue: a render refusal over a failed restore drops the `restore-failed` verdict.

- `ask-user` — a subject inside the checkout's administration directory mutates and restores. The spec closes the hostile-subject list at five reasons, so a sixth reason is a spec change.
- `ask-user` — a build failure in one package beside a failing test in another package reports `bit` at 0 under `--package ./...`. Ticket 02 fixes the precedence as a failing test row first, so the behavior obeys the spec. The reviewer decides whether a build failure anywhere must dominate.
- `ask-user` — a control byte in the subject's base name reaches the copy path. The `preserved[1]` row then refuses, and the caller gets the render error alone at exit 2. The repair-scoped re-review found it. The spec now lists the edge as a Won't handle for veto.
