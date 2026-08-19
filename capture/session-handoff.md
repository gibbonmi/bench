# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `eefc96f4` at the time this was written; the tree holds
`specs/ft227-adoption-smoke/` and this file uncommitted
Spec: `specs/ft227-adoption-smoke/spec.md`, `Status: staged`, awaiting sign-off
Gate: green at HEAD; the uncommitted files are prose and tickets only

## State

`/bench-write-spec` for FT227 ran to its sign-off stop. The spec and its three
tickets are written and the one review round accepted them at the first
iteration with ten partials, all folded (the `Verification log:` line names
them). `bench coverage --check` is green on 15 rows. Nothing is committed yet:
the reviewer has not signed off the spec-and-tickets pair.

The decision source is `roadmap/FT227.md`. The spec's grounding repro was
re-run in a scratch repository this session: `bench setup --yes`, sentinel
removed, `bench gate` through the installed wrapper — red on `HOME`, then red on
the empty inventory; green once a manifest declaring `BENCH_HOME` and `HOME`
existed and the stub's canary call was guarded on `tests/canary` existing.

Decisions that stay closed in the spec: the manifest is a `seed` kind (never
managed); the guard is the stub's, `bench canary` does not change; the sentinel
stays; the journey adopts `owner.repos[1]` and binds a private `BENCH_HOME` on
every launch; every green leg asserts the exact `gate: green` line.

## Next command

After sign-off: `bench commit -m "spec: stage ft227-adoption-smoke" -- specs/ft227-adoption-smoke capture/session-handoff.md`,
then a fresh mid-tier session runs
`/bench-implement-spec specs/ft227-adoption-smoke/spec.md`.

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
