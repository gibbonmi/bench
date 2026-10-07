# Share canonical and cancellation observations

Blocked by: none
Writes: internal/conformance/canonical_path_owner_test.go, internal/conformance/cancel_signal_registrations_test.go, internal/conformance/check_bindings_test.go, internal/conformance/checks_test.go, internal/conformance/tier_live_tree_test.go, specs/conformance-observation-follow-on/assets/
Covers: COF01, COF02, COF03, COF04, COF05, COF06, COF15, COF16, COF17, COF18, COF22, COF23, COF29

## What to build

Planning and slicing may proceed against LTE's approved value contracts.
Before any implementation charge, require complete independently accepted LTE to have landed, including all nine tickets and LTE-C1 through LTE-C8.
After that landing, finish concrete API binding and the landed-tree headroom recheck before any implementation charge.
Refresh the affected caller census against the same source.

If the landed source lacks the accepted contract, amend the plan before writing code.
Do not create a private walker, reader, parser, classifier, or dispatcher.
This is a planning ticket, not implementation admission or evidence that a future test has run.

Migrate canonical-path-owner and cancel-signal-registrations through their actual registered bindings in one green slice.
Keep both root-taking checks as fresh direct adapters to the complete LTE owner.
Keep literal selector matching, cmd/internal non-test selection, owner exclusions, and existing uniqueSorted diagnostics.
Cancellation's signal.Notify byte prefilter must still precede parsing.

Use existing small visitor files for the new dispatch, policy, freshness, and owner-consumption tests.
TestResidualObservationDispatch selects the actual two bindings and independently names their fixture paths and execution membership.
Overlapping complete regular source within ControlRecordLimit requires one physical read per requested path and one parse per requested path/mode.
Each requested directory is observed once per subject snapshot.

TestResidualObservationFreshness proves separate root/kit subjects and clean-red-clean runs through dispatch and direct adapters.
TestResidualObservationBindingsUseSharedOwner traces entries and helper/import aliases, rather than trusting exposed counters alone.
Do not mutate shared bytes, AST nodes, or token positions.

The existing canonical and cancellation canary fixtures remain read-only.
No executable check, registry policy, or universal analyzer is added.
COF-C1 closes independent Standards, Spec, and Coverage review before ticket 02 starts.

## Acceptance

- [ ] Complete LTE acceptance and landing are pinned before build admission; a partial LTE chunk is refused as a dependency.
- [ ] Same-function literal Abs plus EvalSymlinks retains the exact function diagnostic through the registered canonical owner.
- [ ] Notify and NotifyContext retain current argument/spread refusals and exact diagnostics through the registered cancellation owner.
- [ ] Malformed source without signal.Notify stays silent; source in excluded owners and test partitions keeps its current disposition.
- [ ] Actual two-binding dispatch matches independently authored membership and path/operation counts; per-binding snapshots and helper-only invocation fail.
- [ ] Distinct root and kit violations remain separate; reversing policy processing leaves shared observations and each diagnostic order unchanged.
- [ ] Existing TestEveryRetainedFixtureBitesThroughRegisteredOwner and TestCancelSignalRegistrationsBites retain their purposes.
- [ ] Each adapter migrated by this ticket observes direct and dispatch mutation/restoration as clean, red, clean with fresh invocation counts.
- [ ] Each migrated entry and its resolved helper closure consumes the shared owner, with no surviving private walk, read, parse, or alias-hidden replacement.
- [ ] Each migrated policy preserves hostile-path diagnostics, exclusions, ordering, and its existing reader-error posture before this checkpoint.
- [ ] Complete accepted LTE has landed, and concrete API binding plus landed-tree headroom recheck finish before any implementation charge.
- [ ] Each changed file fits the current landed-tree structure limit at this checkpoint, without a new grant or another direct conformance file.
- [ ] Every existing assertion and retained canary meaning remains unchanged.
