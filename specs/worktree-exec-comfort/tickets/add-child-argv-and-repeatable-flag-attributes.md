# Add the child-argv and repeatable-flag grammar attributes

Blocked by: none
Writes: internal/usage/parse.go, internal/usage/parse_test.go, internal/usage/worktree.go

## What to build

The parser gains two grammar attributes, and the exec grammar declares both.
The first attribute declares that every token after `--` is child argv.
`Parse` then passes those tokens through unchanged, and the empty token `''`
reaches the child argv. A grammar without that attribute keeps today's
empty-positional refusal after `--`.

The second attribute marks one flag as repeatable. `Parse` collects a
repeatable flag's values in argv order. A repeated flag without that attribute
stays a usage error that names the flag.

The exec grammar becomes
`bench worktree exec <target> [--env KEY=VALUE]... -- <command> [args...]`.
`--env` sits after the reserved target and before `--`. The flag pass stops at
`--`, so a `--env` token after `--` belongs to the child.

This ticket delivers the parser rows only. The ticket
`pass-argv-stdin-and-env-to-the-exec-child.md` reads the parsed target, the
`--env` values, and the child argv from the `Result` this ticket produces.

## Acceptance

- [ ] X1: `bench worktree exec <target> -- rg -N '' README.md` parses with the
      child argv `rg`, `-N`, `''`, `README.md` and exit 0 from the parser.
- [ ] X2: `bench commit -m x -- ''` still returns the usage line that names
      `""` at exit 2.
- [ ] X8: `--env A=1 --env B=2` parses to both values in argv order.
- [ ] X9: `bench commit -m a -m b` still returns the usage line that names
      `-m` at exit 2.
- [ ] X12: `bench worktree exec <target> -- printf '%s' --env` parses with
      `--env` as the last child argument.
