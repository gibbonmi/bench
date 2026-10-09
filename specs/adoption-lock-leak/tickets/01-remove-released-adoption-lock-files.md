# Remove released adoption lock files

Blocked by: none
Writes: internal/adopt/transaction/lock.go, internal/adopt/transaction/lock_test.go (new), docs/known-issues.md
Covers: none

## What to build

An adoption writer removes each of its lock files when it releases its lease. The writer removes the file while it still holds the lock, and then closes the file. A writer that takes a lock makes sure that the lock file path still names the file that it locked. If the path is gone or names a different file, the writer releases that lock and tries again. Thus the shared lock directory does not grow with each Link, Setup, Upgrade, Unlink, or compatibility repair. Two writers on the same destination continue to exclude each other.

## Acceptance

- [ ] After a lease on many destinations closes, the lock directory holds no lock file for them.
- [ ] After a lease request fails on a held destination, the lock directory holds only the lock file of the holder. After the holder closes its lease, the directory holds no lock file.
- [ ] Concurrent writers on one destination never hold it at the same time while released lock files are removed.
