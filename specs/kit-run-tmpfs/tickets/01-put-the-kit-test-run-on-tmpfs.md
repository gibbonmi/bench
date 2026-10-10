# Put the kit test run on tmpfs

Blocked by: none
Writes: internal/env/kit_run.go, internal/env/kit_run_test.go
Covers: none

## What to build

`OpenKitTestRun` puts its private run directory on tmpfs when the caller sets no TMPDIR. Today the run directory goes to `os.TempDir()`. On a WSL2 host, that directory is on ext4, and fsync takes most of the adopt, repairtest, and system test time. The repairtest package takes 62.7 s with TMPDIR=/tmp and 3.3 s with TMPDIR=/dev/shm.

The base directory comes from this order:

1. A TMPDIR that the caller sets keeps priority. A run directory that the caller TMPDIR refuses is an error, as today.
2. With no caller TMPDIR, the base is `/dev/shm` when it is an existing, writable directory.
3. Otherwise, the base is `os.TempDir()`, as today.

An unexported helper takes the tmpfs candidate as a parameter, so a test can prove each branch without the host's `/dev/shm`. The change adds no environment variable and no exported setting. The gate and `bench test` both call `OpenKitTestRun`, so both get the new base.

## Acceptance

- [ ] With a caller TMPDIR and a writable candidate, the run directory is under the caller TMPDIR.
- [ ] With no caller TMPDIR and a writable candidate directory, the run directory is under the candidate.
- [ ] With no caller TMPDIR, the run directory is under `os.TempDir()` when the candidate is absent, is a regular file, or is not writable.
- [ ] The existing kit run tests stay green and unchanged, including the 16-byte bound on the private TMPDIR growth.
- [ ] The author reports the repairtest package elapsed time on the ticket tip with no TMPDIR set.
- [ ] The author reports whether a socket row in the bounds or landing packages skips under the new base.
