# Add the KitForRoot kit rule to the gate package

Blocked by: none
Writes: internal/gate/kit_root.go (new), internal/gate/kit_source.go, internal/gate/kit_root_test.go (new), internal/gate/kit_value_test.go, cmd/bench/wrapper_kit_test.go (new), cmd/bench/main_test.go, bin/bench.sh, .bench/commitment.json
Covers: R01, R02, R03, R04, R05, R06, R07, R08, R09, R17, R18

## What to build

The gate package gains one exported pure form: `func KitForRoot(kit, root string) string`.
It gives the kit that grades `root`, by the two conditions of the wrapper's `kit_dir` rule in `bin/bench.sh`.
No production caller uses it in this ticket; ticket 02 makes `mergeOwner` call it.

Contract for ticket 02: `KitForRoot` answers `root`, unchanged as the caller spelled it, only when each of these conditions holds:

- `kit` is not empty.
- The canonical spelling of the git top level of `kit` is the canonical spelling of `kit`.
- The canonical spelling of the git common directory of `kit` is the canonical spelling of the git common directory of `root`.

Every other input answers `kit`, unchanged as the caller spelled it.
A git failure or a canonical-path failure at any step answers `kit`.

Put the form in `internal/gate/kit_root.go`, or in `internal/gate/kit_source.go` if that file stays below its budget.
Use the git adapter queries `RootAt` and `CommonDir` for the two git questions.
Use the existing path owner, as `KitSourceCheckoutAtKit` does through `resolvedPath`, for each canonical spelling.
The form reads no environment variable and no working directory.
It spells no `rev-parse` argument, so the `git-plumbing-owner` check stays green.
Do not change `kitRootAt`, `LaneForCommitAtKit`, `KitSourceCheckoutAtKit`, or `KitValue`.

Write one test function for each fixture shape in `internal/gate/kit_root_test.go`, with the names that the coverage map gives for R01 through R08.
Build each repository through `gittest` and the package's `runGitIn` helper (`internal/gate/gate_prose_staged_test.go`).
The ordinary-build census refuses a `worktree` or `init` literal in a `gittest.Run` or `gittest.Output` call in a gate test.
It also refuses an `os/exec` process in a gate test.
Each `.bench` kit in the R03 fixture is a directory that exists, because a missing directory makes git fail and hides the top-level condition.
The R06 test calls `t.Chdir` to the primary checkout of the root's repository, so it does not call `t.Parallel()`.

Add the R09 differential `TestShellWrapperKitMatchesKitForRoot` in the new file `cmd/bench/wrapper_kit_test.go`.
For each shape of `TestShellWrapperResolvesKitForTheInvocation` with no environment override, it runs the wrapper and reads the exported kit.
It then requires `gate.KitForRoot(<wrapper tree>, <working directory>)` to give the same kit.
The five shapes are the linked worktree, the adopted `.bench` kit, the main checkout, the different repository, and the plain directory.
Move the repository fixture and the wrapper call of the existing test out of `cmd/bench/main_test.go` into one shared helper, so one fixture serves both tests.
`cmd/bench/main_test.go` is 779 lines against a 400-line budget with no grant, so its line count only goes down.

Keep the Bash rule in `bin/bench.sh` as it is, because reviewer decision D1 accepts the Bash form and the Go form.
The file is 538 lines with no grant, so an edit to the `kit_dir` comment must add no net line.

The spec-stage close commits this spec's commitment binding inside the range that the plan-only preflight grades, so this ticket's `Writes:` line names `.bench/commitment.json`.
The build itself does not edit that file.

## Acceptance

- [ ] `gate.KitForRoot` has the signature `func KitForRoot(kit, root string) string`, and it answers `root` only under the three conditions above.
- [ ] A linked worktree of the kit's repository, under a parent directory whose name holds a space, gets that worktree (R01).
- [ ] A kit that is a sibling linked worktree gets the target worktree (R02).
- [ ] An adopted kit at `<repo>/.bench` stays the kit for a linked worktree of that repository (R03).
- [ ] A kit from a different repository stays the kit (R04).
- [ ] A root outside any repository keeps the kit (R05).
- [ ] An empty kit gets the empty kit when the working directory is the primary checkout of the root's repository (R06).
- [ ] A kit spelled through a symbolic link to the primary checkout gets the linked worktree (R07).
- [ ] A kit path that does not exist stays the kit (R08).
- [ ] For each of the five shapes, `gate.KitForRoot` gives the kit that the wrapper exports (R09).
- [ ] `TestKitValueIsTheOneBenchKitRead` stays green, so `KitValue` stays the one `BENCH_KIT` read in the gate package (R17).
- [ ] `bench test --check git-plumbing-owner` passes (R18).
- [ ] `TestShellWrapperResolvesKitForTheInvocation` keeps its seven shapes green, and `cmd/bench/main_test.go` has fewer lines than before.
- [ ] A probe that removes the top-level condition from `KitForRoot` reds R03 and R09; the restore turns both green again.
