# FT189 worktree-enumeration hang — probe matrix

Probed 2026-08-14 on this host (git 2.43.0, Linux/WSL2). Method: scratch repo
with one registered worktree, one admin entry at a time replaced by a FIFO,
every probe under `timeout` (exit 124 = hung). Findings are git-version- and
host-dependent; re-run the script below after any git upgrade before trusting
this matrix.

## Repro script

```sh
set -e
d=$(mktemp -d) && cd "$d"
git init -q repo && cd repo && git commit -q --allow-empty -m init
git worktree add -q ../wt HEAD
rm .git/worktrees/wt/gitdir && mkfifo .git/worktrees/wt/gitdir
timeout 5 git worktree list --porcelain; echo "exit=$?"   # 124 = reproduced
```

## Which admin entries hang `git worktree list --porcelain`

| Entry replaced by FIFO | Result |
| --- | --- |
| `gitdir` | hang (exit 124) |
| `HEAD` | hang (exit 124) |
| `commondir` | hang (exit 124) |
| `locked` | hang (exit 124) |
| stray extra FIFO (unread name) | no hang — git opens only named entries |

The mechanism is a blocking `open(2)` for read on a FIFO with no writer, so any
admin entry git opens hangs it; the read set is git's, not Bench's, and can
change across git versions.

## Which git commands hang with a FIFO `gitdir` present

| Command | Result |
| --- | --- |
| `git worktree list --porcelain` | hang |
| `git worktree add` | hang |
| `git worktree lock` / `unlock` | hang |
| `git worktree prune` | hang |
| `git branch --list` | hang (decorates worktree checkouts) |
| `git status --porcelain` | clean |
| `git rev-parse HEAD` | clean |
| `git rev-parse --git-common-dir` | clean — a pre-scan can locate the admin dir safely |
| `git for-each-ref` | clean — `LocalBranches` is unaffected |

`git worktree remove` was not probed (blocked as a destructive operation here);
assume exposed — it reads the same admin entries as `lock`/`prune`.

## Linked-root and repair-state probes (added at spec review)

Re-run 2026-08-14, same host and git 2.43.0:

| Probe | Result |
| --- | --- |
| `git worktree list --porcelain` from the **linked** worktree's root (FIFO `gitdir`) | hang (exit 124) — the admin dir is shared, and the linked root's `.git` is a file |
| `git rev-parse --path-format=absolute --git-common-dir` from the linked root, FIFO present | clean (exit 0), returns the primary `.git` |
| `gitdir` removed from an `<id>` dir | clean (exit 0) |
| `<id>` dir present but empty | clean (exit 0) |
| prunable entry (worktree directory deleted, admin entry intact) | clean (exit 0) |
| `worktrees/` itself replaced by a FIFO | clean (exit 0) — git never opens it as a file |
| `logs/HEAD` under an `<id>` replaced by a FIFO | clean (exit 0) — not in the list read set |
| first-level FIFO directly under `worktrees/` | clean (exit 0) — git skips the non-directory entry |
| first-level symlink (`worktrees/<id>` → elsewhere) | clean (exit 0) — git follows it |
| second-level FIFO at a name outside the read set | clean (exit 0) — matches the original stray-FIFO row |
| `worktrees/` itself a symlink to a healthy admin dir | clean (exit 0) — git follows it |
| `worktrees/` itself a symlink, FIFO `gitdir` behind it | hang (exit 124) — the symlink is a fail-open shape, not an absence |
| `gitdir` a symlink to a regular target | clean (exit 0) — git follows it |
| `gitdir` a symlink to a FIFO target | hang (exit 124) — a stat-based level-2 classifier admits this shape |

## Bench exposure

- Enumeration has one owner: `git.Worktrees` (`internal/git/git.go:138`), the
  sole parser and runner of `worktree list --porcelain -z`. Eight production
  callers reach it: `internal/worktree/{subshell,worktree,classifier,resume,list}.go`,
  `internal/intent/intent.go`, `internal/harness/worktree.go`, and the
  intra-package `PruneLandedBranches` (missed by the original external-caller
  sweep, found at spec review). Discovery runs in
  `bench status`, `bench resume` (session-start hook), `bench worktree *`, and the
  dashboard, so the hang lands before any Bench guard in practically every session.
- Worktree-mutating call sites also exposed (per the command matrix above):
  `internal/worktree/{ownership,lifecycle,reauthorize,snapshot}.go`,
  `internal/gate/engine.go`.
- `git.Worktrees` runs through `Raw` with no deadline; the established bound
  seam is `internal/bounds/bounds.go` — the named-constant policy registry plus
  `bounds.Run` (timeout + process-group SIGKILL).

## What was not read

Upstream git source and upstream tracker were not consulted; the "git may fix
this" retirement trigger rests on the observed 2.43.0 behavior only, not on a
named upstream commit or thread.
