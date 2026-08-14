# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `7ae73f01` plus the FT189 shaping artifacts
Spec: none staged.

## State

FT189 shaping is closed: `decisions/worktree-enumeration-hang.md` is `ready`
with all four tickets resolved. Reviewer decisions (2026-08-14): both arms —
pre-scan refusal (every entry under `<git-common-dir>/worktrees/` must be a
regular file or directory) plus a named `internal/bounds` deadline through
`bounds.Run`; scope is the sole enumeration owner `git.Worktrees` only; the
refusal names the entry path and shape, `bench doctor` reports the same
finding, and no Bench command deletes under `.git/worktrees/` (plan/apply
repair rejected). Mutation-site exposure is parked in `capture/IDEAS.md`.
Probe evidence: `decisions/assets/worktree-enumeration-hang-probe.md`.

`decisions/spec-build-review-gate-cadence.md` is invalid (its Sources cite the
removed `internal/specbuild/checkpoint.go`), which alone keeps `bench maps` at
exit 1; repairing it belongs to that map's own shaping resume, not FT189.

## Next command

`/bench-write-spec` — compile the ready `worktree-enumeration-hang` map into the FT189 build spec.

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
