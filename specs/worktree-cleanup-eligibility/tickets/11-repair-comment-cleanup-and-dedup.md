# Repair: stale-comment cleanup and landedness dedup

Blocked by: 10-repair-landed-set-refusal-ownership.md
Writes: internal/worktree/eligibility.go, internal/worktree/subshell.go, internal/worktree/classifier.go, internal/worktree/landed.go, internal/worktree/eligibility_test.go, internal/worktree/clean_branch_test.go, specs/worktree-cleanup-eligibility/tickets/11-repair-comment-cleanup-and-dedup.md

Source: reviews/worktree-cleanup-eligibility.md, Standards findings 1, 2, 4; Spec finding 2.

## What to build

Two mechanical fixes, both low-risk, do last so comments describe the final
post-repair state rather than needing a second pass:

1. **Stale/falsified comments.** Fix every comment asserting a mechanism ticket
   04+ or ticket 10 (the prior repair) already replaced — e.g. any remaining
   claim that `PlanAutomatic` parses `plan.landed` by prefix, or that "the
   automatic classifier reads only that string." Reread the current code before
   rewriting each — describe what it does now, not its migration history. Drop
   the ~20 stale subshell.go/classifier.go line-range citations in
   eligibility_test.go's comments; keep symbol names, not line numbers.
   Rewrite PR-talk-register comments (narration of the change, "exactly as it
   always has," "where the pre-refactor code left it") into comments describing
   the current invariant, per `bench-craft-comments`. Leave the
   "characterization must move first" warnings on decideExplicit/decideAutomatic
   — those are legitimate, keep them.
2. **Dedup the "proven landed" check.** `kind == landednessProven && landed` is
   spelled out independently at two call sites. Add a `landedness.provenLanded()
   bool` method and use it at both.

## Acceptance

- [x] No comment in the touched files describes removed/superseded behavior.
- [x] No PR-talk register violation remains (no narration, provenance, or self-justifying argument in a comment).
- [x] `landedness.provenLanded()` exists and both call sites use it.
- [x] All existing tests remain green with unmodified assertions — this ticket touches comments and one dedup method only, no behavior change.
