# Model + effort routing for delegated work

Formalize invariant #2's "declare the line" into a routing rubric: a tier-2
session decides, per task, which tier (cheap/mid/top) and effort a delegate
gets — escalating up for genuinely uncertain seams, down for exact-guidance
plumbing. Builds on the existing `Lines` section in `projects/<name>.md`
(tier→model binding + coarse per-work-type routing).

## #1: What form does the tier selector take?

Blocked by: —
Type: Grill

### Question
Numeric scoring rubric, a decision table keyed on a few dominant signals, or
prose heuristics? This shapes everything downstream.

### Answer
Decision table, keyed on a few dominant signals. Numeric rubrics are
unfaithful — an LLM's score becomes post-hoc justification for a choice
already made — while a small table is checkable in review and cheap to apply
on every delegation. Signals and outputs are ticket #2.

## #2: Which signals drive routing, and how does effort couple to tier?

Blocked by: #1
Type: Grill

### Question
Candidate signals: spec precision (does exact guidance exist?), seam
uncertainty (is the answer's shape known?), verification cost (does the gate
catch failure cheaply?), blast radius/reversibility. Which make the cut, and
is effort a coupled output (cheap+low / top+high) or an independent knob at
each tier?

### Answer
— (open)

## #3: Where does the rubric live?

Blocked by: #1
Type: Grill

### Question
Kit skill (generic guidance, e.g. `craft-line`), per-project `Lines` section,
or both? Recommendation: both, following the kit's existing split — generic
rubric as a skill, tier bindings and per-project overrides stay in
`projects/<name>.md`.

### Answer
— (open)

## #4: What can hooks actually enforce, per harness?

Blocked by: —
Type: Research

### Question
The rubric is guidance; the enforceable half is "no undeclared delegation."
There are two delegation surfaces: the in-session Agent tool, and headless
`bench shift` runs through the `.bench/adapters/` `BENCH_AGENT` adapters
(where the model is a flag, so enforcement is trivially scriptable). For the
in-session surface: can a Claude Code PreToolUse hook read an Agent call's
model/effort params and block or warn when no line was declared? Does Codex
have an equivalent? What's the harness-independent backstop, if any? Output:
a short asset mapping the enforcement surface.

### Answer
— (open)

## #5: What triggers a tier move mid-task (escalation ladder)?

Blocked by: #2
Type: Grill

### Question
When a delegate fails the gate N times or self-reports uncertainty, does the
orchestrator retry same-tier with more context, or escalate? Invariant #2
bans silent escalation — is a pre-declared ladder standing approval, or does
each move need the reviewer?

### Answer
— (open)

## #6: How do tiers rebind for the three-tier frontier paradigm?

Blocked by: —
Type: Grill

### Question
Current binding (benchkit, 2026-07-01): cheap=Haiku 4.5, mid=Sonnet 4.6,
top=Opus 4.8, Sonnet 5 excluded by reviewer directive. The new paradigm is
Fable/Opus/Sonnet (OpenAI: Sol/Terra/Luna). Does Fable become top and shift
everything down? Which Sonnet is tier-3, given the Sonnet 5 exclusion? Where
does Haiku land — fourth tier or out?

### Answer
— (open)

## #7: Is a tier-2 orchestrator the right session model?

Blocked by: —
Type: Grill

### Answer
Yes. Routing is judgment plus coordination, which tier-2 handles, and the
always-loaded orchestrator context is where cost concentrates — it should not
run on tier-1. The known weakness — a mid model is worse at knowing what it
doesn't know — is compensated in the rubric, not the model choice: the errors
are asymmetric. Under-escalation produces confidently wrong work on uncertain
seams (expensive to detect); over-delegation down is caught cheaply by the
gate. So the rubric biases up on uncertainty and downgrades only on
cheap-to-verify signals (exact guidance exists, seam known, gate covers it).

## #8: Is the mechanism hook-based?

Blocked by: —
Type: Grill

### Answer
No — split it. Hooks are deterministic scripts: they cannot score a task, only
check that something happened. The routing decision is model-applied guidance
(skill + `Lines` section, per #3); hooks are the optional enforcement layer
that verifies a line was declared before a delegation (feasibility in #4).
This matches the kit's architecture: skills shape generation, hooks enforce.

## #9: Does this port across providers?

Blocked by: —
Type: Grill

### Answer
Already solved — nothing new to build. Tiers stay abstract (cheap/mid/top);
`projects/<name>.md` binds them per harness and `bench models` refreshes the
binding. OpenAI's Luna/Terra/Sol binds the same way.
