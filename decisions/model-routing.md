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
Hybrid. Three signals — spec precision, seam uncertainty, gate coverage —
key the table, and the table only picks the *starting* tier; the escalation
ladder (#5) corrects after gate feedback, so the table needs to be
right-enough, not right. Tier + effort are a joint output per row (invariant
#2 declares them together). Weak gate coverage bumps one tier — the up-bias
that guards against confidently-wrong work sailing through cheap. All
signals are proxies for likelihood-wrong × cost-of-wrong; per-project
`Lines` sections may cache common routings as precomputed rows, but the
taxonomy stays out of the rubric. Slots into the `/bench-implement-spec`
workflow, where the line is already declared before the build.

## #3: Where does the rubric live?

Blocked by: #1
Type: Grill

### Question
Kit skill (generic guidance, e.g. `craft-line`), per-project `Lines` section,
or both?

### Answer
Three homes, each doing what it already does. A new kit skill (working name
`craft-line`) holds the generic decision table and escalation ladder.
`projects/<name>.md` `Lines` keeps the tier→model bindings, cached common
routings, and the top-tier escalation opt-out. `/bench-implement-spec`'s
existing "declare the line" step references the skill instead of restating
it — a one-line change to the command, not a new phase.

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
Enforcement is feasible on both surfaces; effort is enforceable nowhere and
stays declaration discipline.

- **Claude Code (Agent tool):** a PreToolUse hook with matcher `"Agent"`
  fires on every in-session delegation — including headless `claude -p` and
  inside subagent sessions. Its stdin carries the prompt and
  `resolvedModel` (the actual model that will run; effort is not exposed).
  The hook can deny with a reason fed back to the model, ask the reviewer,
  or allow. The robust check is deterministic: resolvedModel must match a
  tier binding in `projects/<name>.md`. (`transcript_path` is available, so
  grepping for a declared line is possible but brittle — not the primary
  check.) Recommended mode: deny with reason, so the model self-corrects by
  re-delegating on a bound tier without interrupting the reviewer.
- **Codex:** hooks support PreToolUse but only a `Bash` matcher, and Codex
  has no Agent-style tool — its delegation flows through shell, i.e.
  through the kit-owned adapter scripts. Enforcement belongs in the
  adapter, not a fragile command-regex hook.
- **`BENCH_AGENT` adapters (harness-independent backstop):** the three
  reference adapters are kit-owned one-liners (`claude -p`, `codex exec`,
  `opencode run`); the contract today is prompt-as-single-arg +
  `BENCH_SHIFT=1`, no model selection at all. Extend the contract with
  `BENCH_MODEL`/`BENCH_EFFORT` env vars that each adapter maps to its
  harness flag (e.g. `claude -p --model …`); the adapter refuses or
  warns-and-logs when they're unset. Works identically on every harness.

## #5: What triggers a tier move mid-task (escalation ladder)?

Blocked by: #2
Type: Grill

### Question
When a delegate fails the gate N times or self-reports uncertainty, does the
orchestrator retry same-tier with more context, or escalate? Invariant #2
bans silent escalation — is a pre-declared ladder standing approval, or does
each move need the reviewer?

### Answer
First red retries same-tier with the gate output fed back as guidance (most
reds are fixable feedback, not capability gaps). Second red at the same tier
escalates one tier. A delegate self-reporting that the seam is more
uncertain than specced escalates immediately, no retry burned. Declaring the
ladder as part of the line makes moves non-silent, and each actual move is
still reported in one line for later audit. **Exception: any bump to the top
tier always requires reviewer approval** — pause-and-ask by default, with a
standing opt-out grantable per project in the `Lines` section. The ladder is
only trusted where the gate is: weak gate coverage routes higher at the
start (per #2) instead of relying on correction, because the gate misses
untested semantics, wrong-seam design, and ungated prose work —
`/bench-review-implementation` is the semantic net behind it.

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
Binding as of 2026-07-01: **top = Fable 5, mid = Opus 4.8 (session model),
cheap = Sonnet 4.6**. Haiku 4.5 leaves the line rotation. Sonnet 5 stays
excluded: Opus-level benchmarks with a higher token burn rate makes it a
poor cheap tier, and its orchestration strength makes it a future *mid*
candidate, not a workhorse. Foreseen next rotation, when the frontier
shifts: Fable-class top, Sonnet 5-class mid/session, Haiku-class cheap —
which is a `Lines` edit plus `bench models` refresh, not a rubric change;
the table routes roles, not models. Three tiers, not four: Haiku 4.5 cannot
take an `effort` value (breaks the tier+effort joint output) and a fourth
tier adds a routing distinction the signals can't reliably make. Sonnet 5
stays out of the mid seat for now — its nominal 40% discount vs Opus 4.8
shrinks to ~25% after its new tokenizer (~30% more tokens for the same
text) and shrinks further when intro pricing ends; the mid seat buys routing
judgment, not raw throughput. **Revisit Sonnet 5 for the mid seat on
2026-09-01 (intro-pricing end) or at the next frontier shift, whichever
comes first.** Binding caveat: the Agent tool's bare `sonnet` alias resolves
to Sonnet 5, so cheap-tier delegates target `claude-sonnet-4-6` explicitly.

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
