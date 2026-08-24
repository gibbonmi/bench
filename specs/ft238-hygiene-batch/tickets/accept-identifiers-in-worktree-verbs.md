# Accept identifiers in the path-taking worktree verbs

Blocked by: emit-absolute-worktree-path.md
Writes: internal/worktree/path.go, internal/worktree/clean.go, internal/worktree/land.go, internal/worktree/reauthorize.go, internal/worktree/worktree.go, tests beside each

## What to build

Every path-taking worktree verb (`clean`, `land`, `reauthorize`, `release`)
also accepts the label or the assignment id. Each verb also accepts an
unambiguous prefix of either, with 8 to 12 characters. The shared resolver already accepts the label and the
exact id; this ticket adds prefix resolution and routes the four verbs through
it. An ambiguous prefix is a refusal that says so. `clean --apply` accepts a
fingerprint prefix under the same ambiguity rule.

## Acceptance

- [ ] Each of the four verbs resolves a label, an id, and a 8–12 character unique prefix to the same worktree as the full path.
- [ ] An ambiguous prefix, and a prefix shorter than 8 characters, are each refused with a named reason.
- [ ] `clean --apply <fingerprint-prefix>` applies the plan the full fingerprint would.
- [ ] A plain path operand keeps its current behavior in all four verbs.
