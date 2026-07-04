# Bench reference

Lookup material split out of `.bench/BENCH.md` to keep the always-loaded operating
guide lean — the file map, the skills index, harness-invocation details, the CLI
command notes, the shift adapter contract, and the hook layers. This file is
**referenced by path, not imported**: it costs no tokens until you open it. Read
it when you need a file's role, how a harness invokes a phase, or how git safety
is layered. The generation-steering rules stay in `.bench/BENCH.md`;
what lives here is reference you consult on demand.

## Files

- `AGENTS.md` contains the project-owned working agreement plus a small
  Bench-managed block.
- `.bench/gate.sh` is the project gate.
- `.bench/learnings.md` is the usage journal for process learnings.
- `.bench/bin/` is the local CLI copy `bench link` installs for hooks, so Stop and
  SessionStart do not depend on a global `bench` on PATH.
- `.agents/commands/` contains portable Bench command phases.
- `.agents/skills/` contains portable Bench skills.
- `.bench/adapters/` contains reference harness adapters for `bench shift`
  (`claude`, `codex`, `opencode`) — point `BENCH_AGENT` at one.
- `.bench/hooks/` contains shared hook scripts used by harness adapters.
- `.bench/lib/` contains shared shell functions the hooks and adapters source
  (`resolve-bench.sh`, the one source of the bench-wrapper search order, lives
  here; the tier-binding parse itself lives in the Go core, `bench resolve-model`
  / `bench check-agent-line`).
- `.claude/` contains Claude Code adapter config. See `.claude/README.md`: Claude
  reads `.claude/skills/` and `.claude/commands/`, and those paths point at the
  portable `.agents/` files. `.claude/skills/` carries only the `bench-craft-*`
  skills — the `$bench-*` phase adapter skills are Codex-only, because Claude
  already has each phase as a command and a same-named skill duplicates the
  slash-menu entry.
- `.codex/` contains Codex adapter config.

## Skills index

Claude Code loads these on its own. On other harnesses, read the file when the
trigger applies — or paste it as context. This block is generated from each
skill's `index:` frontmatter (`.bench/skills-index.sh --write`); edit the skill,
not the list:

<!-- bench:skills-index:start -->
- recording a decision or writing docs → `.agents/skills/bench-craft-adr/SKILL.md`
- building an agent-facing CLI → `.agents/skills/bench-craft-cli/SKILL.md`
- spawning a delegate / verifying a delegate's done-claim → `.agents/skills/bench-craft-delegate/SKILL.md`
- any UI work → `.agents/skills/bench-craft-design-system/SKILL.md` + your project's design source
- adding, weakening, or removing a gate check / authoring the oracle → `.agents/skills/bench-craft-gate/SKILL.md`
- surfacing a decision one question at a time → `.agents/skills/bench-craft-grill/SKILL.md`
- declaring the line / picking a delegate's model or effort → `.agents/skills/bench-craft-line/SKILL.md`
- reviewing a diff / what a finding must cite → `.agents/skills/bench-craft-review/SKILL.md`
- placing a test / designing an interface → `.agents/skills/bench-craft-seams/SKILL.md`
- writing or pruning a skill → `.agents/skills/bench-craft-skills/SKILL.md`
- evaluating a change to the kit itself → `.agents/skills/bench-craft-synthesis/SKILL.md`
- writing tests first → `.agents/skills/bench-craft-tdd/SKILL.md`
<!-- bench:skills-index:end -->

## Harness Invocation

The canonical phase bodies live in `.agents/commands/`. Harnesses may expose those
phases differently:

- **Claude Code:** invoke the phase directly as a slash command, e.g. `/bench-write-spec`.
- **Codex:** invoke the matching explicit skill, e.g. `$bench-write-spec`; each `$bench-*`
  adapter reads the canonical command file and follows it. These adapters are
  explicit-only (`allow_implicit_invocation: false`) because workflow phases are
  reviewer-chosen entry points, not background generation guidance.
  Model-invoked Bench guidance uses visible `craft-*` skill names, leaving `$bench`
  for phases a reviewer deliberately runs.
- **Other AGENTS.md harnesses:** read the phase file in `.agents/commands/` and
  follow it when no native command or skill surface exists.

The rule for translating a recommended phase into *this* harness's invocation
form lives with the communication rules in `.bench/BENCH.md`.

Codex phase adapters installed by Bench:

- `$bench-setup-repo` → `.agents/commands/bench-setup-repo.md`
- `$bench-shape-idea` → `.agents/commands/bench-shape-idea.md`
- `$bench-write-spec` → `.agents/commands/bench-write-spec.md`
- `$bench-debug` → `.agents/commands/bench-debug.md`
- `$bench-implement-spec` → `.agents/commands/bench-implement-spec.md`
- `$bench-review-implementation` → `.agents/commands/bench-review-implementation.md`
- `$bench-final-check` → `.agents/commands/bench-final-check.md`
- `$bench-update-kit` → `.agents/commands/bench-update-kit.md`
- `$bench-integrate-learnings` → `.agents/commands/bench-integrate-learnings.md`

## Command Notes

The canonical CLI inventory lives in `.bench/BENCH.md`, not in `HANDOFF.md`.
Detailed output contracts for the AXI query surfaces live in `projects/benchkit.md`;
hook and adapter plumbing is described in the sections below.

## Harness adapter for the shift loop

`bench shift` drives whatever harness `BENCH_AGENT` names: each iteration it runs
the adapter executable with the generated prompt as its **single positional
argument** and `BENCH_SHIFT=1` armed. There is no default — an unset `BENCH_AGENT`
fails fast before the loop with a configure-your-adapter error. Reference adapters
ship in `.bench/adapters/` (`claude`, `codex`, `opencode`); point `BENCH_AGENT` at
one, or at your own wrapper that maps `$1` to your harness's noninteractive
command. Use an absolute path or an on-`PATH` name; harness flags belong inside
the wrapper — a multi-word `BENCH_AGENT` value is treated as one executable name
and rejected.

The adapters also carry the line (see the `craft-line` skill): `BENCH_MODEL`,
when set, is passed to the harness's model flag. A repo with `.bench/lines.env`
(the tier→model binding) is **routed**: there the reference adapters refuse to
run when `BENCH_MODEL` is unset or names a model outside the binding, so a
headless shift always carries an explicit, bound line. Without `lines.env` the
adapters behave as plain pass-throughs. Effort has no harness flag and stays in
the declared line.

## Hook Layers

Git safety is layered:

- The git `pre-push` hook blocks direct pushes to the default branch.
- Claude Code and Codex hook adapters call the shared scripts in `.bench/hooks/`.
  Codex loads `.codex/hooks.json` only after you trust it once via `/hooks`
  (its project-hook trust step), and only on a Codex build new enough to support
  hooks; an older Codex ignores the file and keeps just the backstops below.
- Linked repos carry a local `.bench/bin/` CLI set for those hooks; a globally
  installed `bench` is convenient for humans, not required for hook execution.
- The `bench shift` loop commits only after the gate is green.

Harness hooks improve ergonomics, but the git hook and gate remain the
harness-independent backstops.
