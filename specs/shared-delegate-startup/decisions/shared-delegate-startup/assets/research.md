# Shared preparation for Bench delegates

Status: research complete; workflow adoption decisions remain in the decision tickets.

Recommendation: fold optional shared preparation into the existing delegation workflow.
Include spec authorship, delegated implementation, and independent review in one adoption decision.
Preserve source checks, retained authorship, and independent review judgment.
Prompt-cache savings remain unmeasured for Bench.

Evidence date: 2026-09-16.
Source baseline: `40d7f0d087047c44b11777cc96b71b549a2397f7`.
Consumed by: specs/shared-delegate-startup/decisions/shared-delegate-startup.md and its subsequent approved specification.
Drift: the cited workflow, charge, line, assessment, or native fork contracts change.
Retire when: the approved specification consumes these findings, or replacement research invalidates them.

## Question graph

1. Where can spec, build, and review delegates share preparation?
2. Which native capabilities preserve context, the selected line, and source isolation?
3. Which evidence distinguishes preparation savings from cache savings and total cost?
4. Which optional workflow fits the existing contracts?

Questions 1 and 2 are independent.
Question 3 combines their constraints with the assessment model.
Question 4 consumes those findings and the reviewer's decisions.

## Q1: Workflow opportunities

### Facts

Spec authoring already forks the approved decision context and inherits the invoking line.
The accepted source is a ready map, a reviewer-confirmed conversation, or a named reviewed artifact.
Source: `.agents/commands/bench-write-spec.md:20` and `.agents/commands/bench-write-spec.md:37`.

Build preparation reads the ticket, spec, build phase, delegate skill, and delegate procedure.
Its source reader compares current files with the pinned revision.
Each ticket requires current source identity and source bytes before action.
Source: `internal/preflight/charge.go:180`, `internal/preflight/charge.go:234`, and `.agents/commands/bench-implement-spec.md:44`.

Review preparation collects the diff, consumers, and coverage once per preparation attempt.
Every axis receives the same combined evidence identity.
Source: `internal/preflight/review.go:50` and `internal/preflight/review.go:155`.

Every review axis independently derives facts from current primary sources.
The review phase separates axis contexts and excludes the orchestrator and authors from delegated review roles.
Source: `.agents/skills/bench-craft-review/SKILL.md:16` and `.agents/commands/bench-review-implementation.md:44`.

### Inferences

| Phase | Opportunity | Limit |
| --- | --- | --- |
| Spec | Preserve approved context through the existing author fork. | One author normally means no sibling preparation to amortize. |
| Build | Share the spec, rules, and common source material before ticket assignments. | Each author validates its current ticket, source, and assignment. |
| Review | Share neutral raw material for one frozen source pair. | Each axis still derives its own facts and findings. |

These opportunities follow from the producers and consumers cited above.
A fresh preparation session cannot replace approved conversational decisions with an unreviewed summary.
Spec preparation must retain the existing decision-source authority.

## Q2: Native capabilities

### Facts

The active Codex native agent tool accepts full-history, recent-turn, and empty-context modes.
Full-history mode inherits the parent model and effort and rejects overrides.
Recent-turn mode accepts overrides but does not promise the complete prepared prefix.
Source: the active `collaboration.spawn_agent` contract, observed 2026-09-16.

The desktop task-fork tool copies completed history into a separate task.
It is a different surface from the native subagent tool.
Source: the active `mcp__codex_app__fork_thread` contract, observed 2026-09-16.

Bench's Claude guard rejects a declared model on a fork when the configured binding is enforceable.
This is a guard contract, not proof of live fork availability.
Source: `internal/lines/lines.go:475`.

Claude documents full history, system prompt, tools, model, and prompt-cache inheritance for forked subagents.
Named subagents start with a separate context and cache.
Claude also documents that a fork cannot spawn another fork.
Source: [Claude subagents](https://code.claude.com/docs/en/sub-agents#how-forks-differ-from-other-subagents), retrieved 2026-09-16.
Drift: Claude changes fork mode, inherited tools, or nested delegation limits.

A forked preparation helper therefore cannot serve as a portable parent for another fork batch.
The workflow needs an eligible parent session or the permitted fresh-delegate route.
This inference follows from the documented nesting limit.

Bench's headless adapters launch stdin-driven sessions and pass a resolved model.
They do not request a context fork or collect child usage records.
Source: `.bench/adapters/codex:14` and `.bench/adapters/claude:14`.

Conversation inheritance does not create a Bench assignment.
Write delegates require isolated worktrees; review delegates require source-bound venues.
Source: `.agents/skills/bench-craft-delegate/SKILL.md:86`.

### Unknowns

No live fork experiment ran during this investigation.
The active Codex schema does not expose the backend cache key or rendered request prefix.
Native capability and parent availability must be checked in the session that dispatches the work.
Child context, line, and assignment compatibility need a live check before a compatibility claim.

## Q3: Evidence and cost

### Facts

OpenAI caches matching rendered prompt prefixes.
Relevant settings, tools, eligible boundaries, and compaction affect reuse.
An existing session does not guarantee a cache hit.
Source: [OpenAI prompt caching](https://developers.openai.com/api/docs/guides/prompt-caching), retrieved 2026-09-16.
Drift: OpenAI changes caching or the selected model's request behavior.

Earlier Bench boot research measured empty Claude delegates and Codex exec sessions.
It did not compare prepared sibling forks with fresh native delegates.
Its Codex no-op decision concerns the tested exec profile change.
Source: `specs/delegate-boot-cost/spec.md:11`, `specs/delegate-boot-cost/spec.md:22`, and `specs/delegate-boot-cost/spec.md:283`.

Assessment records distinguish sessions, attempts, estimated cost, actual charges, and unknown measurements.
The native Codex importer recognizes input, cached input, and output counters.
Source: `internal/assessment/types.go:15`, `internal/assessment/types.go:44`, and `internal/assessment/harness.go:19`.

The rollout reader reports cumulative tokens and separate cached-input totals.
It does not attribute tokens to individual results or collect nested delegate records.
Each child requires its own record and counter-baseline check.
Inherited counters must not count parent preparation twice.
Source: `internal/harnesstranscript/codex.go:59` and `internal/harnesstranscript/codex.go:282`.

The estimate implementation prices uncached input, cached input, and output.
Other estimated charges can carry separately referenced costs.
Provider cache-write charges require correct attribution; ordinary input pricing does not automatically cover them.
Source: `internal/assessment/cost.go:25` and `internal/assessment/types.go:26`.

Claude documents request-level model, child identity, input, and cache counters through its monitoring interface.
This research did not configure monitoring or collect a native Claude record.
Source: [Claude monitoring](https://code.claude.com/docs/en/monitoring-usage#span-attributes), retrieved 2026-09-16.
Drift: Claude changes monitoring fields or producer semantics.

### Inferences

Fewer repeated reads do not prove fewer billed tokens.
A larger inherited context can increase later input cost despite a discounted prefix.
A different model can reduce preparation work while increasing output or repair cost.
These consequences follow from the cache contract and assessment cost components cited above.

### Tested results

This pass inspected source and public command contracts.
It ran no paid comparison and changed no production workflow.
No token, latency, quality, or currency improvement is claimed.

## Q4: Proposed workflow

### Proposals

| Candidate | Classification | Consequence |
| --- | --- | --- |
| Shared preparation in craft-delegate | Fold | One owner defines preparation, eligibility, and fallback. |
| Spec, build, and review integration | Fold | Existing phases consume that rule under their own authority. |
| New evidence collector | Skip | Preflight already owns shared review evidence. |
| New metrics store | Skip | Assessment already owns cost and provenance. |
| Universal fork default | Skip | The reviewer requested an optional workflow first. |
| Extra spec author or preparation session | Recommend against | The existing author fork already transfers approved context. |

Origin: the reviewer's current conversation.
Workflow citations appear in Q1; assessment citations appear in Q3.
These classifications are proposals, not authorization for production edits.

The proposed build preparation starts fresh before ticket-specific work.
The proposed review preparation contains the frozen pair and raw evidence, without author rationale or findings.
The preparation session receives no delegate findings before it completes that batch's forks.
Spec authorship retains its existing approved-source fork.

Each child receives its role, assignment, current source pin, and required supplement after the shared material.
An inherited assignment is context, not permission to write that checkout.
The child validates its own assignment before action.
Dependent tickets keep their existing prerequisite and retained-author rules.

If native inheritance cannot preserve the required contract, the proposed route uses a fresh delegate with explicit prepared evidence.
That fallback does not override spec authoring's mandatory fork contract.
A spec fork that cannot preserve its authority requires a capable-session handoff.
The coordinator records the selected route and any eligible line change.

## Decision ownership

The [decision map](../../shared-delegate-startup.md) indexes the reviewer choices.
Its tickets own scope, activation, model-change evidence, and spec-author treatment.
This report supplies evidence and proposals; it does not duplicate their answers.

## Contradictions and limits

The delegate skill describes ordinary delegates as having no conversation memory, while its fork rule permits inherited context.
The proposed change must distinguish these modes without weakening required charge inputs.
Source: `.agents/skills/bench-craft-delegate/SKILL.md:29` and `.agents/skills/bench-craft-delegate/SKILL.md:36`.

Review requires fresh contexts and independent rereads.
Shared raw material must not substitute for an axis's own derivation.
Source: `.agents/skills/bench-craft-review/SKILL.md:16`.

The staged bounded-charge-evidence spec proposes a different delivery protocol from the current code.
Implementation must consume the canonical preparation owner, without a competing transport.
Source: `specs/bounded-charge-evidence/decisions/bounded-charge-evidence.md:35` and `internal/preflight/charge.go:259`.

The compiled harness record leaves generic measurements unknown, while its explicit Codex record view can report supported observations.
Neither result proves measurements exist for a particular child session.
Source: `internal/harnesses/harnesses.go:217` and `internal/harnesses/command.go:100`.

## Verification record

- [x] Separate facts, inferences, tested results, and proposals.
- [x] Cite the primary evidence for each factual question.
- [x] Retain model inheritance, review independence, and source drift constraints.
- [x] Distinguish unmeasured savings from verified source behavior.
- [x] Verify the harness research return against primary source and the active tool contracts.
- [x] Retain reviewer choices in the decision map.

## Validation plan

1. Consume the confirmed workflow and model-change evidence rule from the decision tickets.
2. Specify native eligibility, permitted fallback, and source validation under existing phase authority.
3. Exercise the chosen native surface in a fresh session before claiming compatibility.
4. Check each child's context, assignment, and declared line.
5. Record preparation and child costs through existing assessment evidence.
6. Preserve unknown telemetry and separate estimated charges from actual charges.
7. Evaluate review quality and repair work alongside startup time.
8. Require separate evidence and reviewer approval before changing the default route.

## Q5. Can parent wait loops distort a fork cost comparison?

### Recommendation and scope

Treat repeated parent inference during quiet waits as a separate possible source of cost.
Do not attribute that cost to conversation forks without bounded telemetry.
This follow-up supports the spec's evidence limits; it authorizes no Codex configuration change or new benchmark.

### Source-backed reports

[Codex issue 35259](https://github.com/openai/codex/issues/35259) reports repeated model calls during agent and terminal waits.
Its author separates genuine usage deltas from copied history and unchanged snapshots.
The report distinguishes raw tokens and rate-card estimates from subscription usage.
Retrieved: 2026-09-16.

[Codex issue 37090](https://github.com/openai/codex/issues/37090) reports repeated compaction, source rereads, and status messages without corresponding progress.
Its author could not establish the account's token breakdown or the cause.
Retrieved: 2026-09-16.

These are first-person issue reports, not maintainer-confirmed diagnoses of this Bench run.
The linked comment was not available in the fetched issue page.
The user supplied its text, including a reported 25-minute wait workaround and separate rollout measurements.
The local installation's support for those settings remains unverified.
No fixed cache lifetime follows from these reports.

### Observations and inference

Several waits in this Bench session returned no new agent update.
That observation establishes repeated waits, not their billed token count or cause.
The user reports that the overall phase felt faster.
Scope, the review cap, inherited context, and wait behavior can each affect that impression.
Their separate contributions remain unknown.

### Validation boundary

No raw rollout telemetry or billing ledger was inspected for this follow-up.
No configuration change or live comparison ran.
Refresh this note when the issues, active runtime, or bounded local telemetry provide new evidence.
A future comparison should account for genuine parent usage deltas and child usage without copied-history duplication.
It should identify timeout-only turns separately and keep token totals distinct from subscription quota.
