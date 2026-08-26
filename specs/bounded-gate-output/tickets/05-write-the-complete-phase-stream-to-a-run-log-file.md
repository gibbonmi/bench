# Write the complete phase stream to a run log file

Blocked by: 03-bound-and-classify-the-red-failure-rows.md
Writes: internal/gate/run_log.go, internal/gate/phase_stream.go (new in 02), internal/gate/run_log_prune_test.go, internal/gate/run_stream_test.go (new), internal/gate/report_test.go (new in 02)

## What to build

The run log owner writes `.logs/gate-<run>.out` beside
`.logs/gate-<run>.jsonl`. The phase buffer writes each line to that file as
the line arrives, so a killed run keeps its stream. Each written line names
its phase, as the old prefix writer did. When the file opens, the engine
prints `gate: stream <path>` once on stderr.

This contract closes the path that ticket 03 opened:

- `(*phaseStreams).path()` returns the `.out` file path once that file opens.

One recognizer owns both names. `gateLogRunFromRecordName` accepts both
suffixes and returns the same run token. The prune counts runs, not files.
It keeps the twenty newest runs and removes each older run's `.out` and
`.jsonl` together. Read one tree fact first: the prune keeps twenty files
today, so a suffix added without this change would keep ten runs.

The stderr line `gate: progress log <path>` stays as it is. When the file
cannot open, the run keeps its existing warning and adds no new diagnosis.
The projection stays bounded, and its more-row says `(stream unavailable)`.
A `.logs` directory that is a symlink receives nothing. This ticket runs in
parallel with ticket 04; the two write no common file except
`report_test.go`, where this ticket adds one test and changes none.

## Acceptance

- [x] A run killed after its first phase leaves that phase's lines in
      `.logs/gate-<run>.out`. (BG24, its end-to-end seam test added by
      ticket 08.)
- [ ] A twenty-first run leaves twenty `.jsonl` files and twenty `.out` files and removes the oldest pair. (BG25)
- [x] A `.logs` the run cannot open its stream through leaves the red table
      at twenty rows plus `+<k> more lines (stream unavailable)`. (BG26,
      retargeted by ticket 08 from a chmod fixture to a symlinked `.logs` —
      both reach the same refusal in `openGateStreamFile`.)
- [ ] A run whose stream file opens prints `gate: stream <path>` once on stderr. (BG35)
- [ ] Every phase line of a completed run appears in the `.out` file in arrival order.
- [ ] A `.logs` symlink writes no stream file and reports the stream unavailable.
