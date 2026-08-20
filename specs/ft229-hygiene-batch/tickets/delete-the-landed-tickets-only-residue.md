# Delete the landed tickets-only residue

Blocked by: close-the-light-path-ticket-on-landing.md, count-tickets-only-folders-in-status.md
Writes: specs/

## What to build

Every tickets-only folder under `specs/` is deleted — 37 of them at HEAD.
Re-count at ticket entry rather than trusting that number, because a light-path
landing between now and then adds one. After this ticket the status row this
build added reports nothing, because there is nothing to report.

No standing check holds the count at zero. The close step requires an in-flight
ticket folder to survive its own landing gate run, so a gate-level zero-count
assertion would red on every future light-path change.

## Acceptance

- [ ] no direct child of `specs/` holds tickets without a `spec.md`.
- [ ] `bench status` renders no tickets-only row at the landing commit (H06).
