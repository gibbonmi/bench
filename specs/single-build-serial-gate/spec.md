# single-build-serial-gate

Status: staged

Decision source: reviewer-confirmed current conversation, 2026-08-08: every local run builds the Bench executable once per exact build environment and reuses it everywhere, and Bench-created resource work has zero overlap.

## Problem

A local Bench gate still creates its own contention. The primary subject has a
published `dist/bench`, but gate toolchain phases invoke the same command through
`go run`; every gate and nested canary entry also invokes
`go run ./internal/freshness/check`; deterministic gate fixtures share a staged
executable only inside one Go test process; and other test helpers call builders
directly. The outer gate launches ready phases concurrently, the primary and
stripped phase sets overlap, canary has its own worker pool, and
`internal/gate` now contains 196 `t.Parallel` calls. Lowering `GOMAXPROCS`
narrows each process without preventing those independent processes from
overlapping.

The earlier fixture fix was implemented as designed, but its design explicitly
stopped at process-local reuse and left scheduling to later FT171 work. That is
too narrow for the required local behavior. The result is repeated compiler,
linker, artifact-write, and test-process demand competing for disk and memory
inside one command.

## Solution

One repository-local artifact store owns a content-addressed executable set,
and one run authority owns a transferable capacity-one resource turn for the
entire process tree. The store resolves the exact build environment under the
requested target, builds one immutable executable for each distinct identity,
and supplies that artifact to every eligible gate, package process, fixture,
and nested child. Sibling package processes in a raw `go test ./...` rendezvous
through the same store and per-identity lock instead of creating process-local
coordinators. Successful artifacts persist for that source identity; output
paths, fixture roots, invocations, and process identities do not create new
identities. Changed target-selected build inputs, toolchains, targets, build
modes, or a registered independent-authorship slot do.

The same authority makes the ordinary gate non-concurrent by construction. It
runs one outer phase lineage at a time, runs the primary and stripped tables as
one ordered schedule, gives Go commands one package slot and one runtime
thread, runs one canary item at a time, and serializes resource creation inside
the remaining tests that must exercise coordination. A synchronous parent
delegates its one turn to a child and becomes quiescent until the child returns;
it never holds a second turn or competes with its descendant. A gate can still
test concurrency semantics, but no two Bench-created resource actors are
runnable at once. External host processes are not part of this claim.

This supersedes the positive-width and work-conserving portions of
`decisions/gate-budget.md` and the retired `gate-test-concurrency` spec. It does
not reopen the requirements that every gate phase remains present, green means
the same thing, failures aggregate honestly, cancellation reaps the process
tree, and compiler-observing proofs perform real builds.

## User stories

1. As a Bench developer, each local gate, focused package run, or multi-package
   Go run produces at most one successful Bench CLI or freshness-verifier
   executable for each exact build environment, and every consumer that does
   not observe compilation uses that immutable artifact instead of invoking
   `go build`, `go run`, or the subject builder again.
   Line: `gpt-5.6-sol` / high. This is a cross-process artifact identity and
   freshness boundary whose cheapest mistake silently executes the wrong
   subject or preserves most of the repeated work.
2. As a Bench developer, an ordinary gate creates no resource contention of
   its own: exactly one phase lineage owns the transferable resource turn, and
   at most one Go package/test slot, canary item, or nested resource-producing
   actor is runnable at a time, on every host and at every nesting depth.
   Line: `gpt-5.6-sol` / high. This changes the oracle scheduler and recursive
   process authority; cancellation and fail-closed propagation make it a
   highest-risk gate seam.

## Implementation decisions

- **One deep authority, two lifetimes.** A concrete coordinator, not a
  caller-defined provider interface, owns a repository-local artifact store
  and a run-local capacity-one turn. Successful artifacts and their identity
  records persist under the repository's common Git directory; the resource
  turn and process authority end only after a run reaps every descendant.
  Non-Git fixtures use the owning kit repository's store. Store discovery,
  locking, record publication, validation, and repair diagnostics stay hidden
  behind the coordinator.
- **Sibling processes rendezvous.** An inherited child joins its parent's run
  authority or fails closed. A focused or multi-package Go test process with no
  inherited authority joins the repository store and its per-identity
  exclusive attempt record; it does not create a private artifact cache.
  Therefore sibling package processes from raw `go test ./...`, separate
  focused invocations, and a gate all converge on one successful artifact for
  the same identity. A repository-local artifact lock coordinates only
  executable authorship; unrelated test execution remains outside story 1.
- **Identity follows target-selected executable inputs.** The reusable key is
  the canonical freshness build-input digest computed with the requested
  `GOOS`, `GOARCH`, `CGO_ENABLED`, toolchain, and build tags, plus the Go
  toolchain executable identity, normalized build flags, package/version stamp,
  artifact class, and build mode. The resolver runs its dependency/file
  enumeration under that target rather than reusing the host's file set.
  Windows-only, Unix-only, architecture-only, and CGO-selected sources are
  required identity cases. The physical root is provenance used to compute and
  recheck the key, not part of it; destination path, worktree spelling, mtime,
  invocation, fixture, and process identity are excluded.
- **One successful artifact and one durable answer per identity.** Resolution
  is cross-process single-flight: check the completed record, acquire the
  per-identity lock, check again, then author. Concurrent and later requests
  receive the same immutable artifact and executable digest. Success is
  atomically durable. Failure or cancellation publishes no success and every
  waiter on that attempt receives the same error; an automatic retry is
  forbidden. A source change creates a new identity, while retrying an
  unchanged failed identity requires the explicit exact-identity artifact
  repair. Repair refuses an active attempt, removes only the named failed or
  corrupt record/artifact, and lets the next request author once. A corrupt or
  missing completed artifact refuses with that repair route rather than
  silently rebuilding.
- **Materialization policy is registered and mutation-safe.** The store artifact
  is read-only and never passed directly to `freshness.Publish`, whose rename
  consumes staging. Each consumer registry row declares either
  `direct-exec`, `immutable-link`, or `mutable-copy`. Direct execution uses
  the store path. An immutable link is permitted only while the destination
  remains non-writable and every mutation API first atomically detaches it;
  a consumer that exposes arbitrary byte or metadata mutation receives a copy.
  Root-local publication and seals remain independently owned. The validator
  refuses a writable multiply-linked inode, and the static ownership check
  rejects direct mutation of an immutable-link destination.
- **The CLI and bootstrap verifier both enroll.** Artifact classes include the
  canonical Bench CLI and a freshness-verifier executable. Preparation builds
  each class once per target identity and publishes both records. The direct
  gate entry executes the trusted verifier artifact to validate `dist/bench`
  before executing the CLI; it never runs `go run
  ./internal/freshness/check`. Nested canary entries inherit that verifier
  identity, eliminating the per-bite compile while preserving refusal before
  untrusted CLI execution.
- **Every ordinary CLI consumer enrolls.** The shared CLI artifact supplies
  gate plumbing that currently uses `go run`, the build phase and every phase
  that executes or copies its output, contract runtime helpers,
  current-subject test helpers, deterministic kit-shaped fixtures, stripped
  subjects whose target-selected build closure matches, prospective
  preparation, preflight unchanged-subject helpers, and nested canary gates.
  A selected path always travels with its artifact record and authenticated
  authority; a bare ambient path is not trusted.
- **Real-build proofs are closed identities, not bypasses.** Changed-source,
  alternate-package, planted-byte, executable-digest, build-authorship,
  prospective-execution, compiler/linker-failure, release-target, and
  reproducibility assertions still request real artifacts. Target and source
  differences are ordinary identities. The registry owns the finite
  independent-authorship slots — for example reproducibility `first` and
  `second` — rather than accepting a caller-supplied nonce. A second request
  for the same slot reuses or returns the same result. Every proof uses the
  coordinator and transferable capacity-one turn.
- **Test binaries are classified separately.** A compiled `go test -c`
  package binary is not the Bench CLI. Canary continues to compile one test
  binary per exact target/package identity and reuse it for that package's
  baselines and bites, but those compiles join the same transferable turn.
- **One structural executable chokepoint.** Only the coordinator's private
  backend constructs an `exec.Cmd` that can author a registered executable.
  The executable registry is code consumed by the backend, caller adapters,
  and conformance. Static conformance follows typed calls/import ownership and
  parsed argv construction, not substring matching: it rejects an assembled
  `go run ./cmd/bench`, raw canonical `go build`, subject-mode
  `scripts/go-build.sh`, or verifier build outside the backend, and it
  rejects a registry consumer with no adapter. Shell and Go fixture sources
  that intentionally contain broken build text are a separate closed fixture
  class, never executable owners.
- **Story 1 prepares before scheduling.** The public entry resolves or verifies
  the CLI and verifier artifact set before constructing the phase schedule.
  The build phase remains present and owns root-local materialization and seal
  publication; phases that read `dist/bench` retain their build edge.
  Toolchain phases execute the already selected CLI artifact directly and need
  no new build edge. Story 1 can therefore land green with the current
  scheduler before story 2 replaces scheduling.
- **A resource turn is a transferable lineage, not a process count.** "Active"
  means the one actor currently entitled to perform Bench-created resource
  work. A composite parent transfers its turn through the authenticated child
  descriptor immediately before a synchronous spawn and becomes quiescent
  while waiting; after the child is reaped, the turn returns exactly once.
  The parent may remain a live process and logical phase, but it performs no
  resource work while the descendant owns the turn. A caller never holds a
  turn while asking a descendant to acquire another one. Missing, duplicated,
  or out-of-lineage transfer state is red.
- **Every resource-producing seam participates.** The registered classes are
  phase execution, stripped-subject materialization and cleanup, Go
  build/test/vet/run/compile, package and fixture materialization, artifact
  publication/sealing, canary compile/baseline/bite, and compiler-observing
  coordination proofs. Non-resource orchestration may prepare metadata, but
  any synchronous child doing a registered class receives the current turn.
  The same registry drives the process-tree active recorder, so an omitted
  class fails enrollment rather than disappearing from the metric.
- **The gate schedule has capacity one.** The ordinary outer table launches
  one ready sibling phase lineage at a time in stable topological order. The
  primary and stripped tables do not get independent schedulers; they settle
  through one ordered schedule and one aggregate verdict. Inner phases execute
  only as descendants of their waiting outer canary lineage. Dependency skips,
  result order, phase summaries, red aggregation, and cancellation reporting
  remain unchanged.
- **The core test phase settles per package.** The core package enumeration
  remains the existing ordered `CoreTestPackages` answer, but the phase invokes
  one package at a time at width one, continues after package reds, aggregates
  the same final red, and publishes one authenticated settlement after each
  package. One package is therefore the largest no-progress unit; a monolithic
  multi-package child cannot hide more than 45 minutes of serial progress from
  the watchdog.
- **Every Go child is width one.** Gate-owned `go test`, `go vet`, `go run`,
  `go build`, and `go test -c` calls receive the applicable `-p=1` and
  `-parallel=1` flags plus exactly one authoritative `GOMAXPROCS=1`. Ambient
  `GOMAXPROCS`, `GOFLAGS`, or test flags cannot widen the run. The build
  backend uses the same width, so a single authoring turn cannot create an
  internal compiler storm.
- **Canary has no worker pool.** Package compiles, baselines, and fixture bites
  run one item at a time in deterministic order. The existing stage ordering
  remains: a fixture cannot run before its package binary and baseline exist.
  `fixtureWorkers`, `bounds.CanaryInnerWidth`, and their width arithmetic are
  retired rather than left as dead policy.
- **Concurrency-positive owners are explicitly replaced.**
  `TestRunnerRunsPhasesConcurrently`,
  `TestSchedulerOverlapsIndependents`,
  `TestSweepRunsFixturesConcurrently`,
  `TestSweepBoundsFixtureConcurrencyAtDerivedWorkerBound`,
  `TestFixtureWorkers`, and
  `TestSweepPinsSingleGOMAXPROCSInInnerEnv` retire with the worker/overlap
  contract and are replaced by capacity-one order, turn-transfer, and exact
  width-one assertions. The gate-entry tests that pin
  `go run ./internal/freshness/check` are amended to require the selected
  verifier artifact and retain their refusal-before-CLI behavior.
- **Parallel tests lose resource authority.** All 196 current
  `internal/gate` `t.Parallel` calls retire. A closed registry names the few
  coordination tests allowed to create simultaneous lightweight callers; an
  AST conformance check rejects any other `t.Parallel` call in the package.
  Registered callers delegate the one resource turn to at most one child, so
  their lock, cancellation, ownership, and atomicity behavior remains real
  without overlapping resource work.
- **Timeout measures stalls, not serialized breadth.** The existing 45-minute
  gate timeout becomes a no-progress deadline, not a whole-run wall deadline.
  It resets only when a phase, canary item, or registered resource turn settles
  successfully or red. A run whose total serialized wall exceeds 45 minutes
  while continuing to settle work may remain green; one actor making no
  progress for 45 minutes reds with the existing timeout code and cancellation
  teardown. Arbitrary heartbeat output never resets the deadline. This prevents
  capacity one from converting legitimate breadth into a timeout
  green-semantics regression.
- **One repository admits one gate lineage.** The public gate lock is anchored
  at the repository common Git directory, so a second gate from another
  worktree of the same repository refuses before phase or artifact work. The
  zero-contention claim is per admitted top-level lineage; unrelated raw tools
  remain external, while two Bench gate lineages for one repository can never
  overlap.
- **No performance escape hatch.** There is no width flag, environment knob,
  automatic host-size formula, or best-effort fallback that can make an
  ordinary local gate concurrent. Changing capacity is a future
  reviewer-owned contract change, not runtime configuration.

### Ownership fences

- Artifact identity, store, backend, registry, verifier, and bootstrap owner:
  `internal/artifactstore/`, `internal/freshness/`,
  `scripts/go-build.sh`, `.bench/gate.sh`,
  `.bench/gate-prospective.sh`,
  `internal/conformance/gate_entry_test.go`, and
  `internal/canary/gate_entry_test.go`.
- Gate artifact-consumer, schedule, turn-transfer, timeout, and package-test
  owner: `internal/gate/` and `internal/bounds/bounds.go`.
- Canary serial-stage and retired-width owner: `internal/canary/` except
  `internal/canary/gate_entry_test.go`, plus
  `internal/conformance/bounds_policy_test.go`.
- Contract, preflight, and release consumer migration owner:
  `internal/contract/`, `internal/preflight/`,
  `scripts/build-artifacts.sh`, `scripts/release-preflight.sh`,
  `scripts/native-proof.sh`, and `scripts/build-offline-archives.sh`.

These are authoring fences, not horizontal tickets. `craft-tickets` derives
independently-green tracer outcomes and their value contracts after sign-off;
no two concurrent writers receive overlapping paths.

## Testing decisions

- Good artifact tests inject a counting builder and drive the real repository
  store through sibling processes, concurrent and repeated requests, failed
  and cancelled attempts, source movement, hostile records, and cross-target
  build-tag mutations. They assert artifact identity and backend call count
  rather than relying on elapsed time.
- Every artifact-store test injects a fresh isolated store root. Fault,
  corruption, repair, and cancellation cases are forbidden from discovering
  or mutating the developer repository's live common-Git-dir store.
- Good scheduler tests inject blocking phase/resource runners with a shared
  active counter. They drive the public schedule, split primary/stripped path,
  parent-to-child turn transfer, nested authority, canary stages, and
  cancellation. They distinguish a logically live parent from the one runnable
  resource actor and assert that the latter's maximum is exactly one.
- Existing freshness publication, hardlink detach, atomic replacement,
  build-attestation, prospective execution, release reproducibility, and
  compiler-failure tests remain the behavior owners for the real-build
  classes. Their setup moves through the coordinator without weakening their
  assertions.
- Gate-entry tests execute the selected freshness-verifier artifact under a
  hostile or stale CLI and preserve the current refusal-before-`gate-phases`
  assertion. A verifier-source mutation produces a new verifier identity and
  makes entry refuse the old artifact before the CLI. A counting verifier
  fixture proves nested canary entries reuse the same current checker identity
  without invoking the Go tool.
- A fake clock drives the no-progress timeout. Multiple settling operations
  whose total exceeds 45 minutes remain live, while one silent actor crossing
  45 minutes receives code 124 and the existing process-group teardown.
- The gate seam is the dev gate's build, test, conformance, contract, and canary
  phases. The implementation gate eventually proves the composed process tree,
  but this spec phase runs no gate; focused red-capable owners and static
  conformance are the implementation tickets' first evidence.
- Total wall time is reported after correctness but is not an acceptance
  threshold. The governing outcomes are exact build counts, maximum runnable
  resource actors, and bounded no-progress time, all independent of aggregate
  machine speed.

### Seam diagrams

    trigger: public gate, focused or sibling package process, or nested child
        │
        ▼
    source root + registered intent + target  ──▶  [ repository artifact store ]  ──▶  durable immutable artifact + identity
                                                         │
                                                         └── per-identity attempt lock ──▶ private builder
                  ◀ tests attach here: sibling processes, counting builder,
                    hostile records, target-tag mutation and cancellation

    trigger: resolved outer or inner gate phase table
        │
        ▼
    primary + stripped ready phases  ──▶  [ capacity-one lineage ]  ──▶  ordered results + one aggregate verdict
                                                   │
                                      parent delegates turn and waits
                                                   │
                                                   ▼
                                      child actor ──▶ nested child actor
                                                   │
                                                   ├── Go: -p=1, -parallel=1, GOMAXPROCS=1
                                                   └── canary: one compile/baseline/bite at a time
                  ◀ tests attach here: blocking runners record runnable actor,
                    turn ownership, launch order, teardown and progress

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| SB1 | 1 | Repeated, concurrent, and later requests for one exact reusable identity across separate processes invoke the builder once and return the same durable artifact and executable digest | repository store with injected counting builder, per-identity lock, and re-exec barrier | new test is red on the current output-path and process-local owners because two sibling processes cannot rendezvous and record more than one build | counts the expensive operation itself across process and invocation lifetimes |
| SB2 | 1 | A target-selected source, toolchain, target, build-mode, artifact-class, or registered proof-slot change creates one distinct identity; output path, root spelling, mtime, invocation, and process do not | target-aware identity resolver plus build-tag fixture matrix | new Windows/Unix and architecture-tag mutation tests are red against a host-derived digest because the non-host source change aliases | prevents stale cross-target reuse while retaining byte-identical cross-root reuse |
| SB3 | 1 | Failure reaches every waiter without automatic retry; cancellation publishes no success; corrupt, missing, symlinked, malformed, or digest-mismatched completed artifacts refuse; and exact-identity repair refuses active work and removes only the named bad record | attempt/result records, validator, and exact repair seam | new fault-injection tests are red if a second author starts, a partial is published, hostile state triggers transparent rebuilding, or repair crosses identity/active-owner scope | forbids retry storms and untrusted execution while keeping recovery bounded and explicit |
| SB4 | 1 | The gate, gate plumbing, phase readers, contract runtime, current-subject helpers, deterministic fixtures, stripped subject, prospective/preflight helpers, and nested canary all consume typed artifact records; only the private backend authors registered executables | typed executable chokepoint, code registry, Go AST call/import audit, and parsed argv audit | the new structural audit is red on today's assembled `GateGoArgv` `go run` argv and direct build adapters; a mutation rebuilding the argv from separate literals remains red | catches constructed commands that a substring scan misses and structurally owns the universal consumer set |
| SB5 | 1 | Changed-source, alternate-artifact, planted-byte, prospective, authorship, release-target, reproducibility, and compiler-failure proofs still observe real construction, but total builds equal their finite registered identities and slots | proof requests through the private backend | existing behavior assertions plus new registry-slot counts are red if a proof is reused incorrectly, invents a nonce, or bypasses the backend | preserves every real-build reason without an open-ended exception |
| SB6 | 1 | A focused package, its re-exec child, and sibling consumer packages from raw `go test ./...` converge on the repository store and one artifact for the same identity | shared consumer-package adapter, common-Git-dir discovery, and sibling re-exec harness | new multi-package subprocess fixture is red before implementation because each current package process owns independent state | closes the process-local `sync.Once` hole for the ordinary developer loop Fable identified |
| SB7 | 1 | Preparation authors one freshness-verifier artifact per identity; direct and nested entries reuse it, verifier-source movement changes identity and refuses the old verifier before the CLI, and no entry runs `go run ./internal/freshness/check` | gate-entry verifier selection, verifier-source mutation, and counting nested-entry fixture | the current command-pin tests, source-mutation case, and build counter are red because today's entry invokes `go run` on every bite and has no durable verifier identity to invalidate | enrolls and currencies the bootstrap compiler route without trusting the CLI it grades |
| SB8 | 1 | Direct-exec, immutable-link, and mutable-copy consumers follow the registry policy; mutating one root can never change the store artifact or another root | materializer plus existing detach/atomic-replacement seams | new policy matrix and existing inode/digest witnesses are red if a mutable consumer receives a shared link, detach is omitted, or a writable multiply-linked inode is accepted | resolves the prior safe-link contradiction with one executable policy source |
| SB9 | 1 | Story 1 alone resolves/verifies artifacts before scheduling, retains the build phase for publication, leaves `dist/bench` reader edges intact, and lets toolchain phases execute the store CLI without a new build edge | artifact preparation to existing phase-table construction | new phase inspection is red on either a missing build phase, a toolchain `go run`, or an added toolchain-to-build edge | proves the artifact story is independently green rather than relying on story 2's serial scheduler |
| ZC1 | 2 | The ordinary outer gate and split primary/stripped path launch one sibling phase lineage at a time; a waiting logical parent may remain live while exactly one descendant owns the runnable resource turn | public scheduler and parent/child turn recorder with blocking fake phases | new sibling-overlap and turn-owner assertions are red against both current concurrent schedulers, but do not count a quiescent parent as a second actor | defines the quantifier so nesting is neither an impossible assertion nor a hidden overlap |
| ZC2 | 2 | Every gate-owned Go build/test/vet/run/compile child is pinned to one package/test slot and one runtime thread, with ambient width inputs unable to widen it | command constructors and closed child environment | new argv/env matrix is red on today's unpinned outer commands, and hostile ambient cases red if duplicate or wider values survive | prevents hidden parallelism inside the one visible phase |
| ZC3 | 2 | Canary compiles, baselines, and fixture bites one at a time in deterministic stage order, with worker/inner-width policy absent | canary item iterator, bounds registry, and blocking fake runner | replacement serial tests are red against `eachIndex`/`fixtureWorkers`; the conformance registry is red if retired `CanaryInnerWidth` remains advertised | removes the second fan-out site and its stale policy rather than weakening a green overlap test silently |
| ZC4 | 2 | A synchronous parent transfers its one turn to a child, performs no resource work while waiting, receives it once after reap, and refuses missing, forged, duplicated, closed, or malformed lineage state | authenticated turn-transfer descriptor and re-exec parent/child harness | hostile transfer tests are red if the child acquires a second turn, the parent continues work, or fallback creates a local pool | directly prevents the nested deadlock and double-owner cases |
| ZC5 | 2 | Phase execution, stripped-subject materialization/cleanup, every Go child, package/fixture materialization, publication/sealing, canary items, and compiler-observing proofs all use the registered turn; the process-tree maximum runnable actor is one | resource-class registry, typed launch adapters, and one process-tree recorder | deleting acquisition at each class in turn makes the active counter reach two or makes the enrollment audit red | enumerates the quantified resource set and catches work omitted from the recorder |
| ZC6 | 2 | Interrupting an owner, delegated child, or waiter reaps every descendant, returns the turn exactly once, emits no false verdict, and lets a later run proceed | cancellation and turn-transfer teardown seam | fault tests hang, duplicate/lose the turn, publish a verdict, or block the next run when reclamation is wrong | capacity one cannot trade contention for deadlock or false completion |
| ZC7 | 2 | An ordinary gate exposes no concurrency knob, contains no unregistered `internal/gate` `t.Parallel`, and stays capacity one under high core counts and hostile `GOMAXPROCS`, `GOFLAGS`, or test flags | public entry, closed coordination-test registry, AST audit, and closed environment | static usage/AST tests plus hostile-env scheduler run are red if any supported widening route or unregistered parallel test exists | makes “every time” a structural contract rather than a recommended invocation |
| ZC8 | 2 | More than 45 minutes of serialized work remains eligible for green while operations settle inside each window; core tests run and settle one enumerated package at a time while continuing after reds; arbitrary output cannot reset the timer; one package/actor stalled for 45 minutes reds with code 124 and full teardown | real core package enumerator plus fake-clock progress watchdog driven only by package/phase/item/turn settlement | new multi-package cumulative, continue-after-red, package-stall, and heartbeat-spam tests are red against today's monolithic child and one-shot whole-run deadline | makes the real longest phase observable to the watchdog without changing its package membership or aggregate verdict |
| ZC9 | 2 | Two gates from worktrees sharing one common Git directory cannot overlap: the second refuses before preparation, phase launch, or store authorship | repository-common gate lock and two-worktree process harness | new cross-worktree overlap test is red if the second lineage reaches an artifact or phase marker | closes same-repository sibling runs rather than hiding them behind the external-process exclusion |

Degenerate-implementation check: story 1's cheapest wrong implementation keeps
the existing process-local `sync.Once` and passes a binary path to one more
helper; SB1, SB4, SB6, and SB7 stay red on sibling build count, structural
owners, and the verifier compile. A host-derived cache key stays red on SB2's
non-host mutation. Story 2's cheapest wrong implementation swaps the outer
scheduler to sequential but leaves the split scheduler, Go package width,
canary workers, nested turn handoff, stripped materialization, package-level
progress, or repository-wide gate admission unchanged; ZC1 through ZC9 isolate
each surviving layer.

Composition check: artifact reuse without capacity one still lets distinct
identities and package binaries contend; capacity one without artifact reuse
serializes the same repeated work and remains unacceptably expensive. SB5 and
ZC5 require the two stories to compose through one authority. SB9 separately
proves story 1's pre-schedule artifact contract stays green before the serial
scheduler lands.

### Edge inventory

- Duplicate/reordered/repeated same-key requests: SB1.
- Separate focused tests, sibling packages from one `go test ./...`, and a
  gate racing for the same identity: SB1/SB6; the repository attempt lock
  admits one author and all others read the completed record.
- Source changes during identity resolution or construction: SB2/SB3; the
  completed artifact is re-verified against the accepted generation before it
  is published, and movement refuses rather than relabeling bytes.
- Target-conditional source selected only by Windows, Unix, architecture, CGO,
  or build tags: SB2; enumeration and the post-build recheck both run under the
  requested target.
- Foreign, empty, duplicated, malformed, or forged artifact/authority
  environment: SB3, SB4, ZC4.
- Artifact-store fault injection and repair: SB3; every test uses an injected
  isolated store root and never discovers the live repository store.
- Output and source paths containing spaces, newlines, leading dashes, or
  symlinks; special files at artifact, authority, or manifest paths: SB2/SB3.
- Read-only filesystem, hardlink refusal, cross-device materialization, and an
  existing destination: SB3/SB8 plus the existing copy/detach/replace
  witnesses.
- Concurrent first request, producer failure, waiter cancellation, parent
  death, closed descriptor, and cleanup after partial construction: SB1, SB3,
  ZC4, ZC6.
- Zero phases, one phase, unsatisfied needs, a red producer, optional absent
  tools, and primary/stripped tables of unequal size: ZC1.
- A live outer canary parent waiting on an inner phase: ZC1/ZC4; the parent is
  logically live but quiescent, and the child owns the transferred turn.
- One-core and many-core hosts; hostile width environment and flags: ZC2/ZC7.
- Empty canary selection, one item, multiple packages, failed compile, empty
  baseline, and interrupted bite: ZC3/ZC6.
- Coordination tests that require simultaneous callers: ZC5; callers may
  overlap, but their resource jobs cannot.
- Cross-platform and cross-mode release builds: SB2/SB5; each target/mode is a
  distinct identity and runs serially.
- Missing/stale freshness verifier, stale CLI, nested gate entry, and verifier
  source movement: SB2/SB3/SB7; no case falls back to `go run`.
- Total wall beyond 45 minutes with continued progress versus one 45-minute
  stall, including per-package progress inside the core test phase: ZC8.
- Two gate entries from distinct worktrees of one repository: ZC9; the
  common-Git-dir owner admits one lineage and refuses the other before work.
- Re-run after an unchanged green: the gate verdict may still be reused under
  its existing subject contract; when work runs, this spec's capacity and
  artifact rules apply.
- **Won't handle:** unrelated user, IDE, CI, or operating-system processes —
  Bench cannot serialize work it did not create or authorize, and this spec
  makes no whole-machine idle claim. A second Bench gate for the same
  repository is not in this exclusion; ZC9 refuses it.
- **Won't handle:** automatically retrying or deleting a failed/corrupt
  unchanged identity — transparent repair can recreate the repeated work and
  hide store corruption, so the explicit repair operation owns that action.
- **Won't handle:** removing lightweight concurrency from tests whose behavior
  is lock exclusion, cancellation, process ownership, or atomicity — serializing
  their resource work preserves the required overlap without weakening those
  claims.

## Out of scope

- A host-wide broker coordinating unrelated Bench repositories or external
  tools: approximately 8–12 edits and 3 gate runs; it requires a security and
  stale-owner policy beyond this repository-owned contract.
- Automatic garbage collection or size-based eviction of obsolete successful
  artifact identities: approximately 6–10 edits and 3 gate runs; content
  addressing makes old generations inert, while deletion/retention policy is a
  separate storage-management capability.
- Reducing, skipping, diff-scoping, or memoizing any gate check: zero checks are
  removed and green semantics remain unchanged.
- Parallel release production across platform targets: targets are distinct
  artifact identities and this spec deliberately runs them one at a time when
  launched by Bench.
