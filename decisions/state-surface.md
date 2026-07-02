# Structured state surface

Expose the state `bench` already computes as machine-queryable subcommands, so
agents stop re-deriving it with hand-assembled greps. Wave 1: `bench learnings`,
`bench maps`, `bench guards` (the guard self-disclosure idea folds in here).
Full findings context: `ASSESSMENT.md` Part 1. The kit-prose consequences
(craft-cli scope clause, benchkit profile, BENCH.md command list, re-pointing
the phase files' grep instructions) are build work for the spec, not decisions.

## #1: Does bench adopt AXI for its own output?

Blocked by: —
Type: Grill

### Question
bench is currently exempt from AXI by a closed decision ("plain text, stderr
errors, documented exit codes"). New agent-facing query surfaces reopen it:
full adoption, keep the exemption, or split?

### Answer
Hybrid. Every existing command keeps its gate-tested plain-text contract
(`status`'s 5-row dashboard, `gate`'s exit-code verdict, the operational
commands). New query surfaces — commands that exist only for agent consumption
— are AXI-conformant: TOON stdout, minimal schemas, structured stdout errors,
definitive empty states, honest exit codes. `craft-cli`'s scope clause and
`projects/benchkit.md` change from a blanket exemption to naming this split.
Primary rationale: one parser per signal ends the two-derivations bug class
(the learnings-counter bug); token savings are secondary.

## #2: How are future codification candidates captured?

Blocked by: —
Type: Grill

### Answer
Through the existing sinks — no new journal. An agent that catches itself
assembling the same ad-hoc check repeatedly appends an `[open]` entry to
`.bench/learnings.md` naming the candidate subcommand; `/bench-integrate-learnings`
is the decide-if-codified gate. One line is added to the learnings charter
naming recurring-ad-hoc-assembly as an entry class. One-off tangents keep
going to `bench idea` / `ROADMAP.md`.

## #3: Query mode on status, or dedicated subcommands?

Blocked by: —
Type: Grill

### Answer
Dedicated subcommands (`bench learnings`, `bench maps`, `bench guards`)
sharing `status`'s parsers. `status` stays exactly the fixed ranked dashboard;
each new command carries its own small AXI contract; no `status --json`.

## #4: What is in wave 1?

Blocked by: #3
Type: Grill

### Answer
Exactly the three subcommands. `status` re-grepping `structure`'s human text is
fixed internally — the two commands share a violations function — with no new
public surface. `models` stays as-is; `.bench/lines.env` is already its
machine-readable form. Second-wave parsers (`diff`, `refs`, `coverage`,
`doctor`, `detect`) stay parked on the roadmap.

## #5: How is TOON emitted from plain shell?

Blocked by: #1
Type: Prototype

### Answer
Hand-rolled bash helper for flat tables only (`name[n]{fields}:` header + CSV
rows) — the exposed data is flat tuples, so the general format isn't needed.
No new runtime dependency. The helper is shared plumbing for all wave-1
commands and gets its own contract test (including the field-escaping edges:
commas, quotes, empty values).

## #6: What counts as a guard, and where does its manifest come from?

Blocked by: —
Type: Grill

### Answer
A guard is a boundary that can deny. Five qualify: `block-dangerous-git.sh`
(PreToolUse:Bash), `check-agent-line.sh` (PreToolUse:Agent), `stop.sh` (Stop,
armed shifts only), the adapters' `_line-guard.sh`, and the bench-managed git
pre-push hook. `session-start.sh` is informational, not a guard.

Manifests are per-script `--describe` modes — each script answers from the
same tables it enforces (the git guard prints its verb classes from the python
structures that deny; the line guard prints the live `lines.env` binding), so
the advertisement cannot drift from the enforcement. The generated pre-push
heredoc carries the same mode. `bench guards` discovers by convention — every
`.bench/hooks/` script plus the two non-hook guards — and reports a
deny-capable script lacking `--describe` as "no manifest," never skipping it
silently. Manifest fields follow `craft-cli` minimal-schema discipline
(name, boundary, denies, why); exact schema is spec detail.

## #7: What does SessionStart inject, within what budget?

Blocked by: #6
Type: Grill

### Question
The roadmap entry wants the block surface known up front instead of learned by
collision. Full `bench guards` output, a one-line-per-guard summary, or only
guards likely to fire given current state? What line budget, and does it join
the existing dashboard hook output or stand alone?

### Answer
— (open)

## #8: Does the build slice at the guards boundary?

Blocked by: #6
Type: Grill

### Question
`learnings` + `maps` + the emitter + the gate conformance layer depend on
nothing open; `guards` waits on #6–#7. Ship as one spec, or slice? Reviewer's
scoping call, recorded here so the map doesn't close around it silently.

### Answer
— (open)

## #9: What does the gate assert about the new surfaces?

Blocked by: #1
Type: Grill

### Answer
A conformance layer mirroring the existing runtime-contract + canary pattern:
for each wave-1 command, assert TOON-shaped stdout, a definitive empty state,
structured error on stdout with exit 1 (and usage errors exit 2), and exit 0
on the empty case. Canary fixtures prove each check still bites. This is the
same discipline the kit prescribes for AXI projects, pointed at itself.
