# Derive the fixture evidence ID from its owner

Blocked by: none
Writes: internal/reviewrecord/ticket.go, internal/reviewrecord/recordcmd/ticket.go, internal/reviewrecord/recordtest/fixture.go, internal/reviewrecord/fold_test.go
Covers: none

## What to build

The review-record fixture makes the evidence ID of a verification from a prefix and a requirement ID. The `record ticket` form makes the same ID with its own composer. The two derivations of one fact can drift, and a change to the production rule does not reach the fixture.

Move the composer into `internal/reviewrecord`, because the fixture and the command package both import that package. The command package calls the moved composer. The fixture calls the moved composer. The fold test that makes a final evidence ID by hand also calls it.

The composer keeps its rule without change, so the fixture now also writes a requirement ID that already holds the prefix only once.

## Acceptance

- [ ] `internal/reviewrecord` holds the one evidence-ID composer, and `recordcmd` holds no copy of it.
- [ ] `recordtest` makes each verification evidence ID with the composer, so no `prefix+"-"+requirement.ID` expression stays in the fixture.
- [ ] A mutation of the composer reaches the fixture records. Before this change, the test binaries of `internal/reviewrecord` and `internal/gate` did not link `recordcmd`, so a probe of the composer was silent there. After this change, a probe that makes the composer return only the prefix bites in `internal/reviewrecord`.
- [ ] `bench test --package ./internal/reviewrecord/...` passes.
