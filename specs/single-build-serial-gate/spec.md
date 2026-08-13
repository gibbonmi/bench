# Single-build serial gate

Status: implemented

Decision source: reviewer-confirmed current conversation on 2026-08-08

## Problem

One logical Bench test or gate run can currently author several host Bench executables. The prospective shell builds one, the ordinary shell compiles a separate freshness verifier, the phase table owns another build, gate-go plumbing uses `go run ./cmd/bench`, and unchanged-host test helpers invoke the subject builder again. A nested canary gate can repeat that work. At the same time, the outer phase scheduler and the primary/stripped split overlap phase processes, multiplying otherwise independent Go compiler demand.

A worktree-lifetime binary is not an exact answer: a ticket's red run, green run, review probe, and promotion gate may grade different source snapshots. Reusing one durable binary across those runs can silently execute old code. The safe unit is therefore one top-level run, not one worktree lifetime.

## Solution

Every non-reused top-level Bench gate or focused `bench test` run owns one private temporary directory and selects one absolute, host-target Bench executable built there through `scripts/go-build.sh` in subject mode. The owner propagates that exact path to every ordinary consumer. Direct and prospective gates use the same owner; nested gates inherit the selected path and have no authoring fallback. The owner reaps the complete child process group and removes the private directory on every terminal outcome. A later run starts over with its own exact-snapshot binary, while Go's normal build cache still avoids recompiling unchanged packages.

Every phase-table invocation runs one phase process at a time. Primary and stripped phases enter one stable topological schedule rather than two overlapping schedules. Existing dependency order, dependent skips, optional-tool skips, red aggregation, and process-group cancellation remain authoritative.

Independent compiler proofs remain independent: cross-target, release, reproducibility, changed-source, alternate-package, and intentional compiler/linker checks still build the distinct artifact their assertion observes. They are enumerated in one census so an ordinary build cannot hide behind the exception policy.

## User stories

1. As a contributor running a ticket's red, green, review, or promotion check, I want that top-level run to build one exact Bench executable and make every ordinary child use it, so the run cannot multiply CLI builds or grade stale code. Line: gpt-5.6-terra / high. The behavior crosses gate ownership, shell entry, Go tests, and nested processes, so correctness needs the high-effort kit line.
2. As a contributor whose gate exercises direct, prospective, stripped, contract, preflight, conformance, coverage, or canary paths, I want every ordinary path to receive the same selected absolute executable and refuse an absent inherited selection, so hidden builders and authoring fallbacks cannot return. Line: gpt-5.6-terra / high. The migration spans several existing seams and a closed exception inventory.
3. As a contributor waiting for a gate, I want all phase tables, including the combined primary/stripped table and inner gates, to execute serially without weakening dependency or cancellation semantics, so compiler demand is bounded without changing what green means. Line: gpt-5.6-terra / high. Scheduler and teardown regressions are high-leverage gate defects.

## Implementation decisions

- The selected executable contract is a process-local value named `BENCH_RUN_BINARY`. Its value is a cleaned absolute path to a regular executable in the current owner's private run directory. A top-level owner overwrites hostile ambient input; an inner gate validates and reuses the inherited value. Symlinks, special files, relative paths, missing files, stale seals, and source-mismatched executables refuse.
- A non-reused gate constructs the owner after subject validation and admission. It builds from the materialized runtime root when that exact subject is the Bench kit; a linked project uses the wrapper-selected kit captured by the owner. Both cases invoke that source root's canonical `scripts/go-build.sh` exactly once in host subject mode. The builder's adjacent seal remains private with the executable.
- `bench test` is the focused-test owner. With no valid inherited selection it creates one private run directory, invokes the canonical builder once, and passes the selected absolute path into its one `go test -count=1` child. With an inherited selection it reuses it and builds zero executables. A red run and a later green run are separate top-level runs and deliberately get separate binaries.
- A top-level gate that reuses an exact green verdict performs no build because it executes no gate. Usage, admission, or subject-validation refusals that occur before owner creation also build zero. Once a private directory exists, every exit path, including a canonical builder that leaves partial output and fails, owns its cleanup.
- `.bench/gate.sh` remains shell. It requires the owner-selected path, executes that binary's `freshness-check` and then the same binary's `gate-phases`, and never calls `go run ./internal/freshness/check`. `.bench/gate-prospective.sh` becomes a no-build pass-through into the same shell route. Neither shell has an authoring fallback.
- Ordinary consumers are exactly: gate-entry freshness and phase routing; the gofmt, test, race, and conformance-suite `gate-go` phase commands; ordinary conformance, coverage, and contract commands exercised by those suites; shellcheck and canary phase launch; unchanged-host helpers in `internal/gate`, `internal/contract`, and `internal/preflight`; stripped-subject consumers; prospective execution; and nested canary gates. These consumers may run `go test` to compile test packages, but they may not invoke another ordinary `go build`, `go run ./cmd/bench`, `go run ./internal/freshness/check`, or subject-mode `scripts/go-build.sh`.
- Wrapper, installation, cache-routing, and publication tests that need Bench bytes copy or link the already selected bytes into their fixture. Tests whose assertion is specifically that a different source, package, target, compiler, linker, release workflow, or reproducibility slot authors different bytes retain their own build.
- The independent-build exception set is closed to: cross-target builds in `scripts/build-artifacts.sh`, `scripts/native-proof.sh`, and `internal/conformance/cross_compile_stress_test.go`; the release build in `scripts/release-preflight.sh` and release-only `internal/preprelease` gate-go invocation; artifact-mode and reproducibility proofs under `internal/contract/surface/artifact/posture/`; changed-source proofs in `internal/contract/freshness_subject_test.go` and `internal/contract/runtime/runtime_gate_freshness_routes_test.go`; alternate-package/compiler attestation in `internal/gate/build_attestation_test.go`; the reduced fixture linker/publication proof in `internal/contract/runtime/runtime_gate_component_boundary_test.go`; and intentional test-executable compilation in `internal/canary/canary.go` and `internal/contract/runtime/runtime_gate_partial_proof_test.go`. Adding a member or constructor requires changing this enumeration and its structural census together.
- The scheduler uses stable declaration order to choose the next ready phase, runs at most one phase process, settles it, then reevaluates readiness. Primary and stripped phases carry their execution root into one combined table. A red or skipped prerequisite still skips only its dependents; unrelated phases still run and all red results aggregate.
- Cancellation continues to signal the active phase's process group, waits through the bounded grace, kills remaining descendants, and reaps before returning. Only after that teardown may the run owner remove the selected executable directory. The same order applies to timeout and interrupt.
- This is a wide refactor and uses `craft-tickets`' expand-migrate-contract sequence. The selection/owner seams expand first, ordinary consumers migrate in independently green cuts, and the exact census plus lifecycle tests contract away every ordinary fallback.

## Testing decisions

- The run-owner seam receives injected builder, filesystem, source-root, and child-runner tests. Tests count canonical builder calls, record executable identity (absolute path, inode, and digest), inject each terminal outcome, and assert cleanup only after descendant settlement.
- Gate execution tests drive the real direct and prospective entry points. A marker command at every ordinary consumer records the selected path; the test requires one identical value and one canonical subject-build trace for the whole run.
- `internal/testreport` tests drive the public `bench test` command with a fake canonical builder and a real child process boundary. They prove one build without inherited selection, zero with one, exact propagation, and a fresh build for the next top-level source snapshot.
- Existing contract and preflight helpers are tested against the real selected executable rather than a process-local template. Wrapper-specific tests observe copied or linked selected bytes; changed-source and compiler-proof tests remain in the enumerated exception set.
- Nested-canary tests execute an actual inner gate process. Missing, relative, stale, symlinked, or changed inherited selection refuses before any builder; a valid path is byte-identical at outer and inner markers.
- Scheduler tests replace the overlap proof with a max-active-process counter and stable-order record. Existing dependency, optional-skip, red-dependent-skip, red aggregation, interrupt, timeout, and orphan-descendant tests stay attached to the real scheduler.
- A structural conformance test owns the exact ordinary-build census and the enumerated independent exception set. It examines assembled argv and builder call sites, not incidental prose. Mutating any ordinary consumer to a forbidden constructor makes the default gate red.

### Seam diagram

    trigger: bench gate, commit gate, prospective gate, or bench test
        │
        ▼
    exact source snapshot ──▶ [ run owner: temp dir + canonical builder once ] ──▶ absolute selected Bench path
                                      │                                              │
                                      │                                   tests attach: count build calls and
                                      │                                   record identity at every consumer
                                      ▼
                           [ ordinary and nested children ] ──▶ terminal outcome + descendants reaped
                                      │
                                      ▼
                              remove private run directory

    trigger: gate-phases outer or inner invocation
        │
        ▼
    primary phases + stripped phases ──▶ [ one stable topological scheduler ] ──▶ ordered results
                                                   │
                                                   └ tests attach: max-active counter,
                                                     dependency/skip/red/cancel record

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| RS1 | 1 | A run owner selects one cleaned absolute host Bench path bound to its exact source snapshot; a later top-level run cannot reuse it. | run-owner builder/selection seam | observed red: `rg -n 'BENCH_[A-Z_]*SELECTED' internal .bench --glob '*.go' --glob '*.sh'` exited 1 | No selected-path contract exists, so current consumers cannot prove identity or run-scoped lifetime. |
| RS2 | 1 | Each non-reused direct or prospective gate builds exactly once through subject-mode `scripts/go-build.sh`; both routes share the owner, while pre-owner and exact-green reuse paths build zero. | real gate execution owner and shell-entry process seam | observed red: `! rg -n 'go run ./internal/freshness/check' .bench/gate.sh` exited 1, and the same assertion for `scripts/go-build.sh` in the prospective shell exited 1 | The current shells independently compile the verifier and prospective CLI, proving ownership is split. |
| RS3 | 1 | A top-level `bench test` builds once and propagates the selected path to its `go test` child; an inherited run builds zero, and the next top-level red/green snapshot gets a new path. | `internal/testreport` command seam | observed red: the RS1 selected-helper query exited 1 while `internal/testreport/testreport.go` launches `go test` with no selected executable | A child identity record and builder counter distinguish one per run from zero propagation or worktree-lifetime reuse. |
| RS4 | 2 | Gate entry, freshness self-check, gate-phases, gofmt, test, race, conformance-suite, ordinary conformance, ordinary contract, shellcheck, and canary phase launch all receive the identical selected path and issue no forbidden ordinary build/run command. | real phase-table argv and gate-entry seam | observed red: `! rg -n 'runPhasesConcurrent' internal/gate/runner.go` exited 1, while the same source assertion found the phase-owned build and gate-go `go run` | Assembled argv exposes the current phase-owned build and gate-go fallback. |
| RS5 | 2 | Every unchanged-host helper under `internal/gate` uses the inherited selected executable and cannot invoke the canonical builder itself. | internal gate helper seam plus structural census | observed red: the ordinary-helper assertion found `internal/gate/runner_serial_test.go` invoking `scripts/go-build.sh` | Counting the real helper's builder trace catches a private test binary even if its bytes happen to match. |
| RS6 | 2 | `internal/contract` coverage/command fixtures and `internal/preflight` unchanged-host helpers use the inherited selected executable; wrapper/cache fixtures copy or link its bytes, while changed-source proofs stay independent. | contract/preflight helper and fixture-command seam | observed red: the ordinary-helper assertion found one contract and five preflight subject-builder call sites | The helper census and real coverage command tests catch repeated current-subject builds while the exception rows keep distinct-source proofs meaningful. |
| RS7 | 2 | A nested canary gate receives the exact outer selected path unchanged and builds zero; missing, relative, stale, symlinked, special, or source-mismatched inheritance refuses before a phase. | real canary outer-to-inner process seam | not TDD-able before RS1: no inherited selected-path seam exists to mutate; the expansion ticket adds it before this migration | A real nested process plus a builder trap distinguishes inheritance from a hidden authoring fallback. |
| RS8 | 1 | Once an owner creates a run directory, canonical-builder failure, green, red, post-owner refusal, timeout, and interrupt all reap any descendants and remove partial output plus the directory; no executable or seal survives. | injected owner lifecycle plus real process-group seam | not TDD-able before RS1: the current tree has no private run-directory owner whose teardown can be observed | Builder and child failure injection at each terminal edge proves both partial-output cleanup and teardown-before-remove ordering. |
| RS9 | 2 | The structural census rejects every ordinary `go build`, `go run ./cmd/bench`, `go run ./internal/freshness/check`, and subject builder outside the single owners, while accepting only the enumerated cross-target, release, artifact-posture/reproducibility, changed-source, alternate-package, linker/publication, and test-executable proofs. | default conformance structural census | observed red: the entrypoint, gate-plumbing, and ordinary-helper source assertions all exited 1 | An exact allowlist makes a new ordinary constructor red instead of letting a category label excuse it. |
| SG1 | 3 | Every outer and inner phase-table invocation runs at most one phase process and chooses ready phases in stable declaration order. | real scheduler max-active seam | observed red: `! rg -n 'TestSchedulerOverlapsIndependents' internal/gate/runner_serial_test.go` exited 1, and the same assertion found concurrent runner symbols | The current positive overlap test and concurrent runner make the cheapest parallel implementation observable. |
| SG2 | 3 | Primary and stripped phases share one serial topological schedule with the correct execution root for each phase. | split-table composition seam | observed red: the SG1 source assertion found the independent `primaryDone` goroutine and stripped schedule | One global max-active counter fails if either table or their composition overlaps. |
| SG3 | 3 | Serial execution preserves needs order, red-dependent skips, optional skips, execution of unrelated phases, and aggregate-red result. | scheduler result seam | already covered by `TestSchedulerRespectsNeeds`, `TestSchedulerSkipsDependentsOfRed`, `TestSchedulerPropagatesOptionalSkip`, and runner red tests | These independent outcomes go red if serialization is implemented as fail-fast or declaration-only execution. |
| SG4 | 3 | Interrupt and timeout cancel the active process group, reap a leader's remaining descendants, return the established code, and only then allow owner cleanup. | scheduler/process-group and owner lifecycle junction | already covered in part by `TestPhasesCommandSignalCancelsRunningPhaseGroups`, `TestRunnerCancelKillsGroup`, and `TestExecuteCancellationKillsDescendantAfterLeaderExits`; RS8 adds cleanup order | Existing process assertions preserve cancellation while the lifecycle marker detects premature directory removal. |

### Edge inventory

- Empty phase table — SG1 records zero active processes and green without constructing another binary.
- One phase and unequal primary/stripped partitions — SG1 and SG2 use the same scheduler and never require both partitions to be non-empty.
- Multiple simultaneously ready phases — SG1 requires stable declaration order with a maximum active count of one.
- Red prerequisite, optional missing tool, unrelated ready phase, and multiple reds — SG3 preserves dependent skips, optional skips, unrelated execution, and aggregate red.
- Interrupt before build, during build, during a phase, after the leader exits, and during cleanup — RS2, RS8, and SG4 bound each transition and require descendant reap before removal when an owner exists.
- Canonical builder creates partial executable or seal output and exits red — RS8 requires no child launch and removes every partial path plus the private directory.
- Timeout during a phase — RS8 and SG4 preserve exit 124, evidence invalidation, process-group teardown, and cleanup.
- Source changes before ownership, during build, or after selection — RS1 and RS2 bind the path to the accepted snapshot; drift refuses the run without a second build.
- Prospective materialization and a linked-project kit — RS2 selects the materialized Bench source for a Bench candidate and the captured wrapper-selected kit for a linked non-Bench subject.
- Repository and temporary paths containing spaces or shell metacharacters — RS1 and RS4 carry argv values without shell interpolation.
- Relative, missing, empty, symlinked, special, non-executable, stale-sealed, source-mismatched, or hostile ambient selected paths — RS1 and RS7 refuse before ordinary consumers.
- Nested gate with a valid inherited path — RS7 records identical absolute path, inode, and digest at both levels and a zero nested builder count.
- Repeated red then green ticket testing — RS3 observes different private paths and exact source identities, never a worktree-lifetime cache entry.
- `go test` package compilation — RS4 permits ordinary test executables while RS9 rejects a second Bench CLI build initiated by those tests.
- Cross-target, release, reproducibility, changed-source, alternate-package, linker/publication, and `go test -c` proofs — RS9 enumerates their exact owners and keeps each independently observable.
- Cleanup after an owner-setup refusal — RS8 removes partial executable, seal, and directory state.

**Won't handle:** durable artifact corruption and repair — no run artifact survives to be discovered or repaired.

**Won't handle:** cross-worktree simultaneous gates — existing admission policy is unchanged; this build bounds only phase execution inside one run.

**Won't handle:** canary fixture-worker overlap — the global canary worker policy is unchanged; only each phase-table invocation is serial.

**Won't handle:** arbitrary subprocess CPU, memory, file, or network resources — the scheduler serializes phase processes, not all descendants or external work.

**Won't handle:** release-target identity unification — independent proof builds keep their current target and authorship semantics.

## Out of scope

- Durable or cross-run artifact storage: 24 edits, 4 gate runs. It is a separate cache/store capability and conflicts with private run lifetime.
- Later-process lookup or reuse of a prior run's executable: 12 edits, 3 gate runs. It requires durable discovery, currency, and repair policy.
- Target-aware proof migration into the selected host executable: 22 edits, 4 gate runs. Those proofs intentionally observe different targets or sources.
- A separately built freshness-verifier artifact: 10 edits, 2 gate runs. The selected Bench binary already owns `freshness-check`.
- Transferable process-tree authority tokens: 18 edits, 3 gate runs. Nested gates carry only the validated selected path and existing gate context.
- The 45-minute progress watchdog: 14 edits, 3 gate runs. This build preserves existing timeout policy.
- Arbitrary subprocess/resource serialization: 20 edits, 3 gate runs. Only phase-table process admission changes.
- Global canary-worker serialization: 14 edits, 3 gate runs. Canary internals retain their independent worker policy.
- Cross-worktree or common-Git-directory admission changes: 16 edits, 3 gate runs. Existing gate lock behavior remains authoritative.
- Converting `.bench/gate.sh` into a Go module or Go entry point: 8 edits, 2 gate runs. The project gate remains shell and delegates to the selected binary.
