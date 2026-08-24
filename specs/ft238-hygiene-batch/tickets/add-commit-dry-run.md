# Add a dry run to bench commit

Blocked by: add-commit-help-example.md
Writes: internal/commit/commit.go, internal/landing/landing.go, tests beside each

## What to build

`bench commit --dry-run -m <msg> -- <path>...` composes the named-path
snapshot, authorizes it through the gate, and reports the verdict. It never
creates a commit and never moves a ref. The caller can then diagnose which
phase a composed set reds on without a junk commit.

## Acceptance

- [ ] A green dry run exits 0, reports the verdict, and `HEAD` does not move.
- [ ] A red dry run exits non-zero with the gate's failure attribution, and `HEAD` does not move.
- [ ] The help text advertises `--dry-run`.
