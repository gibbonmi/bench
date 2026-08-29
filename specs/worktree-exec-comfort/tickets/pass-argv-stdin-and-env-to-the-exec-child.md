# Pass argv, stdin, and env to the exec child

Blocked by: add-child-argv-and-repeatable-flag-attributes.md
Writes: internal/worktree/exec.go, internal/worktree/exec_test.go

## What to build

The exec command reads the parsed `--env` values and applies them to the child
environment. The command applies the caller's environment first, then the
`--env` values, then the verb's own routing variables. `BENCH_HOME`,
`BENCH_WRAPPER`, `BENCH_KIT`, and `BENCH_RUN_BINARY` therefore keep the verb's
values, and an `--env` name among them cannot repoint the child's pool.

A `KEY` matches `[A-Za-z_][A-Za-z0-9_]*`. A value with no `=`, an empty `KEY`,
or a bad `KEY` is a usage error at exit 2 that names the value. Such a refusal
starts no child. An empty `VALUE` sets the name to the empty string.

`runWorktreeChild` gives the caller's stdin to the child, so a heredoc reaches
the child byte for byte. The exec help text carries three lines: the grammar
line, a stdin example that feeds `-- python3 -` from `<<'EOF'`, and the exit-2
rule line. The rule line names `usage: bench worktree exec` as the prefix of a
grammar refusal, so a reader tells that refusal from a child's own exit 2.

The ticket `print-the-worktree-path-on-exec-failures.md` edits the same two
files after this ticket lands, and it owns the child's failure surfaces.

## Acceptance

- [ ] X3: `bench worktree exec --help` prints the grammar line, a line that
      shows `<<'EOF'` feeding `-- python3 -`, and the exit-2 rule line.
- [ ] X4: `runWorktreeChild` gives `-- cat` a reader of three lines that end
      in a NUL byte, and the child emits those bytes unchanged.
- [ ] X6: the exec help's third line names `usage: bench worktree exec` as the
      prefix of an exit-2 grammar refusal.
- [ ] X7: `bench worktree exec <target> --env FOO=bar -- sh -c 'echo $FOO'`
      prints `bar`, and the caller's process has no `FOO`.
- [ ] X10: `--env FOO` and `--env 1X=y` each return the usage line that names
      the value at exit 2, and each starts no child.
- [ ] X11: `--env BENCH_HOME=/x -- sh -c 'echo $BENCH_HOME'` prints the verb's
      resolved home.
- [ ] X13: `--env BENCH_KIT=/x -- sh -c 'echo ${BENCH_KIT-unset}'` prints `unset`.
