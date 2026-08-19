# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `432a018`, clean tree, 12 unpushed commits
Spec: `specs/ft228-debug-restoration/spec.md` (Status: staged)
Gate: green at `ce1be93` — stale, work tree `a3e89fc`

## State

FT228 is implemented and reviewed. All four tickets plus one repair commit are
green on the retained integration source, addressed as `ft228-integration`:
frozen review base `c09a7d64`, source tip `fad39a89`. Three axes ran; Spec and
Coverage returned clean, Standards returned one finding, and the repair commit
carries it. Nothing of FT228 has reached `main` yet — `bench worktree land` is
the remaining step, then `/bench-final-check`.

Three unplanned fixes landed on `main` first (`ba23d877`, `23cdfc29`,
`432a018b`), under `specs/light-path-worktree-exec-kit-leak/`: `bench worktree
exec` leaked the caller's wrapper routing into a differently-rooted child,
`kit_dir` named the wrapper's tree rather than the tree being worked on, and a
reclaim test substring-matched a randomly quoted fingerprint. Together the
first two made any gate run inside a worktree of this kit drop its `race` and
`system` phases and report red over a clean tree. `dist/bench` was rebuilt and the
end-to-end proof ran: `bench worktree exec` now reaches a full-shape green gate
with no environment workaround.

Two open items for the reviewer, neither blocking:

`bench worktree land` refuses a light-path worktree — it requires a reviewed
spec's `spec.md` and ownership fence, and a light-path spec is a tickets-only
directory. The three fixes reached `main` by fast-forward instead.

Ticket 04's stale-row pass gives a narrow canary root a dozen diagnostics beside
its intended red. Coverage confirmed no fixture `EXPECT` collides with that
text, so none passes for the wrong reason; the repair commit states the
behavior where the pass is written rather than scoping the check, which would
make it guess which absences are real.

Undrained: `capture/learnings.md` carries the FT228 spec round's entry.

## Next command

`git push` — the board's leading invocable signal (`git`).

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
