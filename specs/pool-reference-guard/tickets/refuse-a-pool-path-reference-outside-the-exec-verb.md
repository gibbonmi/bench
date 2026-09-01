# Refuse a pool-path reference outside the exec verb

Blocked by: none
Writes: internal/benchguard/benchguard.go, internal/benchguard/benchguard_test.go, cmd/bench/main.go, .bench/hooks/block-bench-follow-on.sh, internal/systemtest/bench_follow_on_test.go
Covers: none

## What to build

The `block-bench-follow-on` hook refuses a Bash command that reaches a Bench worktree
by a route other than `bench worktree exec`. Today the guard refuses only a simple
command whose head word is `cd` with a literal pool-path argument. A shell assignment
that holds the pool path passes, and a `git -C` into the pool path passes.

The guard gains two refusal shapes beside the `cd` shape, and all three share one
refusal line. The first shape is a shell assignment word whose value starts with the
pool prefix. It counts in any command position, including after `export` and `env`. The
second shape is a `git` command whose `-C`, `--git-dir`, or `--work-tree` argument starts
with the pool prefix. The attached form and the separate-word form both count.

A literal pool path as a file argument to a read command stays allowed. `cat`, `sed`,
`rg`, and `head` over a pool path pass. `bench worktree release` and `bench worktree
land` take the pool path by grammar and pass. A `cd` through an unexpanded variable
stays allowed, because the scan reads the command text only. A relative `cd` inside a
wrapper string of an exec child stays allowed.

The refusal line keeps its one-form shape. It names `bench worktree exec "<label>" --
<command>` as the one route, names the three refused shapes, and ends with the target
it read. The hook header's `denies:` line advertises the widened denial in the same
change, so the manifest and the enforcement move together.

Write the regression rows before the predicate. They join the existing pool-target
table in the guard's unit test, one row per shape above, with the allowed rows beside
the refused rows. The exact-message test pins the new refusal line.

## Acceptance

- [ ] The hook exits 2 with the refusal line on `W=<pool>; git -C $W log`.
- [ ] The hook exits 2 on `git -C <pool> log`, on `git --git-dir=<pool>/.git log`, and on `git --work-tree <pool> status`.
- [ ] The hook exits 2 on `export W=<pool>` and on `env W=<pool> git log`.
- [ ] The hook exits 0 on `cat <pool>/specs/x.md` and on `sed -n 1,5p <pool>/f`.
- [ ] The hook exits 0 on `bench worktree release --request r <pool>` and on `bench worktree land --request r --base a --source-tip b -m m <pool>`.
- [ ] The hook exits 0 on `W=/tmp/x; git -C /tmp/x log` and on `cd "$W"`.
- [ ] The hook exits 2 on a bare `cd <pool>`, as before.
- [ ] The hook header `denies:` line names the assignment and the `git -C` shapes.
- [ ] `bench guards --brief` shows the widened denial.
