# Plan reclaimable pool keys

Blocked by: none
Writes: internal/worktree/pool_reclaim.go, internal/worktree/pool_reclaim_test.go, internal/usage/worktree.go, cmd/bench/main.go

## What to build

`bench worktree reclaim`, planning only. It enumerates the direct children of
`$BENCH_HOME/worktrees`, classifies each against one reclaimability predicate,
and prints a TOON table of key, verdict, and reason, an aggregate count, and the
exact `--apply <fingerprint>` invocation that would act on the plan. It removes
nothing.

The predicate lives in `internal/worktree/pool_reclaim.go` and is the only
derivation of this fact in the tree: a key is reclaimable when it holds nothing
at top level, or when every top-level entry is a real directory holding a regular
`.git` file whose `gitdir:` target is provably absent (`os.Lstat` reporting
`IsNotExist`, never any other error). `os.Lstat` throughout — no symlink on the
key, child, or `.git` path is followed. The current repository's key, `Pool(root)`,
is excluded before the predicate runs. Every retention carries a reason naming
what protected it.

The fingerprint is derived with the package's existing `fingerprintParts` helper
over the sorted reclaimable key names and each key's child `gitdir:` targets, so
ticket 02 consumes exactly this value. The command registers in
`cmd/bench/main.go` with its help row and its grammar in
`internal/usage/worktree.go`, staying outside the approved AXI query set and
keeping `bench worktree clean`'s output contract: TOON on stdout, a zero-row
empty state, structured refusals on stdout, exit 0/1/2.

Every test binds its own `BENCH_HOME` under its `t.TempDir`.

## Acceptance

- [ ] a pool holding a dangling-pointer key, an empty key, and a live key plans exactly the first two, as TOON rows carrying key, verdict, and reason (PL1, PL5).
- [ ] the live-pointer, `.git`-is-a-directory, no-`gitdir:`-line, stray-top-level-entry, and one-live-one-dead keys are each retained with a distinct reason (PL2, SH1, SH2, SH3).
- [ ] an unreadable child, and a `gitdir:` target whose stat fails for a reason other than absence, retain the key and are reported (SH4).
- [ ] a symlinked key, a symlinked child, and a symlinked `.git` are each retained unfollowed (SH5).
- [ ] the current repository's own key is retained even when it is empty (SH6).
- [ ] an empty pool and an absent `worktrees` directory each print a zero-row table and exit 0 (PL4).
- [ ] the sorted recursive listing of the pool is byte-identical before and after a bare plan (PL5).
- [ ] the plan prints the apply invocation carrying the fingerprint (PL3), and running the command outside a repository refuses with the not-in-repo error.
