# A cache release seals no generation for a live user, and a Markdown-only commit starts no background grade

Status: staged

Roadmap: FT482

Decision source: `roadmap/FT482.md` (named reviewed artifact), with the reviewer decisions of 2026-11-07 that it records. The reviewer answered one late clarification on 2026-11-07 in the spec-stage conversation: the D2 fix covers the unreaped zombie only.

Verification log: pending.

## Problem

A test run can leak a live descendant that holds a protected-user registration on the build cache. When the cache holder releases its hold, the release seals the cache generation first and grades the users after. The live user keeps the seal in place. Each later `bench test`, gate, or lane commit on that HOME then refuses with "the resource generation is sealed" until the descendant exits. The release destroys nothing, so the seal only blocks work.

A dead process that its parent has not reaped is a zombie. A zombie answers the signal-zero probe as present. So a user whose actor and group are dead, but unreaped, still grades as protected. A scratchpad probe on this host showed that a zombie answers present for its pid and for its group, and that both answer ESRCH after the reap.

A Markdown-only commit in a kit worktree starts a background lane pass. That pass can queue a grade, and the grade holds the cache. When the generation is sealed, the next lane commit refuses. A commit that changes no graded build input needs no background grade.

## Solution

The last release of a cache hold grades every user before it changes admission. It seals and retires the generation only when every user is safe. When any user is protected or uncertain, the release writes nothing, so the generation stays as it was. A clean and a resource removal keep the current seal-first observation.

On Linux, the protection probe treats a zombie as dead. A process counts as present only when one of its tasks is not a zombie. A group counts as present only when one of its members counts as present. Other systems keep the signal-zero answer.

A commit or a merge in a kit worktree starts no background lane pass when it changes at least one path from its first parent, every changed path is Markdown, and it changes no build input.

## User stories

Line: opus / high.
Implementation-line reason: The hardest material chunk is SP-C2, which adds a Linux process-state reader with a hostile parse and a real-zombie test. SP-C1 reorders a shared crash-safe tail under a lock. The spec fixes each predicate, the seams are known, and package tests catch each wrong answer, so the mid tier fits; high effort buys care on the two lock and process chunks.
Harder chunks: SP-C1, SP-C2.

### The cache release

1. As an agent, I want a cache release beside a live user to leave the generation open, so that my next `bench test` on the same HOME starts.
2. As the protection owner, I want a release observation of a generation with a protected user to change no admission, so that a later launch can join that generation.
3. As an agent, I want the release observation to name each unresolved user, so that a retained generation still explains itself.
4. As the protection owner, I want a release of a generation whose every user is safe to seal and retire it, so that a later clean finds no user and a later hold opens a fresh generation.
5. As an agent, I want a release beside an uncertain user to leave the generation open, so that a probe that cannot decide does not block my HOME.
6. As the protection owner, I want a release to keep a seal that a clean already wrote, so that a release never reopens a generation that a cleaner closed.
7. As the protection owner, I want a failed seal publish on the safe release path to return uncertain with every user record kept, so that a crash or a write fault never loses a user.
8. As the protection owner, I want removal and recovery to keep the seal-first observation, so that a resource that an owner wants to delete still stops new users.
9. As an operator, I want `bench cache clean` beside a live oracle to keep its protection refusal, so that a clean never deletes a cache under a live user.
10. As an agent, I want a second `bench test` after a run that leaked a live descendant to run its go test, so that one leak does not block every later run on that HOME.

### The dead user

11. As the protection owner, I want a known user whose actor and group leader are unreaped zombies to grade safe on Linux, so that a release retires its generation.
12. As the protection owner, I want a pending user whose actor is an unreaped zombie to grade safe, so that a launch that died before its start does not hold the resource.
13. As the protection owner, I want a group with a zombie leader and one live member to grade protected, so that a live descendant still keeps its resource.
14. As the protection owner, I want a process with a zombie leader task and one live task to count as present, so that a live thread still keeps its resource.
15. As the protection owner, I want the state read from a stat record to come after the last `)`, so that a hostile command name cannot forge a zombie state.
16. As the protection owner, I want an unreadable or malformed stat record to keep the signal-zero answer, so that a failed read never releases a live user.
17. As the protection owner, I want a group member that exits during the scan to decide nothing, so that a race cannot keep a dead group present.
18. As the protection owner, I want a live non-zombie actor to keep its user protected, so that the reused-PID rule stays in force.
19. As the protection owner, I want a live group member to keep its user protected, so that the zombie rule changes only dead processes.
20. As the protection owner, I want the probe to send signal zero only, so that the zombie rule adds no signal.

### The Markdown-only commit

21. As an agent, I want a kit-worktree commit of only Markdown to start no background lane pass, so that a record commit cannot hold the cache for a grade that it does not need.
22. As an agent, I want a commit of Markdown and one other path to start one lane pass, so that a code change keeps its background grade.
23. As an agent, I want a merge whose fold brings only Markdown to start no lane pass, so that a fold of documentation from `main` holds no cache.
24. As an agent, I want a commit that only deletes Markdown to start no lane pass, so that the rule meets each kind of change.
25. As an agent, I want a commit that renames a non-Markdown file to a Markdown name to start one lane pass, so that rename detection cannot hide the deleted path.
26. As an agent, I want a Markdown path with a space and a double quote to count as Markdown, so that git path quoting cannot change the class.
27. As an agent, I want a commit of a Markdown build input to start one lane pass, so that a change to the binary keeps its background grade.
28. As an agent, I want a root commit to start one lane pass, so that a tip with no first parent keeps the current behavior.
29. As an agent, I want a tip that git cannot compare to start one lane pass, so that a failed read never drops a needed grade.
30. As an agent, I want a commit with no change from its first parent to start one lane pass as today, so that the queue keeps its own rule for a tip whose key it already holds.
31. As an agent, I want a commit or a merge of a non-Markdown path to start one lane pass as before, so that the exemption changes no other commit.
32. As a reader of `.bench/BENCH-reference.md`, I want the `Background grade` section to state the Markdown-only exemption, so that the reference matches the verbs.

### The retirement

33. As a coordinator, I want `docs/known-issues.md` to hold no FT482 row after this spec lands, so that no charge applies a retired workaround.
34. As a teammate who just walked in, I want ADR 0030 to state the release policy and the Linux zombie rule, so that the decided state has one record.

## Implementation decisions

### The release observation

The protection protocol gains one release observation beside `SealAndObserve`. Under the store's protocol lock, it reads the generation strictly and grades every user with the one `grade` predicate. A read failure returns uncertain and writes nothing, as `SealAndObserve` does today (`internal/subprocess/protection/recovery.go:32-35`). When any user is protected or uncertain, the observation returns that disposition with each unresolved user, and it writes no record. When every user is safe, it runs the same tail as `SealAndObserve`. That tail seals admission, removes each user, writes the retired record, and removes the generation, in that order.

The two observations share one body that takes the seal order as its one difference. No second copy of the grade loop or of the retirement tail exists. The release observation never writes admission open. A generation that a clean sealed earlier stays sealed until its users resolve.

The last uncontended release of a cache hold (`internal/gocache/lock.go:161-176`) calls the release observation in place of `SealAndObserve`. Each other caller keeps `SealAndObserve`: `Generation.Release`, `Generation.Remove`, `ObserveRecorded`, `Recover`, the cache clean, and the worktree, release-evidence, prospective, binary, and kit-test owners. So `TestRegistrationSealedAdmission` and the worktree and preprelease `ErrSealed` pins stay green. The comments at `internal/gocache/lock.go:157-160`, `internal/subprocess/protection/store.go:93-95`, and `internal/gocache/protection_test.go:193-194` change to the new state.

The focused-test system row that leaks a live descendant (`internal/systemtest/process_lifetime_callers_test.go:222-226`) drops its private HOME. A new system test runs a second `bench test` on the same HOME after the leak. That test is the end-to-end proof of story 10.

### The zombie rule

The protocol's default observer (`internal/subprocess/protection/store.go:49-54`) changes from `syscall.Kill` to a protection-owned probe. On Linux, the probe sends signal zero first. A non-nil answer returns unchanged. A nil answer to a signal-zero probe of a pid then reads `/proc/<pid>/task/*/stat`. The process counts as absent, with ESRCH, only when every task that the read finds is in state `Z`. A nil answer to a signal-zero probe of a group then scans `/proc/*/stat` for each process whose process group is that group. The group counts as absent only when no member counts as present under the process rule.

The probe parses the state as the first field after the last `)` of the stat record, and the process group as the third field. A record or a task directory that it cannot read or parse keeps the nil answer, so the user stays protected. An entry that is gone (ENOENT) during the scan decides nothing. The probe never sends a signal other than the one that it received.

The reader takes its proc root as a parameter. Production passes `/proc`. Package tests drive the same reader over a fixture directory, so the hostile parse rows need no package variable and no injected-port row. The probe and its reader live in the existing platform files of the package: `root_linux.go` holds the Linux probe, and `root_darwin.go` and `root_other.go` give the other systems the plain `syscall.Kill` probe. The package directory already holds 18 source files against a limit of 12, so SP-C2 adds only the Linux test file `observe_linux_test.go`.

`Registration.Finish` also reads the default observer (`internal/subprocess/protection/registration.go:146`). With the zombie rule, `Finish` adds the group obligation only for a group that has a member that is not a zombie.

`ObserveGroup`, `Classify`, and the lifetime controller keep the signal-zero answer. The controller observes a group only after it reaps the leader (`internal/subprocess/control.go:309-314`). The record schema does not change, so no persisted field joins the record.

### The Markdown-only exemption

`gate.StartBackground` (`internal/gate/background_launch.go:34-41`) is the one launch of a lane pass, and both callers reach it: `internal/commit/background.go:74-79` and `internal/worktree/merge.go:135-137`. After the existing worktree check, the launch reads the change list of `<tip>^1` to `<tip>` through `benchgit.TreeChanges` (`internal/git/tree.go:301`). That helper reads `--raw --no-renames -z`, so each path is its own change and git quotes no path. When the list holds at least one change, every changed path is Markdown, and `freshness.DiffChangesBuildInputs` answers false, the launch starts no pass and prints nothing. An empty list starts the pass as today, and the queue's key rule then decides (BG45). A root commit, a tip that git cannot compare, or a list that `TreeChanges` cannot read starts the pass as today.

"Markdown" has one source: `prose.IsSubjectName` (`internal/prose/walk.go:27`), the `.md` rule that the prose grade keys on. The launch calls it. The lane repeats the same suffix test twice today: in `proseSubject` (`internal/gate/lane_select.go:41`) and in the `markdown` path class (`internal/gate/lane_select.go:147-151`). Both change to call `prose.IsSubjectName`. The class keeps its name and its checks, so the rendered profile lane table does not change. The commit's own outcome line (`background grade: latest …`) still prints before the launch, so the author still sees the latest outcome.

A queued or running grade of an older tip of the same worktree stays in the queue after a Markdown-only commit. The exemption only starts no new pass.

### The retirement

The SP-C2 ticket removes the FT482 row from `docs/known-issues.md`, because that row's retire condition needs SP-C1 and SP-C2 only. It adds one Forbid entry to the ticket-slicing anchors (`internal/anchors/registry_ticket_passes.go`) whose needle is the bare token `FT482`, as the FT405, FT406, and FT438 entries do. Its independent rule in `TestTicketSlicingAnchors` uses the exact row text as its needle:

```text
| A live descendant of a `bench test` run seals the cache generation of its HOME when its CLI releases the cache. Each later `bench test` on that HOME then refuses with `the resource generation is sealed`. Under the gate, a system test that waits for a child mark can show this refusal as a 20 s timeout. | Stop the leaked descendant or wait until it exits, then run the command again. If a system test times out on a child mark, run that test alone before you run the gate again. | FT482: a cache release applies the decided seal policy for a live, leaked, or dead protected user. |
```

ADR 0030 gains the release policy in its protocol paragraph and the zombie rule in its platform scope. `.bench/BENCH-reference.md` `### Background grade` gains the exemption, and `CHANGELOG.md` gains one entry for each chunk.

## Implementation chunks

The build runs the chunks in table order on one integration source. SP-C2 waits for SP-C1, because both change the protection package and its tests. SP-C3 has no blocker, but it shares `CHANGELOG.md` with the other chunks, so the orchestrator runs it last.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| SP-C1 / `1-grade-before-the-cache-release-seals.md` | A cache release beside a live or uncertain user leaves the generation open; a safe release retires it | SP1, SP2, SP3, SP4, SP5, SP6, SP7, SP8, SP9, SP10, SP34 | protection, gocache, and system tests named in the rows | yes |
| SP-C2 / `2-count-a-zombie-user-as-dead.md` | On Linux, a zombie actor and a zombie-only group count as dead | SP11, SP12, SP13, SP14, SP15, SP16, SP17, SP18, SP19, SP20, SP33, SP35 | protection tests named in the rows | yes |
| SP-C3 / `3-start-no-grade-for-a-markdown-only-commit.md` | A Markdown-only commit or merge starts no background lane pass | SP21, SP22, SP23, SP24, SP25, SP26, SP27, SP28, SP29, SP30, SP31, SP32, SP36 | commit, merge, and gate launch tests named in the rows | no |

### Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"SP-C1","tickets":["1-grade-before-the-cache-release-seals.md"],"verification":[{"id":"release-seal-order","command":"bench test --package ./internal/gocache","probe":"seal before the grade in the release observation"},{"id":"protected-release-write","command":"bench test --package ./internal/subprocess/protection","probe":"write admission open on a protected release"}]},{"id":"SP-C2","tickets":["2-count-a-zombie-user-as-dead.md"],"verification":[{"id":"zombie-leader-only","command":"bench test --package ./internal/subprocess/protection","probe":"read only the leader task state"},{"id":"known-issues-row","command":"bench test --package ./internal/anchors","probe":"restore the FT482 known-issues row"}]},{"id":"SP-C3","tickets":["3-start-no-grade-for-a-markdown-only-commit.md"],"verification":[{"id":"rename-detection","command":"bench test --package ./internal/gate","probe":"read the change list with rename detection"},{"id":"commit-markdown","command":"bench test --package ./internal/commit","probe":"start the pass for every Markdown-only tip"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/seal-policy/spec.md"},{"id":"packages","command":"bench test --package ./..."},{"id":"system","command":"bench test --check system"}]}
```

## Testing decisions

- A good test drives the real owner and observes durable state: the store records, a later `Prepare` or hold, or the launch count of the stand-in wrapper.
- The protection rows use the package fixture (`internal/subprocess/protection/fixture_test.go`): `newProtocol`, `startUser`, `nestedHelper`, and `SetProbe` and `SetReplacer` from `export_test.go`. The gocache rows use `protectedHold` and `completedUser` (`internal/gocache/protection_test.go`) and `cleanWith` (`internal/gocache/clean.go:39`).
- The real-zombie rows start a child in its own group with `Setpgid`, do not reap it, and wait until its stat record shows state `Z`. The test reaps the child at cleanup. The hostile-parse rows drive the proc reader over a fixture directory.
- The launch rows use the stand-in wrapper of `newLaunchFixture` (`internal/gate/background_launch_test.go`) and the commit and merge fixtures that BG5 and BG6 use.
- SP27 writes its own minimal module in the gate test: `go.mod`, a `cmd/bench` main package, and a `scripts/go-build.inputs` line that names the Markdown file. The freshness fixture helper is test-only, so the gate test does not call it.
- The gate observes the feature through the package tests of the whole gate and through `bench test --check system` for SP10.

### Seam diagram

    trigger: the last release of a cache hold
        │
        ▼
    store records  ──▶  [ release observation: grade, then seal and retire only when safe ]  ──▶  disposition + records
                      ◀ tests attach here: drive Holder.Release or the protocol, then read the records or Prepare

    trigger: grade() probes an actor pid or a group
        │
        ▼
    pid or -pgid  ──▶  [ protection probe: signal zero, then the Linux task-state reader ]  ──▶  nil or ESRCH
                      ◀ tests attach here: a real unreaped child, or the reader over a fixture proc root

    trigger: bench commit or bench worktree merge publishes a tip in a kit worktree
        │
        ▼
    tip^1..tip  ──▶  [ StartBackground: Markdown-only and no build input → no pass ]  ──▶  launch or no launch
                      ◀ tests attach here: the stand-in wrapper counts launches

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| SP1 | 1 | After a release beside a live pending user, a new hold of the same cache joins the same generation without an error | planned TestReleaseKeepsALiveUsersGenerationOpen in internal/gocache/protection_test.go | The current seal-first release makes the new hold's `Reopen` return `ErrSealed` |
| SP2 | 2 | After a release observation beside a protected user, a `Prepare` on the same descriptor succeeds | planned TestReleaseObservationLeavesAdmissionOpen in internal/subprocess/protection/recovery_test.go | A release that seals first makes `Prepare` return `ErrSealed` |
| SP3 | 3 | The release observation beside a protected user returns the protected state with that user's id | planned TestReleaseObservationNamesTheLiveUser in internal/subprocess/protection/recovery_test.go | A release that returns safe or drops the user list hides the retention |
| SP4 | 4 | A release observation whose every user is safe leaves the retired record and no generation in the store | planned TestReleaseObservationRetiresASafeGeneration in internal/subprocess/protection/recovery_test.go | A release that never seals or never retires leaves the generation record in place |
| SP5 | 5 | After a release observation beside a user whose probe answers EPERM, a `Prepare` on the same descriptor succeeds | planned TestReleaseObservationLeavesAdmissionOpenWhenUncertain in internal/subprocess/protection/recovery_test.go | A release that seals on an uncertain grade makes `Prepare` return `ErrSealed` |
| SP6 | 6 | After a seal-first observation and a release observation beside a live user, a `Prepare` still returns `ErrSealed` | planned TestReleaseObservationKeepsAnEarlierSeal in internal/subprocess/protection/recovery_test.go | A release that writes admission open reopens the clean's seal |
| SP7 | 7 | A failed seal publish on the safe release path returns uncertain and leaves every user record in the store | planned TestReleaseObservationSealFailure in internal/subprocess/protection/recovery_test.go | A tail that removes users before a failed seal loses a user record |
| SP8 | 8 | A seal-first observation beside a live pending user makes the next `Prepare` return `ErrSealed` | `internal/subprocess/protection/registration_test.go` (`TestRegistrationSealedAdmission`) | A change inside the shared seal-first body turns this existing pin red |
| SP9 | 9 | A cache clean while the oracle group of another process lives refuses with `cache protected` | `internal/systemtest/process_lifetime_cache_test.go` (`TestProcessLifetimeCacheSecondProcessClean`) | A clean that adopts the release observation deletes a cache under a live user |
| SP10 | 10 | A second `bench test` on the HOME of a run that leaked a live descendant runs its go test and prints no `the resource generation is sealed` | planned TestFocusedTestAfterALeakedDescendant in internal/systemtest/process_lifetime_callers_test.go | The current release seals the HOME's generation, so the second run refuses |
| SP11 | 11 | A release observation of a known user whose actor and group leader are one unreaped zombie retires the generation on Linux | planned TestReleaseObservationRetiresAZombieUser in internal/subprocess/protection/observe_linux_test.go | Signal zero answers present for a zombie, so the user stays protected |
| SP12 | 12 | A pending user whose actor is an unreaped zombie grades safe on Linux | planned TestZombiePendingActorIsReleased in internal/subprocess/protection/observe_linux_test.go | Signal zero answers present for the zombie actor |
| SP13 | 13 | A user whose group leader is an unreaped zombie and whose group holds one live member grades protected | planned TestZombieLeaderWithALiveMemberIsProtected in internal/subprocess/protection/observe_linux_test.go | A rule that reads only the leader releases a live group |
| SP14 | 14 | The proc reader counts a process present when its leader task is in state `Z` and one other task is in state `S` | planned TestProcLiveTaskBesideAZombieLeader in internal/subprocess/protection/observe_linux_test.go | A rule that reads only the leader task releases a live process |
| SP15 | 15 | The proc reader reads state `S` from a stat record whose command name is `a) Z 1 (b` | planned TestProcStatParsesAHostileCommandName in internal/subprocess/protection/observe_linux_test.go | A parse that splits at the first `)` reads the forged `Z` |
| SP16 | 16 | The proc reader keeps a pid present when its stat record in the fixture proc root is unreadable or malformed | planned TestProcUnreadableStatKeepsPresent in internal/subprocess/protection/observe_linux_test.go | A read failure that counts as dead releases a live user |
| SP17 | 17 | The proc reader counts a group absent when its only entries are one zombie member and one member whose stat record is gone | planned TestProcVanishedMemberIsSkipped in internal/subprocess/protection/observe_linux_test.go | A gone entry that counts as present keeps a dead group protected |
| SP18 | 18 | A pending user whose actor is the live test process stays protected | `internal/subprocess/protection/recovery_test.go` (`TestRecoveryPendingActorAbsent`) | A zombie rule that answers ESRCH for a live process releases the user |
| SP19 | 19 | A known user whose actor is gone and whose group holds a live member stays protected | `internal/subprocess/protection/recovery_test.go` (`TestRecoveryWorkerActorDeath`) | A group rule that ignores live members releases the user |
| SP20 | 20 | A recovery sends signal zero only | `internal/subprocess/protection/recovery_test.go` (`TestRecoveryNoStaleSignal`) | A probe that sends a real signal fails the recorded-signal check |
| SP21 | 21 | A kit-worktree commit whose only change is one Markdown file starts no lane pass | planned TestMarkdownOnlyCommitStartsNoGrade in internal/commit/background_test.go | The current commit starts a lane pass for every published commit |
| SP22 | 22 | A kit-worktree commit of one Markdown file and one text file starts one lane pass | planned TestMixedCommitStartsAGrade in internal/commit/background_test.go | A rule that skips any commit with a Markdown path passes SP21 and fails here |
| SP23 | 23 | A non-current merge whose fold brings only one Markdown file starts no lane pass | planned TestMarkdownOnlyMergeStartsNoGrade in internal/worktree/merge_background_test.go | The merge call site starts a lane pass for each non-current merge |
| SP24 | 24 | A launch of a tip that only deletes one Markdown file starts no pass | planned TestMarkdownDeletionStartsNoGrade in internal/gate/background_launch_test.go | A predicate over the added paths alone starts a pass here |
| SP25 | 25 | A launch of a tip that renames `notes.txt` to `notes.md` starts one pass | planned TestRenameToMarkdownStartsAGrade in internal/gate/background_launch_test.go | Rename detection lists only `notes.md`, so the launch skips the pass |
| SP26 | 26 | A launch of a tip whose only change is `a "b".md` starts no pass | planned TestQuotedMarkdownPathStartsNoGrade in internal/gate/background_launch_test.go | A quoted name list ends the path with a `"`, so the path is not Markdown |
| SP27 | 27 | A launch of a tip whose only change is a Markdown file that the auxiliary build-input manifest names starts one pass | planned TestMarkdownBuildInputStartsAGrade in internal/gate/background_launch_test.go | A rule that ignores build inputs skips a binary change |
| SP28 | 28 | A launch of a root commit that holds only Markdown starts one pass | planned TestRootCommitStartsAGrade in internal/gate/background_launch_test.go | A rule that reads a failed first-parent diff as empty skips the pass |
| SP29 | 29 | A launch of a tip that names no commit starts one pass | planned TestUnreadableTipStartsAGrade in internal/gate/background_launch_test.go | A rule that reads a failed diff as Markdown-only skips the pass |
| SP30 | 30 | A launch of an empty commit while one grade runs queues one grade | `internal/gate/background_launch_test.go` (`TestSecondGradeStartsNoSecondOracle`) | A rule that reads an empty list as Markdown-only starts no pass, so the queue never holds the second grade |
| SP31 | 31 | A kit-worktree commit of `tracked.txt` starts one lane pass | `internal/commit/background_test.go` (`TestCommitStartsABackgroundGrade`) | An exemption that skips each commit turns this existing pin red |
| SP32 | 31 | A non-current merge of a `go.mod` change starts one lane pass | `internal/worktree/merge_background_test.go` (`TestMergeStartsABackgroundGrade`) | An exemption that skips each merge turns this existing pin red |
| SP33 | 33 | `docs/known-issues.md` holds no row whose retire condition names FT482 | `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A new forbidden rule over the row text reds when the row stays |
| SP34 | 34 | ADR 0030 states that the last release of a cache hold seals and retires only a safe generation | review-owned: the SP-C1 review reads the ADR protocol paragraph | No check grades ADR decisions |
| SP35 | 34 | ADR 0030 states that a zombie counts as dead on Linux | review-owned: the SP-C2 review reads the ADR platform scope | No check grades ADR decisions |
| SP36 | 32 | `.bench/BENCH-reference.md` `### Background grade` states that a Markdown-only commit or merge starts no lane pass | review-owned: the SP-C3 review reads the section | No anchor pins the section's sentences |

### Edge inventory

The hostile-input checklist of `projects/benchkit.md` applies through these classes:

- The predicate over a `git diff` path list meets each kind: add (SP21), edit (SP22), delete (SP24), and rename (SP25). The read goes through `benchgit.TreeChanges`, which uses `--no-renames` and `-z` (`internal/git/tree.go:311`), as `freshness.DiffChangesBuildInputs` does (`internal/freshness/freshness_buildinputs.go:132`). `TestDiffChangesBuildInputs` drives a rename and a non-ASCII path through that read (`internal/freshness/decision_test.go:56-57`), and SP25 and SP26 drive the two shapes at the new predicate.
- Absent versus empty: an absent first parent starts a pass (SP28), and an empty change list starts a pass (SP30).
- A path read out of a file: the stat record is resolved by its writer's rule, the last `)` (SP15). An unreadable record is unknown, not absent (SP16).
- A special file in a discovered path does not apply. The reader opens only `stat` records under the proc root, and the kernel owns those entries.
- Each fail-closed error has a row: the read failure before the seal keeps the existing return, the seal publish failure has SP7, the unreadable stat record has SP16, and the unreadable diff has SP29.

These behaviors serve the kit and every linked repository. A linked repository starts no grade today (`TestLinkedRepoCommitStartsNoGrade`), so SP-C3 serves the kit alone.

Won't handle: a reused PID of a dead actor — `TestRecoveryPendingActorAbsent` keeps a live answering actor protected, which also keeps a reused PID protected, and the next release retries after the other process exits.
Won't handle: a zombie on Darwin or another system — the plain signal-zero probe keeps the user protected, and the next release retires the generation after the parent reaps the zombie.
Won't handle: the seal of a refused clean or removal — `bench cache clean` and each owner's removal keep the seal-first observation, because they intend to delete the resource.
Won't handle: an older queued or running grade of the same worktree after a Markdown-only commit — the queue keeps that grade, and `bench worktree release` still stops it.
Won't handle: a Markdown name in another case, such as `README.MD` — the lane's `markdown` class matches `.md` only, so such a commit starts a pass as today.

## Out of scope

- PID-reuse detection through a recorded start time of the actor and the group leader: about 6 edits, 1 gate run.
- Zombie detection on Darwin through the process table: about 3 edits, 1 gate run.
- A release-style observation for the cache clean and the removal owners: about 4 edits, 1 gate run.

## Further notes

### Source-sentence-to-row table

| Source clause in `roadmap/FT482.md` | Rows |
| --- | --- |
| "A release does not seal a generation while a protected user is live." | SP1, SP2, SP3, SP6, SP10 |
| The release still retires a generation with no live user (the current safe tail). | SP4, SP7 |
| "The protection owner treats a user as released when its actor and its process group are dead on the current boot." | SP11, SP12, SP13, SP14, SP15, SP16, SP17 |
| The release policy is a cache release only; other observations stay. | SP8, SP9 |
| "A commit with no graded build input needs no background grade." | SP21, SP23, SP24, SP25, SP26, SP27 |
| The exemption changes no other commit. | SP22, SP28, SP29, SP30, SP31, SP32 |
| The known-issues row names FT482 as its retire condition. | SP33 |

The store root is per boot (`internal/subprocess/protection/root.go:21-34`), so a record from an earlier boot is never read. "On the current boot" therefore needs no new check.

### Flagged additions

- SP5: a release beside an uncertain user also leaves admission open. The decision names a live protected user. An uncertain grade means that the probe cannot decide, and the release deletes nothing, so a seal buys no safety there. Reviewer veto keeps the seal for SP5.
- SP26 and SP15 add hostile-input rows that the checklist demands.
- The zombie rule is Linux-only. The decision names no platform, and ADR 0030 serves Linux and Darwin. Darwin keeps the signal-zero answer, and Out of scope prices Darwin detection. Reviewer veto adds Darwin to SP-C2.
- Consequence, not an addition: a landing whose last commits are Markdown-only finds no background green for its tip. That landing then runs its own whole gate in the foreground. The whole-gate Markdown checks, such as `prose-mechanics` and `roadmap-detail-integrity`, run there and not in the background. The lane still runs `prose`, `docs-currency-workflow`, and `guidance-prose-budgets` at each commit.

### Seam evidence

Each seam beyond the decision source cites one line and one research finding:

- The release observation sits beside `SealAndObserve` (`internal/subprocess/protection/recovery.go:17`). All six readers found that this body seals at lines 36-41 before it grades at lines 42-51.
- The release caller is `heldLock.retire` (`internal/gocache/lock.go:172`). The Q1 readers found that a cache-only change there reds none of the shared `ErrSealed` pins.
- The protection probe replaces the default of `Protocol.observer` (`internal/subprocess/protection/store.go:49-54`). The scratchpad probe on 2026-11-07 showed that a zombie answers present for its pid and its group, and that both answer ESRCH after the reap.
- The Markdown-only predicate sits in `StartBackground` (`internal/gate/background_launch.go:34`). The Q2 and Q3 readers found that both lane-pass callers reach it and that no start path reads a changed path today.
- The Markdown rule is `prose.IsSubjectName` (`internal/prose/walk.go:27`). The SP-C3 slicing readers found it, and found that `proseSubject` (`internal/gate/lane_select.go:41`) and the lane's `markdown` row (`internal/gate/lane_select.go:147-151`) repeat its suffix test.
- The change list is `benchgit.TreeChanges` (`internal/git/tree.go:301`). The SP-C3 slicing reader found that it reads `--raw --no-renames -z` and fails on an unreadable range.

### Pre-review proof checklist

- Cited symbols: `SealAndObserve`, `Recover`, `Reopen`, `grade`, `observer`, `heldLock.retire`, `Holder.Release`, `StartBackground`, `startBackgroundGrade`, `DiffChangesBuildInputs`, `laneClasses`, `ObserveGroup`, `Classify`, and each cited test name resolve in the tree at `7c0ac31a5`.
- Import edges: `internal/gate` already imports `internal/freshness` (`go list -deps ./internal/gate`). No other new edge.
- Source-row clauses and occurrences: the table above quotes each clause. The occurrences are the three lines of `roadmap/FT482.md:21-23`.
- Promised field labels: none.
- Changed-function callers: `SealAndObserve` keeps each caller that the Implementation decisions name. `heldLock.retire` has one caller, `Holder.Release`. `StartBackground` has two production callers and two test callers (`internal/gate/background_launch_test.go:52`, `internal/gate/background_queue_test.go:234`). `Protocol.observer` has two callers, `grade` and `Registration.Finish`.
- Copy survival: the shared observation body leaves one grade loop and one tail. SP8 and SP2 go red if either copy diverges on the seal order.
- Rendered-shape readers: no output shape changes. The refusal text of `ErrSealed` stays.
- Pin operators: SP33 uses the anchor harness's forbidden rule, which reds when the needle is present.
- Entry reads: the proc reader takes its root as a parameter, so no unexported ambient read sits below an entry.
- Derived expectations: none. Each expectation is a fixture fact.
- Consolidated rules: the lane's `proseSubject` and its `markdown` class drop their own `.md` suffix tests and call `prose.IsSubjectName`, which the launch also calls. Each old rule and the new rule accept the same names, so no acceptance cell changes. The class name and checks stay, so the rendered profile table stays.
- Quantified obligations: "every user" in the release observation is the users directory of the store; "every task" is `/proc/<pid>/task`; "every member" is each `/proc/<pid>/stat` whose process group matches; "every listed path" is the NUL-separated diff list.
- Workflow-step writes: none.
- Anchors: `bench anchors` ran on `.bench/BENCH-reference.md`, `docs/known-issues.md`, and ADR 0030. No anchor pins the `### Background grade` section, the FT482 row, or the ADR text.
