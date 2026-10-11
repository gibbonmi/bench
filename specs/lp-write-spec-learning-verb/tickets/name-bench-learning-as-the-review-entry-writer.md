# Name bench learning as the writer of the write-spec review entry

Blocked by: none
Writes: .agents/commands/bench-write-spec.md
Covers: none

## What to build

The `/bench-write-spec` command tells the session to append one `capture/learnings.md` entry when a review round takes more than one iteration to accept. The command names no writer for that entry. The operating guide tells each session to run `bench learning` and never to type the entry heading by hand. Add one sentence to the same physical line of the command that names `bench learning` as the writer. The pinned anchor needles on that line stay intact, and the file keeps its line count of 87, which is its exact prose budget.

## Acceptance

- [ ] `rg -n "bench learning" .agents/commands/bench-write-spec.md` returns one hit, on line 86.
- [ ] `.agents/commands/bench-write-spec.md` holds 87 lines, so the `guidance-prose-budgets` check stays green.
- [ ] The `docs-currency-workflow` check stays green with the multi-iteration learning anchor needles intact.
