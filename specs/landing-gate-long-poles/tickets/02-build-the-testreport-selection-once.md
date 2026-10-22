# Build the testreport owned selection once

Blocked by: none
Writes: internal/testreport/selection_once_test.go (new), internal/testreport/root_conformance_test.go, internal/testreport/check_test.go, internal/testreport/runbinary_test.go, internal/runbinary/runbinary.go, internal/runbinary/shared.go (new), internal/runbinary/shared_test.go (new), internal/worktree/test_run_test.go
Covers: LGP6, LGP7, LGP8, LGP34

## What to build

The two kit-copy tests of the testreport package share one owned run binary for each package run. Today each one builds its own through `installOwnedSelection` (`internal/testreport/root_conformance_test.go:34-38`).

Follow the worktree selector (`internal/worktree/test_run_test.go:16-52`). A package variable holds one `sync.Once` build and a mutex-guarded build log. An inherited `BENCH_RUN_BINARY` wins, and a present invalid selection refuses.

A `TestMain` in the new `internal/testreport/selection_once_test.go` closes the selection after the run. The package has no `TestMain` today.

Put the shared build in `installOwnedSelection` only. `TestEachExecuteBuildsItsOwnSelection` (`internal/testreport/runbinary_test.go:100-124`, PS32) requires two builds for two `Execute` calls under `installTestSelectionFactory`. Keep that helper and its 30 other call sites unchanged. The 31st call is inside `installOwnedSelection` itself (`internal/testreport/root_conformance_test.go:36`), which this ticket rewrites.

The planted defect of `TestRootConformanceCheckPrintsFailures` is in `.agents/commands/bench-drain.md`. That file is outside the run-binary build inputs (`scripts/go-build.inputs`), so both kit copies share one digest. Both tests keep their real children and their kit copies.

`internal/testreport/check_test.go` has 481 lines, over its 400-line limit with no grant. Do not add lines to it. Put the new test in `selection_once_test.go`.

Record `bench test --full --package ./internal/testreport` before and after the change, for LGP34.

## Acceptance

- [ ] The two kit-copy tests together invoke the run-binary builder at most once (LGP6).
- [ ] The planted row-next-grammar defect still reaches its diagnostic through the real child (LGP7).
- [ ] The named check still runs only the registered dev-scope check through the real selection (LGP8).
- [ ] `TestEachExecuteBuildsItsOwnSelection` still sees two builds.
- [ ] The probe "build the selection in each installOwnedSelection call" bites.
- [ ] The verification record holds the package time before and after the change (LGP34).
