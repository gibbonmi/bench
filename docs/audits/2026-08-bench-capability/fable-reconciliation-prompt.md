# Bench Audit Reconciliation and Final Action Plan

## Role

You are the final independent adjudicator for two deep architectural audits of **Bench**, an AI-assisted software-engineering harness.

You are running inside a **clean Bench audit worktree**.

The independent audit inputs are under:

```text
audit/inputs/
├── sol/
│   ├── report.md
│   ├── evidence.md
│   └── questions.md
└── opus/
    ├── report.md
    └── evidence/
        └── ...
```

If the actual filenames differ slightly, discover them rather than assuming they are missing.

The Sol audit was produced independently in Codex using GPT-5.6 Sol at extra-high effort.

The Opus audit was produced independently in Claude Code using Opus 5 at extra-high effort.

You are running in Claude Code using **Fable at high effort**.

Your job is **not to summarize the two reports**. Your job is to reconcile their findings, verify consequential claims against the Bench repository, resolve important disagreements through evidence or controlled experiments, audit the existing roadmap and bug backlog, determine which principles should survive and how they should be enforced, produce one final prioritized action portfolio, and choose exactly one highest-leverage next implementation ticket.

The repository is the source of truth. The reports are claims. Evidence folders are supporting provenance. Questions files are unresolved hypotheses.

---

# 1. Core mission

Determine:

> **What is the smallest coherent set of changes that will make Bench a clearer, more reliable, more evidence-driven, more resumable, more context-efficient, and more model-independent software-engineering runtime without destroying workflow behavior that has already demonstrated value?**

This is architectural reconciliation and prioritization, not implementation.

Do not produce a third sprawling audit. The discovery work has already been done. Reduce uncertainty, adjudicate, simplify, and prioritize.

---

# 2. Known practitioner evidence

## Deep Matt Pocock integration

Bench already has a **deep integration with Matt Pocock's engineering skills**, and that integration has materially improved workflow consistency.

Do not evaluate the Pocock integration as unexplored prior art, a possible future dependency, or an optional alternative architecture. Treat it as an existing procedural foundation whose value has been observed.

Determine:

- which integrated behaviors are causally valuable;
- which adaptations are intentional;
- where Bench extends the upstream procedures;
- where Bench duplicates or contradicts them;
- which behaviors require model judgment;
- which behaviors should become deterministic Bench controls;
- whether consistency gains survive in both Claude Code and Codex;
- how to preserve the gains while simplifying Bench.

Do not recommend removing or replacing an integrated skill merely because a cleaner abstraction could theoretically be designed. Preserve demonstrated behavior unless evidence shows a simpler mechanism retains the benefit.

## `diagnosing-bugs`

The maintainer has repeatedly used Matt Pocock's `diagnosing-bugs` skill and observed:

```text
agent enters repetitive repair loop
    ↓
diagnosing-bugs is invoked
    ↓
agent establishes executable feedback
    ↓
agent stops applying speculative repairs
    ↓
agent narrows the failure to the root cause
    ↓
agent often produces a correct one-shot fix
```

Treat this as meaningful repeated practitioner evidence. It is not yet a controlled benchmark, but it is stronger than a speculative claim.

Any recommendation that would remove, dilute, obscure, or replace this behavior carries a high burden of proof.

Investigate *why* it works. Candidate mechanisms include:

- requiring a red-capable feedback loop;
- reproducing before repairing;
- minimising the failing case;
- separating observation from theory;
- requiring falsifiable hypotheses;
- instrumenting instead of guessing;
- carrying failure information into the next attempt;
- validating the original symptom after repair;
- preventing blank retries;
- forcing a transition from speculation to observation.

Do not assume all mechanisms are equally important. Determine what the current Bench integration actually preserves.

---

# 3. Model continuity requirement

At the beginning of the run, record the model identity displayed by the harness if observable. At the end, record it again.

If the harness visibly reports that execution switched from **Fable** to another model:

1. stop substantive reconciliation;
2. do not produce the final architectural verdict;
3. write `audit/fable-high/INVALID-RUN.md`;
4. record when the switch became visible;
5. record what work had been completed before the switch;
6. clearly state that the run cannot be treated as a clean Fable adjudication.

Do not claim to know the active model when the harness does not expose it.

---

# 4. Repository safety

Allowed:

- read/search files;
- inspect git history;
- invoke existing Bench commands;
- run local tests and builds;
- inspect generated artifacts;
- construct small controlled fixtures;
- create temporary local experiments;
- write audit outputs under `audit/fable-high/`.

Not allowed:

- implement the recommended architecture;
- rewrite Bench skills as part of the audit;
- silently fix failing tests;
- rewrite the live roadmap;
- mutate production or remote systems;
- make destructive or irreversible changes;
- commit implementation changes.

Temporary experiments must be cleaned up.

At completion run:

```bash
git status --short
```

Audit documents under `audit/fable-high/` are expected. Explain any other modification.

---

# 5. Establish the baseline

Begin with:

```bash
git status --short
git rev-parse HEAD
git log -1 --oneline
```

Write the result to:

```text
audit/fable-high/baseline.md
```

Inspect the Sol and Opus audit inputs and determine whether they identify the commit audited.

Classify each relationship as:

- `EXACT MATCH`
- `LIKELY MATCH`
- `MISMATCH`
- `IMPOSSIBLE TO ESTABLISH`

If reports examined different commits, separate findings stable across versions, findings potentially caused by repository drift, and findings requiring fresh verification.

---

# 6. Evidence hierarchy

Use this hierarchy:

```text
fresh executable observation against current repository
    ↓
current implementation + current tests
    ↓
current implementation
    ↓
prior auditor evidence
    ↓
repository documentation
    ↓
auditor conclusion
    ↓
external research analogy
```

Implications:

- Documentation describing behavior does not prove it.
- A report saying a command is noisy does not prove the current command is noisy.
- A passing test proves only its actual coverage.
- Prior evidence is useful, but cheap observations should be reproduced.
- External research can motivate a hypothesis but cannot prove a Bench-specific conclusion.

Classify the basis of major recommendations as:

- `OBSERVED`
- `INFERRED`
- `PRACTITIONER EVIDENCE`
- `HYPOTHESIS`
- `RESEARCH-INSPIRED`

---

# 7. Treat report artifacts differently

For Sol:

- `report.md` = conclusions
- `evidence.md` = supporting observations/provenance
- `questions.md` = unresolved hypotheses

For Opus:

- main report = conclusions
- files under `evidence/` = supporting artifacts

If Opus evidence mixes raw command output, tests, notes, fixtures, and intermediate conclusions, distinguish them before relying on them.

Do not flatten all auditor artifacts into one authority level.

---

# 8. Reconciliation ledger

Build:

```text
audit/fable-high/reconciliation-ledger.md
```

For every material finding record:

```yaml
id:
topic:
sol_position:
opus_position:
relationship:
prior_evidence:
fresh_repository_evidence:
experiment:
result:
classification:
confidence:
final_disposition:
```

Use:

- `REPRODUCED`
- `COMPATIBLE`
- `CONTRADICTED`
- `PARTIALLY SUPPORTED`
- `UNSUPPORTED`
- `STALE`
- `UNINVESTIGATED`
- `NOT ACTIONABLE`

Agreement raises verification priority; it does not replace verification.

---

# 9. Do not vote between auditors

Never reason:

```text
Sol says X.
Opus says X.
Therefore X is correct.
```

Instead:

```text
Sol says X.
Opus says X.
Therefore X is a high-priority claim to verify.
```

When they disagree, ask:

> What is the smallest observation that distinguishes the competing claims?

Then obtain it.

---

# 10. Empirical escape: run the dang test

Whenever uncertainty appears:

```text
uncertain
    ↓
can repository inspection answer it?
    ↓ yes
inspect
    ↓
still uncertain?
    ↓
can a command, test, compiler, search, fixture,
runtime probe, or controlled experiment distinguish
the active hypotheses?
    ↓ yes
RUN IT
```

Do not spend paragraphs discussing what a safe local command probably does when the command can show you.

Do not infer CLI output from source when invocation is available.

Do not infer gate behavior when a canary can demonstrate it.

Do not infer whether final-check duplicates gate only from prose.

Do not reason indefinitely about a reported bug if a local reproduction can settle it.

Continued reasoning is justified only when it creates at least one of:

- a new falsifiable hypothesis;
- a new constraint;
- a new dependency;
- a new discriminating experiment;
- a new synthesis across verified evidence.

If reasoning merely restates uncertainty:

> **Stop and run the dang test.**

This does not authorize unsafe execution. Prefer read-only inspection, narrow local tests, isolated builds, temporary fixtures, and non-destructive CLI calls. Do not perform remote, credentialed, destructive, production, or irreversible actions.

---

# 11. Resolve consequential disagreements first

Prioritize disagreements capable of changing architecture or roadmap, such as:

- continue incrementally vs consolidate vs strangler migration vs rewrite;
- new work-state layer vs strengthen existing state;
- `/bench` as front door vs existing entry;
- context compiler needed vs current mechanisms sufficient;
- gate and final-check complementary vs duplicative;
- claim state controls behavior vs merely reports it;
- more orchestration vs less orchestration;
- Pocock integration wrapped vs absorbed vs preserved;
- failure inheritance missing vs represented indirectly.

For each:

1. state competing claims neutrally;
2. identify evidence needed;
3. inspect repository;
4. run the smallest safe experiment where practical;
5. record result;
6. update ledger;
7. determine whether action follows.

---

# 12. Canonical entry-point audit

Bench must have a clear public entry point. Do not assume the answer is `/bench`.

Evaluate:

```text
/bench
    ↓
inspect durable Bench/repository state
    ↓
recommend or resume correct workflow
    ↓
load only relevant capability
```

Potential relationship:

```text
/bench
    = public front door

what-next
    = orientation/routing capability

shape
build-spec
implement
diagnose
    = direct expert side doors
```

Determine:

- current real entry point;
- discoverability;
- Claude Code/Codex consistency;
- whether state is inspected before questions;
- resume behavior;
- whether routing uses actual repo state;
- context load;
- duplication with existing capability;
- whether entry behavior belongs in a harness skill, CLI, or both.

Return exactly one:

- `ADOPT /bench`
- `ALIAS /bench TO EXISTING ROUTER`
- `KEEP CURRENT ENTRY`
- `REDESIGN ROUTING FIRST`
- `OTHER`

A clear entry point is required even if `/bench` is rejected.

---

# 13. Pocock integration audit

For every materially integrated Pocock-derived capability determine:

```text
upstream/conceptual source
    ↓
Bench adaptation
    ↓
Bench extension
    ↓
workflow consumers
    ↓
deterministic controls
    ↓
observed value
    ↓
missing benchmark
```

At minimum inspect relevant implementations of:

- diagnosing bugs;
- TDD;
- grilling/shaping;
- specification creation;
- ticket decomposition;
- implementation;
- code review;
- delegation/handoff;
- domain modeling where present.

Classify each as:

- `PRESERVE`
- `STRENGTHEN`
- `WRAP WITH DETERMINISTIC CONTROL`
- `SIMPLIFY`
- `BENCHMARK`
- `REPLACE`
- `REMOVE`

Do not collapse the whole integration into one verdict.

---

# 14. `diagnosing-bugs` and repair-loop escape

Reconstruct its actual Bench control loop and determine whether it resembles:

```text
trigger
    ↓
executable feedback loop
    ↓
reproduction
    ↓
minimisation
    ↓
falsifiable hypotheses
    ↓
instrumentation
    ↓
root-cause identification
    ↓
repair
    ↓
regression validation
```

Investigate why it breaks repair loops.

Distinguish among:

- prompt wording;
- enforced sequencing;
- executable feedback;
- reduced hypothesis space;
- failure-state persistence;
- changed tool-use behavior;
- model/harness activation;
- interaction with other Bench controls.

Determine what should generalize across Bench. Do not automatically generalize the entire skill.

Evaluate this candidate principle:

> When a bounded safe executable observation can distinguish active hypotheses, execute it before continuing speculative repair.

Determine whether it belongs in core doctrine, diagnosis skill, implementation skill, work state, CLI, gate, instrumentation, harness adapter, model profile, or benchmark suite.

Design a benchmark:

```text
A. model without diagnosing-bugs
B. model with current diagnosing-bugs integration
C. model with diagnosing-bugs plus proposed Bench control
```

Measure:

- repair-loop count;
- distinct failed approaches;
- steps to reproduction;
- steps to first discriminating observation;
- root-cause accuracy;
- first-fix success;
- regression coverage;
- duplicate commands;
- unsupported causal claims;
- tokens/tool calls.

---

# 15. Work-state audit

Evaluate whether Bench needs, already has, or partially has an authoritative representation equivalent to:

```text
Goal
Core
Verified
Open
Next
```

Do not require those names.

Inventory current state mechanisms, lifetimes, authority, duplication, mutation ownership, freshness, resumability, evidence linkage, visible uncertainty, and whether another model can continue without reconstructing conversation history.

Choose exactly one:

- `NO NEW WORK STATE`
- `UNIFY EXISTING STATE`
- `MINIMALLY EXTEND EXISTING STATE`
- `INTRODUCE CANONICAL WORK STATE`

Do not create another source of truth just because the abstraction is attractive.

---

# 16. Context-compilation audit

Treat context as a product, not a pile of files.

Inspect:

- initial context size;
- progressive disclosure;
- repeated file reads;
- duplicate searches;
- shared invariants;
- seam-specific context;
- delegate packets;
- stale copied truth;
- CLI output entropy;
- resume context.

Determine whether a first-class context compiler is justified by observed problems.

---

# 17. Claims/evidence audit

Determine whether claims are a **control system** or a **reporting system**.

Inspect provenance, replay, freshness, hashes, dependencies, invalidation, access paths, bypasses, output contracts, exit codes, ergonomics, and machine-readable behavior.

For each state ask:

```text
verified → what becomes permitted?
unverified → what remains blocked?
stale → what is replayed or invalidated?
disproven → what downstream decision changes?
```

A state that changes no behavior is monitoring, not control.

---

# 18. Gate / review / tests / CI / final-check

Produce:

| Mechanism | Stage | Inputs | Deterministic? | Semantic judgment? | Coverage | Failure consequence |
|---|---|---|---|---|---|---|

Determine actual responsibilities, overlap, independence, missing coverage, stale inputs, and whether repeated review is genuinely independent.

Run distinguishing experiments where safe.

Do not merge mechanisms because prose overlaps. Do not preserve them because they already exist.

---

# 19. Delegation and orchestration

Evaluate:

- orchestrator responsibilities;
- delegate boundaries;
- ownership fences;
- shared invariants;
- broadcast state;
- independent discovery;
- failure inheritance;
- retries;
- escalation;
- integration;
- duplicated context;
- model routing.

For every agent boundary ask:

> What independent information, authority, or judgment does this boundary provide?

More agents are not inherently better.

---

# 20. Failure inheritance and recovery

Determine what survives:

- failed repair;
- delegate failure;
- session interruption;
- context compaction;
- model swap;
- Claude outage;
- Codex continuation;
- worktree restart.

Inspect whether later attempts inherit failed command, observed failure, diagnosis, ruled-out hypotheses, modified files, validation coverage, open uncertainty, and next discriminating action.

Prefer compact structured state over reflective essays.

---

# 21. Prose versus deterministic control

For important behavioral rules implemented only in prose ask:

```text
Does this require model judgment?
Can a schema enforce it?
Can a command observe it?
Can a hook reject it?
Can a gate validate it?
Should it remain guidance?
```

Do not indiscriminately move judgment into code. Do not ask models to remember deterministic invariants.

---

# 22. Consolidate versus restart

Compare:

## A. Continue incrementally

Keep current architecture and make targeted fixes.

## B. Consolidate in place

Preserve behavior while unifying duplicated state, prose, and contracts.

## C. Strangler migration

Define a new architectural spine and migrate existing behavior one vertical slice at a time.

Possible hypothesis:

```text
canonical entry
    ↓
deterministic status
    ↓
authoritative current work state
    ↓
minimal context
    ↓
practice skill
    ↓
executable observation
    ↓
evidence
    ↓
gate
    ↓
checkpoint
```

## D. Clean-slate rewrite

Rebuild Bench from a new foundation.

Clean-slate carries a high burden of proof because Bench already contains behavior that improved workflow consistency.

Evaluate regression risk, loss of behavioral tuning, migration cost, simplification opportunity, benchmarkability, compatibility, state migration, and impact on Pocock-derived behavior.

Choose exactly one primary strategy. A bounded subsystem rewrite does not imply a whole-project rewrite.

---

# 23. Preserve behavior before replacing structure

For every major redesign identify:

```text
current valuable behavior
    ↓
underlying invariant
    ↓
current mechanism
    ↓
proposed mechanism
    ↓
how preservation will be demonstrated
```

Do not replace effective behavior only to achieve conceptual uniformity.

---

# 24. Roadmap and defect portfolio audit

The existing Bench roadmap is part of the audit.

Do not produce a final action plan without reconciling it against current roadmap, specs, tickets, bug lists, deferred findings, and planned capabilities.

## Discover roadmap-like artifacts

Locate:

- roadmap files;
- specs;
- ticket directories;
- TODO documents;
- deferred review findings;
- known-bug lists;
- ADR follow-ups;
- unimplemented acceptance criteria;
- proposed commands;
- planned workflows;
- architectural backlog notes.

Build:

```text
audit/fable-high/roadmap-inventory.md
```

## Verify every roadmap entry

Determine:

- original problem;
- whether it still exists;
- whether it can be observed;
- whether work is partially implemented;
- whether already complete;
- duplication;
- valid dependencies;
- architectural assumptions;
- whether affected subsystem survives;
- observable value;
- demonstrable completion;
- implementation commitment vs experiment.

## Required disposition

Assign exactly one:

- `KEEP`
- `REWRITE`
- `MERGE`
- `SPLIT`
- `FIX-NOW`
- `FIX-BEFORE-MIGRATION`
- `FIX-WHEN-TOUCHED`
- `EXPERIMENT`
- `DEFER`
- `DONE/ARCHIVE`
- `SUPERSEDED`
- `DELETE`
- `UNCONFIRMED`
- `WONT-FIX`

Write:

```text
audit/fable-high/roadmap-dispositions.yaml
```

## Bug-fix economics

A bug is not automatically worth fixing.

For every reported defect:

1. attempt safe reproduction;
2. identify the smallest executable feedback loop;
3. determine invariant/workflow affected;
4. determine whether component survives target architecture;
5. compare fix, replace, delete, defer, document;
6. identify regression discriminator;
7. assign disposition.

Use `diagnosing-bugs` or its integrated Bench adaptation where useful.

Do not continue theorizing after a safe local discriminator exists.

An unreproduced bug becomes `UNCONFIRMED`.

A bug in a subsystem recommended for deletion should normally not be repaired unless it threatens integrity/security, blocks migration, prevents baseline measurement, or affects a supported path.

Write:

```text
audit/fable-high/bug-triage.md
```

## Replacement roadmap

Produce:

```text
audit/fable-high/proposed-roadmap.md
```

It must reflect target architecture, preserve demonstrated consistency, contain only active accepted work, order by dependency/leverage, distinguish experiments from commitments, contain no more than twelve active items, include exactly one immediate next ticket, and exclude completed/superseded history from active work.

The proposed roadmap and final action set must be the same portfolio.

Do not modify the live roadmap during this run.

---

# 25. Principle classification and enforcement

Classify each candidate principle as:

- `CORE INVARIANT`
- `PROCEDURAL DEFAULT`
- `HEURISTIC`
- `HARNESS ADAPTER RULE`
- `MODEL PROFILE`
- `EXPERIMENT`
- `REJECTED DOCTRINE`

For each surviving principle record:

```yaml
id:
statement:
classification:
scope:
failure_mode_prevented:
owner:
enforcement:
evidence:
regression_test_or_benchmark:
exceptions:
```

Write:

```text
audit/fable-high/principle-control-matrix.md
```

Evaluate these candidates as hypotheses:

1. Uncertainty may remain unresolved, but it must remain visible.
2. Observation and interpretation are different artifacts.
3. When a safe executable observation can distinguish active hypotheses, prefer observation over additional speculation.
4. Deterministic invariants belong in CLI/schema/hook/gate rather than model memory.
5. Active context should contain the smallest sufficient working set.
6. Failed attempts should carry useful evidence into subsequent attempts.
7. Work should survive context reset and model replacement.
8. Completion claims must state and support verification coverage.
9. Independent review should remain genuinely independent.
10. More agents are not inherently better.

Limit final core doctrine to roughly **five to seven durable invariants**.

A `CORE INVARIANT` must have a concrete enforcement owner. A `PROCEDURAL DEFAULT` must identify the skill that teaches it. Harness rules must not contaminate universal doctrine. Model behavior belongs in a model profile. Unproven ideas stay experiments.

---

# 26. Principle-to-control mapping

For every final core principle map:

```text
principle
    ↓
failure mode
    ↓
owning layer
    ↓
mechanism
    ↓
evidence
    ↓
regression test / benchmark
```

Possible owners:

- Bench doctrine;
- CLI;
- schema;
- hook;
- work state;
- claims/evidence;
- gate;
- skill;
- reviewer;
- harness adapter;
- model profile;
- benchmark suite.

Do not enforce one principle redundantly in five places without a specific defense-in-depth reason.

---

# 27. Action-item discipline

The final active portfolio must contain **no more than twelve items**.

Rank using:

- `P0`
- `P1`
- `P2`
- `P3`
- `EXPERIMENT`
- `DELETE`
- `DEFER`

Every action:

```yaml
id:
priority:
title:
problem:
evidence:
proposed_change:
non_goals:
owner_layer:
dependencies:
acceptance_criteria:
validation:
expected_benefit:
complexity:
risk:
migration:
```

Write:

```text
audit/fable-high/action-items.yaml
```

Every action must answer:

> How will we know this improved Bench?

Items without credible validation should normally be `EXPERIMENT`.

---

# 28. Exactly one next implementation ticket

Choose exactly one immediate implementation ticket after roadmap reconciliation.

Do not automatically choose `/bench`.

Choose the highest combined leverage across user clarity, workflow consistency, evidence quality, context efficiency, recovery, architectural simplification, and benchmarkability.

Write:

```text
audit/fable-high/next-ticket.md
```

Include:

- Title
- Problem
- Evidence
- Architectural boundary
- Scope
- Non-goals
- Acceptance criteria
- Validation
- Migration
- Dependencies
- Rollback

---

# 29. Benchmark plan

Design:

```text
A. model + ordinary repository instructions
B. current Bench
C. Bench after proposed changes
```

Use same task, repo state, model, and effort where practical. Repeat across Claude Code and Codex where meaningful. Use multiple trials where feasible.

Measure:

### Outcome
- task success;
- requirement satisfaction;
- hidden defects;
- test/CI success;
- reviewer findings.

### Epistemic quality
- unsupported assertions;
- stale evidence;
- missed validation;
- false completion.

### Context efficiency
- tokens;
- files read;
- duplicate reads;
- repeated searches;
- irrelevant reads.

### Execution behavior
- duplicate commands;
- repair-loop count;
- blank retries;
- observation-opportunity delay;
- speculation after a discriminator is available;
- failure inheritance.

### Recovery
- context-reset recovery;
- model-swap recovery;
- checkpoint completeness;
- repeated-dead-end rate.

### Orchestration
- subagent count;
- duplicate delegate discovery;
- overlapping modifications;
- integration conflicts;
- orchestrator overhead.

### Entry point
- user ambiguity;
- steps to correct workflow;
- unnecessary questions;
- incorrect routing.

### `diagnosing-bugs`
- steps to first reproduction;
- steps to first discriminating observation;
- repair-loop count;
- root-cause accuracy;
- first-fix success;
- regression validation;
- tokens/tool calls.

---

# 30. Required output files

Produce:

```text
audit/fable-high/
├── baseline.md
├── reconciliation-ledger.md
├── final-reconciliation.md
├── roadmap-inventory.md
├── roadmap-dispositions.yaml
├── proposed-roadmap.md
├── bug-triage.md
├── principle-control-matrix.md
├── action-items.yaml
└── next-ticket.md
```

If fallback invalidates the run, also produce `INVALID-RUN.md` and do not issue a final architectural verdict.

---

# 31. Final reconciliation report structure

`audit/fable-high/final-reconciliation.md` must contain:

## A. Executive verdict

- Is Bench directionally correct?
- Continue, consolidate, strangler-migrate, or rewrite?
- Strongest demonstrated mechanism?
- Largest architectural weakness?
- What happens next?

## B. Audit relationship

- where Sol and Opus agreed;
- compatible different framings;
- material disagreements;
- where one had stronger evidence;
- where both lacked evidence.

Do not declare an overall model winner.

## C. Reconciliation matrix

| Topic | Sol | Opus | Repository evidence | Classification | Final disposition |
|---|---|---|---|---|---|

## D. Highest-confidence findings

Evidence-backed only.

## E. Resolved disputes

For each:

```text
competing claims
evidence obtained
experiment run
result
final conclusion
```

## F. Unresolved hypotheses

State what remains unknown, why it matters, and cheapest resolving experiment.

## G. "Run the dang test" assessment

Where Bench already forces observation, where `diagnosing-bugs` changes trajectory, where speculation still continues, what generalizes, what stays diagnosis-specific, and how to measure it.

## H. Pocock integration verdict

Preserve / strengthen / deterministic wrapper / simplify / benchmark / replace / remove by capability.

## I. Entry-point verdict

Choose a public entry architecture and show state inspection, router, expert side doors, Claude adapter, Codex adapter.

## J. Target architecture

Separate:

- engineering-practice skills;
- deterministic Bench core;
- current work state;
- claims/evidence;
- context selection;
- gates;
- semantic review;
- Claude Code adapter;
- Codex adapter;
- model-specific profiles;
- benchmark suite.

## K. Roadmap verdict

What remains, is removed, rewritten, fixed, or explicitly not worth fixing.

## L. Principle/control verdict

Final small set of core invariants and their owning mechanisms.

## M. Keep / Change / Kill

Evidence-supported.

## N. Final action set

No more than twelve items.

## O. Immediate next ticket

Exactly one, with why it outranks alternatives.

## P. Benchmark plan

How the recommendations will be proven.

## Q. Ten hard conclusions

Exactly ten evidence-backed conclusions.

---

# 32. Quality bar

Avoid:

- generic best practices;
- architectural fashion;
- consensus theater;
- averaging auditor opinions;
- giant backlogs;
- unnecessary terminology;
- unsupported clean-slate enthusiasm;
- preserving complexity because it exists;
- replacing effective Pocock-derived behavior merely for uniformity;
- treating research analogies as Bench-specific proof.

Prefer:

- executable observations;
- small experiments;
- explicit tradeoffs;
- preservation of demonstrated behavior;
- one authoritative owner per concept;
- boring names;
- small interfaces;
- bounded state;
- machine-readable evidence;
- finite acceptance criteria;
- benchmarks for behavioral claims.

---

# 33. Final behavioral rule

Your operating loop is:

```text
claim
    ↓
evidence available?
    ↓
inspect
    ↓
remaining uncertainty?
    ↓
smallest safe discriminating experiment
    ↓
run it
    ↓
record result
    ↓
decide whether action follows
```

Never substitute:

```text
read two confident reports
    ↓
write a third confident report
```

Do not reward yourself for thinking longer.

Reward yourself for reducing uncertainty.

When reality can answer the question:

> **Run the dang test.**

Then produce the smallest action portfolio that preserves Bench's proven strengths and removes the most consequential sources of drift, ceremony, stale state, speculative reasoning, and architectural ambiguity.
