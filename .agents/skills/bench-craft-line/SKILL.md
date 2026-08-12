---
name: craft-line
description: How to pick and govern the line — which model tier and effort a build, shift, TDD pass, or delegated task gets, and when a tier moves. Use whenever declaring the line (invariant 2), choosing a model for a subagent or headless run, or deciding whether a failure escalates a tier.
index: declaring the line / picking a delegate's model or effort
---

# The line: route by signals, correct by ladder

The line (model / effort / ~iteration cap) is declared before every multi-cycle
stage: the same three signals every time, not ad hoc judgment, and a ladder
that corrects it once the gate starts talking.

## Resolve the tiers first

Tiers are abstract roles: **cheap / mid / top**, the only identity harnesses
share. The reviewer binds each harness its own column of opaque safe
model-id tokens in `.bench/lines.env`, one `BENCH_<HARNESS>_<TIER>` key per
cell — `BENCH_CLAUDE_MID` and `BENCH_CODEX_MID` answer one tier, neither a
translation of the other. The narrative binding, cached routings, and any
escalation opt-out live in `projects/<name>.md` `Lines`. `bench models`
refreshes candidates but never assigns or validates a tier. Resolve the
tier through the harness you're running in and declare the id you
resolved, not the bare tier. No `lines.env` means the repo is unrouted:
declare from the `Lines` prose and flag the missing binding.

## The decision table picks the starting tier

Three signals, assessed for the stage in front of you:

- **Spec precision** — exact guidance (a spec with seams, a precise
  contract), or only a goal?
- **Seam uncertainty** — is the shape of the answer known, or genuinely open?
- **Gate coverage** — would a wrong answer turn the gate red cheaply, or
  sail through (prose, design, untested semantics)? A scaffold-stub or
  auto-detected gate is weak by definition.

| Spec precision | Seam uncertainty | Gate coverage | → Line |
|---|---|---|---|
| exact | known shape | covered | cheap + low |
| partial | known shape | covered | mid + medium |
| any | genuinely uncertain | any | top + high |
| any | any | weak / uncovered | bump one tier |

Rows read top-down, first match wins: a stage both genuinely uncertain and
weakly gated takes the top row — uncertainty is never demoted to a bump.
Tier and effort are one joint output. The last row is the up-bias and it is
load-bearing: where the gate can't catch a confidently wrong answer, buy
insurance at the start rather than trust the ladder to correct later.
Under-escalation is the expensive error; a wrong downgrade is caught
cheaply by a covered gate.
A fourth signal, **leverage**, overrides the table: an artifact that steers
future generation — a skill, a command phase, shared platform rules — routes
top + high regardless of the rows above. A defect in guidance prose
multiplies through every session that loads it. Check `Lines` for a cached
routing before assessing from scratch; cache hits still get declared.

## Ticketed-build stage defaults

| Stage | Default line |
|---|---|
| Orchestration | mid + medium |
| Ticket implementation | cheap + low |
| Review (axis or falsification) | mid + high |

These are starting defaults under the decision table and ladder, not a flat
rule. A spec's per-story line is a **ceiling, not a binding**: re-run the
decision table per ticket at charge time — a story routed mid for its
uncertain seam often decomposes into tickets that are exact-spec,
known-shape, gate-covered tickets. The leverage override still wins: kit
guidance prose routes top + high in every stage. `craft-delegate` owns
whether authorship runs inline or in an isolated worktree; this skill
routes the chosen author. When stories land as one atomic diff, run the
author at the highest tier any story needs and flag each collapsed line in
the exit report.

## Classify reds before the ladder moves

Pin the inherited baseline at stage start — which reds exist before your
diff touches the tree. Before any retry or escalation, classify every red
the gate reports: **diff-owned** (introduced or exposed by this diff),
**inherited** (already red at baseline, untouched by the diff), or
**spec-predicted** (named as expected at this point). Only diff-owned reds
count toward the ladder; the others are reported but never retried against.
1. **First diff-owned red** — retry the same tier, feeding the gate output
   back as guidance. Most reds are fixable feedback, not capability gaps.
2. **Second diff-owned red at the same tier** — escalate one tier.
3. **Delegate reports the seam is more uncertain than specced** — escalate
   immediately; no retry burned.
4. **A non-shrinking diff-owned red set across an iteration** — stop the
   stage and surface a likely seam or spec contradiction instead of buying
   a more expensive attempt: the ladder corrects wrong-tier work, not a
   wrong seam.
5. **Any bump to the top tier pauses and asks the reviewer** — unless the
   project's `Lines` grants a standing opt-out. Top-tier spend is the
   reviewer's cost decision, not yours.
Every move is reported in one line ("escalated to mid after 2 diff-owned
reds"). The ladder is only trusted where the gate is; it never substitutes
for the up-bias row.

## What is enforced vs. declared

Model membership is enforced: the Agent-tool hook denies delegations to
models outside `lines.env`, and the `BENCH_AGENT` adapters refuse an
undeclared `BENCH_MODEL` in routed repos. Effort has no enforcement surface
anywhere — it exists only in your declaration, which is why it isn't optional.

## The declaration

> Line: <resolved model id> / <effort> / ~<iteration cap> / <fan-out, when the
> stage dispatches delegates>. <one clause: the signals that chose this row.>

Fan-out is declared for the same reason model and effort are — visibility
before spend, and exceeding the declared count is reported like any ladder
move, never done silently.
Derive the cap; don't feel it out: expected red/green cycles plus a margin
for one red. A review-axis delegate is one pass with no fix iteration, so
it prices at ~1 iteration; a shift with a likely red and fix prices higher.

If the stage exhausts its cap, stop and report — never grind, never
escalate silently.
