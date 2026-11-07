# A retired protection store leaves with its resource, and inherited registration has one owner

Status: staged

Roadmap: FT446

Decision source: `roadmap/FT446.md`, the committed outcome FT446-cache-store that the reviewer approved, with the criterion "A retired cache store leaves with its cache, and inherited registration has one owner" (named reviewed artifact).

Verification log: 4 iteration(s) to accept — the spec round and the ticket round each took two iterations; the spec round fixed test isolation, a misstated leak path, and missing pre-review proofs, and the ticket round split the seam chunks and the copy rule.

## Problem

A protection store retires when its last user resolves.
The retirement deletes `resource.json`, which is the only record that names the resource path.
The store directory, its lock file, and its empty users directory stay.
The nested-store sweep reads `resource.json`, so it skips each retired store.
A system run therefore leaves about eight stores in `/dev/shm`, and they stay until a reboot.
A package test that deletes its cache directory also leaves the store of that cache.

The bounds runner and the build-script runner each read the inherited resource vector and build an inherited registration by hand.
`subprocess.RunInherited` owns that same fact, but it always uses the installed protocol.
The build-script runner takes its protocol from its caller, so it cannot call `RunInherited` today.

## Solution

A retired store keeps one small retired record that names its kind, its resource path, and its binding.
`resource.json` stays absent, so every current reader of a store with no generation keeps its answer.
The nested-store sweep and a new absent-resource check read the retired record.
Each deletes a retired store whose resource no longer exists and that holds no user record.
The last release of a cache hold runs the absent-resource check over the cache stores.

`RunInherited` takes the caller's registration.
A nil protocol in that registration means the installed protocol.
The bounds runner and the build-script runner launch through `RunInherited`, and neither builds an inherited registration.

## User stories

Line: opus / high.
Implementation-line reason: RS-C1 and RS-C2 are the hardest chunks, because they change a fail-closed store that every protected launch shares. The spec settles each design point, the seams are existing owners, and package tests drive each row through real stores.
Harder chunks: RS-C1, RS-C2.

Retired store record

1. As a test author, I want a retired store to record the path of its resource, so that a later sweep can find the store after the resource is gone.
2. As a resource owner, I want a store whose generation retired to stay a store with no generation for every current reader, so that recovery, cache clean, and checkout cleanup keep their answers.
3. As a cache user, I want the next hold on a retired store to open a fresh generation, so that a retired record never blocks a new launch.
4. As a maintainer, I want the retired record read under the same strict rules as every store record, so that a malformed, linked, or oversized record keeps its store.

Nested sweep

5. As a kit runner, I want the removal of a run directory to delete the retired store of each resource that the removal deleted, so that no store of the run stays in temporary storage.
6. As a kit runner, I want a retired store whose resource still exists to stay, so that a resource that is still on disk keeps its store.
7. As a kit runner, I want a retired store of a sibling path that only shares the prefix to stay, so that the sweep deletes nothing outside the removed directory.
8. As the user of a nested resource, I want a store that holds a user record never deleted by the retired branch, so that a live user keeps its protection.

Absent-resource check

9. As a test author, I want the last release of a cache hold to delete each retired cache store whose cache directory is gone, so that a cache that its test deleted leaves no store.
10. As a cache user, I want a retired cache store whose directory still exists to stay after a release, so that a reusable cache keeps its store.
11. As a cache user, I want a contended release to delete no store, so that another holder keeps its generation.
12. As a maintainer, I want an unreadable retired record to keep its store without a release error, so that one damaged store does not fail every release.
13. As a maintainer, I want a failed deletion to join the release error, so that a store that cannot go is reported.
14. As a reviewer, I want the protection ADR to state that a retired store leaves with its resource, so that the decided state matches the code.
15. As a reader of the data-handling record, I want the retired record and its lifetime named, so that the record list stays complete.

Inherited registration owner

16. As a maintainer, I want `RunInherited` to register a launch on the protocol that its caller names, with the installed protocol for a nil protocol, so that the build-script runner can call it.
17. As a maintainer, I want the bounds runner and the build-script runner to launch through `RunInherited`, so that inherited registration has one owner.
18. As a caller of the bounds runner, I want a malformed inherited vector to refuse before a start as a start failure, so that an unreadable vector never starts an unprotected child.
19. As a caller of the bounds runner, I want a cancelled parent to stay a cancellation and an expired limit to stay a timeout, both with a KILL, so that the result keeps its current meaning.
20. As a caller of the build-script runner, I want its TERM, its grace, its local resources, and its build error output kept, so that a build keeps its current behavior.
21. As a caller of the build-script runner, I want a malformed inherited vector to refuse the build before it starts, so that the build never runs unprotected.
22. As a reviewer, I want the process-lifetime census to keep both runner entries green, so that the census still sees every lifetime owner.

## Implementation decisions

**D1. The retired store keeps a separate retired record, and `resource.json` stays absent.**
The protection package adds `retired.json` to its store file set, with one strict record type: schema, kind, resource path, and binding.
Reason: four readers key on an absent `resource.json` through `NoGeneration`, and the store join publishes a fresh generation on an absent `resource.json`.
A third admission value in `resource.json` changes each of those readers, the record validator, the registration admission check, and the store observer.
The retired record uses the private mode, the no-follow read, the read limit, and the exact decode of every store record.
Its validator follows the `resourceRecord.valid` pattern and reuses the descriptor's owner checks for the kind, the path, and the binding.

These readers change:

- `SealAndObserve` publishes the retired record after it removes each user and before it removes `resource.json`. A crash between the two steps leaves a sealed generation with no user, which the next recovery retires again.
- `discardNested` reads the retired record of a store whose `resource.json` is absent.
- A new kind-scoped absent-resource check on `Directory` reads the retired record of each store of its kind.

These readers stay unchanged, and their current tests grade them:

- `Recover` and `NoGeneration`: `resource.json` is still absent, so a retired store still reads as a store with no generation.
- The store join: an absent `resource.json` still publishes a fresh generation. A later retirement replaces the retired record.
- `ObserveRecorded`: it still skips a store with no `resource.json`.
- The cache clean, the checkout cleanup, and the process-lifetime artifact pins in the system tests.
- The checkout cleanup evidence: it walks every store file, and a retirement already changed that set when it removed `resource.json`.

**D2. One retired-store predicate serves both deleters.**
A store is deletable through the retired branch only when all of these hold under a non-blocking store lock:

- `resource.json` is absent;
- the retired record reads under the strict rules;
- the users directory holds no entry;
- the recorded resource path is absent.

`discardNested` adds the condition that the recorded path is inside the removed directory.
Before it locks, the predicate skips a store with a `resource.json` entry or with no `retired.json` entry, so it locks only a retired store.
The store lock helper gains a non-blocking form. A contended store lock skips the store with no error, so a sweep never waits on another process.
A store that the predicate cannot lock or read counts as unreadable: it stays, and it adds no error.
The live-generation branch of `discardNested` stays as it is.
Reason: the retired branch must never delete a store that a user record names, and one predicate keeps the two deleters equal.

**D3. The absent-resource check runs at the last release of a cache hold.**
After the hold's own `SealAndObserve` returns without an error, the release runs the check over every store of the cache kind.
A contended release returns before it, as it does today.
A store that the check cannot lock or read stays and adds no error.
A store root that the check cannot list, or a deletion that fails, joins the release error.
The cache owner calls the check through one package variable, `discardRetiredCaches`, so a cache test can record the call.
The injected-port registry admits it as one `seamVariable` row: package `internal/gocache`, variable `discardRetiredCaches`, owner function `heldLock.retire`, purpose "a recording check shows that only the last uncontended release runs the absent-resource check, after its own retirement", `serialIsolation`, and consumer `internal/gocache/protection_test.go::TestReleaseRunsTheRetiredCheck`.
That consumer does not call `t.Parallel`.
Reason: the kit run of a system test holds each cache under its run directory, and the kit run's protected removal deletes that directory before the outer cache release. The RS-C2 sweep therefore removes the system-run leak.
A package test that holds a cache under its own temporary directory deletes that directory with no protected removal. Only a check outside `Generation.Remove` reaches the store of that cache, and the decision source names this check.
The store root is shared by every HOME of the user in the boot. The stores of the last deleted caches stay until the next uncontended cache release by the same user in the boot.
The check tests only the recorded path. A cache that a process creates again at the same path keeps the store of the earlier binding until that cache goes again.

**D4. `RunInherited` takes a `Registration` in place of its variadic local descriptors.**
The registration's resources join the inherited union after the inherited resources, as today.
A nil protocol means the installed protocol.
Reason: the build-script runner already holds a registration with an optional protocol, so it forwards that value unchanged.
The five production callers pass their local descriptors in the registration's resources, and the detached launch passes an empty registration.
The process-lifetime closure check also refuses a hand-built inherited registration in a production function outside `internal/subprocess`.
The rule counts three readers of the process vector: `protection.Inherited()`, `protection.InheritedFrom(os.Environ())`, and `protection.ParseVector(os.Getenv(protection.Env))`.
A function is refused when it calls a counted reader and builds a `subprocess.Registration` literal, either in its own body or in a same-package function that it calls, as the census follows same-package calls.
The rule does not count `InheritedFrom` of another environment, `ParseVector` of another value, or a reader with no `Registration`, so the checkout protection union, the kit run's union, and the release-evidence membership read stay green.
The rule lives in the file diagnostics of the closure check. The closure prefilter parses only files that spell an owner token, so the rule adds its own reader tokens to that prefilter.
Each of the five callers keeps its own signature and wraps its descriptors into the registration at the `RunInherited` call, so their test call sites do not change.
Reason: the census accepts `NewLifetime`, so only this rule makes a surviving copy red.

**D5. Each runner expresses its own cancellation and errors around `RunInherited`.**
The bounds runner keeps its parent check before the call.
It maps any `RunInherited` error to a start failure, which covers the malformed vector and a refused start.
Its cancel request answers an interruption when the parent context ended and a deadline otherwise, each with KILL.
The build-script runner keeps its toolchain check, its environment, and its process attributes before the call.
Its cancel request answers a deadline when the context expired and an interruption otherwise, with TERM.
It keeps its own outcome mapping: the context cause, the start error, the build error with output, and the uncertified outcome.
Reason: `RunInherited` already returns the outcome and leaves the mapping to each caller.
When the environment and the vector both fail, the build-script runner now reports the environment error first. Both refuse before the start.

**D6. The protection ADR records the decided state.**
Its consequence states that a retired store keeps a record of its resource, and that the removal of an enclosing directory or the last release of a cache hold deletes each retired store whose resource no longer exists.
The sentence that keeps a retired store until the next boot goes.

**D7. No new source file goes into a crowded directory.**
The subprocess, protection, run-binary, and cache packages are over the twelve-file default with no grant.
The retired record joins the record owner, the predicate and the check join the directory owner, and each protection, cache, subprocess, and run-binary test joins an existing test file.
The conformance directory holds an accepted grant, so the copy-rule fixtures take one new test file there.
The build-script runner's file shrinks, because its hand-built registration goes.

No intake ledger or review record field changes, so this spec declares no persisted field.

## Implementation chunks

Each chunk is one review checkpoint. RS-C4 has no dependency on RS-C1, RS-C2, or RS-C3.
The ticket review split the chunks: the old RS-C1 is now RS-C1 and RS-C2, the old RS-C2 is now RS-C3, and the old RS-C3 is now RS-C4 with tickets 04 and 05.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| RS-C1 / 01-record-retired-store-path.md | A retired store records its resource, and every reader of a store with no generation keeps its answer | RS1, RS2, RS3, RS25, RS26 | TestRetiredStoreRecordsItsPath, TestRetiredStoreHasNoGeneration, TestProtectionCompletedHold, TestCheckoutProtectionPredicate, the process-lifetime artifact system test | yes |
| RS-C2 / 02-sweep-retired-nested-stores.md | A protected removal deletes the retired stores of the resources that it deleted | RS4, RS5, RS6, RS7, RS8, RS14, RS15, RS31 | TestGenerationRemove, the record tests over each record file, TestRetiredSweepKeepsAUser, TestRetiredSweepSkipsALockedStore | yes |
| RS-C3 / 03-discard-gone-retired-caches.md | The last release of a cache hold deletes each retired cache store whose cache is gone | RS9, RS10, RS11, RS12, RS13, RS27, RS28, RS30 | TestDiscardRetired, TestDiscardRetiredFailsClosed, TestReleaseRunsTheRetiredCheck, the injected-port registry check | no |
| RS-C4 / 04-run-inherited-under-local-protocol.md, 05-refuse-hand-built-inherited-registrations.md | Inherited registration has one owner, both runners launch through it, and the closure check refuses a surviving copy | RS16, RS17, RS18, RS19, RS20, RS21, RS22, RS23, RS24, RS29 | TestRunInheritedRegistersOnLocalProtocol, TestBuildScriptMalformedResources, TestBoundsMalformedResources, TestBoundsCancel, TestBoundsFixedShutdown, TestBuilderCancel, TestBuilderNormalDrain, TestProcessLifetimeCacheSecondProcessClean, TestProcessLifetimeInheritedRegistrationCopy, the process-lifetime closure check | no |

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","review_mode":"unified","run_id":"ft446-build","orchestrator_session":"ft446-orchestrator","author_limit":1,"assignments":{"01-record-retired-store-path.md":[{"session":"ft446-author-01","assignment":"ft446-author-01","model":"opus","effort":"high","source":"9d692cfba956bca51ecb3aefd4e831c6d89108c9","native_ref":"agent:ft446-author-01"}],"02-sweep-retired-nested-stores.md":[{"session":"ft446-author-02","assignment":"ft446-author-02","model":"opus","effort":"high","source":"6d485d5ff5aafea9600312bc389f4da662060405","native_ref":"agent:ft446-author-02"},{"session":"ft446-repair-02-r1","assignment":"ft446-repair-02-r1","model":"opus","effort":"high","source":"482b6d48f011e5006f65ba755eea56717b0bf9d2","native_ref":"agent:ft446-repair-02-r1","predecessor":"ft446-author-02","trigger":"user-directed","stopped":"agent:ft446-author-02","preserved":"bef1597e66827888f6fbdd964ac136d50f035aed"}],"03-discard-gone-retired-caches.md":[{"session":"ft446-author-03","assignment":"ft446-author-03","model":"opus","effort":"medium","source":"d7fd8022d9b20f5c41b27e01d3e5b281bc055906","native_ref":"agent:ft446-author-03"}],"04-run-inherited-under-local-protocol.md":[],"05-refuse-hand-built-inherited-registrations.md":[]}},"chunks":[{"id":"RS-C1","tickets":["01-record-retired-store-path.md"],"verification":[{"id":"retired-record-written","command":"bench test --package ./internal/subprocess/protection","probe":"remove the retired-record publish from SealAndObserve; require TestRetiredStoreRecordsItsPath to fail","ticket":"01-record-retired-store-path.md"},{"id":"cache-no-generation","command":"bench test --package ./internal/gocache","probe":"make Recover return a non-ErrNotExist error when the store holds retired.json; require TestProtectionCompletedHold to fail on its clean after the release","ticket":"01-record-retired-store-path.md"},{"id":"checkout-no-generation","command":"bench test --package ./internal/worktree --run '^TestCheckoutProtectionPredicate$'","probe":"make Recover return a non-ErrNotExist error when the store holds retired.json; require TestCheckoutProtectionPredicate to fail","ticket":"01-record-retired-store-path.md"}]},{"id":"RS-C2","tickets":["02-sweep-retired-nested-stores.md"],"verification":[{"id":"retired-sweep-omission","command":"bench test --package ./internal/subprocess/protection","probe":"make discardNested skip each store with no resource.json, as it does today; require TestGenerationRemove to fail on the retired gone store","ticket":"02-sweep-retired-nested-stores.md"},{"id":"user-guard-omission","command":"bench test --package ./internal/subprocess/protection","probe":"drop the empty-users condition from the retired-store predicate; require TestRetiredSweepKeepsAUser to fail","ticket":"02-sweep-retired-nested-stores.md"},{"id":"lock-skip-omission","command":"bench test --package ./internal/subprocess/protection","probe":"make a contended store lock join the sweep error; require TestRetiredSweepSkipsALockedStore to fail","ticket":"02-sweep-retired-nested-stores.md"}]},{"id":"RS-C3","tickets":["03-discard-gone-retired-caches.md"],"verification":[{"id":"cache-check-omission","command":"bench test --package ./internal/gocache","probe":"remove the absent-resource check call from the cache retire; require TestReleaseRunsTheRetiredCheck to fail","ticket":"03-discard-gone-retired-caches.md"},{"id":"kind-scope-omission","command":"bench test --package ./internal/subprocess/protection","probe":"drop both the store-name prefix filter and the record kind comparison from the absent-resource check; require TestDiscardRetired to fail on the retired store of another kind","ticket":"03-discard-gone-retired-caches.md"},{"id":"registry-row","command":"bench test --package ./internal/conformance","probe":"delete the discardRetiredCaches seamVariable row from the injected-port registry; require the registry check to refuse the unregistered variable","ticket":"03-discard-gone-retired-caches.md"}]},{"id":"RS-C4","tickets":["04-run-inherited-under-local-protocol.md","05-refuse-hand-built-inherited-registrations.md"],"verification":[{"id":"protocol-forwarding","command":"bench test --package ./internal/subprocess","probe":"make RunInherited use the installed protocol whatever the registration names; require TestRunInheritedRegistersOnLocalProtocol to fail","ticket":"04-run-inherited-under-local-protocol.md"},{"id":"build-vector-refusal","command":"bench test --package ./internal/runbinary","probe":"make RunInherited ignore an unreadable inherited vector; require TestBuildScriptMalformedResources to fail","ticket":"04-run-inherited-under-local-protocol.md"},{"id":"runner-callers","command":"bench test --package ./internal/bounds --package ./internal/gate --package ./internal/preprelease --package ./internal/releasepreflight --package ./internal/testreport","probe":"make the bounds cancel request answer a deadline for a cancelled parent; require TestBoundsFixedShutdown to fail","ticket":"04-run-inherited-under-local-protocol.md"},{"id":"copy-survival","command":"bench test --package ./internal/conformance --run '^TestProcessLifetime'","probe":"remove the hand-built inherited registration rule from the closure check; require TestProcessLifetimeInheritedRegistrationCopy to fail","ticket":"05-refuse-hand-built-inherited-registrations.md"}]}],"final_verification":[{"id":"packages","command":"bench test --package ./..."},{"id":"system","command":"bench test --check system"}]}
```

## Testing decisions

- A good test drives a real protection store under a private store base and observes which store directories exist after the operation.
- The protection rows attach to the record owner and the directory owner through `internal/subprocess/protection` tests. `TestGenerationRemove` is the prior art for the sweep.
- The predicate and the kind-scoped check attach to the directory owner through `internal/subprocess/protection` tests on a private store base, under the existing `storeBase` registry row.
- The cache trigger attaches to the cache hold through `internal/gocache` tests that record each call of the check through the admitted package variable. Another process cannot add or remove a call that this process records, so the shared store root cannot make these rows pass or fail. The completed-hold and contended-release tests are the prior art.
- The registration rows attach to `RunInherited` through the start-gate fixtures in `internal/subprocess`. The child environment shows the protocol that registered the launch, because a registration with no protocol leaves the child environment unchanged.
- The runner rows keep their current tests. The process-lifetime closure check and the process-lifetime system tests observe the feature in the gate.
- RS16 sets `BENCH_PROCESS_RESOURCES` to an empty vector with `t.Setenv`, so the launch does not register on the vector of an outer kit run.

### Seam diagram

    trigger: Generation.Remove of a directory, or the last release of a cache hold
        │
        ▼
    store root entries  ──▶  [ retired-store predicate under a non-blocking store lock ]  ──▶  deleted or kept store
                      ◀ tests attach here: real stores under a private store base; os.Stat of each store

    trigger: bounds.Run, runBuildScript, and the five RunInherited callers
        │
        ▼
    command + Registration  ──▶  [ RunInherited ]  ──▶  child environment and outcome
                      ◀ tests attach here: start-gate fixtures read the child environment

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| RS1 | 1 | After a safe retirement, the store holds a retired record whose kind, resource path, and binding equal the retired generation | `internal/subprocess/protection/lifecycle_test.go` (`TestRetiredStoreRecordsItsPath`) | A retirement that writes no retired record leaves the store with no readable path |
| RS2 | 2 | A retired store holds `retired.json`, and `Recover` of that store returns a result that `NoGeneration` accepts | `internal/subprocess/protection/recovery_test.go` (`TestRetiredStoreHasNoGeneration`) | A reader that takes the retired record as a generation makes `NoGeneration` false |
| RS3 | 3 | The next hold on a retired cache store opens a generation with a new resource ID | `internal/gocache/protection_test.go` (`TestProtectionCompletedHold`) | A join that reads the retired record as a published generation refuses the open |
| RS4 | 4 | A retired record that is malformed, a link, at another mode, or over the read limit keeps its store through a sweep | `internal/subprocess/protection/record_test.go` (`TestRecordMalformed`), with the other record tests, run over each record file name through their record planting helpers | A lax reader deletes a store on a forged or damaged path |
| RS5 | 5 | A protected removal deletes the retired store of a nested resource that the removal deleted | `internal/subprocess/protection/lifecycle_test.go` (`TestGenerationRemove`) with a retired gone store | The current sweep skips a store with no `resource.json`, so the retired store stays |
| RS6 | 6 | A retired nested store whose resource still exists stays after the removal | `internal/subprocess/protection/lifecycle_test.go` (`TestGenerationRemove`) with a retired kept store | A sweep that tests only the prefix deletes a store whose resource is on disk |
| RS7 | 7 | A retired store of a sibling path that shares the prefix stays after the removal | `internal/subprocess/protection/lifecycle_test.go` (`TestGenerationRemove`) with a retired sibling store | A prefix test without the separator deletes a store outside the removed directory |
| RS31 | 12 | A retired store whose lock another open file holds stays after a sweep, and the sweep returns no error for it | `internal/subprocess/protection/lifecycle_test.go` (`TestRetiredSweepSkipsALockedStore`) | A blocking lock stalls the sweep, and a lock failure that joins the error fails the removal |
| RS8 | 8 | A retired store whose users directory holds an entry stays after a sweep | `internal/subprocess/protection/lifecycle_test.go` (`TestRetiredSweepKeepsAUser`) | A predicate without the empty-users condition deletes a store that a user names |
| RS9 | 9 | The absent-resource check deletes a retired cache store whose cache directory is gone | `internal/subprocess/protection/directory_test.go` (`TestDiscardRetired`) | Without the check, the retired store stays |
| RS10 | 10 | The absent-resource check keeps a retired cache store whose cache directory exists | `internal/subprocess/protection/directory_test.go` (`TestDiscardRetired`) with a kept cache | A check that deletes every retired store removes a reusable cache's store |
| RS30 | 10 | The cache-kind absent-resource check keeps a retired store of another kind whose resource is gone | `internal/subprocess/protection/directory_test.go` (`TestDiscardRetired`) with a gone binary-kind store | A check without the kind prefix test deletes the store of another kind |
| RS11 | 11 | A contended release makes no call of the absent-resource check | `internal/gocache/protection_test.go` (`TestReleaseRunsTheRetiredCheck`) with a second holder process | A check that runs before the contended return records a call |
| RS12 | 12 | An unreadable retired record keeps its store, and the check returns no error for it | `internal/subprocess/protection/directory_test.go` (`TestDiscardRetiredFailsClosed`) | A check that fails on one damaged store makes every release fail |
| RS13 | 13 | A deletion that fails joins the error that the check returns | `internal/subprocess/protection/directory_test.go` (`TestDiscardRetiredFailsClosed`) with a 0700 store that passes the predicate and holds a non-empty 0500 subdirectory; the test skips through the privilege class of `internal/capability` when that subdirectory stays writable, as the gate's write-denied guard does | A check that drops the deletion error hides a store that cannot go |
| RS14 | 14 | The protection ADR consequence states that a retired store keeps a record of its resource and that a protected removal deletes it with its resource | review-owned: docs/adr/0030 consequence | No test reads ADR prose |
| RS15 | 15 | The data-handling record names the retired record and states that it stays until its resource is gone | review-owned: DATA_HANDLING.md process resource protection section | No test reads that section |
| RS16 | 16 | `RunInherited` registers the launch on the protocol that its registration names | planned TestRunInheritedRegistersOnLocalProtocol in internal/subprocess/startgate_test.go with an empty inherited vector | A `RunInherited` that uses the installed protocol, which tests leave nil, registers nothing, so the child environment lacks the resource |
| RS17 | 16 | A launch through `RunInherited` with a nil protocol registers on the installed protocol | `internal/systemtest/process_lifetime_cache_test.go` (`TestProcessLifetimeCacheSecondProcessClean`) | A nil protocol that registers nothing leaves the gate's launches unprotected, and those tests require a registered user |
| RS18 | 17 | The closure check names each production function outside `internal/subprocess` that reads a counted process-vector reader and builds a `subprocess.Registration` | planned TestProcessLifetimeInheritedRegistrationCopy in internal/conformance/process_lifetime_copy_test.go with one synthetic fixture per counted reader, one with the `Registration` in a same-package helper, and one uncounted reader that stays green | The census accepts `NewLifetime`, so without this rule a surviving hand-built copy stays green |
| RS19 | 18 | A malformed inherited vector makes the bounds runner return a start failure before any child starts | `internal/bounds/lifetime_test.go` (`TestBoundsMalformedResources`) | A runner that starts before the vector read returns another status or a started command |
| RS20 | 19 | A cancelled parent makes the bounds runner return a cancellation after a KILL | `internal/bounds/lifetime_test.go` (`TestBoundsCancel`) | A cancel request that always answers a deadline returns a timeout |
| RS21 | 20 | A cancelled build sends TERM, then KILL after its grace | `internal/runbinary/lifetime_test.go` (`TestBuilderCancel`) | A cancel request with KILL leaves no TERM line in the descendant log |
| RS22 | 20 | A completed build drains its descendant group before the selection returns | `internal/runbinary/lifetime_test.go` (`TestBuilderNormalDrain`) | A launch under another policy returns while the descendant group lives |
| RS23 | 21 | A malformed inherited vector refuses the build script before it starts | planned TestBuildScriptMalformedResources in internal/runbinary/lifetime_resources_test.go | A launch that ignores the vector starts the build |
| RS24 | 22 | The process-lifetime census holds both runner entries green | the process-lifetime-closure check in internal/conformance | A runner that reaches no owner call fails the closure diagnostic |
| RS25 | 2 | The cache clean of a cache with a retired store runs its release | `internal/gocache/protection_test.go` (`TestProtectionCompletedHold`) with a clean after the release | A clean that reads the retired record as a generation refuses or keeps the cache |
| RS26 | 2 | A retired checkout store holds `retired.json`, and the checkout cleanup of that store finds no user | `internal/worktree/protection_test.go` (`TestCheckoutProtectionPredicate`) | A cleanup that reads the retired record as a generation keeps the checkout |
| RS27 | 14 | The protection ADR consequence states that the last uncontended release of a cache hold deletes each retired cache store whose cache is gone | review-owned: docs/adr/0030 consequence | No test reads ADR prose |
| RS28 | 9 | The last uncontended release of a cache hold calls the absent-resource check once, after its own retirement | `internal/gocache/protection_test.go` (`TestReleaseRunsTheRetiredCheck`) | Without the call, the recorded call count is zero |
| RS29 | 17 | The live tree holds no hand-built inherited registration outside `internal/subprocess` | the `process-lifetime-closure` gate check over the live tree | A bounds runner or build-script runner copy that survives the change turns the live check red |

### Edge inventory

RS4 covers a malformed record, a link, another mode, an unknown field, and an oversized record.
RS2, RS25, and RS26 cover each reader of a store with no generation.
RS5 to RS8 cover a gone nested resource, a kept one, a prefix sibling, and a user record in a retired store.
RS9 to RS13 cover a gone cache, a kept cache, a contended release, an unreadable record, and a failed deletion.
RS19 and RS23 cover a malformed inherited vector for each runner.

The hostile-input checklist of `projects/benchkit.md` applies to the retired record path: a relative path, a path with control bytes, and a link at the record all refuse through the record owner's path check and no-follow read, which RS4 grades.

A crash after the retired-record publish and before the `resource.json` removal leaves a sealed generation with no user. The absent-resource check and the retired branch skip that store, because `resource.json` is present. Only a later `Reopen` recovery or the live branch of `discardNested` reaches it.

Won't handle: a cache directory that a process recreates between the check's absence test and its deletion — that open fails with no unprotected launch, and the next hold opens a fresh store.
Won't handle: a cache that a process creates again at the same path — the check tests the path only, so the store of the earlier binding stays until that cache goes again.
Won't handle: the Darwin witness of the builder group policy — the reviewer deferred it until a Mac host is available, and FT446 holds it with no owner.
Won't handle: the FT482 seal policy for a live, leaked, or dead protected user — FT482 is a separate decision row.
Won't handle: the registrations in the checkout exec, the checkout subshell, and the shift session — each passes its own resources and reads no inherited vector through `Inherited`.
Won't handle: the kit run's union and the checkout protection union — each reads its vector from a named environment, not from the process environment that `RunInherited` reads.
Won't handle: the report of a failed deletion to a user — every production caller of the cache release drops its error, and this spec does not change that reporting.
Won't handle: the live-generation branch of `discardNested` — it runs after the outer release, and this spec changes only the retired branch.
Won't handle: the test-side store cleanup in the test-report lifetime tests — it removes stores that those tests own, and its retirement is a separate cleanup.

## Out of scope

A root-wide sweep of retired stores of every kind at each store open is a separate capability: about 6 edits, 1 gate run.

## Further notes

### Seams beyond the decision source

| seam | path:line | research finding |
|---|---|---|
| The retired record as a separate file (D1) | internal/subprocess/protection/recovery.go:90-92 | Four readers key on an absent `resource.json` through `NoGeneration`: internal/gocache/clean.go:104, internal/worktree/protection.go:87, internal/systemtest/process_lifetime_artifacts_test.go:80 and :121 |
| The empty-users condition (D2) | internal/subprocess/protection/directory.go:264-281 | The current sweep deletes on an absent resource without a user check |
| The cache release trigger (D3) | internal/gocache/lock.go:160-172 | The cache retire calls `SealAndObserve` at the last uncontended release, and no protected removal reaches a package-test cache |
| The registration parameter (D4) | internal/subprocess/inherited.go:23-28 | `RunInherited` fixes the installed protocol, and the build-script runner takes its protocol from its caller at internal/runbinary/runbinary.go:331-333 |
| The cache check variable (D3) | internal/conformance/injected_ports_registry_test.go:241-242 | The registry admits `storeBase` for protection tests only, and gocache tests share the real store root of the user, because no gocache test replaces `storeBase` (internal/subprocess/protection/export_test.go:23 is the only replacement) |
| The hand-built registration rule (D4) | internal/conformance/process_lifetime_test.go:24 | The census accepts `NewLifetime` or `RunInherited` by name, so it cannot see a hand-built copy |

### Reader sweep

- `resource.json` readers: each is listed under D1 as changed or unchanged.
- The retired record is new, so it has no reader outside the protection package.
- The injected-port registry gains one `seamVariable` row for the cache check variable in `internal/gocache`. That row names the consuming test, and the registry file joins the RS-C3 writes.
- The closure check rule joins `internal/conformance/process_lifetime_closure_test.go`, and its fixtures join the new `internal/conformance/process_lifetime_copy_test.go`, in the writes of ticket 05.
- `RunInherited` callers: internal/gate/runner.go:66, internal/gate/background_launch.go:226, internal/preprelease/preprelease.go:273, internal/releasepreflight/command.go:308, internal/testreport/command.go:307. The synthetic source at internal/conformance/process_lifetime_test.go:273 is a text fixture for the census and does not compile, so it keeps its text.
- The injected-port registry row for `BENCH_PROCESS_RESOURCES` names TestBoundsMalformedResources, which stays.
- No golden file, testdata fixture, `.mjs` script, workflow, hook, or anchor holds a store record byte or the changed ADR sentence.

### Pre-review proof checklist

- Cited symbols: `SealAndObserve`, `Recover`, `NoGeneration`, `discardNested`, `ObserveRecorded`, `Generation.Remove`, `readRecord`, `heldLock.retire`, `RunInherited`, `Registration`, `bounds.run`, and `runBuildScript` each resolve at 92cfa43e8.
- Import edges: none new. The cache package already imports the protection package, and the bounds and run-binary packages already import the subprocess package.
- Source-row clauses and occurrences: see the source-sentence-to-row table below. The two FT446 occurrences are the FT362 PL-C3 store residue and the PL-C4 Standards finding on the two hand-built registrations.
- Promised field labels: the retired record holds `schema`, `kind`, `resource_path`, and `binding`.
- Changed-function callers, from `bench consumers --production`:
  - `RunInherited`: `startGrade`, `runProcessGroupCommand`, `runStep`, `runner.launch`, `runGoChild`.
  - `SealAndObserve`: `heldLock.retire`, `Directory.ObserveRecorded`, `Generation.Release`, `Protocol.Recover`. Every caller now leaves a retired record on a safe retirement.
  - `discardNested`: `Generation.Remove`.
  - `heldLock.retire`: `Holder.Release`.
  - `bounds.run`: `Run`, `RunOutput`.
  - `runBuildScript`: `BuildSubject`, `canonicalBuild`.
- Copy survival: RS18 and RS29 fail if a hand-built inherited registration survives in either runner.
- Rendered-shape readers: none. No rendered text changes. The ADR and the data-handling record are prose that no test reads.
- Pin operators: the census closure applies exact name membership to the owner calls. The injected-port registry applies exact set equality to its admitted seams.
- Entry reads: `RunInherited` reads the process vector through `protection.Inherited`. RS16 pins that vector with `t.Setenv`. The check reads the store root through `storeRoot`, which the protection tests replace through `storeBase`.
- Derived expectations: each store expectation comes from the test's own fixture paths, not from the code under test.
- Consolidated rules: the inherited registration rule moves from two runners into `RunInherited`. The bounds runner and the build-script runner each keep their cause and signal rule as their cancel request (D5).
- Quantified obligations: "each retired store" covers both deleters through the one predicate. "Inherited registration has one owner" covers the two named runners and the closure rule.
- Workflow-step writes: none. No workflow step changes.

### Flagged additions

- The empty-users condition (D2) is a safety guard that the decision source does not state.
- The cache check variable and its injected-port registry row (D3) are a test seam, admitted under ADR 0034.
- The hand-built registration rule in the closure check (D4) is in scope by coordinator direction.
- The kind scope of the absent-resource check (RS30) keeps the cache release from deleting the stores of other kinds.
- The non-blocking store lock and its prefilter (D2, RS31) keep a sweep from waiting on another process. The store owner joins the writes of ticket 02. Coordinator decision at slicing.
- The release-error Won't handle line records that no caller reports the check's error. Coordinator decision at slicing.
- The RS30 probe drops both kind tests. Coordinator decision at slicing.
- The RS-C3 plan grades the registry row, and the RS-C4 plan grades the bounds runner and the five caller packages. Coordinator decision at slicing.
- RS13 reuses the privilege class of `internal/capability` for its root guard. Coordinator decision at slicing.
- RS2 and RS26 first require `retired.json`, so each has a red path on the current tree. Coordinator decision at slicing.
- RS4 runs the existing record tests over each record file name. Coordinator decision at slicing.
- The closure rule lives in the closure file diagnostics, and its fixtures take one new conformance test file. Coordinator decision at slicing.

### Source sentence to row

| source sentence | rows |
| --- | --- |
| A retired cache store keeps no resource path. | RS1, RS2, RS3, RS4, RS25, RS26 |
| So the nested-store sweep cannot delete it when its test cache directory is gone. | RS5, RS6, RS7, RS8, RS14, RS31 |
| A retired store records its resource path, so the sweep and an absent-resource check can delete it. | RS9, RS10, RS11, RS12, RS13, RS27, RS28, RS30 |
| `bounds` and `runbinary` each build an inherited registration by hand. | RS17, RS18, RS19, RS20, RS21, RS22, RS23, RS24, RS29 |
| The shared inherited run in `subprocess` needs a protocol parameter before `runbinary` can call it. | RS16 |

### Flagged for the reviewer

- `DATA_HANDLING.md` says that each resource owner keeps its store beside the resource. The stores live under the per-user boot root. That drift predates this spec, so this spec does not change that sentence.
- A lifetime with a nil protocol and local resources registers nothing. D4's nil-means-installed rule narrows that fail-open path for `RunInherited`, but the path stays for a direct `NewLifetime` caller.
