# Wire the blast step into the review phase

Blocked by: add-blast-mode.md
Writes: .agents/commands/bench-review-implementation.md, .bench/BENCH-reference.md

## What to build

The review command doc gains the blast step. Before axis dispatch, the
coordinator runs `bench consumers --changed --base <b> --source-tip <t> --full`
over the frozen pair. It attaches the table and its citation to the review
record, and it walks the `touched=false` rows first.
`.bench/BENCH-reference.md` gains the `consumers` command note. This ticket
covers no coverage row; the spec marks story 27 as Not covered, and the
landing gate's prose checks grade the files.

## Acceptance

- [ ] The review command doc names the blast step before axis dispatch,
      with the exact invocation and the walk order.
- [ ] `.bench/BENCH-reference.md` carries the `consumers` command note.
