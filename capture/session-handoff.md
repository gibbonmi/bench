# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `78bd15e`, clean tree, 3 unpushed commits
Spec: `specs/exact-prospective-landing/spec.md` (Status: staged), `specs/ft187-communication-surface-cut/spec.md` (Status: staged), `specs/ft194-project-green-desync/spec.md` (Status: staged), `specs/go-build-cache-footprint/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `097f56e` — current

## State

**Phase reached: `/bench-write-spec` closed for FT194.**
`specs/ft194-project-green-desync/spec.md` is reviewer-approved (2026-08-04)
and committed as `78bd15e`. It repairs the project-green desynchronization
that wedges spec-build promotion after an empty-run fast-forward.

- Decisions closed with the approval: recognition over republication (a marker
  that is an ancestor of both run base and destination is lagging, not
  conflicting; the fast-forward stays marker-free); recognition policy
  single-sourced in the gate-authorization owner, promotion's two publication
  sites routed through the `GateOwner` seam; lagging-aware red-attribution
  stays out of scope (priced in the spec's out-of-scope list).
- The spec's falsification pass ran and its findings are folded in; the
  coverage map validates at 11 rows.
- Build shape recorded in the spec: two tickets — authorization acceptance
  first (fence `internal/gate/authorization`), then publication rerouting
  (fence `internal/specbuild`, `cmd/bench/specbuild.go`). Recognition rows
  must drive the real authorization owner, not the marker-mimicking fakes.
- FT156's spec remains next in the roadmap's recommended sequence once FT194's
  build restores promotion.

## Next command

`/bench-implement-spec specs/ft194-project-green-desync/spec.md` — fresh
session on the mid tier (claude: `opus`).

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
