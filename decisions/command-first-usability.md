# Command-first usability — make Bench feel lightweight from the user surface

> **BOOTSTRAPPED (2026-06-30).** This map records the usability finding from the
> project feedback pass: the reviewer does not use the `bench` CLI directly. The
> user-facing product surface is slash commands plus conversation; the CLI is
> agent-facing infrastructure. The remaining questions decide how the README and
> command handoffs should reflect that.

## Grounding

- Bench's mechanics are lightweight, but the current public explanation starts from
  provenance, invariants, and CLI setup before the user sees the command workflow.
- `decisions/bench-naming.md` already decides that the reviewer does not use the
  CLI directly: the worker drives `bench` underneath.
- The active stale-reference fix is separate. Even after stale names are fixed, the
  onboarding order can still leak implementation details before the command path is
  clear.

## #1: What is the user-facing surface?

Type: Grill

### Question
Is Bench presented to the reviewer as a CLI workflow, or as a command workflow where
the agent runs CLI commands underneath?

### Answer
Resolved from the feedback thread. Bench is presented as a **slash-command workflow
plus conversation**. The reviewer asks for `/bench-setup`, `/bench-ideate`,
`/bench-spec`, `/bench-build`, `/bench-review`, and `/bench-qa`. The worker runs
`bench link`, `bench init`, `bench gate`, `bench status`, and `bench shift` as
infrastructure. Documentation can still describe the CLI, but not as the primary
thing the reviewer operates.

## #2: What should the README lead with?

Blocked by: #1
Type: Grill

### Question
What first-screen onboarding shape makes Bench feel lightweight: a command-first
quick start, the current philosophy-first explanation, or a split between reviewer
and worker setup?

### Answer
— (open)

## #3: How does `/bench-setup` hide or expose `bench link` and `bench init`?

Blocked by: #1
Type: Grill

### Question
If the reviewer never runs the CLI, should install/setup docs tell them to ask the
agent to run `/bench-setup`, with the command responsible for running or checking
`bench link` and `bench init`; or should the CLI bootstrap remain an explicit
maintainer prerequisite before slash commands work?

### Answer
— (open)

## #4: What handoff contract should every slash command satisfy?

Blocked by: #1
Type: Grill

### Question
For each command, what must the reviewer be told at entry and exit so they can stay
oriented without reading the CLI substrate?

### Answer
— (open)

## #5: Where should CLI details live?

Blocked by: #1, #2
Type: Research

### Question
Identify which README and operating-guide sections are reviewer-facing versus
worker-facing, then decide where CLI details such as `bench shift`, `bench status`,
`bench link`, and `bench init` belong.

### Answer
— (open)
