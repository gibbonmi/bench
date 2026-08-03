# Classify decayed checkouts as liveness and delete the prepared-abandon exemption

Blocked by: none
Ownership fence: `internal/specbuild/precondition.go`, `internal/specbuild/abandon_test.go`
Assumptions: `liveCheckout` (`internal/specbuild/precondition.go:229-245`) today collapses every non-absent fault into `errOwnership`; `ownedAssignments` (199-223) treats a missing intent registration (`!found` at 209-212) as identity; the blanket exemption is the `abandon.State == "prepared" && abandon.Result != ""` block at `preconditions` lines 94-99; specbuild abandon tests drive a counting fake owner, so no `internal/worktree` change is needed here; claims re-derived from the tree at pickup

## What to build

Liveness is decided by shape, not probe failure. The checkout probe
classifies: no filesystem entry, a dangling symlink, a non-directory entry
(regular file, FIFO, device node), or a directory with no git metadata entry
(husk) — liveness: abandon proceeds and releases ownership, every other
mutation refuses. A directory with a git metadata entry whose probe fails, a
resolvable checkout whose git common dir differs from the repository's, and
any `Lstat` failure other than not-exist — identity: fatal for every
operation. The probe never opens the path (stat and git only, and git runs
only against a directory carrying git metadata, so a FIFO cannot block). A
missing intent registration reclassifies as liveness. The blanket
prepared-abandon exemption in `preconditions` is deleted — the narrowed inner
classification is the one source of softening — and the enumerated identity
faults (duplicate assignment ID, path, or owner request; owner-request digest
mismatch; registration resolving to a different assignment ID or worktree
path; foreign checkout) refuse abandon even when a prepared abandon operation
with a recorded result exists. Fix the stale comment above the deleted block
to match. Prior art: `abandon_test.go`'s removed-worktree and forged-identity
fixtures.

## Acceptance

- [ ] [DC1] Abandon and ApplyAbandon proceed for a husk (directory present, git metadata gone), releasing the assignment (registration and intent entry released through the owner, run terminal).
- [ ] [DC2] One table test drives husk, dangling symlink, FIFO, and regular file through Abandon (proceed) and Checkpoint (refuse); the probe returns without opening the FIFO.
- [ ] [DC3] A directory carrying git metadata whose probe fails (permissions) refuses every operation, abandon included.
- [ ] [DC4] A foreign checkout at the assignment path still refuses abandon (regression control, unchanged test).
- [ ] [DC5] Abandon proceeds when the intent registration is deleted while the run record still lists the assignment.
- [ ] [DC6] Every enumerated identity class — duplicate ID, duplicate path, duplicate owner request, digest mismatch, registration resolving to a different assignment ID, registration resolving to a different worktree path, foreign checkout — refuses abandon even with a prepared abandon operation and recorded result present (table over the classes).
- [ ] [DC7] The apply re-entry and journal reconcile tests run unchanged and green (the mid-apply resume the exemption existed for).

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DC1 | keep husk classified as `errOwnership` | the husk abandon test | apply, run `go test ./internal/specbuild -run Abandon`, expect the ownership-refusal failure |
| DC2 | classify by directory-existence alone | the dangling-symlink table case | apply, run the table test, expect the symlink failure |
| DC3 | soften every probe failure to absence | the unreadable-with-metadata test | apply, run it, expect the missing-refusal failure |
| DC5 | delete the exemption without reclassifying not-found | the missing-registration test | apply, run it, expect the ownership-refusal failure |
| DC6 | restore the prepared-abandon exemption block | the identity-class table | apply, run it, expect every class to report the missing refusal |
