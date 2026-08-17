# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `7e5e3ae8`, 22 commits ahead of the session-start pin, working tree has only pending capture files uncommitted (see State)
Spec: `specs/worktree-cleanup-eligibility/spec.md` — implemented, landed at `7d799eccd2c6f6452557af8cda4d28fa710aea51`, then retired at `7e5e3ae8`; the folder no longer exists

## State

FT216 (worktree-cleanup-eligibility) is fully landed and retired. One in-process
eligibility module (`internal/worktree/eligibility.go`) now owns both the
explicit and automatic ordered cleanup decisions, plus the shared preservation
refusal both `PlanAutomatic` and the landed-set planner consume; ADR 0005 is
rewritten as resulting-state documentation. 8 build tickets + 4 repair tickets
landed serially on one retained integration source; three-axis review found and
closed a real regression (`--discard-branch` + detached HEAD) and a genuine
spec-compliance gap (CO3/EV2), both independently re-verified clean on
follow-up. FT217/FT218's remaining decisions were re-homed into
`roadmap/FT217.md`/`roadmap/FT218.md` before the spec's decision map was
deleted with retirement.

Mid-build, an unrelated baseline gate red (`TestWorktreeReauthorizeJourney`,
a `git init.defaultBranch=master` host-config dependency in a system-test
fixture) blocked every commit; diagnosed and fixed via `/bench-debug`, landed
separately at `868a4e4e` before FT216 work resumed.

Three files are uncommitted by design, left for a reviewed `/bench-what-next`
capture drain: `capture/learnings.md` (one new open entry — a resolved
`reviews/<slug>.md` left in the tree hard-blocks `bench worktree land`, not
just `bench preflight review`), `capture/retros/worktree-cleanup-eligibility.md`
(the full landing retro), and `capture/agent-performance/claude-models.md`
(refreshed with this landing's Sonnet/Opus evidence — `open-ai-models.md` is
untouched, no OpenAI models ran this session). None of the pre-existing 7
unresolved decision maps or the 62 structure issues `bench status` flags are
from this session — unrelated, pre-existing repo state.

One unrelated, explicitly user-approved change also landed this session:
`.claude/settings.json` now scopes non-bench personal/bundled skills off for
this repo via `skillOverrides` (commit `3cc98077`) — reversible, repo-scoped,
does not touch anything outside this checkout.

`dist/bench` must exist and be reasonably fresh for local `bench` CLI
resolution to work in this checkout — it is not a purely disposable artifact
here; rebuild with `go build -o dist/bench ./cmd/bench` if it's ever removed.

## Next command

`/bench-what-next`

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
