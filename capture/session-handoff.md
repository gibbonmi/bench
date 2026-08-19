# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `5fc748b` at the time this was written, clean but for this drain batch
Spec: none staged; `specs/` is empty
Gate: green

## State

Two specs shipped and retired this session. FT226 stopped the kit's own tests
writing into the operator's `BENCH_HOME` — the reauthorize fixture binds its own
home, and `internal/worktree` runs under a process-private one that reds on
residue and names the leaking test. FT234 added `bench worktree reclaim`, the
first reader over the pool parent, so a key whose source repository is deleted
can be recovered at all. The operator's pool went from 1,719 keys and 91 MB to
one key and 120 KB, and the FT234 acceptance demo ran against that real pool.

This drain reconciled both shipped rows off the board, drained two ideas, four
learnings, and two retros, and opened FT235 from the pool-directory-naming idea.
Every capture source is empty.

The board's most-evidenced open thread is the reviewed-landing spec-byte
question: FT225 (decide whether a review may amend the spec in its own source)
and FT233 (landing refusals name their remedy) each took new occurrences today,
and FT233's wording depends on FT225's decision. That is why the sequence puts
FT225 third despite two HIGH rows above it.

`dist/bench` was a day stale for most of this session; it is rebuilt. Rebuild it
before hand-running any newly landed verb.

## Next command

`/bench-write-spec` — FT227: adoption smoke, so a newly adopted repository's
scaffolded gate can go green.

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
