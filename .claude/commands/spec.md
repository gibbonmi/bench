---
description: Turn the current conversation into a build spec — user stories, pre-agreed test seams, and testing decisions. Use before any build that is more than a trivial change. No interview; synthesize what we've already discussed.
---

# /spec — lock the seams before the loop runs

Synthesize the current conversation and codebase understanding into a spec the
build loop can run against. Do **not** interview me — use what you already know.
If something load-bearing is genuinely missing, ask one question, then proceed.

The point of this command is to decide the **seams and the tests before any code
is written**, so that when `/build` runs TDD it tests at a seam I chose, against a
notion of "correct" I defined — not one the agent invents mid-loop and then
over-fits to. This is the single most important step for keeping an autonomous
loop honest.

## Process

1. **Read the current state.** Explore the repo. Use the project's vocabulary.
   Respect existing ADRs in the area you're touching — read `decisions/` and any
   `projects/<name>.md`.

2. **Pick the seams.** Sketch where this feature will be tested. Prefer an
   existing seam to a new one. Use the highest seam that exercises the real
   behavior — the fewer seams, the better; one is ideal. (See the `seams` skill.)
   State the seams explicitly and check they match my expectation before writing
   the spec.

3. **Price every cut, in your time — not a human's.** Before anything is deferred,
   estimate it in *agent* time, because that's the real cost here. The instinct to
   defer is calibrated to humans who can't spare the afternoon; you can. Two rules
   follow, and they are the point of this step:
   - **No deferral under the threshold — this binds you, not me.** Anything under
     ~30 minutes of your work that introduces no new architectural decision is
     something you do *not* get to propose deferring — it's just part of this build,
     so do it now and state the estimate out loud. But scope is mine: if *I* choose
     to defer something small for my own reasons, record it in Out of scope with its
     estimate and no argument. The rule stops you pitching the dodge; it never
     overrides my call.
   - **A cut must be a separate capability, not the rest of this one.** Something is
     only legitimately out of scope if it has its own future spec — a distinct
     feature. If it's "the rest of *this* feature" (the error cases, the other half
     of the CRUD, the edge handling), it is not out of scope; it's an acceptance
     criterion. Move it into the user stories and testing decisions so the **gate**
     enforces it, rather than leaving it in a prose aside that disappears. You can't
     ticket your way out of the spec's own breadth.

4. **Write `specs/<feature>.md`** using the template below.

## Template

```markdown
# <feature>

## Problem
The problem, from the user's point of view.

## Solution
The solution, from the user's point of view.

## User stories
A long, numbered list. "As a <actor>, I want <feature>, so that <benefit>."
Exhaustive — this list defines the breadth the build must cover, which is the
guard against a loop that does the minimum and stops.

## Implementation decisions
Modules touched, interfaces modified, schema/contract changes, architectural
calls. Decisions, not file paths or snippets — those rot. Exception: a prototype
snippet that encodes a decision more precisely than prose (a state machine, a
schema, a type) may be inlined, trimmed to the decision-rich part.

## Testing decisions
- What a good test is here: exercise external behavior at the seam, not internals.
- Which seams get tested, and the prior art (similar tests already in the repo).
- The gate command this feature must pass (defaults to the project gate).

## Out of scope
Each genuine cut as one line: **what** it is — why it's a *separate capability*
(not just the rest of this feature) — **your** time estimate to build it later.
Anything you can't defend as a separate capability, or that falls under the
threshold, does not belong here — it goes into the user stories above. An empty
section is a fine and common answer; a long one is a signal you're shrinking the
target. (Borrows Pocock's `triage` OUT-OF-SCOPE convention: exclusions are
decisions on the page, not silent omissions.)
```

Every line in this section is something I can read and veto in one pass — which is
the whole purpose. A deferral with a 15-minute estimate next to it usually argues
against itself.

When the spec is written, the build has a fixed target: the user stories set the
breadth, the seams set where tests live, and the gate sets what "done" means.
Offer to run `/build` or kick off a `/shift`.