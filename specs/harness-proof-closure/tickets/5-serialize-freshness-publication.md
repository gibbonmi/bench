# Serialize freshness publication

Blocked by: none
Writes: internal/freshness/freshness_publish.go, internal/freshness/publication_lock.go (new), internal/freshness/freshness_publish_test.go, internal/freshness/publication_lock_test.go (new)
Covers: HP14, HP15, HP16, HP17, HP18

## What to build

Serialize every freshness publication that targets the same broker-manifest
directory. Acquire an exclusive kernel lock on the already-existing manifest
directory before any live executable, seal, or manifest state is read or moved.
Hold it through publication, rollback, residue cleanup, and signal teardown.

Use separate processes and bounded handshakes to prove contention. If the first
publisher is interrupted before its seal lands, it restores its prior triple and
releases the lock before the waiter publishes. A process exit releases the lock
without a persistent lock-file artifact.

## Acceptance

- [ ] A second publisher for one manifest directory cannot enter live mutation while the first holds the publication lock.
- [ ] After two successful contenders finish, the executable, seal, and manifest describe one complete published version.
- [ ] Interruption before the first seal lands restores the prior triple, releases the lock, and lets the waiter complete.
- [ ] Process exit cannot leave a stale lock artifact that blocks the next publication.
- [ ] A missing, wrong-type, or unreadable manifest directory refuses before any member of the live triple changes.
