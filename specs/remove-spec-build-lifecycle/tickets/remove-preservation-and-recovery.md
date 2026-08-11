# Remove preservation refs and the recovery verb

Blocked by: delete-spec-build-core.md

## What to build

`bench worktree recovery` leaves the grammar (unknown-verb structured error,
help omits it). `bench resume` stops authoring preservation refs; its
reconcile deletes refs under `refs/bench/specbuild/` and
`refs/bench/recovery/` only — `refs/bench/green/<branch>` and diagnostic refs
survive byte-identical — and purges lifecycle-typed or unknown-typed
assignments from the intent ledger. Idempotent: a second run is a no-op; an
interrupted run converges on re-entry. `internal/shift`'s recovery-ref
authoring goes with it. Guard test for the green-ref survival is written
first at the new reconcile seam (RM10), then the interruption test (RM9).
Covers RM3, RM4, RM5, RM9, RM10.

## Acceptance

- [ ] `bench worktree --help` omits recovery; the verb answers with the
      structured unknown-verb error.
- [ ] In a throwaway repo seeded with lifecycle refs, recovery refs, a green
      ref, and legacy ledger entries: one resume empties both lifecycle
      namespaces, leaves `refs/bench/green/<branch>` byte-identical, purges
      legacy entries, authors no new refs; a second resume is a no-op; a
      kill-and-rerun mid-reconcile converges.
- [ ] Every surviving ledger entry is worktree-pool-typed with an existing
      worktree path.
- [ ] `go test ./internal/worktree/... ./internal/intent/... ./internal/harness/... ./internal/shift/...` green.
