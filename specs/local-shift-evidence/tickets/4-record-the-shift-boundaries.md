# 4. Record the shift boundaries

Blocked by: 2-redact-and-version-every-record-line.md
Writes: internal/shift/loop.go, internal/shift/session.go, internal/shift/result.go, internal/shift/record.go (new), internal/shift/record_test.go (new), internal/otelrecord/attributes.go, internal/otelrecord/registry.go
Covers: LE22, LE23, LE24, LE25, LE26, LE27, LE28, LE29, LE30, LE31, LE32, LE33, LE34, LE35, LE36, LE37

## What to build

Chunk: LE-B1.

Before this ticket starts, the reviewer decides the structure grant `internal/shift/ 15` in `.bench/structure.budgets`. The ticket must not edit that file.

Start one `shift` span right after the intent entry persists. End it in `finish`, so every exit path ends it, the checkpoint exit through `os.Exit` included. Add the `shift` seam to the registry.

The start line carries `bench.intent.key`, `bench.agent` (the base name of `BENCH_AGENT`), and `bench.shift.cap`. It carries `bench.line.tier` only when `BENCH_MODEL` names a tier in `lines.Tiers`.

The end line carries `bench.shift.outcome`, `bench.work.state`, and `bench.outcome`. It carries `bench.subject.id` as the branch head commit when the shift committed. It carries `bench.recovery.kind`, `bench.recovery.key`, and `bench.cleanup`. The recovery key is the base name of the retained path, so no absolute path enters the record.

The record package owns the work-state and cleanup vocabularies. The shift maps its outcomes: `complete`, `no-op`, and `incomplete` are `completed`; `failed` and `usage` are `failed`; `interrupted` is `interrupted`. Declare every new key in the declared set, because the encoder drops an undeclared key.

## Acceptance

- [ ] A green one-iteration shift writes a start line and an end line for one `shift` span.
- [ ] A shift that exhausts its branch-name retries ends with outcome `usage` and work state `failed`.
- [ ] With `BENCH_AGENT` at `<tmp>/DIRMARK dir/agent`, the span carries `bench.agent` `agent`, and no line holds `DIRMARK`.
- [ ] With `BENCH_MODEL=mid` and `BENCH_MAX_ITERS=2`, the span carries tier `mid` and cap `2`, and with `BENCH_MODEL=custom model` it carries no tier key.
- [ ] The outcome and work-state pairs are `complete` and `completed`, `failed` and `failed`, `incomplete` and `completed`, and `no-op` and `completed`.
- [ ] A shift with one commit carries the branch head commit as its subject id, and a no-op shift carries none.
- [ ] A red shift that retained its worktree carries kind `worktree`, the base-name key, and cleanup `retained`, and no line holds the Bench home path.
- [ ] A green shift carries cleanup `released`.
- [ ] A shift that receives SIGINT during its adapter, run in a helper process, writes an end line with work state `interrupted`.
