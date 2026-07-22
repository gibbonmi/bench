## Standards

1 finding. Worst issue: duplicated durable staged-write policy.

- Hard violation — `stageBytes` in `internal/adopt/transaction.go:19` and
  `stageBeside` in `internal/adopt/link_transaction.go:271` independently own the
  same write → sync → close → cleanup lifecycle. Collapse it to one source per
  the project code standard in `AGENTS.md`.

## Spec

1 finding. Worst issue: story 2 durability is incomplete for symlinked assets.

- `specs/transactional-asset-lifecycle.md:23`, `:32`, and `:103` require staged
  writes to be synced before promotion. Adapter and symlink-mode paths at
  `internal/adopt/link_transaction.go:193-203` call `os.Symlink` directly and
  never sync the containing directory before `promoteAll` renames them.

## Coverage

2 findings. Worst issue: rollback restoration is untested during relink.

- `internal/contract/surface/link_lifecycle_test.go:23` covers fault rollback only
  from a fresh fixture. Add a kit A → kit B relink fault after replacements and a
  dropped clean asset are journaled; assert the exact pre-relink tree and manifest
  return.
- `projects/benchkit.md:104` and `specs/transactional-asset-lifecycle.md:127`
  include repo roots with spaces/glob characters, but
  `internal/contract/surface/link_test.go:232` varies only the kit path. Add a
  link/relink or rollback contract from a repo root such as `repo [1]`.
