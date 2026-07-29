# ft91-artifact-suite — subject-split artifact contracts

Status: staged

Compiled from
`specs/ft91-artifact-suite/decisions/gate-critical-path.md`. The source record
already decided the package split, its performance purpose, the checks that
must remain, and the rejected alternatives. This phase normalized that record
into the map ticket and Handoff shape. The exact package names, test inventory,
private helper-package shape, and topology-enforcement tests are flagged
**[defaulted]** for reviewer veto, as Handoff item 7 requires.
A Claude Opus 5 high-effort falsification pass reviewed the draft; every
verified finding is folded into the coverage map and edge inventory below.

## Problem

The changed-tree dev gate takes about 128 seconds and is one package wide:
the artifact contract package takes about 109 seconds solo, so it floors the
contract phase after every edit. Its work is independent by subject but runs
in one Go test process, where top-level tests are serial. The gate's existing
package scheduler has no independent artifact packages to overlap.

The serial share is real coverage, not removable waste. A fresh measurement on
2026-07-29 put the top-level posture tests at about 48 seconds, offline journeys
and verification at about 25 seconds, distributable-artifact work at about 10
seconds, and prepared-artifact inspection and promotion at about 13 seconds.
Those groups account for about 96 seconds of the prior 109-second solo package
wall; package startup, compilation, cleanup, and host-load variation account
for the remaining observed overhead. Four new processes repeat part of that
overhead, so the map's 40–60-second estimate remains a range to measure, not a
sum asserted as fact. The process topology, not one removable check, is the
wall.

## Solution

Split the artifact contracts into four subject packages — `posture`, `offline`,
`prepared`, and `distributable` — so the existing Go package scheduler can
overlap them. Keep the six prepared-artifact sharers and their singleton in the
same `prepared` process, preserving the prior one-build hoist. Put shared
artifact-fixture behavior behind one private support package, and have every
subject package enter through its one package-main runner so shared-cache
posture remains single-sourced.

All 33 top-level tests remain present exactly once. Behavior-owned canary
fixtures move with the package that owns their named test. No oracle semantic,
gate tier, test assertion, `-count=1` policy, or production artifact behavior
changes. Expected critical path: the artifact suite falls to about 40–60
seconds and the changed-tree gate to about 60–75 seconds, jointly bounded by
the posture package and the existing test phase.

## User stories

1. As a developer changing the kit, I want artifact contracts exposed as four
   independently schedulable subject packages with no tests left in the old
   serial root, so the changed-tree gate overlaps legitimate work instead of
   serializing it. Line: `gpt-5.6-terra` / medium. The package seam is decided
   and gate-observable, while moving the oracle's test topology merits the
   profile's gate-logic routing.
2. As the posture package, I want the hermetic-default, GOPROXY-off,
   stale-record, and checkout-topology contracts kept together, so their
   measured 48-second workload becomes one explicit scheduler unit.
   Line: `gpt-5.6-luna` / low. The exact tests and package seam are fixed, and
   the package's own focused command grades the mechanical migration.
3. As the offline package, I want release-plan projection, native-proof,
   offline-journey, interruption, comparator, and verifier contracts kept
   together, so their measured 25-second workload can overlap the posture
   process. Line: `gpt-5.6-luna` / low. This is a fully enumerated move whose
   existing assertions remain the oracle.
4. As the prepared package, I want all six ruled prepared-artifact sharers,
   their package-scoped singleton, its four belt tests, and the promotion and
   source-state contracts in one process, so the split preserves exactly one
   prepared build rather than multiplying it across packages.
   Line: `gpt-5.6-terra` / medium. The move crosses the singleton's
   package-process lifetime and therefore needs the profile's
   correctness-biased routing.
5. As the distributable package, I want wrapper, platform-artifact, installed
   lifecycle, and special-evidence-input contracts in their own scheduler unit,
   so the remaining 10-second artifact journey overlaps the other subjects.
   Line: `gpt-5.6-luna` / low. Two named top-level tests move together with
   their private assertions unchanged.
6. As every artifact subject package, I want shared source staging,
   package-main cache posture, and common artifact operations owned by one
   private fixture module, so splitting the package does not create four
   derivations of the same test policy. Line: `gpt-5.6-terra` / medium. Placing
   the new helper seam is a one-source-per-fact decision that the gate only
   partially observes.
7. As the canary sweep, I want each of the five behavior-owned artifact
   fixtures moved under its new owning package while keeping its `TEST` and
   `EXPECT` facts unchanged, so per-test bites compile the right binary and
   continue guarding the same mutations. Line: `gpt-5.6-terra` / medium. This
   edits oracle binding data, which follows the cached gate-logic route.
8. As the reviewer, I want an exact inventory assertion proving that all 33
   pre-split top-level tests exist once across the four packages and every
   package enters the one shared package-main runner, so an omitted test,
   forgotten runner call, or inline cache-policy copy cannot buy a misleading
   speedup. Line: `gpt-5.6-terra` / medium. The independently authored
   expectation is needed for named omission and duplication mutations to make
   the gate red.
9. As FT91's stop rule, I want the focused suite, fresh changed-tree gate, and
   dormant outer-width-cap posture re-measured and recorded with the process
   boundaries that prevent further reuse, so the next performance decision
   prices the whole critical path. Line: `gpt-5.6-terra` / medium. Recording
   timings is mechanical, but ruling whether the dormant width cap revives is
   gate-concurrency judgment under the profile's cached routing.

## Implementation decisions

- **[defaulted] Package topology.** The old artifact root becomes a container,
  not a test package. Its four test packages are `posture`, `offline`,
  `prepared`, and `distributable`; a private `internal/fixture` package beneath
  that root owns shared test support. Package names are subject names rather
  than implementation-layer names, matching the map's scheduling unit.
- **[defaulted] Exact inventory.** The 33-test inventory is the omission
  oracle. These are top-level tests; subtests remain with their parent.
  - `posture` (4):
    `TestArtifactBuilderHonorsHermeticDefault`,
    `TestBuildCachePostureUnderGoproxyOff`,
    `TestSharedCacheBuildRemovesStaleReproducibilityRecord`,
    `TestGoBuildIgnoresCheckoutTopology`.
  - `prepared` (13). The six singleton sharers are
    `TestArtifactPromotionIsAtomicAndExclusive`,
    `TestArtifactSourceStagesCommittedHostPlan`,
    `TestSharedCacheBuildPromotesNoRecord`,
    `TestSharedCacheBuildRestoresRecordOnInterruptedPromotion`,
    `TestOfflineArchiveProjection`, and
    `TestPackedArtifactRunsSetupOfflineFromASpacedPrefix`. The four singleton
    belt tests are `TestSharedArtifactSetFailsClosedAfterEarlierStagingFailure`,
    `TestSharedArtifactSetBuildIsLazy`,
    `TestSharedArtifactSetAttributesMutation`, and
    `TestSharedArtifactSetIsReadOnly`. The remaining source-state and staging
    contracts are `TestArtifactBuilderRejectsDirtyAndUntrackedSourceState`,
    `TestArtifactBuilderRefusesMissingBinaryPinManifest`, and
    `TestArtifactSourceSkipsWhenHostTargetIsAbsent`.
  - `offline` (14):
    `TestReleasePlanProjectsDerivedArchiveInventory`,
    `TestReleasePlanProjectsRelocatedPackageEvidence`,
    `TestNativeProofAggregatorRejectsOneTargetOmission`,
    `TestNativeProofAggregatorRejectsDigestMismatch`,
    `TestAuthoritativeNativeProofBehaviorCanary`,
    `TestOfflineRegistryDerivesAcceptedTargetsFromReleasePlan`,
    `TestOfflineArchiveBuildRefusesOutputItCannotAccountFor`,
    `TestOfflineInstructionsVerifyOnlyTargetArchive`,
    `TestOfflineSmokeRunsThePublicJourneyAndAttributesMutations`,
    `TestOfflineNetworkSentinelDeniesUndeclaredEgress`,
    `TestOfflineSmokeRecoversFromEveryStageInterruption`,
    `TestReproducibilityComparatorRejectsArtifactAndEvidenceMutations`,
    `TestOfflineSmokeRequiresApprovedReleaseEvidence`,
    `TestReleaseArtifactVerifierRequiresFullyApprovedEvidence`.
  - `distributable` (2): `TestDistributableArtifactContracts`,
    `TestArtifactBuilderRejectsSpecialReleaseEvidenceInput`.
- **[defaulted] Private helper and runner contract.** The private fixture
  module owns shared release-plan narrowing, artifact-source staging, common
  command setup, and a package-main runner. The runner sets the exact dev
  shared-build-cache token, runs the package, invokes an optional package
  cleanup, and preserves the test exit code. Each subject package has a thin
  `TestMain` call into that runner. Posture tests continue stripping the token
  from subprocess environments when hermetic behavior is their subject.
- `prepared` alone owns the singleton state, package temp directory, cleanup,
  read-only belt, fingerprint belt, build counter, and six consumers. No other
  package imports or recreates that state. General fixture operations may move
  to the private helper; singleton lifetime and cleanup do not. Extend the
  existing laziness test with a fresh-process selection of the six named
  sharers under the build-log knob and assert exactly one log line. The
  accessor's current counter proves one build per state; this subprocess proves
  the six consumers did not acquire six states.
- **[defaulted] Topology enforcement shape.** A test owned by the private
  fixture module parses the four subject packages and compares their top-level
  test functions with the independently authored 33-test expectation. It also
  requires each thin `TestMain` to call the shared runner and rejects an inline
  reference to the shared-cache token. During the build, demonstrate three
  mutations: omit one test, replace one runner call with inline token setup,
  and give each prepared sharer a fresh singleton state. Each must red its named
  assertion before being reverted. The project code standard permits this
  expectation-versus-implementation independence only because these named
  omissions and duplications are recorded and demonstrated; production policy
  remains single-sourced.
- **[defaulted] Canary package mapping.** Canary fixture package paths change
  from `surface/artifact` to
  `surface/artifact/offline` for the three native-proof/offline-verifier owners
  and to `surface/artifact/distributable` for the two wrapper owners. `TEST`,
  `EXPECT`, `BASE`, and mutation content stay unchanged. The canary package
  resolver and compiled-bite runner already accept nested package identifiers;
  no new canary seam is introduced.
- Process boundaries are load-bearing evidence. Go can schedule the four test
  packages concurrently, but each is a distinct binary with its own
  `TestMain`, globals, temp lifecycle, and environment. A package singleton
  cannot be reused by another package, and a behavior-owned canary compiles and
  runs only its owning package binary. Further reuse would require persistent
  cross-process coordination, not another in-memory helper.
- The contract phase keeps its recursive package invocation and `-count=1`.
  There is no package allowlist to maintain and no file-to-test scope map.
  Performance comes solely from exposing independent packages to the existing
  scheduler.

## Testing decisions

- A good test drives the existing package seam. Each subject package's focused
  `go test` command must pass independently; the recursive artifact command
  proves they compose; the dev gate proves contract/canary integration.
- The inventory expectation lists top-level test names by owning package and
  rejects missing, duplicate, extra, or root-owned migrated tests. The same
  private-fixture contract test rejects a missing runner call or an inline
  cache-policy copy; it is outside the four subject inventories, so its own
  top-level tests do not change the 33 migrated-test count. Prior art is the
  conformance suite's exact generated-surface checks.
- The two existing posture contracts are the hermetic-stripping oracle:
  `TestArtifactBuilderHonorsHermeticDefault` requires a reproducibility record
  that a leaked shared-cache token suppresses, and the hermetic-default row of
  `TestBuildCachePostureUnderGoproxyOff` requires the build to fail offline,
  while a leaked token makes it succeed.
- Gate: `.bench/gate.sh` through `bench gate --fresh` after focused tests.
  Focused command:
  `go test -count=1 ./internal/contract/surface/artifact/...`.
- Wall-clock is not a deterministic gate assertion. The green build records
  focused and fresh-gate timings in the compiled decision map and compares
  them with the 2026-07-29 109-second package / 128-second gate baseline.

### Seam diagram

    trigger: contract phase invokes go test -count=1 ./internal/contract/...
        │
        ▼
    artifact package tree ──▶ [ Go package scheduler ]
                              ├─▶ posture process ───────┐
                              ├─▶ offline process ───────┤
                              ├─▶ prepared process ──────┼─▶ package verdicts
                              └─▶ distributable process ┘
                                      │
                                      └─ prepared alone owns one lazy set
    canary fixture ──▶ owning nested package + TEST ──▶ one compiled bite
                     ◀ tests attach here: exact inventory/posture assertions,
                       four focused package commands, recursive suite, and
                       live contract + canary gate phases

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | the old root owns no tests and four nested subject packages are independently runnable | Go package tree | observed red 2026-07-29: each of the four focused `go test` commands failed `directory not found`; the root currently reports 11 test files | an empty cosmetic directory or a split that leaves the serial root intact cannot satisfy both sides |
| 2 | `posture` owns exactly its four enumerated top-level tests and runs them under its package process | `go test -count=1 ./internal/contract/surface/artifact/posture` | observed red 2026-07-29: package directory absent | the package command cannot turn green until the scheduler unit and its named tests exist |
| 3 | `offline` owns exactly its 14 enumerated top-level tests and runs them under its package process | `go test -count=1 ./internal/contract/surface/artifact/offline` | observed red 2026-07-29: package directory absent | the package command and inventory assertion red on an absent or incomplete offline unit |
| 2 | hermetic-subject commands strip the package-main shared-cache token | the two existing posture contract tests | already covered — a leaked token removes the reproducibility record expected by `TestArtifactBuilderHonorsHermeticDefault` and makes the GOPROXY-off hermetic-default case succeed when it must fail | both observable directions fail if the move inherits the package token into a command whose subject is hermetic behavior |
| 4 | `prepared` owns exactly its 13 enumerated tests and the four existing failure/laziness/mutation/read-only belts stay green | `go test -count=1 ./internal/contract/surface/artifact/prepared` | observed red 2026-07-29: package directory absent; the four belt behaviors are already covered by their named tests | moving the belts or their singleton lifetime incorrectly fails the focused package |
| 4 | the six named sharers trigger one prepared build in a fresh process | `TestSharedArtifactSetBuildIsLazy` extended to select all six under the existing build-log knob | observed red required during the build: give every sharer a fresh singleton state, observe six log lines, then revert to one | the current accessor counter proves one build per state; the process log proves all six consumers share that state |
| 5 | `distributable` owns exactly its two enumerated tests and runs them under its package process | `go test -count=1 ./internal/contract/surface/artifact/distributable` | observed red 2026-07-29: package directory absent | an empty package or a partial move cannot provide both named tests |
| 6, 8 | every subject `TestMain` calls the one shared runner and carries no inline shared-cache token | private-fixture topology contract | observed red required during the build: replace one runner call with inline token setup and record the structural refusal before reverting | observing environment state alone accepts four policy copies; call-shape plus token rejection catches the actual one-source failure |
| 7 | all five behavior-owned fixtures compile the new owner package and bite their unchanged owning test | live canary plus owner-binding validation | already covered — unknown package/test bindings refuse before graded runs, and each fixture must bite non-vacuously | moving a test without its fixture path makes the compiled owner absent or the TEST unknown; changing the TEST/EXPECT instead would fail the unchanged-fact review |
| 8 | the union of package inventories is exactly the 33 pre-split top-level tests, with no duplicate, extra, missing, or root-owned test | exact inventory assertion at the package-tree seam | observed red 2026-07-29: all four expected packages are absent; during the build, deleting one expected test must red before revert | the cheapest wrong speedup is omission; exact set equality, not a count alone, catches it |
| 1, 9 | the four package pass intervals overlap and no shared fixture lock serializes them | `go test -json -count=1 ./internal/contract/surface/artifact/...` timing evidence | not TDD-able — scheduler overlap is host-sensitive; the recorded JSON intervals must show overlap, while a serialized trace blocks the performance claim | four cosmetic package names behind one lock pass existence tests but cannot produce overlapping pass intervals |
| 9 | focused suite time, fresh changed-tree gate time, removed serial share including repeated process overhead, remaining critical path, and width-cap verdict are recorded | compiled decision map | not TDD-able — host wall-clock is evidence, not a stable timeout; review compares recorded commands and values with the baseline | a deterministic time assertion would grade host load rather than correctness, while an unrecorded result cannot support the next decision |

### Edge inventory

- **Error path:** a subject package that cannot enter the shared package-main
  runner fails before tests; a prepared-set stage/build/cleanup error retains
  the existing attributed failures. Covered by stories 4, 6, and 8.
- **Empty or absent input:** no tests may remain at the old root and no subject
  package may be empty. Exact inventory equality covers both.
- **Boundary values:** one omitted test, one duplicate test, one extra test, and
  one deleted package-main call are the smallest wrong states. The inventory
  and posture mutations cover them per class.
- **Malformed input:** a behavior-owned fixture with an old package path or an
  unknown `TEST` owner is refused by the existing canary binding contracts.
- **Interrupted or partial state:** each package is a separate process; a hard
  kill can leave its OS temp directory exactly as today. **Won't handle:** host
  temp reaping owns hard-kill residue; this split adds no persistent cache.
- **Re-run idempotency:** each `go test` invocation creates fresh package
  processes and the prepared package creates a fresh singleton directory.
  Existing laziness and cleanup tests stay in the inventory.
- **Cross-process contention:** the only intentionally shared write-heavy
  resources are Go's build and module caches, whose concurrent-process posture
  is already accepted. Each fixture gives `HOME`, `TMPDIR`, npm cache, source
  clone, output directory, and output lock a process-private path. The existing
  atomic-and-exclusive test keeps ownership of same-output contention. The
  recursive focused command is the integration probe; its JSON trace must show
  overlap rather than a hidden shared lock.
- **Hostile environment:** spaced paths, dirty/untracked state, special files,
  missing tools, offline operation, and interrupted promotion remain owned by
  their existing tests after the move. `GOMAXPROCS=2` nested canary runs compile
  and invoke only one owning package, so no assertion may require overlap.
- **Compatibility probe:** the recursive `go test ./...` behavior is the
  official Go package interface already used by the gate. No external wire or
  file format changes.
- **Won't handle:** package-level result caching — the contract phase retains
  `-count=1` because subprocess effects are not in Go's package cache key.
- **Won't handle:** `t.Parallel` inside a subject package — rejected because it
  converts the prepared singleton's ordering analysis into shared-state races.
- **Won't handle:** control bytes, line-sink newlines, missing trailing
  newlines, dangling symlinks, flag-value positionals, TTY state, deep cwd,
  invocation through shipped wrappers, and destructive worktree state — the
  package topology consumes none of these inputs; the moved contracts that do
  consume them retain their existing assertions unchanged.

## Out of scope

- **Persistent cross-process artifact reuse** — a separate cache/coordination
  capability with its own invalidation and cleanup contract; about 12 edits,
  5 gate runs.
- **Reviving the dormant outer conformance/contract width cap** — a separate
  scheduler-policy capability, built only if the post-split measurement shows
  oversubscription; about 6 edits, 3 gate runs.
- Removing `-count=1`, canary input-key skipping, fixture batching, and
  diff-scoped gating are rejected alternatives, not deferred work.
