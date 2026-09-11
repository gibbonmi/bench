# Benchmark orchestration frontier reassessment

This asset preserves historical benchmark evidence and proposals. Resolved decision tickets in the parent map supersede its policy recommendations.

Recommend Sol at high effort as coordinator, with one retained Astra delegate starting at medium effort.
Keep independent review and independent probes.
This trial can test the pairing, but changes coordinator model, delegate count, and implementation effort against the historical run.
Use matched runs before attributing an effect to one factor.

Status: research and unapproved recommendations.
Origin: benchmark learnings and the reviewer's request to assess lower-effort Astra.
Consumed by: the benchmark-orchestration decision map, tickets 4 through 11.
Drift: refresh when the frozen evidence, current owners, or reviewer decisions change.
Retire when: confirmed decisions and the resulting specs replace these proposals.

The reviewer requires flexible implementation effort.
The research does not select a new default, authorize a benchmark run, or approve a delivery split.
Tickets 1 through 3 remain the confirmed decision source.

## Evidence identity and scope

The primary checkout and shaping assignment both start at `cf2caa6d5c195c6b6ed5852fb4478902d1d99fbd`.
The shaping assignment is `benchmark-orchestration-shape`.
Its branch is `bench/assign/ce8853ba95d83532380a0580f5297d1a/42976835b76775651ce813f037d5c551`.

| Candidate | Frozen commit |
|---|---|
| Delegated Astra | `0430539b796572c387ce02d7d85a6b89c61e8ffe` |
| Solo Astra | `17e888bb12d916ed4200301b47f3a88678fa46f6` |
| Delegated Sol | `8278fc9d6bf5b7b0a0095dd3c148369ee24d9019` |
| Delegated Terra | `592d67b051573bc687acaf5d6d14a5966bb91ff2` |

Candidate code citations below name immutable Git objects and source lines.
Current-owner citations refer to the baseline above.
The primary benchmark reports remain authoritative.
The map's earlier snapshot changes their link destinations for local access.

The factual question graph has six nodes, one for each open Round 2 ticket.
Branch structure informs grouping and authority.
Run evidence informs progress, assurance, and cost.
Current owners determine which proposed capabilities already exist.
These answers jointly inform the delivery split and later adoption criteria.

## 4. What do the branches establish about grouping?

### Facts

The delegated Astra run already retained one worker for tickets 2 through 5.
Its record reports retained native thread identities across continuations.
Sol also reused worker 2 across those tickets.
Terra reused its second worker across tickets 3 through 5 after ticket 2.
Sources: [Astra workers](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/astra/result.json:15), [Astra retention](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/astra/result.json:435), [Sol ticket evidence](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/sol/result.json:186), [Terra ticket evidence](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/terra/result.json:51).

| Candidate | Code-backed observation | Consequence for shaping |
|---|---|---|
| Solo Astra | The fingerprint includes staged index entries. A separate reset-ref sweep repeats the existing sweep mechanics. | Retained context can coexist with correct state binding and duplicated policy. |
| Delegated Astra | Reset refs use the shared sweep. Tests cover non-active recorded states and refs moved after listing. | Delegation can preserve shared owners and stronger regression proof. |
| Sol | Both modes share checkout facts. Lock repair unconditionally unlocks an already-unlocked registration. | A shared owner still needs tests for each state that enters it. |
| Terra | Restore uses a separate planner that omits ref and lock identity, nested-state checks, and the cleanup lock. | Related modes can diverge even when one worker retains their tickets. |

Code sources: `17e888bb:internal/worktree/reset.go:183`, `17e888bb:internal/worktree/reconcile.go:110`, `0430539b:internal/worktree/reconcile.go:66`, and `0430539b:internal/worktree/resume_reconcile_test.go:133`.
Sol sources: `8278fc9d:internal/worktree/reset.go:197`, `8278fc9d:internal/worktree/reset_apply.go:61`, and `8278fc9d:internal/worktree/reset_restore.go:97`.
Terra sources: `592d67b0:internal/worktree/reset_restore.go:15`, `:68`, and `:108`.
The [Astra comparison](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-astra-code-review.md:41) and [Sol/Terra assessment](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-two-arms.md:92) retain the external probe results.

### Inference and proposal

These runs do not compare fresh authors per ticket against retained authorship.
They cannot isolate context continuity, coordinator overhead, effort, or assurance cost.
Group related behavior for context reuse, but preserve each ticket's acceptance and commit boundary.
Let the coordinator regroup unstarted work within approved scope and dependencies.
Require a checkpoint and explicit transfer before changing an active author's assignment.
Keep model and effort selection separate from group membership.

## 5. What write authority already exists?

### Facts

Preflight already proposes fixture and registry closure additions.
It names additions outside the spec fence and requests an approved ordering edge for a conflicting pair.
It does not apply those proposals.
The focused closure test confirms that a proposal alone does not make the subsequent charge pass.
Sources: [proposal owner](/home/mgibs/workspace/bench/internal/preflight/proposal.go:19), [ordering owner](/home/mgibs/workspace/bench/internal/preflight/proposal.go:72), [closure test](/home/mgibs/workspace/bench/internal/preflight/proposal_test.go:65).

FT293 still leaves closure-only overlap and test-file derivation for the reviewer.
The current implementation does not settle that policy decision.
Source: [FT293](/home/mgibs/workspace/bench/roadmap/FT293.md:1).

### Proposal

Give a group the union of its approved ticket writes, within the spec's outer fence.
Treat newly derived closure as a proposal until the applicable authority approves it.
Serialize actual conflicting writers and preserve existing dependency edges.
Do not convert closure-only overlap into a new dependency policy through this recommendation.
Do not treat the group's full write set as permission to skip a ticket checkpoint.

## 6. What counts as a stall?

### Facts

Terra's early attempt counts included incremental verification batches.
Later accounting distinguished coherent implementation-and-verification attempts.
Source: [Terra accounting](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/terra/result.json:40).

The current line discipline already classifies inherited and diff-owned failures.
It distinguishes insufficient model capability from a seam or spec contradiction.
Source: [line discipline](/home/mgibs/workspace/bench/.agents/skills/bench-craft-line/SKILL.md:66).

### Proposal

Judge progress against a named acceptance predicate and new discriminating evidence.
Count a completed implementation-and-verification attempt, not each tool call or expected TDD failure.
After repeated non-progress, require a changed hypothesis and a concrete next check.
Use independent diagnosis, a smaller charge, different effort, or a context transfer when the evidence supports it.
Do not make a model increase the automatic response to every stall.
Preserve reviewer decisions, external blockers, and explicit budgets as stop conditions.

The benchmark does not establish an optimum restart threshold.
Recommend reassessment after two completed attempts without progress or new discriminating evidence.
This is a proposed initial trigger, not a measured optimum.
The reviewer still owns this trigger and the authority for the resulting response.

## 7. What assurance can current evidence establish?

### Facts

The current test owner distinguishes failed tests, build failure, no executed test, refusal, and interruption.
It waits for the test child to terminate.
Its public output includes package outcomes, failure rows, and skip reasons.
Sources: [outcome owner](/home/mgibs/workspace/bench/internal/testreport/outcome.go:80), [terminal checks](/home/mgibs/workspace/bench/internal/testreport/command.go:331), [rendered evidence](/home/mgibs/workspace/bench/internal/testreport/testreport.go:165).

The probe owner checks the baseline, restores the subject, and classifies compile failure or no executed test as invalid.
Its selection row reports the test count.
However, any observed test failure yields `bit`.
That label alone cannot prove that the intended assertion failed instead of setup or another test.
Source: [probe owner](/home/mgibs/workspace/bench/internal/probe/probe.go:201).

Preflight collects one shared diff, consumer set, and coverage set per preparation attempt.
It provides source hashes, rejects incomplete consumer evidence, and names full retrieval for compact output.
Each command recollects the evidence.
The full response has no item selector for retrieving only the unread source.
Sources: [review collection](/home/mgibs/workspace/bench/internal/preflight/review.go:47), [packet renderer](/home/mgibs/workspace/bench/internal/preflight/charge.go:72).

Prepared review inputs do not prove that a reviewer executed the review.
Current review guidance writes no artifact for a clean review.
Landing binds a source pair and checks its fence, but that authorization function does not inspect independent reviewer returns.
Sources: [clean-review artifact policy](/home/mgibs/workspace/bench/.agents/commands/bench-review-implementation.md:183), [source authorization](/home/mgibs/workspace/bench/internal/preflight/gather.go:191), [landing call](/home/mgibs/workspace/bench/internal/worktree/land.go:303).

### Proposal

Completion evidence must distinguish author work, independent verification, and the three independent review axes.
Bind each result to the source it examined and retain its terminal outcome.
Preserve review coverage through the initial reviewed pair and subsequent repair reviews.
An uncovered change leaves assurance incomplete.
Do not restart full discovery merely because a repair moves the tip.
The existing [repair-review policy](/home/mgibs/workspace/bench/.agents/commands/bench-review-implementation.md:33) defines that scope.

Missing evidence, invalid probes, skips, and observed defects must remain distinct.
The gate remains the single oracle.
The persistence owner and the effect of missing assurance still require a decision.
FT200's broader landing policy remains open, including its prohibition on a second done-oracle.
Source: [FT200](/home/mgibs/workspace/bench/roadmap/FT200.md:3).

## 8. Which outcomes can ship independently?

### Proposal

Consider four primary delivery outcomes:

1. Retained authorship with bounded group authority and ticket checkpoints.
2. Explicit continuation with progress and reassessment rules.
3. Truthful completion evidence with independent assurance and source coverage.
4. Reproducible assessment records and advisory comparisons under FT231.

Retained authorship and uncapped continuation are independently useful options.
Neither requires the other to operate.
The reviewer must confirm this split.

Keep focused-test diagnostics, assertion-specific probe evidence, and selective evidence retrieval as candidate extensions to their current owners.
Specify each only after a concrete remaining gap is reproduced and its outcome is approved.
Do not create a replacement test runner, review engine, lifecycle, or receipt system from this audit.
The observation-to-spec boundary follows the original [recommendation scope](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/orchestration-recommendations.md:71).

## 9. Is another lower-effort Astra run worth the cost?

### Recorded facts and arithmetic

Delegated Astra's estimated cost was $81.79 under the saved pricing assumptions.
The coordinator accounted for $43.54, or 53.2 percent.
Implementation workers accounted for $35.27, or 43.1 percent.
The original independent review accounted for $2.98.
These are estimates, not observed charges.
Source: [audited role costs](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-evidence/three-arm-audit.json:262).

The two workers recorded 19,638 reasoning output tokens.
At the saved $50 per million output tokens, those tokens cost about $0.98.
Removing all of them would save that amount only if all other usage stayed fixed.
Lower effort can also change tool calls, context size, latency, and repair work.
Those effects remain unmeasured.
Source: [audited session usage and saved rates](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-evidence/three-arm-audit.json:1).

### Coordinator pairing evidence

All three delegated benchmark arms used Astra as coordinator.
The Sol arm changed the implementation workers, not the coordinator.
The saved runs therefore provide no direct Sol-orchestration result.
Source: the coordinator session entries in the [three-arm audit](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-evidence/three-arm-audit.json:1).

Repricing the historical Astra coordinator's identical usage at the saved Sol rates gives $17.42 instead of $43.54.
The arithmetic difference is $26.12, about 32 percent of that run's total estimate.
This is a price-only scenario, not a predicted saving.
Sol can use different amounts of context, require more interventions, or miss an evidence defect.
The saved audit supplies the token counts and rate assumptions.

The official model pages confirm support for the proposed effort settings and the cited standard token rates.
They do not validate this orchestration pairing.
Sources: [OpenAI Docs: Sol](https://developers.openai.com/api/docs/models/gpt-5.6-sol), [OpenAI Docs: Astra](https://developers.openai.com/api/docs/models/gpt-6-astra).
The pages were retrieved on 2026-09-11 UTC.
Refresh this price-only scenario if the rates or execution tier change.

### Recommendation

A single retained Astra delegate at medium effort is worth a bounded exploratory trial.
It tests whether a simpler implementation arrangement can preserve quality with less supervision.
Several lower-effort delegates are a lower priority for this coherent task.
The historical cost suggests that coordinator work deserves at least as much attention as worker reasoning.

Pin the pilot coordinator to `gpt-5.6-sol` at high effort.
Pin one retained implementation delegate to `gpt-6-astra`, starting at medium effort.
Sol owns charges, evidence checks, integration, and decision routing.
Astra owns implementation and related repairs.

Keep reviewer model and effort unchanged during the pilot.
Keep per-ticket checks, independent probes, and all three independent review axes.
Record complete cost, elapsed time, defects, repairs, and any escalation.
If the delegate needs higher effort, retain the medium-effort outcome before recording the recovery.
That run tests an adaptive policy, not a pure medium-effort condition.

A single pilot establishes viability only.
The existing multi-delegate run is a historical reference, not a matched effort control.
Compare single-delegate medium and high effort under identical conditions if the pilot justifies further trials.
Use repeated held-out tasks before changing a default.
FT231's no-Bench, current-Bench, and one-capability comparison remains required for causal kit-capability claims.
Source: [FT231 comparison and ownership boundaries](/home/mgibs/workspace/bench/roadmap/FT231.md:30).

The repeated FT311 task is an authoring case, not held-out evaluation.
Its ignored-path collision still needs the separate behavior decision before it can grade successful completion consistently.
Do not change the frozen candidates to prepare a trial.
Source: [unresolved collision](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-astra-code-review.md:48).

## Additional map uncertainty

Ticket 10 asks where missing or stale assurance blocks a completion claim.
It depends on the assurance contract and delivery split.
It must distinguish phase reporting from FT200's landing chokepoint decision.

Ticket 11 asks what the assessment record retains and what trial rule permits adoption.
It depends on the delivery split and adoption criterion.
It must settle storage, granularity, retention, cost cutoffs, cases, repetition, and quality tolerances within FT231's ownership boundaries.

Group transfer and restart rules remain part of tickets 4 and 6.
Model and effort flexibility is a confirmed reviewer preference, not permission to ignore budgets or alter benchmark conditions silently.

## Verification record and validation plan

The coordinator read the three authoritative reports, their linked snapshot, arm records, current owners, and the cited candidate code.
Two read-only delegates checked the candidate pairs.
The coordinator rechecked the load-bearing code citations and corrected one reversed test attribution before synthesis.
This assessment separates facts, inferences, prior probe results, and unapproved recommendations.
The tables expose the material comparisons without requiring a diagram.

Existing focused tests passed for probe baseline handling, compile-failure classification, unmatched tests, and test counts.
Existing preflight tests passed for stable pinned output, compact retrieval, shared review evidence, and closure proposals.
Both runs reported no skips.
The initial sandboxed probe test could not write the build cache.
Its authorized retry used the normal Bench runner and passed.

No candidate changed, and no new candidate gate or performance trial ran.
This is a targeted shaping assessment, not a replacement formal review of all four branches.
The native response ledgers were not reconstructed again.
Actual charges, causal savings, ideal effort, and optimum group size remain unknown.

Resolve the reviewer frontier before spec authoring.
Reproduce each surviving CLI gap at its existing owner before adding behavior.
Freeze the next experiment's task, authority, assurance, effort policy, and accounting before execution.
Run repeated matched comparisons before a default decision.
