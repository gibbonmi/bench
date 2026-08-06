# Go build cache footprint

Status: implemented

Decision source: reviewer-confirmed conversation of 2026-08-04

## Problem

Bench currently pays for two executables when a local caller asks for one. The canonical Go builder compiles the Bench binary, then `go run`s a second publication helper. Repeated helper executables accounted for about 10.37 GB of a 53–54 GB Go build cache, and the resulting daily cache trim stalled the requested build for 10.724 seconds even though the publisher itself took only 0.573 seconds.

The gate also allows Go to retain successful verdicts for its core-package and filtered-conformance invocations. Filesystem-heavy tests make those verdict records enormous: after the shared cache was cleaned, one cold repository workload produced three reusable `# test log` records totaling more than 75 MB for only three tests. These records spend disk and trim budget even though an oracle run must not reuse a successful verdict.

## Solution

Compile each requested Bench deliverable once. A default local build runs the newly compiled host binary's internal publication operation to create the atomically promoted executable and matching freshness seal. An explicit artifact mode atomically promotes the staged executable without a seal and without executing it, so release and native-proof callers no longer create and delete seals around cross-compiled outputs.

Keep Go's compilation cache, retained component evidence, and the established dev-versus-ship cache posture. Whenever the gate selects either remaining oracle-owned Go test phase for execution, disable only successful test-result caching at its argv constructor by adding `-count=1`. The durable contract is fresh Go test verdicts at those executed gate seams; cache-size inspection is the design and build-verification evidence.

## User stories

1. **A local Bench builder receives one compiled, trusted subject.** The default two-argument `go-build.sh` invocation performs one `go build`, invokes no `go run` publisher, and leaves an executable whose freshness seal matches both its bytes and the current source closure.
   `Line: gpt-5.6-terra / medium.` The expected shape is closed, but self-publication and failure-safe promotion affect the executable the oracle trusts.

2. **Artifact producers receive one unsealed deliverable without host execution.** An explicit semantic artifact mode performs one `go build`, atomically promotes the staged output, emits no freshness seal, and never executes the output; release-artifact and native-proof callers select that mode directly.
   `Line: gpt-5.6-terra / medium.` The script seam is known and strongly observable, but cross-platform behavior and prior-output preservation make a cheap plumbing pass too narrow.

3. **Every executed gate test phase runs current tests while retaining compilation reuse.** Whenever component selection executes the core packages or filtered conformance suite, that Go invocation disables successful test-result caching while compiled packages remain reusable through the ambient or private `GOCACHE` selected by the existing tier policy. Unchanged components may still skip on valid retained gate evidence.
   `Line: gpt-5.6-terra / medium.` The edit is small, but it changes oracle semantics and therefore follows the gate/conformance routing.

## Decision evidence

The 2026-08-04 controlled comparison ran from the requested `ff95c1530f5f3f51ca9e22155593f469a5b564c5` tree after the other live gate owner released. Both caches were fresh `mktemp` directories under `/tmp`, outside the repository; both used the same ambient module cache and `GOMAXPROCS=8`.

Current argv:

```bash
current_cache="$(mktemp -d /tmp/bench-test-cache-current.XXXXXX)"
GOCACHE="$current_cache" GOMAXPROCS=8 go test ./internal/specbuild ./internal/worktree
GOCACHE="$current_cache" GOMAXPROCS=8 go test ./internal/conformance -skip '^(TestRootConformance)$'
rg -0 -l -a -g '*-d' '^# test log$' "$current_cache"
```

The cache contained three `# test log` records totaling 82,810,086 bytes inside a 206,610,102-byte cache:

| bytes | attributed test |
|---:|---|
| 29,738,392 | `TestAbandonPlansForRemovedWorktree` |
| 27,312,779 | `TestBenchShRouteAnchorBites` |
| 25,758,915 | `TestApplyAbandonReleasesHuskWithoutDeletingBytes` |

Cache-disabled argv:

```bash
count1_cache="$(mktemp -d /tmp/bench-test-cache-count1.XXXXXX)"
GOCACHE="$count1_cache" GOMAXPROCS=8 go test -count=1 ./internal/specbuild ./internal/worktree
GOCACHE="$count1_cache" GOMAXPROCS=8 go test -count=1 ./internal/conformance -skip '^(TestRootConformance)$'
rg -0 -l -a -g '*-d' '^# test log$' "$count1_cache"
GOCACHE="$count1_cache" GOMAXPROCS=8 go test -x -count=1 -run '^$' ./internal/specbuild ./internal/worktree ./internal/conformance
```

The `-count=1` cache contained zero `# test log` records and 123,797,749 bytes of compilation artifacts. The warmed `-x` probe executed zero compiler commands, left the cache at exactly 123,797,749 bytes, and still left zero test logs; it relinked the three test binaries. This confirms that `-count=1` disables reusable successful verdicts without disabling compilation caching. Local `go help testflag` independently names `-count=1` as the idiomatic explicit way to disable test caching.

After measurement, each disposable cache was reduced with `GOCACHE=<dir> go clean -cache` and its remaining 228-byte metadata removed. This cleanup is probe hygiene, not product behavior.

## Implementation decisions

**The builder has two semantic modes, not an environment heuristic.** The existing two positional arguments remain the default local-subject invocation. `--mode artifact` is the only added mode selector. Missing operands, a missing mode value, duplicate mode selectors, and unknown modes fail with usage before staging or touching an existing output. `GOOS` and `GOARCH` never choose the mode.

**Default mode always produces a host-runnable freshness subject.** It resolves the toolchain's host target explicitly rather than trusting ambient target overrides, builds one staged Bench executable beside the destination, and invokes that staged executable directly. Ambient `GOOS` or `GOARCH` cannot turn the default into an artifact build or cause a cross-compiled binary to be executed on the build host.

**The staged Bench binary owns sealed publication.** Add one internal plumbing operation whose invoked executable is the staged input and whose explicit arguments are the source root and destination. It delegates digesting, executable validation, atomic promotion, and seal creation to the existing freshness publication owner. Binding the staged input to the running executable avoids a second caller-controlled executable path and removes the need for a separately compiled publisher.

**Artifact mode owns unsealed atomic promotion.** It uses the same beside-destination staging discipline, honors the caller's explicit target environment, validates the staged output, and atomically renames it into place without running it. Success guarantees that no destination seal remains. Any failure before promotion leaves the prior output and its prior seal byte-identical; temporary staging is removed on handled exit or signal.

**Sealed publication remains a single observable operation path.** The builder-visible command surface has exactly one sealed-publication entry — the built binary's subject operation — and the production call graph enumerates every `freshness.Publish` caller in any package as exactly that entry. Artifact mode has no path to that entry: it neither executes a Bench binary nor invokes another Go-built helper. A call-path contract and the builder's command trace jointly distinguish true unsealed promotion from sealed publication followed by seal deletion; final filesystem state alone is not accepted as evidence.

**Artifact callers state their intent once.** Release artifact generation and native proof request `--mode artifact`. They no longer delete seals after the builder returns. Native proof may execute the rebuilt output only later, on the native runner where its canonical target is runnable; publication itself never executes an artifact-mode output.

**The standalone publisher package disappears.** Once every caller uses the built binary or artifact promotion, remove the separate freshness publisher command. No replacement helper executable, `go run` publisher, or post-build `go clean` is introduced.

**Existing cache posture is unchanged.** Development artifact contracts may opt into ambient `GOCACHE` and `GOMODCACHE` only with the exact `BENCH_SHARED_BUILD_CACHE=1` token. Ship-tier artifact generation and native proof retain private disposable caches and independent reproducibility evidence. This implements neither a cache quota nor automatic eviction.

**The gate disables verdict reuse at its two remaining owners.** When selected for execution, `coreTestStep` constructs `go test -count=1 <packages>`. `ConformanceSuiteArgv` constructs `go test -count=1 ./internal/conformance -skip <registry-pattern>`. Package enumeration, tier exclusions, retained component-evidence selection, the registry-owned skip pattern, output streaming, and exit semantics do not move. Individual filesystem-heavy tests are not edited because the controlled probe removed their records at this higher seam.

**The oracle change carries a durable bite.** Focused gate tests assert `-count=1` at both argv owners and emit owner-specific failures. A gate-path recorder supplies a disposable `GOCACHE`, requires the test process to receive that exact path, refuses a `go clean` subprocess, plants a synthetic `# test log` after the fake test execution, and requires that record to survive the owner return. This makes private/disabled cache overrides and post-run record deletion red independently of the `-count=1` assertion. The implementation demonstrates the central omission mutation by removing the argument at each owner in turn and observing its focused contract go red. The real disposable-cache comparison is repeated against the implemented argv to show the observed giant records are gone, but cache byte thresholds do not become the oracle contract.

**The two ownership fences are a deliberate bundle.** Publication/artifact work and gate-verdict work are independently shippable and touch disjoint owners. Splitting them would price at about 12 edits and 1 gate run for publication plus 3 edits and 1 gate run for verdict caching; the reviewed bundle is about 15 edits and 1 composed gate run. The reviewer explicitly expanded this spec to include the gate-record remediation on 2026-08-04, so the extra spec lifecycle and gate run buy no scope clarity.

## Testing decisions

- **External behavior.** Drive the real builder through a recording `go` wrapper that forwards the real toolchain while recording invocations. Assert the exact number and kind of Go commands, then execute or inspect the resulting subject through the existing freshness contract. Drive artifact mode with a deliberately non-host target and an execution marker so any attempt to run the output is red.
- **Primary publication seam.** The builder script is the highest seam: it observes argument grammar, target-environment handling, staging, invocation count, the complete command trace, output replacement, seal presence, interruption cleanup, and caller-visible diagnostics in one journey. The trace records every Go and Bench/helper execution, not only the staged output.
- **Lower publication seam.** The built binary's internal publication operation and the freshness publication owner receive focused tests for invoked-executable binding and seal/output agreement. This lower seam is justified because a script-only assertion could pass with a fake publisher that writes plausible files without using the trusted freshness owner.
- **Artifact caller seam.** Existing artifact-generation and native-proof contract fixtures observe the real scripts, generated binaries, cache selection, reproducibility evidence, and absence of seals. Their existing shared-versus-private cache controls remain positive regressions.
- **Sealed-publication operation seam.** A production call-path contract proves that only the subject operation reaches the freshness owner's sealed publication. Combined with the artifact execution/command recorder, this is the red-capable evidence that artifact mode never creates a seal and deletes it afterward; final seal absence is only the external outcome.
- **Gate seam.** Gate-package tests drive the real `gate-go test` owner through a temporary module and inspect the filtered conformance argv. A fake `go` at that same path records argv and effective `GOCACHE`, refuses cleanup calls, creates a synthetic test-log record, and leaves the owner to return; the test requires the inherited cache path and record to survive. These assertions avoid reimplementing package enumeration or the registry skip pattern.
- **Prior art.** Reuse the freshness subject and route fixtures, the artifact cache-posture and reproducibility fixtures, `GateGoCommand`, and `ConformanceSuiteArgv` tests. Do not create a second fixture harness or a second source for build flags, package lists, skip patterns, or cache policy.
- **Demonstrated bite.** Omit `-count=1` independently from core and conformance construction; each owner-specific test must fail with its own message. Retain the omission while running the disposable-cache repro and observe the corresponding `# test log` record return. Restore the argument and obtain zero records.

### Seam diagrams

```text
trigger: go-build.sh [--mode artifact] <module-root> <output>
    |
    v
validated mode + target + hostile output path
    |
    v
[ one staged go build beside destination ]
    | subject mode                         | artifact mode
    v                                      v
[ invoked Bench publication operation ]   [ unsealed atomic promotion ]
    | uses freshness.Publish               | never executes staged output
    v                                      v
executable + matching seal                 executable + no seal
    ^ tests attach: recording go wrapper, real freshness verification,
      prior-output fingerprints, signal seam, and execution marker
```

```text
trigger: dev gate test phase
    |
    +--> package enumeration --> [ coreTestStep argv + -count=1 ] --> fresh core verdicts
    |
    +--> registry skip pattern -> [ ConformanceSuiteArgv + -count=1 ] --> fresh suite verdict
                                    ^ tests attach: owner argv contracts;
                                      disposable GOCACHE proves no # test log records
```

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Default `go-build.sh <root> <output>` invokes exactly one `go build` and no `go run` or second Go-built publisher | real builder with forwarding Go recorder | observed 2026-08-04: the current builder executes `go build` at line 61 and `go run ./internal/freshness/cmd` at line 62 | Counting the real tool invocation makes retaining or renaming the helper build red |
| 1 | Default output is executable and has a seal matching its bytes and current source closure | real builder followed by freshness verification in a fresh process | existing freshness contracts cover valid and stale subjects; the invocation-count journey is to be observed red before the script change | Skipping freshness or writing a plausible but unrelated seal fails the existing verifier |
| 1 | Ambient cross-target `GOOS` and `GOARCH` do not select artifact behavior: default mode still publishes a host-runnable sealed subject | real builder with non-host ambient target, then execute `version` and verify the seal | to be observed at build time; current default builds the ambient cross target and seals it, leaving a non-runnable local subject | An implementation that infers mode from the target environment or executes a cross binary cannot satisfy both execution and seal assertions |
| 1 | The publication operation binds the staged input to its own invoked executable and delegates to the freshness owner | internal operation test plus real builder journey | to be observed at build time by attempting to pass or substitute a different staged path | A thin replacement that still accepts a helper executable path can preserve the second-build architecture |
| 1 | The standalone freshness publisher command and every production caller are absent | package/build topology contract | current tree contains `internal/freshness/cmd` and the builder calls it | Deleting only the `go run` token while retaining an equivalent helper package fails the topology assertion |
| 2 | Artifact mode invokes exactly one `go build`, invokes no `go run` or Bench/helper publication operation, emits no seal, and never executes the staged output | real builder with complete command recorder and non-host execution marker | current builder has no artifact mode and invokes a Go-built sealed publisher before callers delete its seal | A mode that publishes sealed then cleans up, retains any helper executable, or probes the cross output trips the command trace, operation path, seal, or marker assertion |
| 2 | The production sealed-publication call graph enumerates every production caller in any package as exactly the subject-operation entry; artifact promotion cannot reach it | publication operation call-path contract | current script routes every mode through the standalone sealed publisher | Final-state-only coverage permits create-then-delete; enumerating every production caller in any package makes any artifact reachability, or any new caller, red |
| 2 | Release artifact and native-proof callers request artifact mode and contain no post-build seal deletion | real artifact/native-proof contract plus caller-shape assertion | current callers omit a mode; artifact generation deletes `binary.seal` after every build | Fixing only direct builder tests leaves production callers on the old create-then-delete path |
| 2 | A failed compile, failed validation, or failed promotion leaves the prior output and prior seal byte-identical | builder failure table with fingerprints | to be observed at build time before implementation for each failure site | Direct-to-output build or preemptive seal removal corrupts the previously usable deliverable |
| 1, 2 | SIGINT while compilation or publication is blocked leaves the prior output intact and no staged executable behind | builder subprocess with marker-controlled blocker | to be observed at build time | A trap that handles only ordinary failure leaks cache-bearing scratch executables or destroys the prior output on interruption |
| 1, 2 | Missing operands, missing mode values, duplicate mode selectors, and unknown modes refuse before invoking Go or touching output | builder argument table with Go tally and fingerprints | current script has no mode grammar; new cases are to be observed red before parsing changes | A permissive parser can silently fall back to subject mode and execute an unintended target |
| 1, 2 | Absolute and module-relative outputs containing spaces, glob characters, and dash-led segments are handled literally | real builder hostile-path table | to be observed at build time with exact output and seal paths | Unquoted staging or promotion writes a widened or different path while ordinary paths stay green |
| 1, 2 | Directory, FIFO, socket, device, dangling-symlink, symlink-component, and unwritable output targets refuse without blocking or changing prior bytes | builder hostile-output table | to be observed at build time per available host capability | A read-first or follow-symlink implementation can block, escape the named destination, or clobber foreign bytes |
| 2 | Artifact-mode rerun atomically replaces the output and leaves no stale seal from a prior subject-mode build | two-mode rerun through the real builder | to be observed at build time | Merely declining to create a new seal leaves the old seal beside unrelated artifact bytes |
| 2 | Native proof executes the artifact only after build on the canonical native runner and retains its private cache | native-proof contract | existing native-proof and cache-posture contracts cover runner execution and private `GOCACHE`; extend them to assert artifact mode | Moving execution into the builder or changing cache ownership can pass direct artifact creation while weakening independent proof |
| 2 | `BENCH_SHARED_BUILD_CACHE=1` remains the only dev opt-in to ambient `GOCACHE`/`GOMODCACHE`; absent, empty, and other values remain private | existing artifact cache-posture contracts | already covered; run unchanged after caller-mode edits | It prevents the cache-footprint work from weakening the cost-follows-project-size decision |
| 3 | Whenever selected, core gate tests execute as `go test -count=1 <enumerated packages>` at dev and ship tiers | `gate-go test` argv contract plus existing package-set controls | current `coreTestStep` omits `-count=1`; the omission mutation must restore the focused red | A change applied only to another phase or through ambient `GOFLAGS` cannot satisfy the owner argv contract |
| 3 | Whenever selected, filtered conformance executes as `go test -count=1 ./internal/conformance -skip <registry pattern>` | `ConformanceSuiteArgv` contract plus real suite marker fixture | current argv omits `-count=1`; removing the added argument must fail with the conformance-owner message | It preserves the registry skip source while making verdict freshness explicit at the actual owner |
| 3 | The test owners inherit the caller-selected `GOCACHE` unchanged, invoke no cache cleanup, and leave a recorder-created `# test log` sentinel in place after returning | real `gate-go` owner with fake-Go argv/environment/filesystem recorder | to be observed at build time; the recorder contract does not exist | `GOCACHE=off`, a new private cache, `go clean`, or direct post-run record deletion can coexist with correct argv, but each fails this effective-path control |
| 3 | Valid retained component evidence may still skip an unchanged core or conformance-suite component; `-count=1` governs only a phase selected to execute | existing component-decision and input-identity contracts | already covered; run unchanged after the argv change | It prevents fresh Go verdicts from silently expanding into removal of the gate's separate retained-evidence policy |
| 3 | Core package enumeration, release-only exclusions, conformance skip behavior, output streaming, and red exit semantics remain unchanged | existing gate-go package, red, and marker controls | already covered; run unchanged | A broad argv refactor cannot buy cache behavior by dropping packages, filters, output, or failures |
| 3 | The representative current argv produces the three observed giant records, while the implemented `-count=1` argv produces zero `# test log` records | disposable-cache repro recorded above and repeated at build completion | observed 2026-08-04: current 3 records / 82,810,086 bytes; `-count=1` 0 records / 0 bytes | This falsifies a flag placement that looks correct in a unit assertion but does not reach the Go invocation |
| 3 | `-count=1` preserves compilation caching | warmed disposable-cache `-x -run '^$'` probe | observed 2026-08-04: 0 compiler commands and unchanged 123,797,749-byte cache after the `-count=1` run | Disabling or discarding the whole build cache would remove test logs too, but fails this control |
| 3 | Successful verdicts are not made reusable by post-run cache deletion | owner argv plus effective-cache recorder | current owner argv is cacheable; the new recorder is to be observed red when its planted log is deleted after fake test execution | Adding correct argv and then cleaning records still fails the surviving-record assertion |

The composition degenerate is a builder that removes the helper only for local subjects while artifact callers still create and delete seals, combined with a gate that uses `-count=1` only in one of its two invocations. The publication call-path/command trace, real caller journeys, and independent gate-owner rows make every half red. The global-cache degenerate—`go clean -cache`, `GOCACHE=off`, a new private cache for every dev command, or direct test-log deletion—fails the effective-cache recorder even when it also supplies the correct argv.

### Edge inventory

- Error path — resolved by compile, staged validation, publication, promotion, cache-resolution, and gate-red rows.
- Empty or absent input — resolved by missing argument/mode cases, missing source roots and build-input manifests in existing builder controls, and absent prior output in the first-run journeys.
- Boundary values — resolved by first publication, replacement of an existing subject, artifact-after-subject, and exactly one requested deliverable.
- Malformed input — resolved by unknown/duplicate modes and special, symlinked, and unwritable output targets.
- Interrupted or partial state — resolved by the SIGINT row and prior-output fingerprints. **Won't handle:** uncatchable process or host failure between filesystem operations; the resulting subject may refuse freshness verification, but the build does not claim crash-atomic replacement of two directory entries.
- Re-run idempotency — resolved by subject rebuild, artifact rebuild, and subject-to-artifact replacement rows.
- Process-boundary lifecycle — resolved by verifying and executing the default subject in a fresh process, and by artifact/native-proof callers consuming the builder in separate processes.
- Hostile environment — resolved by ambient cross-target variables, shared/private cache controls, missing tool/build-input controls, literal hostile paths, special files, symlinks, and non-host artifact markers.
- A command whose write changes a fact it reports — **Won't handle:** `go-build.sh` emits no cleanliness, cache-size, or freshness-status report; callers verify the post-write subject independently.
- Control bytes in git-sourced text — **Won't handle:** neither mode reads a Git-sourced path or renders TOON; caller-supplied filesystem path diagnostics remain the shell's existing stderr surface.
- Final line without a newline — existing build-input manifest parsing already accepts the last key without requiring a terminator and remains a regression control.
- Non-TTY stdin — **Won't handle:** the builder and gate invocations never prompt, so stdin mode cannot amputate a caller.
- Host-backed filesystem I/O pressure — resolved only through stage-then-rename and prior-output preservation. **Won't handle:** a cache quota or performance guarantee under host pressure is a separate capability.

## Out of scope

- Other Go build-cache growth, including unrelated command executables or compilation artifacts — a separate attribution-and-remediation capability, about 6 edits and 1 gate run.
- Automatic local cache eviction, `trim.txt` manipulation, global cache disablement, or a local quota daemon — a separate cache-governance capability, about 10 edits and 2 gate runs.
- Historical shared-cache cleanup or a migration command — an operator-maintenance capability, about 4 edits and 1 gate run. The already-authorized manual cleanup is complete.
- GitHub cache quota and retention enforcement — external repository settings, 0 repository edits and 0 gate runs. Making `setup-go` caching explicit in workflow YAML without claiming enforcement would be a separate 2-edit, 1-gate-run clarity change; it does not remediate this local footprint and is not included.
- Refactoring the three filesystem-heavy tests to emit fewer observations — rejected as an alternative, not deferred. The controlled probe shows that disabling reusable gate verdicts removes their giant records without weakening their red behavior.
- Weaker freshness checks, an unsealed default local binary, target-environment mode inference, executing cross-compiled outputs on the build host, or deleting cache records after each gate run are refused by the decision source rather than priced future capabilities.
