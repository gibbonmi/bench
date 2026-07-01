# Dogfood improvements — make the kit safe to install and true to its promises

> **ARCHITECTURE REVIEW TRIAGED.** Install/setup safety is captured in
> `specs/safe-link.md`. Shift ownership is captured in
> `specs/shift-in-worktree.md`. The 2026-07-01 architecture review pulled the next
> concrete defects into specs where the build target is clear, and left only the
> genuinely decision-shaped items open.

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
- **Gate execution leak:** `bench gate` resolves the repo root but runs
  `.bench/gate.sh` and `$BENCH_GATE` in the caller's current directory, unlike the
  auto-detected gate path. The gate seam is therefore shallower than its contract.
- **Worktree cleanup leak:** the pooled worktree release path resets tracked files
  but leaves ignored artifacts behind, so a later `bench worktree` or `bench shift`
  can observe state from a previous run.
- **Git-safety classifier gap:** the destructive-git guard blocks broad reset/clean
  forms but misses path-level checkout/restore commands that can discard user edits.
- **Structure path handling gap:** `bench structure` counts crowded directories by
  piping paths through whitespace-splitting shell tools, so paths with spaces are
  reported incorrectly.
- **Installable-surface drift:** package contents exclude repo-only profiles while a
  shipped pickup doc claimed those profiles ship. The package allowlist, link plan,
  and shipped-doc claims need one installable-surface truth.
- **Line cap honesty:** `BENCH_MAX_TOKENS` is printed as a shift cap but is not
  enforced or measured. That may be acceptable, but the product meaning needs a
  decision before changing behavior.

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

Current known constraint from the architecture review: the loop invokes
`"$AGENT" -p <prompt>`, so a portable contract must either define an adapter command
that accepts `-p` everywhere or move prompt delivery behind a harness-specific
adapter.

## #5: What must the package never publish?

Blocked by: #1
Type: Research

### Question
If packaging hygiene is in scope, define the package allowlist or denylist and add
a gate check that proves local-only files such as `.claude/settings.local.json`
cannot enter the tarball.

### Answer
Resolved by the safe-link slice for the concrete package-local-settings problem:
local-only settings must not enter the npm tarball, and the gate's dry-run package
inspection proves that. The remaining package issue is not "what never publishes" but
"what shipped docs may claim about the installable surface"; that is tracked in #11
and folded into `specs/safe-link.md`.

## #6: Which README claims must be synced after the chosen build?

Blocked by: #1
Type: Research

### Question
After the implementation scope is chosen, identify every README claim affected by
that scope and update only those claims so docs describe the current decided state.

### Answer
Resolved for command-name drift by `specs/doc-command-currency.md`. The architecture
review found a different docs/package mismatch: a shipped pickup doc claimed
repo-only profiles under `projects/` ship even though the npm package excludes them.
That installable-surface claim is tracked in #11.

## #7: What is the gate execution contract?

Blocked by: #1
Type: Research

### Question
`bench gate` knows the repo root, but project gates run in the caller's current
directory while auto-detected gates run at the root. Should the contract be "gates
run from wherever the user invoked Bench" or "all oracle commands run at the repo
root"?

### Answer
All oracle commands run at the repo root. The gate is the highest seam in Bench, and
callers should not need to remember cwd rules to use it correctly. The build target is
`specs/gate-execution-contract.md`: deepen the gate execution module so `.bench/gate.sh`,
`$BENCH_GATE`, auto-detect, `bench shift`, and the Stop hook all share the same root
semantics.

## #8: What does a clean pooled worktree mean?

Blocked by: #3
Type: Research

### Question
The worktree pool currently treats a porcelain-clean tree as reusable, but ignored
artifacts can survive `git clean -fd`. Is "clean" only tracked-file clean, or must it
mean no carryover state a later shift can observe?

### Answer
Clean means no carryover state a later shift or interactive worktree can observe,
including ignored artifacts. This belongs in the existing worktree reuse spec instead
of a parallel spec; `specs/shift-loop-hardening.md` now owns the ignored-artifact
regression case.

## #9: What should the destructive-git guard classify as blocked?

Blocked by: #1
Type: Research

### Question
The guard blocks broad destructive operations, but path-level `git checkout -- file`
and `git restore file` can still discard user edits. Should the guard stay as a small
denylist, or should it classify destructive git intent more deeply?

### Answer
Classify destructive intent more deeply while keeping harmless read commands
allowed. The build target is `specs/git-safety-classifier.md`: make the hook adapter
call one classifier that blocks path-level checkout/restore alongside the existing
push, hard reset, force clean, branch delete, and rebase blocks.

## #10: How should `bench structure` handle pathnames?

Blocked by: #1
Type: Research

### Question
`bench structure` is an agent-facing CLI. Should it assume simple pathnames, or does
the CLI contract require paths with spaces to render correct file and directory
signals?

### Answer
The CLI must handle normal git pathnames, including spaces, without corrupting
directory counts. The build target is `specs/structure-path-safety.md`: deepen the
structure module's file-list handling so source files and crowded-directory counts do
not pass through whitespace-splitting tools.

## #11: How does the installable surface stay truthful?

Blocked by: #5, #6
Type: Research

### Question
Package contents, `bench link`, README layout, and HANDOFF all describe what ships or
installs. Should Bench keep checking only package contents, or also check that shipped
docs do not claim repo-only assets are in the package?

### Answer
Shipped docs may only claim assets that are actually in the installable surface, or
must label repo-only artifacts as local development context. This is part of the
existing safe-link/package contract, not a new surface; `specs/safe-link.md` now owns
the docs/package truth row.

## #12: What should `BENCH_MAX_TOKENS` mean?

Blocked by: #4
Type: Grill

### Question
`BENCH_MAX_TOKENS` is printed as a shift cap but not enforced or measured. Should
Bench enforce it, remove it, rename it to an advisory line value, or move token caps
entirely into the declared line outside the CLI?

### Answer
— (open)
