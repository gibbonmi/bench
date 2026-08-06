# Close adapter blocker metadata

Blocked by: build-exact-landing-owner.md, reuse-exact-green-before-gate-lock.md, preserve-prospective-gate-output.md, resolve-story5-fixture-gitdir.md
Ownership fence: `specs/exact-prospective-landing/tickets/adopt-exact-landing-in-commit.md`
Integration surfaces: dependent blocker declaration→`specs/exact-prospective-landing/tickets/adopt-exact-landing-in-commit.md`; landing-owner producer basename→build-exact-landing-owner.md; gate-reuse producer basename→reuse-exact-green-before-gate-lock.md; gate-output producer basename→preserve-prospective-gate-output.md; story5-fixture producer basename→resolve-story5-fixture-gitdir.md
Contracts: the four producer ticket basenames cross `specs/exact-prospective-landing/tickets/adopt-exact-landing-in-commit.md`→lifecycle blocker resolution, asserted by BR1 against the real sibling producer tickets' `Integration surfaces:` declarations

Accepted exact-candidate review repair for finding S1. Three producing tickets
name `adopt-exact-landing-in-commit.md` as their commit-adapter consumer, but the
dependent ticket's `Blocked by:` field names only the landing-owner producer.

## What to build

Make the adapter ticket's blocker declaration exactly resolve every ticket basename
that produces one of its named integration surfaces. Do not change ticket scope,
code, or lifecycle behavior.

## Acceptance

- [ ] [BR1] `adopt-exact-landing-in-commit.md`'s `Blocked by:` list is exactly the four producer basenames named by the producer tickets' `Integration surfaces:` declarations — `build-exact-landing-owner.md`, `reuse-exact-green-before-gate-lock.md`, `preserve-prospective-gate-output.md`, and `resolve-story5-fixture-gitdir.md`.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| BR1 | omit or corrupt any one producer basename in the adopt ticket's `Blocked by:` list | the four producer tickets' `Integration surfaces:` declarations | read the four producer tickets' `Integration surfaces:` declarations, compare that basename set with the adopt ticket's `Blocked by:` set, and expect the two sets to differ |
