# Bench — handoff

Pick-up doc for continuing this work in a fresh session or harness. The kit is the
source of truth for current state — read `README.md`, then `AGENTS.md`, then this
for objectives and next steps.

## What Bench is

A local agent-development workflow that fuses **Matt Pocock's planning pipeline**
(the cognitive layer) with **kunchenguid's operational substrate** (worktrees + a
gated autonomous loop), bound by four invariants. Pocock decides *what* to do and
where to test it; the substrate is *where work runs and how it's gated*; the
invariants decide who has authority when the agent and the checks disagree. The
point is to make an autonomous build loop safe to leave running — by guaranteeing
the agent never grades its own work.

### The four invariants (the non-negotiables)

1. **The gate is the oracle — the agent never grades its own work.** "Done" = the
   gate exits zero. A Stop hook blocks finishing a shift on a red gate.
2. **Declare the line before a long run** — model + effort + token cap, justified.
3. **Document for the teammate who just walked in** — current-state ADRs; history
   in git.
4. **One small change at a time, repo stays green** — commit on green, never red;
   compose an existing seam before inventing one. The agent has no merge or history
   authority (a git `pre-push` hook + a PreToolUse guard enforce it).

### The layers

- **Enforcement** (hard authority): `.bench/gate.sh`, the git `pre-push` hook, and the
  Claude Code hooks `stop.sh` + `block-dangerous-git.sh`.
- **Generation-shaping** (probabilistic): the skills.
- **Workflow discipline** (canonical phases): the commands.
- Plus the **`bench` CLI** for the operational substrate (worktrees + gated loop).

## Harness portability (this is wired, not aspirational)

`AGENTS.md` is the single source of truth, read by every harness. Claude Code reads
it via a one-line `@AGENTS.md` import in `CLAUDE.md`; Codex/OpenCode/other AGENTS.md
harnesses read it natively. `bench link` wires a repo for all harnesses at once
(skills into `.claude/skills` *and* `.agents/skills`, commands into both, the
Claude hooks, and the git guard), so switching harnesses is just running a different
agent — no reconfiguration.

What degrades on non-Claude harnesses (degrades, doesn't break): slash-command
auto-invocation (commands become "read the phase file") and the interactive Stop /
PreToolUse hooks. What survives intact: the `bench shift` loop (gate-on-green,
commit only on green) and the git `pre-push` guard — both harness-independent.

## Current state

The kit exists as `bench.tar.gz`. Nothing is wired into a real repo yet — it's a
kit, not yet installed. Contents:

- Commands: `/setup` + `/resynthesize` (maintenance), `/map`, `/spec`, `/diagnose`, `/build`, `/review`, `/verify`
- Skills: `seams`, `tdd-at-seams`, `adr`, `axi`, `design-system`,
  `writing-great-skills`, `grill`
- Hooks: `stop.sh` (completion oracle), `block-dangerous-git.sh` (git guard)
- CLI `bin/bench.sh`: `link`, `init`, `gate`, `worktree`, `shift`
- Profiles: `projects/regroup.md`, `projects/gl-axi.md`
- Instructions: `AGENTS.md` (canonical) + `CLAUDE.md` (import shim)

Known stubs by design: `bench shift` assumes a headless agent (`claude -p`); the
auto-detect gate is only a fallback (the real path is an explicit `.bench/gate.sh`);
the gl-axi gate names `axi-conformance` + `bench-glab-delta` and the Regroup gate
names `design-conformance` as contracts you implement.

## Next steps (in order)

1. **Install.** Symlink `bin/bench.sh` onto PATH (a symlink, so it can find its kit),
   then in each repo run `bench link`, `bench init`, and `/setup` (the interactive
   configure step). Exact commands in `README.md`.
2. **Stand up the gate in each repo** — the single load-bearing step; until the gate
   runs real checks the oracle is empty and the system is inert. `/setup` walks you
   through this; or edit each `.bench/gate.sh` directly:
   - Regroup: `mypy regroup && pytest -q && ruff check regroup` (+ design-conformance for UI).
   - gl-axi: `pytest -q && axi-conformance ./gl-axi && bench-glab-delta --assert`.
3. **Implement the contracted gate checks** (referenced, not shipped):
   - gl-axi `axi-conformance` — assert TOON stdout, ≤4-field default schemas,
     structured stdout errors, correct exit codes, definitive empty states.
   - gl-axi `bench-glab-delta --assert` — wrap your existing paired-delta harness to
     exit nonzero on a token/accuracy regression vs raw `glab`.
   - Regroup `design-conformance` — lint against the **design repo**: no raw hex /
     hardcoded px, no `#FFFFFF` near the player, new components import canonical ones.
4. **Write `CONTEXT.md` for each repo.** Ubiquitous language + module map, read
   first by cold sessions. Regroup: phase taxonomy, possession frames,
   `CoordinateProvider`, the shuttle slider. gl-axi: AXI vocabulary, the output
   boundary. Then add "read CONTEXT.md first" to each profile.
5. **Run one real shift end-to-end** on something low-stakes (a gl-axi command
   wrapper): `/spec` it, `bench shift`, confirm the gate gates, the Stop hook blocks
   on red, `/review` runs, you merge.
6. **Optional / when needed:** an AXI SessionStart ambient-context hook
   (cold-session dashboard, pairs with CONTEXT.md); a `/handoff` command if you lose
   mid-task state; run `writing-great-skills` against the kit to prune.

## Design system note

The Regroup design system is a **separate repo** (tokens + canonical components),
consumed here as a submodule/package/path. Any harness reads it. Claude Design is
one optional authoring surface for that repo, not a dependency — under Codex/GPT you
edit the design repo directly. The handoff is repo-to-repo; the `design-conformance`
gate runs against the design repo regardless of harness.

## Discipline carried over

The kit is near the legibility ceiling (6 workflow commands + 2 maintenance commands, 7 skills, 2 hooks, CLI,
profiles). A new piece earns its place only if it fills a gap the existing layers
can't — not because it's good in isolation. Watch the three check surfaces for
sediment: the gate (`/verify`, deterministic, authoritative), `/review` (semantic,
advisory), and the Stop hook. Keep them single-sourced — the gate command lives only
in `.bench/gate.sh`. The urge to add a fourth check surface is the signal to prune.

Skills-to-Bench mapping: see the provenance table in `README.md`.
