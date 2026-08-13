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

## Post-reduction census: decision #20

Measured 2026-08-07 local time on exact commit
`eb6845f275d1978cffaa5a4868d3be26dae612d2`, tree
`cd2ece9233e2c5df8649a5250dd01eea1ffd8805`, in isolated Bench worktree
`152e98b71788880dcffde9a0c9ef1283`. The host exposed 12 online CPUs on an
Intel i7-13620H under WSL2, using go1.25.0 and ambient build cache
`/home/devuser/.cache/go-build`. No other gate was active at launch.

The worktree binary was built through `scripts/go-build.sh`. The process
profile's outer `gate-run` appeared under the main-checkout binary path, but
that binary and the worktree binary were byte-identical
(`sha256:3626f711b029bc97201404b55ad870a8c15e28eb9942dae67b9b9f369105ae4c`),
their seals were byte-identical, and both passed their own freshness check.
`gate-phases` executed the worktree binary against the worktree subject.

### Cache posture

The compile-only warm-up, `go test -count=1 -run '^$' ./...`, was green in
5.43 s wall and 37.87 s CPU. Each measured package then ran alone with
`go test -count=1`, so test-result caching was disabled while compilation
remained warm.

Clearing the ambient Go cache is the wrong preparation for #20. It would mix
cold compilation into the test workload the outer scheduler must price and
would perturb every worktree sharing the host cache. A separate isolated-cache
probe could price cold-start sensitivity, but its figures must not select the
reserve. The sampled fresh gate grew the warm cache from 4,312,994,789 to
4,445,996,316 bytes (133.0 MB) without a pre-run clear.

### Package universe

All 71 packages passed. Serial wall was 767.25 s and CPU was 1110.99 s, an
aggregate CPU/wall ratio of 1.45. The earlier census was 777.0 s wall and
1198.4 s CPU across 70 packages: the shipped reductions changed important
rows, but subsequent test growth left total serial wall only 1.3% lower.

| package | wall (s) | CPU (s) | ratio | tier/shape |
|---|---:|---:|---:|---|
| `internal/preflight` | 175.53 | 378.36 | 2.16 | ship-only, saturating |
| `internal/gate` | 147.34 | 100.65 | 0.68 | dev, mixed |
| `internal/contract/runtime` | 68.80 | 204.17 | 2.97 | dev, saturating |
| `internal/specbuild` | 58.18 | 42.54 | 0.73 | dev, mixed |
| `internal/contract/surface/artifact/posture` | 55.41 | 159.36 | 2.88 | dev, saturating |
| `internal/contract/surface/publication` | 39.73 | 4.84 | 0.12 | dev, idle |
| `internal/worktree` | 30.82 | 25.15 | 0.82 | dev, mixed |
| `internal/publication` | 30.46 | 0.55 | 0.02 | ship-only, idle |
| `internal/contract/surface/artifact/offline` | 27.78 | 23.64 | 0.85 | dev, serial |
| `internal/conformance` | 26.27 | 35.08 | 1.34 | dev, mixed |
| `internal/contract/surface/artifact/prepared` | 25.95 | 23.39 | 0.90 | dev, serial |
| `internal/contract/surface` | 20.48 | 53.42 | 2.61 | dev, saturating |

`internal/conformance` fell from the earlier 85.11 s focused probe to a
three-run mean of 25.83 s after its fixture rows narrowed. The test-seam cut
also reduced the public-document matrix from 35.14 s to 20.11–21.71 s, but
`internal/gate` as a package did not fall: it now carries 241 top-level tests,
and later work replaced the saved wall.

### Repetitions and focused floor

The decision-bearing dev packages ran three times on the exact subject. Range
is `(max - min) / mean`; it states host variance without implying a
distribution from three samples.

| package | min (s) | median (s) | max (s) | mean (s) | range |
|---|---:|---:|---:|---:|---:|
| `internal/gate` | 147.34 | 152.11 | 159.70 | 153.05 | 8.1% |
| `internal/conformance` | 23.63 | 26.27 | 27.59 | 25.83 | 15.3% |
| `internal/contract/runtime` | 67.67 | 67.81 | 68.80 | 68.09 | 1.7% |
| `internal/contract/surface/artifact/posture` | 52.53 | 55.41 | 55.59 | 54.51 | 5.6% |
| `internal/contract/surface` | 18.59 | 20.48 | 21.10 | 20.06 | 12.5% |
| `internal/specbuild` | 58.18 | 76.73 | 84.53 | 73.15 | 36.0% |

`internal/gate` alone exceeds the destination's 120-second whole-gate target
in every focused repetition. The full matrix's 57 changed generations now
each pay one `write-tree` capture and no per-row full-engine execution; even
deleting its entire 20–22 s span would leave only 7–13 s for all gate setup
at the observed package floor. One reduction at that seam is insufficient.

### Exact fresh gate and subprocess fan-out

One fresh gate ran green under the existing process/memory sampler in 246 s.
This run is not variance-bearing; #20 stopped before repeated full gates
because the focused floor already invalidated width pricing. Selected in-gate
package spans were:

| package | focused median (s) | in-gate span (s) | inflation |
|---|---:|---:|---:|
| `internal/gate` | 152.11 | 230.663 | 51.6% |
| `internal/specbuild` | 76.73 | 177.258 | 131.0% |
| `internal/contract/runtime` | 67.81 | 152.126 | 124.3% |
| `internal/contract/surface/artifact/posture` | 55.41 | 114.969 | 107.5% |
| `internal/contract/surface` | 20.48 | 61.260 | 199.1% |

The process sampler walked exact ancestry every two seconds. At the
97-descendant peak, seven dependency-ready outer roots overlapped:
conformance-suite, test, race, vet, root conformance, canary, and contract.
That sample held 21 `go`, five compiler, six linker, 29 `bench`, five
canary-related, and 12 contract-related processes. Across the run, separate
samples reached 23 `go`, 21 compiler, seven linker, 29 `bench`, ten
canary-related, and 29 contract-related processes.

At least 1,682 distinct child PIDs appeared in the two-second samples: 407
`go`, 344 `bash`, 305 `bench`, 157 compiler, 82 linker, and 76 Git processes,
plus test binaries and Node/npm work. This is a lower bound, not an exact spawn
total: `strace`, `perf`, and `bpftrace` were unavailable, and processes shorter
than the sampling interval can disappear between observations. Peak concurrent
width is exact for the sampled instants; cumulative launch count is explicitly
sample-limited.

The gate's peak descendant PSS was 2.02 GB and cgroup current reached 6.93 GB.
Cache-inclusive allocated memory crossed the sampler's 7 GB symptom threshold,
so the sampler exited 3 after retaining the green gate; the gate itself exited
zero.

### Decision #20 result

Contention is current and material, and the subprocess tree confirms that the
unbounded outer scheduler is one source. It still cannot be priced yet. The
longest focused dev package is mixed rather than saturating and exceeds the
whole-gate target before any competing phase starts, so a token pool can bound
peak demand but cannot make the stated destination achievable. Decision #21
must reduce the focused floor before #8 measures candidate reserves and grant
splits.

## Focused demand trace: decision #21

Measured 2026-08-07 local time on the same exact commit `eb6845f` (tree
`cd2ece9`) as decision #20, same 12-online-CPU host, go1.25.0, ambient warm
build cache, main checkout with documentation-only dirt — the Go workload is
byte-identical to #20's subject. Compilation was warmed first (1.16 s). Timing
runs used clean `go test -count=1 -json`; process counts came from separate
runs with a PATH shim logging every `git` argv. Shimmed runs are counting
evidence only — the shim added roughly 10% wall.

### Both packages are strictly serial chains

`internal/gate` ran 241 top-level tests whose elapsed sum, 139.67 s, equals
its 140.80 s package wall (101.82 s CPU). `internal/specbuild` ran 192 tests
at 57.38 s and 59.81 s sums against 57.80 s and 60.05 s walls in two runs.
Nothing overlaps: the gate package carries three `t.Parallel` calls, specbuild
none. The focused floor is therefore the arithmetic sum of every test's
subprocess waits, and a token pool cannot shorten it — the pool can only stop
concurrent phases from inflating it toward #20's 230.663 s in-gate span.

### `internal/gate`: distributed demand, no dominator

Elapsed bands: 18 tests at or above 2 s carry 78.23 s; 22 tests between 1 and
2 s carry 28.58 s; 27 tests between 0.5 and 1 s carry 19.57 s; the remaining
174 tests carry 13.29 s. The eighteen leaders:

| top-level test | elapsed (s) |
|---|---:|
| `TestPublicDocumentClassesProjectTheirExactCheckPartition` | 18.65 |
| `TestFT78Story4ProofLedger` | 6.93 |
| `TestBuildRunsOnEveryUnsoundArtifact` | 5.91 |
| `TestDecisionSiteFailsClosed` | 5.77 |
| `TestExecuteTreeBuildsExactUnpublishedBenchkitSource` | 5.11 |
| `TestDeclaredDocumentInputsInvalidateOwningChecks` | 4.64 |
| `TestFT78Story3ProofLedger` | 3.70 |
| `TestComposedGreenAcceptsOnlyCompleteExactTipEvidence` | 3.34 |
| `TestFreshExecutesEveryComponent` | 3.15 |
| `TestContractRunsWhenConsumerInventoryIsMalformed` | 3.07 |
| `TestFirstRunAndFreshBuildEverything` | 2.94 |
| `TestGateGoRaceRequiresTheTestToRun` | 2.65 |
| `TestGateRunDeadlineTermGraceThenKill` | 2.18 |
| `TestConformanceImplementationSelectsOwningOrAllCanaryFamilies` | 2.08 |
| `TestPhasesCommandNamesStragglersOnTermination` | 2.07 |
| `TestPhasesCommandSignalCancelsRunningPhaseGroups` | 2.02 |
| `TestRunnerCancelDuringNeededPhaseReturns130` | 2.01 |
| `TestRunnerCancelKillsGroup` | 2.01 |

Cost classes inside the package:

- **Fixed deadlines.** Five cancellation tests pay the literal 2 s
  `processGroupCancelGrace` production constant live, 10.29 s together.
  Observing the grace through an injected duration while keeping at least one
  live-cascade proof would save roughly 8 s.
- **Process materialization.** One package run spawns exactly 13,156 Git
  processes (identical in both counted runs): 5,343 `rev-parse`, 2,182
  `read-tree`, 1,817 `cat-file`, 1,324 `add`, 1,253 `write-tree`, 705
  `ls-tree`, 300 `init`, 89 `commit`. The document matrix's 57 generation
  captures explain most of its 18.65 s; synthesizing those generations at the
  decision seam would save roughly 15 s.
- **Toolchain closure derivation.** Every engine resolution re-runs
  `go list -buildvcs=false -json -deps -test ./...` (`moduleTestClosure`,
  memoized only within one resolution). On a fixture-shaped module the warm
  call costs 0.09–0.14 s (0.60 s first). Child `go` spawns could not be
  PATH-counted — `go test` places `$GOROOT/bin` at the head of the test
  process's `PATH` (documented in `cmd/go`'s test help), and no process
  tracer is installed — so the package-wide count is an estimate from
  structure: with engine resolutions in the low hundreds, this contributes an
  estimated 10–20 s.
- **Real-engine integration.** The rest of the ≥0.5 s mass is deliberate
  representative coverage: proof ledgers, unsound-artifact sweeps, and
  full-engine controls retained by decision #18. Fixture construction itself
  is cheap — a construct-only consumer of the kit-shaped fixture (real
  `git init`, real `go build`, seal) completes in 0.21 s.

### `internal/specbuild`: flat process churn

The flattest possible profile: the largest test is 4.28 s / 4.72 s across the
two runs, only six tests reach 1 s (11.93 s / 12.66 s together), and the
median test is under 0.2 s. One package run spawns 35,952–35,957 Git
processes, about 187 per test at roughly 1.6 ms average: 21,531 `rev-parse`,
2,049 `symbolic-ref`, 1,282 `worktree`, 1,223 `add`, 1,162 `show`, 1,158
`update-ref`, 891 `status`, 825 `show-ref`, 717 `config`, 684 `write-tree`,
683 `read-tree`, 678 `for-each-ref`, 669 `commit`. Sixty percent of all
spawns re-derive repository facts through `rev-parse` — the repeated
source-fact pattern decision #15 found in the gate, here spread across
lifecycle operations and test helpers. There is no algorithm to speed up and
no test to fix; the package cost is spawn count times spawn latency.

### Concurrency constraints

Intra-package parallelism is blocked by process-global state, not by test
isolation. Gate tests each own a private root, but production `kitRoot` reads
ambient `BENCH_KIT`, the kit-shaped fixture claims it per test with
`t.Setenv` (51 construction sites, and `t.Setenv` is mutually exclusive with
`t.Parallel`), and two helpers change the process working directory.
`internal/specbuild` has two `t.Setenv` tests and no chdir, so test-only
`t.Parallel` is structurally available there today.

### Decision #21 result

The enumerated serial cuts — grace observation (~8 s), synthesized matrix
generations (~15 s), closure-derivation memoization (~10–20 s, estimate) —
project a 105–125 s `internal/gate` floor from the observed 140.80–159.70 s
range. That still exceeds the 120 s whole-gate target once any setup or
sibling phase is added, so single-test fixes cannot reach the destination.
The remaining structural lever is concurrency inside the package run, which
is a workload-shape and seam decision the map's #22 owns, not a research
finding this asset can select.

## Post-rebuild census: decision #26

Measured 2026-08-13 (12:57–13:05 UTC) on exact commit `a3b599ea`, the same
12-online-CPU host, load average 0.44 at start with nothing else running.
Baseline: the single-build serial gate plus the branch-native rebuild
(`3701c4a0`) — one host binary per top-level run, one phase process at a
time, direct mutation-to-check canaries inside `internal/conformance`, no
stripped-subject reruns, no nested gate drivers. Cache posture per #20's
ruling: ambient shared build cache retained, `-count=1` throughout, and a
compile-only warm-up (`go test -exec=/bin/true ./...`, 2.25 s wall / 13.8 s
CPU) separating build work from test work.

### Serial package census

64 packages, one serial `go test -count=1` per package under `/usr/bin/time`:
105.2 s wall, 46.7 s CPU (28.8 user + 18.0 sys). Decision #20's baseline was
71 packages at 767.25 s wall / 1110.99 s CPU — a 7.3x wall and 23.8x CPU
reduction. The former dominators are gone: `internal/gate` runs its five
in-package top-level tests in 0.01 s, `internal/specbuild` is deleted. The
current top of the table:

| package | wall s | user s | sys s |
|---|---|---|---|
| `internal/publication` | 30.52 | 0.42 | 0.16 |
| `internal/worktree` | 19.49 | 5.43 | 4.03 |
| `internal/conformance` | 15.35 | 6.54 | 3.34 |
| `internal/freshness` | 7.65 | 1.41 | 1.44 |
| `internal/gitguard` | 4.26 | 0.22 | 0.11 |
| `internal/intent` | 3.23 | 0.94 | 0.35 |
| `internal/diff` | 2.84 | 1.05 | 0.81 |

Focused repetitions, three per package: `publication` 30.31–30.71 s,
`worktree` 18.79–19.10 s, `conformance` 14.98–15.16 s, `freshness`
7.03–7.27 s. Every long package is idle- or churn-dominated, not
CPU-dominated: `publication` remains the FT87 unreachable-port wait (30 s
wall, 0.4 s CPU), and `worktree`/`conformance` spend most of their wall in
subprocess latency.

### Fresh gates

Two `bench gate --fresh` runs with 0.5 s descendant sampling: 38.26 s wall
(green, 33.65 s user + 17.65 s sys) and 38.00 s wall (red, 32.01 s user +
17.09 s sys). Phase spans from the gate's own progress log (green run):
gofmt 0.08 s, vet 0.57 s, test 31.90 s, race 2.41 s, system 1.43 s,
shellcheck 0.26 s — 36.63 s of serial phase time plus ~1.6 s of setup
(selection, single subject build, freshness). Peak concurrent descendants:
25 in both runs, against 97 pre-rebuild. Average utilization is ~1.3 of 12
cores; the box is nowhere near saturation.

The test phase compresses the 105.2 s serial census into 31.9 s through `go
test ./...`'s own package parallelism (`-p` defaulting to `GOMAXPROCS`) and
is floored by `publication`'s 30 s idle wait — the phase's critical path is
FT87's timeout, not CPU contention.

The red run's failure was `TestListCommandCheckedInCompletedAssignmentTerminalPair`
in `internal/worktree`; the serial census also had one unattributed
`internal/worktree` red (output not captured), and the three focused
repetitions passed. Two reds in six package runs this session matches
FT203's ~1-in-6–7 profile, but this failing test is not FT203's named
`TestListCommandPublicRowsAndDisclosure` — the flake family is wider than
the roadmap row records.

### Fan-out inventory and the #25 residual

Re-walk of every production `go func` site: the phase runner launches at
most one phase and blocks on its completion (serial by construction);
`guards` and `sessioninspect` wrap serial children in timeout selects;
`models` fans out network calls only; every other site is a single
`cmd.Wait()` await. Bench-owned CPU fan-out sites: zero. The only remaining
width owner is the Go toolchain's package parallelism inside the single test
phase, which #10 already established does not read a token pool.

The `go list -json -deps ./cmd/bench` module-closure derivation
(`internal/freshness.buildInputs`) measures 0.10–0.31 s warm across three
runs, inside the ~1.6 s gate setup. The #25 memoization residual is dead on
this workload.

### Decision #26 result

The destination's premise is structurally gone on this baseline: nothing
oversubscribes the box, the ~123 load-average symptom cannot reproduce from
a serial phase table peaking at 25 descendants, and the whole-gate target is
exceeded at 38 s against 120 s. There is no dev-tier saturating class left
for #8 to certify `r` against — the pool's would-be clients no longer exist.
Whether #8 and the pool destination retire as achieved-by-other-means is a
policy choice routed to the map's #27, not selected here.

## Validity

The 2026-08-04 section remains one box, one tree, and one repetition per
package. Decision #20 adds a new exact tree, three focused repetitions for the
decision-bearing packages, and one process-profiled fresh gate. The latter is
mechanism and inflation evidence, not full-gate variance. Decision #21 adds
per-test attribution and deterministic Git-spawn counts for the two long dev
packages on #20's exact tree; its shimmed runs are counting evidence, never
timing evidence. Decisions #13–#21's figures describe the pre-rebuild
fixture-driven workload and are historical after `3701c4a0`; decision #26's
section is the only one measured on the current baseline, with two fresh-gate
repetitions and three focused repetitions per long package on one host. No
figure here authorizes a width constant; #27 decides whether any constant is
still needed.
