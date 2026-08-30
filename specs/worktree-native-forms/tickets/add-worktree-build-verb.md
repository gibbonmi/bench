# Add the worktree build verb

Blocked by: none
Writes: internal/worktree/build.go (new), internal/worktree/build_test.go (new), internal/worktree/land.go, internal/worktree/exec_test.go, internal/worktree/identifier_operand_test.go, internal/usage/worktree.go, cmd/bench/main.go, cmd/bench/main_test.go, cmd/bench/command_registry_test.go

## What to build

An agent runs `bench worktree build <target>` and gets an executable at
`<worktree>/dist/bench`. The verb resolves the target through `resolveWorktree`.
It refuses an unresolved target through `printTargetRefusal` with the verb name
`bench worktree build`. It builds through a new `build` join on the worktree
`joins` set, whose default is `runbinary.Build`, under `subprocess.NotifyCancel`.

Success prints one `worktree_build[1]{worktree,executable}` table and then one
`next[1]:` line. The next line shell-quotes a line-safe label and carries the
assignment id for a label with a control byte. A build error prints the builder's
sentence and then the `worktree:` line through `nameWorktree`.

This ticket crosses three shared surfaces. It adds `usage.WorktreeBuild` to the
`worktreeCommands` list in `internal/usage/worktree.go`. It adds
`bench worktree build` to the `keptWorktreeGrammars` list in
`cmd/bench/command_registry_test.go`. It adds the `build` row to the exact
`bench help` fixture in `cmd/bench/main_test.go`, because both rows land
together. The create ticket and the gate-form ticket edit the same files, so
the coordinator serializes the three in the order build, create, gate form.

## Acceptance

- [ ] WF1: the verb calls the `build` join once with the canonical worktree path and `<worktree>/dist/bench`.
- [ ] WF2: the production join, against a worktree whose `scripts/go-build.sh` stub writes a marker, leaves that marker byte-equal at `dist/bench` at exit 0.
- [ ] WF3: success prints `worktree_build[1]{worktree,executable}:` with the assignment id and the absolute executable path.
- [ ] WF4: the label `it's a*b` yields the next line `bench worktree exec 'it'\''s a*b' -- ./dist/bench <verb>`, and a newline label yields the id.
- [ ] WF5: a `build` join that returns an error prints `bench worktree build: <error>` and then the `worktree:` line at exit 1.
- [ ] WF6: a `PATH` that holds no `go` refuses at exit 1 with `Go is absent from PATH` and the `worktree:` line.
- [ ] WF7: a `build` join that returns `context.Canceled` exits 130, and stderr ends with the `worktree:` line.
- [ ] WF8: `bench worktree build no-such-label` prints `target is unassigned` and then `next=bench worktree list` at exit 1.
- [ ] WF9: two builds in a row both exit 0, and `dist/bench` holds the second stub's bytes.
- [ ] WF10: `runWorktreeChild` with the argv `./dist/bench version` prints the stub's word on stdout at exit 0.
- [ ] WF11: after the production-join build every untracked path in the worktree sits under `dist/`.
- [ ] WF12: `bench worktree --help` names `bench worktree build <target>`.
- [ ] WF13: `bench help` prints the `build` row, byte-equal to the fixture.
