# Dogfood improvements — make the kit safe to install and true to its promises

> **FIRST SLICE SPEC'D.** Install/setup safety is captured in
> `specs/safe-link.md`. Remaining tickets stay open for later slices.

## Grounding — what the assessment found

- **Strong base:** the gate is green, the gate canary proves covered checks still
  bite, `bench link copy` creates both `.claude` and `.agents` surfaces, and
  `bench worktree` works when `BENCH_HOME` is writable.
- **Release blocker:** `bench link` replaces existing `AGENTS.md`,
  `.claude/skills`, `.claude/commands`, `.agents/skills`, and `.agents/commands`,
  while the README says it will add alongside existing setup.
- **Shift isolation:** `bench shift` now owns a pooled worktree and leaves the main
  checkout untouched; the old promise mismatch is resolved by making the promise true.
- **Portability gap:** the autonomous loop assumes a Claude-style `agent -p
  <prompt>` interface, so non-Claude harness support needs either a wrapper
  contract or clearer documentation.
- **Packaging gap:** `.claude/settings.local.json` is publishable by accident
  because package contents include all of `.claude/`.
- **Docs drift:** the README layout still names old command files (`map`,
  `diagnose`, `review`, `verify`) instead of the current command set.
- **Canonical portable surface:** skills and commands should have one canonical
  portable tree, with `.agents/{skills,commands}` as the natural source and
  `.claude/{skills,commands}` as Claude Code adapters via symlink or generated copy.
  Hooks are different: Claude hooks are adapter-specific, while portable enforcement
  lives in `bench` and git hooks.
- **Codex hook surface:** current Codex supports project-local `.codex/hooks.json`
  or inline `.codex/config.toml` lifecycle hooks, including `PreToolUse` and `Stop`,
  gated by project trust and hook review. Bench can ship a Codex hook adapter, but
  the harness-independent safety still belongs in `bench` and git hooks.

## #1: What is the first dogfood slice?

Type: Grill

### Question
The findings split cleanly between install safety, shift isolation, harness adapter
portability, packaging hygiene, and docs drift. Which slice should the first spec
cover, so we dogfood one small green change instead of turning this into a broad
refactor?

### Answer
First slice is install/setup safety: make the setup path cleanly and safely
incorporate Bench into an existing project. The build should fix the clobbering
contract, add a behavioral gate check/canary for safe incorporation, and update only
the docs affected by that contract. Shift isolation, harness command adapters,
package hygiene, and stale README names stay captured but outside this first slice
unless they are directly required by safe setup.

## #2: What should `bench link` preserve, merge, or refuse to touch?

Blocked by: #1
Type: Grill

### Question
If install safety is in scope, define the contract: should `bench link` preserve an
existing `AGENTS.md` and append/import Bench, merge skills and commands file by
file, or stop with an explicit conflict that the reviewer resolves?

### Answer
Safe incorporation contract:

- `bench link` must not clobber existing project instructions, skills, commands, or
  hooks.
- `AGENTS.md` should receive at most a small Bench block: enough to point agents at
  the gate, the profile, and the Bench operating docs. The full Bench manual must
  not be inlined there.
- The full Bench operating docs live in `.bench/BENCH.md`. Do not use
  `CONTEXT-BENCH.md`, because `CONTEXT.md` already means domain language in this
  kit. Keep `.agents/` focused on portable agent capabilities rather than product
  docs.
- If `AGENTS.md` has custom content and no recognizable Bench block, `bench link`
  should append a small managed block delimited by `<!-- bench:start -->` and
  `<!-- bench:end -->`. If the block already exists, replace only that block. It
  must not rewrite anything outside the managed block or attempt a broad prose
  merge.
- Skills and commands should have one canonical portable tree, with Claude Code
  receiving symlinks or generated copies rather than a second source of truth.
- If a target project already has a skill or command with the same name as a Bench
  one, `bench link` fails with a clear conflict report. It must not silently skip
  or overwrite. A future explicit replace mode can be added, but is out of scope for
  the safe default.
- `bench link` defaults to copy mode. Symlink mode remains available only through an
  explicit `bench link symlink`. This keeps target repos stable after install and
  makes normal installs match npx/dlx behavior.
- Hook scripts should be shared, with adapter config for Claude (`.claude`) and
  Codex (`.codex`) pointing at the same scripts. Git safety remains backed by the
  installed git hook, independent of harness hooks.

## #3: Should `bench shift` own the worktree, or should the docs narrow the promise?

Blocked by: #1
Type: Grill

### Question
If shift isolation is in scope, choose the product truth: make `bench shift` create
and run inside a pooled worktree, or document that users must enter `bench worktree`
before running a shift.

### Answer
Shift owns the worktree. `bench shift` acquires a pooled worktree, runs the gated loop
inside it against the main checkout's committed `HEAD`, and leaves the `bench/shift-<ts>`
branch for review — the main checkout is never switched, reset, or cleaned. Implemented
from `specs/shift-in-worktree.md`.

## #4: What is the harness command contract for autonomous shifts?

Blocked by: #1
Type: Research

### Question
If portability is in scope, determine the smallest command contract that supports
Claude Code, Codex, and OpenCode without hardcoding Claude's `-p` flag into the
loop.

### Answer
— (pending)

## #5: What must the package never publish?

Blocked by: #1
Type: Research

### Question
If packaging hygiene is in scope, define the package allowlist or denylist and add
a gate check that proves local-only files such as `.claude/settings.local.json`
cannot enter the tarball.

### Answer
— (pending)

## #6: Which README claims must be synced after the chosen build?

Blocked by: #1
Type: Research

### Question
After the implementation scope is chosen, identify every README claim affected by
that scope and update only those claims so docs describe the current decided state.

### Answer
— (pending)
