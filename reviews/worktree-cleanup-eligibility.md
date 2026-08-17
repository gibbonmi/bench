# Review: worktree-cleanup-eligibility

Base `868a4e4e60f0f1b2de62e5ac685b485668634160`, tip `ca538cbbc0f5ee7b39f0885d4cda43d4a7b32860` (8 tickets, all committed). Three axes reviewed in parallel, opus/medium, per `craft-review`.

Raw findings: 12 (Standards 7, Spec 3, Coverage 2). De-duplicated repair targets: 6, after collapsing the comment-rot cluster (5 raw findings → 1 target) and the CO3/EV2 landed-set cluster (2 raw findings → 1 target).

## Standards

7 findings. Worst: comments that now assert facts the diff falsified.

1. **Stale/falsified comments** — `eligibility.go:44-46` ("PlanAutomatic still parses by prefix"), `subshell.go:232-238` ("the automatic classifier reads only that string"): both false since ticket 04 migrated `PlanAutomatic` to read `plan.landedTyped`. Plus PR-talk register violations (narration/provenance) in `eligibility.go`/`classifier.go` doc comments, and ~20 stale line-range citations in `eligibility_test.go`. → **auto-fix** (repair target: comment cleanup)
2. **Duplicated Code** — "branch is proven landed" derived independently at `eligibility.go:359` and `landed.go:12`, both spelling `kind == landednessProven && landed`, no shared method. → **auto-fix** (repair target: add `landedness.provenLanded()`)
3. **Lazy Element / Middle Man** — `automaticPreservationTrigger` forwards to `plan.preserves()` and claims to be "the one place that decides," but `lifecycle.go:428` still calls `preserves()` directly, and the actual Action/Reason selection for the landed-set refusal never moved. → **ask-user** (repair target: same as Spec #1 below — one target)
4. **Duplicated fixture harness** — `gatherExplicitFactsForTest` (~110 lines) reimplements `PlanExplicitWithOptions`'s gathering procedure rather than being an independently-authored *expectation* (the spec's sanctioned exception is narrower than this). → **ask-user**

## Spec

3 findings. Worst: CO3/EV2 not fully satisfied.

1. **CO3 / EV2 — landed-set preservation refusal's Action/Reason selection was not moved into the eligibility owner.** Spec: *"The `--landed` selector's 'per-path cleanup is required to preserve work' refusal is part of this owner rather than a post-plan action rewrite."* Only the boolean predicate (`automaticPreservationTrigger`) moved to `eligibility.go`; the triple-assignment `plan.Action, plan.ReasonCode, plan.Reason = ActionRetain, ReasonDirty, "..."` stayed in `clean_landed.go`'s `retainForLandedPreservation`, just relocated within the same file. ADR 0005's new sentence over-claims ("the landed-set selector … resolves that question through one ordered eligibility verdict"). → **ask-user**
2. Same stale-comment issue as Standards #1 (CO1-relevant angle). → **auto-fix** (same repair target as Standards #1)
3. **`CleanupPlan.landed` (string) is now write-only in production** — written at `subshell.go:229`, read only by the new EV1 test; the fingerprint hashes a local var, not the field. Two fields carry one fact. → **ask-user** (removal touches the EV1 assertion)

Coverage audit: 31 of 33 rows confirmed present and correctly mapped. Gaps: EV2, CO3 (finding 1 above). OS1 verified clean (all touched `_test.go` diffs are pure additions).

## Coverage

2 findings. Worst: a real, previously-undetected behavior regression.

1. **`--discard-branch` on a detached-HEAD worktree now wrongly authorizes branch deletion.** At base, the `DiscardBranch` override lived inside `if headRef != "detached"` (subshell.go). Ticket 03's refactor moved it outside that guard. With HEAD detached, `headRef == "detached"` (the literal string), so the plan now sets `deleteBranch=true`, `branchRef="detached"`, changes the fingerprint, and on apply calls `git.DeleteBranchExact(root, "detached", head)`. Reachable via `bench worktree clean --discard-branch <path>` and `--landed --discard-branch`. Not covered by any `TestDiscardBranch*` subtest, `eligibility_test.go`'s AU10, or any mutation probe run during implementation — all existing detached fixtures pass no `DiscardBranch` option. → **auto-fix**, high priority: restore the guard; add a detached+DiscardBranch regression test.
2. `landedness` zero value is `landednessDetached` (iota 0) rather than an explicit "unknown" — verified currently benign (every real construction path sets it explicitly or resolves equivalently to the pre-refactor empty string). → **no-op**

Edge inventory walked; no "Won't handle" item touched. No concurrency risk (value-copy semantics throughout).

## Disposition summary

| Repair target | Axis | Disposition |
|---|---|---|
| Comment cleanup (stale/falsified/PR-talk) | Standards, Spec | auto-fix |
| Add `landedness.provenLanded()` (dedup) | Standards | auto-fix |
| `--discard-branch` + detached HEAD regression | Coverage | auto-fix (priority) |
| Landed-set refusal selection not in eligibility.go (CO3/EV2) | Standards, Spec | ask-user |
| `gatherExplicitFactsForTest` duplicates production gathering | Standards | ask-user |
| `CleanupPlan.landed` dual representation | Spec | ask-user |
| `landedness` zero-value ambiguity | Coverage | no-op |
