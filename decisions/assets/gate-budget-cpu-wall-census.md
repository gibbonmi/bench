# Gate CPU-vs-wall census

Evidence asset for `decisions/gate-budget.md` #8, #13, and #15. Measured
2026-08-04 in the `25385f8` checkout, 12-core box, ambient Go build cache,
go1.25.0. The checkout retained the pre-existing uncommitted
`gate-concurrency` edit and the untracked FT171 shaping documents; no
production code was changed for these probes.

## Method

Compilation was warmed first (`go test -count=1 -run '^$' ./...`) so the census
measures test execution rather than builds. Each of the 70 packages was then run
alone, serially, under `time -p`: `go test -count=1 <pkg>`. Because packages ran
one at a time, these wall figures are **focused spans** — the denominator #6's
span-inflation measure needs — and carry no cross-package contention.

CPU is `user + sys`, which includes reaped child processes. Ratio is CPU/wall:
above 1 means intra-package parallelism (`t.Parallel`, `-p`), below 1 means the
process was waiting rather than computing.

This is a package-universe census, not a dev-gate census. `internal/preflight`,
`internal/releaseevidence`, and `internal/publication` are release-only packages:
the dev tier excludes them and `bench prep-release` runs them at ship tier.

## Totals

| | |
|---|---|
| packages | 70 |
| serial wall | 777.0 s |
| CPU | 1198.4 s |
| aggregate ratio | 1.54 |
| packages under 0.2 s | 16 |

Cost is concentrated: the thirteen packages below carry essentially all of it,
and the remaining 57 are collectively under ten seconds.

## Results

| package | wall (s) | CPU (s) | ratio | shape |
|---|---:|---:|---:|---|
| `internal/preflight` | 168.0 | 369.2 | 2.20 | saturating |
| `internal/gate` | 147.4 | 112.9 | 0.77 | mixed |
| `internal/conformance` | 81.0 | 134.0 | 1.65 | saturating |
| `internal/contract/runtime` | 64.4 | 193.4 | 3.00 | saturating |
| `internal/contract/surface/artifact/posture` | 60.8 | 179.6 | 2.96 | saturating |
| `internal/specbuild` | 43.6 | 31.6 | 0.73 | mixed |
| `internal/contract/surface/publication` | 39.4 | 4.6 | 0.12 | idle |
| `internal/publication` | 30.4 | 0.6 | 0.02 | idle |
| `internal/worktree` | 29.9 | 24.6 | 0.82 | mixed |
| `internal/contract/surface/artifact/prepared` | 24.3 | 23.1 | 0.95 | serial |
| `internal/contract/surface/artifact/offline` | 23.1 | 22.2 | 0.96 | serial |
| `internal/contract/surface` | 18.1 | 50.6 | 2.79 | saturating |
| `internal/contract/surface/artifact/distributable` | 10.7 | 10.7 | 1.00 | serial |
| `internal/contract/axi` | 1.6 | 10.2 | 6.31 | saturating |
| `internal/contract/surface/preprelease` | 1.4 | 7.8 | 5.63 | saturating |

## Three shapes, not one workload

**Saturating (ratio ≥ 2).** `preflight`, `runtime`, `posture`, `contract/surface`,
`axi`, `preprelease`. These are the leafblower: `runtime` alone pulls three
cores' worth. A core-denominated token pool bounds exactly this class.

**Idle (ratio ≤ 0.2).** `contract/surface/publication` at 12% and
`internal/publication` at 2% — the latter spends 30.4 seconds to burn 0.6
seconds of CPU. Located example: `TestReleaseSubmitExitCodeTriple` is 30.2 s of
its 39.4 s package, and its middle subtest points the release client at
`http://127.0.0.1:1` and waits for the connection to give up. A token pool
cannot reclaim any of this; the processes are asleep, not contending.

**Serial (ratio ≈ 1).** `prepared`, `offline`, `distributable`. One core's worth
of genuine work that cannot be parallelized away. A tight reserve hurts this
class without any noise benefit, because it cannot trade width for time.

## What this changes for the map

`internal/preflight` and `internal/gate` carry 40% of this census's serial wall,
but only `internal/gate` is a dev-gate cost. The earlier critical-path figures
missed `internal/gate`; they did not miss 168 seconds of dev-gate preflight.

`r` can only be priced against the saturating class. Measuring span inflation on
an idle-class phase reads "no symptom" from a process that is merely waiting,
which would let #6's acceptance rule certify a bad `r`. A ship-only saturating
package cannot price a dev-gate reserve either.

## Follow-up probes

### Artifact posture cache A/B

`posture` does not depend on a gate-level environment setting: its `TestMain`
calls `fixture.Run`, which sets `BENCH_SHARED_BUILD_CACHE=1` before any test.
The cold-cache premise was stale. An overlay replacing only that `TestMain`
made the counterfactual measurable without changing the tree; both package
variants were precompiled with `-run '^$'` first.

| subject | cache posture | verdict | wall (s) | CPU (s) |
|---|---|---|---:|---:|
| full `posture` package | current shared | green | 61.15 | 182.28 |
| full `posture` package | hermetic overlay | red after completed builds | 89.99 | 332.47 |
| `TestSharedCacheBuildRemovesStaleReproducibilityRecord` | current shared | green | 8.18 | 8.71 |
| same two-build test | hermetic overlay | red after completed builds | 43.31 | 164.65 |

The focused pair is the clean cost comparison: the same two artifact-builder
invocations completed in both runs. The hermetic run then failed its
shared-posture assertion because the reproducibility record correctly remained.
Shared cache cut 94.7% of CPU and 81.1% of wall for those two builds. The full
package's hermetic verdict is not an acceptance result because cache-contract
tests require the package-level opt-in; its timing is supporting evidence only.

Cold recompilation is therefore expensive but already removed from the dev
suite wherever shared posture is intended. The package still exercises the
hermetic default deliberately, so cache posture does not erase `posture` from
the saturating class and does not replace scheduling.

### `internal/preflight` per-test timings

One `go test -json -count=1 ./internal/preflight` run was green at 169.93 seconds
wall and 375.55 seconds CPU (ratio 2.21), reproducing the census. These are all
top-level tests at or above two seconds; the other 21 total 14.46 seconds.

| top-level test | elapsed (s) |
|---|---:|
| `TestReleasePreflightScriptBootstrapsBuiltFullAndFocusedCommands` | 66.06 |
| `TestBuiltCommandRejectsHostilePackageArchivesAndPreservesPriorGeneration` | 12.91 |
| `TestBuiltCommandRejectsProducerRecordMutationsDistinctly` | 11.05 |
| `TestBuiltCommandArchiveBudgetsFailPromptlyAndPreservePriorGeneration` | 7.96 |
| `TestReleaseEvidenceIsDeterministicBoundAndIdempotent` | 7.36 |
| `TestBuiltCommandReleasePolicyFailuresAreRed` | 5.96 |
| `TestBuiltCommandProfileAcceptanceMatrix` | 4.87 |
| `TestBuiltCommandMissingNodeOrNPMPreservesPriorGeneration` | 4.25 |
| `TestFinalTarRejectsHostileEvidenceWithConsistentInventory` | 3.98 |
| `TestBuiltCommandInputDriftPreservesPriorCompleteEvidence` | 3.87 |
| `TestBuiltCommandSourceIdentityDriftPreservesPriorGeneration` | 3.66 |
| `TestReleaseProfilesStayPendingInVerifyAndRedInPublish` | 3.48 |
| `TestGovernanceSchemaRejectsUnknownFieldsAndVersions` | 3.43 |
| `TestConditionalRecordReasonFlowsThroughAuthoritativeEvidence` | 3.25 |
| `TestIndexEncodingFailurePreservesPriorGeneration` | 3.07 |
| `TestBuiltCommandControlByteArchivePathPreservesPriorGeneration` | 2.85 |
| `TestBuiltCommandRejectsConcatenatedGzipArchiveAndPreservesPriorGeneration` | 2.51 |
| `TestBuiltCommandCancellationPreservesPriorCompleteEvidence` | 2.42 |
| `TestBuiltCommandFocusedPublishRunsDiagnosticWithoutAuthorizing` | 2.31 |

The 66.06-second leader is 38.9% of package wall. It copies the source to a
fresh repository, then stubs every full-verify phase except `artifacts`; the
real artifact phase runs without the shared-cache opt-in and the wrapper also
bootstraps a fresh `bench-preflight` binary. Its later focused and symlink calls
reuse that binary. The direct vulnerability-policy test is effectively zero and
the cancellation test is 0.38 seconds, while the leader stubs vulnerability, so
a live vulnerability scan does not explain this package's dominant cost.

This cost is real but ship-only. It must not enter #8's reserve selection or
the dev gate's two-minute target.

### `internal/gate`: repeated gate decisions over one tree

One `go test -json -count=1 ./internal/gate` run was green at 153.40 seconds
wall and 116.33 seconds CPU (ratio 0.76), reproducing the census. The package
held 216 top-level tests; their elapsed sum was 152.68 seconds. Eighteen tests
at or above two seconds carried 96.12 seconds.

| top-level test | elapsed (s) |
|---|---:|
| `TestPublicDocumentClassesProjectTheirExactCheckPartition` | 35.14 |
| `TestFT78Story4ProofLedger` | 7.08 |
| `TestDecisionSiteFailsClosed` | 6.09 |
| `TestBuildRunsOnEveryUnsoundArtifact` | 5.67 |
| `TestExecuteTreeBuildsExactUnpublishedBenchkitSource` | 5.54 |
| `TestDeclaredDocumentInputsInvalidateOwningChecks` | 4.59 |
| `TestFT78Story3ProofLedger` | 4.22 |
| `TestContractRunsWhenConsumerInventoryIsMalformed` | 3.42 |
| `TestComposedGreenAcceptsOnlyCompleteExactTipEvidence` | 3.31 |
| `TestFreshExecutesEveryComponent` | 3.29 |
| `TestFirstRunAndFresh` | 2.73 |
| `TestGateGoRace` | 2.71 |
| `TestGateRunDeadlineTermGraceThenKill` | 2.18 |

Six cancellation/runner cases sit at 2.01 seconds because they pay the real
two-second process-group termination grace. The dominant 35.14-second test is
a 21-row mapping matrix. Each row observes a real gate for a mutation and a
deletion and forces a restoring green gate after each, about 84 gate-engine
executions for the table.

The `bench-debug` fast loop selected one row:

```
go test -count=1 -v ./internal/gate \
  -run '^TestPublicDocumentClassesProjectTheirExactCheckPartition$/^projects$/^benchkit.md$'
```

It was green in 2.90 seconds. PATH wrappers around Git and Go recorded 415 Git
children and no nested Go command beyond the outer test. Counts were 118
`rev-parse`, 94 `read-tree`, 84 `cat-file`, 47 each `add` and `write-tree`, 24
`ls-tree`, and one `init`. The uncommitted fixture's same tree was therefore
materialized 47 times. Its failed `read-tree HEAD` plus `read-tree --empty`
fallback makes the test trace worse than a committed production root, but the
re-derivation itself is production code: identity families, subject forms, and
pre/post inspection each request their own snapshot, and repeated blob reads
are not memoized.

This is a missing gate-evaluation seam. A decision-wide immutable parsed tree
and blob cache is one source for the fact; the current helper-local "snapshot
once" promises do not compose across callers. The mapping matrix can then
grade the resolver exhaustively and retain a small representative end-to-end
bite set instead of invoking the whole engine four times per row.

### `internal/conformance`: broad fixture seams and cold gate entries

One green run took 85.11 seconds wall and 138.49 seconds CPU (ratio 1.63).
Eight tests at or above two seconds carried 71.59 seconds:

| top-level test | elapsed (s) |
|---|---:|
| `TestDocsCurrencyTokenDietAndWorkflowFixturesBite` | 35.14 |
| `TestGateEntryNormalizesIndeterminateFreshnessFailures` | 19.95 |
| `TestLoadValidityMetadataFixturesBite` | 3.97 |
| `TestPackageCoreAndGuardFixturesBite` | 3.09 |
| `TestGateEntryRejectsLegacyBeforeRunningOldTableAndRunsReplacementOnce` | 2.87 |
| `TestGateEntryRefusesUnverifiedBinaryBeforeGatePhases` | 2.27 |
| `TestLineRoutingFixturesBite` | 2.17 |
| `TestSkillsIndexAndCommandAdapterFixturesBite` | 2.13 |

The leader has 83 subtests totaling 35.13 seconds. Every documentation or
workflow mutation calls `RunConformance(..., dev, "")`, the complete dev check
table, even though the canary registry already binds its family to one check
and `RunConformance` already accepts that single-check scope. The freshness
leader has eight variants; each creates a new temporary kit module and reaches
the shell gate entry, whose freshness check uses `go run`. The variants prove
one interface but each pays cold module/process setup.

### Distributed integration-matrix cost

Three other dev packages have no hidden long algorithm; scenario count is the
workload.

| package | wall (s) | CPU (s) | top-level tests | concentration |
|---|---:|---:|---:|---|
| `internal/specbuild` | 43.53 | 31.66 | 145 | only four tests >= 1 s; 9.65 s total |
| `internal/worktree` | 29.65 | 23.61 | 134 | only four tests >= 1 s; 6.48 s total |
| `internal/contract/runtime` | 69.78 | 197.73 | 100 | nine tests >= 2 s; 45.92 s total |

`specbuild` and `worktree` repeatedly create and mutate real temporary Git
repositories across lifecycle/fault matrices. Runtime's leaders are the
58-case FT78 action proof ledger (13.68 s), the four-route freshness-refusal
matrix (12.33 s), and cross-process gate/spec-build proofs. These are real
integration matrices, but they explain why source-line count is unrelated to
their elapsed time and identify where representative end-to-end coverage plus
narrower interface matrices can replace repeated CLI setup.

### Uncontrolled idle wait

`internal/publication` reproduced its 30-second idle row: 30.49 seconds wall,
0.46 seconds CPU. `TestFixtureRegistryStagedOpsNeverCheckToolVersions` alone
was 30.03 seconds. It calls `StageSubmit(context.Background(), ...)` through
`http.DefaultClient` against `http://127.0.0.1:1`; on this host that address is
packet-dropped rather than immediately refused, so the kernel's connect retry
owns the test duration. The dev-tier contract publication matrix uses the same
address and carries the same wait. This is FT87's bounded-network-test defect,
not gate scheduling or Go execution.

### Dev and ship workflow duplication

The dev phase table does not accidentally run an ordinary package suite twice.
Of 70 packages, 12 contract packages go only to the contract phase, conformance
goes to its root/suite split, three release-only packages are absent, and the
remaining 54 go to core tests. The conformance suite skips its root entry test;
dev race repeats only three registered tests under a distinct detector; canary
runs distinct mutated trees. `bench commit` also asks for one verdict and
reuses exact green evidence rather than pre-running a second gate.

Ship is different. `prep-release` requires the exact dev-green verdict and then
runs these surfaces:

| ship surface | repeated work |
|---|---|
| `core-tests-ship` | ship enumeration runs 57 packages: the same 54 dev-core suites plus the three release-only suites its comment says the step exists to add |
| `conformance-ship` | ship is a dev superset, so it reruns all dev conformance checks before the release-evidence probe and cross-compile assertion |
| preflight `gate` | reaches exact-green reuse; no second whole dev gate |
| preflight `vet` | repeats the dev gate's `go vet ./...` exactly |
| preflight `race` | reruns all 70 packages under `-race`; distinct instrumentation, but much broader than the authoritative three-test dev race registry |
| first `artifacts` + preflight `artifacts` | builds the same root into `dist/artifacts` twice |

Artifact work appears five times across the whole ship sequence, but only two
are literal same-root jobs. Conformance ship builds the same snapshot in an
authenticated clone as an independent release probe. The two ship canaries
each build a distinct mutated tree to prove digest and package-evidence
omissions; those cannot be merged, though the current full builder recompiles
and repackages far more than the mutations change.

The direct collapse owners are therefore `CoreTestPackages`/the ship step
(release-only enumeration), pre-release plus release-evidence (one artifact
producer), gate evidence plus release-evidence (inherited exact vet proof), and
the authoritative race registry (explicit ship additions instead of `./...`).
Whether ship must reprove all dev conformance checks is a lifecycle-policy
decision, not an accidental duplicate call.

### Remaining optimization evidence

- `internal/gate` is explained and blocks the target: at 153.40 seconds focused,
  it alone exceeds the two-minute whole-gate target once `-count=1` makes the
  package run real. Collapse its duplicated snapshot work and over-wide matrix
  seam before pricing #8.
- `internal/conformance` has the same over-wide fixture seam; its existing
  single-check interface permits narrowing without merging mutated trees or
  weakening the oracle.
- The idle-class timeout work stays with
  `decisions/bounded-network-resource-cli.md` (FT87).
- After demand reduction lands, re-run this census. #8 must pair those new
  focused spans with in-gate spans for current dev-tier saturating subjects.
  Whole-gate wall remains the global outcome.

## Validity

Measured on one box, one tree, one repetition per package. Ratios are stable
enough to classify shape; the wall figures are not repetition-backed and must
not be quoted as variance-bearing. Re-run before any figure here is used to
pick a constant.
