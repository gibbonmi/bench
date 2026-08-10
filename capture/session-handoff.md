# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — FT173 restructure landed; tree clean at that commit
Spec: `specs/axi-spec-build-complete/spec.md` (Status: staged), `specs/axi-coherent-diff/spec.md` (Status: staged), `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged), `specs/axi-compatibility-oracle/spec.md` (staged, superseded — pending abandon)

## State

FT173 is restructured (reviewer decision 2026-08-10) from the five-spec
byte-preserving foundation to a three-spec forward build, landed as one
spec-only commit: the rewritten decision source
(`decisions/byte-preserving-axi-foundation/ft173-axi-contract.md`), the four
deleted foundation specs, the three staged replacement specs, and the
reconciled `ROADMAP.md`. The three new specs passed six cross-harness codex
review rounds (`gpt-5.6-sol` high, read-only), terminal verdict ACCEPT; all
pass `bench coverage --check`.

The superseded `axi-compatibility-oracle` run is terminal: abandonment applied
(two worktrees released, 137 provisional and 21 recovery refs marked). Its
superseded banner was dropped pre-abandon because lifecycle operations validate
the spec against the run's recorded subject; the supersession record lives in
`ROADMAP.md` FT173. `reclaim` is deferred until `axi-spec-build-complete`'s
tickets are cut, in case one salvages the bounded process observer or registry
census from candidate `9639a81d` — otherwise reclaim and let it go. The spec
folder stays until reclaim completes (reclaim validates recorded spec bytes
like abandon did); `bench spec retire` refuses staged specs, so the folder then
leaves through an ordinary deletion commit.

Build order: `axi-spec-build-complete` (lands `internal/axi` actions +
`help[]`), then `axi-coherent-diff`, then `axi-query-disclosure` (its QD5
harness-log leverage review is the build's first ticket and a reviewer-signed
checkpoint). FT185 composition stays closed: promote's gate payload is
composed when FT185 exists, never re-derived.

## Next command

`/bench-implement-spec` — `specs/axi-spec-build-complete/spec.md`

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
