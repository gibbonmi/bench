# Hold one run binary across several focused runs

Blocked by: 1-repeat-the-probe-package.md
Writes: internal/runbinary/runbinary.go, internal/runbinary/runbinary_test.go, internal/testreport/cancel_test.go, internal/testreport/command.go, internal/testreport/held_runs.go (new), internal/testreport/held_runs_test.go (new), internal/testreport/outcome.go, internal/testreport/runbinary_test.go
Covers: PM27, PM28, PM29, PM31, PM70, PM76

## What to build

`internal/testreport` gains one run-owner seam, the held runs, that serves several focused runs over one prepared request.
The caller passes its cancel context, and each selection, build, and Go child of each held run runs under that context.
The first run selects the run binary through `selectRunBinary`, as `Execute` does today.

Before each later run, the owner validates the held binary through the run-binary owner's existing validation.
That validation runs `freshness.VerifyAt`, so the freshness digest stays the one source of closure membership.
A held binary that still validates serves the run.
When an owned held binary no longer validates, the owner keeps it and builds one fresh run binary for that run only.
The owner removes that fresh binary's private directory when the run returns.
A stale inherited binary refuses that run, and the owner builds nothing.

`internal/runbinary` exports one accessor on `Selection` that reports whether the factory built the binary or inherited it.
The held runs read that accessor, so no code reads `BENCH_RUN_BINARY` a second time.
If the held runs need the private validation of `runbinary`, export that one rule from `runbinary`.
Do not copy it into `testreport`.

A run whose context is cancelled when its validation or its build returns reads the outcome kind `interrupted`, not `refused`.
The owner counts each binary that it builds, and its close removes the held private directory.
A prose request selects no run binary, so its build count is 0.
`Execute` and the single form keep their current behavior.

Drive the seam through a factory-valued test seam with a stub builder and a stub verifier.
Today `installTestSelectionFactory` installs only the `ReuseOrOwn` function value, so extend that seam and keep its current callers green.
The cancelled-build row reuses the parking builder of `internal/testreport/cancel_test.go`.
`internal/testreport` already holds more files than its structure budget, so put the owner and its tests in the two new files only.

## Acceptance

- [ ] After the verifier refuses the held binary, the next held run builds one fresh run binary whose builder reads the changed tree.
- [ ] Each held run after the first calls the verifier on the held path, and an accepting verifier gives reuse with no build.
- [ ] With an inherited `BENCH_RUN_BINARY` that the verifier refuses on the second run, that run reads `refused` and the builder records no call.
- [ ] After a stale run returns, its fresh private directory is gone, and after close no held directory remains.
- [ ] A cancel while a stale run's fresh build parks gives that run `interrupted`, and the parked builder's directory is gone.
- [ ] A cancel while the held validation of an inherited binary runs gives that run `interrupted`, not `refused`.
