# Separate shaping from spec authoring

Blocked by: Bind decision-map integrity to the gate

## What to build

Kit commands, skills, glossary, profile, README, generated guidance, and
changelog consistently make decision maps situational, reserve decision tickets
for reviewer choices, and give spec authoring ownership of engineering seams and
coverage.

## Acceptance

- [x] Spec writing accepts a ready compiled map, a reviewer-confirmed current
  conversation, or a named reviewed artifact and records exactly one Decision
  source line.
- [x] Spec writing re-verifies map Sources and asks at most two late
  clarifications one at a time with a recommendation before routing dependency
  fog back to shaping.
- [x] Shaping uses decision-ticket vocabulary, points to `bench maps --template`,
  removes Handoff, and owns reviewer decisions, constraints, exclusions,
  research objects, and bounded discretion.
- [x] Spec guidance owns engineering seams, tests, acceptance coverage, hostile
  inputs, and gate attachment without copying a second research manifest.
- [x] `craft-tickets` keeps independently-green implementation-ticket semantics,
  while CONTEXT, README, the profile, adapters, indexes, and CHANGELOG use the
  new phase vocabulary consistently.
- [x] Guidance conformance mutations fail when any authorization, ownership, or
  vocabulary anchor is removed or reverted.
