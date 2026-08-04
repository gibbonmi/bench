# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `ebbb1b9`, 9 unpushed commits
Spec: `specs/recovery-discard/spec.md` (Status: staged) is the **active spec build**;
`specs/exact-prospective-landing/spec.md`, `specs/ft187-communication-surface-cut/spec.md`,
and `specs/pre-push-guard-visibility/spec.md` remain staged and untouched.
Gate: not run against the candidate — `promote` is the only gate boundary and has not run.

## State

**Phase reached: implementation complete, review pending.** A `/bench-implement-spec --full`
run built `recovery-discard` through eight tickets. All eight are integrated and released;
the candidate is `5e910bb1`. Every one of the spec's 17 acceptance coverage rows is green.
`bench spec build status recovery-discard` is the authority; the tree wins over this file.

- The reviewer directed that **sol (`gpt-5.6-sol`, the Codex top binding) performs the
  final review**, and that its result becomes the receipt submitted through
  `bench spec build review`. No in-lifecycle review delegate was spawned. Nothing promotes
  until sol's pass exists. Concrete defects it finds re-enter as ownership-fenced repair
  tickets; contestable calls are flagged for veto.
- The reviewer approved the **top binding (`fable`)** for `align-lifecycle-surface-prose`,
  the guidance-prose ticket. Every other ticket ran at `opus` / medium.
- The reviewer also asked that sol's findings drive a follow-up hardening pass over ticket
  and spec prose **and** the kit's guidance surface (`craft-tickets`, `craft-spec`,
  `/bench-implement-spec`), including whether tickets or spec scope should be smaller.
  That pass rides `craft-synthesis` at the top binding and is proposal-only.
- **Two coordinator probes found unasserted safety properties** that the delegates' own
  probes could not reach, both since repaired and re-verified: the recovery fingerprint did
  not pin its domain-tag/effect-string authority change, and the reclamation fingerprint did
  not commit to each ref's disposition. Every delegate mutation was control-flow; both gaps
  were data/constant shaped.
- **Known compromise, deliberate:** RM5's grader (`TestReclaimGrammarBoundsToOneSlug`) lives
  in `internal/specbuild` although it drives `internal/spec`'s parser, because that was the
  only test-capable path inside its ownership fence. It carries a comment saying so.
- **Known gap, disclosed in the checkpoint receipt:** `align-lifecycle-surface-prose`'s AL1
  is done and was enumerated by sweep, but no conformance anchor pins the reclaim wording —
  reverting one document's sentence stays green. Adding an anchor needs
  `internal/conformance/docs_workflow_helpers_test.go`, outside that ticket's fence.
- **Deliberately unrouted:** `internal/spec/build_test.go` still names
  `TestParseBuildExposesExactlyEightOperations` and omits `reclaim` from its operation case
  table. Both are correct improvements a delegate had made and had to revert; the file sits
  outside every ticket's fence and widening one would have forced a recomposition mid-flight.
- `capture/learnings.md` holds one new open entry: the ticket `Assumptions:` field pools
  verifiable preconditions with standing-rule restatements, and its parser splits on commas
  only. Five further authoring findings are **not yet captured** — they are in the reviewer's
  hands as a separate prompt.

## Next command

`/bench-review-implementation` is **not** the route here. The reviewer runs sol over the
exact candidate composition `5e910bb1`, then this session converts its result into the
receipt for `bench spec build review recovery-discard --evidence <receipt>`, and only then
`bench spec build promote recovery-discard` — the single gate and commit boundary. If the
branch tip has moved, `promote` recomposes first and that discards any review.

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
