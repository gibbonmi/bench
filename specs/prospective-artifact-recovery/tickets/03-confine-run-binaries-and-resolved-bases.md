# Confine run binaries and resolved bases

Blocked by: 02-protect-shared-prospective-owners.md
Writes: internal/gate/prospectiveartifact/prospectiveartifact.go, internal/gate/prospectiveartifact/prospectiveartifact_test.go, internal/gate/prospective.go, internal/gate/prospective_owner_test.go

## What to build

Repair review findings P1, S6, C1, S1, C2, P4, P5, and S4 from `reviews/prospective-artifact-recovery.md`.

Set the bundle root on both run-binary branches, so an authored binary never lands outside the bundle.
Resolve the bundle base through symbolic links before the owner creates or sweeps a bundle.
Collapse the probe classification to one `ESRCH` guard.
Strengthen the partial rows PAR21 and PAR33 and add the record-mode and symlinked-base rows.

## Acceptance

- [ ] P1: a candidate without its own kit authors its run binary under the bundle root.
- [ ] C1: a bundle under a symlinked temporary root removes its exact Git registration.
- [ ] C2: a 0644 or 0400 owner record retains its bundle root.
- [ ] PAR21: the hostile-path bundle carries a Git registration and cleanup stays confined.
- [ ] PAR33: a registration-removal failure creates no new bundle or checkout.
- [ ] S1, S4, S6: the one-source repairs land with no behavior change.
