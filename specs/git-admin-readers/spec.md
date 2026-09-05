# Named Git adapter readers

Status: staged

Roadmap: FT302

Decision source: the reviewer-confirmed `/bench-deepen` grill of 2026-09-05 for candidate 1, recorded in `capture/session-handoff.md` under `## main`, with the survey report `capture/architecture-review-20260905T112705.html` as its source evidence

Verification log: 2 iteration(s) to accept — the round on `gpt-6-astra` at low effort blocked at both iterations, and the cap was two. The author folded the eleven findings of iteration 1 and the four of iteration 2. The two behavior folds are the hooks refusal at the link transaction and the doctor's repair, and the pass-through stub modes that reach four probes. The last fold adds the link transaction test's file to the fence. The reviewer's sign-off is the acceptance.

## Problem

Seventeen production sites outside the adapter ask Git for the same two
repository facts. Twelve sites ask for the checkout's administration directory,
and five sites ask for a named file inside it. The twelve count the pin path the
staged pin-removal spec deletes, and the adapter's own checkout probe is a
thirteenth spelling of the directory query. The sites spell the directory query in two ways, and
the file query in two ways.

The failure posture drifted four ways across them.
Two callers require an absolute answer, and two propagate the raw Git error.
One falls back to a different lock file, and three guess a path when Git fails.
The classifier uses one of these answers as ownership evidence under ADR 0005,
so a caller that accepts a symlinked answer weakens that decision.

The Git adapter already owns the sibling fact. `CommonDir` resolves the shared
administration directory, and `Worktrees` validates that answer for shape. No
caller of the checkout's own directory reaches that validator, and `CommonDir`
itself returns the raw answer.

## Solution

The Git adapter gains two named readers. The directory reader answers the
checkout administration directory, and the file reader answers the absolute
path of a named file inside it. Both readers and `CommonDir` run through the
bounded runner and validate their answer the way `Worktrees` does. One typed
error carries a subject field, so every consumer keeps one type check.

Every production site outside the adapter calls a reader. A standing conformance
check reds a production Go file outside the adapter that spells one of the four
Git flags. Three callers keep a decided fallback on the typed error. The hooks
directory and the two relative joins take the typed failure instead of a
guessed path. The glossary gains the term **checkout administration directory**.

## User stories

### Group A — the adapter

Line: opus / medium. The readers carry ownership evidence and a bounded
subprocess, and the scorecard routes gate and adapter logic to the mid tier at
medium effort.

1. As a caller, I want one reader for the checkout administration directory, so that I never spell a Git flag.
2. As a caller, I want one reader for a named file inside that directory, so that I never join a relative answer onto a root.
3. As a caller, I want `CommonDir` to validate its answer the way the readers do, so that one posture covers the three repository facts.
4. As a caller, I want an empty, absent, symlinked, or non-directory administration directory refused with a typed error, so that ownership evidence stays trusted.
5. As a caller, I want a relative answer joined onto the root and a symlink's own path kept, so that I never join.
6. As a caller, I want a hung `git` to return a typed timeout failure under the worktree-list bound, so that no reader stalls a command.
7. As a consumer of the typed error, I want its subject named, so that one type check serves the three readers.
8. As a status reader, I want the common-directory refusal text to keep the words `git common directory`, so that the rendered worktree row stays.
9. As a reviewer, I want each reader's answer to equal an independent `rev-parse` answer over three repository shapes, so that the deepening preserves behavior.
10. As a bare-repository operator, I want the directory reader and `CommonDir` to return one path, so that `IsPrimaryCheckout` keeps its answer.

### Group B — the callers

Line: opus / low. Each migration replaces one exact block at a known seam under
the covering suite and the new check, which is the scorecard's low-effort shape.

11. As a maintainer, I want every production site outside the adapter to call a reader, so that the adapter owns the two facts.
12. As a cleanup operator, I want the registration lock to keep its repository-level fallback when the reader fails, so that a non-worktree target still locks.
13. As a landing operator, I want the merge-pending probe to answer false when the reader fails, so that the commit-and-review route stays.
14. As a gate consumer, I want the composed-green probe to answer false when the reader fails, so that no reuse rests on an unreadable directory.
15. As a consumer, I want the lane, the inspection, the transaction, and the dashboard to keep their refusal outcome, so that nothing changes.
16. As a test author, I want the three tests with four flag-spelling sites to call the reader, so that a test never re-leaks the grammar.
17. As a file owner, I want each migrated file to end at or under its base line count, so that the growth ratchet stays green.

### Group C — the posture changes

Line: opus / medium. Each change alters an observable refusal, and the
scorecard routes a refusal change to medium effort.

18. As a linker, I want an unresolved hooks directory to make `bench link` exit 1 and name the cause, so that no hook lands blindly.
19. As a repairer, I want `bench doctor --fix` to exit 1 the same way, so that the repair never writes under a guessed path.
20. As an unlinker, I want an unresolved hooks directory to make `bench unlink` exit 1 and name the cause, so that no removal runs blindly.
21. As a reader, I want an unresolved hooks directory to read as an absent hook with an empty path, so that no guess is recorded.
22. As a lease reader, I want `LeaseFile` to return the reader's absolute answer or its typed failure, so that no call-site join survives.
23. As a diff reader, I want the index identity to use the reader's absolute answer or its typed failure, so that no call-site join survives.

### Group D — the standing check

Line: opus / medium. A conformance check with a canary fixture is the
scorecard's medium-effort shape.

24. As a maintainer, I want `git-plumbing-owner` to red a non-test Go file outside the adapter that spells a Git flag, so that drift cannot recur.
25. As a maintainer, I want the check to tolerate a flag inside a longer literal, so that the doctor's shell snippet passes.
26. As a maintainer, I want the check's diagnostic to name the file and the flag, so that a red needs no archaeology.
27. As a maintainer, I want one canary fixture that re-types a flag and reds the check with its planted diagnostic, so that the check bites.
28. As a maintainer, I want the check in the registry, the check map, the profile table, and the fixture map, so that each advertisement holds.
29. As a session, I want `bench test --check git-plumbing-owner` to run the check alone, so that a focused run reads one verdict.
32. As a maintainer, I want the check to pass a flag literal outside a `rev-parse` call, so that the guard option tables stay.

### Group E — the glossary

Line: opus / medium. Guidance prose routes to the mid tier at medium effort
under the reviewer's 2026-08-26 rule.

30. As a cold session, I want `CONTEXT.md` to define **checkout administration directory** with its Avoid list, so that the vocabulary does not drift.
31. As an operator, I want the classifier and ownership refusal texts to say `checkout administration directory`, so that the message uses the glossary term.

## Implementation decisions

- The adapter exports a directory reader and a file reader. The directory reader answers the checkout administration directory, which is the directory `git rev-parse --path-format=absolute --git-dir` names. The file reader takes a file name and answers the path `git rev-parse --git-path <name>` names in Git's default path format, because that format never resolves a symlink. If Git answers a relative path, the reader joins it onto the symlink-resolved absolute form of the root it passed with `-C`. Git ran there and resolves `..` physically. The adapter uses the `--path-format=absolute` spelling for the two directory queries, and the tree keeps no `--absolute-git-dir` spelling.
- The directory reader and `CommonDir` validate their answer with the existing common-directory validator. An empty answer, an absent path, a symlink, or a non-directory returns a typed `ResolutionError`. The file reader validates that its answer is non-empty, and it runs no existence check, because absence is the caller's fact.
- The three readers run through the bounded runner under the worktree-list bound. A timeout, a start failure, and a nonzero exit each return a typed `ResolutionError`, as `Worktrees` does today.
- `ResolutionError` gains a subject field with three values: the common directory, the checkout administration directory, and the checkout administration path. The error text opens with the subject's noun. The common-directory noun keeps the words `git common directory`. The action stays `investigate the git failure`. The two `WorktreeFailure` consumers keep their one type check.
- Every production site outside the adapter calls a reader. Six sites are in the worktree package: the classifier, the lifecycle lock, the ownership marker, the cleanup plan, the cleanup capture, and the landing refusal. Five are in the gate package: the fast lane, the verdict inspection, the run transaction, the composed-green probe, and the pin path. The dashboard is the twelfth, and the checkout probe inside the adapter is the thirteenth spelling. After the pin removal lands, eleven directory sites remain outside the adapter.
- The five file sites are the two index reads, the lease file, the diff index identity, and the hooks directory. The enumeration source is the search the check runs, not this list.
- A site the staged pin-removal spec deletes migrates only when it still exists at the build's base.
- Three callers keep a decided fallback on the typed error. The cleanup registration lock opens the repository-level lock file. The merge-pending probe answers false. The composed-green probe answers false. The fast lane, the verdict inspection, the run transaction, and the dashboard keep their current refusal outcome.
- The hooks directory takes the typed failure. The link transaction resolves the hooks directory through the reader before it inspects the hook, and it refuses with exit 1 when the reader fails. The doctor's repair resolves it the same way before it reads the hook state. The hook installer and the hook removal return the error. So `bench link`, `bench doctor --fix`, and `bench unlink` exit 1 with a message that names the root path and the action. The hook health inspection reports the absent state with an empty path, and it adds no state.
- `LeaseFile` and the diff index identity return the reader's absolute answer or its typed failure. The relative joins leave.
- The stub-git harness answers the directory query and the file query. It adds the modes `bad-git-dir`, `empty-git-dir`, `symlink-git-dir`, `file-git-dir`, `block-git-dir`, `block-git-path`, `fail-git-dir`, `fail-git-path`, and `relative-git-path`. The `fail-git-dir` mode exits nonzero on the directory query, and the `fail-git-path` mode exits nonzero on the file query. The `relative-git-path` mode answers the file query with `.git/<name>`. In those three modes the stub passes every other invocation to the real `git`, so a verb that runs many Git commands reaches the query.
- The `git-plumbing-owner` check parses each non-test Go file at the module root, under `cmd/`, and under `internal/`, outside `internal/git/` and `internal/gittest/`. The canary trees under `tests/` are fixtures, not graded code. It reds a string literal that equals one of the four flags inside a call whose arguments also hold the literal `rev-parse`. A flag inside a longer literal passes. A flag literal outside a `rev-parse` call passes, because the two guards parse Git's global options there. The diagnostic reads `<file> spells the Git administration flag <flag> outside internal/git`.
- The check joins the registry with the `go-source` input class, the check map, the profile table, and the fixture map. The fixture `git-flag-retyped` under the `package-core-guard` family re-adds one `--absolute-git-dir` literal to the dashboard and expects that diagnostic. The check lands in the contract ticket after every migration ticket, so no in-flight tree state reds it.
- `CONTEXT.md` gains the term **checkout administration directory**, the directory `rev-parse --git-dir` names. For the primary checkout it is the common directory. For a linked worktree it is that worktree's admin entry directory. Its Avoid list holds `git dir`, `gitdir`, and `private git dir`. The classifier and ownership refusal texts use the term.
- The tests in `internal/git` reuse the package's one repository constructor, because the ordinary-build census allows one.
- The exit proof is: the pre-existing suite passes with test logic unchanged, except the rows this spec names as posture changes. The named exceptions are the hooks and call-site-join rows GR20 to GR26, and the two tests that assert the raw `CommonDir` posture. GR3 and GR36 replace those two tests, which are `TestCommonDirReturnsUnvalidatedOutput` and `TestCommonDirKeepsPlainOutputFailure`. `TestTransactionalLinkAdoptsUnownedAdapterThroughSymlinkParent` keeps its assertions but gains a real repository at its root, because the new refusal needs one. One pre-existing assertion changes under story 21. `TestInspectPrePushRefusesSpecialFilesAndMissingGit` reads the absent state with an empty path when no `git` is on `PATH`, because the guessed fallback leaves.

## Testing decisions

- A good adapter test drives a reader over a real repository in three shapes. It compares the answer with an independent `git rev-parse` run through the test runner. The prior art is `TestWorktreesRejectsBadCommonDirBeforePorcelain` for the hostile shapes and `worktree_admin_enum_test.go` for the real shapes.
- The hostile shapes run through `gittest.StubGit`, which is the one stub harness. The harness gains the directory and file query modes.
- Each decided fallback gets one test at its own function with a non-repository target. The real `git` then fails, and no stub is needed.
- The posture changes run the real verb over a repository with the `fail-git-path` stub on `PATH`. The prior art is `TestPrePushHookAllowProtectedPushConfig` for the hook harness.
- The check gets four unit tests over a temporary tree and one canary fixture. The fixture-bite test proves the fixture through its registered owner. The fixture's mutation restores the old `rev-parse --absolute-git-dir` call in the dashboard, so the planted literal sits inside a `rev-parse` call as the rule requires.
- The posture and probe tests run a real verb or function over a real repository with a pass-through stub on `PATH`. A non-repository root stops at an earlier check in the dashboard, the composed-green probe, and the link transaction, so it never reaches the reader.
- The gate's conformance phase observes the check over the whole tree. The fast lane's structure growth check observes story 17 on each worktree commit.

### Seam diagram

    caller (worktree, gate, dashboard, diff, adopt)
        │
        ▼
    root, name  ──▶  [ Git adapter: directory reader, file reader, CommonDir ]  ──▶  absolute path | *ResolutionError{Subject, Action}
                          ◀ tests attach here: real repository versus an independent rev-parse, and gittest.StubGit modes

    bench gate (conformance phase)
        │
        ▼
    non-test Go files outside internal/git  ──▶  [ git-plumbing-owner ]  ──▶  no diagnostic | `<file> spells the Git administration flag <flag> outside internal/git`
                          ◀ tests attach here: four temporary-tree unit tests and the fixture git-flag-retyped

    bench link | bench doctor --fix | bench unlink
        │
        ▼
    root with the fail-git-path stub  ──▶  [ hooks directory through the file reader ]  ──▶  exit 1 + root path + `investigate the git failure`
                          ◀ tests attach here: the adopt hook test harness

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| GR1 | 1, 9 | over a primary checkout, a linked worktree, and a bare repository, the directory reader returns the path an independent `git rev-parse --path-format=absolute --git-dir` run prints | a new test `TestAdminDirMatchesIndependentRevParse` in `internal/git` | a reader that returns the common directory reds the linked case |
| GR2 | 2, 9 | over the same three repositories, and over a primary checkout reached through a symlinked parent directory, the file reader for `index` returns the path an independent `git rev-parse --path-format=absolute --git-path index` run prints | a new test `TestAdminPathMatchesIndependentRevParse` in `internal/git` | a reader that joins the name onto the root reds the linked case, and a reader that keeps the parent alias reds the symlinked case |
| GR3 | 3, 4 | with the stub answering the directory query with a missing path, an empty line, or a symlink, the directory reader and `CommonDir` each return a `*ResolutionError` whose action is `investigate the git failure` and whose text holds `missing path`, `empty path`, or `symlink` | a new test `TestReadersRefuseBadAdministrationDirectories` over the stub modes `bad-git-dir`, `empty-git-dir`, `symlink-git-dir`, and the three common-directory modes | a reader that returns the raw output accepts the symlink path, and the typed-error assertion reds |
| GR4 | 4 | with the stub answering the directory query with a regular file path, the directory reader returns a `*ResolutionError` whose text holds `non-directory` | the same new test, stub mode `file-git-dir` | a reader that skips the directory check accepts the file path, and the typed-error assertion reds |
| GR5 | 5 | with the stub answering the file query with `.git/<name>`, the file reader returns that answer joined onto the resolved root, for a plain root and for a root that reaches the repository through a symlink followed by `..` | a new test `TestAdminPathJoinsRelativeAnswerOntoRoot`, stub mode `relative-git-path` | a reader that returns the raw answer returns a relative path, and a reader that cleans the root lexically answers a path in the wrong directory |
| GR40 | 5 | over a repository whose `bench-lease` is a symlink to a file outside the repository, the file reader returns the symlink's own path, and `Lstat` on the answer reports a symlink | a new test `TestAdminPathKeepsSymlinkPath` in `internal/git` | a reader on `--path-format=absolute` returns the target, and the symlink assertion reds |
| GR6 | 6 | with the stub blocking on the directory query and the test bound set to 50 milliseconds, the directory reader returns a `*ResolutionError` whose text holds `timed out` before one second elapses | a new test `TestReadersTimeOutUnderTheWorktreeListBound`, stub mode `block-git-dir` under `SetWorktreeListTimeoutForTest` | an unbounded `Output` call blocks the test past its deadline |
| GR38 | 6 | with the stub blocking on the file query and on the common-directory query under the same bound, the file reader and `CommonDir` each return a `*ResolutionError` whose text holds `timed out` before one second elapses | the same new test, stub modes `block-git-path` and `block-rev-parse` | a reader left on the unbounded `Output` blocks the test past its deadline |
| GR39 | 6 | with an empty `PATH`, the three readers each return a `*ResolutionError` whose text holds `rev-parse` and `executable file not found` | a new test `TestReadersTypeStartFailures` in `internal/git` | a reader that returns the raw start error has no subject and no action |
| GR7 | 7 | the text of a `*ResolutionError` opens with `git common directory`, `checkout administration directory`, or `checkout administration path` by its subject | a new test `TestResolutionErrorNamesItsSubject` in `internal/git` | one text for every subject reds the two new subject cases |
| GR8 | 8 | the status worktree row detail holds `git common directory` when the common-directory read fails | `internal/status/status_render_test.go` (`TestAppendWorktreeSurfacesClassifyFailure`) | a renamed common-directory noun reds the assertion |
| GR9 | 10 | over a bare repository, the directory reader and `CommonDir` return one path, and `IsPrimaryCheckout` returns true | a new test `TestBareRepositoryReadersAgree` in `internal/git` | a reader that refuses a bare repository reds |
| GR10 | 11 | the `git-plumbing-owner` check returns no diagnostic over the landing tree | the gate's conformance phase | a surviving spelling outside `internal/git` reds the gate |
| GR11 | 12 | with a target directory that is not a repository, the cleanup registration lock returns a release function, a second non-blocking lock on the repository-level lock file fails while the release is pending, and it succeeds after the release runs | a new test `TestLockCleanupRegistrationFallsBackForNonWorktree` in `internal/worktree` | a caller that propagates the typed failure returns an error, and a no-op release or a lock on another file lets the second lock succeed early |
| GR12 | 13 | with a source directory that is not a repository, the merge-pending probe returns false | a new test `TestSourceMergePendingIsFalseWhenUndecided` in `internal/worktree` | a probe that treats the failure as pending returns true |
| GR13 | 14 | over a real repository with a closed subject and the `fail-git-dir` stub on `PATH`, `ComposedGreen` returns false | a new test `TestComposedGreenIsFalseWhenTheReaderFails` in `internal/gate` | a probe that propagates the failure panics or returns true, and a non-repository root never reaches the reader |
| GR14 | 15 | over a real repository with the `fail-git-dir` stub on `PATH`, the fast lane returns the error `gate: git directory unavailable` | a new test `TestRunLaneRefusesWhenTheReaderFails` in `internal/gate` | a migrated lane that propagates the typed text reds the exact message |
| GR15 | 15 | over a real repository with the `fail-git-dir` stub on `PATH`, the verdict inspection reports the state `Unavailable` with the reason `git directory unavailable` | a new test `TestInspectSubjectReportsUnavailableWhenTheReaderFails` in `internal/gate` | a migrated inspection that returns `Stale` reds |
| GR16 | 15 | over a real repository with the `fail-git-dir` stub on `PATH`, the dashboard file write exits 1 and prints `cannot resolve git dir` | a new test `TestDashboardWriteRefusesWhenTheReaderFails` in `internal/dashboard` | a migrated write that exits 0 reds, and a non-repository root stops at the not-in-repo error before the reader |
| GR17 | 15 | with a planned closed subject, the forced run mode, and the `fail-git-dir` stub on `PATH`, the run transaction returns `ActionExit` 1 with `GateExit` 0 and prints `git directory unavailable` | a new test `TestRunTransactionRefusesWhenTheReaderFails` in `internal/gate`, driving `executeSubjectWithRunBinary` through `plannedEvaluation` as `TestGateTransactionRefusesAnUnheldCache` does | a migrated transaction that starts the gate returns a nonzero `GateExit` |
| GR18 | 16 | the three named tests call the reader at their four sites, and their assertions are unchanged | review-owned: the reviewer reads `TestEligibilityVerdictProjectsWithoutSecondDecision`, `TestRecoveryPreservesEveryGitVisibleLayerWithoutMovingBranchOrIndex`, and `TestReauthorizeCommandRollsBackLockRefreshAndCASLoss` | the check grades no test file, so a survivor passes the gate |
| GR19 | 17 | `bench structure --growth <base>` over the landing source reports no over-budget file that gained lines | the fast lane's structure growth check on each worktree commit | a migrated over-budget file that gained one line reds the lane |
| GR20 | 18 | with the `fail-git-path` stub on `PATH`, `bench link` exits 1 before it stages any file, and it prints the root path and `investigate the git failure` | a new test `TestLinkRefusesUnresolvedHooksDirectory` in the adopt hook tests, driving `transactionalLink` | the old transaction joins `.git/hooks` onto the root, stages the hook there, and exits 0 |
| GR21 | 19 | in a consumer repository with the `fail-git-path` stub on `PATH`, `bench doctor --fix` exits 1 and prints the root path and `investigate the git failure` | a new test `TestDoctorFixRefusesUnresolvedHooksDirectory` in the adopt tests | the old repair reads an absent hook outside the kit checkout and exits 0 |
| GR22 | 20 | with the `fail-git-path` stub on `PATH` and a managed hook present, `bench unlink` exits 1 and prints the root path and `investigate the git failure` | a new test `TestUnlinkRefusesUnresolvedHooksDirectory` in the adopt tests | the old removal joins `.git/hooks` onto the root, removes the managed hook there, and exits 0 |
| GR23 | 21 | with the `fail-git-path` stub on `PATH`, the hook health record has the state `PrePushAbsent` and an empty `Path` | a new test `TestInspectPrePushReportsAbsentWhenHooksDirectoryIsUnresolved` in the adopt hook tests | the old inspection reports the joined `<root>/.git/hooks/pre-push` path |
| GR24 | 22 | over a linked worktree, `LeaseFile` returns the path an independent `git rev-parse --path-format=absolute --git-path bench-lease` run prints | a new test `TestLeaseFileMatchesIndependentRevParse` in `internal/worktree` | a join onto the worktree root reds the linked case |
| GR25 | 22 | with the `fail-git-path` stub on `PATH`, `LeaseFile` returns a `*ResolutionError` | a new test `TestLeaseFileRefusesUnresolvedAnswer` in `internal/worktree` | a lease reader that guesses a path under the root returns one |
| GR26 | 23 | with the `fail-git-path` stub on `PATH`, the index identity returns a `*ResolutionError` | a new test `TestIndexIdentityRefusesUnresolvedAnswer` in `internal/diff` | an index reader that guesses a path under the root reads a file there |
| GR27 | 24, 26 | for each of the four flags, over a temporary tree whose non-test Go file outside `internal/git` holds that flag literal inside a call with the literal `"rev-parse"`, and once more with that file at the module root, the check returns one diagnostic that names the file and the flag | a new test `TestGitPlumbingOwnerRedsARetypedFlag` in `internal/conformance`, one case per flag and one root-level case | a check hard-coded to one flag passes the other three cases, and a walker that visits only `cmd/` and `internal/` passes the root-level case |
| GR28 | 24 | over a temporary tree that holds a `rev-parse` call with `"--git-common-dir"` in a test file, in a file under `internal/git`, and in a file under `internal/gittest`, the check returns no diagnostic | a new test `TestGitPlumbingOwnerSkipsTestsAndTheAdapter` in `internal/conformance` | a check that grades every Go file reds the adapter or the harness |
| GR37 | 32 | over a temporary tree whose Go file holds the literal `"--git-dir"` in a map literal and in a call with no `rev-parse` argument, the check returns no diagnostic | a new test `TestGitPlumbingOwnerIgnoresCallsWithoutRevParse` in `internal/conformance` | a bare-literal check reds the guard option tables |
| GR29 | 25 | over a temporary tree whose Go file holds a flag inside a longer literal, the check returns no diagnostic | a new test `TestGitPlumbingOwnerToleratesEmbeddedFlag` in `internal/conformance` | a substring check reds the doctor snippet |
| GR30 | 27 | the fixture `git-flag-retyped` reds `git-plumbing-owner` with its `EXPECT` diagnostic and goes green on restore | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | a fixture whose mutation targets a test file never reds |
| GR31 | 28 | the registry, the check map, and the profile table agree on `git-plumbing-owner` | the `conformance-meta` and `profile-lane-table` checks | a missing row reds the binding check |
| GR32 | 28 | the fixture map classifies `git-flag-retyped` | `internal/conformance/registry_test.go` (`TestCanaryFixtureRegistryClassifiesEveryFixture`) | an unregistered fixture reds the classification |
| GR33 | 29 | `bench test --help` lists `git-plumbing-owner`, and `bench test --check git-plumbing-owner` runs only the registered dev-scope check | `internal/testreport/check_test.go` (`TestNamedCheckHelpListsEverySupportedCheck` and `TestNamedCheckRunsOnlyRegisteredDevScope`) | both tests derive their expectation from the registry, so the check joins the help inventory and the routed run without an edit, and GR31 reds a registry omission |
| GR34 | 30 | `CONTEXT.md` holds the term `checkout administration directory` with its Avoid list | review-owned: the reviewer reads the entry | the prose mechanics check grades sentences, not terms |
| GR35 | 31 | with the `fail-git-dir` stub on `PATH`, the classifier's owner-evidence refusal and the ownership marker refusal each hold `checkout administration directory` | a new test `TestOwnershipRefusalTextsUseTheGlossaryTerm` in `internal/worktree` | the old texts hold `private worktree administration directory` |
| GR36 | 3 | with the stub answering the common-directory query with a nonzero exit, `CommonDir` returns a `*ResolutionError` whose text holds `rev-parse` | a new test `TestCommonDirTypesRevParseFailure` in `internal/git`, which replaces `TestCommonDirKeepsPlainOutputFailure` | the old plain exit error has no subject and no action |

### Edge inventory

- Error paths: a nonzero `git` exit returns a typed error with the Git stderr (GR3, GR6). A start failure returns a typed error (existing `TestWorktreesTypesStartFailures`, which the readers share).
- Empty input: an empty answer refuses (GR3). An empty file name makes Git answer the administration directory itself, and no caller passes one.
- Boundary values: a bare repository is the one shape where the two directories coincide (GR9). A primary checkout and a linked worktree are the two ordinary shapes (GR1, GR2).
- Hostile paths: a symlinked directory refuses (GR3). A regular file refuses (GR4). A relative answer joins onto the root (GR5), and a symlinked file keeps its own path (GR40). A path with a space passes, because the adapter never splits the answer.
- Currency: the readers read no cache, so no staleness applies.
- Re-run idempotency: the readers write nothing.
- Partial implementation: a build that adds the readers and migrates no caller reds GR10. A build that migrates the callers and keeps `CommonDir` raw reds GR3. A build that skips the check reds GR31.
- Audience: the adapter and the check serve this repository. The verbs `bench link`, `bench doctor --fix`, and `bench unlink` serve every repository that links the kit, so GR20 to GR23 name that audience.
- Package-variable swaps: `SetWorktreeListTimeoutForTest` swaps the bound in the test process (GR6). The reader runs in-process, so the swap reaches it.
- Absent versus empty: an absent administration directory and an empty answer both refuse (GR3). An absent named file passes the file reader, because absence is the caller's fact (existing cleanup evidence tests).
- In-flight tree states of the new check: the check lands in the contract ticket after every migration ticket. The spec's close step flips one Markdown line and creates no Go change.
- Census budget: the new `internal/git` tests reuse the one repository constructor, so the ordinary-build census stays green.

**Won't handle** — the test-side spellings that run through independent test runners in six packages — ADR 0006 keeps an independent expectation as the omission oracle. The `outcomeGit` runner in `internal/gate/run_outcomes_test.go` stays their surviving caller.

**Won't handle** — the doctor's shell snippet that prints the common-directory query for the operator — the snippet runs without `bench` on `PATH`. GR29 proves the check passes it.

**Won't handle** — the wrapper's `rev-parse --git-common-dir` in `bin/bench.sh` — the wrapper runs before any Go binary exists, and `bench doctor` stays its surviving caller.

**Won't handle** — the `--git-path bench-gate-pin` line in the managed hook asset — the staged pin-removal spec deletes it.

**Won't handle** — the common-directory spelling in the `gittest` helper — the helper is the fixture harness, and it cannot import the adapter without a cycle. The check skips `internal/gittest/`, and GR28 proves it.

**Won't handle** — the `--git-dir` entries in the option tables of the two guards — the guards parse Git's global option there, and no call carries `rev-parse`. GR37 proves the check passes them. The guard classifier stays their caller.

**Won't handle** — an existence check in the file reader — the cleanup evidence reader handles absence, and `fileEvidence` stays its surviving caller.

**Won't handle** — a new hook health state — the reviewer decided on 2026-09-05 that the absent state with an empty path carries the case. `bench link` names the cause.

## Ownership fences

- `specs/git-admin-readers/`
- `reviews/git-admin-readers.md`
- `internal/git/worktree_admin.go`
- `internal/git/checkout.go`
- `internal/git/admin_readers_test.go` (new)
- `internal/git/worktree_admin_hostile_test.go`
- `internal/git/worktree_admin_enum_test.go`
- `internal/gittest/gittest.go`
- `internal/worktree/classifier.go`
- `internal/worktree/lifecycle.go`
- `internal/worktree/ownership.go`
- `internal/worktree/subshell.go`
- `internal/worktree/clean.go`
- `internal/worktree/land_refusal.go`
- `internal/worktree/worktree.go`
- `internal/worktree/eligibility_test.go`
- `internal/worktree/lifecycle_acquire_test.go`
- `internal/worktree/reauthorize_test.go`
- `internal/worktree/admin_readers_test.go` (new)
- `internal/worktree/parallel_census_test.go`
- `internal/gate/lane.go`
- `internal/gate/verdict.go`
- `internal/gate/run_transaction.go`
- `internal/gate/composed_green.go`
- `internal/gate/phases.go`
- `internal/gate/admin_readers_test.go` (new)
- `internal/dashboard/dashboard.go`
- `internal/dashboard/dashboard_test.go`
- `internal/diff/snapshot.go`
- `internal/diff/identity_test.go`
- `internal/adopt/link_hook.go`
- `internal/adopt/link_transaction.go`
- `internal/adopt/link_transaction_test.go`
- `internal/adopt/unlink.go`
- `internal/adopt/doctor.go`
- `internal/adopt/link_hook_test.go`
- `internal/adopt/adopt_test.go`
- `internal/conformance/git_plumbing_owner_test.go` (new)
- `internal/conformance/checks_test.go`
- `internal/conformance/registry/registry.go`
- `internal/conformance/registry_test.go`
- `tests/canary/package-core-guard/git-flag-retyped/` (new)
- `projects/benchkit.md`
- `CONTEXT.md`
- `capture/architecture-review-20260905T112705.html`
- `cmd/bench/command_registry.go` — closure headroom only
- `cmd/bench/command_registry_test.go` — closure headroom only
- `cmd/bench/main_test.go` — closure headroom only
- `internal/conformance/axi_query_registry_test.go` — closure headroom only
- `internal/conformance/subcommand_routing_test.go` — closure headroom only
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-coverage-map-term` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-coverage-row-parts` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-decision-map-term` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-reader-sweep-term` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-ticket-vocabulary` — closure headroom only
- `tests/canary/guidance-prose-budgets/over-budget-skill` — closure headroom only
- `tests/canary/line-routing/line-binding-prose-drift` — closure headroom only
- `tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading` — closure headroom only
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-owner` — closure headroom only
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-routing` — closure headroom only
- `tests/canary/workflow-guidance-anchors/benchkit-spec-ownership` — closure headroom only
- `tests/canary/workflow-guidance-anchors/benchkit-system-suite-route` — closure headroom only

A closure headroom entry creates no blocker edge and takes no edit. The five
command registries close the `internal/worktree` and `internal/diff` bindings.
The fifteen fixture directories close the pins on `CONTEXT.md` and
`projects/benchkit.md`.

The fence is the union of the tickets' `Writes:` lines, plus the spec folder, the review pickup, and the report copy. `bench preflight build` closes it over the fixture and registry pins. Six worktree files sit over the structure budget: the classifier, the lifecycle, the ownership, the subshell, the cleanup, and the worktree file. GR19 holds each at or under its base line count.

## Out of scope

- One landing-destination fact for the first run and the resume. Candidate 2 of the same survey, one light-path ticket. 4 edits, 1 gate run.
- Pure policy behind the intent ledger and the composition. Candidates 3 and 4 of the same survey, one spec. 20 edits, 4 gate runs.
- Named readers for the peeled-commit and porcelain-status facts. FT218 requires a fresh census before a ticket is cut. 6 edits, 1 gate run.
- A `bench test --check <name> --fixtures` projection that lists the fixture a check owns. FT290 owns it. 3 edits, 1 gate run.

## Further notes

Flagged additions beyond the decision source:

- The hook health record's absent state with an empty path. The grill named three hooks callers, and the inspection was the fourth. The reviewer decided it on 2026-09-05 after the grill.
- The `--git-common-dir` flag in the check's set. The grill named three flags, and the check reads the adapter's fourth query the same way.
- The stub-git modes. The grill decided the posture, and the harness needs the modes to reach it.
- The check's `rev-parse` call rule and its `internal/gittest/` skip. The ticket slice found three bare flag literals the grill's rule would red: the harness helper and the two guards' option tables.
- The hooks-directory resolution inside the link transaction and the doctor's repair. The review round found that both read the hook record before the installer runs, so the installer alone cannot deliver the refusal.
- The `fail-git-dir` pass-through stub mode. The review round found that a non-repository root stops before the reader in four probes.
- The fixture change in `TestTransactionalLinkAdoptsUnownedAdapterThroughSymlinkParent`. Its root becomes a real repository, and its assertions stay.
- The reader names `AdminDir` and `AdminPath`. The tickets fix one spelling so that sibling tickets compile together.

Build decisions recorded for reviewer veto:

- The check lands in the contract ticket, after every migration, so no in-flight state reds it.
- The `relative-git-path` and `fail-git-path` stubs pass every other invocation to the real `git`. A verb that runs many Git commands must reach the hooks query.
- The file reader uses Git's default path format and joins a relative answer onto the root. The build found on 2026-09-05 that `--path-format=absolute` resolves an existing symlink. Under that spelling `TestReleaseSymlinkLeaseRetainsAsUncertain` reds, and every `Lstat` at a file site is defeated. Rows GR5, GR20 to GR23, GR25, and GR26 now drive a failing file query, and GR40 pins the symlink posture. The reviewer can veto this at the review.
- Review round 1 on 2026-09-05 accepted four repairs. The file reader resolves the root's symlinks before the join (GR2, GR5). The GR11 test proves the lock is held, and the check walks the module root (GR27). The repair commit closed the pickup file. The non-blocking findings stay here for the reviewer's decision:
  - `Worktrees` repeats the common-directory resolution that `CommonDir` owns.
  - `AdminPath` repeats the empty-answer refusal that `validateCommonDir` owns.
  - `stringLiteralValue` in the check repeats `stringLiteral` in the routing check.
  - The root-level walk in the check repeats the file eligibility filter of the recursive walk.
  - The hostile reader table re-types the stub's three planted admin paths.
  - The `StubGitDir` comment inventories the modes the script below encodes.
  - The new test doc comments cite coverage row ids, and two carry the red record the spec owns.
  - The pass-through stub modes match a flag as a substring of the whole argument line. A root whose path text holds a flag is intercepted. The harness is test-only, and no test root carries flag text.
- The two verb tests GR21 and GR22 live in the adopt hook test file, not in the adopt test file the ticket named. That file sits over the structure budget, and the growth ratchet reds a line it gains.
- A site the staged pin-removal spec deletes migrates only when it still exists at the build's base. If it does, the pin-removal landing composes over one migrated line and repairs by merge.
- The two `PATH`-binding tests in `internal/worktree` (GR25 and GR35) join the package's serial set. The pinned serial ceiling in `internal/worktree/parallel_census_test.go` rises by two inside its recorded stub-test class. The fence gains that file.

Source-sentence-to-row table:

| source sentence | rows |
|---|---|
| the adapter exports a directory reader and a file reader with a file-name parameter | GR1, GR2 |
| every one of the 17 sites calls a reader, and no production file outside `internal/git` spells the flags | GR10, GR27, GR28 |
| the directory reader and `CommonDir` run the full validator with a typed `ResolutionError` | GR3, GR4, GR36 |
| one error type carries a subject field, and the common-directory text stays | GR7, GR8 |
| the readers and `CommonDir` run through the bounded runner | GR6, GR38, GR39 |
| the file reader returns an absolute non-empty path and runs no existence check | GR2, GR5, GR40 |
| the cleanup lock, the merge-pending probe, and the composed-green probe keep their fallback | GR11, GR12, GR13 |
| the hooks directory and the two relative joins take the typed failure | GR20, GR21, GR22, GR23, GR24, GR25, GR26 |
| the four tests call the reader with their assertions unchanged | GR18 |
| the exit proof: the suite passes with test logic unchanged, except the named posture rows | GR14, GR15, GR16, GR17, GR19 |
| a standing conformance check with a canary fixture and the fixture-bite registry rows | GR29, GR30, GR31, GR32, GR33, GR37 |
| `CONTEXT.md` gains the term, and the classifier and ownership texts adopt it | GR34, GR35 |
| a bare repository where both readers agree | GR9 |

Pre-review proof checklist:

- Cited symbols: each symbol below resolves in the tree at the spec commit. The tests named `new` are new.
- Cited adapter symbols: `CommonDir`, `validateCommonDir`, `boundedGit`, `worktreeListTimeout`, `SetWorktreeListTimeoutForTest`, `ResolutionError`, `WorktreeFailure`, `Worktrees`, `IsPrimaryCheckout`, `Output`.
- Cited caller symbols: `LeaseFile`, `indexIdentity`, `hooksDir`, `installGitHook`, `removeManagedHook`, `InspectPrePush`, `noPrePushHealth`, `PrePushAbsent`, `transactionalLink`, `lockCleanupRegistration`, `sourceMergePending`, `ComposedGreen`, `runLane`, `inspectSubjectAt`, `executeSubjectWithRunBinary`, `plannedEvaluation`, `forceRun`, `pinPath`.
- Cited harness and registry symbols: `gittest.StubGit`, `gittest.StubGitDir`, `conformanceChecks`, `registry.Checks`, `registry.InputGoSource`, `canaryFixtureRegistry`.
- Cited existing tests: `TestWorktreesRejectsBadCommonDirBeforePorcelain`, `TestWorktreesTypesStartFailures`, `TestAppendWorktreeSurfacesClassifyFailure`, `TestListCommandRendersBoundExpiryAsTypedFailure`, `TestEveryRetainedFixtureBitesThroughRegisteredOwner`, `TestCanaryFixtureRegistryClassifiesEveryFixture`, `TestGateTransactionRefusesAnUnheldCache`, `TestNamedCheckHelpListsEverySupportedCheck`, `TestNamedCheckRunsOnlyRegisteredDevScope`.
- Cited existing tests, continued: `TestPrePushHookAllowProtectedPushConfig`, `TestCommonDirReturnsUnvalidatedOutput`, `TestCommonDirKeepsPlainOutputFailure`, `TestEligibilityVerdictProjectsWithoutSecondDecision`, `TestRecoveryPreservesEveryGitVisibleLayerWithoutMovingBranchOrIndex`, `TestReauthorizeCommandRollsBackLockRefreshAndCASLoss`.
- Import edges: `internal/dashboard`, `internal/diff`, `internal/adopt`, `internal/gate`, and `internal/worktree` already import `internal/git`. No new import edge.
- Source-row clauses and occurrences: the source is the `## main` State of `capture/session-handoff.md` under the heading for candidate 1. The table above lists each clause once.
- Promised field labels: the `ResolutionError` subject field. The diagnostic `<file> spells the Git administration flag <flag> outside internal/git`. The profile table row `git-plumbing-owner` with the input source `go-source`.
- Changed-function callers: `CommonDir` has 14 production callers, and each keeps its call. `hooksDir` has four callers: the health inspection, the installer, the removal, and the doctor's two repair branches through the installer. `LeaseFile` has five test callers and its production callers in the worktree package. `indexIdentity` has one caller in the diff snapshot.
- Copy survival: GR10 reds a surviving spelling in a production Go file outside `internal/git`, and GR27 proves the check reds one.

Reader sweep of the two facts. The readers are:

- the twelve directory sites and the five file sites, which take GR10 and the migration tickets
- the three tests with four flag-spelling sites, which take GR18
- the test-side spellings through independent runners, which the first Won't handle excludes
- the doctor snippet and the wrapper, which the second and third Won't handle exclude
- the managed hook asset, which the pin-removal spec deletes
- the `gittest` helper, which the fifth Won't handle excludes
- the two `WorktreeFailure` consumers in the worktree list and the status board, which keep one type check and take GR7 and GR8
- `CONTEXT.md`, whose worktree admin entry term names `<git-common-dir>` and stays
- `roadmap/FT218.md`, which the retirement folds into FT302 at the drain
- `decisions/assets/gate-pipeline-fixture-inventory.md`, a dated asset pinned to its subject commit, which is excluded

The shipped-surface claim words: the check's diagnostic names `internal/git`, which is a repo-only path, and the claim word `spells` sits beside it. The package-core guard reads the shipped guidance files, and the diagnostic lives in a Go test file, so the guard does not read it.

The trust chain is unchanged. The readers run the `git` on `PATH` as every adapter call does today, and the bounded runner owns the deadline.

The review round runs `codex exec` with the reviewer-named model `gpt-6-astra` at low effort with a cap of two iterations. The model sits outside the tier binding, and the reviewer named it for this run on 2026-09-05. Every subagent runs `opus` at low or medium effort.
