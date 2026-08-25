# Take an explicit home in every verb

Line: opus / high.

The change hoists the one `BENCH_HOME` read to thirteen verb boundaries and
carries the home into every child the verbs start. A leak reaches the
operator's real home, so the top of the mid tier applies.

Blocked by: 11-mark-harness-tests-and-close-census.md
Writes: internal/worktree/*.go, cmd/bench/main.go, internal/sessioninspect/sessioninspect.go, internal/harness/worktree.go, internal/shift/loop.go, internal/shift/session.go, internal/status/status.go, internal/dashboard/dashboard.go, internal/worktree/*_test.go (compile fixes only)

## What to build

`effects.go` exports `Home()`, the one `BENCH_HOME` read, in place of
`benchHome`. Every public verb gains a `home string` parameter after `root`.
The verbs are `CreateCommand`, `ReleaseCommand`, `CleanCommand`,
`ReclaimCommand`, `ListCommand`, `PathCommand`, `ExecCommand`,
`ReauthorizeCommand`, `ResumeCleanCommand`, `LandCommand`,
`ResumeLandCommand`, and `PoolCommand` or `LeaseFileCommand` where it reads
the home. Each caller resolves the
home once at the command boundary through `worktree.Home()`, as it resolves
the root.

Some package functions read the home and serve another package:
`Pool`, `Acquire`, `Create`, `ConservativeCleanup`,
`ClassifyRegisteredWorktrees`, `RetainAndLock`, and `Release`. Each keeps its
form for that caller and gains an `At` form that takes the home. A function
this package alone calls takes the home outright. Inside the package, no function below the boundary calls
`Home()`; the effect census keeps that rule.

A child a verb starts receives `BENCH_HOME=<home>` on its environment, set
from the resolved value: the landing's gate plan and broker children, the
`bench` child `exec` starts, and any git child whose command reads the pool.
Read `internal/gate`'s plan environment seam before you set it; compose the
existing merge, do not invent one.

The tests compile by passing the home they already bind. They keep their
binds in this ticket; the fixture ticket removes them.

## Acceptance

- [ ] WF15 pins `CreateCommand` on an explicit home while the process `BENCH_HOME` names a different directory.
- [ ] WF16 pins a landing whose gate script records `$BENCH_HOME`, and the record names the explicit home.
- [ ] `bench worktree list`, `clean --help`, and `create` print the same bytes and exit codes as before.
- [ ] The effect census stays green: no file below `effects.go` reads the environment.
