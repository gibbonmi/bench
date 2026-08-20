# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `8c8aca6`, 3 dirty paths, 2 unpushed commits
Spec: `specs/worktree-exec-run-binary/spec.md` (Status: staged)
Gate: green at `6d5bd7a` — stale, work tree `fc6abe7`

## State

**This session's build — `worktree-exec-run-binary`, reviewed and ready to land.**
The spec is staged at 22 rows. Three commits sit green on the retained
integration source, addressed by label as `worktree-exec-run-binary`:

- frozen review base `8c8aca66`
- source tip `99bf27ee` — `ebf1ac7b` (exec roots the child at the worktree's own
  wrapper), `290febc7` (the gate entry names its next action), `99bf27ee`
  (closes the review's three accepted findings)

Three-axis review ran and its pickup file is resolved and deleted. A Coverage
re-review of the repair returned zero findings. WX20 is demonstrated by a
controlled A/B: exec driven by the unfixed binary refuses with the reworded
message, exec driven by the fixed binary runs all six phases green. That
evidence still needs recording in the retro.

Main carries two spec commits (`8c8aca66`, `529956b7`) that the source
deliberately does not have — spec edits red `paths-authorized` inside a build
diff. The two file sets are disjoint, so they compose at land.

**Blocking the landing:** the destination is not clean. `capture/learnings.md`,
`capture/agent-performance/claude-models.md`, and the untracked
`capture/retros/` are the prior FT229 session's, uncommitted since roughly
05:00 and waiting on the drain. This session has not touched them.

Two learnings this session owes, still unwritten for the same reason: asserting
a repo convention from artifacts rather than from the check that grades them,
and reaching for raw `git merge --ff-only` because no Bench verb moves a
worktree to a new base. A third is worth capturing: a gate verdict reused from
cache proves nothing about the path under test, so evidence runs need `--fresh`.

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
