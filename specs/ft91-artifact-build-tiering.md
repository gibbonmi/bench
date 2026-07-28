# FT91 — artifact build tiering

Status: staged

Compiled from `decisions/cost-follows-project-size.md`, tickets #8, #9, and #10.

**That map was written in the same session as this draft** (the reviewer-closed
path in `/bench-write-spec`'s entry contract): the three forks were put to the
reviewer and answered on 2026-07-27, then written into the map with a complete
Handoff before this spec was compiled from that file. The map is flagged here for
reviewer veto rather than treated as prior sign-off.

**Four calls are the author's, not the map's**, and are marked in-line below. A
top-tier falsification pass found the last two after the first draft claimed only
two, which is the same-session recall risk in action:

1. The opt-in token's name (story 1) — the Handoff left it unchosen.
2. Using the existing `internal/preprelease` unit seam instead of the new
   conformance check the Handoff named (story 5).
3. *Deleting* a pre-existing `reproducibility.json` (story 4). The Handoff says
   only that an opt-in build must never *emit* a record; active removal of a
   record the reviewer's own `prep-release` produced, plus the rollback
   obligation that comes with it, is an interpretation. It is also this spec's
   only destructive file operation.
4. Setting the opt-in in `TestMain` per package (story 8). The Handoff says only
   that "the dev contract phase observes the opt-in path" and is silent on where
   the variable is set.

## Problem

The dev gate pays for a release-evidence property on every commit.

`scripts/build-artifacts.sh` forces an empty private `GOCACHE` per invocation
(line 39) and a fresh `GOMODCACHE` (line 42), then builds everything a second
time from an independent clone to prove the bytes reproduce (lines 129–143). That
hermeticity is deliberate and correct — the script's own comment says the private
cache exists so ambient state cannot affect release bytes.

It is being charged to the wrong tier. A cold-cache build of this module costs
4.79 s against 0.20 s warm, a 24x penalty, and the two contract suites that drive
the script are 249 s of a ~291 s contract phase against a whole gate of ~4m51s.
`internal/contract/surface/artifact` is 133 s across roughly twenty invocations;
`internal/contract/surface` is 115 s, of which `TestPackageContracts` alone is
112.9 s — it runs the full four-platform generator twice for idempotency, each run
double-building, so sixteen cold cross-compiles for one test. Neither suite was
touched by any of FT91's six shipped arms.

The fresh `GOMODCACHE` has a second consequence that is not about speed: every
build re-downloads `github.com/toon-format/toon-go` from `proxy.golang.org`.
Verified 2026-07-27 — `GOPROXY=off bash scripts/build-artifacts.sh` exits 1 with
`module lookup disabled by GOPROXY=off`. The dev gate cannot run air-gapped, while
the same package asserts an offline posture three files away
(`TestOfflineNetworkSentinelDeniesUndeclaredEgress`).

The test harness already tries to do the right thing and is overridden:
`contract.NewExecFixtureAt` (`internal/contract/command.go:303-307`) forwards the
ambient `GOCACHE` and `GOMODCACHE` into every fixture env, and the script stomps
both.

## Solution

Give `build-artifacts.sh` one explicit posture switch, defaulting to what it does
today.

Absent an opt-in the script is byte-for-byte unchanged: private caches, the
independent second build, and the promoted `reproducibility.json`. That is the
posture `bench prep-release` uses, and the release path keeps proving exactly what
it proves now.

With the opt-in, the script honors the ambient Go build and module caches and
skips the second build. Dev contract suites set it, so they prove the generator's
*logic* — that it emits the planned artifact set, idempotently, with reproducible
pins — while byte-reproducibility across independent builds becomes a
once-per-release claim.

Alongside it, `TestPackageContracts` narrows to the host target using the staged
release-plan rewrite the artifact fixture already owns, promoted so both packages
share one copy.

**No coverage leaves the board.** `internal/releaseevidence.readReproducibility`
(defined at `artifact_proofs.go:12`, called from `release_evidence.go:145`)
requires `dist/reproducibility.json` and refuses it unless the schema is 1, the
status is green, `builds` is 2, and every inspected artifact matches by size and
digest. Two independent ship-tier chains reach it: the `release-evidence-probe`
conformance check (`Tier: Ship` at `registry.go:69`) runs its own hermetic
`build-artifacts.sh` and validates the record it produces
(`native_workflow_test.go:244`), and the release-preflight verify step reads the
record that `prep-release`'s `artifacts` step promotes. The two-build property is
verified on the ship tier today. This spec removes a dev-tier duplicate of a
ship-tier proof — it does not remove the proof.

## User stories

1. As `build-artifacts.sh`, I want one posture switch whose absence means
   hermetic, so that no missing environment variable can silently produce a
   non-hermetic release build. Any value other than the exact opt-in token
   resolves to hermetic. **Author's call, flagged:** the token is
   `BENCH_SHARED_BUILD_CACHE=1`; the Handoff left the name unchosen.
   Line: `gpt-5.6-terra` / medium. This is the fail-posture decision the whole
   spec rests on, and getting the polarity backwards would let a release build
   inherit dev posture without anything going red.

2. As `build-artifacts.sh` under the opt-in, I want the ambient Go build and
   module caches honored, so that a warm cache is reused instead of recompiled.
   Both values are resolved through `go env GOCACHE GOMODCACHE` *before* `HOME`
   is overridden, so an explicitly-passed value wins and an inferred one is still
   correct — `go env` reports the environment's value when set and the computed
   default otherwise, which makes one rule cover both callers.
   Line: `gpt-5.6-terra` / medium. The ordering against the `HOME` override is the
   entire correctness of the story and is invisible in the diff if written wrong.

3. As `build-artifacts.sh` under the opt-in, I want the independent second build
   skipped, so that dev pays one build rather than two. This extends the existing
   `BENCH_REPRO_BUILD` guard at line 129 rather than adding a parallel branch.
   Line: `gpt-5.6-luna` / low. The guard already exists and the change is widening
   one condition at a known seam.

4. As a reader of `dist/`, I want no `reproducibility.json` to outlive the
   artifacts it graded, so that a one-build dev run cannot leave a record claiming
   two-build provenance for bytes it did not compare. Under the opt-in the script
   writes no record and removes a stale one in the promotion directory as part of
   the same atomic promotion that replaces the artifacts, restoring it if that
   promotion fails. **Author's call, flagged:** the Handoff requires only that an
   opt-in build emit no record; active deletion and the rollback obligation are
   this spec's reading, and this is its only destructive file operation. The
   alternative — leave the stale record — was rejected because a record describing
   artifacts that no longer exist is worse than none.
   Line: `gpt-5.6-terra` / medium. This touches the promotion/rollback block, where
   a wrong move loses real release evidence rather than a cache.

5. As the release path, I want it pinned that `prep-release` passes no opt-in, so
   that a later edit cannot make the release build non-hermetic without a test
   going red. **Author's call, flagged:** this asserts on the existing step
   definition at the `internal/preprelease` unit seam, where
   `preprelease_test.go:58` already asserts that step's argv. The Handoff named a
   new conformance check; a conformance check would need its own canary family to
   prove it bites, and invariant 4 prefers composing the existing seam.
   Line: `gpt-5.6-terra` / medium. The deviation from the map is the reviewable
   part, not the assertion itself.

6. As the two contract packages, I want one shared staged-release-plan narrowing
   helper, so that host-only breadth is expressed once rather than copied. The
   narrowing core of `committedHostileArtifactSource` moves to `internal/contract`
   and the artifact fixture calls it; its capability-skip and hostless-plan
   options stay with the artifact fixture, which is the only caller that needs
   them.
   Line: `gpt-5.6-terra` / medium. An extraction that leaves the artifact suite's
   existing options behaving differently is a silent regression across twenty
   call sites.

7. As `TestPackageContracts`, I want to prove generator logic and idempotency at
   the host target, so that dev stops building four platforms twice for a
   property two builds of one target demonstrate. Full four-platform breadth stays
   with `prep-release`.
   Line: `gpt-5.6-terra` / medium. Narrowing the plan must not weaken the
   count-derived assertions, which read platform cardinality from the plan.

8. As the dev contract suites, I want the opt-in set once per package rather than
   at each of the roughly twenty call sites, so that a call site added later
   cannot silently reintroduce a hermetic build. Each of the two packages sets it
   in `TestMain`; both `contract.ProcessEnv` (`command.go:289`) and the
   fixture-driven `mergeEnv` (`command.go:195`) build from `os.Environ()`, so it
   reaches every subprocess either path spawns. **Author's call, flagged:** the
   Handoff says only that the dev contract phase observes the opt-in path and is
   silent on where the variable is set.
   Line: `gpt-5.6-terra` / medium. Per-call-site opt-in is the drift this story
   exists to prevent, so the single-source placement is the decision.

9. As `TestDistributableArtifactContracts`, I want the dev-tier two-build
   assertion removed, so that dev asserts only what the dev posture produces.
   `assertPromotedReproducibility` and its single call at `artifact_test.go:107`
   go; the ship-tier proof named in the Solution is unchanged and untouched.
   Line: `gpt-5.6-terra` / medium. Deleting an assertion is exactly the move
   invariant 1 forbids doing casually, so it carries the argument for why the
   property is still proven elsewhere.

10. As a session reading what dev green claims, I want the narrowing stated where
    the tier split already states its own, so that nobody infers per-commit
    reproducibility from a green gate. One clause in `projects/benchkit.md`
    beside the existing dev-versus-ship description, **plus its anchor** in the
    docs-currency `require` list, matching the existing
    `require("projects/benchkit.md", "hostile-input checklist")` precedent — the
    profile checks are mechanical string presence, so an unanchored clause is
    documentation the gate cannot keep honest.
    Line: `gpt-5.6-sol` / high. Kit guidance prose is the leverage override in
    `craft-line`, and this sentence is what a cold session reads to know what
    green means.

## Implementation decisions

**One posture switch, not two.** The opt-in sets both axes — cache sharing and
second-build skip — because they are one decision ("this is a dev-tier build"),
and independent switches would allow a shared-cache build that still claims
two-build evidence.

**`BENCH_REPRO_BUILD` keeps its current meaning.** It marks "I am the inner
reproducibility build, do not recurse" and continues to force private caches. The
new opt-in is a separate axis. The inner build never runs under the opt-in, but
*not* because line 135's `env` filters it — that call has no `-i`, so the ambient
environment passes straight through. The guarantee comes from the exact-match
resolver instead: the second build only runs when the outer build is hermetic,
which means the token was absent or unrecognized, which means the inner build
resolves hermetic too. The mechanism is the resolver, not the `env` list.

**The opt-in is set by test fixtures, never by the gate.** A gate-wide variable
would apply to any script the gate happens to run, including a ship-tier one. Two
`TestMain` functions localize it to the packages whose tier is dev by
construction, and a developer running the script by hand gets the hermetic
default.

**Story 5's seam cannot see ambient inheritance, and that is acceptable because
the failure is fail-closed.** `prep-release` runs each step with
`append(os.Environ(), step.Env...)` (`preprelease.go:247`), deliberately
inheriting the ambient environment, so a unit assertion on `step.Env` cannot
catch an opt-in token that leaked in from the surrounding shell — which is
exactly what the Handoff's conformance check was for. The substitution loses
nothing only because of a second argument the first draft did not make: an
ambient token makes the artifacts step emit no record, and both ship-tier chains
refuse a missing record, so the run goes red rather than silently shipping
non-hermetic bytes. The unit pin catches the likely regression (someone adds the
token to the step definition); the ship tier catches the unlikely one.

**npm cache, `TMPDIR`, and `HOME` stay private in both tiers.** They were not
measured as cost drivers and npm's cache concurrency posture is untested, so the
narrowing stays scoped to what the measurement justifies.

**Expected effect**, from the measured probe rather than estimated: the artifact
suite fell 133.5 s to 73.2 s (-45%) with a shared cache alone, and story 7's
breadth narrowing is worth roughly 85 s of `TestPackageContracts`'s 112.9 s. The
residual in both is npm pack, node, and git-clone work, which this spec does not
address.

## Testing decisions

A good test here drives `build-artifacts.sh` through its real argv and observes
its filesystem output — artifact names, counts, pins, and the presence or absence
of `reproducibility.json`. Nothing asserts on wall clock: speed is the motive but
timing is not a stable oracle, and a timing assertion would flake under gate
contention.

Prior art: `internal/contract/surface/artifact/artifact_source_state_test.go`
already drives the script with a crafted env and asserts refusals, and
`internal/preprelease/preprelease_test.go` already asserts step argv without
executing the step.

Gate command: `bench gate`.

### Seam diagram

Seam 1 — the posture switch (stories 1–4):

    trigger: contract fixtures (opt-in set) │ prep-release artifacts step (unset)
        │
        ▼
    BENCH_SHARED_BUILD_CACHE  ──▶  [ build-artifacts.sh          ]  ──▶  artifacts/*.tgz
    source root, output dir   ──▶  [   posture resolver          ]  ──▶  reproducibility.json
    ambient GOCACHE/GOMODCACHE──▶  [   cache exports             ]       (hermetic only)
                                   [   second-build guard        ]  ──▶  exit code
                                        ◀ tests attach here: run the real script
                                          with and without the token, then read
                                          the output directory

Seam 2 — the release-path pin (story 5):

    trigger: internal/preprelease unit test
        │
        ▼
    graded root  ──▶  [ preprelease.Steps() ]  ──▶  []Step{Name, Argv, Env}
                           ◀ tests attach here: read the "artifacts" step and
                             assert its Env carries no opt-in token

Seam 3 — host-only breadth (stories 6–7):

    trigger: TestPackageContracts
        │
        ▼
    committed candidate ──▶  [ contract.NarrowReleasePlan ]  ──▶  candidate with
    host GOOS/GOARCH    ──▶  [                            ]       1-target plan
                                  │
                                  ▼
                             [ gen-platform-packages.sh ] ──▶ artifacts + pins
                                  ◀ tests attach here: existing generator
                                    assertions, now over a narrowed plan

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | absent token builds hermetically | seam 1 | `go test ./internal/contract/surface/artifact -run TestArtifactBuilderHonorsHermeticDefault` | Asserts `reproducibility.json` exists with `builds: 2` when no token is set; an inverted default produces no record and the row goes red. |
| 1 | unrecognized token value builds hermetically | seam 1 | same test, `BENCH_SHARED_BUILD_CACHE=yes` case | A truthy-ish parse (`-n`, `!= 0`) would share the cache and drop the record; exact-match keeps it. |
| 1 | hermetic path still refuses offline | seam 1 | `go test ./internal/contract/surface/artifact -run TestBuildCachePostureUnderGoproxyOff` (token unset) | Polarity guard: the hermetic path must still fail under `GOPROXY=off`, so an implementation that shares caches unconditionally reds here instead of quietly passing the opt-in rows. |
| 2 | opt-in honors the ambient module cache | seam 1 | same test, token set, module pre-resolved into the ambient `GOMODCACHE` | An opt-in build under `GOPROXY=off` exits 0 only if it reuses the ambient module cache. An implementation that keeps a private `GOMODCACHE` must download and exits 1 — the exact failure verified against the current script on 2026-07-27. This replaces a cache-directory-inspection signal that could not tell "reused the seed" from "ignored the seed". |
| 2 | caches resolved before the `HOME` override | seam 1 | same test, with `GOCACHE`/`GOMODCACHE` absent from the fixture env | If resolution moves after the override, `go env` computes paths under the private `HOME`, the module is not there, and the `GOPROXY=off` build fails. The ordering the story calls load-bearing is observable through this one lever. |
| 3 | opt-in skips the second build | seam 1 | same test, token set | The inner build at `build-artifacts.sh:135` forces its own private `GOMODCACHE`, so if it still runs under the opt-in the offline build fails. This is the discriminator the record-absence check could not provide: an implementation that suppresses the record but still rebuilds stays green on the record and reds here. |
| 3 | opt-in promotes no record | seam 1 | `go test ./internal/contract/surface/artifact -run TestSharedCacheBuildPromotesNoRecord` | Asserts no `reproducibility.json` appears; paired with the row above so neither "no record" nor "no rebuild" can be satisfied alone. |
| 4 | stale record removed under opt-in | seam 1 | `go test ./internal/contract/surface/artifact -run TestSharedCacheBuildRemovesStaleReproducibilityRecord` | Seeds a `reproducibility.json` beside the output, runs an opt-in build, asserts it is gone; leaving it would let a one-build run present two-build evidence for bytes it did not compare. |
| 4 | failed opt-in promotion restores the record | seam 1 | `go test ./internal/contract/surface/artifact -run TestSharedCacheBuildRestoresRecordOnInterruptedPromotion` | Drives the existing `BENCH_TEST_PROMOTION_READY_FILE` seam, which blocks *inside* the promotion block at line 146, then interrupts — so the failure lands after the record has been moved. The tarball-count refusal at line 127 fires before promotion and can never exercise the trap, which is why it is not the injection used here. |
| 5 | prep-release passes no opt-in | seam 2 | `go test ./internal/preprelease -run TestArtifactStepIsHermetic` | Reads the `artifacts` step and fails if the token appears in its `Env`; a later edit adding it reds here without needing a release run. Ambient inheritance is out of this seam's reach by design — see the fail-closed argument in Implementation decisions. |
| 6 | narrowing helper preserves artifact-suite behavior | seam 3 | `go test ./internal/contract/surface/artifact` (whole package) | already covered — the twenty existing call sites are the regression suite for the extraction; any behavior change in the moved code reds one of them. |
| 7 | generator proves idempotency at host breadth | seam 3 | `go test ./internal/contract/surface -run TestPackageContracts` | The existing idempotency and pin-reproducibility subtests keep their assertions over the narrowed plan; a narrowing that broke plan-derived counts reds the artifact-count subtest. |
| 7 | artifact count still derives from the plan | seam 3 | same test, `platform-package generator failed` subtest | It computes `len(packageReadPlatforms(t))*2 + 1`; a hardcoded count after narrowing would pass for the wrong reason, and reading the narrowed plan is what keeps the derivation honest. |
| 8 | opt-in reaches every fixture-driven call site | seam 1 | `go test ./internal/contract/surface/artifact ./internal/contract/surface` and compare against the recorded baselines | not TDD-able as a hard assertion — no observable distinguishes "every call site opted in" from "most did" except duration, which is not a stable oracle. The `TestMain` placement is what makes it structurally true; the baselines (133.5 s and 115.5 s, measured 2026-07-27) are recorded so a regression shows up in the gate's per-phase timing output. |
| 9 | ship tier still proves reproducibility | — | `bench prep-release` (ship tier only) | already covered, and deliberately not a dev signal — `release-evidence-probe` (`Tier: Ship`) runs a hermetic build and validates the record through `readReproducibility`, which refuses a missing, incomplete, or digest-mismatched record. Named here because story 9's deletion depends on it. The first draft cited `go test ./internal/releaseevidence -run TestReproducibility`, which matches no test and exits 0 — a dead citation of exactly the class FT133 exists to catch. |
| 10 | the profile states what dev green claims | conformance docs anchor | `BENCH_CONFORMANCE_ROOT="$PWD" go test ./internal/conformance -run TestRootConformance` | Story 10 adds the clause *and* its anchor to the docs-currency `require` list, matching the existing `require("projects/benchkit.md", "hostile-input checklist")` precedent, so deleting the clause reds the check. Without the anchor these checks are mechanical string-presence only and grade no tier-claim semantics — the first draft asserted coverage that did not exist. |

### Edge inventory

Walked against the profile's hostile-input checklist for shell CLIs.

- **Paths with spaces or glob characters** — covered: the existing fixtures
  already use `"host-only artifact output [*]"` and `"authenticated candidate"`,
  and the new tests reuse those directories.
- **Absent vs present-but-empty `reproducibility.json`** — covered by story 4's
  rows: absent is the opt-in steady state, present is the stale-record case. A
  zero-byte record is treated as stale and removed, same as a populated one.
- **A command whose own write changes a fact it reports** — covered: the
  promotion is the write, and story 4's rollback row asserts the tracked
  configuration survives a failed promotion.
- **Hand-edited file with no trailing newline** — **Won't handle**: the script
  never parses `reproducibility.json`, only creates or removes it.
- **Special files and dangling symlinks in the output path** — **Won't handle**:
  the output directory's hostility is owned by the existing promotion block and is
  unchanged by this spec; no new path is read.
- **Required tool missing from PATH** — **Won't handle** for `go`: an absent
  toolchain fails the build long before posture matters, and the existing
  absent-toolchain diagnostic owns it.
- **Unquoted multi-word arguments** — covered: the token is compared as a single
  quoted string, and the existing spaced-path fixtures exercise argv quoting.
- **A flag's value read as a positional** — **Won't handle**: the script takes two
  positionals and no flags, so the class has no surface here.
- **Invocation through every shipped surface** — `build-artifacts.sh` has five
  callers, not the three the first draft named. Three are covered by rows:
  `prep-release` (story 5), the dev contract fixtures (story 8), and
  `gen-platform-packages.sh` (story 7). Two more take the hermetic default and are
  **Won't handle** because absence of the token is their whole posture and no row
  can add to that: the ship-tier `release-evidence-probe`
  (`native_workflow_test.go:244`) and the CI artifacts workflow
  (`.github/workflows/native-runtime.yml:52-53`, itself pinned by
  `workflow_checks_test.go:62-72`, which greps for its clone-and-compare).
- **Concurrent invocation** — covered by the existing output-directory lock, which
  this spec does not touch; sharing a Go build cache across concurrent invocations
  is safe because Go's cache takes its own locks.

## Out of scope

- **The artifact suite's ~73 s warm residual** — a separate capability, and the
  first draft's reasoning for cutting it was wrong. Measured 2026-07-27: it is
  not per-tarball packing cost. One warm full-matrix invocation with the second
  build skipped is 19.26 s, of which only ~5.2 s is the four cross-compiles; npm
  pack is 0.61 s per tarball and is *unaffected* by its cache being fresh
  (0.61 s against 0.59 s shared), node script startup is 0.02 s, and a repo clone
  is 0.15 s. The residual is therefore **invocation count, not per-invocation
  cost** — roughly twenty host-only generator runs at ~3.7 s each.

  The lever is the `BENCH_TEST_PREPARED_ARTIFACTS` seam that already exists
  (`build-artifacts.sh:82-85`): `artifactPreparedGeneration` builds a real
  artifact set once and hands it to tests that only need artifacts to *exist*,
  but it is scoped per-test, so `TestArtifactPromotionIsAtomicAndExclusive`
  still pays a full build (7.98 s warm) to produce its own prepared set.
  Hoisting that generation to package scope would let one build serve every test
  that only inspects output.

  It stays out of scope because it is a test-isolation decision, not a
  build-posture one: sharing one artifact set across tests trades each test's
  independence for time, and which tests may safely share is a reviewer call this
  map never asked. Estimate deferred rather than guessed — sizing needs a count
  of how many of the seven build-bound tests only inspect output, which this spec
  did not establish.
- **Vendoring `toon-go`** — a separate capability (dependency-management posture,
  governed by the dependency standard in `AGENTS.md`), and ticket #10 chose the
  shared module cache instead. It would make a bare-machine first build offline,
  which the shared cache does not. Estimate: 3 edits, 2 gate runs.
- **`-count=1` freshness semantics** — still parked in the map's
  `## Not yet specified` with its own trigger; an oracle-semantics decision, not a
  cache one. Estimate: unknown until the decision closes.
