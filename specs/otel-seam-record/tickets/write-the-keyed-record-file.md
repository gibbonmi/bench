# 4. Write the record file under the Bench home

Blocked by: add-the-encoder-and-the-attribute-set.md
Line: opus / medium
Rows: OT8, OT9, OT10, OT11, OT27
Writes: internal/otelrecord/writer.go (new), internal/otelrecord/writer_test.go (new), DATA_HANDLING.md

## What to build

The package appends each encoded line to `otel/<repository key>/traces.jsonl`
under the resolved Bench home. The function `poolkey.Key(root)` supplies the
key, and the record directory is a sibling of the census directory.

Each append is one synchronous `O_APPEND` write, after the census precedent in
`internal/census`. The writer keeps no buffer, and it starts no background
worker. Two writers with separate file handles therefore leave only intact
lines.

The writer refuses the write when the record directory is a symlink, so a
redirected record cannot escape the Bench home.

The writer returns an error when the record directory is unwritable, and it
swallows nothing. The caller decides what to do with that error. Ticket 6
proves that the gate keeps its exit code.

The record file is append-only. Rotation and retention stay with FT71, and the
spec prices them under Out of scope. This ticket writes no size bound.

`DATA_HANDLING.md` gains a paragraph that names this file and its path under
the Bench home. The document already claims that Bench uploads nothing, and the
new paragraph keeps that claim true and complete.

The writer takes the resolved home as a parameter, the census form. It reads no environment variable itself; the provider resolves the home and passes it down.

## Acceptance

- [ ] OT8: with `BENCH_HOME` set to a temporary directory, the writer creates `otel/<repository key>/traces.jsonl` under that directory.
- [ ] OT9: two writers with independent file openers append many lines each, and every written line parses.
- [ ] OT10: the writer refuses the write when the record directory is a symlink.
- [ ] OT11: the writer returns an error when the record directory is unwritable.
- [ ] OT27: `DATA_HANDLING.md` names the record file and its path under the Bench home.
