# Add the worktree show verb

Blocked by: add-child-argv-and-repeatable-flag-attributes.md, refuse-a-missing-tree-and-name-the-next-verb.md
Writes: internal/worktree/show.go (new), internal/worktree/show_test.go (new), internal/worktree/identifier_operand_test.go, internal/usage/worktree.go, cmd/bench/main.go, cmd/bench/main_test.go, cmd/bench/command_registry_test.go

## What to build

`bench worktree show <target> <rev>:<path>` prints a blob at a worktree
revision. The verb resolves the target through the shared resolver, and it
refuses through the shared refusal printer. It then runs
`git cat-file blob <rev>:<path>` in the worktree, and it passes stdout, stderr,
and the exit code through unchanged. A NUL byte in the blob arrives unchanged,
so the verb writes bytes and not lines.

The operand must hold a `:` and must not start with `-`. Either refusal returns
the grammar line at exit 2 and runs no Git. Git's own error and code 128 reach
the caller when the object is missing.

`show` joins the worktree command list, the `bench help` inventory, and the
kept routes. The exec row in both inventories gains `[--env KEY=VALUE]...`, so
the help and the parser agree. That row delivers story 41. The guard ticket
carries stories 42 and 43.

Package invariant for `internal/worktree`: one resolver and one refusal printer
serve `path`, `exec`, and `show`, and the `next=` line has one producer per
reason. Package invariant for `internal/usage`: every grammar declares its
attributes in `internal/usage/worktree.go`, and a grammar without the
child-argv attribute keeps the empty-positional refusal. Package invariant for
`cmd/bench`: the help inventory fixture and the kept-routes list are exact, so a
new verb names itself in both, and the registry row is their one source.

## Acceptance

- [ ] S1: `bench worktree show <target> HEAD:tracked.txt` prints the committed
      bytes of `tracked.txt` on stdout at exit 0.
- [ ] S2: a blob that holds `a\x00b\n` prints exactly those four bytes.
- [ ] S3: `HEAD:no-such-file` exits with Git's code 128 and Git's own stderr,
      and stdout is empty.
- [ ] S4: the operand `tracked.txt` returns the grammar line at exit 2 and runs
      no Git.
- [ ] S5: the operand `--output=/tmp/x:tracked.txt` returns the grammar line at
      exit 2 and runs no Git.
- [ ] S6: `bench worktree show no-such-label HEAD:x` prints
      `bench worktree show: target is unassigned` and then
      `next=bench worktree list`.
- [ ] S7: `bench help` prints the `show` row and the exec row with
      `[--env KEY=VALUE]...`, byte-equal to the fixture.
- [ ] S8: `bench worktree --help` names
      `bench worktree show <target> <rev>:<path>`.
