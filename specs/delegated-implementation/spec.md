# Opt-in delegated implementation

Status: staged

Decision source: Reviewer-confirmed conversation, 2026-09-12, including both grill rounds, confirmation to draft, and the later ticket-authorship and mid-review correction.

Verification log: 3 review iteration(s) — the initial proposal passed after two iterations. One focused review examined the later ticket-author correction. The author folded its limit, distinct-author, and chunk-fence findings; that final fold has no further delegate review.

## Problem

The full workflow binds implementation, chunk verification, and final reconciliation to one session.
A reviewer cannot opt into concurrent ticket authors without contradicting that contract.
Completion evidence cannot distinguish a ticket author from an orchestrator.

## Solution

Add `--delegate` to `$bench-implement-spec --full <spec>` and its existing harness equivalents.
The invoking session selects configured implementation tiers and orchestrates one retained delegate per ticket.
The implementation chunk remains the review scope.
Independent chunks can author concurrently through separate ticket worktrees.
Integration, chunk acceptance, and landing remain serial.
Without `--delegate`, the existing single-author workflow remains unchanged.

## User stories

Line: gpt-5.6-sol / high.
Implementation-line reason: Identity validation is the hardest chunk. Existing seams are known, but source attribution and compatibility need strong synthetic tests.
Harder chunks: DI-C1.
This recommendation applies to implementation of this proposal after approval.
Its own build uses the existing single-author workflow until the proposed capability exists.

1. As a reviewer, I want explicit opt-in, so that ordinary full runs retain their author.
2. As an orchestrator, I want configured tier selection, so that each ticket receives an appropriate model and effort.
3. As an orchestrator, I want concurrent independent authors, so that independent work can proceed together.
4. As a ticket author, I want retained ownership, so that tests, probes, and repairs stay with their implementation.
5. As a reviewer, I want controlled replacement, so that a failed author cannot block all useful progress.
6. As an orchestrator, I want source-bound verification, so that branch-local success cannot certify an integrated change.
7. As a reviewer, I want independent review sessions, so that an author cannot accept its own work.
8. As an orchestrator, I want final verification ownership, so that one session reconciles the complete integrated result.
9. As a reviewer, I want every participant and attempt accounted for, so that delegation cannot hide its cost.
10. As a reviewer, I want costs and unknowns separated, so that estimates cannot become billing claims.
11. As an agent, I want durable identities and evidence, so that interruption cannot silently change authorship or erase failure.
12. As a reviewer, I want synthetic verification, so that this feature requires no paid comparison or model-default change.

## Implementation decisions

### Entry and ownership

The phase accepts `--delegate` only with `--full` and one approved spec with an approved ticket graph.
Missing approval, missing tickets, and invalid flag combinations stop before author dispatch.
This is phase grammar, with no new CLI verb or scheduler service.
The existing preparation owner supplies source facts, fences, dependencies, and charge inputs.
The orchestrator supplies the task-specific model, effort, cap, and probe supplement.

`--delegate` authorizes selection and eligible escalation through all configured tiers, including top.
The invoking harness resolves each tier through its existing binding.
No new ceiling flag or default model change is required.
An unavailable or unbound model stops that dispatch without substitution.
The orchestrator declares an author limit within available harness slots and the user's budget before dispatch.
The limit counts concurrently active ticket writers, not retained idle author sessions.

Native session capacity remains a separate constraint.
Only supported native pause and resume operations can free slots while preserving author identity.
If the harness cannot preserve required authors within its capacity, dispatch stops and reports that limit.

In delegated mode, every review axis uses the invoking harness's configured mid tier at high effort.
This review route applies to cheap, mid, and top implementation tiers alike.
A standing cross-harness review also uses that harness's configured mid tier at high effort.
No-flag runs retain their existing conditional review route.
User budgets, cancellation, and scope decisions retain their stop authority.

The orchestrator owns the integration worktree, plan updates, evidence records, assessment imports, and landing.
Each ticket receives a distinct delegate session.
That ticket author owns its production changes, tests, mutation probes, and repairs through its chunk's acceptance.
Diagnostic helpers receive read-only work.
Existing independent acceptance checks remain verification work and transfer no production authorship.
The orchestrator routes production merge-conflict repairs and later cross-chunk repairs to the relevant author.

### Concurrent work and serial acceptance

The approved completion plan fixes a topological integration order.
Chunks can author concurrently only when their prerequisites have accepted checkpoints and their expected writes do not overlap.
Shared contract changes also require a dependency before concurrent dispatch.
A discovered overlap pauses the affected work until the orchestrator updates its order and charge.
Unrelated worktrees and changes remain outside the run.

Each ticket uses its own Bench assignment and recorded native delegate identity.
A read-only preparation charge obtains that identity before write authority is granted.

Within a chunk, tickets execute serially in dependency order.
The first ticket starts from the accepted prerequisite source.
A later ticket can start from its predecessor's committed green source before the chunk review.
The orchestrator uses existing worktree create-from and merge operations to preserve that source.
It records each delegate identity before sending its write charge.

The orchestrator folds the chunk's clean, committed source through the existing merge owner.
Only one chunk contribution enters the integration source at a time.
Each retained ticket author then performs its required verification on that exact integrated source.
Authors can use their own worktrees after the orchestrator synchronizes them to the integrated tip.
The review covers the full chunk delta after all its tickets are integrated.

Review freezes the predecessor checkpoint tip and the current integrated tip.
The three axes review the whole approved spec and that complete delta.
The existing chunk checkpoint closes acceptance before the next integration step.
A ticket in a dependent chunk starts only after every prerequisite chunk passes its checkpoint.
A same-chunk ticket dependency requires the predecessor ticket's green commit, not an intermediate chunk review.
Authors of independent chunks can continue while another chunk waits for acceptance.

Gate, merge, commit, and integration-check operations respect existing serialization and quiet-tree rules.

Branch-local results remain historical evidence with their original source and performer.
Relabeling their digest or performer cannot satisfy an integrated obligation.
Repairs require current verification and each axis's current result or native reaffirmation.
All ticket authors and recoverable sources remain available for later repairs until the existing landing lifecycle releases them.

### Identity and evidence

Extend the existing completion-plan and review-record owners with an explicit version 2 delegated form.
Version 1 retains its `implementation_session` semantics and acceptance results.
A version 2 plan owns one `execution` declaration with `mode: delegate` and one `run_id`.
It names `orchestrator_session`, `author_limit`, and an ordered assignment history for each ticket basename.
Undispatched tickets can have empty histories.
A dispatched ticket or its verification occurrence requires an effective assignment.

Each assignment names its native session, Bench assignment, selected model, effort, source, and native dispatch reference.
Each replacement names its predecessor, trigger, stopped-writer evidence, and preserved source.
It also includes the reassessment when that trigger requires one.
The latest assignment supplies the ticket's effective author.
One native author session cannot own two different tickets in the run.
The orchestrator and author identities must differ.

The plan is the sole identity authority for acceptance.
The version 2 review record refers to that authority through its existing plan digest.
It does not carry an independently editable `implementation_session` identity.
Mixed versions, missing identities, duplicate assignments, and conflicting ownership refuse acceptance.
An evidence-only edit cannot change who may verify or reconcile completion.

The orchestrator commits the assignment declaration before sending that author's write charge.
Later declarations use existing plan amendments and preserve earlier assignment history.
Unchanged chunk IDs receive identity mappings when the plan digest changes.
Accepted evidence retains its historical plan and source binding.

The current plan's full author history supplies the exclusion set for every review session.
An author selection cannot reuse a session that supplied independent review in this run.
A record cannot downgrade a delegated plan to version 1 acceptance.

Each version 2 chunk-verification requirement names its owning `ticket` in the source-bound plan.
That ticket must belong to the chunk, and every ticket needs at least one verification requirement.
The requirement resolves its author from the frozen plan's ticket assignment.
Final requirements omit `ticket` and retain the orchestrator as verifier.

Historical occurrences remain valid evidence of their original author and source.
New post-replacement occurrences require the successor assignment and fresh verification.
Historical results never satisfy a new source or successor obligation.
Final verification and reconciliation require the orchestrator on the final source.

Version 2 final results use `integration-verification`.
Chunk results retain `author-verification`, and review results retain `independent-review`.
These evidence roles remain distinct from assessment roles.

Each chunk has three distinct review sessions.
Every reviewer differs from the orchestrator and every current or former author in the run.
A reviewer can reaffirm its own axis after repair, but cannot supply another axis for that chunk.
Commands, exits, probe outcomes, restores, source digests, frozen pairs, native excerpts, and excerpt digests remain mandatory.
Missing, pending, skipped, failed, stale, or contradictory evidence cannot close an obligation.

The source-chain, acceptance-row, exact-command, and plan-amendment checks remain in force.
Prospective landing still permits only its existing exact status transform and review-record delta.
New destination content requires integration, review, and verification on the new source.
The installed broker and gate retain execution authority.
A declaration or embedded excerpt cannot authenticate an invented native transcript.

### Replacement and continuation

The orchestrator can replace an author or change its model after one of these triggers:

- Two completed attempts make no progress, followed by a recorded reassessment.
- The author reaches a terminal failure.
- The author exhausts its declared cap.
- The author session is lost.

A terminal failure ends an author assignment with a recorded failure.
An ordinary test red in an active assignment is not a terminal failure.
A planned TDD red, slow tool call, or diagnostic-only action is not a completed failed implementation attempt.

Reassessment states the changed hypothesis and next discriminating check.
The orchestrator confirms the old writer has stopped before assigning the successor.
Without that confirmation, replacement pauses.
A replacement cannot reset the user's budget or discard earlier attempts.

The successor receives preserved source through existing worktree controls and performs fresh verification.
Earlier evidence retains its original author and source.
Replacement invalidates inherited author verification.
An affected accepted chunk reopens through the existing amendment and repair route.

### Accounting

Use one ordinary-work assessment run with the execution declaration's `run_id`.
Add `orchestration` to the existing assessment role vocabulary.
The existing implementation, repair, verification, review, and diagnostic roles retain their meanings.
Record each invocation or coherent attempt when it starts, including failed dispatches.
An unavailable native identity remains explicit unknown evidence with the dispatch reference.

The existing attempts array is the participant and attempt inventory.
At each checkpoint and phase exit, the orchestrator reconciles its native dispatch returns and completion evidence against that inventory.
Include failed, cancelled, superseded, and incomplete work, all review axes, probes, diagnostics, and final verification.
Include orchestrator preparation, integration, and close work.
A native event contributes to only one attempt.
Unknown per-role attribution remains unknown when a native counter cannot separate the work.

Add optional `bench_input_batches` using the existing Bench-input selector type.
This array accepts explicit selectors from several assignments in one run.
The singular `bench_inputs` form remains valid; supplying both forms refuses.
The collection owner enforces assignment identity and unambiguous event-to-attempt mapping across the entire import.
Collect assignment census evidence before release removes it.
Imports remain serialized, orchestrator-owned, and append-preserving.

Reuse existing cost and interval calculations.
Keep estimates, actual charges, currencies, and unknown components separate.
A missing charge is unknown, and a measured zero remains zero.
Concurrent intervals contribute their union to wall time and separate effort to participant time.
Unknown measurements do not block implementation or certify complete accounting.
An omitted known participant prevents a complete-account claim until the assessment update includes it.

## Implementation chunks

Each ticket is one serial green checkpoint and one review chunk for this proposal's build.
The first two capabilities have no data dependency, but this build accepts them in the listed order.

| chunk / ticket | blocked by | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- | --- |
| DI-C1 / `1-bind-delegated-evidence.md` | none | Delegated identities through completion and landing | DI1, DI2, DI3, DI4, DI5, DI6, DI7, DI8, DI9, DI10, DI11, DI12, DI31, DI36, DI37, DI38, DI41 | Evidence and landing fixtures | yes |
| DI-C2 / `2-account-for-participants.md` | none | Complete accounts across assignments | DI13, DI14, DI15, DI16, DI17, DI18, DI19 | Assessment command fixtures | no |
| DI-C3 / `3-enable-delegated-full-runs.md` | 1-bind-delegated-evidence.md, 2-account-for-participants.md | Opt-in orchestration using both capabilities | DI20, DI21, DI22, DI23, DI24, DI25, DI26, DI27, DI28, DI29, DI30, DI32, DI33, DI34, DI35, DI39, DI40, DI42, DI43 | Workflow tripwires and synthetic journey | no |

The build uses this version 1 completion plan until the proposed workflow is delivered.

```bench-completion-plan
{"version":1,"chunks":[{"id":"DI-C1","tickets":["1-bind-delegated-evidence.md"],"verification":[{"id":"identity","command":"bench test --package ./internal/reviewrecord --run Delegated","probe":"omit effective-author comparison"},{"id":"checkpoint","command":"bench test --package ./internal/gate --run Delegated"},{"id":"landing","command":"bench test --package ./internal/landing --run Delegated"},{"id":"preflight","command":"bench test --package ./internal/preflight --run Delegated"}]},{"id":"DI-C2","tickets":["2-account-for-participants.md"],"verification":[{"id":"assessment","command":"bench test --package ./internal/assessment","probe":"omit a selected assignment batch"}]},{"id":"DI-C3","tickets":["3-enable-delegated-full-runs.md"],"verification":[{"id":"workflow","command":"bench test --check docs-currency-workflow","probe":"omit delegated prerequisite-checkpoint instruction"},{"id":"journey","command":"bench test --package ./internal/worktree --run Delegated"}]}],"final_verification":[{"id":"acceptance","command":"bench test"},{"id":"integration","command":"bench test --check system"}]}
```

## Testing decisions

Use the existing review-record, checkpoint, prospective landing, assessment, and worktree seams.
Extend shared fixtures rather than copy their constructors.
No new process seam or external dependency is required.
The enforcement asset records readers and gate attachment.
Test names below are planned additions.

Guidance tripwires prove instruction presence; semantic review grades their meaning.
Synthetic events exercise command owners without launching models.

### Seam diagram

```text
approved invocation -> existing preparation -> isolated ticket authors
 -> serial merge -> author verification + three independent reviews
 -> checkpoint -> orchestrator final verification -> prospective gate and landing
all participants and attempts -> existing assessment record / show / compare
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| DI1 | 1 | Version 1 acceptance matches baseline fixture outcomes | planned TestDelegatedLegacyParity in internal/reviewrecord | Differential success and refusal cases detect default identity drift |
| DI2 | 6, 11 | Delegated verification ownership resolves from the source-bound plan | planned TestDelegatedIdentityAuthority in internal/reviewrecord | An evidence-only performer rewrite cannot change the required author |
| DI3 | 6, 11 | An invalid delegated identity declaration refuses acceptance | planned TestDelegatedIdentityRefusals in internal/reviewrecord | Missing, duplicate, conflicting, and mixed-version fixtures cannot pass |
| DI4 | 6 | A chunk checkpoint requires each ticket's assigned-author verification on the integrated source | planned TestDelegatedChunkVerifier in internal/gate | Orchestrator, foreign-author, and branch-only passes fail the valid obligation |
| DI5 | 7 | A current or former author cannot review any chunk | planned TestDelegatedAuthorReview in internal/gate | Replaced authors and another ticket's author each fail |
| DI6 | 7 | The orchestrator cannot supply independent review | planned TestDelegatedOrchestratorReview in internal/gate | A complete orchestrator-authored axis still refuses |
| DI7 | 7 | One session cannot supply two axes for a chunk | planned TestDelegatedDistinctAxes in internal/gate | Repeated reviewer identities cannot satisfy the axis inventory |
| DI8 | 8 | Final verification requires the orchestrator on the final source | planned TestDelegatedFinalVerifier in internal/landing | A ticket author or older source blocks publication |
| DI9 | 8 | Missing reconciliation blocks delegated landing | planned TestDelegatedReconciliation in internal/landing | Omitting a required row leaves the destination ref unchanged |
| DI10 | 8 | Unreviewed destination content blocks delegated landing | planned TestDelegatedDestinationDelta in internal/landing | Destination additions and deletions retain their refusal |
| DI11 | 5, 11 | Replacement retains historical occurrences under their original frozen assignments | planned TestDelegatedReplacement in internal/reviewrecord | Comparing an older occurrence with the latest author would reject valid history |
| DI12 | 6, 7 | Delegated checkpoints retain source and terminal-evidence refusals | planned TestDelegatedEvidenceRefusals in internal/gate | Missing axes, failed probes, stale plans, chain gaps, and altered excerpts stay red |
| DI13 | 9 | Assessment reports orchestration as its own role | planned TestAssessmentOrchestration in internal/assessment | Omitting the role loses the orchestrator from show and comparison |
| DI14 | 9 | One import collects explicit evidence from multiple assignments | planned TestAssessmentAssignmentBatches in internal/assessment | Two assignments expose a singular-only collector |
| DI15 | 9 | Conflicting event ownership across batches refuses import | planned TestAssessmentBatchConflicts in internal/assessment | A foreign assignment or multiply owned event cannot count |
| DI16 | 10 | Reports separate estimates, actual charges, and unknowns | planned TestAssessmentDelegatedCosts in internal/assessment | Unequal rates, measured zero, and missing charges expose false totals |
| DI17 | 9, 11 | Updates preserve failed and replaced attempts | planned TestAssessmentDelegatedHistory in internal/assessment | Removing a failed author refuses instead of lowering the total |
| DI18 | 10 | Concurrent intervals preserve wall and effort time semantics | planned TestAssessmentDelegatedIntervals in internal/assessment | Fixed overlapping intervals expose summed wall time |
| DI19 | 1 | Singular assessment imports retain baseline results | planned TestAssessmentSingularParity in internal/assessment | Differential fixtures detect legacy collection and cost drift |
| DI20 | 1 | The phase delegates only an approved delegated full run | review-owned: grammar scenarios and workflow tripwires | No-flag, missing-spec, unapproved, and delegate-without-full cases prohibit dispatch |
| DI21 | 2 | The orchestrator selects configured tiers without changing defaults | review-owned: routing scenarios and line-routing tripwires | Top works under opt-in while missing bindings prohibit substitution |
| DI22 | 3 | Authors of independent chunks use separate assignments within the declared limit | review-owned: charges and synthetic journey | Retained idle sessions cannot exhaust the active-writer limit; excess active writers cannot dispatch |
| DI23 | 3, 6 | Concurrent authors reach acceptance through a serial integrated source chain | planned TestDelegatedIntegrationJourney in internal/worktree | A branch-only pass cannot close the second integrated chunk |
| DI24 | 3 | A dependent author waits for every prerequisite checkpoint | review-owned: dispatch scenarios and workflow tripwires | A pending predecessor review prevents dependent dispatch |
| DI25 | 4 | Repairs return to the retained ticket author | review-owned: repair and merge-conflict charge scenarios | An orchestrator production repair violates ownership |
| DI26 | 5 | Two completed no-progress attempts permit replacement or model change only after reassessment | review-owned: replacement scenarios and workflow tripwires | One attempt or absent reassessment cannot trigger transfer |
| DI27 | 9 | Phase exits reconcile known invocations against assessment attempts | review-owned: synthetic dispatch and attempt inventories | An omitted failure, diagnostic, review, or orchestrator interval prevents a complete-account claim |
| DI28 | 11 | Resumption retains identities, source pins, and pending obligations | review-owned: handoff scenarios and workflow tripwires | Resumption cannot invent authors or accept pending evidence |
| DI29 | 12 | Verification launches no paid comparison and changes no model default | review-owned: synthetic journey and complete diff | A trial launch or binding edit violates scope |
| DI30 | 8, 11 | The existing landing owner alone publishes implemented status | planned TestDelegatedIntegrationJourney in internal/worktree | Final verification alone cannot publish completion or release unrelated work |
| DI31 | 3, 11 | An undispatched future ticket accepts an empty assignment history | planned TestDelegatedPendingAssignments in internal/reviewrecord | Requiring all native sessions up front fails a bounded-author fixture |
| DI32 | 5 | A recorded terminal author failure permits replacement or model change | review-owned: terminal-failure scenario and workflow tripwire | An implementation that supports only no-progress retries would block this trigger |
| DI33 | 5 | Exhausting the declared author cap permits replacement or model change | review-owned: cap-exhaustion scenario and workflow tripwire | Ignoring cap exhaustion would prevent authorized transfer |
| DI34 | 5 | A lost author session permits replacement or model change | review-owned: lost-session scenario and workflow tripwire | Requiring another response from the lost session would block transfer |
| DI35 | 5, 11 | Every author transfer requires confirmed termination of the old writer | review-owned: transfer scenarios and workflow tripwire | A valid trigger with an unconfirmed stop still prohibits successor writes |
| DI36 | 5, 6 | Post-replacement acceptance requires the successor's fresh verification | planned TestDelegatedReplacementFreshness in internal/reviewrecord | Relabeling the predecessor pass cannot satisfy the successor obligation |
| DI37 | 4, 6 | One chunk accepts verification from its distinct ticket authors under their own obligations | planned TestDelegatedTicketOwners in internal/reviewrecord | Requiring a single author for the whole chunk rejects a valid two-ticket fixture |
| DI38 | 4, 6 | An invalid ticket-to-verification owner mapping refuses chunk acceptance | planned TestDelegatedTicketObligations in internal/reviewrecord | A foreign ticket, an uncovered ticket, or another ticket's performer cannot satisfy the map |
| DI39 | 7 | Delegated review uses configured mid tier at high effort for every implementation tier | review-owned: cheap, mid, and top author scenarios plus routing tripwire | Retaining the mid-author-to-top-review branch violates the opt-in rule |
| DI40 | 3, 4 | Same-chunk successor tickets can start after their predecessor's green commit | review-owned: two-ticket chunk scenario plus workflow tripwire | Waiting for an intermediate chunk checkpoint deadlocks the chunk |
| DI41 | 4 | One native author session cannot own two ticket histories | planned TestDelegatedDistinctTicketAuthors in internal/reviewrecord | Assigning the same session to two tickets must refuse even with valid verification |
| DI42 | 3, 4 | A same-chunk successor waits while its predecessor commit is pending or red | review-owned: two-ticket chunk scenario plus workflow tripwire | Early dispatch violates the green ticket dependency |
| DI43 | 7 | Chunk acceptance review begins only after every ticket reaches the integrated chunk tip | review-owned: two-ticket chunk scenario plus workflow tripwire | Per-ticket review cannot substitute for the complete chunk review |

### Edge inventory

Both kit and linked-repository callers use the existing supported harnesses.
DI20 covers absent approval, missing specs, absent or empty tickets, and invalid flags.
DI21 covers unavailable and unbound models without substitution.
DI3 covers malformed declarations, unknown versions, duplicate keys, and unsupported fields.
DI4 and DI12 cover missing verification, stale sources, unsuccessful probes, failed restores, and nonterminal reviews.

DI11, DI26, and DI31–DI36 cover pending assignments, replacements, exact transfer triggers, and unavailable stop proof.
DI22–DI25 cover overlapping writes, shared contracts, dependency failure, and merge conflicts.
DI37–DI43 cover ticket ownership, review tier, same-chunk dependencies, and the full-chunk review fence.

DI14–DI19 cover empty batches, both input forms, duplicate selectors, foreign assignments, missing counters, and append conflicts.
DI28 covers cancellation, interruption, stale handoffs, and unavailable native evidence.
Existing bounded, regular-file, and non-symlink checks remain attached to both assessment input forms.
Records never discover input paths or execute reference text.

## Ownership fences

- `internal/reviewrecord`
- `internal/gate/delegated_checkpoint_test.go`
- `internal/landing/delegated_completion_test.go`
- `internal/preflight/delegated_evidence_test.go`
- `internal/assessment`
- `.bench/BENCH.md`
- `.bench/BENCH-reference.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/commands/bench-final-check.md`
- `.agents/skills/bench-craft-line/SKILL.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `.agents/skills/bench-craft-tickets/SKILL.md`
- `CONTEXT.md`
- `projects/benchkit.md`
- `docs/adr/0021-benchmark-workflow-orchestration.md`
- `docs/field-guide.html`
- `CHANGELOG.md`
- `internal/anchors/registry_retained_workflow.go`
- `internal/conformance/retained_workflow_test.go`
- `internal/conformance/docs_workflow_helpers_test.go`
- `internal/worktree/delegated_integration_test.go`
- `tests/canary/workflow-guidance-anchors`
- `reviews/delegated-implementation.md`

- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`

- `tests/canary/docs-currency-token-diet`
- `tests/canary/load-validity-metadata`
- `tests/canary/skills-index-command-adapters`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`

Reviewer disposition: Proposed fences await the spec-and-ticket approval.
No implementation fence authorizes a write during this spec phase.

## Out of scope

Won't handle: independent branch acceptance — concurrent authors use the serial integrated checkpoint chain.
Won't handle: a new scheduler or harness adapter — orchestration uses existing native delegation and worktrees.
Won't handle: paid comparisons or automatic default adoption — ordinary synthetic assessment remains in scope.
Won't handle: a broader workflow rewrite — no-flag runs retain their existing authorship contract.

These exclusions grant no future implementation authority.
No follow-on capability is scheduled or priced by this proposal.

## Further notes

The first two tickets deliver usable evidence and assessment capabilities before the phase exposes the flag.
No superseded workflow spec remains to retire.
Implementation changes ADR 0021 only for the opt-in exception.
Existing model bindings and default decisions remain closed.

The [enforcement read](assets/enforcement.md) records readers, source-to-row mapping, and pre-review checks.
The schema additions cover delegated identities and multi-assignment assessment collection.
The assessment vocabulary gains orchestration only.
Estimates describe proposed work, not incurred charges.
