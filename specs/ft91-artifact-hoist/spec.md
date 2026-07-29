# ft91-artifact-hoist — package-scoped prepared-artifact set

Status: implemented

## Problem

Every test in the artifact contract suite that needs a built artifact set pays
its own full `build-artifacts.sh` generation (~3.7 s warm), even when the test
only inspects or consumes prepared output. Up to 7 of the suite's ~17
generation passes are host-only builds of the identical recipe. The suite runs
~106 s solo and ~158 s under gate parallelism, and five `surface/artifact`
canary fixtures each re-pay the whole suite, so the waste multiplies about
six-fold. This is the last open slice of `gate-critical-path` before the #5
≤60 s re-measurement.

## Solution

A package-scoped lazy singleton inside the artifact test package: the first
sharer stages one host-narrowed committed clone and runs one real shared-cache
build; the six reviewer-ruled sharers consume the staged source and promoted
output read-only. Two belts guard independence: the set is chmod'd read-only
after generation (a mutator fails at write time, attributed), and a digest
fingerprint recorded at build is re-verified before every consumer use,
failing loud as "shared artifact set mutated". Every requester also asserts
exactly one build has occurred, so a rebuild-per-call regression reds the
suite. No sharer selected → no build. Posture, refusal, and uniquely-shaped
tests keep their private builds — their build is the subject.

All stories run in this session on the top tier — the reviewer's explicit
2026-07-29 instruction ("this session builds it"), a declared deviation from
the profile's cached mid-effort routing for gate/conformance-adjacent work.
`fable` is Claude Code's top-tier alias per `.bench/lines.env`.

## User stories

1. As the gate, I want a lazy in-package singleton whose first requester
   stages one host-narrowed committed clone and runs one real shared-cache
   build, exposing the staged source root, the promoted output directory, the
   promoted entry count, and the build fingerprint, so every ruled sharer
   consumes one generation instead of paying its own.
   Line: fable (top alias, this session) / medium, cap 4 focused-loop
   iterations. The reviewer directed the in-session build, and the seam
   contract is fully decided in the map, so medium effort suffices.
2. As a sharer test, I want the shared set (staged source and promoted output)
   made read-only after generation, so a test that writes into it fails at the
   write with attribution instead of corrupting a later test's input.
   Line: fable (top alias, this session) / medium, cap 3 iterations. The
   belt's interaction with git and the build script needs care beyond
   mechanical editing.
3. As a sharer test, I want a digest fingerprint recorded at build and
   re-verified before each consumer use, failing loud with a message naming
   the shared artifact set as mutated, so corruption is attributed to the set
   rather than surfacing as a confusing downstream assertion.
   Line: fable (top alias, this session) / medium, cap 3 iterations. The
   verify path is the failure-attribution contract every sharer depends on.
4. As a developer running `go test -run` selecting only fabricated-fixture
   tests, I want no artifact build to occur, so narrow selections stay cheap
   and the singleton stays lazy rather than moving the cost to `TestMain`.
   Line: fable (top alias, this session) / medium, cap 3 iterations. The
   probe follows an existing in-package subprocess pattern but its oracle
   must survive the degenerate implementations named in the coverage map.
5. As a sharer running after the first requester's staging failed, I want a
   loud red naming the earlier shared-set failure, so a half-staged set can
   never be consumed and the failure never surfaces as a nil dereference.
   Line: fable (top alias, this session) / medium, cap 2 iterations. The
   guard is small but it is the fail-closed posture of the whole seam.
6. As the artifact suite, I want the four inspection/consumption sharers —
   `TestArtifactSourceStagesCommittedHostPlan`,
   `TestSharedCacheBuildPromotesNoRecord`, `TestOfflineArchiveProjection`,
   `TestPackedArtifactRunsSetupOfflineFromASpacedPrefix` — migrated to consume
   the shared set with their existing assertions unchanged, so their per-test
   builds disappear.
   Line: fable (top alias, this session) / low, cap 2 iterations. Mechanical
   migrations at a decided seam, fully graded by the tests' own existing
   assertions.
7. As the artifact suite, I want the two prepared-seam promotion sharers —
   `TestArtifactPromotionIsAtomicAndExclusive`,
   `TestSharedCacheBuildRestoresRecordOnInterruptedPromotion` — migrated to
   use the shared promoted output as their prepared input and its entry count
   as their expected count, so their extra seeding generation disappears while
   their own builds and outputs stay private.
   Line: fable (top alias, this session) / low, cap 2 iterations. Mechanical
   migrations at a decided seam, fully graded by the tests' own existing
   assertions.
8. As the contract helper library, I want the release-plan staging helper
   generalized with a variant that stages the clone under a caller-owned
   directory, with the existing `contract.NarrowReleasePlan` delegating to it,
   so the singleton's clone survives the first requester's `t.TempDir` cleanup
   without duplicating the staging logic.
   Line: fable (top alias, this session) / medium, cap 2 iterations. The one
   edit outside the artifact package; it must keep one source for the staging
   fact.
9. As `TestMain`, I want the package-owned shared directory removed after the
   run — write permission restored, then deleted — so read-only sets do not
   accumulate in the host temp directory across suite runs.
   Line: fable (top alias, this session) / low, cap 2 iterations. A small
   mechanical teardown.
10. As the reviewer, I want the source decision map
    `decisions/gate-critical-path.md` and its map-owned assets
    (`artifact-test-inventory.md`, `gate-critical-path-timeline.md`,
    `behavior-owned-package-bindings.md`) moved — not copied — under
    `specs/ft91-artifact-hoist/decisions/` with every reference updated
    (the referrer set, verified complete by search: `ROADMAP.md`,
    `decisions/gate-pipeline.md`), so the compiled provenance travels with
    the spec.
    Line: fable (top alias, this session) / low, cap 1 iteration. A
    mechanical move with reference fixes.
11. As FT91's stop rule, I want the post-change artifact-suite time and the
    full dev-gate wall-clock measured and recorded in the timeline asset at
    its new spec-local path, so the #5 ≤60 s judgment happens on evidence
    rather than expectation.
    Line: fable (top alias, this session) / low, cap 1 iteration. A
    mechanical measure-and-record step the Handoff names as this slice's
    ship evidence.

## Implementation decisions

- The whole change is test-harness work: the artifact contract test package
  gains the singleton and the migrations; `internal/contract` gains only the
  caller-owned-directory staging variant (story 8). `scripts/build-artifacts.sh`
  and all production packages are untouched. **Flagged for veto:** the map's
  boundary line read "inside `internal/contract/surface/artifact`"; the story-8
  edit is a refinement forced by `t.TempDir` lifetime — `NarrowReleasePlan`
  stages into the first requester's temp directory, which is destroyed when
  that test ends, so the shared clone needs a package-owned home. The variant
  keeps one source for the staging logic instead of copying it into the
  artifact package.
- The singleton is an accessor (working name `requireSharedArtifactSet`) over
  a state struct holding a `sync.Once`, the handles, and a build counter.
  First caller resolves `contract.SubjectRoot`, creates one package-owned
  `os.MkdirTemp` directory (spaced name, per package convention), stages the
  host-narrowed committed clone into it via the story-8 variant, runs one
  real build with the same recipe sharers use today
  (`contract.NewExecFixtureAt(t, root).Run(...)`, inheriting `TestMain`'s
  shared-cache posture) **with the promoted output directory also under the
  package-owned directory** — the output must not live in any test's
  `t.TempDir`, and its parent directory is load-bearing (the build stages
  into it, and `TestSharedCacheBuildPromotesNoRecord` asserts record absence
  beside the output; nothing else may write a reproducibility record there).
  It then chmods the set read-only and records `promotedArtifactDigests`.
- Exactly-one-build is a consumer-side assertion: every call to the accessor
  asserts the build counter equals one after the `Once` resolves, so a
  regression that rebuilds per call reds at the second sharer. The build step
  runs through one function, which is also the laziness observable: when the
  test-only env knob `BENCH_TEST_SHARED_SET_BUILD_LOG` names a file, that
  function appends one line to it. The knob is read by test code only; it is
  the story-4 probe's oracle because it survives the degenerate
  implementations (an eager `TestMain` build has no `*testing.T`, so a log
  call could not; any implementation that builds routes through the one build
  function, so the append happens regardless of who triggered it).
- The accessor's post-`Once` handle resolution is an error-returning function
  on the state struct with a thin `t.Fatal` wrapper: nil state after a
  consumed `Once` returns the named error ("shared artifact set staging
  failed in an earlier test"). `t.Fatal` inside the first requester marks the
  `Once` done via `runtime.Goexit`, so this guard is the fail-closed belt for
  every later requester — and because it is error-returning on the struct, a
  unit test drives a fresh state instance's `Once` with an empty function and
  asserts the exact message without touching the package's real singleton.
- The verify step is likewise an error-returning function with a thin
  `t.Fatal` wrapper, so the mutation probe can assert the exact attributed
  message without redding itself.
- Read-only means files `0o444`, directories `0o555`, walked over both the
  staged source clone and the promoted output. Any chmod or digest error
  during staging fails the first requester loud (helper `t.Fatal`), fail
  closed — later requesters then red through the consumed-`Once` guard.
  Accepted contingency: `build-artifacts.sh` runs `git status` on the source,
  and if that requires a writable `.git`, the source belt narrows to the
  worktree payload and `.git` stays writable — the write-time probe targets
  one promoted tarball **and** one staged-source worktree file (not `.git`),
  so the probe's coverage is unaffected if the contingency fires, and the
  fingerprint belt covers the artifact set either way. The singleton carries
  a comment stating the package's no-`t.Parallel` assumption: sharing
  analysis was done for sequential ordering hazards, and introducing
  parallelism changes it to races.
- `TestMain` captures `m.Run()`'s code, restores write permission over the
  package-owned directory, removes it, then exits — cleanup runs on failure
  too.
- Fresh per `go test` invocation by construction: the directory comes from
  `os.MkdirTemp` and no path is persisted, so there is no cross-run staleness
  and no contact with the parked oracle-semantics levers.
- Non-sharers are untouched: hermetic, GOPROXY-off, and stale-record posture
  tests, both refusal tests, `TestDistributableArtifactContracts`, and
  `TestArtifactBuilderRefusesMissingBinaryPinManifest` keep their private
  builds — their build is the subject. The pin-manifest test stays private
  **by decision** (Handoff item 7's recommendation, adopted): the prepared
  seam skips the whole generation branch including the pin check
  (`build-artifacts.sh` copies prepared artifacts and never reaches the pin
  refusal), so joining would need a new production-script seam, which crosses
  the module boundary above.
- Build and review venue (reviewer-directed, 2026-07-29): this session builds
  the spec on the top tier; the implementation review runs as a Codex CLI
  delegate on `gpt-5.6-sol`, high reasoning effort, yolo mode. Both are
  declared deviations from the profile's cached routings (mid-tier build
  effort for gate-adjacent work; mid model, medium effort for review-axis
  delegates), made on the reviewer's explicit instruction. Overall build cap:
  ~12 focused-loop iterations plus 2 full gate runs; exhaustion stops and
  reports rather than grinding.

## Testing decisions

- A good test here exercises the package's external behavior: the sharers'
  existing assertions are the primary oracle and stay byte-for-byte
  unchanged where possible; new tests assert the two belts, the
  exactly-one-build fact, and laziness at the singleton's observable surface
  (write failure, verify error message, build-log absence), never its
  internals.
- Seam: the singleton accessor inside the artifact test package, plus the
  package suite itself. Prior art: the subprocess re-invocation pattern in
  `TestArtifactSourceSkipsWhenHostTargetIsAbsent` (story 4's probe) and the
  existing `promotedArtifactDigests` helper (story 3's fingerprint).
- Gate: `.bench/gate.sh` (`bench gate`). Focused loop during the build:
  `go test -count=1 ./internal/contract/surface/artifact`.
- Wall-clock is not gate-assertable; ship evidence is the story-11
  measurement recorded in the timeline asset (Handoff item 5).

### Seam diagram

    trigger: first ruled sharer calling requireSharedArtifactSet(t)
        │
        ▼
    committed tree ──▶ [ stage host-narrowed clone → one shared-cache build ]
                       [ chmod read-only → record digest fingerprint        ]
                            │ (build function appends to
                            │  BENCH_TEST_SHARED_SET_BUILD_LOG when set)
                            ▼
                       {source root, output dir, entry count} ──▶ six sharers
    later sharer ──▶  [ verify fingerprint + assert builds == 1 ] ──▶ same
                       handles, or "shared artifact set mutated" red
                      ◀ tests attach here: sharers consume the handles;
                        probes drive the verify error path, the write-time
                        chmod failure, the consumed-Once guard on a fresh
                        state instance, and a -run subprocess asserting an
                        empty build log

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | all six sharers pass against the one shared generation | artifact package suite | already covered — the six migrated tests keep their existing assertions (host-only plan, planned artifact names, record absence, offline archive set, npm install + setup convergence, promotion atomicity/restore) | any divergence between the shared recipe and what each test's private build produced fails that test's own unchanged assertions |
| 1 | exactly one generation occurs across a full suite run | consumer-side build-count assertion in the accessor | observed red required once during the build: demonstrated against a deliberate rebuild-per-call variant, which reds at the second sharer | a singleton that quietly rebuilds per call (the cheapest wrong implementation of the slice's whole payoff) fails the count assertion the moment two sharers run |
| 2 | a write into the shared set fails at write time with a permission error, probed on one promoted tarball and one staged-source worktree file | new probe test in the package | observed red required before green: probe written first, red while the chmod belt is absent because the writes succeed | a successful write is exactly the belt's absence; probing both members means the source-belt contingency cannot silently shrink coverage |
| 3 | fingerprint verify reds with a message naming the shared artifact set as mutated | new probe test (error-returning verify on a deliberately mutated then restored file) | observed red required before green: probe red while verify does not exist or returns nil on a mutated set | asserts the exact attributed message consumers rely on, so mutation can never surface as a downstream assertion |
| 4 | `go test -run` selecting only fabricated-fixture tests invokes no build | subprocess probe re-invoking the test binary with the build-log knob set, asserting the log stays absent | observed red required once during the build: demonstrated against a deliberately eager variant, whose build appends to the log through the one build function | the oracle is intrinsic to the build function, not a courtesy log line, so an eager `TestMain` build (which has no `*testing.T`) still writes it and still reds the probe |
| 5 | a requester after a failed staging reds with the named earlier-failure message | unit test driving a fresh state instance's `Once` with an empty function and asserting the error | observed red required before green: test red while the error-returning guard does not exist or returns nil on a consumed `Once` with nil state | exercises the exact consumed-`Once`-with-nil-state path a first-requester `t.Fatal` produces, without touching the package's real singleton |
| 6, 7 | migrated sharers keep their existing behavior | artifact package suite | already covered — the assertions predate this spec and stay unchanged | the migrations only swap the input source; the assertions grade the swap |
| 8 | the two existing `NarrowReleasePlan` callers (`package_test.go`, `artifact_fixture_test.go`) behave identically through the delegating variant | contract suites (both callers run in the gate's contract phase) | already covered — both callers' suites red on a staging behavior change | a behavior change in the extraction reds the callers that stage plans today |
| 8 | the variant delegates rather than duplicating the staging logic | — | not TDD-able — one-source-per-fact is structural, with no behavioral difference to observe; graded on review against the code standard | — |
| 10 | no reference to the moved map or assets dangles | — | not gate-covered — the conformance sweep checks command tokens, not file paths; graded by an `rg` sweep over the moved paths run at build time and re-checked at review | — |

### Edge inventory

Walked per behavior against the canonical classes and the profile checklist's
applicable rows (this change is Go test harness, not a shell CLI surface; the
inapplicable CLI classes — control bytes, TTY, flag parsing — are owned by the
untouched script's existing tests):

- Error path: singleton staging/build failure → first requester reds via the
  fixture's `RequireExit`; a chmod or digest error during staging → first
  requester reds via helper `t.Fatal`, fail closed; every later requester →
  story 5 row.
- Empty/absent input: no sharer selected → story 4 row; an empty or missing
  prepared directory → existing script refusals, untouched non-sharer tests.
- Boundary: smallest observable mutation (one byte appended to one tarball) →
  story 3 probe exercises exactly this.
- Malformed input: prepared-artifacts path not a directory, or a symlink →
  existing script refusal (`build-artifacts.sh`), already gate-covered.
- Mode propagation: `cp -a` copies the shared output's `0o444` files into the
  story-7 consumers' private staging and promoted outputs; their subsequent
  moves and re-reads need only parent-directory permissions, and the existing
  story-6/7 rows red if this breaks — walked as a build risk, not a coverage
  hole.
- Interrupted/partial state: SIGINT killing the whole test process leaks the
  package-owned temp directory — **Won't handle:** identical exposure to
  today's `t.TempDir` under a hard kill; OS temp reaping owns it.
- Re-run idempotency: fresh per `go test` invocation by construction (no
  persisted path) — **Won't handle:** cross-run staleness is impossible
  without an on-disk cache, and the on-disk cache was a rejected alternative
  (map item 8).
- Hostile environment: spaced paths — the shared directories carry spaces per
  package convention, so every story-1 row runs against them; read-only
  source versus `git status`'s index refresh — accepted contingency recorded
  in Implementation decisions, with the story-2 probe unaffected.
- Cleanup hygiene: `TestMain` teardown of a read-only tree — **Won't handle**
  as a test: no oracle observes host temp hygiene; verified once during the
  build and graded on review.

## Out of scope

- **Oracle-semantics levers** (gate-verdict caching keyed on the pinned
  subject; `-count=1` freshness) — a separate capability with its own parked
  decisions in `cost-follows-project-size.md`; graduates into tickets only if
  story 11's re-measure stays above 60 s (#5's ruling). ~25 edits, ~8 gate
  runs.
- **A prepared-binaries production seam** so the pin-manifest test could
  share — a separate capability crossing into `scripts/build-artifacts.sh`,
  declined per Handoff item 7. ~10 edits, ~4 gate runs.
- **Reviving the outer conformance/contract width cap** — dormant per
  `gate-concurrency.md`; re-check belongs after story 11's measurement, not in
  this diff. ~4 edits, ~2 gate runs.
