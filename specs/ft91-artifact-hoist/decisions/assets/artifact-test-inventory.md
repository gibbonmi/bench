# Artifact-suite test inventory (gate-critical-path #2)

Classification of every test in `internal/contract/surface/artifact`, read
2026-07-29 at `d60f95c`. Buckets follow the map's question: **M** mutates the
artifact set or its environment, **G** asserts on the act of generation itself
(atomicity, promotion, refusals, posture), **I** inspects prepared output only.
"Gen" counts real `build-artifacts.sh` generation passes the test pays today.

## Facts the ruling rests on

- The prepared seam copies, never consumes: `scripts/build-artifacts.sh:110`
  runs `cp -a "$BENCH_TEST_PREPARED_ARTIFACTS/." "$artifacts/"`, and
  `TestArtifactPromotionIsAtomicAndExclusive` reuses one prepared dir across
  winner, loser, interrupt, and rerun with digest-equality asserted
  (`artifact_fixture_test.go:295-326`) — a prepared set is read-only input.
- No test in the package calls `t.Parallel` (rg, 2026-07-29) — execution is
  sequential, so a shared set faces ordering hazards only, not races.
- `TestMain` puts the whole package on the shared-cache build posture
  (`build_cache_posture_test.go:19-24`); posture tests strip the token back
  out via `ambientBuildEnv` per row.
- The hermetic default is a **double** generation: a non-shared-cache build
  spawns a second independent build plus comparison
  (`scripts/build-artifacts.sh:154-167`).
- Every builder-invoking test stages its own committed clone via
  `committedHostileArtifactSource` → `contract.NarrowReleasePlan`
  (`internal/contract/releaseplan.go:42-72`): commit-working-tree + clone +
  narrowed `release-plan.json` per call, host-only unless the test opts into
  more.
- The missing-pin refusal fires after the build matrix
  (`scripts/build-artifacts.sh:139-145` vs `:125`), so that refusal test pays
  one full generation first.

## Bucketed inventory

| Test | Bucket | Gen | Basis |
|---|---|---|---|
| TestDistributableArtifactContracts (`artifact_test.go:61`) | G+I | 1 (host + one cross-compile) | Builds a uniquely-shaped host+non-host set (`:65`), asserts count (`:101`), then pure read-only inspection and installed-lifecycle consumption (`:106-113`) |
| TestArtifactBuilderRejectsSpecialReleaseEvidenceInput (`artifact_test.go:116`) | M+G | 0 (early refusal) | Mutates its private staged source (LICENSE→fifo, `artifact_fixture_test.go:189-215`); refusal precedes generation |
| TestArtifactPromotionIsAtomicAndExclusive (`artifact_concurrency_test.go:14`) | G | 1 (prepared build) | Promotion atomicity via the seam; prepared set consumed read-only |
| TestArtifactBuilderRejectsDirtyAndUntrackedSourceState (`artifact_source_state_test.go:13`) | M+G | 0 | Mutates private staged sources; refusals + stub prepared dir (`:46-48`) skip generation |
| TestArtifactBuilderRefusesMissingBinaryPinManifest (`artifact_source_state_test.go:64`) | G | 1 | Pays a full generation before the pin check refuses (`build-artifacts.sh:139-145`) |
| TestArtifactSourceStagesCommittedHostPlan (`artifact_staging_test.go:13`) | G+I | 1 | Asserts staged plan is exactly host (`:22`) and generated names match the plan (`:27`) — satisfiable against a shared set built by the same recipe |
| TestArtifactSourceSkipsWhenHostTargetIsAbsent (`artifact_staging_test.go:30`) | I | 0 | Subprocess skip-path probe; no artifacts |
| TestArtifactBuilderHonorsHermeticDefault (`build_cache_posture_test.go:44`) | G | 2×2 (hermetic double-builds) | Posture is the subject; strips the shared-cache token per row |
| TestBuildCachePostureUnderGoproxyOff (`build_cache_posture_test.go:70`) | G | 3 (1 refuses early) | Cache posture under GOPROXY=off is the subject |
| TestSharedCacheBuildPromotesNoRecord (`build_cache_posture_test.go:99`) | G+I | 1 | Asserts record absence beside output after a shared-cache build — satisfiable against a shared set's parent |
| TestSharedCacheBuildRemovesStaleReproducibilityRecord (`build_cache_posture_test.go:110`) | M+G | 2 | Seeds a stale record its own build must remove; private by nature |
| TestSharedCacheBuildRestoresRecordOnInterruptedPromotion (`build_cache_posture_test.go:135`) | M+G | 1 (prepared build) | Aborted-promotion record restore via the seam; prepared set read-only, output/record private |
| TestGoBuildIgnoresCheckoutTopology (`reproducibility_test.go:18`) | G | 0 (+2 `go-build.sh` cross builds) | Binary reproducibility of `go-build.sh`, not the artifact suite |
| TestReproducibilityComparatorRejectsArtifactAndEvidenceMutations (`reproducibility_test.go:57`) | I | 0 | Fabricated files; comparator behavior only |
| TestOfflineSmokeRequiresApprovedReleaseEvidence (`reproducibility_test.go:170`) | I | 0 | Empty-dir refusal |
| TestReleaseArtifactVerifierRequiresFullyApprovedEvidence (`reproducibility_test.go:181`) | I | 0 | Fabricated evidence mutations |
| TestOfflineSmokeRunsThePublicJourneyAndAttributesMutations (`offline_journey_test.go:27`) | I | 0 | Fixture tarballs fabricated in Go (`:259-307`); ~13 smoke runs, no builds |
| TestOfflineNetworkSentinelDeniesUndeclaredEgress (`offline_journey_test.go:197`) | I | 0 | Node one-liner |
| TestOfflineSmokeRecoversFromEveryStageInterruption (`offline_journey_test.go:215`) | I | 0 | 4× interrupt + green rerun of smoke on fabricated fixtures |
| TestReleasePlanProjectsDerivedArchiveInventory (`artifact_offline_test.go:19`) | I | 0 | release-plan.mjs projection |
| TestReleasePlanProjectsRelocatedPackageEvidence (`artifact_offline_test.go:32`) | I | 0 | Fixture copy + node assemble |
| TestNativeProofAggregatorRejectsOneTargetOmission / …DigestMismatch (`artifact_offline_test.go:82,111`) | I | 0 | Fabricated proofs |
| TestAuthoritativeNativeProofBehaviorCanary (`artifact_offline_test.go:128`) | I | 0 (+1 nested `go test`) | Spawns a tagged releaseevidence test |
| TestOfflineRegistryDerivesAcceptedTargetsFromReleasePlan (`artifact_offline_test.go:159`) | I | 0 | Node registry fixture |
| TestOfflineArchiveProjection (`artifact_offline_test.go:223`) | G+I | 1 | Real host-only build consumed read-only by `build-offline-archives.sh` |
| TestOfflineArchiveBuildRefusesOutputItCannotAccountFor (`artifact_offline_test.go:248`) | I | 0 | Fabricated output dirs; refusal |
| TestOfflineInstructionsVerifyOnlyTargetArchive (`artifact_offline_test.go:296`) | I | 0 | Fabricated assemble stage |
| TestPackedArtifactRunsSetupOfflineFromASpacedPrefix (`setup_offline_test.go:20`) | G+I | 1 | Real host-only build consumed read-only by npm install + `bench setup` |

## What per-test scoping guarantees today

- **M class:** a fresh committed clone and fresh output dir per test, so source
  corruption (fifo, dirty submodule, broken index) and seeded/interrupted
  records cannot leak anywhere.
- **G class:** every builder run answers for exactly the staging its test just
  performed; nothing compares or reuses bytes across tests.
- **I class:** the isolation buys nothing — the assertions read a set that any
  equivalent generation would satisfy; the per-test build is pure cost.
- Byte-identity across tests is neither guaranteed nor asserted today.

## The collapsible set

Host-only generations from the identical recipe (host-narrowed committed
clone → one shared-cache build): **StagesCommittedHostPlan, PromotesNoRecord,
OfflineArchiveProjection, SetupOffline, Concurrency (prepared),
RestoresRecord (prepared), MissingBinaryPin's pre-refusal build** — up to 7 of
the ~17 generation passes (~3.7 s each warm) collapse to 1, and the saving
multiplies through every canary fixture run that pays the suite.
Non-collapsible by subject: the 4 hermetic double-build passes, the 3
GOPROXY-off passes, the 2 stale-record passes, and Distributable's unique
host+non-host shape (collapsible only if the shared set adopts that shape,
which changes every host-only consumer's count assertions).
