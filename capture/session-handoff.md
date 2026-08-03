# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `5e20751`, tree carries an uncommitted `capture/learnings.md` append (commit it with the next capture landing); the reviewer's parallel-landings grill session may dirty `decisions/` at any time — tree wins, and `bench commit` refuses out-of-set dirt, so name reviewer drafts into the commit only with their say-so (standing precedent from this run: they said carry them).
Spec: `specs/injected-interface-junctions/spec.md` — Status: staged, spec build run active
Gate: green at `5e20751` — current

## State

- Build `injected-interface-junctions`: nine assignments checkpointed and integrated (six feature tickets, one binding follow-up, two sol round-1 repairs). Candidate `e309505844460a29f48ac52c43a6ad7fff40d3ba` on base `9b7aa0a` — but HEAD moved with the round-2 repair-ticket commit, so the run needs recomposition before any assign.
- Sol reviews (codex `gpt-5.6-sol`, high): round 1 FAIL — five findings, all confirmed closed by round 2; round 2 FAIL — four findings, all fixed by the staged ticket `repair-spec-and-comment-consistency.md` (committed, unassigned): SP3 story-2 prose still claims a real-planner unreadable path; SP4 spec contracts four messages, omitting the orphan-row fifth plus its coverage row; S2 fixture comment names the wrong orphan port; S3 a comment embeds review provenance. Terminal repair-pass bound applies: round 3 verifies closure of those four and stops.
- The harness permission layer blocks the agent running `bench spec build promote`; the reviewer runs it via `! bench spec build promote injected-interface-junctions` (recompose now; the final publishing run later).
- Checkpoint receipts are coordinator-authored JSON; schema in `internal/specbuild/checkpoint.go` (validReceipt/validateReceipt). A working builder script: `/tmp/claude-1000/-home-devuser-workspace-bench/8186d7c5-2f61-41f5-9507-e6e38eeef978/scratchpad/build_receipt.py` (reads run state from `.git/bench/specbuild/<sha256(slug)>.json`, computes tree via throwaway index). Review receipt schema in `internal/specbuild/lifecycle.go:365` — version 1, run, candidate, exactly three axes Standards/Spec/Coverage, findings `{id,disposition}` ascending; disposition `accepted` blocks promote, use `resolved`/`endorsed` for closed findings; file must end with a newline, absolute path, outside any assignment worktree.
- Probe discipline (two burns this run): probe evidence must prove execution — `-count=1`, `-v` where skips are possible, `BENCH_CONFORMANCE_ROOT` set for scoped conformance runs; `(cached)` or 3ms greens are not evidence.
- Decisions that stay closed: the spec's stories and seams as reviewed; the round-1/round-2 finding closures; the unreadable-metadata row's pre-composition disposition; sol as final semantic reviewer; per-ticket lines as declared (repairs run `sonnet`/low).
- Retro material beyond learnings: `capture/parallel-session-friction.md` (seven commit/promote serialization points, feeding the reviewer's parallel-landings shape-idea).

## Next command

`! bench spec build promote injected-interface-junctions` (reviewer-run recompose), then in-session:
`/bench-implement-spec specs/injected-interface-junctions/spec.md` — resuming at: assign `repair-spec-and-comment-consistency.md`, one `sonnet`/low delegate, checkpoint + integrate, sol round 3 scoped to SP3/SP4/S2/S3 closure, submit the review receipt, reviewer-run promote, then `/bench-final-check`.

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
