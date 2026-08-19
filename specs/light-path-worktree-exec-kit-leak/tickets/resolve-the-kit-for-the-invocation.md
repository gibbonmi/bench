# Resolve the kit for the invocation

Blocked by: stop-leaking-wrapper-routing-into-a-worktree-child.md
Writes: bin/bench.sh, cmd/bench/main_test.go

## What to build

`kit_dir` resolves the kit as the tree the wrapper file lives in. An operator
who runs the main checkout's `bench` from `PATH` with a worktree of that same
kit as the current directory is working on the worktree, but the wrapper names
the main checkout. The gate then grades the worktree against a kit that is not
it, drops the phases it only runs over its own kit, and reports red over a
clean tree.

Resolve the kit for the invocation instead: when the current directory sits in
a git worktree of the same repository as the wrapper's own tree, that
worktree's top level is the kit. Every other case keeps today's answer, and
every git failure — no repository, an unrelated repository, no `git` at all —
falls back to it rather than erroring.

Two conditions guard the exception, and a linked project repo fails both. The
wrapper's kit must be its own tree's top level, because an adopted repo carries
the kit in a subdirectory of the project. Sameness is identity of the git
common directory, never a path prefix, because two unrelated repositories can
sit under one parent and two working trees of one repository can sit anywhere.

The binary lookup needs no change: `main_tree_kit` already re-anchors a worktree
kit at the main tree, which is where the untracked dev binary lives.

## Acceptance

- [ ] The wrapper invoked with a kit worktree as the current directory resolves the kit as that worktree.
- [ ] The wrapper invoked with a different repository as the current directory resolves the kit as its own tree.
- [ ] The wrapper invoked from an adopted repo, whose kit is a subdirectory, resolves the kit as that subdirectory.
- [ ] The wrapper invoked outside any repository, or with no `git` available, resolves the kit as its own tree.
- [ ] An explicitly set `BENCH_KIT` still wins.
- [ ] Run from a kit worktree, `bash bin/bench.sh gate --fresh` reports green with the race and system phases present and no unstaged environment skip.
