# Review: session resume cleanup

Reviewed commit `d1c96dfbd4c4a06070400c415a5d0977f671bfbf`.

## Standards

**3 findings (2 hard violations, 1 judgment call). Worst:** repository-wide worktree and branch facts are parsed in multiple places, so cleanup, status, and intent lifecycle can drift.

### Hard violations

1. `internal/git/git.go:103` and `internal/intent/intent.go:323` independently parse `git worktree list --porcelain` alongside the typed parser in `internal/worktree/classifier.go:33`. Local branches are likewise enumerated separately in `internal/git/git.go:125`, `internal/intent/intent.go:332`, and `internal/worktree/branches.go:14`. This violates AGENTS.md's **one source per fact** rule against duplicated parsers and count derivations.
2. `bin/bench.sh:270` advertises `bench resume-clean` in general CLI help, while `.bench/BENCH-reference.md:103` classifies it as plumbing and `.bench/BENCH.md` says plumbing commands are driven by hooks and adapters, never typed by sessions.

### Judgment call

1. **Middle Man:** `internal/worktree/branches.go:40` and `internal/worktree/branches.go:55` retain exported one-line forwarders to `git.ResolvedDefault` and `git.LandedInDefault`; the pinned tree shows no external seam that needs this forwarding layer.

## Spec

**3 findings. Worst:** status can falsely report clean when the default branch has unpushed commits.

1. **Default-branch commits are falsely reported pushed.** The spec requires unpushed commit IDs ahead of **each local branch's** configured upstream (`specs/session-resume-cleanup.md:180`; coverage row 15). `internal/git/git.go:131-148` skips `branch == def` before querying its upstream. Reproduction: `main` tracking `origin/main`, one commit ahead, no other signal; `bench status --all` prints `bench: clean — nothing pending` instead of one unpushed commit.
2. **Replayed Claude intent changes bytes and creation time.** Coverage row 2 requires same-key upserts to be byte-idempotent, and the spec makes `tool_use_id` Claude's stable key (`specs/session-resume-cleanup.md:91-95`). `cmd/bench/main.go:205-206` creates a fresh timestamp before assigning that key, and `internal/intent/intent.go:157-167` replaces the existing entry, so replaying the same allowed envelope mutates the ledger.
3. **Cleanup classification failures can return success.** The command contract says classification failures exit 1 (`specs/session-resume-cleanup.md:141-144`). `internal/git/git.go:231-242` collapses failed `rev-list` or `cherry` queries into ordinary “not landed,” and `internal/worktree/resume.go:75-78` discards the distinction and exits 0 after keeping the branch.

## Coverage

**2 findings. Worst:** ignored-only WIP can be irreversibly deleted during SessionStart.

1. **Ignored-only WIP in an out-of-pool worktree is deleted.** `internal/worktree/resume.go:47-63` treats an empty `git status --porcelain` result as clean; ignored files are absent from that output, and `git worktree remove` deletes them. No coverage row inventories ignored-only state, while `internal/worktree/worktree_test.go:214` covers ordinary untracked WIP only. Reproduction: create a registered detached worktree containing only a `.gitignore`-matched file; `bench resume-clean` reports one cleaned worktree and removes the file.
2. **A registered worktree path containing a newline is misparsed.** `internal/worktree/classifier.go:53-56` and `internal/git/git.go:108-112` split non-NUL `git worktree list --porcelain` output on newlines. A newline-bearing path is truncated, so cleanup silently skips it and status reports git state unavailable. No map row or test covers this path even though `projects/benchkit.md` requires control-byte coverage for git-sourced paths.
