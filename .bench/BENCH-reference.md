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
- `capture/learnings.md` is the usage journal for process learnings.
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
- writing or reviewing code comments → `.agents/skills/bench-craft-comments/SKILL.md`
- spawning a delegate / verifying a delegate's done-claim → `.agents/skills/bench-craft-delegate/SKILL.md`
- any UI work → `.agents/skills/bench-craft-design-system/SKILL.md` + your project's design source
- adding, weakening, or removing a gate check / authoring the oracle → `.agents/skills/bench-craft-gate/SKILL.md`
- surfacing a decision one question at a time → `.agents/skills/bench-craft-grill/SKILL.md`
- declaring the line / picking a delegate's model or effort → `.agents/skills/bench-craft-line/SKILL.md`
- reviewing a diff / what a finding must cite → `.agents/skills/bench-craft-review/SKILL.md`
- placing a test / designing an interface → `.agents/skills/bench-craft-seams/SKILL.md`
- writing or pruning a skill → `.agents/skills/bench-craft-skills/SKILL.md`
- coverage-map rows, edge inventories, story sizing, and delegate slicing for a spec → `.agents/skills/bench-craft-spec/SKILL.md`
- evaluating a change to the kit itself → `.agents/skills/bench-craft-synthesis/SKILL.md` (kit-only)
- writing tests first → `.agents/skills/bench-craft-tdd/SKILL.md`
- breaking a build into independently-green tickets → `.agents/skills/bench-craft-tickets/SKILL.md`
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
- `$bench-what-next` → `.agents/commands/bench-what-next.md`
- `$bench-assess` → `.agents/commands/bench-assess.md`

## Command Notes

The canonical CLI inventory lives in `.bench/BENCH.md`, not in `HANDOFF.md`.
Detailed output contracts for the AXI query surfaces live in the project
profile (`projects/<name>.md`);
hook and adapter plumbing is described in the sections below.

### Spec-build lifecycle lookup

The canonical implementation phase is `.agents/commands/bench-implement-spec.md`;
this table is the operation-to-purpose lookup, not a second workflow:

| operation | routing purpose |
|---|---|
| `start` | create or resume the subject-bound run from exact-tip whole-tree green, including a narrow verdict whose inherited evidence still covers every skip |
| `assign` | lease one ownership-fenced ticket worktree |
| `checkpoint` | validate focused evidence and bind a provisional commit |
| `integrate` | compare-and-swap one verified checkpoint into the candidate |
| `review` | bind three-axis evidence to the exact candidate |
| `status` | inspect durable state and retained evidence |
| `promote` | gate and publish the exact reviewed composition; a moved tip recomposes through `promote`, discarding the review |
| `abandon` | plan or apply recoverable cleanup |

## Plumbing subcommands

Driven by hooks and adapters, never typed by sessions — the one enumeration
(the always-loaded inventory in `.bench/BENCH.md` points here): `bench tree-hash`,
`bench gate-run`, `bench freshness-check`, `bench gate-phases`, `bench gate-go`, `bench release-preflight`, `bench guard-git`, `bench resolve-model`,
`bench check-agent-line`, `bench stop-verdict`, `bench session-inspect`, `bench worktree-pool`,
`bench worktree-lease-file`, `bench worktree-hook`, `bench resume-clean`.

## Phase manifest

`.bench/phases.json` in the graded root declares the project's gate as data.
There are exactly two states: a repo that ships a manifest replaces the built-in
phase table entirely — there is no merge, so a partial manifest silently drops
every phase it omits — and a repo that ships none keeps the built-in table. The
Bench kit itself ships no manifest: it runs the built-in table, whose
Go-toolchain phases materialize only when the graded root carries what each step
grades. The manifest is a capability for projects whose gate is not shaped like
the kit's, not the route the kit takes.

The built-in conformance phase receives a gate-authored, registry-ordered
ordinary-check set. Ambient singular and plural selectors are removed before
phase construction; the singular form is restored only for an authenticated
inner canary run. An ordinary dev gate may inherit exact per-check evidence, but
meta checks always execute in the same aggregate process and validate the
complete, disjoint executed/inherited partition. `gate --fresh`, prospective
execution, and ship remain full boundaries. A mixed verdict lists every
executed check and carries the identity and authorship time for every inherited
check; it can compose to the exact tree's landing green but is never reusable as
a later whole-tree green.

The document is one object with a `phases` array. Per phase:

- `name` (required) — the phase's addressable identity: it appears in summary
  lines and output prefixes and targets a single phase for the canary, so it
  must be non-empty with no whitespace or control characters.
- `argv` (required, non-empty) — the command as an argument vector, exec'd
  directly, never through a shell: no interpolation, globbing, or quoting.
- `env` (optional) — string-to-string map set in the phase's environment,
  overriding the gate's own values.
- `needs` (optional) — names of phases that must end green before this one
  starts. A phase whose need ends red or skipped is skipped too — it would
  grade an artifact its need never produced. Every name must exist in the
  manifest.
- `optional` (optional, default false) — when the command is not installed, the
  phase reports skipped ("not installed") instead of red. A present command
  that fails is still red.
- `dir` (optional) — the phase's working directory, a relative path anchored to
  the graded root — which in a linked repo is a different tree from the kit
  checkout, so a root-relative path is the only kind that lands where the
  project's directories actually are. Empty means the root itself.

```json
{
  "phases": [
    {"name": "build", "argv": ["npm", "run", "build"]},
    {"name": "test", "argv": ["npm", "test"], "needs": ["build"], "dir": "web"}
  ]
}
```

The loader fails closed. Anything between the two valid states — a truncated
write or trailing content, a dangling symlink or non-regular file, an unknown
field, a duplicate name, an empty argv, a `dir` that is absolute or escapes the
root, a `needs` edge to a phase that does not exist, a cycle — reds before any
phase runs, naming the path, the defect class, and the offending element. A
manifest with any of these means something its author intended that the loader
cannot know, and running a guessed-at table would grade the tree with the wrong
oracle. If you hit one of these reds, the refusal is deliberate: fix the
manifest — there is no lenient mode.

## Harness adapter for the shift loop

`bench shift` drives whatever harness `BENCH_AGENT` names: each iteration it runs
the adapter executable with the generated prompt written to its **stdin** — no
positional argument — and `BENCH_SHIFT=1` armed. The prompt is multi-line and may
start with a dash; an adapter reads all of stdin, must not re-expose the prompt as
a CLI argument, and exits with its harness's exit code, which the loop takes as
progress evidence (the gate stays the oracle). The adapter launches with a
documented **passlisted environment**, not the parent's full environment; widen it
only by committing extra names under the `[agent]` section (the only section) of
`.bench/env.allow` — a variable the *gate* needs is declared in
`.bench/gate-inputs.json` instead. There is no default — an unset `BENCH_AGENT`
fails fast before the loop with a configure-your-adapter error. Reference adapters
ship in `.bench/adapters/` (`claude`, `codex`, `opencode`); point `BENCH_AGENT` at
one, or at your own wrapper that pipes its stdin to your harness's noninteractive
stdin-reading command (the `claude` and `codex` adapters do exactly this;
`opencode` reads stdin and hands it to `opencode run` positionally after `--`, an
upstream residual until opencode documents a stdin form). Use an absolute path or
an on-`PATH` name; harness flags belong inside the wrapper — a multi-word
`BENCH_AGENT` value is treated as one executable name and rejected.

The adapters also carry the line (see the `craft-line` skill): `BENCH_MODEL`,
when set, is passed to the harness's model flag. A repo with `.bench/lines.env`
(the harness × tier binding matrix) is **routed**: there `BENCH_MODEL` names a
tier — `top`, `mid`, or `cheap` — each adapter asks the core for its own
harness's column, and the adapter refuses to run when the tier is unset or
unknown or its harness's column is unbound, so a headless shift always carries an
explicit, bound line. Without `lines.env` the adapters behave as plain
pass-throughs and forward `BENCH_MODEL` verbatim. Effort has no harness flag and stays in
the declared line.

## Hook Layers

Git safety is layered:

- The git `pre-push` hook blocks direct pushes to the default branch and, when
  a gate pin is recorded (`bench gate pin`), `.bench` drift from the pinned
  tree. Unpinned, the drift check is disarmed; guard discovery reports a static,
  generic deny surface while enforcement remains live.
- Claude Code and Codex hook adapters call the shared scripts in `.bench/hooks/`.
  Codex loads `.codex/hooks.json` only after you trust it once via `/hooks`
  (its project-hook trust step), and only on a Codex build new enough to support
  hooks; an older Codex ignores the file and keeps just the backstops below.
- The agent-line guard (`check-agent-line`) wires on Claude Code only. Codex
  cannot host it: a delegation (`spawn_agent`) never surfaces as a matchable
  `tool_name` on a deny-capable event. `SubagentStart` carries the active model
  through its common input fields, but it cannot deny the spawn:
  `continue: false` does not stop the subagent (Codex hooks docs, checked
  2026-07-11). The line's harness-independent backstop is the shift adapters'
  refusal to run with an unset or unbound `BENCH_MODEL`. Re-check if the Codex
  changelog adds a spawn tool name or a deny-capable SubagentStart.
- Linked repos carry a local `.bench/bin/` CLI set for those hooks; a globally
  installed `bench` is convenient for humans, not required for hook execution.
- The `bench shift` loop commits only after the gate is green.

Harness hooks improve ergonomics, but the git hook and gate remain the
harness-independent backstops.
