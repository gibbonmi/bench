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
Resolved from the feedback thread. Bench is presented as a **harness command workflow
plus conversation**. In Claude Code, the reviewer invokes `/bench-setup-repo`,
`/bench-shape-idea`, `/bench-write-spec`, `/bench-implement-spec`,
`/bench-review-implementation`, and `/bench-final-check`; in Codex, the same phases
use `$bench-*` adapter skills. The worker runs `bench link`, `bench init`,
`bench gate`, `bench status`, and `bench shift` as infrastructure. Documentation can
still describe the CLI, but not as the primary thing the reviewer operates.

## #2: What should the README lead with?

Blocked by: #1
Type: Grill

### Question
What first-screen onboarding shape makes Bench feel lightweight: a command-first
quick start, the current philosophy-first explanation, or a split between reviewer
and worker setup?

### Answer
The README should lead with a **reviewer-first command quick start**: show the
Bench phase commands the reviewer invokes in their harness, then explain that the
worker runs the `bench` CLI underneath. Philosophy, provenance, and CLI mechanics
belong after that first command path, because they are supporting context rather
than the primary operating surface.

## #3: How does the setup phase hide or expose `bench link` and `bench init`?

Blocked by: #1
Type: Grill

### Question
If the reviewer never runs the CLI, should install/setup docs tell them to ask the
agent to run the setup phase, with the command responsible for running or checking
`bench link` and `bench init`; or should the CLI bootstrap remain an explicit
maintainer prerequisite before slash commands work?

### Answer
Reviewer-facing setup is **ask the agent to run `/bench-setup-repo`**. That setup
phase is responsible for checking whether `bench link` and `bench init` have already
run, running or instructing the worker-facing CLI steps when needed, then walking the
reviewer through the repo-specific gate, profile, lines, and optional domain
language. `bench link` and `bench init` remain documented as maintainer/worker
mechanics, not as the first thing a reviewer must learn.

## #4: What handoff contract should every slash command satisfy?

Blocked by: #1
Type: Grill

### Question
For each command, what must the reviewer be told at entry and exit so they can stay
oriented without reading the CLI substrate?

### Answer
Every Bench command has an **entry orientation** and an **exit handoff**.

At entry, the command tells the reviewer what this phase is for, what artifact or
state it will produce, and where it sits in the workflow. At exit, it tells the
reviewer what changed, the current gate or artifact state, and the single
recommended next command. CLI details such as `bench shift`, `bench status`, and
`bench gate` appear only when they affect the next action or explain the reported
state.

## #5: Where should CLI details live?

Blocked by: #1, #2
Type: Research

### Question
Identify which README and operating-guide sections are reviewer-facing versus
worker-facing, then decide where CLI details such as `bench shift`, `bench status`,
`bench link`, and `bench init` belong.

### Answer
README top-level onboarding and command phase descriptions are
**reviewer-facing**: they describe the command workflow, the expected conversation,
and the artifact or state each phase produces. CLI details live in a clearly labeled
**worker and maintainer CLI** section, with `.bench/BENCH.md` serving as the
operational reference for exact `bench link`, `bench init`, `bench status`,
`bench shift`, and `bench gate` mechanics.
