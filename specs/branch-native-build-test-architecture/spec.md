# Branch-native build and test architecture

Status: implemented

Decision source: reviewer-requested clean-slate design in the current conversation on 2026-08-09

## Problem

Bench's ordinary test surface is implemented as an integration harness. Tests that mean
to prove command decisions, gate scheduling, lifecycle transitions, conformance checks,
and preflight policy repeatedly create temporary directories, initialize and commit Git
repositories, copy kit trees, launch shell wrappers, execute Bench subprocesses, and in
some cases launch nested `go test` processes. The current source contains 232
`CommitAll` call sites and 133 direct `exec.Command` call sites across `internal/gate`,
`internal/contract`, and `internal/preflight` tests. Earlier instrumented evidence found
one mapping row launching 415 Git processes and materializing the same fixture tree 47
times.

That work is not Go compilation. A warm sealed Bench build completes in less than a
second on the current host while individual fixture-heavy packages take roughly one to
three minutes. Serializing outer phases bounds contention but makes this accumulated
fixture cost visible end to end; adding concurrency, a shared fixture pool, hardlinks,
or another cache would hide rather than remove the wrong test seam.

The suite also repeats authority. The dev gate drives core tests, contract tests,
conformance entry tests, a separate conformance-suite run, a stripped-subject schedule,
and mutation canaries that can cross another process or test seam. A behavior can
therefore be compiled, materialized, and asserted several times without gaining a new
observable. This makes the gate slow, makes failures expensive to reproduce, and makes
test setup rather than production interfaces the dominant architecture.

## Solution

Replace the fixture-driven ordinary suite with a branch-native, three-layer test
architecture.

Ordinary Go tests exercise production decisions in-process through the highest interface
where their failures remain observable. The current checkout is immutable source input;
tests construct domain values such as accepted snapshots, manifests, plans, lifecycle
states, and command arguments rather than constructing repositories to manufacture those
values. CLI contracts call their production `Command` functions directly. Git,
filesystem, shell, and process-group behavior remain behind their existing narrow
adapters and receive small adapter tests.

One tagged system-test package owns the irreducible end-to-end journeys. It consumes the
top-level run's one selected Bench executable and may materialize at most three disposable
Git repositories across its complete run. It alone proves wrapper routing, freshness,
installation, serialized process state, signals, teardown, and one stripped-distribution
smoke journey. It does not rerun the ordinary package universe against each repository.

The dev gate launches one ordinary `go test -count=1 ./...` driver and lets the Go tool
own package scheduling. Conformance and contract assertions are ordinary Go tests, not
additional phase-table package runs. Mutation canaries invoke their owning check directly
against mutated inputs; one system journey proves top-level canary dispatch. The full
release, cross-target, reproducibility, and artifact workflow remains the explicit ship
tier rather than entering the ordinary branch gate.

This is an atomic replacement, not a compatibility migration. The implementation branch
may use provisional commits, but it promotes only with the old ordinary fixture framework,
duplicated phases, and nested test constructors deleted. It introduces no fixture pool,
old-to-new adapter, hardlink template cache, or durable dual suite.

## User stories

1. As a contributor iterating on Bench code, I want ordinary tests to execute the current
   branch's production interfaces in-process through one Go test driver, so a normal red
   or green run spends its time on behavior rather than rebuilding temporary repositories.
   Line: gpt-5.6-sol / high. This changes the kit's test seams across command, gate,
   contract, and preflight ownership and therefore receives the high-leverage kit line.
2. As a contributor changing shell, Git, installation, freshness, or process-lifecycle
   behavior, I want a small real-binary system suite with explicit materialization and
   process budgets, so irreducible integration behavior stays proven without turning
   every assertion into an end-to-end journey. Line: gpt-5.6-sol / high. The executable,
   repository, and teardown contracts are authority-bearing and require exact integration
   evidence.
3. As a gate maintainer, I want conformance, contract, stripped-subject, and canary proofs
   attached to the narrowest production interface that can still make their defect red,
   so the gate never reruns a package universe or inner gate merely to reach one check.
   Line: gpt-5.6-sol / high. The change must preserve bite while deleting duplicated
   oracle paths.
4. As a reviewer, I want structural enforcement of the new architecture and deletion of
   the old fixture machinery, so a later test cannot silently reintroduce repository
   materialization, nested Go tests, or process-heavy CLI helpers into the ordinary suite.
   Line: gpt-5.6-sol / high. A weak census would allow the cheaper old architecture to
   accrete again one helper at a time.

## Implementation decisions

- Preserve the single-build run owner, `BENCH_RUN_BINARY` propagation, source freshness
  seal, process-group teardown, and `.logs/` major-point records. They are correctness and
  observability contracts, not the source of the fixture cost.
- The ordinary suite is the packages selected by one `go test -count=1 ./...` process.
  The gate does not enumerate packages into separate `go test` commands and no test it
  launches may invoke `go test`, `go run`, or a second Bench build. Go owns ordinary
  package concurrency; outer gate phases remain serial so only one toolchain driver is
  active at a time.
- Production decisions expose deep interfaces before side effects. Gate evaluation and
  scheduling consume accepted snapshot and manifest values; link/setup/upgrade consume
  current-tree and kit inventories and return an operation plan plus report; preflight
  consumes snapshot and evidence values and returns a decision; lifecycle commands
  consume serialized state plus an event and return next state plus actions. Existing
  interfaces are preferred where they already carry these values.
- Do not introduce a general fake Git or fake filesystem. Production adapters load the
  current branch into immutable domain values and apply returned operations. Ordinary
  tests construct those values directly. Adapter tests alone exercise the operating
  system or Git representation.
- Command behavior tests call the production command function with argument slices and
  in-memory stdin/stdout/stderr. A command receives its repository facts or operation
  executor through the same interface production uses. Calling `Fixture.Bench`,
  `Fixture.BenchWrapper`, or another executable helper is not an ordinary command-test
  seam.
- Adapter exceptions are an exact inventory. `internal/git` may create at most one Git
  repository across its complete package-test run. `internal/gate` may launch at most one
  controlled process group to grade descendant signaling and teardown, but may create no
  repository. Shell-wrapper attachment belongs to `internal/systemtest`. Every other
  ordinary package creates zero Git repositories and launches no operating-system process
  from a test. No adapter test invokes a full gate or Go package suite.
- Add one `internal/systemtest` package guarded by a `system` build tag. Its `TestMain`
  receives and validates the inherited selected executable, owns every disposable
  repository, records materialization/process counts, and removes all descendants and
  repositories before returning. Individual tests cannot construct repositories or
  select/build executables themselves.
- The complete system suite has a hard budget of at most three Git repository
  materializations. The intended journeys are: one linked/install repository, one
  lifecycle/state-reload repository, and one stripped-distribution repository. Journeys
  may share an owner but not hidden mutable state; a journey must reset or receive a fresh
  budgeted repository before it starts.
- The stripped proof materializes one stripped subject and runs one real-binary smoke
  journey proving package shape, wrapper resolution, selected-binary use, and refusal of
  an excluded-path dependency. It does not rerun contract, conformance, or the ordinary
  package universe against the stripped tree.
- Canary fixtures become data for their owning production check. A mutation invokes only
  that check, requires its specific diagnostic, restores the input, and requires green.
  Registry tests enforce one owner per fixture and complete fixture coverage. One system
  journey proves `bench canary` dispatch and aggregation; no canary invokes an inner gate
  or nested `go test`.
- Conformance checks that inspect source shape remain direct check functions over one
  immutable generation-scoped snapshot. Their ordinary tests call those functions.
  Conformance is not separately re-entered through `TestRootConformance` after the
  ordinary package driver already compiled and executed the package.
- Race coverage remains one targeted `go test -race` invocation over the authoritative
  concurrency-test registry. It runs after the ordinary driver and does not include the
  system package or fixture-heavy journeys.
- Release preflight, cross-target builds, artifact reproducibility, and publication
  rehearsal remain `bench prep-release` ship-tier capabilities. Their branch tests grade
  planners, validation, and error reporting in-process; only the explicit ship command
  executes the complete release workflow.
- The six decision domains replaced by branch-native interfaces are exact: gate
  acceptance/scheduling/evidence, link/setup/upgrade/unlink planning, preflight
  decision/evidence, spec-build lifecycle transitions, canary selection/aggregation, and
  freshness selection/refusal. The structural test architecture census is a default
  conformance check backed by Go
  syntax/assembled argv, not prose or a raw timing threshold. It rejects ordinary
  `git init`, `git clone`, `git worktree`, `exec.Command` process entry, nested Go tool
  execution, subject builders, duplicated package phases, and repository constructors
  outside the named adapters and `internal/systemtest` owner.
- Structural budgets are exact: one ordinary Go test driver, one selected Bench build,
  zero Git repositories or operating-system processes in decision/command tests, at most
  one ordinary adapter repository in `internal/git`, at most one ordinary controlled
  process group in `internal/gate`, zero nested Go test/run commands, zero inner canary
  gates, at most three system repositories, at most four repositories across the complete
  dev gate, and exactly one stripped system journey. `.logs/` reports elapsed time and
  counts for diagnosis, but elapsed time is not itself the oracle.
- Before deleting an existing test, the implementation records the user-visible behavior
  it proves and maps that behavior to an in-process, adapter, system, ship-tier, or
  intentional-deletion disposition. This behavior ledger exists only for review of the
  replacement branch; it does not become a second permanent registry.
- Delete the broad ordinary fixture framework after replacement coverage exists. In
  particular, no general `Fixture.Bench*`, `CommitAll`, copied-kit constructor, per-test
  repository constructor, or ordinary contract runner remains available for new tests.
  Specialized system and adapter helpers stay private to their owning package.
- This is a wide single-writer build. Ticket slicing may use ownership fences below, but
  no ticket is independently promotable as a shipped compatibility stage; promotion is
  the one composition in which the replacement suite is green and the retired framework
  is absent.

## Testing decisions

- The six named decision domains receive table-driven tests through their public package
  interfaces. Every domain varies accepted values, errors, empty states, hostile values,
  and reruns without constructing a repository or launching a command.
- The production top-level dispatch registry is the sole command enumeration. Every
  registry member receives direct command tests for usage, exit status, stdout, stderr,
  and operation-plan/result translation. Its process-attachment disposition is recorded
  per member and constrained to these exact sets:
  - Direct only: `version`, `worktree`, `resume-clean`, `session-inspect`, `shift`,
    `commit`, `spec`, `gate-go`, `guard-git`, and `check-agent-line`.
  - System attachment: `setup`, `link`, `init`, `doctor`, `unlink`, `upgrade`,
    `worktree-hook`, `gate`, `gate-run`, `gate-pin`, `gate-phases`, `freshness-check`,
    `freshness-publish`, `canary`, and `stop-verdict`.
  - Ship attachment: `release-preflight`, `prep-release`, and `release`.
  Adding, removing, or renaming a top-level entry without assigning its disposition makes
  the registry test red. The direct-only set is about process attachment; its Git and
  lifecycle decisions still receive direct branch-native tests.
- Adapter tests count their own repository and process constructors. A mutation that adds
  a second constructor to one adapter test or calls a full gate makes its package red.
- `internal/systemtest` records selected executable path, inode, and digest at every
  journey and requires one identical value. It records each repository construction and
  process start, fails above budget, and proves cleanup after green, red, interrupt, and
  timeout outcomes.
- The ordinary package driver is captured as assembled argv. Its test requires exactly
  `go test -count=1 ./...` plus only the gate's established cache/environment flags. A
  package loop, a second ordinary test phase, or an argv that names contract/conformance
  packages separately makes the check red.
- The tagged system driver is captured as assembled argv and requires exactly
  `go test -count=1 -tags=system ./internal/systemtest` plus only established
  cache/environment flags. The selected executable environment is mandatory, and no
  other phase may name the package or `system` tag.
- The race driver is captured as assembled argv and requires one
  `go test -race -count=1` invocation over the authoritative `internal/racetests`
  registry. The registry owns every package/test selector; the system package and its
  journeys are forbidden from that argv.
- The structural census parses Go and shell entry sites and owns the exact exception
  inventory. Adding an `exec.Command`, Git repository constructor, nested Go invocation,
  subject builder, system repository, or stripped journey requires an explicit census
  change and reviewer decision.
- Canary mutation tests call the registered owning check directly. A no-op mutation,
  always-green check, generic infrastructure diagnostic, missing owner, duplicate owner,
  inner gate, or nested Go test makes the default suite red.
- The stripped journey tests the shipped shape once with the selected executable. A
  wrapper that reaches the source checkout, an excluded path read, a copied second
  executable, or a second stripped materialization makes it red.
- A composition test drives the actual selected executable through the system owner and
  proves that the same command implementation covered in-process is reached. This catches
  a replacement whose pure command tests and wrapper tests are separately green while the
  installed route dispatches elsewhere.
- A deletion test asserts that the retired fixture APIs and duplicated phase constructors
  are absent. Renaming them or wrapping them behind a new helper still fails the syntax
  census because the forbidden effects, not symbol names alone, are counted.

### Seam diagram

    current branch files + Git facts
                 │
                 ▼
       [ production snapshot adapters ]
                 │ immutable domain values
                 ▼
       [ decision / plan / transition APIs ] ◀── ordinary Go tests
                 │
                 ▼
       [ side-effect adapters + Command ]     ◀── narrow adapter and direct command tests
                 │
                 ▼
    one selected Bench executable
                 │
                 ▼
       [ internal/systemtest owner ]          ◀── at most three disposable repositories
                 │
                 ▼
       wrapper / freshness / process / stripped observations

    dev gate
       ├── format + vet
       ├── one ordinary go test driver
       ├── one targeted race driver
       ├── one selected Bench build
       ├── one system-test package
       └── shellcheck

    ship gate
       └── release, cross-target, reproducibility, and publication workflows

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| BN1 | 1 | Ordinary command behavior is tested through production `Command` interfaces without launching Bench or a wrapper. | direct command interface | observed red: `! rg -n 'func \(f Fixture\) Bench' internal/contract/command.go` exited 1 and named eight process-backed `Bench*` helpers | A fixture helper that still executes the CLI cannot satisfy an in-process command seam merely by hiding behind a new name. |
| BN2 | 1 | Tests for the six named decision domains and every direct command surface materialize zero Git repositories, launch zero operating-system processes, and commit zero fixture trees. | branch-native decision/command structural census | observed red: the zero-budget assertion over `.CommitAll(` exited 1 with `fixture_commit_sites=232` | A repository-backed command or decision test needs commits to manufacture state; counting the constructor effect catches both repeated repos and a shared helper call without contradicting the two named adapter exceptions. |
| BN3 | 1 | The dev gate launches exactly one ordinary `go test -count=1 ./...` driver and does not separately invoke contract or conformance packages. | assembled gate phase argv | observed red: `! rg -n -e conformancePhaseName -e 'canary\.PhaseContract' -e PhaseConformanceSuite internal/gate/phases.go` exited 1 and named all three extra phase owners | A renamed phase cannot pass because assembled argv must contain one package-universe invocation and no separately named ordinary package run. |
| BN4 | 1 | Go owns ordinary package scheduling inside its one driver; the gate has no per-package serial loop or nested test command. | ordinary-driver execution seam plus structural census | observed red: `! rg -n 'exec\.Command\("go", "test"' internal/contract internal/canary --glob '*_test.go'` exited 1 and named nested test constructors | A manual loop or nested Go driver creates an additional constructor the census rejects even when it runs serially. |
| BN5 | 1 | Gate acceptance/scheduling/evidence, link/setup/upgrade/unlink planning, preflight decision/evidence, spec-build lifecycle transitions, canary selection/aggregation, and freshness selection/refusal consume immutable domain values and express success, refusal, and planned effects without repository or process setup. | six named decision interfaces | not TDD-able as one row before interface selection: current decisions are distributed across command and fixture setup; BN1 and BN2 are the observed composition red | Naming all six domains makes omission checkable and prevents a pass-through abstraction that still forces any listed test around the module. |
| SY1 | 2 | One tagged `internal/systemtest` owner validates and propagates the exact inherited selected executable to every real-binary journey. | system `TestMain` and selected-binary process seam | observed red: `test -d internal/systemtest` exited 1 | A directory/owner absence proves there is no singular system surface; identity records make a private build or alternate path observable once it exists. |
| SY2 | 2 | The complete system suite materializes at most three repositories, all through its owner, and cleans repositories plus descendants on every terminal outcome. | system materialization/process owner | not TDD-able before SY1: no singular system owner exists to count across journeys; the current ordinary fixture-commit count is 232 | One global counter catches per-test constructors that individually look bounded, while terminal markers prove teardown rather than only allocation count. |
| SY3 | 2 | Wrapper routing, installation, freshness refusal, serialized state reload, interrupt, timeout, and process-group teardown are each proven once through the real selected executable. | named system journeys | currently distributed across process-backed contract fixtures; SY1's missing owner is the red signal until the named journeys exist | Enumerating the process-only classes prevents a cheap system suite that tests only `version` while claiming end-to-end coverage. |
| SY4 | 2 | `internal/git` is the sole ordinary repository exception and constructs at most one repository per package run; `internal/gate` is the sole ordinary process exception and launches at most one controlled process group; both invoke zero full gates or Go package suites. | exact ordinary adapter exception inventory | observed red: the current fixture framework supplies repository and process constructors across gate, contract, and preflight packages; BN2 records the constructor red | Naming both packages, effects, and limits distinguishes necessary representation tests from a broadly exempt `integration` label. |
| SY5 | 2 | The dev gate invokes the tagged system owner exactly once as `go test -count=1 -tags=system ./internal/systemtest`, with the inherited selected executable and no second phase naming the package or tag. | assembled system phase argv | observed red: `! rg -n -e internal/systemtest -e tags=system internal/gate .bench/gate.sh` exits 0 because no system driver is wired | A package that exists but is never called would otherwise satisfy owner and journey tests while providing no gate evidence. |
| NC1 | 3 | Each canary mutation invokes only its registered owning check and requires a mutation-specific red followed by restored green. | canary registry-to-check interface | observed red: the nested-Go-test assertion named publication and release-evidence canary constructors | A direct owning-check call fails if the check is always green, ownership is wrong, or the mutation only triggers generic infrastructure noise. |
| NC2 | 3 | One system journey proves top-level canary dispatch and aggregation; no canary launches an inner gate or nested Go test. | canary command system seam plus structural census | current canaries retain process and nested-test paths; the NC1 source assertion exits 1 | The dispatch journey preserves the shipped surface while the zero-constructor census prevents recreating the whole gate behind a focused selector. |
| ST1 | 3 | One stripped subject is materialized and receives one selected-binary distribution smoke journey, with no stripped contract/conformance package rerun. | stripped system journey and phase argv | observed red: `! rg -n -e contractSubtree -e stripped internal/gate/phases.go internal/gate/stripped_worktree.go` exited 1 and showed the contract subtree entering the split schedule | A smoke that silently retains the old stripped package schedule still exposes the forbidden argv and second phase-table construction. |
| ST2 | 3 | The stripped journey proves package shape, wrapper resolution, selected-binary identity, and refusal of excluded-path dependency. | real stripped subject process seam | current full stripped scheduling covers these indirectly; cannot start as a focused red until SY1 owns the journey | Enumerated observations prevent the degenerate fix of deleting stripped verification to meet the one-materialization budget. |
| RG1 | 4 | A default syntax/argv census enforces one build, one ordinary test driver, zero decision/command repositories and processes, the one-repository `internal/git` exception, the one-process-group `internal/gate` exception, zero nested Go drivers, zero inner canary gates, at most three system repositories, at most four dev-gate repositories total, and one stripped journey. | default conformance structural check | not TDD-able before its exact registries and owners exist; BN2, BN3, BN4, NC1, ST1, and SY4 are the observed reds it must contract | One exact census makes every forbidden constructor and every exception reviewer-visible instead of relying on a timing symptom. |
| RG2 | 4 | The retired general fixture APIs, copied-kit constructors, per-test repository constructors, and duplicated phase constructors are absent at promotion. | deletion test plus syntax census | observed red: BN1-BN4 and ST1 source assertions all exit 1 on the current tree | Deletion prevents an atomic replacement from becoming a permanent second framework layered beside the first. |
| RG3 | 4 | Release, cross-target, reproducibility, and publication workflows remain ship-tier only while their decisions retain ordinary in-process coverage. | dev-phase argv and ship command interfaces | existing tier separation covers the workflow execution side; the new branch tests cannot start until the decision interfaces in BN5 exist | Checking both phase absence and direct decision coverage prevents either pulling release work into the dev gate or deleting its branch-level validation. |
| RC1 | 1, 4 | The dev gate retains exactly one targeted race driver whose argv is `go test -race -count=1` plus selectors derived only from `internal/racetests`; it excludes `internal/systemtest`. | authoritative race registry and assembled phase argv | existing control: `internal/gate/gate_go_test.go` `TestGateGoRaceRequiresTheTestToRun` proves every registry member runs, but the redesigned phase still needs the exact one-driver/exclusion assertion | The execution sentinel catches an argv that skips a registered test, while the new argv assertion catches duplicate race drivers or accidental system-suite instrumentation. |
| CP1 | 1, 2, 3 | The actual selected executable reaches the same command implementation covered in-process, while one system owner accounts for every repository and process it starts. | command-to-system composition seam | not TDD-able before BN5 and SY1 exist; current separately green fixture surfaces can still fail together, as the selected-binary migration exposed | The real producer catches the composition degenerate where pure command tests use one path and wrapper/system tests dispatch another. |

### Degenerate implementations the map rejects

- Rename `Fixture.Bench` and `CommitAll` while retaining their process and Git effects.
  BN1, BN2, and RG1 count the effects and remain red.
- Add a shared fixture pool or hardlinked repository template while leaving ordinary tests
  repository-shaped. BN2's zero ordinary materialization budget remains red.
- Add `internal/systemtest` but allow every subtest to create its own repository. SY2's
  suite-wide owner and counter fail above three.
- Launch one outer `go test ./...` whose tests launch nested `go test` commands. BN4 and
  RG1 reject the nested constructors.
- Replace nested canary gates with nested owning-package `go test` commands. NC1 and NC2
  require direct check calls and zero nested Go drivers.
- Remove the stripped schedule without proving installed package shape and excluded-path
  refusal. ST2 remains uncovered/red.
- Make in-process command tests and wrapper tests independently green while the wrapper
  dispatches a different implementation. CP1 drives the real selected executable and
  compares the reached command observation.

### Edge inventory

- Error path — decision interfaces represent malformed state, adapter failure, missing
  tool, and refused operation without requiring a repository; adapter and system rows
  retain the actual Git/process error translation.
- Empty or absent input — empty package universe, absent Git metadata, absent manifest,
  no system journeys, and no canary fixtures have explicit green/refusal behavior and do
  not manufacture placeholder repositories.
- Limit values — zero and three system repositories are accepted, a fourth is red;
  zero and one stripped journeys are distinguished; an empty ordinary package set launches
  no second test driver.
- Malformed input — malformed serialized lifecycle state, manifests, seals, and command
  arguments are exercised through decision interfaces; malformed on-disk representation
  receives one adapter/system case where parsing alone cannot observe the failure.
- Interrupted or partial state — the system owner proves interrupt and timeout teardown,
  including a leader that exits before its descendant; no repository or selected binary
  is removed before descendants settle.
- Re-run idempotency — decision tests repeat planning over the same immutable input and
  system journeys repeat the operation where persistence or installation is the behavior
  under test.
- Process lifecycle — only named adapter/system tests create process groups; each records
  start, terminal result, descendant drain, and cleanup.
- Hostile environment — system journeys cover an absolute path with spaces/glob
  characters, stripped `PATH`, hostile ambient `BENCH_RUN_BINARY`, symlinked wrapper path,
  and missing optional tool without multiplying repositories per value.
- Dirty or foreign working state — ordinary tests never mutate the working branch;
  destructive lifecycle behavior runs only in an owner-created system repository and
  cannot stage, reset, or commit the contributor's checkout.
- Cold and warm Go caches — both must preserve identical verdicts and selected-binary
  identity. Timing is logged for diagnosis, never averaged into correctness.

### Won't handle

- Eliminating all temporary directories. Non-repository temporary output remains valid
  where the output file or atomic promotion is the behavior being tested; RG1 governs
  repository and process constructors, not `t.TempDir` by spelling alone.
- Running destructive lifecycle journeys directly against the contributor's branch. The
  branch is immutable ordinary input; the bounded system repositories are the deliberate
  safety exception.
- Moving the full release matrix into the dev gate. `bench prep-release` is already a
  separate ship capability and retains its own exact evidence contract.
- Enforcing a fixed wall-clock pass/fail threshold. Host filesystem and toolchain latency
  vary; exact construction and invocation budgets are the stable oracle, with `.logs/`
  retaining elapsed telemetry.
- Preserving compatibility for internal test helpers. They are not a shipped API and this
  replacement explicitly deletes rather than adapts them.

## Ownership fence

This atomic replacement has one writer. Ticket slices, if used for ordering, do not create
additional writers and cannot be promoted independently. The writer may edit only these
exact paths:

- `cmd/bench/**`
- `internal/**`
- `.bench/gate.sh`
- `.bench/gate-prospective.sh`
- `bin/bench.sh`
- `scripts/go-build.sh`
- `scripts/go-build.inputs`
- `tests/canary/**`
- `projects/benchkit.md`
- `CHANGELOG.md`
- `specs/branch-native-build-test-architecture/spec.md`

`ROADMAP.md`, `capture/learnings.md`, `capture/session-handoff.md`, every other spec, and
all paths outside the list remain foreign. The implementation must stop for authorization
if a required change falls outside this fence.
