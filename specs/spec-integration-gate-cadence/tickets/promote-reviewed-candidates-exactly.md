# Promote reviewed candidates exactly

Blocked by: Retain subject-addressed gate evidence, Record review and recover runs

Ownership fence: `internal/specbuild`, `tests/canary/spec-build-lifecycle`
Assumptions: the gate owner can authorize an unpublished exact prospective tree

## What to build

Complete `internal/specbuild` promotion and drift handling. Construct the
canonical staged-to-implemented prospective tree without editing the live spec,
ask the gate owner to authorize that exact subject, publish a byte-identical
squash child and then its branch-scoped green marker, retain provisional
evidence, and classify red outcomes without restoring ticket gates.

## Acceptance

- [ ] [R29-R35] Promotion gates the unpublished implemented tree, publishes nothing on red, lands one byte-identical squash outside provisional ancestry on green, retains inspectable evidence, and recovers branch-before-marker interruption.
- [ ] [R37] Working-base advancement recomposes before mutation; conflicts preserve refs and changed candidates require review again.
- [ ] [R39-R40] Candidate-owned reds route to attributed repair, inherited/infrastructure reds route to diagnosis or retry, and cap exhaustion leaves the run provisional.
- [ ] [R58] Repeated prospective status construction produces the same tree object.
- [ ] Canary mutations for receipt bypass, marker-before-branch, and a non-promote gate caller each make the gate red with their own diagnostic.
