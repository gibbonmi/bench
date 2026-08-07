# Repair adapter symlink convergence

Blocked by: none
Ownership fence: `internal/adopt/link_transaction.go`, `internal/contract/surface/link_lifecycle_test.go`
Integration surfaces: convergence predicate→`internal/adopt/link_transaction.go` + ASC1; link lifecycle contract→`internal/contract/surface/link_lifecycle_test.go` + ASC1
Contracts: none crosses
Closure: ASC1/adapter-symlink-converged, ASC2/adapter-symlink-drifted

## What to build

Make `convergedFingerprint` symlink-aware. An adapter entry stages a symlink to
its canonical `.agents` file, so when a repo satisfies adapters with a
directory-level symlink (`.claude/commands -> ../.agents/commands`), the
destination is a regular file whose content is byte-identical to the staged
link's target — converged in every observable sense — yet the `Lstat`
type-and-permission comparison can never match a symlink against a file, so
every adapter reads "needs a write" and the symlink-parent rule hard-aborts
`bench link` on any such repo, including this one. Converged means the
destination resolves to the same effective content the staged entry would
produce: compare a staged symlink by its resolved target content (and skip the
permission comparison that is meaningless for links), while keeping the
existing byte-and-mode comparison for staged regular files. A genuinely drifted
adapter destination under a symlink parent must still hard-abort.

## Acceptance

- [ ] [ASC1] (covers local) `bench link` completes with exit 0 and no destination writes on a fixture whose adapter directory is a symlink to the canonical directory and every entry is content-identical.
- [ ] [ASC2] (covers local) The same fixture with one adapter's canonical content changed still hard-aborts exit 1 with the symlink-parent conflict.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| ASC1/adapter-symlink-converged | compare a staged symlink with `Lstat` type-and-perm equality instead of resolved content | link lifecycle contract | link the symlinked-adapter fixture, expect the exit-0 assertion red on the symlink-parent abort |
| ASC2/adapter-symlink-drifted | treat any resolvable staged symlink as converged regardless of content | link lifecycle contract | change one canonical file, link, expect the hard-abort assertion red on exit 0 |
