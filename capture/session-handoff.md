# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `8683312`, 1 dirty path, 0 unpushed commits
Spec: `specs/ft187-communication-surface-cut/spec.md` (Status: staged), `specs/go-build-cache-footprint/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `0807d43` — stale, work tree `00bcc1c`

## State

**Phase reached: FT195 ticket derivation complete; implementation has not started.**

The read-only seam audit at `8683312` found no spec drift. The build is three tracer tickets: sealed local publication and fresh gate verdicts form the first frontier; artifact-mode callers depend on sealed local publication. The closed spec decisions remain the two semantic builder modes, one compiled subject, sealed default publication through the freshness owner, unsealed artifact promotion without host execution, and `-count=1` only at the two executed gate-test owners.

## Next command

`$bench-implement-spec --full FT195 specs/go-build-cache-footprint/spec.md`

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
