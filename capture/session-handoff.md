# Session handoff

Repository: `8a4b5b1767e3a4edaad6e8fa8d4762ee-c7cd342784c033ee75f6f9003dc1dc7b` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/8a4b5b1767e3a4edaad6e8fa8d4762ee-c7cd342784c033ee75f6f9003dc1dc7b`
Branch: `bench/assign/8a4b5b1767e3a4edaad6e8fa8d4762ee/c7cd342784c033ee75f6f9003dc1dc7b` — HEAD `b2e4bec`, clean tree, 11 unpushed commits
Spec: `specs/module-size-split/spec.md` (Status: staged)
Gate: green at `393fc6b` — current

## State

Batch 2 of `module-size-split` holds tickets 03 to 12 as ten green commits on the
integration source `mss-batch2` (base `0585b31f`, tip `b2e4bec1`). Every ticket is
a same-package move with a proven identical body-line multiset, except ticket 07.
That ticket consolidates six throwaway-root builders into one `throwawayRoot` in
`internal/conformance/harness_test.go` and grants four dated caps. `bench structure`
fell from 71 to 56 issues. The review over that range returned Standards 0, Spec 1
(the R13 shortfall below), Coverage 3 raw findings, all `no-op`; no `reviews/`
artifact exists.

Ticket 13 stops as a material acceptance shortfall. Its premise is stale:
`internal/worktree/eligibility_test.go` has held 209 lines since before the spec's
base, so there is no split to make. Row R13 asks for at most 55 issues, but the
scoped files can reach only 56. The reviewer decides: amend R13 to 56 and close
ticket 13 as done-by-tree, or widen the scope by one file from the 400-to-700
remainder. Until then the spec stays `staged` and the batch lands spec-less.

Three calls are flagged for veto. Ticket 07 granted `checks_test.go` a cap at 692,
under the spec's over-700 phrasing, because R14 binds to the 400-line scan. Ticket
06's `buildLandingBinary` fixture never existed in the package. The `.bench/lines.env`
plant in `writeHostileSkillRoot` is not graded by any test; the base had the same gap.

One diff-owned red occurred. The build census pins `owner_test.go` by path, and
ticket 10's first cut moved the stripped-journey marker out. The delegate repaired it
inside its fence.

## Next command

`/bench-write-spec specs/module-size-split/spec.md`

## Shape

Rewritten in full at every phase close, pruned rather than accreted. A fresh
session pays for every line it reads cold; drop anything it would not act on.

Operational gotchas are placed by lifetime, not copied here. One that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes. One scoped to a build
belongs instead in that spec's coverage rows.

This file names at most when you'll hit one, never the command — a second copy
drifts from the source.

Keep the three sections above. **State** holds what is true now, including anything
uncommitted. **Next command** holds the exact harness-native invocation, not a
description of it. This section is the third.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
