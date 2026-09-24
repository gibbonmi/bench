---
name: craft-line
description: How to pick and govern the line — which model tier and effort a build, shift, TDD pass, or delegated task gets, and when a tier moves. Use whenever you declare the line (invariant 2), choose a model for a subagent, or escalate a tier.
index: declaring the line / picking a delegate's model or effort
---

# The line: route by signals, correct by ladder

Declare the line before every multi-cycle stage. Name the model, effort, and iteration policy. Select the starting line with the three signals below, and correct it with the gate-feedback ladder.

## Resolve the tiers first

Tiers are abstract roles: **cheap / mid / top**, the only identity harnesses share. The reviewer binds each harness to opaque model IDs in `.bench/lines.env`. One `BENCH_<HARNESS>_<TIER>` key names each cell. Cached routes live in `projects/<name>.md` `Lines`.

`bench models` refreshes candidates but never assigns a tier. Resolve the tier through the harness. No `lines.env` means the repo is unrouted; use the `Lines` prose and flag the missing binding.

When `capture/agent-performance/` has a scorecard for the provider, read its `Current routing` and `Current decisions` before you select the line.
Use the scorecard to suggest or validate the line. It does not indicate outcomes for the current task. User direction, project `Lines`, and current-task signals remain authoritative.

Before a planned fork, compare the coordinator's model with the recorded line. If they differ, use a fresh writer on the recorded model.
Flag the venue change for reviewer veto before dispatch. This resolution takes precedence over default fork inheritance.
Apply [craft-delegate](../bench-craft-delegate/SKILL.md) to any transfer of an existing author.

Raise each judgment call that the gate cannot settle as a read-only consultation question when it appears.
Apply [craft-delegate](../bench-craft-delegate/SKILL.md) to that consultation.

## The decision table picks the starting tier

Assess three signals for the stage in front of you:

- **Spec precision** — exact guidance, or only a goal?
- **Seam uncertainty** — is the shape of the answer known, or genuinely open?
- **Gate coverage** — does a wrong answer turn the gate red cheaply? A scaffold stub or an auto-detected gate is weak.

| Spec precision | Seam uncertainty | Gate coverage | → Line |
|---|---|---|---|
| exact | known shape | covered | cheap + low |
| partial | known shape | covered | mid + medium |
| any | genuinely uncertain | any | mid + high; offer the top tier to the reviewer |
| any | any | weak / uncovered | bump one tier |

Rows read top-down, first match wins. A stage that is uncertain and weakly gated takes the uncertain row.
Tier and effort are one joint output. The last row is the load-bearing up-bias. Where the gate cannot catch a wrong answer, buy insurance at the start. Under-escalation is expensive; a covered gate catches a wrong downgrade cheaply.

A fourth signal, **leverage**, overrides the table. An artifact that steers future generation routes mid + high. A guidance defect multiplies through every session. The top tier implements nothing unless the reviewer names it. Read `Lines` for a cached route before you assess from scratch.

## Ticketed-build stage defaults

| Stage | Default line |
|---|---|
| Orchestration | mid + medium |
| Ticket implementation | cheap + low |
| Review axis | conditional + high |

The conditional review line uses high effort and ~1 iteration. A Codex mid implementation sends each axis to the Codex top binding. Every other implementation sends each axis to the invoking harness's mid binding. A different implementation model, or a session beyond the planned ticket authors, requires user direction. The ticket author can adjust effort in its own session and reports the change.

These are starting defaults, not a flat rule. Every ticket author runs on the spec's declared `Line:`.

A post-review repair runs at low effort. The author raises that effort only when the repair is at risk of failure at low effort, and reports the raise in one line. The repair session keeps the ticket's declared model.

The leverage override still wins for orchestration and implementation. The conditional review line owns review.
`craft-delegate` owns the author venue; this skill routes the author.

## Delegated author lines

`--delegate` authorizes selection and eligible escalation through every configured tier, including top.
An unavailable or unbound model stops that dispatch without a substitution.
The declared author limit counts concurrently active ticket writers, not retained idle author sessions.
Every delegated review axis uses the invoking harness's configured mid binding at high effort.
That route applies to cheap, mid, and top authors alike. A no-flag run keeps the conditional review line above.

## Classify reds before the ladder moves

Pin the inherited baseline at stage start. Classify every red before retry or escalation.
A red is **diff-owned**, **inherited**, or **spec-predicted**. Only diff-owned
reds count toward the ladder. Report the others without retrying against them.
1. **First diff-owned red** — retry the same tier, feeding the gate output
   back as guidance. Most reds are fixable feedback, not capability gaps.
2. **Second diff-owned red at the same tier** — without `--delegate`, raise the effort one step and resume.
   Outside `--delegate`, a tier move of a fresh ticket author asks the reviewer first. A `--delegate` run moves the tier without that stop.
3. **Delegate reports the seam is more uncertain than specced** —
   escalate immediately under the step 2 tier-move rule; no retry burned.
4. **A non-shrinking diff-owned red set across an iteration** — stop the
   stage and surface a likely seam or spec contradiction instead of buying
   a more expensive attempt: the ladder corrects wrong-tier work, not a
   wrong seam.
5. **Outside `--delegate`, a bump to the top tier pauses and asks the reviewer** —
   unless the project's `Lines` grants a standing opt-out. Top-tier spend is
   the reviewer's cost decision, not yours.
Report every move in one line ("escalated to mid after 2 diff-owned reds").
Trust the ladder only where you trust the gate; it never replaces the up-bias row.

Known-flaky retry stops are in `craft-delegate`'s delegation discipline.

## What is enforced vs. declared

The Agent-tool hook and `BENCH_AGENT` adapters enforce model membership. Effort has no enforcement surface, so the declaration must name it.

## The declaration

> Line: <model id> / <effort> / <iteration policy> / <fan-out when used>.
> Expected repair rounds: <count> / confidence <0-10>. `references/calibration-score.md` owns the score.
> <one clause: the signals that selected this row.>

The iteration policy is a numeric cap or an explicit `uncapped` policy. Declare fan-out for visibility before spend. Report an overrun like a ladder move. Derive a numeric cap from expected cycles plus one red. Price a likely shift repair higher.

## Retained implementation continuation

At initial review or before post-review repairs, including resumed work, read [the bounded repair policy](references/bounded-repair-policy.md). It governs post-review repairs in every implementation mode.
The continuation rules below govern the pre-review work of each ticket author; the policy governs their use after initial review.
A verified acceptance improvement is progress. New useful evidence is progress when it changes the next action.
Continue while progress holds inside the approved scope.

An attempt contains one coherent hypothesis, implementation change, and verification.
An expected TDD red is not an attempt. An individual tool call is not an attempt.
A diagnostic-only action is not a completed implementation-and-verification attempt.

After two completed attempts with no progress, reassess before the next implementation attempt.
State the changed hypothesis and the next discriminating check. After reassessment, the ticket author can invoke `$bench-debug`.

Stop when the run exhausts a selected numeric cap.
Stop dependent implementation for a required user decision.
Stop for an external blocker when no independent work remains.
Stop when the run reaches an explicit user budget.

Stop immediately when the user cancels. If no useful next check remains, report the unresolved blocker and the smallest decision or evidence needed.

A diagnostic route does not change the ticket author's session. A diagnostic route does not change the ticket author's model.

An uncapped implementation has no artificial iteration stop within the approved spec. It retains the user budget, approval boundaries, external blockers, and cancellation stops.
