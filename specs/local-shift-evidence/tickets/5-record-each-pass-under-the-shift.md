# 5. Record each pass under the shift

Blocked by: 4-record-the-shift-boundaries.md
Writes: internal/shift/loop.go, internal/shift/session.go, internal/shift/record.go (new), internal/shift/record_test.go (new), internal/otelrecord/attributes.go, internal/otelrecord/registry.go
Covers: LE37, LE38, LE39, LE40, LE41, LE42, LE43, LE44

## What to build

Chunk: LE-B2.

Open one `shift.iteration` span for each main iteration and one `shift.refactor` span for each refactor pass, under the `shift` span. Add both seams to the registry and declare the new keys.

Pass the pass span's context to the adapter run and to `gate.RunAndRecordContext`, so the gate span is the child of the pass span. The pass span carries `bench.adapter.result` (`exited` or `spawn-failed`) and `bench.adapter.exit` when the adapter exited. A pass that committed carries that commit as `bench.subject.id`. No span carries a pass index or a committed count, because a reader derives both from the pass spans.

The session keeps the open pass span. `finish` ends it before it ends the `shift` span, so the checkpoint exit leaves both end lines.

The SIGINT rows need a real process exit, because the checkpoint exits through `os.Exit`. Re-exec the shift test binary into one role test of the shift package, and kill that process group on cleanup. This is the Go helper-process idiom that each package authors for its own roles, as `internal/gocache/lock_test.go` does. Copy nothing from the worktree package's `descendant`.

The spawn-failure row cannot use an adapter that is not executable at start, because `requireAdapter` refuses it before the span starts. Its adapter commits a change in iteration 1 and then runs `chmod -x` on its own file, so iteration 2 fails to spawn.

## Acceptance

- [ ] A two-iteration shift writes exactly two `shift.iteration` spans under the `shift` span.
- [ ] A pass whose adapter exits 3 carries result `exited` and exit `3`.
- [ ] A pass whose adapter removed its own execute bit in iteration 1 carries result `spawn-failed` and no exit key in iteration 2.
- [ ] A refactor pass writes one `shift.refactor` span under the `shift` span.
- [ ] Each pass span has exactly one child `gate` span with a subject id and an outcome.
- [ ] A pass that committed carries the commit that the branch gained in that pass.
- [ ] A shift that receives SIGINT during its adapter writes an end line for its `shift` span with work state `interrupted`.
- [ ] That shift writes the end line of the open pass span before the end line of the `shift` span.
