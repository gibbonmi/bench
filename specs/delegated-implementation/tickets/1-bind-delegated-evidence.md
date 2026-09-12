# Bind delegated identities to completion evidence

Blocked by: none
Writes: internal/reviewrecord, internal/gate/delegated_checkpoint_test.go (new), internal/gate/review_checkpoint_test.go, internal/landing/delegated_completion_test.go (new), internal/landing/completion_evidence_test.go, internal/preflight/delegated_evidence_test.go (new)
Covers: DI1, DI2, DI3, DI4, DI5, DI6, DI7, DI8, DI9, DI10, DI11, DI12, DI31, DI36, DI37, DI38, DI41

## What to build

Extend the existing plan and record owners with the delegated identity contract in the spec.
Make the existing chunk checkpoint and final landing accept valid delegated evidence and refuse invalid ownership.
Keep the version 1 path unchanged.
Use the existing recordtest fixture owner for source histories and evidence.
Exercise preflight through its existing completion-evidence projection.
This ticket delivers a usable evidence capability before phase guidance exposes the flag.

Its consumer receives a version 2 source-bound plan and record accepted by the ordinary checkpoint and landing owners.
No separate identity registry or verification engine is introduced.

## Acceptance

- [ ] The version 1 fixture matrix matches its pre-change outcomes.
- [ ] Valid delegated evidence reaches a green chunk checkpoint and final landing.
- [ ] Identity edits outside the source-bound plan cannot change verification ownership.
- [ ] All three reviewers exclude the orchestrator and every current or former author.
- [ ] One session cannot satisfy two axes for a chunk.
- [ ] Undispatched future tickets accept empty histories within a bounded author limit.
- [ ] Historical occurrences retain their original frozen assignments after replacement.
- [ ] New post-replacement occurrences require the successor and fresh verification.
- [ ] The existing source-chain, native-result, probe, reconciliation, and destination-delta refusals remain effective.
- [ ] Preflight accepts and describes both valid record versions without weakening invalid-record diagnostics.
- [ ] A two-ticket chunk accepts each distinct author for its own verification requirements.
- [ ] Foreign, missing, and incorrectly performed ticket obligations refuse chunk acceptance.
- [ ] Reusing one native author session across two ticket histories refuses acceptance.
