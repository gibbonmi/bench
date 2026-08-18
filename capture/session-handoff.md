# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `d884020`, 8 dirty paths, 0 unpushed commits
Spec: `specs/bench-front-door/spec.md` (Status: staged)
Gate: green at `764af7e` — stale, work tree `d514682`

## State

**Reviewer override standing for this repo:** the 2026-08 capability audit's own priority
order (`docs/audits/2026-08-bench-capability/results-fable-high/proposed-roadmap.md`)
supersedes `ROADMAP.md`'s `## Recommended sequence` until A1–A11 are exhausted; recorded
in `ROADMAP.md`. A1 and A2 are landed.

**A3 spec staged, sign-off given 2026-08-18.** `specs/bench-front-door/spec.md` (42
stories, 47 rows, seven serial tickets under `tickets/`) is the reviewed build spec;
decision source is `action-items.yaml` entry A3. Reviewer-visible calls the build may
still contest are collected in the spec's Further notes (the flagged `decisions: ready`
state, gate red → `/bench-debug`, inventory into the binary's `help` verb, staged-spec
severity tie rule, the destination-side "awaiting land" answer). Landed alongside it as a
light path: the stale-command sweep's `Introduces commands:` allowance
(`specs/light-path-introduces-commands-allowance/tickets/`), which is what lets a staged
spec name the phase it ships.

## Next command

`/bench-implement-spec specs/bench-front-door/spec.md`

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
