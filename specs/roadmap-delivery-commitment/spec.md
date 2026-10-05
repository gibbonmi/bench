# Protect the delivery commitment and close verified work

Status: staged

Decision source: `specs/roadmap-delivery-commitment/decisions/roadmap-delivery-commitment.md`, the ready map approved on 2026-10-04.

Verification log: 2 iteration(s) to accept the slices — independent GPT-6.1 Sol xhigh review accepted the spec in one round and the slices in two rounds.
The author corrected ordinary CLI trust wording after spec review. Slice repairs added the charge-evidence fixture and root shift help to their ticket fences.
The final accounting correction adds the root command package check to DC-C4. Both independent reviews are accepted; routine engineering and breakdown decisions were user-delegated.

## Problem

A useful finding can replace previously selected work without a separate displacement decision. The current drain rewrites the recommended sequence on each pass.
Its sequence rule lives in `.agents/commands/bench-drain.md`, under `Refresh the sequence`.
A completed row can also survive its implementation and spec retirement. The final-check command currently assigns roadmap closure to a later drain.

The current flow report shows `target_met=false`, with 22 opened rows, 7 retired rows, and a net increase of 15.
This report describes three addition-bearing commits. It does not prove that the new findings are unnecessary.
The retained diagnosis verifies the growth and the separation between discovery, selection, and completion.

## Solution

Protect one finite milestone and its ordered outcomes through Bench commands. Keep assessed findings outside that commitment until the reviewer explicitly admits them.
Require a separate decision for displacement, activation, switching, and parallel outcomes. Close completed roadmap references in the verified delivery publication.
Complete the milestone only after evidence verifies its stated outcome.

This is one coherent prerequisite for the approved quality milestone. Internal checkpoints do not create more roadmap work or separate feature releases.
The first milestone covers FT376, FT373, and FT349, in that order. Its implementation follows this prerequisite.

## User stories

Line: opus / medium for tickets 05 through 10, one fresh author for each ticket.
Tickets 01 through 04 used the current Codex session with unknown runtime model and effort metadata.

Implementation-line reason: the broker closure is the hardest part, because its exact transform and publication proof must agree. The guidance changes also invoke the leverage override. Existing command, transaction, and gate seams make the design testable.
Harder chunks: assign the broker closure and admission concurrency at ticket slicing.

### Protect a finite commitment

1. As the reviewer, I want one active milestone with ordered outcomes, so that the project has a finite delivery target.
2. As the reviewer, I want future milestones to remain planned, so that they do not compete with the active milestone.
3. As the reviewer, I want each change to name delayed and removed work, so that I can approve its actual cost.
4. As the reviewer, I want approval bound to one exact proposal, so that later edits cannot reuse that approval.
5. As the worker, I want new findings retained as uncommitted work, so that discovery does not displace delivery.
6. As the reviewer, I want planning commands to preserve committed obligations, so that a drain cannot bypass protection.
7. As the worker, I want actionable refusals, so that I know which decision permits the requested work.

### Admit and finish the right work

8. As the worker, I want only eligible committed outcomes to start, so that the agreed order governs delivery.
9. As the worker, I want a recorded blocker to preserve its obligation, so that the next independent committed outcome can proceed.
10. As the reviewer, I want one active outcome by default, so that parallel work needs my explicit direction.
11. As the orchestrator, I want ticket authors to share their outcome binding, so that existing delegation remains available.
12. As the worker, I want every supported start route to check the commitment, so that an alternate command grants no exemption.
13. As the worker, I want commands to recheck current authority, so that a stale assignment cannot publish displaced work.
14. As a planner, I want decision maps and staged specs to remain publishable, so that planning and adoption can proceed.
15. As the reviewer, I want a planning assignment to refuse production changes, so that its purpose is not a bypass.
16. As the worker, I want claims to survive interruption, so that a restart cannot start a competing outcome.

### Close verified delivery

17. As the worker, I want a verified delivery to close its completed roadmap row, so that the next session sees remaining work.
18. As the worker, I want closure to remove its sequence references, so that a retired spec is not the next command.
19. As the reviewer, I want incomplete obligations retained, so that partial implementation cannot close the whole row.
20. As the reviewer, I want closure and delivery in one publication, so that neither can succeed alone.
21. As the worker, I want interrupted closure to resume safely, so that recovery cannot publish twice or lose obligations.
22. As the reviewer, I want milestone evidence beyond an empty work list, so that closure means the stated outcome holds.
23. As the reviewer, I want unmet milestone criteria to remain open, so that a green code gate cannot hide missing delivery.
24. As the worker, I want the active commitment to supply the next action, so that an old recommendation cannot change priorities.

### Adopt without losing authorized work

25. As a project owner, I want an explicit initial commitment, so that Bench never adopts an old sequence as approval.
26. As the worker, I want existing authorized runs to have a concrete continuation route, so that adoption preserves their approved scope.
27. As the project owner, I want the same command checks in a linked project, so that protection does not depend on kit-only tests.
28. As the reviewer, I want this prerequisite to ship before activation, so that adoption does not require an unavailable command.
29. As the reviewer, I want the named quality work to follow this prerequisite, so that enforcement does not become a new expanding milestone.
30. As the reviewer, I want defects before refactors before features when proposing work, so that classification serves the approved purpose.

### Preserve authority under hostile inputs

31. As the worker, I want malformed or unreadable commitment input refused, so that missing evidence grants no authority.
32. As the worker, I want exact path and argument handling, so that spaces or flag-like values cannot change the requested operation.
33. As the reviewer, I want concurrent changes to serialize, so that two commands cannot each claim the only available slot.
34. As the reviewer, I want the installed broker to enforce publication, so that candidate code cannot approve its own change.
35. As the maintainer, I want one owner for each policy fact, so that commands and their projections cannot drift.
36. As the project owner, I want local records documented, so that their contents and retention are explicit.

## Implementation decisions

### Audience and authority

The behavior applies to each Bench repository, including projects that link the kit. Linked worktrees share the repository's commitment and runtime claims.
Independent Git clones remain separate Bench repositories. This feature adds no remote scheduler or coordination service.

The approval command records explicit reviewer direction. The agent must obtain that direction before it invokes the command.
A decision reference identifies that direction for review. It does not authenticate a human or grant security against a malicious same-user process.
This preserves the trust model in `SECURITY.md` and the command scope in `.bench/BENCH.md`.

The selected ordinary CLI reads authority from the resolved default branch and the repository's local receipt store. Its execution is an operator trust assumption.
Candidate policy files cannot provide approval. Ordinary CLI selection is not authenticated by the landing manifest.
The installed promotion broker checks the same authority before publication. A candidate can contain the proposed state, but cannot authorize that state itself.

The operator-trusted installed wrapper resolves its own installation manifest before it launches the broker.
It checks the executable path, version, and digest through the existing `land_route` implementation.
The authenticated broker authorizes the candidate gate subject before publication. Candidate policy bytes supply no executable or receipt authority.

The new start checks run inside the selected ordinary CLI before the adapter or candidate workload executes.
Its existing route accepts `BENCH_RUN_BINARY` or a repository build. Only publication uses the authenticated installed-owner route described above.
DC25 and DC63 retain marker-based proof at these two execution points.

### One commitment owner

Add `internal/commitment` as the owner of policy validation, proposal effects, admission, planning scope, and completion decisions.
Its decisions consume immutable values. Its repository adapter reads the default-branch policy and composes with `intent.Transact` for shared local records.
The domain imports no roadmap, worktree, landing, status, or dashboard package. Those consumers pass facts or read its projection.
A thin command adapter composes the domain with existing repository readers.

Reuse `internal/jsonfile`, `internal/bounds`, `internal/git`, `internal/toon`, and the existing intent transaction.
Do not add a second JSON decoder, roadmap parser, lock protocol, command inventory, or approval ledger inside a consumer.
The roadmap owner remains the source of row identity, detail ownership, occurrence fields, and sequence spans.

The tracked policy is `.bench/commitment.json`. It is project-owned data, not a copied kit payload.
The policy carries a version, milestone identities, outcome criteria, ordered outcome identities, source references, dependency references, and exact parallel grants.
It also carries one active milestone reference and broker-authored delivery facts. Planned milestones have no admission effect.

The policy does not copy the text of roadmap detail files. Each outcome binds its source references and their approved content identities.

An outcome can own several roadmap rows. A row has at most one outcome owner in a commitment.
An outcome can instead name an approved deliverable without a roadmap row. This permits the bounded prerequisite and ordinary project work.
Each delivery binding names one approved spec or one tickets-only folder and the obligations that it completely satisfies.
A partial spec cannot declare the whole row satisfied merely because its `Roadmap:` field names that row.

The local intent record carries proposal receipts, assignment bindings, active outcome claims, blockers, and narrow legacy continuations.
These records are repository-scoped and use the existing transaction lock. They do not form a second source for milestone order or outcome scope.
The schema change preserves old assignment records as inventory. Their presence alone supplies no delivery approval.

### Explicit proposals and approvals

The new command family has these operations. Its form table owns help and dispatch together.

| operation | behavior |
| --- | --- |
| `bench commitment show` | Show adoption state, the active milestone, active or blocked outcomes, and the next eligible committed outcome. |
| `bench commitment inventory` | Show current roadmap obligations, staged deliverables, and existing run identities for an adoption proposal. |
| `bench commitment plan --input <file>` | Validate the proposed policy against current authority and print its exact effect set and plan identity. |
| `bench commitment approve --plan <id> --decision <reference> --delayed <ids-or-none> --removed <ids-or-none>` | Record explicit direction for that exact plan and stage its policy and sequence projection in the owning planning worktree. |
| `bench commitment start --outcome <id> --request <request> --deliverable <path>` | Bind an existing owned assignment to the eligible outcome and claim it atomically. |
| `bench commitment block --outcome <id> --reason <text>` | Preserve the active obligation with a blocker and release its active slot. |
| `bench commitment unblock --outcome <id>` | Clear its blocker without displacing another active outcome. |
| `bench commitment verify --milestone <id> --evidence <file>` | Validate outcome evidence and issue a receipt for a later explicit milestone-completion proposal. |

These are executable Bench operations, not new harness phases. All grammar answers remain non-interactive.
Use exit 0 for success or an identical replay, 1 for unsatisfied intent, and 2 for invalid grammar.
Refusals name the observed commitment state, the blocked action, and the exact applicable command token.
No refusal offers a force flag or suggests that a drain approval permits displacement.

The plan identity binds the predecessor policy identity, proposed bytes, affected source identities, and the computed effects.
Effects include added, reordered, delayed, removed, activated, switched, and parallel-authorized outcomes.
A new earlier outcome delays each unfinished predecessor that it passes. Removal names each removed unfinished outcome.
A switch names the unfinished outcomes that stop receiving work. The explicit delayed and removed operands must exactly match the computed sets.

Approval requires a distinct plan operation and a nonempty decision reference. Findings, a routine drain receipt, or an arbitrary candidate marker cannot replace it.
Identical approval replay changes nothing. A changed proposal, protected source, or predecessor policy invalidates the receipt.
An unrelated default-branch advance does not invalidate approval when every bound authority fact remains identical.
The publisher still checks its exact destination revision through the existing landing protocol.

Approval stages tracked changes in its Bench worktree. Normal lane and landing checks publish them.
The receipt permits only its exact policy transition and derived sequence. A raw policy edit without that receipt refuses at commit and landing.

### Planning and delivery admission

A newly created assignment without an outcome binding is a planning assignment. Inspection, diagnosis, and planning remain available before adoption.
The planning scope permits named decision artifacts, staged spec artifacts, research documents, capture records, and roadmap maintenance.
It also permits glossary or durable decision promotion that the named planning artifact requires. These remain documentation, not executable workflow changes.

The commitment owner declares this path classification once. Commit and landing consume that declaration.
Production source, executable files, gate configuration, and shared agent instructions are outside planning scope.

A planning commit can add uncommitted roadmap rows. It cannot delete, rename, or change a protected row's obligation without exact approval.
Occurrence evidence can change through the roadmap parser's existing event fields. It cannot change protected requirement prose under an event label.
The protected comparison also covers row ownership and the recommended sequence. A planner cannot reset the sequence to an unrelated feature.

A staged spec can commit and land before delivery starts. Its approved source remains visible to the later start command.
Add `--plan-only` to build preflight for this authoring check. It validates the spec and tickets without granting delivery admission.

Ordinary build preflight requires a current outcome binding. Its refusal names `bench commitment start` or `bench commitment plan`, as applicable.
The plan-only result cannot satisfy an implementation admission check, commit a production change, or authorize delivery landing.

The start command checks assignment identity, current policy, deliverable identity, membership, dependencies, order, blockers, and existing claims.
It records the assignment binding and the outcome claim in one intent transaction. A failure leaves neither record behind.
Multiple ticket assignments within one outcome share that claim. They do not count as parallel outcomes.

A sibling created from an already-bound integration source inherits only that same outcome and its approved deliverable.
An arbitrary `--from` value cannot grant a different outcome.

The first eligible unfinished outcome can start. A blocked predecessor stays in the commitment while the next independent outcome becomes eligible.
An unmet dependency still blocks its dependent outcome. An unrelated roadmap row remains ineligible, even when every committed outcome is blocked.
Unblocking does not interrupt a later active outcome. The earlier obligation returns to normal order after that active outcome finishes or blocks.

A second outcome requires a grant naming those exact parallel outcomes. A numeric limit alone does not authorize an arbitrary pair.
Two racing starts cannot each occupy the default single slot. A process exit, released worktree, or expired lease never proves outcome completion.
A blocked or interrupted outcome retains its obligation and recovery information until an authorized transition resolves it.

The command adapters apply the owner at these operation points:

| route | enforcement point |
| --- | --- |
| `bench commitment start` | Before an assignment gains delivery authority. |
| `bench worktree create --from <target>` | Before it inherits a delivery binding from the source assignment. |
| `bench shift --outcome <id> <objective>` | Before refresh, intent publication, worktree acquisition, or adapter execution. |
| `bench preflight build <slug>` | Before it reports delivery readiness or prepares a build charge. |
| `bench commit` and its `--preflight-build` form | Before formatting, lane execution, or publication of production changes. |
| `bench worktree land` and `--resume` | Before authorization and publication, with current authority checked again after the gate. |

Usage and existing identity refusals keep their precedence. Once identity is established, commitment refusal precedes the route's mutable effects.
Inspection, `--help`, planning, and the existing recovery commands remain reachable. A cleanup command cannot mark an outcome delivered.

### Verified delivery and roadmap closure

Extend the installed broker's existing prospective transform. The transform includes the exact spec-status change or tickets-only close, the satisfied obligations, and their roadmap closure.
It removes completed row headings and matching detail owners together. It also removes their recommended sequence references and satisfied dependency references.
Every unaffected row, detail body, and relative order remains unchanged.

The broker derives closure from the approved outcome binding and current completion evidence. It never derives completion from deleted rows or an agent's finished label.
For a spec, retain the existing complete checkpoint, acceptance reconciliation, and final integration evidence.
For a tickets-only delivery, retain its approved ticket acceptance and whole-project gate, with a closure binding that names the fully satisfied obligation.
A delivery without roadmap owners updates only its commitment obligation.

The gate grades the exact prospective delivery and closure tree. Its completion comparison permits only the broker-derived transform and the existing review-record allowance.
An omitted row removal, surviving sequence reference, extra closure, or unrelated changed byte refuses. A general path allowlist is insufficient.
The authorizing owner, graded candidate, and published tree remain separate identities.

A verified publication marks its satisfied outcomes delivered in the tracked policy. The fact refers to the reviewed source and retained completion evidence.
It does not embed the publication's own commit or tree identity into itself. The landing result binds those identities after publication.
The broker reconciles local claims from that published fact. A local write failure cannot make the verified obligation open again.

Before publication, a failure or interruption changes neither the default branch nor its roadmap. After publication, the existing resume route completes pending local effects.
A resume verifies the original publication and performs no second delivery. Changed authority still refuses an unpublished retry.

This closure also applies to explicitly authorized legacy deliveries. Uncommitted or partially satisfied work stays open.
Spec retirement can still promote durable content and delete a delivered spec. It no longer schedules completed roadmap cleanup for a later drain.
FT283 and FT284 keep any requirements that this delivery does not satisfy. This spec does not retire either row by association.

### Milestone verification

Each milestone has explicit outcome criteria before activation. Criteria have stable identities and observable statements.
The verification input provides one result and native evidence reference for each criterion at one exact delivered default-branch revision.
It also names that revision's retained green gate evidence. Each result is `verified`, `unmet`, or `blocked`.
Unknown, missing, duplicate, stale, or contradicted results cannot produce a completion receipt.

The command validates evidence identity, criterion coverage, revision, and terminal results. It does not claim to infer the truth of arbitrary prose.
The reviewer assesses the native evidence and explicitly approves the completion proposal. This is the semantic outcome-verification step.
A green code gate, absent rows, or the word `verified` without criterion evidence is insufficient.

A successful verification receipt binds the milestone criteria, evidence, and examined revision. The completion plan and approval consume that exact receipt.
All required delivery obligations must already be satisfied. A completion receipt cannot silently add, remove, or weaken a criterion.
Completing one milestone does not activate a planned successor. Activation or switching remains an explicit approved transition.

### Initial adoption and the quality milestone

An absent policy reports `adoption-required`. An empty, malformed, unsupported, or unreadable policy reports its own refusal.
Neither state infers approval from the old recommended sequence. Help, inventory, planning, and adoption operations remain available.
New delivery starts refuse until initial adoption. Existing assignments remain inspectable and recoverable.

The inventory identifies each existing run and its approved deliverable or explicit scope evidence. The initial proposal lists the already-authorized runs allowed to finish.
Approval binds each continuation to its assignment identity, request, and existing scope. An unlisted old assignment receives no permission.
A listed run can resume and land that scope through the same verification and closure owner. It cannot start a new outcome or expand its scope.

A continuation ends when its scope is delivered or when its run is no longer active, so a finished or released run holds no slot. Existing active continuations occupy the active-work allowance. New delivery waits unless the reviewer explicitly grants additional parallel work. A parallel grant that admits new delivery beside a continuation names that continuation's assignment identity.

The prerequisite's initial implementation uses the pre-adoption installed owner. The final source does not require an already-active policy to install its own admission capability.
After publication, its adoption proposal inventories actual delivery state. This creates the first real commitment through the new commands.
This bootstrap is the existing-run continuation route, not a permanent kit exception or a feature flag that disables protection.

For this repository, the approved successor contains FT376, FT373, and FT349, in that order.
The implementation verifies delivered obligations before it proposes the remaining ordered list. A stale row is not proof of remaining work.
The September 29 survey rows, including the staged FT358 spec, remain uncommitted intake. This prerequisite does not implement or reslice that spec.

The proposal generally places confirmed defects before refactors, then features. Literal dependencies and explicit approved order govern execution.
Class labels do not override the actual purpose of the work. Minimal support required by a committed fix or refactor stays within that outcome.
New findings remain uncommitted. No periodic score, age rule, drain, or severity label can displace the approved list automatically.

### Projections and shared guidance

The commitment owner supplies the current order and eligibility facts. Roadmap, status, and dashboard projections consume those facts.
They show blocked obligations and uncommitted findings separately from the next eligible work. An all-blocked milestone shows its blockers and required decision.
The legacy sequence is visible as unapproved input before adoption. It is never an active commitment.

After adoption, the roadmap owner renders the recommended sequence from the committed policy. Approval and closure update that derived section.
A checked-in sequence that disagrees with the projection refuses a Bench commit or landing. It cannot become a second priority source.
The roadmap context command includes commitment facts without changing capture ownership or the flow calculation.

Shared guidance routes phase starts through admission and sends intake to uncommitted work. It removes unconditional implement-now and sequence-rewrite directions.
The final-check guidance names verified closure as part of delivery. Drain retains reconciliation as a backstop for historical or residual work.
The canonical rule lives in the operating guide. Other phases invoke its commands without restating the policy.

Update the local-data inventory for approval references, proposal identities, run bindings, blocker text, and native evidence references.
Use the existing intent record location, permissions, and atomic transaction. Published policy and source references remain reviewable project data.
Unconsumed receipts remain available for recovery. Consumed receipts can be removed after their publication is verified, without removing live claims or published facts.
No network call, hosted service, new third-party dependency, or transcript archive is required.

## Implementation chunks

The independent spec review accepted the design in one round. Independent slice review accepted this graph in two rounds.
Each ticket is one green checkpoint on the retained integration source. All checkpoints publish together as the bounded prerequisite.
No intermediate checkpoint is a separately adopted enforcement release. The current installed owner governs this build until the complete prerequisite publishes.

The chain is serial because later commands consume earlier owner contracts and share their state files. Each seam-producing ticket has its own review chunk.
DC-C7 groups two independent delivery consumers after the closure seam review. Ticket 08 follows ticket 07 for shared-write order and its rowless delivery facts.
Every ticket must be verifiable while its successors remain unbuilt. No successor test is credited to an earlier coverage owner.

| ticket | blocked by | delivered outcome | chunk |
| --- | --- | --- | --- |
| `01-plan-exact-policy-changes.md` | `none` | Plan and approve exact commitment changes | DC-C1 |
| `02-admit-committed-outcomes.md` | `01-plan-exact-policy-changes.md` | Admit only eligible committed outcomes | DC-C2 |
| `03-protect-planning-and-commits.md` | `02-admit-committed-outcomes.md` | Protect planning commits and build charges | DC-C3 |
| `04-bind-worktrees-and-shifts.md` | `03-protect-planning-and-commits.md` | Bind worktrees and shifts to one outcome | DC-C4 |
| `05-authorize-current-publication.md` | `04-bind-worktrees-and-shifts.md` | Authorize publication against current commitment | DC-C5 |
| `06-close-verified-spec-delivery.md` | `05-authorize-current-publication.md` | Close verified spec delivery in its publication | DC-C6 |
| `07-close-light-delivery.md` | `06-close-verified-spec-delivery.md` | Close verified light-path and rowless delivery | DC-C7 |
| `08-verify-milestone-outcomes.md` | `07-close-light-delivery.md` | Verify milestone criteria before completion | DC-C7 |
| `09-project-commitment-guidance.md` | `08-verify-milestone-outcomes.md` | Project commitment state through readers and guidance | DC-C8 |
| `10-qualify-linked-adoption.md` | `09-project-commitment-guidance.md` | Qualify installed adoption and prepare the finite milestone | DC-C9 |

The ticket `Covers:` rows are the chunk's acceptance inventory. Each row has exactly one ticket owner.
The commands in the completion plan are the checkpoint evidence. Ticket acceptance adds the exact refusal, side-effect, or publication observation to record.

The version 2 execution plan assigns one fresh author to each ticket from ticket 05, before dispatch.
The run retains one integration source and permits one active author.
Independent review uses Claude Opus at high effort, as the reviewer requested.
Ticket 08 records its intent package check separately from ticket 07.

Ticket 02 reuses the spec owner for tickets-only classification and folder identity.
The landing API delegates to that owner. Ticket 07 consumes the same contract for closure.
Its repair coverage includes changed and deleted deliverables, plus both supported deliverable types.

```bench-completion-plan
{"version":2,"chunks":[{"id":"DC-C1","tickets":["01-plan-exact-policy-changes.md"],"verification":[{"id":"commitment","command":"bench test --package ./internal/commitment","ticket":"01-plan-exact-policy-changes.md"},{"id":"intent","command":"bench test --package ./internal/intent","ticket":"01-plan-exact-policy-changes.md"},{"id":"roadmap","command":"bench test --package ./internal/roadmap","ticket":"01-plan-exact-policy-changes.md"},{"id":"bench","command":"bench test --package ./cmd/bench","ticket":"01-plan-exact-policy-changes.md"},{"id":"jsonfile","command":"bench test --package ./internal/jsonfile","ticket":"01-plan-exact-policy-changes.md"}]},{"id":"DC-C2","tickets":["02-admit-committed-outcomes.md"],"verification":[{"id":"commitment","command":"bench test --package ./internal/commitment","probe":"Ignore deliverable validation errors. TestCommitmentStartPublishedIdentity must fail on the reported admission result and preserve an exact source restore.","ticket":"02-admit-committed-outcomes.md"},{"id":"intent","command":"bench test --package ./internal/intent","ticket":"02-admit-committed-outcomes.md"},{"id":"spec","command":"bench test --package ./internal/spec","ticket":"02-admit-committed-outcomes.md"},{"id":"landing","command":"bench test --package ./internal/landing","ticket":"02-admit-committed-outcomes.md"},{"id":"bench","command":"bench test --package ./cmd/bench","ticket":"02-admit-committed-outcomes.md"},{"id":"conformance","command":"bench test --package ./internal/conformance","ticket":"02-admit-committed-outcomes.md"}]},{"id":"DC-C3","tickets":["03-protect-planning-and-commits.md"],"verification":[{"id":"commit","command":"bench test --package ./internal/commit","ticket":"03-protect-planning-and-commits.md"},{"id":"preflight","command":"bench test --package ./internal/preflight","ticket":"03-protect-planning-and-commits.md"},{"id":"roadmap","command":"bench test --package ./internal/roadmap","ticket":"03-protect-planning-and-commits.md"},{"id":"commitment","command":"bench test --package ./internal/commitment","ticket":"03-protect-planning-and-commits.md"},{"id":"charge-evidence-system","command":"bench test --check system","ticket":"03-protect-planning-and-commits.md"},{"id":"landing","command":"bench test --package ./internal/landing","ticket":"03-protect-planning-and-commits.md"},{"id":"bench","command":"bench test --package ./cmd/bench","ticket":"03-protect-planning-and-commits.md"},{"id":"conformance","command":"bench test --package ./internal/conformance","ticket":"03-protect-planning-and-commits.md"},{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"03-protect-planning-and-commits.md"}]},{"id":"DC-C4","tickets":["04-bind-worktrees-and-shifts.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"04-bind-worktrees-and-shifts.md"},{"id":"shift","command":"bench test --package ./internal/shift","ticket":"04-bind-worktrees-and-shifts.md"},{"id":"intent","command":"bench test --package ./internal/intent","ticket":"04-bind-worktrees-and-shifts.md"},{"id":"root-help","command":"bench test --package ./cmd/bench","ticket":"04-bind-worktrees-and-shifts.md"},{"id":"commitment","command":"bench test --package ./internal/commitment","ticket":"04-bind-worktrees-and-shifts.md"},{"id":"conformance","command":"bench test --package ./internal/conformance","ticket":"04-bind-worktrees-and-shifts.md"},{"id":"system","command":"bench test --check system","ticket":"04-bind-worktrees-and-shifts.md"}]},{"id":"DC-C5","tickets":["05-authorize-current-publication.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"05-authorize-current-publication.md"},{"id":"landing","command":"bench test --package ./internal/landing","ticket":"05-authorize-current-publication.md"},{"id":"commitment","command":"bench test --package ./internal/commitment","ticket":"05-authorize-current-publication.md"},{"id":"intent","command":"bench test --package ./internal/intent","ticket":"05-authorize-current-publication.md"},{"id":"bench","command":"bench test --package ./cmd/bench","ticket":"05-authorize-current-publication.md"},{"id":"conformance","command":"bench test --package ./internal/conformance","ticket":"05-authorize-current-publication.md"},{"id":"system","command":"bench test --check system","ticket":"05-authorize-current-publication.md"},{"id":"repository","command":"bench test --package ./internal/commitment/repository","ticket":"05-authorize-current-publication.md"}]},{"id":"DC-C6","tickets":["06-close-verified-spec-delivery.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"06-close-verified-spec-delivery.md"},{"id":"landing","command":"bench test --package ./internal/landing","ticket":"06-close-verified-spec-delivery.md"},{"id":"gate","command":"bench test --package ./internal/gate","ticket":"06-close-verified-spec-delivery.md","probe":"Omit one sequence removal from the exact allowed transform. The exact-transform check must fail, then pass after the restore."},{"id":"roadmap","command":"bench test --package ./internal/roadmap","ticket":"06-close-verified-spec-delivery.md"},{"id":"repository","command":"bench test --package ./internal/commitment/repository","ticket":"06-close-verified-spec-delivery.md"},{"id":"commitment","command":"bench test --package ./internal/commitment","ticket":"06-close-verified-spec-delivery.md"}]},{"id":"DC-C7","tickets":["07-close-light-delivery.md","08-verify-milestone-outcomes.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"07-close-light-delivery.md"},{"id":"landing","command":"bench test --package ./internal/landing","ticket":"07-close-light-delivery.md"},{"id":"spec","command":"bench test --package ./internal/spec","ticket":"07-close-light-delivery.md"},{"id":"commitment","command":"bench test --package ./internal/commitment","ticket":"08-verify-milestone-outcomes.md"},{"id":"intent","command":"bench test --package ./internal/intent","ticket":"08-verify-milestone-outcomes.md"},{"id":"repository","command":"bench test --package ./internal/commitment/repository","ticket":"07-close-light-delivery.md"},{"id":"bench","command":"bench test --package ./cmd/bench","ticket":"08-verify-milestone-outcomes.md"},{"id":"conformance","command":"bench test --package ./internal/conformance","ticket":"08-verify-milestone-outcomes.md"},{"id":"milestone-repository","command":"bench test --package ./internal/commitment/repository","ticket":"08-verify-milestone-outcomes.md"},{"id":"gate","command":"bench test --package ./internal/gate","probe":"Omit the folder removal proof from the tickets-only oracle. TestCommitmentTicketsOnlyTransform must fail, then pass after the restore.","ticket":"07-close-light-delivery.md"}]},{"id":"DC-C8","tickets":["09-project-commitment-guidance.md"],"verification":[{"id":"roadmap","command":"bench test --package ./internal/roadmap","ticket":"09-project-commitment-guidance.md"},{"id":"status","command":"bench test --package ./internal/status","ticket":"09-project-commitment-guidance.md"},{"id":"dashboard","command":"bench test --package ./internal/dashboard","ticket":"09-project-commitment-guidance.md"},{"id":"bench","command":"bench test --package ./cmd/bench","ticket":"09-project-commitment-guidance.md"},{"id":"anchors","command":"bench test --package ./internal/anchors","ticket":"09-project-commitment-guidance.md"},{"id":"conformance","command":"bench test --package ./internal/conformance","ticket":"09-project-commitment-guidance.md"},{"id":"commitment","command":"bench test --package ./internal/commitment","ticket":"09-project-commitment-guidance.md"},{"id":"intent","command":"bench test --package ./internal/intent","ticket":"09-project-commitment-guidance.md"},{"id":"usage","command":"bench test --package ./internal/usage","ticket":"09-project-commitment-guidance.md"},{"id":"repository","command":"bench test --package ./internal/commitment/repository","probe":"Omit the continuation write in approval. TestCommitmentContinuationApproval must fail, then pass after the restore.","ticket":"09-project-commitment-guidance.md"}]},{"id":"DC-C9","tickets":["10-qualify-linked-adoption.md"],"verification":[{"id":"installed-adoption","command":"bench test --check system","ticket":"10-qualify-linked-adoption.md"},{"id":"finite-adoption-review","command":"bench roadmap --context","ticket":"10-qualify-linked-adoption.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/roadmap-delivery-commitment/spec.md"},{"id":"installed-adoption","command":"bench test --check system"},{"id":"route-inventory","command":"bench test --package ./cmd/bench"},{"id":"completion-oracle","command":"bench test --package ./internal/gate"},{"id":"guidance","command":"bench test --package ./internal/conformance"}],"execution":{"mode":"delegate","run_id":"dc-full-20261004","orchestrator_session":"claude:session_01U6xYL2GjVjH4DvNmn18cMy","author_limit":1,"assignments":{"01-plan-exact-policy-changes.md":[],"02-admit-committed-outcomes.md":[],"03-protect-planning-and-commits.md":[],"04-bind-worktrees-and-shifts.md":[],"05-authorize-current-publication.md":[{"session":"claude:dc_t05","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"medium","source":"7adeee6df6e03efa44871c36f3e8a52b1777934f","native_ref":"claude-agent:dc_t05"},{"session":"claude:dc_r05_1","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"medium","source":"92c7ba0bdaa459c0e6b26057d2af081ad970a072","native_ref":"claude-agent:dc_r05_1","predecessor":"claude:dc_t05","trigger":"user-directed","stopped":"The ticket 05 author returned its final report with no live command, test, or write.","preserved":"6664d59e0431888c754d2b5c0bdb01887e1646fe"},{"session":"claude:dc_r05_2","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"medium","source":"d4c106b1f69347b4ce48b6f9ebfa46c0d7f0da9b","native_ref":"claude-agent:dc_r05_2","predecessor":"claude:dc_r05_1","trigger":"user-directed","stopped":"Repair session 1 returned its final report with no live command, test, or write.","preserved":"1642decabd17dcd1271847abad848a2d7a88a3aa"}],"06-close-verified-spec-delivery.md":[{"session":"claude:dc_t06","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"medium","source":"24f2f2d012cf0f83332c1de0858a6066e868873a","native_ref":"claude-agent:dc_t06"},{"session":"claude:dc_r06_1","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"medium","source":"5904aa60ecb964ad032b30d75db456017bd7f6c7","native_ref":"claude-agent:dc_r06_1","predecessor":"claude:dc_t06","trigger":"user-directed","stopped":"The ticket 06 author returned its final report with no live command, test, or write.","preserved":"7a0e9080652f89f4d968e7e92e98f7a3046ff376"},{"session":"claude:dc_r06_2","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"medium","source":"edab452c49fef1b4e4912ef33cf60d4b73205cf1","native_ref":"claude-agent:dc_r06_2","predecessor":"claude:dc_r06_1","trigger":"user-directed","stopped":"Repair session 1 returned its final report with no live command, test, or write.","preserved":"9f9e2d6bf64139818743ce96ca6ddfe2d90ff4a3"}],"07-close-light-delivery.md":[{"session":"claude:dc_t07","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"medium","source":"55c6f9ccf5f6db7a48e531aac0959e2adb534f4c","native_ref":"claude-agent:dc_t07"},{"session":"claude:dc_r07_1","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"medium","source":"036a508ab383fe37c980742e8de1f141157a84c9","native_ref":"claude-agent:dc_r07_1","predecessor":"claude:dc_t07","trigger":"user-directed","stopped":"The ticket 07 author returned its final report with no live command, test, or write.","preserved":"d074aaf143bb52b9cdce36ba78388e993af0fd21"}],"08-verify-milestone-outcomes.md":[{"session":"claude:dc_t08","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"medium","source":"d074aaf143bb52b9cdce36ba78388e993af0fd21","native_ref":"claude-agent:dc_t08"},{"session":"claude:dc_r08_1","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"medium","source":"a8d352d36f9a4ccc5f457b29defa1e813db1fc76","native_ref":"claude-agent:dc_r08_1","predecessor":"claude:dc_t08","trigger":"user-directed","stopped":"The ticket 08 author returned its final report with no live command, test, or write.","preserved":"b1296652eae56ea7b941516529aa24e0e57c5260"}],"09-project-commitment-guidance.md":[{"session":"claude:dc_t09","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"medium","source":"1a803d5c9782ff25ec58efa9ddab9a5a5bc2bc91","native_ref":"claude-agent:dc_t09"},{"session":"claude:dc_t09_high","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"high","source":"d8650f234a3687058d761b603854b905df0a5aa1","native_ref":"claude-agent:dc_t09_high","predecessor":"claude:dc_t09","trigger":"user-directed","stopped":"The orchestrator stopped the medium-effort author; no test, build, or write process remained live.","preserved":"The uncommitted worktree changes at d8650f234a3687058d761b603854b905df0a5aa1, with a patch and tar copy in the orchestrator scratchpad."},{"session":"claude:dc_r09_1","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"medium","source":"5a590b5aaffb173b4575409bd47bf96fbb60fbc7","native_ref":"claude-agent:dc_r09_1","predecessor":"claude:dc_t09_high","trigger":"user-directed","stopped":"The ticket 09 author returned its final report with no live command, test, or write.","preserved":"0f7fd8768eae657ab9175a96e5e9b4afe81f29d1"},{"session":"claude:dc_r09_2","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"medium","source":"dae5402257715dcd87317646e21e9a4ef22d8c52","native_ref":"claude-agent:dc_r09_2","predecessor":"claude:dc_r09_1","trigger":"user-directed","stopped":"The cycle 1 repair session returned its final report with a clean tree and no live command, test, or write.","preserved":"dae5402257715dcd87317646e21e9a4ef22d8c52"},{"session":"claude:dc_r09_3","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"medium","source":"c9c51f64603961220233abdc905e06d4730ba2fd","native_ref":"claude-agent:dc_r09_3","predecessor":"claude:dc_r09_2","trigger":"user-directed","stopped":"The cycle 2 repair session recorded its verification and returned with no live command, test, or write.","preserved":"047e47c304d96d40b746dc98cd08c994c5e8cd82"}],"10-qualify-linked-adoption.md":[{"session":"claude:dc_t10","assignment":"af49ac1c59888c026ae64548dc90b756","model":"opus","effort":"medium","source":"ffc100d510e93a3d73fe3c5086af2a94d7a7e343","native_ref":"claude-agent:dc_t10"}]}}}
```

## Testing decisions

Use command-level observations for admission, approval, and closure. Assert both the reported refusal and unchanged authority or publication state.
Use the existing intent transaction seam for competing claims and persistence failures. Use the real installed broker journey for the publication guarantee.
Policy tests attach to the new owner, which multiple real command consumers use. They do not reach around a consumer into private helpers.

The ordinary gate executes package tests and conformance owners. The tagged system journey runs through `bench test --check system` with `BENCH_KIT`.
Extend the existing disposable adoption journey for linked-project evidence. Do not create another standalone system harness or nested test runner.

Use TDD for policy transitions, competing starts, and the broker closure transform. Plan a red from each exact wrong result before its production change.
Observe an omission mutation at the central property: remove the installed landing admission call and attempt a displaced source publication.
Also omit one sequence removal from a successful closure transform. Both changes must turn their named behavioral checks red.
An independent coordinator probe uses a different site and mutation kind.

### Seam diagram

```text
reviewer direction + exact proposal
        |
        v
commitment commands -> commitment owner -> staged policy + local receipt
                              ^
                              | same policy decision
worktree / shift / commit / preflight / installed landing
        |
        v
current authority + assignment + operation -> admit or actionable refusal

reviewed source + complete evidence + approved outcome
        |
        v
installed broker -> exact delivery-and-closure tree -> whole-project gate
        |                                                   |
        +------------------ publish only on green <---------+

criterion evidence + current published revision
        |
        v
milestone verification -> exact receipt -> explicit completion approval
```

### Acceptance coverage map

Every named test below is planned. Its fixture description is the behavior cell. DC53 is review-owned.
The test package owns the production seam it drives. System rows name the installed command journey explicitly.

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| DC1 | 1 | Activating M1 yields exactly its ordered outcomes A then B | planned TestCommitmentActivate in internal/commitment/command_test.go | An unordered or implicit selection produces another order. |
| DC2 | 2 | Adding planned M2 leaves M1's eligible outcome unchanged | planned TestCommitmentPlannedMilestone in internal/commitment/command_test.go | A second active selector changes the next action. |
| DC3 | 3 | Inserting C before A and B requires delayed operands A and B | planned TestCommitmentDisplacementEffects in internal/commitment/authority_test.go | A generic approval hides the two displaced obligations. |
| DC4 | 3 | Removing unfinished B requires removed operand B | planned TestCommitmentRemovalEffects in internal/commitment/authority_test.go | An incomplete effect set can silently drop work. |
| DC5 | 3 | Switching M1 to M2 requires the unfinished M1 outcomes in the delayed set | planned TestCommitmentSwitchEffects in internal/commitment/authority_test.go | A switch can otherwise evade displacement accounting. |
| DC6 | 4 | Editing the proposed order after approval refuses commit with `bench commitment plan` | planned TestCommitmentProposalIdentity in internal/commit/commitment_test.go | A receipt for earlier bytes cannot approve later bytes. |
| DC7 | 4 | A changed predecessor policy invalidates its outstanding approval | planned TestCommitmentStaleApproval in internal/commitment/authority_test.go | A stale plan can overwrite a later decision. |
| DC8 | 4 | An identical approval replay preserves the policy and receipt bytes | planned TestCommitmentApprovalReplay in internal/commitment/authority_test.go | Replay must not create another transition. |
| DC9 | 4 | An unrelated default-branch advance preserves approval whose bound facts are identical | planned TestCommitmentUnrelatedAdvance in internal/commitment/authority_test.go | Commit identity alone would demand another priority decision. |
| DC10 | 5 | A planning commit adds uncommitted C while preserving the A then B commitment | planned TestCommitmentUncommittedAppend in internal/commit/commitment_test.go | Intake cannot authorize C or reorder A. |
| DC11 | 6 | A planning commit that deletes active A's row refuses | planned TestCommitmentProtectedDeletion in internal/commit/commitment_test.go | Deletion cannot serve as implicit completion. |
| DC12 | 6 | A planning landing that renames active A's row without approval refuses | planned TestCommitmentProtectedRename in internal/worktree/land_specless_test.go | A renamed owner cannot escape the protected comparison. |
| DC13 | 6 | A planning commit that replaces A's requirement prose refuses | planned TestCommitmentProtectedScope in internal/commit/commitment_test.go | Keeping the row ID alone cannot preserve its obligation. |
| DC14 | 6 | A planning landing that recommends unrelated C first refuses | planned TestCommitmentProtectedSequence in internal/worktree/land_specless_test.go | A sequence-only change cannot bypass membership protection. |
| DC15 | 5 | Adding a valid occurrence to A preserves its requirement identity | planned TestCommitmentOccurrenceUpdate in internal/roadmap/commitment_test.go | Event evidence must not require a scope rewrite. |
| DC16 | 7 | Starting uncommitted C refuses with `bench commitment plan` and no claim | planned TestCommitmentUncommittedStart in internal/commitment/admission_test.go | A command that only warns still admits unrelated work. |
| DC17 | 8 | Starting B while independent A is eligible refuses | planned TestCommitmentOrderedStart in internal/commitment/admission_test.go | Membership alone does not protect order. |
| DC18 | 9 | Blocking A permits independent B while retaining A and its reason | planned TestCommitmentBlockedSuccessor in internal/commitment/admission_test.go | Skipping A must not remove the obligation. |
| DC19 | 9 | Blocking A does not permit B when B depends on A | planned TestCommitmentBlockedDependency in internal/commitment/admission_test.go | A blocked dependency is not a satisfied dependency. |
| DC20 | 9 | An all-blocked commitment refuses unrelated C | planned TestCommitmentAllBlocked in internal/commitment/admission_test.go | No eligible committed work is not general permission. |
| DC21 | 10 | A second outcome refuses while A owns the default active slot | planned TestCommitmentSingleOutcome in internal/commitment/admission_test.go | Two distinct outcomes must not share one default allowance. |
| DC22 | 10 | A grant for A and B permits B but still refuses C | planned TestCommitmentExactParallelGrant in internal/commitment/admission_test.go | A numeric allowance would authorize the wrong pair. |
| DC23 | 11 | Two ticket assignments for A share one active outcome claim | planned TestCommitmentTicketAssignments in internal/worktree/pool_root_test.go | Ticket delegation must not consume another outcome slot. |
| DC24 | 11 | A sibling cannot inherit B's authority from an integration source bound to A | planned TestCommitmentSiblingBinding in internal/worktree/tree_target_test.go | The sibling route must preserve exact outcome identity. |
| DC25 | 12 | A refused shift writes no adapter marker or shift intent | planned TestCommitmentShiftBeforeEffects in internal/shift/refresh_test.go | The check must run before acquisition and execution. |
| DC26 | 12 | Unbound production commit refuses before formatting its Go file | planned TestCommitmentCommitBeforeEffects in internal/commit/commitment_test.go | A late check has already changed the candidate. |
| DC27 | 12 | Ordinary build preflight refuses an unbound deliverable | planned TestCommitmentBuildPreflight in internal/preflight/commitment_test.go | A prepared build charge must not grant missing admission. |
| DC28 | 13 | A displaced assignment refuses landing after its current binding becomes invalid | planned TestCommitmentStaleLanding in internal/worktree/land_identity_test.go | Start-time approval alone cannot authorize later publication. |
| DC29 | 13 | A policy change during the gate prevents the now-disallowed publication | planned TestCommitmentGateRace in internal/worktree/land_identity_test.go | A pre-gate check alone loses the authority race. |
| DC30 | 14 | A planning assignment commits and lands a decision map and staged spec before adoption | planned TestCommitmentPlanningBootstrap in internal/worktree/land_specless_test.go | Adoption cannot require already-admitted planning work. |
| DC31 | 14 | Plan-only preflight validates a staged spec without creating a delivery claim | planned TestCommitmentPlanOnly in internal/preflight/commitment_test.go | Authoring validation must not activate delivery. |
| DC32 | 15 | A planning assignment refuses a production file through commit and landing | planned TestCommitmentPlanningFence in internal/worktree/land_specless_test.go | A purpose label must not become an unrestricted write route. |
| DC33 | 16 | Restarting a process retains A's claim until delivery or an authorized blocker transition | planned TestCommitmentRetainedClaim in internal/commitment/admission_test.go | Process liveness is not outcome completion. |
| DC34 | 17 | Verified delivery of A removes exactly A's row and detail owner | planned TestCommitmentDeliveryClosure in internal/worktree/commitment_landing_test.go | A completed spec alone leaves stale roadmap work. |
| DC35 | 18 | Verified delivery removes every sequence reference to A or its completed spec | planned TestCommitmentSequenceClosure in internal/roadmap/commitment_test.go | Removing only the board row retains a dead next command. |
| DC36 | 19 | A partial delivery retains its residual roadmap obligation | planned TestCommitmentPartialDelivery in internal/worktree/commitment_landing_test.go | A `Roadmap:` label does not prove the full row delivered. |
| DC37 | 20 | A red prospective gate publishes neither delivery nor closure | planned TestCommitmentClosureGateRed in internal/worktree/commitment_landing_test.go | A separate closure write can falsely remove unfinished work. |
| DC38 | 20 | The gate refuses a transform that keeps A's sequence reference | planned TestCommitmentExactTransform in internal/gate/commitment_completion_test.go | A closure path allowlist would accept the omission. |
| DC39 | 20 | The gate refuses an extra deletion of unrelated B | planned TestCommitmentExtraClosure in internal/gate/commitment_completion_test.go | The allowed transform must identify exact obligations. |
| DC40 | 21 | A pre-oracle persistence failure leaves authority and destination unchanged | planned TestCommitmentPersistenceBeforeGate in internal/worktree/commitment_landing_test.go | A partially written intent must not grant publication. |
| DC41 | 21 | Interruption inside the oracle retains the unpublished obligation | planned TestCommitmentInterruptedGate in internal/worktree/commitment_landing_test.go | Gate interruption is not delivery evidence. |
| DC42 | 21 | A terminal local-write failure resumes the original publication without publishing again | planned TestCommitmentClosureResume in internal/worktree/commitment_landing_test.go | A failed local effect must not replay delivery. |
| DC43 | 22 | Empty roadmap rows without criterion evidence cannot complete M1 | planned TestCommitmentEmptyRowsNotComplete in internal/commitment/verification_test.go | Row absence alone cannot prove the milestone outcome. |
| DC44 | 23 | One unmet criterion refuses M1 completion despite a green gate | planned TestCommitmentUnmetCriterion in internal/commitment/verification_test.go | Code health cannot substitute for outcome evidence. |
| DC45 | 22 | Complete current evidence and explicit approval complete M1 without activating M2 | planned TestCommitmentMilestoneCompletion in internal/commitment/verification_test.go | Completion cannot infer successor approval. |
| DC46 | 23 | Evidence for another criteria identity or revision refuses verification | planned TestCommitmentStaleEvidence in internal/commitment/verification_test.go | Reused evidence must not verify a changed target. |
| DC47 | 24 | Roadmap, status, and dashboard name B when A is blocked and B is independent | planned TestCommitmentNextProjection in cmd/bench/commitment_test.go | A surviving legacy reader would still select A or unrelated work. |
| DC48 | 25 | An absent policy refuses new delivery with `bench commitment plan` | planned TestCommitmentAdoptionRequired in internal/commitment/admission_test.go | An old sequence cannot silently become authority. |
| DC49 | 26 | An explicitly listed legacy run can finish its existing scope after adoption | planned TestCommitmentLegacyContinuation in internal/worktree/land_spec_amendment_test.go | Adoption must preserve previously authorized delivery. |
| DC50 | 26 | An unlisted old assignment cannot use another run's continuation | planned TestCommitmentLegacyIdentity in internal/worktree/pool_root_test.go | Assignment age alone is not authorization. |
| DC51 | 27 | An adopted linked project refuses an uncommitted production start through its installed wrapper | planned TestCommitmentLinkedAdoption in internal/systemtest/adoption_test.go | Kit-only checks cannot protect linked projects. |
| DC52 | 28 | The prerequisite installs before an initial policy exists and then admits only an approved adoption | planned TestCommitmentBootstrapInstall in internal/systemtest/adoption_test.go | The feature cannot require its own unavailable authority during installation. |
| DC53 | 29 | The kit adoption proposal includes only the named quality owners with remaining verified obligations | review-owned: compare the proposal with decision 14 and verified delivery history | A broader defect inventory would expand the approved milestone. |
| DC54 | 30 | Guidance requires purpose-based defect and refactor priority without automatic displacement | planned TestCommitmentGuidance in internal/conformance/commitment_guidance_test.go | A new severity sort would revive the original priority problem. |
| DC55 | 31 | Empty policy bytes refuse instead of reporting adoption-required | planned TestCommitmentEmptyPolicy in internal/commitment/parse_test.go | Empty corruption must not grant bootstrap behavior. |
| DC56 | 31 | Duplicate fields, unknown fields, duplicate IDs, and dependency cycles each refuse policy input | planned TestCommitmentMalformedPolicy in internal/commitment/parse_test.go | Ambiguous input cannot supply an approval target. |
| DC57 | 31 | A live symlink, dangling symlink, or FIFO at a policy input refuses before reading | planned TestCommitmentUnsafeInput in internal/commitment/command_test.go | Following or opening the path can import authority or block. |
| DC58 | 32 | An input path containing spaces and glob bytes resolves as one literal operand | planned TestCommitmentLiteralInput in internal/commitment/command_test.go | Shell reinterpretation would read another proposal. |
| DC59 | 32 | A flag value resembling an operation cannot dispatch that operation | planned TestCommitmentGrammar in internal/commitment/command_test.go | Positional guessing can select the wrong command. |
| DC60 | 32 | Control-bearing decision text refuses without changing its receipt | planned TestCommitmentControlInput in internal/commitment/command_test.go | A split field can change the recorded decision. |
| DC61 | 33 | Racing A and B starts produce exactly one default active outcome | planned TestCommitmentConcurrentStarts in internal/commitment/concurrent_test.go | Separate check and write steps can both succeed. |
| DC62 | 33 | A concurrent policy replacement cannot consume a stale plan | planned TestCommitmentConcurrentApproval in internal/commitment/store_test.go | The receipt check and transition must share one transaction. |
| DC63 | 34 | A candidate that removes its admission call still cannot land displaced work through the installed broker | planned TestCommitmentInstalledAuthority in internal/systemtest/adoption_test.go | The candidate cannot be its own publication authority. |
| DC64 | 35 | Omitting one routed admission consumer fails the command-route test | planned TestCommitmentRouteInventory in cmd/bench/commitment_test.go | A shared owner without every caller leaves a bypass. |
| DC65 | 36 | The data inventory names commitment record contents and their local retention | planned TestCommitmentDataInventory in internal/conformance/data_handling_test.go | A new durable record must not have undocumented contents. |
| DC66 | 9 | Unblocking A leaves active B in place until B finishes or blocks | planned TestCommitmentUnblockOrder in internal/commitment/admission_test.go | An unblock operation must not displace active work. |
| DC67 | 20 | Tickets-only delivery closes its approved complete roadmap obligation in the same publication | planned TestCommitmentTicketsOnlyClosure in internal/worktree/commitment_light_landing_test.go | Closure must not depend on a full spec file. |
| DC68 | 20 | A delivery without roadmap owners completes only its bound obligation | planned TestCommitmentNoRoadmapOwner in internal/worktree/commitment_light_landing_test.go | Ordinary linked projects need no invented roadmap row. |
| DC69 | 12 | A competing identity error precedes commitment refusal before any mutable effect | planned TestCommitmentRefusalOrder in internal/worktree/tree_target_test.go | A reordered check can misattribute an unknown assignment. |
| DC70 | 22 | Missing, duplicate, blocked, or unknown criterion results each refuse milestone completion | planned TestCommitmentCriterionCoverage in internal/commitment/verification_test.go | An incomplete evidence list cannot become a verified outcome. |
| DC71 | 32 | A valid proposal whose last line has no newline plans the same transition | planned TestCommitmentInputFraming in internal/commitment/command_test.go | Human-authored input must not require persisted-record framing. |
| DC72 | 33 | A competing blocker cannot change admission between the final check and ref publication | planned TestCommitmentPublishLock in internal/worktree/land_identity_test.go | A last read without a lock retains a runtime authority race. |
| DC73 | 8 | A changed or deleted approved deliverable refuses start without changing bindings or claims | planned TestCommitmentStartPublishedIdentity in internal/commitment/admission_test.go | Path equality cannot prove current source identity. |
| DC74 | 8 | Start accepts an approved spec or tickets-only folder and refuses an ordinary source file | planned TestCommitmentDeliverableTypes in internal/commitment/deliverable_test.go | A generic regular-file reader rejects folders and admits the wrong file type. |
| DC75 | 18 | Verified delivery removes satisfied dependency references to A from the board dependency tables | planned TestCommitmentDependencyClosure in internal/roadmap/commitment_test.go | A closed row named as a blocker keeps dependent work blocked. |
| DC76 | 26 | A listed legacy run closes its delivered scope and keeps a partly delivered scope open | planned TestCommitmentLegacyClosure in internal/worktree/commitment_light_landing_test.go | A continuation without a binding otherwise lands with no closure. |
| DC77 | 17 | Retirement after a verified closure schedules no roadmap row cleanup | planned TestCommitmentRetireClosedRow in internal/spec/spec_test.go | A retire hint for a closed row sends the next session to remove absent work. |
| DC78 | 20 | The gate accepts the exact tickets-only close and refuses a tree that keeps the closed folder, the satisfied row's detail owner, or the delivered sequence entry | planned TestCommitmentTicketsOnlyTransform in internal/gate/commitment_completion_test.go | Admission alone grades only the policy bytes of a tickets-only close. |
| DC79 | 20 | The gate refuses a completion tree that changes one byte outside the broker transform | planned TestCommitmentUnrelatedByte in internal/gate/commitment_completion_test.go | A path allowlist accepts an unrelated edit inside a closed path. |
| DC80 | 22 | A policy whose active milestone has no criteria refuses | planned TestCommitmentCriteriaBeforeActivation in internal/commitment/parse_test.go | A milestone without criteria can complete on an empty list. |
| DC81 | 23 | Completion refuses while an outcome is undelivered and when the proposal changes an examined criterion | planned TestCommitmentCompletionProposal in internal/commitment/verification_refusal_test.go | A changed or partial milestone must not inherit earlier evidence. |
| DC82 | 26 | A listed legacy run lands a tickets-only folder in its scope, and the publication closes that delivery | planned TestCommitmentLegacyClosure in internal/worktree/commitment_light_landing_test.go | A continuation without a spec path otherwise cannot close its folder. |
| DC83 | 26 | Approval records a continuation for exactly the listed runs, each with its assignment, request, and scope | planned TestCommitmentContinuationApproval in internal/commitment/repository/publication_test.go | Without a writer, no real adoption can authorize a legacy run, so DC49 and DC76 prove only seeded state. |
| DC84 | 26 | An open listed continuation holds the default active slot until a parallel grant names it, its scope is delivered, or its run is no longer active | planned TestCommitmentContinuationOccupiesSlot in internal/commitment/continuation_test.go | A continuation outside the slot count lets adoption start a second concurrent delivery with no grant. |

### Edge inventory

The profile's shell checklist applies to the new command operands, persisted records, and publication route.
DC55 to DC60 cover absent versus empty input, strict decoding, special files, symlinks, literal paths, grammar, and control bytes.
The existing strict JSON document reader accepts a final input line without a newline. Machine-written receipts retain its existing newline framing.
JSON field names and identity tokens use exact grammar. Human descriptions cannot carry line controls.

DC61 and DC62 cover competing writers. DC40 to DC42 cover persistence before the oracle, interruption within it, and terminal persistence.
DC28 and DC29 cover stale and changing authority. DC34 to DC39 cover exact closure, partial obligations, omitted removal, and unrelated deletion.

DC30 and DC31 keep planning and adoption reachable. DC49 and DC50 distinguish an approved legacy continuation from an old unapproved assignment.
All ordinary route tests use the existing package seams. Any process-global substitution stays in its owning test process.

Won't handle: cryptographic human identity — the surviving operator approval command uses Bench's documented same-user trust posture.
Won't handle: coordination across independent clones — the surviving start command serializes all worktrees of one Bench repository.
Won't handle: a raw editor or raw Git bypass — the surviving Bench commit and installed landing enforce their command scope.
Won't handle: automatic proof of arbitrary natural-language criteria — the surviving verifier binds native evidence for explicit reviewer outcome assessment.
Won't handle: broad roadmap restructuring — the surviving intake command retains findings outside the commitment.

## Ownership fences

This fence is the exact union of ticket writes, excluding spec-local artifacts, capture, and the review pickup.
The review pickup is `reviews/roadmap-delivery-commitment.md` (new). Its source-bound evidence accompanies final implementation review.
Directory fences include focused sibling extraction destinations for the named owner. They do not authorize unrelated refactoring.
An author who discovers another destination updates the ticket and this union, then reruns closure preflight before writing there.

- `.agents/commands/bench-debug.md`
- `.agents/commands/bench-drain.md`
- `.agents/commands/bench-final-check.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-setup-repo.md`
- `.agents/commands/bench-what-next.md`
- `.agents/commands/bench-write-spec.md`
- `.agents/commands/bench.md`
- `.bench/BENCH-reference.md`
- `.bench/BENCH.md`
- `CHANGELOG.md`
- `DATA_HANDLING.md`
- `README.md`
- `cmd/bench`
- `cmd/bench/main.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `docs/adr/0015-the-landing-verb-is-the-one-author-of-the-spec-flip.md`
- `internal/anchors`
- `internal/commit`
- `internal/commitment` (new)
- `internal/conformance`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/conformance/help_inventory_single_source_test.go`
- `internal/dashboard`
- `internal/gate/checkpoint.go`
- `internal/gate/commitment_completion_test.go` (new)
- `internal/gate/completion.go`
- `internal/gate/completion_test.go`
- `internal/gate/evaluation.go`
- `internal/gittest`
- `internal/intent`
- `internal/jsonfile`
- `internal/landing`
- `internal/landing/attribution.go`
- `internal/landing/close.go`
- `internal/preflight`
- `internal/roadmap`
- `internal/shift`
- `internal/spec`
- `internal/status`
- `internal/systemtest/adoption_test.go`
- `internal/systemtest/charge_evidence_test.go`
- `internal/systemtest/owner_landing_fixture_test.go`
- `internal/systemtest/owner_land_race_test.go`
- `internal/systemtest/otel_verbs_test.go`
- `internal/tickets/registry_data.go`
- `internal/usage`
- `internal/worktree`
- `internal/worktree/land_fixtures_test.go`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/workflow-guidance-anchors`
- `tests/canary/data-handling-derivation/undocumented-passlist-var`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/docs-currency-token-diet/dogfood-referent-shipped`
- `tests/canary/docs-currency-token-diet/introduces-undeclared-command`
- `tests/canary/docs-currency-token-diet/missing-cli-inventory`
- `tests/canary/docs-currency-token-diet/readme-command-first`
- `tests/canary/docs-currency-token-diet/stale-cli-doc-reference`
- `tests/canary/docs-currency-token-diet/stale-command-reference`
- `tests/canary/load-validity-metadata/readme-shared-rule-drift`
- `tests/canary/load-validity-metadata/shared-rule-drift`
- `tests/canary/row-next-grammar/token-table-lacks-kit-edit`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/dangling-index`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/skills-index-command-adapters/missing-index-field`
- `tests/canary/skills-index-command-adapters/stale-index-wording`
- `tests/canary/skills-index-command-adapters/unindexed-skill`


## Out of scope

The following capabilities are independent future work. The estimates describe their smallest plausible separate specification, implementation, and gate exercise.

| separate capability | estimate | retained limit |
| --- | --- | --- |
| Distributed scheduling across independent clones | 8 edits, 3 gate runs | One repository's linked worktrees share the current commitment. |
| Cryptographic reviewer authentication | 6 edits, 3 gate runs | The operator records explicit direction under the existing trust model. |
| Automatic global roadmap scoring or age-based promotion | 4 edits, 2 gate runs | New findings stay uncommitted until an explicit decision. |

The named quality implementations are the already-approved successor milestone, not a cut from this enforcement feature.
This spec adds no release qualification, new survey, unrelated refactor, or new roadmap row.

## Further notes

### Authoring authority and review order

On 2026-10-04, the reviewer delegated routine spec and slice decisions to the author. This supersedes the map's earlier discretion limitation for engineering choices.
It does not reopen the 17 product decisions or expand the quality milestone.
The reviewer requested a spec-only Sol 6.1 xhigh review, followed by slicing, then a Sol 6.1 xhigh review of the slices.

The author preserved that sequence. Spec review accepted the design, and the separate slice review accepted the completed graph.
The approval table below records the delegated engineering choices and accepted review boundaries.

### Source verification

The author read the approved map and all 17 answers at the shaping source, then verified the landed source at `24b92a555ee5a6f479d305cb5b574bee4b2adfdf`.
Each structured source remains in place. The compiled map and its complete ticket folder moved together under this spec.
The current drain rule, final-check closure rule, FT373 scope, roadmap membership, and FT358 staged status still support the map's claims.

The research report remains historical evidence at its stated source commit. The live flow command reproduced its quoted values during this phase.
The FT358 source was read for its status, scope, and existing implementation decisions. Its unrelated Markdown coverage rows do not authorize this prerequisite.

### Source clauses and coverage

| decision | exact source clause | coverage |
| --- | --- | --- |
| 1 | The rule applies to the shared Bench workflow, including linked projects. | DC51, DC63 |
| 2 | The rule protects a finite milestone and its ordered work list. | DC1, DC17 |
| 3 | The decision names the work that gets delayed or removed. | DC3, DC4, DC5, DC6, DC7 |
| 4 | A milestone completes only when its stated outcome is verified. | DC43, DC44, DC45, DC46, DC70 |
| 5 | A finding cannot join or reorder the delivery commitment without an explicit reviewer decision. | DC10, DC11, DC12, DC13, DC14, DC16 |
| 6 | The worker reports the blocker and preserves the blocked obligation. | DC18, DC19, DC20, DC66 |
| 7 | Parallel committed outcomes require explicit reviewer authorization. | DC21, DC22, DC23, DC24, DC61, DC84 |
| 8 | The reviewer explicitly approves activation or a switch. | DC1, DC2, DC5, DC45 |
| 9 | They refuse unauthorized work starts and commitment changes, and show the decision needed to proceed. | DC6, DC16, DC25, DC26, DC27, DC28, DC29, DC32, DC64 |
| 10 | Verified delivery closes its completed roadmap row and sequence references. | DC34, DC35, DC36, DC37, DC38, DC39, DC67, DC68 |
| 11 | Already-authorized runs may finish. | DC30, DC48, DC49, DC50, DC83, DC84 |
| 12 | Commitment protection and reliable roadmap closure ship as one coherent delivery. | DC37, DC51, DC52 |
| 13 | A class label alone does not establish the work's purpose. | DC54 |
| 14 | Verify which obligations have already shipped before scheduling the remaining work. | DC53 |
| 15 | Dependencies and the explicitly committed order govern execution. | DC17, DC19, DC54 |
| 16 | Procedural instructions alone do not satisfy the command-enforcement requirement. | DC25, DC26, DC28, DC63, DC64 |
| 17 | The commitment enforcement prerequisite runs first. | DC52, DC53 |

The occurrences of these clauses are their numbered decision tickets and the compiled index summaries. The implementation changes their command consumers, not their approved meaning.

### Reader and writer sweep

The sweep included hidden directories, command guidance, agent definitions, harness documentation, tests, scripts, and `.github` workflows.
It searched roadmap identity, `Recommended sequence`, spec completion, assignment lookup, and each changed command route.
No workflow file directly calls the new command family. Existing workflow build and test commands retain their behavior.

| fact or route | inspected owners and consumers | disposition |
| --- | --- | --- |
| Roadmap rows and sequence | `internal/roadmap/roadmap.go`, `tree.go`, `context_parse.go`, `context_render.go` | Reuse the parser and project commitment order. |
| Ambient next action | `internal/status/status.go`, `internal/dashboard/dashboard.go`, `render.go` | Consume the shared commitment projection. |
| Command inventory and grammar | `cmd/bench/main.go`, `command_registry.go`, `worktree_leaves.go`, `internal/usage` | Register the family and changed flags once. |
| Worktree start and identity | `internal/worktree/worktree.go`, `land_identity.go`, `internal/intent/assignment.go` | Compose admission with the current identity checks. |
| Shared writes | `internal/intent/transaction.go`, `ledger/ledger.go`, `ledger_aliases.go` | Extend the existing transaction and schema owner. |
| Shift execution | `internal/shift/shift.go`, `loop.go` | Refuse before mutable effects and adapter execution. |
| Commit publication | `internal/commit/commit.go`, `internal/landing/landing.go` | Enforce planning scope and current delivery binding. |
| Prepared builds | `internal/preflight/command.go`, `gather.go`, `decision.go`, `plan.go` | Separate plan validation from delivery readiness. |
| Completion evidence | `internal/reviewrecord/completion.go`, `recordcmd/command.go`, `internal/gate/completion.go` | Preserve evidence validation and extend only the exact broker transform. |
| Spec metadata and retirement | `internal/spec/spec.go`, `cmd/bench/spec_retire_listing.go` | Retain metadata ownership and remove deferred completed-row cleanup. |
| Guidance authority | `.bench/BENCH.md`, `.bench/BENCH-reference.md`, the eight fenced phase documents | Replace the contradictory start, sequence, and closure directions. |
| Policy enforcement advertisement | `internal/anchors/registry_data.go`, `internal/conformance/docs_workflow_checks_test.go`, `registry/checks.go` | Update existing anchor owners and executable bite proofs. |
| Package and command closure | `internal/tickets/registry_data.go`, `internal/conformance/injected_ports_registry_test.go` | Keep registry and injected-port obligations in ticket fences. |
| Public and local data claims | `README.md`, `SECURITY.md`, `DATA_HANDLING.md`, `package.json` | Preserve the same-user trust model and project-owned policy. |

The three shared writers are policy approval, runtime admission, and broker closure. They use exact identities and the intent transaction.
A candidate policy change is provisional until its approved landing. Runtime claims always recheck the current published policy.
A landing rechecks that authority before publication and reconciles local state from the published result.

The final admission check and ref publication share the intent lock. The gate runs outside that lock.
A concurrent blocker or approval cannot change runtime authority between the final check and publication.

### Pre-review proof checklist

- Cited symbols: `intent.Transact`, `roadmap.LoadTree`, `roadmap.ParseDocument`, `landing.Owner.LandReviewed`, `gate.WithCompletion`, and `spec.Implemented` resolve in the inspected owners.
- Import edges: the new owner imports no consumer package; the planned adapters import the owner. Existing consumers already reach intent, Git, and gate owners.
- Source-row clauses and occurrences: the source table above names each decision and its compiled occurrence.
- Promised field labels: `adoption-required`, `active_milestone`, `outcome`, `blocked`, `delayed`, `removed`, and the command operands identify the proposed public contract.
- Changed-function callers: registry dispatch calls the worktree and shift entries; commit calls the landing owner; worktree joins call `LandReviewed`; landing calls `WithCompletion`.
- Copy survival: DC47 and DC64 expose a surviving legacy selector or omitted command consumer. The owner supplies policy facts once.
- Rendered-shape readers: the changed sequence reaches roadmap context, dashboard, status, their tests, and command envelope tests. Ticket slicing assigns each exact reader.

The spec-only checkpoint does not claim final ticket closure. The author checks current anchors and dependency directions before spec review.
Slicing will complete per-ticket reader closure and write proposals before its independent review.
Fixture constructors requiring approved state include the worktree landing fixtures, commit landing fixture, shift collision fixture, preflight fixture, and system adoption journey.
Each constructor will either seed a real approved state or retain an explicit adoption-refusal case. No blanket test-only admission switch is allowed.

Current large owners include worktree creation, preflight gathering and decisions, roadmap parsing, command dispatch, and status.
Each affected ticket must create headroom in its own checkpoint. No later ticket can repair a file-budget increase.
The author will inspect directory headroom before choosing a split. This spec grants no structure-budget exception.

### Flagged additions

None to the approved product scope. Exact identities, constrained planning, atomic claims, and literal input handling implement the approved command-enforcement guarantee.
The engineering choices above do not add a priority score, automatic admission rule, remote coordinator, or new roadmap program.

### Slice ownership and approval

| source clauses or route | ticket owner | decision state |
| --- | --- | --- |
| Decisions 1 to 5 and 8: policy identity, one milestone, protected admission | 01, 02, 03 | Approved source; engineering delegated |
| Decisions 6 and 7: blockers, independent successors, one active outcome | 02, 04 | Approved source; engineering delegated |
| Decision 9: command enforcement and planning protection | 03, 04, 05, 09, 10 | Approved source; engineering delegated |
| Decisions 4, 10, and 12: verified closure and coherent delivery | 06, 07, 08 | Approved source; engineering delegated |
| Decision 11: explicit adoption and named legacy continuations | 04, 05, 10 | Approved source; actual adoption follows publication |
| Decisions 13 to 17: purpose-based priority and finite successor milestone | 09, 10 | Approved source; no added roadmap intake |
| Ordinary startup execution | 02, 04 | Operator-trusted CLI assumption corrected after spec review |
| Authenticated installed publication | 05, 06, 07, 10 | Existing manifest broker authority preserved |
| Independent spec review | Complete | GPT-6.1 Sol xhigh accepted round 1 |
| Routine engineering and ticket breakdown | Delegated and complete | Current user direction; no implementation authorization inferred |
| Independent slice review | Complete | GPT-6.1 Sol xhigh accepted round 2; final command accounting folded |

Each ticket names its producer contract, executable checkpoint, and read surface. Existing package fixtures remain the single source for setup.
For admission posture changes, the implementing ticket enumerates calls to its affected fixture helper before editing that helper.
The package-directory fences cover those current callers; a caller outside the fence requires a recorded expansion before use.

The caller census is `specs/roadmap-delivery-commitment/fixture-census.md`. It names the current fixture definitions and each direct call site.
The new policy shape cannot have an existing executable red during this planning phase.

Ticket 10 prepares adoption evidence only. The initial policy is created through the installed commands after this prerequisite publishes.
FT283 and FT284 remain open unless implementation evidence proves their entire distinct outcomes. Association with this feature does not close either row.

### Ticket 01 help projection closure

Ticket 01 includes the help conformance owner that registers command form projections.
The command form table remains the single source for help rows.
The author observed the root conformance refusal before this fence expansion.

### Current-session author transfer

The reviewer directs implementation in the current session on 2026-10-04.
This direction supersedes the fresh-author and orchestrator-only rules for this build.
The first ticket author stopped before any ticket commit, with no live command or test.
The current session preserves that work and performs fresh verification.
Independent Standards, Spec, and Coverage reviews still use separate GPT-6.1 Sol sessions at high effort.

Previous author: `/root/dc_t01`, GPT-5.6 Sol, high effort.
Transfer source: `1e9a2f97fa0de89e014c1affc64c17cff3fb178e`, with the uncommitted ticket 01 files preserved.
Stopped evidence: The native author return confirms no live commands, tests, or later writes.
Transfer trigger: user-directed.

### Exact policy field validation

Ticket 01 extends the existing JSON owner to enforce exact policy field names.
Its checkpoint retains the JSON package check with all prior checks.
The document scanner remains the one parser for duplicate fields and document framing.

### Delegated author direction

The reviewer directs fresh delegated authors on Claude Opus at medium effort from ticket 05 on 2026-10-04.
This direction supersedes the current-session author transfer for tickets 05 through 10.
Tickets 01 through 04 keep their accepted version 1 plans, so their checkpoints keep the current Codex session as the author.
Those tickets declare an empty assignment history in the version 2 plan, because no delegated author wrote them.

The orchestrator is the current Claude Code session. It writes plans and records only.
Independent Standards, Spec, and Coverage reviews use separate Claude Sonnet sessions at high effort.
