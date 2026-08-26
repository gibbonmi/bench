# Derive the census home beside the worktree pool

Blocked by: none
Writes: internal/poolkey/poolkey.go (new), internal/poolkey/poolkey_test.go (new), internal/census/census.go (new), internal/census/census_test.go (new), internal/worktree/worktree.go, internal/worktree/worktree_test.go, internal/worktree/pool_reclaim_test.go, internal/benchguard/benchguard.go, internal/benchguard/benchguard_test.go

## What to build

This ticket is the prefactor. It opens the seams that every later ticket
composes, and it fixes where the census records live.

`internal/census` needs the repository pool key. `internal/worktree` needs the
census drop. A direct import in both directions makes an import cycle. Move
the pool-key derivation and the canonical-root derivation out of
`internal/worktree/worktree.go` into a new low package, `internal/poolkey`.
`poolAt` and every `canonicalRoot` caller then call that package, and the
pool path stays byte-identical for a canonical root.

`Key` canonicalizes any root first. `bench worktree-pool` inside a linked
worktree therefore names the primary pool, not a pool keyed by the worktree
path. State that change in the ticket's landing note. The checksum helper
and its POSIX `cksum` goldens in `internal/worktree/worktree_test.go` move
with the key, so one source derives the key and one test pins it.

Add `internal/census` with the records' home only. The records sit under
`<home>/census/<repo-key>/`, a sibling of `<home>/worktrees/<repo-key>/`.
`bench worktree reclaim` reads only `<home>/worktrees`, so the sibling
directory never changes a reclaim plan.

Export the guard's wrapper-aware `bench` test, because the census must apply
the same test. The test walks each simple command and one wrapper level
(`sh -c`, `bash -c`, `zsh -c`), as `Classify` does today. Keep every current
caller.

These contracts cross into tickets 02 to 06:

- `poolkey.Canonical(root string) string` returns the primary checkout root for any root inside the repository.
- `poolkey.Key(root string) string` returns `<base>-<crc32>` for `Canonical(root)`.
- `census.Dir(home, root string) string` returns `<home>/census/<poolkey.Key(root)>`.
- `benchguard.InvokesBench(command string, r Resolver) bool` reports whether any simple command, or a wrapped string one level deep, invokes Bench.

Give each new package a package comment. Pass the home and the root
explicitly. Call `t.Parallel()` in each eligible new worktree test, because
the package census grades that call.

## Acceptance

- [ ] `poolkey.Key` returns the key the pool path used before this ticket, for a canonical root.
- [ ] `poolkey.Key` returns the primary repository's key for a linked worktree's root, and `bench worktree-pool` inside a linked worktree names the primary pool.
- [ ] The `cksum` goldens pass in `internal/poolkey` and no longer live in `internal/worktree`.
- [ ] `census.Dir` returns a path under `<home>/census/` and never under `<home>/worktrees/`. (EC26)
- [ ] `benchguard.InvokesBench` reports true for `bash -c 'bench status'` and false for `ls`.
- [ ] The follow-on guard reports the verdicts it reported before the export.
