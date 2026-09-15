# Debug structure transfer research

Recommendation: Preserve debug's local procedure and adapt its action, evidence, and stop sequence to the three build phases.
Scope: Spec writing, implementation, semantic review, and the already-requested debug authorship repair.
Evidence status: Source reads support the structural comparison. No Bench timing experiment establishes a speed gain.
Consumed by: specs/debug-loop-guidance/decisions/debug-loop-guidance.md
Drift: Refresh after the cited guidance, proposed scope, or published research changes.
Retire when: The replacement spec compiles this map and preserves its evidence.

## Question graph

1. Which current debug steps constrain the next action?
2. Which constraints exist in spec writing, implementation, and review?
3. What primary research supports transfer, and what remains unknown?

Questions 1 and 2 use local sources. Question 3 uses an independent primary-source read.
The recommendation depends on all three answers.

## Source identity

The local baseline is main commit `2f26631291785fc65728cd1377e4f96c334bc7d6`.
The shaping worktree starts at `fa917c098719380d774c1996b53e334e37148b1e`.
Its extra commits contain the existing trial spec, not implemented guidance.

The original discussion is Claude session `813a009c-11dd-443c-8cba-7be2c4cd33d8`, dated 2026-09-15.
The coordinator read its user observation, local skill comparison, and proposed trial in the preceding session turn.
Those entries contain no identified external study or measured comparison.
The current conversation supplies direct-adoption authority. The historical trial request does not override it.

## Question 1: What debug actually supplies

### Source facts

The debug command places six stages in execution order.
Phase 1 requires a command already run against the reported symptom.
It names four properties and refuses advancement without reproduction evidence.
Phase 2 checks symptom fidelity and reduces the reproduction.
Phases 3 and 4 rank falsifiable hypotheses and isolate experimental changes.
Phases 5 and 6 preserve the original reproduction through repair and close.

Sources: `.agents/commands/bench-debug.md:35`, `:49`, `:65`, `:67`, `:87`, `:94`, `:100`, `:116`.

The command also tells the agent to improve the reproduction's speed, precision, and determinism.
Its loop-constructions reference supplies practical methods when the next action is unclear.

Source: `.agents/commands/bench-debug.md:37` and `:43`.

### Inference

The distinctive structure combines a concrete first action with evidence required before the next action.
A phase heading alone does not supply that structure.
The original symptom constrains the repair target, which prevents an easier substitute from becoming the completion test.
The local evidence supports this mechanism description, not its measured contribution to user-observed efficiency.

```text
reported symptom -> reproduce -> discriminate -> repair -> rerun original symptom
                         |             |                         |
                    missing evidence   no useful change          regression check
                         |             |                         |
                         +------ stop or change approach --------+
```

## Question 2: What transfers to each phase

### Source facts

| Phase | Existing mechanism | Relevant source |
| --- | --- | --- |
| Spec writing | Source reads, seam sketch, wrong-implementation challenge, and coverage map | `.agents/skills/bench-craft-spec/SKILL.md:11`, `:17`, `:34`, `:76` |
| Seam selection | Existing-seam preference, bounded read budget, and alternative designs for genuine uncertainty | `.agents/skills/bench-craft-seams/SKILL.md:55`, `:72`, `:82` |
| Implementation | Approved ticket order, fixed design, focused checks, and wrong-spec exit | `.agents/commands/bench-implement-spec.md:36`, `:40`, `:44`, `:68` |
| TDD | One behavioral red, minimal declarations when needed, and row classifications | `.agents/skills/bench-craft-tdd/SKILL.md:47`, `:56`, `:85`, `:89` |
| Review | Frozen source, independent derivation, refutation, and repair handoff | `.agents/commands/bench-review-implementation.md:50`, `:186`; `.agents/skills/bench-craft-review/SKILL.md:16` |
| Repair | Existing repair allowance and mandatory-standard blockers without automated checks | `.agents/skills/bench-craft-line/references/bounded-repair-policy.md:5`, `:30` |

### Proposed adaptations

These are proposals for reviewer disposition, not new policy.

| Phase | Concrete next action | Evidence before advancement | Stop or return condition |
| --- | --- | --- | --- |
| Spec writing | Take one source outcome and inspect its current behavior | State an exact scenario and how the wrong behavior would be exposed | Unresolved intended behavior returns to the reviewer |
| Implementation | Take the next approved slice and identify its existing verification route | Run the row's check or use the existing TDD classification | A contradicted spec returns through the current stop-short route |
| Review | Take one candidate finding against the frozen source | Derive the requirement and attempt refutation with appropriate evidence | An unsupported concern stays uncertain; a confirmed finding reaches repair routing |

For spec writing, a planned test need not already execute against a nonexistent feature.
For implementation, setup may require a test or minimal declarations before a meaningful red exists.
For review, exact source evidence can establish a mandatory-standard violation without an automated test.
These distinctions follow the existing source rules in the table.

### Conflicts and scope hazards

- The old proposal requires a red before the first edit. Existing TDD permits minimal declarations before behavioral failure.
- The old proposal caps seam candidates. Existing seam guidance requires broad design exploration when uncertainty warrants it.
- Existing finding guidance requires a real refutation run. Repair policy still retains mandatory standards without automated checks.
- Debug forbids delegate invocation and routes repair authorship elsewhere. The current conversation explicitly requests the opposite ownership behavior.

Sources: `.agents/skills/bench-craft-tdd/SKILL.md:56`; `.agents/skills/bench-craft-seams/SKILL.md:82`.
Sources: `.agents/skills/bench-craft-review/references/finding-discipline.md:24`; `.agents/skills/bench-craft-line/references/bounded-repair-policy.md:30`.
Sources: `.agents/commands/bench-debug.md:141`, `:161`, `:167`.

The replacement must reconcile these conflicts at their existing owners.
It must not invent another repair allowance or turn review into a repair author.
The project gate remains the oracle. Phase evidence does not replace it.

Sources: `.agents/skills/bench-craft-line/references/bounded-repair-policy.md:5`; `.agents/commands/bench-review-implementation.md:186`; `CONTEXT.md:29`.

## Question 3: What external evidence supports transfer?

### Published results

[ReAct](https://arxiv.org/html/2210.03629), sections 2 and 3.3, interleaves reasoning, action, and observation.
It reports improved grounding, but also repetitive failures and reduced reasoning flexibility in some tasks.
Its design varies the frequency of reasoning with the task.
This supports task-specific adaptation rather than a rigid universal sequence.
Retrieved: 2026-09-15; studied version: v3.

[SWE-agent](https://arxiv.org/html/2405.15793), section 5.1 and Table 3, tests interface variations on software repair tasks.
Summarized search, bounded file views, and lint feedback outperform the corresponding alternatives in those experiments.
The intervention includes tool interfaces and context management. It does not isolate a prose instruction.
Retrieved: 2026-09-15.

[Agentless](https://arxiv.org/html/2407.01489), sections 3 and 5.2.3, uses explicit localization, repair, and validation stages.
Its patch-selection experiment improves with regression tests and generated reproduction tests.
Section 5.1.3 also reports an important limitation: only 94 of 213 selected reproduction tests recognized the reference repair.
An agent-generated red therefore requires scrutiny before it can support a repair claim.
Retrieved: 2026-09-15; studied version: v2.

### Inference and limits

The three sources support deliberate action selection and external evidence as useful design principles.
They do not establish that Bench's prose structure causes faster completion.
They do not test this proposed adaptation to spec writing or semantic review.
No result justifies a universal three-candidate cap, a red before every edit, or a two-second loop for every phase.

## Option comparison

| Option | Consequence | Basis |
| --- | --- | --- |
| Extract a universal loop | Reduces repeated general prose, but risks hiding debug's decisive local instructions | Debug Phase 1 and its local refusal at `.agents/commands/bench-debug.md:35` |
| Adapt each phase at its current owner | Exposes the next action and evidence without replacing phase-specific authority | The source inventory in Question 2 |
| Add candidate caps and citations only | Small change, but leaves the broader execution sequence implicit | The source inventory in Question 2 |

Recommendation: Adapt each phase at its current owner. Preserve debug's core sequence and practical reproduction guidance.
The already-requested authorship repair remains independently useful and needs its own reviewable outcome.
The reviewer decides whether both outcomes share one spec.

## Residual unknowns

- The relative contribution of debug's structure, examples, task type, and model effort remains unknown.
- Existing red-first and refutation guidance may already supply much of the proposed benefit.
- The exact amount of useful local guidance requires a fresh-session task check after implementation.
- No comparative trial is planned, so adoption cannot establish a causal speed advantage.

## Verification record

- [x] Opened with recommendation, scope, and evidence status.
- [x] Separated source facts, inferences, published results, and proposals.
- [x] Answered each question or retained its unknowns.
- [x] Compared options and retained conflicts.
- [x] Included a diagram for the evidence sequence.
- [x] Reopened all three primary papers after the delegate return.
- [x] Cited local source locations and source invalidation conditions.
- [x] Recorded that no fresh-session behavioral validation ran during research.

## Validation plan

During implementation, preserve the current red-capable debug reproduction and the whole-project gate.
Check one fresh-session task per changed phase for its first action, required evidence, stop behavior, and final handoff.
Include a new feature without an existing test and a review finding about a mandatory prose standard.
Include a bug outside a delegate's fence to check the diagnostic handoff without unauthorized edits.

Use existing review evidence to retain observations. Do not add trial flags or a new measurement schema.
The spec author chooses the engineering seams and checks after the reviewer resolves the scope.
