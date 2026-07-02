---
name: craft-line
description: How to pick and govern the line — which model tier and effort a build, shift, TDD pass, or delegated task gets, and when a tier moves. Use whenever declaring the line (invariant #2), choosing a model for a subagent or headless run, or deciding whether a failure escalates a tier. Reach for this before any multi-cycle stage.
---

# The line: route by signals, correct by ladder

The line (model / effort / ~token cap) is declared before every multi-cycle
stage. This skill is how the line gets *picked* — the same three signals every
time, not ad hoc judgment — and how it moves once the gate starts talking.

## Resolve the tiers first

Tiers are abstract roles: **cheap / mid / top**. They bind to model ids in
`.bench/lines.env` (`BENCH_TIER_TOP`, `BENCH_TIER_MID`, `BENCH_TIER_CHEAP`) —
the machine-readable source enforcement also reads. The narrative binding,
cached routings, and any escalation opt-out live in `projects/<name>.md`
`Lines`. Refresh candidates with `bench models`. Declare the line with the
resolved model id, never the bare tier. No `lines.env` means the repo is
unrouted: declare the line from the `Lines` prose and flag the missing binding.

## The decision table picks the starting tier

Three signals, assessed for the stage in front of you:

- **Spec precision** — does exact guidance exist (a spec with seams, a precise
  contract), or only a goal?
- **Seam uncertainty** — is the shape of the answer known, or genuinely open?
- **Gate coverage** — would a wrong answer turn the gate red cheaply, or could
  it sail through (prose, design, untested semantics)?

| Spec precision | Seam uncertainty | Gate coverage | → Line |
|---|---|---|---|
| exact | known shape | covered | cheap + low |
| partial | known shape | covered | mid + medium |
| any | genuinely uncertain | any | top + high |
| any | any | weak / uncovered | bump one tier |

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

> Line: <resolved model id> / <effort> / ~<token cap>. <one clause: the
> signals that chose this row.>

If the stage blows past its cap, stop and report — never grind, never
escalate silently.
