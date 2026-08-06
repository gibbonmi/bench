# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — FT195 promoted at `af02997`, then `spec-retire: go-build-cache-footprint` and a capture commit; clean tree, ~19 unpushed commits
Spec: `specs/ft187-communication-surface-cut/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at the capture commit — current

## State

**Phase reached: FT195 (go-build-cache-footprint) promoted, retired, and retro'd.**

The lifecycle run is terminal: candidate `caaa6fb0`, promotion squash
`af02997`, retro at `capture/retros/go-build-cache-footprint.md`. Three
review findings were left flagged (not accepted) on the run's receipts for
reviewer veto: overlapping malformed-selector tables in the posture package,
hostile-output refusals exercised only through artifact mode, and a cosmetic
prose reflow in a retired ticket. Reopen any by reviewer request.

Closed reviewer decisions from the run: the Go publication child owns the
sealed transaction (backup, atomic pair replacement, signal-safe restore);
the topology invariant is the one-entry all-package enumeration; artifact
promotion is recorded by staged-entry consumption, not a flag.

A reviewer-initiated kit change is queued for a separate session (Codex):
simplify craft-tickets sizing to thinnest-independently-green-slice. Its
decision source is the two 2026-08-06 sizing entries in
`capture/learnings.md`; the prepared prompt gates on this run being terminal,
which it now is. Four open learnings entries and one pending retro await the
drain.

Unpushed commits await the reviewer's push decision.

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
