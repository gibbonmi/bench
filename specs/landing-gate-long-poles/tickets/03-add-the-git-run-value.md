# Add the git run value

Blocked by: none
Writes: internal/git/run.go (new), internal/git/run_cache_test.go (new), internal/git/worktree_admin.go, internal/git/git.go, internal/poolkey/poolkey.go, internal/poolkey/run_test.go (new), internal/gittest/gittest.go, internal/gittest/gittest_test.go
Covers: LGP9, LGP10, LGP11, LGP12, LGP13, LGP14, LGP15, LGP34

## What to build

`internal/git` gains a run value that caches the stable location answers of one command run. A command creates it at its entry, and later tickets pass it down.

The run value caches the raw answers of the common-directory, top-level, and admin-directory queries. The key is the root spelling that the caller passes. On each read it runs the existing validation again, such as `validateCommonDir` (`internal/git/worktree_admin.go:245-259`). The cache saves the git start, not the filesystem check.

The run value has its own `Worktrees` form. That form reads the common directory through the run value, and it keeps the admin scan and the porcelain start live.

The run value holds no ref, HEAD, registration, index, default-branch, or ancestry answer. No package variable holds a cache.

The package-level `CommonDir`, `AdminDir`, `Worktrees`, `Root`, and `RootAt` keep their signatures and stay uncached. About 100 callers depend on them, and the hostile tests in `internal/git/worktree_admin_hostile_test.go` call them directly.

Put the run value in a new `internal/git/run.go`, because `worktree_admin.go` has 348 lines. Spell the four plumbing flags only inside `internal/git`, so the git plumbing owner check stays green.

`internal/poolkey` gains a form of `Canonical` that reads the common directory through the run value. `Canonical` keeps its signature for its other callers. Put the poolkey test in a new `internal/poolkey/run_test.go`, because a canary pins a line of `poolkey_test.go`.

Count git starts with the argument log of `gittest.StubGitDir` (`internal/gittest/gittest.go:37-138`). `internal/gittest` does not import `internal/git`, so the import adds no cycle. Count only the lines of the query under test.

## Acceptance

- [ ] Two common-directory reads in one run start git once (LGP9).
- [ ] Two top-level reads and two admin-directory reads in one run start git once each (LGP10).
- [ ] A common directory replaced by a symlink after the first read refuses on the second read (LGP11).
- [ ] A branch moved by `update-ref` inside one run reads its new commit (LGP12).
- [ ] `TestWorktreesRejectsBadCommonDirBeforePorcelain` passes unchanged (LGP13).
- [ ] Two run values over one root each start their own common-directory query (LGP14).
- [ ] `TestGitPlumbingOwnerRedsARetypedFlag` passes, and the git plumbing owner check stays green (LGP15).
- [ ] The probe "skip the validation on a cached common directory" bites.
- [ ] The verification record holds the `internal/git` package time before and after the change (LGP34).
