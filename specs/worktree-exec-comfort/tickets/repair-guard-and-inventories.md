# Repair the guard API and the inventories the review found

Blocked by: none (repair ticket; runs after the review at `1ccdb060`)
Writes: internal/benchguard/benchguard.go, internal/benchguard/benchguard_test.go, internal/shellcommand/shellcommand.go, internal/worktree/exec.go, cmd/bench/main.go, cmd/bench/command_registry_test.go

## What to build

`benchguard.Classify` becomes the verdict function: it takes the name `Judge`
holds now and returns the `Verdict`. The bool wrapper goes away, and every
caller reads `.Blocked`. The spec's seam map already names `Classify`.

`internal/shellcommand` exports one predicate for a shell assignment word, and
`internal/worktree/exec.go` reads the `--env` shape through it instead of a
second copy of the regex.

`TestKeptWorktreeOperationsKeepTheirGrammar` asserts the full
`bench worktree show <target> <rev>:<path>` grammar, not the prefix.

## Acceptance

- [ ] `rg -n 'func Judge' internal/benchguard` finds nothing, and `rg -n 'func Classify' internal/benchguard` returns a `Verdict`.
- [ ] G1–G14 stay green at the unit seam and the hook seam.
- [ ] `rg -n '\[A-Za-z_\]\[A-Za-z0-9_\]\*=' internal --glob '*.go'` finds one production site.
- [ ] X10 stays green.
- [ ] S8's test names the full `show` grammar and stays green.
