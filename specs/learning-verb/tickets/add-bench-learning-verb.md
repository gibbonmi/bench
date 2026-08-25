# Add the bench learning verb

Blocked by: none
Writes: internal/learnings/entry.go (new), internal/roadmap/learning.go (new), internal/roadmap/roadmap.go, internal/adopt/init.go, cmd/bench/main.go, .bench/BENCH.md

## What to build

An agent captures a process learning with `bench learning "<title>" --what <text> --right <text> [--rule <text>]`. The verb appends one entry to `capture/learnings.md` in the shape the journal parser reads back as open. The `[open]` marker is the drain's state flag. The verb routes the write like `bench idea`: an ignored journal writes the primary checkout's copy, and a tracked journal refuses the primary checkout. One formatter in `internal/learnings` is the single source of the entry shape for the verb, the adopt scaffold's worked example, and the parser. The rule in `.bench/BENCH.md` names the verb instead of a hand append.

## Acceptance

- [ ] `bench learning` appends an entry that `learnings.Parse` returns as one open entry with no malformed record.
- [ ] A title without both `--what` and `--right` exits 2 and creates no file.
- [ ] From a linked worktree with an ignored journal, the entry lands in the primary checkout's copy.
- [ ] On a primary checkout with a tracked journal, the verb returns the shared refusal on exit 1.
- [ ] `bench help` lists `bench learning`, and `bench learning --help` prints its usage.
