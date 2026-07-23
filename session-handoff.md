# Session handoff

Repository: /home/mgibs/workspace/bench — branch `main`. Everything below is
executable from a cold start; no conversation history is needed.

## State

- **FT87 slice 3 is spec'd and approved** (`795dad1`).
  `specs/cli-grammar-and-capability-evidence.md`, `Status: staged`, compiled from
  the closed `decisions/bounded-network-resource-cli.md` map (tickets #7 and #8).
  13 stories, 4 seams, 27 coverage rows — `bench coverage --check` green. The
  reviewer signed it off; the build has not started.
- **Two decisions that stay closed.** The grammar helper lives in
  `internal/usage` (the map named the boundary, not the package; `internal/usage`
  already owns CLI usage strings and cycles cleanly with `internal/toon`). The
  step-9 falsification pass was deliberately skipped: its only firing trigger was
  the Handoff's uncertainty flags, and all three belong to slices 1 and 2.
- **Story 7 deviates from the profile's cheap CLI-plumbing row on purpose** and
  routes to the mid tier, because it widens the `bench commit` attribution guard
  — the code that makes a green gate describe exactly the diff that lands.
- **The FT87 roadmap row names the spec path**, so `bench status` cross-checks it
  against the tree from here on.
- **Known advisory debt:** `bench structure` reports 10 violations (crowded
  `internal/adopt/` and `internal/contract/surface/`, plus seven over-length
  files). The gate is green with them and no roadmap row covers them.
- **Uncommitted:** `IDEAS.md` carries one parked idea awaiting a drain. Nothing
  else is dirty.
- **Unpushed:** `main` is ahead of origin. Pushing is the reviewer's call.

## Next command

`/bench-implement-spec specs/cli-grammar-and-capability-evidence.md` in a fresh
**mid-tier** session, interactively.

Not `bench shift`: stories 2, 7, 10, 11, and 13 route to the mid tier, so the
spec fails `craft-line`'s venue-routing test, which requires every story's line
to be cheap.

Two facts the build needs and would otherwise rediscover: `git add -A --
:(literal)<dir>` already stages a directory's children, so story 7's change
belongs in the block-check and `stagePlan`, not in staging; and a bare `t.Skip`
message is invisible under non-verbose `go test`, so story 9's helper must write
its structured line to stdout *before* skipping.

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
