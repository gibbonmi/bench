# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — pre-commit HEAD `09dbc35`
Spec: `specs/ft156-anchor-registry/spec.md` (Status: staged), plus the four earlier staged specs (`exact-prospective-landing`, `ft187-communication-surface-cut`, `go-build-cache-footprint`, `pre-push-guard-visibility`)
Gate: green expected for a documents-only change; the spec commit is the verdict

## State

**Phase reached: `$bench-write-spec` closed over FT156.**
The ready map compiled to `specs/ft156-anchor-registry/decisions/ft156-anchor-registry.md`
and the spec stages the map's whole program as one spec, reviewer-pre-approved:
two pulled-forward red-capability closures (the generic section-scoped
`.bench/BENCH.md` fixture; the ticket template's `Contracts:`-line anchor), the
declarative registry extraction carrying `bench anchors <path>`, comment-strip on
the shared matcher, and the honest-strength prose sweep. The Codex falsification
pass (`gpt-5.6-sol` / high) ran; its accepted findings are already folded in.

Decisions that stay closed: stronger-than-substring matching stays deferred
(map, 2026-08-02); no structure-grant raise (`.bench/structure.budgets` is
reviewer-owned); registry data lives in a non-test package below the
conformance import edge. Story 5's `fable / high` line rides the leverage
override and is flagged veto surface.

## Next command

`/bench-implement-spec specs/ft156-anchor-registry/spec.md` — in a fresh
mid-tier session; the coverage map carries not-gate-observable rows, so this is
not a `bench shift` candidate.

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
