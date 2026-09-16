# Bounded charge evidence: research and design comparison

Recommendation: preserve a complete prepared evidence set, then deliver it through bounded, deterministic reads.
Use the existing preflight owner for build and review.
Prefer an immutable evidence artifact with a stateless reader over a new session-progress service.
The reviewer selects the artifact and stateless-read design.

The response budget, consumer trust contract, and removal of charge --full are approved.
The reviewer approved the resource policy, logical identity rules, and command flow.
Tickets 12 and 14 fix the canonical TOON manifest and the uncompressed Bench packfile.
Ticket 15 records final approval of the shaped scope.
Decision tickets own the approved predicates.

Scope: the charge output defect and its required consumers.
Evidence status: current-tree trace and a live build reproduction are complete.
Two disposable probes have run on Linux.
Their limited observations are recorded in tickets 5 and 7 and do not establish production compatibility.
The write-spec comparison establishes an adjacent contract gap, not the same observed defect.

Baseline: `9deb0a7af31712427ff47d6fd0515e8458dbde0d`
Consumed by: `specs/bounded-charge-evidence/decisions/bounded-charge-evidence.md`
Drift: refresh after charge producers, collectors, phase consumers, source policy, or relevant storage owners change
Retire when: an approved spec consumes this research and its behavior lands

## Question graph and evidence boundaries

1. Which consumers need charge facts, and which action follows each fact?
2. What preparation, identity, and completeness guarantees does the executable provide?
3. Why did full retrieval become required, and what failure does it prevent?
4. Does write-spec share that failure or only an interaction pattern?
5. Which existing owner can support evidence after another process or worktree takes over?
6. Given those facts, which candidate preserves the invariant with the least new behavior?

Questions 1 and 3 share the consumer/history read set.
Questions 2 and 5 share the preparation/storage read set.
Two read-only Sol/high delegates examined those independent sets.
The coordinator reopened the load-bearing code, phase instructions, tests, and historical spec.
Questions 4 and 6 follow the combined findings.

All factual evidence comes from this repository or its Git history.
The user supplied the original reproduction measurements.
The coordinator ran a new measurement in assignment `bounded-charge-shaping`.
No active debug-loop-guidance assignment supplied files or received changes.

## Derived core problem

A charge consumer cannot obtain every required pinned input through Bench while keeping each response within the consumer's output limit.
The available compact response names the evidence but routes retrieval to the response that exceeds that limit.
This blocks a compliant start or encourages incomplete, manually reconstructed context.

| Aspect | Current finding |
|---|---|
| Observed symptom | Full build output contains 77,162 bytes and exceeds the reproduced 60,000-byte limit. |
| User failure | The author cannot follow both required evidence retrieval and bounded-output rules. |
| Protected invariant | One task receives all required mechanical evidence from one verified source state before action. |
| Architectural cause | The shared renderer makes evidence completeness depend on embedding every source body in one response. |
| Core classification | An output-projection defect exposes an unclear completeness contract. Transport must support the corrected contract. |
| Not established | Source duplication, missing preparation checks, or an identical write-spec failure. |
| Adjacent opportunities | Fork handoff clarity, general context reuse, guidance reduction, and global query pagination. |

These statements do not depend on a bundle, a retrieval verb, or progressive disclosure.
The current renderer contains each required body once.
Its `complete` field changes with `full`, rather than a verified consumer state.
See `internal/preflight/charge.go:72`, `internal/preflight/charge.go:100`, and `.bench/BENCH.md:18`.

## Tested result

The coordinator invoked the supplied command from the dedicated assignment before writing its shaping artifacts.
The first sandboxed attempt stopped during Go-cache access and did not reach successful charge rendering.
The retry used ordinary cache access and changed no repository source.

| View | Exit | Bytes | Lines | SHA-256 |
|---|---|---|---|---|
| Full | 0 | 77162 | 24 | `25fece27feed060fd97b50d4e9e338ffe055a30066cdacbe1b8c1d2ec3666080` |
| Compact | 0 | 4696 | 26 | `98a1296c8cbe3d0cf3771992c48524038c8af5dbdfec794365c6ff2dda8f91f5` |

The captures are `/tmp/bounded-charge-shaping-full.toon` and `/tmp/bounded-charge-shaping-compact.toon`.
They are diagnostic captures, not a proposed retrieval interface.
The full result fails the supplied output-size predicate.
The compact result passes it but directs the consumer back to full retrieval.

The full digest differs from the supplied reproduction.
The packet includes the assignment identifier and absolute checkout path, which differ between reproductions.
Those fields come from `internal/preflight/charge.go:155`.
The coordinator also verified the five raw source bodies total 71,814 bytes at the baseline.
TOON encoding and the packet fields account for the larger response.

## Q1: consumers, decisions, and required visibility

| Consumer | Decision or action | Required facts before that action |
|---|---|---|
| Retained build author | Start one approved ticket | Selected ticket, approval, supplement, source identity, fence, dependencies, coverage, checks, and required source bodies. |
| Delegated-run orchestrator | Dispatch one ticket author | The same mechanical evidence plus the recorded author, venue, and approved line. |
| Ticket author | Implement and verify the ticket | The complete task-specific charge and supplement, independently checked against the assigned source. |
| Review coordinator | Inspect the change and dispatch axes | Frozen pair, shared diff, consumers, coverage, axis charges, and required guidance. |
| Standards, Spec, Coverage delegates | Derive findings independently | Their charge, the same prepared shared evidence, and each axis's primary sources. |
| Successor review harness | Resume authorized dispatch | Repository, assignment, frozen pair, preserved evidence, and the exact native continuation. |

Sources: `.agents/commands/bench-implement-spec.md:22`, `.agents/commands/bench-implement-spec.md:40`, and `.agents/commands/bench-implement-spec.md:78`.
Review sources: `.agents/commands/bench-review-implementation.md:48`, `:73`, `:86`, and `:122`.
The packet is a mechanical supplement to phase judgment.
It does not confer reviewer approval or prove comprehension.

The first response must identify the selected task, frozen pair, evidence identity, validation disposition, completeness state, and deterministic next action.
It must identify the fence, checks, coverage, and required-source inventories, including any undisplayed entries.
An action-capable response or verified consumer state must contain their complete values.
A long inventory can use required detail pages while the initial response remains explicitly nonterminal.

No source body inherently needs to share a response with another body.
They must share a verifiable membership relationship and source identity.
The complete ticket, spec, and required guidance must remain available before build action.
Review additionally requires the exact prepared generated evidence.
This inference follows the task requirements and the renderer's distinct source identities in `internal/preflight/charge.go:85`.

## Q2: current preparation and completeness

Current flow:

```text
approved spec/ticket or authorized review
  -> active assignment and explicit base/tip
  -> movement-checked preparation
  -> source checks and preflight verdict
  -> build source set OR review source set plus collectors
  -> compact identities + omitted list + full retrieval action
     OR full source bodies + complete=true
  -> phase checks approval/supplement/fence/dependencies
  -> author action or independent review dispatch
```

The source read compares no-follow checkout bytes with `git show source-tip:path`.
It refuses a mismatch, absent or empty source, unsupported file kind, or unrepresentable body.
The source identity is SHA-256 over exact body bytes.
The path is a separate field.
See `internal/preflight/charge.go:28`, `:180`, `:202`, and `:234`.

The packet currently has no overall digest.
The source digest alone does not bind semantic role, path, or source tip.
The surrounding row and the verified execution provide those relationships today.
This matters when a later process receives only a stored identifier.

Preparation requires an active assignment and a clean checkout.
One surrounding movement check includes preparation and rendering.
It permits one retry, discards the failed attempt, and refuses continued movement.
See `internal/preflight/preparation.go:21`, `:76`, and `internal/preflight/charge_snapshot_test.go:12`.
The snapshot covers repository state, not every external collector dependency.
See `internal/diff/snapshot.go:232` and `:337`.

Build requires five sources in a declared order.
Review requires five canonical sources plus generated diff, consumers, and coverage streams.
The collectors run once per preparation attempt, including compact review preparation.
All three axes receive the same shared evidence identity.
See `internal/preflight/charge.go:170`, `internal/preflight/review.go:50`, and `:130`.

The review shared identity hashes `diff || consumers || coverage` in fixed order.
Separate source digests identify each stream.
The separately rendered completion-evidence table does not join that shared identity.
It contains source and plan digests plus the current review-record state.
See `internal/preflight/review.go:35`, `:180`, and `:197`.
A new prepared-set identity must explicitly account for this generated metadata.

Today, complete means that one output includes every required retrievable source body.
It does not mean approved, read, understood, or authorized to act.
The compact response is a successful preparation with incomplete delivery.
The code does not maintain a consumer receipt or prevent arbitrary editor use.
See `internal/preflight/charge.go:72` and `.agents/commands/bench-implement-spec.md:30`.

## Q3: intended protection and enforcement

The originating FT311 spec identifies omitted inputs and mixed revisions from handwritten coordinator reconstruction.
It calls preparation evidence readiness, separate from permission to write.
It requires bounded output with complete retrieval, then selects compact identities plus full retrieval.
Source: `aafeb009:specs/ft311-delegate-preparation/spec.md:9`, `:26`, `:35`, and `:99`.

The build implementation arrived in `eb9f9785`, and review support arrived in `cd97a787`.
The review guidance correction in `062aaafd` requires consumers to use the prepared evidence.
The current shared renderer prevents build and review projection rules from drifting.
Its compact/full tests enforce the present route, but the large-ticket test bounds only compact output.
See `internal/preflight/charge_test.go:140` and `internal/preflight/charge_snapshot_test.go:69`.

Several uncoordinated file reads would discard the complete required-input inventory and its single-snapshot verification.
They could mix a current ticket with older guidance or omit generated review evidence.
Several coordinated, identity-checked reads need not violate that invariant.
The violation concerns coherence and coverage, not the number of calls.

The phase-level `implement-spec --full` is a different control.
It orchestrates a complete implementation run and records phase handoffs.
It must not be renamed merely because charge retrieval changes.
See `.agents/commands/bench-implement-spec.md:71`.

## Q4: write-spec comparison and consumer routes

| Concern | Write-spec today | Implement-spec today |
|---|---|---|
| Orientation | An approved decision source starts spec and ticket authoring. | An approved spec and ticket start implementation. |
| Repository/assignment | All writes belong in the phase worktree. | Charge preparation requires an active clean assignment. |
| Source identity | Exactly one named decision source; map sources require re-verification. | Explicit base/tip plus byte identities. |
| Approval | Ready map, confirmed conversation, or named reviewed artifact. | Reviewer approval stays separate from staged artifacts and generated charge. |
| Source selection | The phase accepts one of three source forms. | The phase selects an existing spec. |
| Ticket selection | The author creates the ticket breakdown. | The caller selects an approved existing ticket. |
| Ownership fence | The author derives the future build fence from ticket writes. | Preflight projects the existing declared fence and selected writes. |
| Line | The author fork inherits the invoking line. | The supplement names model, effort, cap, and mutation. |
| Coverage | The author creates and checks coverage rows. | The charge supplies selected row identities and full source bodies. |
| Guidance | Canonical phase, spec, domain, tickets, and line skills. | Canonical sources enter the mechanical evidence set. |
| Full retrieval | No compact/full evidence command exists. | Full charge retrieval is expressly required. |
| Next action | Author, check, slice, and present one approval table. | Retrieve evidence, verify prerequisites, then work the ticket. |
| Fresh session | The author fork inherits approved conversation context. | Source facts must be recovered and verified before action. |
| Resume | Reuse the compiled map in place; general handoff preserves phase state. | Revalidate sources before each ticket; retain delegated identities and obligations. |
| Movement | Re-verify map sources and update references when moving the map. | Reject moved source pins and regenerate the charge. |

Sources: `.agents/commands/bench-write-spec.md:18`, `:32`, `:37`, `:45`, and `:57`.
Implementation sources: `.agents/commands/bench-implement-spec.md:22`, `:30`, `:40`, and `:78`.
Authoring requirements: `.agents/skills/bench-craft-spec/SKILL.md:55`.
The preflight grammar accepts build and review only: `internal/preflight/command.go:18`.

The write-spec author-fork rule leaves the exact inherited payload and interrupted-fork recovery implicit.
The anchor checks the fork instruction, not successful native payload delivery.
See `internal/anchors/registry_data.go:183` and commit `acd787a5`.
This is a candidate improvement, not proof that write-spec loses inputs or exceeds a response bound.
No live fork failure was reproduced in this research.

The common experience can expose orientation, source identity, missing prerequisites, and the next action.
The completeness predicates differ: settled decisions authorize spec synthesis; verified mechanical inputs permit an approved build action.
There is no demonstrated common retrieval owner.
Recommend no write-spec implementation changes for the charge defect.

| Route or consumer surface | Current responsibility | Migration consequence |
|---|---|---|
| Codex phase adapters | Read the canonical command file. | Keep them thin; no independent retrieval sequence. |
| Claude command/skill links | Resolve to the canonical command and skill sources. | Change the canonical source once. |
| Other AGENTS harnesses | Follow the same canonical files. | Require no harness-specific transport feature. |
| Front-door phase | Follow the action from status. | Preserve routing; update only changed executable actions. |
| Maps and status | Route ready decisions to spec authoring. | No charge protocol ownership. |
| SessionStart | Provide ambient inspection and status. | Do not inject all evidence into session startup. |
| Agent-line hook | Check model membership, with a fail-open adapter rim. | It proves neither retrieval nor comprehension. |
| Stop hook | Grade an armed shift at stop. | It does not authorize charge consumption. |
| Review fallback | Preserve charge inputs for a capable native harness. | Give the successor a stable evidence identity and retrieval route. |
| Historical/spec-less review | Use their existing diff entry points. | Explicitly outside this charge migration. |

Sources: `.bench/BENCH-reference.md:102`, `.agents/commands/bench.md:7`, and `internal/maps/maps.go:247`.
Hook sources: `.bench/hooks/session-start.sh:40`, `.bench/hooks/check-agent-line.sh:21`, and `.bench/hooks/stop.sh:11`.
Fallback sources: `.agents/commands/bench-review-implementation.md:62` and `:122`.
The hidden-directory reader sweep includes adapters, hooks, command prose, anchors, and canary mutations.

## Q5: storage ownership

| Existing owner | Why it does not own charge persistence today |
|---|---|
| Tracked review record | Stores results and native excerpts. Charge-time writes would dirty the required clean source. |
| Gate evidence | Repository-common storage has useful mechanics, but its authority is tree/oracle verdicts and verdict expiry. |
| Census | Stores assignment call counts and deletes them when the assignment retires. |
| OTel record | Explicitly excludes payload bytes. |
| Prospective artifact | Owns temporary gate checkouts and teardown. |
| Release evidence | Owns release-preflight output and checkout-local replacement. |
| Go build cache | Owns compilation data and independent cleanup. |
| Local-note routing | Chooses the primary note path; it supplies no evidence schema or lifecycle. |

Sources: `internal/reviewrecord/files.go:17`, `internal/reviewrecord/record.go:44`, and `internal/preflight/preparation.go:76`.
Store sources: `internal/gate/engine.go:107`, `:142`, `:181`, and `internal/git/worktree_admin.go:191`.
Other sources: `internal/census/census.go:357`, `internal/worktree/lifecycle.go:426`, and `internal/otelrecord/attributes.go:3`.
Additional owners: `internal/gate/prospectiveartifact/prospectiveartifact.go:71`, `internal/releaseevidence/release_evidence.go:18`, `internal/gocache/gocache.go:23`, and `internal/git/localnote.go:5`.

Inference: no existing storage owner accepts prepared charges without a new domain responsibility.
Repository-common placement can reuse path-resolution mechanics without placing charge data in the gate's record domain.
Checkout-local storage fails post-release retrieval.
Bench-home storage needs another repository-key lifecycle and provides no necessary advantage for this repository-bound evidence.

## Canonical terms for the selected design

| Term | Meaning | Avoid |
|---|---|---|
| Prepared evidence set | The closed required-input inventory, pinned descriptors, exact bodies, and generated facts from one successful preparation. | Complete response, prompt bundle. |
| Evidence identity | A digest over the versioned set contract and its exact content. | Bare file digest, checkout path, approval. |
| Source identity | The exact-byte digest for one named source. | Role identity, permission. |
| Delivery state | The required evidence pages verified as delivered to this consumer. | Understood, approved, ready without qualification. |
| Action prerequisites | Current source validation, complete required context, reviewer approval, and a complete task supplement. | One complete boolean. |
| Evidence page | A bounded, explicitly indexed portion of one prepared set. | Preview when full required bytes are included. |
| Evidence artifact | Immutable storage beneath the retrieval contract. | A second user-facing workflow. |

The existing glossary owns delegate charge, worktree, line, result projection, and acceptance row.
See `CONTEXT.md`.
These design terms do not change the shipped glossary during shaping.

## Q6: distinct designs

Each design keeps preparation's current refusal rules.
Each design must account for unbounded metadata as well as source bodies.
The shared source policy remains in one executable owner.

| Design | Mechanism | Benefit | Cost or failure | Assessment |
|---|---|---|---|---|
| A. Inline charge plus digest-addressed bundle | Prepare one immutable artifact; return task facts and its digest. Read artifact content separately. | Freezes generated review evidence and survives assignment release. | A whole-file read merely moves the overflow. Requires bounded reads and lifecycle ownership. | Viable transport, incomplete interface alone. |
| B. Manifest plus stateless details | Reconstruct each selected source from pinned Git objects; return a bounded page with membership proof. | Build can avoid durable copies and progress state. | Review collectors cannot be recomputed as the original frozen evidence. Durable generated bytes remain necessary. | Viable for build alone; hybrid required for review. |
| C. State-driven disclosure | Persist a consumer cursor and required-layer receipts; return the next missing layer. | Strong resume accounting and simple next actions. | Adds sessions, concurrent consumers, invalidation, locks, and receipt recovery. Cannot prove comprehension. | More state than this defect demonstrates. |
| D. Existing pinned-source projection | Return Git object locators and commands using existing repository sources. | Smallest build-only change; no new artifact store. | Existing raw reads and full collectors remain unbounded. Release, generated evidence, and complete accounting remain unresolved. | Does not meet the whole charge contract alone. |
| E. Prepared artifact with stateless bounded pages | Freeze one self-contained set; provide manifest and byte pages without server-side consumer progress. | Preserves review evidence and portable retrieval with one protocol. | Adds an evidence artifact lifecycle and local verification; needs a portable round-trip probe. | Recommended candidate. |

The comparison follows the body-only identities and generated review sources found in Q2.
The lifecycle costs follow the owner inventory in Q5.
No option can avoid a source-size problem merely by moving the ticket inline.
No option can prove that an agent understands bytes because a command returned them.

## Selected architecture and container contract

The artifact and stateless-read architecture is approved.
Decision tickets own the accepted grammar, logical format, operational policies, and trust predicates.
Tickets 12 and 14 own the selected encoding and physical-container contract.

The existing preflight command seam owns preparation, source policy, and deterministic disclosure.
A small artifact owner beneath it owns immutable publication, safe reads, and cleanup.
Build and review use the same identity, framing, pagination, and verification rules.
Their required-source sets and phase prerequisites remain distinct.

Ticket 13 owns the complete command grammar and its state contract.
The charge flag explicitly prepares persistent evidence.
The first response gives the evidence identity and the exact first retrieval command.
Each result carries one deterministic next action or an explicit terminal state.

The source selector accepts an identifier from the prepared manifest, never an arbitrary filesystem path.
Byte pages preserve original bytes, including missing final newlines.
Section names can be optional navigation metadata, but do not replace exact-byte delivery.
A section parser is unnecessary for the core fix.

### Layers, state, and the action barrier

| Layer | Required content | State after this layer | Next action |
|---|---|---|---|
| Orientation | Task selectors, pins, evidence identity, validation result, inventory counts, fence/check references, explicit pending state. | Set prepared; delivery incomplete. | Read the first manifest page. |
| Manifest | All source descriptors, required roles, byte sizes, metadata sections, and their identities. | Inventory complete; bodies may remain missing. | Read the first missing required body page. |
| Required details | Full fence/check/coverage metadata and every required source body. | Delivery complete only when all required page ranges verify. | Verify the set and current action prerequisites. |
| Verification | Set digest, source membership, byte lengths, and current task binding. | Evidence verified; approval and supplement remain separate. | Follow the phase action only after its independent prerequisites hold. |

Default traversal orders metadata first, then canonical sources, then generated evidence.
The producer derives that order from the same source descriptors used for preparation.
Each page identifies its set, source, offset, total length, page digest, and successor.
The last page marks stream termination, not proof that earlier pages were consumed.

The consumer verifies complete, gap-free page coverage against the manifest.
The consumer checks actual page results against the complete manifest.
The complete verified page ranges establish delivery coverage.
A remembered claim of completion or a transferred final cursor supplies no such evidence.

The `--verify` result proves artifact integrity and membership; it must not claim that the consumer read every page.
The phase requires both complete verified page coverage and the required context before action.
The candidate stores no consumer reading log and introduces no server progress session.
The selected design adds no hook that claims to police arbitrary editor use.

The selected barrier matches current honest-agent trust: fail closed within Bench and prohibit early action in canonical phase guidance.
Ticket 9 excludes a universal enforced edit barrier.
An acknowledgment alone is not evidence of delivery or comprehension.

A retained consumer can reuse unchanged verified content that remains available in its required context.
Reuse requires the same source identity, role, requiredness, and an established consumer context containing that evidence.
A fresh consumer retrieves its required context even when artifact creation and hashing are reusable.
A receipt without available content cannot replace context.
This avoids treating conversation memory as an evidence source.

### Identity and deterministic content

Ticket 12 selects a versioned Bench profile of flat-table TOON for the logical manifest.
The format fixes field order, source order, integer spelling, and string escapes.
It includes no incidental whitespace, timestamp, or absolute checkout location.
Source bytes remain exact; the format does not normalize them.
The expected evidence identity hashes the complete canonical manifest bytes.
The manifest binds each ordered source descriptor, source digest, page offset, page length, and page digest.

The logical set covers phase mode, selectors, base, source tip, required roles, paths, source kinds, source identities, and source bytes.
It also covers fence/check/coverage metadata, generated completion facts, and collector identity inputs.
No semantic role can change while the evidence identity stays unchanged.
The metadata and body inventory derive from one prepared structure.
Stored representations are verified derivations, not a competing policy source.

Assignment and absolute checkout path belong to the current action binding.
They do not alter identical logical evidence from sibling worktrees.
Preparation returns their verified binding to the evidence identity.
A later action revalidates that binding rather than treating an origin path as authority.
Collector provenance records the declared producer and invocation context.
Captured bytes remain part of the identity even when external inputs are not fully closed.

Ticket 11 excludes reproducible collector execution across different environments.
The current loader inherits the local Go environment through `internal/consumers/loader.go:26`.

The expected digest comes from the trusted preparation response or a retained verified handoff.
A bundle's self-reported digest cannot authenticate itself.
An independent reader hashes the exact manifest bytes and compares the expected identity.
The reader then validates source descriptors and their page digests.
For repository sources, the reader can additionally compare the pinned Git object.
Generated bytes retain their prepared identity and provenance without recomputation.

Each read verifies the complete manifest against the expected evidence identity.
It then verifies the requested page against the manifest before returning that page.
An independent consumer can retrieve the manifest through bounded pages and hash its exact reconstructed bytes.
The consumer can then check each evidence page against the verified manifest.
A page digest alone supplies no membership proof.
The expected identity still comes from trusted preparation or a verified handoff.

A page result certifies only that page and its descriptor membership.
It does not certify the remaining stored bytes or the consumer's delivery coverage.
Full artifact verification checks every required stored page and source digest.
This separate operation refuses corruption in an unread page.
Neither operation certifies comprehension or approval.

### Bounds and failure output

Ticket 8 fixes the response budget at 48,000 encoded stdout bytes, below the reproduced 60,000-byte ceiling.
It applies after TOON escaping, including metadata, paths, identifiers, errors, and next actions.
The producer reserves envelope space before selecting page content.
No source, ticket, fence list, coverage table, or diagnostic receives an unbounded inline exemption.

The initial view carries counts and stable references when inventories do not fit.
It explicitly states which required metadata remains undisplayed.
An oversized scalar becomes a required byte-paged value.
A malformed selector produces a bounded refusal without echoing its entire contents.
An optional preview names its total byte size and exact omission range.
The default candidate needs no preview.

### Lifecycle and hostile storage

Proposed location: a dedicated `bench-charge-evidence` namespace beneath Git's verified common directory.
The evidence owner resolves the root; a caller never supplies an arbitrary artifact path.
A fixed-format digest selects one immutable regular file.
No-follow resolution rejects symlinks, devices, FIFOs, sockets, unexpected directories, and invalid identifier syntax.

The writer builds and verifies a private temporary artifact in that namespace.
It publishes only after the final source-movement check succeeds.
Publication must be atomic and must not overwrite a different or corrupt artifact silently.
Identical concurrent writers verify and reuse the same published identity.
Different identities publish independently.
An interrupted write cannot expose a ready handle.

A reader validates the opened object, the complete manifest, and the requested page before returning that page.
A replaced, truncated, corrupt, or symlinked artifact returns a typed refusal.
It does not fall back to current checkout bytes or recompute generated review evidence.
The evidence identity remains readable from retained, sibling, and review worktrees in the same repository.
Assignment release removes no evidence artifact.
Moving the repository preserves its common-directory evidence.

Historical retrieval does not authorize action against a moved source or released assignment.
The current action binding must pass the original clean-assignment and source-tip checks.
A new source tip creates a new set identity.
Unchanged sources can reuse byte acquisition after their roles and membership are reverified in the new manifest.

A new preparation collects generated review evidence again.
Reads of an existing artifact never rerun its collectors.
A new capture can establish unchanged content through matching source identities.
A source-tip move always invalidates the previous permission to act.

The candidate uses explicit cleanup with a plan fingerprint, not automatic age eviction.
Cleanup acquires an exclusive repository evidence lock and refuses while a reader or writer holds a shared lock.
The cleanup plan names exact identities and consequences; apply rechecks the plan.
An idle retained handle may become unavailable after explicitly approved cleanup.
Retrieval then refuses and names preparation as recovery, rather than substituting bytes.

The owner reports capacity and refuses publication before exceeding its approved store quota.
It does not delete existing evidence to make room implicitly.
Ticket 10 sets the default store quota and permits an explicit override.
It sets no smaller per-artifact cap.
An explicit preparation override can raise the quota for a large required set.
The default is an approved operational policy, not a limit derived from the timing sample.

Interrupted temporary files are recoverable cleanup targets after the owner proves no writer holds them.
The bundle can exceed a consumer response limit because no command returns the whole bundle unbounded.

## Required scenario comparison

The letters refer to the distinct designs above.
The outcomes for E describe the selected contract.
Prototype observations cover only the checks and limits recorded in tickets 5 and 7.
Production acceptance remains unimplemented.

| # | Scenario | Required observable outcome for E | Consequence for alternatives |
|---|---|---|---|
| 1 | 40 KB spec with small other sources | All bytes arrive across bounded pages with one set identity. | A needs a bounded reader; B and C can page; D retains raw-read risk. |
| 2 | Ticket exceeds the whole inline budget | Ticket uses required pages; orientation stays bounded and nonterminal. | Every design must stop assuming an inline ticket. |
| 3 | Very large review diff or consumer report | Freeze once, then page the exact generated bytes. | B needs frozen generated storage; D alone fails; A and C can preserve it. |
| 4 | Two worktrees request identical evidence | One logical identity survives concurrent publication; neither writer exposes partial data. | A/C need publication coordination; B/D avoid build storage but still face review storage. |
| 5 | Tip changes after orientation | Old evidence stays identifiable; current action validation refuses the old binding. | B/D must use Git objects rather than current paths; A/C must separate history from current authority. |
| 6 | Another worktree uses the handle | Same-repository retrieval succeeds after independent identity verification. | Checkout-local A/C fail; B/D need reachable objects and generated evidence. |
| 7 | Origin assignment releases | Historical evidence remains readable; build action requires a new valid assignment binding. | A/C need repository-level lifetime; B/D cannot retain generated bytes without storage. |
| 8 | Artifact is replaced, truncated, or symlinked | Refuse without following links or returning substituted bytes. | All stored designs need safe-open and digest checks; B/D retain equivalent Git-source checks. |
| 9 | Digest looks numeric to TOON | Decode it as the exact string supplied by the producer. | All designs must use the shared encoder and decoded-cell assertions. |
| 10 | Path contains spaces, globs, Unicode, or Git quoting | Preserve representable path bytes and pass selectors as data. Refuse unsupported bytes and traversal. | All designs need source identities independent of shell interpretation. |
| 11 | No final newline | Reassembled bytes and digest exactly match the source. | All designs must avoid line-based framing that invents a newline. |
| 12 | Required source absent, empty, or special | Refuse preparation before a handle can claim a complete set. | All designs retain the current refusal policy. |
| 13 | One snapshot movement | Discard the first preparation, retry once, and publish only the verified attempt. | A/C must delay publication; B/D cannot mix attempts in detail results. |
| 14 | Identical pinned inputs | Identical complete logical inputs and captured bytes yield one identity; a Git pair alone does not close collector inputs. | A/C must exclude incidental timestamps and locations; B/D must fix producer dependencies. |
| 15 | Cleanup overlaps reader or writer | Cleanup refuses without deletion while the shared lock is held. | A/C require lifecycle rules; B/D avoid new cleanup only for repository-native sources. |
| 16 | Consumer stops after orientation | Verified page coverage is incomplete and the phase prohibits author action. | C can enforce its own transition; no design proves comprehension or blocks every editor. |
| 17 | One layer read, then tip moves | Invalidate action binding; verify unchanged layer reuse against the new manifest. | B must not silently reread HEAD; A/C need separate historical and active states. |
| 18 | Phase names an obsolete route | Guidance conformance and the stale-route mutation fail. | Every design must migrate phase owners and their exact anchors together. |
| 19 | Initial metadata nears the bound | Emit references and required metadata pages before crossing the encoded limit. | A/B/C cannot leave manifests unbounded; D needs more than terse body locators. |
| 20 | Role changes but body digest does not | Evidence identity changes; old role membership is insufficient. | Bare content-addressing in A/B/D fails without descriptor binding; C also needs it. |
| 21 | Fresh write-spec from approved conversation | Continue the existing author-fork route; no fictitious Git pin for conversation bytes. | A-D do not justify a new write-spec transport. |
| 22 | Resumed write-spec changes ticket breakdown | Revalidate the authored plan under write-spec's own approval contract. | Source reuse cannot imply approval of revised tickets in any design. |
| 23 | Build has approved spec but no selected ticket | Return selection-required state or the existing usage refusal; never claim a prepared ticket charge. | All designs need an actual selector before task-specific preparation. |
| 24 | Shared vocabulary, different predicates | Explain pending inputs consistently while retaining distinct phase prerequisites. | A shared session machine would add unsupported semantics. |
| 25 | Charge-only fix succeeds | Build and review meet the contract with write-spec behavior unchanged. | D alone remains insufficient; A/B hybrid or E can fit that scope. |

Current regression evidence includes `internal/preflight/charge_edges_test.go:181` and `internal/preflight/charge_file_kinds_test.go:75`.
Byte cases appear in `internal/preflight/charge_text_test.go:11`, `:30`, `:61`, and `:81`.
Snapshot cases appear in `internal/preflight/charge_snapshot_test.go:12` and `internal/preflight/review_charge_test.go:322`.
Those source reads do not claim that the selected protocol passes those scenarios.

## Observable acceptance and migration contract

1. The replacement charge route supplies bounded actions to every required byte, and charge --full receives the ordinary bounded usage error.
2. The consumer reconstructs the original five build source bodies without omission or mutation.
3. Review detail retrieval returns prepared collector bytes without a second collector run.
4. Every metadata and body response stays within the approved encoded-byte budget.
5. Required content that exceeds one page remains explicitly incomplete until delivery verifies.
6. A changed role, path, selector, pin, requiredness, or body changes the relevant evidence identity.
7. Dirty or moved sources and unsafe required files never produce a current action-ready state.
8. Concurrent preparation and cleanup cannot publish or return partial evidence.
9. Fresh and resumed consumers receive deterministic actions and preserve approval separation.
10. The canonical phase and executable output agree on retrieval without a second sequence in adapters.

These are behavioral inputs to later spec authoring, not a completed acceptance coverage map.
The spec author owns engineering seams, fixtures, mutation attachment, and gate attachment.
The existing public preflight command seam is the smallest demonstrated consumer seam.
The source-set descriptors should remain the one policy owner for renderer, artifact, verifier, and required-next-action derivation.

Ticket 6 removes charge `--full` without a compatibility route.
Unsupported uses receive the ordinary bounded usage error.
The phase-level full-run flag remains unchanged.

Migrate the build and review command prose plus their exact guidance anchors and canary mutations.
Relevant owners include `internal/anchors/registry_ft311_preparation.go:9` and `internal/anchors/registry_ft311_review_dispatch.go:10`.
The matching prepared-build-approval and prepared-review-shared-evidence canaries exercise those contracts.
Adapters continue to load canonical prose and own no source list.
Tests should decode TOON and derive shared fixture data from its existing owner.
An independent expectation needs its named biting mutation under the repository's duplication exception.

## Scope and final confirmation

Decision tickets 2 through 9 and ticket 11 own the accepted product predicates.
The scope covers build and review charges.
Write-spec behavior stays unchanged.
The change preserves exact review captures without adding reproducible collector execution across machines.

Ticket 10 owns the approved quota and explicit cleanup policy.
Ticket 12 owns the approved canonical TOON manifest and page-verification contract.
Ticket 13 owns the approved prepare, retrieval, verification, reuse, and current-validation flow.
Tickets 12 and 14 own the selected encoding and physical-container contract.

Global CLI pagination and guidance minimization are adjacent opportunities, not required outcomes of this charge change.
The selected design adds no previews, source summaries, or general session-state service.
The map retains hard choices as reviewer-owned tickets.

## Verification record and remaining evidence

- [x] The report separates current facts, measured results, inferences, and proposals.
- [x] Each load-bearing current claim cites a source or a historical source location.
- [x] The question graph has a synthesized answer per question.
- [x] The option table names consequences and their source basis.
- [x] The current flow diagram exposes preparation, delivery, and action as separate steps.
- [x] The write-spec inference does not claim an unobserved runtime failure.
- [x] The report retains the mismatch between current completeness prose and full-output behavior.
- [x] The report identifies its drift triggers and consuming map.
- [x] A disposable byte-delivery probe passed its stated checks; ticket 5 retains its limits.
- [ ] Bounded manifest delivery and the production verification interface still need implementation validation.
- [x] A disposable Linux lifecycle probe passed its stated checks; ticket 7 retains its limits.
- [ ] Cross-harness byte transport has not been tested.
- [x] Tickets 10 and 11 record the approved resource policy and collector scope.

The existing unit and system suites were read, not rerun as a claim that this defect is fixed.
The research delegates changed no files and reported a clean primary checkout.
Only the shaping assignment contains the map and this report.

The reviewer accepts both disposable prototype results with their stated limits.
All design tickets are resolved, and ticket 15 records final approval of the complete shaped scope.
The byte-delivery prototype reconstructs hostile fixtures through separately invoked bounded reads.
The lifecycle prototype checks sibling access, release, movement, corruption, concurrent publication, and cleanup exclusion.
No prototype result can approve the product contract.
After final confirmation and map validation, route to `$bench-write-spec`.

## Prototype-driven design question

Ticket 7 measured 52.07 ms to verify one 64 MiB artifact on this host.
A stateless reader that rehashes the whole artifact for every page multiplies that work by the page count.
A manifest can instead bind ordered source descriptors and fixed-size page digests.
Its identity binds page bytes transitively through those digests.
That design preserves exact-byte verification without whole-body rehashing on every retrieval.

Ticket 5 records the successful byte-delivery probe for this candidate.
Ticket 12 records the selected logical identity rules and TOON encoding.
Production validation must prove bounded manifest delivery and independent verification.

## Physical container comparison

The reviewer permits a Bench-specific format and asks whether it improves storage.
The logical manifest and the physical container have separate purposes.
The manifest defines the expected evidence identity and page membership.
The container locates exact stored bytes under that identity.
The container must not become another source-policy owner.

| Container | Benefit | Cost | Assessment |
|---|---|---|---|
| Directory with a manifest and one file per page | Each page has a simple filesystem address. | Each page adds a filesystem object, and publication must expose the complete directory atomically. | Feasible, but unnecessary object count for one immutable set. |
| ZIP with a manifest and page entries | Go supplies the reader, writer, random-access input, and stored or compressed entries. | Bench must still constrain names, entries, sizes, and semantic identity. Archive fields exceed the required evidence contract. | A credible standard-container alternative. |
| Minimal Bench packfile | One immutable file contains the manifest and contiguous exact source bytes. | Bench owns a small parser, format version, and hostile-length checks. | Recommended for the selected evidence contract. |

Go's ZIP reader accepts an `io.ReaderAt`, and its entry methods expose individual file content.
The standard package supports both stored entries and DEFLATE.
These APIs make ZIP viable; no benchmark establishes that a custom format is faster.
Source: [Go archive/zip documentation](https://pkg.go.dev/archive/zip).
The existing release archive reader has domain-specific extraction and validation semantics.
Source: `internal/releaseevidence/package_artifact.go:34`.
It does not supply a reusable charge-store owner.

Ticket 14 selects a packfile with three ordered regions:

1. A fixed header contains a format marker, a container version, and the manifest length.
2. The manifest contains the exact canonical bytes defined by ticket 12.
3. The payload contains exact uncompressed source bytes in manifest order.

Logical page offsets remain relative to their source.
The reader derives physical source offsets from the validated header, manifest length, and preceding source lengths.
A second independently authored offset registry is unnecessary.
The reader checks arithmetic, every range, and the exact expected file length before it returns bytes.
It refuses unsupported versions, overlapping or invalid page ranges, truncation, and trailing data.

A reader can request bytes at a specified offset through Go's `ReaderAt` contract.
Source: [Go io.ReaderAt documentation](https://pkg.go.dev/io#ReaderAt).
That mechanism supports a page read without a full payload scan.
The complete manifest and requested-page checks remain those defined by ticket 12.
Full verification still reads and verifies every required page and source.

The first container version uses no compression.
This choice needs no compression codec, compressed-offset index, or decompression limit.
A later container version can change physical representation without changing the identity of unchanged logical evidence.
Such a version must still reconstruct the exact manifest-bound bytes and preserve independent verification.
The selected container adds no general database, cross-artifact deduplication, or archive extraction workflow.

This comparison is a design assessment, not a production benchmark or compatibility result.
No custom-container prototype or production implementation has run.
The final spec must define the header bytes and format checks from one format owner.
Ticket 14 records approval of this physical-container choice.

## Manifest encoding comparison

The reviewer asks whether TOON reduces the token cost of the manifest.
Stored bytes consume model tokens only when they enter model context.
The chosen protocol exposes the manifest through bounded responses, so its text representation affects that cost.
The source bodies remain raw bytes inside the selected packfile.
A manifest encoding change does not reduce the required source content or replace the response bound.

TOON declares the fields of a uniform table once instead of repeating the keys for each row.
This shape fits source descriptors and page descriptors.
The standard also permits format options, so Bench must fix its own versioned canonical profile for digest stability.
Source: [TOON specification](https://github.com/toon-format/spec/blob/main/SPEC.md).
The benefit depends on the manifest shape and the tokenizer.
No manifest token-count comparison has run in this shaping phase.

Bench already emits flat TOON tables through one adapter.
The adapter fixes two-space indentation, comma delimiters, typed scalar handling, and trailing newlines.
It preserves a numeric-looking string as a string.
Sources: `internal/toon/toon.go:32`, `internal/toon/toon.go:74`, and `internal/toon/toon.go:90`.
The package tests cover quoting and mixed scalar types.
Sources: `internal/toon/toon_test.go:13` and `internal/toon/toon_test.go:155`.

The selected format uses that existing flat-table subset for the stored canonical manifest and its bounded delivery.
It does not assume support for newer upstream TOON features.
One schema owner fixes the block order, fields, types, row order, and supported manifest version.
The existing encoder remains the only owner of TOON escaping and cell representation.
An encoder change that alters canonical bytes requires explicit format compatibility handling.
The reader compares the hash of exact canonical bytes with the expected evidence identity.

The selected container contains a versioned binary header, a canonical TOON manifest, and exact uncompressed source bytes.
The bounded read protocol still identifies every fragment and verifies complete coverage.
A large scalar or manifest row still needs bounded fragments; table syntax alone does not establish the response limit.
Production validation must cover deterministic encoding, strict schema checks, hostile cells, numeric-looking digests, and manifest reconstruction.
The existing byte-delivery probe tested TOON body transport, but its local manifest used JSON.
It does not prove the complete selected TOON manifest format.
