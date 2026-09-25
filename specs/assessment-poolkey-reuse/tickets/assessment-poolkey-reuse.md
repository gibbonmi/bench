# Use the poolkey predicate for the expected assignment id

Blocked by: none
Writes: internal/assessment/collection.go, internal/census/events_test.go
Covers: none

## What to build

The assessment collector makes sure that each expected assignment id is valid. It
builds a pool segment from the id twice, then splits the segment again. The poolkey
package exports `IsAssignmentID`, which is the one owner of this rule. The collector
calls that predicate instead.

The census event reader refuses a name that is not an assignment id. No test covers
that refusal, so a probe of the check stays silent. A new test in the census package
covers a stray file name and a name that climbs out of the census directory.

## Acceptance

- [ ] The collector calls `poolkey.IsAssignmentID` for each expected assignment id, and the collector no longer composes a pool segment for that check.
- [ ] `ReadEvents` refuses a stray file name and a climbing name, and it returns no events and no problems for them.
- [ ] A probe that removes the id check from `ReadEvents` turns the new census test red, and the restore returns it to green.
- [ ] `bench test --package ./internal/assessment` and `bench test --package ./internal/census` pass, and the gate is green at the landing.
