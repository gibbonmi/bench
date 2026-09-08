# FT311 delegate preparation and review dispatch

Status: staged

Decision source: specs/ft311-delegate-preparation/decisions/ft311-coordinator-work.md

Verification log: 0 iteration(s) to accept — draft awaits the required native review.

## Problem

Coordinators reconstruct revisions, ticket fences, evidence handles, and repeated charge instructions before each delegation.
Review preparation repeats this work for independent axes.
Handwritten reconstruction can omit an input or mix evidence from different revisions.
The existing preflight checks validate artifacts but do not prepare the complete mechanical charge.

## Solution

Extend preflight with opt-in charge preparation and read-only fence proposals.
The active phase supplies task judgment and verifies reviewer approval before dispatch.
An authorized review phase launches the independent native axes after preparation succeeds.
A mid-tier delegate can classify repair evidence and draft a repair charge for coordinator acceptance.

This is the first of five FT311 specs.
It serves this kit and repositories that link it.
It changes neither tier defaults nor the landing gate.
Successful preparation is evidence readiness, not permission to write or proof that the work passes.

## User stories

### Prepare a bounded assignment

Line: gpt-5.6-terra / medium.
The CLI composes existing ticket, snapshot, and preflight owners.

1. As a coordinator, I want a charge for one selected ticket, so that its mechanical inputs need no reconstruction.
2. As a delegate, I want the exact assignment and frozen revision, so that I work on the intended tree.
3. As a delegate, I want the ticket and its coverage rows, so that I retain the approved task requirements.
4. As a coordinator, I want missing inputs to refuse preparation, so that an incomplete charge cannot authorize dispatch.
5. As a coordinator, I want stale or mixed evidence refused, so that one charge describes one source.
6. As an existing caller, I want unchanged plain preflight behavior, so that the new forms do not break phase entry.
7. As a coordinator, I want bounded output with complete retrieval, so that a large charge loses no required evidence.
8. As a repository owner, I want input text treated as data, so that charge preparation cannot execute ticket prose.

### Propose ownership closure

Line: gpt-5.6-terra / medium.
Existing closure registries supply the proposed paths.

9. As a coordinator, I want missing fixture and registry paths proposed, so that I can repair a ticket fence precisely.
10. As a reviewer, I want proposals to grant no authority, so that new writes still require my approval.
11. As a coordinator, I want affected dependency edges identified, so that an approved fence change cannot create concurrent writers.
12. As a delegate, I want approved paths retained without repeat approval, so that preparation does not add a redundant pause.

### Prepare independent review evidence

Line: gpt-5.6-terra / medium.
Frozen diff, consumers, and coverage already have public owners.

13. As a coordinator, I want shared review evidence gathered once, so that each axis receives the same complete inputs.
14. As a reviewer, I want each axis to receive its canonical charge, so that generated packets do not redefine review policy.
15. As a coordinator, I want failed evidence collection to stop dispatch, so that omitted consumers cannot appear as a complete review.

### Consume charges in native phases

Line: gpt-5.6-terra / high.
Kit guidance and its enforcement require the guidance-authoring effort override.

16. As a reviewer, I want the phase to establish my approval, so that staged metadata cannot impersonate a decision.
17. As a coordinator, I want an authorized review to launch independent native axes automatically, so that preparation removes repeated dispatch work.
18. As a reviewer, I want separate axis returns, so that the coordinator can assess each finding against its source.
19. As a coordinator, I want a capable-harness handoff when native dispatch is unavailable, so that review retains its independence.
20. As a coordinator, I want stale prepared inputs checked before dispatch, so that delayed execution cannot consume an obsolete charge.
21. As a coordinator, I want a bounded mid-tier repair triage, so that ambiguous probe results become a reviewable repair proposal.
22. As a reviewer, I want repair acceptance and independent verification retained, so that delegation does not lower the evidence bar.
23. As a maintainer, I want canonical phase guidance without surviving fallback contradictions, so that linked harnesses follow the same authority boundary.
24. As a reviewer, I want the remaining FT311 capabilities kept separate, so that each approved spec has a clear outcome.

## Implementation decisions

### Charge forms

Keep the existing `review|build <slug>` entry modes.
Add `--charge` for generation, with explicit `--base` and `--source-tip` values.
The build form also requires `--ticket <basename>`.
The review form rejects `--ticket` and prepares the existing spec-backed axes.
New forms require an active assignment and a clean source at the supplied tip.
Existing forms retain their current grammar, checks, output, and exit meanings.

Add `--propose-writes` as a separate build form with the same ticket and revision arguments.
It cannot combine with `--charge`.
It presents repairs even when fixture or registry closure is red.
Invalid ticket grammar, unresolved revisions, or an unreadable registry still refuse the proposal.
The spelling replaces the roadmap's tentative `--fix-writes` name because the approved operation cannot apply edits.

The new forms use structured TOON output and exit 0 for successful preparation, 1 for refusal, and 2 for invocation errors.
Preflight retains its existing operational classification in the command registry.
Use the shared encoder, bounds policy, and error owners.

Default charge output contains identity, input handles, and completeness accounting.
`--full` retrieves the complete generated charge and evidence from the same explicit revisions.
A compact projection is never dispatch-ready until every required input has been retrieved and validated.
A refusal names its failed input and a concrete repair action.

### One source per charge fact

The charge contains these fields or sections.
The field labels below are the output contract for new forms.

| Label | Source and meaning |
|---|---|
| assignment | The current assignment owner and checkout identity |
| base | The resolved frozen base commit |
| source_tip | The resolved source commit |
| ticket | The selected ticket's path and content identity, or the reviewed spec for axis charges |
| writes | The ticket parser's exact ownership entries, or an explicitly read-only axis |
| evidence | Frozen source handles, required source text, and matched coverage rows |
| checks | Ticket acceptance text and canonical phase verification requirements |
| return | The canonical delegate return requirements, plus the selected axis requirements when applicable |
| complete | Whether the response includes every required input without truncation or refusal |
| next | A concrete retrieval or repair action when preparation is incomplete |

Do not parse a second ticket grammar.
Use the ticket enumeration and parser for fields, and retain the complete ticket bytes as required evidence.
An existing optional delegate-charge section survives in those bytes without a new required ticket field.
Checks are the stated acceptance and verification requirements, not guessed test selectors.
The coordinator adds executable selectors, the model and effort, the iteration cap, and the task-specific mutation before dispatch.
The phase refuses dispatch when that task-specific supplement is incomplete.

The canonical craft-delegate and craft-review sources own procedural requirements and axis meanings.
Generated sections cite their pinned source text rather than maintaining a parallel policy registry.
Coverage comes from the coverage owner, including each row's behavior, seam, and failure explanation.
The generator never invents an acceptance row, approval record, or test result.
A staged spec proves artifact state only.
The active phase must bind the packet to the reviewer-approved ticket and graph in the session's authorized decision source.

### Closure proposals

Compose the existing fixture and registry closure derivations.
Return each missing path once with its owning fixture or registry citation.
Resolve proposed paths against the spec fence and identify any required spec expansion.
List ticket pairs whose Writes would overlap after the proposal.
Keep existing dependency ordering when it already serializes those pairs.
For unordered pairs, report that an approved ordering edge is required without inventing its direction.

A proposal changes no ticket, spec, assignment, or approval state.
The coordinator obtains approval only for new authority, updates the approved ticket and affected edges, and reruns preparation.
Already-approved paths need no repeated permission.
A charge remains refused while required closure is absent.
Missing closure and malformed grammar together report the grammar refusal first.

### Frozen review evidence

The review form composes the existing diff, consumers, and coverage owners.
Resolve the pair once and collect the shared evidence once per preparation attempt.
Carry the diff snapshot and full patch, consumer tables with their citation, and complete coverage rows.
Preserve touched and deleted consumer distinctions.
Inject the executable's existing version into consumer citation generation through the command registration owner.
Do not spawn another Bench executable to reconstruct these facts.

Buffer the response until the final snapshot check succeeds.
A moved tip, dirty source, mismatched pair, collector refusal, or incomplete consumer projection prevents a complete charge.
Use existing movement-check behavior and its bounded retry policy.
A retry is a new preparation attempt, not a mixture of successful fragments from different attempts.
Use the existing source identities and content hashes for handles.
No durable packet cache or new assignment state machine is required.

Prepare Standards, Spec, and Coverage inputs from their canonical review sources.
Each charge references the shared evidence instead of collecting it again.
Axis-specific source reads remain independent review work.
The spec-backed form requires a spec and coverage map.
Historical and spec-less reviews retain their existing collection entry points.
Their native-dispatch authority rule follows the same updated review phase.

### Phase authority and native dispatch

The build phase consumes generated inputs only after the existing spec-and-ticket approval.
It verifies the selected ticket, fence, dependency readiness, expected source tip, and task-specific supplement before the write delegate starts.
A new charge does not establish that blockers have completed.
The coordinator establishes completion from the retained integration source and verified ticket returns.

The review phase establishes review authorization before it prepares inputs.
When preparation and runtime capability checks pass, it dispatches each independent axis through the native agent surface.
It does not ask for another approval within that already-authorized phase.
Each axis has its own context and the existing isolated read-only venue.
Independent source derivation, model routing, and coordinator verification remain mandatory.
The coordinator collects each return, accepts findings, and directs repairs.

Missing or failed returns remain incomplete review, never a clean finding set.

Native availability is a fact of this active session.
A compiled harness record cannot prove that a native tool exists here.
When the tool is unavailable or delegation is prohibited, preserve the prepared charges and stop dispatch.
Emit the repository, assignment, frozen pair, charge inputs, destination harness, and exact native continuation command.
Do not substitute a same-family CLI launcher or combine the axes in coordinator context.
Keep the separate cross-harness falsification route and its existing authority.

Before dispatch, revalidate the source identity and required source contents against the charge.
A changed ticket or source forces preparation again.
This is a workflow check against accidental drift, not an authentication protocol against a malicious harness.
The existing installed entrypoint and toolchain remain the execution trust assumptions.
Preparation does not execute generated prose or launch models.
No new publisher, bootstrap verifier, or independent executable authentication claim is introduced.

### Repair triage

Use a read-only mid-tier delegate for interpretation that needs triage.
Resolve the harness's bound mid model at medium effort for one iteration.
Use high effort only when the seam requires it under the existing line discipline.
Routine deterministic failure projection remains CLI work.
Give triage the failed command, selected checks, actual execution evidence, findings, and frozen source.
Missing or contradictory evidence remains unknown.

The triage return distinguishes an invalid probe, no executed test, missing coverage, and a scope defect when evidence supports that classification.
It proposes a bounded repair charge and cites the supporting evidence.
The coordinator accepts the diagnosis and routes approved repairs under existing authoring and verification rules.
Triage cannot mark a finding repaired, expand a fence, waive a probe, or change the default tier.
Kit guidance authorship remains mid/high.

## Testing decisions

The primary executable seam is preflight's public command over real temporary Git repositories.
Use the existing same-package command harness and explicit-base tests as precedents.
Exercise real ticket parsing, closure registries, snapshot ownership, and rendered output through that seam.
Do not copy their decision rules into a fixture implementation.

New command tests run in the ordinary Go phase of the project gate.
Guidance pins run through `docs-currency-workflow` and its existing workflow-guidance canary family.
Register new anchor storage with that family's existing owner list.
Each required procedural clause receives an omission mutation that makes the conformance check red.
Record the observed red for any independently authored expectation required by the one-source exception.

Native tool use is review-owned evidence.
Exercise the updated phase in a capable native harness and retain actual dispatch and return records.
Also exercise an unavailable surface without substituting another launcher.
Static anchors prove that required instructions remain present, not that a model obeyed them.
A missing live exercise is an explicit acceptance shortfall.
The full landing gate remains the final code oracle.

### Seam diagram

    trigger: authorized build or review phase
        |
        v
    frozen source + ticket --> [ preflight command ] --> prepared inputs
                                  ^ tests attach here       |
                                  real Git fixtures         v
                                                    [ native phase ]
                                                     |            |
                                                   delegates    handoff
                                                     ^
                                          review-owned live exercise

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| DP1 | 1, 2 | A build charge identifies the selected assignment and exact revision pair. | New TestChargeIdentity through preflight Command | A charge from the primary checkout or another assignment fails the identity assertion. |
| DP2 | 3 | Every selected ticket input survives in the complete charge. | New TestChargeTicketEvidence through preflight Command | Independently omit ticket bytes, one Writes entry, or one covered row to make the test red. |
| DP3 | 3 | Canonical checks and return requirements remain source-backed. | New TestChargeCanonicalRequirements through preflight Command | A changed canonical source must change its generated section without a second policy edit. |
| DP4 | 4 | An absent or invalid required input refuses a build charge with a repair action. | New TestChargeRequiredInputs through preflight Command | Missing assignment, spec, ticket, coverage, or procedural source cannot produce complete=true. |
| DP5 | 5 | A moved or dirty source refuses a complete charge. | New TestChargeSnapshotMovement through preflight Command | Change the tip or index between reads using the existing snapshot test seam. |
| DP6 | 6 | Legacy preflight responses retain their existing semantics. | New TestLegacyPreflightDifferential through preflight Command | Compare baseline and candidate results over the same enumerated legacy fixtures. |
| DP7 | 7 | Compact output identifies omitted content and its exact full retrieval command. | New TestChargeProjection through preflight Command | An oversized ticket cannot appear complete in the compact response. |
| DP8 | 7 | Full retrieval contains every required input from the pinned source. | New TestChargeFullRetrieval through preflight Command | Remove one evidence block while retaining its handle to make the test red. |
| DP9 | 8 | Hostile input remains data or produces a named refusal. | New TestChargeHostileInputs through preflight Command | Shell-looking prose cannot create a sentinel or split an output field. |
| DP10 | 9 | A closure proposal names each missing fixture and registry path once. | New TestWritesProposalClosure through preflight Command | Omit either closure producer's contribution to make the fixture red. |
| DP11 | 10 | A closure proposal leaves every tracked and local authority input unchanged. | New TestWritesProposalReadOnly through preflight Command | Compare ticket, spec, index, worktree, and assignment records before and after. |
| DP12 | 11 | Proposed overlap identifies every newly unordered ticket pair. | New TestWritesProposalEdges through preflight Command | A second writer added by fixture closure must appear as requiring an ordering decision. |
| DP13 | 12 | An already-covered path produces no new authority request. | New TestWritesProposalAlreadyCovered through preflight Command | Repeating the same approved closure yields no proposed addition. |
| DP14 | 9, 10 | Invalid grammar wins over a concurrent closure defect. | New TestWritesProposalRefusalOrder through preflight Command | Both defects reach the entrypoint and the response must name grammar first. |
| DP15 | 13 | Review preparation gathers shared evidence once per stable attempt. | New TestReviewChargeSharedEvidence through preflight Command | The real collectors' results are shared across axes rather than recollected per axis. |
| DP16 | 13 | Review charges preserve complete diff, consumer, and coverage evidence. | New TestReviewChargeEvidence through preflight Command | Use a touched=false consumer and a deleted declaration to expose lost evidence. |
| DP17 | 14 | Each axis charge references its current canonical source. | New TestReviewChargeAxes through preflight Command | A changed axis source must change its input identity without editing a second registry. |
| DP18 | 15 | Failed or incomplete evidence prevents a complete review charge. | New TestReviewChargeRefusals through preflight Command | Exercise collector refusal, mismatched pair, and omitted unrepresentable consumer rows separately. |
| DP19 | 16 | Dispatch requires the approved ticket and complete task-specific supplement. | Review-owned native build exercise plus docs-currency-workflow canaries | A staged but unapproved ticket must not start a writer. |
| DP20 | 17 | An authorized prepared review starts each native axis without another approval. | Review-owned native review exercise plus docs-currency-workflow canaries | Dispatch records show independent contexts without a second approval turn. |
| DP21 | 18 | An absent or failed axis return leaves review incomplete. | Review-owned native review exercise plus docs-currency-workflow canaries | Withhold one return and verify no clean aggregate is claimed. |
| DP22 | 19 | Missing native capability produces a capable-harness handoff. | Review-owned unavailable-surface exercise plus docs-currency-workflow canaries | The record retains charges and contains neither a substitute CLI launch nor inline axes. |
| DP23 | 20 | Changed prepared inputs stop dispatch until preparation repeats. | Review-owned native phase exercise plus docs-currency-workflow canaries | Move the source after preparation and verify no stale charge starts. |
| DP24 | 21 | Mid-tier triage returns an evidence-backed repair proposal. | Review-owned triage exercise plus docs-currency-workflow canaries | Invalid, empty, missing-coverage, and scope-defect evidence receive distinct supported dispositions. |
| DP25 | 22 | Repair acceptance retains coordinator authority and independent verification. | Review-owned repair exercise plus docs-currency-workflow canaries | A triage done-claim cannot close a finding or authorize a fence expansion. |
| DP26 | 23 | Canonical guidance has no surviving same-family CLI or inline-axis fallback for review. | New TestPreparedReviewGuidance through docs-currency-workflow | Swap the route while retaining command tokens to expose a contradictory fallback. |
| DP27 | 24 | This spec preserves the other four capability boundaries. | Review-owned spec and implementation scope audit | A tier-default, landing, recovery, or diagnostic behavior edit is outside this fence. |

### Edge inventory

The canonical walk covers empty input, boundaries, errors, repetition, process boundaries, and hostile environments.
The attached profile is the shell CLI checklist in projects/benchkit.md.
All new input checks serve linked repositories as well as the kit.

| Class | Concrete disposition | Rows |
|---|---|---|
| Absent versus empty | Distinguish missing and empty spec, tickets directory, ticket, and required guidance source in diagnostics. Neither state dispatches. | DP4 |
| Grammar boundaries | Reject duplicate options, unknown flags, missing values, invalid mode combinations, and ticket traversal. | DP4, DP9 |
| Paths | Exercise spaces, glob characters, trailing spaces, and Git-quoted patch paths through existing owners. | DP9, DP16 |
| Text sinks | Exercise ESC, BEL, tab, newline, return, and numeric-looking hashes through the shared TOON encoder. | DP9 |
| File kinds | Refuse live links, dangling links, FIFOs, and other special required inputs before reading their content. | DP4, DP9 |
| Text completion | Accept valid ticket text without its final newline. Preserve Unicode text without introducing a new whitespace parser. | DP2, DP9 |
| Repetition | Identical pinned inputs yield identical content identities. A repeated approved closure yields no new proposal. | DP8, DP13 |
| Snapshot movement | Change HEAD, index, or required input bytes during preparation and between preparation and dispatch. | DP5, DP18, DP23 |
| Partial evidence | A collector refusal or incomplete projection cannot stand for an authoritative empty result. | DP18 |
| Dependency graph | Exercise no overlap, existing order, transitive order, and a newly unordered pair. No direction is silently chosen. | DP12 |
| Competing failures | Malformed ticket plus missing closure reports the malformed source first. | DP14 |
| Process boundary | Read the complete output in a fresh harness context without coordinator memory. | DP19, DP20 |
| Dispatch interruption | Missing capability and missing returns retain incomplete state with an actionable continuation. | DP21, DP22 |
| Authority | Staged metadata, proposed paths, and triage recommendations confer no new authority. | DP11, DP19, DP25 |
| Self-changing output | Preparation performs no authority-input writes, including a repeated tracked-input run. | DP11 |

**Won't handle:** A durable dispatch scheduler or crash-resume protocol belongs to FT305. The active phase remains the in-scope caller.

**Won't handle:** Cross-harness launch design belongs to FT312. The review phase retains its existing falsification caller.

**Won't handle:** A malicious installed executable or harness needs an independent trust design. This preflight caller retains the existing installation boundary.

**Won't handle:** New historical or spec-less packet forms are separate CLI capabilities. Their existing review entry points remain available.

No new whitespace predicate, Git flag, package-variable substitution, deletion, or restore operation is required by this design.
If implementation needs one, its ticket must attach the corresponding hostile case before claiming its row.

## Ownership fences

These entries are the union of ticket Writes.
They authorize implementation only after the reviewer approves this spec and its ticket graph.

- `internal/preflight/`
- `cmd/bench/main.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/anchors/`
- `internal/conformance/registry_test.go`
- `internal/conformance/docs_workflow_helpers_test.go`
- `internal/conformance/ft311_preparation_test.go`
- `tests/canary/workflow-guidance-anchors/`
- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `.agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md`
- `reviews/ft311-delegate-preparation.md`

The build cannot edit this spec, its acceptance rows, or its ticket graph without the existing spec-change authority.
The anchor prefix permits a separate registry file rather than growing the over-budget registry body.
The workflow canary prefix supplies the family closure for its pinned owners.
Existing generated adapters consume the canonical commands without copied content.

## Out of scope

These are planning estimates of file edits, not measured implementation costs.
Each future capability gets its own specification and one whole-project landing gate.

| Capability | Estimated price | Derivation |
|---|---|---|
| Diagnostics | 12 edits, 1 gate runs | Probe, prose, and anchors owners plus tests and public help closure |
| Landing completion | 16 edits, 1 gate runs | Landing effects, broker refresh, cleanup, handoff state input, and retrospective facts with tests |
| Recoverable reset | 10 edits, 1 gate runs | Worktree grammar, checkpoint preservation, plan validation, restore, and closure tests |
| Lower-tier trials | 6 edits, 1 gate runs | Task definitions, twelve matched runs' evidence projection, comparison, and reviewer decision record |

The remaining order is diagnostics, landing completion, recovery, then trials.
The landing-completion spec owns the original handoff state-file and retrospective fact folds.
No original FT311 fold disappears into an unnamed remainder.
FT254, FT258, FT265, FT290, and FT312 retain their neighboring subjects.
The existing merge lane and the landing gate remain unchanged.
The trial spec owns any measured savings claim and any later recommendation to change tier defaults.

## Further notes

### Decision provenance

The whole ready map and its topic folder moved together into this spec's decisions directory.
Future FT311 specs consume this same source in place.
Before this spec is retired, promote the shared decisions and repair all remaining consumers in the same change.
Do not delete provenance still required by another FT311 spec.
Do not retire the FT311 roadmap row when only this first capability lands.

The three structured map sources were reread on 2026-09-08.
The research report and FT311 roadmap were checked against the current preflight, review phase, closure owners, and tier profile.
The source's historical corpus counts remain observational evidence, not a performance target for this spec.
No outside source or current model-pricing claim is required here.

### Source-sentence-to-row accounting

References below identify resolved tickets in the single Decision source.
Each table entry paraphrases its assigned clause without creating another decision source.

| Source clause | Disposition |
|---|---|
| Ticket 1: existing folds plus charges, review preparation, and repair triage | DP1–DP27 cover this capability. The four priced specs own the other folds. |
| Ticket 2: recorded landing effects | Landing-completion spec |
| Ticket 3: CLI derives facts and delegates interpret | DP1–DP3, DP10, DP15–DP18, DP24 |
| Ticket 4: bounded cheap trials retain guidance and review routes | DP24–DP27 retain existing routes. Trial spec owns cheap trials. |
| Ticket 5: assignment, revision, ticket, fence, evidence, checks, return, and refusal | DP1–DP9, DP19, DP23 |
| Ticket 6: automatic native axes after preparation with separate returns | DP15–DP18, DP20, DP21, DP25 |
| Tickets 7 and 8: cleanup scope and interrupted landing effects | Landing-completion spec |
| Ticket 9: comparable cost with no additional defects or authority violations | Trial spec |
| Ticket 10: five specs with one source | DP27 and Decision provenance |
| Ticket 11: recoverable checkpoint reset | Recovery spec |
| Ticket 12: unknown facts and uncertain identity | DP18, DP23–DP25 here. Landing-completion spec owns draft closeout facts. |
| Ticket 13: native unavailability retains charges and stops dispatch | DP22, DP26 |
| Tickets 14 and 15: twelve runs and complete cost accounting | Trial spec |
| Ticket 16: proposals, approval, dependency edges, and no repeat permission | DP10–DP14, DP19, DP25 |
| Ticket 17: staged bytes and truthful diagnostic selection | Diagnostics spec |
| Tickets 18 and 19: order and final confirmed scope | DP27 and Out of scope |

### Reader sweep and enforcement reads

The repository sweep used hidden-file coverage and included scripts, workflow files, fixtures, and Go tests.
Release-preflight matches describe a different command and are excluded.
These are the relevant readers and their disposition.

| Reader or owner | Read evidence and disposition |
|---|---|
| preflight command, GatherPinned, Decide, and assignmentTarget | Existing phase-entry and closure owners. DP1–DP18 extend their package. |
| tickets.ParseTicket and tickets.Enumerate | Canonical field and file owners. Preserve grammar and reuse complete bytes. |
| tickets.BoundFiles and canary fixture pins | Canonical closure producers. DP10–DP14 consume them without registry copies. |
| diff snapshot and movement-checked read | Existing explicit-pair owner. DP5 and DP15–DP18 compose it. |
| consumers.CommandWithVersion and changed-mode collector | Existing blast owner and citation version injection. DP15–DP18 compose it. |
| coverage.ParseSpec and coverage.Command | Existing row owner. DP2 and DP16 retain complete rows. |
| cmd/bench/main.go and command_registry.go | Preflight registration and version injection. DP6 preserves the existing public inventory meaning. |
| command_registry_test.go and help_inventory_test.go | Exact public inventory consumers. Ticket 3 retains their closure. |
| axi_query_registry_test.go and subcommand_routing_table_test.go | Operational classification and route ownership. No AXI inventory expansion. |
| bin/bench.sh preflight route | Existing argument-forwarding route. No new shell dispatch or wrapper grammar. |
| anchors registry, evaluator, and docs_workflow_helpers_test.go | Existing guidance enforcement. DP19–DP26 use the same executed owner. |
| conformance registry_test.go workflow-guidance family | Registers every family implementation path. Ticket 4 adds the new anchor owner here. |
| workflow-guidance canary fixtures | Existing implement preflight and review fallback needles require current-state updates where changed. |
| craft-delegate and its two delegation references | Canonical return, isolation, and fallback owners. DP19–DP26 remove the contradictory review fallback. |
| build_subject_mode_test.go and preflight legacy tests | Existing build and frozen-pair semantics remain under DP6. |
| diff range diagnostics, historical audits, roadmap details, changelog, and legacy handoff fixture | Their legacy command references remain valid. No rewrite is required. |
| harnesses records and generated phase adapters | Static capabilities are not native runtime evidence. Canonical phase edits flow through existing adapters. |

The public command registry, routing tests, and help fixtures retain the existing command name.
No new shipped-surface claim or phase command is introduced.
The new check is reached through the existing ordinary Go test root.
Guidance canaries reach TestRootConformance through the docs-currency-workflow registry entry and workflow-guidance family.
No independent gate invocation is added inside a test constructor.

### Pre-review proof checklist

- Cited symbols: The command, parser, snapshot, assignment, closure, and rendering owners above were read in the current tree.
- Import edges: preflight already imports diff, tickets, and coverage. New composition of consumers stays below the command registry and avoids conformance imports.
- Source-row clauses and occurrences: The source accounting table enumerates every resolved ticket. No source clause is rewritten by this spec.
- Promised field labels: assignment, base, source_tip, ticket, writes, evidence, checks, return, complete, next.
- Changed-function callers: Existing preflight Command is registered in cmd/bench/main.go and called by its same-package tests. Existing GatherPinned and Decide retain their legacy callers.
- Copy survival: DP3 and DP17 change canonical source inputs. DP26 removes contradictory native-review fallback instructions across the named guidance readers.

### Flagged additions and engineering choices

No new product promise beyond the reviewed map is intended.
The opt-in flag names, structured labels, explicit-pair requirement, and existing-owner composition are spec-author engineering choices.
The coordinator supplement makes the source's task-judgment boundary explicit.
It does not require a new structured approval database or ticket format.
The proposed seams were presented for reviewer confirmation before drafting.

### Implementation ticket approval table

| Ticket | Blocked by | Delivered outcome |
|---|---|---|
| 1. Prepare one build charge | none | One validated selected-ticket charge through the existing CLI |
| 2. Propose ticket ownership closure | 1.md | Read-only fixture, registry, and dependency repair proposals |
| 3. Prepare shared review evidence | 2.md | Complete frozen evidence and independent axis inputs |
| 4. Consume build charges and triage repairs | 3.md | Native build preparation and bounded mid-tier repair proposals |
| 5. Dispatch prepared native reviews | 4.md | Automatic independent axes with capability and return handling |

The serial edges reflect shared preflight and guidance owners.
Each ticket includes its public behavior, focused tests, and required registry closure.
Ticket 3 carries the final preflight invariant.
Ticket 5 carries the final shared-guidance invariant and the implementation review pickup.

### Approval surface

| Subject | Disposition requested |
|---|---|
| Stories and lines | Approve the four outcome groups and their bound model efforts. |
| Seams | Approve public preflight tests and review-owned native exercises. |
| Acceptance and edges | Approve DP1–DP27 and the explicit exclusions. |
| Ownership fences | Approve the exact union above for implementation. |
| Scope and tickets | Approve the five-ticket graph within the first FT311 capability. |
