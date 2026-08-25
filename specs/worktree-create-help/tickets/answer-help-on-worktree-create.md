# Answer --help on bench worktree create

Blocked by: none
Writes: internal/worktree/worktree.go, internal/worktree/worktree_test.go, cmd/bench/command_registry_test.go

## What to build

`bench worktree create --help` prints the create grammar on stdout and exits
zero. Today the verb hand-parses its flags, so a help spelling exits 2 with
the usage line on stderr. Its sibling `bench worktree reauthorize` parses
through `usage.Parse`, which owns the help spellings, the repeated-flag rule,
and the empty-value rule. Move create onto the same seam. Declare a
`createGrammar` with `--request` and `--label` as required non-empty value
flags and `--refresh` as a boolean flag. Parse through `usage.Parse`.

Keep `refreshop.Consume` as the one owner of the refresh fetch. Parse first,
then hand the original args to `refreshop.Consume` only when the parse
succeeded, so a help request does not fetch. The success output and the
`next[2]` block do not change.

## Acceptance

- [ ] `bench worktree create --help`, `-h`, and `help` print `usage: ` plus `usage.WorktreeCreate` on stdout and exit 0.
- [ ] `bench worktree create --request x --help` exits 0 with the same help and performs no refresh.
- [ ] `bench worktree create --request x` (no `--label`) exits 2 with the usage line, and `--request ""` exits 2 naming the empty value.
- [ ] The `command_registry_test.go` help table lists `worktree create --help`, so the registry sweep proves the route.
- [ ] `bench worktree create --request r --label l` still prints the `worktree_create` table and the `next[2]` block unchanged.
