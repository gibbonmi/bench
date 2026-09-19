# Test determinism on small runners

Status: staged

Decision source: reviewer-confirmed current conversation, 2026-09-19. The reviewer chose items 1 to 3 of the runner answer:

1. No wall-clock limit as a correctness signal.
2. One hermetic test-environment owner.
3. No shared files between parallel tests.

The nightly stress job is out of scope and parked in the idea inbox. Bigger runners and `-p 1` are rejected. The causes of CI families B and C stay unconfirmed. Linked repositories keep their own environment.

Verification log: 2 iteration(s) to accept — the Fable 5.1 medium review returned seven findings, and pass 2 folded all of them. It adds the switch census and the subprocess switch removal (TD45), the discovery window owner (TD46), and merge composition for the ft290 overlap. It also adds the second-accessor rule for the three moved windows. It adds the build-script ticket for the confirmed `dist/` writer (TD47, TD48). It adds the named review artifact for TD19 and TD22, and one resolved line. It folds the switch-name recommendation too.

## Problem

The native-runtime pipeline goes red on tests that pass on a developer machine. The CI runner is a small virtual machine, and the race phase runs about four times slower there than on a developer host. Three classes of test defect turn that slowness, or a new runner image, into a red:

- A production time bound expires on a slow runner and changes a verdict. The test then fails on a correct build.
- A test sees the ambient environment of the runner. Family A was an example: a new git default ran a background prune that deleted a test fixture.
- Two tests that run in parallel share one file. One test writes the file while the other test reads it or asserts on it.

Each red costs a debug session. A red also teaches the team to re-run the pipeline, and that habit hides real failures.

## Solution

A kit test run gets one hermetic environment from one owner. The environment has a private `HOME`, a private `TMPDIR`, no global or system git configuration, and the git maintenance policy. The operator's Go caches and Go settings stay pinned, so no build becomes cold and no module download starts. The gate's kit phases, `bench test` in the kit, and the release preflight phases use that owner.

In a kit test run, a production verdict bound does not expire. A test that wants a timeout sets its own bound, and that bound still expires. A test that runs a Bench process out of process gives that child the policy bounds again. A test that hangs fails with Go's own test timeout, which names the test. Every timed wait in production code reads its window through the bounds package, so a new wall-clock guess cannot enter without a visible choice.

The gate goes red when a kit test phase changes the live checkout outside a declared set. The failure names each changed path. Test helpers that share a fixed temporary name give each call its own directory. The tests that run the kit's build scripts, and the named-check test, work on a private copy of the kit.

## User stories

Line: opus / high.
Implementation-line reason: TD-C3 is the hardest chunk, because the switch changes the timing posture of the whole suite. A wrong window choice hangs only under load, where the gate catches it late.
Harder chunks: TD-C2, TD-C3.

Hermetic environment:

1. As a kit maintainer, I want kit test runs to ignore the global git configuration, so that developer settings cannot change a verdict.
2. As a kit maintainer, I want kit test runs to ignore the system git configuration, so that a runner image cannot change a verdict.
3. As a kit maintainer, I want a kit test run to keep git auto-maintenance off, so that the family A fix stays in force.
4. As a kit maintainer, I want a private `HOME` per kit test run, so that a test cannot touch operator files or another run.
5. As a kit maintainer, I want a private `TMPDIR` per kit test run, so that a leaked temporary file stays with that run.
6. As a kit maintainer, I want the Go child to resolve `GOMODCACHE`, `GOPATH`, and `GOENV` as the operator's environment does. Then a private `HOME` starts no module download and drops no Go setting.
7. As a kit maintainer, I want the Go child to keep the operator's Bench build cache, so that the build cache stays warm.
8. As a kit maintainer, I want the whole run directory removed after the child exits, so that a run leaves no residue. This includes a directory that a test left unreadable.
9. As a kit maintainer, I want the private `TMPDIR` to add at most 16 bytes to the operator's `TMPDIR`. Then test socket paths stay under the platform limit and need no capability skip.
10. As an operator, I want a run to refuse before the child starts when `HOME` is absent, empty, or relative. The refusal names the cause, so I do not read a failed build instead.
11. As an operator, I want a run to refuse before the child starts when `TMPDIR` cannot hold the run directory. The refusal names the cause, so I do not read a failed test instead.
12. As an operator, I want a `TMPDIR` path with a space to work, so that a normal host path does not break the run.
13. As a linked-repository operator, I want my repository's tests to keep my environment, so that the kit's test policy never reaches my project.
14. As a kit maintainer, I want three runners to carry the entries of one owner, so that the runners cannot drift. The three runners are the gate's kit phases, `bench test` in the kit, and the release preflight phases.
15. As a kit maintainer, I want the new owner to replace the git test configuration function, so that the git policy has one source.
16. As a developer, I want the owner's documentation to say that a plain `go test` is outside the policy. Then nobody assumes coverage that does not exist.

Verdict bounds:

17. As a kit maintainer, I want no production verdict bound to expire in a kit test run, so that slow runners cannot change verdicts.
18. As a kit maintainer, I want a bound that a test sets to expire in a kit test run, so that timeout behavior stays tested.
19. As a kit maintainer, I want the bounds switch to accept only the exact value `1`, so that a stray value cannot turn bounds off.
20. As an operator, I want every bound unchanged when the switch is absent, so that production keeps its fail-safes.
21. As a kit maintainer, I want the hermetic environment to set the bounds switch, so that all three runners carry it.
22. As a kit maintainer, I want an unbounded wait to create no deadline at all, so that a sentinel window can never expire at once.
23. As a kit maintainer, I want every production timed wait to read its window through the bounds package. Then a new wall-clock guess needs a visible choice.
24. As a kit maintainer, I want cancel graces and poll intervals fixed under the switch, so that kill escalation and sampling keep their timing.
25. As a kit maintainer, I want a hang in a kit test run to fail as Go's test timeout. That red names the test that hung.
26. As a kit maintainer, I want the whole gate green with the switch on, so that the posture change lands green. I also want the spec to name each test that needed its own bound.
27. As a kit maintainer, I want a test that runs a Bench process out of process to remove the switch from that child. The child then keeps the policy bounds that the test measures.
28. As a kit maintainer, I want the two session inspection windows to read the first accessor, so that the owner table covers every verdict window. The two windows are the discovery window and the provider window.

Shared files:

29. As a kit maintainer, I want each conformance probe to get its own Bench home, so that no two probes share a fixed temporary name.
30. As a kit maintainer, I want the npm cache to stay the one declared shared probe cache. Offline smoke probes then keep a warm cache.
31. As a kit maintainer, I want the build-script tests to work on a private copy of the kit. Then they leave the live `dist/` and the live broker manifest untouched.
32. As a kit maintainer, I want the named-check test to grade a root that no other package writes. Then a concurrent conformance run cannot change its verdict.
33. As a kit maintainer, I want the gate to go red when a kit test phase adds a path to the live checkout. Then a test that writes shared files cannot land.
34. As a kit maintainer, I want the gate to go red when a kit test phase changes or removes an untracked or ignored checkout path. Then no write to the checkout goes unseen.
35. As a kit maintainer, I want a red gate when a kit test phase adds or changes a Bench-owned file in the git directory. Then no shared record changes under a parallel reader.
36. As a kit maintainer, I want the guard to permit the conformance timing file and the gate's own run records. Then the guard reds only undeclared writes.
37. As a kit maintainer, I want the guard failure to name each changed path in byte order, so that I can find the writer.
38. As a kit maintainer, I want the guard to go red when it cannot read the checkout state, so that an unreadable state never passes.
39. As a linked-repository operator, I want no checkout guard on my repository's phases, so that the kit's test policy never reaches my project.
40. As a kit maintainer, I want the current suite to pass the guard, so that the guard lands green.

Reviewed exclusions:

41. As a kit maintainer, I want the record-age windows left unchanged, so that record freshness and staleness keep their meaning.
42. As a kit maintainer, I want `bench test` and the preflight race phase to run without the checkout guard. The gate alone owns that guard.

## Implementation decisions

### The kit test run

The environment package owns one kit test run. A runner opens the run from its base environment before the child starts. The run creates one run directory directly under the operator's `TMPDIR`, with a short name. The run directory holds a private home directory and a private temporary directory.

The runner reads the run's entries, merges them over its child environment, and closes the run after the child exits. Close removes the whole run directory. Close first restores owner permissions on each directory inside it, because a test can leave a directory unreadable.

The entries are these:

- `HOME` names the private home directory.
- `TMPDIR` names the private temporary directory.
- `GIT_CONFIG_GLOBAL` names the null device, and `GIT_CONFIG_NOSYSTEM` is `1`.
- The git configuration count, key, and value entries turn `maintenance.auto` off. These entries replace any configuration the base environment passes to git the same way. This policy moves here from the existing git test configuration function, and that function goes away.
- `GOMODCACHE`, `GOPATH`, and `GOENV` hold the values that Go resolves under the base environment. The run asks Go for the values. It does not copy Go's default rules.
- The bounds switch is `1` (from TD-C3).

`GOCACHE` is not an entry. Each runner already derives the Bench build cache from the operator's `HOME` before it merges the entries, and that derivation stays first.

Open refuses before it creates anything when the base `HOME` is absent, empty, or relative. It refuses when it cannot create the run directory under the base `TMPDIR`. Each refusal names the variable.

The gate decides kit-only scope with its existing kit-root predicate. The gate's kit phases, `bench test` when its Go child runs in the kit, and every release preflight external phase open a run. Release preflight grades the kit alone. A plain `go test` opens no run. No single owner exists for it, because Go has no module-wide test entry. The owner's documentation states that gap.

### Verdict bounds

The bounds package owns one switch. Its name is `BENCH_TEST_UNBOUNDED_WAITS`, and the `TEST` in the name marks it as a test posture. Its only accepted value is exactly `1`. The bounds package documentation states that the switch removes every verdict window, the 45-minute gate timeout included.

The package owns two accessors for a wait window. The first accessor returns the policy window, or the bounds-owned unbounded value when the switch is on. The second accessor returns its argument unchanged. It marks a window, such as a cancel grace, a poll interval, or an operator-set wall limit, that the switch must not change.

Every bounds wait function treats the unbounded value as no deadline. No caller outside the bounds package adds a window to a time, because a very large window can overflow into the past.

Each package variable that holds a verdict policy window initializes through the first accessor. The existing setter functions for tests assign a raw window, so a window that a test sets still expires under the switch.

The verdict policy windows are these seven:

- the provider timeout
- the environment discovery timeout
- the git refresh timeout
- the worktree list timeout
- the guard scan timeout
- the gate timeout
- the package load timeout

Session inspection reads the provider timeout and the environment discovery timeout as constants today. It gains one package variable for each window, both through the first accessor, and a setter for tests. Ticket 8 moves three local windows into the bounds registry and adds them to this set. They are the ref-check timeout, the handoff lock deadline, and the capture lock wait.

The bounds-policy check grows one rule. In production code outside the bounds package, a timed wait passes its window through one of the two accessors. The rule covers each context timeout or deadline, each timer, each after-channel, each after-function, and each deadline computed from the current time. The check's owner table names the first accessor around each verdict window, and its required list names each verdict window. The owner table is the inventory of verdict windows.

A kit test run has no timed backstop from the bounds package. Go's own package test timeout is the backstop. The gate's test argv sets no timeout, so Go's default applies. A hang then fails with Go's timeout panic, and the panic names the running test.

### Subprocess children

A test that runs a Bench process out of process removes the switch from that child's environment. The child then has the policy bounds, and a wait that the test derives from a policy window again contains the child's inner bound. In the system suite, one base-environment helper removes the switch. Three sites build a child environment from the process environment, and each calls that helper. They are the owner's child environment seam, the start of a selected executable, and the start of an artifact landing. The helper names the switch through the bounds package, so the name has one source.

### The live-checkout guard

The gate's phase command reads the checkout state of its root before the kit phases start and after they end. The state is the path list of untracked and ignored files with the size and modification time of each file. It also includes each Bench-owned file in the root's git directory with its size and modification time. When the two states differ outside the declared set, the gate goes red. One failure row names each changed path, in byte order.

The declared set holds two kinds of path. The first is the conformance timing file, which the root conformance entry writes by design and the gate prints. The second is the gate's own run records: the run log, the run stream, the gate lock, and the gate owner record. Each declared path comes from the function that owns it, so the guard restates no path.

A root whose checkout state cannot be read makes the gate red with one row that names the unreadable state. A linked root gets no guard. A new declared path needs a reviewer decision.

### Private kit copies

The shared git test scaffold owns one kit copy helper. It copies each file that git lists as tracked, or as untracked and not ignored, from the kit root into a fresh temporary directory. It then initializes a git repository there with one commit. The copy holds no ignored file, so it holds no `dist/` and no broker manifest.

The two tests that run the kit's build scripts work on a copy: the subject-mode build test and the release preflight build test. Each runs the real script inside the copy, so each still grades the operands that the script hands the builder. The named-check test grades a copy too.

### Shared temporary names

The conformance subprocess environment gives each call its own Bench home under the call's own temporary directory. It keeps the npm cache as the one declared shared cache under the run's `TMPDIR`.

### Timestamps

Gate records truncate their times to whole seconds. The one ordering comparison between two record times compares two truncated values. This spec changes no timestamp.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| TD-C1a / `1-open-kit-test-run.md` | The kit test run owner exists, and the gate's kit phases carry its entries. | TD1, TD2, TD3, TD4, TD5, TD6, TD8, TD9, TD10, TD11, TD12, TD13, TD14, TD15 | `bench test --package ./internal/env`, `bench test --package ./internal/gate` | no |
| TD-C1b / `2-compose-kit-test-run.md` | `bench test` in the kit and the release preflight phases carry the same entries, and the git policy has one source. | TD7, TD16, TD17, TD18, TD19 | `bench test --package ./internal/testreport`, `bench test --package ./internal/releasepreflight`, `bench test --package ./internal/env` | no |
| TD-C2 / `3-isolate-conformance-probe-home.md`, `4-run-build-scripts-on-kit-copy.md`, `5-grade-named-check-on-private-root.md`, `6-guard-live-checkout.md` | Parallel tests share no file, and the gate reds a test that writes the live checkout. | TD20, TD21, TD22, TD23, TD24, TD25, TD26, TD27, TD28, TD29, TD30, TD31, TD32, TD47, TD48 | `bench test --package ./internal/conformance`, `bench test --package ./internal/gittest`, `bench test --package ./cmd/bench`, `bench test --package ./internal/testreport`, `bench test --package ./internal/gate`, `bench gate` | yes |
| TD-C3 / `7-switch-verdict-windows.md`, `8-name-every-production-wait.md` | No verdict bound expires in a kit test run unless the test set it, a subprocess child keeps its bounds, and every production wait names its window. | TD33, TD34, TD35, TD36, TD37, TD38, TD39, TD40, TD41, TD42, TD43, TD44, TD45, TD46 | `bench test --package ./internal/bounds`, `bench test --package ./internal/git`, `bench test --package ./internal/sessioninspect`, `bench test --package ./internal/worktree --run TestListCommandRendersBoundExpiryAsTypedFailure`, `bench test --check bounds-policy`, `bench test --check system`, `bench gate` | yes |

Chunk IDs keep their pass-1 values. TD-C2 gains ticket 4. The pass-1 tickets 4, 5, 6, and 7 are now tickets 5, 6, 7, and 8.

## Testing decisions

- A good test drives a runner's real entry and observes the child's environment through a probe script that the child runs. The probe prints the first violated property and exits nonzero. The maintenance probe in the shared git test scaffold is the precedent. `TestKitPhaseGitStartsNoAutoMaintenance`, `TestKitRunEnvironmentGitStartsNoAutoMaintenance`, and `TestExternalPhaseGitStartsNoAutoMaintenance` are the three runner precedents.
- The owner gets unit tests in the environment package for each entry, each refusal, and close.
- The guard gets fixture-phase tests through the gate's phase command. The manifest fixture and `runFixturePhases` are the precedent.
- The switch gets unit tests in the bounds package. The bounds-policy rule gets canaries in the `package-core-guard` family. `bounds-duplicate-owner` is the precedent.
- The subprocess switch removal is observed at the system suite's hook child, through `TestSessionStartTE15BoundsDiscoveryAndContinues` with the switch on in the test process.
- The gate observes the whole feature: the landing gate runs every kit test phase under the new environment, the switch, and the guard.

### Posture change: fixtures that the new behavior reds

- `fixturePhaseRoot` builds a root with no git directory. The guard reds that root as unreadable. So the ticket that adds the guard gives the helper a git directory. Its callers in the gate package are `TestFixturePhaseRedRunExitsOne`, `TestFixturePhaseGreenRunExitsZero`, `TestFixturePhaseCancelledRunExitsOneHundredThirty`, `TestPhaseElapsedCellEqualsItsFinishRecord`, `TestFixturePhaseLinesReachTheRunsStreamFile`, and `TestKitPhaseGitStartsNoAutoMaintenance`.
- The hermetic census run on 2026-09-19 found no new red. It ran the whole module with a private `HOME` and `TMPDIR`, pinned Go caches, and no global or system git configuration. The same two `internal/adopt` setup tests failed with and without the emulation. The `internal/adopt` package took 27.7 seconds against 13.2 seconds, because its nested builds derived a cold Bench cache from the private `HOME`.

### Switch census

The author read each test below on 2026-09-19. The census sweep searched every test file for a comparison against a verdict window and for a blocking stub without a setter.

| test | what the switch does to it | disposition |
| --- | --- | --- |
| `TestSessionStartTE15BoundsDiscoveryAndContinues` in the system suite | The hook runs out of process. It asserts that discovery stops near the 2-second discovery window. A setter cannot reach the child. A child that inherits the switch runs discovery with no window, and the test's outer kill deadline ends the hook. | The system suite removes the switch from each child (TD45). |
| `TestSessionStartTE16KillsDiscoveryDescendants` and the other session-start rows | The fixture's run derives its kill deadline from the discovery window. | The same removal covers them. |
| `awaitStartedSpan` in `TestOtelCrashKeepsStartedPhaseLine` | Its wait window derives from the worktree list window of the child. The window stays finite, but a switched child has no inner bound for the window to contain. | The same removal restores the inner bound. |
| `awaitArtifactBarrier` in `TestProspectiveArtifactRecoveryAfterKilledLanding` | The same shape as `awaitStartedSpan`. | The same removal restores the inner bound. |
| `TestEnvironmentPhaseTE15StopsAtDiscoveryBound` in session inspection | It runs in process. It asserts that discovery stops near the discovery window, but session inspection reads the constant, so no setter exists. | Ticket 7 gives session inspection a variable and a setter, and the test sets its own window. |
| `TestWorktreeListTimeoutDefaultUsesPolicy` in the git package | It asserts that the variable equals the policy constant. That equality is false under the switch. | Ticket 7 changes the expected value to the first accessor's result. |
| `TestUpdateRefusesALockAnotherWriterHolds` in the handoff store | It relies on the two-second lock window to expire. | Ticket 8 gives the test its own window through a setter. |

The capture transaction has no lock-held test, so its lock wait needs no test change.

One candidate is not in the census: the refusal test of the consumers command at line 157. Its wait is a test-side window over a constant, and the consumers package has no production bound. The switch does not reach it.

### Seam diagram

    trigger: `bench gate`, `bench test`, release preflight
        │
        ▼
    runner base env ──▶ [ kit test run: open, entries, close ] ──▶ Go child env
                              ◀ tests attach here: owner unit tests; a probe
                                script run as the child of each runner
        │
        ▼
    gate phase command ──▶ [ checkout state before, kit phases, state after ] ──▶ verdict rows
                              ◀ tests attach here: manifest fixture phases that
                                write into the fixture root
        │
        ▼
    production wait ──▶ [ bounds window accessor, switch ] ──▶ deadline or none
                              ◀ tests attach here: bounds unit tests; the
                                bounds-policy canaries; the system suite's
                                hook child with the switch on

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| TD1 | 1 | Under the run's entries, `git config --global --get probe.marker` exits 1 while the base environment's `GIT_CONFIG_GLOBAL` file sets `probe.marker` | planned owner test in internal/env | An owner that drops the global entry leaves the marker readable, so the command exits 0. |
| TD2 | 2 | Under the run's entries, `git config --get probe.system` exits 1 while the base environment's `GIT_CONFIG_SYSTEM` file sets `probe.system` | planned owner test in internal/env | An owner that drops `GIT_CONFIG_NOSYSTEM` leaves the system marker readable. |
| TD3 | 3 | The maintenance probe exits 0 under the run's entries | planned owner test in internal/env | An owner that drops the maintenance entries lets the commit start auto-maintenance, and the probe exits 3. |
| TD4 | 4 | The `HOME` entry names an existing empty directory inside the run directory, and it differs from the base `HOME` | planned owner test in internal/env | An owner that keeps the operator's `HOME` fails the inside-the-run-directory check. |
| TD5 | 5 | The `TMPDIR` entry names an existing empty directory inside the run directory, and it differs from the base `TMPDIR` | planned owner test in internal/env | An owner that keeps the operator's `TMPDIR` fails the inside-the-run-directory check. |
| TD6 | 6 | `go env GOMODCACHE GOPATH GOENV` prints the same three lines under the run's entries as under the base environment | planned owner test in internal/env | A `HOME` swap without pins moves all three paths under the private home. |
| TD7 | 7 | The `bench test` child environment carries `GOCACHE` equal to the Bench cache that the operator's `HOME` derives | planned test in internal/testreport, extending `TestKitRunEnvironmentGitStartsNoAutoMaintenance` | A runner that merges the private `HOME` before it derives the cache hands the child a cold cache under the private home. |
| TD8 | 8 | After close, the run directory is absent, although a test left a mode-0 directory that holds a file inside the private `TMPDIR` | planned owner test in internal/env | A plain recursive removal cannot list the mode-0 directory, so the run directory survives. |
| TD9 | 9 | The `TMPDIR` entry is at most 16 bytes longer than the base `TMPDIR` | planned owner test in internal/env | A run directory with a long name or a deep path exceeds the 16-byte bound. |
| TD10 | 10 | Open with no `HOME` in the base environment returns an error that names `HOME`, and the base `TMPDIR` holds no new entry | planned owner test in internal/env | An owner that creates the run directory first leaves an entry, and an owner that accepts the base continues with no error. |
| TD11 | 10 | Open with a relative `HOME` returns an error that names `HOME` | planned owner test in internal/env | An owner that tests only for absence accepts the relative value. |
| TD12 | 11 | Open with a base `TMPDIR` that names a regular file returns an error that names `TMPDIR` | planned owner test in internal/env | An owner that falls back to the system temporary directory returns no error. |
| TD13 | 12 | Open, the probe, and close succeed with a base `TMPDIR` whose path holds a space | planned owner test in internal/env | An owner that splits or unquotes the path fails to create or to remove the run directory. |
| TD14 | 13, 39 | `KitTestEnv` for a linked root returns no entry, and the gate phases of a linked root carry no `HOME` entry | `TestKitTestEnvSkipsALinkedRoot`, extended in internal/gate | A kit-only predicate that is dropped hands the linked root the private `HOME`. |
| TD15 | 14 | A gate kit fixture phase runs the kit-run probe and the gate exits 0 | planned `TestKitPhaseRunsInTheKitTestRun` in internal/gate, through `runFixturePhases` | A gate composition that misses any entry makes the probe exit nonzero, and the gate goes red. |
| TD16 | 14 | The `bench test` child environment for the kit runs the kit-run probe with exit 0 | planned test in internal/testreport | A runner that composes only the git entries fails the probe on `HOME` or `TMPDIR`. |
| TD17 | 14 | A release preflight external phase runs the kit-run probe with exit 0 | planned test in internal/releasepreflight, extending `TestExternalPhaseGitStartsNoAutoMaintenance` | A preflight phase that composes only the git entries fails the probe. |
| TD18 | 15 | No non-test Go file outside the owner's file names `maintenance.auto` | planned `TestGitPolicyHasOneSource` in internal/env | A surviving copy of the git policy in a runner makes the scan find the name. |
| TD19 | 16 | The owner's documentation states that a plain `go test` is outside the policy | review-owned: reviews/test-determinism.md | No executable check can grade the sentence. |
| TD20 | 29 | Two calls of `conformanceSubprocessEnv` give two different Bench home values, and neither value is the fixed name under the run's `TMPDIR` | planned test in internal/conformance | The old fixed name gives both calls the same value. |
| TD21 | 30 | `conformanceSubprocessEnv` with no npm cache in the base environment sets the npm cache to the fixed shared name under `TMPDIR` | planned test in internal/conformance | A change that gives the npm cache a per-call name fails the fixed-name comparison. |
| TD22 | 32 | `TestNamedCheckRunsOnlyRegisteredDevScope` reads no file under the live checkout's git directory | review-owned: reviews/test-determinism.md | The test is the seam, so no second test can grade it. The guard and review catch a return to the live read. |
| TD23 | 33 | A kit fixture phase that creates `stray` in the root makes the gate exit 1 with one failure row that names `stray` | planned `TestCheckoutGuardRedsAnAddedPath` in internal/gate | A missing guard lets the gate exit 0. |
| TD24 | 34 | A kit fixture phase that rewrites the bytes of an existing ignored file makes the gate exit 1 with a row that names the file | planned test in internal/gate | A guard that compares only the path list misses a changed file. |
| TD25 | 34 | A kit fixture phase that removes an existing untracked file makes the gate exit 1 with a row that names the file | planned test in internal/gate | A guard that reports only new paths misses the removal. |
| TD26 | 35 | A kit fixture phase that creates a new Bench-owned file in the root's git directory makes the gate exit 1 with a row that names the file | planned test in internal/gate | A guard that reads only the worktree misses the git directory. |
| TD27 | 36 | A kit fixture phase that writes the conformance timing file of the root makes the gate exit 0 | planned test in internal/gate | A guard with no declared set reds the timing write. |
| TD28 | 36 | A kit fixture phase run under an open gate run log makes the gate exit 0 while the run log and stream grow | planned test in internal/gate, through `beginGateRunLog` | A guard that does not declare the run records reds the gate's own log. |
| TD29 | 37 | A kit fixture phase that creates `b-stray` and `a-stray` makes the gate print the `a-stray` row before the `b-stray` row | planned test in internal/gate | A guard that prints paths in map order fails the order check. |
| TD30 | 38 | A kit fixture root whose git directory is removed during the phase makes the gate exit 1 with one row that names the unreadable checkout state | planned test in internal/gate | A guard that treats a failed read as an empty state reports a pass or a list of removed paths. |
| TD31 | 39 | A linked-root fixture phase that creates `stray` in the root makes the gate exit 0 | planned test in internal/gate | A guard with no kit-only predicate reds the linked root. |
| TD32 | 40 | The landing gate on the reconciled source is green with the guard in the kit phases | the landing gate | A surviving undeclared writer in the suite makes the guard red the landing. |
| TD33 | 17, 21 | With the switch at `1`, the first accessor returns the unbounded value for each of the seven verdict windows that the Verdict bounds decision lists | planned test in internal/bounds | An accessor that ignores the switch returns the policy window. |
| TD34 | 20 | With the switch absent, the first accessor returns the policy window for each of the seven verdict windows | planned test in internal/bounds | An accessor that always returns the unbounded value turns off production fail-safes. |
| TD35 | 19 | With the switch at `0`, `true`, ` 1`, or `1` plus a newline, the first accessor returns the policy window | planned test in internal/bounds | An accessor that parses a boolean or trims the value accepts a stray value. |
| TD36 | 22 | `bounds.Run` with the unbounded value runs a child that sleeps 200 ms to the complete status, and `bounds.Context` with it returns a context with no deadline | planned test in internal/bounds | A sentinel that reaches a timeout call as a duration expires at once or sets a deadline. |
| TD37 | 18 | `TestListCommandRendersBoundExpiryAsTypedFailure` passes with the switch at `1` in its process | existing test in internal/worktree, run under the switch | A setter that routes through the first accessor becomes unbounded, and the blocking stub hangs the test. |
| TD38 | 18 | A test-set gate timeout of 100 ms still expires with the switch at `1` | existing `TestGateRunTimeoutInvalidatesOldEvidence` in internal/gate, run under the switch | A gate timeout that reads the accessor at use time ignores the assigned window. |
| TD39 | 21 | The kit test run's entries carry `BENCH_TEST_UNBOUNDED_WAITS=1` | planned owner test in internal/env | An owner without the switch leaves the verdict bounds on in every runner. |
| TD40 | 23 | The bounds-policy check reds a production file that passes a raw local duration to `context.WithTimeout` | new canary in the `package-core-guard` family that restores the local ref-check window in internal/git | A check without the new rule stays green on the mutation. |
| TD41 | 23 | The bounds-policy check reds a production file that computes a deadline from the current time with a raw local duration | new canary in the `package-core-guard` family that restores the local handoff lock deadline | A rule that reads only call arguments misses the deadline sum. |
| TD42 | 23 | The bounds-policy check reds an owner-table file whose verdict window variable does not initialize through the first accessor | new canary in the `package-core-guard` family that removes the accessor from the worktree list window | An owner table that still asks only for the raw entry name accepts the unswitched window. |
| TD43 | 24 | With the switch at `1`, the second accessor returns its argument for the gate's process-group cancel grace | planned test in internal/bounds | A grace that routes through the first accessor becomes unbounded, and a cancelled child is never killed. |
| TD44 | 25, 26 | The landing gate on the reconciled source is green with the switch on in every kit test phase | the landing gate | A test that relies on a policy window to expire hangs until Go's test timeout, and the gate goes red. |
| TD45 | 27 | `TestSessionStartTE15BoundsDiscoveryAndContinues` passes with `BENCH_TEST_UNBOUNDED_WAITS=1` in the test process | existing test in internal/systemtest, at the hook child, run through `bench test --check system` under the kit test run | A hook child that inherits the switch runs discovery with no window, so the hook passes the discovery window and the elapsed check is red. |
| TD46 | 28 | The bounds-policy check reds a session inspection file whose discovery window reads `bounds.EnvironmentDiscoveryTimeout` without the first accessor | new canary in the `package-core-guard` family that restores the constant read in internal/sessioninspect | An owner table and a required list without the discovery window accept the constant read. |
| TD47 | 31 | The kit copy of the kit root holds every tracked file and every untracked file that is not ignored, its own git directory, and no `dist/` path | planned test in internal/gittest | A copy that uses the live root or copies ignored files holds the live `dist/` or shares the live git directory. |
| TD48 | 31 | After `TestReleasePreflightBuildDoesNotRebindThePromotionBroker` and `TestGoBuildSubjectModePublishesTheStampedVersion` run, the live checkout holds no new `dist/` path and an unchanged `bin/bench-broker.manifest` | review-owned: reviews/test-determinism.md, with the recorded probe run; the landing gate's guard (TD32) after ticket 6 | A test that still runs the script in the live root recreates `dist/bench-preflight`, and a restore rewrites the live manifest. |

Not covered: story 41 — the reviewed exclusion changes no behavior, and the Won't handle line below records it.
Not covered: story 42 — the reviewed exclusion changes no behavior, and the Won't handle line below records it.

### Edge inventory

The canonical edge classes and the profile's hostile-input checklist, walked at each seam:

- Paths with spaces: TD13 covers a `TMPDIR` with a space. Glob characters reach no shell, because the owner passes each path as one environment value.
- Control bytes in git-sourced text: the guard prints changed paths. A path with a control byte passes through the row renderer, which refuses or escapes it. TD23 to TD30 use plain names.
- Absent versus empty: TD10 covers an absent `HOME`, and TD11 covers a relative `HOME`. An empty `HOME` is relative, so TD11's predicate covers it.
- Special files: TD12 covers a `TMPDIR` that is a regular file. A FIFO or a device at that path fails the same directory creation.
- A dangling symlink as `TMPDIR` fails the same directory creation as TD12.
- A command whose write changes a fact it reports: the guard excludes the gate's own run records through TD28.
- Concurrent runs: each run creates a unique run directory, so two gate runs never share one. The gate execution lock already refuses a second gate run in one repository.
- A test that swaps a package variable: TD37 and TD38 are the setter precedents. The switch census lists every test that the switch reaches.
- A child process that inherits the switch: TD45 covers it.
- Fail-closed cleanup: TD8 covers a mode-0 directory. A removal failure that needs host privileges has no reachable fixture.

**Won't handle** lines:

- Record-age windows: the verdict freshness window (60 minutes), the assignment staleness window, and the operator's wall limit. Go's default package test timeout of 10 minutes ends a kit test run before any of them can expire. The in-scope verdict bounds of TD33 survive.
- The lease staleness window of one minute. It decides only a lease with no owner PID. Process liveness judges a lease with a PID. Flagged for review. The in-scope verdict bounds of TD33 survive.
- A plain `go test` — Go has no module-wide test entry, so no single owner can open a run for it. TD19 records the gap, and the three runners of TD15 to TD17 survive.
- The checkout guard in `bench test` and the preflight race phase — the gate is the oracle, and the gate's kit phases run the same tests. The gate guard of TD23 survives.
- A run-directory removal failure that needs host privileges — no fixture reaches it without root. TD8 survives.
- Second-truncated timestamps — the one ordering comparison uses two truncated times, and no defect exists. The verdict bounds of TD33 survive.
- The causes of CI families B and C — they stay unconfirmed. This spec removes classes of defect, and the family B diagnostic line gives the next evidence.

## Ownership fences

- `reviews/test-determinism.md`
- `internal/env/`
- `internal/gate/phases.go`
- `internal/gate/phases_test.go`
- `internal/gate/checkout_guard.go`
- `internal/gate/checkout_guard_test.go`
- `internal/gate/runner.go`
- `internal/gate/gate.go`
- `internal/gittest/gittest.go`
- `internal/gittest/gittest_test.go`
- `internal/testreport/environment.go`
- `internal/testreport/environment_test.go`
- `internal/testreport/check_test.go`
- `internal/testreport/command.go`
- `internal/releasepreflight/command.go`
- `internal/releasepreflight/external_test.go`
- `internal/conformance/checks_test.go`
- `internal/conformance/conformance_env_test.go`
- `internal/conformance/bounds_policy_test.go`
- `internal/conformance/fixture_bite_test.go`
- `tests/canary/package-core-guard/`
- `internal/bounds/bounds.go`
- `internal/bounds/bounds_test.go`
- `internal/git/git.go`
- `internal/git/worktree_admin_enum_test.go`
- `internal/handoffdoc/`
- `internal/capturetx/`
- `internal/runbinary/runbinary.go`
- `internal/worktree/subshell.go`
- `internal/worktree/exec.go`
- `internal/freshness/freshness_publish.go`
- `internal/releaseevidence/release_evidence.go`
- `internal/shift/loop.go`
- `internal/models/models.go`
- `internal/sessioninspect/sessioninspect.go`
- `internal/sessioninspect/sessioninspect_test.go`
- `internal/systemtest/owner_test.go`
- `internal/systemtest/owner_land_race_test.go`
- `internal/systemtest/owner_artifact_recovery_test.go`
- `internal/guards/guards.go`
- `internal/coverage/citation_execution.go`
- `internal/refresh/refresh.go`
- `cmd/bench/build_subject_mode_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`

Reviewer disposition: pending.

## Ticket graph

| ticket | Blocked by | chunk |
| --- | --- | --- |
| `1-open-kit-test-run.md` | none | TD-C1a |
| `2-compose-kit-test-run.md` | `1-open-kit-test-run.md` | TD-C1b |
| `3-isolate-conformance-probe-home.md` | none | TD-C2 |
| `4-run-build-scripts-on-kit-copy.md` | none | TD-C2 |
| `5-grade-named-check-on-private-root.md` | `4-run-build-scripts-on-kit-copy.md` | TD-C2 |
| `6-guard-live-checkout.md` | `1-open-kit-test-run.md`, `4-run-build-scripts-on-kit-copy.md` | TD-C2 |
| `7-switch-verdict-windows.md` | `2-compose-kit-test-run.md` | TD-C3 |
| `8-name-every-production-wait.md` | `7-switch-verdict-windows.md` | TD-C3 |

## Out of scope

- A nightly stress job with a CPU quota, `-shuffle=on`, and `-cpu=1,2` — about 3 edits and 1 gate run. It is parked in the idea inbox.
- A shared Bench build cache for nested builds under the private `HOME` — about 4 edits and 1 gate run. The census measured its cost in one package only.

## Further notes

### Source trace

| source sentence | rows or exclusion |
| --- | --- |
| Tests use no wall-clock limit as a correctness signal. | TD33 to TD46 |
| The seam lets tests supply the clock and the limits for the gate's timeouts. | TD33, TD37, TD38 |
| The seam covers freshness windows. | Won't handle: record-age windows |
| The seam covers second-truncated timestamps. | Won't handle: second-truncated timestamps |
| The seam covers git helper bounds. | TD33, TD37, TD40 |
| One hermetic test-environment owner extends the git test configuration. | TD1 to TD13, TD18 |
| The owner sets `GIT_CONFIG_GLOBAL`, `GIT_CONFIG_NOSYSTEM`, `HOME`, and `TMPDIR`, and keeps the maintenance policy. | TD1, TD2, TD3, TD4, TD5 |
| The same three runners compose the owner. | TD15, TD16, TD17 |
| Each test writes only under its own temporary directory. | TD20, TD22, TD23 to TD32, TD47, TD48 |
| The live checkout, `.logs`, timing records, and the Bench cache home are the named suspects. | TD4, TD22, TD23, TD27, TD28, TD48 |
| Linked repositories keep their own environment. | TD14, TD31 |
| A plain `go test` stays outside unless an honest single owner exists. | TD19, Won't handle: plain `go test` |

### Reader sweep and proof checklist

- Cited symbols: each symbol below resolves in the tree at `1b030b1c`. `GitTestConfig` resolves, and ticket 2 retires it.
  - Environment and gate: `KitTestEnv`, `withKitTestEnv`, `phasesCommandAtKitWithSelection`, `runFixturePhases`, `fixturePhaseRoot`, `beginGateRunLog`, `gateRunStreamFile`, and `inheritGateRunLog`.
  - Runners and conformance: `selectedRunEnvironment`, `testEnvironment`, `runExternal`, `conformanceSubprocessEnv`, `RunConformanceSelection`, `registry.TimingPath`, and `registry.NewTimingWriter`.
  - Caches and bounds: `gocache.Apply`, `gocache.Dir`, `bounds.Run`, `bounds.RunOutput`, `bounds.Context`, `bounds.ContextCause`, and `bounds.TestDeadline`.
  - Bounds policy: `SetWorktreeListTimeoutForTest`, `checkBoundsPolicy`, `checkBoundCaller`, and `processGroupGrace`.
  - System suite: `childEnvironment`, `mergeEnvironment`, `systemStartSelected`, `startArtifactLand`, `awaitStartedSpan`, `awaitArtifactBarrier`, and `preserveWrapperManifest`.
  - Tests: `TestSessionStartTE15BoundsDiscoveryAndContinues`, `TestWorktreeListTimeoutDefaultUsesPolicy`, `TestReleasePreflightBuildDoesNotRebindThePromotionBroker`, and `TestGoBuildSubjectModePublishesTheStampedVersion`.
- Import edges: the environment package imports the bounds package for the switch name (new edge; the bounds package imports only the standard library). The gate, the testreport, and the release preflight packages already import the environment package. The system suite already imports the bounds package.
- Source-row clauses and occurrences: the `maintenance.auto` policy occurs once in production code, in the environment package. The fixed temporary names occur once each, in `conformanceSubprocessEnv`. The timing file path has one owner, `registry.TimingPath`. The system suite builds a child environment from the process environment at three sites: `childEnvironment`, `systemStartSelected`, and `startArtifactLand`.
- Promised field labels: none. The guard rows reuse the gate's existing `failures[N]{phase,line}` table.
- Changed-function callers: `GitTestConfig` has two callers, `KitTestEnv` and `runExternal`. `KitTestEnv` has two callers, `withKitTestEnv` and `selectedRunEnvironment`. `selectedRunEnvironment` has four callers in the testreport package. `conformanceSubprocessEnv` has callers in the conformance package only. `preserveWrapperManifest` has two callers, the two build-script tests. Each verdict window variable in the owner table has one initializer.
- Copy survival: TD18 reds a surviving git policy copy. TD42 and TD46 red a verdict window that skips the accessor.

### Hostile edges from the census

The census on 2026-09-19 ran the whole module in this worktree and compared the checkout state before and after. The plain run created `dist/` in the checkout and wrote the conformance timing file in the git directory. The timing file is declared. The census compared path lists only for the checkout, so it cannot see a rewrite of an existing ignored file. Ticket 6's first gate run with the guard is the complete census.

The author confirmed the `dist/` writer on 2026-09-19. With `dist/` removed, the author ran the freshness package alone, and `dist/` stayed absent. The author then ran `TestReleasePreflightBuildDoesNotRebindThePromotionBroker` alone, and `dist/` came back with `bench-preflight`, `bench-preflight.seal`, and `bench-broker.manifest`. That test and `TestGoBuildSubjectModePublishesTheStampedVersion` also rewrite the live `bin/bench-broker.manifest` through `preserveWrapperManifest`.

The hermetic run left `.cache/bench` and `.config/go` in the private home, and `bench-npm-cache` and `bench-conformance-home` in the private temporary directory.

### Fence disposition

The command registry, its two tests, and the two conformance registry tests join the fence through the binding registry closure only. The guards, coverage, and worktree packages, and the `cmd/bench` package, are bound packages. The build expects no edit to those five files.

### Collision with a staged spec

The staged `ft290-test-projection` spec writes all of `internal/testreport/`, because each of its tickets writes there. This spec writes that package in ticket 2 (`environment.go` and its test), ticket 5 (`check_test.go`), and ticket 8 (`command.go`). The ft290 ticket 1 still plans a move of the named-check owner that `f877c158` already landed. The reviewer decides which spec builds first. The second build composes the first build's landed tip with `bench worktree merge --from`. The workflow rejects a rebase.

### Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"TD-C1a","tickets":["1-open-kit-test-run.md"],"verification":[{"id":"env","command":"bench test --package ./internal/env"},{"id":"gate","command":"bench test --package ./internal/gate"}]},{"id":"TD-C1b","tickets":["2-compose-kit-test-run.md"],"verification":[{"id":"testreport","command":"bench test --package ./internal/testreport"},{"id":"releasepreflight","command":"bench test --package ./internal/releasepreflight"},{"id":"env","command":"bench test --package ./internal/env"}]},{"id":"TD-C2","tickets":["3-isolate-conformance-probe-home.md","4-run-build-scripts-on-kit-copy.md","5-grade-named-check-on-private-root.md","6-guard-live-checkout.md"],"verification":[{"id":"conformance","command":"bench test --package ./internal/conformance"},{"id":"gittest","command":"bench test --package ./internal/gittest"},{"id":"cmd","command":"bench test --package ./cmd/bench"},{"id":"testreport","command":"bench test --package ./internal/testreport"},{"id":"gate","command":"bench test --package ./internal/gate"}]},{"id":"TD-C3","tickets":["7-switch-verdict-windows.md","8-name-every-production-wait.md"],"verification":[{"id":"bounds","command":"bench test --package ./internal/bounds"},{"id":"git","command":"bench test --package ./internal/git"},{"id":"sessioninspect","command":"bench test --package ./internal/sessioninspect"},{"id":"worktree-bound","command":"bench test --package ./internal/worktree --run TestListCommandRendersBoundExpiryAsTypedFailure"},{"id":"bounds-policy","command":"bench test --check bounds-policy"},{"id":"system","command":"bench test --check system"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/test-determinism/spec.md"},{"id":"env","command":"bench test --package ./internal/env"},{"id":"gate","command":"bench test --package ./internal/gate"},{"id":"bounds","command":"bench test --package ./internal/bounds"},{"id":"bounds-policy","command":"bench test --check bounds-policy"},{"id":"system","command":"bench test --check system"}]}
```

### Flagged additions

- The 16-byte bound on the private `TMPDIR` (story 9, TD9) is a number that the spec chose. The source named `TMPDIR` but gave no bound. The bound protects socket-path tests from new capability skips.
- The second window accessor (story 24, TD43) is a mechanism that the source did not name. It keeps graces, polls, and the operator's wall limit fixed under the switch.
- The guard's red on an unreadable checkout state (story 38, TD30) follows the fail-closed rule for a new standing check.
- The kit copy helper (story 31, TD47) comes from review finding 5. It is the one owner of a private kit tree for tickets 4 and 5.
