## Standards

Finding count: 0. Worst: none.

The frozen candidate has no Standards findings.

## Spec

Finding count: 2. Worst: the scanner accepts a non-regular `bench-lease` despite
the approved direct-entry predicate.

- `auto-fix` — `specs/worktree-enumeration-hang/spec.md:61-65` requires every
  direct admin entry to be regular or a directory. The name-based `bench-lease`
  exemption in `internal/git/git.go:233-237` accepts a symlink and conflicts
  with that contract. Remove the exemption and retain the typed scanner refusal.
- `auto-fix` — `specs/worktree-enumeration-hang/spec.md:84-100` keeps
  `CommonDir` on the plain `git.Output` runner and limits validation to
  `Worktrees`. Move common-dir shape validation out of `CommonDir` and into
  `Worktrees`.

## Coverage

Finding count: 4. Worst: the undocumented `bench-lease` exception.

- `auto-fix` — The non-symlink, non-directory `worktrees/` root-as-absence
  contract at `specs/worktree-enumeration-hang/spec.md:65-69` lacks a regular
  file probe; add it.
- `auto-fix` — The direct-entry regular-file-or-directory predicate at
  `specs/worktree-enumeration-hang/spec.md:61-65` lacks a regular first-level
  `worktrees/<id>` probe; add it.
- `auto-fix` — The hostile linked-root producer family is not combined with a
  malformed shared admin entry; add an enumeration path containing hostile
  shell characters and prove the scanner refuses it safely.
- `auto-fix` — The `bench-lease` symlink/FIFO allowance is unreviewed and is
  covered by the Spec repair above; remove the exemption and update the
  conflicting lifecycle expectation to the scanner's typed refusal behavior.
