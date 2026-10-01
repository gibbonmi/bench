# Prove that unlink removes only clean Bench-owned assets

Blocked by: none
Writes: internal/adopt/link_transaction_test.go, internal/adopt/adopt_test.go
Covers: none

## What to build

The quality survey of 2026-09-29, card 07, found that `bench unlink` is a destructive path with almost no test. The only unlink test in `internal/adopt` checks the refusal for an unresolved hooks directory. The package has 63.7% statement coverage.

The project profile gives the promise under test. Unlink removes only clean manifest-owned assets. Modified assets and project-owned collisions stay in place. Unlink reports each residual with the same machine-readable partial result as link.

The `Unlink` contract adds three more promises. With no manifest, unlink exits 1. A dry run writes nothing. A manifest row that escapes the repository is refused and never removed.

Add tests through `Unlink`, with real temporary repositories that a real `Link` run prepares, for these uncovered cases:

- A clean link and unlink removes each manifest row, the manifest, and the managed hook, and exits 0. A project file outside the manifest and the project prose in `AGENTS.md` stay.
- A modified managed asset stays with its bytes, and a project-owned collision at a managed path stays with its bytes. Unlink exits 3, removes the other clean rows, keeps the manifest, and prints the `residuals` table with the modified row.
- A dry run over the same repository exits 3 and removes nothing.
- A manifest row that names a path outside the repository is refused, and the file at that path stays. Unlink exits 3.
- With no manifest, unlink exits 1 and names the missing manifest.

Put the tests in `link_transaction_test.go`, because the `internal/adopt/` directory has no file-count headroom and `adopt_test.go` is above its line budget. Move the consumer-repository setup into one helper there, and make the relink test in `adopt_test.go` use it. Change no production code.

## Acceptance

- [ ] Each case above has a test that passes on the current tree.
- [ ] A `bench probe` on the keep-or-remove decision in `unlink.go` turns at least two of the tests red, and the restored file turns them green.
- [ ] The relink test and the new tests share one consumer-repository helper.
- [ ] `go vet ./...` and `bench test --changed` pass, and `internal/adopt` statement coverage goes up.
