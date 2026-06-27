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

3. **Write `specs/<feature>.md`** using the template below.

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
What this explicitly does not cover.
```

When the spec is written, the build has a fixed target: the user stories set the
breadth, the seams set where tests live, and the gate sets what "done" means.
Offer to run `/build` or kick off a `/shift`.
