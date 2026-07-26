# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — FT91 merged, green, **unpushed**.
Spec: `specs/ft91-gate-tier-split.md` — `Status: implemented`.
Review findings: `reviews/ft91-gate-tier-split.md` — **9 still open**.

## State

- **FT91 landed.** The gate now splits into a dev tier and a ship tier
  (`bench prep-release`), with per-check timing printed on every run. Whole gate
  5m50s → 4m36s; conformance phase 328.8s → 106.9s.
- **`/bench-review-implementation` ran and found 12 issues; 3 are closed.** The
  inner canary gate now pins the tier being swept — before the fix the two
  release-evidence fixtures could never bite, so `bench prep-release` could not
  have reached green. Also closed: an untiered registry check silently left the
  dev gate, and an empty `CHECK` file now names that condition.
- **The next work is a gate refactor, not more FT91.** The reviewer rejected the
  ship canary's cost outright. The measured picture: `package-core-guard` is
  ~86s of the gate because it runs seven toolchain steps serially inside one Go
  test function, and inside the conformance phase every canary fixture runs all
  16 registry checks to observe the one its `CHECK` file names. `decisions/gate-concurrency.md`
  reaches the same conclusion independently — at k≤2 the long pole is the
  conformance phase, and further canary tuning cannot move wall-clock.
- **`bench prep-release` is shelved,** and cannot reach green on this host anyway:
  a real data race in `internal/guards.Scan`, and `govulncheck` is not installed.
  Both are parked in `IDEAS.md`.
- **Decisions that stay closed:** ship is a superset of dev; `internal/conformance`
  is excluded from the unfiltered inner run at both tiers; the release-only package
  tests are owned by the ship-tier conformance run; `prep-release`'s verify-mode
  evidence deliberately cannot satisfy `VerifyPublishAuthority`.

## Next command

`/bench-shape-idea` on the gate pipeline refactor. A fresh-session prompt for it
was generated at the end of the previous session; if it is lost, the inputs are
the last three entries in `IDEAS.md` plus `decisions/gate-concurrency.md`.

Before that: `main` is unpushed, and the branch/worktree sweep was proposed but
not executed — 23 non-`main` branches and 19 worktrees, all of whose work is
verified present in `main`. Deleting them is the reviewer's call.

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
