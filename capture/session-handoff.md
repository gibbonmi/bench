# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `c3cb914`, 11 dirty paths (the FT227 drain's staged
edits plus this spec), 6 unpushed commits
Spec: `specs/ft228-debug-restoration/spec.md` — staged, awaiting sign-off.
Gate: stale (last green `cff00da`).

## State

`/bench-write-spec` for FT228 is authored and reviewed: 24 stories, a 22-row
coverage map (`bench coverage --check` green), four serial tickets, and a
2-iteration verification log. The round's learnings entry is appended to
`capture/learnings.md` (open, undrained). The FT227 drain's edits
(`CHANGELOG.md`, `ROADMAP.md`, roadmap detail files) are staged but
uncommitted from the prior phase.

Decisions closed in this phase, dated 2026-08-19: `$bench-debug` becomes
implicitly invocable on Codex — the yaml flips, the adapter description gets a
symptom-bearing trigger, and a per-phase invocation-policy table graded on
both harness surfaces records every other phase at today's posture.

Contestable calls left standing for veto at sign-off: the prose-budget row for
`.agents/commands/bench-debug.md` (story 12 / DP3) is a promise beyond the
decision source, flagged in the spec's Further notes; the ticket-03/04 split
of the trigger work follows the review round's quiz.

## Next command

Sign off on the spec-and-tickets pair, then in a fresh mid-tier session:
`/bench-implement-spec --full specs/ft228-debug-restoration/spec.md`

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
