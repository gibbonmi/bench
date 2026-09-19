# 1. Open a kit test run for the gate's kit phases

Blocked by: none
Writes: internal/env/, internal/gate/phases.go, internal/gate/phases_test.go, internal/gittest/gittest.go
Covers: TD1, TD2, TD3, TD4, TD5, TD6, TD8, TD9, TD10, TD11, TD12, TD13, TD14, TD15

## What to build

Chunk: TD-C1a.

Add the kit test run to the environment package. A runner opens the run from its base environment. Open creates one run directory directly under the base `TMPDIR`, with a private home directory and a private temporary directory inside it. The run's entries set `HOME`, `TMPDIR`, `GIT_CONFIG_GLOBAL` to the null device, `GIT_CONFIG_NOSYSTEM` to `1`, and the git maintenance policy. The entries also pin `GOMODCACHE`, `GOPATH`, and `GOENV` to the values that Go resolves under the base environment. Close restores owner permissions inside the run directory and removes it.

Open refuses before it creates anything when the base `HOME` is absent, empty, or relative. It refuses when it cannot create the run directory under the base `TMPDIR`. Each refusal names the variable.

The gate's phase command opens one run for the kit phases of a kit root. It merges the entries into each phase and closes the run after the phases end. A linked root opens no run. The existing git test configuration function stays until ticket 2 moves its last callers.

Add a kit-run probe beside the maintenance probe in the shared git test scaffold. The probe takes the expected operator values, checks each property of the entries in the child, and prints the first violated property.

## Acceptance

- [ ] Under the run's entries, a global git marker and a system git marker from the base environment are not readable.
- [ ] Under the run's entries, a commit starts no git auto-maintenance.
- [ ] `HOME` and `TMPDIR` name empty directories inside the run directory, and the `TMPDIR` entry is at most 16 bytes longer than the base value.
- [ ] `go env GOMODCACHE GOPATH GOENV` prints the same lines under the entries as under the base environment.
- [ ] Close removes the run directory although a mode-0 directory with a file is inside it.
- [ ] Open refuses an absent `HOME`, a relative `HOME`, and a `TMPDIR` that is a regular file, and each refusal names the variable.
- [ ] A base `TMPDIR` with a space in its path works through open, the probe, and close.
- [ ] A gate kit fixture phase runs the kit-run probe and the gate exits 0, and a linked root's phases carry no `HOME` entry.
