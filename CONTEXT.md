# Context — ubiquitous language for benchkit

The terms below have one canonical name each. Use these; don't invent synonyms. A
cold session reads this first to avoid drifting the vocabulary.

## Core terms

- **gate** — the executable oracle at `.bench/gate.sh` (or `$BENCH_GATE`, or an
  auto-detected default). Exit 0 means shippable. It is the *only* thing that can
  call work done. Not "the checks", not "CI", not "the suite" — the gate.
- **oracle** — the authority that decides done-ness. The gate is this repo's oracle.
  The model is never the oracle; it does not grade its own work.
- **shift** — one run of the gated loop (`bench shift "<objective>"`): iterate, run
  the gate, commit only on green. The unit of autonomous work.
- **seam** — a stable interface where a test attaches and where you compose rather
  than invent. Listed per repo in `projects/<name>.md`. Not "boundary", not
  "interface layer" — seam.
- **line** — the declared model + effort + rough token cap for a stage, with one
  clause of justification. "Declare the line" = state this before a long run. Not
  "budget" alone — the line is the whole routing decision.
- **worktree** — a warm, isolated, reusable git worktree from the pool
  (`bench worktree`). Where a shift runs without touching the main checkout.
- **invariant** — one of the four non-negotiable rules in `AGENTS.md` that override
  convenience. Not "guideline", not "best practice" — invariant.
- **harness** — the agent runtime that reads `AGENTS.md` (Claude Code, Codex,
  OpenCode, …). The kit is harness-agnostic by design.
- **kit** — benchkit itself: the shipped workflow (CLI + working agreement + skills +
  commands + hooks). Not "framework", not "tool" — kit.
- **profile** — the per-repo file at `projects/<name>.md` a cold session reads to
  learn the seams, lines, gate command, and (for UI repos) the design source.
- **skill** — probabilistic guidance that shapes *how* the model generates
  (`.claude/skills/*/SKILL.md`). Reached for when the task matches; not a rule.
- **command** — a canonical phase of the workflow (`/map`, `/spec`, `/diagnose`,
  `/build`, `/review`, `/verify`, plus `/setup`, `/resynthesize`). Not "slash command
  template" in prose — command.

## Avoid

- "CI" / "the build" when you mean **the gate**.
- "task" / "session" when you mean a **shift**.
- "boundary" / "abstraction point" when you mean a **seam**.
- "framework" / "tooling" when you mean the **kit**.
