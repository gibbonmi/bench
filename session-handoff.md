# Session handoff

Repository: /home/mgibs/workspace/bench — branch `main`. Everything below is
executable from a cold start; no conversation history is needed.

## State

- **FT87 slice 3 is built and green.** All 13 stories of
  `specs/cli-grammar-and-capability-evidence.md` are implemented across six
  gated commits (`108a170`, `aff7cb8`, `3fda6e3`, `e83ac02`, `d65b315`,
  `e52da23`). `bench coverage --check` is green on all 27 rows. The spec is
  still `Status: staged` on purpose — the landing commit and the
  `Status: implemented` flip belong to `/bench-final-check`, not to the build.
- **Semantic review has not run.** That is the next phase.
- **Two reviewer decisions made mid-build, both overriding the spec.** They are
  closed; do not reopen them, but the spec text still describes the old designs
  and should be corrected when the spec is next touched.
  1. **Capability evidence travels by side-channel file, not phase stdout.** The
     spec's "the gate aggregates from phase output" and "a collector tees that
     stream" are dead. `go test` without `-v` discards a passing or skipping
     package's stdout *and* stderr, so a collector teeing that stream would see
     nothing forever. The gate now sets `BENCH_SKIP_LOG` per run; skips append
     one atomic line; the gate reads it after the phases join.
  2. **A repeated flag is a usage error at exit 2.** This extends story 1's
     enumerated grammar, which said nothing about duplicates. Without it,
     routing `diff` and `roadmap` would have silently flipped two gate-asserted
     exit-2 contracts to exit 0.
- **Two calls left open for reviewer veto.** Neither blocks review.
  - The marker-wait conformance check grades only the *slow* deadline argument
    of package-qualified `WaitForTwoLegMarkers` calls; the fast leg is bounded
    by no named policy and the helper's own package tests a fake clock where a
    literal is the subject.
  - `capability.Capability`/`Environment` take a local `capability.TB`
    interface rather than `testing.TB`, so `internal/gate` can import the line
    shape without linking `testing` into `dist/bench`.
- **One history defect, unrepaired, reviewer's call.** Commit `c82ba1f` is
  labelled "capture: park the gate phase-timeout headroom idea" but also
  contains the entire stories 10–11 slice (649 insertions, 11 files). Cause: a
  bare `git commit` after `git merge --squash` commits the whole index. A later
  full gate ran green on that exact tree, so the content is verified; only the
  history is wrong. `main` is unpushed, so a split is risk-free.
- **Known advisory debt.** `bench structure` reports 14 issues, up from 10 at
  slice start; the new one is `internal/gate/` at 18 source files against a
  budget of 16, from the collector's two new files. Gate is green with them.
- **Four ideas parked this session** (`bench idea`), all found while building
  and none in scope here: the mutation-revert papercut against
  `block-dangerous-git`; the gate's conformance phase inheriting `go test`'s
  10-minute default with ~4 minutes of headroom; **real data races in
  `guards.Scan`** that fail under `-race` on `main` today and that the gate
  never runs; and `waitForPIDFile`'s hardcoded 2s literal deadline — the same
  defect class story 13 fixed, at a call site the spec did not name.
- **Three learnings logged** in `.bench/learnings.md`, all with proposed rule
  changes: a spec's load-bearing tooling claim went unverified into three
  stories; a delegate charge's verification list became the delegate's ceiling
  and let a canary regression through; and the squash-merge index trap above.
- **Unpushed:** `main` is well ahead of origin. Pushing is the reviewer's call.

## Next command

`/bench-review-implementation` in a fresh session.

Review base is `9732ebe` (the commit before the build began), so `bench diff`
resolves the whole slice. Two things worth the reviewer's attention beyond the
usual three axes: whether the side-channel transport is the right shape now
that it is built rather than merely decided, and whether the seven flat
subcommands newly routed through `usage.Parse` kept their output contracts —
their behavior changed in ways no coverage row pinned (trailing garbage now
rejected, help now uniformly exit 0).

Then `/bench-final-check` for the landing commit and the `Status: implemented`
flip.

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
