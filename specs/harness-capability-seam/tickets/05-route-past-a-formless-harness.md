# Route past a formless harness

Blocked by: 03-derive-the-status-harness-grammar-from-the-record.md
Writes: internal/status/route.go, internal/status/route_test.go, internal/status/handoff.go, internal/status/handoff_test.go, internal/handoff/handoff.go, internal/handoff/render.go, internal/handoff/render_test.go, internal/status/status_command_test.go

## What to build

A shell operator and an OpenCode agent both get a runnable next step, or an
honest empty cell.

A phase action is invocable for a harness only when that harness has a phase
form. `Route` skips a phase action for a formless harness the way it skips a
prose action today. So a board with a phase signal above a `git push` signal
routes `--harness none` to `git push`, and `--harness opencode` routes the
same way. A board with only phase signals routes `--harness none` to the
first signal with an empty command and `NoCommand` true.

The lead's `state` and `why` never depend on the harness. `Route` does not
re-rank the board per harness; only the command cell differs. A harness with
a phase form keeps today's translation, so `--harness codex` on a
phase-led board prints `$bench-` in the command cell.

`bench handoff --harness none` accepts the value and writes the routed shell
command under `## Next command`. A board led by `git push` therefore writes
`git push` there, so a cold session without a harness resumes.

This ticket reads the prefix table ticket 03 derives from the record. It
shares `internal/status/route.go` with ticket 03, so the two run in
sequence.

## Acceptance

- [ ] A board with a phase signal above a `git push` signal routes `--harness none` to `git push`. (covers HC17)
- [ ] The same board routes `--harness opencode` to `git push`. (covers HC18)
- [ ] A board with only phase signals routes `--harness none` to the first signal with an empty command and `NoCommand`. (covers HC19)
- [ ] The lead's `state` and `why` are equal across all four harnesses for one board. (covers HC20)
- [ ] `bench handoff --harness none` on a board led by `git push` writes `git push` under `## Next command`. (covers HC22)
- [ ] `bench status --route --harness codex` on a board led by a phase prints `$bench-` in the command cell. (covers HC50)
- [ ] A formless harness leaves a phase signal out of the runners-up.
