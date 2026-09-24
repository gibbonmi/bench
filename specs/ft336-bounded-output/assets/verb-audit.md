# Verb audit: default response size

The spec author ran this audit on 2026-09-24 at `main` commit `a519ed97`. Each verb
ran through `bench worktree exec ft336-bounded-output -- sh -s` in the phase
worktree. A shell function captured the stdout and the stderr of each verb and
counted the lines and the bytes. The pool held 47 worktree rows at that time.

The cost of a verb is the number of requests multiplied by the context that
each response adds. So the audit records what the agent acts on in each
response, and the fix keeps that part in the same response. A fix that moves
the needed part behind a second call adds one request for each use.

| verb | exit | lines | bytes | what the agent acts on | FT336 disposition |
| --- | --- | --- | --- | --- | --- |
| `bench help` | 0 | 69 | 8348 | the grammar of the next verb | backstop spill |
| `bench status` | 0 | 7 | 373 | the lead signal and its next command | inside the bound |
| `bench worktree list` | 0 | 142 | 14489 | the target identity and its state | help slot rows, then backstop spill |
| `bench worktree path <target>` | 0 | 2 | 228 | the absolute path | inside the bound |
| `bench worktree --help` | 0 | 16 | 1324 | the grammar of one leaf | backstop spill |
| `bench worktree show <target> <rev>:<path>` | 0 | 117 | 6343 | the blob bytes | backstop spill |
| `bench roadmap` | 0 | 20 | 1985 | the top row and the drain state | backstop spill |
| `bench learnings` | 0 | 2 | 43 | the open entries | inside the bound |
| `bench maps` | 0 | 54 | 6330 | the frontier and ready maps | backstop spill |
| `bench guards` | 0 | 10 | 819 | a stale or unwired guard | inside the bound |
| `bench guards --brief` | 0 | 8 | 660 | a stale or unwired guard | inside the bound |
| `bench diff` | 0 | 9 | 475 | the base and the changed files | inside the bound |
| `bench coverage <spec>` | 0 | 92 | 15537 | the coverage state and a red row | backstop spill |
| `bench preflight review <slug>` | 1 | 17 | 535 | the red check and its next command | green summary line |
| `bench preflight build <slug>` | 0 | 18 | 535 | the verdict | green summary line |
| `bench anchors <path>` | 0 | 18 | 1216 | the needles that pin the path | backstop spill |
| `bench outline` | 0 | 339 | 47651 | the candidate seam | backstop spill |
| `bench consumers <symbol>` | 1 | 1 | 130 | the reference edges | inside the bound |
| `bench spec history <slug>` | 0 | 1 | 35 | the retire commits | inside the bound |
| `bench harnesses` | 0 | 7 | 366 | the harness cells | inside the bound |
| `bench commands --brief` | 0 | 3 | 31 | the liveness answer | inside the bound |
| `bench version` | 0 | 1 | 25 | the version | inside the bound |
| `bench cache` | 0 | 2 | 156 | the cache footprint | inside the bound |
| `bench structure` | 1 | 108 | 8617 | the first structural issue | backstop spill |
| `bench skills-index` | 0 | 0 | 0 | the drift, if any | inside the bound |
| `bench test --package ./internal/census` | 0 | 4 | 155 | the verdict | inside the bound |

## Verbs the audit did not measure

The audit did not run these verbs. Each one changes state, runs for a long
time, or needs an artifact that the phase did not have:

- `bench worktree create`, `release`, `reauthorize`, `merge`, and `land`
- `bench commit`, `bench gate`, `bench shift`, and `bench probe`
- `bench handoff`, `bench retro`, `bench idea`, and `bench learning`
- `bench preflight review --charge`, `bench preflight build --charge`, and `bench preflight evidence`
- each `--full` form

The dispatcher backstop covers each of these verbs, because the bound applies
to every public command. The census output record measures them after the
build.

## Finding

Eleven of the 26 measured verbs print more than 10 lines by default. The
largest is `bench outline` at 339 lines. The per-row help block of
`bench worktree list` holds 93 of its 142 lines.
