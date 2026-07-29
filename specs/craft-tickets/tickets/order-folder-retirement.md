# Order folder retirement

Blocked by: none

## What to build

Replace whole-folder `RemoveAll` retirement with the specified recovery order:
review pickup, `tickets/`, `spec.md`, then the now-empty spec folder. Preserve
the accepted re-run outcomes for every partial state and make the `RelTo`
comment describe current folder output.

## Acceptance

- [ ] Story 20 and its acceptance-coverage row are green.
- [ ] Fault-seam or black-box cases prove the named recoverable and terminal boundaries.
- [ ] No retirement comment advertises the old flat live-spec output.
