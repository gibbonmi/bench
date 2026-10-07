# Production test-seam admission (FT343; C09)

Status: ready

## Destination

Make the existing injected-port audit cover actual production test-seam forms.
Admit justified fault seams while preventing invisible bypasses of production behavior.

## Notes

The reviewer delegated architecture recommendations on 2026-10-06.
A test seam supplies a controlled dependency or fault at a production interface.
Avoid: a test-only success path that skips the subject.
Use craft-seams, craft-spec, and craft-gate for the policy spec.

## Decisions so far

- [Which seam forms should the audit admit?](architecture-test-seams/tickets/1.md): prefer instance-local seams and register justified exceptions.
- [Which current hook claims need correction?](architecture-test-seams/tickets/2.md): retain proven consumers and classify exact unused hooks.
- [What must prove that the policy is enforced?](architecture-test-seams/tickets/3.md): require omission tests for each admitted seam form.

## Not yet specified

## Spec-writer discretion

- Audit representation within the existing owner.
- Instance-local collaborator names that preserve the production entry point.

## Out of scope

- A blanket deletion of every setter or environment hook.
- A blanket registry exemption for an entire package.
- A weaker omission oracle or replacement of local process behavior with mocks.
- Treating a restore closure as proof of concurrency safety.

## Sources

- Path: `roadmap/FT343.md`
  Supports: the seam-policy question and stale unused-hook premise.
  Drift: a change to the row.
- Path: `internal/conformance/injected_ports_test.go`
  Supports: the audit's current named-type derivation.
  Drift: a change to derivedInjectedPorts.
- Path: `internal/git/git.go`
  Supports: existing package setters.
  Drift: a change to their scope or tests.
- Path: `internal/otelrecord/provider.go`
  Supports: the fallback-home test-runtime branch.
  Drift: a change to recordHome.
- Path: `internal/releasepreflight/identity_test.go`
  Supports: a real REF hook consumer.
  Drift: a change to identityFixture.
- Path: `internal/conformance/native_workflow_test.go`
  Supports: dynamic phase override consumers.
  Drift: a change to the release probe.

