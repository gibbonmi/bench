# Exercise the hostile environment and process cases

Blocked by: exercise-hostile-argv-grammar-cases.md
Ownership fence: `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`, `specs/axi-compatibility-oracle/testdata`
Integration surfaces: hostile argv matrix, comparator, and candidate rebuild→exercise-hostile-argv-grammar-cases.md; hostile environment fixtures→`specs/axi-compatibility-oracle/testdata`; environment and process case derivation, plus the inventory-class totality check→`internal/axi/compatibility`; exact executables and wrapper process identity→`scripts/go-build.sh` and `bin/bench.sh` exercised unchanged by every HE1 row; outside-repository refusal owner→`internal/toon/toon.go` (`NotInRepo`) exercised unchanged by HE1/outside-repository
Contracts: the hostile environment case IDs `<member>-hostile-<class>` and the process facts — argv vector, cwd path, `PATH` value, invocation path, child deadline, and post-run descendant set — cross the environment fixtures→`cmd/bench/axi_compatibility_test.go`; their type is one baseline observation record plus one teardown assertion per case, membership is the eight declared environment and process classes over every member the inventory marks for them, order is setup then run then teardown, and a case whose teardown leaves a descendant fails rather than being skipped; asserted by HE1 against really spawned bounded processes
Closure: HE1/outside-repository, HE1/deep-cwd, HE1/path-with-spaces-and-globs, HE1/stripped-path, HE1/missing-optional-tool, HE1/symlink-invocation, HE1/control-bearing-git-text, HE1/bounded-interrupt, HE1/descendant-teardown, HE1/inventory-class-totality

## What to build

The last hostile classes cannot be observed in process: they are about where the
executable runs, what it can find, how it is invoked, and what it leaves behind.
Each case here spawns a fresh baseline and candidate process from the exact
executables, with an explicit cwd, an explicit `PATH`, and an explicit
`exec.CommandContext` deadline, and compares all four observations plus one
teardown fact — that no descendant of the case survives its deadline.

This ticket also closes CO6: after its cases land, every hostile class the
command-help inventory declares resolves to at least one derived case, and
`HE1/inventory-class-totality` keeps that closure honest rather than leaving it as
a claim in prose.

Mutations are applied to a scratch copy of the tree from which the candidate
executable is rebuilt through `scripts/go-build.sh`, except the two harness facts —
the child deadline and the descendant sweep — which are mutated in the case runner
itself. The rebuild is bounded at 180s, each child at 30s, and the interruption
case at a 10s post-signal wait, so every row's red is a bounded failure.

## Acceptance

- [ ] [HE1] (covers CO6) every declared hostile environment and process class compares byte-exact on stdout, stderr, exit, and acceptance from a fresh bounded process, every case leaves no surviving descendant, and every hostile class the command-help inventory declares resolves to at least one derived case.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| HE1/outside-repository | resolve the repository root by walking up past the process cwd to the filesystem root in the candidate rebuild | the hostile-environment test | run `go test ./cmd/bench -run TestHostileEnvironmentObservationsAreExact/outside_repository -timeout 900s`; it fails at the raw stderr equality assertion for case `root-diff-hostile-outside-repo`, which loses the `toon.NotInRepo` line, and at the exit assertion reporting 0 against the baseline's non-zero; the child runs with cwd set to a fresh non-repository directory under a 30s `exec.CommandContext` |
| HE1/deep-cwd | resolve relative inputs against the repository root instead of the process cwd in the candidate rebuild | the hostile-environment test | run `go test ./cmd/bench -run TestHostileEnvironmentObservationsAreExact/deep_cwd -timeout 900s`; it fails at the raw stdout equality assertion for case `root-anchors-hostile-deep-cwd`, invoked from a six-level subdirectory, reporting the root-relative path where the baseline echoes the cwd-relative one; bounded by the 30s child deadline |
| HE1/path-with-spaces-and-globs | quote repository paths with `strconv.Quote` before emitting them in the candidate rebuild | the hostile-environment test | run `go test ./cmd/bench -run TestHostileEnvironmentObservationsAreExact/spaces_and_globs -timeout 900s`; it fails at the raw stdout equality assertion for case `root-outline-hostile-space-glob-path`, whose fixture directory is named `a b*c`, printing the added quotes; bounded by the 30s child deadline |
| HE1/stripped-path | resolve a helper tool from the ambient `PATH` instead of the absolute path the command already holds in the candidate rebuild | the hostile-environment test | run `go test ./cmd/bench -run TestHostileEnvironmentObservationsAreExact/stripped_path -timeout 900s`; it fails at the raw stdout equality assertion for case `root-structure-hostile-stripped-path`, run with `PATH=/nonexistent`, reporting the tool-not-found refusal where the baseline succeeds; bounded by the 30s child deadline |
| HE1/missing-optional-tool | treat a missing optional tool as a hard failure instead of the baseline's degraded result in the candidate rebuild | the hostile-environment test | run `go test ./cmd/bench -run TestHostileEnvironmentObservationsAreExact/missing_optional_tool -timeout 900s`; it fails at the exit assertion for case `root-doctor-hostile-missing-tool`, reporting 1 against the baseline 0 with the degraded body replaced; bounded by the 30s child deadline |
| HE1/symlink-invocation | canonicalize `os.Args[0]` with `filepath.EvalSymlinks` before recording the executable identity in the candidate rebuild | the hostile-environment test | run `go test ./cmd/bench -run TestHostileEnvironmentObservationsAreExact/symlink_invocation -timeout 900s`; it fails at the raw stdout equality assertion for case `root-version-hostile-symlink-invocation`, invoked through a symlink named `bench-alias`, reporting the resolved target path where the baseline echoes the invocation path; bounded by the 30s child deadline |
| HE1/control-bearing-git-text | stop escaping control runes in Git-sourced text before rendering it in the candidate rebuild | the hostile-environment test | run `go test ./cmd/bench -run TestHostileEnvironmentObservationsAreExact/control_bearing_git_text -timeout 900s`; it fails at the raw stdout equality assertion for case `root-diff-hostile-control-git-text`, whose fixture commit subject carries a bare carriage return, reporting the raw control byte against the baseline's escape; bounded by the 30s child deadline |
| HE1/bounded-interrupt | remove the `exec.CommandContext` deadline from the interruption case runner so the signalled child is waited on indefinitely | the interruption test | run `go test ./cmd/bench -run TestHostileProcessesTerminateWithinTheirDeadline/interrupt -timeout 900s`; it fails at the assertion that case `root-shift-hostile-interrupt` returned a terminal observation within its 10s post-signal wait, reporting the still-running child rather than blocking, because the runner keeps the outer `context.WithTimeout` as the wait's own bound |
| HE1/descendant-teardown | skip the post-case descendant sweep and report the case as complete when the direct child exits | the process-teardown test | run `go test ./cmd/bench -run TestHostileProcessesLeaveNoDescendant -timeout 900s`; it fails at the assertion that the case's process group is empty after teardown, naming case `root-shift-hostile-interrupt` and the surviving grandchild PID; the sweep itself runs under a 10s deadline |
| HE1/inventory-class-totality | derive hostile cases only for members whose inventory row marks `R`, skipping the `St` and `Rc` rows | the hostile-class totality test | run `go test ./cmd/bench -run TestEveryHostileInventoryClassHasACase -timeout 900s`; it fails at the unresolved-class assertion naming the first uncovered pair — member `nested-worktree-recovery`, class `Rc` — and the count of remaining uncovered pairs; the inventory parse runs under a 15s deadline |
