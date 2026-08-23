# Make a published-but-unreconciled commit exit 3 and name its remainder

Blocked by: retire-the-commit-route-flip-and-close.md
Writes: internal/commit/commit.go, internal/commit/close_test.go (renamed by the blocker), internal/landing/landing.go, internal/landing/landing_test.go

## What to build

When the ref update succeeded and the reconcile failed, `bench commit` exits 3
and prints one stdout record in the landing verb's `name{key=value,...}`
grammar: `committed{published_commit=<sha>,path=<path>,next=<command>}`.
Every return of the reconcile step carries the failed path in a typed error,
so the command renders the path without text parsing. The `next=` value is one
`git restore --source=<sha> --staged --worktree --` followed by every named
path, shell-quoted, in the owner's sorted and deduplicated order.

`path=` renders through the control-byte sanitizer. When a named path is not
line-safe, `next=` replaces the path list with the placeholder
`<named-paths>`. A refusal before publication keeps exit 1 with its `error:`
prefix, and a grammar error keeps exit 2. The help text names the exit-3
meaning. FB9 drives the production reconcile to a failure with no injected
function, for example by a mutation the stub gate performs on the checkout.

## Acceptance

- [ ] FB1: a `bench commit` whose ref update succeeded and whose reconcile fails exits 3.
- [ ] FB2: that record's `published_commit=` equals the new HEAD.
- [ ] FB3: when the second of two named paths fails reconcile under an injected reconcile, `path=` names the second path.
- [ ] FB9: the production reconcile, with no injected function, failing on a named path yields `path=` that names it.
- [ ] FB4: that record's `next=` is `git restore --source=<sha> --staged --worktree --` followed by every named path, shell-quoted, in sorted order, with one path carrying a space.
- [ ] FB5: a red gate exits 1 and moves no ref.
- [ ] FB6: a missing `-m` exits 2.
- [ ] FB7: a failed path that carries an ESC byte renders `path=` as the sanitizer's spelling and `next=` with `<named-paths>`.
- [ ] FB8: `bench commit --help` names exit 3.
