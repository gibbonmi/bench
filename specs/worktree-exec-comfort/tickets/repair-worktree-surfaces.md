# Repair the worktree surfaces the review found

Blocked by: none (repair ticket; runs after the review at `1ccdb060`)
Writes: internal/worktree/show.go, internal/worktree/show_test.go, internal/worktree/list.go, internal/worktree/list_actions_test.go, internal/worktree/identifier_operand_test.go, internal/worktree/identity_component_test.go, internal/worktree/exec_test.go, specs/worktree-exec-comfort/spec.md, CHANGELOG.md

## What to build

`ShowCommand` parses its two positionals through `usage.Parse` with a grammar
that reserves both slots. So the sole `--help` spelling prints the grammar
line at exit 0. A dash operand still refuses at exit 2 with the grammar line.
The operand also takes the `lineSafe` check, so a control byte refuses at exit
2 before any Git runs. The `runShowChild` comment drops the claim that an
interrupt reports as `exec` does. The spec's Edge inventory records the
`show` cancel code as Won't handle.

`missingTreeRecovery.line()` quotes the path with `axi.ShellQuote`. So the
`next=` line and the `list` help row render one way, and a `'` in the path
stays paste-safe. The spec's implementation decision and rows F9 and F11 say
"shell-quoted" instead of "quoted". The test for F9 gains a path with a `'`.

`identity_component_test.go` reads the `next=` default from `nextList`. The
X8 child test sets two `--env` values and reads both from the child.
`CHANGELOG.md` gains one `Unreleased` entry for `show`, `--env`, the `next=`
line, and the `worktree:` line, and the spec's fence gains `CHANGELOG.md`.

## Acceptance

- [ ] `bench worktree show --help` prints `usage: bench worktree show <target> <rev>:<path>` at exit 0.
- [ ] `show` with an operand that holds a control byte exits 2 and runs no Git.
- [ ] S4 and S5 stay green.
- [ ] F9 with a path that holds `'` prints a `next=` line whose path is `axi.ShellQuote`d.
- [ ] F11's help row and F9's `next=` line quote one path the same way.
- [ ] X8: `--env A=1 --env B=2 -- sh -c 'echo $A$B'` prints `12`.
- [ ] `identity_component_test.go` names no literal `next=` default.
- [ ] `CHANGELOG.md` holds the entry, and `bench preflight build worktree-exec-comfort` is green.
