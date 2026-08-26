# Show the census row on the ambient dashboard

Blocked by: 03-name-the-real-verb-head.md
Writes: internal/census/census.go, internal/census/census_test.go, internal/status/status.go, internal/status/status_signals_test.go, internal/status/status_render_test.go, CONTEXT.md

## What to build

`bench status` shows the count ambiently, so the reviewer reads it without a
question.

This contract crosses into ticket 06:

- `census.Counts(home, root string) (map[string]int, error)` returns the record count per assignment id.

An absent file and an empty file both read as zero. A refused file type reads
as zero. The reader never fails the board.

`internal/status` gains `appendCensus`, beside `appendGuards` in `SignalsWith`.
It joins the counts to the ledger's active assignments and sums them. When the
sum is above zero, it appends one row with severity 3, the signal `census`,
and the detail `<n> raw call(s) across <k> worktree(s)`. The board's own
`Plural` renders both counts, so one source states the words. The row names
no action, because no command is the remedy.

`bench status --all` expands the row to one row per active assignment, detail
`<label> <n> raw call(s)`. The expander is a sibling of the `intent` expander
in the same `--all` render branch. The five-row budget therefore holds with
many worktrees.

The row appears only for active assignments. An assignment with no records
adds nothing. A count with no active ledger entry adds nothing, so a dead
assignment never reaches the board.

The label passes through `sanitize.Controls` before the board renders it, so
a control byte cannot split a fixed-width row. Status resolves the Bench home
at its own boundary and passes it down.

Add `census` to the `**signal**` enumeration in `CONTEXT.md`, because the
signal-vocabulary conformance check reads that enumeration and reds a signal
name it does not list. Change nothing else in the glossary.

## Acceptance

- [ ] Two active assignments with two and one records return one `census` row detailed `3 raw calls across 2 worktrees`. (EC17)
- [ ] An active assignment with no file and a released assignment with two records each add no row. (EC18)
- [ ] The `census` row carries severity 3 and renders no action. (EC19)
- [ ] A label with an ESC byte renders in the expanded `census` row with the byte replaced. (EC25)
- [ ] `bench status --all` with two active assignments renders two `census` rows, each `<label> <n> raw call(s)`. (EC34)
- [ ] One assignment with one record renders `1 raw call across 1 worktree`.
- [ ] The signal-vocabulary check passes with `census` emitted.
