# Close adapter blocker metadata

Accepted exact-candidate review repair for finding S1. Three producing tickets
name `adopt-exact-landing-in-commit.md` as their commit-adapter consumer, but the
dependent ticket's `Blocked by:` field names only the landing-owner producer.

## What to build

Make the adapter ticket's blocker declaration exactly resolve every ticket basename
that produces one of its named integration surfaces. Do not change ticket scope,
code, or lifecycle behavior.

## Acceptance

- [ ] [BR1] `adopt-exact-landing-in-commit.md` names `build-exact-landing-owner.md`, `reuse-exact-green-before-gate-lock.md`, `preserve-prospective-gate-output.md`, and `resolve-story5-fixture-gitdir.md` in `Blocked by:`, with no duplicate or unresolved basename.
- [ ] [BR2] The spec's ticket graph and conformance example-agreement checks accept the repaired metadata.

Ownership fence: specs/exact-prospective-landing/tickets/adopt-exact-landing-in-commit.md

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| BR1, BR2 | omit any one producer named by the adapter's integration surfaces | ticket dependency validator | run the focused ticket-graph/conformance check and expect the dependent-resolution rule to fail |
