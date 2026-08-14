# Enumerate specs without treating roots as glob patterns

Blocked by: none
Writes: internal/spec, internal/preflight

## What to build

Replace live folder-spec enumeration that interpolates the repository root into
a filesystem glob with literal directory traversal owned by `internal/spec`.
Repository and worktree path bytes must never become pattern syntax. Preserve
the current behavior for directory symlinks by classifying candidate entries
through `os.Stat`, and replace preflight's obsolete glob-specific diagnostic
with the literal folder-spec enumeration term.

## Acceptance

- [ ] `Facts` resolves staged folder specs when the repository root contains literal `[` and `]` characters.
- [ ] Enumeration retains deterministic slug order, skips a non-directory entry under `specs/` and a directory with no `spec.md`, preserves directory-symlink traversal, and still retains unreadable or special-file `spec.md` candidates as empty-metadata evidence rows.
- [ ] Restoring glob-based enumeration makes the focused bracketed-root `Facts` test red.
- [ ] Preflight's unresolved-status detail names folder-spec enumeration rather than a removed glob mechanism.
