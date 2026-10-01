# Joins seam classification (FT356 ticket 9)

Consumed by: decision ticket 10 and the seam reduction spec.
Drift: a field added to or removed from `joins`, or a change to a cited default or test site.
Retire when: the seam reduction spec lands.

## Recommendation, scope, and evidence status

Keep 14 of the 29 joins fields. Remove 13 fields, and move `now` and `home` to the
verb entry. The scope is every field of `joins` in `internal/worktree/joins.go:27-88`
at commit `8ab52861`.

Five read-only delegates read the code, one for each field family. The coordinator
re-opened the load-bearing citations and ran the git probes in the scratchpad. No
delegate ran a test, so each verdict is a claim until the spec build shows its red.

Five verdicts are contestable, and ticket 10 decides them: `restoreClean`,
`creationLockAttempt`, `advanceLandingMarker`, `buildSubject`, and `resetEnvelope`.

## Facts: the classification

KEEP-a means the field injects a fault or an interleaving that no real fixture
produces. KEEP-b means its default runs a compile or a gate run. GO means real
fixtures replace it. Confidence is the delegate's frozen score from 0 to 10.

| Field | Default | Test sites | Verdict | Conf | Reason |
|---|---|---|---|---|---|
| landReviewed | `landing.New().LandReviewed` (joins.go:94-96) | 5 | KEEP-b | 9 | The default runs the gate through `gate.ExecuteTree` (landing/landing.go:254, gate/authorization/authorization.go:151). |
| advanceLandingMarker | `authorization.AdvanceMarker` (joins.go:97) | 12 | GO, contestable | 5 | A fixture gate script can move `refs/bench/green/<branch>` inside the window, so the swap fails. |
| reconcileLanding | `reconcileLandingDestination` (joins.go:98) | 4 | GO | 7 | A nested repository in the destination fails the residue guard (land_identity.go:322-333). |
| releaseLandingAssignment | `releaseCommandWith` (joins.go:99) | 15 | KEEP-a | 6 | A sibling created after publication needs the published hash (land_empty_sibling_test.go:26). |
| pruneLandedBranches | `intent.PruneUnclaimedLandedBranches` (joins.go:100) | 1 | GO | 7 | A stale `refs/heads/<name>.lock` fails the branch delete. |
| authorizeLandingSource | `preflight.AuthorizeReviewedSource` (joins.go:101) | 4 | GO | 8 | No fault, no compile, no gate; `publicLandingFixture` supplies the fences. |
| cleanupBoundary | nil `Fault` (joins.go:38-40) | 31 | KEEP-a | 9 | It fails a fixed transaction step inside a lock-held window. |
| cleanupLockAttempt | no-op (joins.go:102) | 3 | KEEP-a | 7 | Only the hook signals that the second applier passed the receipt read (resume.go:109-135). |
| creationLockAttempt | no-op (joins.go:103) | 1 | KEEP-a, contestable | 6 | The hook makes a removed-lock mutation fail every time, but it does not change the outcome. |
| claimTakeoverGap | no-op (joins.go:104) | 2 | KEEP-a | 9 | The judged bytes must change between the read and the rename (lifecycle.go:135-141). |
| claimStealGap | no-op (joins.go:105) | 1 | KEEP-a | 9 | The slot exists only after this process's own rename (lifecycle.go:141-149). |
| restoreClean | `restoreCleanCheckout` (joins.go:106) | 1 | KEEP-a, contestable | 5 | A claim must run mid-cleanup; probe P5 shows a `post-checkout` hook can do that. |
| chmodPool | `os.Chmod` (joins.go:107) | 3 | KEEP-a | 8 | No temp-directory fixture lets the Lstat pass and the chmod fail (pool_root.go:27-43). |
| ignoredLstat | `os.Lstat` (joins.go:108) | 1 | GO | 5 | Probe P6 shows git lists a file whose Lstat then fails. |
| resolveRunningBinary | `os.Executable` (joins.go:109) | 2 | KEEP-a | 8 | The fail-safe branches need an unresolvable or unstattable running binary (live_binary.go:40-56). |
| liveBinaryWarnings | `os.Stderr` (joins.go:110) | 2 | GO | 6 | The verb entry threads its own stderr down instead. |
| planLandedExplicit | `planExplicitWith` (joins.go:111) | 1 | GO | 8 | It is a spy; the test already asserts the shape reason (clean_landed_hostile_test.go:129). |
| reauthorizeUnlock | `unlockWorktree` (joins.go:113) | 1 | GO | 6 | Probe P1 shows a write-denied admin directory fails the unlock. |
| reauthorizeLock | `lockWorktree` (joins.go:114) | 1 | KEEP-a | 8 | The fault fails only the new-reason lock after a good unlock (reauthorize.go:131-138). |
| reauthorizeBeforeCAS | nil (absent from joins.go:92-123) | 1 | KEEP-a | 9 | The ledger lock holds across the decision, so no real writer wins the swap (intent/transaction.go:51-60). |
| mergeLane | `gate.LaneForCommit` (joins.go:115) | 7 | GO | 8 | Its only seam reason is the `BENCH_KIT` read (gate/phases.go:224-229). |
| kitSourceCheckout | `gate.KitSourceCheckout` (joins.go:116) | 2 | GO | 8 | Its only seam reason is the `BENCH_KIT` read (gate/kit_source.go:15-24). |
| mergeReconcile | `reconcileMergeCheckout` (joins.go:117) | 2 | GO | 7 | A planted `index.lock` fails `reset --merge` (merge.go:35-38, probe P3). |
| resetMove | `moveResetCheckout` (joins.go:118) | 3 | GO | 6 | Probe P2 shows a stale `HEAD.lock` fails `symbolic-ref`. |
| resetLayers | `restoreResetLayers` (joins.go:119) | 2 | KEEP-a | 7 | The recapture mismatch needs a write between the restore and the recapture. |
| resetEnvelope | `writeResetEnvelope` (joins.go:120) | 1 | KEEP-a, contestable | 6 | Only a ref move between write and verify trips the check; a `reference-transaction` hook might do it. |
| buildSubject | `runbinary.BuildSubject` (joins.go:121) | 5 + about 17 | GO, contestable | 6 | Five tests already run the default against stub build scripts (build_test.go:118-136). |
| now | `currentTime` (joins.go:122) | 2 | verb entry | — | Ticket 4 moves the clock read to the verb entry. |
| home | `Home()` (joins.go:112) | 0 | verb entry | — | Ticket 4 moves the home read to the verb entry. |

## Tested results: the git probes

The coordinator ran these probes as uid 1000 with git 2.43.0 in a scratch repository.
The guard hooks refuse `git reset` and `git update-ref -d` by command text, so the
coordinator did not run those two commands.

| Probe | Setup | Result | Supports |
|---|---|---|---|
| P1 | admin directory at mode 0500, then `git worktree unlock` | exit 255, "unable to unlink" | reauthorizeUnlock GO |
| P2 | stale `HEAD.lock` in the admin directory, then `symbolic-ref HEAD` | exit 1; `rev-parse HEAD` still exits 0 | resetMove GO |
| P3 | stale `index.lock`, then `status` and an index write | `status` exits 0; `update-index --refresh` exits 128 | mergeReconcile GO, by analogy |
| P5 | `post-checkout` hook in the common directory, then `switch --detach` in the worktree | the hook ran | restoreClean GO is possible |
| P6 | ignored file in a directory at mode 0600 | `ls-files --others --ignored` lists it; `stat` fails | ignoredLstat GO |

git 2.43.0 has no `--show-ref-format`, so its ref backend is files. A ref lock file
is therefore a valid fault for `pruneLandedBranches`, but no probe ran the delete.

## Inferences

The landing fields couple through one helper. `stubLandJoins` returns a fake commit
(land_fixtures_test.go:181-186), and that fake forces stubs for four other fields.

```text
stubLandJoins: landReviewed -> fake commit "aaaa..."
  -> advanceLandingMarker, reconcileLanding, authorizeLandingSource,
     releaseLandingAssignment must be stubs, because no real commit exists
GO for any of the four -> the test moves to publicLandingFixture
                          with the real landReviewed and the shell gate
```

The `advanceLandingMarker` GO route needs the real `landReviewed`, because only the
fixture gate script runs inside the marker window.

## Contradictions

- The roadmap row names the merge reconcile and the build as seams to keep. This
  research finds both fields GO. The cleanup boundary stays.
- The `mergeReconcile` comment says no fixture can make a bare reset fail
  (joins.go:71-73). Probe P3 shows an index lock fails an index write.
- The `planLandedExplicit` comment says the seam avoids the real fixture
  (joins.go:50-51). The only test passes through to the real planner.
- `clean_landed.go:21-23` declares `planLandedExplicitWithOptions`, and nothing
  reads it. It is a dead second copy of one seam.
- `reconcileLandingDestination` passes its joins argument to a function that
  ignores it (land_resume.go:119-123), against joins.go:25-26.

## Defects found

- `resumeCleanCommandWith` never sets `j.home` (resume.go:488-497). Its cleanup
  drops census records under `Home()`, not under the verb's home. Production
  passes the same value, so the defect is latent. The binds at
  resume_test.go:38 and :89 hide it.
- `resumeCleanCommandWith` reads `currentTime()` at resume.go:497, and its callee
  reads it again at resume.go:371 and :375. One verb reads the clock twice.
- `defaultJoins` reads `Home()` eagerly (joins.go:112). Each `defaultJoins()` call
  below a verb entry is a hidden home read: worktree.go:656, subshell.go:67, :81,
  :140, :154, and list.go:34.

## The kit root and the census edge

Two kit-root derivations exist, and their fallbacks differ on purpose
(gate/kit_source.go:10-14). `kitRoot` falls back to the graded root
(gate/phases.go:224-229). `KitDir` falls back to the executable's parent
(gate/kit_source.go:15-24). A single read at the verb entry must carry the raw
`BENCH_KIT` value, and each consumer keeps its own fallback.

`gate.LaneForCommit` has a second caller outside the package (commit/commit.go:145).
`KitDir` has six (command_registry.go:263, adopt/doctor.go:278 and :419,
adopt/link.go:51, doctor_rows.go:75 and :176).

Six exported functions are not verb entries, and they read `Home()` or
`currentTime()` themselves. They are `Create` (ownership.go:137), `Acquire`
(lifecycle.go:173), `PlanAutomatic` (classifier.go:308), `ClaimRecordedLease`
(snapshot.go:49), `ClassifyRegisteredWorktrees` (worktree.go:79), and `Pool`
(worktree.go:25). Packages `shift`, `harness`, and the dashboard call them.

The two `BENCH_KIT` binds (land_journey_test.go:364, exec_test.go:28) grade child
processes, so a kit parameter drops neither one. Of the 11 `BENCH_HOME` binds,
about four drop after the move (worktree_test.go:535, resume_test.go:38 and :89,
list_actions_test.go:182 and :242).

## Tests that cannot convert

These tests matter only if their field goes.

- `buildSubject` WF7, `TestBuildCancelExitsOneHundredThirty` (build_test.go:208).
  The cancel comes from a process signal (subprocess/cancel.go:25), which reaches
  every parallel test. The delegate proposes a direct `buildExitCode` test instead.
- `buildSubject` sealed-broker rows (land_effects_test.go:249, :358, and
  land_effects_cleanup_test.go:214, :272). They convert only if
  `freshness.Digest(root)` stays stable across the landing; that is unknown.

The tests that cannot convert under a KEEP field stay as they are:
land_flags_test.go:157, the ordinary case at land_empty_sibling_test.go:26, five
`resolveRunningBinary` sites, and reset_restore_test.go:82.

## Proposals

- Fold `cleanupLockAttempt` and `creationLockAttempt` into `Fault` steps on the
  existing fault channel. Fewer fields then carry the same interleavings.
- Rename `liveBinaryWarnings` if it stays as a writer, because retirement also uses
  it (lifecycle.go:462-466).
- Convert the spies at land_flags_test.go:202, :255, :322,
  land_identity_test.go:25, and land_reauthorization_test.go:206 to assertions on
  real state, as land_identity_test.go:108 and land_reauthorization_test.go:230 do.

## Unknowns

- Whether `git reset --merge` fails on a stale `index.lock`. Probe P3 tested an
  index write, not the reset itself.
- Whether `git update-ref -d` fails on a stale ref lock file.
- Whether a gate script that moves the green marker keeps the green verdict
  (gate/authorization/authorization.go:150).
- Whether `freshness.Digest(root)` stays stable across the broker landing.
- Whether Bench's own ignored listing flags (layers.go:24-25) list a file in a
  no-search directory, as probe P6 did with the standard flags.
- Whether the test suite ever runs as root, which defeats the chmod fixtures.

## Validation plan

1. The seam reduction spec gives each GO field one acceptance row for each test
   site. The build shows each converted test red under a named mutation.
2. The build runs the five unknown git and fixture behaviors as `bench probe` runs
   before it removes the matching field.
3. If a probe fails, the field stays KEEP-a, and the build records the change for
   reviewer veto.

## Verification record

- [x] The output opens with the recommendation, the scope, and the evidence status.
- [x] Facts, tested results, inferences, and proposals sit in separate sections.
- [x] One synthesized section answers the one question of ticket 9.
- [x] The classification table records each option's consequence with a citation.
- [x] Every contradiction and residual unknown is kept.
- [x] A diagram shows the landing-field coupling.
- [x] The unverified items are stated.
- [x] The output ends with a validation plan.
- [x] Each load-bearing claim cites a path and line at `8ab52861`.
