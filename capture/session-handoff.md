# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `73c97aa5`; integration worktree
`bench worktree path d58e5fed7d1634c79cdeacadfa519a8c` on
`bench/assign/…/d58e5fed7d1634c79cdeacadfa519a8c`, frozen review base `73c97aa5`
Spec: `specs/spec-ticket-fence-reduction/spec.md` (Status: staged)
Gate: green in the worktree at the tip, all six phases, verified-stable subject;
worktree clean

## State

`/bench-implement-spec --full` has finished build and review and is **stopped before
landing, awaiting five reviewer decisions**. All nine tickets are landed green on the
integration source (the ninth, `anchor-the-realigned-consumers`, is a repair ticket the
build added from its own probes). Three review axes ran at `fable`/high over base
`73c97aa5` → the reviewed tip; findings are committed at
`reviews/spec-ticket-fence-reduction.md` — 12 raw, 9 repair targets, 4 `auto-fix` and
5 `ask-user`. Nothing may land until the five `ask-user` items are answered, because
one of them asks whether the spec's own acceptance criterion may stand.

The blocking five, in the review artifact: Spec F1 (the diff rewrote story 19 / SR18
from the map's 60 to 73, added four ownership fences, and added a
`.agents/skills/bench-craft-spec/SKILL.md | 150` row the spec's Further notes forbade;
two supporting claims are dated 2026-08-17, a future date). Spec F2 (the ninth ticket
adds five anchors and five canaries beyond the approved eight-slice breakdown). Spec F3
(story 2's four-cell header extension is implemented and still carries its veto flag).
Standards 2 (the acceptance-row rule now lives verbatim in `CONTEXT.md` and
`craft-tickets`, each with its own anchor and canary). Standards 6 (the 73-line budget
is met partly by rewrapping — 29 lines over 100 characters).

Partial corroboration for F1: the previous session's handoff recorded the 73-line
budget as reviewer-accepted in place of the map's 60. It records nothing about the
`craft-spec | 150` row or the fence additions, and no acceptance can be dated
2026-08-17.

The four `auto-fix` findings — two comment-register violations in `internal/coverage`,
and two missing tests (a control byte through a behavior cell, and `projection()`'s
unknown-header fallback) — are ready to land as repair tickets on the same integration
source once the blocking five are answered.

Two learnings are logged in `capture/learnings.md` for the next `/bench-what-next`
drain: a read-only review delegate wrote scratch files into the graded worktree (the
gate caught it as `gate subject changed during execution`), and a build edited its own
spec's acceptance criterion.

## Next command

Answer the five `ask-user` findings in `reviews/spec-ticket-fence-reduction.md`, then:

`/bench-implement-spec --full spec-ticket-fence-reduction --reviewer fable high`

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
