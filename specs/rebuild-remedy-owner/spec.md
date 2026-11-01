# One owner for the rebuild remedy of a pool worktree

Status: staged

Roadmap: FT344

Decision source: named reviewed artifact — the bound criterion `FT344-remedy.one-owner` in `.bench/commitment.json`, read with `roadmap/FT344.md` at its pinned identity `sha256:b927fd97a505b7ffdb373a8fb7eff0ebce1a5274ead4f6777a835efb8278ab04`, plus the reviewer decision of 2026-10-31 in the `/bench-write-spec FT344` conversation (D1: the shell copy keeps the `cd` form).

Verification log: 3 iteration(s) to accept — the spec round took 2 iterations: iteration 1 rejected on RB12, which sat in the RB-C2 acceptance cell while RB-C1 claimed it, and iteration 2 accepted the repair. The ticket round took 1 iteration and folded five prose findings. Reviewer sign-off 2026-10-31.

## Problem

When a Bench binary is stale, Bench prints a rebuild remedy.
One owner composes that remedy for every Go caller: `freshness.RebuildAction` returns `cd <root> && bash scripts/go-build.sh <root> <root>/dist/bench` (`internal/freshness/freshness_verify.go:120-122`).
`(*freshness.Refusal).Error` appends the same remedy after `; rebuild with ` (`internal/freshness/freshness_verify.go:106-108`).

The hooks refuse a `cd` into the worktree pool path.
The follow-on guard runs `benchguard.PoolReference` over each Bash command before any other verdict (`cmd/bench/guards.go:171`).
That function refuses a `cd` argument that starts with `<Bench home>/worktrees/` (`internal/benchguard/pool.go:18-40`, `internal/poolkey/poolkey.go:44-46`).
So when the stale binary belongs to a pool worktree, the printed remedy cannot run.

Two callers already work around this.
The `binary-seal` row of the build preflight prints `Refusal.Reason` and names `bench worktree build <assignment id>` (`internal/preflight/binary_seal.go:55-68`).
The tree-target build refusal names `bench worktree build <label>` (`internal/worktree/tree_target.go:90-93`).
Both fill the one build grammar through `usage.WorktreeBuildFor` (`internal/usage/worktree.go:75-78`).
But the owner itself cannot tell a pool root from the primary checkout, so every other caller still prints the `cd` form for a pool root:

- the command entry: `bench commands --brief` prints `RebuildAction(git.Root())` for a stale executable inside the cwd checkout (`cmd/bench/commands.go:32-33`);
- the live-binary check: a release or clean that removes the running `dist/bench` of a pool worktree prints `RebuildAction(plan.Target)` (`internal/worktree/live_binary.go:22-29`, `internal/worktree/clean.go:304`);
- each printer of `Refusal.Error`: the doctor seal row (`internal/adopt/doctor_rows.go:257-258`), the `freshness-check` verb (`cmd/bench/freshness_publish.go:26-27`), and the run-binary selection, whose refusal reaches the gate and `bench test` (`internal/runbinary/runbinary.go:174-183`, `internal/testreport/command.go:50`).

The `binary-seal` row decides the root kind a second way: by assignment ownership (`internal/preflight/gather.go:378-384`), not by the pool path that the hook grades.
That is a second source for one fact.

The landing caller, which FT344 names, cannot print a pool remedy.
`brokerChangeNotice` names the rebuild only for the kit source checkout (`internal/worktree/land.go:387-396`).
The landing destination must be checked out on the default branch (`internal/worktree/land_identity.go:51-53`), so its root is never a pool worktree.
It already reads the owner, so it stays correct without a change.

## Solution

The owner decides the root kind with the hook's own pool test.
For a root whose `cd` remedy the pool guard refuses, the owner returns `bench worktree build '<root>'`.
For every other root, the owner returns the current `cd` form, unchanged.
Every Go caller reads the remedy from that owner, so each pool-root remedy names the build verb.

The build verb accepts an absolute worktree path as its target and resolves it to the assignment that owns that path (`internal/worktree/path.go:171-175`, `internal/worktree/path.go:300-317`).
The pool guard allows a pool path as a plain argument (`internal/benchguard/pool.go:15-17`), and the follow-on rule allows one bare Bench call.
So the new remedy runs past both guards.

The `binary-seal` row reads the root kind and the remedy from the owner and stops deciding by ownership.

The shell copy `bench_rebuild_action` keeps the `cd` form (D1).
It prints only when the Go core is unreachable or stale (`.bench/lib/resolve-bench.sh:354-360`).
In both states the follow-on hook does not run the pool guard: an unreachable core allows every call (`.bench/hooks/block-bench-follow-on.sh:24`), and a stale core passes every non-Bench call and blocks every Bench call (`.bench/hooks/block-bench-follow-on.sh:42-53`).
So in those states the `cd` form runs, and `bench worktree build` does not.

## User stories

Line: opus / medium.
Implementation-line reason: the hardest material chunk is RB-C1, which changes the one owner that about twenty callers and pins read, and which must keep the system suite green when the kit under test is itself a pool worktree. The spec fixes the decision rule, the exact remedy text, and every seam, and each seam has a precedent. The owner rows grade the remedy against the real hook predicates, so a wrong form turns the ordinary Go test phase red. RB-C2 migrates two callers at known seams.
Harder chunks: RB-C1.

### The owner names the build verb for a pool root

1. As an agent whose pool worktree binary is stale, I want the rebuild remedy to name `bench worktree build '<root>'`, so that I can run the remedy that the hooks allow.
2. As an operator in the primary checkout, I want the remedy for my checkout to stay the `cd` form, so that a root that no assignment owns still gets a remedy that runs.
3. As an agent, I want the pool-root remedy to pass the pool guard and the follow-on rule, so that the hook does not refuse the remedy that it caused.
4. As a kit maintainer, I want the owner to decide the root kind with the pool guard's own test, so that the remedy changes exactly where the hook refuses the `cd` form.
5. As an agent, I want a sibling directory such as `<home>/worktrees-old/x` to keep the `cd` form, so that only a true pool path takes the build verb.
6. As an agent, I want an empty Bench home to give the `cd` form, so that an unknown home never classifies a root as a pool root.
7. As an agent with a hostile worktree path, I want the build remedy to keep the root as one shell word, so that a space, a quote, or a glob character cannot split the target.
8. As an agent, I want the owner to read the process Bench home, so that the remedy and the hook grade the same pool path.
9. As a caller with an injected Bench home, I want an owner form that takes the home, so that a test or a verb with its own home gets the same answer as the hook in that home.
10. As an agent, I want a stale-binary refusal for a pool root to carry the build remedy once, so that each printer of the refusal names a remedy that runs.
11. As a kit maintainer, I want the `cd` form to stay available as its own owner function, so that the shell copy has one Go source to mirror.

### Each caller reads the owner

12. As an agent who runs `bench commands --brief` in a stale pool worktree, I want the refusal to name `bench worktree build '<root>'`, so that the liveness probe gives a remedy that runs.
13. As an agent who releases a pool worktree whose `dist/bench` answers `bench`, I want the live-binary warning to name the build remedy, so that the warning names no refused `cd`.
14. As a kit maintainer, I want the live-binary warning to read the remedy with the verb's own Bench home, so that the warning and the release grade the same pool.
15. As an orchestrator, I want the `binary-seal` row of a pool root to print the refusal reason and name the owner's build remedy as its next step, so that the row keeps its current shape.
16. As an orchestrator, I want the `binary-seal` row of any other root to print the whole refusal with no next step, so that the row keeps its current answer there.
17. As a kit maintainer, I want the `binary-seal` row to stop deciding the root kind by assignment ownership, so that the root kind has one source.
18. As an orchestrator, I want the preflight gather to read the remedy with the Bench home that it reads at its boundary, so that the row and the hook grade the same pool.
19. As a kit maintainer, I want the landing notice for the kit source checkout to keep the `cd` form, so that the destination, which is never a pool root, keeps its remedy.

### The shell copy keeps the `cd` form

20. As an agent whose Bench core is stale in a pool-rooted session, I want the follow-on hook to print the `cd` form, so that I get the one remedy that the stale hook passes.
21. As an agent, I want the stale follow-on hook to pass the printed `cd` remedy itself, so that the remedy runs in the state that printed it.
22. As a kit maintainer, I want the shell-copy pins to compare the hook output with the `cd`-form owner, so that the system suite stays green when the kit under test is a pool worktree.
23. As a kit maintainer, I want the shell-copy comment to state why the copy keeps the `cd` form, so that nobody adds a pool test that the stale hook would block.

### The suite stays green in a pool worktree

24. As a kit maintainer, I want each system pin that compares a child's remedy with the owner to grade it with the child's Bench home, so that a run inside a pool worktree does not red on a home mismatch.

### Records

25. As a reader of the release notes, I want one `CHANGELOG.md` entry that names the pool-root remedy, so that the change is visible after the row file retires.

### Reviewed exclusions

26. As a reviewer, I want the shell copy to keep the `cd` form, so that D1 holds.
27. As a reviewer, I want no change to the landing caller, so that the spec does not edit a caller whose root is never a pool root.

## Implementation decisions

### The owner

`internal/freshness` keeps the one owner of the rebuild remedy in `freshness_verify.go`.
The package gets these exports:

- `CheckoutRebuildAction(root string) string` returns the current `cd` form, byte for byte. The current body of `RebuildAction` moves here.
- `Rebuild` is a struct with two fields: `Action string`, the copy-paste remedy, and `Worktree bool`, which is true when `Action` is the build verb.
- `RebuildIn(home, root string) Rebuild` decides the kind. `Worktree` is true when `benchguard.PoolReference(CheckoutRebuildAction(root), poolkey.Pools(home))` is not empty and `home` is not empty. Then `Action` is `usage.WorktreeBuildFor(shellQuote(root))`. Otherwise `Action` is `CheckoutRebuildAction(root)`.
- `RebuildAction(root string) string` returns `RebuildIn(benchhome.Dir(), root).Action`. It is the one ambient read of the home in the owner.

`Refusal.Error` keeps its text shape: `Reason() + "; rebuild with " + RebuildAction(r.RepairRoot)`.
The `Refusal` doc comment states that `Error` carries the owner's remedy for every root, and that `Reason` serves a caller that prints the remedy in a separate cell.

The decision rule asks the pool guard whether it refuses the `cd` form.
The guard reads command text and resolves no path, so the owner and the guard agree on every spelling of a root, a symbolic-link spelling included.
`benchguard` imports `shellcommand` and the standard library alone, and `poolkey` imports `git`. `go list -deps` of both names no `freshness` or `usage` package, so the new edges `freshness` to `benchguard`, `poolkey`, `usage`, and `benchhome` make no cycle.

The build remedy carries the absolute root as its target.
The owner holds no ledger read, so it cannot name the assignment id or label.
The build verb resolves a path target to the assignment that owns that path (`internal/worktree/path.go:171-175`).

### The callers

`cmd/bench/commands.go` keeps its call to `RebuildAction`. The command reads the process home, as the hook does.

`warnBeforeRemovingLiveBinary` calls `freshness.RebuildIn(a.home, root).Action` (`internal/worktree/live_binary.go:22-29`, `internal/worktree/effects.go:39`).
The verb's injected home is the home in which the release grades its pool.

The landing, the doctor row, `freshness-check`, and the run-binary selection keep their calls. Each one reads the owner through `RebuildAction` or `Refusal.Error`.

### The `binary-seal` row

`binarySealFacts` takes the Bench home and also returns `freshness.RebuildIn(home, root)` for the graded root.
The gatherer passes the home that it reads at its boundary through `benchhome.Dir()`.
`Facts` gets one field, `BinarySealRebuild freshness.Rebuild`, beside `BinarySealRefusal` (`internal/preflight/decision.go:164-168`).

`binarySealCheck` keeps its not-applicable and green branches.
For a refusal, a `Worktree` rebuild gives `red(Reason())` with `Next` set to the rebuild `Action`.
Any other rebuild gives `red(Error())` with an empty `Next`.
The check no longer reads `AssignmentTarget`, and `binary_seal.go` no longer imports `usage`.
`AssignmentTarget` stays, because the stale-base remedy, the charge pack, and the stage step read it (`internal/preflight/command.go:110`, `internal/preflight/charge_pack.go:131`, `internal/preflight/preparation.go:58`).
The comment above `binarySealCheck` states the new rule.

The row's `Next` cell changes from `bench worktree build <assignment id>` to `bench worktree build '<root>'`. Both forms run. This is the contestable call of this spec, and it follows from the one-source rule.

### The shell copy

`bench_rebuild_action` keeps its `printf` line, so `checkPublishedExecutablePath` keeps its regular expression (`internal/conformance/build_contracts_test.go:29`).
The comment above `bench_shell_quote`, which covers both shell functions, names `freshness.CheckoutRebuildAction` as the Go source and states that the copy keeps the `cd` form because it prints only while the follow-on guard is off.

### The system pins

`TestSessionStartNamesTheBuildCommandWhenTheCoreIsUnreachable` and `TestBenchFollowOnHookDegradedRim` read the expectation from `freshness.CheckoutRebuildAction`, because they pin the shell copy.
`TestSystemOwnerRefusesAStaleSelectedSeal` reads the expectation from `freshness.RebuildIn` with the Bench home that its child reads, because the child's owner composes the refusal in that home.

### Records

The `## [Unreleased]` section of `CHANGELOG.md` gets one entry that names the pool-root remedy and the `binary-seal` next step.

## Implementation chunks

RB-C1 ships the owner, its guard-graded rows, the command-entry row, and the shell-copy pins. The owner is a seam that RB-C2 consumes, so RB-C1 is its own small chunk, and its review closes before RB-C2 starts.
RB-C2 moves the live-binary check and the `binary-seal` row onto the owner.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| RB-C1 / 01-pool-rebuild-owner.md | The owner names `bench worktree build '<root>'` for each pool root, every `RebuildAction` and `Refusal.Error` printer follows, and the shell copy keeps the `cd` form with pins that stay green in a pool worktree | RB01, RB02, RB03, RB04, RB05, RB06, RB07, RB08, RB09, RB10, RB11, RB12, RB17, RB18, RB19, RB20, RB21, RB22, RB26 | TestRebuildRemedyPassesTheHooks, TestRebuildRemedyKind, TestRebuildRemedyHostileRoot, TestRebuildActionReadsTheProcessHome, TestRefusalCarriesThePoolRemedy, TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths, TestCommandsBriefGradesTheRunningExecutable, TestBenchFollowOnHookDegradedRim, TestSessionStartNamesTheBuildCommandWhenTheCoreIsUnreachable, TestSystemOwnerRefusesAStaleSelectedSeal | yes |
| RB-C2 / 02-callers-read-owner.md | The live-binary warning and the `binary-seal` row read the remedy and the root kind from the owner | RB13, RB14, RB15, RB16, RB23, RB24, RB25 | TestResidueGuardWarnsBeforeRemovingTheLiveBinary, TestBinarySealRemedyNamesTheWorktreeBuildVerb, TestBinarySealFactsGradeARootOnDisk, TestLandCommandReportsInstallStepForABrokerChangingDiff | no |

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft344-build-20261031","orchestrator_session":"session_01JK5CPcwpGroiAPWzzXxytc","author_limit":1,"assignments":{"01-pool-rebuild-owner.md":[{"session":"ft344-author-01","assignment":"efbed6125ab26685cc058f28028b791b","model":"opus","effort":"medium","source":"fba082fb602a1e280c730c41c8e87c759f872903","native_ref":"claude-agent:bench-writer-medium:ft344-author-01"},{"session":"ft344-author-01-rep1","assignment":"efbed6125ab26685cc058f28028b791b","model":"opus","effort":"medium","source":"7b7f22cc7cc677d925d6657ab7562d47b5b1c415","native_ref":"claude-agent:bench-writer-medium:ft344-author-01-rep1","predecessor":"ft344-author-01","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:ft344-author-01","preserved":"7b7f22cc7cc677d925d6657ab7562d47b5b1c415"}],"02-callers-read-owner.md":[]}},"chunks":[{"id":"RB-C1","tickets":["01-pool-rebuild-owner.md"],"verification":[{"id":"owner-kind","command":"bench test --package ./internal/benchguard --package ./internal/freshness --run 'RebuildRemedy|RebuildAction|RefusalCarries|RefusalRebuildAction'","probe":"Make RebuildIn return CheckoutRebuildAction for every root, and record that TestRebuildRemedyPassesTheHooks turns red on the pool-root case.","ticket":"01-pool-rebuild-owner.md"},{"id":"command-entry","command":"bench test --package ./cmd/bench --run 'CommandsBriefGradesTheRunningExecutable'","probe":"Make RebuildAction read an empty home in place of benchhome.Dir(), and record that the pool-root subtest of TestCommandsBriefGradesTheRunningExecutable turns red.","ticket":"01-pool-rebuild-owner.md"},{"id":"shell-copy","command":"bench test --check system","probe":"Add a pool-path case to bench_rebuild_action that prints bench worktree build, and record that TestBenchFollowOnHookDegradedRim turns red on the pool-kit case.","ticket":"01-pool-rebuild-owner.md"}]},{"id":"RB-C2","tickets":["02-callers-read-owner.md"],"verification":[{"id":"live-binary","command":"bench test --package ./internal/worktree --run 'ResidueGuardWarnsBeforeRemovingTheLiveBinary|LandCommandReportsInstallStepForABrokerChangingDiff'","probe":"Call freshness.RebuildAction(root) in warnBeforeRemovingLiveBinary in place of the injected home form, and record that TestResidueGuardWarnsBeforeRemovingTheLiveBinary turns red.","ticket":"02-callers-read-owner.md"},{"id":"binary-seal","command":"bench test --package ./internal/preflight --run 'BinarySeal'","probe":"Restore the AssignmentTarget branch in binarySealCheck, and record that TestBinarySealRemedyNamesTheWorktreeBuildVerb turns red on the owned-checkout case.","ticket":"02-callers-read-owner.md"}]}],"final_verification":[{"id":"owner-packages","command":"bench test --package ./internal/freshness/... --package ./internal/benchguard --package ./internal/preflight/... --package ./internal/worktree --package ./cmd/bench"},{"id":"system","command":"bench test --check system"},{"id":"coverage","command":"bench coverage --check rebuild-remedy-owner"}]}
```

## Testing decisions

- A good owner test calls `freshness.RebuildIn` with an explicit home and grades the returned remedy with the real hook predicates: `benchguard.PoolReference` with `poolkey.Pools(home)` and `benchguard.Classify`. The refusal-route guard test is the precedent (`internal/conformance/refusal_route_guard_test.go:100-135`).
- The owner tests live in a new external test file, `internal/benchguard/rebuild_remedy_test.go` (`package benchguard_test`). `bench structure` grades `internal/freshness` at its 12-file directory maximum today, and it grades `freshness_verify_test.go` within 7 lines of the 400-line file maximum.
- Each expected remedy is a literal in the test, composed from the fixture path and the fixed words `bench worktree build`. No expectation calls `usage.WorktreeBuildFor` or `CheckoutRebuildAction` for a row that grades the new form.
- The command-entry row adds a subtest to `TestCommandsBriefGradesTheRunningExecutable`: the synthetic checkout lies under `<home>/worktrees/<key>/<id>`, and the child runs with `BENCH_HOME=<home>`.
- The live-binary row reuses `newResidueGuardFixture`, whose creation path lies under the fixture home (`internal/worktree/live_binary_test.go:19-31`).
- The `binary-seal` rows drive `Decide` over `Facts`, and the gather row drives `binarySealFacts` with an explicit home, as `TestBinarySealFactsGradeARootOnDisk` does (`internal/preflight/binary_seal_test.go:133`).
- The shell-copy rows extend `TestBenchFollowOnHookDegradedRim` (`internal/systemtest/bench_follow_on_test.go:222`) with a kit path under the hook's `BENCH_HOME`.
- The ordinary Go test phase of the gate observes RB01 to RB16 and RB23 to RB25. The system phase observes RB17, RB18, and RB20. RB19, RB21, RB22, and RB26 are review-owned.

### Seam diagram

    trigger: a stale dist/bench at root (refusal, probe, warning, preflight row)
        │
        ▼
    root, home  ──▶  [ freshness.RebuildIn: pool guard over CheckoutRebuildAction(root) ]  ──▶  Rebuild{Action, Worktree}
                         ◀ tests attach here: rebuild_remedy_test.go calls RebuildIn and grades Action
                           with benchguard.PoolReference and benchguard.Classify
        │
        ├──▶ RebuildAction / Refusal.Error ──▶ commands --brief, doctor, freshness-check, run-binary, landing
        │        ◀ commands_brief_test.go runs the real binary in a pool-path checkout
        ├──▶ live-binary warning (injected home)   ◀ live_binary_test.go release fixture
        └──▶ binary-seal row (gathered Rebuild)     ◀ binary_seal_test.go over Facts

    trigger: the follow-on hook with a stale core
        │
        ▼
    kit_root ──▶ [ bench_rebuild_action (shell copy, cd form) ] ──▶ BLOCKED or WARNING line
                     ◀ tests attach here: bench_follow_on_test.go compares with CheckoutRebuildAction

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| RB01 | 1 | `RebuildIn(home, home+"/worktrees/k/x")` returns `Action` equal to `bench worktree build '<home>/worktrees/k/x'` and `Worktree` true | `internal/benchguard/rebuild_remedy_test.go` (`TestRebuildRemedyKind`) | Today's owner returns the `cd` form for every root |
| RB02 | 2 | `RebuildIn(home, root)` for a root outside `<home>/worktrees/` returns `Action` equal to `cd '<root>' && bash scripts/go-build.sh '<root>' '<root>/dist/bench'` and `Worktree` false | `internal/benchguard/rebuild_remedy_test.go` (`TestRebuildRemedyKind`) | An owner that prints the build verb for every root gives the primary checkout a target that no assignment owns |
| RB03 | 3 | For the RB01 root, `benchguard.PoolReference(Action, poolkey.Pools(home))` is empty and `benchguard.Classify(Action, ...)` does not block | `internal/benchguard/rebuild_remedy_test.go` (`TestRebuildRemedyPassesTheHooks`) | A build remedy that keeps a `cd` lead or adds a follow-on is refused by the hook it must pass |
| RB04 | 4 | For the RB01 root, `benchguard.PoolReference(CheckoutRebuildAction(root), poolkey.Pools(home))` is not empty | `internal/benchguard/rebuild_remedy_test.go` (`TestRebuildRemedyPassesTheHooks`) | A fixture that the guard does not refuse proves nothing about the switch |
| RB05 | 5 | `RebuildIn(home, home+"/worktrees-old/x")` returns the `cd` form and `Worktree` false | `internal/benchguard/rebuild_remedy_test.go` (`TestRebuildRemedyKind`) | A prefix test without the path separator calls a sibling directory a pool root |
| RB06 | 6 | `RebuildIn("", "/worktrees/k/x")` returns the `cd` form and `Worktree` false | `internal/benchguard/rebuild_remedy_test.go` (`TestRebuildRemedyKind`) | A rule that builds the pool path from an empty home classifies a root under `/worktrees/` as a pool root |
| RB07 | 7 | For a pool root whose last segment holds a space, a single quote, and `*`, the `Action` parses through `shellcommand.Parse` to exactly the four words `bench`, `worktree`, `build`, and the root | `internal/benchguard/rebuild_remedy_test.go` (`TestRebuildRemedyHostileRoot`) | An unquoted fill splits the target into several words |
| RB08 | 8 | With `BENCH_HOME` set to `home` in the test process, `RebuildAction(home+"/worktrees/k/x")` returns the RB01 `Action` | `internal/benchguard/rebuild_remedy_test.go` (`TestRebuildActionReadsTheProcessHome`) | An owner that reads a fixed home or ignores `BENCH_HOME` grades a pool that the hook does not use |
| RB09 | 10 | With `BENCH_HOME` set to `home`, a `Refusal` whose `RepairRoot` is a pool root renders `Error` with exactly one `; rebuild with bench worktree build '<root>'` and no `cd ` | `internal/benchguard/rebuild_remedy_test.go` (`TestRefusalCarriesThePoolRemedy`) | An `Error` that keeps `CheckoutRebuildAction` leaves every refusal printer on the refused form |
| RB10 | 11 | `CheckoutRebuildAction(root)` returns the `cd` form for the RB01 pool root | `internal/benchguard/rebuild_remedy_test.go` (`TestRebuildRemedyKind`) | A `cd`-form owner that also switches leaves the shell copy with no Go source to mirror |
| RB11 | 2 | `Refusal.Error` for a root outside the pool still runs under `bash -c` with `$PWD`, `$1`, and `$2` equal to the root, the root, and `<root>/dist/bench` | `internal/freshness/freshness_verify_test.go` (`TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths`) | A kept-rule witness: the checkout remedy must stay copy-paste safe for a hostile path |
| RB12 | 12 | `bench commands --brief` under a mismatched seal in a checkout at `<home>/worktrees/k/x`, run with `BENCH_HOME=<home>`, prints `bench worktree build '<root>'` and no `cd ` | `cmd/bench/commands_brief_test.go` (`TestCommandsBriefGradesTheRunningExecutable`) | A command that composes its own remedy, or reads another home, prints the refused form |
| RB13 | 13 | A release of the residue-guard fixture whose `dist/bench` is the running binary warns with `bench worktree build '<creation path>'` and with no `cd ` and no `scripts/go-build.sh` | `internal/worktree/live_binary_test.go` (`TestResidueGuardWarnsBeforeRemovingTheLiveBinary`) | Today the warning prints the `cd` form for the pool root |
| RB14 | 14 | The RB13 warning holds when the test process `BENCH_HOME` differs from the fixture home | `internal/worktree/live_binary_test.go` (`TestResidueGuardWarnsBeforeRemovingTheLiveBinary`) | A warning that calls `RebuildAction` grades the process home and prints the `cd` form for the fixture pool |
| RB15 | 15 | `Decide` over a refused seal with a `Worktree` rebuild fact and an empty `AssignmentTarget` renders `binary-seal` red with `Detail` equal to `Reason()` and `Next` equal to the fact's `Action`, and no cell holds `cd ` or `scripts/go-build.sh` | `internal/preflight/binary_seal_test.go` (`TestBinarySealRemedyNamesTheWorktreeBuildVerb`) | A row that prints `Error()` for a pool root repeats the refused form |
| RB16 | 16 | `Decide` over a refused seal with a non-worktree rebuild fact renders `binary-seal` red with `Detail` equal to `Error()` and an empty `Next` | `internal/preflight/binary_seal_test.go` (`TestBinarySealRemedyNamesTheWorktreeBuildVerb`) | A row that names the build verb for every root gives the primary checkout a dead next step |
| RB23 | 17 | `Decide` with `AssignmentTarget` set and a non-worktree rebuild fact still renders `Detail` equal to `Error()` and an empty `Next` | `internal/preflight/binary_seal_test.go` (`TestBinarySealRemedyNamesTheWorktreeBuildVerb`) | A row that keeps the ownership branch decides the root kind from a second source |
| RB24 | 18 | `binarySealFacts(home, root)` for a stale seal at `<home>/worktrees/k/x` returns a rebuild with `Worktree` true and `Action` equal to `bench worktree build '<root>'` | `internal/preflight/binary_seal_test.go` (`TestBinarySealFactsGradeARootOnDisk`) | A gather that reads the rebuild with another home, or skips it, leaves the row on the checkout branch |
| RB25 | 19 | The broker notice for the kit source checkout still holds `CheckoutRebuildAction(root)` and `bench doctor --fix` | `internal/worktree/land_broker_notice_test.go` (`TestLandCommandReportsInstallStepForABrokerChangingDiff`) | A kept-rule witness: the landing destination is never a pool root, so its remedy must not change |
| RB17 | 20 | The stale-core follow-on hook with `BENCH_KIT` at `<home>/worktrees/k/x` and `BENCH_HOME=<home>` prints `CheckoutRebuildAction(kit)` in its BLOCKED line for `bench gate` | `internal/systemtest/bench_follow_on_test.go` (`TestBenchFollowOnHookDegradedRim`) | A shell copy that gains a pool test prints a Bench call that the stale hook blocks |
| RB18 | 21 | The same stale hook passes the printed `CheckoutRebuildAction(kit)` command at exit 0 | `internal/systemtest/bench_follow_on_test.go` (`TestBenchFollowOnHookDegradedRim`) | A remedy that the stale hook refuses is a dead remedy in the one state that prints it |
| RB19 | 22 | `TestBenchFollowOnHookDegradedRim` and `TestSessionStartNamesTheBuildCommandWhenTheCoreIsUnreachable` read their expected remedy from `freshness.CheckoutRebuildAction` | review-owned: a search for RebuildAction calls in the two shell-pin test files names only CheckoutRebuildAction | A shell pin that reads `RebuildAction` reds whenever the kit under test lies in the pool |
| RB20 | 24 | `TestSystemOwnerRefusesAStaleSelectedSeal` reads its expected remedy from `freshness.RebuildIn` with the Bench home that its child reads | `internal/systemtest/owner_stale_seal_test.go` (`TestSystemOwnerRefusesAStaleSelectedSeal`) | A pin that grades with the test-process home reds when the kit under test lies in the pool |
| RB21 | 23 | The comment above `bench_shell_quote` names `freshness.CheckoutRebuildAction` and states that the copy keeps the `cd` form while the follow-on guard is off | review-owned: `.bench/lib/resolve-bench.sh:354-360` | A comment that names the switching owner invites a pool test in the shell copy |
| RB22 | 26 | `bench_rebuild_action` keeps its one `printf` line, and `PublishedExecutable` keeps its spelling, so the `published-executable-path` check stays green on the tree | review-owned: the `published-executable-path` conformance check stays green on the tree (`internal/conformance/registry/checks.go:27`) | A kept-rule witness: a rewritten shell copy or a moved `PublishedExecutable` breaks the published-path drift guard |
| RB26 | 4 | `RebuildIn` decides `Worktree` by a call to `benchguard.PoolReference` over `CheckoutRebuildAction(root)`, and no production file under `internal/freshness` holds the literal `worktrees` | review-owned: `rg -n 'worktrees' internal/freshness --glob '!*_test.go'` returns no hit | A raw prefix test of the root passes RB01 to RB10, because quotes fold in the guard's parse, and it is a second source for the pool rule |

Not covered: story 9 — the injected-home form is the seam of RB01 to RB07, RB14, and RB24, and it has no behavior apart from those rows.
Not covered: story 25 — the `CHANGELOG.md` entry is review-owned in the `## [Unreleased]` section and takes no behavior row.
Not covered: story 27 — the reviewed exclusion of a landing change takes RB25 as its kept witness and no new behavior.

### Edge inventory

These behaviors serve the kit and every linked repository, because `internal/freshness`, the hooks, and `bench worktree build` ship.
RB05 and RB06 cover the sibling directory and the empty home.
RB07 covers a root with a space, a quote, and a glob character.
A root with a control character cannot reach the build verb, because the verb refuses a target that no output line can carry (`internal/worktree/path.go:79-81`), and a pool directory name is a generated key and id.
No test swaps a package variable. RB08 and RB09 set `BENCH_HOME` with `t.Setenv`, so those tests do not run in parallel.

Won't handle: a pool root whose assignment is not active, such as a released tree — `bench worktree build` refuses with the assignment state, and the `cd` form is refused by the hook as well, so no runnable remedy is lost.
Won't handle: a pool root of another repository's pool — the build verb reads the ledger of the cwd repository and refuses an unassigned target; the operator runs the remedy from that repository.
Won't handle: a root that reaches the pool only through a symbolic-link spelling — the guard compares text too, so it passes the `cd` form, and the owner prints the `cd` form, which runs.
Won't handle: the live-binary warning after the remove — the discard runs only for a remove action, and the worktree goes away right after it (`internal/worktree/lifecycle.go:548-570`), so neither remedy runs on that tree afterward. This spec makes the warning read the owner; a better target for that warning is a separate capability.
Won't handle: the shell copy in the session-start hook — D1 keeps the `cd` form, and RB19 pins the hook output to `CheckoutRebuildAction`.
Won't handle: the stale `dist/bench` after a ticket commit (`docs/known-issues.md:10`) — FT406 owns that row and its retirement.

## Out of scope

- A live-binary warning that names the rebuild of the checkout that the wrapper resolves after the remove is a separate capability. Estimate: 4 edits, 1 gate run.
- A rebuild remedy that names the assignment id or label in place of the path needs a ledger read in the owner, and is a separate capability. Estimate: 5 edits, 1 gate run.

## Further notes

### Source-sentence-to-row table

| source sentence | rows |
| --- | --- |
| Criterion: "Every rebuild remedy for a pool worktree names bench worktree build from one owner." | RB01, RB03, RB08, RB09, RB12, RB13, RB15 |
| Row: "The rebuild owner cannot tell a pool root from the primary checkout, so the fix needs a root-kind seam at that owner." | RB01, RB02, RB04, RB05, RB06, RB26 |
| Row: "Then each caller reads the remedy form from the owner." | RB12, RB13, RB14, RB15, RB16, RB23, RB24 |
| Row: "The live-binary check, the command entry, and the landing still print the `cd` form." | RB12, RB13, RB25 |
| Row: "The binary-seal row of the build preflight and the tree-target build refusal name `bench worktree build` through `usage.WorktreeBuildFor`." | RB15, RB23 |
| D1 (2026-10-31): the shell copy keeps the `cd` form | RB10, RB17, RB18, RB19, RB21, RB22 |

The landing clause of the row does not hold in the current tree: the landing destination is never a pool root (`internal/worktree/land_identity.go:51-53`). The spec records this premise and keeps the landing unchanged, with RB25 as its witness.

### Flagged additions

- RB15, RB16, RB23, and RB24 go beyond the named callers. The one-source rule of `AGENTS.md` requires one root-kind decision, so the `binary-seal` row reads the owner. The row's `Next` cell then names the path, not the assignment id. This is the contestable call of this spec.
- RB19 and RB20 go beyond the source sentences. A kit under test that lies in the pool makes the test-process home and the child home disagree, so the system pins must name the home they grade.
- RB07 comes from the hostile-input checklist for the new fill.
- RB11, RB22, and RB25 are kept-rule witnesses, not new promises.
- The tree-target build refusal keeps its label target, because it makes no root-kind decision (`internal/worktree/tree_target.go:90-93`).

### Reader sweep

- Production readers of `RebuildAction`: `cmd/bench/commands.go:33`, `internal/worktree/live_binary.go:28`, `internal/worktree/land.go:393`, and `internal/freshness/freshness_verify.go:107`. Readers of `Refusal.Error`: `internal/preflight/binary_seal.go:63`, `internal/adopt/doctor_rows.go:257-258`, `cmd/bench/freshness_publish.go:26-27`, `internal/freshness/check/main.go:15-17`, `internal/runbinary/runbinary.go:178-183` with `internal/gate/prospective.go:35`, `internal/gate/lane.go:316`, and `internal/testreport/command.go:50`, and `internal/worktree/land_effects.go:169-170`. Each one follows the owner with no edit except the live-binary check and the `binary-seal` row.
- The one production reader of `Refusal.Reason` is `internal/preflight/binary_seal.go:65`.
- Test pins that derive from `RebuildAction`: `internal/freshness/freshness_verify_test.go:206,377`, `cmd/bench/commands_brief_test.go:53`, `cmd/bench/freshness_publish_test.go:32`, `internal/worktree/live_binary_test.go:57`, `internal/worktree/land_broker_notice_test.go:57,131`, `internal/runbinary/runbinary_test.go:169`, `internal/gate/prospective_owner_test.go:109-112`, `internal/adopt/broker_test.go:169`, `internal/preflight/binary_seal_test.go:262`, `internal/systemtest/session_start_route_test.go:31`, `internal/systemtest/bench_follow_on_test.go:284`, and `internal/systemtest/owner_stale_seal_test.go:31`. The roots of the unit pins are temporary directories outside any pool, so they keep their answer. The system pins can see a kit inside the pool, so RB19 and RB20 move them.
- Literal pins of the `cd` form: `internal/freshness/freshness_verify_test.go:203,258,381` count or slice `rebuild with `, which `Refusal.Error` keeps. `internal/systemtest/testdata/guard-projection-family-d.golden:109` holds the shell copy's line for a kit outside the pool, which D1 keeps.
- `internal/conformance/build_contracts_test.go:29` reads the `printf` line of `bench_rebuild_action`, which D1 keeps.
- `internal/conformance/refusal_route_bypass_test.go:32` scans `internal/worktree` literals for `next=`, `next[`, `run bench `, and `"; then "`. The live-binary edit adds no literal.
- Prose readers: `projects/benchkit.md:696-702` already names `bench worktree build <target>` inside an assignment worktree, `docs/known-issues.md:10` names the same form, and `README.md:247-260` names a kit-clone install. None of them changes.
- `specs/shell-argument-ownership/spec.md:580` cites `TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths`, which keeps its name.
- `.github/workflows/release.yml:62` runs `scripts/go-build.sh` and prints no remedy. `bin/bench.sh` prints no rebuild remedy.
- Structure: `bench structure` is the source of each budget below. `internal/freshness` sits at its 12-file directory maximum, so the owner stays in `freshness_verify.go`, and its tests go to `internal/benchguard`, which sits well under its maximum. `internal/preflight/gather.go` sits within 5 lines of the 400-line file maximum, so it grows by at most 3 lines. `internal/preflight/decision.go` sits within 16 lines of that maximum, so it grows by at most 5 lines.

### Pre-review proof checklist

- Cited symbols: `RebuildAction` (`internal/freshness/freshness_verify.go:120`), `Refusal` (`internal/freshness/freshness_verify.go:96`), `Refusal.Reason` (`internal/freshness/freshness_verify.go:102`), `Refusal.Error` (`internal/freshness/freshness_verify.go:106`), `PublishedExecutable` (`internal/freshness/freshness_verify.go:115`), `shellQuote` (`internal/freshness/freshness_verify.go:124`), `benchguard.PoolReference` (`internal/benchguard/pool.go:18`), `benchguard.Classify` (`internal/benchguard/benchguard.go:108`), `poolkey.Pools` (`internal/poolkey/poolkey.go:44`), `benchhome.Dir` (`internal/benchhome`), `usage.WorktreeBuildFor` (`internal/usage/worktree.go:78`), `warnBeforeRemovingLiveBinary` (`internal/worktree/live_binary.go:22`), `discardIgnored` (`internal/worktree/clean.go:294`), `ambient.home` (`internal/worktree/effects.go:39`), `brokerChangeNotice` (`internal/worktree/land.go:387`), `landingDestinationIdentity` (`internal/worktree/land_identity.go:46`), `selectAssignment` (`internal/worktree/path.go:170`), `targetPath` (`internal/worktree/path.go:300`), `binarySealFacts` (`internal/preflight/binary_seal.go:26`), `binarySealCheck` (`internal/preflight/binary_seal.go:55`), `assignmentTarget` (`internal/preflight/gather.go:378`), `commandsCommand` (`cmd/bench/commands.go:24`), and `bench_rebuild_action` (`.bench/lib/resolve-bench.sh:379`).
- Import edges: `internal/freshness` gains direct edges to `internal/benchguard`, `internal/poolkey`, `internal/usage`, and `internal/benchhome`. `go list -deps` of `./internal/benchguard`, `./internal/poolkey`, `./internal/usage`, and `./internal/benchhome` names no `internal/freshness` package, so no edge makes a cycle. `internal/preflight` gains a direct edge to `internal/benchhome`. The external test package `benchguard_test` imports `internal/freshness`.
- Source-row clauses and occurrences: the table above quotes each clause. The row records one occurrence, the 2026-09-25 binary-seal-remedy light path (`roadmap/FT344.md:14`).
- Promised field labels: none. No intent-ledger or review-record field changes, so no `Persisted fields:` line applies.
- Changed-function callers: the reader sweep lists every caller of `RebuildAction`, `Refusal.Error`, `Refusal.Reason`, `warnBeforeRemovingLiveBinary`, `binarySealFacts`, and `binarySealCheck`. `binarySealFacts` has one production caller (`internal/preflight/gather.go:104`) and three test callers (`internal/preflight/binary_seal_test.go:187,215,242`).
- Copy survival: RB26 reviews the owner for a second pool rule. RB23 reds a `binary-seal` row that keeps the ownership branch. RB13 and RB14 red a live-binary warning that composes its remedy from another home.
- Rendered-shape readers: the old-text needle is `bash scripts/go-build.sh`. Its readers are in the reader sweep, and each keeps its root outside the pool or moves under RB19 and RB20.
- Pin operators: the owner rows compare the whole `Action` string. RB12 and RB13 use substring containment plus the absence of `cd `, as the current pins do (`cmd/bench/commands_brief_test.go:53`, `internal/worktree/live_binary_test.go:59-60`).
- Entry reads: `RebuildAction` reads `BENCH_HOME` through `benchhome.Dir`, and RB08 grades that read. `RebuildIn` is its internal form with an injected home. The preflight gather reads the home at its boundary and passes it to `binarySealFacts`, and RB24 drives that form.
- Derived expectations: each new-form expectation is a literal built from the fixture path. RB10, RB17 to RB19, and RB25 read `CheckoutRebuildAction` on purpose, because they pin the kept `cd` form and not the switch.
- Consolidated rules: the `binary-seal` row's ownership rule and the owner's pool rule become one rule. Old rule: an active assignment owns the root. New rule: the pool guard refuses the `cd` form of the root. RB15, RB16, and RB23 map the change.
- Quantified obligations: "every rebuild remedy for a pool worktree" quantifies over the Go printers in the reader sweep. The owner change reaches each one, and RB09 grades the shared `Refusal.Error` path. The shell copy is the reviewed exclusion D1.
- Workflow-step writes: none.
