# A command run caches its pure git queries

Status: staged

Decision source: reviewer-confirmed current conversation, 2026-10-28, which holds the FT402-git-cache outcome criteria in `.bench/commitment.json` with `roadmap/FT402.md`, and decisions D1 and D2

Verification log: 4 iteration(s) to accept — the spec round took 2, because row GQ34 compared a start count with "fewer than", so its probe could not bite; the ticket round took 2, because ticket 01 left `internal/git/git.go` unfenced and the seam ticket shared a chunk with its consumer.

Implementation approval: approved by the reviewer on 2026-10-28, including the five-ticket breakdown, the chunk split into GQ-C2A and GQ-C2B, and the declared implementation line.

## Problem

The `internal/worktree` test package is a long pole of the gate's test phase. The gate logs of 2026-10-27 and 2026-10-28 record 85.7 to 241.9 s for it.
Its tests start git 263,108 times in one run.
Of those starts, 102,320 repeat an earlier pure query of the same test. A pure query is a location read or a read of full object ids.

`git.Run` already caches three location answers for one verb run. Only four call sites read the verb's run, so most repeats start git again.
Thirteen other worktree sites build their own run. The ledger address, the lease, the build window, and the commitment store start their queries on each call.
The research of 2026-10-28 cites these sites: `internal/worktree/effects.go:53`, `internal/intent/intent.go:34-40`, `internal/commitment/repository/lease.go:28`, and `internal/worktree/identity_component.go:102`.

### The measurement before any code change

The outcome requires a measurement before any code change. The spec stage measured the package at `main` 1d6c96fc1 and with a throwaway prototype cache.
Each run used `env GOFLAGS="-p=2 -parallel=2" bench test --in <worktree> --package ./internal/worktree` on a quiet host. `TMPDIR` was unset, so the kit run used `/dev/shm`.

The spec reads "-p=2" as the background-grade width `-p=2 -parallel=2` (`internal/gate/width.go:12-22`, `.bench/BENCH-reference.md:794`). The reviewer can veto this reading.

| arm | runs, in ms | range | median |
|---|---|---|---:|
| baseline at `main` | 153812, 152437, 148268, 149179, 148502, 146198 | 146.2 to 153.8 s | 148.8 s |
| prototype cache | 130783, 130410, 130489, 130201, 131114 | 130.2 to 131.1 s | 130.5 s |

The sixth baseline run is a control after the prototype runs. The baseline spread is 7.6 s wide.
The median gain is 18.4 s, which is 12.3 percent. The smallest gain, the baseline minimum minus the prototype maximum, is 15.1 s.
That smallest gain is about twice the spread, so the gain criterion `FT402-git-cache.gain` holds.

A counting shim on `PATH` logged each git start of one baseline run. The table gives the starts that repeat inside one test.

| class | starts | repeats inside one test |
|---|---:|---:|
| location: common directory, admin directory, git path, top level | 53,783 | 47,875 |
| reads whose revision operands are all full object ids | 79,403 | 54,445 |
| state: refs, HEAD, index, status, config, registrations | 129,922 | 80,124 |

The largest repeat groups are 34,555 common-directory reads, 13,515 `cat-file blob` reads, and 16,292 single-path `ls-tree` reads. Those tree reads name `.bench/commitment.json`, `reviews/<slug>.md`, `specs/<slug>/spec.md`, and `ROADMAP.md`.
The prototype was an uncommitted patch in the worktree `FT402-cache-prototype`, built by an opus/high write delegate. The spec-stage close discards it.

## Solution

`git.Run` stays the one per-command cache owner. It gains cached forms for the administration path and for reads whose revision operands are all full object ids.
Each verb passes its one run to the ledger, the commitment store, the roadmap revision reads, and its worktree helpers. A nil run caches nothing.

A read whose operand is a ref, HEAD, or a short id starts git on every call. The worktree porcelain and the default branch stay live.
The package-level git functions stay uncached, so each hostile-path test keeps its own starts. One predicate names a full object id.

The build measures again after its first vertical slice and at the end, with the same command and settings. No wall-clock threshold enters the gate.

## User stories

Line: opus / xhigh.
Implementation-line reason: GQ-C2A and GQ-C2B thread one run across the intent, commitment, roadmap, and worktree owner seams. Its 34 store literals and its many ledger callers make a repair round expensive. The prototype settles the shape, and the cache rules take exact start-count rows. The prototype probe of the ref guard stayed silent, so the purity rows are new evidence. The scorecard gives xhigh to a code ticket that crosses owner seams.
Harder chunks: GQ-C2A, GQ-C2B, GQ-C3.

### A command run reads each pure answer once

1. As a Bench operator, I want one command run to start each pure git query once, so that a landing spends less time in git.
2. As a maintainer, I want the cache to hold only location answers and full-id reads, so that no answer that a write changes is reused.
3. As a maintainer, I want a ref, HEAD, or short-id operand read live on every call, so that a write inside the run is seen.
4. As a maintainer, I want a failed read to store nothing, so that a later read inside the run starts git again.
5. As a maintainer, I want each stored answer returned as a copy, so that a caller that edits the bytes cannot change a later answer.
6. As a maintainer, I want the key to hold the root spelling and the exact argv, so that two query forms never share an answer.
7. As a maintainer, I want a nil run to cache nothing, so that a caller without a run keeps the current behavior.
8. As a maintainer, I want the package-level git functions to stay uncached, so that each hostile-path test still sees its own starts.
9. As a reviewer, I want each directory answer and each administration path validated again on each read, so that a removed directory still refuses.
10. As a maintainer, I want the porcelain and the default branch read live, so that a registration or a branch move is seen.
11. As a maintainer, I want an ancestry answer of two full ids stored for yes and no, so that a repeated check starts no git.
12. As a maintainer, I want one predicate to name a full object id, so that the cache and the worktree identity checks cannot disagree.
13. As a maintainer of a SHA-256 repository, I want 64-digit ids accepted as full ids, so that the cache serves both object formats.
14. As a maintainer, I want an empty listing at a full tree id stored, so that an absent path costs one start.

### The verb's run reaches every reader in the command

15. As a Bench operator, I want the ledger reads of one verb run to share its run, so that the common directory resolves once.
16. As a Bench operator, I want the commitment store to read through the verb's run, so that its repeated tree reads start git once.
17. As a Bench operator, I want the store's roadmap revision reads to use the same run, so that the planning reads share their answers.
18. As a Bench operator, I want one landing to pass its run to each reader, so that the landing starts each pure query once.
19. As a Bench operator, I want the cleanup, release, reset, and lifecycle paths to pass the run down, so that each verb gains the saving.
20. As a maintainer, I want no unexported worktree helper to build its own run, so that a new helper cannot silently split a command's cache.
21. As a maintainer, I want each exported root-taking form to delegate to its run form, so that one implementation owns each read.
22. As a maintainer, I want two verb runs to share no run, so that no answer crosses a command boundary.

### The gain is measured and kept

23. As a reviewer, I want the spec-stage measurement recorded with its runs, command, and base, so that the premise of the build is visible.
24. As a reviewer, I want a measurement after the first vertical slice, so that a gain inside the spread returns the build to the spec.
25. As a reviewer, I want a final measurement with the same command and settings, so that the retro records the delivered gain.
26. As a maintainer, I want each changed start-count pin to keep its purpose, so that a faster test keeps its failure detection.

### Documentation

27. As a teammate who just walked in, I want one ADR to record the cache rule, so that the scope lives outside the comments.
28. As a maintainer, I want each comment that names the old cache scope corrected, so that the code states the current rule.

### Reviewed exclusions

29. As a reviewer, I want the gate, review-record, preflight, and landing-package readers left uncached, so that this build stays inside the measured paths.
30. As a reviewer, I want a worktree removed and re-created inside one run left out of scope, so that the cache needs no re-creation rule.

## Implementation decisions

### Decisions closed on 2026-10-28

- D1: the spec stage measured with a throwaway prototype before any code change. The prototype never lands.
- D2: "pure" means the location answers plus each read whose revision operands are all full object ids. Refs, HEAD, the index, status, config, and registrations stay live.
- The cache scope is one command run. No package variable holds a run, and two runs share nothing. `internal/git/run.go:10-14` states this rule today.

### The cache owner

`git.Run` keeps its three location answers and their revalidation. It gains these cached forms, each with today's package function as its live form:

| form | cached when | live when |
|---|---|---|
| `AdminPath(root, name)` | always, keyed by root spelling and name | never; the read checks the answer's parent directory again |
| `ReadControlBlob(root, object)` | the object is a full id | any other operand |
| `ListTree(root, args...)` | exactly one revision operand before `--`, and it is a full id with an optional `^{tree}` | a ref, HEAD, a short id, or two revisions |
| `ReadTreeFile(root, tree, path)` | through `ListTree` and `ReadControlBlob` | a ref tree lists live, and its blob id caches |
| `ResolveCommit(root, rev, flags...)` | `rev` is a full id; the key holds the flags | any other operand |
| `IsAncestor(root, a, d)` | both operands are full ids; a yes and a no both store | an error stores nothing |
| `PrimaryCheckout(root)` | the common directory, through the run | the porcelain starts on every call |

The object key is the exact argv, `-C` and the root spelling included. A stored answer goes back as a copy. A failed read stores nothing.
A nil `*Run` starts the query on every call, as the package function does. The run holds its mutex while git runs, as the location answers do now.
`Run.Worktrees` keeps its live porcelain and admin scan. No form caches the default branch, HEAD, status, a ref, or the index.

The prototype's run file reached 312 lines, so `internal/git/run.go` stays under its 400-line budget. The object-read tests go to a new test file, because `run_cache_test.go` has 237 lines of headroom.

### One predicate for a full object id

`git.FullObjectID` reports 40 or 64 hexadecimal digits, of either case. Git resolves an uppercase full id to its object, by an observed run on 2026-10-28.
The shim log held no short-hex operand and no uppercase id, so the predicate keeps every logged pure shape.

| site | old rule | resulting rule | row |
|---|---|---|---|
| prototype `git.FullObjectID` | 40 or 64 lowercase digits | 40 or 64 digits of either case | GQ25 |
| `fullCommitIdentity` (`internal/worktree/worktree.go:572`) | 40 digits of either case | calls `git.FullObjectID`, so a 64-digit id also passes | GQ26 |
| `abbreviatedIdentity` (`internal/worktree/land_identity.go:290-292`), a caller of `fullCommitIdentity` | expands a 4- to 39-digit prefix only to a 40-digit id | also expands such a prefix to a 64-digit id; the 39-digit prefix bound stays | flagged addition |
| `hexIdentity` (`internal/worktree/worktree.go:574`) | a hex prefix that a user types | unchanged; it grades a prefix, not a full id | none |
| `lowercaseHex` (`internal/worktree/worktree.go:136`) | a Bench sha256 fingerprint | unchanged; it grades a fingerprint, not an object id | none |
| `fullCommitID` (`internal/git/checkout.go:73-79`) | a lowercase hex line that `rev-parse` answers with the same bytes | unchanged; it grades a merge-head line by a round trip and has no length rule | none |

### The run reaches every reader

The verb's run is the ambient run that `newAmbient` builds (`internal/worktree/effects.go:51-54`). The build threads it to these owners:

- **Ledger:** `intent` gains a `Repo{Run, Root}` value and its constructor `intent.In(run, root)`. The value's methods are the ledger reads and writes. Each exported root-taking function becomes a one-line delegation with a nil run. The slicer enumerates the ledger callers with `bench consumers` on each exported `intent` function.
- **Commitment store:** `repository.Store` gains a `Run` field. Its ledger, lease, window, and tree reads go through that run. The 34 production `Store{Root` literals that build no run stay uncached.
- **Roadmap revision reads:** `RevisionIndex` and `RevisionDocument` gain run forms that the store calls. The root forms delegate with a nil run. Their raw `git.Output` listings (`internal/roadmap/tree_validation.go:55` and `:73`) move onto `Run.ListTree`, with the argv that `ReadTreeFile` uses.
- **Worktree verbs:** the landing, release, cleanup, reset, lifecycle, and creation paths pass the verb's run to their helpers. `planExplicitWith` reads the admin directory and the index path through the run that it holds.

An unexported worktree helper builds no run. A run construction sits only in `newAmbient` or in an exported entry point, which is one command run.
At base 1d6c96fc1, the package builds 14 runs; 13 sit outside `newAmbient`, and four of those sit in exported functions.
Each prototype `X` and `XIn` pair of unexported helpers collapses to the run form, and each test calls that form with a run or nil.
A yes-or-no ancestry check that moves onto the run calls `IsAncestor` and reads an error as no, as `git.OK` does today. The form table gains no yes-or-no form.
The landed-set cleanup that each landing runs threads with the landing, because the landing row counts its reads.

`Store.policyAt` lists through the run's `ListTree` in place of a raw `git.Output` listing. A malformed listing record then refuses instead of reading as present.

### Seam citations beyond the decision source

| seam | `path:line` | research finding |
|---|---|---|
| object-read forms on `git.Run` | `internal/git/tree.go:154`, `:224`, `:236` | shim record: `cat-file blob` repeats 13,515 times, and single-path `ls-tree` of full ids repeats 16,292 times |
| `policyAt` listing through the run | `internal/commitment/repository/repository.go:250` | slices A and B: the policy read lists through a raw `git.Output` call, then reads the file again |
| `intent.Repo` | `internal/intent/intent.go:34-40` | slice B: each ledger read starts the common-directory query |
| `Store.Run` | `internal/commitment/repository/repository.go:21` | slices A and B: the store is a root-only value, and its policy reads repeat |
| roadmap run forms | `internal/roadmap/tree_validation.go:59`, `:87` | prototype report: the store reads revisions through these two readers |
| run-construction census | `internal/worktree/effect_census_test.go:21` | slice A: 13 helper sites build a fresh run |
| logging pass-through stub mode | `internal/gittest/gittest.go:85-96` | slice C: the stub modes that pass through still refuse one query |

### Documentation

ADR 0033 records the decision: a command run caches its pure git queries, and the ADR states the pure scope and the live set. It names no file path and no code.
The build corrects the cache-scope comments at `internal/git/run.go:5-14` and `:52-54`, and at `internal/git/worktree_admin.go:191-194`.
It also corrects `internal/worktree/effects.go:30-36` and `:48-50`, and `internal/worktree/journey_test.go:245-249` and `:266-270`. Review owns this rule wording, so no row anchors it.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| GQ-C1 / 01-cache-full-object-id-reads.md | `git.Run` caches the pure reads under one predicate, and one ADR records the rule | GQ1–GQ10, GQ13–GQ26, GQ43–GQ45, GQ47 | TestRunReadsRefsLive, planned TestRunObjectReads, TestRunObjectReadsStayLive, TestRunObjectFailures, TestRunNilCachesNothing, TestFullObjectID, TestRunAdminPathRevalidates, TestReauthorizeArgAcceptsAFullObjectID | no |
| GQ-C2A / 02-thread-the-run-through-the-ledger-and-store.md | The ledger, the commitment store, and the roadmap revision reads take a run, and their root forms delegate with a nil run | GQ27–GQ33 | planned TestRepoReadsTheCommonDirOnce, TestStoreReadsAPolicyOnce, TestRevisionIndexReadsOnce | yes |
| GQ-C2B / 03-thread-the-landing.md | One landing, its landed-set cleanup included, shares its run with the ledger, the store, and its helpers | GQ34, GQ41 | planned TestLandingReadsTheDestinationCommonDirOnce | yes |
| GQ-C3 / 04-thread-the-cleanup-and-release-paths.md, 05-thread-reset-lifecycle-and-create.md | Every worktree verb shares its one run, and no helper builds its own | GQ35–GQ40, GQ42, GQ46 | TestExplicitPlanReadsEachLocationOnce, TestOwnerMarkerReadsTheCommonDirOnce, TestPackageTestCountPin, planned TestRunConstructionCensus | yes |

GQ-C2A needs GQ-C1, because it calls the new run forms. GQ-C2B needs GQ-C2A, because the landing calls `intent.In` and `Store.Run`.
GQ-C2A is its own review chunk, because tickets 03, 04, and 05 consume its seam. GQ-C3 needs GQ-C2B, because its helpers share the landing's files.
GQ-C2B is the first vertical slice that a timing can see, so measurement GQ41 runs after its checkpoint. The slicer may rename the tickets; each probe property stays.
### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft402-git-cache-full-20261028","orchestrator_session":"claude:gq-orchestrator-20261028","author_limit":3,"assignments":{"01-cache-full-object-id-reads.md":[{"session":"claude:gq-t01-author-20261028","assignment":"75d15bed1584e842dbdd823f8169be0e","model":"opus","effort":"medium","source":"f503721a88a2137e150db1689d1e7fb89089b810","native_ref":"claude-agent:bench-writer-medium:gq-t01-author-20261028"},{"session":"claude:gq-t01-repair1-20261029","assignment":"75d15bed1584e842dbdd823f8169be0e","model":"opus","effort":"medium","source":"05c5593e1cd3b66e30a284610ff9f273eb5733c1","native_ref":"claude-agent:bench-writer-medium:gq-t01-repair1-20261029","predecessor":"claude:gq-t01-author-20261028","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:gq-t01-author-20261028","preserved":"05c5593e1cd3b66e30a284610ff9f273eb5733c1"}],"02-thread-the-run-through-the-ledger-and-store.md":[{"session":"claude:gq-t02-author-20261029","assignment":"75d15bed1584e842dbdd823f8169be0e","model":"opus","effort":"xhigh","source":"6e11f649fa6d1abeef4bb435dae15ae6854da5bc","native_ref":"claude-agent:bench-writer-xhigh:gq-t02-author-20261029"},{"session":"claude:gq-t02-repair1-20261029","assignment":"75d15bed1584e842dbdd823f8169be0e","model":"opus","effort":"xhigh","source":"2b1619c73c3375442ad6b9e988d55e66797d91d3","native_ref":"claude-agent:bench-writer-xhigh:gq-t02-repair1-20261029","predecessor":"claude:gq-t02-author-20261029","trigger":"user-directed","stopped":"claude-agent:bench-writer-xhigh:gq-t02-author-20261029","preserved":"2b1619c73c3375442ad6b9e988d55e66797d91d3"}],"03-thread-the-landing.md":[{"session":"claude:gq-t03-author-20261029","assignment":"75d15bed1584e842dbdd823f8169be0e","model":"opus","effort":"xhigh","source":"4da7aaa154b4f0676ecfe1200cfc1970ec3a90bb","native_ref":"claude-agent:bench-writer-xhigh:gq-t03-author-20261029"},{"session":"claude:gq-t03-repair1-20261029","assignment":"75d15bed1584e842dbdd823f8169be0e","model":"opus","effort":"xhigh","source":"21f23f5e2470cdfc11f6e873967b4b12e6c675e5","native_ref":"claude-agent:bench-writer-xhigh:gq-t03-repair1-20261029","predecessor":"claude:gq-t03-author-20261029","trigger":"user-directed","stopped":"claude-agent:bench-writer-xhigh:gq-t03-author-20261029","preserved":"21f23f5e2470cdfc11f6e873967b4b12e6c675e5"}],"04-thread-the-cleanup-and-release-paths.md":[{"session":"claude:gq-t04-author-20261029","assignment":"75d15bed1584e842dbdd823f8169be0e","model":"opus","effort":"xhigh","source":"0e7ade0ade73f6296f79ecced372f3b3f9032a42","native_ref":"claude-agent:bench-writer-xhigh:gq-t04-author-20261029"}],"05-thread-reset-lifecycle-and-create.md":[{"session":"claude:gq-t05-author-20261029","assignment":"75d15bed1584e842dbdd823f8169be0e","model":"opus","effort":"xhigh","source":"556fab303dd546badbe362e4a3e97a68a9a09e26","native_ref":"claude-agent:bench-writer-xhigh:gq-t05-author-20261029"},{"session":"claude:gq-t05-repair1-20261029","assignment":"75d15bed1584e842dbdd823f8169be0e","model":"opus","effort":"xhigh","source":"459408aa267239e2ad34cf0801ec614522b46071","native_ref":"claude-agent:bench-writer-xhigh:gq-t05-repair1-20261029","predecessor":"claude:gq-t05-author-20261029","trigger":"user-directed","stopped":"claude-agent:bench-writer-xhigh:gq-t05-author-20261029","preserved":"459408aa267239e2ad34cf0801ec614522b46071"}]}},"chunks":[{"id":"GQ-C1","tickets":["01-cache-full-object-id-reads.md"],"verification":[{"id":"probe-01-hit","command":"bench probe internal/git/run.go --swap 'if raw, ok := r.objects[key]; ok {' --with 'if raw, ok := r.objects[key]; false && ok {' --package ./internal/git --run TestRunObjectReads","probe":"Never answer an object read from the store. Each second read then starts git, so the object-read test fails.","ticket":"01-cache-full-object-id-reads.md"},{"id":"probe-01-live","command":"bench probe internal/git/run.go --swap 'if !FullObjectID(object) {' --with 'if false {' --package ./internal/git --run TestRunObjectReadsStayLive","probe":"Cache a blob read whose operand is a ref. The second read of the moved ref then returns the old bytes, so the live-read test fails.","ticket":"01-cache-full-object-id-reads.md"}]},{"id":"GQ-C2A","tickets":["02-thread-the-run-through-the-ledger-and-store.md"],"verification":[{"id":"probe-02","command":"bench probe internal/intent/intent.go --swap 'common, err := r.Run.CommonDir(r.Root)' --with 'common, err := git.CommonDir(r.Root)' --package ./internal/intent --run TestRepoReadsTheCommonDirOnce","probe":"Read the ledger address without the run. Two ledger reads then start the common-directory query twice, so the ledger test fails.","ticket":"02-thread-the-run-through-the-ledger-and-store.md"}]},{"id":"GQ-C2B","tickets":["03-thread-the-landing.md"],"verification":[{"id":"probe-03","command":"bench probe internal/worktree/land.go --swap 'commitrepo.Store{Root: root, Run: a.run}).AdmitWindow(publication)' --with 'commitrepo.Store{Root: root}).AdmitWindow(publication)' --package ./internal/worktree --run TestLandingReadsTheDestinationCommonDirOnce","probe":"Build the window store without the verb's run. The landing then starts more destination common-directory queries, so the landing test fails.","ticket":"03-thread-the-landing.md"}]},{"id":"GQ-C3","tickets":["04-thread-the-cleanup-and-release-paths.md","05-thread-reset-lifecycle-and-create.md"],"verification":[{"id":"probe-04","command":"bench probe internal/worktree/subshell.go --swap 'ledgerBytes, err := intent.In(run, root).LifecycleEvidence()' --with 'ledgerBytes, err := intent.In(nil, root).LifecycleEvidence()' --package ./internal/worktree --run TestExplicitPlanReadsEachLocationOnce","probe":"Read the lifecycle evidence without the run. The explicit plan then starts the common-directory query twice, so the plan test fails.","ticket":"04-thread-the-cleanup-and-release-paths.md"},{"id":"probe-05","command":"bench probe internal/worktree/reset.go --insert-before 'func planReset(run *git.Run, ' --with 'func probeFreshRun() *git.Run { return &git.Run{} }; ' --package ./internal/worktree --run TestRunConstructionCensus","probe":"Add an unexported helper that builds its own run. The census then names the helper, so the census test fails.","ticket":"05-thread-reset-lifecycle-and-create.md"}]}],"final_verification":[{"id":"whole-project","command":"bench gate"},{"id":"coverage","command":"bench coverage --check specs/git-query-cache/spec.md"}]}
```

Each probe needle quotes the prototype text. When the landed text differs, the author amends the needle in the plan and keeps the probe property.

## Testing decisions

A good test drives a real repository through the real `git.Run` form or the real verb helper. It counts git starts through the stub log of `gittest.StubGitDir`.
The stub's `fail-git-path` mode passes every other query to the real git, so it serves the object-read rows in `internal/git` (precedent `internal/worktree/journey_test.go:253`).
That mode refuses each `--git-path` query (`internal/gittest/gittest.go:91-96`), so the `AdminPath` rows GQ22 and GQ23 and the landing row need a mode that refuses no query.
The stub gains a `log` mode that records each start and runs the real git. Ticket 01 adds the mode, and the landing ticket uses it.

- The cache rules attach at `git.Run` in `internal/git`. The precedent is `internal/git/run_cache_test.go`, which counts starts with `gittest.StubGitStarts`.
- The ledger rows attach at `intent.Repo` in `internal/intent`. The store rows attach at `repository.Store` in `internal/commitment/repository`.
- The landing row attaches at the in-process landing that `internal/worktree/land_journey_test.go` drives, with the stub on the process `PATH`.
- The census row reads the package's production files with `go/ast`, as `internal/worktree/effect_census_test.go` reads them with a pattern.

A test that binds `PATH` is serial in `internal/worktree`. The landing row adds one serial test, so `worktreeSerialCeiling` rises from 21 to 22.
Rows GQ26, GQ34, and GQ38 add three top-level tests, so `worktreeTestCount` rises from 877 to 880. Both pins sit in `internal/worktree/parallel_census_test.go:906` and `:926`.
Each ticket that adds a test raises the count pin in its own commit. The pin reads 878 after ticket 01, 879 after ticket 03, and 880 after ticket 05.
Row GQ45 extends `TestRunReadsRefsLive`, which today reads the ref through the package `ResolveCommit` only.
The start-count pin at `internal/worktree/journey_test.go:278` reads 3 today and 1 after GQ-C3 (GQ35). The pin at `:257` stays 1.
The rows in `internal/intent` and `internal/commitment/repository` meet no serial or count pin.
The rows in `internal/git` meet one cap: `internal/conformance/ordinary_build_census_test.go:70-84` allows one repository or process site in the package's test files.
That one site is `internal/git/testhelpers_test.go:18`. So each new `internal/git` test uses the existing repository helpers and the `gittest` stub, and it starts no command of its own.

The gate observes every row through its package tests. The timing rows are review-owned, and no row adds a gate check.

### Seam diagram

    trigger: a worktree verb run (land, clean, release, reset)
        │
        ▼
    newAmbient builds one *git.Run  ──▶  [ ledger Repo, commitment Store, roadmap reads, worktree helpers ]  ──▶  git starts
                      ◀ tests attach here: the stub log counts each start for one root and one argv

    trigger: a git.Run form call
        │
        ▼
    root spelling + argv  ──▶  [ git.Run: purity predicate, store, copy, revalidation ]  ──▶  answer or live start
                      ◀ tests attach here: two calls on one run, with a ref move or a failure between them

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| GQ1 | 1, 2 | Two `ReadControlBlob` reads of one full blob id on one run start `cat-file` once | `internal/git/run_object_test.go` (`TestRunObjectReads`) | A form that never stores starts two reads |
| GQ2 | 3 | Two `ReadControlBlob` reads of `HEAD:<path>` on one run, with a commit between them, return the new bytes | `internal/git/run_object_test.go` (`TestRunObjectReadsStayLive`) | A form that caches a ref operand returns the old bytes |
| GQ3 | 1, 2 | Two `ListTree` reads of one full tree id with one path start `ls-tree` once | `internal/git/run_object_test.go` (`TestRunObjectReads`) | A form that never stores starts two listings |
| GQ4 | 3 | Two `ListTree` reads of `HEAD` with a commit between them list the new tree | `internal/git/run_object_test.go` (`TestRunObjectReadsStayLive`) | A form that caches a ref operand lists the old tree |
| GQ5 | 3 | Two `ListTree` reads of a seven-digit tree prefix start `ls-tree` twice | `internal/git/run_object_test.go` (`TestRunObjectReadsStayLive`) | A predicate that accepts any hex string caches the short id |
| GQ6 | 14 | Two `ListTree` reads of an absent path at one full tree id return an empty listing and start `ls-tree` once | `internal/git/run_object_test.go` (`TestRunObjectReads`) | A store that skips an empty answer starts two listings |
| GQ7 | 1 | Two `ReadTreeFile` reads of one path at one full tree id start `ls-tree` once and `cat-file` once | `internal/git/run_object_test.go` (`TestRunObjectReads`) | A file read that bypasses the run forms starts two of each |
| GQ8 | 1, 2 | Two `ResolveCommit` calls with one full commit id and one flag set start `rev-parse` once | `internal/git/run_object_test.go` (`TestRunObjectReads`) | A form that never stores starts two resolutions |
| GQ9 | 3 | Two `ResolveCommit` calls with a branch name, with an `update-ref` between them, return the moved commit | `internal/git/run_object_test.go` (`TestRunObjectReadsStayLive`) | The guard mutation of the prototype probe returns the old commit |
| GQ10 | 6 | A `ResolveCommit` call with `--quiet` and one without it, for one full id, start `rev-parse` twice | `internal/git/run_object_test.go` (`TestRunObjectReads`) | A key without the flags answers the second call from the first |
| GQ13 | 11 | Two `IsAncestor` checks of one ancestor pair of full ids return yes and start `merge-base` once | `internal/git/run_object_test.go` (`TestRunObjectReads`) | A form that never stores a yes starts two checks |
| GQ14 | 11 | Two `IsAncestor` checks of one non-ancestor pair of full ids return no and start `merge-base` once | `internal/git/run_object_test.go` (`TestRunObjectReads`) | A form that stores only a yes starts two checks |
| GQ15 | 3 | Two `IsAncestor` checks with a branch operand, with an `update-ref` between them, return the changed answer | `internal/git/run_object_test.go` (`TestRunObjectReadsStayLive`) | A form that caches a ref operand returns the old answer |
| GQ16 | 4 | Two `IsAncestor` checks that name a missing full id return an error and start `merge-base` twice | `internal/git/run_object_test.go` (`TestRunObjectFailures`) | A store of the failure starts one check |
| GQ17 | 4 | Two `ReadControlBlob` reads of a missing full id return an error and start `cat-file` twice | `internal/git/run_object_test.go` (`TestRunObjectFailures`) | A store of the failure starts one read |
| GQ18 | 5 | A change to the bytes of a first `ReadControlBlob` answer leaves the second answer equal to the blob | `internal/git/run_object_test.go` (`TestRunObjectReads`) | A form that returns the stored slice gives the changed bytes |
| GQ19 | 7 | Two calls of each cached form on a nil run start the query twice | `internal/git/run_object_test.go` (`TestRunNilCachesNothing`) | A nil run that stores panics or starts once |
| GQ20 | 8 | Two package-level `ReadControlBlob` reads of one full id start `cat-file` twice | `internal/git/run_object_test.go` (`TestRunNilCachesNothing`) | A package function that caches through a shared run starts once |
| GQ21 | 22 | Two run values each start one `cat-file` for one full id | `internal/git/run_object_test.go` (`TestRunObjectReads`) | A cache shared across runs starts once |
| GQ22 | 9 | Two `AdminPath` reads of one name on one run start `--git-path` once | `internal/git/run_object_test.go` (`TestRunAdminPathRevalidates`) | A form that never stores starts two queries |
| GQ23 | 9 | After the parent directory of an `AdminPath` answer is removed, the second read refuses and starts no query | `internal/git/run_object_test.go` (`TestRunAdminPathRevalidates`) | A form with no revalidation returns the removed path |
| GQ24 | 10 | Two `PrimaryCheckout` calls on one run start the common-directory query once and the porcelain twice | `internal/git/run_object_test.go` (`TestRunObjectReadsStayLive`) | A form that caches the porcelain starts it once |
| GQ25 | 12, 13 | `FullObjectID` accepts 40 and 64 hexadecimal digits of either case | `internal/git/run_object_test.go` (`TestFullObjectID`) | A 40-only or lowercase-only predicate refuses a valid id |
| GQ47 | 12 | `FullObjectID` refuses 39, 41, 63, and 65 digits, a value with one non-hex byte, and the empty value | `internal/git/run_object_test.go` (`TestFullObjectID`) | A length-only or prefix-only predicate accepts a short id |
| GQ26 | 12 | The reauthorize hint quotes a 64-digit full id instead of the placeholder | `internal/worktree/reauthorize_test.go` (`TestReauthorizeArgAcceptsAFullObjectID`) | A surviving 40-only `fullCommitIdentity` copy prints the placeholder |
| GQ27 | 15 | Two ledger reads through one `intent.Repo` with a run start the common-directory query once | `internal/intent/repo_test.go` (`TestRepoReadsTheCommonDirOnce`) | A ledger address that reads without the run starts twice |
| GQ28 | 21 | Two root-taking `intent.Read` calls start the common-directory query twice | `internal/intent/repo_test.go` (`TestRepoReadsTheCommonDirOnce`) | A root form that holds a shared run starts once |
| GQ29 | 16 | Two policy reads through one `repository.Store` with a run, at one full tree id, start `ls-tree` once and `cat-file` once | `internal/commitment/repository/run_test.go` (`TestStoreReadsAPolicyOnce`) | A store read that bypasses its run starts two of each |
| GQ30 | 16 | Two ledger reads through one `repository.Store` with a run start the common-directory query once | `internal/commitment/repository/run_test.go` (`TestStoreReadsAPolicyOnce`) | A store ledger that drops its run starts twice |
| GQ31 | 21 | Two policy reads through a `repository.Store` with no run start `ls-tree` twice | `internal/commitment/repository/run_test.go` (`TestStoreReadsAPolicyOnce`) | A store that builds a hidden run starts once |
| GQ32 | 16 | A policy read at a tree that lacks `.bench/commitment.json` reports absence through the run's listing | `internal/commitment/repository/run_test.go` (`TestStoreReadsAPolicyOnce`) | A parsed listing that reads empty as present gives a policy |
| GQ33 | 17 | Two `RevisionIndex` reads through one run at one full tree id start `ls-tree` once | `internal/roadmap/tree_validation_test.go` (`TestRevisionIndexReadsOnce`) | A roadmap read that bypasses the run starts twice |
| GQ34 | 18 | One in-process landing under the `log` stub starts the destination common-directory query exactly the threaded count that the test pins with `==` | `internal/worktree/land_run_test.go` (`TestLandingReadsTheDestinationCommonDirOnce`) | A landing that drops the run from any one store or helper that the fixture reaches starts more queries than the pin |
| GQ35 | 19 | One explicit plan for an owned target starts the common-directory query once, where the pin reads 3 today | `internal/worktree/journey_test.go` (`TestExplicitPlanReadsEachLocationOnce`) | A ledger read outside the run raises the count above one |
| GQ36 | 19 | One explicit plan for an owned target starts the admin-directory query once for the target | `internal/worktree/journey_test.go` (`TestExplicitPlanReadsEachLocationOnce`) | The package `git.AdminDir` call in `planExplicitWith` starts a second query |
| GQ37 | 22 | Two verb runs hold two distinct run values | `internal/worktree/journey_test.go` (`TestOwnerMarkerReadsTheCommonDirOnce`) | A run held in a package variable gives one pointer |
| GQ38 | 20 | A production file of `internal/worktree` builds a run, as `&git.Run{}` or `new(git.Run)`, only in `newAmbient` or in an exported function | `internal/worktree/effect_census_test.go` (`TestRunConstructionCensus`) | A surviving helper with either run form is named by the census |
| GQ39 | 20 | The census names the file and line of a planted run construction of each form in an unexported function | `internal/worktree/effect_census_test.go` (`TestRunConstructionCensus`) | A census that reads no function scope passes the planted line |
| GQ40 | 26 | At the GQ-C3 checkpoint, `worktreeTestCount` reads 880 and `worktreeSerialCeiling` reads 22 | `internal/worktree/parallel_census_test.go` (`TestPackageTestCountPin`) | A pin left at its old value reds the census tests |
| GQ41 | 24 | After the GQ-C2B checkpoint, five timed runs give a smallest gain outside the 7.6 s baseline spread | review-owned: the coordinator's five runs with the spec command | A slice whose gain falls inside the spread returns the build to the spec |
| GQ42 | 25 | At completion, five timed runs with the same command give a smallest gain outside the baseline spread | review-owned: the coordinator's five runs, recorded in the retro | A build that loses the gain shows a median inside the spread |
| GQ43 | 28 | The `internal/git` comments named in Documentation state the run's current cache scope | review-owned: the Standards axis reads each named comment | A comment that names only the location answers survives |
| GQ44 | 27 | ADR 0033 states the pure scope, the live set, and the one-run rule, with no file path and no code | review-owned: the Standards axis and the `prose` check | An ADR that names a path or omits the live set fails review |
| GQ45 | 10 | Inside one run, an `update-ref` of the default branch is seen by `ResolvedDefault` and by the run's `ResolveCommit` of the branch | `internal/git/run_cache_test.go` (`TestRunReadsRefsLive`) | A default-branch or ref answer held in the run returns the old tip |
| GQ46 | 28 | The `internal/worktree` comments named in Documentation state the run's current cache scope | review-owned: the Standards axis reads each named comment | A comment that says a ledger read stays live survives |

Not covered: story 23 — the spec itself records the spec-stage measurement in Problem, so no build row grades it
Not covered: story 29 — a reviewed exclusion; the edge inventory holds its Won't handle line
Not covered: story 30 — a reviewed exclusion; the edge inventory holds its Won't handle line

### Edge inventory

The hostile-input classes of `projects/benchkit.md` apply as follows.

- A root spelling with spaces, glob bytes, or a symbolic-link component keys as exact bytes. Two spellings of one repository miss each other and start twice, which is correct.
- A dangling symlink or a special file at a cached directory refuses on the next read. The revalidation of GQ23 and LGP11 runs on each read.
- A write inside the run that moves a ref, HEAD, or the index stays visible. Each such operand is live (GQ2, GQ4, GQ9, GQ15, GQ45).
- An uppercase full id is pure, because git resolves it to the same object (GQ25).
- An absent file and an empty listing are distinct from a failure: the empty listing stores (GQ6), and the failure stores nothing (GQ16, GQ17).
- A destructive worktree state, a removed admin entry, refuses on the next `AdminPath` read (GQ23).
- `ReadTreeFile` composes two locked reads in sequence, so no cached form re-enters the run's mutex.

**Won't handle:** a worktree removed and re-created at one root inside one run — creation builds its own run in its exported entry (GQ38).
**Won't handle:** an object pruned or replaced inside one run — no Bench verb writes `refs/replace/` or prunes objects during a verb run.
**Won't handle:** a `PATH` swap inside one run — the run stays one command's value, and each stub test builds a fresh run (GQ21, GQ37).
**Won't handle:** the gate, review-record, preflight, and landing-package readers of a landing — `TestLandingReadsTheDestinationCommonDirOnce` still grades the worktree side.
**Won't handle:** `show <oid>:<path>`, `diff <oid> <oid>`, `rev-parse <oid>:<path>`, and `ls-tree --name-only` of a full id — their callers keep live starts.

## Out of scope

- The serial worktree tests face of FT402, which moves serial tests into the parallel phase — 6 edits, 2 gate runs.
- A run for the gate, review-record, preflight, and landing-package readers that a landing calls — 18 edits, 3 gate runs.
- Cached forms for `show`, `diff`, `rev-parse <oid>:<path>`, and `ls-tree --name-only` of full ids — 6 edits, 2 gate runs.
- Any cache of a state read, such as a ref, HEAD, or the porcelain — 8 edits, 2 gate runs, and a reviewer decision.

## Further notes

### Source and caller inventory

- Fresh run constructions in production at base 1d6c96fc1: 18. Three sit in `internal/git/worktree_admin.go`, one in `internal/poolkey/poolkey.go`, one in `newAmbient`, and 13 in other `internal/worktree` helpers.
- The exported construction sites are `Descriptor` and `CheckoutUnprotected` in `internal/worktree/protection.go`, and `PlanExplicit` and `PlanExplicitWithOptions` in `internal/worktree/subshell.go`.
- The unexported construction sites are `planReset`, `planAutomaticAt`, `identityBundleRefusal`, `planLandedAssignment`, `validateCreationBundle`, `validateOwnerMarker`, `acquireAt`, `releaseWith`, and `planExplicitCleanupRow`.
- Package-level location calls in `internal/worktree` production files: 22, by `rg 'git\.(CommonDir|AdminDir|AdminPath|RootAt|Root|Worktrees|PrimaryCheckout)\('`.
- Production `Store{Root` literals: 34, by `rg 'Store\{Root' --glob '!*_test.go'`. The ledger callers come from `bench consumers` on each exported `intent` function, which the slicer runs.
- Start-count readers: `internal/git/run_cache_test.go:36`, `:56`, `:99`, `:129`, `:144`, `:160`; `internal/poolkey/run_test.go:24`; `internal/worktree/journey_test.go:257` and `:278`; `internal/diff/identity_test.go:202-208`.
- Recorded-argv absence readers: `internal/git/worktree_admin_hostile_test.go:155`, `:229`, `:270`, `:286`; `internal/worktree/clean_landed_hostile_test.go:182`, `:186`; `internal/systemtest/land_route_test.go:51`; `internal/adopt/link_hook_test.go:170-189`. The package functions stay uncached, so these keep their counts.
- Other enforcement: `internal/conformance/git_plumbing_owner_test.go:23` keeps the location flags inside `internal/git`, and `internal/conformance/bounds_policy_test.go:45` pins two tokens in `internal/git/git.go`.

The build reruns `bench consumers` on each changed function before each charge, unexported functions included.

### Pre-review proof checklist

- Cited symbols: each cited symbol resolves at base 1d6c96fc1, except the new run forms, `intent.Repo`, `intent.In`, `Store.Run`, `git.FullObjectID`, the `log` stub mode, and the planned tests.
- Import edges: `internal/intent`, `internal/commitment/repository`, and `internal/roadmap` already import `internal/git`. The build adds no new import edge.
- Source-row clauses and occurrences: the source-to-row table below.
- Promised field labels: `Run`, `Repo`, `In`, `FullObjectID`, and the stub mode `log`.
- Changed-function callers: the inventory above. The slicer reruns `bench consumers` for each changed function.
- Copy survival: GQ26 reds a surviving 40-only `fullCommitIdentity`, and GQ38 reds a surviving helper run.
- Rendered-shape readers: none. No output shape changes.
- Pin operators: `TestPackageTestCountPin` compares with `<` and `>`, and `TestSerialSetStaysBelowTheCeiling` requires equality. Every start-count row, GQ34 included, compares with `==`.
- Entry reads: none. The run comes from `newAmbient`, which already owns the ambient reads.
- Derived expectations: each start count comes from the fixture's own calls. GQ34 pins the threaded landing's observed count, and its test doc records the larger main-tip count as context only.
  A GQ-C3 ticket that lowers the landing's count sets the GQ34 pin to the new observed count in its own commit.
- Consolidated rules: the predicate consolidation table above.
- Quantified obligations: GQ19 covers each cached form of the form table, and GQ38 covers each production file of `internal/worktree`.
- Workflow-step writes: none. The cache writes no file.

### Source-to-row table

| source sentence | rows |
|---|---|
| Before any code change, a measurement compares the package time with a per-command cache against the recorded spread. | Problem, GQ41, GQ42 |
| When the measured gain falls outside the spread, a command run caches its pure git queries. | GQ1–GQ10, GQ13–GQ40, GQ43–GQ47 |
| D2: pure means the location answers plus each read whose revision operands are all full object ids. | GQ1–GQ10, GQ13–GQ25, GQ45, GQ47 |
| D2: refs, HEAD, the index, status, config, and registrations stay live. | GQ2, GQ4, GQ5, GQ9, GQ15, GQ24, GQ45 |
| FT402: a faster test keeps its purpose and its failure detection. | GQ35, GQ36, GQ40 |
| FT402: a timing gain inside the recorded spread counts as no change. | GQ41, GQ42 |

### Flagged additions

- ADR 0033 records the cache rule. Invariant 3 needs the decided scope outside code comments.
- The `log` mode of the git stub. The landing row needs a stub that refuses no query.
- The run-construction census `TestRunConstructionCensus`. It makes the one-run rule red-capable.
- The 64-digit widening of the reauthorize hint. It follows from the one-predicate consolidation.
- The 64-digit widening of `abbreviatedIdentity`: a short prefix in a SHA-256 repository now expands to its full id. It follows from the one-predicate consolidation.
- The parent-directory check of each cached `AdminPath` read. A removed administration entry then refuses, as each cached location answer does.
- The pin changes `worktreeSerialCeiling` 21 to 22 and `worktreeTestCount` 877 to 880.
- The `policyAt` listing through `ListTree`. A malformed listing record then refuses instead of reading as present.
- The build-entry staleness pass reruns the baseline when `internal/git` or `internal/worktree` changed on `main` after 1d6c96fc1.
- New test files go into directories above their file-count budget: `internal/git`, `internal/intent`, `internal/commitment/repository`, `internal/roadmap`, and `internal/worktree`. The lane grades file growth only, so this debt does not go red.

### Claims marked uncertain

- The prototype names of the planned tests and the probe needles are uncertain until each ticket lands; the plan amendment keeps each probe property.
- The out-of-scope prices are uncertain estimates from the prototype report's uncached-site list.
