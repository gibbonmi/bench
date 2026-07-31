# Retain subject-addressed gate evidence

Blocked by: none

Ownership fence: `internal/gate`
Assumptions: the current closed-subject builder remains the one oracle-identity owner

## What to build

Extend the `internal/gate` seam so callers can inspect and execute an unpublished
prospective tree as one closed subject, retain exact green evidence by subject,
and validate the branch-scoped project-green marker without parsing the latest
gate-cache projection. Keep the subject builder and verdict store single-sourced;
restructure the crowded package by responsibility rather than adding another
parallel derivation.

## Acceptance

- [ ] [R03] Bootstrap inspection imports only exact evidence validated by the gate owner and never runs the gate.
- [ ] [R05] Project-green requires the branch tip, marker, and retained closed subject to agree.
- [ ] [R30] Prospective execution drives the real gate wrapper and launcher closure from the unpublished tree.
- [ ] [R36] Tree, oracle, launcher, tools, environment, policy, and freshness drift invalidate evidence; history-only drift does not.
- [ ] [R38] Reds are not reusable and the latest cache remains a projection rather than the retained authority.
