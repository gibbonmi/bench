# Cache the full-id reads and the administration path on the run

Blocked by: none
Writes: docs/adr/0033-a-command-run-caches-its-pure-git-queries.md (new), internal/git/run.go, internal/git/run_cache_test.go, internal/git/run_object_test.go (new), internal/git/tree.go, internal/git/worktree_admin.go, internal/gittest/gittest.go, internal/worktree/land_identity.go, internal/worktree/parallel_census_test.go, internal/worktree/reauthorize_test.go, internal/worktree/worktree.go
Covers: GQ1, GQ2, GQ3, GQ4, GQ5, GQ6, GQ7, GQ8, GQ9, GQ10, GQ13, GQ14, GQ15, GQ16, GQ17, GQ18, GQ19, GQ20, GQ21, GQ22, GQ23, GQ24, GQ25, GQ26, GQ43, GQ44, GQ45, GQ47

## What to build

Give `git.Run` the cached forms of the spec's form table. Each package function stays the live form, and it starts git on every call.
Add `git.FullObjectID`. It accepts 40 or 64 hexadecimal digits of either case. The cache and the worktree identity checks use only this predicate.

Key each object answer by its exact argv, with `-C` and the root spelling included. Return a copy of each stored answer, and store no failed answer.
A yes and a no of `IsAncestor` both store. An error stores nothing. A ref, HEAD, or short-id operand reads live.

Make each run form safe on a nil `*Run`. A nil run starts the query on every call and stores nothing. Tickets 02 and 03 rely on this contract.
`AdminPath` caches by root spelling and name. Each read checks the parent directory of the answer again, as the location answers do now.
`PrimaryCheckout` reads the common directory through the run, and it starts the porcelain on every call.
`ReadTreeFile` composes the run's `ListTree` and `ReadControlBlob` in sequence, so no form locks the run's mutex twice.

Split the record parser of `ListTree` in `internal/git/tree.go`, so that the run form and the package form share it. Keep each error text byte for byte.
Three tests quote the `ReadTreeFile` refusal text: `internal/commitment/head_binding_test.go:236`, `internal/preflight/charge_test.go:363`, and `internal/gate/checkpoint_prepare_test.go:148`.

Add a `log` mode to the stub in `internal/gittest/gittest.go`. It records each start and passes every query to the real git. Ticket 03 uses it.
The `fail-git-path` mode refuses each `--git-path` query (`internal/gittest/gittest.go:91-96`), so the `AdminPath` rows use the `log` mode.

Replace the 40-digit rule of `fullCommitIdentity` (`internal/worktree/worktree.go:572`) with a call to `git.FullObjectID`.
Then `abbreviatedIdentity` (`internal/worktree/land_identity.go:290-292`) also expands a short prefix to a 64-digit id. The 39-digit prefix bound stays.
Leave `hexIdentity`, `lowercaseHex`, and `fullCommitID` unchanged, because each grades a different fact.

Correct the cache-scope comments at `internal/git/run.go:5-14` and `:52-54`, and at `internal/git/worktree_admin.go:191-194`.
Write ADR 0033. It states the pure scope, the live set, and the one-run rule. It names no file path and no code.

Put the new `internal/git` tests in `run_object_test.go`. Use `newRepo` and the `gittest` stub, and start no command of your own.
The cap at `internal/conformance/ordinary_build_census_test.go:70-84` allows one repository site in the package tests, and `testhelpers_test.go:18` holds it.
`run.go` stays under its 400-line budget. `worktree_admin.go` has 49 lines of headroom.
Raise `worktreeTestCount` from 877 to 878 for the new reauthorize test. `parallel_census_test.go` is over its budget, so change the pin line in place and add no line.

## Acceptance

- [ ] Two `ReadControlBlob` reads of one full blob id on one run start `cat-file` once (GQ1).
- [ ] Two `ReadControlBlob` reads of `HEAD:<path>`, with a commit between them, return the new bytes (GQ2).
- [ ] Two `ListTree` reads of one full tree id with one path start `ls-tree` once (GQ3).
- [ ] Two `ListTree` reads of `HEAD`, with a commit between them, list the new tree (GQ4).
- [ ] Two `ListTree` reads of a seven-digit tree prefix start `ls-tree` twice (GQ5).
- [ ] Two `ListTree` reads of an absent path at one full tree id return an empty listing and start `ls-tree` once (GQ6).
- [ ] Two `ReadTreeFile` reads of one path at one full tree id start `ls-tree` once and `cat-file` once (GQ7).
- [ ] Two `ResolveCommit` calls with one full commit id and one flag set start `rev-parse` once (GQ8).
- [ ] Two `ResolveCommit` calls with a branch name, with an `update-ref` between them, return the moved commit (GQ9).
- [ ] A `ResolveCommit` call with `--quiet` and one without it, for one full id, start `rev-parse` twice (GQ10).
- [ ] Two `IsAncestor` checks of one ancestor pair of full ids return yes and start `merge-base` once (GQ13).
- [ ] Two `IsAncestor` checks of one non-ancestor pair of full ids return no and start `merge-base` once (GQ14).
- [ ] Two `IsAncestor` checks with a branch operand, with an `update-ref` between them, return the changed answer (GQ15).
- [ ] Two `IsAncestor` checks that name a missing full id return an error and start `merge-base` twice (GQ16).
- [ ] Two `ReadControlBlob` reads of a missing full id return an error and start `cat-file` twice (GQ17).
- [ ] A change to the bytes of a first `ReadControlBlob` answer leaves the second answer equal to the blob (GQ18).
- [ ] Two calls of each cached form on a nil run start the query twice (GQ19).
- [ ] Two package-level `ReadControlBlob` reads of one full id start `cat-file` twice (GQ20).
- [ ] Two run values each start one `cat-file` for one full id (GQ21).
- [ ] Two `AdminPath` reads of one name on one run start `--git-path` once (GQ22).
- [ ] After the parent directory of an `AdminPath` answer is removed, the second read refuses and starts no query (GQ23).
- [ ] Two `PrimaryCheckout` calls on one run start the common-directory query once and the porcelain twice (GQ24).
- [ ] `FullObjectID` accepts 40 and 64 hexadecimal digits of either case (GQ25).
- [ ] `FullObjectID` refuses 39, 41, 63, and 65 digits, a value with one non-hex byte, and the empty value (GQ47).
- [ ] The reauthorize hint quotes a 64-digit full id instead of the placeholder (GQ26).
- [ ] Inside one run, an `update-ref` of the default branch is seen by `ResolvedDefault` and by the run's `ResolveCommit` of the branch (GQ45).
- [ ] The named `internal/git` comments state the run's current cache scope (GQ43, review-owned).
- [ ] ADR 0033 states the pure scope, the live set, and the one-run rule, with no file path and no code (GQ44, review-owned).
- [ ] `worktreeTestCount` reads 878.

## Checkpoint verification

- `bench test --package ./internal/git`
- `bench test --package ./internal/gittest`
- `bench test --package ./internal/commitment`
- `bench test --package ./internal/preflight`
- `bench test --package ./internal/gate`
- `bench test --package ./internal/worktree`
- `bench test --package ./internal/conformance`

The lane grades the ADR prose at the ticket commit. Record one red for each independent start count before the green, with the named omission that turns it red.
