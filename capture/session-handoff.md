# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `f2a925ee`, drain batch still dirty, 0 unpushed commits
Spec: `specs/gate-test-concurrency/spec.md` — staged, three tickets written and review-passed
Gate: green at `f2a925ee`

## State

**Phase reached: spec staged and committed; implementation is the next phase.**

`gate-test-concurrency` (gate-budget #23) is staged with tickets
`inject-kit-root-below-entries` → `retire-fixture-kit-pins` →
`adopt-gate-test-parallelism`, one `internal/gate` fence, sequential blockers,
covering map rows KC1–KC5. Both the spec falsification pass and the
ticket-breakdown review ran on the mid tier; all findings are folded in.
Reviewer decisions that stay closed: route one (kit-root injection seam plus
`t.Parallel`) over a package split; the chdir tests keep their chdir and stay
serial (#22 answer amended); serial cuts land after this spec (#25); the
census re-run (#26) gates #8.

The implementation session also lands gate-budget **#24** — test-only
`t.Parallel` in `internal/specbuild`, its two `t.Setenv` tests staying
serial — as a light-path ticket before or alongside the spec build; its ticket
text derives from the map's #24 entry, and it deliberately lives outside
`specs/gate-test-concurrency/tickets/` so the lifecycle's totality never sees
it.

Still pending from an earlier session, untouched by this one: the roadmap
drain batch (ROADMAP.md, capture files, the `pre-push-guard-visibility` spec
deletions) sits uncommitted awaiting reviewer approval; on approval it lands
in one green commit whose subject ends `spec-retire: pre-push-guard-visibility`.

## Next command

`/bench-implement-spec --full specs/gate-test-concurrency/spec.md`

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
