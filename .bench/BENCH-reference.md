# Bench reference

This file holds lookup material split out of `.bench/BENCH.md`, so the always-loaded
operating guide stays lean. It holds the file map, the skills index,
harness-invocation details, the CLI command notes, the shift adapter contract,
and the hook layers. This file is **referenced by path, not imported**: it costs
no tokens until you open it.

Read it when you need a file's role, how a harness
invokes a phase, or how Bench layers git safety. The generation-steering rules stay
in `.bench/BENCH.md`.

**`bench`** is the operational layer over the Go core. Read `CONTEXT.md` for
the mental model and `projects/<name>.md` for seams, gate command, and line
assignments; the file map, adapter contracts, and hook layers live below.

## How the pieces fit

- **Skills** shape *how* you generate — guidance, not rules — in
  `.agents/skills/` (and `.claude/skills/` for Claude Code);
  `.bench/BENCH-reference.md` indexes them.
- **Commands:** `/bench` (Claude Code) and `$bench` (Codex) route observed state;
  `/bench-setup-repo` handles adoption and `/bench-drain` pending capture.
  `/bench-update-kit`, `/bench-assess`, and `craft-synthesis`
  ship only in the Bench kit repository; a linked repo upgrades with `bench upgrade`.
  The skills index marks omitted rows.
- **The gate and the hooks** enforce, with authority you do not have.
  `bench shift` gates every iteration and commits only on green; a `pre-push`
  hook protects the default branch.

## Files

- `AGENTS.md` contains the project-owned working agreement plus a small
  Bench-managed block.
- `.bench/gate.sh` is the project gate.
- `bin/bench.sh` is the CLI launcher and strangler router to the Go core.
- `.bench/hooks/session-start.sh` renders ambient repository status when a
  harness session starts.
  The `census` signal counts raw calls per assignment from `$BENCH_HOME/census/<repo-key>/`.
- `projects/<name>.md` holds the project's seams, gate shape, and line bindings.
- `capture/learnings.md` is the usage journal for process learnings.
- `capture/IDEAS.md` is the parked-idea sink `bench idea` writes. If
  `bench` isn't on PATH, append the dated line (`- YYYY-MM-DD  <text>`) to
  `capture/IDEAS.md` yourself.
- `capture/retros/` holds one retro per spec: `/bench-final-check` writes
  `capture/retros/<spec-slug>.md` and refreshes affected
  `capture/agent-performance/` scorecards, and `/bench-drain` owns their
  reviewed drain and its capture commit.
- `ROADMAP.md` is the working roadmap's index. It holds board prose plus one
  heading line per row, with no bodies. `roadmap/` holds one detail owner per
  row, `roadmap/FT<n>.md`, carrying that row's body, `Occurrence:` ledger, and
  `Sources:` line.
- `.bench/bin/` is the local CLI copy `bench link` installs for hooks, so Stop
  and SessionStart do not depend on a global `bench` on PATH.
- `.agents/commands/` contains portable Bench command phases.
- `.agents/skills/` contains portable Bench skills.
- `.bench/adapters/` contains reference harness adapters for `bench shift`
  (`claude`, `codex`, `opencode`). Point `BENCH_AGENT` at one of them.
  `internal/harnesses` holds one row per harness, and each row names that
  harness's headless adapter.
- `.bench/hooks/` contains shared hook scripts used by harness adapters.
- `.bench/lib/` contains shared shell functions the hooks and adapters source.
  `resolve-bench.sh` is the one source of the bench-wrapper search order. The
  tier-binding parse itself lives in the Go core, in `bench resolve-model` and
  `bench check-agent-line`.
- `.claude/` contains Claude Code adapter config. See `.claude/README.md`:
  Claude reads `.claude/skills/` and `.claude/commands/`, and those paths point
  at the portable `.agents/` files. `.claude/skills/` carries only the
  `bench-craft-*` skills. The `$bench-*` phase adapter skills stay Codex-only,
  because Claude already has each phase as a command, and a same-named skill
  would duplicate the slash-menu entry.
- `.codex/` contains Codex adapter config.

## Skills index

Claude Code loads these skills on its own. On another harness, read this file
when the trigger applies, or paste the file as context. Bench generates this
block from each skill's `index:` frontmatter (`bench skills-index --write`).
Edit the skill, not the list below:

<!-- bench:skills-index:start -->
- recording a decision or writing docs → `.agents/skills/bench-craft-adr/SKILL.md`
- building an agent-facing CLI → `.agents/skills/bench-craft-cli/SKILL.md`
- writing or reviewing code comments → `.agents/skills/bench-craft-comments/SKILL.md`
- spawning a delegate / verifying a delegate's done-claim → `.agents/skills/bench-craft-delegate/SKILL.md`
- any UI work → `.agents/skills/bench-craft-design-system/SKILL.md` + your project's design source
- pinning domain terms / enumerating concept-edge scenarios → `.agents/skills/bench-craft-domain/SKILL.md`
- adding, weakening, or removing a gate check / authoring the oracle → `.agents/skills/bench-craft-gate/SKILL.md`
- surfacing decisions in numbered frontier rounds → `.agents/skills/bench-craft-grill/SKILL.md`
- declaring the line / picking a delegate's model or effort → `.agents/skills/bench-craft-line/SKILL.md`
- reviewing a diff / what a finding must cite → `.agents/skills/bench-craft-review/SKILL.md`
- placing a test / designing an interface → `.agents/skills/bench-craft-seams/SKILL.md`
- writing or pruning a skill → `.agents/skills/bench-craft-skills/SKILL.md`
- coverage-map rows, edge inventories, story sizing, and delegate slicing for a spec → `.agents/skills/bench-craft-spec/SKILL.md`
- evaluating a change to the kit itself → `.agents/skills/bench-craft-synthesis/SKILL.md` (kit-only)
- writing tests first → `.agents/skills/bench-craft-tdd/SKILL.md`
- breaking a build into tracer-bullet tickets → `.agents/skills/bench-craft-tickets/SKILL.md`
- spiking a disposable prototype → `.agents/skills/prototype/SKILL.md`
<!-- bench:skills-index:end -->

## Harness Invocation

The canonical phase bodies live in `.agents/commands/`. Each harness may expose
these phases in a different way:

- **Claude Code:** invoke `/bench` to route from observed state, or invoke a phase
  directly as a slash command, e.g. `/bench-write-spec`.
- **Codex:** invoke `$bench` to route from observed state, or invoke the matching
  skill, e.g. `$bench-write-spec`. Each `$bench-*` adapter reads the canonical
  command file and follows it.
- **Other AGENTS.md harnesses:** read the phase file in `.agents/commands/` and
  follow it when no native command or skill surface exists.

A reviewer decides each phase's trigger on each harness; the trigger is not
ambient frontmatter. One invocation-policy table records, per phase, whether
the Claude model may reach for the command on its own. It also records
whether Codex may invoke the adapter implicitly. The gate grades both
surfaces against the table.

Most phases are reviewer-chosen entry points that
a model does not start unbidden. The maintenance phases stay off the model's
reach entirely on Claude; read a phase's own row for which case applies.
`$bench-debug` is
the exception: Codex may invoke it implicitly, because a reported symptom
should route to the bug path without the operator remembering the phase name.
Model-invoked Bench guidance otherwise uses the visible `craft-*` skill names,
leaving `$bench` for phases a reviewer deliberately runs.

The rule for translating a recommended phase into *this* harness's invocation
form lives with the communication rules in `.bench/BENCH.md`.

Codex phase adapters installed by Bench:

- `$bench` → `.agents/commands/bench.md` (front door; routes with `--harness codex`)
- `$bench-setup-repo` → `.agents/commands/bench-setup-repo.md`
- `$bench-shape-idea` → `.agents/commands/bench-shape-idea.md`
- `$bench-write-spec` → `.agents/commands/bench-write-spec.md`
- `$bench-debug` → `.agents/commands/bench-debug.md`
- `$bench-implement-spec` → `.agents/commands/bench-implement-spec.md`
- `$bench-review-implementation` → `.agents/commands/bench-review-implementation.md`
- `$bench-final-check` → `.agents/commands/bench-final-check.md`
- `$bench-update-kit` → `.agents/commands/bench-update-kit.md`
- `$bench-drain` → `.agents/commands/bench-drain.md`
- `/bench-what-next`, `$bench-what-next` → `.agents/commands/bench-drain.md` (one-release alias)
- `$bench-assess` → `.agents/commands/bench-assess.md`
- `$bench-deepen` → `.agents/commands/bench-deepen.md`

## Command Notes

Bench renders `bench help` from the Go `commandRegistry`; it is the executable
inventory. `.bench/BENCH-reference.md` gives category-level operational guidance,
not a second command list. The project profile (`projects/<name>.md`) carries the
detailed output contracts for the AXI query surfaces. The sections below
describe the hook and adapter plumbing.

- Setup and adoption connect a repository to the kit and maintain that installation.
- Context commands expose current state, navigation, capture, and planning evidence.
- Oracle commands inspect or enforce readiness from development through release.
- Work commands own isolated execution, gated changes, and spec lifecycle operations.

`bench consumers` is the resolved-reference query for a Go symbol. With
`--changed`, the same verb is the review blast over a frozen base and source
tip. Each success response ends with one citation row before the help
envelope, so a reviewer replays the answer. The review phase runs the blast
step before it dispatches the axes.

A reviewed spec-backed build keeps its serial ticket commits in one retained
integration source. Semantic review freezes the explicit base and source tip.
Accepted findings commit there on the same cadence.

A worktree `bench commit`
runs the fast lane on a private checkout of the composed snapshot, and a lane pass
publishes onto the worktree branch. The lane writes a lane record and no gate
verdict. From the destination, `bench worktree land` is the operational handoff: it
composes and runs the one whole-project gate on that pair before publication and
source release. Executable help owns its flags and positional grammar.

The spec is optional on the landing and on its resume: a spec-less phase lands
with no `--spec`, and a tickets-only `--spec` closes its folder. Every phase
lands this way; the rule is guidance, not a hook, so `bench commit` still works
on any branch. An abbreviated commit identity expands to the exact commit
before any proof runs. One preflight prints every refusal the caller must
clear, and each refusal names its paths. A stale Bench executable is rebuilt,
and the landing re-runs under it.

A conflicted `capture/` path composes by a rule table with three verbs:
`source`, `destination`, and `union`. `capture/session-handoff.md` takes
`source`; `capture/learnings.md` and `capture/IDEAS.md` take `union`; every
other `capture/` path takes `destination`. The landing discloses each
resolution, and `capture/` is authorized for every reviewed range. Any other
conflict refuses and names every path, and its `next=` names the repair in
order:

- merge the destination into the source worktree with raw Git, because no
  Bench verb moves a retained worktree onto the destination yet (FT238)
- commit the repair with `bench commit`
- review the new range
- re-run the landing with the new source tip

Its exit meanings follow the publication boundary:

- `0` — the source was released
- `1` — a refusal before publication
- `2` — invalid command usage
- `3` — publication succeeded but marker, checkout reconciliation, or source
  release remains incomplete; the exit-3 record carries the
  `bench worktree land --resume` invocation you need

## Plumbing subcommands

The registry classifies some definitions as internal inventory; they support
hooks and adapters, not interactive sessions. The registry owns their exact set, so a visibility change
cannot drift from `bench help`. Inspect the registry when you maintain those
callers.

## Phase manifest

`.bench/phases.json` in the graded root declares the project's gate as data.
Exactly two states exist. A repo that ships a manifest replaces the built-in
phase table entirely: there is no merge, so a partial manifest silently drops
every phase it omits. A repo that ships no manifest keeps the built-in table.
The Bench kit itself ships no manifest: it runs the built-in table. Its
Go-toolchain phases materialize only when the graded root carries what each
step grades. The manifest is a capability for a project whose gate is not
shaped like the kit's; it is not the route the kit itself takes.

The ordinary test phase carries the graded root and the dev tier to the
conformance entry point. The conformance registry owns check order,
subject, and tier. That one ordinary test run executes the complete dev-tier
check set, meta checks included, and so validates the registry itself. There
is no separate conformance phase, driver, or per-check evidence partition.
`gate --fresh`, prospective execution, and ship remain full boundaries.

The document is one object with a `phases` array. Each phase carries:

- `name` (required) — the phase's addressable identity. It fills the `phase`
  cell of the bounded output table, so it must be non-empty, with no whitespace
  or control character. It also prefixes the phase's lines in
  `.logs/gate-<run>.out`.
- `argv` (required, non-empty) — the command as an argument vector. Bench
  execs it directly, never through a shell, so it allows no interpolation,
  globbing, or quoting.
- `env` (optional) — a string-to-string map set in the phase's environment; it
  overrides the gate's own values.
- `needs` (optional) — names of phases that must end green before this phase
  starts. When a needed phase ends red or skipped, Bench skips this phase too,
  because it would grade an artifact its need never produced. Every name
  named here must exist in the manifest.
- `optional` (optional, default false) — when the command is not installed,
  the phase reports skipped ("not installed") instead of red. A present
  command that fails still reports red.
- `dir` (optional) — the phase's working directory, a relative path anchored
  to the graded root. In a linked repo the graded root is a different tree
  from the kit checkout. Only a root-relative path lands where the
  project's directories actually are. An empty value means the root itself.

```json
{
  "phases": [
    {"name": "build", "argv": ["npm", "run", "build"]},
    {"name": "test", "argv": ["npm", "test"], "needs": ["build"], "dir": "web"}
  ]
}
```

The loader fails closed. Anything between the two valid states reds before
any phase runs. It names the path, the defect class, and the offending
element. The offending element is one of:

- a truncated write or trailing content
- a dangling symlink or non-regular file
- an unknown field
- a duplicate name
- an empty argv
- a `dir` that is absolute or escapes the root
- a `needs` edge to a phase that does not exist, or a cycle

A manifest with any of these defects means the loader cannot know what its
author intended. A guessed-at table would grade the tree with the wrong
oracle. If you hit one of these reds, the refusal is deliberate: fix the
manifest. There is no lenient mode.

### Gate output

A red run prints one `failures[N]{phase,line}` table, and then `gate: red`.
The table holds failure rows only. Each phase gives at most twenty rows, and one
more row names the file that holds the complete stream.

A green run prints one `phases[N]{phase,verdict,elapsed_ms}` table, one
`capability-skips` line, and then `gate: green`. Above seven phases the table
collapses to one `phases: N/N green` row.

The complete phase stream goes to `.logs/gate-<run>.out`. That file sits beside
the `.logs/gate-<run>.jsonl` progress log, under the same twenty-run retention.

## Harness adapter for the shift loop

`bench shift` drives whatever harness `BENCH_AGENT` names. Each iteration it
runs the adapter executable with the generated prompt written to its
**stdin** — no positional argument — and with `BENCH_SHIFT=1` armed. The
prompt is multi-line and may start with a dash. An adapter must read all of
stdin, must not re-expose the prompt as a CLI argument, and exits with its
harness's exit code. The loop takes that exit code as progress evidence, and
the gate stays the oracle.

The adapter launches with a documented
**passlisted environment**, not the parent's full environment. Widen it only
by committing extra names under the `[agent]` section (the only section) of
`.bench/env.allow`. Declare a variable the *gate* needs in
`.bench/gate-inputs.json` instead; `bench setup` seeds that file with the
names the installed wrapper needs (`BENCH_HOME` and `HOME`). There is no
default: an unset `BENCH_AGENT` fails fast before the loop with a
configure-your-adapter error.

Reference adapters ship in `.bench/adapters/`
(`claude`, `codex`, `opencode`). Point `BENCH_AGENT` at one, or at your own
wrapper that pipes its stdin to your harness's noninteractive stdin-reading
command. The `claude` and `codex` adapters do exactly this. `opencode` reads
stdin and hands it to `opencode run` positionally after `--`, an upstream
residual until opencode documents a stdin form. Use an absolute path or an
on-`PATH` name. Put harness flags inside the wrapper, because Bench treats a
multi-word `BENCH_AGENT` value as one executable name and rejects it.

The adapters also carry the line (see the `craft-line` skill). When
`BENCH_MODEL` is set, the adapter passes it to the harness's model flag. A repo with
`.bench/lines.env` (the harness × tier binding matrix) is **routed**. There
`BENCH_MODEL` names a tier — `top`, `mid`, or `cheap` — and each adapter asks
the core for its own harness's column.

The adapter refuses to run when the
tier is unset or unknown, or when its harness's column is unbound. A
headless shift always carries an explicit, bound line. Without `lines.env`
the adapters behave as plain pass-throughs and forward `BENCH_MODEL`
verbatim. Effort has no harness flag and stays in the declared line.

## Hook Layers

Bench layers git safety:

- The git `pre-push` hook blocks a direct push to the default branch. When a
  gate pin exists (`bench gate pin`), it also blocks `.bench` drift from the
  pinned tree. Without a pin, the drift check stays disarmed; guard discovery
  reports a static, generic deny surface while enforcement stays live.
- Claude Code and Codex hook adapters call the shared scripts in
  `.bench/hooks/`. Codex loads `.codex/hooks.json` only after you trust it
  once through `/hooks` (its project-hook trust step). It loads only on a
  Codex build new enough to support hooks. An older Codex ignores the file
  and keeps only the backstops below.
- The agent-line guard (`check-agent-line`) wires on Claude Code only. Codex
  cannot host it, because a delegation (`spawn_agent`) never surfaces as a
  matchable `tool_name` on a deny-capable event. Run `bench harnesses codex`
  for that verdict with its source and its date. The line's harness-independent
  backstop is the shift adapters'
  refusal to run with an unset or unbound `BENCH_MODEL`. If the Codex
  changelog adds a spawn tool name or a deny-capable SubagentStart, revisit
  this rule.
- Linked repos carry a local `.bench/bin/` CLI set for those hooks. A globally
  installed `bench` is convenient for a human; hook execution does not need
  it.
- The `bench shift` loop commits only after the gate turns green.

Harness hooks improve ergonomics, but the git hook and the gate remain the
harness-independent backstops.
