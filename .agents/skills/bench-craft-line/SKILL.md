---
name: craft-line
description: How to pick and govern the line — which model tier and effort a build, shift, TDD pass, or delegated task gets, and when a tier moves. Use whenever declaring the line (invariant 2), choosing a model for a subagent or headless run, or deciding whether a failure escalates a tier.
index: declaring the line / picking a delegate's model or effort
---

# The line: route by signals, correct by ladder

The line (model / effort / ~iteration cap) is declared before every multi-cycle
stage. This skill is how the line gets *picked* — the same three signals every
time, not ad hoc judgment — and how it moves once the gate starts talking.

## Resolve the tiers first

Tiers are abstract roles: **cheap / mid / top**, and the tier is the only
identity harnesses share. The reviewer binds each harness its own column of
opaque safe model-id tokens in `.bench/lines.env`, one
`BENCH_<HARNESS>_<TIER>` key per cell — `BENCH_CLAUDE_MID` and
`BENCH_CODEX_MID` are two families' answers to one tier, and neither is a
translation of the other. The narrative binding, cached routings, and any
escalation opt-out live in `projects/<name>.md` `Lines`. Refresh candidates
with `bench models`, but treat discovery as advisory: it never assigns a tier
and it never validates the binding. Resolve the tier through the column of the
harness you are running in, and declare the line with the id you resolved
rather than the bare tier: that id is what this harness's own surfaces accept,
and reaching into another column is a schema error, not a stricter
declaration. No `lines.env` means the repo is unrouted: declare the line from
the `Lines` prose and flag the missing binding.

## The decision table picks the starting tier

Three signals, assessed for the stage in front of you:

- **Spec precision** — does exact guidance exist (a spec with seams, a precise
  contract), or only a goal?
- **Seam uncertainty** — is the shape of the answer known, or genuinely open?
- **Gate coverage** — would a wrong answer turn the gate red cheaply, or could
  it sail through (prose, design, untested semantics)? A scaffold-stub or
  auto-detected gate is weak by definition: a repo where `/bench-setup-repo`
  has not authored the gate never counts as covered.

| Spec precision | Seam uncertainty | Gate coverage | → Line |
|---|---|---|---|
| exact | known shape | covered | cheap + low |
| partial | known shape | covered | mid + medium |
| any | genuinely uncertain | any | top + high |
| any | any | weak / uncovered | bump one tier |

Rows read top-down and the first match wins: a stage that is both genuinely
uncertain and weakly gated takes the top row — genuine uncertainty is never
demoted to a bump. "Bump one tier" means one tier above whatever the other
two signals would have selected: a weakly gated stage that would have run
cheap runs mid.

Tier and effort are one joint output — declared together, moved together. The
last row is the up-bias and it is load-bearing: where the gate can't catch a
confidently wrong answer, buy insurance at the start instead of trusting the
ladder to correct later. Under-escalation is the expensive error; a wrong
downgrade is caught cheaply by a covered gate.

A fourth signal, **leverage**, overrides the table: an artifact that steers
future generation — a skill, a command phase, shared platform rules — routes
top + high regardless of the rows above. Gate coverage on guidance prose is
structural at best, a defect is invisible and multiplies through every session
that loads the artifact, and the edit itself costs few tokens. This is the one
place where routing cheap is the expensive choice.

Check `Lines` for a cached routing before assessing from scratch — common
work types in a project get precomputed rows. Cache hits still get declared.

## Ticketed-build stage defaults

| Stage | Default line |
|---|---|
| Orchestration | mid + medium |
| Ticket implementation | cheap + low |
| Review (axis or falsification) | mid + high |

These are starting defaults under the decision table and ladder, not a flat
rule. A failed done-claim or an uncleared red moves through the same declared
ladder below. The leverage override still wins: kit guidance prose routes
top + high in every stage.

A spec's per-story line is a **ceiling, not a binding**: re-run the decision
table per ticket at charge time, because a story routed mid for its uncertain
seam often decomposes into tickets that are individually exact-spec,
known-shape, and gate-covered — the cheap row. Charging every ticket at the
story's tier is the quiet over-provisioning this paragraph exists to stop;
the re-route is declared like any line, and the ladder corrects a wrong
cheap start for one red's cost.

## Route the venue, not just the tier

The same signals pick the model and effort for a spec'd build. `craft-delegate`
owns whether code authorship runs inline or in an isolated worktree; this skill
routes the chosen author without redefining that threshold. Per-story tiers
presume separable slices; when stories land as one atomic diff, run the author at
the highest tier any story needs and flag each collapsed line in the exit report.
Use top only for the uncertain decision, then move gate-observable iteration to
the cheapest line that preserves isolation and verification.

## The ladder corrects the start

The table only has to be right-enough; gate feedback does the rest:

1. **First red** — retry the same tier, feeding the gate output back as added
   guidance. Most reds are fixable feedback, not capability gaps.
2. **Second red at the same tier** — escalate one tier.
3. **Delegate reports the seam is more uncertain than specced** — escalate
   immediately; no retry burned.
4. **Any bump to the top tier pauses and asks the reviewer** — unless the
   project's `Lines` grants a standing opt-out. Top-tier spend is the
   reviewer's cost decision, not yours.

Every move is reported in one line ("escalated to mid after 2 reds") — a
declared ladder makes moves non-silent, but each move still gets said.
The ladder is only trusted where the gate is: it corrects work the gate can
observe. It never substitutes for the up-bias row.

## What is enforced vs. declared

Model membership is enforced: the Agent-tool hook denies delegations to
models outside `lines.env`, and the `BENCH_AGENT` adapters refuse an
undeclared `BENCH_MODEL` in routed repos. Effort has no enforcement surface
anywhere — it exists only in your declaration, which is exactly why the
declaration is not optional.

## The declaration

> Line: <resolved model id> / <effort> / ~<iteration cap> / <fan-out, when the
> stage dispatches delegates>. <one clause: the signals that chose this row.>

Fan-out is declared for the same reason model and effort are — visibility
before spend, not a cost cap. A stage that plans delegates says how many
before dispatching, and exceeding the declared count is reported like any
ladder move, never done silently.

Derive the cap; don't feel it out: expected red/green cycles, plus a margin for
one red. A review-axis delegate is one pass with no fix iteration, so it prices
at ~1 iteration; a shift with a likely red and fix prices higher. Add a token
estimate only as a sizing note for a charge, not as the stop condition. The
project's cached routings hold figures like this — a stage with no basis for an
iteration count is a stage you don't understand yet, which is itself a routing
signal.

If the stage exhausts its cap, stop and report — never grind, never
escalate silently.
