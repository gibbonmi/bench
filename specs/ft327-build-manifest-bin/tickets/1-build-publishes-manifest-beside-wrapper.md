# The worktree build publishes the manifest beside the wrapper

Blocked by: none
Writes: internal/worktree/build.go, internal/worktree/build_test.go, internal/worktree/parallel_census_test.go, CHANGELOG.md
Covers: none

## What to build

FT327 has a build face. `bench worktree build` builds the worktree's own published
executable at `dist/bench`. The verb called the private build join. That join passes
`--manifest-dir` with the directory of the output, so the build script wrote the broker
manifest into `dist/`. The doctor row and the landing read the manifest beside the
wrapper in `bin/`. Thus the manifest that they read did not change after a build.

The build package has a subject form for this case. That form passes no manifest
directory, so the script uses its own default, which is `bin/`. The landing rebuild and
the wrapper's land rebuild already use that form, and `bench doctor --fix` also writes
the manifest beside the wrapper. The verb now calls the subject join.

## Acceptance

- [x] `bench worktree build` gives the build script no `--manifest-dir` operand, so the script publishes the manifest at its `bin/` default.
- [x] The build rows that use a recorder replace the subject join, which is the join that the verb calls.
- [x] A change of the verb back to the private build join turns the new test red.

## Verification

`TestBuildLeavesTheManifestDirectoryToTheScript` plants the stub build script and runs
the verb through the production joins. The stub records its arguments. Before the fix,
the test was red because the arguments held `--manifest-dir` with the worktree's `dist/`
path. After the fix, the test is green. The `bench test --package ./internal/worktree/...`
run passes; its skips are capability skips for sockets on this host. A `bench probe` swap
of `j.buildSubject` back to `j.build` in `build.go` bit the test, and the probe restored
the file.
