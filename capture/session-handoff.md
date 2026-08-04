# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `99b1c72`, 11 unpushed commits
Spec: `specs/recovery-discard/spec.md` (Status: staged) is the **active spec build**;
`specs/exact-prospective-landing/spec.md`, `specs/ft187-communication-surface-cut/spec.md`,
and `specs/pre-push-guard-visibility/spec.md` remain staged and untouched.
Gate: not run against the candidate — `promote` is the only gate boundary and has not run.

## State

**Phase reached: repair round complete, awaiting a second authoritative review.** Thirteen
tickets have integrated — eight building the feature, five repairing it. The candidate is
`38062e15`. All 17 acceptance coverage rows are green.
`bench spec build status recovery-discard` is the authority; the tree wins over this file.

- **The first review found a critical defect and it is fixed.** `--discard` accepted any
  ref that existed, so it would delete an ordinary branch. `PlanRecovery` now checks the
  recovery namespace before existence and emits a `foreign` verdict no verb authorizes.
  Also fixed: the recovery fingerprint now seals the plan's change summary; reclaim now
  enumerates prior terminal runs from `record.History`; an interrupted reclaim reports the
  refs it spent and converges on a fresh plan; the reclaim prose gained conformance
  anchors; the parser's operation set is single-sourced again with a bidirectional check;
  `CHANGELOG.md` records both verbs and the `--apply` behaviour change.
- The reviewer directed that **sol (`gpt-5.6-sol`, Codex top binding) is the authoritative
  review**, and its result becomes the receipt for `bench spec build review`. Invoke with
  `codex exec -m gpt-5.6-sol -c model_reasoning_effort="high" -c service_tier="fast"
  -C <repo> --dangerously-bypass-approvals-and-sandbox`, charged read-only.
- **Two findings are queued for that review, deliberately not built** — the repair-pass
  bound is one round then one review. First: `--discard` authorizes `RecoveryRetain`, a
  catch-all that also covers "could not classify" (ambiguous rows, `verifyRecovery`
  failures) — the same authority shape as the critical bug. Second: `CheckpointRef` is
  persisted without validation against `digest(Run + assignment ID)`, so a hand-edited
  state file could name another run's checkpoint.
- Sol's six judgment dispositions from round one stay closed unless the reviewer reopens
  them: it kept the lifecycle framing, the non-`recovered` refusal, the lifecycle-family
  exclusion, and the `s.resolve` bypass; it vetoed the misplaced grammar test and the
  unanchored prose, both of which the repair round fixed.
- **Deferred with evidence:** `git.DeleteBranchExact(root, ref, "")` deletes
  unconditionally — confirmed by sol in a throwaway repo, unreachable from this candidate.
  `ParseBuild`'s empty-args branch still hand-lists the operations as a literal string, a
  second restatement of the grammar table.
- The reviewer asked that the review's findings drive a follow-up hardening pass over
  ticket and spec prose **and** the kit's guidance surface, including whether tickets or
  spec scope should be smaller. Sol's round-one verdict: this should have been **two
  specs** (recovery discard and reclamation have disjoint package sets) and about **ten
  tickets**. A copy-paste prompt for that session was handed to the reviewer separately.
- `capture/learnings.md` holds one open entry (the ticket `Assumptions:` field). Five
  further authoring findings are with the reviewer in that prompt, not yet captured.

## Next command

The reviewer runs sol over candidate `38062e15`, carrying the two queued findings above.
Its result becomes the receipt for `bench spec build review recovery-discard --evidence
<receipt>`, and only then `bench spec build promote recovery-discard` — the single gate and
commit boundary. If the branch tip moves first, `promote` recomposes and that discards the
review, so review last.

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
