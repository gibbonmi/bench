You are conducting an adversarial, evidence-driven architectural audit of **Bench**, an AI-assisted software-engineering harness.

This is not merely:

- a code review;
- a prose review;
- a prompt review;
- a skill review;
- a CLI review;
- a workflow review;
- or a feature brainstorm.

It is a review of the **entire system by which Bench attempts to convert model capability into reliable, efficient, understandable, verifiable, resumable software-engineering work**.

Your governing question is:

> Does Bench demonstrably improve the realization of model capability, or has it accumulated ceremony that sounds rigorous without reliably changing outcomes?

Be skeptical of everything, including the framing in this prompt.

Do not preserve an existing mechanism merely because significant effort has gone into it.

Do not remove a mechanism merely because it is complex.

Do not assume a skill is useful because its prose is compelling.

Do not assume a skill is redundant because an upstream skill has a similar name.

Do not assume a CLI command is useful because it exists.

Do not assume a warning is a control.

Do not assume a green test proves the requirement.

Do not assume another reviewer creates independent review.

Do not assume more agents produce better work.

Do not assume more context creates better active cognition.

Do not assume more reasoning creates more understanding.

And especially:

> **When the repository, CLI, compiler, test suite, runtime, search tools, git history, or a small controlled experiment can answer the question, stop theorizing and run the dang test.**

---

## 1. Declare the audit environment

Before making any substantive claim, record:

```text
AUDIT MODE:
- in-repository dogfood
- external cold observer
- paired audit

HARNESS:
- Claude Code
- Codex
- other

MODEL:
EFFORT / REASONING SETTING:
BENCH REPOSITORY PATH:
BENCH COMMIT SHA:
BENCH BRANCH:
WORKTREE CLEAN/DIRTY:
AUTO-LOADED INSTRUCTIONS:
AUTO-LOADED SKILLS:
AVAILABLE TOOLS:
OUTPUT DIRECTORY:
```

Use an output directory outside the audited repository when practical, for example:

```text
../bench-audit-results/<commit>/<harness>/
```

Do not commit, push, or publish audit changes.

Record the initial and final git status.

Use read-only inspection by default. Controlled fixtures and experiments are permitted only when they are safe, isolated, reversible, and cleaned up.

---

## 2. Known baseline facts that must not be erased

The following are inputs to the audit. They are not universal scientific conclusions, but they are materially relevant evidence.

### 2.1 Deep Matt Pocock integration already exists

Bench has already performed a deeper integration with Matt Pocock’s engineering skills. The maintainer reports that this integration made Bench’s behavior more consistent.

Therefore:

- do not evaluate the Pocock integration as a hypothetical future dependency;
- do not treat it as merely optional prior art;
- do not recommend replacing it wholesale without comparative evidence;
- do map exactly what Bench inherited, adapted, composed, constrained, or superseded;
- do determine whether Bench preserves its consistency benefits across Claude Code and Codex;
- do identify accidental divergence, duplicate prose, and missing mechanical enforcement.

The correct question is not:

> Should Bench use Matt Pocock’s skills?

The correct questions are:

> Which Pocock-derived procedures are now foundational to Bench?

> What does Bench add around them?

> Which adaptations improve consistency?

> Where has Bench damaged or duplicated them?

> Which prompt-level procedures should remain model judgment, and which should now become deterministic Bench control?

### 2.2 `diagnosing-bugs` has repeated practitioner evidence

The maintainer has used Matt Pocock’s `diagnosing-bugs` skill many times and repeatedly observed it:

- break agents out of repair loops;
- stop repeated patching of symptoms;
- force construction of an executable feedback loop;
- drive the agent toward the root cause;
- and often enable a correct first fix after diagnosis.

Treat this as **repeated practitioner evidence**.

It is not proof that the skill improves every model, repository, or bug class. Controlled benchmarking is still required.

However, the burden of proof is on any recommendation that would remove, dilute, bypass, or replace the behavior that produced those results.

In particular, preserve and study the causal sequence:

```text
construct a red-capable feedback loop
→ reproduce the user’s actual symptom
→ minimise the reproducer
→ generate falsifiable hypotheses
→ instrument discriminating boundaries
→ identify root cause
→ write the regression test at the correct seam
→ apply the fix
→ rerun the original unminimised feedback loop
→ clean up instrumentation
```

Do not reduce this to a vague instruction to “debug carefully.”

### 2.3 Bench also draws from Kunchenguid’s work

Bench has drawn from Kunchenguid’s validation and agent-interface ideas.

Audit the current relationship to relevant projects, especially:

- `no-mistakes`
- `axi`
- agent-native CLI design
- isolated validation worktrees
- deterministic gates
- structured findings
- compact machine-readable output
- model/harness fallbacks

Do not assume Bench still matches the current upstream projects. Pin the exact revisions inspected and classify the relationship.

### 2.4 A clear entry point is a first-class product requirement

The maintainer wants one clear, obvious place for a user to begin.

Treat entry-point clarity as a core requirement, not a cosmetic documentation issue.

Evaluate whether the canonical entry should be:

- `/bench`;
- the existing `what-next` capability;
- a thin `/bench` front door that invokes `what-next`;
- a CLI command;
- or another design.

Do not assume `/bench` is correct merely because it is obvious.

Push back if another design is measurably clearer or more robust.

The final recommendation must choose **one canonical normal-user entry point**, while allowing explicit expert entry points where justified.

---

## 3. Evidence labels

Classify every major finding using one of these labels:

### OBSERVED

Directly demonstrated by repository contents, command output, test execution, runtime behavior, or a controlled experiment performed during this audit.

### REPEATED PRACTITIONER EVIDENCE

A behavior repeatedly observed in real use by the maintainer, but not yet established through a controlled benchmark.

Use this label for the existing `diagnosing-bugs` experience unless you reproduce it experimentally.

### INFERRED

A strong conclusion derived from multiple observations, with the inference made explicit.

### HYPOTHESIS

Plausible, testable, and not yet established.

State the cheapest experiment that could resolve it.

### RESEARCH

Suggested by external research or prior art, but not established for Bench.

Do not allow RESEARCH to silently become OBSERVED.

Do not allow confident prose to silently become evidence.

---

## 4. Research package

Read enough of the following primary material to understand the relevant claims. Do not read every source linearly before touching the repository. Use progressive disclosure: read a source when it becomes relevant to an active audit question.

Do not accept any source as doctrine.

For each source used, record:

- what it claims;
- what evidence it presents;
- what is directly relevant to Bench;
- what does not transfer;
- what Bench evidence supports or contradicts it;
- what remains unproven.

### 4.1 Active workspace and J-space

#### Verbalizable Representations Form a Global Workspace in Language Models

Transformer Circuits:

https://transformer-circuits.pub/2026/workspace/index.html

Paper:

https://arxiv.org/abs/2607.15495

Focus on:

- limited active workspace;
- selectivity;
- broadcast representations;
- reportable versus automatic processing;
- deliberately maintained information;
- the difference between information existing in context and information being active in computation.

Do not confuse neural J-space with a textual work-state file.

#### J-CoT: Chain-of-Thought in J-Space

https://arxiv.org/abs/2607.21981

Focus on the distinction between:

- propagating all prior natural-language reasoning;
- propagating an entire dense hidden state;
- and selectively carrying the intermediate information needed by the next computation.

Do not infer that Bench requires latent reasoning.

Use this only as a lens for asking whether Bench carries too much prior prose and too little active state.

#### DeepSeek V4 × J-Space Capability Realization Report

https://github.com/Tiger3807861189/DeepSeek-V4-J-Space-Capability-Realization-Report

#### J-Space Cognition Suite

https://github.com/Tiger3807861189/J-Space-Cognition-Suite-V3.6

Treat benchmark claims cautiously.

Focus on engineering hypotheses such as:

- capability-realization loss;
- harness and tool-schema sensitivity;
- path dependence;
- selective workspace loading;
- persistent state;
- checkpointing;
- empirical escape;
- blank retries;
- failure inheritance;
- verification coverage;
- `Goal / Core / Verified / Open / Next`.

Explicitly distinguish:

```text
neural J-space research
≠
J-CoT latent reasoning
≠
a text-level agent-control protocol
```

### 4.2 Reasoning and action

#### ReAct

https://arxiv.org/abs/2210.03629

Focus on:

```text
reason
→ act
→ observe
→ update
```

Challenge workflows that allow:

```text
reason
→ reason
→ restate
→ speculate
→ patch
```

without encountering reality.

#### Reflexion

https://arxiv.org/abs/2303.11366

Focus on:

- feedback across attempts;
- episodic memory;
- learning from failed trajectories.

Then challenge whether model-generated reflective prose should be replaced by smaller structured failure state.

### 4.3 Agent-computer interfaces and repository work

#### SWE-agent

https://arxiv.org/abs/2405.15793

Focus on the Agent-Computer Interface argument.

Ask whether Bench’s CLI and harness surface are designed for agents rather than inherited from human workflows.

#### SWE-bench

https://arxiv.org/abs/2310.06770

Use as background for realistic repository-level software-engineering tasks.

#### Agentless

https://arxiv.org/abs/2407.01489

Use as an adversarial simplicity baseline.

Repeatedly ask:

> Could this capability be replaced by localization → repair → validation without losing the guarantees Bench cares about?

Do not assume agent orchestration is valuable merely because it is sophisticated.

### 4.4 Context, state, and efficiency

#### Lost in the Middle

https://arxiv.org/abs/2307.03172

Challenge the assumption:

> The information is somewhere in the prompt, therefore the model can use it reliably.

Distinguish:

- available context;
- retrieved context;
- authoritative context;
- active context;
- stale context.

#### Coconut: Training Large Language Models to Reason in a Continuous Latent Space

https://arxiv.org/abs/2412.06769

Lower priority.

Use only to reinforce the conceptual distinction between carrying all textual reasoning and carrying only useful intermediate state.

#### ToM-SWE

https://arxiv.org/abs/2510.21903

Consider persistent user intent and stateful software-engineering interaction. Do not automatically import a second user-model agent into Bench.

#### SWE-Effi

https://arxiv.org/abs/2509.09853

Use to challenge token snowballing, expensive failures, and scaffolds that improve correctness but destroy resource efficiency.

### 4.5 Whether skills actually help

#### SWE-Skills-Bench

https://arxiv.org/abs/2603.15401

This is a critical counterweight.

Treat every skill as a hypothesis about improved behavior.

Evaluate marginal value through controlled comparisons rather than prose quality.

Pay attention to:

- specialized versus generic skills;
- pass-rate changes;
- token overhead;
- version mismatch;
- guidance conflicting with repository context.

### 4.6 Matt Pocock engineering skills

Repository:

https://github.com/mattpocock/skills

Engineering overview:

https://github.com/mattpocock/skills/blob/main/skills/engineering/README.md

Inspect the exact current or pinned revisions of at least:

- Ask Matt  
  https://github.com/mattpocock/skills/blob/main/skills/engineering/ask-matt/SKILL.md
- Grill with Docs  
  https://github.com/mattpocock/skills/blob/main/skills/engineering/grill-with-docs/SKILL.md
- Domain Modeling  
  https://github.com/mattpocock/skills/blob/main/skills/engineering/domain-modeling/SKILL.md
- To Spec  
  https://github.com/mattpocock/skills/blob/main/skills/engineering/to-spec/SKILL.md
- To Tickets  
  https://github.com/mattpocock/skills/blob/main/skills/engineering/to-tickets/SKILL.md
- Implement  
  https://github.com/mattpocock/skills/blob/main/skills/engineering/implement/SKILL.md
- TDD  
  https://github.com/mattpocock/skills/blob/main/skills/engineering/tdd/SKILL.md
- Diagnosing Bugs  
  https://github.com/mattpocock/skills/blob/main/skills/engineering/diagnosing-bugs/SKILL.md
- Code Review  
  https://github.com/mattpocock/skills/blob/main/skills/engineering/code-review/SKILL.md
- Codebase Design  
  https://github.com/mattpocock/skills/blob/main/skills/engineering/codebase-design/SKILL.md
- Wayfinder  
  https://github.com/mattpocock/skills/blob/main/skills/engineering/wayfinder/SKILL.md
- Setup Matt Pocock Skills  
  https://github.com/mattpocock/skills/blob/main/skills/engineering/setup-matt-pocock-skills/SKILL.md

Pin the commit inspected.

Do not compare by name alone. Compare:

- trigger;
- phase ordering;
- invariants;
- output contract;
- composition;
- harness metadata;
- model-invocation policy;
- state assumptions;
- failure behavior;
- measurable outcome.

### 4.7 Kunchenguid prior art

#### no-mistakes

https://github.com/kunchenguid/no-mistakes

Focus on:

- disposable worktrees;
- validation before push;
- structured findings;
- safe mechanical fixes versus human judgment;
- agent-agnostic execution;
- gating;
- clear entry points.

#### AXI

https://github.com/kunchenguid/axi

Focus on:

- token-efficient output;
- minimal default schemas;
- content truncation;
- precomputed aggregates;
- definitive empty states;
- structured errors and exit codes;
- ambient context;
- content-first defaults;
- contextual disclosure;
- consistent help.

Evaluate Bench’s CLI against these principles through actual execution, not by reading doctrine alone.

---

## 5. Central audit model

Reconstruct the actual path from user intent to verified completion.

Do not start from documentation diagrams.

Derive it from implementation, configuration, skills, hooks, tests, commands, state, and actual runs.

The system may resemble:

```text
user intent
→ entry point
→ workflow routing
→ shaping
→ specification
→ tickets/work units
→ context assembly
→ orchestration
→ delegation
→ repository interaction
→ implementation
→ executable feedback
→ evidence/claims
→ gate
→ review
→ final check
→ integration/promotion
→ handoff/resume
```

But the repository is authoritative.

At every stage identify:

```text
input
→ transformation
→ state change
→ output
→ consumer
→ enforcement
→ failure path
```

Identify every point where information is:

- introduced;
- copied;
- transformed;
- summarized;
- inferred;
- dropped;
- reloaded;
- made authoritative;
- allowed to become stale;
- broadcast;
- duplicated;
- verified;
- forgotten.

---

# Audit Phases

## Phase 0: Preflight and safety

1. Record the environment defined above.
2. Record all automatically loaded instructions and skills.
3. Record git status and commit SHA.
4. Identify commands that could mutate external systems.
5. Establish safe execution boundaries:
   - read-only inspection;
   - local tests;
   - local builds;
   - isolated worktree mutation;
   - external side effects;
   - destructive operations;
   - production actions.
6. Prohibit commits, pushes, deployments, secret reads, and destructive actions.
7. Redact secrets in all captured output.
8. Define the audit output directory.
9. Create an evidence log.
10. Run the repository’s cheapest baseline checks if safe.

Do not “run the dang test” across an unsafe boundary.

Empirical escape must remain bounded by explicit safety policy.

---

## Phase 1: Cold reconstruction

During this phase, do not use Bench’s own high-level workflow conclusions as premises.

If the harness auto-loaded Bench instructions, record them, then treat them as claims to be tested.

Begin with:

- repository tree;
- executable entry points;
- CLI command registry;
- internal packages;
- skills;
- hooks;
- schemas;
- templates;
- tests;
- fixtures;
- examples;
- generated files;
- docs;
- specs;
- tickets;
- roadmap artifacts;
- state directories.

Before reading explanatory prose for a major component:

1. inspect implementation;
2. infer its contract;
3. execute safe representative behavior;
4. record the derived model;
5. then compare it to the documentation.

For each component classify:

- **REPRODUCED**: implementation and prose agree;
- **DIVERGED**: observable behavior differs;
- **UNREPRODUCIBLE**: claim could not be demonstrated;
- **UNSOURCED**: important behavior is asserted but lacks implementation or evidence.

Apply cold re-derivation especially to:

- the entry point;
- `what-next`;
- Shape;
- Build Spec;
- Implement;
- Diagnose;
- Gate;
- Final Check;
- claims/evidence;
- ticket generation;
- spec-build lifecycle;
- delegation;
- integration;
- checkpoints;
- model routing;
- context loading.

Do not assume these names exist exactly. Discover the repository’s real terms.

---

## Phase 2: Provenance and integration map

Build a provenance map for every meaningful skill, workflow, doctrine, and CLI pattern.

Classify origin as:

- Pocock-derived unchanged;
- Pocock-derived adapted;
- Kunchenguid-derived unchanged;
- Kunchenguid-derived adapted;
- Bench-native;
- common industry pattern;
- unknown.

For each Pocock-derived component identify:

```text
upstream source
upstream commit
Bench location
intentional changes
unintentional drift
added controls
removed behavior
composition changes
harness-specific metadata
observed value
remaining duplication
```

Do not assume divergence is a defect.

Do not assume upstream is better.

Do not recommend “use upstream directly” without testing whether Bench’s adaptation is responsible for the observed consistency improvement.

For `diagnosing-bugs`, perform an especially deep comparison.

Determine whether Bench preserves:

- no hypothesis before a red-capable loop;
- exact-symptom reproduction;
- loop tightening;
- deterministic or high-rate reproduction;
- minimisation;
- ranked falsifiable hypotheses;
- one-variable-at-a-time instrumentation;
- regression test at the correct seam;
- original-loop rerun after the fix;
- cleanup;
- root-cause recording.

Identify any Bench prose or orchestration that weakens these constraints.

---

## Phase 3: Entry-point audit

A clear entry point is mandatory.

Determine the actual current first-run and daily-use experience.

Test at least these scenarios:

1. New user in a repository with Bench installed.
2. Existing user with no active work.
3. User with an unshaped idea.
4. User with a shaped idea but no spec.
5. User with a spec but no bounded work units.
6. User with a ready ticket.
7. User with interrupted implementation.
8. User with failing tests or a reported bug.
9. User with unverified changes.
10. User with completed local work awaiting review/integration.
11. User who knows exactly which expert workflow they want.
12. User switching from Claude Code to Codex.

For each scenario measure:

- commands or prompts required;
- whether the correct path is discoverable;
- whether Bench asks questions it could answer from repository state;
- whether state is inspected before routing;
- whether two entry paths conflict;
- whether a user can accidentally bypass required controls;
- whether the experience differs between harnesses.

Evaluate this candidate design without assuming it is correct:

```text
/bench
→ inspect repository and Bench state
→ invoke the routing capability
→ recommend one legal next action
→ load only the needed workflow
```

Possible product model:

```text
normal user:
  /bench

expert direct entry:
  /shape
  /build-spec
  /implement
  /diagnose
  other justified direct paths

internal composition:
  specialized skills and deterministic CLI operations
```

Also evaluate whether setup must be separate:

```text
bench init
```

The final report must choose one canonical normal-user entry point and explain:

- why it is clearer than alternatives;
- how it discovers state;
- how it handles ambiguity;
- how it resumes work;
- how it behaves in Claude Code;
- how it behaves in Codex;
- what remains directly invokable by experts.

---

## Phase 4: Full skill sweep

Locate every skill. Do not restrict the sweep to known names.

For every skill answer:

### Identity

- What is its exact name?
- Is it user-invoked, model-invoked, or both?
- How is invocation policy represented in Claude Code?
- How is invocation policy represented in Codex?
- Is the skill actually discoverable in each harness?

### Purpose

- What problem does it solve?
- Can the problem be stated in one sentence?
- What failure becomes more common if the skill disappears?

### Trigger

- When should it run?
- Is the trigger observable?
- Can two skills believe they own the same transition?
- Can the agent fail to invoke it when needed?
- Can it invoke it too broadly?

### Inputs

- What repository state does it need?
- What user intent does it need?
- What context does it load?
- Which inputs are authoritative?
- Which inputs are copied prose?
- Which inputs can become stale?

### Procedure

- What sequence does it prescribe?
- Which steps require judgment?
- Which steps are mechanical?
- Where can the model skip or reorder steps?
- What happens when the model becomes uncertain?

### Output

- What artifact or state transition occurs?
- Who consumes it?
- Is the output bounded and machine-readable?
- Can another model resume from it?

### Enforcement

- What prevents bypass?
- Which requirements exist only in prose?
- Which should become CLI, schema, hook, test, or gate behavior?

### Failure and recovery

- What happens on tool failure?
- What state survives?
- Does a retry inherit the diagnosis?
- Can another model resume?
- Is there an escape hatch?

### Composition

- Which skills call or rely on it?
- Is the dependency explicit?
- Can the harness see the dependency?
- Does composition preserve the upstream invariant?
- Does composition create duplicated instructions?

### Provenance

- Pocock-derived?
- Kunchenguid-derived?
- Bench-native?
- What changed and why?

### Marginal value

- What evidence shows it improves outcomes?
- What does it cost in tokens, latency, and complexity?
- What controlled comparison would establish its value?

Classify each skill:

- essential;
- high-value;
- useful;
- unproven;
- redundant;
- harmful;
- merge candidate;
- deletion candidate;
- benchmark required.

Do not delete a skill with repeated practitioner evidence merely because a simpler description exists.

---

## Phase 5: Prose and instruction audit

Perform a repository-wide audit of active and passive prose:

- skills;
- `CLAUDE.md`;
- `AGENTS.md`;
- system or harness instructions;
- templates;
- tickets;
- specs;
- doctrine;
- comments used as agent instructions;
- examples;
- generated prompts;
- onboarding docs;
- hooks that emit instructions;
- error messages;
- CLI help.

Distinguish:

```text
active instruction
reference documentation
generated context
historical rationale
example
normative contract
```

Look for:

### Repetition

The same rule expressed in multiple places.

Recommend one authoritative source and generated or linked projections where possible.

### Contradiction

Two instructions can reasonably produce different behavior.

Demonstrate with a concrete scenario.

### Negative-instruction overload

Examples:

- never;
- always;
- do not;
- remember;
- make sure;
- be careful.

For each important prohibition ask:

> Why is this not enforced mechanically?

### Vague verbs

Flag terms such as:

- carefully;
- thoroughly;
- appropriately;
- robustly;
- properly;
- relevant;
- sufficient;
- complete;
- clean;
- safe.

A term is acceptable only when context or a contract makes it observable.

### Cargo-cult reasoning prompts

Find instructions that mainly ask the model to think harder, reflect more, or be meticulous without adding evidence, state, constraints, or actions.

### Context dumping

Find instructions that load everything because it might become relevant.

### Hidden policy

Find important behavior that exists only in prose.

### Stale prose

Find documentation describing behavior the implementation no longer has.

### Model-specific contamination

Find Claude-specific or Codex-specific hacks embedded in universal Bench doctrine.

### Output bloat

Find prose that consumes active context but does not change decisions.

For every major prose finding recommend:

- keep;
- shorten;
- centralize;
- generate;
- convert to code;
- convert to schema;
- move to reference docs;
- delete.

---

## Phase 6: Workflow and state-machine reconstruction

Derive the actual Bench state machine.

It may contain states resembling:

```text
idea
→ shaped
→ specified
→ ticketed
→ assigned
→ implementing
→ locally verified
→ integrated
→ reviewed
→ promoted
```

Do not use this as the answer.

For every state identify:

- entry event;
- required artifact;
- required evidence;
- legal next states;
- illegal transitions;
- regression path;
- repair path;
- ownership;
- persistence;
- resume behavior;
- cross-model visibility.

For every transition ask:

- Is it mechanically represented?
- Is it merely described?
- Can it be skipped?
- Can stale evidence satisfy it?
- Can a local verifier accidentally authorize global completion?
- Can the next model reconstruct why the transition occurred?

Trace representative requirements end to end:

```text
idea
→ shape
→ spec
→ ticket
→ implementation
→ test
→ evidence
→ review
→ completion
```

Try to make a requirement disappear silently.

Try to satisfy a ticket while violating the originating requirement.

Try to pass a local check while leaving whole-work coverage incomplete.

---

## Phase 7: The “run the dang test” audit

This is a first-class major audit, not a small subsection.

The governing rule is:

> **If a safe, bounded executable observation can cheaply distinguish the active hypotheses, perform the observation before continuing speculative reasoning.**

Use this rule on yourself throughout the audit.

### 7.1 Mandatory uncertainty loop

Whenever you are uncertain:

```text
state the exact uncertainty
→ list the live hypotheses
→ identify the cheapest discriminating observation
→ run it if safe
→ record the result
→ eliminate or update hypotheses
```

Do not produce several paragraphs of “probably,” “likely,” or “seems” when a command can answer the question.

### 7.2 Information-gain rule

Continued reasoning is justified only when it produces at least one of:

- a new falsifiable hypothesis;
- a new constraint;
- a new dependency;
- a new discriminating experiment;
- a synthesis that changes the next action.

If reasoning merely restates uncertainty, stop.

Act.

### 7.3 Blank retries are prohibited

A retry must inherit:

- the previous command;
- the observed failure;
- the current diagnosis;
- ruled-out hypotheses;
- what changed;
- the expected discriminating outcome.

This is invalid:

```text
attempt
→ failure
→ retry approximately the same thing
```

This is required:

```text
attempt
→ failure observation
→ diagnosis
→ changed hypothesis or variable
→ discriminating retry
```

### 7.4 Repair-loop tripwire

If two repair attempts encounter the same failure class without materially new evidence:

1. stop patching;
2. enter the `diagnosing-bugs` discipline or equivalent;
3. construct or tighten the feedback loop;
4. reproduce the exact symptom;
5. minimise;
6. form falsifiable hypotheses;
7. instrument;
8. identify root cause before another repair.

Do not let “fix-forward momentum” bypass diagnosis.

### 7.5 Original-signal requirement

After any bug fix:

- run the regression test;
- run the original unminimised feedback loop;
- verify the user’s exact symptom is gone;
- remove temporary instrumentation.

A new unit test passing is not enough if it does not exercise the original failure path.

### 7.6 Audit targets

Search for places Bench permits or encourages:

- hypothesis before reproduction;
- code reading before a red-capable signal exists;
- patching before root cause;
- retry without diagnosis;
- broad test suites before a narrow discriminator;
- local green mistaken for whole-work completion;
- continued discussion after a canary or test is available;
- reviewer speculation that could be resolved by execution;
- CLI behavior inferred from prose rather than invoked;
- dependency direction inferred rather than searched;
- stale claims trusted rather than replayed.

### 7.7 Metrics

Measure or design measurements for:

#### Observation Opportunity Delay

Number of model/tool steps between:

```text
a useful executable discriminator becoming available
```

and:

```text
the discriminator being executed
```

#### Speculation After Discriminator

Count instances where the agent continues theorizing after identifying a safe useful test, command, or probe.

#### Blank Retry Rate

```text
retries without materially new evidence
÷
total retries
```

#### Failure Inheritance Rate

```text
retries explicitly incorporating prior observed failure
÷
total retries
```

#### Red-Loop Acquisition

Time, turns, commands, and tokens required to establish a red-capable feedback loop for a bug.

#### Root-Cause Before Fix Rate

Percentage of repair attempts where the causal mechanism is demonstrated before code is changed.

#### Original-Signal Rerun Rate

Percentage of fixes that rerun the original failure signal after the patch.

#### First-Fix Success After Diagnosis

Percentage of diagnosed bugs fixed correctly on the first post-root-cause patch.

This metric is especially important given the maintainer’s repeated experience with `diagnosing-bugs`.

### 7.8 Required deliverable

Create a dedicated section titled:

> **Run the Dang Test Findings**

For each pattern include:

```text
current behavior
available observation
why observation is superior
current trigger
recommended trigger
recommended enforcement
measurement
```

Propose a concise Bench doctrine for empirical escape that is short enough to remain active in working context.

Do not create a giant new skill unless the repository evidence shows one is necessary.

---

## Phase 8: `diagnosing-bugs` preservation and generalization audit

This phase receives special treatment because it has repeated practitioner evidence.

### 8.1 Determine the causal mechanism

Do not merely say the skill “encourages discipline.”

Determine which specific constraints appear to break repair loops:

- feedback loop before theory;
- red-capable exact-symptom signal;
- fast/deterministic loop;
- minimisation;
- multiple ranked hypotheses;
- falsifiable predictions;
- one-variable probes;
- regression at the right seam;
- original-loop rerun;
- cleanup.

Identify which are load-bearing.

### 8.2 Test for dilution

Inspect whether Bench’s wrappers, orchestration, or prose:

- shorten the feedback-loop phase;
- allow a hypothesis early;
- allow a patch before reproduction;
- replace exact-symptom assertion with “does not crash”;
- accept a flaky or slow loop too early;
- lose the original repro during minimisation;
- skip the original-loop rerun;
- let a delegate return a fix without root-cause evidence.

### 8.3 Test across harnesses

Compare Claude Code and Codex on the same representative bugs:

- no skill;
- current Bench path;
- upstream `diagnosing-bugs`;
- Bench-integrated `diagnosing-bugs`.

Use repeated trials where practical.

Do not optimize the benchmark to favor the skill.

### 8.4 Generalization question

Determine whether the underlying pattern should become a general Bench invariant:

> When a bounded executable observation can distinguish the current hypotheses, execute it before further speculation.

Evaluate how this applies beyond bugs:

- design canaries;
- CLI behavior;
- integration conflicts;
- performance regressions;
- requirement interpretation;
- gate semantics;
- review disagreements;
- migration behavior.

Preserve domain-specific diagnosis depth. Do not flatten `diagnosing-bugs` into a generic slogan.

### 8.5 Removal burden

Any recommendation to remove, replace, or substantially simplify this skill must show:

- controlled outcome equivalence or improvement;
- no regression in repair-loop escape;
- no regression in root-cause identification;
- no regression in first-fix success;
- lower or equal cost;
- preserved exact-symptom verification.

Absent that evidence, recommend preservation and targeted integration improvement.

---

## Phase 9: Work state and active workspace

Evaluate whether Bench needs or already has a canonical representation of the active work.

Candidate shape:

```yaml
goal:
  ...

core:
  - ...

verified:
  - evidence-reference

open:
  - ...

next:
  ...
```

Possible additions only if justified:

```yaml
scope:
attempts:
coverage:
checkpoint:
dependencies:
feedback_loop:
```

Do not automatically adopt these fields.

Ask:

### Goal

Can a new model identify the exact objective without reconstructing conversation history?

### Core

Can it identify the small set of constraints that matter now?

### Verified

Are verified statements backed by current evidence references?

### Open

Can unresolved questions survive a context reset without becoming accidental facts?

### Next

Is the next action executable and bounded?

### Attempts

Do failed attempts preserve useful diagnosis without carrying verbose reflection?

### Coverage

Does state record what the verifier actually covered?

### Feedback loop

For bugs and risky changes, does state preserve the command that can go red and green?

Determine whether equivalent state already exists.

Do not create a second source of truth.

Recommend one of:

- adopt a new work-state layer;
- consolidate existing state;
- strengthen an existing mechanism;
- reject the idea.

---

## Phase 10: Context compiler and progressive disclosure

Treat context as a compiled product, not a pile of files.

Determine whether Bench should compile:

```text
persistent project truth
+
current work unit
+
fresh evidence
+
active dependencies
+
harness/model adapter
=
smallest sufficient active context
```

Audit what agents actually receive.

A bounded implementation context may need:

- objective;
- active requirement IDs;
- relevant invariants;
- current seam;
- ownership boundaries;
- accepted decisions;
- verified dependencies;
- unresolved questions;
- feedback command;
- required validation;
- immediate next action.

It probably does not need:

- every project document;
- every ticket;
- all historical reasoning;
- all claims;
- all roadmap items;
- unrelated architecture.

For each context item ask:

> If removed, does success decrease?

For each omitted item ask:

> Can it be retrieved progressively when needed?

Measure or estimate:

- total context loaded;
- duplicate file reads;
- repeated searches;
- irrelevant reads;
- re-reading after compaction;
- context reconstructed by delegates;
- active constraints lost in long prompts.

Audit progressive disclosure across:

```text
repository
→ subsystem
→ seam
→ file
→ symbol
→ executable feedback
```

Do not optimize for arbitrary small files.

Optimize for semantic locality, cohesion, testability, and discoverability.

---

## Phase 11: Claims, evidence, freshness, and control

Audit the claim/evidence system deeply.

Inspect:

- command surface;
- storage;
- provenance;
- labels;
- excerpts;
- hashes;
- timestamps;
- freshness;
- replay;
- dependencies;
- exit codes;
- machine-readable output;
- human-readable output;
- bypasses;
- direct file access rules;
- gate integration.

Test the doctrine:

```text
CLI establishes observation
→ agent supplies interpretation
→ claim records epistemic state
→ gate enforces required evidence
→ replay detects staleness
```

Try to break it.

Test:

- unsupported assertion;
- missing evidence;
- stale source;
- changed command output;
- deleted file;
- failed replay;
- ambiguous empty result;
- evidence without interpretation;
- interpretation masquerading as observation;
- downstream decision depending on stale evidence;
- direct reads bypassing the CLI;
- claim dependency invalidation.

For every signal ask:

> What behavior changes because this state exists?

Examples:

```text
claim = unknown
→ what is blocked?

claim = stale
→ what is replayed or prohibited?

claim = disproven
→ what dependent decisions are invalidated?

claim = verified
→ what becomes legal?
```

If nothing changes, classify it as monitoring without control.

---

## Phase 12: CLI and Agent Experience Interface audit

Treat the Bench CLI as an interface whose users are models.

Run commands. Do not infer output from implementation alone.

For every command evaluate:

- discoverability;
- no-argument behavior;
- help;
- argument clarity;
- safe defaults;
- deterministic ordering;
- output size;
- structured output;
- TOON/JSON/plain-text tradeoffs;
- empty-state clarity;
- exit codes;
- error structure;
- idempotency;
- truncation;
- `--full` or paging escape hatch;
- precomputed aggregates;
- next-step guidance;
- composability;
- ambient repository context;
- external side-effect safety.

Compare against AXI principles, but do not copy them blindly.

Measure:

- token cost of representative output;
- round trips needed to answer common questions;
- unnecessary fields;
- ambiguous emptiness;
- stack-trace noise;
- unstable ordering;
- commands that force filesystem archaeology;
- commands that duplicate reasoning the model already does well.

Produce:

| Current behavior | Agent cost | Failure mode | Proposed owner | Proposed change |
|---|---:|---|---|---|

---

## Phase 13: CLI-first and judgment-boundary audit

Apply this rule:

> If a correctness property can be codified reliably, the model should not be responsible for remembering it.

Search for prompt-level behavior that could become:

- CLI;
- schema validation;
- static analysis;
- hooks;
- deterministic tests;
- state transition checks;
- generated context;
- gate checks.

Candidates may include:

- ownership fences;
- forbidden paths;
- dependency checks;
- requirement coverage;
- evidence freshness;
- claim lookup;
- test selection hints;
- context compilation;
- checkpoint creation;
- changed-file detection;
- stale references;
- ticket format;
- transition legality;
- output bounds.

Also identify the reverse mistake:

> Software is making a semantic decision that genuinely requires model judgment.

For each boundary recommend:

```text
model judgment
CLI observation
schema
hook
gate
reviewer
human
```

Do not encode ambiguous product judgment into brittle deterministic rules.

---

## Phase 14: Gate, review, final-check, tests, and CI

Reconstruct each independently.

For:

- tests;
- local validation;
- Gate;
- Review;
- Final Check;
- CI;
- integration review;
- promotion;

identify:

```text
purpose
inputs
coverage
timing
owner
failure effect
bypass path
independence
```

Create a responsibility matrix.

Test adversarial cases that:

- pass tests but violate a requirement;
- pass Gate but fail Final Check;
- pass Final Check but fail CI;
- produce a review finding that Gate cannot detect;
- contain stale evidence;
- have partial coverage;
- claim completion based on local success.

Useful redundancy must defend against a distinct failure mode.

Bad redundancy is the same model rereading the same prose and expressing confidence twice.

For reviewers, test independence:

- What context do they receive?
- Are they anchored by implementer rationale?
- Can they cold-derive the contract from requirements and diff?
- Is a separate model actually independent if it receives the same summary?
- Are findings converted into bounded repair work?
- Can unresolved findings be bypassed?

Track:

```text
verifier
+
coverage
+
result
```

Never treat result alone as sufficient evidence.

---

## Phase 15: Ticket and spec-build audit

Inspect all ticket and spec machinery.

For tickets evaluate:

- scope;
- ownership;
- dependencies;
- assumptions;
- requirement IDs;
- acceptance criteria;
- implementation guidance;
- required tests;
- validation;
- escalation;
- output bounds;
- duplicated global truth.

A ticket should contain enough for bounded work without becoming a second copy of the entire project.

Trace ticket drift from spec.

If a lifecycle exists resembling:

```text
start
→ assign
→ checkpoint
→ integrate
→ review
→ promote
```

derive the actual legal transitions from implementation.

Test:

- partial work;
- provisional commits;
- abandoned work;
- repair tickets;
- overlapping ownership;
- integration order;
- promotion;
- blocked work;
- stale dependencies;
- inherited findings.

Determine whether whole-project acceptance exists at exactly one authoritative transition or is ambiguously duplicated.

---

## Phase 16: Delegation and orchestration

For every delegation path identify:

```text
who delegates
why
what state is transmitted
what is omitted
what scope is owned
what is forbidden
what output is required
what evidence is required
how failure returns
how integration occurs
```

Challenge unnecessary agent boundaries.

For each subagent ask:

> What independent information, parallelism, or judgment does this boundary buy?

Weak answer:

> Another model looks at it.

Strong answers may include:

- independent cold review;
- isolated ownership;
- parallel non-overlapping work;
- specialized tool access;
- separate verification;
- bounded research.

Audit broadcast state.

Shared facts that may deserve one authoritative projection:

- invariants;
- accepted decisions;
- interface contracts;
- ownership;
- verified dependencies.

Conclusions that may need independent rediscovery:

- implementation correctness;
- absence of bugs;
- requirement satisfaction;
- review verdict.

Do not broadcast confidence as fact.

Audit the orchestrator.

It should not become:

- giant-context implementer;
- passive message router;
- duplicate reviewer;
- summary bottleneck;
- agent-count maximizer.

Determine which responsibilities require global model judgment and which belong to a deterministic scheduler or CLI.

---

## Phase 17: Failure inheritance, checkpoints, and resume

Simulate loss of the current model.

Can another model continue?

A useful checkpoint may need:

- objective;
- completed work;
- modified files;
- verified evidence;
- feedback command;
- validation coverage;
- unresolved questions;
- failed attempts;
- diagnosis;
- ruled-out approaches;
- next action.

Test:

- context compaction;
- terminal restart;
- Claude outage;
- Codex takeover;
- model upgrade;
- delegate crash;
- partial implementation;
- failed validation;
- stale evidence.

Bad retry state:

```text
task
+
“previous attempt failed”
```

Useful retry state:

```text
failed command
+
exact observed failure
+
diagnosis
+
ruled-out hypotheses
+
changed variable
+
next discriminator
```

Do not store verbose self-reflection when smaller structured state is enough.

---

## Phase 18: Model and harness independence

For every behavior classify it as:

### Bench core

Stable engineering semantics independent of model and harness.

### Claude Code adapter

Skill metadata, hooks, subagent behavior, or command integration specific to Claude Code.

### Codex adapter

Equivalent Codex-specific integration.

### Model-specific profile

A temporary optimization for a particular model/version.

Do not contaminate Bench core with model-specific prompt tricks.

Test:

- skill discovery;
- implicit invocation;
- user invocation;
- nested skill composition;
- subagent context;
- tool schemas;
- output parsing;
- checkpoint resume;
- model routing.

The final architecture must show explicit boundaries among these layers.

---

## Phase 19: Dogfood pass

After cold reconstruction, use Bench as intended.

Do not merely inspect it.

Exercise representative workflows:

- enter through the normal front door;
- shape a bounded idea;
- build or inspect a spec;
- create or inspect work units;
- implement a controlled small task if safe;
- diagnose a controlled bug;
- run Gate;
- run Review;
- run Final Check;
- create or inspect a checkpoint;
- resume;
- inspect claims;
- replay evidence;
- perform a handoff.

Record friction.

If Bench’s own tooling makes the audit difficult, that is evidence.

Compare the cold-derived architecture with the dogfood experience.

Classify divergences.

---

## Phase 20: Required controlled experiments

Perform as many safe experiments as the repository supports.

At minimum attempt to design or execute:

### Entry-point routing

Can the canonical entry distinguish:

- new idea;
- ready work;
- failing work;
- interrupted work;
- unverified changes;
- completed work?

### Gate differentiation

Construct or locate a case Gate catches.

### Final-check differentiation

Construct or locate a case Final Check catches that earlier controls do not.

### Review independence

Compare reviewer findings with and without implementer narrative.

### Stale evidence

Change an observed source and test whether freshness changes.

### Claim bypass

Attempt a direct or unsupported assertion path.

### Ownership violation

Touch or simulate an unowned path and measure detection timing.

### Requirement loss

Attempt to drop a requirement between spec and implementation.

### Context recovery

Provide a checkpoint to a fresh session and measure resume quality.

### Model swap

Start with Claude Code and resume with Codex, or reverse.

### Failure inheritance

Cause a controlled failure and determine what the next attempt receives.

### Diagnosing-bugs comparison

Use representative bug tasks and compare:

```text
no specialized skill
current Bench default
upstream diagnosing-bugs
Bench-integrated diagnosing-bugs
```

Measure:

- repair attempts;
- time to red-capable loop;
- root-cause accuracy;
- blank retries;
- tokens;
- tool calls;
- first-fix success;
- original-signal rerun;
- regression quality.

Use multiple trials where practical.

Do not draw conclusions from one lucky trajectory.

---

## Phase 21: Start-over versus consolidation assessment

Explicitly answer:

> Should Bench be rebuilt from scratch, structurally refactored in place, or migrated through a new architectural spine?

Do not let greenfield elegance erase observed behavior.

Evaluate three strategies:

### A. Continue incrementally

Keep current architecture and make local improvements.

### B. Architectural consolidation / strangler migration

Create a small new spine, route workflows through it, and migrate proven behavior behind explicit contracts.

Candidate spine:

```text
bench init
→ canonical entry
→ deterministic status
→ work state
→ context compilation
→ selected practice skill
→ executable feedback
→ evidence
→ gate
→ checkpoint
```

### C. Full rewrite

Rebuild from zero.

For each compare:

- risk to demonstrated consistency;
- migration complexity;
- duplicated sources of truth;
- testability;
- harness compatibility;
- rollback;
- time to user value;
- ability to benchmark old versus new.

The default burden of proof should favor preserving demonstrated useful behavior while replacing accidental structure.

Do not recommend a rewrite because the current system is hard to explain.

Do not reject a rewrite if the current system cannot be made coherent without preserving harmful coupling.

---

## Phase 22: Capability-realization loss map

Construct:

| Stage | Potential capability loss | Current Bench mitigation | Evidence it works | Cost | Remaining gap |
|---|---|---|---|---:|---|
| entry | | | | | |
| intent | | | | | |
| shaping | | | | | |
| specification | | | | | |
| ticketing | | | | | |
| active context | | | | | |
| implementation | | | | | |
| debugging | | | | | |
| tools | | | | | |
| evidence | | | | | |
| gate | | | | | |
| review | | | | | |
| integration | | | | | |
| resume | | | | | |
| model swap | | | | | |

Identify the highest-leverage losses.

---

## Phase 23: Complexity budget and kill list

Every mechanism has costs:

- tokens;
- latency;
- code complexity;
- maintenance;
- context load;
- state synchronization;
- failure surface;
- model dependence;
- human learning;
- debugging difficulty.

For every major capability estimate:

```text
benefit
÷
complexity
```

Classify:

- keep;
- strengthen;
- simplify;
- merge;
- move to CLI;
- move out of CLI;
- move to reference docs;
- benchmark;
- delete.

Use Agentless as a hostile baseline.

Anything Bench adds beyond:

```text
understand
→ locate
→ change
→ execute feedback
→ inspect diff
→ validate requirement
```

must justify itself through:

- correctness;
- evidence;
- recovery;
- context efficiency;
- delegation;
- safety;
- reproducibility;
- model independence;
- reduced human oversight.

---

## Phase 24: Benchmark design

Design an evaluation suite comparing:

```text
A. model + repository instructions, no Bench workflow
B. current Bench
C. proposed Bench architecture
```

Add focused variants where relevant:

```text
D. current Bench without diagnosing-bugs
E. current Bench with diagnosing-bugs
F. upstream skill versus Bench-integrated skill
```

Control:

- repository commit;
- task;
- model;
- effort;
- harness;
- tool permissions;
- time budget;
- token budget;
- starting context.

Repeat trials.

Use representative Bench tasks, not toy tasks only.

Measure:

### Outcome

- task success;
- requirement satisfaction;
- hidden defects;
- review findings;
- CI result;
- first-fix success.

### Epistemic quality

- unsupported assertions;
- stale evidence;
- incorrect completion;
- missing validation;
- claim/evidence mismatch.

### Context efficiency

- tokens;
- files read;
- duplicate reads;
- irrelevant reads;
- repeated searches;
- context loaded before relevance.

### Execution efficiency

- tool calls;
- duplicate commands;
- observation opportunity delay;
- speculation after discriminator;
- blank retries;
- time to red-capable loop;
- time to root cause.

### Recovery

- checkpoint completeness;
- cross-session resume;
- cross-model resume;
- repeated dead ends;
- failure inheritance.

### Orchestration

- subagent count;
- duplicated discovery;
- overlapping changes;
- integration conflicts;
- orchestrator tokens.

### Complexity

- Bench operations required;
- latency;
- new failure modes;
- human interventions.

Predefine success criteria.

Do not claim Bench is worth using merely because it catches more things if the cost becomes unacceptable.

---

# Required Final Deliverable

Produce a single primary report plus an evidence appendix or evidence directory.

The primary report must use this structure.

## A. Executive verdict

Maximum approximately one page.

Answer:

- Is Bench directionally correct?
- What is its strongest architectural idea?
- What demonstrated behavior must be preserved?
- What is its largest weakness?
- What should be removed?
- What should be strengthened?
- Should it be rewritten, consolidated, or incrementally improved?
- What should be built next?

## B. Audit environment

Include commit, harness, model, effort, instructions, skills, tools, and limitations.

## C. Bench reconstructed

Show the actual end-to-end architecture and state transitions.

## D. Entry-point verdict

Choose exactly one canonical normal-user entry point.

Include direct expert entry points and setup behavior.

Explain the Claude Code and Codex experience.

## E. Provenance map

Map Pocock, Kunchenguid, Bench-native, and other sources.

Pin exact upstream commits.

## F. Component inventory

For every meaningful:

- skill;
- CLI subsystem;
- hook;
- workflow;
- state artifact;
- gate/check;
- agent role;

include:

```text
purpose
input
output
consumer
enforcement
provenance
observed value
cost
problems
recommendation
```

## G. Pocock integration assessment

Explain:

- what was preserved;
- what was improved;
- what drifted;
- what was duplicated;
- what should remain prompt-level;
- what should become deterministic;
- whether consistency gains survive across harnesses.

## H. `diagnosing-bugs` assessment

This must be a dedicated major section.

Include:

- causal mechanisms;
- evidence from the audit;
- practitioner evidence;
- dilution risks;
- cross-harness behavior;
- benchmark proposal;
- preserve/change/replace recommendation.

## I. Run the Dang Test Findings

Use the required format:

```text
current behavior
available observation
why observation is superior
recommended trigger
recommended enforcement
measurement
```

Propose a concise empirical-escape doctrine.

## J. Prose audit

List duplication, contradiction, vague rules, hidden policy, context bloat, model-specific contamination, and deletion candidates.

## K. Work-state assessment

Evaluate `Goal / Core / Verified / Open / Next` or a better alternative.

Recommend adopt, consolidate, strengthen, or reject.

## L. Context compiler assessment

Show the smallest sufficient context for representative work.

## M. Claim/evidence assessment

Cover provenance, freshness, replay, dependencies, control transitions, bypasses, output ergonomics, and failure behavior.

## N. Gate / Review / Final Check / CI matrix

Show responsibility, coverage, timing, and overlap.

## O. CLI / AXI assessment

Evaluate actual agent ergonomics and token efficiency.

## P. Orchestration assessment

Cover orchestrator role, delegates, ownership, broadcast state, independent review, failure inheritance, and integration.

## Q. Model/harness boundaries

Separate:

```text
Bench core
Claude Code adapter
Codex adapter
model-specific profile
```

## R. Start-over decision

Compare incremental improvement, strangler consolidation, and full rewrite.

State a clear recommendation.

## S. Complexity kill list

Be direct.

For each candidate say:

- delete;
- merge;
- simplify;
- move;
- benchmark first.

## T. Missing capabilities

Only after simplification.

Rank by leverage.

## U. Recommended architecture

Show the smallest coherent future architecture.

## V. Prioritized roadmap

Use:

```text
P0
P1
P2
P3
EXPERIMENT
DELETE
```

For every recommendation include:

- problem;
- evidence;
- proposed change;
- benefit;
- cost;
- dependency;
- proof of success.

Do not create dozens of equal priorities.

## W. Benchmark plan

Include A/B/C and focused `diagnosing-bugs` comparisons.

## X. Research reconciliation

For every research source used:

```text
claim
Bench support
Bench contradiction
transfer limit
unproven question
```

## Y. Ten hardest truths

Exactly ten.

Use the most important findings, not theatrical phrasing.

## Z. One next ticket

Choose exactly one implementation ticket.

It must have the highest leverage and a concrete proof plan.

---

# Required Adversarial Questions

Answer all of these explicitly.

1. What does Bench ask modern Claude or Codex to do that they already do reliably without help?
2. What does Bench make more consistent?
3. Which consistency gains are merely claimed, and which are observed?
4. Which Pocock-derived behaviors are load-bearing?
5. Where has Bench weakened an upstream skill through composition?
6. Where has Bench improved an upstream skill through state or enforcement?
7. Where does Bench ask a model to remember something software should enforce?
8. Where does Bench ask software to decide something that requires judgment?
9. Where does Bench confuse available context with active context?
10. Where does Bench confuse a plan with current state?
11. Where does Bench confuse assertion with observation?
12. Where does Bench confuse a passing test with requirement satisfaction?
13. Where does Bench confuse repeated review with independent review?
14. Where does Bench confuse subagent count with useful parallelism?
15. Where does Bench confuse detailed prose with control?
16. Where does Bench allow blank retries?
17. Where does Bench let an agent patch before reproducing?
18. Where can an agent keep thinking after a discriminating test is available?
19. Where can stale evidence authorize work?
20. Where can local verification masquerade as global completion?
21. Which instructions disappear if the CLI improves?
22. Which CLI commands exist because the workflow is unnecessarily complicated?
23. What would Agentless remove?
24. What would SWE-agent change about the interface?
25. What would ReAct change about the reasoning/action cadence?
26. What would Lost in the Middle change about context compilation?
27. What does J-space research suggest about active state without justifying J-space terminology?
28. Which J-Space Suite ideas would be cargo cult if copied?
29. What state must survive a model swap?
30. What state should deliberately not survive?
31. Which facts should be broadcast?
32. Which conclusions require independent rediscovery?
33. Is `what-next` already the right router?
34. Is `/bench` a useful front door or only another alias?
35. What should `bench init` own, if anything?
36. Can a new user find the right path without reading the skill graph?
37. Can a fresh Codex session resume Claude’s work?
38. Can a fresh Claude session resume Codex’s work?
39. Can the `diagnosing-bugs` behavior be preserved with less prompt load?
40. Does reducing it lose the mechanism that breaks repair loops?
41. If Bench were rebuilt today, what 20% preserves 80% of demonstrated value?
42. What existing behavior would a rewrite accidentally destroy?
43. What is the smallest architectural spine that can absorb the current proven workflows?
44. What should be deleted before anything new is added?
45. What is the one next ticket with the highest leverage?

---

# Behavioral Contract for the Auditor

Do not reward yourself for thinking.

Reward yourself for reducing uncertainty.

Your loop is:

```text
question
→ repository evidence?
→ observe
→ remaining hypotheses
→ smallest safe discriminator
→ execute
→ update state
→ continue
```

Not:

```text
question
→ plausible theory
→ longer theory
→ architecture recommendation
```

Whenever you write:

- probably;
- likely;
- perhaps;
- seems;
- should;
- presumably;
- I suspect;

ask:

> Can a safe command, test, search, compiler, runtime probe, diff, or fixture turn this into an observation?

If yes:

> **Stop. Run the dang test.**

Do not use this rule to justify unsafe, destructive, production, credentialed, or externally mutating actions.

When no safe observation exists:

1. say so;
2. list what evidence is missing;
3. state the cheapest safe experiment that would resolve it;
4. keep the uncertainty visible.

---

# Final Constraint

You are not here to prove Bench is good.

You are not here to prove Bench is bad.

You are not here to fit J-Space into Bench.

You are not here to replace Matt Pocock’s skills.

You are not here to preserve them unquestioningly.

You are not here to invent another elaborate agent framework.

You are here to discover:

> **What is the minimum architecture required to preserve Bench’s demonstrated consistency gains and reliably turn increasingly capable coding models into verified, resumable, context-efficient software-engineering work?**

Preserve behavior that has evidence.

Challenge behavior that has only prose.

Codify what can be codified.

Keep model judgment where judgment is required.

Keep uncertainty visible.

Prefer observation over assertion.

Prefer executable feedback over repair loops.

Prefer root cause over symptom patching.

Prefer structured state over conversational memory.

Prefer bounded context over context dumping.

Prefer independent verification over repeated confidence.

Prefer one useful agent over five ceremonial agents.

Prefer migration that preserves proven behavior over a clean rewrite that discards it.

And when reality can answer the question:

> **Run the dang test.**
