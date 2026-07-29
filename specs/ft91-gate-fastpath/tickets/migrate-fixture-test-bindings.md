# Migrate the 33 behavior-owned fixtures to TEST bindings

Blocked by: Carry per-test bite scoping through the canary runner

## What to build

Story 7 of `specs/ft91-gate-fastpath/spec.md`: every behavior-owned fixture
under `tests/canary/behavior-owned/` gains a `TEST` file naming its EXPECT's
owning contract test, so the sweep scopes each fixture to one owning test.

## Acceptance

- [ ] All 33 behavior-owned fixtures carry a `TEST` file naming the owning
      `Test*` function of their EXPECT's failure message.
- [ ] The full canary sweep is green with scoping active (every fixture bites,
      none vacuous) — proven at the landing commit's gate.
- [ ] One deliberately wrong owner demonstrated red (did-not-bite) during the
      build and reverted, recorded in the return log.
