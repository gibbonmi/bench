# Consume release evidence directly

Blocked by: none
Writes: internal/releasepreflight/types.go, internal/releasepreflight/command.go, internal/releasepreflight/decision.go, internal/releasepreflight/identity.go, internal/releasepreflight/vulnerability.go, internal/releasepreflight/decision_test.go, internal/releasepreflight/external_test.go, internal/releasepreflight/identity_test.go, internal/releasepreflight/evidence_contract_test.go (new), internal/publication/approved.go, specs/release-evidence-forwards/decisions/architecture-release.md, specs/release-evidence-forwards/assets/migration-reference.json (new)
Covers: RF01, RF02, RF03, RF04, RF05, RF06, RF07, RF08, RF09, RF10, RF11, RF12, RF13, RF14, RF15, RF16, RF17, RF18, RF19, RF20, RF21, RF22, RF23, RF24, RF25, RF26, RF27, RF28, RF29

## What to build

Deliver the complete preflight caller migration to releaseevidence and remove the redundant forwarding file in one green checkpoint.
The evidence owner already exists and remains unchanged; there is no preparatory owner ticket or new facade.
Keep releasepreflight.Command, Decide, their signatures, argument policy, phase execution order, and all public command/export surfaces.
The implementation line stays gpt-5.6-sol/high, with one fresh author and independent RF1 chunk review before final integration.
Refresh the exact caller census, closure, interfaces, and headroom against the implementation source before the charge.

Before qualification or deletion, capture the existing forwarding API's complete canned compatibility family in the fenced migration-reference.json asset.
Use the temporary same-package driver in evidence_contract_test.go, then remove that driver after comparison and retain only the focused command regression.
Keep the independently reviewed reference separate from candidate results; do not create a permanent compatibility harness or elapsed-time comparison.
Record the actual reference source, input family, error classifications, inventories, and raw record bytes.

The differential covers PhaseNames for both modes and every returned PhaseDefinition, TerminalStatus for all six status families, and existing Decide inputs.
It covers ReadRegular and ReadPackageVersion on all accepted reader cases and PromoteEvidence with fixed manifests and green, red, and interrupted results.
FinalizeEvidence covers focused verify, focused public publish, focused bank publish, and an already-cancelled context.
Use identical values and relative paths with normalized temporary-root labels; preserve every non-root error byte, record byte, and framing newline.
Compare complete promoted file inventories and bytes, exact errors, errors.Is missing-file/cancellation identity, and errors.As concrete ReleaseIntentError identity.

Migrate the six active functions at every enumerated caller: PhaseNames, phaseDefinition, terminalStatus, FinalizeEvidence, readPackageVersion, and readRegular.
Use the existing owner names PhaseNames, PhaseDefinitionFor, TerminalStatus, FinalizeEvidence, ReadPackageVersion, and ReadRegular.
Qualify DecisionInput, Decision, runner fields, signatures, literals, and every caller test with the existing owner types and constants.
Remove all ten dormant functions, eighteen aliases, ten constants, and the requirements snapshot with types.go after that complete migration.
Do not remove owner implementations, recreate dormant exports elsewhere, or count owner-origin references as additional forwarding consumers.

In Command, recognize *releaseevidence.ReleaseIntentError directly with errors.As, without a new wrapper or replacement error.
Add TestCommandFocusedPublishRequiresEvidenceAuthority using the existing same-package identityFixture and its agreeing tag, versions, committed HEAD, and toolchain.
Set and restore the working directory and environment on failure; do not run this test in parallel or copy its Git fixture.
Invoke --mode publish --profile public --phase identity so identity succeeds before the existing focused-scope refusal.
Require exit 1 and exactly {"kind":"requirement","message":"focused publish runs cannot authorize publication"} followed by its existing newline.
A generic evidence diagnostic or an earlier identity/tool failure does not satisfy this witness.

Keep Command's FinalizeEvidence call effective through the existing release-evidence-probe.
That registered ship check authenticates the committed source, builds its selected preflight executable, and drives scripts/release-preflight.sh.
Its phase stubs still reach the real finalizer, generated release-index.json, offline consumer, and attributed native-proof mutations.
Do not copy the producer harness, replace the consumer with field-presence assertions, or treat local npm stubs as live qualification.

The author omission substitutes a nil finalizer result while retaining the assembled run reference so the command compiles.
The actual release-evidence-probe must refuse with release evidence probe did not generate release-index.json.
The independent coordinator omission forces only concrete ReleaseIntentError recognition false while retaining a compilable Command.
The focused command witness must fail because its requirement diagnostic becomes the generic evidence diagnostic.
For each omission, pin the exact source and diagnostic before mutation, demonstrate behavioral red, restore byte-identically, and require the same focused green.
Compilation, setup, skipped, unattributed, and restore-failed failures close no proof.

When deleting types.go, replace only its compiled-map structured source entry with internal/releaseevidence/types.go and correct that entry's Supports and Drift.
Preserve all resolved answers, historical citations, and relative topic links; run decision-map-integrity on the resulting map.
All four moved-map preimages remain byte-identical during planning; RF27's record remains separate from this authorized future source-entry correction.
The root coordinator owns C11 navigation reconciliation outside this implementation ticket.
Correct only the requiredness owner name above releaseIndexAuthority in approved.go to releaseevidence.
No publication runtime body, registry schema, evidence artifact path, command binding, count, fixture, anchor, or injected port changes.

Retain every existing assertion in the migrated caller tests and all unchanged evidence and publication tests.
Reuse TestReleaseIndexBindsComponentManifestDigest, TestFirstPublicationRecordsPlatformsBeforeWrapper, the npm held-lock staging refusal, and the existing resume/promotion/rollback journeys.
Keep RunFirstPublication, RunStagedPublication, RunPromotion, RunRollback, runSubmit, and FixtureRegistry implementations byte-identical to the pre-charge source.
The npm staged refusal still wins before AcquireReleaseLock and registry construction; fixture StageSubmit and Approve remain implemented.
Public CLI/export use stays unknown outside the tree; public staging, broader FT368, identity alias-table deletion, and FT142/FT306 qualification stay excluded.

All first-use caller, reader, reference, error, assertion, and ownership obligations close in this ticket's own green checkpoint.
Canonical FixturePins, ReferencingFiles, and BoundFiles currently require no additional path co-names for its twelve Writes entries.
Requery these owners on the fresh implementation source; stale or missing closure needs an in-scope plan amendment before writes.
Keep command.go and the new contract test within 400 lines without a grant; preserve the nine-file preflight inventory after one deletion and one addition.
Do not grow inherited oversized publication test or state-machine files; only approved.go's owner comment changes in publication.

## Acceptance

- [ ] Every active production and test caller consumes the existing evidence owner directly, while Command and Decide retain their callable signatures and behavior.
- [ ] The complete six-active/ten-dormant/eighteen-alias/ten-constant/one-snapshot census closes after types.go deletion, without a new facade or surviving duplicate.
- [ ] All existing caller assertions remain effective and the hidden-tree sweep retains the generated and registered Command consumers without editing their bindings.
- [ ] Verify/publish phase order, immutable decision verdicts, invalid-mode refusal, and interrupted-over-red precedence match the independently retained canned reference.
- [ ] The identityFixture command scenario reaches the focused-publish refusal and emits the exact requirement JSON diagnostic at exit 1.
- [ ] Public and bank focused refusals retain their exact messages and concrete ReleaseIntentError identity; cancellation and missing files retain errors.Is identity.
- [ ] Every absent, empty, ordinary, directory, and symlink reader case, version case, space path, and non-ASCII path matches the reference bytes and errors.
- [ ] Fixed green/red/interrupted promotion retains complete file inventories, nullable fields, final newlines, and every encoded byte for absent and existing real dist.
- [ ] Non-directory and symlink dist retain their refusal; focused verify and already-cancelled finalization retain their effects and error identity.
- [ ] The authenticated release-evidence-probe still generates an index accepted by its real offline consumer and attributes its existing native-proof mutations.
- [ ] The compiling finalizer omission reds that registered probe specifically for missing release-index.json; byte-identical restoration returns the same focused green.
- [ ] The compiling errors.As-recognition omission reds the focused Command diagnostic; byte-identical restoration returns the same focused green.
- [ ] No invalid, setup-only, compilation-only, skipped, unattributed, or restore-failed proof is retained as behavioral evidence.
- [ ] The compiled map resolves after deletion with the sole source-entry correction, while every resolved answer and historical citation remains unchanged.
- [ ] All four moved-map planning preimages remain byte-identical to the accepted source, and the C11 root index remains outside this ticket.
- [ ] Only approved.go's requiredness owner comment changes; the entire evidence owner, registries, public surfaces, and publication runtime bytes remain unchanged.
- [ ] Npm staged refusal beats the held lock before registry effects; fixture staging remains implemented without a complete staged-qualification claim.
- [ ] Existing first-publication ordering, resume, adapter-selected promotion and rollback, component digest, and authority assertions remain unchanged and pass.
- [ ] All 29 rows close in RF1 without later-ticket repair, a new fixture/scanner/port, a live qualification claim, or retirement of FT142/FT306/broader FT368.
- [ ] Actual complete-delta file growth fits existing caps, the new test stays within 400 lines, and the preflight package remains at nine files without a grant.
