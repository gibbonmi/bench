# Record one raw call per Bash call

Blocked by: 01-derive-the-census-home-beside-the-pool.md
Writes: internal/census/census.go, internal/census/census_test.go, internal/racetests/racetests.go

## What to build

`internal/census` appends one line per raw call. A raw call is one Bash command
text that names a path under the repository pool, where no simple command
invokes `bench`.

This contract crosses into tickets 03 to 06:

- `census.Record(command, root, home string, now time.Time) error` appends one line.

The caller ignores the error, so a failed write never changes a verdict.

The match reads the raw command text. A path matches when the text holds the
pool prefix and a separator directly after the pool. A sibling directory
`<pool>x/` does not match. A relative path never matches, because the match
reads absolute text only. A quoted path with a space or a glob character folds
into one word and matches.

The `bench` test is `benchguard.InvokesBench`. A text where any simple
command invokes `bench`, at the top level or one wrapper level deep, records
nothing. A `bench worktree exec` text therefore records nothing, and so does
a `bash -c` string around it, and so does a text that spells the wrapper by
an absolute path.

The record key is the assignment id. Take the path segment directly after the
pool and read the part after the owner id and the hyphen. A segment without the
owner-and-id shape records nothing. The writer never reads the ledger, so an
unknown or dead id still records under its own name.

Each record line holds the time and the verb head. It never holds the command
text and never the label; the ledger owns the label at render time. The
writer opens the file for append and writes one line per call, so two
concurrent writers leave two intact lines. Register the concurrent test in
`internal/racetests.Tests`, because the ordinary suite cannot observe the
race. A line without a trailing newline still counts as one record.

An unwritable home, a symlinked census directory, and a refused file type each
return an error and write nothing.

For this ticket the verb head is the head of the first simple command whose
words name a pool path. Ticket 03 refines that head.

## Acceptance

- [ ] A text `sed -i s/a/b/ <pool>/<owner>-<id>/x` appends one record under `<id>`. (EC01)
- [ ] A chain that names a pool path in two commands appends one record. (EC02)
- [ ] A text `bench worktree exec a -- sed -i x <pool>/<owner>-<id>/b` appends no record. (EC03)
- [ ] A text `/tmp/r/bin/bench.sh status <pool>/<owner>-<id>` appends no record, through a stub resolver. (EC04)
- [ ] A record line holds a time and a verb head, and holds no command text. (EC06)
- [ ] A text `bash -c 'bench worktree exec a -- sed -i x <pool>/<owner>-<id>/b'` appends no record. (EC11)
- [ ] A text that names an unknown assignment id appends one record under that id. (EC14)
- [ ] Two concurrent writers leave two intact lines under `-race`. (EC15)
- [ ] A text that names `<pool>/<owner>-<id>/x` records under `<id>`, not `<owner>-<id>`. (EC16)
- [ ] An unwritable home returns an error and leaves no file.
