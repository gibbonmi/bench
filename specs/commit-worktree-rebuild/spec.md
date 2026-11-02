# A commit in a kit worktree leaves no stale worktree build

Status: staged

Roadmap: FT406

Decision source: named reviewed artifact — `roadmap/FT406.md` at its pinned identity `sha256:aed754faf293126519efe9467ae9d00323fe84b5238b6bf57e71747e34a45335`, read with the bound criterion `FT406-rebuild.fresh` in `.bench/commitment.json` and the reviewer decision of 2026-10-10 in that row, plus the reviewer decision D1 of 2026-11-02 in the `/bench-write-spec FT406` conversation (D1: the trigger is the seal check; when no `dist/bench` exists, the commit builds only when its diff names a build input).

Verification log: <n> iteration(s) to accept — <note>

## Problem

A kit worktree publishes its own Bench executable at `dist/bench`, with a seal at `dist/bench.seal`.
The seal records the digest of the build inputs (`internal/freshness/freshness.go:35-61`).
The build inputs are the files of `go list -deps ./cmd/bench` inside the root, `go.mod`, `go.sum`, and the entries of `scripts/go-build.inputs` (`internal/freshness/freshness_buildinputs.go:50-140`).

An author edits a Go file and commits it with `bench commit`.
The commit composes, runs the lane, and publishes, but it builds nothing (`internal/commit/commit.go:81-110`).
So the worktree's `dist/bench` keeps the previous source, and its seal no longer matches the tree.

Two readers grade that seal with the same predicate, `freshness.Verify` or `freshness.VerifyAt`:

- The `binary-seal` row of the build preflight reds (`internal/preflight/binary_seal.go:28-71`).
- A `--in` verb such as `bench test --in` or `bench anchors --in` takes the first sealed build that matches the tree, and it refuses with `worktree build does not match the tree` when no candidate matches (`internal/treetarget/build.go:20-51`).

The author or the orchestrator then runs `bench worktree build <target>` by hand.
`docs/known-issues.md:10` records that workaround, and FT406 is its retire owner.

The row also reports that a Markdown-only commit made `bench test --in` and `bench anchors --in` refuse.
That premise does not hold in the tree.
No Markdown file is a build input: the only embedded files are `.bench/consumer-payload.json`, `internal/releaseevidence/registry.json`, `internal/releaseevidence/requirements.json`, and `internal/adopt/prepush.sh` (research C, `go list -deps ./cmd/bench`).
The digest inputs did not change after 8ae24503f (2026-09-30).
On 2026-10-10 the `--in` route accepted only the worktree's own build; dffdef899 (2026-10-28) added the primary build and the running executable as candidates (research C, `git show dffdef899 -- internal/treetarget/build.go`).
The build-directions tickets of that day also edited Go files, for example 823bca727 and 6303bb6f1 in `internal/anchors/`.
So the build key already ignores a Markdown-only change, as the row's decision states.

## Solution

After a `bench commit` publishes in a kit worktree, the commit makes the worktree's `dist/bench` match the tree.
The rule is D1:

- When the worktree declares no build inputs, the commit does nothing more. A linked repository is such a worktree.
- When `dist/bench` exists, the commit rebuilds it only if `freshness.Verify` refuses it.
- When `dist/bench` does not exist, the commit builds it only if the commit's diff names a build input.

A deleted build input counts as a named build input, because it was an input of the parent commit.
The commit then reads `freshness.Verify` a second time, and that read decides the outcome.
The step runs before the background grade starts, so the grade's wrapper finds the fresh build.
Each commit in a kit worktree now runs at least one `go list -deps`, through `Verify` or through the diff test, even when it builds nothing.
The step prints on stderr only, and it never changes the exit of the commit.
A failed step prints the cause and the rebuild remedy of the FT344 owner.

Two existing landing rules become shared owners in `internal/freshness`, so each fact keeps one source:

- `freshness.DiffChangesBuildInputs` decides whether a diff names a build input, a deleted input included. It replaces `brokerSourceChanged` (`internal/worktree/land.go:350-375`).
- `freshness.Refresh` runs the seal check, the build, and the second seal read. The landing refresh (`internal/worktree/land_effects.go:148-174`) and the commit step both call it.

The known-issues row retires, and an anchor forbids its return.

## User stories

Line: opus / medium.
Implementation-line reason: the hardest material chunk is CW-C2, which adds a post-publication step to a core write verb and must keep the stdout contract, the exit contract, and the order before the background grade. The spec fixes the predicate, each output line, and each seam, and each seam has a precedent: `refreshBroker` for the build and the second read, and `brokerSourceChanged` for the diff test. The rows drive the step with an injected build over a kit fixture, so a wrong predicate turns the ordinary Go test phase red at low cost. CW-C1 moves two existing rules into owners with no behavior change.
Harder chunks: CW-C2.

### The commit leaves a fresh worktree build

1. As a ticket author, I want a lane-passed `bench commit` in a kit worktree to rebuild a stale `dist/bench`, so that the next build preflight shows `binary-seal` green.
2. As an orchestrator, I want the rebuilt `dist/bench` to pass `freshness.Verify`, so that the next `--in` verb takes the worktree's own build.
3. As a ticket author in a worktree with no `dist/bench`, I want a commit whose diff names a build input to build one, so that a `--in` verb finds a build that matches my Go change.
4. As a ticket author, I want the rebuild to run after the commit publishes, so that the build seals the bytes that the commit formatted and published.
5. As a ticket author, I want the rebuild to finish before the background grade starts, so that the grade's wrapper resolves the fresh build.
6. As a ticket author whose `dist/bench` already verifies, I want no rebuild, so that a commit pays no build that it does not need.
7. As a spec author in a worktree with no `dist/bench`, I want a Markdown-only commit to build nothing, so that a Markdown commit costs no Go build.
8. As a ticket author whose `dist/bench` went stale before a Markdown-only commit, I want that commit to rebuild it, so that the commit leaves no stale build, as D1 decides.
9. As a ticket author who deletes a Go file, I want the commit to rebuild a present `dist/bench`, so that a deletion leaves no stale build.
10. As an orchestrator, I want a second seal read to decide the outcome, so that a build that exits 0 with no valid seal reports a failure.

### The commit keeps its output and its exit

11. As a commit-chain caller, I want `committed{published_commit=…,path_count=…}` to stay the last stdout line, so that each reader of that line keeps its answer.
12. As a ticket author, I want one stderr line after a rebuild, so that I see that the commit rebuilt the worktree build.
13. As a ticket author, I want one stderr line with the cause and the rebuild remedy after a failed rebuild, so that I can run the remedy by hand.
14. As a commit-chain caller, I want a failed rebuild to keep the commit's exit 0, so that exit 0 still means the commit published.
15. As a ticket author without Go on PATH, I want the failure line to name the absent Go, so that I know which tool to install.
16. As a ticket author, I want a build error with a newline or a control byte to print as one line, so that the cause cannot split or forge an output line.
17. As a ticket author, I want a commit that skips the rebuild to print no rebuild line, so that a skip adds no noise.
18. As a ticket author in a worktree whose path holds a space, I want the failure remedy to keep the root as one shell word, so that the remedy runs as printed.

### The commit rebuilds only where it may

19. As an operator of a linked repository, I want a commit to run no rebuild, so that a repository with no build inputs keeps its current commit.
20. As a ticket author, I want a refused commit to run no rebuild, so that a commit that published nothing builds nothing.
21. As a ticket author, I want a dry run to run no rebuild, so that a dry run writes nothing.
22. As a ticket author, I want a commit that published but did not reconcile its checkout to run no rebuild, so that no build seals an unreconciled tree.
23. As a ticket author, I want the rebuild outputs to stay out of the tracked tree, so that a rebuild leaves `git status` clean.
24. As an orchestrator of a declared build, I want the commit step to read no held-install predicate, so that a worktree whose manifest binds its own build still rebuilds.

### One owner for each rule

25. As a kit maintainer, I want one predicate that decides whether a diff names a build input, so that the landing and the commit cannot disagree.
26. As a kit maintainer, I want one owner of the seal check, the build, and the second seal read, so that the landing refresh and the commit step obey one rule.
27. As an orchestrator, I want the landing's broker notice and refresh to keep their current output, so that the landing behaves as before.

### Kept behavior

28. As a commit-chain caller, I want `--preflight-build` to keep its three steps, so that a plan commit keeps its build and its preflight.
29. As a ticket author, I want a publishing commit to keep starting its background grade, so that the grade still runs.
30. As an orchestrator, I want the `binary-seal` row and the `--in` route to keep their code, so that they grade the rebuilt build with the same predicate.

### Records

31. As a coordinator who writes a delegate charge, I want the stale-build row gone from `docs/known-issues.md`, so that no charge names a workaround that the commit now makes needless.
32. As a reader of the Command Notes, I want one sentence that states the commit rebuild, so that a session in any harness knows when a worktree build is fresh.
33. As a kit maintainer, I want the profile sentence about the worktree's own build to name the commit rebuild, so that the profile states each source of that build.
34. As a reader of the release notes, I want one `CHANGELOG.md` entry, so that the change is visible after the row retires.

### Reviewed exclusions

35. As a reviewer, I want an absent `dist/bench` to build only when the diff names a build input, so that D1 holds.
36. As a reviewer, I want no change to the build key, so that the decision "the build key ignores a Markdown-only change" holds as the tree already does.

### Deleted inputs and declared builds

37. As a ticket author in a worktree with no `dist/bench`, I want a commit that deletes a Go build input to build one, so that the next `--in` verb finds a build that matches the tree.
38. As an orchestrator, I want a landing that only deletes a Go build input to report the install step of a broker-changing landing, so that the shared diff test gives the landing the same answer as the commit.
39. As an orchestrator of a declared build, I want a commit in a worktree whose own broker manifest binds its stale build to rebuild that build, so that an open declared build does not stop the rebuild.

## Implementation decisions

### The shared owners

`internal/freshness` sits above its directory maximum, so both owners go into existing files.

`DiffChangesBuildInputs(root, base, tip string) bool` goes into `freshness_buildinputs.go`.
Its body is the current body of `brokerSourceChanged`: a root that does not declare build inputs answers false, an input set that does not resolve answers false, and a diff that does not read answers false.
Otherwise it reads the diff once, as `git diff --no-renames --name-status base tip`, and it answers true when that diff names a path of `BuildInputs(root)`. The status `D` marks a deleted path for the second test below.
The `--no-renames` flag makes a rename print its old path as a deletion and its new path as an addition. Porcelain `git diff` detects renames by default and then prints only the new path. The kit already reads diffs this way (`internal/diff/snapshot.go:55`).
It also answers true when the diff deletes a `.go` file that is not a `_test.go` file, and the directory of that file is the directory of a package that `go list -deps ./cmd/bench` lists for `root`.
The post-commit input set cannot hold a deleted path, so this second test covers the deleted input that D1 names, and through `--no-renames` it covers the old path of a renamed input.
A deletion that empties a whole package also changes an importer of that package, and a deleted file of a single-file `//go:embed` pattern also changes its directive, so the first test covers those two cases. Each embed in the build closure today names a single file (research C). A deleted file under a directory or glob embed pattern is a Won't handle.
The rule reads the package directories from the same `go list` call that `buildInputs` makes, so the input set keeps one source (`internal/freshness/freshness_buildinputs.go:55-76`).
The rule also counts a deleted file that a build constraint excluded, such as a `_darwin.go` file. In the commit step that costs one extra build. In the landing it gives an install-step notice for a diff that does not change the broker.
The landing's broker notice reads the same owner, so a landing that only deletes a Go build input now reports the install step. That change is a flagged addition.
It reads the diff through `internal/git`. `internal/freshness` already depends on `internal/git` through `internal/poolkey`, and `go list -deps ./internal/git` names no `freshness` package, so the edge makes no cycle.
`brokerSourceChanged` goes away. Its two callers, `internal/worktree/land.go:193` and `internal/worktree/land_stack.go:85`, call the owner.

`Refresh(ctx context.Context, root string, build func(context.Context, string, string) error) (rebuilt bool, buildErr, err error)` goes into `freshness_verify.go`.
When `Verify(root, PublishedExecutable(root))` passes, it returns `false, nil, nil` and does not call `build`.
Otherwise it calls `build(ctx, root, PublishedExecutable(root))` once and keeps that error as `buildErr`.
It then reads `Verify` a second time, returns that error as `err`, and returns `rebuilt` true only when `err` is nil.
`refreshBroker` keeps its not-applicable branch, its held-install branch, its two stderr texts, and its effect results. It calls `Refresh` in place of its own three steps, and it prints `landing refresh build failed: ` before a non-nil `buildErr` and `landing refresh failed: ` before a non-nil `err`, as now.

### The commit step

`internal/commit/background.go` gets the step, because `internal/commit/` sits above its directory maximum and `background.go` holds the other post-publication step.

- A package variable `worktreeBuild` holds `runbinary.BuildSubject` (`internal/runbinary/runbinary.go:279-281`). The step does not call `worktree.BuildCommand`, which prints a table and an exec hint and resolves an assignment (`internal/worktree/build.go:36-95`).
- `refreshWorktreeBuild(root, published, home string, build func(context.Context, string, string) error, stderr io.Writer)` is the internal form. It reads no ambient value.
- `Run` calls `refreshWorktreeBuild(root, outcome.Published, benchhome.Dir(), worktreeBuild, stderr)` only when `exit` is 0 and `outcome.Published` is not empty, and it calls it before `startBackgroundGrade`.

`refreshWorktreeBuild` decides in this order:

1. When `freshness.DeclaresBuildInputs(root)` is false, it returns and prints nothing.
2. When `os.Lstat(freshness.PublishedExecutable(root))` reports `fs.ErrNotExist` and `freshness.DiffChangesBuildInputs(root, published+"^", published)` is false, it returns and prints nothing. Any other `Lstat` answer counts as present, as `binarySealFacts` counts it (`internal/preflight/binary_seal.go:30-34`).
3. It calls `freshness.Refresh(context.Background(), root, build)`.
4. When `Refresh` passes with no build, it prints nothing.
5. When `Refresh` reports `rebuilt`, it prints `worktree build: rebuilt dist/bench` on stderr. The second seal read decides, so this line prints even when `buildErr` is not nil.
6. When `Refresh` reports a non-nil `err`, it prints one stderr line: `worktree build: failed: <cause>; rebuild with <remedy>`. The cause is `buildErr.Error()` when `buildErr` is not nil. Otherwise the cause is `Refusal.Reason()` when `errors.As` finds a `*freshness.Refusal` in `err`, and `err.Error()` for any other error. `Refusal.Error` already ends with `; rebuild with ` and the remedy (`internal/freshness/freshness_verify.go:111-113`), so the step never prints it as the cause. The step passes the cause through `sanitize.Controls`. The remedy is `freshness.RebuildIn(home, root).Action` (`internal/freshness/freshness_verify.go:142-148`), so the line holds the remedy exactly once.

The step never changes `exit`.
The step passes `context.Background()` to `Refresh`, as `refreshBroker` passes it to its build (`internal/worktree/land_effects.go:161`).
No new literal holds `next=`, `next[`, `run bench `, or `"; then "`, so the route-bypass scan stays green (`internal/conformance/refusal_route_bypass_test.go:32`; the scanned packages are at `internal/conformance/registry/route_bypass.go:16-27`, and they include `internal/commit`).
The step reads no held-install predicate. During a declared build the worktree's own broker manifest binds the worktree build, so `readHeldInstall` would hold in every built worktree (`internal/worktree/land_effects.go:176-200`). `Refresh` holds no held-install branch either; `refreshBroker` keeps that branch before its call to `Refresh`.

`Digest` reads the working tree (`internal/freshness/freshness.go:35-61`).
The commit's gofmt step rewrites the named Go files before publication (`internal/commit/commit.go:159-167`), so the rebuild reads the formatted bytes.
An unnamed dirty Go file is sealed too; the `binary-seal` row and the `--in` route grade the same working tree.

### Records

- `docs/known-issues.md` loses the stale-build row at line 10. `docs/known-issues.md:12` stays, because its Go edit is not committed.
- `internal/anchors/registry_ticket_passes.go` gets one `Forbid` row on the known-issues file with the needle `FT406`, as the `FT405` and `FT438` rows do (`internal/anchors/registry_ticket_passes.go:15-24`). `internal/anchors/registry_ticket_passes_test.go` gets the matching independent rule with the retired row's bytes.
- `.bench/BENCH-reference.md` gets one sentence after the worktree-commit paragraph of the Command Notes: in a kit worktree, a commit that publishes rebuilds the worktree's `dist/bench` when its seal does not verify, or when no build exists and the commit changes a build input.
- `projects/benchkit.md:713` keeps each predicate of its sentence and adds the commit rebuild as a second source of the worktree's own build.
- `CHANGELOG.md` gets one entry under `## [Unreleased]`.

## Implementation chunks

CW-C1 moves the two landing rules into shared owners. The commit step consumes both owners, so CW-C1 is its own small chunk, and its review closes before CW-C2 starts.
CW-C2 adds the commit step, its output, and the records.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| CW-C1 / 01-shared-freshness-owners.md | `freshness.DiffChangesBuildInputs` and `freshness.Refresh` own the diff test and the refresh rule, the diff test counts a deleted input, and the landing reads both owners with its refresh unchanged | CW01, CW02, CW03, CW04, CW05, CW06, CW07, CW08, CW09, CW10, CW11, CW12, CW41, CW42, CW44, CW46, CW47, CW49 | TestDiffChangesBuildInputs, TestRefreshBuildsOnlyAStaleExecutable, TestRefreshSecondReadDecides, TestLandCommandReportsInstallStepForABrokerChangingDiff, TestLandReportsInstallStepForADeletedBrokerInput, TestLandRefreshesTheBrokerAfterPublication, TestLandReportsAFailedRefresh, TestLandSkipsAFreshBroker, TestRefreshHoldsAnInstalledCandidate | no |
| CW-C2 / 02-commit-rebuilds-worktree-build.md | A lane-passed commit in a kit worktree leaves a fresh `dist/bench` before its background grade, with stderr-only output, an unchanged exit, and the stale-build known issue retired | CW13, CW14, CW15, CW16, CW17, CW18, CW19, CW20, CW21, CW22, CW23, CW24, CW25, CW26, CW27, CW28, CW29, CW30, CW31, CW32, CW33, CW34, CW35, CW36, CW37, CW38, CW39, CW40, CW43, CW45, CW48, CW50 | TestWorktreeRebuildRepairsAStaleBuild, TestWorktreeRebuildBuildsAnAbsentBuildForAnInputChange, TestWorktreeRebuildBuildsAnAbsentBuildForADeletedInput, TestWorktreeRebuildSkipsAFreshBuild, TestWorktreeRebuildSkipsAMarkdownCommitWithoutABuild, TestWorktreeRebuildRepairsAStaleBuildAfterAMarkdownCommit, TestWorktreeRebuildRepairsADeletedInput, TestWorktreeRebuildReportsAFailure, TestWorktreeRebuildNamesAnAbsentGo, TestWorktreeRebuildSkipsALinkedRepository, TestWorktreeRebuildIgnoresAHeldInstall, TestCommitRebuildsBeforeItsGrade, TestCommitKeepsExitZeroOnAFailedRebuild, TestUnpublishedCommitRebuildsNothing, TestUnreconciledCommitRebuildsNothing, TestKitIgnoresTheRebuildOutputs, TestCommitWithoutPreflightBuildRunsCommitAlone, TestCommitChainRunsThreeSteps, TestCommitStartsABackgroundGrade, TestBinarySealFactsGradeARootOnDisk, TestRunKitWorktreeBuild, TestTicketSlicingAnchors | yes |

### Completion plan

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "CW-C1",
      "tickets": ["01-shared-freshness-owners.md"],
      "verification": [
        {"id": "owners", "command": "bench test --package ./internal/freshness --run 'DiffChangesBuildInputs|Refresh'", "probe": "Make Refresh return rebuilt true after a nil build error without the second Verify read, and record that TestRefreshSecondReadDecides turns red."},
        {"id": "deleted-input", "command": "bench test --package ./internal/freshness --run 'DiffChangesBuildInputs'", "probe": "Remove the deleted-file test from DiffChangesBuildInputs so that it reads only the post-commit input set, and record that TestDiffChangesBuildInputs turns red on the deleted-input case."},
        {"id": "renamed-input", "command": "bench test --package ./internal/freshness --run 'DiffChangesBuildInputs'", "probe": "Remove --no-renames from the diff that DiffChangesBuildInputs reads, and record that TestDiffChangesBuildInputs turns red on the renamed-input case."},
        {"id": "landing", "command": "bench test --package ./internal/worktree --run 'LandCommandReportsInstallStepForABrokerChangingDiff|LandReportsInstallStepForADeletedBrokerInput|LandRefreshesTheBrokerAfterPublication|LandReportsAFailedRefresh|LandSkipsAFreshBroker|RefreshHoldsAnInstalledCandidate'", "probe": "Make DiffChangesBuildInputs return false for every diff, and record that TestLandCommandReportsInstallStepForABrokerChangingDiff turns red."}
      ]
    },
    {
      "id": "CW-C2",
      "tickets": ["02-commit-rebuilds-worktree-build.md"],
      "verification": [
        {"id": "commit-step", "command": "bench test --package ./internal/commit --run 'WorktreeRebuild|CommitRebuildsBeforeItsGrade|CommitKeepsExitZeroOnAFailedRebuild|UnpublishedCommitRebuildsNothing|UnreconciledCommitRebuildsNothing|KitIgnoresTheRebuildOutputs|CommitStartsABackgroundGrade|GreenLandingPublishesTheNamedPath'", "probe": "Move the refreshWorktreeBuild call in Run after startBackgroundGrade, and record that TestCommitRebuildsBeforeItsGrade turns red."},
        {"id": "failure-line", "command": "bench test --package ./internal/commit --run 'WorktreeRebuildReportsAFailure'", "probe": "Print err.Error() as the cause for every error, and record that TestWorktreeRebuildReportsAFailure turns red on the seal-refusal case."},
        {"id": "absent-predicate", "command": "bench test --package ./internal/commit --run 'WorktreeRebuild'", "probe": "Remove the DiffChangesBuildInputs test for an absent executable so that the step always calls Refresh, and record that TestWorktreeRebuildSkipsAMarkdownCommitWithoutABuild turns red."},
        {"id": "chain", "command": "bench test --package ./cmd/bench --run 'CommitChain|CommitWithoutPreflightBuild'", "probe": "Make commitChainCommand call steps.Build for a plain commit, and record that TestCommitWithoutPreflightBuildRunsCommitAlone turns red."},
        {"id": "known-issues", "command": "bench test --package ./internal/anchors --run 'TicketSlicingAnchors'", "probe": "Insert the retired FT406 row into docs/known-issues.md, and record that TestTicketSlicingAnchors turns red."}
      ]
    }
  ],
  "final_verification": [
    {"id": "owner-packages", "command": "bench test --package ./internal/freshness/... --package ./internal/commit --package ./internal/worktree --package ./internal/anchors --package ./internal/preflight/... --package ./internal/treetarget/... --package ./cmd/bench"},
    {"id": "system", "command": "bench test --check system"},
    {"id": "coverage", "command": "bench coverage --check commit-worktree-rebuild"}
  ]
}
```

## Testing decisions

- A good owner test calls `freshness.DiffChangesBuildInputs` and `freshness.Refresh` over a fixture root, and it grades the answer and the count of build calls. The owner tests go into `internal/freshness/decision_test.go` (`package freshness`, 31 lines). They build their roots with the package's own fixture helpers `writeBuildFixture` (`internal/freshness/freshness_digest_test.go:282`) and `writePublishedFixture` (`internal/freshness/freshness_verify_test.go:338`). `kittest` imports `internal/freshness`, so an internal test of that package cannot import it.
- The fake build of an owner test publishes with one `freshness.Publish` call, as `internal/treetarget/build_test.go:52` does. A second fake returns nil and publishes nothing, and a third fake returns an error.
- The landing rows reuse the landing tests as they stand. They inject `buildSubject` through `refreshJoins` (`internal/worktree/land_effects_test.go:90`).
- The commit-step rows call `refreshWorktreeBuild` with an explicit home and an injected build. Each fixture is a Git repository whose tree comes from `kittest.WriteTree` (`internal/treetarget/kittest/kittest.go:26`), with a first commit and a published build. `kittest.EditBuildInput` makes a published build stale (`internal/treetarget/kittest/kittest.go:50`).
- The commit-level rows call `Run` through `runCommand`, which runs in the test process (`internal/commit/landing_test.go:21`). They swap the package variable `worktreeBuild` in that process, so the swap reaches `Run`. Those tests already set `BENCH_HOME` and `BENCH_KIT`, so they do not run in parallel. The launch fixture of `internal/commit/background_test.go` gets the kit files of `kittest.WriteTree` in the worktree, so `DeclaresBuildInputs` answers true there.
- The order row reads the stand-in grade record: the fake build sees `f.launches.Launched(t)` empty when it runs.
- The tests go into `internal/commit/background_test.go` (206 lines). `landing_test.go` (374 lines) and `lane_test.go` (390 lines) sit near the 400-line maximum, and `internal/commit/` takes no new file.
- Each expected line is a literal in the test, except the remedy, which the test reads from the FT344 owner `freshness.RebuildIn(home, root).Action`. That owner is not under test here.
- The ordinary Go test phase observes CW01 to CW09, CW11 to CW38, and CW41 to CW50. CW10, CW39, and CW40 are review-owned.
- CW45 needs an open declared build in the fixture repository's intent ledger. The commit test builds it from exported helpers, as `openInstallBuild` does in `internal/worktree` (`internal/worktree/land_effects_test.go:233-243`): it registers and admits the assignment as `landingRepo` does (`internal/commit/commitment_test.go:21-31`), commits a spec whose header holds a `Persisted fields:` line on `main` in the primary repository, because `commitrepo.SourceIdentity` and the declaration read take the default branch (`internal/commitment/repository/window.go:49-67`), reads its identity with `commitrepo.SourceIdentity`, and rebinds the binding that `Admit` wrote to that deliverable and identity with `commitmenttest.EditCommitment`, with a claim on its outcome, as `openInstallBuild` does (`internal/worktree/land_effects_test.go:236-244`). The fixture then writes the worktree's broker manifest over the stale published build with `brokermanifest.Write`, as `internal/worktree/land_effects_test.go:296` does. Before it calls the step, the test asserts its precondition: `commitrepo.Store{Root: root}.DeclaredBuild()` answers declared, and the manifest binds the published executable. Without that assert, CW45 passes over a fixture that holds no open declared build.
- The CW49 fixture renames the file with its content unchanged, so that Git's rename detection fires without `--no-renames`.
- CW43 and CW21 need a kit fixture whose `go list -deps ./cmd/bench` lists a second Go file in the `cmd/bench` package, so the deletion changes the digest. The ticket adds that file to the fixture of the test, not to `kittest.WriteTree`. The system suite observes no new row, and `bench test --check system` stays a final verification because the commit step runs in every kit worktree commit.

### Seam diagram

    trigger: a lane-passed bench commit in a kit worktree (exit 0, published commit)
        │
        ▼
    root, published, home  ──▶  [ refreshWorktreeBuild ]  ──▶  stderr line or none; exit unchanged
                                   │ DeclaresBuildInputs(root)?
                                   │ dist/bench absent → DiffChangesBuildInputs(root, published^, published)?
                                   ▼
                              [ freshness.Refresh: Verify → build → Verify ]  ──▶  rebuilt, buildErr, err
                                   ◀ tests attach here: background_test.go calls refreshWorktreeBuild
                                     with a kittest repository and a fake build that publishes or fails
        │
        └──▶ startBackgroundGrade (after the step)   ◀ background_test.go: the fake build sees no launch

    trigger: a landing (broker notice and refresh effect)
        │
        ▼
    worktree, base, tip ──▶ [ freshness.DiffChangesBuildInputs ] ──▶ broker notice
    root ──▶ [ refreshBroker: held-install branch, then freshness.Refresh ] ──▶ effect result
              ◀ land_broker_notice_test.go and land_effects_test.go, unchanged

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| CW01 | 25 | `DiffChangesBuildInputs` returns true for a diff that names `scripts/go-build.sh` in a fixture root that declares build inputs | planned TestDiffChangesBuildInputs in internal/freshness/decision_test.go | An owner that drops the intersection never reports a broker change |
| CW02 | 25, 36 | `DiffChangesBuildInputs` returns false for a diff that names only `README.md` | planned TestDiffChangesBuildInputs in internal/freshness/decision_test.go | An owner that reports any changed path builds for a Markdown commit |
| CW03 | 19, 25 | `DiffChangesBuildInputs` returns false for a root with no `scripts/go-build.inputs` | planned TestDiffChangesBuildInputs in internal/freshness/decision_test.go | An owner without the manifest test runs `go list` in a linked repository and answers from an error |
| CW04 | 26 | `Refresh` over a root whose published build verifies returns `rebuilt` false and nil errors, with zero build calls | planned TestRefreshBuildsOnlyAStaleExecutable in internal/freshness/decision_test.go | An owner that skips the first read builds on every call |
| CW05 | 26 | `Refresh` over a root whose published build is stale calls the build once and returns `rebuilt` true with nil errors | planned TestRefreshBuildsOnlyAStaleExecutable in internal/freshness/decision_test.go | An owner that does not build leaves the stale seal in place |
| CW06 | 10, 26 | `Refresh` with a build that returns nil and publishes nothing returns `rebuilt` false and a non-nil `err` | planned TestRefreshSecondReadDecides in internal/freshness/decision_test.go | An owner that trusts the build's exit reports a rebuild that no seal authenticates |
| CW07 | 26 | `Refresh` with a build that returns an error returns that error as `buildErr` | planned TestRefreshSecondReadDecides in internal/freshness/decision_test.go | An owner that drops the build error leaves the caller only the seal reason |
| CW08 | 3, 26 | `Refresh` over a root with no published build calls the build once | planned TestRefreshBuildsOnlyAStaleExecutable in internal/freshness/decision_test.go | An owner that treats an absent build as fresh never builds the first build |
| CW09 | 27 | A landing of a diff that changes `scripts/go-build.sh` still reports the install step of a broker-changing landing | `internal/worktree/land_broker_notice_test.go` (`TestLandCommandReportsInstallStepForABrokerChangingDiff`) | A caller that lost the owner's answer drops the broker notice |
| CW10 | 25 | `internal/worktree` holds no `diff --name-only` intersection with `BuildInputs` after the move | review-owned: `rg -n 'name-only' internal/worktree/land.go internal/worktree/land_stack.go` | A surviving local copy is a second source for one fact |
| CW11 | 27 | A landing refresh over a stale destination build calls the build once and reports `complete` | `internal/worktree/land_effects_test.go` (`TestLandRefreshesTheBrokerAfterPublication`) | A `refreshBroker` that calls `Refresh` wrongly loses the effect result |
| CW12 | 27 | A landing refresh of an installed candidate stays `held` and calls no build | `internal/worktree/land_effects_test.go` (`TestRefreshHoldsAnInstalledCandidate`) | A `refreshBroker` that moves the held-install branch into `Refresh` rebuilds over a candidate |
| CW13 | 1 | `refreshWorktreeBuild` over a committed repository whose published build is stale calls the build once | planned TestWorktreeRebuildRepairsAStaleBuild in internal/commit/background_test.go | Today's commit calls no build |
| CW14 | 2 | After CW13, `freshness.Verify(root, freshness.PublishedExecutable(root))` returns nil | planned TestWorktreeRebuildRepairsAStaleBuild in internal/commit/background_test.go | A step that builds to another path leaves the graded build stale |
| CW15 | 12 | After CW13, stderr is exactly `worktree build: rebuilt dist/bench` and one newline | planned TestWorktreeRebuildRepairsAStaleBuild in internal/commit/background_test.go | A silent rebuild gives the author no sign, and a chatty one breaks the one-line contract |
| CW16 | 3 | `refreshWorktreeBuild` over a repository with no `dist/bench`, whose published commit changes `scripts/go-build.sh`, calls the build once | planned TestWorktreeRebuildBuildsAnAbsentBuildForAnInputChange in internal/commit/background_test.go | A step that skips every absent build leaves a Go commit with no matching build |
| CW17 | 6 | `refreshWorktreeBuild` over a repository whose published build verifies calls no build | planned TestWorktreeRebuildSkipsAFreshBuild in internal/commit/background_test.go | A step that always builds costs a Go build on every commit |
| CW18 | 17 | After CW17, stderr is empty | planned TestWorktreeRebuildSkipsAFreshBuild in internal/commit/background_test.go | A step that prints on a skip adds noise to every commit |
| CW19 | 7, 35 | `refreshWorktreeBuild` over a repository with no `dist/bench`, whose published commit changes only a Markdown file, calls no build and prints nothing | planned TestWorktreeRebuildSkipsAMarkdownCommitWithoutABuild in internal/commit/background_test.go | A step that treats an absent build as stale builds for each Markdown commit |
| CW20 | 8 | `refreshWorktreeBuild` over a repository whose published build went stale before a Markdown-only commit calls the build once | planned TestWorktreeRebuildRepairsAStaleBuildAfterAMarkdownCommit in internal/commit/background_test.go | A step keyed on the diff alone leaves the stale build that D1 repairs |
| CW21 | 9 | `refreshWorktreeBuild` over a repository whose published commit deletes a sealed `cmd/bench` Go file calls the build once | planned TestWorktreeRebuildRepairsADeletedInput in internal/commit/background_test.go | A step keyed on the post-commit input set misses a deleted input |
| CW22 | 13, 18 | A build that returns `boom` gives exactly one stderr line `worktree build: failed: boom; rebuild with ` followed by `freshness.RebuildIn(home, root).Action` for a root whose path holds a space | planned TestWorktreeRebuildReportsAFailure in internal/commit/background_test.go | A step that omits the cause or composes its own remedy gives a line that does not run |
| CW23 | 16 | A build error that holds a newline and an ESC byte gives exactly one stderr line with no ESC byte | planned TestWorktreeRebuildReportsAFailure in internal/commit/background_test.go | An unsanitized cause splits the line or carries a terminal escape |
| CW24 | 10, 13 | A build that returns nil and publishes nothing gives stderr exactly `worktree build: failed: ` plus `Refusal.Reason()` of the second seal read, then `; rebuild with ` plus `freshness.RebuildIn(home, root).Action`, and one newline | planned TestWorktreeRebuildReportsAFailure in internal/commit/background_test.go | A step that trusts the build's exit reports a rebuild that no seal authenticates, and a step that prints `Refusal.Error` as the cause prints the remedy twice |
| CW25 | 15 | With `runbinary.BuildSubject` as the build and a `PATH` with no `go`, the stderr line holds `Go is absent from PATH` | planned TestWorktreeRebuildNamesAnAbsentGo in internal/commit/background_test.go | A step that swallows the build error hides the missing tool |
| CW26 | 19 | `refreshWorktreeBuild` over a repository with no `scripts/go-build.inputs` calls no build and prints nothing | planned TestWorktreeRebuildSkipsALinkedRepository in internal/commit/background_test.go | A step without the manifest test runs a build script that a linked repository lacks |
| CW27 | 5 | In a kit launch fixture with a stale build and a fake build, the fake build sees zero grade launches when it runs | planned TestCommitRebuildsBeforeItsGrade in internal/commit/background_test.go | A step that runs after `startBackgroundGrade` lets the grade resolve the stale build |
| CW28 | 4 | In the CW27 fixture, the fake build runs once after `Run` publishes, and the published commit is `HEAD` when it runs | planned TestCommitRebuildsBeforeItsGrade in internal/commit/background_test.go | A step that runs before the publication seals a tree that the commit may refuse |
| CW29 | 11 | In the CW27 fixture, the last stdout line starts with `committed{published_commit=` and stdout holds no `worktree build` text | planned TestCommitRebuildsBeforeItsGrade in internal/commit/background_test.go | A step that prints on stdout breaks the last-line contract of the commit record |
| CW30 | 14 | In a kit launch fixture with a stale build and a fake build that fails, `Run` exits 0 | planned TestCommitKeepsExitZeroOnAFailedRebuild in internal/commit/background_test.go | A step that changes the exit turns a published commit into a refusal |
| CW31 | 20, 21 | A lane-refused commit and a dry run in a kit launch fixture with a stale build each call no build | planned TestUnpublishedCommitRebuildsNothing in internal/commit/background_test.go | A step keyed on the fixture instead of the publication builds for a commit that published nothing |
| CW32 | 22 | A commit that exits 3 in a kit fixture with a stale build calls no build | planned TestUnreconciledCommitRebuildsNothing in internal/commit/background_test.go | A step keyed on a published subject alone builds over an unreconciled checkout |
| CW33 | 23 | At the module root, `git check-ignore` reports `dist/bench`, `dist/bench.seal`, and `bin/bench-broker.manifest` as ignored | planned TestKitIgnoresTheRebuildOutputs in internal/commit/background_test.go | A kit that tracks a rebuild output leaves `git status` dirty after each rebuild |
| CW34 | 28 | A plain commit through the chain still calls no chain build step | `cmd/bench/commit_chain_test.go` (`TestCommitWithoutPreflightBuildRunsCommitAlone`) | A rebuild added as a chain step changes the chain contract |
| CW35 | 28 | A `--preflight-build` commit still prints the three chain steps | `cmd/bench/commit_chain_test.go` (`TestCommitChainRunsThreeSteps`) | A chain change drops the build or the preflight step |
| CW36 | 29 | A publishing commit in a kit worktree still starts one background grade | `internal/commit/background_test.go` (`TestCommitStartsABackgroundGrade`) | A step that returns early from `Run` skips the grade |
| CW37 | 30 | The `binary-seal` gather still grades a stale and a fresh root on disk | `internal/preflight/binary_seal_test.go` (`TestBinarySealFactsGradeARootOnDisk`) | An edit of the row's code changes the predicate that the rebuild must satisfy |
| CW38 | 31 | `docs/known-issues.md` holds no `FT406` text | `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingAnchors`) | A restored stale-build row tells each charge to run a needless manual build |
| CW39 | 32 | `.bench/BENCH-reference.md` states the commit rebuild in one sentence of the Command Notes | review-owned | Review grades the rule wording |
| CW40 | 33 | `projects/benchkit.md` names the commit rebuild as a source of the worktree's own build and keeps each predicate of the old sentence | review-owned | Review grades the reword rule |
| CW50 | 10, 12 | A build that returns `boom` after it publishes a valid seal gives stderr exactly `worktree build: rebuilt dist/bench` and one newline | planned TestWorktreeRebuildReportsAFailure in internal/commit/background_test.go | A step that lets `buildErr` win over the second seal read reports a failure for a build that the seal authenticates |
| CW41 | 25, 37 | `DiffChangesBuildInputs` returns true for a diff that only deletes a non-test `.go` file from the directory of a listed `cmd/bench` package | planned TestDiffChangesBuildInputs in internal/freshness/decision_test.go | An owner that reads only the post-commit input set never sees a deleted input |
| CW49 | 25, 37 | `DiffChangesBuildInputs` returns true for a diff that only renames a non-test `.go` file of a listed `cmd/bench` package to a `_test.go` name in the same directory | planned TestDiffChangesBuildInputs in internal/freshness/decision_test.go | An owner that reads the diff with rename detection sees only the new `_test.go` path and misses the removed input |
| CW42 | 25 | `DiffChangesBuildInputs` returns false for a diff that only deletes a `_test.go` file from that directory | planned TestDiffChangesBuildInputs in internal/freshness/decision_test.go | An owner that counts every deleted `.go` file builds for a test-only deletion that the digest never read |
| CW43 | 37 | `refreshWorktreeBuild` over a repository with no `dist/bench`, whose published commit only deletes a sealed `cmd/bench` Go file, calls the build once | planned TestWorktreeRebuildBuildsAnAbsentBuildForADeletedInput in internal/commit/background_test.go | A step whose diff test misses a deletion leaves the next `--in` verb with no build that matches the tree |
| CW44 | 38 | A landing whose reviewed diff only deletes a non-test Go file of a listed package reports the install step of a broker-changing landing | planned TestLandReportsInstallStepForADeletedBrokerInput in internal/worktree/land_broker_notice_test.go | A landing that keeps a local post-tip intersection drops the notice for a deleted broker input |
| CW45 | 24, 39 | `refreshWorktreeBuild` over a repository whose intent ledger holds an open declared build, and whose own broker manifest binds its stale published build, calls the build once | planned TestWorktreeRebuildIgnoresAHeldInstall in internal/commit/background_test.go | A step that copies the held-install branch of `refreshBroker` never rebuilds in a declared-build worktree |
| CW46 | 27 | A landing refresh whose build fails still reports `failed` with the text `landing refresh build failed: ` | `internal/worktree/land_effects_test.go` (`TestLandReportsAFailedRefresh`) | A `refreshBroker` that calls `Refresh` and drops `buildErr` loses the build failure text |
| CW47 | 27 | A landing refresh over a fresh destination build calls no build and reports `complete` | `internal/worktree/land_effects_test.go` (`TestLandSkipsAFreshBroker`) | A `refreshBroker` that skips the first seal read rebuilds a fresh broker |
| CW48 | 30 | A `--in` child in a kit worktree takes the worktree's own build when that build verifies | `internal/treetarget/build_test.go` (`TestRunKitWorktreeBuild`) | An edit of the `--in` selection changes the reader that the rebuilt build must satisfy |

Not covered: story 34 — the `CHANGELOG.md` entry is review-owned in the `## [Unreleased]` section and takes no behavior row.

### Edge inventory

These behaviors serve the kit only. A linked repository declares no build inputs, so CW03 and CW26 keep its commit unchanged.
CW19 and CW16 cover the absent build, and CW13 and CW17 cover the present build. CW22 covers a root path with a space; the FT344 owner quotes the root (`internal/freshness/freshness_verify.go:142-148`).
A dangling or non-regular `dist/bench` counts as present, so `Verify` refuses it and the step rebuilds it through the publication's backup path (`internal/freshness/freshness_publish.go:308-340`).
A merge-parent commit reads its first parent, `published^` (`internal/commit/merge_parent_test.go`).
The step writes only ignored paths, and CW33 pins that fact for the kit's own `.gitignore`.
CW25 swaps `PATH` with `t.Setenv`, and the commit-level rows swap `worktreeBuild` in the test process, so those tests do not run in parallel.

Won't handle: `bench commit --in <label>` after an uncommitted Go edit — the tree-target route refuses before `commit.Run` runs (`internal/treetarget/run.go:63-66`), and FT341 owns that route.
Won't handle: a reader in the window between the executable rename and the seal rename — the publication lock and rollback own that window (`internal/freshness/freshness_publish.go:68-94`), and FT398 owns the flake.
Won't handle: the `candidate-not-installed` refusal after a Go commit in a declared build — the guard compares the running executable with the primary manifest (`internal/commitment/repository/window.go:78-143`), and `.agents/commands/bench-implement-spec.md:46` already requires the install again after a Go change. A worktree that never had a build now meets the same guard after its first Go commit, because the wrapper then prefers the new worktree build (`bin/bench.sh:241-244`) and the guarded entries read the running executable (`internal/gate/checkpoint.go:234`, `internal/reviewrecord/recordcmd/refusal.go:45`). The install rule covers that worktree too, because the install runs after each Go change.
Won't handle: a signal during the rebuild — the commit has already published, and the publication rollback restores the prior build (`internal/freshness/publication_lock_test.go` (`TestPublishWaiterSeesRestoredTripleAfterInterruption`)).
Won't handle: a deleted file under a directory or glob `//go:embed` pattern in a worktree with no `dist/bench` — no embed in the build closure uses such a pattern today (research C), and a present `dist/bench` still rebuilds through the seal check.
Won't handle: a build that hangs — the step passes `context.Background()`, as the landing refresh does (`internal/worktree/land_effects.go:161`), so a hung build blocks the commit's return the way a hung `bench worktree build` blocks; the operator interrupts it, and the publication rollback keeps the prior build.
Won't handle: a published commit with no parent — a Bench worktree branch starts at a commit of the default branch, so `published^` always resolves; an unresolved diff answers false and builds nothing.

## Out of scope

- The same rebuild after a `bench worktree merge` fold that brings Go changes is a separate capability on another write verb. Estimate: 4 edits, 1 gate run.
- Dropping the build step of the `--preflight-build` chain, which a commit rebuild can make redundant, is a separate capability. Estimate: 3 edits, 1 gate run.

## Further notes

### Source-sentence-to-row table

| source sentence | rows |
| --- | --- |
| Criterion: "A commit in a build worktree leaves no stale dist/bench for the next preflight or --in verb." | CW13, CW14, CW16, CW20, CW21, CW27, CW37 |
| Row: "A ticket commit that changes a Go build input leaves the worktree binary at the previous source." | CW13, CW16, CW21 |
| Row: "The next `binary-seal` row then reds, and the author or the orchestrator runs `bench worktree build` by hand." | CW14, CW37 |
| Row: "`docs/known-issues.md` holds the workaround, and this row is its retire owner." | CW38 |
| Row: "A Markdown-only commit also makes `bench test --in` and `bench anchors --in` refuse with a build mismatch." | CW02, CW19 (premise contradicted; see below) |
| Row: "The reviewer decided on 2026-10-10 that `bench commit` rebuilds the worktree binary when the commit changes a Go build input." | CW13, CW16, CW21 |
| Row: "The build key ignores a Markdown-only change." | CW02, CW19 |
| Row: "FT344 owns the remedy wording, and FT398 owns a test that reads a stale binary." | CW22 |
| D1 (2026-11-02): seal-check trigger, and a diff test only for an absent build | CW16, CW17, CW19, CW20, CW21, CW41, CW43, CW49 |
| Spec review, iteration 1 (B3): the step reads no held-install predicate | CW45 |

The Markdown sentence of the row contradicts the tree. No Markdown file is a build input, and the digest inputs did not change after 8ae24503f. The 2026-10-10 refusals came from the pre-dffdef899 `--in` route and from Go edits in the same tickets. This is a non-behavioral contradiction: the spec obeys the tree convention, keeps the build key unchanged, and flags the premise for reviewer veto.

### Flagged additions

- CW01 to CW12 go beyond the source. The one-source rule of `AGENTS.md` requires one owner for the diff test and one owner for the refresh rule, because the commit step needs both rules that the landing owns today.
- CW20 and CW21 come from D1, not from the 2026-10-10 sentence.
- CW22 to CW25 come from the hostile-input checklist and the failure path of the build.
- CW27, CW28, and CW29 fix the order and the stream of the new step against the background grade and the commit record.
- CW33 comes from the checklist class of a command whose own write changes a fact that it reports.
- CW39 adds one sentence to the shipped Command Notes, so that a session outside this kit's profile knows the rule.
- CW34 to CW37, CW46, CW47, and CW48 are kept-rule witnesses, not new promises.
- CW41 to CW43 read D1's "the diff names a build input" over the parent commit too, so a deleted input counts. CW42 bounds that rule to the files that the digest read.
- CW44 is a behavior change of the landing: a landing that only deletes a Go build input now reports the install step. The shared owner makes it, and the one-source rule of `AGENTS.md` requires that owner. A deleted file that a build constraint excluded also gives that notice, although the broker does not change. This is the contestable call of this spec.
- CW50 comes from the spec review of iteration 3: it pins the stated rule that the second seal read wins over a build error.
- CW49 comes from the spec review of iteration 2: `--no-renames` keeps the old path of a renamed input in the diff.
- CW45 comes from the spec review of iteration 1: it reds a step that copies the held-install branch.

### Reader sweep

- Readers of a worktree `dist/bench` and its seal (research B): the `binary-seal` row (`internal/preflight/binary_seal.go:28-71`), the `--in` route (`internal/treetarget/build.go:20-51`), the doctor seal row (`internal/adopt/doctor_rows.go:249-261`), `bench commands --brief` (`cmd/bench/commands.go:49-72`), `freshness-check` (`cmd/bench/freshness_publish.go:26`), the wrapper (`bin/bench.sh:236-256`), and the background grade through the wrapper (`internal/gate/background_launch.go:216-232`). Each reads the rebuilt build with no edit. The wrapper prefers a worktree `dist/bench`, so a first build makes a bare `bench` in that worktree run the worktree build; `projects/benchkit.md:709-713` states that routing.
- Readers that never use a worktree build: the gate and `bench test` build a private run binary (`internal/runbinary/runbinary.go:174-214`), and the landing refresh grades only the destination (`internal/worktree/land_effects.go:148-174`).
- Readers of the commit output: `internal/commit/landing_test.go` (`TestGreenLandingPublishesTheNamedPath`), `internal/commit/merge_parent_test.go` (`TestCommitRecordsTheMergeParent`), `internal/worktree/commit_route_test.go` (`TestCommitExitThreeRouteReconcilesTheCheckout`), and `cmd/bench/commit_chain_test.go`. The step writes no stdout, so each keeps its answer.
- Guidance sweep for `worktree build` (research B, `rg --hidden -g '!.git' -n 'worktree build' .`): `docs/known-issues.md:10` retires. `docs/known-issues.md:12`, `.agents/commands/bench-final-check.md:31` (anchor-pinned at `internal/anchors/registry_retained_workflow.go:173`), `projects/benchkit.md:700`, `CONTEXT.md:212-216`, and the staged tickets of `specs/guidance-weight/` name a build after an uncommitted edit or a rehearsal, and they stay. The install and restore sentences of `.bench/BENCH-reference.md:416` and `:453` and of `.agents/commands/bench-implement-spec.md:46` are unrelated.
- No check grades known-issues rows against roadmap rows. The anchors registry is the one grader (`internal/anchors/registry_ticket_passes.go:15-24`).
- `package-core-guard` reds a claim word beside `specs/`, `decisions/`, or `tests/` in shipped Markdown (`internal/conformance/package_core_checks_test.go:196-222`). `CHANGELOG.md` and `.bench/BENCH-reference.md` ship, so their new sentences name no such path.
- Structure: `internal/commit/` holds 18 source files and `internal/freshness/` sits above its maximum (`bench structure`), so this spec adds no file there. `internal/worktree/land.go` is 396 lines and loses the predicate. `internal/freshness/freshness_verify.go` is 158 lines, `internal/freshness/freshness_buildinputs.go` is 233 lines, and `internal/commit/background.go` is 49 lines.

### Pre-review proof checklist

- Cited symbols: `commit.Run` (`internal/commit/commit.go:81`), `commitAttributed` (`internal/commit/commit.go:115`), `publicationRemainder` (`internal/commit/commit.go:267`), `startBackgroundGrade` (`internal/commit/background.go:26`), `runCommand` (`internal/commit/landing_test.go:21`), `freshness.Verify` (`internal/freshness/freshness_verify.go:20`), `freshness.RebuildIn` (`internal/freshness/freshness_verify.go:142`), `freshness.PublishedExecutable` (`internal/freshness/freshness_verify.go:120`), `freshness.DeclaresBuildInputs` (`internal/freshness/freshness.go:29`), `freshness.BuildInputs` (`internal/freshness/freshness.go:65`), `freshness.Publish` (`internal/freshness/freshness_publish.go:42`), `runbinary.BuildSubject` (`internal/runbinary/runbinary.go:279`), `brokerSourceChanged` (`internal/worktree/land.go:353`), `refreshBroker` (`internal/worktree/land_effects.go:148`), `readHeldInstall` (`internal/worktree/land_effects.go:180`), `binarySealFacts` (`internal/preflight/binary_seal.go:28`), `kittest.WriteTree` and `kittest.EditBuildInput` (`internal/treetarget/kittest/kittest.go:26,50`), `writeBuildFixture` (`internal/freshness/freshness_digest_test.go:282`), and `writePublishedFixture` (`internal/freshness/freshness_verify_test.go:338`).
- Import edges: `internal/freshness` gains a direct edge to `internal/git`; `go list -deps ./internal/git` names no `freshness` package. `internal/commit` gains direct edges to `internal/runbinary` and `internal/sanitize`; `go list -deps ./internal/commit` already names `internal/freshness` and `internal/runbinary` (research A). The `internal/commit` tests import `internal/treetarget/kittest`, which imports only `internal/freshness`.
- Source-row clauses and occurrences: the table above quotes each clause. The row records two occurrences, the 2026-10-08 FT349 build (`roadmap/FT406.md:14`) and the 2026-10-10 delegates (`roadmap/FT406.md:15`).
- Promised field labels: none. No intent-ledger or review-record field changes, so no `Persisted fields:` line applies.
- Changed-function callers: `brokerSourceChanged` has two production callers (`internal/worktree/land.go:193`, `internal/worktree/land_stack.go:85`) and no direct test caller. `refreshBroker` has its effect caller in `internal/worktree/land_effects.go` and the tests in `internal/worktree/land_effects_test.go`. `commit.Run` has the callers `cmd/bench/commit_chain.go` and `internal/reviewrecord/recordcmd` (research A).
- Copy survival: CW10 reviews `internal/worktree` for a surviving diff intersection. CW12 reds a `refreshBroker` that loses its held-install branch.
- Rendered-shape readers: the step adds two stderr lines and changes no existing line. The old-text needles `landing refresh build failed: ` and `landing refresh failed: ` keep their bytes.
- Pin operators: CW15 compares the whole stderr. CW22 compares the whole line. CW23 counts lines and tests for the absent ESC byte. CW29 tests the prefix of the last stdout line and the absence of `worktree build`.
- Entry reads: `Run` reads `benchhome.Dir()` and the package variable `worktreeBuild` at its boundary. `refreshWorktreeBuild` is the internal form with the home and the build injected.
- Derived expectations: each expected line is a literal, except the remedy, which comes from the FT344 owner outside this change.
- Consolidated rules: the diff test of the landing becomes `DiffChangesBuildInputs`. Old rule: the diff names a path of the post-tip input set. New rule: the diff names such a path, or it deletes a non-test Go file in the directory of a listed package. CW09, CW41, CW42, and CW44 map that change. The verify-build-verify rule of `refreshBroker` becomes `Refresh` with the same rule. CW11, CW12, CW46, and CW47 map that move.
- Quantified obligations: "a commit in a build worktree" quantifies over every lane-passed commit with a published subject and exit 0 in a root that declares build inputs. CW19, CW26, CW31, and CW32 bound the excluded commits.
- Workflow-step writes: the step writes `dist/bench`, `dist/bench.seal`, and `bin/bench-broker.manifest` in the worktree, all ignored (CW33). No checkpoint digest reads them: the chunk checkpoint and the landing grade the committed tree.
