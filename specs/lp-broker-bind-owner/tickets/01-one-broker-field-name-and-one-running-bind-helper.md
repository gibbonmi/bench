# One broker field name and one running bind helper

Blocked by: none
Writes: internal/brokermanifest/brokermanifest.go, internal/adopt/broker.go, internal/adopt/doctor_rows.go, internal/worktree/land_effects.go, internal/worktree/build.go, internal/commitment/commitmenttest/admission.go, internal/gate/complete_checkpoint_broker_test.go, internal/worktree/land_reuse_test.go
Covers: none

## What to build

The `brokermanifest` package owns the field names of the broker manifest.
Some production readers outside that package spell the path, version, and digest fields as literals, so each of these field names has more than one source.
`brokermanifest` now also declares `VersionField` beside `PathField` and `DigestField`.
Each production reader of these fields in `internal/adopt` and `internal/worktree` now uses the constants.
Tests that write or read the manifest file keep their literals, because they pin the format on disk.

Two test packages bind the running test binary in a broker manifest.
Each package pairs `os.Executable()` with `commitmenttest.BindManifest` in its own helper.
A new `commitmenttest.BindRunning` helper is now the one source of this bind step.
The gate tests and the worktree landing tests call it.

The persisted checkpoint test in the external gate test package still reads `gate.InstallRunning`.
Another ticket owns that file, so `gate.InstallRunning` stays as a name for `commitmenttest.BindRunning`.

The candidate check test in `internal/commitment/repository` also calls `os.Executable()`.
That test publishes the manifest through the production writer and does not call `BindManifest`, so it keeps its own flow.

## Acceptance

- [ ] No production reader outside `internal/brokermanifest` spells a manifest key as a literal.
- [ ] `brokermanifest` declares `VersionField`, and its writer and its bound field list use it.
- [ ] `commitmenttest.BindRunning` is the only place that pairs `os.Executable()` with `BindManifest`.
- [ ] The gate tests and the worktree landing tests bind the running test binary through `commitmenttest.BindRunning`.
- [ ] `bench test` is green for each touched package.
