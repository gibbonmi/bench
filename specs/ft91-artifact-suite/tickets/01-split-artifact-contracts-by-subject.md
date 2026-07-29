# Split artifact contracts by subject

Blocked by: none

## What to build

Expose the 33 artifact contract tests as four independently schedulable subject
packages, with shared fixture policy owned once and the prepared-artifact
singleton retained in the prepared package.

## Acceptance

- [ ] The old artifact root owns no migrated tests, and the posture, offline,
  prepared, and distributable packages own exactly the spec's 4/14/13/2
  inventories.
- [ ] Every subject package enters through the shared package-main runner
  without copying the shared-cache token policy.
- [ ] The six prepared-artifact sharers select one package-scoped singleton in
  a fresh process, and the existing failure, mutation, and read-only belts stay
  green.
- [ ] All five behavior-owned artifact canaries bind to their new owner package
  without changing their TEST, EXPECT, BASE, or mutation content.
- [ ] The omission, inline-runner-policy, and per-sharer-state mutations each
  make their named assertion red before the correct topology is restored.
