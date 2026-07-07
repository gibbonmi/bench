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
- **invariant** — one of the four non-negotiable rules (canonical in `.bench/BENCH.md`)
  that override convenience. Not "guideline", not "best practice" — invariant.
- **harness** — the agent runtime that reads `AGENTS.md` (Claude Code, Codex,
  OpenCode, …). The kit is harness-agnostic by design.
- **kit** — benchkit itself: the shipped workflow (CLI + working agreement + skills +
  commands + hooks). Not "framework", not "tool" — kit.
- **profile** — the per-repo file at `projects/<name>.md` a cold session reads to
  learn the seams, lines, gate command, and (for UI repos) the design source.
- **skill** — probabilistic guidance that shapes *how* the model generates
  (`.agents/skills/*/SKILL.md`; `.claude/skills/` mirrors them for Claude Code).
  Reached for when the task matches; not a rule.
- **command** — a canonical phase of the workflow (`/bench-shape-idea`, `/bench-write-spec`, `/bench-debug`,
  `/bench-implement-spec`, `/bench-review-implementation`, `/bench-final-check`, plus `/bench-setup-repo`, `/bench-update-kit`, `/bench-what-next`). Not "slash command
  template" in prose — command.
- **roadmap** — the working prioritization document at `ROADMAP.md` (repo root):
  assessed open work in priority order, ending in a `## Recommended sequence`
  that names the next actions. Printed with `bench roadmap`. A row leaves when
  the work ships or a reconcile removes it. Not "icebox", not "backlog" — roadmap.
- **ideas inbox** — the capture-and-forget sink at `IDEAS.md` (repo root):
  out-of-scope ideas parked with `bench idea`, committing to nothing.
  Append-only, no status or lifecycle; drained to zero into the **roadmap** by
  the maintenance phase.
- **park** — to capture an idea in the **ideas inbox** without committing to it
  (`bench idea "<text>"`). A *parked idea* graduates into committed work only when
  `/bench-shape-idea` pulls it into a decision map. Not "stash", not "file" — park.
- **ambient dashboard** — what `bench status` prints: the cold-session + on-demand view
  of what needs attention. The feature that renders it is the *ambient-feedback surface*;
  the printed thing is the ambient dashboard. Not "status report", not "summary".
- **signal** — one ranked line on the **ambient dashboard** (gate, uncommitted, worktree,
  learnings, structure, decision map). Shown only when it fires. Not "check", not "alert".
- **severity ladder** — the fixed rank order that decides which **signal** leads the
  dashboard and which drop under the five-row budget. Not "priority queue".
- **gate cache** — the last gate verdict (`<status> <sha> <iso8601>`) the Stop hook writes
  to the git dir, so the dashboard reads gate state without a cold run. Not "gate log".

## Avoid

- "CI" / "the build" when you mean **the gate**.
- "task" / "session" when you mean a **shift**.
- "boundary" / "abstraction point" when you mean a **seam**.
- "framework" / "tooling" when you mean the **kit**.
- "icebox" / "backlog" when you mean the **roadmap**.
