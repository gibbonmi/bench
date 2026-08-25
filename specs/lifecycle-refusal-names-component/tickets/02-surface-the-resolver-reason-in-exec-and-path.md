# Surface the resolver reason in exec and path

Blocked by: 01-name-the-landing-identity-component.md
Writes: internal/worktree/ownership.go, internal/worktree/path.go, internal/worktree/exec.go, internal/worktree/identifier_operand_test.go, cmd/bench/command_registry_test.go, CHANGELOG.md

## What to build

An operator who runs `bench worktree exec` or `bench worktree path` on a
target that does not resolve reads why. The reasons are no match, a label
collision with the colliding ids, an inactive state, or the failed bundle
component. The
creation bundle validator returns the ticket-01 component error in place of
its two merged sentences. The selector's ambiguous-target sentence gains the
colliding ids. The resolver passes each error through unchanged.

Both verbs print the verb prefix and the resolver's error on stderr, and the
two verbs share one printer so their bytes cannot drift. The changelog
carries one user-facing line for the new refusal sentences.

## Acceptance

- [ ] LR11 prints the selector's no-match sentence for an unknown target.
- [ ] LR12 prints the selector's ambiguity sentence with the colliding ids.
- [ ] LR13 prints the not-active component with the assignment id for a released assignment.
- [ ] LR14 prints the owner-marker component after a marker rewrite.
- [ ] LR15 proves `worktree path` and `worktree exec` print byte-identical stderr for one broken target.
- [ ] LR16 turns red on a registry entry without a producing fixture.
- [ ] LR17 turns red on a production detail sentence that joins two component names with ` or `.
