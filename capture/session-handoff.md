# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `ca6c35e`, clean tree, 49 unpushed commits
Spec: `specs/axi-coherent-diff/spec.md` (Status: staged), `specs/axi-compatibility-oracle/spec.md` (Status: staged), `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/axi-spec-build-complete/spec.md` (Status: staged), `specs/checkpoint-scoped-review/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged), `specs/ticket-bundle-refusal/spec.md` (Status: staged)
Gate: green at `f6e02cb` — current

## State

`axi-spec-build-complete`'s run was abandoned this session (2026-08-11) on a
validated `/bench-debug` receipt: assignments `a69293c0…` and `a931c8ea…` were
mutually stranded — the CF2 canary fixture edit and its registry expectation in
`internal/conformance/fixture_bite_test.go` sat in opposite ownership fences,
so each worktree's focused `TestSpecTicketHandoffWorkflowFixturesAreComplete`
red was caused by the other's fence, with no lifecycle route to retire either.
The abandon retained 22 recovery refs (`bench worktree recovery`), including
both stranded worktrees' in-fence dirty work; candidate `9146b5ed…` and 25
integrated checkpoints are gone with the run.

Landed since:

- Spec revision committed (`4721c2d0`): story 5, SB9/SB10, and the one-fence
  proportionate ticket-evidence contract. Spec status stays `staged`.
- Combined SB9/SB10 ticket
  `specs/axi-spec-build-complete/tickets/adopt-proportionate-ticket-evidence-contract.md`
  committed (`ca6c35ea`), replacing the deleted
  `remove-mandatory-mutation-guidance.md` and
  `repair-guidance-canary-fixture-anchors.md`. Breakdown-reviewed pre-assign
  (read-only delegate); its two covers-honesty findings repaired (MG2→SB9,
  CF2/CF3→SB10). Blocked by `permanent-optional-ticket-inventory.md`.

Closed decisions: the abandon (receipt-backed, applied), the combined
one-fence ticket per spec.md line 40, and the spec revision content.
Open for reviewer: the round-8 P1/C1 tab/CR-vs-FR6 conflict from the abandoned
run's review receipt still needs a call before any re-landed repair; the new
learnings entry proposing abandon's staged-spec precondition exemption awaits
the `/bench-what-next` drain. Build order unchanged: `axi-spec-build-complete`
(restart), then `axi-coherent-diff`, then `axi-query-disclosure`; the
`checkpoint-scoped-review` and `ticket-bundle-refusal` specs (staged,
reviewer-signed) implement after the AXI runs in that coupled order.

## Next command

`/bench-implement-spec` for `specs/axi-spec-build-complete/spec.md` — restart
the build with the recomposed ticket set, recovering preserved work from the
retained recovery refs.

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
