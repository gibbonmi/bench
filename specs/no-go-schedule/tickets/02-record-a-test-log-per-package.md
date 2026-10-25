# Record one Go test log for each test package of a whole kit run

Blocked by: none
Writes: internal/gate/testlog_exec.go (new), internal/gate/testlog_exec_test.go (new), internal/gate/gate_go.go, internal/gate/phases.go, internal/subprocess/replace.go (new), internal/subprocess/replace_test.go (new), internal/conformance/ordinary_build_census_test.go, internal/gate/branch_native_phases_test.go, projects/benchkit.md
Covers: NG5, NG6, NG7, NG8, NG9, NG10, NG11, NG12, NG72

## What to build

`bench gate-go exec <binary> <args...>` is the trace shim. When `BENCH_GATE_TESTLOG` names a directory, the shim creates `<dir>/<package key>.log` first. It then replaces its image with `<binary> -test.testlogfile=<file> <args...>`. Without `BENCH_GATE_TESTLOG`, it replaces its image with `<binary> <args...>`. The `gate-go` usage line names the new step.

The package key is the test binary's working directory relative to the graded root. Encode it so that two different directories never share one name, for example `a/b` and `a__b`.

One new exported function of `internal/subprocess` replaces the image, beside the start gate that owns the other image replacement (`internal/subprocess/startgate.go:46`).

`toolchainPhases` gives the kit's whole `test` phase the traced argv only when the graded root is the kit. The argv is these seven elements: `go`, `test`, `-trimpath`, `-count=1`, `-timeout=10m0s`, `-exec=<bench-run-binary> gate-go exec`, and `./...`. A linked root, the delta set, and every other test phase keep their argv.

`withRunBinary` replaces the run-binary token inside an `-exec=` element with the selection path in quoted form, so that `cmd/go` keeps a path with a space as one word. Its lane caller (`internal/gate/lane.go:306`) carries no `-exec=` element and keeps its argv.

Update the three readers of the old argv: the `ordinary-build-census` check's expected `test` argv, the expectation at `internal/gate/branch_native_phases_test.go:104`, and the phase-table row at `projects/benchkit.md:346`. The census check keeps its exact `slices.Equal` comparison.

`internal/gate/phases.go` holds 388 lines. Put the shim, the key encoding, and the traced phase builder in the new file, so `phases.go` stays at 400 lines or fewer.

## Acceptance

- [ ] With `BENCH_GATE_TESTLOG` set, a probe binary run through the shim records `-test.testlogfile=<dir>/<package key>.log` before its other arguments.
- [ ] Without `BENCH_GATE_TESTLOG`, the probe binary records its arguments unchanged.
- [ ] A probe binary that exits at once leaves an empty log file.
- [ ] The package directories `a/b` and `a__b` map to two different log names.
- [ ] The `ordinary-build-census` check accepts the seven-element argv and reds the old argv and a split `-exec` element.
- [ ] `withRunBinary` turns `-exec=<bench-run-binary> gate-go exec` into an element whose first quoted word is the selection path.
- [ ] A linked root's `test` phase argv carries no `-exec=` element.
- [ ] A phase-child run of the whole kit table over a fixture module, with `BENCH_GATE_TESTLOG` set, leaves one log that starts with `# test log` for each test package.
- [ ] A real `go test` run with a run-binary path that holds a space starts the shim at the whole path.
- [ ] The profile's phase-table row shows the seven-element argv.
