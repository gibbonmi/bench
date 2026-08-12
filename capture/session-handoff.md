# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `a354ced`, 3 dirty paths, 41 unpushed commits
Spec: `specs/axi-coherent-diff/spec.md` (Status: staged), `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `c93bc38` — stale, work tree `bf91238`

## State

Pocock-guidance-doctrine (Spec C) tickets 01–09 are landed:
`craft-domain`, the TDD/seams reference leaves, and `prototype`;
frontier-round grilling and the light-path TDD seam gate; slim
`craft-tickets`/`bench-implement-spec`; slim `craft-delegate` and
`craft-line`'s owned-red convergence; re-derive-then-compare review with
dispositions and committed pickup; `.bench/BENCH.md` at 150 lines with
`AGENTS.md`'s shell rules; the spec-reread rule, seams read budget, and
`/bench-debug` entry isolation; and the fail-closed prose-budget conformance
check with its canary. Ticket 10 (whole-artifact reread and this record pass)
is in flight in a write-delegate worktree, diff ready, not yet committed.

Closed decisions that stay closed: byte-identical and divergent foreign
`bench link` adapter targets keep their landed hard refusals — ticket 10 does
not touch adoption behavior. The `tests/canary/guidance-prose-budgets` fence
addition, `prototype`'s missing `disable-model-invocation` flag, and
`.bench/BENCH.md`'s ~135-column wrap are flagged veto items from the review
pass, not open questions — they ride the batch approval unless vetoed.

Ticket 10's sweep found dangling cross-references to retired ticket-schema
ceremony outside its ownership fence — reported, not patched:
`.agents/skills/bench-craft-spec/SKILL.md` points at a `craft-tickets`
section (`Discover the contracts before writing files`) that no longer
exists; `.agents/skills/bench-craft-review/SKILL.md` cites a `craft-tickets`
"prefactor rule" and "integration-surface discovery" that no longer exist
there; `.agents/commands/bench-debug.md` still describes the retired
delegate "debug receipt", "assignment ID", and an `assign --refresh` verb,
and points at a `bench-implement-spec.md` section ("When a delegate is
blocked outside its fence") that the ticket-04 slimming removed.
`.agents/commands/bench-shape-idea.md` names the "Prototype" decision-ticket
type but never points at the `prototype` skill by name — a discoverability
gap, not a contradiction. These four files are outside ticket 10's fence.

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
