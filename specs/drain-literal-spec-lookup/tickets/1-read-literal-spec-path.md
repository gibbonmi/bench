# Read the literal spec path

Blocked by: none
Writes: internal/spec/tickets_only.go, internal/landing/close_test.go, CHANGELOG.md
Covers: none

## What to build

The committed-tree tickets-only reader treats glob syntax in a slug as literal
text when it checks for `spec.md`.

## Acceptance

- [ ] A literal glob slug with no `spec.md` qualifies as tickets-only in a commit tree.
- [ ] A literal glob slug with its own `spec.md` does not qualify as tickets-only.
- [ ] Existing normal, absent, file, and traversal cases keep their current results.
